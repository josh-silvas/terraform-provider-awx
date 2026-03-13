package awx

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagExecutionEnvironmentTitle = "Execution Environment"

func resourceExecutionEnvironment() *schema.Resource {
	return &schema.Resource{
		Description: "Execution Environment is a configuration that defines the runtime environment for a job template. " +
			"This includes the container image, organization, and credential used to run the job.",
		CreateContext: resourceExecutionEnvironmentsCreate,
		ReadContext:   resourceExecutionEnvironmentsRead,
		UpdateContext: resourceExecutionEnvironmentsUpdate,
		DeleteContext: resourceExecutionEnvironmentsDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the execution environment.",
			},
			"image": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The container image used for the execution environment.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The description of the execution environment.",
			},
			"organization": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The organization that the execution environment belongs to.",
			},
			"credential": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The credential used to access the execution environment.",
			},
			"pull": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				Description:  "Pull image before running?",
				ValidateFunc: validation.StringInSlice([]string{"", "always", "missing", "never"}, false),
			},
		},
	}
}

func resourceExecutionEnvironmentsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.ExecutionEnvironmentsService

	result, err := awxService.CreateExecutionEnvironment(map[string]interface{}{
		"name":         d.Get("name").(string),
		"image":        d.Get("image").(string),
		"description":  d.Get("description").(string),
		"organization": utils.AtoiDefault(d.Get("organization").(string), nil),
		"credential":   utils.AtoiDefault(d.Get("credential").(string), nil),
		"pull":         d.Get("pull").(string),
	}, map[string]string{})
	if err != nil {
		log.Printf("Fail to Create ExecutionEnvironment %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create ExecutionEnvironments",
			Detail:   fmt.Sprintf("ExecutionEnvironments with name %s, failed to create %s", d.Get("name").(string), err),
		})
		return diags
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceExecutionEnvironmentsRead(ctx, d, m)
}

func resourceExecutionEnvironmentsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Update ExecutionEnvironments", d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)

	if _, err := client.ExecutionEnvironmentsService.GetExecutionEnvironmentByID(id, params); err != nil {
		return utils.DiagNotFound(diagExecutionEnvironmentTitle, id, err)
	}

	if _, err := client.ExecutionEnvironmentsService.UpdateExecutionEnvironment(id, map[string]interface{}{
		"name":         d.Get("name").(string),
		"image":        d.Get("image").(string),
		"description":  d.Get("description").(string),
		"organization": utils.AtoiDefault(d.Get("organization").(string), nil),
		"credential":   utils.AtoiDefault(d.Get("credential").(string), nil),
		"pull":         d.Get("pull").(string),
	}, map[string]string{}); err != nil {
		return utils.DiagUpdate(diagExecutionEnvironmentTitle, id, err)
	}

	return resourceExecutionEnvironmentsRead(ctx, d, m)
}

func resourceExecutionEnvironmentsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.ExecutionEnvironmentsService
	id, diags := utils.StateIDToInt("Read ExecutionEnvironments", d)
	if diags.HasError() {
		return diags
	}

	res, err := awxService.GetExecutionEnvironmentByID(id, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound(diagExecutionEnvironmentTitle, id, err)

	}
	d = setExecutionEnvironmentsResourceData(d, res)
	return nil
}

func resourceExecutionEnvironmentsDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Delete ExecutionEnvironment", d)
	if diags.HasError() {
		return diags
	}

	if _, err := client.ExecutionEnvironmentsService.DeleteExecutionEnvironment(id); err != nil {
		return utils.DiagDelete(diagExecutionEnvironmentTitle, id, err)
	}
	d.SetId("")
	return diag.Diagnostics{}
}

func setExecutionEnvironmentsResourceData(d *schema.ResourceData, r *awx.ExecutionEnvironment) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	if err := d.Set("image", r.Image); err != nil {
		fmt.Println("Error setting image", err)
	}
	if err := d.Set("description", r.Description); err != nil {
		fmt.Println("Error setting description", err)
	}
	if err := d.Set("organization", r.Organization); err != nil {
		fmt.Println("Error setting organization", err)
	}
	if err := d.Set("credential", r.Credential); err != nil {
		fmt.Println("Error setting credential", err)
	}
	if err := d.Set("pull", r.Pull); err != nil {
		fmt.Println("Error setting pull", err)
	}
	d.SetId(strconv.Itoa(r.ID))
	return d
}
