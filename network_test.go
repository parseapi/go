package parseapi

import (
	"context"
	"testing"
)

func TestNetworkRecords(t *testing.T) {
	c, _ := newTestClient(t, okJSON(`{"asn":4294967295,"name":null,"country":null,"country_name":null,"future":true}`))
	asn, err := c.ASN(context.Background(), "4294967295")
	if err != nil || asn.ASN != 4294967295 || asn.Name != nil || asn.Country != nil || asn.CountryName != nil {
		t.Fatalf("ASN nullable decode: %#v, %v", asn, err)
	}
}
