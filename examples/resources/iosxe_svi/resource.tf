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