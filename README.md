# parseapi-go

Official parseAPI client for Go.

```bash
go get github.com/parseapi/parseapi-go
```

```go
import parseapi "github.com/parseapi/parseapi-go"

parse, err := parseapi.New("your-api-key")
country, err := parse.Country(ctx, "US")
```

Get a key at [parseapi.com](https://parseapi.com). An empty key falls back to the `PARSEAPI_KEY` environment variable.

## Calls

One method per endpoint, named after the route. Every method takes a context first. Optional parameters ride in a per-method options struct, nil means defaults.

```go
parse.IP(ctx, "8.8.8.8", nil)
parse.IPSelf(ctx, nil)
parse.Email(ctx, "hello@gmail.com", nil)
parse.Phone(ctx, "+14155552671", nil)
parse.Postal(ctx, "28202", "US")
parse.PostalNearby(ctx, "28202", "US", &parseapi.PostalNearbyOptions{Radius: 40})
parse.PostalDistance(ctx, "28202", "10001", "US")
parse.City(ctx, "charlotte", &parseapi.CityOptions{Country: "US"})
parse.CityID(ctx, "city_mb8mbqrkz8zb")
parse.CitySearch(ctx, "char", &parseapi.CitySearchOptions{Country: "US", Limit: 10})
parse.CityNearest(ctx, 35.2271, -80.8431)
parse.Country(ctx, "US")
parse.CountryStates(ctx, "US")
parse.State(ctx, "NC", "US")
parse.StateDistricts(ctx, "NC", "US")
parse.District(ctx, "37081", nil)
parse.Continent(ctx, "NA")
parse.ContinentCountries(ctx, "NA")
parse.Currency(ctx, "USD")
parse.CurrencyRate(ctx, "USD", "EUR")
parse.Language(ctx, "en")
parse.Timezone(ctx, "America/New_York", nil)
parse.Holiday(ctx, "US", &parseapi.HolidayOptions{Year: 2026})
parse.HolidayDate(ctx, "US", "2026-12-25")
parse.Elevation(ctx, 35.2271, -80.8431)
parse.Point(ctx, 36.0726, -79.792, nil)
parse.Weather(ctx, 40.7128, -74.006, nil)
parse.Domain(ctx, "example.com", nil)
parse.MX(ctx, "example.com")
parse.Useragent(ctx, uaString, nil)
parse.Emoji(ctx, "rocket")
parse.EmojiSearch(ctx, "fire", nil)
```

Every response is a typed struct. Nullable fields are pointers.

## Deep

Pass deep options to include the nested deep object with richer fields.

```go
ip, err := parse.IP(ctx, "52.94.76.10", &parseapi.DeepOptions{Deep: true})
if ip.Deep != nil && ip.Deep.Datacenter != nil && *ip.Deep.Datacenter {
	// datacenter IP
}
```

## Errors

Every non-2xx response returns a `*parseapi.Error` with `Status`, `Code`, `Docs`, and `RequestID`. Branch on `Code`.

```go
_, err := parse.City(ctx, "atlantis", nil)
var apiErr *parseapi.Error
if errors.As(err, &apiErr) && apiErr.Code == "not_found" {
	// no such city
}
```

## Options

```go
parse, err := parseapi.New("your-api-key",
	parseapi.WithTimeout(10*time.Second), // per-attempt timeout
	parseapi.WithRetries(2),              // automatic retries on network errors, 429, and 5xx
)
```

Requires Go 1.21 or later. Standard library only, zero dependencies.

## Docs

Full field reference for every endpoint: [parseapi.com/docs](https://parseapi.com/docs)
