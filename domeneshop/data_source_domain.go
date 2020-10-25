package domeneshop

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/VegarM/domeneshop-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"strconv"
	"terraform-provider-domeneshop/domeneshop/api"
)

func dataSourceDomain() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"expiry_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"nameservers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"registrant": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"renew": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"services_dns": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"services_email": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"services_registrar": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"services_webhotell": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"registered_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDomainRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(*http.Client)

	name := d.Get("domain").(string)

	domains, err := getDomains(client)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, domain := range domains {
		if domain.Domain == name {
			diags = append(diags, setDomainData(&domain, d)...)
			id := strconv.Itoa(int(domain.Id))
			d.SetId(id)
			break
		}
	}

	return diags
}

func setDomainData(domain *domeneshop.Domain, d *schema.ResourceData) diag.Diagnostics {
	var errs []error

	errs = append(errs, d.Set("domain", domain.Domain))
	errs = append(errs, d.Set("expiry_date", domain.ExpiryDate))
	errs = append(errs, d.Set("id", domain.Id))
	errs = append(errs, d.Set("nameservers", domain.Nameservers))
	errs = append(errs, d.Set("registrant", domain.Registrant))
	errs = append(errs, d.Set("renew", domain.Renew))
	errs = append(errs, d.Set("services_dns", domain.Services.Dns))
	errs = append(errs, d.Set("services_email", domain.Services.Email))
	errs = append(errs, d.Set("services_registrar", domain.Services.Registrar))
	errs = append(errs, d.Set("services_webhotell", domain.Services.Webhotel))
	errs = append(errs, d.Set("status", domain.Status))
	errs = append(errs, d.Set("registered_date", domain.RegisteredDate))

	var diags diag.Diagnostics
	for _, err := range errs {
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

func getDomains(client *http.Client) ([]domeneshop.Domain, error) {
	response, err := client.Get(api.Domains())
	if err != nil {
		return nil, fmt.Errorf("HTTP get domains: %w", err)
	}
	defer closeBody(response.Body)

	var domains []domeneshop.Domain
	err = json.NewDecoder(response.Body).Decode(&domains)

	if err != nil {
		return nil, fmt.Errorf("decoding domains: %w", err)
	}

	return domains, nil
}
