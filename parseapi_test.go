package parseapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type capture struct {
	path   string
	rawQry string
	header http.Header
}

func newTestClient(t *testing.T, handler http.HandlerFunc, opts ...Option) (*Client, *capture) {
	t.Helper()
	captured := &capture{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured.path = r.URL.EscapedPath()
		captured.rawQry = r.URL.RawQuery
		captured.header = r.Header.Clone()
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	client, err := New("test_key_123", append([]Option{WithBaseURL(server.URL), WithRetries(0)}, opts...)...)
	if err != nil {
		t.Fatal(err)
	}
	return client, captured
}

func okJSON(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}
}

func TestURLMapping(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name      string
		invoke    func(*Client) error
		wantPath  string
		wantQuery string
	}{
		{"ip", func(c *Client) error { _, err := c.IP(ctx, "8.8.8.8", nil); return err }, "/ip/8.8.8.8", ""},
		{"ip me", func(c *Client) error { _, err := c.IPMe(ctx, nil); return err }, "/ip/me", ""},
		{"ip deep", func(c *Client) error { _, err := c.IP(ctx, "8.8.8.8", &DeepOptions{Deep: true}); return err }, "/ip/8.8.8.8", "deep=true"},
		{"continent", func(c *Client) error { _, err := c.Continent(ctx, "NA"); return err }, "/continent/NA", ""},
		{"continent countries", func(c *Client) error { _, err := c.ContinentCountries(ctx, "NA"); return err }, "/continent/NA/countries", ""},
		{"country", func(c *Client) error { _, err := c.Country(ctx, "US"); return err }, "/country/US", ""},
		{"country states", func(c *Client) error { _, err := c.CountryStates(ctx, "US"); return err }, "/country/US/states", ""},
		{"state", func(c *Client) error { _, err := c.State(ctx, "NC", "US"); return err }, "/state/NC", "country=US"},
		{"state districts", func(c *Client) error { _, err := c.StateDistricts(ctx, "NC", "US"); return err }, "/state/NC/districts", "country=US"},
		{"district", func(c *Client) error { _, err := c.District(ctx, "37081", nil); return err }, "/district/37081", ""},
		{"city", func(c *Client) error { _, err := c.City(ctx, "charlotte", &CityOptions{State: "NC"}); return err }, "/city/charlotte", "state=NC"},
		{"city id", func(c *Client) error { _, err := c.CityID(ctx, "city_mb8mbqrkz8zb"); return err }, "/city/id/city_mb8mbqrkz8zb", ""},
		{"city search", func(c *Client) error {
			_, err := c.CitySearch(ctx, "char", &CitySearchOptions{Country: "US", Limit: 10})
			return err
		}, "/city", "country=US&limit=10&q=char"},
		{"city nearest", func(c *Client) error { _, err := c.CityNearest(ctx, 35.2271, -80.8431); return err }, "/city", "lat=35.2271&lon=-80.8431"},
		{"postal", func(c *Client) error { _, err := c.Postal(ctx, "28202", "US"); return err }, "/postal/28202", "country=US"},
		{"postal nearby", func(c *Client) error {
			_, err := c.PostalNearby(ctx, "28202", "US", &PostalNearbyOptions{Radius: 40, Unit: "km"})
			return err
		}, "/postal/28202/nearby", "country=US&radius=40&unit=km"},
		{"postal distance", func(c *Client) error { _, err := c.PostalDistance(ctx, "28202", "10001", "US"); return err }, "/postal/28202/distance/10001", "country=US"},
		{"email", func(c *Client) error { _, err := c.Email(ctx, "a@b.com", nil); return err }, "/email/a@b.com", ""},
		{"phone", func(c *Client) error { _, err := c.Phone(ctx, "+14155552671", &PhoneOptions{Deep: true}); return err }, "/phone/+14155552671", "deep=true"},
		{"domain", func(c *Client) error { _, err := c.Domain(ctx, "example.com", nil); return err }, "/domain/example.com", ""},
		{"mx", func(c *Client) error { _, err := c.MX(ctx, "example.com"); return err }, "/mx/example.com", ""},
		{"useragent", func(c *Client) error { _, err := c.Useragent(ctx, "TestUA/1.0", nil); return err }, "/useragent", ""},
		{"currency", func(c *Client) error { _, err := c.Currency(ctx, "USD"); return err }, "/currency/USD", ""},
		{"currency rate", func(c *Client) error { _, err := c.CurrencyRate(ctx, "USD", "EUR"); return err }, "/currency/USD/EUR", ""},
		{"language", func(c *Client) error { _, err := c.Language(ctx, "en"); return err }, "/language/en", ""},
		{"timezone encodes slash", func(c *Client) error { _, err := c.Timezone(ctx, "America/New_York", nil); return err }, "/timezone/America%2FNew_York", ""},
		{"holiday", func(c *Client) error { _, err := c.Holiday(ctx, "US", &HolidayOptions{Year: 1955}); return err }, "/holiday/US", "year=1955"},
		{"holiday date", func(c *Client) error { _, err := c.HolidayDate(ctx, "US", "2026-12-25"); return err }, "/holiday/US/2026-12-25", ""},
		{"elevation", func(c *Client) error { _, err := c.Elevation(ctx, 35.2, -80.8); return err }, "/elevation", "lat=35.2&lon=-80.8"},
		{"point deep", func(c *Client) error { _, err := c.Point(ctx, 36.0726, -79.792, &DeepOptions{Deep: true}); return err }, "/point", "deep=true&lat=36.0726&lon=-79.792"},
		{"weather", func(c *Client) error {
			_, err := c.Weather(ctx, 40.7128, -74.006, &DeepOptions{Deep: true})
			return err
		}, "/weather", "deep=true&lat=40.7128&lon=-74.006"},
		{"emoji", func(c *Client) error { _, err := c.Emoji(ctx, "rocket"); return err }, "/emoji/rocket", ""},
		{"emoji search", func(c *Client) error { _, err := c.EmojiSearch(ctx, "fire", &EmojiSearchOptions{Limit: 20}); return err }, "/emoji", "limit=20&q=fire"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client, captured := newTestClient(t, okJSON(`{}`))
			if err := tc.invoke(client); err != nil {
				t.Fatal(err)
			}
			if captured.path != tc.wantPath {
				t.Errorf("path = %q, want %q", captured.path, tc.wantPath)
			}
			if captured.rawQry != tc.wantQuery {
				t.Errorf("query = %q, want %q", captured.rawQry, tc.wantQuery)
			}
		})
	}
}

func TestHeaders(t *testing.T) {
	client, captured := newTestClient(t, okJSON(`{}`))
	if _, err := client.Country(context.Background(), "US"); err != nil {
		t.Fatal(err)
	}
	if got := captured.header.Get("X-API-Key"); got != "test_key_123" {
		t.Errorf("X-API-Key = %q", got)
	}
	if got := captured.header.Get("User-Agent"); got != "parseapi-go/"+version {
		t.Errorf("User-Agent = %q", got)
	}
}

func TestUseragentHeaderOverride(t *testing.T) {
	client, captured := newTestClient(t, okJSON(`{}`))
	if _, err := client.Useragent(context.Background(), "Mozilla/5.0 (Test)", nil); err != nil {
		t.Fatal(err)
	}
	if got := captured.header.Get("User-Agent"); got != "Mozilla/5.0 (Test)" {
		t.Errorf("User-Agent = %q", got)
	}
}

func TestMissingKey(t *testing.T) {
	t.Setenv("PARSEAPI_KEY", "")
	if _, err := New(""); err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestEnvKey(t *testing.T) {
	t.Setenv("PARSEAPI_KEY", "env_key_456")
	client, err := New("")
	if err != nil {
		t.Fatal(err)
	}
	if client.apiKey != "env_key_456" {
		t.Errorf("apiKey = %q", client.apiKey)
	}
}

func TestErrorShape(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":       "not_found",
			"message":    "City not found",
			"docs":       "https://parseapi.com/docs#not_found",
			"request_id": "req_abc",
		})
	})
	_, err := client.City(context.Background(), "notarealcityxyz", nil)
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if apiErr.Status != 404 || apiErr.Code != "not_found" || apiErr.Message != "City not found" ||
		apiErr.Docs != "https://parseapi.com/docs#not_found" || apiErr.RequestID != "req_abc" {
		t.Errorf("unexpected error fields: %+v", apiErr)
	}
}

func TestNonJSONErrorBody(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("gateway timeout"))
	})
	_, err := client.Country(context.Background(), "US")
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if apiErr.Code != "unknown_error" {
		t.Errorf("Code = %q", apiErr.Code)
	}
}

func TestRetryThenSuccess(t *testing.T) {
	attempts := 0
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"code":"server_error","message":"boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"country":"us"}`))
	}, WithRetries(2))
	result, err := client.Country(context.Background(), "US")
	if err != nil {
		t.Fatal(err)
	}
	if result.Country != "us" {
		t.Errorf("Country = %q", result.Country)
	}
	if attempts != 2 {
		t.Errorf("attempts = %d", attempts)
	}
}

func TestNoRetryOn404(t *testing.T) {
	attempts := 0
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":"not_found","message":"nope"}`))
	}, WithRetries(2))
	if _, err := client.Country(context.Background(), "XX"); err == nil {
		t.Fatal("expected error")
	}
	if attempts != 1 {
		t.Errorf("attempts = %d", attempts)
	}
}

func TestGivesUpAfterRetries(t *testing.T) {
	attempts := 0
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"code":"rate_limited","message":"slow down"}`))
	}, WithRetries(2))
	_, err := client.Country(context.Background(), "US")
	var apiErr *Error
	if !errors.As(err, &apiErr) || apiErr.Code != "rate_limited" {
		t.Fatalf("unexpected error: %v", err)
	}
	if attempts != 3 {
		t.Errorf("attempts = %d", attempts)
	}
}

func TestDeepTriadDecoding(t *testing.T) {
	// deep omitted -> nil, deep {} -> empty struct pointer, populated -> fields set
	client, _ := newTestClient(t, okJSON(`{"ip":"1.2.3.4"}`))
	result, err := client.IP(context.Background(), "1.2.3.4", nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.Deep != nil {
		t.Error("expected nil Deep when omitted")
	}

	client2, _ := newTestClient(t, okJSON(`{"ip":"1.2.3.4","deep":{}}`))
	result2, err := client2.IP(context.Background(), "1.2.3.4", &DeepOptions{Deep: true})
	if err != nil {
		t.Fatal(err)
	}
	if result2.Deep == nil {
		t.Error("expected non-nil Deep for {}")
	}
	if result2.Deep != nil && result2.Deep.Datacenter != nil {
		t.Error("expected nil Datacenter in empty deep")
	}
}

func TestMain(m *testing.M) {
	// Keep env overrides from leaking into URL assertions.
	os.Unsetenv("PARSEAPI_BASE_URL")
	os.Exit(m.Run())
}
