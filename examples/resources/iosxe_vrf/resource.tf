resource "iosxe_vrf" "example" {
  name        = "FOOBAR"
  description = "TEST-VRF"
  rd          = "566:4560"

  address_family {
    ip_version = 4
  }

  address_family {
    ip_version = 6
  }

  route_target {
    community = "export"
    rt        = "6969:111"
  }

  route_target {
    community = "import"
    rt        = "778:420"
  }

  route_target {
    community = "import"
    rt        = "778:421"
  }

}

output "debug" {
  value = iosxe_vrf.example
}
