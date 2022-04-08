resource "iosxe_vrf" "example" {
  name        = "FOOBAR"
  description = "ACC-TEST"
  rd          = "566:4560"

  address_family {
    ip_version = 4
  }
}


resource "iosxe_interface_port_channel_subinterface" "example" {
  name        = "69.421"
  vlanid      = 421
  description = "totallyterraformed"
  ip          = "192.1.1.1/29"
  shutdown    = false
  vrf         = iosxe_vrf.example.name

  secondary_ip {
    ip = "10.55.6.1/30"
  }
}

output "debug" {
  value = iosxe_interface_port_channel_subinterface.example
}
