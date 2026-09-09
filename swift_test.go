package parseapi

import (
	"context"
	"testing"
)

func TestSWIFTUnknownIdentityAndEncoding(t *testing.T) {
	client, call := newTestClient(t, okJSON(`{"swift":"ZZZZUS33","valid":true,"country":"US","name":null,"future":true}`))
	result, err := client.SWIFT(context.Background(), " bofa/us3n? ")
	if err != nil || !result.Valid || result.Name != nil || result.Country == nil || *result.Country != "US" {
		t.Fatalf("lost unknown identity: %#v %v", result, err)
	}
	if call.path != "/swift/%20bofa%2Fus3n%3F%20" || call.rawQry != "" {
		t.Fatalf("bad request: %#v", call)
	}
	if _, err := client.SWIFT(context.Background(), "BOFAUS3N", SWIFTOptions{}, SWIFTOptions{}); err == nil {
		t.Fatal("accepted multiple options")
	}
}
