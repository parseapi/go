// Live smoke against the edge. Canary-ready: env-driven, clean exit codes.
//
//	PARSEAPI_KEY       required
//	PARSEAPI_BASE_URL  optional override
//
// Run: go run ./smoke
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	parseapi "github.com/parseapi/go"
)

var (
	failures = 0
	total    = 0
)

func check(name string, ok bool, detail string) {
	total++
	status := "ok  "
	if !ok {
		failures++
		status = "FAIL"
	}
	if detail != "" {
		fmt.Printf("%s %s (%s)\n", status, name, detail)
	} else {
		fmt.Printf("%s %s\n", status, name)
	}
}

func expectOk[T any](name string, result *T, err error, assert func(*T) string) {
	if err != nil {
		var apiErr *parseapi.Error
		if errors.As(err, &apiErr) {
			check(name, false, fmt.Sprintf("%d %s", apiErr.Status, apiErr.Code))
		} else {
			check(name, false, err.Error())
		}
		return
	}
	problem := ""
	if assert != nil {
		problem = assert(result)
	}
	check(name, problem == "", problem)
}

func expectError(name string, err error, code string) {
	if err == nil {
		check(name, false, "expected error, got 200")
		return
	}
	var apiErr *parseapi.Error
	if errors.As(err, &apiErr) {
		check(name, apiErr.Code == code, "got "+apiErr.Code)
		return
	}
	check(name, false, err.Error())
}

const testUA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

func str(value *string) string {
	if value == nil {
		return "<nil>"
	}
	return *value
}

func main() {
	ctx := context.Background()
	parse, err := parseapi.New("")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ip, err := parse.IP(ctx, "8.8.8.8")
	expectOk("ip", ip, err, func(r *parseapi.IP) string {
		if r.IP != "8.8.8.8" {
			return "wrong ip"
		}
		return ""
	})

	me, err := parse.IPSelf(ctx)
	expectOk("ip self", me, err, func(r *parseapi.IP) string {
		if r.IP == "" {
			return "no ip"
		}
		return ""
	})

	continent, err := parse.Continent(ctx, "NA")
	expectOk("continent", continent, err, func(r *parseapi.Continent) string {
		if r.Name != "North America" {
			return "wrong name"
		}
		return ""
	})

	continentCountries, err := parse.ContinentCountries(ctx, "NA")
	expectOk("continent countries", continentCountries, err, func(r *parseapi.ContinentCountries) string {
		if len(r.Countries) == 0 {
			return "empty"
		}
		return ""
	})

	bloc, err := parse.Bloc(ctx, "EU")
	expectOk("bloc", bloc, err, func(r *parseapi.Bloc) string {
		if r.Name != "European Union" || r.Members != 27 {
			return "wrong bloc"
		}
		return ""
	})

	blocCountries, err := parse.BlocCountries(ctx, "SCHENGEN")
	expectOk("bloc countries", blocCountries, err, func(r *parseapi.BlocCountries) string {
		if len(r.Countries) != 29 {
			return "wrong members"
		}
		return ""
	})

	country, err := parse.Country(ctx, "US")
	expectOk("country", country, err, func(r *parseapi.Country) string {
		if r.ISO3 != "USA" {
			return "wrong iso3"
		}
		return ""
	})

	states, err := parse.CountryStates(ctx, "US")
	expectOk("country states", states, err, func(r *parseapi.CountryStates) string {
		if len(r.States) < 50 {
			return "too few states"
		}
		return ""
	})

	state, err := parse.State(ctx, "NC", parseapi.StateOptions{Country: "US"})
	expectOk("state", state, err, func(r *parseapi.State) string {
		if r.Name != "North Carolina" {
			return "wrong name"
		}
		return ""
	})

	districts, err := parse.StateDistricts(ctx, "NC", parseapi.StateDistrictsOptions{Country: "US"})
	expectOk("state districts", districts, err, func(r *parseapi.StateDistricts) string {
		if len(r.Districts) == 0 {
			return "empty"
		}
		return ""
	})

	district, err := parse.District(ctx, "37081")
	expectOk("district", district, err, func(r *parseapi.District) string {
		if r.Name == "" {
			return "no name"
		}
		return ""
	})

	city, err := parse.City(ctx, "charlotte", parseapi.CityOptions{Country: "US"})
	expectOk("city", city, err, func(r *parseapi.City) string {
		if r.Name != "Charlotte" {
			return "wrong city"
		}
		if r.ID == "" || len(r.ID) < 5 || r.ID[:5] != "city_" {
			return "missing id"
		}
		return ""
	})

	if city != nil && city.ID != "" {
		byID, err := parse.CityID(ctx, city.ID)
		expectOk("city id", byID, err, func(r *parseapi.City) string {
			if r.ID != city.ID || r.Name != "Charlotte" {
				return "id mismatch"
			}
			return ""
		})
	}

	citySearch, err := parse.CitySearch(ctx, "char", parseapi.CitySearchOptions{Country: "US", Limit: 5})
	expectOk("city search", citySearch, err, func(r *parseapi.CitySearch) string {
		if len(r.Cities) == 0 {
			return "empty"
		}
		return ""
	})

	cityNearest, err := parse.CityNearest(ctx, 35.2271, -80.8431)
	expectOk("city nearest", cityNearest, err, func(r *parseapi.CityNearest) string {
		if r.Name == "" {
			return "no city"
		}
		return ""
	})

	postal, err := parse.Postal(ctx, "28202", parseapi.PostalOptions{Country: "US"})
	expectOk("postal", postal, err, func(r *parseapi.Postal) string {
		if str(r.City) != "Charlotte" {
			return "city " + str(r.City)
		}
		return ""
	})

	nearby, err := parse.PostalNearby(ctx, "28202", parseapi.PostalNearbyOptions{Country: "US", Radius: 40})
	expectOk("postal nearby", nearby, err, func(r *parseapi.PostalNearby) string {
		if len(r.Nearby) == 0 {
			return "empty"
		}
		return ""
	})

	distance, err := parse.PostalDistance(ctx, "28202", "10001", parseapi.PostalDistanceOptions{Country: "US"})
	expectOk("postal distance", distance, err, func(r *parseapi.PostalDistance) string {
		if r.Distance < 800 || r.Distance > 1000 {
			return fmt.Sprintf("distance %.1f", r.Distance)
		}
		return ""
	})

	email, err := parse.Email(ctx, "hello@gmail.com")
	expectOk("email", email, err, func(r *parseapi.Email) string {
		if !r.Valid {
			return "not valid"
		}
		return ""
	})

	vat, err := parse.VAT(ctx, "DE136695976")
	expectOk("vat", vat, err, func(r *parseapi.VAT) string {
		if !r.Valid || str(r.Country) != "DE" {
			return "not valid DE"
		}
		return ""
	})

	iban, err := parse.IBAN(ctx, "DE89370400440532013000")
	expectOk("iban", iban, err, func(r *parseapi.IBAN) string {
		if !r.Valid || str(r.Country) != "DE" || str(r.Bank) != "37040044" {
			return "not valid DE"
		}
		return ""
	})

	ibanJunk, err := parse.IBAN(ctx, "hello")
	expectOk("iban junk", ibanJunk, err, func(r *parseapi.IBAN) string {
		if r.Valid {
			return "expected invalid"
		}
		return ""
	})

	npi, err := parse.NPI(ctx, "1881018208")
	expectOk("npi", npi, err, func(r *parseapi.NPI) string {
		if !r.Valid || r.Registered == nil || !*r.Registered {
			return "not registered"
		}
		return ""
	})

	npiJunk, err := parse.NPI(ctx, "hello")
	expectOk("npi junk", npiJunk, err, func(r *parseapi.NPI) string {
		if r.Valid {
			return "expected invalid"
		}
		return ""
	})

	phone, err := parse.Phone(ctx, "+14155552671")
	expectOk("phone", phone, err, func(r *parseapi.Phone) string {
		if str(r.Phone) != "+14155552671" {
			return "phone " + str(r.Phone)
		}
		return ""
	})

	// Metered core siblings: junk numbers return valid false without a metered lookup.
	carrier, err := parse.Carrier(ctx, "555-0100")
	expectOk("carrier junk free", carrier, err, func(r *parseapi.Carrier) string {
		if r.Valid {
			return "expected invalid"
		}
		return ""
	})

	caller, err := parse.Caller(ctx, "555-0100")
	expectOk("caller junk free", caller, err, func(r *parseapi.Caller) string {
		if r.Valid {
			return "expected invalid"
		}
		return ""
	})

	hlr, err := parse.HLR(ctx, "555-0100")
	expectOk("hlr junk free", hlr, err, func(r *parseapi.HLR) string {
		if r.Valid {
			return "expected invalid"
		}
		return ""
	})

	domain, err := parse.Domain(ctx, "gmail.com")
	expectOk("domain", domain, err, func(r *parseapi.Domain) string {
		if r.Available {
			return "gmail available?"
		}
		return ""
	})

	asn, err := parse.ASN(ctx, "AS13335")
	expectOk("asn", asn, err, func(r *parseapi.ASN) string {
		if r.ASN != 13335 {
			return "wrong ASN"
		}
		return ""
	})
	mac, err := parse.MAC(ctx, "00:1B:63:84:45:E6")
	expectOk("mac", mac, err, func(r *parseapi.MAC) string {
		if !r.Valid || r.MAC != "00:1B:63:84:45:E6" || r.Local == nil || *r.Local || r.Multicast == nil || *r.Multicast {
			return "wrong MAC"
		}
		return ""
	})

	mx, err := parse.MX(ctx, "gmail.com")
	expectOk("mx", mx, err, func(r *parseapi.MX) string {
		if len(r.MX) == 0 {
			return "no mx"
		}
		return ""
	})

	ua, err := parse.UserAgent(ctx, testUA)
	expectOk("useragent", ua, err, func(r *parseapi.UserAgent) string {
		if str(r.Browser) != "Chrome" {
			return "browser " + str(r.Browser)
		}
		return ""
	})

	vin, err := parse.VIN(ctx, "1HGCM82633A004352")
	expectOk("vin", vin, err, func(r *parseapi.VIN) string {
		if !r.Valid || str(r.Make) != "Honda" {
			return "make " + str(r.Make)
		}
		return ""
	})

	vinJunk, err := parse.VIN(ctx, "1HGCM82613A004352")
	expectOk("vin junk", vinJunk, err, func(r *parseapi.VIN) string {
		if r.Valid {
			return "expected invalid"
		}
		return ""
	})

	currency, err := parse.Currency(ctx, "USD")
	expectOk("currency", currency, err, func(r *parseapi.Currency) string {
		if str(r.Symbol) != "$" {
			return "symbol " + str(r.Symbol)
		}
		return ""
	})

	rate, err := parse.CurrencyRate(ctx, "USD", "EUR")
	expectOk("currency rate", rate, err, func(r *parseapi.CurrencyRate) string {
		if r.Rate <= 0 || r.Rate >= 10 {
			return fmt.Sprintf("rate %f", r.Rate)
		}
		return ""
	})

	language, err := parse.Language(ctx, "en")
	expectOk("language", language, err, func(r *parseapi.Language) string {
		if r.Language != "en" || r.Name != "English" {
			return "wrong language"
		}
		return ""
	})

	name, err := parse.Name(ctx, "BILLY O'SHALL")
	expectOk("name", name, err, func(r *parseapi.Name) string {
		if r.Name != "Billy O'Shall" || !r.Valid || r.Gender == nil || *r.Gender != "male" {
			return "wrong name"
		}
		return ""
	})

	timezone, err := parse.Timezone(ctx, "America/New_York")
	expectOk("timezone", timezone, err, func(r *parseapi.Timezone) string {
		if r.OffsetMinutes == nil || (*r.OffsetMinutes != -240 && *r.OffsetMinutes != -300) {
			return fmt.Sprintf("offset %d", r.OffsetMinutes)
		}
		return ""
	})

	holidays, err := parse.Holiday(ctx, "US")
	expectOk("holiday", holidays, err, func(r *parseapi.HolidayYear) string {
		if len(r.Holidays) <= 5 {
			return "too few"
		}
		return ""
	})

	christmas, err := parse.HolidayDate(ctx, "US", "2026-12-25")
	expectOk("holiday date", christmas, err, func(r *parseapi.HolidayDate) string {
		if r.Holiday == nil || r.Holiday.Name != "Christmas Day" {
			return "not christmas"
		}
		return ""
	})

	notHoliday, err := parse.HolidayDate(ctx, "US", "2026-08-12")
	expectOk("holiday null", notHoliday, err, func(r *parseapi.HolidayDate) string {
		if r.Holiday != nil {
			return "expected nil"
		}
		return ""
	})

	elevation, err := parse.Elevation(ctx, 35.2271, -80.8431)
	expectOk("elevation", elevation, err, func(r *parseapi.Elevation) string {
		if r.Elevation == nil {
			return "no elevation"
		}
		return ""
	})

	point, err := parse.Point(ctx, 36.0726, -79.792)
	expectOk("point", point, err, func(r *parseapi.Point) string {
		if str(r.Country) != "US" {
			return "country " + str(r.Country)
		}
		return ""
	})

	weather, err := parse.Weather(ctx, 40.7128, -74.006)
	expectOk("weather", weather, err, func(r *parseapi.Weather) string {
		if r.Station == nil || r.Station.ID == "" {
			return "no station"
		}
		return ""
	})

	emoji, err := parse.Emoji(ctx, "rocket")
	expectOk("emoji", emoji, err, func(r *parseapi.Emoji) string {
		if r.Emoji != "\U0001F680" {
			return "wrong emoji"
		}
		return ""
	})

	emojiSearch, err := parse.EmojiSearch(ctx, "fire", parseapi.EmojiSearchOptions{Limit: 5})
	expectOk("emoji search", emojiSearch, err, func(r *parseapi.EmojiSearch) string {
		if len(r.Emojis) == 0 {
			return "empty"
		}
		return ""
	})

	pointDeep, err := parse.Point(ctx, 36.0726, -79.792, parseapi.PointOptions{Deep: true})
	expectOk("point deep triad", pointDeep, err, func(r *parseapi.Point) string {
		if r.Deep == nil {
			return "deep missing"
		}
		return ""
	})

	_, err = parse.City(ctx, "notarealcityxyz")
	expectError("honest 404", err, "not_found")

	bogus, _ := parseapi.New("bogus_key_123", parseapi.WithRetries(0))
	_, err = bogus.Country(ctx, "US")
	expectError("bogus key 401", err, "invalid_api_key")

	fmt.Printf("\n%d/%d passed\n", total-failures, total)
	if failures > 0 {
		os.Exit(1)
	}
}
