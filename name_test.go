package parseapi

import (
	"context"
	"encoding/json"
	"testing"
)

func TestNameCountryAndKnown(t *testing.T) {
	client, call := newTestClient(t, okJSON(`{"name":"王","valid":true,"future":true,"deep":{"known":true,"gender":null}}`))
	result, err := client.Name(context.Background(), "王", NameOptions{Country: "CN", Deep: true})
	if err != nil {
		t.Fatal(err)
	}
	if (result.Deep.Known == nil || !*result.Deep.Known) || result.Deep.Gender != nil {
		t.Fatalf("lost name facts: %#v", result)
	}
	if call.path != "/name/%E7%8E%8B" || call.rawQry != "country=CN&deep=true" {
		t.Fatalf("bad country request: %#v", call)
	}
	if _, err := client.Name(context.Background(), "Andrea"); err != nil {
		t.Fatal(err)
	}
	if call.rawQry != "" {
		t.Fatal("country leaked into old call")
	}
	var old Name
	if err := json.Unmarshal([]byte(`{"name":"Andrea","valid":true,"deep":{"gender":null}}`), &old); err != nil {
		t.Fatal(err)
	}
	if old.Deep.Known != nil {
		t.Fatal("unexpected legacy membership")
	}
}
