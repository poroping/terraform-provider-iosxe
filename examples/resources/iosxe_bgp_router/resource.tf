resource "iosxe_bgp_router" "example" {
  as                   = 65420
  log_neighbor_changes = true
}

output "debug" {
  value = iosxe_bgp_router.example
}