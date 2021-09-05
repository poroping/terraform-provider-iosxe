resource "iosxe_vlan" "example" {
  vlanid = 420
  name   = "IoT"
}

output "debug" {
  value = iosxe_vlan.example
}