// Package parseapi is the official ParseAPI client for Go.
// One key, minimal JSON, fast. https://parseapi.com
package parseapi

import (
	"bytes"
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
	version         = "1.10.7"
	defaultBaseURL  = "https://api.parseapi.com"
	defaultTimeout  = 10 * time.Second
	defaultRetries  = 2
	retryAfterCapMs = 5000
)

// The response types' wire contract. Changes require a reviewed major SDK release.
const apiVersion = "2.0.0"

var retryStatus = map[int]bool{429: true, 500: true, 502: true, 503: true, 504: true}

// Error is every non-2xx response from the API. Branch on Code, never on Message.
type Error struct {
	_         [0]func()
	Status    int
	Code      string
	Message   string
	Docs      string
	RequestID string
	// RetryAfter is the original HTTP Retry-After header, or nil when absent.
	RetryAfter *string
	// Candidates is populated only for structured HTTP409 ambiguous_city errors.
	Candidates []CityAmbiguityCandidate
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
	timeoutSet bool
	httpClient *http.Client
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides https://api.parseapi.com (tests, canaries).
func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.baseURL = strings.TrimRight(baseURL, "/") }
}

// WithTimeout sets the per-attempt timeout for every operation. Defaults are 35s for Stack and 10s otherwise.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeoutSet = true
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
		c.timeoutSet = true
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

// A negative delay declines an automatic retry that would exceed its wait budget.
func retryDelay(attempt int, retryAfter string) time.Duration {
	if retryAfter != "" {
		value := strings.TrimSpace(retryAfter)
		numeric := value != ""
		dots := 0
		for i, ch := range value {
			if ch == '.' && i > 0 && i < len(value)-1 && dots == 0 {
				dots++
				continue
			}
			if ch < '0' || ch > '9' {
				numeric = false
				break
			}
		}
		if numeric {
			seconds, _ := strconv.ParseFloat(value, 64)
			// Compare before converting, including positive float overflow.
			if seconds > float64(retryAfterCapMs)/1000 {
				return -1
			}
			return time.Duration(seconds * float64(time.Second))
		}
		if at, err := http.ParseTime(value); err == nil {
			delay := time.Until(at)
			if delay < 0 {
				return 0
			}
			if delay > 5*time.Second {
				return -1
			}
			return delay
		}
	}
	return time.Duration(rand.Float64()*math.Min(250*math.Pow(2, math.Min(float64(attempt), 16)), retryAfterCapMs)) * time.Millisecond
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
	return c.request(ctx, http.MethodGet, path, query, headers, nil, out)
}

func (c *Client) post(ctx context.Context, path string, input any, out any) error {
	body, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return c.request(ctx, http.MethodPost, path, nil, nil, body, out)
}

func (c *Client) request(ctx context.Context, method, path string, query url.Values, headers map[string]string, body []byte, out any) error {
	// A request-local copy preserves concurrent lookups and caller-owned clients.
	httpClient := c.httpClient
	if !c.timeoutSet && strings.HasPrefix(path, "/stack/") {
		stackClient := *httpClient
		stackClient.Timeout = 35 * time.Second
		httpClient = &stackClient
	}
	retries := c.retries
	if !c.retriesSet && meteredRequest(path, query) {
		retries = 0
	}
	target := c.baseURL + path
	if len(query) > 0 {
		target += "?" + query.Encode()
	}

	for attempt := 0; ; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, target, bytes.NewReader(body))
		if err != nil {
			return err
		}
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		req.Header.Set("X-API-Key", c.apiKey)
		req.Header.Set("Parse-Version", apiVersion)
		req.Header.Set("User-Agent", "parseapi-go/"+version)
		for name, value := range headers {
			req.Header.Set(name, value)
		}

		res, err := httpClient.Do(req)
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
			if delay := retryDelay(attempt, retryAfter); delay >= 0 {
				res.Body.Close()
				if err := waitRetry(ctx, delay); err != nil {
					return err
				}
				continue
			}
		}

		apiErr := &Error{
			Status:  res.StatusCode,
			Code:    "unknown_error",
			Message: fmt.Sprintf("Request failed with status %d", res.StatusCode),
		}
		if values, present := res.Header[http.CanonicalHeaderKey("Retry-After")]; present && len(values) > 0 {
			value := values[0]
			apiErr.RetryAfter = &value
		}
		var body struct {
			Code       string            `json:"code"`
			Message    string            `json:"message"`
			Docs       string            `json:"docs"`
			RequestID  string            `json:"request_id"`
			Candidates []json.RawMessage `json:"candidates"`
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
			if apiErr.Status == 409 && apiErr.Code == "ambiguous_city" && body.Candidates != nil {
				apiErr.Candidates = cityErrorCandidates(body.Candidates)
			}
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
	// Lang selects translated display names for this request.
	Lang string
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
	if err := c.get(ctx, "/ip/"+seg(ip), values("lang", opts.Lang, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// IPSelfOptions configures IPSelf. Omitted fields use API defaults.
type IPSelfOptions struct {
	// Lang selects translated display names for this request.
	Lang string
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
	if err := c.get(ctx, "/ip", values("lang", opts.Lang, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// ContinentOptions reserves optional settings for Continent.
type ContinentOptions struct {
	// Lang selects translated display names for this request.
	Lang string
	_    [0]func()
}

// Continent calls /continent/{code}.
func (c *Client) Continent(ctx context.Context, code string, options ...ContinentOptions) (*Continent, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &Continent{}
	if err := c.get(ctx, "/continent/"+seg(code), values("lang", opts.Lang), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// ContinentCountriesOptions reserves optional settings for ContinentCountries.
type ContinentCountriesOptions struct {
	// Lang selects translated display names for this request.
	Lang string
	_    [0]func()
}

// ContinentCountries calls /continent/{code}/countries.
func (c *Client) ContinentCountries(ctx context.Context, code string, options ...ContinentCountriesOptions) (*ContinentCountries, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &ContinentCountries{}
	if err := c.get(ctx, "/continent/"+seg(code)+"/countries", values("lang", opts.Lang), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}









// CountryOptions reserves optional settings for Country.
type CountryOptions struct {
	// Lang selects translated display names for this request.
	Lang string
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
	if err := c.get(ctx, "/country/"+seg(code), values("lang", opts.Lang, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CountryStatesOptions reserves optional settings for CountryStates.
type CountryStatesOptions struct {
	// Lang selects translated display names for this request.
	Lang string
	_    [0]func()
}

// CountryStates calls /country/{code}/states.
func (c *Client) CountryStates(ctx context.Context, code string, options ...CountryStatesOptions) (*CountryStates, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &CountryStates{}
	if err := c.get(ctx, "/country/"+seg(code)+"/states", values("lang", opts.Lang), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// StateOptions configures State. Omitted fields use API defaults.
type StateOptions struct {
	// Lang selects translated display names for this request.
	Lang    string
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
	if err := c.get(ctx, "/state/"+seg(code), values("lang", opts.Lang, "country", opts.Country, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// StateDistrictsOptions configures StateDistricts. Omitted fields use API defaults.
type StateDistrictsOptions struct {
	// Lang selects translated display names for this request.
	Lang    string
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
	if err := c.get(ctx, "/state/"+seg(code)+"/districts", values("lang", opts.Lang, "country", opts.Country, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// DistrictOptions configures District. Omitted fields use API defaults.
type DistrictOptions struct {
	// Lang selects translated display names for this request.
	Lang    string
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
	if err := c.get(ctx, "/district/"+seg(code), values("lang", opts.Lang, "country", opts.Country, "state", opts.State, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CityOptions configures City. Omitted fields use API defaults.
type CityOptions struct {
	// Lang selects translated display names for this request.
	Lang      string
	_         [0]func()
	Country   string
	State     string
	Deep      bool
	Type      string
	Grain     string
	Hierarchy bool
	Aliases   bool
	Ambiguity string
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
	if err := c.get(ctx, "/city/"+seg(name), values("lang", opts.Lang, "country", opts.Country, "state", opts.State, "deep", deep, "type", opts.Type, "grain", opts.Grain, "hierarchy", cityFlag(opts.Hierarchy), "aliases", cityFlag(opts.Aliases), "ambiguity", opts.Ambiguity), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CityIDOptions reserves optional settings for CityID.
type CityIDOptions struct {
	// Lang selects translated display names for this request.
	Lang      string
	_         [0]func()
	Deep      bool
	Hierarchy bool
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
	if err := c.get(ctx, "/city/id/"+seg(id), values("lang", opts.Lang, "deep", deep, "hierarchy", cityFlag(opts.Hierarchy)), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CitySearchOptions configures CitySearch. Omitted fields use API defaults.
type CitySearchOptions struct {
	// Lang selects translated display names for this request.
	Lang             string
	_                [0]func()
	Country          string
	State            string
	Limit            int
	Deep             bool
	Type             string
	Grain            string
	Hierarchy        bool
	Aliases          bool
	MinPopulation    *int64
	MaxPopulation    *int64
	Timezone         string
	BBox             string
	Sort             string
	Order            string
	Page             bool
	Cursor           string
	Countries        []string
	ExcludeCountries []string
	Timezones        []string
	MinElevation     *float64
	MaxElevation     *float64
	Match            string
	MaxDistance      int
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
	params := values("lang", opts.Lang, "q", query, "country", opts.Country, "state", opts.State, "limit", limit, "deep", deep, "type", opts.Type, "grain", opts.Grain, "hierarchy", cityFlag(opts.Hierarchy), "aliases", cityFlag(opts.Aliases))
	cityCollectionQuery(params, opts.MinPopulation, opts.MaxPopulation, opts.Timezone, opts.BBox, opts.Sort, opts.Order, opts.Page, opts.Cursor)
	cityExtraQuery(params, opts.Countries, opts.ExcludeCountries, opts.Timezones, opts.MinElevation, opts.MaxElevation)
	if opts.Match != "" {
		params.Set("match", opts.Match)
	}
	if opts.MaxDistance != 0 {
		params.Set("max_distance", strconv.Itoa(opts.MaxDistance))
	}
	if opts.Match != "" || opts.MaxDistance != 0 {
		var page struct {
			CitySearch
			Data []City `json:"data"`
		}
		if err := c.get(ctx, "/city/search", params, nil, &page); err != nil {
			return nil, err
		}
		page.Cities = page.Data
		return &page.CitySearch, nil
	}
	if err := c.get(ctx, "/city", params, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CityNearestOptions reserves optional settings for CityNearest.
type CityNearestOptions struct {
	// Lang selects translated display names for this request.
	Lang      string
	_         [0]func()
	Deep      bool
	Country   string
	State     string
	Type      string
	Grain     string
	Hierarchy bool
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
	if err := c.get(ctx, "/city", values("lang", opts.Lang, "lat", f(lat), "lon", f(lon), "deep", deep, "country", opts.Country, "state", opts.State, "type", opts.Type, "grain", opts.Grain, "hierarchy", cityFlag(opts.Hierarchy)), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CityNearbyOptions configures CityNearby. Omitted fields use API defaults.
type CityNearbyOptions struct {
	// Lang selects translated display names for this request.
	Lang             string
	_                [0]func()
	Country          string
	State            string
	Radius           float64
	Unit             string
	Limit            int
	Deep             bool
	Type             string
	Grain            string
	Hierarchy        bool
	Aliases          bool
	Ambiguity        string
	MinPopulation    *int64
	MaxPopulation    *int64
	Timezone         string
	BBox             string
	Sort             string
	Order            string
	Page             bool
	Cursor           string
	Countries        []string
	ExcludeCountries []string
	Timezones        []string
	MinElevation     *float64
	MaxElevation     *float64
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
	query := values("lang", opts.Lang, "country", opts.Country, "state", opts.State, "radius", radius, "unit", opts.Unit, "limit", limit, "deep", deep, "type", opts.Type, "grain", opts.Grain, "hierarchy", cityFlag(opts.Hierarchy), "aliases", cityFlag(opts.Aliases), "ambiguity", opts.Ambiguity)
	cityCollectionQuery(query, opts.MinPopulation, opts.MaxPopulation, opts.Timezone, opts.BBox, opts.Sort, opts.Order, opts.Page, opts.Cursor)
	cityExtraQuery(query, opts.Countries, opts.ExcludeCountries, opts.Timezones, opts.MinElevation, opts.MaxElevation)
	if err := c.get(ctx, "/city/"+seg(name)+"/nearby", query, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// PostalOptions configures Postal. Omitted fields use API defaults.
type PostalOptions struct {
	// Lang selects translated display names for this request.
	Lang    string
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
	if err := c.get(ctx, "/postal/"+seg(code), values("lang", opts.Lang, "country", opts.Country, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// PostalNearbyOptions configures PostalNearby. Omitted fields use API defaults.
type PostalNearbyOptions struct {
	// Lang selects translated display names for this request.
	Lang    string
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
	if err := c.get(ctx, "/postal/"+seg(code)+"/nearby", values("lang", opts.Lang, "country", opts.Country, "radius", radius, "unit", opts.Unit, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// PostalDistanceOptions configures PostalDistance. Omitted fields use API defaults.
type PostalDistanceOptions struct {
	// Lang selects translated display names for this request.
	Lang    string
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
	if err := c.get(ctx, "/postal/"+seg(code)+"/distance/"+seg(other), values("lang", opts.Lang, "country", opts.Country, "deep", deep), nil, out); err != nil {
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

// StackOptions configures website technology observation.
type StackOptions struct {
	_      [0]func()
	Deep   bool
	Pretty bool
}

// Stack observes technologies on a public hostname without a scheme or path.
// Nil collections mean the check did not complete.
func (c *Client) Stack(ctx context.Context, domain string, options ...StackOptions) (*Stack, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep, pretty := "", ""
	if opts.Deep {
		deep = "true"
	}
	if opts.Pretty {
		pretty = "true"
	}
	out := &Stack{}
	if err := c.get(ctx, "/stack/"+seg(domain), values("deep", deep, "pretty", pretty), nil, out); err != nil {
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
	// Lang selects translated display names for this request.
	Lang string
	_    [0]func()
}

// ASN calls /asn/{asn}.
func (c *Client) ASN(ctx context.Context, asn string, options ...ASNOptions) (*ASN, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &ASN{}
	if err := c.get(ctx, "/asn/"+seg(asn), values("lang", opts.Lang), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}


// Card looks up a 2-11 digit card prefix. Preserve leading zeros in the string.
func (c *Client) Card(ctx context.Context, bin string) (*Card, error) {
	return c.CardWithOptions(ctx, bin, CardOptions{})
}

// CardOptions requests optional recorded issuer details, included on every plan.
type CardOptions struct {
	_    [0]func()
	Deep bool
}

// CardWithOptions adds recorded issuer details without changing core identity.
func (c *Client) CardWithOptions(ctx context.Context, bin string, opts CardOptions) (*Card, error) {
	digits := 0
	valid := len(bin) <= 64
	for _, ch := range bin {
		if ch >= '0' && ch <= '9' {
			digits++
		} else if ch != ' ' && ch != '\t' && ch != '\r' && ch != '\n' && ch != '-' {
			valid = false
		}
	}
	if !valid || digits < 2 || digits > 11 {
		return nil, errors.New("parseapi: Card requires a string containing 2 to 11 digits. Send a prefix only.")
	}
	out := &Card{}
	query := url.Values{}
	if opts.Deep {
		query.Set("deep", "true")
	}
	if err := c.get(ctx, "/card/"+seg(bin), query, nil, out); err != nil {
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




// NAICSOptions reserves optional settings for Industry.
type NAICSOptions struct {
	_    [0]func()
	Deep bool
}

// Industry looks up a US NAICS 2022 code and its hierarchy.
func (c *Client) Industry(ctx context.Context, code string, options ...NAICSOptions) (*Industry, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &Industry{}
	if err := c.get(ctx, "/industry/"+seg(code), values("deep", deep), nil, out); err != nil {
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

// IndustrySearch searches US NAICS 2022 industry names and activity terms.
func (c *Client) IndustrySearch(ctx context.Context, query string, options ...NAICSSearchOptions) (*IndustrySearch, error) {
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
	out := &IndustrySearch{}
	if err := c.get(ctx, "/industry", values("q", query, "limit", limit, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}




// LanguageOptions reserves optional settings for Language.
type LanguageOptions struct {
	// Lang selects translated display names for this request.
	Lang string
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
	if err := c.get(ctx, "/language/"+seg(code), values("lang", opts.Lang, "deep", deep), nil, out); err != nil {
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



// TimeOptions configures Time. With To, offsetless At is source wall time.
type TimeOptions struct {
	IP       string
	City     string
	Country  string
	State    string
	IATA     string
	ICAO     string
	UNLOCODE string
	Address  string

	// Lang selects translated display names for this request.
	Lang string
	_    [0]func()
	At   string
	To   string
	Deep bool
	// Disambiguation selects compatible (default), earlier, later, or reject for offsetless conversion at clock changes.
	Disambiguation string
	// Targets converts to 1-10 zones, preserving order and duplicates. Mutually exclusive with To.
	Targets []string
}

// Time returns local time and timezone facts. An empty timezone selects UTC.
func (c *Client) Time(ctx context.Context, timezone string, options ...TimeOptions) (*Time, error) {
	if source := strings.ToLower(strings.TrimSpace(timezone)); source == "zones" || source == "help" {
		return nil, errors.New("Time source must be an IANA timezone ID. Use timezone discovery to list IDs.")
	}
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	if err := timeSource(timezone, opts); err != nil {
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
	targets, err := timeTargets(opts.Targets, opts.To)
	if err != nil {
		return nil, err
	}
	out := &Time{}
	if err := c.get(ctx, path, values("ip", opts.IP, "city", opts.City, "country", opts.Country, "state", opts.State, "iata", opts.IATA, "icao", opts.ICAO, "unlocode", opts.UNLOCODE, "address", opts.Address, "lang", opts.Lang, "at", opts.At, "to", opts.To, "targets", targets, "disambiguation", opts.Disambiguation, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// TimeAtOptions configures TimeAt. Omitted fields use API defaults.
type TimeAtOptions struct {
	// Lang selects translated display names for this request.
	Lang string
	_    [0]func()
	At   string
	To   string
	Deep bool
	// Disambiguation selects compatible (default), earlier, later, or reject for offsetless conversion at clock changes.
	Disambiguation string
	// Targets converts to 1-10 zones, preserving order and duplicates. Mutually exclusive with To.
	Targets []string
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
	targets, err := timeTargets(opts.Targets, opts.To)
	if err != nil {
		return nil, err
	}
	out := &Time{}
	if err := c.get(ctx, "/time", values("lang", opts.Lang, "lat", f(lat), "lon", f(lon), "at", opts.At, "to", opts.To, "targets", targets, "disambiguation", opts.Disambiguation, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// TimeZonesOptions configures TimeZones.
type TimeZonesOptions struct {
	_            [0]func()
	Country      string
	Area         string
	Offset       string
	Abbreviation string
	DST          *bool
	ObservesDST  *bool
	At           string
	Details      bool
	Sort         string
}

// TimeZones searches serving timezone IDs. An empty query lists all.
func (c *Client) TimeZones(ctx context.Context, query string, options ...TimeZonesOptions) (*TimeZones, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	dst, observesDST, details := "", "", ""
	if opts.DST != nil {
		dst = strconv.FormatBool(*opts.DST)
	}
	if opts.ObservesDST != nil {
		observesDST = strconv.FormatBool(*opts.ObservesDST)
	}
	if opts.Details {
		details = "true"
	}
	out := &TimeZones{}
	if err := c.get(ctx, "/time/zones", values("q", query, "country", opts.Country, "area", opts.Area, "offset", opts.Offset, "abbreviation", opts.Abbreviation, "dst", dst, "observes_dst", observesDST, "at", opts.At, "details", details, "sort", opts.Sort), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func timeSource(timezone string, opts TimeOptions) error {
	primary := 0
	for _, value := range []string{opts.IP, opts.City, opts.IATA, opts.ICAO, opts.UNLOCODE, opts.Address} {
		if value != "" {
			primary++
		}
	}
	for _, value := range []string{opts.IP, opts.City, opts.Country, opts.State, opts.IATA, opts.ICAO, opts.UNLOCODE, opts.Address} {
		if value != "" && strings.TrimSpace(value) == "" {
			return errors.New("Pass one Time source, using country only with city or address and state only with city or address and country.")
		}
	}
	if (timezone != "" && (primary > 0 || opts.Country != "" || opts.State != "")) || primary > 1 ||
		(opts.Country != "" && primary > 0 && opts.City == "" && opts.Address == "") ||
		(opts.State != "" && ((opts.City == "" && opts.Address == "") || opts.Country == "")) || (opts.Address != "" && opts.Country == "") {
		return errors.New("Pass one Time source, using country only with city or address and state only with city or address and country.")
	}
	return nil
}

func timeTargets(targets []string, to string) (string, error) {
	if targets == nil {
		return "", nil
	}
	if to != "" || len(targets) < 1 || len(targets) > 10 {
		return "", fmt.Errorf("parseapi: Time targets requires 1 to 10 timezone IDs and cannot be combined with to")
	}
	for _, zone := range targets {
		if strings.TrimSpace(zone) == "" || strings.Contains(zone, ",") {
			return "", fmt.Errorf("parseapi: Time targets requires 1 to 10 timezone IDs and cannot be combined with to")
		}
	}
	return strings.Join(targets, ","), nil
}

// TimezoneOptions configures Timezone. Omitted fields use API defaults.
type TimezoneOptions struct {
	// Lang selects translated display names for this request.
	Lang string
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
	if err := c.get(ctx, "/timezone/"+seg(timezone), values("lang", opts.Lang, "at", opts.At, "to", opts.To, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// TimezoneAtOptions configures TimezoneAt. Omitted fields use API defaults.
type TimezoneAtOptions struct {
	// Lang selects translated display names for this request.
	Lang string
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
	if err := c.get(ctx, "/timezone", values("lang", opts.Lang, "lat", f(lat), "lon", f(lon), "at", opts.At, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// DateOptions configures Date. Omitted fields use API defaults.
type DateOptions struct {
	// Lang selects translated display names for this request.
	Lang   string
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
	if err := c.get(ctx, "/date/"+seg(date), values("lang", opts.Lang, "format", opts.Format, "to", opts.To, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// DateTodayOptions configures DateToday. Omitted fields use API defaults.
type DateTodayOptions struct {
	// Lang selects translated display names for this request.
	Lang string
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
	if err := c.get(ctx, "/date", values("lang", opts.Lang, "to", opts.To, "deep", deep), nil, out); err != nil {
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
	// Lang selects translated display names for this request.
	Lang string
	_    [0]func()
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
	if err := c.get(ctx, "/point", values("lang", opts.Lang, "lat", f(lat), "lon", f(lon), "deep", deep), nil, out); err != nil {
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
	ShortcodeSource string
	// Lang selects translated display names for this request.
	Lang string
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
	if err := c.get(ctx, "/emoji/"+seg(emoji), values("lang", opts.Lang, "deep", deep, "shortcode_source", opts.ShortcodeSource), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// EmojiSearchOptions configures EmojiSearch. Omitted fields use API defaults.
type EmojiSearchOptions struct {
	Category         string
	Subcategory      string
	Version          string
	MaxVersion       string
	IncludeSkinTones *bool
	Offset           *int
	// Lang selects translated display names for this request.
	Lang  string
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
	params := emojiQuery(opts)
	params.Set("q", query)
	out := &EmojiSearch{}
	if err := c.get(ctx, "/emoji", params, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// EmojiBrowseOptions configures catalog browsing. Omitted fields use API defaults.
type EmojiBrowseOptions = EmojiSearchOptions

// EmojiBrowse lists the catalog without a search query.
func (c *Client) EmojiBrowse(ctx context.Context, options ...EmojiBrowseOptions) (*EmojiBrowse, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &EmojiBrowse{}
	if err := c.get(ctx, "/emoji", emojiQuery(opts), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// EmojiCatalogOptions configures the complete JSON export. Deep defaults true.
type EmojiCatalogOptions struct {
	_                [0]func()
	Profile          string
	Deep             *bool
	Lang             string
	Category         string
	Subcategory      string
	Version          string
	MaxVersion       string
	IncludeSkinTones *bool
}

func (c *Client) EmojiCatalog(ctx context.Context, options ...EmojiCatalogOptions) (*EmojiCatalog, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	params := emojiQuery(EmojiSearchOptions{Lang: opts.Lang, Category: opts.Category, Subcategory: opts.Subcategory, Version: opts.Version, MaxVersion: opts.MaxVersion, IncludeSkinTones: opts.IncludeSkinTones})
	deep := opts.Deep == nil || *opts.Deep
	params.Set("deep", strconv.FormatBool(deep))
	if opts.Profile != "" {
		params.Set("profile", opts.Profile)
	}
	out := &EmojiCatalog{}
	if err := c.get(ctx, "/emoji/catalog", params, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) EmojiConvert(ctx context.Context, text, direction, shortcodeSource string) (*EmojiConversion, error) {
	body, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return nil, err
	}
	out := &EmojiConversion{}
	if err := c.request(ctx, http.MethodPost, "/emoji/convert", values("direction", direction, "shortcode_source", shortcodeSource), nil, body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// EmojiMetadata retrieves the corresponding catalog resource.

func (c *Client) EmojiMetadata(ctx context.Context) (*EmojiMetadata, error) {
	out := &EmojiMetadata{}
	if err := c.get(ctx, "/emoji/metadata", nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// EmojiCategories retrieves the corresponding catalog resource.

func (c *Client) EmojiCategories(ctx context.Context) (*EmojiCategories, error) {
	out := &EmojiCategories{}
	if err := c.get(ctx, "/emoji/categories", nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

type EmojiParseOptions struct {
	_    [0]func()
	Deep bool
	Lang string
}

// EmojiParse extracts complete sequences. Positions are UTF-16 code units; end is exclusive.
func (c *Client) EmojiParse(ctx context.Context, text string, options ...EmojiParseOptions) (*EmojiParse, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	body, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return nil, err
	}
	out := &EmojiParse{}
	if err := c.request(ctx, http.MethodPost, "/emoji/parse", values("deep", deep, "lang", opts.Lang), nil, body, out); err != nil {
		return nil, err
	}
	return out, nil
}

func emojiQuery(opts EmojiSearchOptions) url.Values {
	params := values("lang", opts.Lang, "category", opts.Category, "subcategory", opts.Subcategory, "version", opts.Version, "max_version", opts.MaxVersion)
	if opts.Deep {
		params.Set("deep", "true")
	}
	if opts.Limit != 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset != nil {
		params.Set("offset", strconv.Itoa(*opts.Offset))
	}
	if opts.IncludeSkinTones != nil {
		params.Set("include_skin_tones", strconv.FormatBool(*opts.IncludeSkinTones))
	}
	return params
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
	// Lang selects translated display names for this request.
	Lang    string
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
	if err := c.get(ctx, "/company/"+seg(number), values("lang", opts.Lang, "country", opts.Country, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}


type IndustryOptions = NAICSOptions
type IndustrySearchOptions = NAICSSearchOptions

// NAICS is the compatibility name for Industry.
func (c *Client) NAICS(ctx context.Context, code string, options ...NAICSOptions) (*NAICS, error) {
	return c.Industry(ctx, code, options...)
}

// NAICSSearch is the compatibility name for IndustrySearch.
func (c *Client) NAICSSearch(ctx context.Context, query string, options ...NAICSSearchOptions) (*NAICSSearch, error) {
	return c.IndustrySearch(ctx, query, options...)
}






type BINOptions struct {
	_    [0]func()
	Deep bool
}





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

// CityListOptions configures collection browsing. Cursor values are opaque.
type CityListOptions struct {
	Romanized        bool
	_                [0]func()
	Country          string
	State            string
	Limit            int
	Deep             bool
	Lang             string
	Type             string
	Grain            string
	Hierarchy        bool
	MinPopulation    *int64
	MaxPopulation    *int64
	Timezone         string
	BBox             string
	Sort             string
	Order            string
	Cursor           string
	Countries        []string
	ExcludeCountries []string
	Timezones        []string
	MinElevation     *float64
	MaxElevation     *float64
}

// CityList browses /city. The response always includes pagination metadata.
func (c *Client) CityList(ctx context.Context, options ...CityListOptions) (*CityList, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	limit := ""
	if opts.Limit != 0 {
		limit = strconv.Itoa(opts.Limit)
	}
	query := values("browse", "true", "country", opts.Country, "state", opts.State, "limit", limit, "deep", cityFlag(opts.Deep), "lang", opts.Lang, "type", opts.Type, "grain", opts.Grain, "hierarchy", cityFlag(opts.Hierarchy))
	cityCollectionQuery(query, opts.MinPopulation, opts.MaxPopulation, opts.Timezone, opts.BBox, opts.Sort, opts.Order, false, opts.Cursor)
	cityExtraQuery(query, opts.Countries, opts.ExcludeCountries, opts.Timezones, opts.MinElevation, opts.MaxElevation)
	cityResponseQuery(query, "list", opts.Romanized)
	out := &CityList{}
	if err := c.get(ctx, "/city", query, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CityNearbyCoordsOptions filters cities around a coordinate. BBox uses west,south,east,north.
type CityNearbyCoordsOptions struct {
	Romanized        bool
	_                [0]func()
	Country          string
	State            string
	Radius           float64
	Unit             string
	Limit            int
	Deep             bool
	Lang             string
	Type             string
	Grain            string
	Hierarchy        bool
	MinPopulation    *int64
	MaxPopulation    *int64
	Timezone         string
	BBox             string
	Sort             string
	Order            string
	Cursor           string
	Countries        []string
	ExcludeCountries []string
	Timezones        []string
	MinElevation     *float64
	MaxElevation     *float64
}

// CityNearbyCoords lists cities around a coordinate and returns opaque pagination metadata.
func (c *Client) CityNearbyCoords(ctx context.Context, lat float64, lon float64, options ...CityNearbyCoordsOptions) (*CityNearbyCoords, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	radius, limit := "", ""
	if opts.Radius != 0 {
		radius = f(opts.Radius)
	}
	if opts.Limit != 0 {
		limit = strconv.Itoa(opts.Limit)
	}
	query := values("nearby", "true", "lat", f(lat), "lon", f(lon), "radius", radius, "unit", opts.Unit, "country", opts.Country, "state", opts.State, "limit", limit, "deep", cityFlag(opts.Deep), "lang", opts.Lang, "type", opts.Type, "grain", opts.Grain, "hierarchy", cityFlag(opts.Hierarchy))
	cityCollectionQuery(query, opts.MinPopulation, opts.MaxPopulation, opts.Timezone, opts.BBox, opts.Sort, opts.Order, false, opts.Cursor)
	cityExtraQuery(query, opts.Countries, opts.ExcludeCountries, opts.Timezones, opts.MinElevation, opts.MaxElevation)
	cityResponseQuery(query, "list", opts.Romanized)
	out := &CityNearbyCoords{}
	if err := c.get(ctx, "/city", query, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func cityCollectionQuery(query url.Values, minPopulation *int64, maxPopulation *int64, timezone, bbox, sort, order string, page bool, cursor string) {
	if minPopulation != nil {
		query.Set("min_population", strconv.FormatInt(*minPopulation, 10))
	}
	if maxPopulation != nil {
		query.Set("max_population", strconv.FormatInt(*maxPopulation, 10))
	}
	for key, value := range values("timezone", timezone, "bbox", bbox, "sort", sort, "order", order, "page", cityFlag(page), "cursor", cursor) {
		query[key] = value
	}
}

// cityFlag omits default options so existing City requests retain their URLs.
func cityFlag(value bool) string {
	if value {
		return "true"
	}
	return ""
}

type CityResourceOptions struct {
	Romanized bool
	// Lang selects translated display names for this request.
	Lang      string
	_         [0]func()
	Country   string
	State     string
	Deep      bool
	Type      string
	Grain     string
	Hierarchy bool
	Aliases   bool
	Ambiguity string
}

type CityIDResourceOptions struct {
	Romanized bool
	// Lang selects translated display names for this request.
	Lang      string
	_         [0]func()
	Deep      bool
	Hierarchy bool
}

type CityNearestResourceOptions struct {
	Romanized bool
	// Lang selects translated display names for this request.
	Lang      string
	_         [0]func()
	Deep      bool
	Country   string
	State     string
	Type      string
	Grain     string
	Hierarchy bool
}

type CitySearchListOptions struct {
	Romanized bool
	// Lang selects translated display names for this request.
	Lang             string
	_                [0]func()
	Country          string
	State            string
	Limit            int
	Deep             bool
	Type             string
	Grain            string
	Hierarchy        bool
	Aliases          bool
	MinPopulation    *int64
	MaxPopulation    *int64
	Timezone         string
	BBox             string
	Sort             string
	Order            string
	Page             bool
	Cursor           string
	Countries        []string
	ExcludeCountries []string
	Timezones        []string
	MinElevation     *float64
	MaxElevation     *float64
	Match            string
	MaxDistance      int
}

type CityNearbyListOptions struct {
	Romanized bool
	// Lang selects translated display names for this request.
	Lang             string
	_                [0]func()
	Country          string
	State            string
	Radius           float64
	Unit             string
	Limit            int
	Deep             bool
	Type             string
	Grain            string
	Hierarchy        bool
	Aliases          bool
	Ambiguity        string
	MinPopulation    *int64
	MaxPopulation    *int64
	Timezone         string
	BBox             string
	Sort             string
	Order            string
	Page             bool
	Cursor           string
	Countries        []string
	ExcludeCountries []string
	Timezones        []string
	MinElevation     *float64
	MaxElevation     *float64
}

func (c *Client) CityResource(ctx context.Context, name string, options ...CityResourceOptions) (*CityResource, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &CityResource{}
	if err := c.get(ctx, "/city/"+seg(name), cityResponseQuery(values("lang", opts.Lang, "country", opts.Country, "state", opts.State, "deep", deep, "type", opts.Type, "grain", opts.Grain, "hierarchy", cityFlag(opts.Hierarchy), "aliases", cityFlag(opts.Aliases), "ambiguity", opts.Ambiguity), "resource", opts.Romanized), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) CityIDResource(ctx context.Context, id string, options ...CityIDResourceOptions) (*CityResource, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &CityResource{}
	if err := c.get(ctx, "/city/id/"+seg(id), cityResponseQuery(values("lang", opts.Lang, "deep", deep, "hierarchy", cityFlag(opts.Hierarchy)), "resource", opts.Romanized), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) CityNearestResource(ctx context.Context, lat float64, lon float64, options ...CityNearestResourceOptions) (*CityResource, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &CityResource{}
	if err := c.get(ctx, "/city", cityResponseQuery(values("lang", opts.Lang, "lat", f(lat), "lon", f(lon), "deep", deep, "country", opts.Country, "state", opts.State, "type", opts.Type, "grain", opts.Grain, "hierarchy", cityFlag(opts.Hierarchy)), "resource", opts.Romanized), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) CitySearchList(ctx context.Context, query string, options ...CitySearchListOptions) (*CityList, error) {
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
	out := &CityList{}
	params := values("lang", opts.Lang, "q", query, "country", opts.Country, "state", opts.State, "limit", limit, "deep", deep, "type", opts.Type, "grain", opts.Grain, "hierarchy", cityFlag(opts.Hierarchy), "aliases", cityFlag(opts.Aliases))
	cityCollectionQuery(params, opts.MinPopulation, opts.MaxPopulation, opts.Timezone, opts.BBox, opts.Sort, opts.Order, opts.Page, opts.Cursor)
	cityExtraQuery(params, opts.Countries, opts.ExcludeCountries, opts.Timezones, opts.MinElevation, opts.MaxElevation)
	cityResponseQuery(params, "list", opts.Romanized)
	if opts.Match != "" {
		params.Set("match", opts.Match)
	}
	if opts.MaxDistance != 0 {
		params.Set("max_distance", strconv.Itoa(opts.MaxDistance))
	}
	endpoint := "/city"
	if opts.Match != "" || opts.MaxDistance != 0 {
		endpoint = "/city/search"
	}
	if err := c.get(ctx, endpoint, params, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) CityNearbyList(ctx context.Context, name string, options ...CityNearbyListOptions) (*CityList, error) {
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
	out := &CityList{}
	query := values("lang", opts.Lang, "country", opts.Country, "state", opts.State, "radius", radius, "unit", opts.Unit, "limit", limit, "deep", deep, "type", opts.Type, "grain", opts.Grain, "hierarchy", cityFlag(opts.Hierarchy), "aliases", cityFlag(opts.Aliases), "ambiguity", opts.Ambiguity)
	cityCollectionQuery(query, opts.MinPopulation, opts.MaxPopulation, opts.Timezone, opts.BBox, opts.Sort, opts.Order, opts.Page, opts.Cursor)
	cityExtraQuery(query, opts.Countries, opts.ExcludeCountries, opts.Timezones, opts.MinElevation, opts.MaxElevation)
	cityResponseQuery(query, "list", opts.Romanized)
	if err := c.get(ctx, "/city/"+seg(name)+"/nearby", query, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CityChildren returns reviewed direct children. Coverage is partial.
func (c *Client) CityChildren(ctx context.Context, id string, options ...CityRelationsOptions) (*CityRelations, error) {
	return c.cityRelations(ctx, id, "children", options)
}

func (c *Client) CitySiblings(ctx context.Context, id string, options ...CityRelationsOptions) (*CityRelations, error) {
	return c.cityRelations(ctx, id, "siblings", options)
}

func (c *Client) cityRelations(ctx context.Context, id, relation string, options []CityRelationsOptions) (*CityRelations, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	q := values("cursor", opts.Cursor, "hierarchy", cityFlag(opts.Hierarchy), "deep", cityFlag(opts.Deep), "lang", opts.Lang)
	if opts.Limit != 0 {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}
	out := &CityRelations{}
	if opts.Romanized {
		q.Set("romanized", "true")
	}
	if err := c.get(ctx, "/city/id/"+seg(id)+"/"+relation, q, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) CityAncestors(ctx context.Context, id string, options ...CityAncestorsOptions) (*CityRelations, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &CityRelations{}
	if err := c.get(ctx, "/city/id/"+seg(id)+"/ancestors", values("hierarchy", cityFlag(opts.Hierarchy), "deep", cityFlag(opts.Deep), "lang", opts.Lang, "romanized", cityFlag(opts.Romanized)), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CityDistance compares stored centers by great-circle distance. Unknown coordinates remain nil.
func (c *Client) CityDistance(ctx context.Context, from, to string) (*CityDistance, error) {
	out := &CityDistance{}
	if err := c.get(ctx, "/city/id/"+seg(from)+"/distance/"+seg(to), nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

type CityRelationsOptions struct {
	Romanized bool
	_         [0]func()
	Limit     int
	Cursor    string
	Hierarchy bool
	Deep      bool
	Lang      string
}

type CityAncestorsOptions struct {
	Romanized bool
	_         [0]func()
	Hierarchy bool
	Deep      bool
	Lang      string
}

func cityResponseQuery(q url.Values, format string, romanized ...bool) url.Values {
	if len(romanized) > 0 && romanized[0] {
		q.Set("romanized", "true")
	}
	q.Set("format", format)
	return q
}

func cityExtraQuery(q url.Values, countries, excluded, timezones []string, minElevation, maxElevation *float64) {
	for k, v := range map[string][]string{"countries": countries, "exclude_countries": excluded, "timezones": timezones} {
		if v != nil {
			q.Set(k, strings.Join(v, ","))
		}
	}
	if minElevation != nil {
		q.Set("min_elevation", f(*minElevation))
	}
	if maxElevation != nil {
		q.Set("max_elevation", f(*maxElevation))
	}
}

func cityErrorCandidates(rows []json.RawMessage) []CityAmbiguityCandidate {
	out := make([]CityAmbiguityCandidate, 0, 2)
	if len(rows) > 2 {
		rows = rows[:2]
	}
	for _, raw := range rows {
		var row map[string]json.RawMessage
		if json.Unmarshal(raw, &row) != nil {
			continue
		}
		var id, object string
		if json.Unmarshal(row["id"], &id) != nil || len(id) != 17 || !strings.HasPrefix(id, "city_") || json.Unmarshal(row["object"], &object) != nil || object != "city" {
			continue
		}
		valid := true
		for _, ch := range id[5:] {
			if !(ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9') {
				valid = false
			}
		}
		if !valid {
			continue
		}
		safe := map[string]json.RawMessage{"id": row["id"], "object": row["object"]}
		for _, key := range []string{"name", "name_local", "country", "country_name", "state", "state_name", "type"} {
			if value, ok := row[key]; ok {
				var text *string
				if json.Unmarshal(value, &text) == nil && (text == nil || len([]rune(*text)) <= 255) {
					safe[key] = value
				}
			}
		}
		data, _ := json.Marshal(safe)
		var candidate CityAmbiguityCandidate
		if json.Unmarshal(data, &candidate) == nil {
			out = append(out, candidate)
		}
	}
	return out
}
