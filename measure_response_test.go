package parseapi

import (
    "context"
    "encoding/json"
    "os"
    "testing"
)

func TestMeasureCleanAndLegacyResponses(t *testing.T) {
    data, err := os.ReadFile("testdata/measure-response.json")
    if err != nil { t.Fatal(err) }
    var cases []json.RawMessage
    if err := json.Unmarshal(data, &cases); err != nil { t.Fatal(err) }
    for _, body := range cases {
        var expected Measure
        if err := json.Unmarshal(body, &expected); err != nil { t.Fatal(err) }
        client, _ := newTestClient(t, okJSON(string(body)))
        answer, err := client.Measure(context.Background(), expected.Measure)
        if err != nil { t.Fatal(err) }
        if (answer.Amount == nil) != (expected.Amount == nil) || (answer.Amount != nil && *answer.Amount != *expected.Amount) { t.Fatal("lost decimal amount") }
        if (answer.Reason == nil) != (expected.Reason == nil) || (answer.Reason != nil && *answer.Reason != *expected.Reason) { t.Fatal("lost reason") }
        if len(answer.Choices) != len(expected.Choices) { t.Fatal("lost choices") }
        if answer.Valid {
            encoded, err := json.Marshal(answer)
            if err != nil { t.Fatal(err) }
            var fields map[string]json.RawMessage
            if err := json.Unmarshal(encoded, &fields); err != nil { t.Fatal(err) }
            if len(fields) != 5 { t.Fatalf("successful serialization contains noise: %s", encoded) }
        }
    }
}
