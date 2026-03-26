package awx

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func dataSourceCredentialVault() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialVaultRead,
		Description: "Data source for Vault credentials in AWX.",
		Schema: map[string]*schema.Schema{
			"credential_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the credential.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the credential.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the credential.",
			},
			"organization_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The organization ID of the credential.",
			},
			"vault_password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Vault Password.",
			},
			"vault_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The vault identity to use.",
			},
		},
	}
}

func dataSourceCredentialVaultRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	id, _ := d.Get("credential_id").(int)
	cred, err := client.CredentialsService.GetCredentialsByID(id, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credentials",
			Detail:   fmt.Sprintf("Unable to fetch credentials with id %d: %s", id, err.Error()),
		})
		return diags
	}

	if err := d.Set("name", cred.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", cred.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("organization_id", cred.OrganizationID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vault_password", d.Get("vault_password").(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vault_id", cred.Inputs["vault_id"]); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
