---
page_title: "iosxe_vlan Resource - terraform-provider-iosxe"
subcategory: ""
description: |-
  Manage a VLAN.
---

# Resource `iosxe_vlan`

Manage a VLAN.

## Example Usage

```terraform
resource "iosxe_vlan" "example" {
  vlanid = 420
  name   = "IoT"
}

output "debug" {
  value = iosxe_vlan.example
}
```

## Schema

### Optional

- **vlanid** (Int, Required) VLAN ID.
- **name** (String, Optional) VLAN name.


