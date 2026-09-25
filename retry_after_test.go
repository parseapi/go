package parseapi

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestLongRetryAfterReturnsOriginalErrorWithoutWaiting(t *testing.T) {
	for _, header := range []string{"60", strings.Repeat("9", 400), time.Now().Add(time.Minute).UTC().Format(http.TimeFormat)} {
		var calls atomic.Int32
		client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			calls.Add(1)
			w.Header().Set("Retry-After", header)
			w.WriteHeader(429)
			_, _ = w.Write([]byte(`{"code":"rate_limited","message":"Wait for reset","docs":"https://parseapi.com/docs","request_id":"req_fixture"}`))
		}, WithRetries(2))
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err := client.Card(ctx, "001234")
		cancel()
		var apiErr *Error
		if !errors.As(err, &apiErr) || apiErr.Status != 429 || apiErr.Code != "rate_limited" || apiErr.Message != "Wait for reset" || apiErr.RequestID != "req_fixture" || apiErr.Docs != "https://parseapi.com/docs" || apiErr.RetryAfter == nil || *apiErr.RetryAfter != header || calls.Load() != 1 {
			t.Fatalf("error=%#v calls=%d", err, calls.Load())
		}
	}
}

func TestRetryAfterShortMalformedAndMissing(t *testing.T) {
	if got := retryDelay(0, "2"); got != 2*time.Second {
		t.Fatalf("short wait: %v", got)
	}
	if got := retryDelay(0, "0"); got != 0 {
		t.Fatalf("zero: %v", got)
	}
	if got := retryDelay(0, time.Now().Add(3*time.Second).UTC().Format(http.TimeFormat)); got < time.Second || got > 3*time.Second {
		t.Fatalf("date: %v", got)
	}
	for _, header := range []string{"", "nonsense", "-1", "NaN", "Infinity", "1e999"} {
		if got := retryDelay(0, header); got < 0 || got > 250*time.Millisecond {
			t.Fatalf("malformed %q: %v", header, got)
		}
	}
	if got := retryDelay(100000, ""); got < 0 || got > 5*time.Second {
		t.Fatalf("backoff budget: %v", got)
	}
}

func TestRetryAfterMetadataOnExhaustedOrDisabledRetries(t *testing.T) {
	for _, retries := range []int{0, 1} {
		var calls atomic.Int32
		client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			calls.Add(1)
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(429)
			_, _ = w.Write([]byte(`{"code":"rate_limited"}`))
		}, WithRetries(retries))
		_, err := client.Card(context.Background(), "001234")
		var apiErr *Error
		if !errors.As(err, &apiErr) || apiErr.RetryAfter == nil || *apiErr.RetryAfter != "0" || calls.Load() != int32(retries+1) {
			t.Fatalf("%#v calls=%d", err, calls.Load())
		}
	}
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(400) })
	_, err := client.Card(context.Background(), "001234")
	var apiErr *Error
	if !errors.As(err, &apiErr) || apiErr.RetryAfter != nil {
		t.Fatalf("absent header: %#v", err)
	}
}
