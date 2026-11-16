package parseapi

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestCountyRoutesPin22AndRetainMethodAliases(t *testing.T) {
	const body = `{"county":"08031","name":"Denver County","country":"US","type":"county","timezones":["America/Denver"],"population":0}`
	client, request := newTestClient(t, okJSON(body))
	county, err := client.County(context.Background(), "Denver County", CountyOptions{Country: "US", State: "CO", Deep: true})
	if err != nil || county.County != "08031" || county.District != "08031" || county.Population == nil || *county.Population != 0 {
		t.Fatalf("county: %#v %v", county, err)
	}
	if request.path != "/county/Denver%20County" || request.rawQry != "country=US&deep=true&state=CO" || request.header.Get("X-API-Version") != "2.2.0" {
		t.Fatalf("request: %#v", request)
	}
	old, err := client.District(context.Background(), "08031")
	if err != nil || old.County != "08031" || request.path != "/county/08031" {
		t.Fatalf("legacy method: %#v %v", old, err)
	}
	listClient, listRequest := newTestClient(t, okJSON(`{"state":"CO","country":"US","counties":[{"county":"08031","name":"Denver County","type":"county","population":0}]}`))
	listing, err := listClient.StateCounties(context.Background(), "CO", StateCountiesOptions{Country: "US"})
	if err != nil || len(listing.Counties) != 1 || len(listing.Districts) != 1 || listRequest.path != "/state/CO/counties" {
		t.Fatalf("listing: %#v %v", listing, err)
	}
}
func TestCountyWireReadersKeepOldKeysNullsAndOneEncodedFact(t *testing.T) {
	var postal Postal
	if err := json.Unmarshal([]byte(`{"postal":"08001","country":"US","district":"08031","district_name":"Denver","district_name_local":null}`), &postal); err != nil {
		t.Fatal(err)
	}
	if postal.County == nil || *postal.County != "08031" || postal.District != postal.County || postal.CountyNameLocal != nil {
		t.Fatalf("old wire: %#v", postal)
	}
	encoded, err := json.Marshal(postal)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), `"district`) || !strings.Contains(string(encoded), `"county":"08031"`) {
		t.Fatalf("duplicate/old wire: %s", encoded)
	}
	var ip IP
	if err := json.Unmarshal([]byte(`{"ip":"8.8.8.8","county":null,"district":"bad","county_name":null,"district_name":"bad"}`), &ip); err != nil {
		t.Fatal(err)
	}
	if ip.County != nil || ip.District != nil || ip.CountyName != nil {
		t.Fatal("canonical null lost")
	}
	var nearest CityNearest
	if err := json.Unmarshal([]byte(`{"name":"Paris","country":"FR","district":"75","district_name":"Paris","type":"district","capital_of":"county","distance":0.25,"distance_mi":0.16}`), &nearest); err != nil {
		t.Fatal(err)
	}
	if nearest.County == nil || *nearest.County != "75" || nearest.Distance != 0.25 || nearest.Type == nil || *nearest.Type != "district" || nearest.CapitalOf == nil || *nearest.CapitalOf != "county" {
		t.Fatalf("nearest fields: %#v", nearest)
	}
}
