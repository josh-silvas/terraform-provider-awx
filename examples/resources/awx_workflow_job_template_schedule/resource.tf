data "awx_organization" "example" {
  name = "Default"
}

resource "awx_inventory" "example" {
  name            = "Example Inventory"
  organization_id = data.awx_organization.example.id
}

resource "awx_workflow_job_template" "example" {
  name            = "workflow-job"
  organization_id = data.awx_organization.example.id
  inventory_id    = awx_inventory.example.id
}

resource "awx_workflow_job_template_schedule" "default" {
  workflow_job_template_id = awx_workflow_job_template.example.id
  name                     = "schedule-test"
  rrule                    = "DTSTART;TZID=Europe/Paris:20211214T120000 RRULE:INTERVAL=1;FREQ=DAILY"
  limit                    = "webservers"
  extra_data               = <<EOL
organization_name: testorg
EOL
}
