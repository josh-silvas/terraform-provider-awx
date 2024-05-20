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

func resourceOrganization() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource Organization is used to manage the organization in AWX",
		CreateContext: resourceOrganizationsCreate,
		ReadContext:   resourceOrganizationsRead,
		UpdateContext: resourceOrganizationsUpdate,
		DeleteContext: resourceOrganizationsDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the organization",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The description of the organization",
			},
			"max_hosts": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Maximum number of hosts allowed to be managed by this organization",
			},
			"custom_virtualenv": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Local absolute file path containing a custom Python virtualenv to use",
			},
			"default_environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The default execution environment for jobs run by this organization.",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		//
		//Timeouts: &schema.ResourceTimeout{
		//	Create: schema.DefaultTimeout(1 * time.Minute),
		//	Update: schema.DefaultTimeout(1 * time.Minute),
		//	Delete: schema.DefaultTimeout(1 * time.Minute),
		//},
	}
}

func resourceOrganizationsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	result, err := client.OrganizationsService.CreateOrganization(map[string]interface{}{
		"name":                d.Get("name").(string),
		"description":         d.Get("description").(string),
		"max_hosts":           d.Get("max_hosts").(int),
		"custom_virtualenv":   d.Get("description").(string),
		"default_environment": d.Get("default_environment").(string),
	}, map[string]string{})
	if err != nil {
		return utils.DiagCreate(diagOrganizationTitle, err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceOrganizationsRead(ctx, d, m)
}

func resourceOrganizationsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Update Organizations", d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)
	if _, err := client.OrganizationsService.GetOrganizationsByID(id, params); err != nil {
		return utils.DiagNotFound(diagOrganizationTitle, id, err)
	}

	if _, err := client.OrganizationsService.UpdateOrganization(id, map[string]interface{}{
		"name":                d.Get("name").(string),
		"description":         d.Get("description").(string),
		"max_hosts":           d.Get("max_hosts").(int),
		"custom_virtualenv":   d.Get("description").(string),
		"default_environment": d.Get("default_environment").(string),
	}, map[string]string{}); err != nil {
		return utils.DiagUpdate(diagOrganizationTitle, id, err)
	}

	return resourceOrganizationsRead(ctx, d, m)
}

func resourceOrganizationsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Read Organizations", d)
	if diags.HasError() {
		return diags
	}

	res, err := client.OrganizationsService.GetOrganizationsByID(id, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound(diagOrganizationTitle, id, err)

	}
	d = setOrganizationsResourceData(d, res)
	return nil
}

func resourceOrganizationsDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Delete Organization", d)
	if diags.HasError() {
		return diags
	}

	if _, err := client.OrganizationsService.DeleteOrganization(id); err != nil {
		return utils.DiagDelete(diagOrganizationTitle, id, err)
	}
	d.SetId("")
	return diags
}

func setOrganizationsResourceData(d *schema.ResourceData, r *awx.Organization) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	if err := d.Set("description", r.Description); err != nil {
		fmt.Println("Error setting description", err)
	}
	if err := d.Set("max_hosts", r.MaxHosts); err != nil {
		fmt.Println("Error setting max_hosts", err)
	}
	if err := d.Set("custom_virtualenv", r.CustomVirtualenv); err != nil {
		fmt.Println("Error setting custom_virtualenv", err)
	}
	d.SetId(strconv.Itoa(r.ID))
	return d
}
