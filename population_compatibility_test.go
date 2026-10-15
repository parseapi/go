package parseapi

import (
 "encoding/json"
 "strings"
 "testing"
)

func TestRetiredPopulationPeriodSourceCompatibilityAndNoWireField(t *testing.T) {
 period := "2024"
 values := []any{CountryDeep{PopulationPeriod:&period},StateDeep{PopulationPeriod:&period},DistrictDeep{PopulationPeriod:&period},PostalDeep{PopulationPeriod:&period},CityDeep{PopulationPeriod:&period},StateDistrictItemDeep{PopulationPeriod:&period}}
 for _, value := range values {
  payload, err := json.Marshal(value)
  if err != nil { t.Fatal(err) }
  if strings.Contains(string(payload), "population_period") { t.Fatalf("legacy member returned in JSON: %s",payload) }
 }
 var city CityDeep
 if err := json.Unmarshal([]byte(`{"population":0,"population_period":"2024"}`),&city); err != nil { t.Fatal(err) }
 if city.PopulationPeriod != nil || city.Population == nil || *city.Population != 0 { t.Fatalf("retired member must remain nil while preserving zero: %#v",city) }
}
