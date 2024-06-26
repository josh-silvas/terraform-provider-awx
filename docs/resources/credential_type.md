---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "awx_credential_type Resource - terraform-provider-awx"
subcategory: ""
description: |-
  Resource awx_credential_type manages credential types within an AWX instance.
---

# awx_credential_type (Resource)

Resource `awx_credential_type` manages credential types within an AWX instance.

## Example Usage

```terraform
resource "awx_credential_type" "example" {
  name        = "Example Credential Type"
  description = "Example Credential Type"
  kind        = "cloud"
  inputs      = "{\"fields\": [{\"id\": \"password\", \"label\": \"Password\", \"type\": \"string\", \"secret\": true}]}"
  injectors   = "{\"extra_vars\": {\"password\": \"password\"}}"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `injectors` (String) Injectors for this credential type.
- `inputs` (String) Inputs for this credential type.
- `name` (String) Name of this credential type.

### Optional

- `description` (String) Optional description of this credential type.
- `kind` (String) Can be one of: `cloud` or `net`

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# Order can be imported by specifying the numeric identifier.
terraform import awx_credential_type.example 580
```
