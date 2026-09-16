// Package parseapi is the official ParseAPI client for Go.
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
	version         = "0.4.0"
	defaultBaseURL  = "https://api.parseapi.com"
	defaultTimeout  = 10 * time.Second
	defaultRetries  = 2
	retryAfterCapMs = 5000
)

var retryStatus = map[int]bool{429: true, 500: true, 502: true, 503: true, 504: true}

// Error is every non-2xx response from the API. Branch on Code, never on Message.
type Error struct {
	_         [0]func()
	Status    int
	Code      string
	Message   string
	Docs      string
	RequestID string
}

func (e *Error) Error() string {
	return fmt.Sprintf("parseapi: %s (%s)", e.Message, e.Code)
}

// Client is a ParseAPI client. Create one with New and share it.
type Client struct {
	apiKey     string
	baseURL    string
	retries    int
	retriesSet bool
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
	return func(c *Client) {
		if c.httpClient != nil {
			// Options must not mutate an HTTP client shared by the caller.
			configured := *c.httpClient
			configured.Timeout = timeout
			c.httpClient = &configured
		}
	}
}

// WithRetries overrides retries for every operation. Ordinary lookups default
// to two retries, while metered operations default to none. Additional attempts
// can be billed. Zero disables all automatic retries.
func WithRetries(retries int) Option {
	return func(c *Client) { c.retries = retries; c.retriesSet = true }
}

// WithHTTPClient replaces the underlying http.Client (instrumentation, proxies).
// The client is copied and redirects remain disabled so API keys stay on the requested origin.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient == nil {
			c.httpClient = nil
			return
		}
		configured := *httpClient
		configured.CheckRedirect = func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }
		c.httpClient = &configured
	}
}

// New creates a Client. An empty apiKey falls back to the PARSEAPI_KEY env var.
func New(apiKey string, opts ...Option) (*Client, error) {
	// You found Dev. https://parseapi.com/dev
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
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		retries: defaultRetries,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
	for _, opt := range opts {
		if opt == nil {
			return nil, errors.New("parseapi: client option must not be nil")
		}
		opt(c)
	}
	if c.httpClient == nil {
		return nil, errors.New("parseapi: HTTP client must not be nil")
	}
	if c.retries < 0 {
		return nil, errors.New("parseapi: retries must be zero or greater")
	}
	if c.httpClient.Timeout < 0 {
		return nil, errors.New("parseapi: timeout must be zero or greater")
	}
	return c, nil
}

func retryDelay(attempt int, retryAfter string) time.Duration {
	if retryAfter != "" {
		if seconds, err := strconv.ParseFloat(retryAfter, 64); err == nil && seconds >= 0 {
			return time.Duration(math.Min(seconds*1000, retryAfterCapMs)) * time.Millisecond
		}
		if at, err := http.ParseTime(retryAfter); err == nil {
			delay := time.Until(at)
			if delay < 0 {
				return 0
			}
			if delay > 5*time.Second {
				return 5 * time.Second
			}
			return delay
		}
	}
	return time.Duration(rand.Float64()*250*math.Pow(2, float64(attempt))) * time.Millisecond
}

func waitRetry(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return ctx.Err()
	}
}

// String describes the client without exposing its API key.
func (c Client) String() string {
	return fmt.Sprintf("parseapi.Client{apiKey:[REDACTED], retries:%d}", c.retries)
}

// GoString keeps API keys out of Go-syntax debug output.
func (c Client) GoString() string { return c.String() }

func meteredRequest(path string, query url.Values) bool {
	for _, product := range []string{"carrier", "caller", "hlr", "litigator", "reassigned"} {
		if strings.HasPrefix(path, "/"+product+"/") {
			return true
		}
	}
	return query.Get("deep") == "true" && (strings.HasPrefix(path, "/email/") || strings.HasPrefix(path, "/vat/") || strings.HasPrefix(path, "/address/"))
}

func (c *Client) get(ctx context.Context, path string, query url.Values, headers map[string]string, out any) error {
	retries := c.retries
	if !c.retriesSet && meteredRequest(path, query) {
		retries = 0
	}
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
			if attempt < retries && ctx.Err() == nil {
				if err := waitRetry(ctx, retryDelay(attempt, "")); err != nil {
					return err
				}
				continue
			}
			return err
		}

		if res.StatusCode >= 200 && res.StatusCode < 300 {
			defer res.Body.Close()
			return json.NewDecoder(res.Body).Decode(out)
		}

		if retryStatus[res.StatusCode] && attempt < retries {
			retryAfter := res.Header.Get("Retry-After")
			res.Body.Close()
			if err := waitRetry(ctx, retryDelay(attempt, retryAfter)); err != nil {
				return err
			}
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

// oneOption accepts zero or one option value. Multiple values are an error.
func oneOption[T any](options []T) (T, error) {
	var zero T
	if len(options) > 1 {
		return zero, errors.New("parseapi: pass at most one options value")
	}
	if len(options) == 1 {
		return options[0], nil
	}
	return zero, nil
}

// IPOptions configures IP. Omitted fields use API defaults.
type IPOptions struct {
	_    [0]func()
	Deep bool
}

// IP looks up an IP. Deep enrichment is included with a paid plan, without a separate check meter.
func (c *Client) IP(ctx context.Context, ip string, options ...IPOptions) (*IP, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &IP{}
	if err := c.get(ctx, "/ip/"+seg(ip), values("deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// IPSelfOptions configures IPSelf. Omitted fields use API defaults.
type IPSelfOptions struct {
	_    [0]func()
	Deep bool
}

// IPSelf looks up the public IP making this request. On a server, this is the server's IP.
func (c *Client) IPSelf(ctx context.Context, options ...IPSelfOptions) (*IP, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &IP{}
	if err := c.get(ctx, "/ip", values("deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// ContinentOptions reserves optional settings for Continent.
type ContinentOptions struct {
	_ [0]func()
}

// Continent calls /continent/{code}.
func (c *Client) Continent(ctx context.Context, code string, options ...ContinentOptions) (*Continent, error) {
	_, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &Continent{}
	if err := c.get(ctx, "/continent/"+seg(code), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// ContinentCountriesOptions reserves optional settings for ContinentCountries.
type ContinentCountriesOptions struct {
	_ [0]func()
}

// ContinentCountries calls /continent/{code}/countries.
func (c *Client) ContinentCountries(ctx context.Context, code string, options ...ContinentCountriesOptions) (*ContinentCountries, error) {
	_, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &ContinentCountries{}
	if err := c.get(ctx, "/continent/"+seg(code)+"/countries", nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// BlocOptions reserves optional settings for Bloc.
type BlocOptions struct {
	_ [0]func()
}

// Bloc calls /bloc/{code}.
func (c *Client) Bloc(ctx context.Context, code string, options ...BlocOptions) (*Bloc, error) {
	_, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &Bloc{}
	if err := c.get(ctx, "/bloc/"+seg(code), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// BlocCountriesOptions reserves optional settings for BlocCountries.
type BlocCountriesOptions struct {
	_ [0]func()
}

// BlocCountries calls /bloc/{code}/countries.
func (c *Client) BlocCountries(ctx context.Context, code string, options ...BlocCountriesOptions) (*BlocCountries, error) {
	_, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &BlocCountries{}
	if err := c.get(ctx, "/bloc/"+seg(code)+"/countries", nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CountryOptions reserves optional settings for Country.
type CountryOptions struct {
	_    [0]func()
	Deep bool
}

// Country calls /country/{code}.
func (c *Client) Country(ctx context.Context, code string, options ...CountryOptions) (*Country, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Country{}
	if err := c.get(ctx, "/country/"+seg(code), values("deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CountryStatesOptions reserves optional settings for CountryStates.
type CountryStatesOptions struct {
	_ [0]func()
}

// CountryStates calls /country/{code}/states.
func (c *Client) CountryStates(ctx context.Context, code string, options ...CountryStatesOptions) (*CountryStates, error) {
	_, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &CountryStates{}
	if err := c.get(ctx, "/country/"+seg(code)+"/states", nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// StateOptions configures State. Omitted fields use API defaults.
type StateOptions struct {
	_       [0]func()
	Country string
	Deep    bool
}

// State calls /state/{code}.
func (c *Client) State(ctx context.Context, code string, options ...StateOptions) (*State, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &State{}
	if err := c.get(ctx, "/state/"+seg(code), values("country", opts.Country, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// StateDistrictsOptions configures StateDistricts. Omitted fields use API defaults.
type StateDistrictsOptions struct {
	_       [0]func()
	Country string
	Deep    bool
}

// StateDistricts calls /state/{code}/districts.
func (c *Client) StateDistricts(ctx context.Context, code string, options ...StateDistrictsOptions) (*StateDistricts, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &StateDistricts{}
	if err := c.get(ctx, "/state/"+seg(code)+"/districts", values("country", opts.Country, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// DistrictOptions configures District. Omitted fields use API defaults.
type DistrictOptions struct {
	_       [0]func()
	Country string
	State   string
	Deep    bool
}

// District calls /district/{code}.
func (c *Client) District(ctx context.Context, code string, options ...DistrictOptions) (*District, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &District{}
	if err := c.get(ctx, "/district/"+seg(code), values("country", opts.Country, "state", opts.State, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CityOptions configures City. Omitted fields use API defaults.
type CityOptions struct {
	_       [0]func()
	Country string
	State   string
	Deep    bool
}

// City calls /city/{name}.
func (c *Client) City(ctx context.Context, name string, options ...CityOptions) (*City, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &City{}
	if err := c.get(ctx, "/city/"+seg(name), values("country", opts.Country, "state", opts.State, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CityIDOptions reserves optional settings for CityID.
type CityIDOptions struct {
	_    [0]func()
	Deep bool
}

// CityID calls /city/id/{id}.
func (c *Client) CityID(ctx context.Context, id string, options ...CityIDOptions) (*City, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &City{}
	if err := c.get(ctx, "/city/id/"+seg(id), values("deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CitySearchOptions configures CitySearch. Omitted fields use API defaults.
type CitySearchOptions struct {
	_       [0]func()
	Country string
	State   string
	Limit   int
	Deep    bool
}

// CitySearch calls /city.
func (c *Client) CitySearch(ctx context.Context, query string, options ...CitySearchOptions) (*CitySearch, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	limit := ""
	if opts.Limit != 0 {
		limit = strconv.Itoa(opts.Limit)
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &CitySearch{}
	if err := c.get(ctx, "/city", values("q", query, "country", opts.Country, "state", opts.State, "limit", limit, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CityNearestOptions reserves optional settings for CityNearest.
type CityNearestOptions struct {
	_    [0]func()
	Deep bool
}

// CityNearest calls /city.
func (c *Client) CityNearest(ctx context.Context, lat float64, lon float64, options ...CityNearestOptions) (*CityNearest, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &CityNearest{}
	if err := c.get(ctx, "/city", values("lat", f(lat), "lon", f(lon), "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CityNearbyOptions configures CityNearby. Omitted fields use API defaults.
type CityNearbyOptions struct {
	_       [0]func()
	Country string
	State   string
	Radius  float64
	Unit    string
	Limit   int
	Deep    bool
}

// CityNearby calls /city/{name}/nearby.
func (c *Client) CityNearby(ctx context.Context, name string, options ...CityNearbyOptions) (*CityNearby, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	radius := ""
	if opts.Radius != 0 {
		radius = f(opts.Radius)
	}
	limit := ""
	if opts.Limit != 0 {
		limit = strconv.Itoa(opts.Limit)
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &CityNearby{}
	if err := c.get(ctx, "/city/"+seg(name)+"/nearby", values("country", opts.Country, "state", opts.State, "radius", radius, "unit", opts.Unit, "limit", limit, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// PostalOptions configures Postal. Omitted fields use API defaults.
type PostalOptions struct {
	_       [0]func()
	Country string
	Deep    bool
}

// Postal looks up a postal area. Pass country when known. Check nullable coordinates before
// another location lookup.
func (c *Client) Postal(ctx context.Context, code string, options ...PostalOptions) (*Postal, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Postal{}
	if err := c.get(ctx, "/postal/"+seg(code), values("country", opts.Country, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// PostalNearbyOptions configures PostalNearby. Omitted fields use API defaults.
type PostalNearbyOptions struct {
	_       [0]func()
	Country string
	Radius  float64
	Unit    string
	Deep    bool
}

// PostalNearby calls /postal/{code}/nearby.
func (c *Client) PostalNearby(ctx context.Context, code string, options ...PostalNearbyOptions) (*PostalNearby, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	radius := ""
	if opts.Radius != 0 {
		radius = f(opts.Radius)
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &PostalNearby{}
	if err := c.get(ctx, "/postal/"+seg(code)+"/nearby", values("country", opts.Country, "radius", radius, "unit", opts.Unit, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// PostalDistanceOptions configures PostalDistance. Omitted fields use API defaults.
type PostalDistanceOptions struct {
	_       [0]func()
	Country string
	Deep    bool
}

// PostalDistance calls /postal/{code}/distance/{other}.
func (c *Client) PostalDistance(ctx context.Context, code string, other string, options ...PostalDistanceOptions) (*PostalDistance, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &PostalDistance{}
	if err := c.get(ctx, "/postal/"+seg(code)+"/distance/"+seg(other), values("country", opts.Country, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// EmailOptions configures Email. Omitted fields use API defaults.
type EmailOptions struct {
	_    [0]func()
	Deep bool
}

// Email parses an email and check its format and domain. Deep explicitly requests a metered
// deliverability check. Deep checks use one attempt by default. An explicit retry count can repeat
// paid usage.
func (c *Client) Email(ctx context.Context, email string, options ...EmailOptions) (*Email, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Email{}
	if err := c.get(ctx, "/email/"+seg(email), values("deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// VATOptions configures VAT. Omitted fields use API defaults.
type VATOptions struct {
	_       [0]func()
	Country string
	From    string
	Deep    bool
}

// VAT checks VAT format and checksum. Deep requests a metered registry check where supported. Deep
// checks use one attempt by default. Supply your own VAT number for a consultation reference when
// supported.
func (c *Client) VAT(ctx context.Context, number string, options ...VATOptions) (*VAT, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &VAT{}
	if err := c.get(ctx, "/vat/"+seg(number), values("country", opts.Country, "from", opts.From, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// IBANOptions configures IBAN. Omitted fields use API defaults.
type IBANOptions struct {
	_       [0]func()
	Country string
	Deep    bool
}

// IBAN calls /iban/{iban}.
func (c *Client) IBAN(ctx context.Context, iban string, options ...IBANOptions) (*IBAN, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &IBAN{}
	if err := c.get(ctx, "/iban/"+seg(iban), values("country", opts.Country, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// NPIOptions configures NPI. Omitted fields use API defaults.
type NPIOptions struct {
	_    [0]func()
	Deep bool
}

// NPI calls /npi/{npi}.
func (c *Client) NPI(ctx context.Context, npi string, options ...NPIOptions) (*NPI, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &NPI{}
	if err := c.get(ctx, "/npi/"+seg(npi), values("deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// PhoneOptions configures Phone. Omitted fields use API defaults.
type PhoneOptions struct {
	_       [0]func()
	Country string
	Deep    bool
}

// Phone parses a phone number and its formats. Pass country for national numbers when needed. Deep
// reveals numbering-plan location and timezone on every plan. Carrier, caller, and HLR are separate metered lookups.
func (c *Client) Phone(ctx context.Context, number string, options ...PhoneOptions) (*Phone, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Phone{}
	if err := c.get(ctx, "/phone/"+seg(number), values("country", opts.Country, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CarrierOptions configures Carrier. Omitted fields use API defaults.
type CarrierOptions struct {
	_       [0]func()
	Country string
	Deep    bool
}

// Carrier requests a metered carrier lookup. No automatic retries by default.
func (c *Client) Carrier(ctx context.Context, number string, options ...CarrierOptions) (*Carrier, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Carrier{}
	if err := c.get(ctx, "/carrier/"+seg(number), values("country", opts.Country, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CallerOptions configures Caller. Omitted fields use API defaults.
type CallerOptions struct {
	_       [0]func()
	Country string
}

// Caller requests a metered caller-name lookup for a NANP number. No automatic retries by default.
func (c *Client) Caller(ctx context.Context, number string, options ...CallerOptions) (*Caller, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &Caller{}
	if err := c.get(ctx, "/caller/"+seg(number), values("country", opts.Country), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// HLROptions configures HLR. Omitted fields use API defaults.
type HLROptions struct {
	_       [0]func()
	Country string
	Deep    bool
}

// HLR looks up phone status at the last check. Live means assigned and connected means reachable at
// that check. Cached results may be returned. Null means unconfirmed. Deep adds network
// diagnostics within the same metered lookup. No automatic retries by default.
func (c *Client) HLR(ctx context.Context, number string, options ...HLROptions) (*HLR, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &HLR{}
	if err := c.get(ctx, "/hlr/"+seg(number), values("country", opts.Country, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// DomainOptions configures Domain. Omitted fields use API defaults.
type DomainOptions struct {
	_    [0]func()
	Deep bool
}

// Domain checks whether a domain is registered.
// Deep adds registration dates, registrar, status and DNSSEC on paid plans.
func (c *Client) Domain(ctx context.Context, domain string, options ...DomainOptions) (*Domain, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Domain{}
	if err := c.get(ctx, "/domain/"+seg(domain), values("deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// ASNOptions reserves optional settings for ASN.
type ASNOptions struct {
	_ [0]func()
}

// ASN calls /asn/{asn}.
func (c *Client) ASN(ctx context.Context, asn string, options ...ASNOptions) (*ASN, error) {
	_, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &ASN{}
	if err := c.get(ctx, "/asn/"+seg(asn), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// MACOptions reserves optional settings for MAC.
type MACOptions struct {
	_ [0]func()
}

// MAC calls /mac/{mac}.
func (c *Client) MAC(ctx context.Context, mac string, options ...MACOptions) (*MAC, error) {
	_, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &MAC{}
	if err := c.get(ctx, "/mac/"+seg(mac), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// BINOptions configures BIN. Deep requests an empty object on every plan.
type BINOptions struct {
	_    [0]func()
	Deep bool
}

// BIN looks up a 6-11 digit card prefix. Preserve leading zeros in the string.
func (c *Client) BIN(ctx context.Context, bin string, options ...BINOptions) (*BIN, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &BIN{}
	if err := c.get(ctx, "/bin/"+seg(bin), values("deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}


// DNSOptions configures DNS. Type selects the question, including its CNAME chain.
// Omit Type to check all ten supported record types.
type DNSOptions struct {
	_    [0]func()
	Type string
}

// DNS returns published DNS records with TTLs. Pooled on every plan.
func (c *Client) DNS(ctx context.Context, domain string, options ...DNSOptions) (*DNS, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &DNS{}
	if err := c.get(ctx, "/dns/"+seg(domain), values("type", opts.Type), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// MXOptions reserves optional settings for MX.
type MXOptions struct {
	_ [0]func()
}

// MX calls /mx/{domain}.
func (c *Client) MX(ctx context.Context, domain string, options ...MXOptions) (*MX, error) {
	_, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &MX{}
	if err := c.get(ctx, "/mx/"+seg(domain), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// UserAgentOptions configures UserAgent. Omitted fields use API defaults.
type UserAgentOptions struct {
	_    [0]func()
	Deep bool
}

// UserAgent calls /useragent.
func (c *Client) UserAgent(ctx context.Context, ua string, options ...UserAgentOptions) (*UserAgent, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &UserAgent{}
	if err := c.get(ctx, "/useragent", values("deep", deep), map[string]string{"User-Agent": ua}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// VINOptions configures VIN. Omitted fields use API defaults.
type VINOptions struct {
	_    [0]func()
	Deep bool
}

// VIN calls /vin/{vin}.
func (c *Client) VIN(ctx context.Context, vin string, options ...VINOptions) (*VIN, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &VIN{}
	if err := c.get(ctx, "/vin/"+seg(vin), values("deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// NAICSOptions reserves optional settings for NAICS.
type NAICSOptions struct {
	_    [0]func()
	Deep bool
}

// NAICS looks up a US NAICS 2022 code and its hierarchy.
func (c *Client) NAICS(ctx context.Context, code string, options ...NAICSOptions) (*NAICS, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &NAICS{}
	if err := c.get(ctx, "/naics/"+seg(code), values("deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// NAICSSearchOptions configures keyword search. Limit defaults to 10 and accepts 1-50.
type NAICSSearchOptions struct {
	_     [0]func()
	Limit int
	Deep  bool
}

// NAICSSearch searches US NAICS 2022 industry names and activity terms.
func (c *Client) NAICSSearch(ctx context.Context, query string, options ...NAICSSearchOptions) (*NAICSSearch, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	limit := ""
	if opts.Limit != 0 {
		limit = strconv.Itoa(opts.Limit)
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &NAICSSearch{}
	if err := c.get(ctx, "/naics", values("q", query, "limit", limit, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// TariffOptions configures Tariff. Omitted fields use API defaults.
type TariffOptions struct {
	_ [0]func()
	// Add units and the special and other schedule columns on paid plans.
	Deep bool
	// ISO 3166-1 alpha-2 origin. With paid deep, resolves country-specific measures. Optional for schedule detail.
	Origin string
}

// Tariff looks up the general US duty schedule line. Paid deep adds units and the special and other
// schedule columns. Add origin with deep to resolve country-specific measures. Without origin,
// schedule detail remains available and origin-dependent fields are null. A null effective rate is
// not a zero rate.
func (c *Client) Tariff(ctx context.Context, code string, options ...TariffOptions) (*Tariff, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Tariff{}
	if err := c.get(ctx, "/tariff/"+seg(code), values("deep", deep, "origin", opts.Origin), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// TariffSearchOptions reserves optional settings for TariffSearch.
type TariffSearchOptions struct {
	_ [0]func()
}

// TariffSearch calls /tariff.
func (c *Client) TariffSearch(ctx context.Context, query string, options ...TariffSearchOptions) (*TariffSearch, error) {
	_, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &TariffSearch{}
	if err := c.get(ctx, "/tariff", values("q", query), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CurrencyOptions reserves optional settings for Currency.
type CurrencyOptions struct {
	_    [0]func()
	Deep bool
}

// Currency calls /currency/{code}.
func (c *Client) Currency(ctx context.Context, code string, options ...CurrencyOptions) (*Currency, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Currency{}
	if err := c.get(ctx, "/currency/"+seg(code), values("deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// LanguageOptions reserves optional settings for Language.
type LanguageOptions struct {
	_    [0]func()
	Deep bool
}

// Language calls /language/{code}.
func (c *Client) Language(ctx context.Context, code string, options ...LanguageOptions) (*Language, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Language{}
	if err := c.get(ctx, "/language/"+seg(code), values("deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// NameOptions configures Name. Country is an ISO2 gender context.
type NameOptions struct {
	_       [0]func()
	Country string
	Deep    bool
	// NameLocale selects CLDR formatting rules, defaulting to en. Parsing stays unchanged.
	NameLocale string
}

// Name calls /name/{name}.
func (c *Client) Name(ctx context.Context, name string, options ...NameOptions) (*Name, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Name{}
	if err := c.get(ctx, "/name/"+seg(name), values("country", opts.Country, "deep", deep, "name_locale", opts.NameLocale), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CurrencyRateOptions configures CurrencyRate. Omitted fields use API defaults.
type CurrencyRateOptions struct {
	_      [0]func()
	Date   string
	Amount *float64
}

// CurrencyRate calls /currency/{base}/{quote}.
func (c *Client) CurrencyRate(ctx context.Context, base string, quote string, options ...CurrencyRateOptions) (*CurrencyRate, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	amount := ""
	if opts.Amount != nil {
		amount = f(*opts.Amount)
	}
	out := &CurrencyRate{}
	if err := c.get(ctx, "/currency/"+seg(base)+"/"+seg(quote), values("date", opts.Date, "amount", amount), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// TimeOptions configures Time. With To, offsetless At is source wall time.
type TimeOptions struct {
	_    [0]func()
	At   string
	To   string
	Deep bool
}

// Time returns local time and timezone facts. An empty timezone selects UTC.
func (c *Client) Time(ctx context.Context, timezone string, options ...TimeOptions) (*Time, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	path := "/time"
	if timezone != "" {
		path += "/" + seg(timezone)
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Time{}
	if err := c.get(ctx, path, values("at", opts.At, "to", opts.To, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// TimeAtOptions configures TimeAt. Omitted fields use API defaults.
type TimeAtOptions struct {
	_    [0]func()
	At   string
	To   string
	Deep bool
}

// TimeAt returns local time at the coordinates, optionally converted with To.
func (c *Client) TimeAt(ctx context.Context, lat float64, lon float64, options ...TimeAtOptions) (*Time, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Time{}
	if err := c.get(ctx, "/time", values("lat", f(lat), "lon", f(lon), "at", opts.At, "to", opts.To, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// TimezoneOptions configures Timezone. Omitted fields use API defaults.
type TimezoneOptions struct {
	_    [0]func()
	At   string
	To   string
	Deep bool
}

// Timezone calls /timezone/{timezone}.
func (c *Client) Timezone(ctx context.Context, timezone string, options ...TimezoneOptions) (*Timezone, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Timezone{}
	if err := c.get(ctx, "/timezone/"+seg(timezone), values("at", opts.At, "to", opts.To, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// TimezoneAtOptions configures TimezoneAt. Omitted fields use API defaults.
type TimezoneAtOptions struct {
	_    [0]func()
	At   string
	Deep bool
}

// TimezoneAt calls /timezone.
func (c *Client) TimezoneAt(ctx context.Context, lat float64, lon float64, options ...TimezoneAtOptions) (*Timezone, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Timezone{}
	if err := c.get(ctx, "/timezone", values("lat", f(lat), "lon", f(lon), "at", opts.At, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// DateOptions configures Date. Omitted fields use API defaults.
type DateOptions struct {
	_      [0]func()
	Format string
	To     string
	Deep   bool
}

// Date calls /date/{date}.
func (c *Client) Date(ctx context.Context, date string, options ...DateOptions) (*DateInfo, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &DateInfo{}
	if err := c.get(ctx, "/date/"+seg(date), values("format", opts.Format, "to", opts.To, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// DateTodayOptions configures DateToday. Omitted fields use API defaults.
type DateTodayOptions struct {
	_    [0]func()
	To   string
	Deep bool
}

// DateToday calls /date.
func (c *Client) DateToday(ctx context.Context, options ...DateTodayOptions) (*DateInfo, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &DateInfo{}
	if err := c.get(ctx, "/date", values("to", opts.To, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// HolidayOptions configures Holiday. Omitted fields use API defaults.
type HolidayOptions struct {
	_    [0]func()
	Year int
}

// Holiday calls /holiday/{country}.
func (c *Client) Holiday(ctx context.Context, country string, options ...HolidayOptions) (*HolidayYear, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	year := ""
	if opts.Year != 0 {
		year = strconv.Itoa(opts.Year)
	}
	out := &HolidayYear{}
	if err := c.get(ctx, "/holiday/"+seg(country), values("year", year), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// HolidayDateOptions reserves optional settings for HolidayDate.
type HolidayDateOptions struct {
	_ [0]func()
}

// HolidayDate calls /holiday/{country}/{date}.
func (c *Client) HolidayDate(ctx context.Context, country string, date string, options ...HolidayDateOptions) (*HolidayDate, error) {
	_, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &HolidayDate{}
	if err := c.get(ctx, "/holiday/"+seg(country)+"/"+seg(date), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// ElevationOptions reserves optional settings for Elevation.
type ElevationOptions struct {
	_ [0]func()
}

// Elevation calls /elevation.
func (c *Client) Elevation(ctx context.Context, lat float64, lon float64, options ...ElevationOptions) (*Elevation, error) {
	_, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &Elevation{}
	if err := c.get(ctx, "/elevation", values("lat", f(lat), "lon", f(lon)), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// PointOptions configures Point. Omitted fields use API defaults.
type PointOptions struct {
	_ [0]func()
	// Add terrain and compact nearest-city context on every plan. The timezone ID stays in core.
	Deep bool
}

// Point resolves the country, state, district and timezone at coordinates. Deep adds terrain and compact
// nearest-city context on every plan. The timezone ID stays in core. The nearest city is null when
// none is within 200 km.
func (c *Client) Point(ctx context.Context, lat float64, lon float64, options ...PointOptions) (*Point, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Point{}
	if err := c.get(ctx, "/point", values("lat", f(lat), "lon", f(lon), "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// WeatherOptions configures Weather. Omitted fields use API defaults.
type WeatherOptions struct {
	_ [0]func()
	// Add specialist current measurements, forecasts and related detail on paid plans.
	Deep bool
	// Past UTC day (YYYY-MM-DD). Requires paid deep and adds deep.history alongside current conditions.
	Date string
}

// Weather gets current conditions in metric and imperial units. Paid deep adds specialist current
// measurements, forecasts and related detail. With deep, date selects a past UTC day (YYYY-MM-DD)
// in deep.history alongside current conditions. Date alone does not request history.
func (c *Client) Weather(ctx context.Context, lat float64, lon float64, options ...WeatherOptions) (*Weather, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Weather{}
	if err := c.get(ctx, "/weather", values("lat", f(lat), "lon", f(lon), "deep", deep, "date", opts.Date), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// EmojiOptions reserves optional settings for Emoji.
type EmojiOptions struct {
	_    [0]func()
	Deep bool
}

// Emoji calls /emoji/{emoji}.
func (c *Client) Emoji(ctx context.Context, emoji string, options ...EmojiOptions) (*Emoji, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Emoji{}
	if err := c.get(ctx, "/emoji/"+seg(emoji), values("deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// EmojiSearchOptions configures EmojiSearch. Omitted fields use API defaults.
type EmojiSearchOptions struct {
	_     [0]func()
	Limit int
	Deep  bool
}

// EmojiSearch calls /emoji.
func (c *Client) EmojiSearch(ctx context.Context, query string, options ...EmojiSearchOptions) (*EmojiSearch, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	limit := ""
	if opts.Limit != 0 {
		limit = strconv.Itoa(opts.Limit)
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &EmojiSearch{}
	if err := c.get(ctx, "/emoji", values("q", query, "limit", limit, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// AddressOptions configures Address. Omitted fields use API defaults.
type AddressOptions struct {
	_       [0]func()
	Country string
	Deep    bool
}

// Address calls /address/{address}.
func (c *Client) Address(ctx context.Context, address string, options ...AddressOptions) (*Address, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Address{}
	if err := c.get(ctx, "/address/"+seg(address), values("country", opts.Country, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// AddressSearchOptions configures AddressSearch. Omitted fields use API defaults.
type AddressSearchOptions struct {
	_       [0]func()
	Country string
	Postal  string
	City    string
	State   string
	IP      string
}

// AddressSearch finds address suggestions using the context supplied. Prefer postal, or city and state, from the
// form; ip is an optional end-user locality hint for server-side calls. An empty result has reason
// more_input, missing_context or no_matches. Suggestions have reason null. Operational failures
// are errors.
func (c *Client) AddressSearch(ctx context.Context, query string, options ...AddressSearchOptions) (*AddressSearch, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &AddressSearch{}
	if err := c.get(ctx, "/address", values("q", query, "country", opts.Country, "postal", opts.Postal, "city", opts.City, "state", opts.State, "ip", opts.IP), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CompanyOptions configures Company. Omitted fields use API defaults.
type CompanyOptions struct {
	_       [0]func()
	Country string
	Deep    bool
}

// Company calls /company/{number}.
func (c *Client) Company(ctx context.Context, number string, options ...CompanyOptions) (*Company, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Company{}
	if err := c.get(ctx, "/company/"+seg(number), values("country", opts.Country, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// MeasureOptions configures measurement conversion and explicit ambiguity resolution.
type MeasureOptions struct {
	_      [0]func()
	To     string
	Locale string
	System string
}

// Measure parses or converts a measurement. Amount is a decimal string. Without To,
// the API uses the type's canonical unit. System accepts us or imperial.
func (c *Client) Measure(ctx context.Context, measure string, options ...MeasureOptions) (*Measure, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &Measure{}
	if err := c.get(ctx, "/measure/"+seg(measure), values("to", opts.To, "locale", opts.Locale, "system", opts.System), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// MeasureUnitsOptions filters the reviewed unit catalog. Unit selects compatible targets.
type MeasureUnitsOptions struct {
	_     [0]func()
	Query string
	Type  string
	Unit  string
}

// MeasureUnits discovers reviewed units and their accepted aliases.
func (c *Client) MeasureUnits(ctx context.Context, options ...MeasureUnitsOptions) (*MeasureUnits, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &MeasureUnits{}
	if err := c.get(ctx, "/measure/units", values("q", opts.Query, "type", opts.Type, "unit", opts.Unit), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
