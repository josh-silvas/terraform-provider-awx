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

func dataSourceCredentialSCM() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialSCMRead,
		Description: "Data source for Source Control credentials in AWX.",
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
			"ssh_key_unlock": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The SSH key unlock passphrase for the credential.",
			},
		},
	}
}

func dataSourceCredentialSCMRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	if err := d.Set("ssh_key_unlock", d.Get("ssh_key_unlock").(string)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
