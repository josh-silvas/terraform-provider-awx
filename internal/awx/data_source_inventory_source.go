package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

func dataSourceInventorySource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInventorySourceRead,
		Description: "Data source for Inventory Sources in AWX.",
		Schema: map[string]*schema.Schema{
			"inventory_source_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the inventory source.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the inventory source.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the inventory source.",
			},
			"inventory_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The inventory the inventory source belongs to.",
			},
			"credential_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The credential used for the inventory source.",
			},
			"source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The source of the inventory source.",
			},
			"source_vars": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The variables for the inventory source.",
			},
			"source_project_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The source project for the inventory source.",
			},
			"source_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The source path for the inventory source.",
			},
			"enabled_var": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The variable that determines if the inventory source is enabled.",
			},
			"enabled_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The value of the variable that determines if the inventory source is enabled.",
			},
			"host_filter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The host filter for the inventory source.",
			},
			"overwrite": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to overwrite the inventory source.",
			},
			"overwrite_vars": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to overwrite the inventory source variables.",
			},
			"update_on_launch": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to update the inventory source on launch.",
			},
			"update_cache_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The update cache timeout for the inventory source.",
			},
			"verbosity": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The verbosity for the inventory source.",
			},
			"execution_environment": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The execution environment for the inventory source.",
			},
		},
	}
}

func dataSourceInventorySourceRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id := d.Get("inventory_source_id").(int)

	res, err := client.InventorySourcesService.GetInventorySourceByID(id, map[string]string{})
	if err != nil {
		return utils.DiagFetch(diagInventorySourceTitle, id, err)
	}

	d = setInventorySourceResourceData(d, res)
	if err := d.Set("inventory_source_id", id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(id))
	return nil
}
