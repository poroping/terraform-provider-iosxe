---
page_title: "iosxe_interface_port_channel Resource - terraform-provider-iosxe"
subcategory: ""
description: |-
  Manage a Port Channel interface.
---

# Resource `iosxe_interface_port_channel`

Manage a Port Channel interface.

## Example Usage

```terraform
resource "iosxe_interface_port_channel" "example" {
  name        = "56"
  description = "totallyterraformed"
}

output "debug" {
  value = iosxe_interface_port_channel.example
}

```

## Argument Reference

- **name** (String, Required) Interface name.
- **description** (String, Optional) Interface description.


