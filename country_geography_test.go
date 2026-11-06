package parseapi

import (
	"encoding/json"
	"testing"
)

func TestCountryGeographyValues(t *testing.T) {
	var country Country
	err := json.Unmarshal([]byte(`{"country":"XX","gdp_usd":30769700000000,"gdp_per_capita_usd":90026.5163,"population_growth_percent":-0.25,"urban_population_percent":0,"internet_users_percent":95.5,"land_area":10010.5,"water_area":0,"coastline":0,"elevation":0,"lowest_point":{"name":null,"elevation":-430.5},"highest_point":{"name":"Summit","elevation":8848.86}}`), &country)
	if err != nil {
		t.Fatal(err)
	}
	c := &country
	if c == nil || c.GDPUSD == nil || *c.GDPUSD != 30769700000000 || c.GDPPerCapitaUSD == nil || *c.GDPPerCapitaUSD != 90026.5163 || c.PopulationGrowthPercent == nil || *c.PopulationGrowthPercent != -0.25 || c.UrbanPopulationPercent == nil || *c.UrbanPopulationPercent != 0 || c.InternetUsersPercent == nil || *c.InternetUsersPercent != 95.5 {
		t.Fatalf("country statistics: %#v", c)
	}
	if c == nil || c.LandArea == nil || *c.LandArea != 10010.5 || c.WaterArea == nil || *c.WaterArea != 0 || c.Coastline == nil || *c.Coastline != 0 || c.Elevation == nil || *c.Elevation != 0 {
		t.Fatalf("country geography: %#v", c)
	}
	if c.LowestPoint == nil || c.LowestPoint.Name != nil || c.LowestPoint.Elevation != -430.5 || c.HighestPoint == nil || c.HighestPoint.Name == nil || *c.HighestPoint.Name != "Summit" || c.HighestPoint.Elevation != 8848.86 {
		t.Fatalf("country extremes: %#v / %#v", c.LowestPoint, c.HighestPoint)
	}
}

func TestCountryGeographyMissingAndNull(t *testing.T) {
	for _, body := range []string{`{}`, `{}`, `{"land_area":null,"water_area":null,"coastline":null,"elevation":null,"lowest_point":null,"highest_point":null}`} {
		var country Country
		if err := json.Unmarshal([]byte(body), &country); err != nil {
			t.Fatal(err)
		}
		c := &country
		if c != nil && (c.LandArea != nil || c.WaterArea != nil || c.Coastline != nil || c.Elevation != nil || c.LowestPoint != nil || c.HighestPoint != nil) {
			t.Fatalf("unknown country geography became known: %#v", c)
		}
	}
}
