package parseapi

import (
	"encoding/json"
	"os"
	"testing"
)

func TestPostalAdditiveModels(t *testing.T) {
	bytes, err := os.ReadFile("testdata/postal-additive.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases map[string]Postal
	if err = json.Unmarshal(bytes, &cases); err != nil {
		t.Fatal(err)
	}
	for name, postal := range cases {
		if postal.Deep == nil || len(postal.Deep.Metros) != 1 {
			t.Fatalf("%s: metro missing", name)
		}
		metro := postal.Deep.Metros[0]
		if metro.Share == nil || *metro.Share != 0.125 {
			t.Fatalf("%s: legacy share changed", name)
		}
		if name == "missing" || name == "null" {
			if postal.Type != nil || metro.SharePercent != nil {
				t.Fatalf("%s: unknown became a value", name)
			}
		}
	}
	if *cases["current"].Type != "standard" || *cases["current"].Deep.Metros[0].SharePercent != 12.5 {
		t.Fatal("current fields missing")
	}
	if *cases["zero"].Deep.Metros[0].SharePercent != 0 {
		t.Fatal("zero lost")
	}
	if *cases["future"].Type != "future_custom" || *cases["future"].Deep.Metros[0].SharePercent != 100 {
		t.Fatal("future type or endpoint percentage lost")
	}
	encoded, err := json.Marshal(cases["missing"])
	if err != nil {
		t.Fatal(err)
	}
	var roundtrip Postal
	if err = json.Unmarshal(encoded, &roundtrip); err != nil {
		t.Fatal(err)
	}
	if *roundtrip.Deep.Metros[0].Share != 0.125 {
		t.Fatal("legacy roundtrip changed")
	}
}
