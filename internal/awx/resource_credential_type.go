package awx

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

func resourceCredentialType() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource `awx_credential_type` manages credential types within an AWX instance.",
		CreateContext: resourceCredentialTypeCreate,
		ReadContext:   resourceCredentialTypeRead,
		UpdateContext: resourceCredentialTypeUpdate,
		DeleteContext: resourceCredentialTypeDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this credential type.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional description of this credential type.",
			},
			"kind": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "cloud",
				Description: "Can be one of: `cloud` or `net`",
			},
			"inputs": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Inputs for this credential type.",
			},
			"injectors": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Injectors for this credential type.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCredentialTypeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	inputs := d.Get("inputs").(string)
	inputsMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(inputs), &inputsMap); err != nil {
		return utils.DiagCreate("Credential Type", err)
	}

	injectors := d.Get("injectors").(string)
	injectorsMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(injectors), &injectorsMap); err != nil {
		return utils.DiagCreate("Credential Type", err)
	}

	newCredentialType := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"kind":        d.Get("kind").(string),
		"inputs":      inputsMap,
		"injectors":   injectorsMap,
	}

	client := m.(*awx.AWX)
	credType, err := client.CredentialTypeService.CreateCredentialType(newCredentialType, map[string]string{})
	if err != nil {
		return utils.DiagCreate("Credential Type", err)
	}

	d.SetId(strconv.Itoa(credType.ID))
	resourceCredentialTypeRead(ctx, d, m)

	return diag.Diagnostics{}
}

func resourceCredentialTypeRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return utils.DiagFetch("Credential Type", id, err)
	}
	credType, err := client.CredentialTypeService.GetCredentialTypeByID(id, map[string]string{})
	if err != nil {
		return utils.DiagFetch("Credential Type", id, err)
	}

	if err := d.Set("name", credType.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", credType.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("kind", credType.Kind); err != nil {
		return diag.FromErr(err)
	}

	inputsBytes, err := json.Marshal(credType.Inputs)
	if err != nil {
		return diag.FromErr(err)
	}
	injectorBytes, err := json.Marshal(credType.Injectors)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("inputs", string(inputsBytes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("injectors", string(injectorBytes)); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceCredentialTypeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	keys := []string{
		"name",
		"description",
		"kind",
		"inputs",
		"injectors",
	}

	if d.HasChanges(keys...) {
		inputs := d.Get("inputs").(string)
		inputsMap := make(map[string]interface{})
		if err := json.Unmarshal([]byte(inputs), &inputsMap); err != nil {
			return utils.DiagUpdate("Credential Type", d.Id(), err)
		}

		injectors := d.Get("injectors").(string)
		injectorsMap := make(map[string]interface{})
		if err := json.Unmarshal([]byte(injectors), &injectorsMap); err != nil {
			return utils.DiagUpdate("Credential Type", d.Id(), err)
		}

		id, err := strconv.Atoi(d.Id())
		if err != nil {
			return utils.DiagUpdate("Credential Type", id, err)
		}
		payload := map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
			"kind":        d.Get("kind").(string),
			"inputs":      inputsMap,
			"injectors":   injectorsMap,
		}

		client := m.(*awx.AWX)
		if _, err = client.CredentialTypeService.UpdateCredentialTypeByID(id, payload, map[string]string{}); err != nil {
			return utils.DiagUpdate("Credential Type", id, err)
		}
	}

	return resourceCredentialTypeRead(ctx, d, m)
}

func resourceCredentialTypeDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return utils.DiagDelete("Credential Type", id, err)
	}
	client := m.(*awx.AWX)
	if err := client.CredentialTypeService.DeleteCredentialTypeByID(id, map[string]string{}); err != nil {
		return utils.DiagDelete("Credential Type", id, err)
	}
	return diag.Diagnostics{}
}
