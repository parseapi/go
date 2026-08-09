// Package parseapi is the official parseAPI client for Go.
// One key, minimal JSON, fast. https://parseapi.com
package parseapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	version         = "0.1.0"
	defaultBaseURL  = "https://api.parseapi.com"
	defaultTimeout  = 10 * time.Second
	defaultRetries  = 2
	retryAfterCapMs = 5000
)

var retryStatus = map[int]bool{429: true, 500: true, 502: true, 503: true, 504: true}

// Error is every non-2xx response from the API. Branch on Code, never on Message.
type Error struct {
	Status    int
	Code      string
	Message   string
	Docs      string
	RequestID string
}

func (e *Error) Error() string {
	return fmt.Sprintf("parseapi: %s (%s)", e.Message, e.Code)
}

// Client is a parseAPI client. Create one with New and share it.
type Client struct {
	apiKey     string
	baseURL    string
	retries    int
	httpClient *http.Client
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides https://api.parseapi.com (tests, canaries).
func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.baseURL = strings.TrimRight(baseURL, "/") }
}

// WithTimeout sets the per-attempt timeout. Default 10s.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) { c.httpClient.Timeout = timeout }
}

// WithRetries sets retries after the first attempt on network errors / 429 / 5xx.
// Default 2, 0 disables.
func WithRetries(retries int) Option {
	return func(c *Client) { c.retries = retries }
}

// WithHTTPClient replaces the underlying http.Client (instrumentation, proxies).
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) { c.httpClient = httpClient }
}

// New creates a Client. An empty apiKey falls back to the PARSEAPI_KEY env var.
func New(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		apiKey = os.Getenv("PARSEAPI_KEY")
	}
	if apiKey == "" {
		return nil, errors.New("parseapi: missing API key, pass one or set PARSEAPI_KEY")
	}
	baseURL := os.Getenv("PARSEAPI_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	c := &Client{
		apiKey:     apiKey,
		baseURL:    strings.TrimRight(baseURL, "/"),
		retries:    defaultRetries,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

func retryDelay(attempt int, retryAfter string) time.Duration {
	if retryAfter != "" {
		if seconds, err := strconv.ParseFloat(retryAfter, 64); err == nil && seconds >= 0 {
			return time.Duration(math.Min(seconds*1000, retryAfterCapMs)) * time.Millisecond
		}
	}
	return time.Duration(rand.Float64()*250*math.Pow(2, float64(attempt))) * time.Millisecond
}

func (c *Client) get(ctx context.Context, path string, query url.Values, headers map[string]string, out any) error {
	target := c.baseURL + path
	if len(query) > 0 {
		target += "?" + query.Encode()
	}

	for attempt := 0; ; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
		if err != nil {
			return err
		}
		req.Header.Set("X-API-Key", c.apiKey)
		req.Header.Set("User-Agent", "parseapi-go/"+version)
		for name, value := range headers {
			req.Header.Set(name, value)
		}

		res, err := c.httpClient.Do(req)
		if err != nil {
			if attempt < c.retries && ctx.Err() == nil {
				time.Sleep(retryDelay(attempt, ""))
				continue
			}
			return err
		}

		if res.StatusCode >= 200 && res.StatusCode < 300 {
			defer res.Body.Close()
			return json.NewDecoder(res.Body).Decode(out)
		}

		if retryStatus[res.StatusCode] && attempt < c.retries {
			retryAfter := res.Header.Get("Retry-After")
			res.Body.Close()
			time.Sleep(retryDelay(attempt, retryAfter))
			continue
		}

		apiErr := &Error{
			Status:  res.StatusCode,
			Code:    "unknown_error",
			Message: fmt.Sprintf("Request failed with status %d", res.StatusCode),
		}
		var body struct {
			Code      string `json:"code"`
			Message   string `json:"message"`
			Docs      string `json:"docs"`
			RequestID string `json:"request_id"`
		}
		if json.NewDecoder(res.Body).Decode(&body) == nil {
			if body.Code != "" {
				apiErr.Code = body.Code
			}
			if body.Message != "" {
				apiErr.Message = body.Message
			}
			apiErr.Docs = body.Docs
			apiErr.RequestID = body.RequestID
		}
		res.Body.Close()
		return apiErr
	}
}

func seg(value string) string { return url.PathEscape(value) }

func values(pairs ...string) url.Values {
	query := url.Values{}
	for i := 0; i+1 < len(pairs); i += 2 {
		if pairs[i+1] != "" {
			query.Set(pairs[i], pairs[i+1])
		}
	}
	return query
}

func f(value float64) string { return strconv.FormatFloat(value, 'f', -1, 64) }

// DeepOptions requests the nested deep object. Paid on most endpoints.
type DeepOptions struct {
	Deep bool
}

func deepValue(opts *DeepOptions) string {
	if opts != nil && opts.Deep {
		return "true"
	}
	return ""
}

// IP looks up an IP address.
func (c *Client) IP(ctx context.Context, ip string, opts *DeepOptions) (*IP, error) {
	out := &IP{}
	return out, c.get(ctx, "/ip/"+seg(ip), values("deep", deepValue(opts)), nil, out)
}

// IPMe looks up the caller's IP.
func (c *Client) IPMe(ctx context.Context, opts *DeepOptions) (*IP, error) {
	out := &IP{}
	return out, c.get(ctx, "/ip/me", values("deep", deepValue(opts)), nil, out)
}

// Continent looks up a continent by code (NA, EU, ...).
func (c *Client) Continent(ctx context.Context, code string) (*Continent, error) {
	out := &Continent{}
	return out, c.get(ctx, "/continent/"+seg(code), nil, nil, out)
}

// ContinentCountries lists countries in a continent.
func (c *Client) ContinentCountries(ctx context.Context, code string) (*ContinentCountries, error) {
	out := &ContinentCountries{}
	return out, c.get(ctx, "/continent/"+seg(code)+"/countries", nil, nil, out)
}

// Country looks up a country by ISO code.
func (c *Client) Country(ctx context.Context, code string) (*Country, error) {
	out := &Country{}
	return out, c.get(ctx, "/country/"+seg(code), nil, nil, out)
}

// CountryStates lists states in a country.
func (c *Client) CountryStates(ctx context.Context, code string) (*CountryStates, error) {
	out := &CountryStates{}
	return out, c.get(ctx, "/country/"+seg(code)+"/states", nil, nil, out)
}

// State looks up a state or province. Country is required because state codes
// are not globally unique.
func (c *Client) State(ctx context.Context, code, country string) (*State, error) {
	out := &State{}
	return out, c.get(ctx, "/state/"+seg(code), values("country", country), nil, out)
}

// StateDistricts lists districts under a state.
func (c *Client) StateDistricts(ctx context.Context, code, country string) (*StateDistricts, error) {
	out := &StateDistricts{}
	return out, c.get(ctx, "/state/"+seg(code)+"/districts", values("country", country), nil, out)
}

// DistrictOptions narrows a district lookup.
type DistrictOptions struct {
	Country string
}

// District looks up a district (ADM2) by code.
func (c *Client) District(ctx context.Context, code string, opts *DistrictOptions) (*District, error) {
	country := ""
	if opts != nil {
		country = opts.Country
	}
	out := &District{}
	return out, c.get(ctx, "/district/"+seg(code), values("country", country), nil, out)
}

// CityOptions narrows a city lookup.
type CityOptions struct {
	Country string
	State   string
}

// City looks up a city by name.
func (c *Client) City(ctx context.Context, name string, opts *CityOptions) (*City, error) {
	country, state := "", ""
	if opts != nil {
		country, state = opts.Country, opts.State
	}
	out := &City{}
	return out, c.get(ctx, "/city/"+seg(name), values("country", country, "state", state), nil, out)
}

// CityID fetches a city by its minted parse id (city_…).
func (c *Client) CityID(ctx context.Context, id string) (*City, error) {
	out := &City{}
	return out, c.get(ctx, "/city/id/"+seg(id), nil, nil, out)
}

// CitySearchOptions narrows a city search.
type CitySearchOptions struct {
	Country string
	State   string
	Limit   int
}

// CitySearch searches cities by name prefix.
func (c *Client) CitySearch(ctx context.Context, q string, opts *CitySearchOptions) (*CitySearch, error) {
	country, state, limit := "", "", ""
	if opts != nil {
		country, state = opts.Country, opts.State
		if opts.Limit > 0 {
			limit = strconv.Itoa(opts.Limit)
		}
	}
	out := &CitySearch{}
	return out, c.get(ctx, "/city", values("q", q, "country", country, "state", state, "limit", limit), nil, out)
}

// CityNearest finds the nearest city to a point.
func (c *Client) CityNearest(ctx context.Context, lat, lon float64) (*CityNearest, error) {
	out := &CityNearest{}
	return out, c.get(ctx, "/city", values("lat", f(lat), "lon", f(lon)), nil, out)
}

// Postal looks up a postal or ZIP code. Country is required.
func (c *Client) Postal(ctx context.Context, code, country string) (*Postal, error) {
	out := &Postal{}
	return out, c.get(ctx, "/postal/"+seg(code), values("country", country), nil, out)
}

// PostalNearbyOptions tunes a nearby search. Radius in the unit ("km" default, "mi").
type PostalNearbyOptions struct {
	Radius float64
	Unit   string
}

// PostalNearby lists postal codes near one.
func (c *Client) PostalNearby(ctx context.Context, code, country string, opts *PostalNearbyOptions) (*PostalNearby, error) {
	radius, unit := "", ""
	if opts != nil {
		if opts.Radius > 0 {
			radius = f(opts.Radius)
		}
		unit = opts.Unit
	}
	out := &PostalNearby{}
	return out, c.get(ctx, "/postal/"+seg(code)+"/nearby", values("country", country, "radius", radius, "unit", unit), nil, out)
}

// PostalDistance measures the distance between two postal codes.
func (c *Client) PostalDistance(ctx context.Context, from, to, country string) (*PostalDistance, error) {
	out := &PostalDistance{}
	return out, c.get(ctx, "/postal/"+seg(from)+"/distance/"+seg(to), values("country", country), nil, out)
}

// Email validates an email address.
func (c *Client) Email(ctx context.Context, email string, opts *DeepOptions) (*Email, error) {
	out := &Email{}
	return out, c.get(ctx, "/email/"+seg(email), values("deep", deepValue(opts)), nil, out)
}

// PhoneOptions narrows a phone lookup. Country is the default region for
// national formats without a leading plus.
type PhoneOptions struct {
	Country string
	Deep    bool
}

// Phone validates and formats a phone number.
func (c *Client) Phone(ctx context.Context, number string, opts *PhoneOptions) (*Phone, error) {
	country, deep := "", ""
	if opts != nil {
		country = opts.Country
		if opts.Deep {
			deep = "true"
		}
	}
	out := &Phone{}
	return out, c.get(ctx, "/phone/"+seg(number), values("country", country, "deep", deep), nil, out)
}

// Domain checks if a domain is available to register.
func (c *Client) Domain(ctx context.Context, domain string, opts *DeepOptions) (*Domain, error) {
	out := &Domain{}
	return out, c.get(ctx, "/domain/"+seg(domain), values("deep", deepValue(opts)), nil, out)
}

// MX returns MX records for a domain.
func (c *Client) MX(ctx context.Context, domain string) (*MX, error) {
	out := &MX{}
	return out, c.get(ctx, "/mx/"+seg(domain), nil, nil, out)
}

// Useragent parses a User-Agent string.
func (c *Client) Useragent(ctx context.Context, ua string, opts *DeepOptions) (*Useragent, error) {
	out := &Useragent{}
	return out, c.get(ctx, "/useragent", values("deep", deepValue(opts)), map[string]string{"User-Agent": ua}, out)
}

// Currency looks up a currency by ISO 4217 code.
func (c *Client) Currency(ctx context.Context, code string) (*Currency, error) {
	out := &Currency{}
	return out, c.get(ctx, "/currency/"+seg(code), nil, nil, out)
}

// Language looks up a language by BCP 47 shortest code or ISO 639-3.
func (c *Client) Language(ctx context.Context, code string) (*Language, error) {
	out := &Language{}
	return out, c.get(ctx, "/language/"+seg(code), nil, nil, out)
}

// CurrencyRate returns the daily official reference rate for a currency pair.
func (c *Client) CurrencyRate(ctx context.Context, base, quote string) (*CurrencyRate, error) {
	out := &CurrencyRate{}
	return out, c.get(ctx, "/currency/"+seg(base)+"/"+seg(quote), nil, nil, out)
}

// TimezoneOptions evaluates the zone at an optional ISO-8601 instant.
type TimezoneOptions struct {
	At string
}

// Timezone looks up an IANA timezone.
func (c *Client) Timezone(ctx context.Context, id string, opts *TimezoneOptions) (*Timezone, error) {
	at := ""
	if opts != nil {
		at = opts.At
	}
	out := &Timezone{}
	return out, c.get(ctx, "/timezone/"+seg(id), values("at", at), nil, out)
}

// HolidayOptions selects a year. Zero means the current UTC year.
type HolidayOptions struct {
	Year int
}

// Holiday lists public holidays for a country and year.
func (c *Client) Holiday(ctx context.Context, country string, opts *HolidayOptions) (*HolidayYear, error) {
	year := ""
	if opts != nil && opts.Year > 0 {
		year = strconv.Itoa(opts.Year)
	}
	out := &HolidayYear{}
	return out, c.get(ctx, "/holiday/"+seg(country), values("year", year), nil, out)
}

// HolidayDate checks one date (YYYY-MM-DD). Holiday is nil when the date is
// not a holiday.
func (c *Client) HolidayDate(ctx context.Context, country, date string) (*HolidayDate, error) {
	out := &HolidayDate{}
	return out, c.get(ctx, "/holiday/"+seg(country)+"/"+seg(date), nil, nil, out)
}

// Elevation returns the elevation at a point.
func (c *Client) Elevation(ctx context.Context, lat, lon float64) (*Elevation, error) {
	out := &Elevation{}
	return out, c.get(ctx, "/elevation", values("lat", f(lat), "lon", f(lon)), nil, out)
}

// Point returns everything at a point: elevation plus the admin place.
func (c *Client) Point(ctx context.Context, lat, lon float64, opts *DeepOptions) (*Point, error) {
	out := &Point{}
	return out, c.get(ctx, "/point", values("lat", f(lat), "lon", f(lon), "deep", deepValue(opts)), nil, out)
}

// WeatherOptions tunes a weather lookup. Unit "imperial" flips wind, pressure,
// and distance.
type WeatherOptions struct {
	Unit string
	Deep bool
}

// Weather returns current conditions at a point from the nearest official station.
func (c *Client) Weather(ctx context.Context, lat, lon float64, opts *WeatherOptions) (*Weather, error) {
	unit, deep := "", ""
	if opts != nil {
		unit = opts.Unit
		if opts.Deep {
			deep = "true"
		}
	}
	out := &Weather{}
	return out, c.get(ctx, "/weather", values("lat", f(lat), "lon", f(lon), "unit", unit, "deep", deep), nil, out)
}

// Emoji resolves an emoji by character, shortcode, or name.
func (c *Client) Emoji(ctx context.Context, query string) (*Emoji, error) {
	out := &Emoji{}
	return out, c.get(ctx, "/emoji/"+seg(query), nil, nil, out)
}

// EmojiSearchOptions caps the result count.
type EmojiSearchOptions struct {
	Limit int
}

// EmojiSearch searches emoji by keyword.
func (c *Client) EmojiSearch(ctx context.Context, q string, opts *EmojiSearchOptions) (*EmojiSearch, error) {
	limit := ""
	if opts != nil && opts.Limit > 0 {
		limit = strconv.Itoa(opts.Limit)
	}
	out := &EmojiSearch{}
	return out, c.get(ctx, "/emoji", values("q", q, "limit", limit), nil, out)
}
