data "awx_credential_machine" "example" {
  credential_id = 10
}

output "credential_name" {
  value = data.awx_credential_machine.example.name
}
