// Package awx provides a Terraform provider for Ansible AWX.
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

func dataSourceCredentialAzure() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialAzureRead,
		Description: "Data source for Azure Key Vault credentials in AWX. " +
			"See: https://docs.ansible.com/ansible-tower/latest/html/towerapi/api_ref.html#/Credentials/Credentials_credentials_read",
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
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"tenant": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCredentialAzureRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	id, _ := d.Get("credential_id").(int)
	cred, err := client.CredentialsService.GetCredentialsByID(id, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credentials",
			Detail:   fmt.Sprintf("Unable to credentials with id %d: %s", id, err.Error()),
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
	if err := d.Set("url", cred.Inputs["url"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("client", cred.Inputs["client"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("secret", d.Get("secret").(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant", cred.Inputs["tenant"]); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
