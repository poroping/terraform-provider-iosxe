---
page_title: "Cisco IOS-XE Provider"
subcategory: ""
description: |-
  
---

# Cisco IOS-XE Provider

Pre-release. Not stable! Models/signatures will change. Bug reports/pulls/request for additional resources welcome on git repo.



Enable restconf on your device. [DevNetDocs](https://developer.cisco.com/docs/ios-xe/#!enabling-restconf-on-ios-xe)

```
conf t
ip http secure-server
restconf
end
wr
```


## Example Usage

```terraform
provider "iosxe" {
  host     = "https://192.168.1.1"
  username = "cisco"
  password = "cisco"
  insecure = true
}
```

## Example L3 VLAN

```terraform
provider "iosxe" {
  host     = "https://192.168.1.1"
  username = "cisco"
  password = "cisco"
  insecure = true
}

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
```
PS. I hate yang

<!-- ## Schema -->
