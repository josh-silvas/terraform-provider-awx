data "awx_inventory_source" "example" {
  inventory_source_id = 10
}

output "inventory_source_name" {
  value = data.awx_inventory_source.example.name
}
