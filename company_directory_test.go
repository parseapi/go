package parseapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

const companyDirectoryProfile = `{"id":"co_222222222222","name":"Example","country":"US","website":null,"listings":[{"exchange":"Future Exchange","symbol":"A/B"}],"address":null,"deep":{"legal_name":"Example Inc.","aliases":[],"jurisdiction":{"country":"US","state":null},"status":"future-status","websites":[{"domain":"example.com","url":null}],"identifiers":[{"type":"registration","authority":"future:registry","value":"0000123"}],"incorporated":null,"addresses":[{"type":"future-role","street":null,"city":"Example City","state":null,"postal":null,"country":"US"}],"industries":[{"type":"future-scheme","code":"001","name":null}],"parent":null,"description":null,"logo":null,"socials":[],"founded":{"value":"2006","precision":"year"},"sources":[{"type":"website","url":"https://example.com/","fields":["founded"],"observed_at":"2026-09-23T17:35:06.956Z","updated_at":null,"future":true}],"future":true},"future":true}`

const companyDirectorySearch = `{"companies":[{"id":"co_222222222222","name":"Example","country":"US","website":null,"listings":[{"exchange":"Future Exchange","symbol":"A/B"}],"address":null,"deep":{"legal_name":"Example Inc.","aliases":[],"jurisdiction":{"country":"US","state":null},"status":"future-status","websites":[{"domain":"example.com","url":null}],"identifiers":[{"type":"registration","authority":"future:registry","value":"0000123"}],"incorporated":null,"addresses":[{"type":"future-role","street":null,"city":"Example City","state":null,"postal":null,"country":"US"}],"industries":[{"type":"future-scheme","code":"001","name":null}],"parent":null,"description":null,"logo":null,"socials":[],"founded":{"value":"2006","precision":"year"},"sources":[{"type":"website","url":"https://example.com/","fields":["founded"],"observed_at":"2026-09-23T17:35:06.956Z","updated_at":null,"future":true}],"future":true},"future":true,"match":{"field":"future-field","value":"0000123","type":"future-scheme","authority":null,"exchange":null,"future":true}}],"next":"opaque+/="}`

const companyDirectoryCoverage = `{"scope":"sample","label":"Company directory","description":"Edition profiles","snapshot_at":"2026-09-23T17:35:06.956Z","companies":0,"countries":[],"with_website":0,"with_listings":0,"with_address":0,"future":true}`

func TestCompanyDirectoryMethodsAndQueryEncoding(t *testing.T) {
	ctx := context.Background()
	client, call := newTestClient(t, okJSON(companyDirectoryProfile))
	profile, err := client.CompanyID(ctx, "co_/ ?", CompanyIDOptions{Deep: true})
	if err != nil {
		t.Fatal(err)
	}
	if call.path != "/company/id/co_%2F%20%3F" || call.rawQry != "deep=true" {
		t.Fatalf("bad ID route: %#v", call)
	}
	if profile.Deep.Identifiers[0].Value != "0000123" || *profile.Deep.Status != "future-status" || profile.Deep.Socials == nil || len(profile.Deep.Socials) != 0 {
		t.Fatalf("lost profile values: %#v", profile)
	}
	if profile.Deep.Founded.Precision != "year" || profile.Deep.Sources[0].UpdatedAt != nil || profile.Deep.Jurisdiction.State != nil {
		t.Fatal("lost precision/nulls")
	}
	if call.header.Get("Parse-Version") != "2.0.0" || call.header.Get("X-API-Key") != "test_key_123" {
		t.Fatal("request contract changed")
	}
	_, err = client.CompanyID(ctx, "co_222222222222")
	if err != nil || call.rawQry != "" {
		t.Fatal("unexpected default deep", err)
	}
	client, call = newTestClient(t, okJSON(companyDirectorySearch))
	cases := []struct {
		opts CompanySearchOptions
		want url.Values
	}{
		{CompanySearchOptions{Query: "A & B", Country: "US", Limit: 2, Cursor: "opaque+/=", Deep: true}, url.Values{"q": {"A & B"}, "country": {"US"}, "limit": {"2"}, "cursor": {"opaque+/="}, "deep": {"true"}}},
		{CompanySearchOptions{Domain: "https://sub.example.com/a?b=1"}, url.Values{"domain": {"https://sub.example.com/a?b=1"}}},
		{CompanySearchOptions{Ticker: "A/B", Exchange: "Future Exchange"}, url.Values{"ticker": {"A/B"}, "exchange": {"Future Exchange"}}},
		{CompanySearchOptions{Identifier: "0000123", Authority: "future:registry"}, url.Values{"identifier": {"0000123"}, "authority": {"future:registry"}}},
		{CompanySearchOptions{Country: "US"}, url.Values{"country": {"US"}}},
		{CompanySearchOptions{RegistrationAuthority: "ra000599"}, url.Values{"registration_authority": {"ra000599"}}},
		{CompanySearchOptions{Country: "US", Industry: "0700", IndustryType: "sic", RegistrationAuthority: "RA000599", RegistrationForm: "DPC", RegistrationStatus: " Good Standing ", Limit: 2, Cursor: "opaque+/=", Deep: true}, url.Values{"country": {"US"}, "industry": {"0700"}, "industry_type": {"sic"}, "registration_authority": {"RA000599"}, "registration_form": {"DPC"}, "registration_status": {" Good Standing "}, "limit": {"2"}, "cursor": {"opaque+/="}, "deep": {"true"}}},
		{CompanySearchOptions{Identifier: "00001", Authority: "SEC", RegistrationAuthority: "RA000599", RegistrationForm: "future/Form", RegistrationStatus: "future+& status"}, url.Values{"identifier": {"00001"}, "authority": {"SEC"}, "registration_authority": {"RA000599"}, "registration_form": {"future/Form"}, "registration_status": {"future+& status"}}},
		{CompanySearchOptions{Industry: "0700", IndustryType: "sic"}, url.Values{"industry": {"0700"}, "industry_type": {"sic"}}},
		{CompanySearchOptions{Country: "US", Industry: "0700", IndustryType: "sic", Cursor: "opaque+/=", Limit: 2, Deep: true}, url.Values{"country": {"US"}, "industry": {"0700"}, "industry_type": {"sic"}, "cursor": {"opaque+/="}, "limit": {"2"}, "deep": {"true"}}},
		{CompanySearchOptions{Query: "Example", Industry: "0700", IndustryType: "sic"}, url.Values{"q": {"Example"}, "industry": {"0700"}, "industry_type": {"sic"}}},
	}
	for _, tc := range cases {
		result, err := client.CompanySearch(ctx, tc.opts)
		if err != nil {
			t.Fatal(err)
		}
		got, _ := url.ParseQuery(call.rawQry)
		if call.path != "/company" || !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("bad search: %#v", call)
		}
		if result.Next == nil || *result.Next != "opaque+/=" || *result.Companies[0].Match.Field != "future-field" || result.Companies[0].Deep == nil {
			t.Fatal("lost cursor/match/deep")
		}
	}
	client, call = newTestClient(t, okJSON(companyDirectoryCoverage))
	counts, err := client.CompanyCoverage(ctx)
	if err != nil || counts.Companies != 0 || counts.Countries == nil || call.path != "/company/directory/coverage" || call.rawQry != "" {
		t.Fatal("bad coverage", err)
	}
}

func TestCompanyDirectoryEmptyPartialAndFutureResponses(t *testing.T) {
	for _, suffix := range []string{"", `,"deep":{}`, `,"deep":{"legal_name":"Old name"}`, `,"deep":{"description":null,"socials":null,"sources":null,"employees":null}`} {
		var value CompanyProfile
		if err := json.Unmarshal([]byte(`{"id":"co_222222222222","name":"Example","country":null,"website":null,"listings":[],"address":null`+suffix+`}`), &value); err != nil {
			t.Fatal(err)
		}
		if suffix == "" && value.Deep != nil {
			t.Fatal("invented deep")
		}
		if value.Deep != nil && (value.Deep.Socials != nil || value.Deep.Sources != nil || value.Deep.Founded != nil || value.Deep.Employees != nil) {
			t.Fatal("invented rich observation")
		}
	}
	var value CompanyProfile
	if err := json.Unmarshal([]byte(`{"id":"co_222222222222","name":"Example","deep":{"socials":[],"sources":[],"founded":{"value":"spring 2006","precision":"season"}}}`), &value); err != nil {
		t.Fatal(err)
	}
	if value.Deep.Socials == nil || value.Deep.Sources == nil || value.Deep.Founded.Precision != "season" {
		t.Fatal("lost empty/open fields")
	}
	var empty CompanySearch
	if err := json.Unmarshal([]byte(`{"companies":[],"next":null}`), &empty); err != nil || len(empty.Companies) != 0 || empty.Next != nil {
		t.Fatal("bad empty search", err)
	}
}

func TestCompanyDirectoryServerValidationAndOptionGuards(t *testing.T) {
	client, call := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"code":"invalid_request","message":"Choose one selector"}`))
	})
	_, err := client.CompanySearch(context.Background(), CompanySearchOptions{Query: "Example", Domain: "example.com"})
	var apiErr *Error
	if !errors.As(err, &apiErr) || apiErr.Code != "invalid_request" || call.path != "/company" {
		t.Fatal("server validation changed", err)
	}
	if _, err := client.CompanyID(context.Background(), "id", CompanyIDOptions{}, CompanyIDOptions{}); err == nil {
		t.Fatal("multiple ID options accepted")
	}
	if _, err := client.CompanySearch(context.Background(), CompanySearchOptions{}, CompanySearchOptions{}); err == nil {
		t.Fatal("multiple search options accepted")
	}
	if _, err := client.CompanyCoverage(context.Background(), CompanyCoverageOptions{}, CompanyCoverageOptions{}); err == nil {
		t.Fatal("multiple coverage options accepted")
	}
}

func TestCompanyDirectoryEmployeeObservations(t *testing.T) {
	cases := []struct {
		json                string
		count               int64
		date, scope, method string
		approximate         bool
	}{
		{`{"count":0,"as_of":"2025-12-31","scope":"legal_entity","method":"reported","approximate":false}`, 0, "2025-12-31", "legal_entity", "reported", false},
		{`{"count":12500,"as_of":"2026-06-30","scope":"consolidated_group","method":"reported","approximate":true}`, 12500, "2026-06-30", "consolidated_group", "reported", true},
		{`{"count":7,"as_of":"2026-01-15","scope":"future_scope","method":"future_method","approximate":false,"future":null}`, 7, "2026-01-15", "future_scope", "future_method", false},
	}
	for _, tc := range cases {
		profileJSON := `{"id":"co_222222222222","name":"Example","deep":{"employees":` + tc.json + `}}`
		client, _ := newTestClient(t, okJSON(profileJSON))
		profile, err := client.CompanyID(context.Background(), "co_222222222222", CompanyIDOptions{Deep: true})
		if err != nil {
			t.Fatal(err)
		}
		client, _ = newTestClient(t, okJSON(`{"companies":[`+profileJSON[:len(profileJSON)-1]+`,"match":{"field":"name"}}],"next":null}`))
		page, err := client.CompanySearch(context.Background(), CompanySearchOptions{Query: "Example", Deep: true})
		if err != nil {
			t.Fatal(err)
		}
		for _, got := range []*CompanyProfileEmployees{profile.Deep.Employees, page.Companies[0].Deep.Employees} {
			if got == nil || got.Count != tc.count || got.AsOf != tc.date || got.Scope != tc.scope || got.Method != tc.method || got.Approximate != tc.approximate {
				t.Fatalf("lost employee observation: %#v", got)
			}
		}
	}
}

func TestCompanyRegistrationsPreserveSourceRolesAndNullableFacts(t *testing.T) {
	const row = `{"authority":"RA000599","number":"0001234567","jurisdiction":{"country":"US","state":"CO"},"role":"domestic","legal_form":{"code":"DNC","name":"Domestic Non-profit Corporation"},"status":"Good Standing","formation_date":"2004-02-29","address":{"kind":"principal","line1":"12 Main St.","line2":"Suite 2","city":"Example","state":"CO","postal":"00123-0001","country_raw":"US"},"future":"retained"}`
	for _, deep := range []string{`{}`, `{"registrations":null}`, `{"registrations":[]}`, `{"registrations":[` + row + `]}`} {
		profileJSON := `{"id":"co_222222222222","name":"Example","deep":` + deep + `}`
		client, _ := newTestClient(t, okJSON(profileJSON))
		profile, err := client.CompanyID(context.Background(), "co_222222222222", CompanyIDOptions{Deep: true})
		if err != nil {
			t.Fatal(err)
		}
		client, _ = newTestClient(t, okJSON(`{"companies":[`+profileJSON[:len(profileJSON)-1]+`,"match":{"field":"identifier","value":"0001234567"}}],"next":null}`))
		page, err := client.CompanySearch(context.Background(), CompanySearchOptions{Identifier: "0001234567", Authority: "RA000599", Deep: true})
		if err != nil {
			t.Fatal(err)
		}
		for _, d := range []*CompanyProfileDeep{profile.Deep, page.Companies[0].Deep} {
			if len(d.Registrations) > 0 {
				r := d.Registrations[0]
				if r.Number != "0001234567" || r.Jurisdiction.Country != "US" || r.FormationDate == nil || *r.FormationDate != "2004-02-29" || r.Address == nil || *r.Address.Postal != "00123-0001" || r.Address.Kind != "principal" {
					t.Fatalf("lost registry facts: %#v", r)
				}
			} else if deep == `{"registrations":[]}` && d.Registrations == nil {
				t.Fatal("empty registrations collapsed to nil")
			}
		}
	}
	var future CompanyProfileDeep
	if err := json.Unmarshal([]byte(`{"registrations":[{"authority":"future","number":"0000","jurisdiction":{"country":"US","state":null},"role":"future_role","legal_form":{"code":"future","name":"Future"},"status":"future_status","formation_date":null,"address":null}]}`), &future); err != nil {
		t.Fatal(err)
	}
	if future.Registrations[0].Role != "future_role" || future.Registrations[0].Address != nil || future.Registrations[0].FormationDate != nil {
		t.Fatal("future/null registry fact changed")
	}
}
