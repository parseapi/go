package parseapi

import (
	"context"
	"encoding/json"
	"net/url"
	"reflect"
	"testing"
)

func TestCityQualityOptionsPreserveDefaults(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name   string
		invoke func(*Client) error
		path   string
		query  string
	}{
		{"old name", func(c *Client) error { _, e := c.City(ctx, "Wien"); return e }, "/city/Wien", ""},
		{"old id", func(c *Client) error { _, e := c.CityID(ctx, "city_pin"); return e }, "/city/id/city_pin", ""},
		{"old search", func(c *Client) error { _, e := c.CitySearch(ctx, "Wie"); return e }, "/city", "q=Wie"},
		{"old nearest", func(c *Client) error { _, e := c.CityNearest(ctx, 48, 16); return e }, "/city", "lat=48&lon=16"},
		{"old nearby", func(c *Client) error { _, e := c.CityNearby(ctx, "Wien"); return e }, "/city/Wien/nearby", ""},
		{"name", func(c *Client) error {
			_, e := c.City(ctx, "Wien", CityOptions{Country: "AT", State: "AT-9", Type: "city", Grain: "municipality", Aliases: true, Hierarchy: true, Ambiguity: "error"})
			return e
		}, "/city/Wien", "country=AT&state=AT-9&type=city&grain=municipality&aliases=true&hierarchy=true&ambiguity=error"},
		{"id", func(c *Client) error { _, e := c.CityID(ctx, "city_pin", CityIDOptions{Hierarchy: true}); return e }, "/city/id/city_pin", "hierarchy=true"},
		{"search", func(c *Client) error {
			_, e := c.CitySearch(ctx, "Wie", CitySearchOptions{Country: "AT", State: "9", Type: "city", Grain: "municipality", Aliases: true, Hierarchy: true})
			return e
		}, "/city", "q=Wie&country=AT&state=9&type=city&grain=municipality&aliases=true&hierarchy=true"},
		{"nearest", func(c *Client) error {
			_, e := c.CityNearest(ctx, 48, 16, CityNearestOptions{Country: "AT", State: "9", Type: "city", Grain: "municipality", Hierarchy: true})
			return e
		}, "/city", "lat=48&lon=16&country=AT&state=9&type=city&grain=municipality&hierarchy=true"},
		{"nearby", func(c *Client) error {
			_, e := c.CityNearby(ctx, "Wien", CityNearbyOptions{Country: "AT", State: "9", Type: "city", Grain: "municipality", Aliases: true, Hierarchy: true, Ambiguity: "error"})
			return e
		}, "/city/Wien/nearby", "country=AT&state=9&type=city&grain=municipality&aliases=true&hierarchy=true&ambiguity=error"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client, call := newTestClient(t, okJSON(`{}`))
			if err := tc.invoke(client); err != nil {
				t.Fatal(err)
			}
			want, _ := url.ParseQuery(tc.query)
			got, _ := url.ParseQuery(call.rawQry)
			if call.path != tc.path || !reflect.DeepEqual(got, want) {
				t.Fatalf("got %s?%s, want %s?%s", call.path, call.rawQry, tc.path, tc.query)
			}
		})
	}
}

func TestCityQualityHierarchyDecoding(t *testing.T) {
	child := `{"id":"city_child","name":"Ward","country":"JP","grain":"subdivision","parent":{"id":"city_parent","name":null,"name_local":"堺市","country":"JP","state":"JP-27","type":"city","grain":"municipality","future":true},"future":true}`
	var city City
	if err := json.Unmarshal([]byte(child), &city); err != nil {
		t.Fatal(err)
	}
	if city.Grain == nil || *city.Grain != "subdivision" || city.Parent == nil || city.Parent.ID != "city_parent" || city.Parent.Name != nil || city.Parent.NameLocal == nil || *city.Parent.NameLocal != "堺市" {
		t.Fatalf("lost nullable hierarchy: %#v", city)
	}
	var search CitySearch
	if err := json.Unmarshal([]byte(`{"cities":[`+child+`]}`), &search); err != nil {
		t.Fatal(err)
	}
	if search.Cities[0].Parent.ID != "city_parent" {
		t.Fatal("search hierarchy lost")
	}
	var nearby CityNearby
	if err := json.Unmarshal([]byte(`{"nearby":[`+child+`]}`), &nearby); err != nil {
		t.Fatal(err)
	}
	if nearby.Nearby[0].Parent.ID != "city_parent" {
		t.Fatal("nearby hierarchy lost")
	}
	for _, body := range []string{`{"id":"city_old","name":"Old","country":"AT"}`, `{"id":"city_old","name":"Old","country":"AT","grain":null,"parent":null}`} {
		var old City
		if err := json.Unmarshal([]byte(body), &old); err != nil {
			t.Fatal(err)
		}
		if old.Grain != nil || old.Parent != nil {
			t.Fatal("unknown hierarchy inferred")
		}
	}
}

func TestCityQualityDatedPopulationPreservesUnknown(t *testing.T) {
	var city City
	if err := json.Unmarshal([]byte(`{"id":"city_pin","deep":{"population":0,"population_period":"2025","population_observed_on":"2025-12-31","population_basis":"registered_population"}}`), &city); err != nil {
		t.Fatal(err)
	}
	if city.Deep.Population == nil || *city.Deep.Population != 0 || city.Deep.PopulationObservedOn == nil || *city.Deep.PopulationObservedOn != "2025-12-31" || city.Deep.PopulationBasis == nil || *city.Deep.PopulationBasis != "registered_population" {
		t.Fatal("dated population metadata lost")
	}
	for _, body := range []string{`{"deep":{}}`, `{"deep":{"population_observed_on":null,"population_basis":null}}`} {
		var old City
		if err := json.Unmarshal([]byte(body), &old); err != nil {
			t.Fatal(err)
		}
		if old.Deep.PopulationObservedOn != nil || old.Deep.PopulationBasis != nil {
			t.Fatal("unknown observation metadata inferred")
		}
	}
}
