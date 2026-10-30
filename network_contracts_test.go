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
		// Encode the actual response fields, keeping unknown values honest.
		for _, legacy := range []string{"scope", "pages", "partial"} {
			if _, present := fields[legacy]; present {
				t.Fatal("unsupported field encoded")
			}
		}
	}
	var state State
	if err := json.Unmarshal(fixtures["state_current"], &state); err != nil || len(state.Timezones) != 1 || state.Timezones[0] != "America/Los_Angeles" {
		t.Fatalf("state timezones lost: %v", err)
	}
	var district District
	if err := json.Unmarshal(fixtures["district_current"], &district); err != nil || len(district.Timezones) != 1 || district.Timezones[0] != "America/New_York" {
		t.Fatalf("district timezones lost: %v", err)
	}
	var domain Domain
	if err := json.Unmarshal(fixtures["domain_unregistered"], &domain); err != nil || !domain.Available || domain.Deep == nil || domain.Deep.Registration == nil || domain.Deep.Registration.Registered {
		t.Fatalf("partial registration rejected: %v", err)
	}
}
