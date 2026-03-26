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

func dataSourceCredentialMachine() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialMachineRead,
		Description: "Data source for Machine credentials in AWX.",
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
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The username for the credential.",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The password for the credential.",
			},
			"ssh_key_data": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The SSH key data for the credential.",
			},
			"ssh_public_key_data": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SSH public key data for the credential.",
			},
			"ssh_key_unlock": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The SSH key unlock passphrase for the credential.",
			},
			"become_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The become method for the credential.",
			},
			"become_username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The become username for the credential.",
			},
			"become_password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The become password for the credential.",
			},
		},
	}
}

func dataSourceCredentialMachineRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	if err := d.Set("username", cred.Inputs["username"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("password", d.Get("password").(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ssh_key_data", d.Get("ssh_key_data").(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ssh_public_key_data", cred.Inputs["ssh_public_key_data"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ssh_key_unlock", d.Get("ssh_key_unlock").(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("become_method", cred.Inputs["become_method"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("become_username", cred.Inputs["become_username"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("become_password", d.Get("become_password").(string)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
