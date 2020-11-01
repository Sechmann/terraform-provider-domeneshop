# DNS Record Resource

Manage a DNS record

## Example Usage

```hcl
data "domeneshop_domain" "example_com" {
  domain = "example.com"
}

resource "domeneshop_dns_record" "wildcard" {
  domain_id = data.domeneshop_domain.example_com

  type = "A"
  host = "*"
  data = "11.22.33.44"
  ttl  = 300
}
```

## Argument Reference
* `type` - (Required) One of: [`A`,`AAAA`,`CNAME`,`MX`,`SRV`,`TXT`]
* `host` - (Required) The subdomain this record is for, `@` for the top level.
* `data` - (Required) Contents of the record, depends on TYPE. i.e. for type `A`, `"11.22.33.44"`
* `ttl`  - (Required) Time to live in seconds, i.e. `300`
* `priority` - (Optional) For Only applicable when type is `SRV`/`MX`,
* `weight` - (Optional) Only applicable when type is `SRV`
* `port` - (Optional) Only applicable when type is `SRV`

## Attribute Reference
* `id` - The id of this dns record

## Import

Domeneshop DNS Record can be imported using the domain id and record id, e.g.

```
$ terraform import domeneshop_dns_record.record  1337/1338
```

