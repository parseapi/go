package parseapi

import (
	"context"
	"testing"
)

func TestNAICSHierarchyAndSearch(t *testing.T) {
	client, call := newTestClient(t, okJSON(`{"naics":"31-33","name":"Manufacturing","description":null,"level":2,"parent":null,"parent_name":null,"children":[{"naics":"311","name":"Food Manufacturing"}],"year":2022,"country":"US","future":true}`))
	result, err := client.NAICS(context.Background(), "31-33")
	if err != nil {
		t.Fatal(err)
	}
	if result.NAICS != "31-33" || result.Parent != nil || result.Description != nil || len(result.Children) != 1 {
		t.Fatalf("lost hierarchy: %#v", result)
	}
	if call.path != "/naics/31-33" || call.rawQry != "" {
		t.Fatalf("bad lookup: %#v", call)
	}
	_, _ = client.NAICS(context.Background(), "54/11")
	if call.path != "/naics/54%2F11" {
		t.Fatalf("bad encoding: %#v", call)
	}
	if _, err := client.NAICS(context.Background(), "54", NAICSOptions{}, NAICSOptions{}); err == nil {
		t.Fatal("accepted multiple options")
	}
	client, call = newTestClient(t, okJSON(`{"q":"coffee & tea","year":2022,"country":"US","results":[]}`))
	search, err := client.NAICSSearch(context.Background(), "coffee & tea", NAICSSearchOptions{Limit: 5})
	if err != nil {
		t.Fatal(err)
	}
	if search.Year != 2022 || search.Results == nil || len(search.Results) != 0 {
		t.Fatalf("bad search: %#v", search)
	}
	if call.path != "/naics" || call.rawQry != "limit=5&q=coffee+%26+tea" {
		t.Fatalf("bad search request: %#v", call)
	}
	_, _ = client.NAICSSearch(context.Background(), "plumbing")
	if call.rawQry != "q=plumbing" {
		t.Fatalf("unexpected defaults: %#v", call)
	}
}

func TestNAICSExclusionsAndMatchCompatibility(t *testing.T) {
	client, _ := newTestClient(t, okJSON(`{"q":"sofware","year":2022,"country":"US","results":[{"naics":"541511","name":"Custom Computer Programming Services","description":null,"level":6,"parent":"54151","parent_name":"Computer Systems Design and Related Services","children":[],"year":2022,"country":"US"},{"naics":"541511","name":"Custom Computer Programming Services","description":null,"level":6,"parent":"54151","parent_name":"Computer Systems Design and Related Services","children":[],"year":2022,"country":"US","exclusions":null,"match":null},{"naics":"541511","name":"Custom Computer Programming Services","description":null,"level":6,"parent":"54151","parent_name":"Computer Systems Design and Related Services","children":[],"year":2022,"country":"US","exclusions":[],"match":{"field":"future-field","text":"Future matching evidence","corrections":[],"future":true}},{"naics":"541511","name":"Custom Computer Programming Services","description":null,"level":6,"parent":"54151","parent_name":"Computer Systems Design and Related Services","children":[],"year":2022,"country":"US","exclusions":[{"description":"Designing integrated computer systems","codes":[{"naics":"541512","name":"Computer Systems Design Services"}]},{"description":"Activities classified elsewhere","codes":[]}],"match":{"field":"term","text":"Computer software programming services","corrections":[{"from":"sofware","to":"software"}]},"future":true}]}`))
	search, err := client.NAICSSearch(context.Background(), "sofware")
	if err != nil {
		t.Fatal(err)
	}
	var results []NAICS = search.Results // Preserve existing consumer collection type.
	if results[0].Exclusions != nil || results[0].Match != nil || results[1].Exclusions != nil || results[1].Match != nil {
		t.Fatal("older/null fields must remain unknown")
	}
	if results[2].Exclusions == nil || len(results[2].Exclusions) != 0 || results[2].Match.Field != "future-field" || len(results[2].Match.Corrections) != 0 {
		t.Fatal("lost empty arrays or future match field")
	}
	if results[3].Exclusions[0].Codes[0].NAICS != "541512" || results[3].Exclusions[1].Description != "Activities classified elsewhere" || len(results[3].Exclusions[1].Codes) != 0 {
		t.Fatal("lost exclusions without linked codes")
	}
	correction := results[3].Match.Corrections[0]
	if correction.From != "sofware" || correction.To != "software" || results[3].Match.Text != "Computer software programming services" {
		t.Fatal("lost actual match evidence")
	}
}
