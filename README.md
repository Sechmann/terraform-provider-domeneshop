Domeneshop terraform provider
---

Data Sources:
- `domeneshop_domain`

Resources:
- `domeneshop_dns_record`

### Usage
```terraform
data "domeneshop_domain" "desperate_solutions" {
  domain = "desperate.solutions"
}
```

```terraform
# Add k8s.desperate.solutions A record pointing to 13.37.13.37
resource "domeneshop_dns_record" "k8s" {
  domain_id = data.domeneshop_domain.desperate_solutions.id

  type = "A"
  host = "test"
  data = "13.37.13.37"
  ttl  = 360
}
```


