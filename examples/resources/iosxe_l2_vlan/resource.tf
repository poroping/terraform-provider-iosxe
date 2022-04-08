resource "iosxe_l2_vlan" "example" {
  vlanid = 420
  name   = "IoT"
}

output "debug" {
  value = iosxe_l2_vlan.example
}