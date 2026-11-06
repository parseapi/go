package parseapi_test

import (
	"encoding/json"
	parseapi "github.com/parseapi/go"
	"testing"
)

func TestEmailEnrichmentCompatibility(t *testing.T) {
	var enriched parseapi.Email
	if err := json.Unmarshal([]byte(`{"email":"jane.doe+news@example.com","valid":true,"free":false,"role":false,"disposable":false,"future":true,"first_name":"Jane","no_reply":false,"tag":"news","mail_provider":"future-provider","deliverable":true,"catchall":false,"status":"future-status","reason":"future_reason"}`), &enriched); err != nil {
		t.Fatal(err)
	}
	if enriched.FirstName == nil || *enriched.FirstName != "Jane" || enriched.NoReply == nil || *enriched.NoReply || enriched.Tag == nil || *enriched.Tag != "news" || enriched.MailProvider == nil || *enriched.MailProvider != "future-provider" {
		t.Fatalf("lost address hints: %#v", enriched)
	}
	if enriched.Status == nil || *enriched.Status != "future-status" || enriched.Reason == nil || *enriched.Reason != "future_reason" {
		t.Fatalf("lost deep enrichment: %#v", enriched)
	}
	for _, body := range []string{`{"email":"a@example.com"}`, `{"email":"a@example.com"}`, `{"email":"a@example.com","first_name":null,"no_reply":null,"tag":null,"mail_provider":null,"status":null,"reason":null}`, `{"email":"a@example.com","deliverable":false,"catchall":true}`} {
		var result parseapi.Email
		if err := json.Unmarshal([]byte(body), &result); err != nil {
			t.Fatal(err)
		}
		if result.FirstName != nil || result.NoReply != nil || result.Tag != nil || result.MailProvider != nil {
			t.Fatalf("invented enrichment: %#v", result)
		}
		if result.Status != nil || result.Reason != nil {
			t.Fatalf("invented deep enrichment: %#v", result)
		}
	}
}
