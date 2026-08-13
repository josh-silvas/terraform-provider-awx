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

func resourceWorkflowJobTemplateNode() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource WorkflowJobTemplateNode manages the workflow job template node in AWX.",
		CreateContext: resourceWorkflowJobTemplateNodeCreate,
		ReadContext:   resourceWorkflowJobTemplateNodeRead,
		UpdateContext: resourceWorkflowJobTemplateNodeUpdate,
		DeleteContext: resourceWorkflowJobTemplateNodeDelete,

		Schema: map[string]*schema.Schema{
			"extra_data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Extra data for the workflow job template node.",
			},
			"inventory_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Inventory applied as a prompt, assuming job template prompts for inventory.",
			},
			"scm_branch": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "SCM branch to use for the job template.",
			},
			"job_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "run",
				Description: "Type of job to run.",
			},
			"job_tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Tags to use for the job template.",
			},
			"skip_tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Tags to skip for the job template.",
			},
			"limit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Limit the job template to a specific host or group.",
			},
			"diff_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable diff mode for the job template.",
			},
			"verbosity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Verbosity level for the job template. One of 0, 1, 2, 3, 4 or 5.",
			},
			"workflow_job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Workflow job template ID to use for the workflow job template node.",
			},
			"unified_job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unified job template ID to use for the workflow job template node.",
			},
			"all_parents_must_converge": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"identifier": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for the workflow job template node.",
			},
		},
	}
}

func resourceWorkflowJobTemplateNodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateNodeService

	result, err := awxService.CreateWorkflowJobTemplateNode(map[string]interface{}{
		"extra_data":                d.Get("extra_data"),
		"inventory":                 d.Get("inventory_id").(int),
		"scm_branch":                d.Get("scm_branch").(string),
		"skip_tags":                 d.Get("skip_tags").(string),
		"job_type":                  d.Get("job_type").(string),
		"job_tags":                  d.Get("job_tags").(string),
		"limit":                     d.Get("limit").(string),
		"diff_mode":                 d.Get("diff_mode").(bool),
		"verbosity":                 d.Get("verbosity").(int),
		"workflow_job_template":     d.Get("workflow_job_template_id").(int),
		"unified_job_template":      d.Get("unified_job_template_id").(int),
		"all_parents_must_converge": d.Get("all_parents_must_converge").(bool),
		"identifier":                d.Get("identifier").(string),
	}, map[string]string{})
	if err != nil {
		log.Printf("Fail to Create Template %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create WorkflowJobTemplateNode",
			Detail:   fmt.Sprintf("WorkflowJobTemplateNode with JobTemplateID %d and WorkflowID: %d failed to create %s", d.Get("unified_job_template_id").(int), d.Get("workflow_job_template_id").(int), err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceWorkflowJobTemplateNodeRead(ctx, d, m)
}

func resourceWorkflowJobTemplateNodeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Update WorkflowJobTemplateNode", d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)
	if _, err := client.WorkflowJobTemplateNodeService.GetWorkflowJobTemplateNodeByID(id, params); err != nil {
		return utils.DiagNotFound("workflow job template node", id, err)
	}

	if _, err := client.WorkflowJobTemplateNodeService.UpdateWorkflowJobTemplateNode(id, map[string]interface{}{
		"extra_data":                d.Get("extra_data"),
		"inventory":                 d.Get("inventory_id").(int),
		"scm_branch":                d.Get("scm_branch").(string),
		"skip_tags":                 d.Get("skip_tags").(string),
		"job_type":                  d.Get("job_type").(string),
		"job_tags":                  d.Get("job_tags").(string),
		"limit":                     d.Get("limit").(string),
		"diff_mode":                 d.Get("diff_mode").(bool),
		"verbosity":                 d.Get("verbosity").(int),
		"workflow_job_template":     d.Get("workflow_job_template_id").(int),
		"unified_job_template":      d.Get("unified_job_template_id").(int),
		"all_parents_must_converge": d.Get("all_parents_must_converge").(bool),
		"identifier":                d.Get("identifier").(string),
	}, map[string]string{}); err != nil {
		return utils.DiagUpdate("workflow job template node", d.Get("name").(string), err)
	}

	return resourceWorkflowJobTemplateNodeRead(ctx, d, m)
}

func resourceWorkflowJobTemplateNodeRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Read WorkflowJobTemplateNode", d)
	if diags.HasError() {
		return diags
	}

	res, err := client.WorkflowJobTemplateNodeService.GetWorkflowJobTemplateNodeByID(id, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound("workflow job template node", id, err)

	}
	d = setWorkflowJobTemplateNodeResourceData(d, res)
	return nil
}

func resourceWorkflowJobTemplateNodeDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Workflow Job Template Node", d)
	if diags.HasError() {
		return diags
	}

	if _, err := client.WorkflowJobTemplateNodeService.DeleteWorkflowJobTemplateNode(id); err != nil {
		return utils.DiagDelete("workflow job template node", id, err)
	}
	d.SetId("")
	return nil
}

func setWorkflowJobTemplateNodeResourceData(d *schema.ResourceData, r *awx.WorkflowJobTemplateNode) *schema.ResourceData {

	if err := d.Set("extra_data", r.ExtraData); err != nil {
		fmt.Println("Error setting extra_data", err)
	}
	if err := d.Set("inventory_id", r.Inventory); err != nil {
		fmt.Println("Error setting inventory_id", err)
	}
	if err := d.Set("scm_branch", r.ScmBranch); err != nil {
		fmt.Println("Error setting scm_branch", err)
	}
	if err := d.Set("job_type", r.JobType); err != nil {
		fmt.Println("Error setting job_type", err)
	}
	if err := d.Set("job_tags", r.JobTags); err != nil {
		fmt.Println("Error setting job_tags", err)
	}
	if err := d.Set("skip_tags", r.SkipTags); err != nil {
		fmt.Println("Error setting skip_tags", err)
	}
	if err := d.Set("limit", r.Limit); err != nil {
		fmt.Println("Error setting limit", err)
	}
	if err := d.Set("diff_mode", r.DiffMode); err != nil {
		fmt.Println("Error setting diff_mode", err)
	}
	if err := d.Set("verbosity", r.Verbosity); err != nil {
		fmt.Println("Error setting verbosity", err)
	}

	if err := d.Set("workflow_job_template_id", r.WorkflowJobTemplate); err != nil {
		fmt.Println("Error setting workflow_job_template_id", err)
	}
	if err := d.Set("unified_job_template_id", r.UnifiedJobTemplate); err != nil {
		fmt.Println("Error setting unified_job_template_id", err)
	}
	if err := d.Set("all_parents_must_converge", r.AllParentsMustConverge); err != nil {
		fmt.Println("Error setting all_parents_must_converge", err)
	}
	if err := d.Set("identifier", r.Identifier); err != nil {
		fmt.Println("Error setting identifier", err)
	}

	d.SetId(strconv.Itoa(r.ID))
	return d
}
