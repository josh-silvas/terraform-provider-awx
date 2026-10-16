package awx

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagCredentialTitle = "Credential"

func resourceCredential() *schema.Resource {
	return &schema.Resource{
		Description:   "The `awx_credential` resource allows you to create and manage credentials in Ansible Tower.",
		CreateContext: resourceCredentialCreate,
		ReadContext:   resourceCredentialRead,
		UpdateContext: resourceCredentialUpdate,
		DeleteContext: resourceCredentialDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the credential",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the credential",
			},
			"organization_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The organization ID that the credential belongs to",
			},
			"credential_type_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Specify the type of credential you want to create. Refer to the Ansible Tower documentation for details on each type",
			},
			"inputs": {
				Type:             schema.TypeMap,
				Optional:         true,
				Sensitive:        true,
				Description:      "The inputs to be created with the credential.",
				DiffSuppressFunc: suppressValueDiff,
			},
			"secret_inputs_wo": {
				Type:         schema.TypeString,
				Optional:     true,
				WriteOnly:    true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsJSON,
				Description: "Secret inputs as a JSON object, e.g. `jsonencode({ password = var.password })`, merged with `inputs` " +
					"before being sent to AWX. Keys must not also be set in `inputs`. Write-only: never stored in state. Requires Terraform >= 1.11.",
			},
			"secret_inputs_wo_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"secret_inputs_wo"},
				Description:  "Version of `secret_inputs_wo`. Change it to send a new value of `secret_inputs_wo` to AWX.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCredentialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	inputsMap, err := credentialInputs(d)
	if err != nil {
		return utils.DiagCreate(diagCredentialTitle, err)
	}

	payload := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"organization":    d.Get("organization_id").(int),
		"credential_type": d.Get("credential_type_id").(int),
		"inputs":          inputsMap,
	}

	client := m.(*awx.AWX)
	cred, err := client.CredentialsService.CreateCredentials(payload, map[string]string{})
	if err != nil {
		return utils.DiagCreate(diagCredentialTitle, err)
	}

	d.SetId(strconv.Itoa(cred.ID))
	resourceCredentialRead(ctx, d, m)
	return diag.Diagnostics{}
}

func resourceCredentialRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return utils.DiagFetch(diagCredentialTitle, d.Id(), err)
	}
	cred, err := client.CredentialsService.GetCredentialsByID(id, map[string]string{})
	if err != nil {
		return utils.DiagFetch(diagCredentialTitle, d.Id(), err)
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
	if err := d.Set("credential_type_id", cred.CredentialTypeID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("inputs", stateCredentialInputs(d.Get("inputs").(map[string]interface{}), cred.Inputs)); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceCredentialUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	keys := []string{
		"name",
		"description",
		"organization_id",
		"inputs",
		"secret_inputs_wo_version",
	}

	if d.HasChanges(keys...) {
		inputsMap, err := credentialInputs(d)
		if err != nil {
			return utils.DiagUpdate(diagCredentialTitle, d.Id(), err)
		}

		id, err := strconv.Atoi(d.Id())
		if err != nil {
			return utils.DiagUpdate(diagCredentialTitle, d.Id(), err)
		}
		update := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			"credential_type": d.Get("credential_type_id"),
			"inputs":          inputsMap,
		}

		client := m.(*awx.AWX)
		if _, err = client.CredentialsService.UpdateCredentialsByID(id, update, map[string]string{}); err != nil {
			return utils.DiagUpdate(diagCredentialTitle, d.Id(), err)
		}
	}

	return resourceCredentialRead(ctx, d, m)
}

func resourceCredentialDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return utils.DiagDelete(diagCredentialTitle, d.Id(), err)
	}
	client := m.(*awx.AWX)
	if err := client.CredentialsService.DeleteCredentialsByID(id, map[string]string{}); err != nil {
		return utils.DiagDelete(diagCredentialTitle, d.Id(), err)
	}

	return diag.Diagnostics{}
}

// credentialInputs merges "inputs" with the write-only "secret_inputs_wo" JSON object from the configuration.
func credentialInputs(d *schema.ResourceData) (map[string]interface{}, error) {
	inputs := map[string]interface{}{}
	for k, v := range d.Get("inputs").(map[string]interface{}) {
		inputs[k] = v
	}

	raw, diags := d.GetRawConfigAt(cty.GetAttrPath("secret_inputs_wo"))
	if diags.HasError() {
		return inputs, nil
	}
	secretJSON, ok := stringFromCty(raw)
	if !ok || secretJSON == "" {
		return inputs, nil
	}

	return mergeSecretInputs(inputs, secretJSON)
}

func mergeSecretInputs(inputs map[string]interface{}, secretJSON string) (map[string]interface{}, error) {
	var secrets map[string]interface{}
	if err := json.Unmarshal([]byte(secretJSON), &secrets); err != nil {
		return nil, fmt.Errorf("secret_inputs_wo must be a JSON object: %w", err)
	}
	for k, v := range secrets {
		if _, ok := inputs[k]; ok {
			return nil, fmt.Errorf("input %q is set in both inputs and secret_inputs_wo", k)
		}
		inputs[k] = v
	}

	return inputs, nil
}

// stateCredentialInputs returns the inputs to store in state. AWX returns every secret input as '$encrypted$';
// those not already tracked in state (e.g. sent through secret_inputs_wo) are left out to avoid a permanent diff.
func stateCredentialInputs(prior, remote map[string]interface{}) map[string]interface{} {
	inputs := make(map[string]interface{}, len(remote))
	for k, v := range remote {
		if _, tracked := prior[k]; v == "$encrypted$" && !tracked {
			continue
		}
		inputs[k] = v
	}

	return inputs
}
