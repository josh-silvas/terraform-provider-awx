package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func dataSourceInstanceGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get the instance group ID by name, " +
			"which can be used in resources such as awx_job_template_instance_groups.",
		ReadContext: dataSourceInstanceGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the instance group to look up.",
			},
			"is_container_group": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the instance group is a container group.",
			},
			"capacity": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The capacity of the instance group.",
			},
		},
	}
}

func dataSourceInstanceGroupRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	name := d.Get("name").(string)

	groups, _, err := client.InstanceGroupsService.ListInstanceGroups(map[string]string{
		"name": name,
	})
	if err != nil {
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "Unable to fetch instance groups",
				Detail:   fmt.Sprintf("Unable to fetch instance groups from AWX API: %s", err),
			},
		}
	}

	for _, g := range groups {
		if g.Name == name {
			d.SetId(strconv.Itoa(g.ID))
			if err := d.Set("is_container_group", g.IsContainerGroup); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("capacity", g.Capacity); err != nil {
				return diag.FromErr(err)
			}
			return nil
		}
	}

	return diag.Diagnostics{
		{
			Severity: diag.Error,
			Summary:  "Instance group not found",
			Detail:   fmt.Sprintf("No instance group with name %q was found.", name),
		},
	}
}
