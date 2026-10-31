package parseapi

import (
	"encoding/json"
	"testing"
)

func TestCapitalCoordinatesWireCompatibility(t *testing.T) {
	for _, input := range []string{
		`{"capital":{"name":"Example","lat":0,"lon":null}}`,
		`{"capital":{"name":"Example","latitude":0,"longitude":null}}`,
		`{"capital":{"name":"Example","lat":99,"lon":99,"latitude":0,"longitude":null}}`,
		`{"capital":{"name":"Example","latitude":0,"longitude":null,"lat":99,"lon":99}}`,
	} {
		var country Country
		var state State
		if err := json.Unmarshal([]byte(input), &country); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal([]byte(input), &state); err != nil {
			t.Fatal(err)
		}
		if country.Capital.Latitude == nil || *country.Capital.Latitude != 0 || country.Capital.Longitude != nil || state.Capital.Latitude == nil || *state.Capital.Latitude != 0 || state.Capital.Longitude != nil {
			t.Fatalf("lost zero/null: %s", input)
		}
		for _, cap := range []any{country.Capital, state.Capital} {
			encoded, err := json.Marshal(cap)
			if err != nil {
				t.Fatal(err)
			}
			var raw map[string]any
			if err = json.Unmarshal(encoded, &raw); err != nil {
				t.Fatal(err)
			}
			if len(raw) != 3 || raw["latitude"] != float64(0) || raw["longitude"] != nil {
				t.Fatalf("noncanonical capital: %s", encoded)
			}
			if _, found := raw["lat"]; found {
				t.Fatal("old latitude emitted")
			}
			if _, found := raw["lon"]; found {
				t.Fatal("old longitude emitted")
			}
		}
	}
}
