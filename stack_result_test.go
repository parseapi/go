package parseapi

import (
	"context"
	"encoding/json"
	"os"
	"testing"
)

func TestCurrentStackFromCapturedWire(t *testing.T) {
	body, err := os.ReadFile("testdata/current-stack.json")
	if err != nil {
		t.Fatal(err)
	}
	client, request := newTestClient(t, okJSON(string(body)))
	result, err := client.Stack(context.Background(), "stripe.com", StackOptions{Pretty: true})
	if err != nil {
		t.Fatal(err)
	}
	if request.path != "/stack/stripe.com" || request.rawQry != "pretty=true" || request.header.Get("Parse-Version") != "2.0.0" {
		t.Fatalf("unexpected request: %#v", request)
	}
	if result.CheckedAt == nil || result.CMS == nil || len(result.CMS) != 0 || len(result.Servers) != 1 || result.Servers[0].Version != nil {
		t.Fatalf("wire meaning changed: %#v", result)
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &fields); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"scope", "pages", "partial"} {
		if _, ok := fields[key]; ok {
			t.Fatalf("fabricated %s", key)
		}
	}
	result.CMS = nil
	result.Frameworks = []StackTechnology{}
	encoded, _ = json.Marshal(result)
	json.Unmarshal(encoded, &fields)
	if string(fields["cms"]) != "null" || string(fields["frameworks"]) != "[]" {
		t.Fatal("unknown and empty conflated")
	}
}
