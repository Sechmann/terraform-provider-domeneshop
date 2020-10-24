package domeneshop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"terraform-provider-domeneshop/domeneshop/api"
	"terraform-provider-domeneshop/domeneshop/model"
	"time"
)

type IdResponse struct {
	Id int `json:"id"`
}

func resourceDNSRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSRecordCreate,
		ReadContext:   resourceDNSRecordRead,
		UpdateContext: resourceDNSUpdate,
		DeleteContext: resourceDNSRecordDelete,
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data": {
				Type:     schema.TypeString,
				Required: true,
			},
			"priority": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceDNSRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client := m.(*http.Client)

	domainId := d.Get("domain_id").(int)

	record, err := dnsRecordFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

	buffer := new(bytes.Buffer)
	err = json.NewEncoder(buffer).Encode(record)
	if err != nil {
		return diag.FromErr(err)
	}
	response, err := client.Post(api.DNSRecords(domainId), "application/json", buffer)
	if err != nil {
		return diag.FromErr(err)
	}
	defer closeBody(response.Body)

	switch response.StatusCode {
	case 201:
		var parsed IdResponse
		err := json.NewDecoder(response.Body).Decode(&parsed)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(strconv.Itoa(parsed.Id))

		// refresh state
		diags = append(diags, resourceDNSRecordRead(ctx, d, m)...)
	default:
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return diag.FromErr(err)
		}
		return []diag.Diagnostic{{
			Severity: diag.Error,
			Summary:  "unexpected status code from create operation",

			Detail: fmt.Sprintf("expected 201, got %d. Response : %v", response.StatusCode, string(b)),
		}}
	}

	return diags
}

func resourceDNSRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*http.Client)

	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if domainId, ok := d.GetOk("domain_id"); ok {
		response, err := client.Get(fmt.Sprintf("https://api.domeneshop.no/v0/domains/%d/dns/%d", domainId, recordId))
		if err != nil {
			return diag.FromErr(err)
		}
		defer closeBody(response.Body)

		var record model.DnsRecord
		err = json.NewDecoder(response.Body).Decode(&record)
		if err != nil {
			return diag.FromErr(err)
		}

		var errs []error
		errs = append(errs, d.Set("type", record.Type))
		errs = append(errs, d.Set("host", record.Host))
		errs = append(errs, d.Set("data", record.Data))
		errs = append(errs, d.Set("ttl", record.Ttl))

		switch record.Type {
		case "SRV":
			errs = append(errs, d.Set("priority", record.Priority))
			errs = append(errs, d.Set("port", record.Port))
			errs = append(errs, d.Set("weight", record.Weight))
		case "MX":
			errs = append(errs, d.Set("priority", record.Priority))
		}

		for _, err = range errs {
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
		
		if diags != nil && diags.HasError(){
			return diags
		}
	} else {
		return diag.FromErr(fmt.Errorf("domain_id is required"))
	}

	return diags
}

func resourceDNSUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*http.Client)
	domainId := d.Get("domain_id").(int)
	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("type", "data", "priority", "weight", "host", "ttl", "port") {
		dnsRecord, err := dnsRecordFromSchema(d)
		if err != nil {
			return diag.FromErr(err)
		}

		buffer := new(bytes.Buffer)
		err = json.NewEncoder(buffer).Encode(dnsRecord)
		if err != nil {
			return diag.FromErr(err)
		}

		request, err := http.NewRequest("PUT", api.DNSRecord(domainId, recordId), buffer)
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := client.Do(request)
		if err != nil {
			return diag.FromErr(err)
		}
		defer closeBody(response.Body)

		if response.StatusCode != 204 {
			b, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return diag.FromErr(err)
			}
			return []diag.Diagnostic{{
				Severity: diag.Error,
				Summary:  "unexpected status code during update operation",

				Detail: fmt.Sprintf("expected 204, got %d. Response : %v", response.StatusCode, string(b)),
			}}
		}

		err = d.Set("last_updated", time.Now().Format(time.RFC850))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDNSRecordRead(ctx, d, m)
}

func resourceDNSRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	domainId := d.Get("domain_id").(int)
	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := m.(*http.Client)
	request, err := http.NewRequest("DELETE", api.DNSRecord(domainId, recordId), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := client.Do(request)
	if err != nil {
		return diag.FromErr(err)
	}

	if response.StatusCode != 204 {
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return diag.FromErr(err)
		}
		return []diag.Diagnostic{{
			Severity: diag.Error,
			Summary:  "unexpected status code during delete operation",

			Detail: fmt.Sprintf("expected 204, got %d. Response : %v", response.StatusCode, string(b)),
		}}
	}

	return diags
}

func dnsRecordFromSchema(d *schema.ResourceData) (*model.DnsRecord, error) {
	recordType := d.Get("type").(string)
	record := model.DnsRecord{
		Type: recordType,
		Host: d.Get("host").(string),
		Ttl:  d.Get("ttl").(int),
		Data: d.Get("data").(string),
	}

	switch recordType {
	case "SRV":
		if priority, ok := d.GetOk("priority"); ok {
			record.Priority = priority.(string)
		} else {
			return nil, fmt.Errorf("%s is required for %s record", "priority", recordType)
		}

		if weight, ok := d.GetOk("weight"); ok {
			record.Weight = weight.(int)
		} else {
			return nil, fmt.Errorf("%s is required for %s record", "weight", recordType)
		}

		if port, ok := d.GetOk("port"); ok {
			record.Port = port.(int)
		} else {
			return nil, fmt.Errorf("%s is required for %s record", "port", recordType)
		}
	case "MX":
		if priority, ok := d.GetOk("priority"); ok {
			record.Priority = priority.(string)
		} else {
			return nil, fmt.Errorf("%s is required for %s record", "priority", recordType)
		}
	}

	return &record, nil
}

func closeBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		log.Printf("closing body: %v\n", err)
	}
}
