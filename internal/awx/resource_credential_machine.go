package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

const machineCredentialTypeName = "Machine" //nolint:gosec

// credentialMachineSecrets are the secret inputs that also accept a write-only "<field>_wo" variant.
var credentialMachineSecrets = map[string]string{
	"password":        "The password for the credential.",
	"ssh_key_data":    "The SSH key data for the credential.",
	"ssh_key_unlock":  "The SSH key unlock for the credential.",
	"become_password": "The become password for the credential.",
}

//nolint:funlen
func resourceCredentialMachine() *schema.Resource {
	r := &schema.Resource{
		Description:   "The `awx_credential_machine` resource allows creation and management of machine credentials within an AWX instance.",
		CreateContext: resourceCredentialMachineCreate,
		ReadContext:   resourceCredentialMachineRead,
		UpdateContext: resourceCredentialMachineUpdate,
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
				Description: "The organization ID that the credential belongs to.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The username for the credential.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The password for the credential.",
			},
			"ssh_key_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The SSH key data for the credential.",
			},
			"ssh_public_key_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The SSH public key data for the credential.",
			},
			"ssh_key_unlock": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The SSH key unlock for the credential.",
			},
			"become_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The become method for the credential.",
			},
			"become_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The become username for the credential.",
			},
			"become_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The become password for the credential.",
			},
		},
	}

	for field, description := range credentialMachineSecrets {
		addWriteOnlySecret(r.Schema, field, description)
	}

	return r
}

func resourceCredentialMachineCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	client := m.(*awx.AWX)
	credType, err := client.CredentialTypeService.GetCredentialTypeByName(machineCredentialTypeName, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credential type",
			Detail:   fmt.Sprintf("Unable to fetch credential type with Name: %s. Error: %s", machineCredentialTypeName, err.Error()),
		})
		return diags
	}

	newCredential := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"organization":    d.Get("organization_id").(int),
		"credential_type": credType.ID,
		"inputs": map[string]interface{}{
			"username":            d.Get("username").(string),
			"password":            getSecret(d, "password"),
			"ssh_key_data":        getSecret(d, "ssh_key_data"),
			"ssh_public_key_data": d.Get("ssh_public_key_data").(string),
			"ssh_key_unlock":      getSecret(d, "ssh_key_unlock"),
			"become_method":       d.Get("become_method").(string),
			"become_username":     d.Get("become_username").(string),
			"become_password":     getSecret(d, "become_password"),
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
	resourceCredentialMachineRead(ctx, d, m)

	return diags
}

func resourceCredentialMachineRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	if err := d.Set("username", cred.Inputs["username"]); err != nil {
		return diag.FromErr(err)
	}
	if err := setSanitizedEncryptedCredential(d, "password", cred); err != nil {
		return diag.FromErr(err)
	}
	if err := setSanitizedEncryptedCredential(d, "ssh_key_data", cred); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ssh_public_key_data", cred.Inputs["ssh_public_key_data"]); err != nil {
		return diag.FromErr(err)
	}
	if err := setSanitizedEncryptedCredential(d, "ssh_key_unlock", cred); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("become_method", cred.Inputs["become_method"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("become_username", cred.Inputs["become_username"]); err != nil {
		return diag.FromErr(err)
	}
	if err := setSanitizedEncryptedCredential(d, "become_password", cred); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("organization_id", cred.OrganizationID); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCredentialMachineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	keys := []string{
		"name",
		"description",
		"username",
		"password",
		"ssh_key_data",
		"ssh_public_key_data",
		"ssh_key_unlock",
		"become_method",
		"become_username",
		"become_password",
		"organization_id",
		"team_id",
		"owner_id",
	}
	keys = append(keys, writeOnlyVersionKeys(credentialMachineSecrets)...)

	if d.HasChanges(keys...) {
		var err error

		id, _ := strconv.Atoi(d.Id())
		client := m.(*awx.AWX)
		credType, err := client.CredentialTypeService.GetCredentialTypeByName(machineCredentialTypeName, map[string]string{})
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to fetch credential type",
				Detail:   fmt.Sprintf("Unable to fetch credential type with Name: %s. Error: %s", machineCredentialTypeName, err.Error()),
			})
			return diags
		}

		updatedCredential := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			"credential_type": credType.ID,
			"inputs": map[string]interface{}{
				"username":            d.Get("username").(string),
				"password":            getSecret(d, "password"),
				"ssh_key_data":        getSecret(d, "ssh_key_data"),
				"ssh_public_key_data": d.Get("ssh_public_key_data").(string),
				"ssh_key_unlock":      getSecret(d, "ssh_key_unlock"),
				"become_method":       d.Get("become_method").(string),
				"become_username":     d.Get("become_username").(string),
				"become_password":     getSecret(d, "become_password"),
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

	return resourceCredentialMachineRead(ctx, d, m)
}
