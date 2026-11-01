package parseapi

import (
 "encoding/json"
 "testing"
)

func TestPointDeepRetainsMetricAndElectoralFacts(t *testing.T) {
	for _, body := range []string{`{"latitude":0,"longitude":0}`, `{"latitude":0,"longitude":0,"deep":{}}`, `{"latitude":0,"longitude":0,"deep":{"elevation":null,"electoral":null}}`} {
		var point Point
		if err := json.Unmarshal([]byte(body), &point); err != nil {
			t.Fatal(err)
		}
		if point.Deep != nil && (point.Deep.Elevation != nil || point.Deep.Electoral != nil) {
			t.Fatal("unknown detail became a fact")
		}
	}
	var zero Point
	if err := json.Unmarshal([]byte(`{"latitude":0,"longitude":0,"deep":{"elevation":0,"elevation_ft":0,"electoral":{"code":"0801","name":null,"congress":119,"future":true}}}`), &zero); err != nil {
		t.Fatal(err)
	}
	if zero.Deep.Elevation == nil || *zero.Deep.Elevation != 0 || zero.Deep.ElevationFt == nil || *zero.Deep.ElevationFt != 0 {
		t.Fatal("zero terrain was lost")
	}
	if zero.Deep.Electoral == nil || zero.Deep.Electoral.Code != "0801" || zero.Deep.Electoral.Name != nil || zero.Deep.Electoral.Congress == nil || *zero.Deep.Electoral.Congress != 119 {
		t.Fatal("electoral data changed")
	}
	var old Point
	if err := json.Unmarshal([]byte(`{"deep":{"elevation":-1.5,"resolution":30,"electoral":{"code":"0200","name":"At-large","congress":null}}}`), &old); err != nil {
		t.Fatal(err)
	}
	if old.Deep.Elevation == nil || *old.Deep.Elevation != -1.5 || old.Deep.Resolution == nil || *old.Deep.Resolution != 30 || old.Deep.Electoral.Congress != nil {
		t.Fatal("nullable or historical detail changed")
	}
}
