package parseapi

import (
	"context"
	"errors"
	"net/url"
	"strings"
)

func currencyBool(value bool) string {
	if value {
		return "true"
	}
	return ""
}

func currencyQuotes(quotes []string) (string, error) {
	if len(quotes) < 1 || len(quotes) > 50 {
		return "", errors.New("parseapi: Currency quotes must contain 1 to 50 distinct three-letter codes")
	}
	seen := make(map[string]bool, len(quotes))
	for _, code := range quotes {
		if len(code) != 3 {
			return "", errors.New("parseapi: Currency quotes must contain three-letter codes")
		}
		for _, letter := range code {
			if !((letter >= 'A' && letter <= 'Z') || (letter >= 'a' && letter <= 'z')) {
				return "", errors.New("parseapi: Currency quotes must contain three-letter codes")
			}
		}
		upper := strings.ToUpper(code)
		if seen[upper] {
			return "", errors.New("parseapi: Currency quotes must be distinct")
		}
		seen[upper] = true
	}
	return strings.Join(quotes, ","), nil
}

func currencyHistoryQuery(from, to string, amount *float64, pretty bool) url.Values {
	value := ""
	if amount != nil {
		value = f(*amount)
	}
	return values("from", from, "to", to, "amount", value, "pretty", currencyBool(pretty))
}

// CurrencyRatesOptions configures a multi-quote lookup keyed by currency code.
type CurrencyRatesOptions struct {
	_      [0]func()
	Date   string
	Amount *float64
	Pretty bool
}

// CurrencyRates accepts nil quotes for all available maintained rates, or 1–50 distinct codes.
func (c *Client) CurrencyRates(ctx context.Context, base string, quotes []string, options ...CurrencyRatesOptions) (*CurrencyRates, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	query := currencyHistoryQuery("", "", opts.Amount, opts.Pretty)
	if quotes != nil {
		codes, err := currencyQuotes(quotes)
		if err != nil {
			return nil, err
		}
		query.Set("quotes", codes)
	}
	if opts.Date != "" {
		query.Set("date", opts.Date)
	}
	out := &CurrencyRates{}
	if err := c.get(ctx, "/currency/"+seg(base)+"/rates", query, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CurrencyHistoryOptions configures CurrencyHistory.
type CurrencyHistoryOptions struct {
	_        [0]func()
	Amount   *float64
	Pretty   bool
	Interval string
}

// CurrencyHistory reads daily held history, or week/month samples anchored to from with to included.
func (c *Client) CurrencyHistory(ctx context.Context, base string, quote string, from string, to string, options ...CurrencyHistoryOptions) (*CurrencyHistory, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	query := currencyHistoryQuery(from, to, opts.Amount, opts.Pretty)
	if opts.Interval != "" {
		query.Set("interval", opts.Interval)
	}
	out := &CurrencyHistory{}
	if err := c.get(ctx, "/currency/"+seg(base)+"/"+seg(quote)+"/history", query, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CurrencyRatesHistoryOptions configures CurrencyRatesHistory.
type CurrencyRatesHistoryOptions struct {
	_        [0]func()
	Amount   *float64
	Pretty   bool
	Interval string
}

// CurrencyRatesHistory reads up to 5000 sampled quote-dates including gaps under a bounded complete read.
func (c *Client) CurrencyRatesHistory(ctx context.Context, base string, quotes []string, from string, to string, options ...CurrencyRatesHistoryOptions) (*CurrencyRatesHistory, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	codes, err := currencyQuotes(quotes)
	if err != nil {
		return nil, err
	}
	query := currencyHistoryQuery(from, to, opts.Amount, opts.Pretty)
	if opts.Interval != "" {
		query.Set("interval", opts.Interval)
	}
	query.Set("quotes", codes)
	out := &CurrencyRatesHistory{}
	if err := c.get(ctx, "/currency/"+seg(base)+"/rates/history", query, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}





// CurrencyListOptions filters maintained currency metadata; rate coverage is separate.
type CurrencyListOptions struct {
	_       [0]func()
	Q       string
	Country string
	Deep    bool
	Pretty  bool
	// Lang selects translated display names for this request.
	Lang string
}

func (c *Client) CurrencyList(ctx context.Context, options ...CurrencyListOptions) (*CurrencyList, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	out := &CurrencyList{}
	if err := c.get(ctx, "/currency", values("q", opts.Q, "country", opts.Country, "deep", currencyBool(opts.Deep), "pretty", currencyBool(opts.Pretty), "lang", opts.Lang), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CurrencyRatesSummaryOptions configures observation summaries.
type CurrencyRatesSummaryOptions struct {
	_        [0]func()
	Interval string
	Pretty   bool
}

// CurrencyRatesSummary returns period averages, extrema, observation counts and calendar gaps.
func (c *Client) CurrencyRatesSummary(ctx context.Context, base string, quotes []string, from string, to string, options ...CurrencyRatesSummaryOptions) (*CurrencyRatesSummary, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	codes, err := currencyQuotes(quotes)
	if err != nil {
		return nil, err
	}
	query := values("quotes", codes, "from", from, "to", to, "interval", opts.Interval, "pretty", currencyBool(opts.Pretty))
	out := &CurrencyRatesSummary{}
	if err := c.get(ctx, "/currency/"+seg(base)+"/rates/summary", query, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
