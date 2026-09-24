package parseapi

import (
	"context"
	"strings"
	"testing"
)

func TestCardPreservesLeadingZerosAndNullableFields(t *testing.T) {
	client, call := newTestClient(t, okJSON(`{"bin":"00123456","prefix":"001234","country":null,"issuer":"Fixture Bank","brand":"future-brand","type":null,"prepaid":false,"deep":{},"future":true}`))
	result, err := client.Card(context.Background(), "00 1234-56")
	if err != nil {
		t.Fatal(err)
	}
	if result.BIN != "00123456" || result.Prefix == nil || *result.Prefix != "001234" || result.Country != nil || result.Prepaid == nil || *result.Prepaid {
		t.Fatalf("lost BIN semantics: %#v", result)
	}
	if call.path != "/card/00%201234-56" || call.rawQry != "" {
		t.Fatalf("bad request: %#v", call)
	}
}

func TestCardUnknown(t *testing.T) {
	client, call := newTestClient(t, okJSON(`{"bin":"000000","prefix":null,"country":null,"issuer":null,"brand":null,"type":null,"prepaid":null}`))
	result, err := client.Card(context.Background(), "000000")
	if err != nil {
		t.Fatal(err)
	}
	if result.Prefix != nil || result.Prepaid != nil || call.path != "/card/000000" || call.rawQry != "" {
		t.Fatalf("unknown semantics: %#v %#v", result, call)
	}
}

func TestCardValidatesBeforeDispatchWithoutRepairingInput(t *testing.T) {
	for _, input := range []string{"", "12345", "123456789012", "4242424242424242", "００１２３４", "00\u00a01234", "00\v1234", "00%201234", "00+1234", "00/1234", "00_1234", strings.Repeat(" ", 59) + "001234"} {
		client, call := newTestClient(t, okJSON(`{}`))
		_, err := client.Card(context.Background(), input)
		if err == nil || call.path != "" || err.Error() != "parseapi: Card requires a string containing 6 to 11 digits. Send a prefix only." {
			t.Fatalf("guard: %v, %#v", err, call)
		}
	}
	for _, input := range []string{"001234", "00123456789", "00 1234-56", "00\t12\r34\n-56", strings.Repeat(" ", 58) + "001234"} {
		client, call := newTestClient(t, okJSON(`{"bin":"001234"}`))
		if _, err := client.Card(context.Background(), input); err != nil {
			t.Fatal(err)
		}
		if call.path != "/card/"+seg(input) || call.rawQry != "" {
			t.Fatalf("input was changed: %#v", call)
		}
	}
}
