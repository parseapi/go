package parseapi

import (
	"context"
	"testing"
)

func TestDNSRecordsAndQuestion(t *testing.T) {
	client, call := newTestClient(t, okJSON(`{"domain":"example.com","records":[{"name":"example.com.","type":"TXT","ttl":0,"value":"\"one\" \"two\"","future":true},{"name":"alias.example.","type":"CNAME","ttl":300,"value":"target.example."}],"future":true}`))
	result, err := client.DNS(context.Background(), "_dmarc.bücher.example.", DNSOptions{Type: "txt"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 2 || result.Records[0].TTL != 0 || result.Records[0].Value != `"one" "two"` || result.Records[1].Type != "CNAME" {
		t.Fatalf("lost DNS presentation: %#v", result)
	}
	if call.path != "/dns/_dmarc.b%C3%BCcher.example." || call.rawQry != "type=txt" {
		t.Fatalf("bad DNS request: %#v", call)
	}
	_, _ = client.DNS(context.Background(), "example.com")
	if call.rawQry != "" {
		t.Fatalf("unexpected default: %#v", call)
	}
	if _, err := client.DNS(context.Background(), "example.com", DNSOptions{}, DNSOptions{}); err == nil {
		t.Fatal("accepted multiple options")
	}
	client, _ = newTestClient(t, okJSON(`{"domain":"example.com","records":[]}`))
	empty, err := client.DNS(context.Background(), "example.com")
	if err != nil || empty.Records == nil || len(empty.Records) != 0 {
		t.Fatalf("lost empty records: %#v %v", empty, err)
	}
}
