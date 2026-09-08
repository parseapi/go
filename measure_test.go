package parseapi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestMeasureWireAndDecimal(t *testing.T) {
	body := `{"measure":"5 ft 11 in","valid":true,"type":"future-type","amount":"180.34000000000000000001","unit":"cm","reason":null,"choices":[],"future":null}`
	client, call := newTestClient(t, okJSON(body))
	result, err := client.Measure(context.Background(), "5 ft 11 in", MeasureOptions{To: "cm", Locale: "en-US", System: "us"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Amount == nil || *result.Amount != "180.34000000000000000001" || result.Type == nil || *result.Type != "future-type" {
		t.Fatalf("lost precision or vocabulary: %#v", result)
	}
	if call.path != "/measure/5%20ft%2011%20in" || call.rawQry != "locale=en-US&system=us&to=cm" {
		t.Fatalf("unexpected request: %#v", call)
	}
	if _, err := client.Measure(context.Background(), "1 kg/m^3", MeasureOptions{To: "g/L"}); err != nil {
		t.Fatal(err)
	}
	if call.path != "/measure/1%20kg%2Fm%5E3" || call.rawQry != "to=g%2FL" {
		t.Fatalf("bad compound encoding: %#v", call)
	}
	if _, err := client.Measure(context.Background(), "1 m", MeasureOptions{}, MeasureOptions{}); err == nil {
		t.Fatal("accepted multiple options")
	}
}

func TestMeasureCatalogAndUnknown(t *testing.T) {
	client, call := newTestClient(t, okJSON(`{"units":[{"unit":"m","name":"metre","type":"length","aliases":["meter"],"future":true}]}`))
	result, err := client.MeasureUnits(context.Background(), MeasureUnitsOptions{Query: "US gallon", Type: "volume", Unit: "L"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Units) != 1 || result.Units[0].Aliases[0] != "meter" {
		t.Fatalf("bad catalog: %#v", result)
	}
	if call.path != "/measure/units" || call.rawQry != "q=US+gallon&type=volume&unit=L" {
		t.Fatalf("bad catalog request: %#v", call)
	}
	_, _ = client.MeasureUnits(context.Background())
	if call.rawQry != "" {
		t.Fatal("unexpected discovery defaults")
	}
	var invalid Measure
	if err := json.Unmarshal([]byte(`{"measure":"1 gallon","valid":false,"type":null,"amount":null,"unit":null,"reason":"ambiguous_unit","choices":[{"unit":"us_gal","name":"US liquid gallon"}]}`), &invalid); err != nil {
		t.Fatal(err)
	}
	if invalid.Valid || invalid.Amount != nil || invalid.Type != nil || invalid.Reason == nil || len(invalid.Choices) != 1 {
		t.Fatalf("lost ambiguity: %#v", invalid)
	}
}

func TestMeasureBadTargetUsesAPIError(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"code":"bad_request","message":"Incompatible units","request_id":"req_measure"}`))
	})
	_, err := client.Measure(context.Background(), "1 m", MeasureOptions{To: "kg"})
	if err == nil {
		t.Fatal("expected API error")
	}
}
