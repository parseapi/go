package parseapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestAPIContractPinSurvivesRetriesAndSelfLookup(t *testing.T) {
	var requests []capture
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, capture{path: r.URL.Path, rawQry: r.URL.RawQuery, header: r.Header.Clone()})
		if len(requests) == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = fmt.Fprint(w, `{"code":"unavailable","message":"Try again"}`)
			return
		}
		_, _ = fmt.Fprint(w, `{"ip":"192.0.2.1","country":null,"deep":{"datacenter":null},"future":true}`)
	}, WithRetries(1))
	result, err := client.IPSelf(context.Background(), IPSelfOptions{Deep: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.IP != "192.0.2.1" || result.Country != nil || result.Deep == nil || result.Deep.Datacenter != nil {
		t.Fatalf("nullable response changed: %#v", result)
	}
	if len(requests) != 2 {
		t.Fatalf("attempts = %d, want 2", len(requests))
	}
	for _, request := range requests {
		if request.path != "/ip" || request.rawQry != "deep=true" {
			t.Errorf("request = %s?%s", request.path, request.rawQry)
		}
		if request.header.Get("Parse-Version") != "2.0.0" || request.header.Get("X-API-Key") != "test_key_123" || request.header.Get("User-Agent") != "parseapi-go/"+version {
			t.Errorf("request headers changed: %v", request.header)
		}
	}
}

func TestAPIContractPinWithUserAgentInput(t *testing.T) {
	client, captured := newTestClient(t, okJSON(`{"useragent":"Example/1.0","bot":false,"mobile":false,"deep":{}}`))
	if _, err := client.UserAgent(context.Background(), "Example/1.0"); err != nil {
		t.Fatal(err)
	}
	if captured.header.Get("Parse-Version") != "2.0.0" || captured.header.Get("User-Agent") != "Example/1.0" || captured.header.Get("X-API-Key") != "test_key_123" {
		t.Fatalf("request headers changed: %v", captured.header)
	}
}

func TestAPIVersionErrorsDoNotRetryOrFallBack(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusGone} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			attempts := 0
			client, captured := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				attempts++
				w.WriteHeader(status)
				_, _ = fmt.Fprint(w, `{"code":"invalid_request","message":"Unsupported API version","request_id":"req_version"}`)
			}, WithRetries(2))
			_, err := client.IPSelf(context.Background())
			var apiErr *Error
			if !errors.As(err, &apiErr) || apiErr.Status != status || apiErr.Code != "invalid_request" || apiErr.RequestID != "req_version" {
				t.Fatalf("unexpected API error: %v", err)
			}
			if attempts != 1 || captured.header.Get("Parse-Version") != "2.0.0" {
				t.Fatalf("version error fell back or retried: %d attempts", attempts)
			}
		})
	}
}
