package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

const gceCredentialTypeName = "Google Compute Engine" //nolint:gosec

func resourceCredentialGoogleComputeEngine() *schema.Resource {
	return &schema.Resource{
		Description:   "`awx_credential_google_compute_engine` manages Google Compute Engine credentials in AWX.",
		CreateContext: resourceCredentialGoogleComputeEngineCreate,
		ReadContext:   resourceCredentialGoogleComputeEngineRead,
		UpdateContext: resourceCredentialGoogleComputeEngineUpdate,
		DeleteContext: resourceCredentialDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the credential.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the credential.",
			},
			"organization_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The organization ID this credential belongs to.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username to use for the credential.",
			},
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The project to use for the credential.",
			},
			"ssh_key_data": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The SSH key data to use for the credential.",
			},
		},
	}
}

func resourceCredentialGoogleComputeEngineCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	client := m.(*awx.AWX)
	gceCredType, err := client.CredentialTypeService.GetCredentialTypeByName(gceCredentialTypeName, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create new credentials",
			Detail:   fmt.Sprintf("Unable to fetch credential type with Name: %s. Error: %s", gceCredentialTypeName, err.Error()),
		})
		return diags
	}

	newCredential := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"organization":    d.Get("organization_id").(int),
		"credential_type": gceCredType.ID,
		"inputs": map[string]interface{}{
			"username":     d.Get("username").(string),
			"project":      d.Get("project").(string),
			"ssh_key_data": d.Get("ssh_key_data").(string),
		},
	}

	cred, err := client.CredentialsService.CreateCredentials(newCredential, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create new credentials",
			Detail:   fmt.Sprintf("Unable to create new credentials: %s", err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(cred.ID))
	resourceCredentialGoogleComputeEngineRead(ctx, d, m)

	return diags
}

func resourceCredentialGoogleComputeEngineRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	id, _ := strconv.Atoi(d.Id())
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
	if err := d.Set("username", cred.Inputs["username"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project", cred.Inputs["project"]); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCredentialGoogleComputeEngineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	keys := []string{
		"name",
		"description",
		"username",
		"project",
		"ssh_key_data",
	}

	client := m.(*awx.AWX)

	if d.HasChanges(keys...) {
		var err error

		gceCredType, err := client.CredentialTypeService.GetCredentialTypeByName(gceCredentialTypeName, map[string]string{})
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create new credentials",
				Detail:   fmt.Sprintf("Unable to fetch credential type with Name: %s. Error: %s", gceCredentialTypeName, err.Error()),
			})
			return diags
		}

		id, _ := strconv.Atoi(d.Id())
		updatedCredential := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			"credential_type": gceCredType.ID,
			"inputs": map[string]interface{}{
				"username":     d.Get("username").(string),
				"project":      d.Get("project").(string),
				"ssh_key_data": d.Get("ssh_key_data").(string),
			},
		}

		_, err = client.CredentialsService.UpdateCredentialsByID(id, updatedCredential, map[string]string{})
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to update existing credentials",
				Detail:   fmt.Sprintf("Unable to update existing credentials with id %d: %s", id, err.Error()),
			})
			return diags
		}
	}

	return resourceCredentialGoogleComputeEngineRead(ctx, d, m)
}
