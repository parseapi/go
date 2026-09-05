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
	c, _ = newTestClient(t, okJSON(`{"mac":"junk","valid":false,"vendor":null,"local":null,"multicast":null,"future":true}`))
	mac, err := c.MAC(context.Background(), "junk")
	if err != nil || mac.MAC != "junk" || mac.Valid || mac.Vendor != nil || mac.Local != nil || mac.Multicast != nil {
		t.Fatalf("MAC invalid decode: %#v, %v", mac, err)
	}
	c, _ = newTestClient(t, okJSON(`{"mac":"02:00:00:00:00:01","valid":true,"vendor":null,"local":true,"multicast":false}`))
	mac, err = c.MAC(context.Background(), "02:00:00:00:00:01")
	if err != nil || !mac.Valid || mac.Local == nil || !*mac.Local || mac.Multicast == nil || *mac.Multicast {
		t.Fatalf("MAC flag decode: %#v, %v", mac, err)
	}
}
