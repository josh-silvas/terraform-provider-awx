---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "awx_notification_template Data Source - terraform-provider-awx"
subcategory: ""
description: |-
  Data source for AWX Notification Template
---

# awx_notification_template (Data Source)

Data source for AWX Notification Template

## Example Usage

```terraform
data "awx_notification_template" "example" {
  name = "my-notification-template"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (Number) The ID of the Notification Template
- `name` (String) The name of the Notification Template
