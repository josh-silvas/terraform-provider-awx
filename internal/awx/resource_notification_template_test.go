package awx

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func notificationTemplateConfig(url string, extra map[string]interface{}) map[string]interface{} {
	config := map[string]interface{}{
		"url":         url,
		"http_method": "POST",
		"username":    "admin",
	}
	for k, v := range extra {
		config[k] = v
	}
	return map[string]interface{}{
		"name":                       "test",
		"organization_id":            "1",
		"notification_type":          "webhook",
		"notification_configuration": []interface{}{config},
	}
}

func notificationTemplateRemote(url string) *awx.NotificationTemplate {
	return &awx.NotificationTemplate{
		ID:               1,
		Name:             "test",
		NotificationType: "webhook",
		NotificationConfiguration: map[string]interface{}{
			"url":                      url,
			"http_method":              "POST",
			"username":                 "admin",
			"password":                 "$encrypted$",
			"disable_ssl_verification": false,
		},
	}
}

func notificationConfigurationSentOnUpdate(t *testing.T, stateConfig, newConfig map[string]interface{}, remote *awx.NotificationTemplate) map[string]interface{} {
	t.Helper()

	r := resourceNotificationTemplate()
	d := schema.TestResourceDataRaw(t, r.Schema, stateConfig)
	d.SetId("1")
	setNotificationTemplateResourceData(d, remote)

	state := d.State()
	diff, err := schema.InternalMap(r.Schema).Diff(context.Background(), state, terraform.NewResourceConfigRaw(newConfig), nil, nil, false)
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	if diff == nil {
		t.Fatal("expected a diff, got none")
	}

	updateData, err := schema.InternalMap(r.Schema).Data(state, diff)
	if err != nil {
		t.Fatalf("data: %v", err)
	}
	sent := updateData.Get("notification_configuration").([]interface{})
	if len(sent) != 1 {
		t.Fatalf("expected 1 notification_configuration element, got %d", len(sent))
	}
	return sent[0].(map[string]interface{})
}

func TestNotificationTemplateUpdateSendsNewValue(t *testing.T) {
	sent := notificationConfigurationSentOnUpdate(t,
		notificationTemplateConfig("https://old.example.com", map[string]interface{}{"password": "secret"}),
		notificationTemplateConfig("https://new.example.com", map[string]interface{}{"password": "secret"}),
		notificationTemplateRemote("https://old.example.com"),
	)

	if sent["url"] != "https://new.example.com" {
		t.Errorf("url = %q, want https://new.example.com", sent["url"])
	}
	if sent["password"] != "secret" {
		t.Errorf("password not preserved")
	}
}

func TestNotificationTemplateUpdateSendsNewValueWithUnmanagedSecret(t *testing.T) {
	sent := notificationConfigurationSentOnUpdate(t,
		notificationTemplateConfig("https://old.example.com", nil),
		notificationTemplateConfig("https://new.example.com", nil),
		notificationTemplateRemote("https://old.example.com"),
	)

	if sent["url"] != "https://new.example.com" {
		t.Errorf("url = %q, want https://new.example.com", sent["url"])
	}
	if sent["http_method"] != "POST" {
		t.Errorf("http_method = %q, want POST", sent["http_method"])
	}
}

func TestNotificationTemplateReadWithoutMessages(t *testing.T) {
	r := resourceNotificationTemplate()
	d := schema.TestResourceDataRaw(t, r.Schema, notificationTemplateConfig("https://example.com", nil))

	setNotificationTemplateResourceData(d, notificationTemplateRemote("https://example.com"))

	if messages := d.Get("messages").([]interface{}); len(messages) != 0 {
		t.Errorf("messages = %v, want empty", messages)
	}
}
