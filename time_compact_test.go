package parseapi

import (
 "encoding/json"
 "os"
 "testing"
)

func TestApprovedCompactTimePlacesAndChoices(t *testing.T) {
 raw, err := os.ReadFile("testdata/time-compact.json"); if err != nil { t.Fatal(err) }
 var rows []struct { ID string `json:"id"`; Body json.RawMessage `json:"body"` }
 if err = json.Unmarshal(raw, &rows); err != nil { t.Fatal(err) }
 for _, row := range rows { t.Run(row.ID, func(t *testing.T) {
  var result Time; if err := json.Unmarshal(row.Body, &result); err != nil { t.Fatal(err) }
  var source map[string]any; json.Unmarshal(row.Body, &source)
  place, present := source["location"].(map[string]any)
  if !present { if result.Location != nil { t.Fatal("invented place") }; return }
  if result.Location == nil { t.Fatal("lost location") }
  encoded, _ := json.Marshal(result.Location); var actual map[string]any; json.Unmarshal(encoded, &actual)
  for _, field := range []string{"id", "name", "country", "state", "latitude", "longitude"} {
   if expected, ok := place[field]; ok && actual[field] != expected { t.Fatalf("%s: got %#v want %#v", field, actual[field], expected) }
  }
  if status, ok := place["status"]; ok && result.Location.Status != status { t.Fatal("lost unresolved status") }
 }) }
}
