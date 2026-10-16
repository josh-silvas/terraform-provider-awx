package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

const gitlabCredentialTypeName = "GitLab Personal Access Token" //nolint:gosec

// credentialGitlabSecrets are the secret inputs that also accept a write-only "<field>_wo" variant.
var credentialGitlabSecrets = map[string]string{
	"token": "The GitLab Personal Access Token.",
}

func resourceCredentialGitlab() *schema.Resource {
	r := &schema.Resource{
		Description:   "`awx_credential_gitlab` manages GitLab Personal Access Token credentials in AWX.",
		CreateContext: resourceCredentialGitlabCreate,
		ReadContext:   resourceCredentialGitlabRead,
		UpdateContext: resourceCredentialGitlabUpdate,
		DeleteContext: resourceCredentialDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the credential.",
			},
			"organization_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The organization ID that this credential belongs to.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the credential.",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The GitLab Personal Access Token.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}

	for field, description := range credentialGitlabSecrets {
		addWriteOnlySecret(r.Schema, field, description)
	}

	return r
}

func resourceCredentialGitlabCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	client := m.(*awx.AWX)
	credType, err := client.CredentialTypeService.GetCredentialTypeByName(gitlabCredentialTypeName, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credential type",
			Detail:   fmt.Sprintf("Unable to fetch credential type with Name: %s. Error: %s", gitlabCredentialTypeName, err.Error()),
		})
		return diags
	}

	newCredential := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"organization":    d.Get("organization_id").(int),
		"credential_type": credType.ID,
		"inputs": map[string]interface{}{
			"token": getSecret(d, "token"),
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
	resourceCredentialGitlabRead(ctx, d, m)

	return diags
}

func resourceCredentialGitlabRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	if err := setSanitizedEncryptedCredential(d, "token", cred); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("organization_id", cred.OrganizationID); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCredentialGitlabUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	keys := []string{
		"name",
		"description",
		"token",
		"organization_id",
	}
	keys = append(keys, writeOnlyVersionKeys(credentialGitlabSecrets)...)

	if d.HasChanges(keys...) {
		var err error

		id, _ := strconv.Atoi(d.Id())
		client := m.(*awx.AWX)
		credType, err := client.CredentialTypeService.GetCredentialTypeByName(gitlabCredentialTypeName, map[string]string{})
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to fetch credential type",
				Detail:   fmt.Sprintf("Unable to fetch credential type with Name: %s. Error: %s", gitlabCredentialTypeName, err.Error()),
			})
			return diags
		}

		updatedCredential := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			"credential_type": credType.ID,
			"inputs": map[string]interface{}{
				"token": getSecret(d, "token"),
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

	return resourceCredentialGitlabRead(ctx, d, m)
}
