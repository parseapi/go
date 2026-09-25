package parseapi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestTimeLocationCatalogAndSeasonalEvidence(t *testing.T) {
	calls := 0
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		q := r.URL.Query()
		if r.URL.Path == "/time/zones" {
			for k, v := range map[string]string{"country": "US", "area": "America", "offset": "+00:00", "abbreviation": "UTC", "dst": "false", "observes_dst": "false", "at": "1970-01-01T00:00:00Z", "details": "true", "sort": "offset"} {
				if q.Get(k) != v {
					t.Errorf("%s=%q", k, q.Get(k))
				}
			}
			_, _ = w.Write([]byte(`{"timezone_database_version": "2026c", "timezones": ["UTC"], "at": "1970-01-01T00:00:00.000Z", "zones": [{"timezone": "UTC", "countries": [], "area": null, "abbreviation": "UTC", "offset": "+00:00", "offset_seconds": 0, "dst": false, "observes_dst": false}]}`))
		} else {
			if q.Get("city") != "Springfield" || q.Get("country") != "US" || q.Get("state") != "IL" {
				t.Error("source query lost")
			}
			_, _ = w.Write([]byte(`{"timezone": null, "targets": null, "location": {"input": {"type": "city", "value": "Springfield"}, "status": "ambiguous", "candidates": [{"id": "city_a", "name": "Springfield", "country": "US", "state": "IL", "timezone": "America/Chicago", "latitude": 0, "longitude": 0}], "truncated": false, "source": "city_reference"}, "deep": {"standard_offset": "+01:00", "standard_offset_seconds": 3600, "dst_offset_seconds": -3600, "season": {"start": {"at": "2026-10-25T01:00:00Z", "before": {"offset_seconds": 3600, "dst": false}, "after": {"offset_seconds": 0, "dst": true}, "change_seconds": -3600}, "end": null}}}`))
		}
	})
	no := false
	zones, err := client.TimeZones(context.Background(), "", TimeZonesOptions{Country: "US", Area: "America", Offset: "+00:00", Abbreviation: "UTC", DST: &no, ObservesDST: &no, At: "1970-01-01T00:00:00Z", Details: true, Sort: "offset"})
	if err != nil || len(zones.Zones) != 1 || zones.Zones[0].OffsetSeconds != 0 || zones.Zones[0].DST {
		t.Fatalf("catalog: %v", err)
	}
	result, err := client.Time(context.Background(), "", TimeOptions{City: "Springfield", Country: "US", State: "IL", Targets: []string{"UTC"}, Deep: true})
	if err != nil || result.Timezone != nil || result.Location.Status != "ambiguous" || *result.Location.Candidates[0].Latitude != 0 || *result.Deep.DSTOffsetSeconds != -3600 || *result.Deep.Season.Start.ChangeSeconds != -3600 || *result.Deep.Season.Start.Before.DST {
		t.Fatalf("observations: %v %+v", err, result)
	}
	for _, opts := range []TimeOptions{{IP: "8.8.8.8", City: "Paris"}, {IP: "8.8.8.8", Country: "US"}, {State: "NY"}, {City: "Paris", State: "IDF"}, {Address: "a"}} {
		if _, err := client.Time(context.Background(), "", opts); err == nil {
			t.Error("invalid selection accepted")
		}
	}
	if _, err := client.Time(context.Background(), "UTC", TimeOptions{City: "Paris"}); err == nil {
		t.Error("mixed timezone source accepted")
	}
	if calls != 2 {
		t.Fatal("invalid sources dispatched")
	}
	for _, raw := range []string{`{}`, `{"deep":{}}`, `{"deep":{"season":null}}`} {
		var old Time
		if err := json.Unmarshal([]byte(raw), &old); err != nil {
			t.Fatal(err)
		}
		if old.Location != nil {
			t.Fatal("invented source")
		}
	}
}
