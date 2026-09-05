```bash
go get github.com/parseapi/go
```

```go
parse, err := parseapi.New("your-api-key")
if err != nil {
    return err
}
country, err := parse.Country(ctx, "US")
```

Import `parseapi "github.com/parseapi/go"`. Every call takes a `context.Context` first and returns a typed result plus an error. Check the error before using the result. Get a key at [parseapi.com](https://parseapi.com). An empty key reads `PARSEAPI_KEY`.

## Calls

Choose the operation and pass what you have. Related operations are separate direct calls, and results are plain data.

```go
country, err := parse.Country(ctx, "US")
states, err := parse.CountryStates(ctx, "US")
postal, err := parse.Postal(ctx, "28202", parseapi.PostalOptions{Country: "US"})
city, err := parse.City(ctx, "charlotte", parseapi.CityOptions{State: "NC", Country: "US"})
phone, err := parse.Phone(ctx, "+14155552671")
```

Optional query inputs use one options value. Omit it to use defaults. Use named fields, such as `PostalOptions{Country: "US"}`. Passing more than one options value returns an error before making a request. Every operation reserves its own options type, including operations whose options are currently empty. New optional fields can be added without changing calls or stored method signatures.

```go
parse.IP(ctx, "8.8.8.8", parseapi.IPOptions{Deep: true})
parse.Email(ctx, "hello@example.com", parseapi.EmailOptions{Deep: true})
parse.VAT(ctx, "DE136695976", parseapi.VATOptions{Deep: true})
parse.IBAN(ctx, "DE89370400440532013000")
parse.NPI(ctx, "1881018208")
parse.ASN(ctx, "AS13335")
parse.MAC(ctx, "00:1B:63:84:45:E6")
parse.VIN(ctx, "1HGCM82633A004352")
parse.Carrier(ctx, "+14155552671")
parse.Caller(ctx, "+18004633339")
parse.HLR(ctx, "+447712345678")
parse.UserAgent(ctx, "Mozilla/5.0")
parse.Tariff(ctx, "8471.30.01.00", parseapi.TariffOptions{Origin: "DE", Deep: true})
parse.Address(ctx, "123 Main St", parseapi.AddressOptions{Country: "US"})
parse.AddressSearch(ctx, "123 Main", parseapi.AddressSearchOptions{Country: "US", State: "NC"})
parse.Company(ctx, "123456789", parseapi.CompanyOptions{Country: "FR"})
parse.Date(ctx, "03/04/2026", parseapi.DateOptions{Format: "mdy"})
parse.DateToday(ctx, parseapi.DateTodayOptions{To: "2026-12-25"})
parse.Timezone(ctx, "America/New_York", parseapi.TimezoneOptions{At: "2026-09-05T15:00:00", To: "Asia/Tokyo"})
parse.TimezoneAt(ctx, 40.7128, -74.006)
parse.Weather(ctx, 40.7128, -74.006, parseapi.WeatherOptions{Deep: true, Date: "2026-09-01"})
```

Use named fields when constructing response values for fixtures too. Response and options structs reserve room for future fields and cannot be compared with `==`. Nullable values are pointers. Unknown JSON fields are accepted. An omitted `deep` is nil, a requested empty `deep` is a non-nil object, and unknown fields within it stay nil. Nullable arrays use nil slices.

## Errors

Every non-2xx response returns `*parseapi.Error` with `Status`, `Code`, `Message`, `Docs`, and `RequestID`. Branch on `Code`.

```go
_, err := parse.City(ctx, "atlantis")
var apiErr *parseapi.Error
if errors.As(err, &apiErr) && apiErr.Code == "not_found" {
    // No matching city.
}
```

Network failures, invalid JSON, and context cancellation retain their native Go error types. A failed call returns a nil result.

## Requests and retries

Create one client and share it across goroutines. A context deadline covers the whole call, including retry waits. The client timeout defaults to 10 seconds per attempt. Cancellation stops the request and any retry wait.

Ordinary lookups retry twice on network failures, 429, 500, 502, 503, and 504. Carrier, caller, and HLR calls use one attempt by default. Deep email, VAT, and address calls also use one attempt, reserving that behavior for address verification. Address deep currently returns an empty object. An explicit retry setting applies to every call, including metered operations. Additional attempts can be billed.

```go
parse, err := parseapi.New("your-api-key",
    parseapi.WithTimeout(5*time.Second),
    parseapi.WithRetries(0),
)
```

`WithRetries(0)` disables all automatic retries. Both numeric and HTTP-date `Retry-After` values are honored, capped at five seconds. Built-in requests do not follow redirects. `WithHTTPClient` copies your client and keeps redirects disabled. Custom transports must also keep credentials on the requested origin.

Requires Go 1.21 or later. Standard library only.

[Full endpoint and field reference](https://parseapi.com/docs)
