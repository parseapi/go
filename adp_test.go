package parseapi

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestADPModelsPreserveCoreAndDeepTriad(t *testing.T) {
	cases := []struct {
		name   string
		core   string
		detail string
		model  any
	}{
		{"Country", `{"country":"US","name":"United States","continent":"NA"}`, `{"iso3":"USA","population":0,"tax_rate":0,"plugs":[]}`, &Country{}},
		{"State", `{"state":"NC","name":"North Carolina","country":"US"}`, `{"population":0,"area":0.25,"tax_rate":0}`, &State{}},
		{"District", `{"district":"37081","name":"Guilford","country":"US"}`, `{"population":0,"water_area":0.25}`, &District{}},
		{"City", `{"name":"Charlotte","country":"US","id":"city_123"}`, `{"population":0,"area":0.25}`, &City{}},
		{"Postal", `{"postal":"28202","country":"US"}`, `{"metros":[],"water_area":0.25,"tax_rate":0}`, &Postal{}},
		{"Bank", `{"iban":"DE89370400440532013000","valid":true}`, `{"checksum":"89","branch":null,"account":"0532013000"}`, &Bank{}},
		{"NPI", `{"npi":"1881018208","valid":true,"excluded":true,"credential":"MD","state_name":"Minnesota"}`, `{"deactivated_at":"2026-09-01","enrollments":[]}`, &Provider{}},
		{"VIN", `{"vin":"1HGCM82633A004352","valid":true,"make":"Honda"}`, `{"horsepower":240.5,"recalls":[]}`, &VIN{}},
		{"Phone", `{"phone":"+14155552671","valid":true}`, `{"state":"CA","timezone":"America/Los_Angeles"}`, &Phone{}},
		{"Carrier", `{"phone":"+14155552671","valid":true,"carrier":"Example"}`, `{"city":"San Francisco","state":"CA"}`, &Carrier{}},
		{"HLR", `{"phone":"+14155552671","valid":true,"live":true,"connected":false}`, `{"roaming":false,"mcc":"310","mnc":"01"}`, &HLR{}},
		{"Tariff", `{"hts":"8471.30.01.00","description":"Portable computers","revision":"2026"}`, `{"units":[],"special":"Free","origin":null,"effective_rate":null,"measures":null}`, &Tariff{}},
		{"NAICS", `{"naics":"541511","name":"Programming","level":6,"parent":"54151","year":2022,"country":"US"}`, `{"description":"Definition","children":[],"exclusions":[]}`, &NAICS{}},
		{"Company", `{"company":"01234567","valid":true,"name":"Example"}`, `{"activity":"6201","gst":false,"vat":null}`, &Company{}},
		{"Currency", `{"currency":"USD","name":"US Dollar"}`, `{"numeric":840,"countries":[]}`, &Currency{}},
		{"Language", `{"language":"en","name":"English","direction":"ltr"}`, `{"iso3":"eng","countries":[]}`, &Language{}},
		{"Name", `{"name":"Andrea","valid":true,"first":"Andrea"}`, `{"gender":null,"salutation":null}`, &Name{}},
		{"Timezone", `{"timezone":"UTC","unix":0,"at":"1970-01-01T00:00:00+00:00","offset":"+00:00","dst":false}`, `{"name":"UTC","offset_seconds":0,"offset_minutes":0,"next_dst":null}`, &Timezone{}},
		{"DateInfo", `{"date":"1970-01-01","valid":true,"unix":0}`, `{"year":1970,"weekday":4,"leap":false}`, &DateInfo{}},
		{"Emoji", `{"emoji":"😀","name":"grinning face","shortcodes":[]}`, `{"hex":"1F600","skins":[]}`, &Emoji{}},
		{"Point", `{"latitude":0,"longitude":0,"timezone":"Etc/GMT"}`, `{"elevation":0.25,"city":{"name":"Place","id":"city_123","type":"city","distance":0.5}}`, &Point{}},
		{"Weather", `{"latitude":0,"longitude":0,"current":{"temperature":0,"observed_at":"2026-09-08T12:00:00Z"}}`, `{"current":{"pressure":1012.5,"wind_gust":0},"forecast":[]}`, &Weather{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(tc.core), tc.model); err != nil {
				t.Fatal(err)
			}
			value := reflect.ValueOf(tc.model).Elem()
			if !value.FieldByName("Deep").IsNil() {
				t.Fatal("omitted deep became populated")
			}
			plain, _ := json.Marshal(tc.model)
			var body map[string]json.RawMessage
			_ = json.Unmarshal([]byte(tc.core), &body)
			for _, detail := range []string{"{}", tc.detail} {
				body["deep"] = json.RawMessage(detail)
				wire, _ := json.Marshal(body)
				out := reflect.New(value.Type()).Interface()
				if err := json.Unmarshal(wire, out); err != nil {
					t.Fatal(err)
				}
				deep := reflect.ValueOf(out).Elem().FieldByName("Deep")
				if deep.IsNil() {
					t.Fatal("requested deep was erased")
				}
				deep.Set(reflect.Zero(deep.Type()))
				core, _ := json.Marshal(out)
				if string(core) != string(plain) {
					t.Fatalf("deep changed core: %s versus %s", core, plain)
				}
			}
		})
	}
}

func TestADPOptionsSendOneExplicitDisclosure(t *testing.T) {
	ctx := context.Background()
	calls := []func(*Client) error{
		func(c *Client) error { _, err := c.Country(ctx, "US", CountryOptions{Deep: true}); return err },
		func(c *Client) error { _, err := c.State(ctx, "NC", StateOptions{Deep: true}); return err },
		func(c *Client) error {
			_, err := c.StateDistricts(ctx, "NC", StateDistrictsOptions{Deep: true})
			return err
		},
		func(c *Client) error { _, err := c.District(ctx, "37081", DistrictOptions{Deep: true}); return err },
		func(c *Client) error { _, err := c.City(ctx, "Charlotte", CityOptions{Deep: true}); return err },
		func(c *Client) error { _, err := c.CityID(ctx, "city_123", CityIDOptions{Deep: true}); return err },
		func(c *Client) error { _, err := c.CitySearch(ctx, "char", CitySearchOptions{Deep: true}); return err },
		func(c *Client) error { _, err := c.CityNearest(ctx, 0, 0, CityNearestOptions{Deep: true}); return err },
		func(c *Client) error {
			_, err := c.CityNearby(ctx, "Charlotte", CityNearbyOptions{Deep: true})
			return err
		},
		func(c *Client) error { _, err := c.Postal(ctx, "28202", PostalOptions{Deep: true}); return err },
		func(c *Client) error {
			_, err := c.PostalNearby(ctx, "28202", PostalNearbyOptions{Deep: true})
			return err
		},
		func(c *Client) error {
			_, err := c.PostalDistance(ctx, "28202", "10001", PostalDistanceOptions{Deep: true})
			return err
		},
		func(c *Client) error {
			_, err := c.Bank(ctx, "DE89370400440532013000", BankOptions{Deep: true})
			return err
		},
		func(c *Client) error {
			_, err := c.Carrier(ctx, "+14155552671", CarrierOptions{Deep: true})
			return err
		},
		func(c *Client) error { _, err := c.HLR(ctx, "+14155552671", HLROptions{Deep: true}); return err },
		func(c *Client) error { _, err := c.NAICS(ctx, "541511", NAICSOptions{Deep: true}); return err },
		func(c *Client) error {
			_, err := c.NAICSSearch(ctx, "software", NAICSSearchOptions{Deep: true})
			return err
		},
		func(c *Client) error { _, err := c.Currency(ctx, "USD", CurrencyOptions{Deep: true}); return err },
		func(c *Client) error { _, err := c.Language(ctx, "en", LanguageOptions{Deep: true}); return err },
		func(c *Client) error { _, err := c.Name(ctx, "Andrea", NameOptions{Deep: true}); return err },
		func(c *Client) error { _, err := c.Time(ctx, "UTC", TimeOptions{Deep: true}); return err },
		func(c *Client) error { _, err := c.TimeAt(ctx, 0, 0, TimeAtOptions{Deep: true}); return err },
		func(c *Client) error { _, err := c.Timezone(ctx, "UTC", TimezoneOptions{Deep: true}); return err },
		func(c *Client) error { _, err := c.TimezoneAt(ctx, 0, 0, TimezoneAtOptions{Deep: true}); return err },
		func(c *Client) error { _, err := c.Date(ctx, "1970-01-01", DateOptions{Deep: true}); return err },
		func(c *Client) error { _, err := c.DateToday(ctx, DateTodayOptions{Deep: true}); return err },
		func(c *Client) error { _, err := c.Emoji(ctx, "😀", EmojiOptions{Deep: true}); return err },
		func(c *Client) error {
			_, err := c.EmojiSearch(ctx, "smile", EmojiSearchOptions{Deep: true})
			return err
		},
	}
	for _, call := range calls {
		client, got := newTestClient(t, okJSON(`{}`))
		if err := call(client); err != nil {
			t.Fatal(err)
		}
		if got.path == "/bank" {
			var body map[string]any
			if json.Unmarshal(got.body, &body) != nil || body["deep"] != true {
				t.Fatal("missing Bank body deep")
			}
			continue
		}
		if got.rawQry != "deep=true" && !strings.Contains(got.rawQry, "deep=true&") && !strings.Contains(got.rawQry, "&deep=true") {
			t.Fatalf("missing explicit disclosure: %s", got.rawQry)
		}
	}
}

func TestADPChildProfilesAndConversionDepth(t *testing.T) {
	var listing StateDistricts
	if err := json.Unmarshal([]byte(`{"state":"NC","country":"US","districts":[{"district":"37081","name":"Guilford","deep":{"population":0}}]}`), &listing); err != nil {
		t.Fatal(err)
	}
	if d := listing.Districts[0].Deep; d == nil || d.Population == nil || *d.Population != 0 {
		t.Fatal("lost district zero")
	}
	var clock Time
	if err := json.Unmarshal([]byte(`{"timezone":"UTC","unix":0,"deep":{"offset_seconds":0,"next_dst":null},"to":{"timezone":"UTC","unix":0,"deep":{"offset_seconds":0}}}`), &clock); err != nil {
		t.Fatal(err)
	}
	if clock.To == nil || clock.To.Deep == nil || clock.To.Deep.OffsetSeconds == nil || *clock.To.Deep.OffsetSeconds != 0 {
		t.Fatal("lost target depth")
	}
	if _, found := reflect.TypeOf(*clock.To.Deep).FieldByName("NextDST"); found {
		t.Fatal("target invents a next change")
	}
}
