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

Version 1.3.0 sends `Parse-Version: 2.0.0` on every request, including retries. Its response types match API `2.0.0`, and the client selects that contract automatically. No extra constructor setting or key change is needed. This behavior requires the matching API request-version release.

The team setting in [Dashboard API version](https://parseapi.com/dashboard/versions) is the default for requests without a version header. This SDK's header takes precedence without changing that saved default. Existing published packages keep their documented behavior.

Test the new SDK dependency in staging, then deploy the same locked dependency with your application code and existing production key. Future major SDK upgrades can deliberately select a newer API contract, so review their migration notes before upgrading. Rolling back the code and dependency restores the contract selected by that SDK release. If the older SDK does not send a version header, its requests use the team default, which must stay unchanged through that rollback window.

Keep the package version locked in your dependency configuration or lockfile. The selected API contract stays fixed across releases within this planned SDK major. See [API versions and migration](https://parseapi.com/docs/versioning).

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
parse.Bank(ctx, "DE89370400440532013000")
parse.Card(ctx, "424242")
parse.Provider(ctx, "1881018208")
parse.ASN(ctx, "AS13335")
parse.MAC(ctx, "00:1B:63:84:45:E6")
parse.Name(ctx, "Andrea", parseapi.NameOptions{Country: "IT"})
parse.Name(ctx, "Robert James Smith", parseapi.NameOptions{Deep: true, NameLocale: "en"})
parse.Vehicle(ctx, "1HGCM82633A004352")
parse.Carrier(ctx, "+14155552671")
parse.Caller(ctx, "+18004633339")
parse.HLR(ctx, "+447712345678")
parse.UserAgent(ctx, "Mozilla/5.0")
parse.DNS(ctx, "example.com")
parse.DNS(ctx, "_dmarc.example.com", parseapi.DNSOptions{Type: "TXT"})
parse.Industry(ctx, "541511")
parse.IndustrySearch(ctx, "coffee shop", parseapi.IndustrySearchOptions{Limit: 5})
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

Paid Industry `deep` includes full definitions, child categories and classification `exclusions`, each with a description and linked codes. Generic exclusions can have no linked codes. Omitted or null exclusions in older responses remain unknown. Search results keep `country` and `year` on the envelope and optional depth on each result. They also include core `match`: the matched `field` (`name`, `term` or `naics`) and `text`, plus `corrections` with `from` and `to` tokens for typo fallback. Corrections are empty for exact, plural and prefix matches. Direct code lookups omit `match`. Older responses may omit it.

Use named fields when constructing response values for fixtures too. Response and options structs reserve room for future fields and cannot be compared with `==`. Nullable values are pointers. Unknown JSON fields are accepted. An omitted `deep` is nil, a requested empty `deep` is a non-nil object, and unknown fields within it stay nil. Nullable arrays use nil slices.

DNS uses pooled requests on every plan. Omit `type` to check A, AAAA, CNAME, MX, NS, TXT, SOA, CAA, SRV and PTR. Records contain `name`, `type`, `ttl` in seconds and a DNS presentation `value`. TXT values retain quoting and chunk boundaries. A selected question can include its CNAME chain. Empty records mean no records. Lookup failures remain errors.

## Company directory

`CompanyID`, `CompanySearch`, `CompanyCoverage` retrieve directory profiles, search candidates and edition coverage. Existing national company-number validation stays unchanged. Use at most one of query (sent as `q`), domain, ticker or identifier, or discover by country or exact industry; country filters the candidates, exchange narrows a ticker and authority narrows an identifier. The API validates combinations. Each call makes one request and does not automatically resolve candidates or fetch linked assets.

```go
page, err := parse.CompanySearch(ctx, parseapi.CompanySearchOptions{Domain: "cloudflare.com", Deep: true})
if err != nil { return err }
if len(page.Companies) > 0 {
    profile, err := parse.CompanyID(ctx, page.Companies[0].ID, parseapi.CompanyIDOptions{Deep: true})
    if err != nil { return err }
    _ = profile.Deep // *CompanyProfileDeep; fields may be nil
}
coverage, err := parse.CompanyCoverage(ctx)
_ = coverage
_ = err

// Continue a name search with the same selector, filters and limit.
first, err := parse.CompanySearch(ctx, parseapi.CompanySearchOptions{Query: "Example", Country: "US", Limit: 10})
if err != nil { return err }
if first.Next != nil {
    next, err := parse.CompanySearch(ctx, parseapi.CompanySearchOptions{Query: "Example", Country: "US", Limit: 10, Cursor: *first.Next})
    _ = next
    _ = err
}
```

Directory `deep` adds detail to each profile in the same pooled request on every plan. Omitted deep, empty deep and partial detail remain distinct; description, logo, socials, founded, employees, registrations and sources may be absent on older server releases. Sources attribute only their listed selected enrichment fields. Founding precision is preserved separately from incorporation. Missing listings do not establish private ownership, and a domain match does not prove legal identity. Coverage describes the returned edition, not every company worldwide.

Discovery example: `parse.CompanySearch(ctx, parseapi.CompanySearchOptions{Country: "US", Industry: "0700", IndustryType: "sic"})`

Supply `Industry` and `IndustryType` together. The supported namespace is `sic`, with an exact four-digit string such as `0700`; leading zeros are meaningful. Country-only discovery is also supported. Filters intersect and may narrow an existing selector. Country matches the profile country, not a headquarters or operating-presence claim. Unknown values do not match a requested filter. Filter-only candidates use `match.field: "filters"` and `match.value: null`; reuse the same filters and limit with a returned cursor. Counts describe this directory edition, not complete country coverage.

Reviewed `deep.registrations` retain the registry authority and exact number, registration jurisdiction, domestic role, legal form, administrative status and source-scoped formation date. Principal addresses keep their role and recorded text; they are not headquarters. Registration does not establish current operations or tax exemption. `[]` means no admitted registration facts; older responses may omit the field. Sources use `business_register` for these facts and preserve the original observation time; unknown record update times remain null.

## Time

`time` returns local ISO `at` with its UTC offset and integer Unix seconds in `unix`. Request `deep` for the display name, exact `offset_seconds`, whole `offset_minutes` and next clock change. A conversion target has its own optional `deep` without a next-change field. Historical offsets and ISO times can include offset seconds. Omitted `at` means now. With `to` or `targets`, an offsetless `at` is source wall time. Otherwise it is UTC. Include an offset for repeated local times around a clock change. Current time and conversion use pooled requests on every plan. Coordinate clock fields can be null when the timezone is unknown. Existing `timezone` methods remain supported.

For an offsetless `at` with `to` or `targets`, choose how to handle a clock change with `disambiguation`. It applies to named-zone and coordinate Time calls.

| Value | Repeated time | Skipped time |
| --- | --- | --- |
| `compatible` (default) | Earlier occurrence | Shift forward by the clock change |
| `earlier` | Earlier occurrence | Shift backward by the clock change |
| `later` | Later occurrence | Shift forward by the clock change |
| `reject` | `400 ambiguous_time` | `400 nonexistent_time` |

An explicit UTC offset selects an instant directly. For example, `2026-11-01T01:30:00-04:00` and `2026-11-01T01:30:00-05:00` identify the two New York occurrences. A valid `disambiguation` value has no effect on explicit instants, current-time requests or lookups without `to` or `targets`. For user-entered appointment times, start with `reject`. Handle `ambiguous_time` or `nonexistent_time` by collecting an explicit offset or an earlier/later choice from the user. Other malformed input still uses `invalid_request`.

```go
result, err := parse.Time(ctx, "America/New_York", parseapi.TimeOptions{
    At: "2026-11-01T01:30:00", To: "UTC", Disambiguation: "later",
})
if err != nil { return err }
fmt.Println(result.To.At) // 2026-11-01T06:30:00+00:00
```

Canonical Time `deep` includes the pinned rule edition in `Deep.TimezoneDatabaseVersion` and source-wall resolution in `Deep.Resolution`. Resolution records `kind` (`unique`, `overlap` or `gap`), the selected policy, signed `adjustment_seconds`, and chronological alternatives with exact `at`, Unix seconds and UTC offset. Unique times have an empty alternatives list. Explicit instants, current time and lookups without conversion have null resolution. Destination detail stays compact.

Search serving IANA IDs by city or region, or omit the query to list all (Go and Rust use an empty string). Discovery returns `timezone_database_version` and sorted `timezones`. No search matches returns `timezones: []`.

Pass `targets` to convert one instant to 1-10 zones in a single pooled request. The native list preserves order and duplicates. Use `targets` instead of `to`. The response adds `targets`, with optional detail inside each target. Unknown source coordinates return `targets: null`. An unknown destination rejects the whole request with `not_found`. Omission keeps the original response shape.

```go
zones, err := parse.TimeZones(ctx, "New York")
if err != nil { return err }
fmt.Println(zones.Timezones)
result, err := parse.Time(ctx, "UTC", parseapi.TimeOptions{
    At: "2026-09-24T12:00:00Z", Targets: []string{"America/New_York", "Asia/Tokyo"},
})
if err != nil { return err }
for _, target := range result.Targets {
    fmt.Println(target.Timezone, target.At)
}
```

### Location inputs and timezone filters

`parse.Time(ctx, "", parseapi.TimeOptions{IATA: "JFK", Deep: true})` and `parse.TimeZones(ctx, "", parseapi.TimeZonesOptions{Country: "US", Details: true})`. `DST` and `ObservesDST` are `*bool` so an explicit false is distinct from an omitted filter.

Choose one explicit location input: IP, exact city name or stable city ID, country, IATA airport, ICAO airport, port UN/LOCODE, or address. Country and state can narrow a city or address. State requires country. Address lookup requires US country context and a strict address-point match. Port lookup covers the reviewed port subset, not every assigned UN/LOCODE. IP lookup always uses the supplied IP.

Location calls add `location` with `status`, `candidates`, `truncated`, `source` and the typed input. Check `status` before using the clock: ambiguous or missing locations retain null time fields. Candidate coordinates and IDs can also be null. A country with multiple timezones does not silently choose one. Named-zone and coordinate calls retain their existing signatures.

Timezone discovery accepts country, IANA area, exact signed offset, abbreviation, DST-at-instant and observes-DST-during-year filters. `at` selects the common instant, `sort` selects timezone or offset order, and `details` adds `zones` rows plus the evaluation `at`. The default `timezones` list stays compact. False DST filters are sent explicitly. An abbreviation returns candidate zones rather than choosing one. Observes-DST uses the UTC calendar year containing `at`.

Source deep adds `standard_offset`, `standard_offset_seconds`, signed `dst_offset_seconds` and `season`. Seasonal adjustments can be negative. `season` describes the current DST-flag interval, or the next within 400 days, with actual before/after transition facts and signed `change_seconds`. Unknown boundaries remain null. These fields are optional and nullable, and destination deep stays compact.

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

Australian postal lookups include core `localities` with suburb choices (`city`, `state`, `state_name`) on every plan. Null or an omitted field means unknown, while `[]` means the reviewed reference has no eligible choices. `city` stays null when the source is ambiguous, even if there is only one eligible choice. Let the user select their suburb and keep manual entry available. These are geographic choices, not mailing-address verification. [G-NAF source, adaptations and licence](https://parseapi.com/legal/attribution#postal-au).

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
Tariff lookup and search accept an optional `edition` fingerprint and `date` (`YYYY-MM-DD`). The edition pins exact immutable source bytes. A date is accepted only when verified source coverage exists. An edition without a date returns undated schedule context (`date: null`). Default requests use today. Paid detail exposes an open-string `reason` when `effective_rate` is null, including `incomplete_coverage`. A null rate never means zero. Explicit selections fail with `tariff_selection_mismatch` if an older server ignores the requested scope.


Origin means where the goods originate, not where they ship from. The effective rate covers matched stored schedule measures only; it is not complete duty or landed cost.

Codes contain 4, 6, 8 or 10 ASCII digits; dots and whitespace are optional. Search returns up to 20 description matches with parent `lineage` so a result named “Other” has context. Search is not product classification. In deep, `measures: null` means origin-dependent measures were not resolved; `measures: []` means the resolved lookup found none.

```go
parse.Tariff(ctx, "8471.30.01.00")
parse.Tariff(ctx, "8471.30.01.00", parseapi.TariffOptions{Deep: true})
parse.Tariff(ctx, "8471.30.01.00", parseapi.TariffOptions{Deep: true, Origin: "CN"})
```

Address search uses context from the form: prefer postal, or city and state. An optional end-user `ip` is a locality hint for server-side calls. An empty result explains itself with `reason`: `more_input`, `missing_context` or `no_matches`. With suggestions, reason is null. Older responses may omit it, and future reasons remain strings. Catalog and lookup failures use the existing API errors.

HLR reports status at the last check. `live` means assigned and `connected` means reachable at that check. Cached results may be returned. Null means unconfirmed. Deep diagnostics stay within the same metered lookup.

Bank returns core `checks` for input, country, length, structure, checksum and national rules, plus an `issues` list. States are `passed`, `failed`, `not_checked` or `not_supported`. Unsupported national checking is not a failure. `valid` covers the implemented format and checksum rules, not account existence, ownership or payment reachability. Directory names and BICs may be null independently. Older responses may omit `checks` and `issues`, and future states and issue codes remain strings. Pass the original input unchanged so the API can report invalid characters. Deep `account` remains the BBAN remainder.

Bank inputs use `POST /bank` JSON bodies, keeping IBAN and account values out of request URLs. Pass original strings; the server owns normalization and validation. Avoid logging request bodies. IBAN deep can include `directory` with the immutable `edition`, resolved `country` and actual `match` grain (`bank`, `branch`, `prefix` or `none`); it is absent if no directory lookup ran. A match does not prove complete country coverage or payment reachability.

Use country requirements to build supported input fields. US ACH has an explicit helper with no deep option. It checks the routing format/ABA checksum and account-field syntax; `account_checksum` is `not_supported`. It preserves account characters and leading zeros. A nullable bank name is routing-directory identity, not account existence, ownership or ACH eligibility. The examples below are synthetic test inputs, not payment instructions.

```go
parse.BankRequirements(ctx, "US", parseapi.BankRequirementsOptions{Format: "us_ach"})
parse.BankUSACH(ctx, parseapi.BankUSACHInput{Routing: "011000015", Account: "0001234567"})
```

## Provider lookup

```go
provider, err := parse.Provider(ctx, "1881018208")
profile, err := parse.Provider(ctx, "1881018208", parseapi.ProviderOptions{Deep: true})
```

Pass the original NPI as a string. `valid` checks its format and checksum; `registered` means a match in the stored NPPES snapshot. `active` reflects recorded NPI deactivation, not licensure. `excluded` is an NPI-only OIG LEIE match; `false` is not a complete exclusion clearance. These directory facts do not verify credentials, current practice contact or payment eligibility.

Invalid input returns `valid: false` with unknown provider fields. A checksum-valid number missing from the snapshot returns `registered: false`; unavailable storage remains an API error. Preserve `null` as unknown.

The default pooled lookup includes provider identity, specialty and practice contact where held. Paid `deep` adds `deactivated_at`, `medicare`, `opt_out` and `enrollments` from stored source files, with no separate check meter or live verification. `enrollments: null` means unavailable; `[]` means no enrollment rows are returned. The API omits unrequested `deep` and returns `{}` when requested on Free.

Paid Deep also returns `taxonomies` in published order, with taxonomy code, specialty label, primary flag and provider-reported license number/state, plus `enumerated_at`, `updated_at` and `reactivated_at` record dates. Reported licenses are not verified licenses. Null lists mean unavailable; empty lists mean the edition contains no entries. Core `sources` is available on every plan: NPPES, LEIE, PECOS and opt-out each have nullable edition metadata (`edition`, `published_at`, `through`, `imported_at`). Provider record dates are separate from source publication and completed import dates. Older responses may omit these additions. Edition details remain null until a verified source is served.

## Deep

The default call returns the common answer. Request more detail with `parse.Country(ctx, "US", parseapi.CountryOptions{Deep: true})`. Read those fields from the optional deep member; this does not change the core answer.

| Operation | What `deep` requests |
|---|---|
| NPI | All published taxonomies, reported license details, provider record dates and Medicare detail on paid plans. Primary specialty, exclusion flag and source metadata stay core. |
| IP | Richer IP fields included with a paid plan. No separate check meter. |
| Domain | Registration dates, registrar, status and DNSSEC, included with a paid plan. Use `DNS` for DNS records and `MX` for mail routing. |
| Email | A metered mailbox check with deliverability, catch-all, status, reason and address hints, using included email checks or enabled on-demand usage. |
| VAT | A metered registry check where supported, using included VAT checks or enabled on-demand usage. |
| Country, State, City, District, Postal | Reference profiles included with a paid plan; place identity and coordinates stay core. |
| VIN, Industry, Company (national number) | Paid technical or registration profiles. Industry hierarchy stays core. |
| Tariff | Paid schedule columns and units; add origin for applicable measures. |
| Name, Weather | Paid name context or weather detail; parsing and current conditions stay core. |
| Phone, Bank | Numbering-plan or bank structure detail in the same pooled request on every plan. |
| Time, Date, Currency, Language, Emoji, Point | Optional reference detail in the same pooled request on every plan. |
| Carrier, HLR | Available place or network detail from the same metered core unit, including Free included units. |

Email deep includes mailbox status and the reason for the result, plus a suggested first name, no-reply flag, plus-address tag and mail service. The suggested name is not a verified identity. Unavailable details are null.

Reasons include `accepted`, `invalid_format`, `invalid_domain`, `no_mail_server`, `mailbox_not_found`, `mailbox_disabled`, `mailbox_full`, `catchall`, `disposable`, `temporary_failure`, `rejected` and `unconfirmed`.

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

Create one client and share it across goroutines. A context deadline covers the whole call, including retry waits. The client timeout defaults to 35 seconds for Stack and 10 seconds for other operations. Cancellation stops the request and any retry wait.

Ordinary lookups retry twice on network failures, 429, 500, 502, 503, and 504. Carrier, caller, and HLR calls use one attempt by default. Deep email, VAT, and address calls also use one attempt, reserving that behavior for address verification. Address deep currently returns an empty object. An explicit retry setting applies to every call, including metered operations. Additional attempts can be billed.

```go
parse, err := parseapi.New("your-api-key",
    parseapi.WithTimeout(5*time.Second),
    parseapi.WithRetries(0),
)
```

`WithRetries(0)` disables all automatic retries. Numeric and HTTP-date `Retry-After` values within five seconds are honored. Longer waits return the original API error immediately. `RetryAfter` carries the original header, or nil when absent. Built-in requests do not follow redirects. `WithHTTPClient` copies your client and keeps redirects disabled. Custom transports must also keep credentials on the requested origin.

Requires Go 1.21 or later. Standard library only.

[Full endpoint and field reference](https://parseapi.com/docs)

## Card

Send 2–11 leading digits as a string. Core returns `bin`, `brand`, `brand_name`
and a CDN SVG `logo`. Brand detection uses reviewed network rules independently
of issuer records. Unknown or ambiguous prefixes return null brand fields and a
generic logo; a known network without reviewed artwork also uses the generic logo.

Optional Deep adds `prefix`, `issuer`, `country`, `type` and `prepaid`, included
in the same pooled request on every plan. Six or more digits enable directory
matching. Fewer digits return all-null Deep fields. Compare `deep.prefix` with
`bin`: equal is an exact recorded match; shorter is broader; null is no match.
The longest row wins, including null fields. `prepaid: null` means unknown, not
false. This is partial reference data, not card validity or payment acceptance.

```go
card, err := parse.Card(ctx, "51")
if err != nil { return err }
fmt.Println(card.Logo)
details, err := parse.CardWithOptions(ctx, "43737400", parseapi.CardOptions{Deep: true})
if err != nil { return err }
if details.Deep != nil { fmt.Println(details.Deep.Prefix) }
```

Leading zeros are preserved. Only ASCII spaces, tabs, CR, LF and hyphens are
removed; raw input is limited to 64 characters. Invalid prefixes are rejected
before dispatch, accepted input is forwarded unchanged. Never send a full card number.

## Stack

```go
result, err := parse.Stack(ctx, "example.com")
if err != nil {
    return err
}
```

Pass a public hostname without a scheme, path, port or IP address. Stack returns the homepage URL and `checked_at` time, then eight technology arrays: `cms`, `servers`, `frameworks`, `ecommerce`, `analytics`, `chat`, `payments` and `hosting`. Each entry contains a `technology` code, name and nullable version. Multiple CMSs or servers remain separate entries. Empty arrays mean no matches in the checked pages. An unsuccessful check returns null arrays and a null `checked_at`.

`scope` identifies `homepage` or `site` coverage. `pages` counts successfully checked HTML pages. `partial` is true for a homepage-only or incomplete bounded site check, false when the known in-scope candidates finished, and null when no check succeeded. False does not guarantee that every page on the website was discovered.

The complete technology result is included in the core response. The generic `deep=true` option adds only an empty object and is unnecessary for Stack. Successful checks may be reused for up to 24 hours. `pretty` optionally formats the wire JSON. Each lookup uses one request and API version 2.0.0 selected by this client.

Stack defaults to 35 seconds per attempt so a first scan has time to finish. Other lookups retain their 10-second default. An explicit client timeout takes precedence.

Vehicle lookups use `vin` as the input and response field. Existing VIN methods remain available for compatibility.
