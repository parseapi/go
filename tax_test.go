package parseapi

import (
	"encoding/json"
	"testing"
)

func TestTaxReferenceFields(t *testing.T) {
	var country Country
	if err := json.Unmarshal([]byte(`{"country":"DE","deep":{"tax":"VAT","tax_rate":19,"tax_id_format":"DE999999999","tax_id_regex":"^DE[0-9]{9}$"}}`), &country); err != nil {
		t.Fatal(err)
	}
	c := country.Deep
	if c == nil || c.Tax == nil || *c.Tax != "VAT" || c.TaxRate == nil || *c.TaxRate != 19 || c.TaxIDFormat == nil || *c.TaxIDFormat != "DE999999999" || c.TaxIDRegex == nil || *c.TaxIDRegex != "^DE[0-9]{9}$" {
		t.Fatalf("country tax: %#v", c)
	}
	var postal Postal
	if err := json.Unmarshal([]byte(`{"postal":"12345","country":"US","deep":{"tax":"Sales tax","tax_rate":7.9,"tax_rate_state":5,"tax_rate_county":0,"tax_rate_city":null,"tax_rate_special":2.9}}`), &postal); err != nil {
		t.Fatal(err)
	}
	p := postal.Deep
	if p == nil || p.TaxRate == nil || *p.TaxRate != 7.9 || p.TaxRateState == nil || *p.TaxRateState != 5 || p.TaxRateCounty == nil || *p.TaxRateCounty != 0 || p.TaxRateCity != nil || p.TaxRateSpecial == nil || *p.TaxRateSpecial != 2.9 {
		t.Fatalf("postal tax: %#v", p)
	}
}

func TestMissingAndNullTaxRemainUnknown(t *testing.T) {
	for _, body := range []string{`{}`, `{"deep":{}}`, `{"deep":{"tax":null,"tax_rate":null}}`} {
		var country Country
		var postal Postal
		if err := json.Unmarshal([]byte(body), &country); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal([]byte(body), &postal); err != nil {
			t.Fatal(err)
		}
		if country.Deep != nil && (country.Deep.Tax != nil || country.Deep.TaxRate != nil || country.Deep.TaxIDFormat != nil || country.Deep.TaxIDRegex != nil) {
			t.Fatal("unknown country tax became a known value")
		}
		if postal.Deep != nil && (postal.Deep.Tax != nil || postal.Deep.TaxRate != nil || postal.Deep.TaxRateSpecial != nil) {
			t.Fatal("unknown postal tax became a known value")
		}
	}
}
