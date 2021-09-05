---
page_title: "iosxe_bgp_router Resource - terraform-provider-iosxe"
subcategory: ""
description: |-
  Manage a BGP router instance.
---

# Resource `iosxe_bgp_router`

Manage a BGP router instance.

## Example Usage

```terraform
resource "iosxe_bgp_router" "example" {
  as                   = 65420
  log_neighbor_changes = true
}

output "debug" {
  value = iosxe_bgp_router.example
}
```

## Argument Reference

- **as** (Int, Required) ASN.
- **log_neighbor_changes** (Bool, Optional) Log neighbor changes.

## Attribute Reference

In addition to all the above arguments, the following attributes are exported:
- **id** - resource identifier.


