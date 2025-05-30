package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

const containerRegistryCredentialTypeName = "Container Registry" //nolint:gosec

func resourceCredentialContainerRegistry() *schema.Resource {
	return &schema.Resource{
		Description:   "`awx_credential_container_registry` manages container registry credentials in AWX.",
		CreateContext: resourceCredentialContainerRegistryCreate,
		ReadContext:   resourceCredentialContainerRegistryRead,
		UpdateContext: resourceCredentialContainerRegistryUpdate,
		DeleteContext: resourceCredentialDelete,
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
				Optional:    true,
				Description: "The username to use for the credential.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "A password or token used to authenticate with.",
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Authentication endpoint for the container registry.",
			},
			"verify_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Verify SSL",
			},
		},
	}
}

func resourceCredentialContainerRegistryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	client := m.(*awx.AWX)
	containerRegistryCredType, err := client.CredentialTypeService.GetCredentialTypeByName(containerRegistryCredentialTypeName, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create new credentials",
			Detail:   fmt.Sprintf("Unable to fetch credential type with Name: %s. Error: %s", containerRegistryCredentialTypeName, err.Error()),
		})
		return diags
	}

	newCredential := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"organization":    d.Get("organization_id").(int),
		"credential_type": containerRegistryCredType.ID,
		"inputs": map[string]interface{}{
			"username":   d.Get("username").(string),
			"password":   d.Get("password").(string),
			"host":       d.Get("host").(string),
			"verify_ssl": d.Get("verify_ssl").(bool),
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
	resourceCredentialContainerRegistryRead(ctx, d, m)

	return diags
}

func resourceCredentialContainerRegistryRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	if err := setSanitizedEncryptedCredential(d, "password", cred); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host", cred.Inputs["host"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("verify_ssl", cred.Inputs["verify_ssl"]); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCredentialContainerRegistryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	keys := []string{
		"name",
		"description",
		"username",
		"password",
		"host",
		"verify_ssl",
	}

	client := m.(*awx.AWX)

	if d.HasChanges(keys...) {
		var err error

		containerRegistryCredType, err := client.CredentialTypeService.GetCredentialTypeByName(containerRegistryCredentialTypeName, map[string]string{})
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create new credentials",
				Detail:   fmt.Sprintf("Unable to fetch credential type with Name: %s. Error: %s", containerRegistryCredentialTypeName, err.Error()),
			})
			return diags
		}

		id, _ := strconv.Atoi(d.Id())
		updatedCredential := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			"credential_type": containerRegistryCredType.ID,
			"inputs": map[string]interface{}{
				"username":   d.Get("username").(string),
				"password":   d.Get("password").(string),
				"host":       d.Get("host").(string),
				"verify_ssl": d.Get("verify_ssl").(bool),
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

	return resourceCredentialContainerRegistryRead(ctx, d, m)
}
