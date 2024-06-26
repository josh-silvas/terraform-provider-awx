---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "awx_credentials Data Source - terraform-provider-awx"
subcategory: ""
description: |-
  Data source for fetching credentials from AWX
---

# awx_credentials (Data Source)

Data source for fetching credentials from AWX

## Example Usage

```terraform
data "awx_credentials" "example" {
  name = "example"
  id   = 10
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `credentials` (List of Object) (see [below for nested schema](#nestedatt--credentials))
- `id` (String) The ID of this resource.

<a id="nestedatt--credentials"></a>
### Nested Schema for `credentials`

Read-Only:

- `description` (String)
- `id` (Number)
- `kind` (String)
- `name` (String)
- `username` (String)
