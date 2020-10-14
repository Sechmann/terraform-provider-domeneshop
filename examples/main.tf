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

data "domeneshop_domains" "all" {}

# Returns all coffees
output "all_domains" {
  value = data.domeneshop_domains.all
}