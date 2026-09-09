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
