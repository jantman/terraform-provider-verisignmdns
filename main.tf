resource "verisignmdns_rr" "foo" {
  record_name = "foo.example.com"
  record_type = "A"
  record_data = "1.2.3.4"
}
