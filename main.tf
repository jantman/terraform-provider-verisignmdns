resource "verisignmdns_rr" "foo" {
  account_id  = "9999999"
  zone_name   = "example.com"
  record_name = "foo.example.com"
  record_type = "A"
  record_data = "1.2.3.4"
}
