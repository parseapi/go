package parseapi

import (
	"context"
	"encoding/json"
	"testing"
)

func TestNameCountryAndKnown(t *testing.T) {
	client, call := newTestClient(t, okJSON(`{"name":"王","valid":true,"known":true,"countries":["CN","TW"],"gender":null,"future":true}`))
	result, err := client.Name(context.Background(), "王", NameOptions{Country: "CN"})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Known || len(result.Countries) != 2 || result.Gender != nil {
		t.Fatalf("lost name facts: %#v", result)
	}
	if call.path != "/name/%E7%8E%8B" || call.rawQry != "country=CN" {
		t.Fatalf("bad country request: %#v", call)
	}
	if _, err := client.Name(context.Background(), "Andrea"); err != nil {
		t.Fatal(err)
	}
	if call.rawQry != "" {
		t.Fatal("country leaked into old call")
	}
	var old Name
	if err := json.Unmarshal([]byte(`{"name":"Andrea","valid":true,"gender":null}`), &old); err != nil {
		t.Fatal(err)
	}
	if old.Known || len(old.Countries) != 0 {
		t.Fatal("unexpected legacy membership")
	}
}
