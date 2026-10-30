package parseapi

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"
)

func TestDomainResultUnknownPartialAndOneRequest(t *testing.T) {
	data, err := os.ReadFile("testdata/domain-result.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]json.RawMessage
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	for name, body := range fixtures {
		t.Run(name, func(t *testing.T) {
			calls := 0
			client, captured := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { calls++; okJSON(string(body))(w, r) })
			result, err := client.DomainResult(context.Background(), "bücher.example", DomainOptions{Deep: true})
			if err != nil {
				t.Fatal(err)
			}
			if calls != 1 || captured.path != "/domain/b%C3%BCcher.example" || captured.rawQry != "deep=true" || captured.header.Get("Parse-Version") != "2.0.0" || captured.header.Get("X-API-Key") != "test_key_123" {
				t.Fatal("lookup transport or request count changed")
			}
			if name == "unknown" || name == "unknown_deep" {
				if result.Available != nil {
					t.Fatal("unknown became a boolean")
				}
			} else if result.Available == nil {
				t.Fatal("known boolean lost")
			} else if *result.Available != (name == "available" || name == "unregistered_partial") {
				t.Fatal("known boolean changed")
			}
			if name == "available" || name == "registered" || name == "unknown" {
				if result.Deep != nil {
					t.Fatal("unrequested Deep added")
				}
			} else if result.Deep == nil {
				t.Fatal("requested Deep omitted")
			}
			if name == "locked" || name == "unknown_deep" || name == "null_registration" {
				if result.Deep.Registration != nil {
					t.Fatal("empty Deep invented registration")
				}
			}
			if name == "registered_false_empty" {
				r := result.Deep.Registration
				if r.DNSSEC == nil || *r.DNSSEC || r.Status == nil || len(r.Status) != 0 {
					t.Fatal("false/empty lost")
				}
			}
			if name == "missing_registered" || name == "null_registered" {
				if result.Deep.Registration == nil || result.Deep.Registration.Registered != nil {
					t.Fatal("unknown registered became false")
				}
			}
			if name == "unregistered_partial" {
				if result.Deep.Registration.Registered == nil || *result.Deep.Registration.Registered {
					t.Fatal("known false registered lost")
				}
			}
			if name == "registered_unknown_detail" || name == "unregistered_partial" {
				r := result.Deep.Registration
				if r.DNSSEC != nil || r.Status != nil {
					t.Fatal("partial detail invented a negative or empty list")
				}
			}
			if name == "registered_true_status" {
				r := result.Deep.Registration
				if r.DNSSEC == nil || !*r.DNSSEC || len(r.Status) != 1 || r.Status[0] != "client transfer prohibited" || r.Created == nil || *r.Created != "2020-01-02T03:04:05Z" {
					t.Fatal("known detail lost")
				}
			}
			if _, err := client.DomainResult(context.Background(), "example.com", DomainOptions{}, DomainOptions{}); err == nil || calls != 1 {
				t.Fatal("multiple options must fail before another request")
			}
		})
	}
}

func TestExistingDomainConsumerStillCompiles(t *testing.T) {
	client, _ := newTestClient(t, okJSON(`{"domain":"example.com","available":true}`))
	var oldCall func(context.Context, string, ...DomainOptions) (*Domain, error) = client.Domain
	old, err := oldCall(context.Background(), "example.com")
	if err != nil {
		t.Fatal(err)
	}
	var available bool = old.Available
	var dnssec bool = DomainRegistration{}.DNSSEC
	var status []string = DomainRegistration{}.Status
	if !available || dnssec || status != nil {
		t.Fatal("existing fields changed")
	}
}
