package parseapi_test

import (
	"encoding/json"
	parseapi "github.com/parseapi/go"
	"testing"
)

func TestEmailEnrichmentCompatibility(t *testing.T) {
	var enriched parseapi.Email
	if err := json.Unmarshal([]byte(`{"email":"jane.doe+news@example.com","valid":true,"free":false,"role":false,"disposable":false,"deep": {"first_name":"Jane","no_reply":false,"tag":"news","mail_provider":"future-provider","deliverable":true,"catchall":false,"status":"future-status","reason":"future_reason"},"future":true}`), &enriched); err != nil {
		t.Fatal(err)
	}
	if enriched.Deep == nil || enriched.Deep.FirstName == nil || *enriched.Deep.FirstName != "Jane" || enriched.Deep.NoReply == nil || *enriched.Deep.NoReply || enriched.Deep.Tag == nil || *enriched.Deep.Tag != "news" || enriched.Deep.MailProvider == nil || *enriched.Deep.MailProvider != "future-provider" {
		t.Fatalf("lost address hints: %#v", enriched)
	}
	if enriched.Deep == nil || enriched.Deep.Status == nil || *enriched.Deep.Status != "future-status" || enriched.Deep.Reason == nil || *enriched.Deep.Reason != "future_reason" {
		t.Fatalf("lost deep enrichment: %#v", enriched.Deep)
	}
	for _, body := range []string{`{"email":"a@example.com"}`, `{"email":"a@example.com","deep":{}}`, `{"email":"a@example.com","deep": {"first_name":null,"no_reply":null,"tag":null,"mail_provider":null,"status":null,"reason":null}}`, `{"email":"a@example.com","deep":{"deliverable":false,"catchall":true}}`} {
		var result parseapi.Email
		if err := json.Unmarshal([]byte(body), &result); err != nil {
			t.Fatal(err)
		}
		if result.Deep != nil && (result.Deep.FirstName != nil || result.Deep.NoReply != nil || result.Deep.Tag != nil || result.Deep.MailProvider != nil) {
			t.Fatalf("invented enrichment: %#v", result)
		}
		if result.Deep != nil && (result.Deep.Status != nil || result.Deep.Reason != nil) {
			t.Fatalf("invented deep enrichment: %#v", result.Deep)
		}
	}
}
