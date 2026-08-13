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

func resourceInventory() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource Inventory is used to define an inventory in AWX",
		CreateContext: resourceInventoryCreate,
		ReadContext:   resourceInventoryRead,
		DeleteContext: resourceInventoryDelete,
		UpdateContext: resourceInventoryUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the inventory",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The description of the inventory",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The organization id of the inventory",
			},
			"kind": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The kind of the inventory",
			},
			"host_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The host filter of the inventory",
			},
			"variables": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				StateFunc:   utils.Normalize,
				Description: "The variables of the inventory",
			},
			"prevent_instance_group_fallback": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Only run jobs on instance groups assigned to this inventory, without falling back to the organization instance groups",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceInventoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	result, err := client.InventoriesService.CreateInventory(map[string]interface{}{
		"name":                            d.Get("name").(string),
		"organization":                    d.Get("organization_id").(string),
		"description":                     d.Get("description").(string),
		"kind":                            d.Get("kind").(string),
		"host_filter":                     d.Get("host_filter").(string),
		"variables":                       d.Get("variables").(string),
		"prevent_instance_group_fallback": d.Get("prevent_instance_group_fallback").(bool),
	}, map[string]string{})
	if err != nil {
		return utils.DiagCreate(diagInventoryTitle, err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceInventoryRead(ctx, d, m)

}

func resourceInventoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagInventoryTitle, d)
	if diags.HasError() {
		return diags
	}
	if _, err := client.InventoriesService.UpdateInventory(id, map[string]interface{}{
		"name":                            d.Get("name").(string),
		"organization":                    d.Get("organization_id").(string),
		"description":                     d.Get("description").(string),
		"kind":                            d.Get("kind").(string),
		"host_filter":                     d.Get("host_filter").(string),
		"variables":                       d.Get("variables").(string),
		"prevent_instance_group_fallback": d.Get("prevent_instance_group_fallback").(bool),
	}, nil); err != nil {
		return utils.DiagUpdate(diagInventoryTitle, id, err)
	}

	return resourceInventoryRead(ctx, d, m)

}

func resourceInventoryRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, err := strconv.Atoi(d.Id())
	id, diags := utils.StateIDToInt(diagInventoryTitle, d)
	if diags.HasError() {
		return diags
	}
	r, err := client.InventoriesService.GetInventory(id, map[string]string{})
	if err != nil {
		return utils.DiagFetch(diagInventoryTitle, id, err)
	}
	d = setInventoryResourceData(d, r)
	return nil
}

func resourceInventoryDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventoriesService
	id, diags := utils.StateIDToInt(diagInventoryTitle, d)
	if diags.HasError() {
		return diags
	}
	if _, err := awxService.DeleteInventory(id); err != nil {
		return utils.DiagDelete(diagInventoryTitle, id, err)
	}
	d.SetId("")
	return nil
}

func setInventoryResourceData(d *schema.ResourceData, r *awx.Inventory) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	if err := d.Set("organization_id", strconv.Itoa(r.Organization)); err != nil {
		fmt.Println("Error setting organization_id", err)
	}
	if err := d.Set("description", r.Description); err != nil {
		fmt.Println("Error setting description", err)
	}
	if err := d.Set("kind", r.Kind); err != nil {
		fmt.Println("Error setting kind", err)
	}
	if err := d.Set("host_filter", r.HostFilter); err != nil {
		fmt.Println("Error setting host_filter", err)
	}
	if err := d.Set("variables", utils.Normalize(r.Variables)); err != nil {
		fmt.Println("Error setting variables", err)
	}
	if err := d.Set("prevent_instance_group_fallback", r.PreventInstanceGroupFallback); err != nil {
		fmt.Println("Error setting prevent_instance_group_fallback", err)
	}
	d.SetId(strconv.Itoa(r.ID))
	return d
}
