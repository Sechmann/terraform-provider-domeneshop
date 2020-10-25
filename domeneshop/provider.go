package domeneshop

import (
	"context"
	"encoding/base64"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"time"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOMENESHOP_TOKEN", nil),
			},
			"secret": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("DOMENESHOP_SECRET", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"domeneshop_dns_record": resourceDNSRecord(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"domeneshop_domain": dataSourceDomain(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("token").(string)
	secret := d.Get("secret").(string)

	var diags diag.Diagnostics

	if len(token) == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create domeneshop client",
			Detail:   "Unable to auth user for authenticated domeneshop client: `token` not set",
		})
	}

	if len(secret) == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create domeneshop client",
			Detail:   "Unable to auth user for authenticated domeneshop client: `secret` not set",
		})
	}

	if len(diags) > 0 {
		return nil, diags
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
		Transport: &AddHeaderTransport{
			T: http.DefaultTransport,
			Headers: map[string]string{
				"Authorization": basicAuth(token, secret),
			},
		},
	}

	return client, diags
}

type AddHeaderTransport struct {
	T       http.RoundTripper
	Headers map[string]string
}

func (adt *AddHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range adt.Headers {
		req.Header.Add(key, value)
	}
	return adt.T.RoundTrip(req)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
