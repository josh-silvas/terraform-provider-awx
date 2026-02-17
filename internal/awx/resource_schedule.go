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

func resourceSchedule() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource Schedule for AWX (Ansible Tower)",
		CreateContext: resourceScheduleCreate,
		ReadContext:   resourceScheduleRead,
		UpdateContext: resourceScheduleUpdate,
		DeleteContext: resourceScheduleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the schedule",
			},
			"rrule": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "RRule for the schedule. See https://github.com/ansible/awx/blob/devel/awx/api/templates/api/_schedule_detail.md for more information. ",
			},
			"unified_job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the Unified Job Template to be scheduled",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the schedule",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable or disable the schedule",
			},
			"inventory": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the Inventory to be used for the schedule",
			},
			"extra_data": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				Description:      "Extra data to be pass for the schedule (JSON format)",
				DiffSuppressFunc: SuppressEquivalentJSONDiffs,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceScheduleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.ScheduleService

	scheduleData := map[string]interface{}{
		"name":                 d.Get("name").(string),
		"rrule":                d.Get("rrule").(string),
		"unified_job_template": d.Get("unified_job_template_id").(int),
		"description":          d.Get("description").(string),
		"enabled":              d.Get("enabled").(bool),
		"extra_data":           utils.UnmarshalJSON(d.Get("extra_data").(string)),
	}
	if _, ok := d.GetOk("inventory"); ok {
		scheduleData["inventory"] = d.Get("inventory").(int)
	}

	result, err := awxService.Create(scheduleData, map[string]string{})
	if err != nil {
		log.Printf("Fail to Create Schedule %v", err)
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

func resourceScheduleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Update Schedule", d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)
	if _, err := client.ScheduleService.GetByID(id, params); err != nil {
		return utils.DiagNotFound("Schedule", id, err)
	}

	payload := map[string]interface{}{
		"name":                 d.Get("name").(string),
		"rrule":                d.Get("rrule").(string),
		"unified_job_template": d.Get("unified_job_template_id").(int),
		"description":          d.Get("description").(string),
		"enabled":              d.Get("enabled").(bool),
		"extra_data":           utils.UnmarshalJSON(d.Get("extra_data").(string)),
	}
	if _, ok := d.GetOk("inventory"); ok {
		payload["inventory"] = d.Get("inventory").(int)
	}

	if _, err := client.ScheduleService.Update(id, payload, map[string]string{}); err != nil {
		return utils.DiagUpdate("Schedule", id, err)
	}

	return resourceScheduleRead(ctx, d, m)
}

func resourceScheduleRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Read schedule", d)
	if diags.HasError() {
		return diags
	}

	res, err := client.ScheduleService.GetByID(id, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound("Schedule", id, err)

	}
	d = setScheduleResourceData(d, res)
	return nil
}

func resourceScheduleDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagHostTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := client.ScheduleService.Delete(id); err != nil {
		return utils.DiagDelete("Schedule", id, err)
	}
	d.SetId("")
	return nil
}

func setScheduleResourceData(d *schema.ResourceData, r *awx.Schedule) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	if err := d.Set("rrule", r.Rrule); err != nil {
		fmt.Println("Error setting rrule", err)
	}
	if err := d.Set("unified_job_template_id", r.UnifiedJobTemplate); err != nil {
		fmt.Println("Error setting unified_job_template_id", err)
	}
	if err := d.Set("description", r.Description); err != nil {
		fmt.Println("Error setting description", err)
	}
	if err := d.Set("enabled", r.Enabled); err != nil {
		fmt.Println("Error setting enabled", err)
	}
	if err := d.Set("inventory", r.Inventory); err != nil {
		fmt.Println("Error setting inventory", err)
	}
	if err := d.Set("extra_data", utils.MarshalJSON(r.ExtraData)); err != nil {
		fmt.Println("Error setting extra_data", err)
	}
	d.SetId(strconv.Itoa(r.ID))
	return d
}
