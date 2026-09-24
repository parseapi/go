package parseapi

import (
	"encoding/json"
	"testing"
)

func TestPostalLocalitiesPreserveObservationAndAmbiguity(t *testing.T) {
	for _, tc := range []struct {
		name, field string
		observed    bool
		count       int
	}{
		{"missing", "", false, 0},
		{"null", `,"localities":null`, false, 0},
		{"empty", `,"localities":[]`, true, 0},
		{"one", `,"localities":[{"city":"SYDNEY","state":"NSW","state_name":"New South Wales","future":true}]`, true, 1},
		{"multiple", `,"localities":[{"city":"SYDNEY","state":"NSW","state_name":"New South Wales"},{"city":"HAYMARKET","state":"NSW","state_name":"New South Wales"}]`, true, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var postal Postal
			if err := json.Unmarshal([]byte(`{"postal":"2000","country":"AU","city":null`+tc.field+`}`), &postal); err != nil {
				t.Fatal(err)
			}
			if postal.City != nil || (postal.Localities != nil) != tc.observed || len(postal.Localities) != tc.count {
				t.Fatalf("lost observation or inferred city: %#v", postal)
			}
			if tc.count > 0 && (postal.Localities[0].City != "SYDNEY" || postal.Localities[0].State != "NSW" || postal.Localities[0].StateName != "New South Wales") {
				t.Fatalf("wrong choice: %#v", postal.Localities[0])
			}
			encoded, err := json.Marshal(postal)
			if err != nil {
				t.Fatal(err)
			}
			var again Postal
			if err := json.Unmarshal(encoded, &again); err != nil {
				t.Fatal(err)
			}
			if (again.Localities != nil) != tc.observed || len(again.Localities) != tc.count {
				t.Fatal("round trip erased null versus empty")
			}
		})
	}
}
