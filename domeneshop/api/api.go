package api

import "fmt"

const (
	apiURL = "https://api.domeneshop.no/v0"
)

func DNSRecords(domainId int) string {
	return fmt.Sprintf("%s/domains/%d/dns", Domains(), domainId)
}

func DNSRecord(domainId, recordId int) string {
	return fmt.Sprintf("%s/%d", DNSRecords(domainId), recordId)
}

func Domains() string {
	return fmt.Sprintf("%s/domains", apiURL)
}

func Domain(domainId  int) string {
	return fmt.Sprintf("%s/%d", Domains(), domainId)
}
