---
page_title: "iosxe_l2_vlan Resource - terraform-provider-iosxe"
subcategory: ""
description: |-
  Manage a L2 VLAN.
---

# Resource `iosxe_l2_vlan`

Manage a L2 VLAN.

## Example Usage

```terraform
resource "iosxe_l2_vlan" "example" {
  vlanid = 420
  name   = "IoT"
}

output "debug" {
  value = iosxe_l2_vlan.example
}
```

## Argument Reference

- **vlanid** (Int, Required) VLAN ID.
- **name** (String, Optional) VLAN name.


