---
page_title: "iosxe_svi Resource - terraform-provider-iosxe"
subcategory: ""
description: |-
  Manage a SVI.
---

# Resource `iosxe_svi`

Manage a SVI.

## Example Usage

```terraform
resource "iosxe_svi" "example" {
  vlanid      = 666
  description = "totallyterraformed"
  ip          = "192.168.66.6/24"
  shutdown    = false

  secondary_ip {
    ip = "10.1.2.1/24"
  }
}

output "debug" {
  value = iosxe_svi.example
}
```

## Argument Reference

- **vlanid** (Int, Required) VLAN ID.
- **description** (String, Optional) Interface description.
- **ip** (String, Required) IP in CIDR notation.
- **secondary_ip** (Optional) Block defined below.
- **shutdown** (Bool, Optional) Interface status.

The **secondary_ip** block contains:

- **ip** (String, Optional) IP in CIDR notation.

## Attribute Reference

In addition to all the above arguments, the following attributes are exported:
- **id** - resource identifier.
- **name** - interface name.


