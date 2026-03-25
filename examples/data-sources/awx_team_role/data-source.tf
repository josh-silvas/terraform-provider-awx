data "awx_organization" "example" {
  name = "example"
}

data "awx_team" "example" {
  name            = "example"
  organization_id = data.awx_organization.example.id
}

data "awx_team_role" "example" {
  name    = "Member"
  team_id = awx_team.example.id
}
