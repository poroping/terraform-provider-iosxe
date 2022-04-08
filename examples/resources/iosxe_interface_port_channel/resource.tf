resource "iosxe_interface_port_channel" "example" {
  name        = "56"
  description = "totallyterraformed"
}

output "debug" {
  value = iosxe_interface_port_channel.example
}
