package parseapi

import (
	"context"
	"encoding/json"
	"testing"
)

func TestNameCountryAndNullableEvidence(t *testing.T) {
	client, call := newTestClient(t, okJSON(`{"name":"王","valid":true,"future":true,"gender":null,"salutation":null}`))
	result, err := client.Name(context.Background(), "王", NameOptions{Country: "CN", Deep: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.Gender != nil || result.Salutation != nil {
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
	if err := json.Unmarshal([]byte(`{"name":"Andrea","valid":true,"known":true,"gender":null,"future":true}`), &old); err != nil {
		t.Fatal(err)
	}
	if old.Gender != nil || old.Salutation != nil {
		t.Fatal("lost nullable legacy evidence")
	}
}

func TestNameFormattingLocaleAndNullableResults(t *testing.T) {
	client, call := newTestClient(t, okJSON(`{"name":"Robert James Smith","gender":"male","salutation":"Mr","short":"R.J. Smith","directory":"Smith, Robert James","initials":"RJS"}`))
	result, err := client.Name(context.Background(), "Robert James Smith", NameOptions{Country: "US", Deep: true, NameLocale: "en-GB"})
	if err != nil {
		t.Fatal(err)
	}
	if call.path != "/name/Robert%20James%20Smith" || call.rawQry != "country=US&deep=true&name_locale=en-GB" {
		t.Fatalf("bad formatting request: %#v", call)
	}
	if result.Short == nil || *result.Short != "R.J. Smith" || result.Directory == nil || *result.Directory != "Smith, Robert James" || result.Initials == nil || *result.Initials != "RJS" {
		t.Fatalf("lost formatting: %#v", result)
	}
	for _, body := range []string{`{"short":null,"directory":null,"initials":null}`, `{"gender":null,"salutation":null}`, `{}`} {
		var old Name
		if err := json.Unmarshal([]byte(body), &old); err != nil {
			t.Fatal(err)
		}
		if old.Short != nil || old.Directory != nil || old.Initials != nil {
			t.Fatalf("invented formatting: %#v", old)
		}
	}
}
