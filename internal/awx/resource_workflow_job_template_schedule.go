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

func resourceWorkflowJobTemplateSchedule() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource `awx_workflow_job_template_schedule` manages workflow job template schedule within AWX.",
		CreateContext: resourceWorkflowJobTemplateScheduleCreate,
		ReadContext:   resourceScheduleRead,
		UpdateContext: resourceScheduleUpdate,
		DeleteContext: resourceScheduleDelete,
		Schema: map[string]*schema.Schema{

			"workflow_job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The workflow_job_template id for this schedule",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the schedule",
			},
			"rrule": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The recurrence rule for the schedule. See https://github.com/ansible/awx/blob/devel/awx/api/templates/api/_schedule_detail.md for more information.",
			},
			"unified_job_template_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The unified job template id for this schedule",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the schedule",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the schedule is enabled or not",
			},
			"inventory": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Inventory applied as a prompt, assuming job template prompts for inventory (id, default=``)",
			},
			"extra_data": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				Description:      "Extra data to be pass for the schedule (YAML format)",
				DiffSuppressFunc: suppressYAMLDiff,
			},
		},
	}
}

func resourceWorkflowJobTemplateScheduleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateScheduleService

	workflowJobTemplateID := d.Get("workflow_job_template_id").(int)

	result, err := awxService.CreateWorkflowJobTemplateSchedule(workflowJobTemplateID, map[string]interface{}{
		"name":        d.Get("name").(string),
		"rrule":       d.Get("rrule").(string),
		"description": d.Get("description").(string),
		"enabled":     d.Get("enabled").(bool),
		"inventory":   utils.AtoiDefault(d.Get("inventory").(string), nil),
		"extra_data":  utils.UnmarshalYAML(d.Get("extra_data").(string)),
	}, map[string]string{})
	if err != nil {
		log.Printf("Fail to Create Schedule for WorkflowJobTemplate %d: %v", workflowJobTemplateID, err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Schedule",
			Detail:   fmt.Sprintf("Schedule failed to create %s", err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceScheduleRead(ctx, d, m)
}
