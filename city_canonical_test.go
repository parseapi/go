package parseapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestCityCanonicalSparseTransportAndFilters(t *testing.T) {
	ctx := context.Background()
	id := "city_7aj899nvwyku"
	c, call := newTestClient(t, okJSON(`{"object":"list","data":[{"id":"city_7aj899nvwyku","object":"city","name":"Vienna"}],"has_more":true,"url":"/city","next_cursor":"opaque+/="}`))
	low, zero := float64(-10), float64(0)
	pop := int64(0)
	result, err := c.CityList(ctx, CityListOptions{Countries: []string{"AT", "CH"}, ExcludeCountries: []string{"JP"}, Timezones: []string{"Europe/Vienna"}, MinElevation: &low, MaxElevation: &zero, MinPopulation: &pop, Sort: "-population,+name", Cursor: "opaque+/="})
	if err != nil {
		t.Fatal(err)
	}
	if result.Object != "list" || result.NextCursor == nil || result.Data[0].ID != id || result.Data[0].Country != nil || result.Data[0].Name == nil {
		t.Fatalf("sparse result %#v", result)
	}
	q, _ := url.ParseQuery(call.rawQry)
	for k, v := range map[string]string{"format": "list", "countries": "AT,CH", "exclude_countries": "JP", "timezones": "Europe/Vienna", "min_elevation": "-10", "max_elevation": "0", "min_population": "0", "sort": "-population,+name", "cursor": "opaque+/="} {
		if q.Get(k) != v {
			t.Fatalf("%s=%q want%q", k, q.Get(k), v)
		}
	}
	c, call = newTestClient(t, okJSON(`{"id":"city_7aj899nvwyku","object":"city","name_local":"Wien"}`))
	resource, err := c.CityResource(ctx, "Vienna", CityResourceOptions{Country: "AT"})
	if err != nil {
		t.Fatal(err)
	}
	if resource.NameLocal == nil || *resource.NameLocal != "Wien" || resource.Country != nil {
		t.Fatal("lost sparse native fields")
	}
	q, _ = url.ParseQuery(call.rawQry)
	if q.Get("format") != "resource" || q.Has("fields") {
		t.Fatal(q)
	}
	_, err = c.City(ctx, "Vienna")
	if err != nil {
		t.Fatal(err)
	}
	if call.rawQry != "" {
		t.Fatal("legacy default changed", call.rawQry)
	}
}

func TestCityCanonicalRelationsDistanceAndNewWrappers(t *testing.T) {
	ctx := context.Background()
	id := "city_7aj899nvwyku"
	c, call := newTestClient(t, okJSON(`{"object":"list","data":[],"has_more":false,"url":"/city","next_cursor":null,"city_id":"city_7aj899nvwyku","relationship":"children","coverage":"partial"}`))
	r, err := c.CityChildren(ctx, id, CityRelationsOptions{Cursor: "opaque+/=", Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if r.Coverage != "partial" || call.path != "/city/id/"+id+"/children" {
		t.Fatal(r, call.path)
	}
	if _, err = c.CityAncestors(ctx, id); err != nil {
		t.Fatal(err)
	}
	if call.path != "/city/id/"+id+"/ancestors" {
		t.Fatal(call.path)
	}
	if _, err = c.CitySiblings(ctx, id); err != nil {
		t.Fatal(err)
	}
	if _, err = c.CitySearchList(ctx, "Vi", CitySearchListOptions{}); err != nil {
		t.Fatal(err)
	}
	q, _ := url.ParseQuery(call.rawQry)
	if q.Get("format") != "list" {
		t.Fatal(q)
	}
	if _, err = c.CityNearbyList(ctx, "Vienna"); err != nil {
		t.Fatal(err)
	}
	if call.path != "/city/Vienna/nearby" {
		t.Fatal(call.path)
	}
	for _, value := range []string{"0", "null"} {
		c, call := newTestClient(t, okJSON(fmt.Sprintf(`{"object":"city_distance","from":"%s","to":"%s","distance":%s,"distance_mi":%s,"method":"great_circle"}`, id, id, value, value)))
		d, err := c.CityDistance(ctx, id, id)
		if err != nil {
			t.Fatal(err)
		}
		if (value == "null") != (d.Distance == nil) {
			t.Fatal(d)
		}
		if d.Distance != nil && *d.Distance != 0 {
			t.Fatal(d)
		}
		if call.rawQry != "" {
			t.Fatal(call.rawQry)
		}
	}
}

func TestCityCanonicalBounded409CandidatesAndOrdinaryErrors(t *testing.T) {
	safe := map[string]any{"id": "city_7aj899nvwyku", "object": "city", "name": "Vienna", "deep": map[string]any{"population": 99}}
	for _, tc := range []struct {
		status  int
		code    string
		present bool
	}{{409, "ambiguous_city", true}, {400, "ambiguous_city", false}, {409, "invalid_request", false}, {429, "rate_limited", false}} {
		calls := 0
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			calls++
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(tc.status)
			json.NewEncoder(w).Encode(map[string]any{"code": tc.code, "message": "Choose", "candidates": []any{safe, safe, safe}})
		})
		_, err := c.City(context.Background(), "Vienna", CityOptions{Ambiguity: "error"})
		var e *Error
		if !errors.As(err, &e) || e.Status != tc.status || e.Code != tc.code || e.RetryAfter == nil || *e.RetryAfter != "60" || calls != 1 {
			t.Fatal(err, calls)
		}
		if tc.present != (e.Candidates != nil) {
			t.Fatal(e)
		}
		if tc.present && (len(e.Candidates) != 2 || e.Candidates[0].Name == nil || *e.Candidates[0].Name != "Vienna") {
			t.Fatal(e)
		}
	}
}

func TestCityExplicitMatchingRetainsSuggestionEvidence(t *testing.T) {
	c, call := newTestClient(t, okJSON(`{"object":"list","data":[{"id":"city_7aj899nvwyku","object":"city","match":{"type":"fuzzy","source":"name","value":"Vienna","language":null,"edit_distance":1}}],"has_more":false,"url":"/city","next_cursor":null}`))
	page, err := c.CitySearchList(context.Background(), "Viena", CitySearchListOptions{Match: "fuzzy", MaxDistance: 1})
	if err != nil {
		t.Fatal(err)
	}
	if page.Data[0].Match == nil || page.Data[0].Match.EditDistance == nil || *page.Data[0].Match.EditDistance != 1 {
		t.Fatal("lost suggestion evidence")
	}
	if call.path != "/city/search" {
		t.Fatal("explicit matching endpoint", call.path)
	}
	q, _ := url.ParseQuery(call.rawQry)
	if q.Get("match") != "fuzzy" || q.Get("max_distance") != "1" || q.Get("format") != "list" {
		t.Fatal(q)
	}
	_, err = c.CitySearch(context.Background(), "Vi")
	if err != nil {
		t.Fatal(err)
	}
	if call.rawQry != "q=Vi" {
		t.Fatal("classic default changed", call.rawQry)
	}
}

func TestCityReviewedRomanizationOptIn(t *testing.T) {
	c, call := newTestClient(t, okJSON(`{"id":"city_7aj899nvwyku","object":"city","name_romanized":"Kitakyūshū","romanization_language":"ja-Latn"}`))
	row, e := c.CityIDResource(context.Background(), "city_7aj899nvwyku", CityIDResourceOptions{Romanized: true})
	if e != nil {
		t.Fatal(e)
	}
	if row.NameRomanized == nil || *row.NameRomanized != "Kitakyūshū" || row.RomanizationLanguage == nil || *row.RomanizationLanguage != "ja-Latn" {
		t.Fatal(row)
	}
	q, _ := url.ParseQuery(call.rawQry)
	if q.Get("romanized") != "true" {
		t.Fatal(q)
	}
	c, _ = newTestClient(t, okJSON(`{"id":"city_7aj899nvwyku","object":"city","name_romanized":null,"romanization_language":null}`))
	row, e = c.CityIDResource(context.Background(), "city_7aj899nvwyku", CityIDResourceOptions{Romanized: true})
	if e != nil || row.NameRomanized != nil || row.RomanizationLanguage != nil {
		t.Fatal(row, e)
	}
}

func TestCityClassicExplicitMatchAdaptsCanonicalEnvelope(t *testing.T) {
	c, call := newTestClient(t, okJSON(`{"object":"list","q":"Viena","data":[{"id":"city_7aj899nvwyku","name":"Vienna","country":"AT"}],"has_more":false,"next_cursor":null}`))
	page, err := c.CitySearch(context.Background(), "Viena", CitySearchOptions{Country: "AT", Match: "fuzzy", MaxDistance: 1})
	if err != nil || len(page.Cities) != 1 || page.Cities[0].ID != "city_7aj899nvwyku" || page.Q != "Viena" || call.path != "/city/search" {
		t.Fatalf("page%+v error%v endpoint%s", page, err, call.path)
	}
}
