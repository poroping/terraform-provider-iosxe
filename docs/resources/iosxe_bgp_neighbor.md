---
page_title: "iosxe_bgp_neighbor Resource - terraform-provider-iosxe"
subcategory: ""
description: |-
  Manage a BGP Neighbor.
---

# Resource `iosxe_bgp_neighbor`

Manage a BGP Neighbor.

## Example Usage

```terraform
resource "iosxe_bgp_router" "example" {
  as                   = 65420
  log_neighbor_changes = true
}

resource "iosxe_bgp_neighbor" "example" {
  as                = iosxe_bgp_router.example.as
  ip                = "7.7.7.7"
  remote_as         = 8899
  default_originate = true

  prefix_list {
    name      = "pl_test"
    direction = "in"
  }
}

output "debug" {
  value = iosxe_bgp_neighbor.example
}
```

## Argument Reference

- **as** (Int, Required) ASN.
- **ip** (String, Required) IP address of BGP peer.
- **remote_as** (String, Required) Remote peer ASN.
- **default_originate** (Bool, Optional) Originate default route.
- **description** (String, Optional) Description.
- **ebgp_multihop** (Int, Optional) EBG multi-hop.
- **local_as** (Int Optional) Override local ASN.
- **prefix_list** (Optional) Block defined below.
- **remove_private_as** (Bool, Optional) Remove private ASNs.
- **shutdown** (Bool, Optional) Neighbor status.
- **soft_reconfiguration** (String, Optional) Soft reconfiguration.
- **timers** (Optional) Block defined below.

The **prefix_list** block contains:

- **direction** (String, Required) Direction list applied.
- **name** (String, Required) Name of prefix-list.

The **timers** block contains:

- **keepalive_interval** (Int, Optional) Keepalive interval.
- **holdtime** (Int, Optional) Hold down time.
- **minimum_neighbor_hold** (Int, Optional) Min hold time from neighbor.

## Attribute Reference

In addition to all the above arguments, the following attributes are exported:
- **id** - resource identifier.


