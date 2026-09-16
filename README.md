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

## API versions

Choose your team's API version in [Dashboard → API version](https://parseapi.com/dashboard/versions). One setting applies to every key, including new and replacement keys. Existing teams keep `1.0.0`; new teams start on `2.0.0`. Keep the same keys and lookup URLs. Installing or upgrading the package does not change the team's setting.

SDK `0.3.2` targets API `1.0.0`. SDK `0.4.0` and the examples and response types in this source tree target API `2.0.0`. Use a package release documented for your team's version. These types do not model every historical response; moving to `2.0.0` may require updating code that reads renamed, moved or removed fields.

Test the target contract in a separate development team before changing your production team's version. A change applies to every integration in that team. See [API versions and migration](https://parseapi.com/docs/versioning).

## Weather from a postal code

Start with the postal code, then pass its coordinates to weather. Reuse the client from the example above.

```go
place, err := parse.Postal(ctx, "28202", parseapi.PostalOptions{Country: "US"})
if err != nil {
    return err
}
if place.Latitude != nil && place.Longitude != nil {
    weather, err := parse.Weather(ctx, *place.Latitude, *place.Longitude)
    if err != nil {
        return err
    }
    fmt.Println(weather)
}
```

The coordinates represent the postal area. Weather is for that point. Missing coordinates skip the weather lookup. This composition performs two ordinary lookups when coordinates are available, with the retry policy below.

Import `fmt` for this example. Place it inside a function that returns `error`.

## Supply the context you know

Pass `country` when a postal code or national phone number needs disambiguation. A complete international phone number already carries its country context. For a numeric date such as `03/04/2026`, supply the intended `format`. Defaults resolve what the input establishes. Ambiguous input needs your context.

Results are plain data. Pass a returned code or coordinate to another operation when the task needs it. Check nullable values before composing the next call.

Name paid deep includes flat `short`, `directory`, and `initials` fields beside `gender` and `salutation`. `NameLocale` selects CLDR formatting rules and defaults to `en`. It changes formatting only. Country remains gender context, and unavailable formatting is null. Older responses may omit these fields.

## Display language

This source candidate accepts an optional language for supported display fields.
It requires the matching API localization release and data.

```go
country, err := parse.Country(ctx, "DE", parseapi.CountryOptions{Lang: "fr"})
if err != nil {
    return err
}
fmt.Println(country.Name) // Allemagne
```

`lang` applies to this request. The next call uses its usual default unless it
also supplies a language. Codes, native names, numeric facts and response
structure stay unchanged. Missing translations keep the API's documented
fallback. Existing `deep` rules still apply; Date `format` and Measure input
`locale` retain their parsing meanings.

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
parse.BIN(ctx, "424242")
parse.NPI(ctx, "1881018208")
parse.ASN(ctx, "AS13335")
parse.MAC(ctx, "00:1B:63:84:45:E6")
parse.Name(ctx, "Andrea", parseapi.NameOptions{Country: "IT"})
parse.Name(ctx, "Robert James Smith", parseapi.NameOptions{Deep: true, NameLocale: "en"})
parse.VIN(ctx, "1HGCM82633A004352")
parse.Carrier(ctx, "+14155552671")
parse.Caller(ctx, "+18004633339")
parse.HLR(ctx, "+447712345678")
parse.UserAgent(ctx, "Mozilla/5.0")
parse.DNS(ctx, "example.com")
parse.DNS(ctx, "_dmarc.example.com", parseapi.DNSOptions{Type: "TXT"})
parse.NAICS(ctx, "541511")
parse.NAICSSearch(ctx, "coffee shop", parseapi.NAICSSearchOptions{Limit: 5})
parse.Tariff(ctx, "8471.30.01.00", parseapi.TariffOptions{Origin: "DE", Deep: true})
parse.Address(ctx, "123 Main St", parseapi.AddressOptions{Country: "US"})
parse.AddressSearch(ctx, "123 Main", parseapi.AddressSearchOptions{Country: "US", State: "NC"})
parse.Company(ctx, "123456789", parseapi.CompanyOptions{Country: "FR"})
parse.Date(ctx, "03/04/2026", parseapi.DateOptions{Format: "mdy"})
parse.DateToday(ctx, parseapi.DateTodayOptions{To: "2026-12-25"})
parse.Time(ctx, "") // UTC now
parse.Time(ctx, "America/New_York", parseapi.TimeOptions{At: "2026-09-05T15:00:00", To: "Asia/Tokyo"})
parse.TimeAt(ctx, 40.7128, -74.006)
parse.Weather(ctx, 40.7128, -74.006, parseapi.WeatherOptions{Deep: true, Date: "2026-09-01"})
```

Paid NAICS `deep` includes full definitions, child categories and classification `exclusions`, each with a description and linked codes. Generic exclusions can have no linked codes. Omitted or null exclusions in older responses remain unknown. Search results keep `country` and `year` on the envelope and optional depth on each result. They also include core `match`: the matched `field` (`name`, `term` or `naics`) and `text`, plus `corrections` with `from` and `to` tokens for typo fallback. Corrections are empty for exact, plural and prefix matches. Direct code lookups omit `match`. Older responses may omit it.

Use named fields when constructing response values for fixtures too. Response and options structs reserve room for future fields and cannot be compared with `==`. Nullable values are pointers. Unknown JSON fields are accepted. An omitted `deep` is nil, a requested empty `deep` is a non-nil object, and unknown fields within it stay nil. Nullable arrays use nil slices.

DNS uses pooled requests on every plan. Omit `type` to check A, AAAA, CNAME, MX, NS, TXT, SOA, CAA, SRV and PTR. Records contain `name`, `type`, `ttl` in seconds and a DNS presentation `value`. TXT values retain quoting and chunk boundaries. A selected question can include its CNAME chain. Empty records mean no records. Lookup failures remain errors.

## Time

`time` returns local ISO `at` with its UTC offset and integer Unix seconds in `unix`. Request `deep` for the display name, exact `offset_seconds`, whole `offset_minutes` and next clock change. A conversion target has its own optional `deep` without a next-change field. Historical offsets and ISO times can include offset seconds. Omitted `at` means now. With `to`, an offsetless `at` is source wall time. Otherwise it is UTC. Include an offset for repeated local times around a clock change. Current time and conversion use pooled requests on every plan. Coordinate clock fields can be null when the timezone is unknown. Existing `timezone` methods remain supported.

## Measurements

```go
result, err := client.Measure(ctx, "5 ft 11 in", parseapi.MeasureOptions{To: "cm"})
if err != nil { return err }
units, err := client.MeasureUnits(ctx, parseapi.MeasureUnitsOptions{Unit: "m"})
if err != nil { return err }
```

`amount` is a decimal string, such as `"180.34"`. Without `to`, the API returns the canonical unit for the measurement type. Pass `locale` for number formatting and `system` (`us` or `imperial`) when a customary unit needs context. Ambiguous input returns `valid: false`, a `reason`, and available `choices`. Invalid or incompatible target units use the normal API error.

Unit discovery accepts optional `query`, `type`, and `unit` filters. `unit` selects compatible targets. Omit the filters for the reviewed catalog. Both operations use pooled requests.

## Place statistics and optional detail

Postal and District paid profiles include `deep.property_tax` where supported. It contains `annual_median`, `currency` and `period`: median annual property tax payable on owner-occupied homes in the statistical area. The amount is adjusted to the final year of the reporting period (`YYYY-YYYY`). This is an area statistic, not a rate or an individual property bill. Unsupported, missing and censored estimates are null.

```go
place, err := parse.Postal(ctx, "28202", parseapi.PostalOptions{Country: "US", Deep: true})
if err != nil {
    log.Fatal(err)
}
if place.Deep != nil && place.Deep.PropertyTax != nil {
    fmt.Println(place.Deep.PropertyTax.AnnualMedian, place.Deep.PropertyTax.Currency, place.Deep.PropertyTax.Period)
}
```

Read `population_period` alongside `population`: a reporting year (`YYYY`) or period (`YYYY-YYYY`), null when unknown or unverifiable. Keep missing or null values unknown and preserve a known zero. These fields belong to full place profiles. State district lists include each district's population and period. Postal nearby and distance detail remains metropolitan associations only. Continent population stays in core.

Point returns the timezone ID with the core location. Its optional deep detail adds terrain and compact nearest-city context on every plan. A nearest city is null when none is within 200 km.

Weather returns current conditions by default. Paid deep adds specialist current measurements, forecasts and related detail. A past `date` is a UTC day and requires deep: it adds `deep.history` alongside current conditions. Date alone does not request history.

```go
parse.Weather(ctx, 40.7128, -74.006, parseapi.WeatherOptions{Deep: true, Date: "2026-08-15"})
```

Tariff starts with the general schedule line. Paid deep adds units and the special and other schedule columns. An optional origin then resolves country-specific measures. The three calls below show those successive choices. Without origin, schedule detail is still returned and origin-dependent fields are null. A null effective rate is not a zero rate.

```go
parse.Tariff(ctx, "8471.30.01.00")
parse.Tariff(ctx, "8471.30.01.00", parseapi.TariffOptions{Deep: true})
parse.Tariff(ctx, "8471.30.01.00", parseapi.TariffOptions{Deep: true, Origin: "CN"})
```

Address search uses context from the form: prefer postal, or city and state. An optional end-user `ip` is a locality hint for server-side calls. An empty result explains itself with `reason`: `more_input`, `missing_context` or `no_matches`. With suggestions, reason is null. Older responses may omit it, and future reasons remain strings. Catalog and lookup failures use the existing API errors.

HLR reports status at the last check. `live` means assigned and `connected` means reachable at that check. Cached results may be returned. Null means unconfirmed. Deep diagnostics stay within the same metered lookup.

## Deep

The default call returns the common answer. Request more detail with `parse.Country(ctx, "US", parseapi.CountryOptions{Deep: true})`. Read those fields from the optional deep member; this does not change the core answer.

| Operation | What `deep` requests |
|---|---|
| IP | Richer IP fields included with a paid plan. No separate check meter. |
| Domain | Registration dates, registrar, status and DNSSEC, included with a paid plan. Use `DNS` for DNS records and `MX` for mail routing. |
| Email | A metered deliverability check, using included email checks or enabled on-demand usage. |
| VAT | A metered registry check where supported, using included VAT checks or enabled on-demand usage. |
| Country, State, City, District, Postal | Reference profiles included with a paid plan; place identity and coordinates stay core. |
| VIN, NPI, NAICS, Company | Paid technical or registration profiles. NPI exclusion status and NAICS hierarchy stay core. |
| Tariff | Paid schedule columns and units; add origin for applicable measures. |
| Name, Weather | Paid name context or weather detail; parsing and current conditions stay core. |
| Phone, IBAN | Numbering-plan or bank structure detail in the same pooled request on every plan. |
| Time, Date, Currency, Language, Emoji, Point | Optional reference detail in the same pooled request on every plan. |
| Carrier, HLR | Available place or network detail from the same metered core unit, including Free included units. |

Carrier, caller, and HLR are separate metered operations. Choose them explicitly when you need their answers. Ordinary lookups retry twice by default. Metered checks use one attempt by default. Setting retries explicitly can repeat paid usage.

Without `deep`, the response omits that key. When requested, it is an empty object if access is locked or the operation has no deep fields. Otherwise it contains the available fields. A missing or null field means unknown.

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

BIN lookup accepts 6-11 digits as a string, including leading zeros. Spaces and hyphens are accepted. `prefix` is the actual longest match and can be shorter than the input. Unknown reference fields are null. `deep` adds an empty object on every plan.
