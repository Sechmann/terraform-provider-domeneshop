package domeneshop

import (
	"context"
	domeneshopapi "github.com/VegarM/domeneshop-go"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	baseContext context.Context
)

func init() {
	token, found := os.LookupEnv("DOMENESHOP_TOKEN")
	if !found {
		log.Fatal("Reading env DOMENESHOP_TOKEN: not set")
	}

	secret, found := os.LookupEnv("DOMENESHOP_SECRET")
	if !found {
		log.Fatal("Reading env DOMENESHOP_SECRET: not set")
	}

	baseContext = context.WithValue(context.Background(), domeneshopapi.ContextBasicAuth, domeneshopapi.BasicAuth{
		UserName: token,
		Password: secret,
	})
}

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

func baseContextWithTimeout(duration time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(baseContext, duration)
}

func dataSourceDomainsRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	cfg := domeneshopapi.NewConfiguration()
	client := domeneshopapi.NewAPIClient(cfg)
	ctx, cancel := baseContextWithTimeout(time.Second * 15)
	defer cancel()

	domains, response, err := client.DomainsApi.GetDomains(ctx, nil)
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

func flattenDomains(domains *[]domeneshopapi.Domain) []interface{} {
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
