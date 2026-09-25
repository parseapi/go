package parseapi

import (
	"context"
	"net/url"
	"reflect"
	"testing"
)

// A translated call must preserve existing input/depth controls, and the next
// default call must not inherit its language. Frozen method types live in contract_test.go.
func TestDisplayLanguagePerRequest(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name, path, query string
		invoke            func(*Client, string) error
	}{
		{"IP", "/ip/8.8.8.8", "deep=true", func(c *Client, lang string) error {
			_, err := c.IP(ctx, "8.8.8.8", IPOptions{Lang: lang, Deep: true})
			return err
		}},
		{"IPSelf", "/ip", "deep=true", func(c *Client, lang string) error {
			_, err := c.IPSelf(ctx, IPSelfOptions{Lang: lang, Deep: true})
			return err
		}},
		{"Continent", "/continent/EU", "", func(c *Client, lang string) error {
			_, err := c.Continent(ctx, "EU", ContinentOptions{Lang: lang})
			return err
		}},
		{"ContinentCountries", "/continent/EU/countries", "", func(c *Client, lang string) error {
			_, err := c.ContinentCountries(ctx, "EU", ContinentCountriesOptions{Lang: lang})
			return err
		}},
		{"BlocCountries", "/bloc/EU/countries", "", func(c *Client, lang string) error {
			_, err := c.BlocCountries(ctx, "EU", BlocCountriesOptions{Lang: lang})
			return err
		}},
		{"Country", "/country/DE", "deep=true", func(c *Client, lang string) error {
			_, err := c.Country(ctx, "DE", CountryOptions{Lang: lang, Deep: true})
			return err
		}},
		{"CountryStates", "/country/DE/states", "", func(c *Client, lang string) error {
			_, err := c.CountryStates(ctx, "DE", CountryStatesOptions{Lang: lang})
			return err
		}},
		{"State", "/state/CA", "country=US", func(c *Client, lang string) error {
			_, err := c.State(ctx, "CA", StateOptions{Lang: lang, Country: "US"})
			return err
		}},
		{"StateDistricts", "/state/CA/districts", "country=US&deep=true", func(c *Client, lang string) error {
			_, err := c.StateDistricts(ctx, "CA", StateDistrictsOptions{Lang: lang, Country: "US", Deep: true})
			return err
		}},
		{"District", "/district/37081", "country=US&state=NC", func(c *Client, lang string) error {
			_, err := c.District(ctx, "37081", DistrictOptions{Lang: lang, Country: "US", State: "NC"})
			return err
		}},
		{"City", "/city/M%C3%BCnchen", "country=DE", func(c *Client, lang string) error {
			_, err := c.City(ctx, "München", CityOptions{Lang: lang, Country: "DE"})
			return err
		}},
		{"CityID", "/city/id/city_fixture", "deep=true", func(c *Client, lang string) error {
			_, err := c.CityID(ctx, "city_fixture", CityIDOptions{Lang: lang, Deep: true})
			return err
		}},
		{"CitySearch", "/city", "limit=2&q=M%C3%BCn", func(c *Client, lang string) error {
			_, err := c.CitySearch(ctx, "Mün", CitySearchOptions{Lang: lang, Limit: 2})
			return err
		}},
		{"CityNearest", "/city", "lat=0&lon=0", func(c *Client, lang string) error {
			_, err := c.CityNearest(ctx, 0, 0, CityNearestOptions{Lang: lang})
			return err
		}},
		{"CityNearby", "/city/M%C3%BCnchen/nearby", "radius=8&unit=km", func(c *Client, lang string) error {
			_, err := c.CityNearby(ctx, "München", CityNearbyOptions{Lang: lang, Radius: 8, Unit: "km"})
			return err
		}},
		{"Postal", "/postal/SW1A%201AA", "country=GB", func(c *Client, lang string) error {
			_, err := c.Postal(ctx, "SW1A 1AA", PostalOptions{Lang: lang, Country: "GB"})
			return err
		}},
		{"PostalNearby", "/postal/28202/nearby", "country=US&radius=8", func(c *Client, lang string) error {
			_, err := c.PostalNearby(ctx, "28202", PostalNearbyOptions{Lang: lang, Country: "US", Radius: 8})
			return err
		}},
		{"PostalDistance", "/postal/28202/distance/10001", "country=US", func(c *Client, lang string) error {
			_, err := c.PostalDistance(ctx, "28202", "10001", PostalDistanceOptions{Lang: lang, Country: "US"})
			return err
		}},
		{"Company", "/company/732829320", "country=FR&deep=true", func(c *Client, lang string) error {
			_, err := c.Company(ctx, "732829320", CompanyOptions{Lang: lang, Country: "FR", Deep: true})
			return err
		}},
		{"NPI", "/provider/1881018208", "deep=true", func(c *Client, lang string) error {
			_, err := c.Provider(ctx, "1881018208", ProviderOptions{Lang: lang, Deep: true})
			return err
		}},
		{"ASN", "/asn/AS13335", "", func(c *Client, lang string) error {
			_, err := c.ASN(ctx, "AS13335", ASNOptions{Lang: lang})
			return err
		}},
		{"Currency", "/currency/USD", "deep=true", func(c *Client, lang string) error {
			_, err := c.Currency(ctx, "USD", CurrencyOptions{Lang: lang, Deep: true})
			return err
		}},
		{"Language", "/language/ja", "", func(c *Client, lang string) error {
			_, err := c.Language(ctx, "ja", LanguageOptions{Lang: lang})
			return err
		}},
		{"Time", "/time/America%2FNew_York", "at=2026-01-01T12%3A00&deep=true&to=UTC", func(c *Client, lang string) error {
			_, err := c.Time(ctx, "America/New_York", TimeOptions{Lang: lang, At: "2026-01-01T12:00", To: "UTC", Deep: true})
			return err
		}},
		{"TimeAt", "/time", "at=2026-01-01T12%3A00Z&lat=0&lon=0", func(c *Client, lang string) error {
			_, err := c.TimeAt(ctx, 0, 0, TimeAtOptions{Lang: lang, At: "2026-01-01T12:00Z"})
			return err
		}},
		{"Timezone", "/timezone/UTC", "deep=true", func(c *Client, lang string) error {
			_, err := c.Timezone(ctx, "UTC", TimezoneOptions{Lang: lang, Deep: true})
			return err
		}},
		{"TimezoneAt", "/timezone", "deep=true&lat=0&lon=0", func(c *Client, lang string) error {
			_, err := c.TimezoneAt(ctx, 0, 0, TimezoneAtOptions{Lang: lang, Deep: true})
			return err
		}},
		{"Date", "/date/03%2F04%2F2026", "deep=true&format=dmy&to=2026-05-01", func(c *Client, lang string) error {
			_, err := c.Date(ctx, "03/04/2026", DateOptions{Lang: lang, Format: "dmy", To: "2026-05-01", Deep: true})
			return err
		}},
		{"DateToday", "/date", "to=2026-05-01", func(c *Client, lang string) error {
			_, err := c.DateToday(ctx, DateTodayOptions{Lang: lang, To: "2026-05-01"})
			return err
		}},
		{"Point", "/point", "deep=true&lat=0&lon=0", func(c *Client, lang string) error {
			_, err := c.Point(ctx, 0, 0, PointOptions{Lang: lang, Deep: true})
			return err
		}},
		{"Emoji", "/emoji/%F0%9F%98%80", "deep=true", func(c *Client, lang string) error {
			_, err := c.Emoji(ctx, "😀", EmojiOptions{Lang: lang, Deep: true})
			return err
		}},
		{"EmojiSearch", "/emoji", "limit=2&q=visage", func(c *Client, lang string) error {
			_, err := c.EmojiSearch(ctx, "visage", EmojiSearchOptions{Lang: lang, Limit: 2})
			return err
		}},
		{"MeasureUnits", "/measure/units", "q=meter&unit=m", func(c *Client, lang string) error {
			_, err := c.MeasureUnits(ctx, MeasureUnitsOptions{Lang: lang, Query: "meter", Unit: "m"})
			return err
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client, call := newTestClient(t, okJSON(`{"future":null}`))
			for _, lang := range []string{"fr-CA", ""} {
				if err := tc.invoke(client, lang); err != nil {
					t.Fatal(err)
				}
				want, err := url.ParseQuery(tc.query)
				if err != nil {
					t.Fatal(err)
				}
				if lang != "" {
					want.Set("lang", lang)
				}
				if call.path != tc.path || call.rawQry != want.Encode() {
					t.Fatalf("lang=%q: got %s?%s, want %s?%s", lang, call.path, call.rawQry, tc.path, want.Encode())
				}
			}
		})
	}
}

func TestDisplayLanguageKeepsNativeNullAndInputOptions(t *testing.T) {
	ctx := context.Background()
	client, call := newTestClient(t, okJSON(`{"country":"DE","name":"Allemagne","name_local":"Deutschland","currency_name":null,"future":true}`))
	result, err := client.Country(ctx, "DE", CountryOptions{Lang: "fr"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Country != "DE" || result.Name != "Allemagne" || result.NameLocal == nil || *result.NameLocal != "Deutschland" || result.CurrencyName != nil || result.Deep != nil {
		t.Fatalf("changed localized/native/null fields: %#v", result)
	}
	if _, err := client.Date(ctx, "03/04/2026", DateOptions{Format: "dmy", Lang: "en-US"}); err != nil {
		t.Fatal(err)
	}
	if call.path != "/date/03%2F04%2F2026" || call.rawQry != "format=dmy&lang=en-US" {
		t.Fatalf("changed date parsing: %#v", call)
	}
	if _, err := client.Measure(ctx, "1,5 m", MeasureOptions{Locale: "de-DE", To: "cm"}); err != nil {
		t.Fatal(err)
	}
	if call.path != "/measure/1%2C5%20m" || call.rawQry != "locale=de-DE&to=cm" {
		t.Fatalf("changed measure input or leaked lang: %#v", call)
	}
	for _, options := range []any{BlocOptions{}, CurrencyRateOptions{}, MeasureOptions{}, HolidayOptions{}, HolidayDateOptions{}, NameOptions{}, EmailOptions{}, PhoneOptions{}, AddressOptions{}} {
		if _, ok := reflect.TypeOf(options).FieldByName("Lang"); ok {
			t.Fatalf("unsupported language option on %T", options)
		}
	}
}
