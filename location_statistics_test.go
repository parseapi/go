package parseapi

import (
	"encoding/json"
	"testing"
)

func TestLocationStatistics(t *testing.T) {
	for _, body := range []string{`{}`, `{"population":null,"population_period":null,"property_tax":null}`} {
		var place PostalDeep
		if err := json.Unmarshal([]byte(body), &place); err != nil {
			t.Fatal(err)
		}
		if place.PropertyTax != nil || place.PopulationPeriod != nil || place.Population != nil {
			t.Fatalf("unknown statistic: %+v", place)
		}
	}
	body := []byte(`{"population":0,"population_period":"2020-2024","property_tax":{"annual_median":0,"currency":"USD","period":"2020-2024"}}`)
	var postal PostalDeep
	var district DistrictDeep
	for _, value := range []any{&postal, &district} {
		if err := json.Unmarshal(body, value); err != nil {
			t.Fatal(err)
		}
	}
	if postal.Population == nil || *postal.Population != 0 || *postal.PopulationPeriod != "2020-2024" || postal.PropertyTax == nil || postal.PropertyTax.AnnualMedian != 0 || postal.PropertyTax.Currency != "USD" || postal.PropertyTax.Period != "2020-2024" {
		t.Fatalf("lost statistic: %+v", postal)
	}
	if district.PropertyTax == nil || district.PropertyTax.Period != "2020-2024" {
		t.Fatal("district property statistic missing")
	}
	for _, body := range []string{`{"q":"a","addresses":[]}`, `{"q":"a","addresses":[],"reason":null}`, `{"q":"a","addresses":[],"reason":"future_reason"}`} {
		var result AddressSearch
		if err := json.Unmarshal([]byte(body), &result); err != nil {
			t.Fatal(err)
		}
		if result.Reason != nil && *result.Reason != "future_reason" {
			t.Fatal("reason vocabulary was constrained")
		}
	}
}
