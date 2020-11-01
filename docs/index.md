# Domeneshop Provider
Unofficial, alpha. Use at your own risk.

This provider manages DNS records for the domeneshop registrar.

https://api.domeneshop.no/docs/

## Example Usage

```hcl
provider "domeneshop" {
  version = "0.3"
}
```

## Argument Reference

Env vars `DOMENESHOP_TOKEN` and `DOMENESHOP_SECRET` must be set, get yours here: https://domene.shop/admin?view=api
