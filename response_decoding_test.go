package parseapi

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestOrdinaryResponseAndRetry(t *testing.T) {
 attempts:=0
 client,_:=newTestClient(t,func(w http.ResponseWriter,r *http.Request){
  attempts++
  if r.Header.Get("Parse-Response")!="" || r.Header.Get("Parse-Version")!="2.0.0" {t.Fatal("unexpected request headers")}
  if attempts==1 {w.Header().Set("Retry-After","0");w.WriteHeader(503);return}
  body,err:=os.ReadFile("testdata/response-fixtures/clean-timezones.json");if err!=nil {t.Fatal(err)};_,_=w.Write(body)
 },WithRetries(1))
 result,err:=client.TimeZones(context.Background(),"")
 if err!=nil || len(result.Timezones)==0 || attempts!=2 {t.Fatalf("decode/retry: %v %#v",err,result)}
}
func TestResponseDecodingActualHandlersDecode(t *testing.T) {
	for _, profile := range []string{"standard", "clean"} {
		var zones TimeZones
		var time Time
		var tariff Tariff
		var search TariffSearch
		var metadata EmojiMetadata
		var catalog EmojiCatalog
		var coverage CompanyCoverage
		var company CompanyProfile
		rows := map[string]interface{}{"timezones": &zones, "time": &time, "tariff": &tariff, "tariff-search": &search, "emoji-metadata": &metadata, "emoji-catalog": &catalog, "company-coverage": &coverage, "company-profile": &company}
		for name, target := range rows {
			body, err := os.ReadFile("testdata/response-fixtures/" + profile + "-" + name + ".json")
			if err != nil {
				t.Fatal(err)
			}
			if err = json.Unmarshal(body, target); err != nil {
				t.Fatalf("%s/%s: %v", profile, name, err)
			}
		}
		getters := []string{zones.TimezoneDatabaseVersion, time.Location.Source, tariff.Revision, search.Revision, metadata.SourceVersion, coverage.SnapshotAt, company.Deep.Employees.Method}
		_ = getters
		if company.Deep.Employees.Count != 12 || company.Deep.Employees.Approximate {
			t.Fatal("known count/false lost")
		}
		wantMethod := "reported"
		if profile == "clean" {
			wantMethod = ""
		}
		if company.Deep.Employees.Method != wantMethod {
			t.Fatal("published string decoder semantics changed")
		}
	}
}

// Compare receiver merge/null/error behavior directly with the standard decoder
// for the previously published struct, rather than assuming a fresh destination.
func TestResponseDecodingEmployeeDecoderPreservesExistingReceiver(t *testing.T) {
	type publishedEmployees CompanyProfileEmployees
	seed := CompanyProfileEmployees{Count: 12, AsOf: "2025-12-31", Scope: "consolidated_group", Method: "future_method", Approximate: true}
	for _, data := range []string{`{"count":0}`, `null`, `{"method":null}`, `{"count":"invalid","scope":"updated"}`} {
		got := seed
		want := publishedEmployees(seed)
		err := json.Unmarshal([]byte(data), &got)
		oldErr := json.Unmarshal([]byte(data), &want)
		if (err == nil) != (oldErr == nil) || !reflect.DeepEqual(got, CompanyProfileEmployees(want)) {
			t.Fatalf("receiver changed for %s: got=%#v want=%#v err=%v oldErr=%v", data, got, want, err, oldErr)
		}
	}
	var fresh CompanyProfileEmployees
	if err := json.Unmarshal([]byte(`{"count":12,"as_of":"2025-12-31","scope":"consolidated_group","approximate":false}`), &fresh); err != nil || fresh.Method != "" {
		t.Fatalf("clean fresh getter: %#v, %v", fresh, err)
	}
	var empty CompanyProfileEmployees
	if err := json.Unmarshal([]byte(`null`), &empty); err != nil || empty.Method != "" {
		t.Fatalf("whole null invented metadata: %#v, %v", empty, err)
	}
}

func TestResponseDecodingEmployeesKeepPublishedEmbeddingMethodSet(t *testing.T) {
	type embedded struct {
		CompanyProfileEmployees
		Extra string `json:"extra"`
	}
	var result embedded
	if err := json.Unmarshal([]byte(`{"count":12,"method":"reported","extra":"kept"}`), &result); err != nil {
		t.Fatal(err)
	}
	if result.Extra != "kept" || result.Count != 12 || result.Method != "reported" {
		t.Fatalf("embedded response changed: %#v", result)
	}
	if _, implements := interface{}(&result.CompanyProfileEmployees).(json.Unmarshaler); implements {
		t.Fatal("published employee method set changed")
	}
}
