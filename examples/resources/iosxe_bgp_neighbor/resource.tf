resource "iosxe_bgp_router" "example" {
  as                   = 65420
  log_neighbor_changes = true
}

resource "iosxe_bgp_neighbor" "example" {
  as                = iosxe_bgp_router.example.as
  ip                = "7.7.7.7"
  remote_as         = 8899
  default_originate = true

  prefix_list {
    name      = "pl_test"
    direction = "in"
  }
}

output "debug" {
  value = iosxe_bgp_neighbor.example
}