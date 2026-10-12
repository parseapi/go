package parseapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"reflect"
	"testing"
)

func currencyFixture(t *testing.T) map[string]json.RawMessage {
	t.Helper()
	data, err := os.ReadFile("testdata/currency-reference.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture map[string]json.RawMessage
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Fatal(err)
	}
	return fixture
}

func TestCurrencyOldCallsAndNamedNullableReference(t *testing.T) {
	ctx := context.Background()
	c, call := newTestClient(t, okJSON(`{"base":"USD","quote":"EUR","rate":0.9,"date":"2026-10-09","amount":0,"converted":0,"future":true}`))
	zero := 0.0
	old := c.CurrencyRate
	rate, err := old(ctx, "USD", "EUR", CurrencyRateOptions{Date: "2026-10-09", Amount: &zero})
	if err != nil || rate.Converted == nil || *rate.Converted != 0 || rate.Stale != nil || call.rawQry != "amount=0&date=2026-10-09" {
		t.Fatal(err, rate, call)
	}
	fixture := currencyFixture(t)
	c, call = newTestClient(t, okJSON(string(fixture["pair"])))
	one := 1.0
	rate, err = c.CurrencyRate(ctx, "USD", "ROL", CurrencyRateOptions{Date: "2004-01-05", Amount: &one})
	if err != nil || rate.Converted != nil || rate.Date != "2004-01-05" || rate.Stale != nil {
		t.Fatal(err, rate)
	}
	if call.path != "/currency/USD/ROL" || call.rawQry != "amount=1&date=2004-01-05" {
		t.Fatal(call)
	}
	if rate.Dates != nil {
		t.Fatal(rate.Dates)
	}
}

func TestCurrencySourcesBatchAndHeldHistory(t *testing.T) {
	fixture := currencyFixture(t)
	ctx := context.Background()
	zero := 0.0
	quotes := []string{"EUR", "JPY", "KPW"}
	c, call := newTestClient(t, okJSON(string(fixture["batch"])))
	batch, err := c.CurrencyRates(ctx, "USD", quotes, CurrencyRatesOptions{Date: "2026-10-09", Amount: &zero})
	if err != nil || len(batch.Rates) != 2 || len(batch.Missing) != 1 || batch.Missing[0].Quote != "KPW" || batch.Rates["EUR"].Converted == nil || *batch.Rates["EUR"].Converted != 0 {
		t.Fatal(err, batch)
	}
	params, _ := url.ParseQuery(call.rawQry)
	if call.path != "/currency/USD/rates" || !reflect.DeepEqual(params, url.Values{"quotes": {"EUR,JPY,KPW"}, "date": {"2026-10-09"}, "amount": {"0"}}) {
		t.Fatal(call, params)
	}
	c, call = newTestClient(t, okJSON(string(fixture["history"])))
	history, err := c.CurrencyHistory(ctx, "USD", "EUR", "2026-10-09", "2026-10-10", CurrencyHistoryOptions{})
	if err != nil || len(history.Data) != 2 || history.Data[1].RequestedDate != "2026-10-10" || history.Data[1].Date != "2026-10-09" || call.path != "/currency/USD/EUR/history" {
		t.Fatal(err, history, call)
	}
	c, call = newTestClient(t, okJSON(string(fixture["rates_history"])))
	multi, err := c.CurrencyRatesHistory(ctx, "USD", quotes, "2026-10-09", "2026-10-10", CurrencyRatesHistoryOptions{Amount: &zero})
	if err != nil || len(multi.Data) != 4 || len(multi.Missing) != 2 || multi.Missing[1].RequestedDate != "2026-10-10" || multi.Missing[1].Reason != "no_rate" || call.path != "/currency/USD/rates/history" {
		t.Fatal(err, multi, call)
	}
	if !reflect.DeepEqual(quotes, []string{"EUR", "JPY", "KPW"}) {
		t.Fatal("mutated quote input")
	}
}

func TestCurrencyUnknownSourceHistoryAndEvidenceStayNil(t *testing.T) {
	var rate CurrencyRate
	if err := json.Unmarshal([]byte(`{"base":"USD","quote":"EUR","rate":1,"date":"2026-10-09","stale":false,"future":true}`), &rate); err != nil {
		t.Fatal(err)
	}
	if rate.Stale == nil || *rate.Stale {
		t.Fatal(rate)
	}
}

func TestCurrencyInvalidInputsMultipleOptionsAndCancellation(t *testing.T) {
	c, call := newTestClient(t, okJSON(`{}`))
	ctx := context.Background()
	for _, quotes := range [][]string{{}, {"USD", "usd"}, {"US"}, {"USD/JPY"}, make([]string, 51)} {
		if _, err := c.CurrencyRates(ctx, "USD", quotes); err == nil {
			t.Fatal("accepted invalid quotes", quotes)
		}
	}
	if _, err := c.CurrencyRates(ctx, "USD", []string{"EUR"}, CurrencyRatesOptions{}, CurrencyRatesOptions{}); err == nil {
		t.Fatal("accepted multiple rates options")
	}
	if _, err := c.CurrencyHistory(ctx, "USD", "EUR", "2026-10-09", "2026-10-10", CurrencyHistoryOptions{}, CurrencyHistoryOptions{}); err == nil {
		t.Fatal("accepted multiple history options")
	}
	if _, err := c.CurrencyRatesHistory(ctx, "USD", []string{"EUR"}, "2026-10-09", "2026-10-10", CurrencyRatesHistoryOptions{}, CurrencyRatesHistoryOptions{}); err == nil {
		t.Fatal("accepted multiple multi options")
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := c.CurrencyRatesHistory(cancelled, "USD", []string{"JPY"}, "2026-10-09", "2026-10-10"); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	if call.path != "" {
		t.Fatal("rejected call was dispatched", call)
	}
}

func TestCurrencyAllQuotesCatalogAnd(t *testing.T) {
	ctx := context.Background()
	fixture := currencyFixture(t)
	c, call := newTestClient(t, okJSON(string(fixture["batch"])))
	rates, err := c.CurrencyRates(ctx, "USD", nil, CurrencyRatesOptions{})
	if err != nil || len(rates.Rates) != 2 || rates.Rates["EUR"].Date != "2026-10-09" || call.rawQry != "" {
		t.Fatal(err, rates, call)
	}
	if _, err := c.CurrencyRates(ctx, "USD", []string{}); err == nil {
		t.Fatal("explicit empty quotes must not mean all")
	}
	bytes, err := os.ReadFile("testdata/currency-catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	var extra map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &extra); err != nil {
		t.Fatal(err)
	}
	c, call = newTestClient(t, okJSON(string(extra["currencies"])))
	catalog, err := c.CurrencyList(ctx, CurrencyListOptions{Q: "gold", Country: "UA", Deep: true, Pretty: true, Lang: "fr"})
	if err != nil || len(catalog.Data) != 0 || call.path != "/currency" || call.rawQry != "country=UA&deep=true&lang=fr&pretty=true&q=gold" {
		t.Fatal(err, catalog, call)
	}
}

func TestCurrencyCatalogNonCashUnits(t *testing.T) {
	fixture := currencyFixture(t)
	c, _ := newTestClient(t, okJSON(string(fixture["currencies"])))
	result, err := c.CurrencyList(context.Background(), CurrencyListOptions{Q: "gold", Deep: true})
	if err != nil || len(result.Data) != 2 || result.Data[0].Unit != "troy_ounce" || result.Data[0].Digits != nil {
		t.Fatal(err, result)
	}
	ordinary := CurrencyResource{Currency: Currency{Currency: "USD"}}
	encoded, err := json.Marshal(ordinary)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &fields); err != nil {
		t.Fatal(err)
	}
	if _, present := fields["unit"]; present {
		t.Fatal("ordinary currency must omit unit", string(encoded))
	}
}

func TestCurrencyIndexedUnits(t *testing.T) {
	fixture := currencyFixture(t)
	c, _ := newTestClient(t, okJSON(string(fixture["indexed_pair"])))
	one := 1.0
	pair, err := c.CurrencyRate(context.Background(), "USD", "MXV", CurrencyRateOptions{Amount: &one})
	if err != nil || pair.Converted == nil || *pair.Converted != 2.06 || pair.Units.Quote != "indexed_unit" {
		t.Fatal(err, pair)
	}
}

func TestCurrencySampledHistoryAndCompletePublicationSummary(t *testing.T) {
	data, err := os.ReadFile("testdata/currency-summary.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture map[string]json.RawMessage
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	c, call := newTestClient(t, okJSON(string(fixture["sampled_history"])))
	history, err := c.CurrencyHistory(ctx, "USD", "JPY", "1999-01-01", "2026-10-09", CurrencyHistoryOptions{Interval: "month"})
	if err != nil || history.Interval != "month" || len(history.Data) != 334 || !reflect.DeepEqual(history.Missing, []string{"1999-01-01"}) {
		t.Fatal(err, history)
	}
	params, _ := url.ParseQuery(call.rawQry)
	if params.Get("interval") != "month" {
		t.Fatal(call)
	}
	c, call = newTestClient(t, okJSON(string(fixture["sampled_rates_history"])))
	_, err = c.CurrencyRatesHistory(ctx, "USD", []string{"JPY"}, "1999-01-01", "2026-10-09", CurrencyRatesHistoryOptions{Interval: "week"})
	if err != nil {
		t.Fatal(err)
	}
	params, _ = url.ParseQuery(call.rawQry)
	if params.Get("interval") != "week" || call.path != "/currency/USD/rates/history" {
		t.Fatal(call)
	}
	c, call = newTestClient(t, okJSON(string(fixture["summary"])))
	result, err := c.CurrencyRatesSummary(ctx, "USD", []string{"EUR"}, "1999-01-01", "1999-01-31", CurrencyRatesSummaryOptions{Interval: "month"})
	var expected CurrencyRatesSummary
	if err := json.Unmarshal(fixture["summary"], &expected); err != nil {
		t.Fatal(err)
	}
	if err != nil || !reflect.DeepEqual(*result, expected) || result.Data[0].ObservationCount != 20 || result.Data[0].MissingDayCount != 3 || result.Data[0].Stale != nil {
		t.Fatal(err, result)
	}
	params, _ = url.ParseQuery(call.rawQry)
	if call.path != "/currency/USD/rates/summary" || !reflect.DeepEqual(params, url.Values{"quotes": {"EUR"}, "from": {"1999-01-01"}, "to": {"1999-01-31"}, "interval": {"month"}}) {
		t.Fatal(call)
	}
}
