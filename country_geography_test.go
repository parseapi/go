package parseapi

import (
	"encoding/json"
	"testing"
)

func TestCountryGeographyValues(t *testing.T) {
	var country Country
	err := json.Unmarshal([]byte(`{"country":"XX","deep":{"land_area":10010.5,"water_area":0,"coastline":0,"elevation":0,"lowest_point":{"name":null,"elevation":-430.5},"highest_point":{"name":"Summit","elevation":8848.86}}}`), &country)
	if err != nil {
		t.Fatal(err)
	}
	c := country.Deep
	if c == nil || c.LandArea == nil || *c.LandArea != 10010.5 || c.WaterArea == nil || *c.WaterArea != 0 || c.Coastline == nil || *c.Coastline != 0 || c.Elevation == nil || *c.Elevation != 0 {
		t.Fatalf("country geography: %#v", c)
	}
	if c.LowestPoint == nil || c.LowestPoint.Name != nil || c.LowestPoint.Elevation != -430.5 || c.HighestPoint == nil || c.HighestPoint.Name == nil || *c.HighestPoint.Name != "Summit" || c.HighestPoint.Elevation != 8848.86 {
		t.Fatalf("country extremes: %#v / %#v", c.LowestPoint, c.HighestPoint)
	}
}

func TestCountryGeographyMissingAndNull(t *testing.T) {
	for _, body := range []string{`{}`, `{"deep":{}}`, `{"deep":{"land_area":null,"water_area":null,"coastline":null,"elevation":null,"lowest_point":null,"highest_point":null}}`} {
		var country Country
		if err := json.Unmarshal([]byte(body), &country); err != nil {
			t.Fatal(err)
		}
		if (country.Deep == nil) != (body == `{}`) {
			t.Fatal("omitted and locked deep were conflated")
		}
		c := country.Deep
		if c != nil && (c.LandArea != nil || c.WaterArea != nil || c.Coastline != nil || c.Elevation != nil || c.LowestPoint != nil || c.HighestPoint != nil) {
			t.Fatalf("unknown country geography became known: %#v", c)
		}
	}
}
