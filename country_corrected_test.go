package parseapi

import (
	"encoding/json"
	"os"
	"testing"
)

func TestCorrectedCountryResponse(t *testing.T) {
	b, e := os.ReadFile("testdata/country-corrected.json")
	if e != nil {
		t.Fatal(e)
	}
	var c Country
	if e = json.Unmarshal(b, &c); e != nil {
		t.Fatal(e)
	}
	if c.NameFull == nil || *c.NameFull != "United States of America" || c.FullName != c.NameFull || c.Demonym == nil || *c.Demonym != "American" || c.Numeric == nil || *c.Numeric != 840 {
		t.Fatal("Core name family missing")
	}
	if c.Flag == nil || *c.Flag != "🇺🇸" || c.Emoji != c.Flag || c.Capital == nil || c.Capital.Name == nil || *c.Capital.Name != "Washington, D.C." || c.Capital.Latitude == nil || *c.Capital.Latitude != 38.89 || c.Capital.Longitude == nil || *c.Capital.Longitude != -77.05 {
		t.Fatal("Country flag/capital group missing")
	}
	if c.Deep == nil || c.Deep.Tax == nil || c.Deep.Tax.Name != nil || c.Deep.Tax.Rate != nil || c.Deep.Demonym != nil || c.Deep.FullName != nil || c.Deep.GDPPerCapitaUSD == nil || *c.Deep.GDPPerCapitaUSD != 90026.52 || c.Deep.ElevationAverage == nil || *c.Deep.ElevationAverage != 760 || c.Deep.Population == nil || c.Deep.Area == nil {
		t.Fatal("Country Deep contract changed")
	}
	encoded, e := json.Marshal(c)
	if e != nil {
		t.Fatal(e)
	}
	var raw map[string]json.RawMessage
	if e = json.Unmarshal(encoded, &raw); e != nil {
		t.Fatal(e)
	}
	for _, key := range []string{"emoji", "capital_lat", "capital_lon", "full_name"} {
		if _, ok := raw[key]; ok {
			t.Fatal("Legacy key emitted", key)
		}
	}
	if _, ok := raw["full_name"]; ok {
		t.Fatal("Legacy key emitted")
	}
	if _, ok := raw["name_full"]; !ok {
		t.Fatal("Canonical key missing")
	}
	var deep map[string]json.RawMessage
	json.Unmarshal(raw["deep"], &deep)
	for _, key := range []string{"full_name", "name_full", "demonym", "iso3", "numeric", "capital", "borders", "gdp", "elevation", "tax_rate", "tax_id_format", "tax_id_regex"} {
		if _, ok := deep[key]; ok {
			t.Fatal("Deep alias emitted", key)
		}
	}
}
func TestCountryNamesLegacyAndNull(t *testing.T) {
	for _, input := range []string{`{"country":"US","name_full":"Full"}`, `{"country":"US","full_name":"Full"}`} {
		var c Country
		if e := json.Unmarshal([]byte(input), &c); e != nil {
			t.Fatal(e)
		}
		if c.NameFull == nil || *c.NameFull != "Full" || c.FullName != c.NameFull {
			t.Fatal("Legacy name adapter")
		}
	}
	var c Country
	if e := json.Unmarshal([]byte(`{"name_full":null,"full_name":"Legacy"}`), &c); e != nil {
		t.Fatal(e)
	}
	if c.NameFull != nil || c.FullName != nil {
		t.Fatal("Canonical null lost")
	}
	if e := json.Unmarshal([]byte(`{"deep":{"full_name":"Earlier","demonym":"Earlier"}}`), &c); e != nil {
		t.Fatal(e)
	}
	if c.Deep.FullName == nil || *c.Deep.FullName != "Earlier" || c.Deep.Demonym == nil {
		t.Fatal("Historical Deep decoder lost")
	}
}
func TestCountryActualAntarctica(t *testing.T) {
	b, e := os.ReadFile("testdata/country-aq.json")
	if e != nil {
		t.Fatal(e)
	}
	var c Country
	if e = json.Unmarshal(b, &c); e != nil {
		t.Fatal(e)
	}
	if c.Currency != nil || c.CurrencyName != nil || c.CurrencySymbol != nil || c.Languages != nil || c.Deep == nil || c.Deep.Population == nil || *c.Deep.Population != 0 || c.NameFull == nil || *c.NameFull != "Antarctica" {
		t.Fatal("Actual nullable Country fields changed")
	}
}

func TestCountryTaxGroupsAndZeroes(t *testing.T) {
	b, err := os.ReadFile("testdata/country-de.json")
	if err != nil {
		t.Fatal(err)
	}
	var c Country
	if err = json.Unmarshal(b, &c); err != nil {
		t.Fatal(err)
	}
	tax := c.Deep.Tax
	if tax == nil || tax.Name == nil || *tax.Name != "VAT" || tax.Rate == nil || *tax.Rate != 19 {
		t.Fatal("Actual German tax reference lost")
	}
	body := `{"flag":null,"emoji":"OLD","capital":{"name":null,"latitude":0,"longitude":null},"capital_lat":99,"capital_lon":99,"deep":{"tax":{"name":null,"rate":0},"tax_rate":99,"tax_id_format":"OLD"}}`
	if err = json.Unmarshal([]byte(body), &c); err != nil {
		t.Fatal(err)
	}
	if c.Flag != nil || c.Emoji != nil || c.Capital.Name != nil || c.Capital.Latitude == nil || *c.Capital.Latitude != 0 || c.Capital.Longitude != nil || c.Deep.Tax.Name != nil || c.Deep.Tax.Rate == nil || *c.Deep.Tax.Rate != 0 || c.Deep.TaxIDFormat != nil {
		t.Fatal("Canonical nested null/zero was overwritten by a legacy value")
	}
	encoded, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]json.RawMessage
	json.Unmarshal(encoded, &raw)
	for _, key := range []string{"emoji", "capital_lat", "capital_lon"} {
		if _, ok := raw[key]; ok {
			t.Fatal("Legacy root emitted", key)
		}
	}
	var deep map[string]json.RawMessage
	json.Unmarshal(raw["deep"], &deep)
	var taxChildren map[string]json.RawMessage
	json.Unmarshal(deep["tax"], &taxChildren)
	if len(taxChildren) != 2 || taxChildren["name"] == nil || taxChildren["rate"] == nil {
		t.Fatal("Tax summary must contain only name and rate")
	}
	for _, key := range []string{"tax_rate", "tax_id_format", "tax_id_regex"} {
		if _, ok := deep[key]; ok {
			t.Fatal("Legacy tax emitted", key)
		}
	}
}
func TestCountryHistoricalShapeAdapter(t *testing.T) {
	for _, body := range []string{
		`{"emoji":"🇺🇸","capital":"Washington, D.C.","capital_lat":0,"capital_lon":-77.05,"deep":{"tax":"VAT","tax_rate":0,"tax_id_format":"FORMAT","tax_id_regex":"REGEX"}}`,
		`{"emoji":"🇺🇸","capital":null,"capital_lat":0,"capital_lon":-77.05,"deep":{"tax":null,"tax_rate":0,"tax_id_format":"FORMAT","tax_id_regex":"REGEX"}}`,
	} {
		var c Country
		if err := json.Unmarshal([]byte(body), &c); err != nil {
			t.Fatal(err)
		}
		if c.Flag == nil || *c.Flag != "🇺🇸" || c.Capital == nil || c.Capital.Latitude == nil || *c.Capital.Latitude != 0 || c.Deep == nil || c.Deep.Tax == nil || c.Deep.Tax.Rate == nil || *c.Deep.Tax.Rate != 0 || c.Deep.TaxIDFormat == nil || *c.Deep.TaxIDFormat != "FORMAT" {
			t.Fatal("Historical shape adapter lost zero or facts")
		}
	}
}
func TestCountryDeepModesAndCompactFlags(t *testing.T) {
	for _, body := range []string{`{"country":"US","flag":"🇺🇸","capital":{"name":null,"latitude":null,"longitude":null}}`, `{"country":"US","flag":"🇺🇸","capital":{"name":null,"latitude":null,"longitude":null},"deep":{}}`} {
		var c Country
		if err := json.Unmarshal([]byte(body), &c); err != nil {
			t.Fatal(err)
		}
		var raw map[string]json.RawMessage
		json.Unmarshal([]byte(body), &raw)
		_, requested := raw["deep"]
		if requested != (c.Deep != nil) {
			t.Fatal("Deep triad lost")
		}
		if c.Deep != nil && c.Deep.Tax != nil {
			t.Fatal("Locked Deep fabricated tax")
		}
	}
	for _, body := range []string{`{"country":"US","flag":"🇺🇸"}`, `{"country":"US","emoji":"🇺🇸"}`} {
		var c ContinentCountryItem
		if err := json.Unmarshal([]byte(body), &c); err != nil {
			t.Fatal(err)
		}
		if c.Flag == nil || *c.Flag != "🇺🇸" || c.Emoji != c.Flag {
			t.Fatal("Compact Country flag lost")
		}
		encoded, _ := json.Marshal(c)
		var raw map[string]json.RawMessage
		json.Unmarshal(encoded, &raw)
		if _, ok := raw["emoji"]; ok {
			t.Fatal("Compact legacy key emitted")
		}
	}
}
