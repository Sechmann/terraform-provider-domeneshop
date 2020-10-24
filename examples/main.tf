terraform {
  required_providers {
    domeneshop = {
      versions = [
        "0.2"
      ]
      source = "hashicorp.com/edu/domeneshop"
    }
  }
}

provider "domeneshop" {
  version = "0.2"
}

#data "domeneshop_domains" "all" {
#}
#
#data "domeneshop_dns" "all" {
#  for_each = toset([for domain in data.domeneshop_domains.all.domains: tostring(domain.id)])
#  domain_id = tonumber(each.key)
#}
#
#output "all_domains" {
#  value = data.domeneshop_domains.all
#}
#
#output "all_records" {
#  value = data.domeneshop_dns.all
#}

resource "domeneshop_dns_record" "test" {
  domain_id = 1171517
  type = "A"
  host = "test"
  data = "13.37.13.37"
  ttl = 360
}
