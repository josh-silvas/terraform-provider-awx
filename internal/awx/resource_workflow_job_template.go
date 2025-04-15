package awx

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

func resourceWorkflowJobTemplate() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource `awx_workflow_job_template` manages workflow job templates within AWX.",
		CreateContext: resourceWorkflowJobTemplateCreate,
		ReadContext:   resourceWorkflowJobTemplateRead,
		UpdateContext: resourceWorkflowJobTemplateUpdate,
		DeleteContext: resourceWorkflowJobTemplateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this workflow job template. (string, required)",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Optional description of this workflow job template.",
			},
			"variables": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				StateFunc:   utils.Normalize,
				Description: "Extra variables used by Ansible in YAML or JSON format. (string, default=``)",
			},
			"organization_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The organization used to determine access to this template. (id, default=``)",
			},
			"survey_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"allow_simultaneous": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_variables_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"inventory_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Inventory applied as a prompt, assuming job template prompts for inventory. (id, default=``)",
				Default:     "",
			},
			"limit": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"scm_branch": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ask_inventory_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_scm_branch_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_limit_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"webhook_service": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"webhook_credential": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceWorkflowJobTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateService

	result, err := awxService.CreateWorkflowJobTemplate(map[string]interface{}{
		"name":                     d.Get("name").(string),
		"description":              d.Get("description").(string),
		"organization":             d.Get("organization_id").(int),
		"inventory":                utils.AtoiDefault(d.Get("inventory_id").(string), nil),
		"extra_vars":               d.Get("variables").(string),
		"survey_enabled":           d.Get("survey_enabled").(bool),
		"allow_simultaneous":       d.Get("allow_simultaneous").(bool),
		"ask_variables_on_launch":  d.Get("ask_variables_on_launch").(bool),
		"limit":                    d.Get("limit").(string),
		"scm_branch":               d.Get("scm_branch").(string),
		"ask_inventory_on_launch":  d.Get("ask_inventory_on_launch").(bool),
		"ask_scm_branch_on_launch": d.Get("ask_scm_branch_on_launch").(bool),
		"ask_limit_on_launch":      d.Get("ask_limit_on_launch").(bool),
		"webhook_service":          d.Get("webhook_service").(string),
		"webhook_credential":       d.Get("webhook_credential").(string),
	}, map[string]string{})
	if err != nil {
		log.Printf("Fail to Create Template %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create WorkflowJobTemplate",
			Detail:   fmt.Sprintf("WorkflowJobTemplate with name %s failed to create %s", d.Get("name").(string), err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceWorkflowJobTemplateRead(ctx, d, m)
}

func resourceWorkflowJobTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Update WorkflowJobTemplate", d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)
	if _, err := client.WorkflowJobTemplateService.GetWorkflowJobTemplateByID(id, params); err != nil {
		return utils.DiagNotFound("job Workflow template", id, err)
	}

	if _, err := client.WorkflowJobTemplateService.UpdateWorkflowJobTemplate(id, map[string]interface{}{
		"name":                     d.Get("name").(string),
		"description":              d.Get("description").(string),
		"organization":             d.Get("organization_id").(int),
		"inventory":                utils.AtoiDefault(d.Get("inventory_id").(string), nil),
		"extra_vars":               d.Get("variables").(string),
		"survey_enabled":           d.Get("survey_enabled").(bool),
		"allow_simultaneous":       d.Get("allow_simultaneous").(bool),
		"ask_variables_on_launch":  d.Get("ask_variables_on_launch").(bool),
		"limit":                    d.Get("limit").(string),
		"scm_branch":               d.Get("scm_branch").(string),
		"ask_inventory_on_launch":  d.Get("ask_inventory_on_launch").(bool),
		"ask_scm_branch_on_launch": d.Get("ask_scm_branch_on_launch").(bool),
		"ask_limit_on_launch":      d.Get("ask_limit_on_launch").(bool),
		"webhook_service":          d.Get("webhook_service").(string),
		"webhook_credential":       d.Get("webhook_credential").(string),
	}, map[string]string{}); err != nil {
		return utils.DiagUpdate("Job Workflow template", d.Get("name").(string), err)
	}

	return resourceWorkflowJobTemplateRead(ctx, d, m)
}

func resourceWorkflowJobTemplateRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Read WorkflowJobTemplate", d)
	if diags.HasError() {
		return diags
	}

	res, err := client.WorkflowJobTemplateService.GetWorkflowJobTemplateByID(id, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound("workflow job template", id, err)

	}
	d = setWorkflowJobTemplateResourceData(d, res)
	return nil
}

func resourceWorkflowJobTemplateDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Workflow Job Template", d)
	if diags.HasError() {
		return diags
	}

	if _, err := client.WorkflowJobTemplateService.DeleteWorkflowJobTemplate(id); err != nil {
		return utils.DiagDelete("Workflow Job Template", id, err)
	}
	d.SetId("")
	return nil
}

func setWorkflowJobTemplateResourceData(d *schema.ResourceData, r *awx.WorkflowJobTemplate) *schema.ResourceData {

	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	if err := d.Set("description", r.Description); err != nil {
		fmt.Println("Error setting description", err)
	}
	if err := d.Set("organization_id", strconv.Itoa(r.Organization)); err != nil {
		fmt.Println("Error setting organization_id", err)
	}
	if err := d.Set("inventory_id", utils.ItoaDefault(r.Inventory, "")); err != nil {
		fmt.Println("Error setting inventory_id", err)
	}
	if err := d.Set("survey_enabled", r.SurveyEnabled); err != nil {
		fmt.Println("Error setting survey_enabled", err)
	}
	if err := d.Set("allow_simultaneous", r.AllowSimultaneous); err != nil {
		fmt.Println("Error setting allow_simultaneous", err)
	}
	if err := d.Set("ask_variables_on_launch", r.AskVariablesOnLaunch); err != nil {
		fmt.Println("Error setting ask_variables_on_launch", err)
	}
	if err := d.Set("limit", r.Limit); err != nil {
		fmt.Println("Error setting limit", err)
	}
	if err := d.Set("scm_branch", r.ScmBranch); err != nil {
		fmt.Println("Error setting scm_branch", err)
	}
	if err := d.Set("ask_inventory_on_launch", r.AskInventoryOnLaunch); err != nil {
		fmt.Println("Error setting ask_inventory_on_launch", err)
	}
	if err := d.Set("ask_scm_branch_on_launch", r.AskScmBranchOnLaunch); err != nil {
		fmt.Println("Error setting ask_scm_branch_on_launch", err)
	}
	if err := d.Set("ask_limit_on_launch", r.AskLimitOnLaunch); err != nil {
		fmt.Println("Error setting ask_limit_on_launch", err)
	}
	if err := d.Set("webhook_service", r.WebhookService); err != nil {
		fmt.Println("Error setting webhook_service", err)
	}
	if err := d.Set("webhook_credential", r.WebhookCredential); err != nil {
		fmt.Println("Error setting webhook_credential", err)
	}
	if err := d.Set("variables", utils.Normalize(r.ExtraVars)); err != nil {
		fmt.Println("Error setting variables", err)
	}

	d.SetId(strconv.Itoa(r.ID))
	return d
}
