package parseapi

import (
	"context"
	"encoding/json"
	"testing"
)

func TestNameCountryAndNullableEvidence(t *testing.T) {
	client, call := newTestClient(t, okJSON(`{"name":"王","valid":true,"future":true,"deep":{"gender":null,"salutation":null}}`))
	result, err := client.Name(context.Background(), "王", NameOptions{Country: "CN", Deep: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.Deep == nil || result.Deep.Gender != nil || result.Deep.Salutation != nil {
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
	// Historical and future fields remain harmless when decoding older responses.
	if err := json.Unmarshal([]byte(`{"name":"Andrea","valid":true,"deep":{"known":true,"gender":null,"future":true}}`), &old); err != nil {
		t.Fatal(err)
	}
	if old.Deep == nil || old.Deep.Gender != nil || old.Deep.Salutation != nil {
		t.Fatal("lost nullable legacy evidence")
	}
}

func TestNameFormattingLocaleAndNullableResults(t *testing.T) {
	client, call := newTestClient(t, okJSON(`{"name":"Robert James Smith","deep":{"gender":"male","salutation":"Mr","short":"R.J. Smith","directory":"Smith, Robert James","initials":"RJS"}}`))
	result, err := client.Name(context.Background(), "Robert James Smith", NameOptions{Country: "US", Deep: true, NameLocale: "en-GB"})
	if err != nil {
		t.Fatal(err)
	}
	if call.path != "/name/Robert%20James%20Smith" || call.rawQry != "country=US&deep=true&name_locale=en-GB" {
		t.Fatalf("bad formatting request: %#v", call)
	}
	if result.Deep == nil || result.Deep.Short == nil || *result.Deep.Short != "R.J. Smith" || result.Deep.Directory == nil || *result.Deep.Directory != "Smith, Robert James" || result.Deep.Initials == nil || *result.Deep.Initials != "RJS" {
		t.Fatalf("lost formatting: %#v", result.Deep)
	}
	for _, body := range []string{`{"deep":{"short":null,"directory":null,"initials":null}}`, `{"deep":{"gender":null,"salutation":null}}`, `{"deep":{}}`} {
		var old Name
		if err := json.Unmarshal([]byte(body), &old); err != nil {
			t.Fatal(err)
		}
		if old.Deep == nil || old.Deep.Short != nil || old.Deep.Directory != nil || old.Deep.Initials != nil {
			t.Fatalf("invented formatting: %#v", old.Deep)
		}
	}
}
