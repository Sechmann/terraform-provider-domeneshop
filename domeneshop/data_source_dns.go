package domeneshop

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"terraform-provider-domeneshop/domeneshop/model"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDNS() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSRead,
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDNSRead(ctx context.Context, data *schema.ResourceData, providerClient interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := providerClient.(*http.Client)

	domainId, ok := data.GetOk("domain_id")
	if !ok {
		return diag.Errorf("required attribute domain_id not set")
	}

	response, err := client.Get(fmt.Sprintf("https://api.domeneshop.no/v0/domains/%d/dns", domainId))
	if err != nil {
		return diag.FromErr(err)
	}
	defer response.Body.Close()

	var records []model.DnsRecord
	err = json.NewDecoder(response.Body).Decode(&records)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to get domain info",
				Detail:   fmt.Sprintf("Unable to get domain records, error: %v", err),
			}}
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return diag.Errorf("listing records gave unsuccessful status code: %v", response)
	}

	err = data.Set("records", flattenDNSRecords(&records))
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenDNSRecord(record *model.DnsRecord) interface{} {
	flatRecord := make(map[string]interface{})

	flatRecord["id"] = record.Id
	flatRecord["host"] = record.Host
	flatRecord["ttl"] = record.Ttl
	flatRecord["type"] = record.Type
	flatRecord["data"] = record.Data

	switch record.Type {
	case "MX":
		flatRecord["priority"] = record.Priority
	case "SRV":
		flatRecord["priority"] = record.Priority
		flatRecord["weight"] = record.Weight
		flatRecord["port"] = record.Port
	}
	return flatRecord
}

func flattenDNSRecords(records *[]model.DnsRecord) []interface{} {
	if records != nil {
		flatRecords := make([]interface{}, len(*records), len(*records))

		for i, record := range *records {
			flatRecords[i] = flattenDNSRecord(&record)
		}

		return flatRecords
	}

	return make([]interface{}, 0)
}
