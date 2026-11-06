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
	if c.NameFull == nil || *c.NameFull != "United States of America" || c.Demonym == nil || *c.Demonym != "American" || c.Numeric == nil || *c.Numeric != 840 {
		t.Fatal("Core name family missing")
	}
	if c.Flag == nil || *c.Flag != "🇺🇸" || c.Capital == nil || c.Capital.Name == nil || *c.Capital.Name != "Washington, D.C." || c.Capital.Latitude == nil || *c.Capital.Latitude != 38.89 || c.Capital.Longitude == nil || *c.Capital.Longitude != -77.05 {
		t.Fatal("Country flag/capital group missing")
	}
	if c.Tax == nil || c.Tax.Name != nil || c.Tax.Rate != nil || c.GDPPerCapitaUSD == nil || *c.GDPPerCapitaUSD != 90026.52 || c.ElevationAverage == nil || *c.ElevationAverage != 760 || c.Population == nil || c.Area == nil {
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
	if _, ok := raw["deep"]; ok { t.Fatal("2.1 must not emit deep") }
}
func TestCountryCanonicalNamesAndNull(t *testing.T) {
 var c Country
 if err:=json.Unmarshal([]byte(`{"country":"US","name_full":"Full","demonym":"American"}`),&c);err!=nil{t.Fatal(err)}
 if c.NameFull==nil || *c.NameFull!="Full" || c.Demonym==nil{t.Fatal("Canonical name missing")}
 if err:=json.Unmarshal([]byte(`{"name_full":null,"demonym":null,"population":0}`),&c);err!=nil{t.Fatal(err)}
 if c.NameFull!=nil || c.Demonym!=nil || c.Population==nil || *c.Population!=0{t.Fatal("Null or zero changed")}
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
	if c.Currency != nil || c.CurrencyName != nil || c.CurrencySymbol != nil || len(c.Languages) != 0 || c.Population == nil || *c.Population != 0 || c.NameFull == nil || *c.NameFull != "Antarctica" {
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
	tax := c.Tax
	if tax == nil || tax.Name == nil || *tax.Name != "VAT" || tax.Rate == nil || *tax.Rate != 19 {
		t.Fatal("Actual German tax reference lost")
	}
	body := `{"flag":null,"emoji":"OLD","capital":{"name":null,"latitude":0,"longitude":null},"capital_lat":99,"capital_lon":99,"tax":{"name":null,"rate":0},"tax_rate":99,"tax_id_format":"OLD"}`
	if err = json.Unmarshal([]byte(body), &c); err != nil {
		t.Fatal(err)
	}
	if c.Flag != nil || c.Emoji != nil || c.Capital.Name != nil || c.Capital.Latitude == nil || *c.Capital.Latitude != 0 || c.Capital.Longitude != nil || c.Tax.Name != nil || c.Tax.Rate == nil || *c.Tax.Rate != 0 {
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
	deep = raw
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
func TestCountryDeepModesAndCompactFlags(t *testing.T) {
	for _, body := range []string{`{"country":"US","flag":"🇺🇸","capital":{"name":null,"latitude":null,"longitude":null}}`, `{"country":"US","flag":"🇺🇸","capital":{"name":null,"latitude":null,"longitude":null}}`} {
		var c Country
		if err := json.Unmarshal([]byte(body), &c); err != nil {
			t.Fatal(err)
		}
		var raw map[string]json.RawMessage
		json.Unmarshal([]byte(body), &raw)
		if _, exists := raw["deep"]; exists {
			t.Fatal("2.1 fixture retained an obsolete wrapper")
		}

		if c.Tax != nil {
			t.Fatal("Locked Deep fabricated tax")
		}
	}
	for _, body := range []string{`{"country":"US","flag":"🇺🇸"}`, `{"country":"US","emoji":"🇺🇸"}`} {
		var c ContinentCountryItem
		if err := json.Unmarshal([]byte(body), &c); err != nil {
			t.Fatal(err)
		}
		if c.Flag == nil || *c.Flag != "🇺🇸" {
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
