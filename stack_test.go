package parseapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestStackInventoryShapesAndEncoding(t *testing.T) {
	var records []json.RawMessage
	if err := json.Unmarshal([]byte(`[{"domain":"xn--bcher-kva.example","url":"https://xn--bcher-kva.example/","checked_at":null,"scope":"homepage","pages":0,"partial":null,"cms":null,"servers":null,"frameworks":null,"ecommerce":null,"analytics":null,"chat":null,"payments":null,"hosting":null,"future":true},{"domain":"xn--bcher-kva.example","url":"https://xn--bcher-kva.example/","checked_at":"2026-09-21T12:00:00Z","scope":"site","pages":2,"partial":false,"cms":[],"servers":[],"frameworks":[],"ecommerce":[],"analytics":[],"chat":[],"payments":[],"hosting":[],"future":true,"deep":{}},{"domain":"xn--bcher-kva.example","url":"https://xn--bcher-kva.example/","checked_at":"2026-09-21T12:00:00Z","scope":"site","pages":3,"partial":true,"cms":[{"technology":"wordpress","name":"WordPress","version":"6.8"},{"technology":"ghost","name":"Ghost","version":null}],"servers":[{"technology":"nginx","name":"nginx","version":"1.26.2"},{"technology":"apache","name":"Apache","version":null}],"frameworks":[{"technology":"nextjs","name":"Next.js","version":null,"future":true},{"technology":"react","name":"React","version":"19.1"}],"ecommerce":[{"technology":"woocommerce","name":"WooCommerce","version":null}],"analytics":[{"technology":"google-analytics","name":"Google Analytics","version":null}],"chat":[{"technology":"intercom","name":"Intercom","version":null}],"payments":[{"technology":"stripe","name":"Stripe","version":null}],"hosting":[{"technology":"vercel","name":"Vercel","version":null}],"future":true},{"domain":"xn--bcher-kva.example","url":"https://xn--bcher-kva.example/","checked_at":null,"scope":"future-scope","pages":0,"partial":null,"cms":null,"servers":null,"frameworks":null,"ecommerce":null,"analytics":null,"chat":null,"payments":null,"hosting":null,"future":true,"deep":{}},{"domain":"xn--bcher-kva.example","url":"https://xn--bcher-kva.example/","checked_at":"2026-09-21T12:00:00Z","scope":"homepage","pages":1,"partial":true,"cms":[],"servers":[],"frameworks":[{"technology":"nextjs","name":"Next.js","version":null,"future":true}],"ecommerce":[],"analytics":[],"chat":[],"payments":[],"hosting":[],"future":true,"deep":{}}]`), &records); err != nil {
		t.Fatal(err)
	}
	for i, body := range records {
		client, captured := newTestClient(t, okJSON(string(body)))
		result, err := client.Stack(context.Background(), "bücher.example", StackOptions{Deep: true, Pretty: true})
		if err != nil {
			t.Fatal(err)
		}
		if captured.path != "/stack/b%C3%BCcher.example" || captured.rawQry != "deep=true&pretty=true" || captured.header.Get("Parse-Version") != "2.0.0" {
			t.Fatalf("wrong request: %#v", captured)
		}
		if i == 0 || i == 3 {
			if result.CMS != nil || result.Servers != nil || result.Pages != 0 || result.Partial != nil {
				t.Fatal("unknown inventory changed")
			}
		} else if i != 2 && (result.CMS == nil || len(result.CMS) != 0 || result.Servers == nil || len(result.Servers) != 0) {
			t.Fatal("empty inventory groups changed")
		}
		for _, group := range [][]StackTechnology{result.Ecommerce, result.Analytics, result.Chat, result.Payments, result.Hosting} {
			if i == 0 || i == 3 {
				if group != nil {
					t.Fatal("unknown collection changed")
				}
			} else if i == 2 {
				if len(group) != 1 || group[0].Technology == "" {
					t.Fatal("group entry lost")
				}
			} else if group == nil || len(group) != 0 {
				t.Fatal("empty collection changed")
			}
		}
		switch i {
		case 0:
			if result.Frameworks != nil || result.CheckedAt != nil || result.Deep != nil {
				t.Fatal("unknown or omitted values changed")
			}
		case 1:
			if result.Scope != "site" || result.Pages != 2 || result.Partial == nil || *result.Partial || result.Frameworks == nil || len(result.Frameworks) != 0 || result.Deep == nil || len(result.Deep) != 0 {
				t.Fatal("empty result or generic deep lost")
			}
		case 2:
			if result.Scope != "site" || result.Pages != 3 || result.Partial == nil || !*result.Partial || result.CheckedAt == nil || len(result.CMS) != 2 || len(result.Servers) != 2 || len(result.Frameworks) != 2 || result.Deep != nil {
				t.Fatal("site inventory changed")
			}
			if result.CMS[0].Technology != "wordpress" || result.CMS[0].Name != "WordPress" || result.CMS[0].Version == nil || *result.CMS[0].Version != "6.8" || result.CMS[1].Technology != "ghost" || result.CMS[1].Version != nil || result.Servers[0].Version == nil || *result.Servers[0].Version != "1.26.2" || result.Servers[1].Technology != "apache" || result.Frameworks[0].Version != nil {
				t.Fatal("inventory technology detail changed")
			}
		case 3:
			if result.Scope != "future-scope" || result.Frameworks != nil || result.Deep == nil || len(result.Deep) != 0 {
				t.Fatal("future/null fields changed")
			}
		case 4:
			if result.Scope != "homepage" || result.Pages != 1 || result.Partial == nil || !*result.Partial || result.Deep == nil || len(result.Deep) != 0 || result.Frameworks[0].Technology != "nextjs" || result.Frameworks[0].Version != nil {
				t.Fatal("generic deep changed core")
			}
		}
		if _, err := client.Stack(context.Background(), "example.com"); err != nil {
			t.Fatal(err)
		}
		if captured.rawQry != "" {
			t.Fatal("default options must be omitted")
		}
	}
}

type stackDeadlineTransport func(*http.Request) (*http.Response, error)

func (f stackDeadlineTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestStackDeadlineDefaultsAndExplicitSettings(t *testing.T) {
	for _, test := range []struct {
		name            string
		options         []Option
		stack, ordinary time.Duration
	}{
		{"default", nil, 35 * time.Second, 10 * time.Second},
		{"explicit default", []Option{WithTimeout(10 * time.Second)}, 10 * time.Second, 10 * time.Second},
		{"explicit short", []Option{WithTimeout(1200 * time.Millisecond)}, 1200 * time.Millisecond, 1200 * time.Millisecond},
		{"custom client", []Option{WithHTTPClient(&http.Client{Timeout: 8 * time.Second})}, 8 * time.Second, 8 * time.Second},
		{"custom unlimited", []Option{WithHTTPClient(&http.Client{})}, 0, 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			client, err := New("fixture", test.options...)
			if err != nil {
				t.Fatal(err)
			}
			var deadlines []time.Duration
			client.httpClient.Transport = stackDeadlineTransport(func(r *http.Request) (*http.Response, error) {
				remaining := time.Duration(0)
				if at, ok := r.Context().Deadline(); ok {
					remaining = time.Until(at)
				}
				deadlines = append(deadlines, remaining)
				return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"domain":"example.com","url":"https://example.com/","checked_at":null,"scope":"homepage","pages":0,"partial":null,"cms":null,"servers":null,"frameworks":null,"ecommerce":null,"analytics":null,"chat":null,"payments":null,"hosting":null,"deep":{}}`))}, nil
			})
			result, err := client.Stack(context.Background(), "example.com")
			if err != nil || result.Frameworks != nil || result.Deep == nil || len(result.Deep) != 0 {
				t.Fatalf("unavailable detail changed: %v %v", result, err)
			}
			if _, err := client.Domain(context.Background(), "example.com"); err != nil {
				t.Fatal(err)
			}
			if _, err := client.Stack(context.Background(), "example.com", StackOptions{Deep: true}); err != nil {
				t.Fatal(err)
			}
			for i, want := range []time.Duration{test.stack, test.ordinary, test.stack} {
				if deadlines[i] > want || deadlines[i] < want-500*time.Millisecond {
					t.Fatalf("deadline %d = %v; want %v", i, deadlines[i], want)
				}
			}
			if client.httpClient.Timeout != test.ordinary {
				t.Fatal("operation changed client timeout")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			if _, err := client.Stack(ctx, "example.com"); err != nil {
				t.Fatal(err)
			}
			if deadlines[3] > 100*time.Millisecond {
				t.Fatal("caller deadline was extended")
			}
		})
	}
}
