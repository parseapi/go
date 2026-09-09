package parseapi

import (
	"context"
	"testing"
)

func TestBINPreservesLeadingZerosAndNullableFields(t *testing.T) {
	client, call := newTestClient(t, okJSON(`{"bin":"00123456","prefix":"001234","country":null,"issuer":"Fixture Bank","brand":"future-brand","type":null,"prepaid":false,"deep":{},"future":true}`))
	result, err := client.BIN(context.Background(), "00 1234-56", BINOptions{Deep: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.BIN != "00123456" || result.Prefix == nil || *result.Prefix != "001234" || result.Country != nil || result.Prepaid == nil || *result.Prepaid || result.Deep == nil {
		t.Fatalf("lost BIN semantics: %#v", result)
	}
	if call.path != "/bin/00%201234-56" || call.rawQry != "deep=true" {
		t.Fatalf("bad request: %#v", call)
	}
	if _, err := client.BIN(context.Background(), "001234", BINOptions{}, BINOptions{}); err == nil {
		t.Fatal("accepted multiple options")
	}
}

func TestBINUnknown(t *testing.T) {
	client, call := newTestClient(t, okJSON(`{"bin":"000000","prefix":null,"country":null,"issuer":null,"brand":null,"type":null,"prepaid":null}`))
	result, err := client.BIN(context.Background(), "000000")
	if err != nil {
		t.Fatal(err)
	}
	if result.Prefix != nil || result.Prepaid != nil || result.Deep != nil || call.path != "/bin/000000" || call.rawQry != "" {
		t.Fatalf("unknown semantics: %#v %#v", result, call)
	}
}
