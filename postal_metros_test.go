package parseapi

import (
	"encoding/json"
	"testing"
)

func TestPostalMetrosObservationStates(t *testing.T) {
	for _, tc := range []struct {
		name, field string
		observed    bool
		count       int
	}{
		{"missing", "", false, 0},
		{"null", `,"metros":null`, false, 0},
		{"outside", `,"metros":[]`, true, 0},
		{"populated", `,"metros":[{"code":"12345","name":"Example area","type":"future-area-type","share":0.75,"residential_share":0,"business_share":1,"other_share":null,"future":true}]`, true, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			member := `{"postal":"12345","country":"US","city":null,"future":true` + tc.field + `}`
			var postal Postal
			var nearby PostalNearby
			var distance PostalDistance
			for _, input := range []struct {
				body  string
				value any
			}{
				{member, &postal},
				{`{"postal":"12345","country":"US","nearby":[` + member + `]` + tc.field + `}`, &nearby},
				{`{"country":"US","from":` + member + `,"to":` + member + `}`, &distance},
			} {
				if err := json.Unmarshal([]byte(input.body), input.value); err != nil {
					t.Fatal(err)
				}
			}
			for _, metros := range [][]PostalMetro{postal.Metros, nearby.Metros, nearby.Nearby[0].Metros, distance.From.Metros, distance.To.Metros} {
				if (metros != nil) != tc.observed || len(metros) != tc.count {
					t.Fatalf("wrong observation state: %#v", metros)
				}
				if len(metros) > 0 {
					m := metros[0]
					if m.Code != "12345" || m.Name != "Example area" || m.Type != "future-area-type" || m.Share == nil || *m.Share != 0.75 || m.ResidentialShare == nil || *m.ResidentialShare != 0 || m.BusinessShare == nil || *m.BusinessShare != 1 || m.OtherShare != nil {
						t.Fatalf("wrong association: %#v", m)
					}
				}
			}
			encoded, err := json.Marshal(postal)
			if err != nil {
				t.Fatal(err)
			}
			var roundtrip Postal
			if err := json.Unmarshal(encoded, &roundtrip); err != nil {
				t.Fatal(err)
			}
			if (roundtrip.Metros != nil) != tc.observed {
				t.Fatal("serialization erased null versus empty")
			}
		})
	}
}

func TestPostalMetroMissingSharesStayUnknown(t *testing.T) {
	var metro PostalMetro
	if err := json.Unmarshal([]byte(`{"code":"12345","name":"Example area","type":"metropolitan"}`), &metro); err != nil {
		t.Fatal(err)
	}
	if metro.Share != nil || metro.ResidentialShare != nil || metro.BusinessShare != nil || metro.OtherShare != nil {
		t.Fatal("missing share became a known zero")
	}
}
