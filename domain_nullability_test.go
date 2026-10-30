package parseapi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestOrdinaryDomainPreservesUnknownAndKnownStates(t *testing.T) {
	for _, available := range []any{nil, false, true} {
		for _, dnssec := range []any{"absent", nil, false, true} {
			registration := map[string]any{"registered": true}
			if dnssec != "absent" {
				registration["dnssec"] = dnssec
			}
			body := map[string]any{"domain": "example.com", "available": available}
			if available == true {
				body["deep"] = map[string]any{"registration": map[string]any{"registered": false}}
			} else if available == nil {
				body["deep"] = map[string]any{}
			} else {
				body["deep"] = map[string]any{"registration": registration}
			}
			encoded, err := json.Marshal(body)
			if err != nil {
				t.Fatal(err)
			}
			calls := 0
			client, captured := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { calls++; okJSON(string(encoded))(w, r) })
			result, err := client.Domain(context.Background(), "example.com", DomainOptions{Deep: true})
			if err != nil {
				t.Fatal(err)
			}
			if calls != 1 || captured.path != "/domain/example.com" || captured.rawQry != "deep=true" {
				t.Fatal("ordinary lookup request changed")
			}
			if available == nil {
				if result.Available != nil {
					t.Fatal("unknown became false")
				}
			} else if result.Available == nil || *result.Available != available {
				t.Fatal("known availability lost")
			}
			if available == false {
				got := result.Deep.Registration.DNSSEC
				if dnssec == nil || dnssec == "absent" {
					if got != nil {
						t.Fatal("missing DNSSEC became false")
					}
				} else if got == nil || *got != dnssec {
					t.Fatal("known DNSSEC lost")
				}
			}
		}
	}
}

func TestOrdinaryDomainDecodeErrorDoesNotRetry(t *testing.T) {
	calls := 0
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		okJSON(`{"domain":"example.com","available":"false"}`)(w, r)
	})
	if _, err := client.Domain(context.Background(), "example.com"); err == nil || calls != 1 {
		t.Fatal("malformed JSON type must fail once")
	}
}
