package parseapi

import (
	"encoding/json"
	"os"
	"testing"
)

func TestNetworkContractsCurrentAndFrozen(t *testing.T) {
	contents, err := os.ReadFile("testdata/network-contracts.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]json.RawMessage
	if err := json.Unmarshal(contents, &fixtures); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"ip_current", "ip_frozen", "ip_locked"} {
		var result IP
		if err := json.Unmarshal(fixtures[name], &result); err != nil {
			t.Fatal(err)
		}
		if result.Deep == nil {
			t.Fatal("requested Deep lost")
		}
		if name == "ip_current" {
			if result.Deep.District == nil || *result.Deep.District != "27053" || result.Deep.DistrictName == nil || *result.Deep.DistrictName != "Hennepin County" || result.Deep.CensusDivision == nil || *result.Deep.CensusDivision != "West North Central" || result.Deep.CensusRegion == nil || *result.Deep.CensusRegion != "Midwest" {
				t.Fatal("supported geography lost")
			}
		} else if result.Deep.District != nil || result.Deep.CensusRegion != nil {
			t.Fatal("unavailable geography invented")
		}
	}
	for _, name := range []string{"stack_success", "stack_unavailable", "stack_empty"} {
		var result Stack
		if err := json.Unmarshal(fixtures[name], &result); err != nil {
			t.Fatal(err)
		}
		var scope string = result.Scope
		var pages int = result.Pages
		if scope != "" || pages != 0 || result.Partial != nil {
			t.Fatal("legacy fields inferred")
		}
		if name == "stack_unavailable" {
			if result.CMS != nil || result.CheckedAt != nil {
				t.Fatal("unknown became empty")
			}
		} else if result.CMS == nil || result.Frameworks == nil {
			t.Fatal("checked inventory lost")
		}
		encoded, err := json.Marshal(result)
		if err != nil {
			t.Fatal(err)
		}
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(encoded, &fields); err != nil {
			t.Fatal(err)
		}
		// Retain the released Go model's JSON encoding for existing callers.
		for _, legacy := range []string{"scope", "pages", "partial"} {
			if _, present := fields[legacy]; !present {
				t.Fatal("published Go encoding changed")
			}
		}
	}
	var domain Domain
	if err := json.Unmarshal(fixtures["domain_unregistered"], &domain); err != nil || domain.Available == nil || !*domain.Available || domain.Deep == nil || domain.Deep.Registration == nil || domain.Deep.Registration.Registered {
		t.Fatalf("partial registration rejected: %v", err)
	}
}
