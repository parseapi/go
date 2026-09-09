package parseapi

import (
	"encoding/json"
	"testing"
)

func TestNameLocalWireFields(t *testing.T) {
	for _, body := range []string{`{"name_local":"München"}`, `{"name_local":null}`, `{}`} {
		var country Country
		var state State
		var city City
		var nearest CityNearest
		var language Language
		var holiday Holiday
		var pointCity PointCity
		for _, value := range []any{&country, &state, &city, &nearest, &language, &holiday, &pointCity} {
			if err := json.Unmarshal([]byte(body), value); err != nil {
				t.Fatal(err)
			}
		}
		for _, name := range []*string{country.NameLocal, state.NameLocal, city.NameLocal, nearest.NameLocal, language.NameLocal, holiday.NameLocal, pointCity.NameLocal} {
			if body == `{"name_local":"München"}` {
				if name == nil || *name != "München" {
					t.Fatalf("lost native name: %v", name)
				}
			} else if name != nil {
				t.Fatalf("expected unknown native name, got %q", *name)
			}
		}
	}
}
