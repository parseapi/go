package parseapi

import (
	"encoding/json"
	"testing"
)

func TestNPIPublishedDetails(t *testing.T) {
	var value Provider
	err := json.Unmarshal([]byte(`{"npi":"1881018208","valid":true,"sources":{"nppes":{"edition":"example","published_at":null,"through":"2026-09-20","imported_at":"2026-09-24T12:00:00.000Z"}},"deep":{"updated_at":"2026-09-18","taxonomies":[{"taxonomy":"207Q00000X","primary":true,"license":"000123"}]}}`), &value)
	if err != nil {
		t.Fatal(err)
	}
	if value.Sources.Nppes.PublishedAt != nil || *value.Deep.Taxonomies[0].License != "000123" || *value.Deep.UpdatedAt != "2026-09-18" {
		t.Fatalf("lost source details: %+v", value)
	}
	for _, pair := range []struct {
		payload string
		absent  bool
	}{{`{"deep":{"taxonomies":null}}`, true}, {`{"deep":{"taxonomies":[]}}`, false}} {
		if err := json.Unmarshal([]byte(pair.payload), &value); err != nil {
			t.Fatal(err)
		}
		if (value.Deep.Taxonomies == nil) != pair.absent {
			t.Fatal("null and empty conflated")
		}
	}
}
