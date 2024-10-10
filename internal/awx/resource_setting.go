package awx

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

func resourceSetting() *schema.Resource {
	return &schema.Resource{
		Description: "This resource configure generic AWX settings.\n" +
			"Please note that resource deletion only delete object from terraform state and do not reset setting to his initial value.\n\n" +
			"See available settings list here: https://docs.ansible.com/ansible-tower/latest/html/towerapi/api_ref.html#/Settings/Settings_settings_update",
		CreateContext: resourceSettingUpdate,
		ReadContext:   resourceSettingRead,
		DeleteContext: resourceSettingDelete,
		UpdateContext: resourceSettingUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of setting to modify",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Value to be modified for given setting.",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceSettingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)

	if _, err := client.SettingService.GetSettingsBySlug("all", make(map[string]string)); err != nil {
		return utils.DiagCreate("Settings Update", err)
	}

	var (
		mapDecoded     map[string]interface{}
		arrayDecoded   []interface{}
		formattedValue interface{} // nolint:typecheck
	)

	name := d.Get("name").(string)
	value := d.Get("value").(string)

	// Attempt to unmarshall string into a map
	if err := json.Unmarshal([]byte(value), &mapDecoded); err != nil {
		// Attempt to unmarshall string into an array
		if err = json.Unmarshal([]byte(value), &arrayDecoded); err != nil {
			formattedValue = value
		} else {
			formattedValue = arrayDecoded
		}
	} else {
		formattedValue = mapDecoded
	}

	payload := map[string]interface{}{
		name: formattedValue,
	}

	if _, err := client.SettingService.UpdateSettings("all", payload, make(map[string]string)); err != nil {
		return utils.DiagUpdate("Settings Update", formattedValue, err)
	}

	d.SetId(name)
	return resourceSettingRead(ctx, d, m)
}

func resourceSettingRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	res, err := client.SettingService.GetSettingsBySlug("all", make(map[string]string))
	if err != nil {
		return utils.DiagFetch("Settings Read", "all", err)
	}

	if err := d.Set("name", d.Id()); err != nil {
		return diag.FromErr(err)
	}

	val, err := (*res)[d.Id()].MarshalJSON()
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("value", string(val)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSettingDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
