---
page_title: "iosxe_vlan Data Source - terraform-provider-iosxe"
subcategory: ""
description: |-
  Get VLAN information.
---

# Data Source `iosxe_vlan`

Get VLAN information.

## Example Usage

```terraform
data "iosxe_vlan" "example" {
  vlanid = 666
}

output "debug" {
  value = data.iosxe_vlan.example
}
```

## Schema

### Required

- **vlanid** (String, Required) VLAN ID.


