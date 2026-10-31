package parseapi

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestISOWireNamesAcceptReleasedAndCanonicalResponses(t *testing.T) {
	for _, body := range []string{`{"iso3":"USA","numeric":840}`, `{"iso_alpha3":"USA","iso_numeric":840}`} {
		var country Country
		if err := json.Unmarshal([]byte(body), &country); err != nil {
			t.Fatal(err)
		}
		if country.ISO3 == nil || *country.ISO3 != "USA" || country.Numeric == nil || *country.Numeric != 840 {
			t.Fatalf("Country ISO fields lost: %s", body)
		}
		encoded, err := json.Marshal(country)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(encoded), `"iso_alpha3":"USA"`) || !strings.Contains(string(encoded), `"iso_numeric":840`) || strings.Contains(string(encoded), `"iso3":`) || strings.Contains(string(encoded), `"numeric":`) {
			t.Fatalf("Country wire aliases leaked: %s", encoded)
		}
		var deep CountryDeep
		if err = json.Unmarshal([]byte(body), &deep); err != nil || deep.ISO3 == nil || *deep.ISO3 != "USA" || deep.Numeric == nil || *deep.Numeric != 840 {
			t.Fatalf("Historical Country Deep lost: %s %v", body, err)
		}
	}
	for _, body := range []string{`{"iso3":"hrv"}`, `{"iso_alpha3":"hrv"}`} {
		var value Language
		if err := json.Unmarshal([]byte(body), &value); err != nil || value.ISO3 == nil || *value.ISO3 != "hrv" {
			t.Fatalf("Language: %s %v", body, err)
		}
		encoded, _ := json.Marshal(value)
		if !strings.Contains(string(encoded), `"iso_alpha3":"hrv"`) || strings.Contains(string(encoded), `"iso3":`) {
			t.Fatalf("Language wire: %s", encoded)
		}
	}
	for _, body := range []string{`{"iso":"US-CA"}`, `{"iso_3166_2":"US-CA"}`, `{"iso_subdivision":"US-CA"}`} {
		var value State
		if err := json.Unmarshal([]byte(body), &value); err != nil || value.ISO == nil || *value.ISO != "US-CA" {
			t.Fatalf("State: %s %v", body, err)
		}
		encoded, _ := json.Marshal(value)
		if !strings.Contains(string(encoded), `"iso_subdivision":"US-CA"`) || strings.Contains(string(encoded), `"iso":`) || strings.Contains(string(encoded), `"iso_3166_2":`) {
			t.Fatalf("State wire: %s", encoded)
		}
	}
}

func TestISOCandidatesRespectExplicitCanonicalNull(t *testing.T) {
	for _, body := range []string{`{"iso3":"USA","numeric":840,"iso_alpha3":null,"iso_numeric":null}`, `{"iso_alpha3":null,"iso_numeric":null,"iso3":"USA","numeric":840}`} {
		var value Country
		var deep CountryDeep
		if err := json.Unmarshal([]byte(body), &value); err != nil || value.ISO3 != nil || value.Numeric != nil {
			t.Fatalf("Country explicit null: %s %v", body, err)
		}
		if err := json.Unmarshal([]byte(body), &deep); err != nil || deep.ISO3 != nil || deep.Numeric != nil {
			t.Fatalf("Country Deep explicit null: %s %v", body, err)
		}
	}
	for _, body := range []string{`{"iso_alpha3":null,"iso3":"hrv"}`, `{"iso3":"hrv","iso_alpha3":null}`} {
		var value Language
		if err := json.Unmarshal([]byte(body), &value); err != nil || value.ISO3 != nil {
			t.Fatalf("Language explicit null: %s %v", body, err)
		}
	}
	for _, body := range []string{`{"iso_subdivision":null,"iso":"US-CA","iso_3166_2":"US-TX"}`, `{"iso":"US-CA","iso_3166_2":"US-TX","iso_subdivision":null}`} {
		var value State
		if err := json.Unmarshal([]byte(body), &value); err != nil || value.ISO != nil {
			t.Fatalf("State explicit null: %s %v", body, err)
		}
	}
}

func TestHistoricalDeepReadMembersRemainAvailable(t *testing.T) {
	var language LanguageDeep
	if err := json.Unmarshal([]byte(`{"iso3":"eng","countries":["US"]}`), &language); err != nil || language.ISO3 == nil || *language.ISO3 != "eng" || len(language.Countries) != 1 {
		t.Fatalf("Language Deep: %#v %v", language, err)
	}
	var state StateDeep
	if err := json.Unmarshal([]byte(`{"fips":"06","capital":"Sacramento","area_codes":["916"],"tax":"Sales tax","tax_rate":0}`), &state); err != nil || state.FIPS == nil || *state.FIPS != "06" || state.Capital == nil || *state.Capital != "Sacramento" || len(state.AreaCodes) != 1 || state.TaxRate == nil || *state.TaxRate != 0 || state.Tax == nil || state.Tax.Name == nil || *state.Tax.Name != "Sales tax" {
		t.Fatalf("State Deep: %#v %v", state, err)
	}
	body, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{`"fips":`, `"capital":`, `"area_codes":`, `"tax_rate":`} {
		if strings.Contains(string(body), key) {
			t.Fatalf("Historical Deep alias leaked: %s", body)
		}
	}
}

func TestHistoricalCountryCoordinatesRemainReadableWithoutCapitalName(t *testing.T) {
	var country Country
	if err := json.Unmarshal([]byte(`{"country":"US","iso3":"USA","capital_lat":0,"capital_lon":null}`), &country); err != nil {
		t.Fatal(err)
	}
	if country.CapitalLat == nil || *country.CapitalLat != 0 || country.CapitalLon != nil || country.Capital == nil || country.Capital.Name != nil {
		t.Fatalf("Historical standalone coordinates changed: %#v", country)
	}
}
