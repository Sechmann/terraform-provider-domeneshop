package domeneshop

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"terraform-provider-domeneshop/domeneshop/model"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDomains() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainsRead,
		Schema: map[string]*schema.Schema{
			"domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
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
						"transferred_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDomainsRead(ctx context.Context, data *schema.ResourceData, providerClient interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := providerClient.(*http.Client)

	response, err := client.Get("https://api.domeneshop.no/v0/domains")
	if err != nil {
		return diag.FromErr(err)
	}

	defer response.Body.Close()
	var domains []model.Domain
	err = json.NewDecoder(response.Body).Decode(&domains)
	if err != nil {
		return diag.FromErr(err)
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return diag.Errorf("listing domains gave unsuccessful status code: %v", response)
	}

	err = data.Set("domains", flattenDomains(&domains))
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenDomains(domains *[]model.Domain) []interface{} {
	if domains != nil {
		flatDomains := make([]interface{}, len(*domains), len(*domains))

		for i, domain := range *domains {
			flatDomain := make(map[string]interface{})

			flatDomain["domain"] = domain.Domain
			flatDomain["expiry_date"] = domain.ExpiryDate
			flatDomain["id"] = domain.Id
			flatDomain["nameservers"] = domain.Nameservers
			flatDomain["registrant"] = domain.Registrant
			flatDomain["renew"] = domain.Renew
			flatDomain["services_dns"] = domain.Services.Dns
			flatDomain["services_email"] = domain.Services.Email
			flatDomain["services_registrar"] = domain.Services.Registrar
			flatDomain["services_webhotell"] = domain.Services.Webhotel
			flatDomain["status"] = domain.Status

			flatDomains[i] = flatDomain
		}

		return flatDomains
	}

	return make([]interface{}, 0)
}

func logIfError(err error) {
	if err != nil {
		log.Printf("closing response body: %v\n", err)
	}
}
