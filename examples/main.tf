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

data "domeneshop_domains" "all" {
}

data "domeneshop_dns" "all" {
  for_each = toset([for domain in data.domeneshop_domains.all.domains: tostring(domain.id)])
  domain_id = tonumber(each.value)
}

# Returns all coffees
output "all_domains" {
  value = data.domeneshop_domains.all
}

output "all_records" {
  value = data.domeneshop_dns.all
}
