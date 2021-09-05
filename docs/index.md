---
page_title: "Cisco IOS-XE Provider"
subcategory: ""
description: |-
  
---

# Cisco IOS-XE Provider

Pre-release. Models/signatures could change. Bug reports/pulls/request for additional resources welcome on git repo.

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

PS. I hate yang

<!-- ## Schema -->
