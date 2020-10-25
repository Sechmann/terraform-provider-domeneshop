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

data "domeneshop_domain" "desperate_solutions" {
  domain = "desperate.solutions"
}

resource "domeneshop_dns_record" "test" {
  domain_id = data.domeneshop_domain.desperate_solutions.id
  type = "A"
  host = "test"
  data = "13.37.13.37"
  ttl = 360
}
