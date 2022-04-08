resource "iosxe_vrf" "example" {
  name        = "FOOBAR"
  description = "ACC-TEST"
  rd          = "566:4560"

  address_family {
    ip_version = 4
  }
}

resource "iosxe_l2_vlan" "example" {
  vlanid = 666
  name   = "IoT"
}

resource "iosxe_interface_vlan" "example" {
  vlanid      = iosxe_l2_vlan.example.vlanid
  description = "totallyterraformed"
  ip          = "192.168.66.6/24"
  shutdown    = false
  vrf         = iosxe_vrf.example.name

  secondary_ip {
    ip = "10.55.2.1/30"
  }
}

output "debug" {
  value = iosxe_interface_vlan.example
}
