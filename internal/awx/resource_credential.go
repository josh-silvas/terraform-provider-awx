package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Required:         true,
				Sensitive:        true,
				Description:      "The inputs to be created with the credential.",
				DiffSuppressFunc: suppressValueDiff,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCredentialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	inputsMap := d.Get("inputs")

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
	if err := d.Set("inputs", cred.Inputs); err != nil {
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
	}

	if d.HasChanges(keys...) {
		var err error

		inputsMap := d.Get("inputs")

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

// suppressCredentialInputsDiff suppresses diffs when AWX returns encrypted placeholders
//func suppressCredentialInputsDiff(k, old, new string, d *schema.ResourceData) bool {
//	// AWX returns "$encrypted$" for sensitive credential inputs
//	// Suppress diff if old value contains this placeholder
//	if old == "$encrypted$" {
//		return true
//	}
//	return false
//}
