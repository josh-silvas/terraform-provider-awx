package awx

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagJobTemplateLaunchTitle = "Job Template Launch"

//nolint:funlen
func resourceJobTemplateLaunch() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource `awx_job_template_launch` manages job template launch.",
		CreateContext: resourceJobTemplateLaunchCreate,
		ReadContext:   resourceJobRead,
		DeleteContext: resourceJobDelete,

		Schema: map[string]*schema.Schema{
			"job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Job template ID",
				ForceNew:    true,
			},
			"limit": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "List of comma delimited hosts to limit job execution. Required ask_limit_on_launch set on job_template.",
			},
			"inventory_id": {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Computed:    true,
				Description: "Override Inventory ID. Required ask_inventory_on_launch set on job_template.",
				ForceNew:    true,
			},
			"extra_vars": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Override job template variables. YAML or JSON values are supported.",
				ForceNew:    true,
				StateFunc:   utils.Normalize,
			},
			"wait_for_completion": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Default:     false,
				Description: "Resource creation will wait for job completion.",
				ForceNew:    true,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func statusInstanceState(_ context.Context, svc *awx.JobService, id int) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := svc.GetJob(id, map[string]string{})
		return output, output.Status, err
	}
}

func jobTemplateLaunchWait(ctx context.Context, svc *awx.JobService, job *awx.JobLaunch, timeout time.Duration) error {

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"new", "pending", "waiting", "running"},
		Target:     []string{"successful"},
		Refresh:    statusInstanceState(ctx, svc, job.ID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

// JobTemplateLaunchData provides payload data used by the JobTemplateLaunch method
type JobTemplateLaunchData struct {
	Limit       string `json:"limit,omitempty"`
	InventoryID int    `json:"inventory,omitempty"`
	ExtraVars   string `json:"extra_vars,omitempty"`
}

func resourceJobTemplateLaunchCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)

	jobTemplateID := d.Get("job_template_id").(int)
	if _, err := client.JobTemplateService.GetJobTemplateByID(jobTemplateID, make(map[string]string)); err != nil {
		return utils.DiagNotFound(diagJobTemplateLaunchTitle, jobTemplateID, err)
	}

	data := JobTemplateLaunchData{
		Limit:       d.Get("limit").(string),
		InventoryID: d.Get("inventory_id").(int),
		ExtraVars:   d.Get("extra_vars").(string),
	}

	var iData map[string]interface{}
	idata, err := json.Marshal(data)
	if err != nil {
		return utils.DiagCreate(diagJobTemplateLaunchTitle, err)
	}
	if err := json.Unmarshal(idata, &iData); err != nil {
		return utils.DiagCreate(diagJobTemplateLaunchTitle, err)
	}

	res, err := client.JobTemplateService.Launch(jobTemplateID, iData, map[string]string{})
	if err != nil {
		return utils.DiagCreate(diagJobTemplateLaunchTitle, err)
	}

	// return resourceJobRead(ctx, d, m)
	d.SetId(strconv.Itoa(res.ID))

	if d.Get("wait_for_completion").(bool) {
		err = jobTemplateLaunchWait(ctx, client.JobService, res, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return utils.Diagf(
				"JobTemplate execution failure",
				"JobTemplateLaunch with ID %d and template ID %d, failed to complete %s", res.ID, d.Get("job_template_id").(int), err.Error(),
			)
		}
	}
	return nil
}

func resourceJobRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceJobDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	jobID, diags := utils.StateIDToInt("Delete Job", d)
	if diags.HasError() {
		return diags
	}
	if _, err := client.JobService.GetJob(jobID, map[string]string{}); err != nil {
		return utils.DiagNotFound(diagJobTemplateLaunchTitle, jobID, err)
	}

	d.SetId("")
	return nil
}
