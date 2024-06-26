---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "awx_credential_google_compute_engine Resource - terraform-provider-awx"
subcategory: ""
description: |-
  awx_credential_google_compute_engine manages Google Compute Engine credentials in AWX.
---

# awx_credential_google_compute_engine (Resource)

`awx_credential_google_compute_engine` manages Google Compute Engine credentials in AWX.

## Example Usage

```terraform
resource "awx_organization" "example" {
  name = "example"
}

resource "awx_credential_google_compute_engine" "example" {
  name            = "awx-gce-credential"
  organization_id = awx_organization.example.id
  description     = "This is a GCE credential"
  ssh_key_data    = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDQ6"
  username        = "admin"
  project         = "my-project"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the credential.
- `organization_id` (Number) The organization ID this credential belongs to.
- `project` (String) The project to use for the credential.
- `ssh_key_data` (String, Sensitive) The SSH key data to use for the credential.
- `username` (String) The username to use for the credential.

### Optional

- `description` (String) The description of the credential.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# Order can be imported by specifying the numeric identifier.
terraform import awx_credential_google_compute_engine.example 540
```
