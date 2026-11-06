package parseapi

import (
	"encoding/json"
	"testing"
)

func TestPostalCompactPropertyTaxAndDistrictFullRecord(t *testing.T) {
	for _, body := range []string{`{"postal":"10001","country":"US"}`, `{"postal":"10001","country":"US"}`, `{"postal":"10001","country":"US","property_tax":null}`} {
		var postal Postal
		if err := json.Unmarshal([]byte(body), &postal); err != nil {
			t.Fatal(err)
		}
		if postal.PropertyTax != nil {
			t.Fatal("unknown tax became a value")
		}
	}
	for _, amount := range []string{"0", "1234.5"} {
		var postal Postal
		body := `{"postal":"10001","country":"US","currency":"USD","property_tax":{"annual_median":` + amount + `},"future":true}`
		if err := json.Unmarshal([]byte(body), &postal); err != nil {
			t.Fatal(err)
		}
		if postal.PropertyTax == nil || postal.Currency == nil || *postal.Currency != "USD" {
			t.Fatal("compact Postal tax was lost")
		}
		want := 0.0
		if amount != "0" {
			want = 1234.5
		}
		if postal.PropertyTax.AnnualMedian != want {
			t.Fatal("amount changed")
		}
		if postal.PropertyTax.Currency != nil || postal.PropertyTax.Period != nil {
			t.Fatal("compact record invented legacy metadata")
		}
		wire, err := json.Marshal(postal.PropertyTax)
		if err != nil {
			t.Fatal(err)
		}
		var roundtrip map[string]any
		if err := json.Unmarshal(wire, &roundtrip); err != nil {
			t.Fatal(err)
		}
		if len(roundtrip) != 1 || roundtrip["annual_median"] != want {
			t.Fatalf("compact tax gained fields: %s", wire)
		}
	}
	var legacy PostalDeep
	if err := json.Unmarshal([]byte(`{"property_tax":{"annual_median":1234.5,"currency":"USD","period":"2020-2024"}}`), &legacy); err != nil {
		t.Fatal(err)
	}
	if legacy.PropertyTax == nil || legacy.PropertyTax.Currency == nil || *legacy.PropertyTax.Currency != "USD" || legacy.PropertyTax.Period == nil || *legacy.PropertyTax.Period != "2020-2024" {
		t.Fatal("older richer Postal metadata was lost")
	}
	var district DistrictDeep
	if err := json.Unmarshal([]byte(`{"property_tax":{"annual_median":0,"currency":"USD","period":"2020-2024"}}`), &district); err != nil {
		t.Fatal(err)
	}
	if district.PropertyTax == nil || district.PropertyTax.Currency != "USD" || district.PropertyTax.Period != "2020-2024" {
		t.Fatal("District's full record changed")
	}
}

func TestPointNaturalRetainsMetricAndElectoralFacts(t *testing.T) {
	for _, body := range []string{`{"latitude":0,"longitude":0}`, `{"latitude":0,"longitude":0}`, `{"latitude":0,"longitude":0,"elevation":null,"electoral":null}`} {
		var point Point
		if err := json.Unmarshal([]byte(body), &point); err != nil {
			t.Fatal(err)
		}
		if point.Elevation != nil || point.Electoral != nil {
			t.Fatal("unknown detail became a fact")
		}
	}
	var zero Point
	if err := json.Unmarshal([]byte(`{"latitude":0,"longitude":0,"elevation":0,"elevation_ft":0,"electoral":{"code":"0801","name":null,"congress":119,"future":true}}`), &zero); err != nil {
		t.Fatal(err)
	}
	if zero.Elevation == nil || *zero.Elevation != 0 || zero.ElevationFt == nil || *zero.ElevationFt != 0 {
		t.Fatal("zero terrain was lost")
	}
	if zero.Electoral == nil || zero.Electoral.Code != "0801" || zero.Electoral.Name != nil || zero.Electoral.Congress == nil || *zero.Electoral.Congress != 119 {
		t.Fatal("electoral data changed")
	}
	var old Point
	if err := json.Unmarshal([]byte(`{"elevation":-1.5,"resolution":30,"electoral":{"code":"0200","name":"At-large","congress":null}}`), &old); err != nil {
		t.Fatal(err)
	}
	if old.Elevation == nil || *old.Elevation != -1.5 || old.Resolution == nil || *old.Resolution != 30 || old.Electoral.Congress != nil {
		t.Fatal("nullable or historical detail changed")
	}
}
