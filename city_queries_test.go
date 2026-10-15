package parseapi

import (
	"context"
	"encoding/json"
	"net/url"
	"reflect"
	"testing"
)

func TestCityQueriesEncodeFiltersAndOpaqueCursor(t *testing.T) {
	ctx := context.Background()
	zero, maximum := int64(0), int64(5000000)
	cursor := "opaque+/=cursor?&"
	base := "min_population=0&max_population=5000000&timezone=Asia%2FTokyo&bbox=170%2C-20%2C-170%2C20&sort=population&order=desc&cursor=opaque%2B%2F%3Dcursor%3F%26"
	cases := []struct {
		name, path, query string
		invoke            func(*Client) error
	}{
		{"browse", "/city", "browse=true&format=list&country=JP&limit=5&hierarchy=true&" + base, func(c *Client) error {
			_, err := c.CityList(ctx, CityListOptions{Country: "JP", Limit: 5, Hierarchy: true, MinPopulation: &zero, MaxPopulation: &maximum, Timezone: "Asia/Tokyo", BBox: "170,-20,-170,20", Sort: "population", Order: "desc", Cursor: cursor})
			return err
		}},
		{"coordinate nearby", "/city", "nearby=true&format=list&lat=0&lon=0&radius=20&unit=mi&country=JP&state=27&" + base, func(c *Client) error {
			_, err := c.CityNearbyCoords(ctx, 0, 0, CityNearbyCoordsOptions{Radius: 20, Unit: "mi", Country: "JP", State: "27", MinPopulation: &zero, MaxPopulation: &maximum, Timezone: "Asia/Tokyo", BBox: "170,-20,-170,20", Sort: "population", Order: "desc", Cursor: cursor})
			return err
		}},
		{"search page", "/city", "q=Sak&page=true&" + base, func(c *Client) error {
			_, err := c.CitySearch(ctx, "Sak", CitySearchOptions{Page: true, MinPopulation: &zero, MaxPopulation: &maximum, Timezone: "Asia/Tokyo", BBox: "170,-20,-170,20", Sort: "population", Order: "desc", Cursor: cursor})
			return err
		}},
		{"named nearby page", "/city/Sakai/nearby", "page=true&" + base, func(c *Client) error {
			_, err := c.CityNearby(ctx, "Sakai", CityNearbyOptions{Page: true, MinPopulation: &zero, MaxPopulation: &maximum, Timezone: "Asia/Tokyo", BBox: "170,-20,-170,20", Sort: "population", Order: "desc", Cursor: cursor})
			return err
		}},
		{"browse defaults", "/city", "browse=true&format=list", func(c *Client) error { _, err := c.CityList(ctx); return err }},
		{"coordinate defaults", "/city", "nearby=true&format=list&lat=0&lon=0", func(c *Client) error { _, err := c.CityNearbyCoords(ctx, 0, 0); return err }},
		{"cursor alone", "/city", "q=Sak&cursor=opaque%2B%2F%3Dcursor%3F%26", func(c *Client) error {
			_, err := c.CitySearch(ctx, "Sak", CitySearchOptions{Cursor: cursor})
			return err
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client, call := newTestClient(t, okJSON(`{}`))
			if err := tc.invoke(client); err != nil {
				t.Fatal(err)
			}
			got, _ := url.ParseQuery(call.rawQry)
			want, _ := url.ParseQuery(tc.query)
			if call.path != tc.path || !reflect.DeepEqual(got, want) {
				t.Fatalf("got %s?%s, want %s?%s", call.path, call.rawQry, tc.path, tc.query)
			}
		})
	}
}

func TestCityQueriesDecodePaginationWithoutInventingCursor(t *testing.T) {
	var list CityList
	if err := json.Unmarshal([]byte(`{"object":"list","has_more":false,"url":"/city","data":[{"object":"city","id":"city_one","name":"One","country":"JP"}],"country":"JP","next_cursor":"opaque","future":true}`), &list); err != nil {
		t.Fatal(err)
	}
	if len(list.Data) != 1 || list.NextCursor == nil || *list.NextCursor != "opaque" || list.Country == nil || *list.Country != "JP" {
		t.Fatalf("lost collection metadata: %#v", list)
	}
	var coords CityNearbyCoords
	if err := json.Unmarshal([]byte(`{"latitude":0,"longitude":0,"radius":20,"unit":"mi","object":"list","has_more":false,"url":"/city","data":[{"object":"city","id":"city_one","name":"One","country":"JP","distance":4,"distance_mi":2.485}],"next_cursor":null}`), &coords); err != nil {
		t.Fatal(err)
	}
	if coords.Latitude != 0 || coords.Longitude != 0 || coords.Radius != 20 || len(coords.Data) != 1 || coords.Data[0].Distance == nil || *coords.Data[0].Distance != 4 || coords.NextCursor != nil {
		t.Fatalf("lost coordinate collection: %#v", coords)
	}
	for _, body := range []string{`{"cities":[]}`, `{"cities":[],"next_cursor":null}`} {
		var legacy CitySearch
		if err := json.Unmarshal([]byte(body), &legacy); err != nil {
			t.Fatal(err)
		}
		if legacy.NextCursor != nil {
			t.Fatal("invented cursor")
		}
	}
	for _, body := range []string{`{"nearby":[]}`, `{"nearby":[],"next_cursor":null}`} {
		var legacy CityNearby
		if err := json.Unmarshal([]byte(body), &legacy); err != nil {
			t.Fatal(err)
		}
		if legacy.NextCursor != nil {
			t.Fatal("invented cursor")
		}
	}
}
