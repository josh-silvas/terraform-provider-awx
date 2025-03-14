---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "awx_job_template Resource - terraform-provider-awx"
subcategory: ""
description: |-
  Resource awx_job_template manages job templates within AWX.
---

# awx_job_template (Resource)

Resource `awx_job_template` manages job templates within AWX.

## Example Usage

```terraform
data "awx_organization" "default" {
  name = "Default"
}

resource "awx_project" "example" {
  name            = "example-ansible-main"
  organization_id = data.awx_organization.default.id
  scm_type        = "git"
  scm_url         = "git@github.com/josh-silvas/example-ansible.git"
  scm_branch      = "main"
}

resource "awx_inventory" "example" {
  name            = "Example Inventory"
  organization_id = data.awx_organization.default.id
}

resource "awx_job_template" "example" {
  name           = "baseconfig"
  job_type       = "run"
  inventory_id   = awx_inventory.example.id
  project_id     = awx_project.example.id
  playbook       = "master-configure-system.yml"
  become_enabled = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `job_type` (String) Can be one of: `run`, `check`, or `scan`
- `name` (String) The name of the job template.
- `project_id` (Number) The project ID to associate with the job template.

### Optional

- `allow_simultaneous` (Boolean)
- `ask_credential_on_launch` (Boolean)
- `ask_diff_mode_on_launch` (Boolean)
- `ask_execution_environment_on_launch` (Boolean)
- `ask_forks_on_launch` (Boolean)
- `ask_instance_group_on_launch` (Boolean)
- `ask_inventory_on_launch` (Boolean) Defaults to false. Whether to ask for inventory on launch. If set to false, `inventory_id` must be set.
- `ask_job_slice_count_on_launch` (Boolean)
- `ask_job_type_on_launch` (Boolean)
- `ask_labels_on_launch` (Boolean)
- `ask_limit_on_launch` (Boolean)
- `ask_scm_branch_on_launch` (Boolean)
- `ask_skip_tags_on_launch` (Boolean)
- `ask_tags_on_launch` (Boolean)
- `ask_timeout_on_launch` (Boolean)
- `ask_variables_on_launch` (Boolean)
- `ask_verbosity_on_launch` (Boolean)
- `become_enabled` (Boolean)
- `custom_virtualenv` (String)
- `description` (String) The description of the job template.
- `diff_mode` (Boolean)
- `execution_environment` (String) The selected execution environment that this playbook will be run in.
- `extra_vars` (String) The extra variables to associate with the job template.
- `force_handlers` (Boolean) Force handlers to run on the job template.
- `forks` (Number) The number of forks to associate with the job template.
- `host_config_key` (String)
- `inventory_id` (String) The inventory ID to associate with the job template. If not set, `ask_inventory_on_launch` must be true.
- `job_slice_count` (Number)
- `job_tags` (String) The job tags to associate with the job template.
- `limit` (String) The limit to apply to filter hosts that run on this job template.
- `playbook` (String) The playbook to associate with the job template.
- `scm_branch` (String)
- `skip_tags` (String) The tags to skip on the job template.
- `start_at_task` (String) The task to start at on the job template.
- `survey_enabled` (Boolean)
- `timeout` (Number) The timeout to associate with the job template. Default is 0
- `use_fact_cache` (Boolean) Use the fact cache on the job template.
- `verbosity` (Number) One of 0,1,2,3,4,5

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# Order can be imported by specifying the numeric identifier.
terraform import awx_job_template.example 650
```
