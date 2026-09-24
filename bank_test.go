package parseapi

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestBankChecksAndRawInput(t *testing.T) {
	data, err := os.ReadFile("testdata/bank-fixtures.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		Records []json.RawMessage `json:"records"`
		Inputs  []string          `json:"inputs"`
	}
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Fatal(err)
	}
	for i, body := range fixture.Records {
		var expected Bank
		if err := json.Unmarshal(body, &expected); err != nil {
			t.Fatal(err)
		}
		switch i {
		case 0:
			if expected.Checks != nil || expected.Issues != nil {
				t.Fatal("older fields became populated")
			}
		case 1:
			if expected.Checks == nil || expected.Checks.National != "not_supported" || expected.Issues == nil || len(expected.Issues) != 0 {
				t.Fatal("current checks or empty issues lost")
			}
			if expected.Deep == nil || expected.Deep.Account == nil || *expected.Deep.Account != "0532013000" {
				t.Fatal("leading zero lost")
			}
		case 2:
			if expected.Valid || expected.Checks.Input != "failed" || expected.Checks.Checksum != "not_checked" || expected.Issues[0].Code != "invalid_characters" {
				t.Fatal("invalid result changed")
			}
		case 3:
			if expected.Checks.National != "future-state" || expected.Issues[0].Field != "future-field" || expected.Issues[0].Code != "future-code" {
				t.Fatal("future string values lost")
			}
		}
		client, calls := newTestClient(t, okJSON(string(body)))
		for _, input := range fixture.Inputs {
			result, err := client.Bank(context.Background(), input, BankOptions{Deep: true})
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(result, &expected) {
				t.Fatalf("response mismatch: %#v", result)
			}
			var sent map[string]any
			if err := json.Unmarshal(calls.body, &sent); err != nil {
				t.Fatal(err)
			}
			if calls.path != "/bank" || calls.rawQry != "" || calls.method != "POST" || sent["iban"] != input || sent["deep"] != true || calls.header.Get("Content-Type") != "application/json" || calls.header.Get("Parse-Version") != "2.0.0" {
				t.Fatal("Bank did not preserve its private JSON request")
			}
		}
	}
}

func TestBankDomesticRequirementsAndDirectory(t *testing.T) {
	ctx := context.Background()
	body := `{"format":"us_ach","country":"US","routing":"011000015","account":"00aB-9","valid":true,"bank_name":null,"checks":{"routing_format":"passed","routing_checksum":"passed","account_format":"future-state","account_checksum":"not_supported"},"issues":[{"field":"future-field","code":"future-code","message":"Static issue"}],"future":true}`
	client, calls := newTestClient(t, okJSON(body))
	raw := BankUSACHInput{Routing: "\t011-000-015", Account: " 00aB-%20\uFEFF"}
	result, err := client.BankUSACH(ctx, raw)
	if err != nil {
		t.Fatal(err)
	}
	var sent map[string]string
	if err := json.Unmarshal(calls.body, &sent); err != nil {
		t.Fatal(err)
	}
	if calls.path != "/bank" || calls.rawQry != "" || calls.method != "POST" || sent["account"] != raw.Account || sent["routing"] != raw.Routing || sent["country"] != "US" || sent["format"] != "us_ach" {
		t.Fatal("domestic body changed")
	}
	if result.Checks.AccountFormat != "future-state" || result.BankName != nil || result.Issues[0].Code != "future-code" {
		t.Fatal("response semantics lost")
	}
	client, calls = newTestClient(t, okJSON(`{"country":"US","format":"future-format","supported":false,"fields":[],"checks":{"future":"not_supported"},"limitations":["No verification"]}`))
	requirements, err := client.BankRequirements(ctx, "US", BankRequirementsOptions{Format: "future-format"})
	if err != nil {
		t.Fatal(err)
	}
	if calls.path != "/bank/requirements" || calls.rawQry != "country=US&format=future-format" || calls.method != "GET" || len(calls.body) != 0 || requirements.Checks["future"] != "not_supported" || requirements.Supported {
		t.Fatal("requirements changed")
	}
	var bank Bank
	if err := json.Unmarshal([]byte(`{"deep":{"directory":{"edition":"edition","country":"FR","match":"future-grain"}}}`), &bank); err != nil {
		t.Fatal(err)
	}
	if bank.Deep.Directory.Match != "future-grain" {
		t.Fatal("directory evidence lost")
	}
}
