package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

//nolint:funlen
func resourceCredentialMachine() *schema.Resource {
	return &schema.Resource{
		Description:   "The `awx_credential_machine` resource allows creation and management of machine credentials within an AWX instance.",
		CreateContext: resourceCredentialMachineCreate,
		ReadContext:   resourceCredentialMachineRead,
		UpdateContext: resourceCredentialMachineUpdate,
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
}

func resourceCredentialMachineCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	newCredential := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"organization":    d.Get("organization_id").(int),
		"credential_type": 1, // SSH
		"inputs": map[string]interface{}{
			"username":            d.Get("username").(string),
			"password":            d.Get("password").(string),
			"ssh_key_data":        d.Get("ssh_key_data").(string),
			"ssh_public_key_data": d.Get("ssh_public_key_data").(string),
			"ssh_key_unlock":      d.Get("ssh_key_unlock").(string),
			"become_method":       d.Get("become_method").(string),
			"become_username":     d.Get("become_username").(string),
			"become_password":     d.Get("become_password").(string),
		},
	}

	client := m.(*awx.AWX)
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
	if err := setSanitizedEncryptedCredential(d, "ssh_public_key_data", cred); err != nil {
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

	if d.HasChanges(keys...) {
		var err error

		id, _ := strconv.Atoi(d.Id())
		updatedCredential := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			"credential_type": 1, // SSH
			"inputs": map[string]interface{}{
				"username":            d.Get("username").(string),
				"password":            d.Get("password").(string),
				"ssh_key_data":        d.Get("ssh_key_data").(string),
				"ssh_public_key_data": d.Get("ssh_public_key_data").(string),
				"ssh_key_unlock":      d.Get("ssh_key_unlock").(string),
				"become_method":       d.Get("become_method").(string),
				"become_username":     d.Get("become_username").(string),
				"become_password":     d.Get("become_password").(string),
			},
		}

		client := m.(*awx.AWX)
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
