---
page_title: "iosxe_interface_port_channel_subinterface Resource - terraform-provider-iosxe"
subcategory: ""
description: |-
  Manage a L3 Port-Channel Subinterface interface.
---

# Resource `iosxe_interface_port_channel_subinterface`

Manage a L3 Port-Channel Subinterface interface.

## Example Usage

```terraform
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

```

## Argument Reference

- **vlanid** (Int, Required) Dot1q encapsulation VLAN.
- **description** (String, Optional) Interface description.
- **ip** (String, Required) IP in CIDR notation.
- **secondary_ip** (Optional) Block defined below.
- **shutdown** (Bool, Optional) Interface status.
- **name** (String, Required) Interface name.

The **secondary_ip** block contains:

- **ip** (String, Optional) IP in CIDR notation.

## Attribute Reference

In addition to all the above arguments, the following attributes are exported:
- **id** - resource identifier.


