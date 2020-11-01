package domeneshop_test

import (
	"terraform-provider-domeneshop/domeneshop"
	"testing"
)

func TestProvider(t *testing.T) {
	if err := domeneshop.Provider().InternalValidate(); err != nil {
		t.Fatalf("internal validate err: %v", err)
	}
}
