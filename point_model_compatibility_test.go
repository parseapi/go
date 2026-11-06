package parseapi

import (
	"encoding/json"
	"testing"
)

func TestPointDeepRetainsMetricAndElectoralFacts(t *testing.T) {
	for _, body := range []string{`{"latitude":0,"longitude":0}`, `{"latitude":0,"longitude":0}`, `{"latitude":0,"longitude":0,"elevation":null,"electoral":null}`} {
		var point Point
		if err := json.Unmarshal([]byte(body), &point); err != nil {
			t.Fatal(err)
		}
		if point.Elevation != nil || point.Electoral != nil {
			t.Fatal("unknown detail became a fact")
		}
	}
	var zero Point
	if err := json.Unmarshal([]byte(`{"latitude":0,"longitude":0,"elevation":0,"elevation_ft":0,"electoral":{"code":"0801","name":null,"congress":119,"future":true}}`), &zero); err != nil {
		t.Fatal(err)
	}
	if zero.Elevation == nil || *zero.Elevation != 0 || zero.ElevationFt == nil || *zero.ElevationFt != 0 {
		t.Fatal("zero terrain was lost")
	}
	if zero.Electoral == nil || zero.Electoral.Code != "0801" || zero.Electoral.Name != nil || zero.Electoral.Congress == nil || *zero.Electoral.Congress != 119 {
		t.Fatal("electoral data changed")
	}
	var old Point
	if err := json.Unmarshal([]byte(`{"elevation":-1.5,"resolution":30,"electoral":{"code":"0200","name":"At-large","congress":null}}`), &old); err != nil {
		t.Fatal(err)
	}
	if old.Elevation == nil || *old.Elevation != -1.5 || old.Resolution == nil || *old.Resolution != 30 || old.Electoral.Congress != nil {
		t.Fatal("nullable or historical detail changed")
	}
}
