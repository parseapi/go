package parseapi_test

import (
	"encoding/json"
	parseapi "github.com/parseapi/go"
	"testing"
)

func TestStackMarshalPreservesUnknownAndEmpty(t *testing.T) {
	value := parseapi.Stack{Domain: "example.com", URL: "https://example.com/",
		CMS: []parseapi.StackTechnology{}, Servers: nil, Frameworks: []parseapi.StackTechnology{{Technology: "react", Name: "React"}}}
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &fields); err != nil {
		t.Fatal(err)
	}
	for key, expected := range map[string]string{"cms": `[]`, "servers": `null`} {
		if string(fields[key]) != expected {
			t.Fatalf("%s = %s, want %s", key, fields[key], expected)
		}
	}
	var technologies []parseapi.StackTechnology
	if err := json.Unmarshal(fields["frameworks"], &technologies); err != nil || len(technologies) != 1 || technologies[0].Technology != "react" {
		t.Fatalf("populated group changed: %s %v", fields["frameworks"], err)
	}
}
