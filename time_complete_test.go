package parseapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestTimeTargetCollections(t *testing.T) {
	calls := 0
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { calls++; _, _ = w.Write([]byte(`{}`)) })
	for _, targets := range [][]string{{}, {""}, {"UTC,UTC"}, {"UTC", "UTC", "UTC", "UTC", "UTC", "UTC", "UTC", "UTC", "UTC", "UTC", "UTC"}} {
		if _, err := client.Time(context.Background(), "UTC", TimeOptions{Targets: targets}); err == nil {
			t.Fatal("expected local error")
		}
		if _, err := client.TimeAt(context.Background(), 0, 0, TimeAtOptions{Targets: targets}); err == nil {
			t.Fatal("expected local error")
		}
	}
	if _, err := client.Time(context.Background(), "UTC", TimeOptions{Targets: []string{"UTC"}, To: "UTC"}); err == nil {
		t.Fatal("expected conflict")
	}
	if calls != 0 {
		t.Fatal("invalid targets dispatched")
	}
	for _, raw := range []string{`{}`, `{"targets":null}`} {
		var result Time
		if err := json.Unmarshal([]byte(raw), &result); err != nil || result.Targets != nil {
			t.Fatalf("unknown targets: %v", err)
		}
	}
	var result Time
	if err := json.Unmarshal([]byte(`{"targets":[{"timezone":"UTC","unix":0},{"timezone":"UTC","unix":0}]}`), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.Targets) != 2 || result.Targets[0].Unix == nil || *result.Targets[0].Unix != 0 {
		t.Fatal("lost duplicate or zero")
	}
	var zones TimeZones
	if err := json.Unmarshal([]byte(`{"timezone_database_version":"2026c","timezones":[]}`), &zones); err != nil || zones.TimezoneDatabaseVersion != "2026c" || zones.Timezones == nil {
		t.Fatal("bad zone discovery")
	}
}

func TestTimeResolutionAndReservedSources(t *testing.T) {
	client, err := New("fixture", WithBaseURL("http://127.0.0.1:9"))
	if err != nil {
		t.Fatal(err)
	}
	for _, zone := range []string{"zones", "help", " ZONES ", "Help"} {
		if _, err := client.Time(context.Background(), zone); err == nil || !strings.Contains(err.Error(), "IANA timezone ID") {
			t.Fatalf("reserved source: %v", err)
		}
	}
	var result Time
	if err := json.Unmarshal([]byte(`{"deep":{"timezone_database_version":"2026c","resolution":{"kind":"gap","policy":"earlier","adjustment_seconds":-1800,"alternatives":[{"at":"1970-01-01T00:00:00.123+00:00","unix":0,"offset":"+00:00"},{"at":"1970-01-01T00:30:00.123+00:00","unix":1800,"offset":"+00:00"}],"future":true}}}`), &result); err != nil {
		t.Fatal(err)
	}
	if result.Deep == nil || *result.Deep.TimezoneDatabaseVersion != "2026c" || *result.Deep.Resolution.AdjustmentSeconds != -1800 || *result.Deep.Resolution.Alternatives[0].Unix != 0 {
		t.Fatalf("resolution not preserved: %+v", result.Deep)
	}
	for _, raw := range []string{`{"deep":{}}`, `{"deep":{"resolution":null}}`} {
		var result Time
		if err := json.Unmarshal([]byte(raw), &result); err != nil || result.Deep.Resolution != nil {
			t.Fatalf("unknown resolution: %v", err)
		}
	}
	if err := json.Unmarshal([]byte(`{"deep":{"resolution":{"kind":"unique","alternatives":[]}}}`), &result); err != nil || result.Deep.Resolution.Alternatives == nil {
		t.Fatal("known empty alternatives lost")
	}
}
