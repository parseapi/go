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

This SDK sends `Parse-Version: 2.0.0` on every request, including retries. Its response types match API `2.0.0`, and the client selects that contract automatically. No extra constructor setting or key change is needed. This behavior requires the matching API request-version release.

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

Directory `deep` adds detail to each profile in the same pooled request on every plan. Omitted deep, empty deep and partial detail remain distinct; description, logo, social_profiles, phone_numbers, email_addresses, domains, founded, employees, registrations and sources may be absent on older server releases. Sources attribute only their listed selected enrichment fields. Founding precision is preserved separately from incorporation. Missing listings do not establish private ownership, and a domain match does not prove legal identity. Coverage describes the returned edition, not every company worldwide.

Company directory social profiles contain `platform`, `url` and `handle`; unknown metadata stays null. Phone and email records use `type` for the source-reported purpose. Earlier response fields remain readable by the client.

Discovery example: `parse.CompanySearch(ctx, parseapi.CompanySearchOptions{Country: "US", Industry: "0700", IndustryType: "sic"})`

Supply `Industry` and `IndustryType` together. The supported namespace is `sic`, with an exact four-digit string such as `0700`; leading zeros are meaningful. Country-only discovery is also supported. Filters intersect and may narrow an existing selector. Country matches the profile country, not a headquarters or operating-presence claim. Unknown values do not match a requested filter. Filter-only candidates use `match.field: "filters"` and `match.value: null`; reuse the same filters and limit with a returned cursor. Counts describe this directory edition, not complete country coverage.


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

Successful HTTP JSON contains `measure`, `valid`, `type`, `amount`, and `unit`. Invalid input adds `reason`. `choices` appears only when unit alternatives are available.

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

Read `population` as the source’s count for that place boundary. Keep missing or null values unknown and preserve a known zero. State district lists include each district’s population. Postal nearby and distance detail remains metropolitan associations only. Continent population stays in core.

Point returns the timezone ID with the core location. Its optional deep detail adds terrain and compact nearest-city context on every plan. A nearest city is null when none is within 200 km.

Weather returns current conditions by default. Paid deep adds specialist current measurements, forecasts and related detail. A past `date` is a UTC day and requires deep: it adds `deep.history` alongside current conditions. Date alone does not request history.

```go
parse.Weather(ctx, 40.7128, -74.006, parseapi.WeatherOptions{Deep: true, Date: "2026-08-15"})
```

Tariff starts with the general schedule line. Paid deep adds units and the special and other schedule columns. An optional origin then resolves country-specific measures. The three calls below show those successive choices. Without origin, schedule detail is still returned and origin-dependent fields are null. A null effective rate is not a zero rate.

Tariff lookup and search accept the `edition` ID returned by Tariff, such as `ted_3`, and an optional `date` (`YYYY-MM-DD`). The edition pins exact immutable source bytes. A date is accepted only when verified source coverage exists. An edition without a date returns undated schedule context (`date: null`). Default requests use today. Paid detail exposes an open-string `reason` when `effective_rate` is null, including `incomplete_coverage`. A null rate never means zero. Explicit selections fail with `tariff_selection_mismatch` if an older server ignores the requested scope.

Origin means where the goods originate, not where they ship from. The effective rate covers matched stored schedule measures only. It is not complete duty or landed cost.

Codes contain 4, 6, 8 or 10 ASCII digits; dots and whitespace are optional. Search returns up to 20 description matches with parent `lineage` so a result named "Other" has context. Search is not product classification. In deep, `measures: null` means origin-dependent measures were not resolved. `measures: []` means the resolved lookup found none.

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

Paid Deep also returns `taxonomies` in published order, with taxonomy code, specialty label, primary flag and provider-reported license number/state, plus `enumerated_at`, `updated_at` and `reactivated_at` record dates. Reported licenses are not verified licenses. Null lists mean unavailable; empty lists mean the edition contains no entries. Provider record dates describe the provider record. Unknown dates stay null.

## Deep

The default call returns the common answer. Request more detail with `parse.Country(ctx, "US", parseapi.CountryOptions{Deep: true})`. Read those fields from the optional deep member; this does not change the core answer.

| Operation | What `deep` requests |
|---|---|
| NPI | All published taxonomies, reported license details, provider record dates and Medicare detail on paid plans. Primary specialty and exclusion flag stay core. |
| IP | Richer IP fields included with a paid plan. No separate check meter. |
| Domain | Registration dates, registrar, status and DNSSEC, included with a paid plan. Use `DNS` for DNS records and `MX` for mail routing. |
| Email | A metered mailbox check with deliverability, catch-all, status, reason and address hints, using included email checks or enabled on-demand usage. |
| VAT | A metered registry check where supported, using included VAT checks or enabled on-demand usage. |
| Country, State, City, District, Postal | Reference profiles included with a paid plan; place identity and coordinates stay core. |
| Industry, Company (national number) | Paid technical or registration profiles. Industry hierarchy stays core. |
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

## Emoji catalog additions

Resolve Unicode emoji or a named shortcode, search or browse a bounded catalog,
parse emoji positions in text, and convert between Unicode and a selected shortcode source.
The ordinary lookup remains four fields:

```http
GET /emoji/rocket
```

```json
{"emoji":"🚀","name":"Rocket","shortcodes":[":rocket:","rocket"],"category":"Travel & Places"}
```

Use `deep=true` for Unicode facts, keywords, tone families and available artwork.
Artwork is keyed by style: `deep.artwork.openmoji.svg` is the SVG URL, and
`deep.artwork.openmoji.png["72"]` is the 72×72 PNG URL (use your language’s map
access syntax). The same shape applies to `noto` and `twemoji`. Missing styles
and formats are omitted. URLs pin immutable artwork editions and retain their bytes.
Core values stay unchanged. Text parsing returns `emojis` with UTF-16 `positions`;
conversion returns `text` and unresolved tokens with positions and reasons.

The compact JSON catalog contains offline identity, qualification-form and
shortcode indexes. Request `profile=full` when artwork is needed and `lang` when
localized names and keywords are needed. `/emoji/metadata` carries dataset
coverage and shared source notices once. Artwork supports Google Noto, Twemoji
and OpenMoji, with real SVG/PNG availability and immutable first-party URLs.
A device's emoji font remains the ordinary display.

These additions are a local release candidate until the matching SDK and
backend are published. Existing lookup and search contracts are preserved.


## City collections

```go
minimum := int64(100000)
page, err := parse.CityList(ctx, parseapi.CityListOptions{
    Country: "JP", MinPopulation: &minimum, Sort: "population", Limit: 20,
})
nearby, err := parse.CityNearbyCoords(ctx, 34.57, 135.48,
    parseapi.CityNearbyCoordsOptions{Radius: 20, Unit: "km", Country: "JP"})
// Existing name search supports opt-in paging.
search, err := parse.CitySearch(ctx, "Sak", parseapi.CitySearchOptions{Page: true})
```

Browse the curated City catalog or list cities around a supplied coordinate. These operations return `next_cursor`, which is null at the end. Treat it as opaque and reuse all filters, limit, language, hierarchy and deep settings for the next page. Existing name searches and named nearby calls include pagination metadata only when paging or a cursor is requested.

Collection options include minimum and maximum population, an exact IANA timezone, and `bbox` in `west,south,east,north` order. A west bound greater than east crosses the dateline. Population filtering and sorting use the same count as paid Deep: an admitted observation when available, otherwise the stored count. Missing counts do not match a population range and sort last; legacy counts keep unknown dates. Coverage and observation dates vary by place.

Collections support `name`, `population`, `country` and `elevation`; search also supports `relevance`, and nearby supports `distance`. Combine up to three signed terms, for example `sort=-population,name`. A signed multi-key sort does not also take `order`. Population defaults to descending order; other sorts default to ascending. Set `order` to `asc` or `desc`, except relevance only supports ascending. Named nearby uses country and state to select the anchor; population, timezone and bounding box filter the returned neighbors. Coordinate nearby applies country and state to its returned cities.



## Canonical City operations

Classic lookup, search, nearest and named-nearby methods retain their default envelopes and required-field models. Browse and coordinate-nearby use canonical lists: `object`, `data` and nullable `next_cursor` (camelCase properties in Swift/Kotlin). Keep all query options unchanged when continuing a page; no client fetches pages automatically.

```go
page, err := parse.CitySearchList(ctx, "Viena", parseapi.CitySearchListOptions{
    Match: "fuzzy", MaxDistance: 1, })
children, err := parse.CityChildren(ctx, "city_7aj899nvwyku")
distance, err := parse.CityDistance(ctx, "city_7aj899nvwyku", "city_7aj899nvwyku")
```

Explicit resource and search-list/nearby-list methods return complete City resources with `id` and `object`. Existing classic methods keep their published return types. Country inclusion/exclusion sets, timezone sets, inclusive population and meter-elevation bounds, stored-center `bbox` and signed multi-key sorting are available for collections. Unknown numeric facts do not match range filters.

Search defaults to prefix matching. Exact and fuzzy matching are explicit; fuzzy accepts maximum edit distance 1 or 2 and returns `match` evidence (`type`, `source`, `value`, `language`, `edit_distance`). Review a fuzzy choice before storing its ID.

Reviewed children, ancestors and siblings use canonical lists with partial relationship coverage. Empty data does not prove there are no geographic relatives. Children/siblings accept page controls; ancestors do not. Distance reports great-circle kilometers and miles between stored centers: zero is known zero, null is unknown, and it does not imply travel time or containment.

Strict `ambiguity=error` can return HTTP 409 `ambiguous_city` with up to two safe candidate identities. The existing error constructor and unrelated error metadata are preserved; candidates are retained only for that status/code and never contain Deep. GeoJSON is available through HTTP; these typed City list methods do not return GeoJSON.


Reviewed name forms:

```go
city, err := parse.CityIDResource(ctx, "city_7aj899nvwyku", parseapi.CityIDResourceOptions{Provenance: true, Romanized: true})
```

Coverage by country and available fields is described at https://parseapi.com/city#coverage.

`romanized=true` requests reviewed official name forms as `name_romanized` and `romanization_language` (camelCase in Swift/Kotlin). Current reviewed forms cover 20 Japanese whole-city aliases in ja-Latn. Missing or conflicting forms stay null; ASCII folding is not substituted for an official form.


Conditional `If-None-Match` reads and GeoJSON currently use HTTP directly. A matching City ETag returns 304 without a body; ordinary billing and rate limits apply. These JSON client methods do not supply conditional-cache handling.
