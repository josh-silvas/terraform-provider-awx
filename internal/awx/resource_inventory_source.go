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

const diagInventorySourceTitle = "Inventory Source"

//nolint:funlen
func resourceInventorySource() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource Inventory Source is used to manage inventory sources in AWX.",
		CreateContext: resourceInventorySourceCreate,
		ReadContext:   resourceInventorySourceRead,
		UpdateContext: resourceInventorySourceUpdate,
		DeleteContext: resourceInventorySourceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the inventory source.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the inventory source.",
			},
			"enabled_var": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The variable that determines if the inventory source is enabled.",
			},
			"enabled_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The value of the variable that determines if the inventory source is enabled.",
			},
			"overwrite": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to overwrite the inventory source.",
			},
			"overwrite_vars": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to overwrite the inventory source variables.",
			},
			"update_on_launch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to update the inventory source on launch.",
			},
			"inventory_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The inventory to use for the inventory source.",
			},
			"credential_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The credential to use for the inventory source.",
			},
			"source": {
				Type:        schema.TypeString,
				Default:     "scm",
				Optional:    true,
				Description: "The source of the inventory source.",
			},
			"source_vars": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The variables for the inventory source.",
			},
			"host_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The host filter for the inventory source.",
			},
			"update_cache_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "The update cache timeout for the inventory source.",
			},
			"verbosity": {
				Type:        schema.TypeInt,
				Default:     1,
				Optional:    true,
				Description: "The verbosity for the inventory source. [0,1,2,3]",
			},
			"execution_environment": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The selected execution environment that this inventory will be run in.",
			},
			// obsolete schema added so terraform doesn't break
			// these don't do anything in later versions of AWX! Update your code.
			"source_regions": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "[Obsolete] The source regions for the inventory source.",
			},
			"instance_filters": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "[Obsolete] The instance filters for the inventory source.",
			},
			"group_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "[Obsolete] The group by for the inventory source.",
			},
			"source_project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "[Obsolete] The source project for the inventory source.",
			},
			"source_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "[Obsolete] The source path for the inventory source.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceInventorySourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)

	payload := map[string]interface{}{
		"name":                 d.Get("name").(string),
		"description":          d.Get("description").(string),
		"enabled_var":          d.Get("enabled_var").(string),
		"enabled_value":        d.Get("enabled_value").(string),
		"overwrite":            d.Get("overwrite").(bool),
		"overwrite_vars":       d.Get("overwrite_vars").(bool),
		"update_on_launch":     d.Get("update_on_launch").(bool),
		"inventory":            d.Get("inventory_id").(int),
		"source":               d.Get("source").(string),
		"source_vars":          d.Get("source_vars").(string),
		"host_filter":          d.Get("host_filter").(string),
		"update_cache_timeout": d.Get("update_cache_timeout").(int),
		"verbosity":            d.Get("verbosity").(int),
		// obsolete schema added so terraform doesn't break
		// these don't do anything in later versions of AWX! Update your code.
		"source_regions":   d.Get("source_regions").(string),
		"instance_filters": d.Get("instance_filters").(string),
		"group_by":         d.Get("group_by").(string),
		"source_path":      d.Get("source_path").(string),
	}
	if _, ok := d.GetOk("credential_id"); ok {
		payload["credential"] = d.Get("credential_id").(int)
	}
	if _, ok := d.GetOk("source_project_id"); ok {
		payload["source_project"] = d.Get("source_project_id").(int)
	}
	if _, ok := d.GetOk("execution_environment"); ok {
		payload["execution_environment"] = d.Get("execution_environment").(int)
	}

	result, err := client.InventorySourcesService.CreateInventorySource(payload, map[string]string{})
	if err != nil {
		return utils.DiagCreate(diagInventorySourceTitle, err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceInventorySourceRead(ctx, d, m)

}

func resourceInventorySourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventorySourcesService
	id, diags := utils.StateIDToInt(diagInventorySourceTitle, d)
	if diags.HasError() {
		return diags
	}

	payload := map[string]interface{}{
		"name":                 d.Get("name").(string),
		"description":          d.Get("description").(string),
		"enabled_var":          d.Get("enabled_var").(string),
		"enabled_value":        d.Get("enabled_value").(string),
		"overwrite":            d.Get("overwrite").(bool),
		"overwrite_vars":       d.Get("overwrite_vars").(bool),
		"update_on_launch":     d.Get("update_on_launch").(bool),
		"inventory":            d.Get("inventory_id").(int),
		"source":               d.Get("source").(string),
		"source_vars":          d.Get("source_vars").(string),
		"host_filter":          d.Get("host_filter").(string),
		"update_cache_timeout": d.Get("update_cache_timeout").(int),
		"verbosity":            d.Get("verbosity").(int),
		// obsolete schema added so terraform doesn't break
		// these don't do anything in later versions of AWX! Update your code.
		"source_regions":   d.Get("source_regions").(string),
		"instance_filters": d.Get("instance_filters").(string),
		"group_by":         d.Get("group_by").(string),
		"source_path":      d.Get("source_path").(string),
	}
	if _, ok := d.GetOk("credential_id"); ok {
		payload["credential"] = d.Get("credential_id").(int)
	}
	if _, ok := d.GetOk("source_project_id"); ok {
		payload["source_project"] = d.Get("source_project_id").(int)
	}
	if _, ok := d.GetOk("execution_environment"); ok {
		payload["execution_environment"] = d.Get("execution_environment").(int)
	}

	if _, err := awxService.UpdateInventorySource(id, payload, nil); err != nil {
		return utils.DiagUpdate(diagInventorySourceTitle, id, err)
	}

	return resourceInventorySourceRead(ctx, d, m)
}

func resourceInventorySourceDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagInventorySourceTitle, d)
	if diags.HasError() {
		return diags
	}
	if _, err := client.InventorySourcesService.DeleteInventorySource(id); err != nil {
		return utils.DiagDelete(diagInventorySourceTitle, id, err)
	}
	d.SetId("")
	return nil
}

func resourceInventorySourceRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagInventorySourceTitle, d)
	if diags.HasError() {
		return diags
	}
	res, err := client.InventorySourcesService.GetInventorySourceByID(id, make(map[string]string))
	if err != nil {
		return utils.DiagFetch(diagInventorySourceTitle, id, err)
	}
	d = setInventorySourceResourceData(d, res)
	return nil
}

func setInventorySourceResourceData(d *schema.ResourceData, r *awx.InventorySource) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	if err := d.Set("description", r.Description); err != nil {
		fmt.Println("Error setting description", err)
	}
	if err := d.Set("enabled_var", r.EnabledVar); err != nil {
		fmt.Println("Error setting enabled_var", err)
	}
	if err := d.Set("enabled_value", r.EnabledValue); err != nil {
		fmt.Println("Error setting enabled_value", err)
	}
	if err := d.Set("overwrite", r.Overwrite); err != nil {
		fmt.Println("Error setting overwrite", err)
	}
	if err := d.Set("overwrite_vars", r.OverwriteVars); err != nil {
		fmt.Println("Error setting overwrite_vars", err)
	}
	if err := d.Set("update_on_launch", r.UpdateOnLaunch); err != nil {
		fmt.Println("Error setting update_on_launch", err)
	}
	if err := d.Set("inventory_id", r.Inventory); err != nil {
		fmt.Println("Error setting inventory_id", err)
	}
	if err := d.Set("credential_id", r.Credential); err != nil {
		fmt.Println("Error setting credential_id", err)
	}
	if err := d.Set("source", r.Source); err != nil {
		fmt.Println("Error setting source", err)
	}
	if err := d.Set("source_vars", utils.Normalize(r.SourceVars)); err != nil {
		fmt.Println("Error setting source_vars", err)
	}
	if err := d.Set("host_filter", r.HostFilter); err != nil {
		fmt.Println("Error setting host_filter", err)
	}
	if err := d.Set("update_cache_timeout", r.UpdateCacheTimeout); err != nil {
		fmt.Println("Error setting update_cache_timeout", err)
	}
	if err := d.Set("verbosity", r.Verbosity); err != nil {
		fmt.Println("Error setting verbosity", err)
	}
	if err := d.Set("execution_environment", r.ExecutionEnvironment); err != nil {
		fmt.Println("Error setting execution_environment", err)
	}
	if err := d.Set("source_project_id", r.SourceProject); err != nil {
		fmt.Println("Error setting source_project_id", err)
	}
	if err := d.Set("source_path", r.SourcePath); err != nil {
		fmt.Println("Error setting source_path", err)
	}
	// obsolete schema added so terraform doesn't break
	// these don't do anything in later versions of AWX! Update your code.
	if err := d.Set("source_regions", r.SourceRegions); err != nil {
		fmt.Println("Error setting source_regions", err)
	}
	if err := d.Set("instance_filters", r.InstanceFilters); err != nil {
		fmt.Println("Error setting instance_filters", err)
	}
	if err := d.Set("group_by", r.GroupBy); err != nil {
		fmt.Println("Error setting group_by", err)
	}
	if err := d.Set("source_project_id", r.SourceProject); err != nil {
		fmt.Println("Error setting source_project_id", err)
	}
	if err := d.Set("source_path", r.SourcePath); err != nil {
		fmt.Println("Error setting source_path", err)
	}

	return d
}
