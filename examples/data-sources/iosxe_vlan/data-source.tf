data "iosxe_vlan" "example" {
    vlanid = 666
}

output "debug" {
    value = data.iosxe_vlan.example
}