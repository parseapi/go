package parseapi

import (
	"encoding/json"
	"os"
	"testing"
)

func TestCoreConnectorTypesPreserveDeepAndOlderResponses(t *testing.T) {
	data, err := os.ReadFile("testdata/core-connectors.json")
	if err != nil {
		t.Fatal(err)
	}
	var f map[string]json.RawMessage
	if err = json.Unmarshal(data, &f); err != nil {
		t.Fatal(err)
	}
	var country Country
	var language Language
	var currency Currency
	var postal Postal
	var industry Industry
	var emoji Emoji
	var company CompanyProfile
	var point Point
	for key, target := range map[string]interface{}{"country": &country, "language": &language, "currency": &currency, "postal": &postal, "industry": &industry, "emoji": &emoji, "company": &company, "point": &point} {
		if err = json.Unmarshal(f[key], target); err != nil {
			t.Fatalf("%s: %v", key, err)
		}
	}
	if country.ISO3 == nil || *country.ISO3 != "USA" || country.Capital == nil || country.Capital.Latitude == nil || *country.Capital.Latitude != 0 || len(country.Borders) != 2 {
		t.Fatal("country refs lost")
	}
	if language.ISO3 == nil || *language.ISO3 != "eng" || len(language.Countries) != 2 || len(currency.Countries) != 1 || postal.Currency == nil || *postal.Currency != "USD" {
		t.Fatal("reference codes lost")
	}
	if industry.Children[0].NAICS != "111" {
		t.Fatal("industry Core/Deep types mixed")
	}
	var search NAICSSearch
	if err = json.Unmarshal(f["industry_search"], &search); err != nil || search.Results[0].NAICS != "722515" || search.Results[0].Match.Text != "Coffee shops, on premise brewing" {
		t.Fatal("Industry search identity or evidence changed")
	}
	if emoji.Base == nil || *emoji.Base != "👍" || emoji.Variants[0].Emoji != "👍🏽" {
		t.Fatal("emoji Core/Deep types mixed")
	}
	if company.Identifiers[0].Value != "0000320193" || company.Industries[0].Code != "3571" || company.PhoneNumbers[0].E164 == nil || *company.PhoneNumbers[0].E164 != "+14089961010" || company.PhoneNumbers[0].Country == nil || *company.PhoneNumbers[0].Country != "US" {
		t.Fatal("company refs lost")
	}
	if point.City.ID == nil || *point.City.ID != "city_7aj899nvwyku" {
		t.Fatal("Point Deep city identity changed")
	}
	if postal.Country != "US" || len(postal.Neighbors) != 1 || postal.Neighbors[0] != "00544" {
		t.Fatal("Postal neighbor country, zeros or Deep changed")
	}
	for _, body := range []string{`{}`, `{"neighbors":null}`} {
		var older PostalDeep
		if err = json.Unmarshal([]byte(body), &older); err != nil || older.Neighbors != nil {
			t.Fatal("missing/null neighbor metadata changed")
		}
	}
	var emptyPostal PostalDeep
	if err = json.Unmarshal([]byte(`{"neighbors":[]}`), &emptyPostal); err != nil || emptyPostal.Neighbors == nil {
		t.Fatal("known empty neighbors changed")
	}
	var oldPoint Point
	if err = json.Unmarshal([]byte(`{"city_id":"unpublished-alias"}`), &oldPoint); err != nil || oldPoint.City != nil {
		t.Fatal("older Point changed")
	}
	for _, target := range []interface{}{&Country{}, &Language{}, &Currency{}, &Postal{}, &Industry{}, &Emoji{}, &CompanyProfile{}} {
		if err = json.Unmarshal([]byte(`{}`), target); err != nil {
			t.Fatal(err)
		}
	}
	var unknown Language
	if err = json.Unmarshal([]byte(`{"countries":null}`), &unknown); err != nil || unknown.Countries != nil {
		t.Fatal("unknown relationships changed")
	}
	var empty Language
	if err = json.Unmarshal([]byte(`{"countries":[]}`), &empty); err != nil || empty.Countries == nil {
		t.Fatal("known empty relationships changed")
	}
}
