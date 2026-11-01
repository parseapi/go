package parseapi

import (
	"encoding/json"
	"testing"
)

func TestPostalCompactEstimateAndLegacyMetadata(t *testing.T) {
	for _, body := range []string{`{}`, `{"deep":{}}`, `{"deep":{"property_tax":null}}`} {
		var postal Postal
		if err := json.Unmarshal([]byte(body), &postal); err != nil {
			t.Fatal(err)
		}
		if postal.Deep != nil && postal.Deep.PropertyTax != nil {
			t.Fatal("unknown tax became a value")
		}
	}
	for _, amount := range []string{"0", "1234.5"} {
		var postal Postal
		body := `{"deep":{"currency":"USD","property_tax":{"annual_median":` + amount + `},"future":true}}`
		if err := json.Unmarshal([]byte(body), &postal); err != nil {
			t.Fatal(err)
		}
		if postal.Deep == nil || postal.Deep.PropertyTax == nil || postal.Deep.Currency == nil || *postal.Deep.Currency != "USD" {
			t.Fatal("compact record lost")
		}
		tax := postal.Deep.PropertyTax
		if tax.Currency != nil || tax.Period != nil {
			t.Fatal("compact record invented metadata")
		}
		wire, err := json.Marshal(tax)
		if err != nil {
			t.Fatal(err)
		}
		if string(wire) != `{"annual_median":`+amount+`}` {
			t.Fatalf("unexpected compact JSON: %s", wire)
		}
	}
	var postal PostalDeep
	var district DistrictDeep
	body := []byte(`{"property_tax":{"annual_median":0,"currency":"USD","period":"2020-2024"}}`)
	if err := json.Unmarshal(body, &postal); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(body, &district); err != nil {
		t.Fatal(err)
	}
	if postal.PropertyTax.Currency == nil || *postal.PropertyTax.Currency != "USD" || postal.PropertyTax.Period == nil || *postal.PropertyTax.Period != "2020-2024" {
		t.Fatal("legacy Postal metadata lost")
	}
	if district.PropertyTax.Currency != "USD" || district.PropertyTax.Period != "2020-2024" {
		t.Fatal("District changed")
	}
}

func TestWeatherMinuteCanonicalPresenceAndLegacyPayloads(t *testing.T) {
	for _, fixture := range []struct {
		body         string
		rate, inches *float64
	}{
		{`{}`, nil, nil},
		{`{"precipitation_rate":1.8,"precipitation_rate_in":0.071}`, floatPointer(1.8), floatPointer(0.071)},
		{`{"precipitation":1.8,"precipitation_in":0.071}`, floatPointer(1.8), floatPointer(0.071)},
		{`{"precipitation_rate":0,"precipitation_rate_in":0,"precipitation":9,"precipitation_in":9}`, floatPointer(0), floatPointer(0)},
		{`{"precipitation_rate":null,"precipitation_rate_in":null,"precipitation":9,"precipitation_in":9}`, nil, nil},
		{`{"precipitation_rate":null,"precipitation_in":0.071}`, nil, floatPointer(0.071)},
	} {
		var minute WeatherMinute
		if err := json.Unmarshal([]byte(fixture.body), &minute); err != nil {
			t.Fatal(err)
		}
		for _, pair := range [][2]*float64{{minute.PrecipitationRate, fixture.rate}, {minute.PrecipitationRateIn, fixture.inches}} {
			if (pair[0] == nil) != (pair[1] == nil) || pair[0] != nil && *pair[0] != *pair[1] {
				t.Fatalf("wrong rate for %s", fixture.body)
			}
		}
		wire, err := json.Marshal(minute)
		if err != nil {
			t.Fatal(err)
		}
		var keys map[string]json.RawMessage
		if err := json.Unmarshal(wire, &keys); err != nil {
			t.Fatal(err)
		}
		if _, present := keys["precipitation"]; present {
			t.Fatal("legacy key serialized")
		}
		if _, present := keys["precipitation_in"]; present {
			t.Fatal("legacy key serialized")
		}
	}
	var reused WeatherMinute
	if err := json.Unmarshal([]byte(`{"precipitation_rate":1}`), &reused); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(`{}`), &reused); err != nil {
		t.Fatal(err)
	}
	if reused.PrecipitationRate != nil {
		t.Fatal("missing key retained stale value")
	}
	var weather Weather
	if err := json.Unmarshal([]byte(`{"deep":{"history":{"precipitation":12.3,"precipitation_in":0.484}}}`), &weather); err != nil {
		t.Fatal(err)
	}
	if weather.Deep.History.Precipitation == nil || *weather.Deep.History.Precipitation != 12.3 {
		t.Fatal("history accumulation changed")
	}
}

func floatPointer(value float64) *float64 { return &value }
