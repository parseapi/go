package parseapi

import (
	"encoding/json"
	"testing"
)

func TestIPDeepExpansion(t *testing.T) {
	var value IP
	if err := json.Unmarshal([]byte(`{"ip":"198.51.100.10","country":"XX","country_name":"Example","continent":"NA","asn":"AS64512","asn_name":"Fixture","city":"Example","state":"EX","state_name":"Example state","postal":"00123","latitude":0,"longitude":-0.000001,"timezone":"Etc/UTC","network":"198.51.100.0/24","asn_type":"education","registry":"arin","datacenter":false,"relay":false,"tor":false,"vpn":false,"proxy":false,"provider":null}`), &value); err != nil {
		t.Fatal(err)
	}
	d := &value
	if d == nil || d.Latitude == nil || *d.Latitude != 0 || d.Longitude == nil || *d.Longitude != -0.000001 || d.Postal == nil || *d.Postal != "00123" || d.Proxy == nil || *d.Proxy || d.StateName == nil || *d.StateName != "Example state" || d.Timezone == nil || *d.Timezone != "Etc/UTC" || d.Network == nil || *d.Network != "198.51.100.0/24" || d.ASNType == nil || *d.ASNType != "education" {
		t.Fatal("lost new Deep facts")
	}
	for _, text := range []string{`{"ip":"198.51.100.10"}`, `{"ip":"198.51.100.10"}`, `{"ip":"198.51.100.10","latitude":null,"longitude":null,"proxy":null}`} {
		var old IP
		if err := json.Unmarshal([]byte(text), &old); err != nil {
			t.Fatal(err)
		}
		if old.Latitude != nil || old.Proxy != nil {
			t.Fatal("unknown became a fact")
		}
	}
}
