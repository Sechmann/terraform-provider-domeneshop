package main

import (
	"context"
	"github.com/VegarM/domeneshop-go"
	"log"
	"os"
	"time"
)

var baseContext context.Context

func init() {
	token, found := os.LookupEnv("DOMENESHOP_TOKEN")
	if !found {
		log.Fatal("Reading env DOMENESHOP_TOKEN: not set")
	}

	secret, found := os.LookupEnv("DOMENESHOP_SECRET")
	if !found {
		log.Fatal("Reading env DOMENESHOP_SECRET: not set")
	}

	baseContext = context.WithValue(context.Background(), domeneshop.ContextBasicAuth, domeneshop.BasicAuth{
		UserName: token,
		Password: secret,
	})
}

func baseContextWithTimeout(duration time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(baseContext, duration)
}

func main() {
	cfg := domeneshop.NewConfiguration()
	client := domeneshop.NewAPIClient(cfg)
	ctx, cancel := baseContextWithTimeout(time.Second*15)
	defer cancel()
	domains, i, err := client.DomainsApi.GetDomains(ctx, nil)
	if err != nil {
		log.Fatalf("listing domains: %v", err)
	}

	log.Printf("domains: %+v\ni: %+v", domains, i)
}
