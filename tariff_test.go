package parseapi

import (
	"context"
	"encoding/json"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestTariffSearchPreservesParentContextAndOlderResponses(t *testing.T) {
	for _, extra := range []string{"", `,"lineage":null`} {
		var result TariffSearch
		if err := json.Unmarshal([]byte(`{"q":"horses","revision":"fixture","lines":[{"hts":"0101.29.00.90","description":"Other","general":null`+extra+`}]}`), &result); err != nil {
			t.Fatal(err)
		}
		if result.Lines[0].Lineage != nil {
			t.Fatal("missing or null lineage became observed context")
		}
	}
	var result TariffSearch
	if err := json.Unmarshal([]byte(`{"q":"horses","revision":"fixture","lines":[{"hts":"0101.29.00.90","description":"Other","general":null,"lineage":["Live horses","Other horses"],"future":true},{"hts":"0101","description":"Live horses","general":null,"lineage":[]}]}`), &result); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result.Lines[0].Lineage, []string{"Live horses", "Other horses"}) {
		t.Fatal("parent descriptions were lost")
	}
	if result.Lines[1].Lineage == nil || len(result.Lines[1].Lineage) != 0 {
		t.Fatal("empty lineage became unknown")
	}
}

func TestTariffEditionDateRoundtrip(t *testing.T) {
	edition := strings.Repeat("a", 64)
	payload := `{"hts":"0101","edition":"` + edition + `","date":"2026-09-15","deep":{"effective_rate":null,"reason":"future_reason","measures":[]}}`
	client, capture := newTestClient(t, okJSON(payload))
	result, err := client.Tariff(context.Background(), "0101", TariffOptions{Deep: true, Origin: "CA", Edition: edition, Date: "2026-09-15"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Edition == nil || *result.Edition != edition || (result.Date == nil || *result.Date != "2026-09-15") || result.Deep.Reason == nil || *result.Deep.Reason != "future_reason" || result.Deep.EffectiveRate != nil {
		t.Fatal("edition/date/reason lost")
	}
	query, _ := url.ParseQuery(capture.rawQry)
	if query.Get("edition") != edition || query.Get("date") != "2026-09-15" || query.Get("origin") != "CA" || query.Get("deep") != "true" {
		t.Fatal(query)
	}
	search, err := client.TariffSearch(context.Background(), "horses", TariffSearchOptions{Edition: edition, Date: "2026-09-15"})
	if err != nil {
		t.Fatal(err)
	}
	if search.Edition == nil || *search.Edition != edition || (search.Date == nil || *search.Date != "2026-09-15") {
		t.Fatal("search edition lost")
	}
	query, _ = url.ParseQuery(capture.rawQry)
	if query.Get("q") != "horses" || query.Get("edition") != edition || query.Get("date") != "2026-09-15" {
		t.Fatal(query)
	}
	client, capture = newTestClient(t, okJSON(strings.Replace(payload, `"date":"2026-09-15"`, `"date":null`, 1)))
	_, err = client.Tariff(context.Background(), "0101", TariffOptions{Edition: edition})
	if err != nil {
		t.Fatal(err)
	}
	query, _ = url.ParseQuery(capture.rawQry)
	if query.Has("date") {
		t.Fatal("invented date")
	}
	var older Tariff
	if err := json.Unmarshal([]byte(`{"hts":"0101","deep":{}}`), &older); err != nil {
		t.Fatal(err)
	}
	if older.Edition != nil || older.Date != nil || older.Deep.Reason != nil {
		t.Fatal("invented older metadata")
	}
}

func TestTariffRejectsIgnoredSelection(t *testing.T) {
	client, _ := newTestClient(t, okJSON(`{"revision":"old"}`))
	_, err := client.Tariff(context.Background(), "0101", TariffOptions{Edition: strings.Repeat("a", 64)})
	if e, ok := err.(*Error); !ok || e.Code != "tariff_selection_mismatch" || e.Status != 0 {
		t.Fatal(err)
	}
	_, err = client.TariffSearch(context.Background(), "horses", TariffSearchOptions{Date: "2026-09-15"})
	if e, ok := err.(*Error); !ok || e.Code != "tariff_selection_mismatch" {
		t.Fatal(err)
	}
}

func TestTariffDateSelectionRejectsInvalidReturnedEdition(t *testing.T) {
	for _, edition := range []string{"legacy", "", strings.Repeat("A", 64), strings.Repeat("a", 64) + "\n"} {
		body, err := json.Marshal(map[string]string{"edition": edition, "date": "2026-09-15"})
		if err != nil {
			t.Fatal(err)
		}
		client, _ := newTestClient(t, okJSON(string(body)))
		_, err = client.Tariff(context.Background(), "0101", TariffOptions{Date: "2026-09-15"})
		if e, ok := err.(*Error); !ok || e.Code != "tariff_selection_mismatch" || e.Status != 0 {
			t.Fatalf("lookup accepted invalid edition %q: %v", edition, err)
		}
		_, err = client.TariffSearch(context.Background(), "horses", TariffSearchOptions{Date: "2026-09-15"})
		if e, ok := err.(*Error); !ok || e.Code != "tariff_selection_mismatch" || e.Status != 0 {
			t.Fatalf("search accepted invalid edition %q: %v", edition, err)
		}
	}
}
