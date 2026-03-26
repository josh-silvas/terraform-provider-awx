data "awx_credential_vault" "example" {
  credential_id = 10
}

output "credential_name" {
  value = data.awx_credential_vault.example.name
}
