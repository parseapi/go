package parseapi

import (
 "encoding/json"
 "os"
 "testing"
)

func TestApprovedCityMatchWithoutMechanics(t *testing.T) {
 raw, err := os.ReadFile("testdata/noise-approved.json"); if err != nil { t.Fatal(err) }
 var rows []struct { Body struct { Data []json.RawMessage `json:"data"` } `json:"body"` }; if err = json.Unmarshal(raw, &rows); err != nil { t.Fatal(err) }
 for _, row := range rows { for _, rawCity := range row.Body.Data {
  var city CityResource; if err = json.Unmarshal(rawCity, &city); err != nil { t.Fatal(err) }
  if city.Match == nil || city.Match.Value == "" || city.Match.Type == "" { t.Fatal("lost useful match evidence") }
  if city.Match.Source != "" || city.Match.EditDistance != nil { t.Fatal("invented match mechanics") }
  actual, _ := json.Marshal(city.Match); var fields map[string]any; json.Unmarshal(actual, &fields)
  if _, ok := fields["source"]; ok { t.Fatal("invented source") }; if _, ok := fields["edit_distance"]; ok { t.Fatal("invented score") }
 } }
}
