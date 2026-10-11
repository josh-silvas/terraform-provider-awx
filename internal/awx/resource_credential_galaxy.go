package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

const galaxyCredentialTypeName = "Ansible Galaxy/Automation Hub API Token" //nolint:gosec

// credentialGalaxySecrets are the secret inputs that also accept a write-only "<field>_wo" variant.
var credentialGalaxySecrets = map[string]string{
	"token": "The API token for the Ansible Galaxy/Automation Hub API.",
}

func resourceCredentialGalaxy() *schema.Resource {
	r := &schema.Resource{
		Description:   "`awx_credential_galaxy` manages Ansible Galaxy/Automation Hub API Token credentials in AWX.",
		CreateContext: resourceCredentialGalaxyCreate,
		ReadContext:   resourceCredentialGalaxyRead,
		UpdateContext: resourceCredentialGalaxyUpdate,
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
				Description: "The organization that the credential belongs to.",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL of the Ansible Galaxy/Automation Hub API.",
			},
			"auth_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL of the Ansible Galaxy/Automation Hub API authentication endpoint.",
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The API token for the Ansible Galaxy/Automation Hub API.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}

	for field, description := range credentialGalaxySecrets {
		addWriteOnlySecret(r.Schema, field, description)
	}

	return r
}

func resourceCredentialGalaxyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	client := m.(*awx.AWX)
	credType, err := client.CredentialTypeService.GetCredentialTypeByName(galaxyCredentialTypeName, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credential type",
			Detail:   fmt.Sprintf("Unable to fetch credential type with Name: %s. Error: %s", galaxyCredentialTypeName, err.Error()),
		})
		return diags
	}

	newCredential := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"organization":    d.Get("organization_id").(int),
		"credential_type": credType.ID,
		"inputs": map[string]interface{}{
			"url":      d.Get("url").(string),
			"auth_url": d.Get("auth_url").(string),
			"token":    getSecret(d, "token"),
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
	resourceCredentialGalaxyRead(ctx, d, m)

	return diags
}

func resourceCredentialGalaxyRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	if err := d.Set("url", cred.Inputs["url"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("auth_url", cred.Inputs["auth_url"]); err != nil {
		return diag.FromErr(err)
	}
	if err := setSanitizedEncryptedCredential(d, "token", cred); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("organization_id", cred.OrganizationID); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCredentialGalaxyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	keys := []string{
		"name",
		"description",
		"url",
		"auth_url",
		"token",
		"organization_id",
		"team_id",
		"owner_id",
	}
	keys = append(keys, writeOnlyVersionKeys(credentialGalaxySecrets)...)

	if d.HasChanges(keys...) {
		var err error

		id, _ := strconv.Atoi(d.Id())
		client := m.(*awx.AWX)
		credType, err := client.CredentialTypeService.GetCredentialTypeByName(galaxyCredentialTypeName, map[string]string{})
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to fetch credential type",
				Detail:   fmt.Sprintf("Unable to fetch credential type with Name: %s. Error: %s", galaxyCredentialTypeName, err.Error()),
			})
			return diags
		}

		updatedCredential := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			"credential_type": credType.ID,
			"inputs": map[string]interface{}{
				"url":      d.Get("url").(string),
				"auth_url": d.Get("auth_url").(string),
				"token":    getSecret(d, "token"),
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

	return resourceCredentialGalaxyRead(ctx, d, m)
}
