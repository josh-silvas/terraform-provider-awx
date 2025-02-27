package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagJobTemplateTitle = "Job Template"

//nolint:funlen
func resourceJobTemplate() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource `awx_job_template` manages job templates within AWX.",
		CreateContext: resourceJobTemplateCreate,
		ReadContext:   resourceJobTemplateRead,
		UpdateContext: resourceJobTemplateUpdate,
		DeleteContext: resourceJobTemplateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the job template.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The description of the job template.",
			},
			// Run, Check, Scan
			"job_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Can be one of: `run`, `check`, or `scan`",
			},
			"inventory_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The inventory ID to associate with the job template. If not set, `ask_inventory_on_launch` must be true.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The project ID to associate with the job template.",
			},
			"playbook": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The playbook to associate with the job template.",
			},
			"scm_branch": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"forks": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The number of forks to associate with the job template.",
			},
			"limit": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The limit to apply to filter hosts that run on this job template.",
			},
			//0,1,2,3,4,5
			"verbosity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "One of 0,1,2,3,4,5",
			},
			"extra_vars": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The extra variables to associate with the job template.",
			},
			"job_tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The job tags to associate with the job template.",
			},
			"force_handlers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Force handlers to run on the job template.",
			},
			"skip_tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The tags to skip on the job template.",
			},
			"start_at_task": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The task to start at on the job template.",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The timeout to associate with the job template. Default is 0",
			},
			"use_fact_cache": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Use the fact cache on the job template.",
			},
			"host_config_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ask_diff_mode_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_limit_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_tags_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_verbosity_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_inventory_on_launch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Defaults to false. Whether to ask for inventory on launch. If set to false, `inventory_id` must be set.",
			},
			"ask_variables_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_credential_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_scm_branch_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_execution_environment_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_labels_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_forks_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_job_slice_count_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_timeout_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_instance_group_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"survey_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"become_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"diff_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_skip_tags_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"allow_simultaneous": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"custom_virtualenv": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ask_job_type_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"execution_environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The selected execution environment that this playbook will be run in.",
			},
			"job_slice_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceJobTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	result, err := client.JobTemplateService.CreateJobTemplate(map[string]interface{}{
		"name":                                d.Get("name").(string),
		"description":                         d.Get("description").(string),
		"job_type":                            d.Get("job_type").(string),
		"inventory":                           utils.AtoiDefault(d.Get("inventory_id").(string), nil),
		"project":                             d.Get("project_id").(int),
		"playbook":                            d.Get("playbook").(string),
		"scm_branch":                          d.Get("scm_branch").(string),
		"forks":                               d.Get("forks").(int),
		"limit":                               d.Get("limit").(string),
		"verbosity":                           d.Get("verbosity").(int),
		"extra_vars":                          d.Get("extra_vars").(string),
		"job_tags":                            d.Get("job_tags").(string),
		"force_handlers":                      d.Get("force_handlers").(bool),
		"skip_tags":                           d.Get("skip_tags").(string),
		"start_at_task":                       d.Get("start_at_task").(string),
		"timeout":                             d.Get("timeout").(int),
		"use_fact_cache":                      d.Get("use_fact_cache").(bool),
		"host_config_key":                     d.Get("host_config_key").(string),
		"ask_scm_branch_on_launch":            d.Get("ask_scm_branch_on_launch").(bool),
		"ask_diff_mode_on_launch":             d.Get("ask_diff_mode_on_launch").(bool),
		"ask_variables_on_launch":             d.Get("ask_variables_on_launch").(bool),
		"ask_limit_on_launch":                 d.Get("ask_limit_on_launch").(bool),
		"ask_tags_on_launch":                  d.Get("ask_tags_on_launch").(bool),
		"ask_skip_tags_on_launch":             d.Get("ask_skip_tags_on_launch").(bool),
		"ask_job_type_on_launch":              d.Get("ask_job_type_on_launch").(bool),
		"ask_verbosity_on_launch":             d.Get("ask_verbosity_on_launch").(bool),
		"ask_inventory_on_launch":             d.Get("ask_inventory_on_launch").(bool),
		"ask_credential_on_launch":            d.Get("ask_credential_on_launch").(bool),
		"ask_execution_environment_on_launch": d.Get("ask_execution_environment_on_launch").(bool),
		"ask_labels_on_launch":                d.Get("ask_labels_on_launch").(bool),
		"ask_forks_on_launch":                 d.Get("ask_forks_on_launch").(bool),
		"ask_job_slice_count_on_launch":       d.Get("ask_job_slice_count_on_launch").(bool),
		"ask_timeout_on_launch":               d.Get("ask_timeout_on_launch").(bool),
		"ask_instance_group_on_launch":        d.Get("ask_instance_group_on_launch").(bool),
		"survey_enabled":                      d.Get("survey_enabled").(bool),
		"become_enabled":                      d.Get("become_enabled").(bool),
		"diff_mode":                           d.Get("diff_mode").(bool),
		"allow_simultaneous":                  d.Get("allow_simultaneous").(bool),
		"custom_virtualenv":                   utils.AtoiDefault(d.Get("custom_virtualenv").(string), nil),
		"execution_environment":               utils.AtoiDefault(d.Get("execution_environment").(string), nil),
		"job_slice_count":                     d.Get("job_slice_count").(int),
	}, map[string]string{})
	if err != nil {
		return utils.DiagCreate(diagJobTemplateTitle, err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceJobTemplateRead(ctx, d, m)
}

func resourceJobTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagJobTemplateTitle, d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)
	if _, err := client.JobTemplateService.GetJobTemplateByID(id, params); err != nil {
		return utils.DiagNotFound(diagJobTemplateTitle, id, err)
	}

	if _, err := client.JobTemplateService.UpdateJobTemplate(id, map[string]interface{}{
		"name":                                d.Get("name").(string),
		"description":                         d.Get("description").(string),
		"job_type":                            d.Get("job_type").(string),
		"inventory":                           utils.AtoiDefault(d.Get("inventory_id").(string), nil),
		"project":                             d.Get("project_id").(int),
		"playbook":                            d.Get("playbook").(string),
		"scm_branch":                          d.Get("scm_branch").(string),
		"forks":                               d.Get("forks").(int),
		"limit":                               d.Get("limit").(string),
		"verbosity":                           d.Get("verbosity").(int),
		"extra_vars":                          d.Get("extra_vars").(string),
		"job_tags":                            d.Get("job_tags").(string),
		"force_handlers":                      d.Get("force_handlers").(bool),
		"skip_tags":                           d.Get("skip_tags").(string),
		"start_at_task":                       d.Get("start_at_task").(string),
		"timeout":                             d.Get("timeout").(int),
		"use_fact_cache":                      d.Get("use_fact_cache").(bool),
		"host_config_key":                     d.Get("host_config_key").(string),
		"ask_scm_branch_on_launch":            d.Get("ask_scm_branch_on_launch").(bool),
		"ask_diff_mode_on_launch":             d.Get("ask_diff_mode_on_launch").(bool),
		"ask_variables_on_launch":             d.Get("ask_variables_on_launch").(bool),
		"ask_limit_on_launch":                 d.Get("ask_limit_on_launch").(bool),
		"ask_tags_on_launch":                  d.Get("ask_tags_on_launch").(bool),
		"ask_skip_tags_on_launch":             d.Get("ask_skip_tags_on_launch").(bool),
		"ask_job_type_on_launch":              d.Get("ask_job_type_on_launch").(bool),
		"ask_verbosity_on_launch":             d.Get("ask_verbosity_on_launch").(bool),
		"ask_inventory_on_launch":             d.Get("ask_inventory_on_launch").(bool),
		"ask_credential_on_launch":            d.Get("ask_credential_on_launch").(bool),
		"ask_execution_environment_on_launch": d.Get("ask_execution_environment_on_launch").(bool),
		"ask_labels_on_launch":                d.Get("ask_labels_on_launch").(bool),
		"ask_forks_on_launch":                 d.Get("ask_forks_on_launch").(bool),
		"ask_job_slice_count_on_launch":       d.Get("ask_job_slice_count_on_launch").(bool),
		"ask_timeout_on_launch":               d.Get("ask_timeout_on_launch").(bool),
		"ask_instance_group_on_launch":        d.Get("ask_instance_group_on_launch").(bool),
		"survey_enabled":                      d.Get("survey_enabled").(bool),
		"become_enabled":                      d.Get("become_enabled").(bool),
		"diff_mode":                           d.Get("diff_mode").(bool),
		"allow_simultaneous":                  d.Get("allow_simultaneous").(bool),
		"custom_virtualenv":                   utils.AtoiDefault(d.Get("custom_virtualenv").(string), nil),
		"execution_environment":               utils.AtoiDefault(d.Get("execution_environment").(string), nil),
		"job_slice_count":                     d.Get("job_slice_count").(int)}, map[string]string{}); err != nil {
		return utils.DiagUpdate(diagJobTemplateTitle, id, err)
	}

	return resourceJobTemplateRead(ctx, d, m)
}

func resourceJobTemplateRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagJobTemplateTitle, d)
	if diags.HasError() {
		return diags
	}

	res, err := client.JobTemplateService.GetJobTemplateByID(id, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound(diagJobTemplateTitle, id, err)
	}
	if res.ExtraVars != "" {
		res.ExtraVars = utils.Normalize(res.ExtraVars)
	}
	d = setJobTemplateResourceData(d, res)
	return nil
}

func resourceJobTemplateDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagJobTemplateTitle, d)
	if diags.HasError() {
		return diags
	}
	if _, err := client.JobTemplateService.DeleteJobTemplate(id); err != nil {
		return utils.DiagDelete(diagJobTemplateTitle, id, err)
	}
	d.SetId("")
	return nil
}

// nolint: gocyclo
func setJobTemplateResourceData(d *schema.ResourceData, r *awx.JobTemplate) *schema.ResourceData {
	if err := d.Set("allow_simultaneous", r.AllowSimultaneous); err != nil {
		fmt.Println("Error setting allow_simultaneous", err)
	}
	if err := d.Set("ask_credential_on_launch", r.AskCredentialOnLaunch); err != nil {
		fmt.Println("Error setting ask_credential_on_launch", err)
	}
	if err := d.Set("ask_job_type_on_launch", r.AskJobTypeOnLaunch); err != nil {
		fmt.Println("Error setting ask_job_type_on_launch", err)
	}
	if err := d.Set("ask_limit_on_launch", r.AskLimitOnLaunch); err != nil {
		fmt.Println("Error setting ask_limit_on_launch", err)
	}
	if err := d.Set("ask_skip_tags_on_launch", r.AskSkipTagsOnLaunch); err != nil {
		fmt.Println("Error setting ask_skip_tags_on_launch", err)
	}
	if err := d.Set("ask_tags_on_launch", r.AskTagsOnLaunch); err != nil {
		fmt.Println("Error setting ask_tags_on_launch", err)
	}
	if err := d.Set("ask_variables_on_launch", r.AskVariablesOnLaunch); err != nil {
		fmt.Println("Error setting ask_variables_on_launch", err)
	}
	if err := d.Set("ask_scm_branch_on_launch", r.AskScmBranchOnLaunch); err != nil {
		fmt.Println("Error setting ask_scm_branch_on_launch", err)
	}
	if err := d.Set("ask_diff_mode_on_launch", r.AskDiffModeOnLaunch); err != nil {
		fmt.Println("Error setting ask_diff_mode_on_launch", err)
	}
	if err := d.Set("ask_limit_on_launch", r.AskLimitOnLaunch); err != nil {
		fmt.Println("Error setting ask_limit_on_launch", err)
	}
	if err := d.Set("ask_tags_on_launch", r.AskTagsOnLaunch); err != nil {
		fmt.Println("Error setting ask_tags_on_launch", err)
	}
	if err := d.Set("ask_verbosity_on_launch", r.AskVerbosityOnLaunch); err != nil {
		fmt.Println("Error setting ask_verbosity_on_launch", err)
	}
	if err := d.Set("ask_inventory_on_launch", r.AskInventoryOnLaunch); err != nil {
		fmt.Println("Error setting ask_inventory_on_launch", err)
	}
	if err := d.Set("ask_execution_environment_on_launch", r.AskExecutionEnvironmentOnLaunch); err != nil {
		fmt.Println("Error setting ask_execution_environment_on_launch", err)
	}
	if err := d.Set("ask_labels_on_launch", r.AskLabelsOnLaunch); err != nil {
		fmt.Println("Error setting ask_labels_on_launch", err)
	}
	if err := d.Set("ask_forks_on_launch", r.AskForksOnLaunch); err != nil {
		fmt.Println("Error setting ask_forks_on_launch", err)
	}
	if err := d.Set("ask_job_slice_count_on_launch", r.AskJobSliceCountOnLaunch); err != nil {
		fmt.Println("Error setting ask_job_slice_count_on_launch", err)
	}
	if err := d.Set("ask_timeout_on_launch", r.AskTimeoutOnLaunch); err != nil {
		fmt.Println("Error setting ask_timeout_on_launch", err)
	}
	if err := d.Set("ask_instance_groups_on_launch", r.AskInstanceGroupsOnLaunch); err != nil {
		fmt.Println("Error setting ask_instance_groups_on_launch", err)
	}
	if err := d.Set("description", r.Description); err != nil {
		fmt.Println("Error setting description", err)
	}
	if err := d.Set("extra_vars", utils.Normalize(r.ExtraVars)); err != nil {
		fmt.Println("Error setting extra_vars", err)
	}
	if err := d.Set("force_handlers", r.ForceHandlers); err != nil {
		fmt.Println("Error setting force_handlers", err)
	}
	if err := d.Set("forks", r.Forks); err != nil {
		fmt.Println("Error setting forks", err)
	}
	if err := d.Set("host_config_key", r.HostConfigKey); err != nil {
		fmt.Println("Error setting host_config_key", err)
	}
	if err := d.Set("inventory_id", r.Inventory); err != nil {
		fmt.Println("Error setting inventory_id", err)
	}
	if err := d.Set("job_tags", r.JobTags); err != nil {
		fmt.Println("Error setting job_tags", err)
	}
	if err := d.Set("job_type", r.JobType); err != nil {
		fmt.Println("Error setting job_type", err)
	}
	if err := d.Set("job_slice_count", r.JobSliceCount); err != nil {
		fmt.Println("Error setting job_slice_count", err)
	}
	if err := d.Set("diff_mode", r.DiffMode); err != nil {
		fmt.Println("Error setting diff_mode", err)
	}
	if err := d.Set("custom_virtualenv", r.CustomVirtualenv); err != nil {
		fmt.Println("Error setting custom_virtualenv", err)
	}
	if err := d.Set("limit", r.Limit); err != nil {
		fmt.Println("Error setting limit", err)
	}
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	if err := d.Set("become_enabled", r.BecomeEnabled); err != nil {
		fmt.Println("Error setting become_enabled", err)
	}
	if err := d.Set("use_fact_cache", r.UseFactCache); err != nil {
		fmt.Println("Error setting use_fact_cache", err)
	}
	if err := d.Set("playbook", r.Playbook); err != nil {
		fmt.Println("Error setting playbook", err)
	}
	if err := d.Set("project_id", r.Project); err != nil {
		fmt.Println("Error setting project_id", err)
	}
	if err := d.Set("skip_tags", r.SkipTags); err != nil {
		fmt.Println("Error setting skip_tags", err)
	}
	if err := d.Set("start_at_task", r.StartAtTask); err != nil {
		fmt.Println("Error setting start_at_task", err)
	}
	if err := d.Set("survey_enabled", r.SurveyEnabled); err != nil {
		fmt.Println("Error setting survey_enabled", err)
	}
	if err := d.Set("verbosity", r.Verbosity); err != nil {
		fmt.Println("Error setting verbosity", err)
	}
	d.SetId(strconv.Itoa(r.ID))
	return d
}
