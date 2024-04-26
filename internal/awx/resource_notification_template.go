package awx

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

const diagNotificationTemplateTitle = "Notification Template"

func resourceNotificationTemplate() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource `awx_notification_template` manages notification templates within an AWX organization.",
		CreateContext: resourceNotificationTemplateCreate,
		ReadContext:   resourceNotificationTemplateRead,
		UpdateContext: resourceNotificationTemplateUpdate,
		DeleteContext: resourceNotificationTemplateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the notification template.",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The organization ID to associate with the notification template.",
			},
			"notification_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of notification template.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the notification template.",
			},
			"notification_configuration": {
				Type:        schema.TypeSet,
				Optional:    true,
				Default:     nil,
				Description: "Notification configuration settings based on the notification type.",
				// documented at OPTIONS /api/v2/notification_templates/
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// generic settings (re-used across notification types)
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HTTP or SMTP password.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "SMTP or IRC server port.",
						},
						"token": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Slack or PagerDuty authentication token.",
						},
						"username": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HTTP or SMTP username.",
						},
						"use_ssl": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Use SSL for IRC or SMTP connections.",
						},
						// email notification-type settings
						"host": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SMTP server hostname.",
						},
						"recipients": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The email address(es) to send notifications to.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"sender": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The email address to send notifications from.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "SMTP server timeout.",
							Default:     30,
						},
						"use_tls": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Use TLS for SMTP connections.",
						},
						// slack notification-type settings
						"channels": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The Slack channel(s) to send notifications to.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"hex_color": {
							Type:     schema.TypeString,
							Optional: true,
						},
						// twilio notification-type settings
						"account_sid": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Twilio account SID",
						},
						"account_token": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Twilio account token",
						},
						"from_number": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Twilio from number",
						},
						"to_numbers": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Twilio to numbers",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						// pagerduty notification-type settings
						"client_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "PagerDuty client name",
						},
						"service_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "PagerDuty service key",
						},
						"subdomain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "PagerDuty subdomain",
						},
						// grafana notification-type settings
						"grafana_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Grafana API key",
						},
						"grafana_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Grafana URL",
						},
						// webhook notification-type settings
						"disable_ssl_verification": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"headers": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "The headers to include in the webhook request.",
						},
						"http_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HTTP method to use when sending the webhook request.",
						},
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HTTP webhook URL.",
						},
						// mattermost notification-type settings
						"mattermost_no_verify_ssl": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"mattermost_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Mattermost URL.",
						},
						// rocketchat notification-type settings
						"rocketchat_no_verify_ssl": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"rocketchat_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The RocketChat URL.",
						},
						// irc notification-type settings
						"server": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The IRC server hostname.",
						},
						"nickname": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"targets": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The IRC channel(s) to send notifications to.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"messages": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The description of the notification template. Options are `started`, `success`, `error`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"started": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "The message to send when the job starts.",
						},
						"success": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "The message to send when the job starts.",
						},
						"error": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "The message to send when the job starts.",
						},
					},
				},
			},
		},
	}
}

func resourceNotificationTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	payload := map[string]interface{}{
		"name":              d.Get("name").(string),
		"description":       d.Get("description").(string),
		"organization":      d.Get("organization_id").(string),
		"notification_type": d.Get("notification_type").(string),
	}

	notificationConfig := d.Get("notification_configuration").(*schema.Set).List()
	if len(notificationConfig) != 0 {
		payload["notification_configuration"] = notificationConfig[0].(map[string]interface{})
	}

	messages := d.Get("messages").(*schema.Set).List()
	if len(messages) != 0 {
		payload["messages"] = messages[0].(map[string]interface{})
	}
	result, err := client.NotificationTemplatesService.Create(payload, map[string]string{})
	if err != nil {
		return utils.DiagCreate("NotificationTemplate", err)
	}

	d.SetId(strconv.Itoa(result.ID))
	time.Sleep(time.Second * 3)
	return resourceNotificationTemplateRead(ctx, d, m)
}

func resourceNotificationTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Update NotificationTemplate", d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)
	if _, err := client.NotificationTemplatesService.GetByID(id, params); err != nil {
		return utils.DiagNotFound(diagNotificationTemplateTitle, id, err)
	}
	payload := map[string]interface{}{
		"name":              d.Get("name").(string),
		"description":       d.Get("description").(string),
		"organization":      d.Get("organization_id").(string),
		"notification_type": d.Get("notification_type").(string),
	}

	notificationConfig := d.Get("notification_configuration").(*schema.Set).List()
	if len(notificationConfig) != 0 {
		payload["notification_configuration"] = notificationConfig[0].(map[string]interface{})
	}

	messages := d.Get("messages").(*schema.Set).List()
	if len(messages) != 0 {
		payload["messages"] = messages[0].(map[string]interface{})
	}
	if _, err := client.NotificationTemplatesService.Update(id, payload, map[string]string{}); err != nil {
		return utils.DiagUpdate(diagNotificationTemplateTitle, id, err)
	}
	time.Sleep(time.Second * 3)
	return resourceNotificationTemplateRead(ctx, d, m)
}

func resourceNotificationTemplateRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt("Read notification_template", d)
	if diags.HasError() {
		return diags
	}

	res, err := client.NotificationTemplatesService.GetByID(id, make(map[string]string))
	if err != nil {
		return utils.DiagNotFound(diagNotificationTemplateTitle, id, err)

	}
	d = setNotificationTemplateResourceData(d, res)
	return nil
}

func resourceNotificationTemplateDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	id, diags := utils.StateIDToInt(diagNotificationTemplateTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := client.NotificationTemplatesService.Delete(id); err != nil {
		return utils.DiagDelete(diagNotificationTemplateTitle, id, err)
	}
	d.SetId("")
	return nil
}

func setNotificationTemplateResourceData(d *schema.ResourceData, r *awx.NotificationTemplate) *schema.ResourceData {
	if err := d.Set("name", r.Name); err != nil {
		fmt.Println("Error setting name", err)
	}
	if err := d.Set("description", r.Description); err != nil {
		fmt.Println("Error setting description", err)
	}
	if err := d.Set("organization", r.Organization); err != nil {
		fmt.Println("Error setting organization", err)
	}
	if err := d.Set("notification_type", r.NotificationType); err != nil {
		fmt.Println("Error setting notification_type", err)
	}
	if err := d.Set("notification_configuration", schema.NewSet(func(i interface{}) int { return len(i.(map[string]interface{})) }, []interface{}{r.NotificationConfiguration})); err != nil {
		fmt.Println("Error setting notification_configuration", err)
	}
	if err := d.Set("messages", schema.NewSet(func(i interface{}) int { return len(i.(map[string]interface{})) }, []interface{}{r.Messages})); err != nil {
		fmt.Println("Error setting messages", err)
	}
	d.SetId(strconv.Itoa(r.ID))
	return d
}
