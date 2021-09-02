resource "iosxe_svi" "example" {
    vlanid = 666
    description = "totallyterraformed"
    ip = "192.168.66.6/24"
}

output "debug" {
    value = iosxe_svi.example
}