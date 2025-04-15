// Package awx provides a Terraform provider for Ansible AWX.
package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func dataSourceCredentialByID() *schema.Resource {
	return &schema.Resource{
		Description: "This data source provides a credential object back to Terraform. " +
			"See: https://docs.ansible.com/ansible-tower/latest/html/towerapi/api_ref.html#/Credentials/Credentials_credentials_read",
		ReadContext: dataSourceCredentialByIDRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the credential to fetch",
			},
			"tower_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for tower.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The username for the credential",
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The kind of credential",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the credential",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the credential",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func dataSourceCredentialByIDRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	id := d.Get("id").(int)
	cred, err := client.CredentialsService.GetCredentialsByID(id, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credential",
			Detail:   "The given credential ID is invalid or malformed",
		})
	}

	if err := d.Set("username", cred.Inputs["username"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("kind", cred.Kind); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tower_id", id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", cred.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", cred.Name); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(id))
	// d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
