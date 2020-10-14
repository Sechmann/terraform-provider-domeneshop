package domeneshop

import (
	"context"
	"encoding/base64"
	domeneshopapi "github.com/VegarM/domeneshop-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"domeneshop_domains": dataSourceDomains(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
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

	cfg := domeneshopapi.NewConfiguration()
	cfg.AddDefaultHeader("Authorization", basicAuth(token, secret))

	client := domeneshopapi.NewAPIClient(cfg)

	return client, diags
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
