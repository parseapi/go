package parseapi

import (
	"encoding/json"
	"testing"
)

func TestTaxReferenceFields(t *testing.T) {
	var country Country
	if err := json.Unmarshal([]byte(`{"country":"DE","tax":{"name":"VAT","rate":19}}`), &country); err != nil {
		t.Fatal(err)
	}
	c := &country
	if c.Tax == nil || c.Tax.Name == nil || *c.Tax.Name != "VAT" || c.Tax.Rate == nil || *c.Tax.Rate != 19 {
		t.Fatalf("country tax: %#v", c)
	}
	var postal Postal
	if err := json.Unmarshal([]byte(`{"postal":"12345","country":"US","tax":{"name":"Sales tax","rate":7.9,"rate_state":5,"rate_county":0,"rate_city":null,"rate_special":2.9}}`), &postal); err != nil {
		t.Fatal(err)
	}
	p := &postal
	if p.Tax == nil || p.Tax.Rate == nil || *p.Tax.Rate != 7.9 || p.Tax.RateState == nil || *p.Tax.RateState != 5 || p.Tax.RateCounty == nil || *p.Tax.RateCounty != 0 || p.Tax.RateCity != nil || p.Tax.RateSpecial == nil || *p.Tax.RateSpecial != 2.9 {
		t.Fatalf("postal tax: %#v", p)
	}
}

func TestMissingAndNullTaxRemainUnknown(t *testing.T) {
	for _, body := range []string{`{}`, `{}`, `{"tax":null,"tax_rate":null}`} {
		var country Country
		var postal Postal
		if err := json.Unmarshal([]byte(body), &country); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal([]byte(body), &postal); err != nil {
			t.Fatal(err)
		}
		if country.Tax != nil && (country.Tax.Name != nil || country.Tax.Rate != nil) {
			t.Fatal("unknown country tax became a known value")
		}
		if postal.Tax != nil {
			t.Fatal("unknown postal tax became a known value")
		}
	}
}
