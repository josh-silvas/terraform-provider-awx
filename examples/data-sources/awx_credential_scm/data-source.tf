data "awx_credential_scm" "example" {
  credential_id = 10
}

output "credential_name" {
  value = data.awx_credential_scm.example.name
}
