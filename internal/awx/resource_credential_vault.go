package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

const vaultCredentialTypeName = "Vault" //nolint:gosec

func resourceCredentialVault() *schema.Resource {
	return &schema.Resource{
		Description:   "`awx_credential_vault` manages vault credentials in AWX.",
		CreateContext: resourceCredentialVaultCreate,
		ReadContext:   resourceCredentialVaultRead,
		UpdateContext: resourceCredentialVaultUpdate,
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
				Optional:    true,
				Description: "The organization ID this credential belongs to.",
			},
			"vault_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Vault Password.",
			},
			"vault_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The vault identity to use.",
			},
		},
	}
}

func resourceCredentialVaultCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	client := m.(*awx.AWX)
	vaultCredType, err := client.CredentialTypeService.GetCredentialTypeByName(vaultCredentialTypeName, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create new credentials",
			Detail:   fmt.Sprintf("Unable to fetch credential type with Name: %s. Error: %s", vaultCredentialTypeName, err.Error()),
		})
		return diags
	}

	newCredential := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"organization":    d.Get("organization_id").(int),
		"credential_type": vaultCredType.ID,
		"inputs": map[string]interface{}{
			"vault_password": d.Get("vault_password").(string),
			"vault_id":       d.Get("vault_id").(string),
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
	resourceCredentialVaultRead(ctx, d, m)

	return diags
}

func resourceCredentialVaultRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	id, _ := strconv.Atoi(d.Id())
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
	if err := setSanitizedEncryptedCredential(d, "vault_password", cred); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vault_id", cred.Inputs["vault_id"]); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCredentialVaultUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	keys := []string{
		"name",
		"description",
		"organization_id",
		"vault_password",
		"vault_id",
	}

	client := m.(*awx.AWX)

	if d.HasChanges(keys...) {
		var err error

		vaultCredType, err := client.CredentialTypeService.GetCredentialTypeByName(vaultCredentialTypeName, map[string]string{})
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create new credentials",
				Detail:   fmt.Sprintf("Unable to fetch credential type with Name: %s. Error: %s", vaultCredentialTypeName, err.Error()),
			})
			return diags
		}

		id, _ := strconv.Atoi(d.Id())
		updatedCredential := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			"credential_type": vaultCredType.ID,
			"inputs": map[string]interface{}{
				"vault_password": d.Get("vault_password").(string),
				"vault_id":       d.Get("vault_id").(string),
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

	return resourceCredentialVaultRead(ctx, d, m)
}
