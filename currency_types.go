package parseapi

// CurrencyRateDates qualifies original observations when they differ from the result date.
type CurrencyRateDates struct {
	_     [0]func()
	Base  []string `json:"base,omitempty"`
	Quote []string `json:"quote,omitempty"`
}
type CurrencyRateUnits struct {
	_     [0]func()
	Base  string `json:"base,omitempty"`
	Quote string `json:"quote,omitempty"`
}
type CurrencyRateMissing struct {
	_      [0]func()
	Quote  string `json:"quote"`
	Reason string `json:"reason"`
}
type CurrencyQuote struct {
	_         [0]func()
	Rate      float64            `json:"rate"`
	Date      string             `json:"date"`
	Converted *float64           `json:"converted,omitempty"`
	Stale     *bool              `json:"stale,omitempty"`
	Dates     *CurrencyRateDates `json:"dates,omitempty"`
	Units     *CurrencyRateUnits `json:"units,omitempty"`
}
type CurrencyObservation struct {
	_             [0]func()
	RequestedDate string             `json:"requested_date"`
	Rate          float64            `json:"rate"`
	Date          string             `json:"date"`
	Converted     *float64           `json:"converted,omitempty"`
	Stale         *bool              `json:"stale,omitempty"`
	Dates         *CurrencyRateDates `json:"dates,omitempty"`
	Units         *CurrencyRateUnits `json:"units,omitempty"`
}
type CurrencyQuoteObservation struct {
	CurrencyObservation
	Quote string `json:"quote"`
}
type CurrencyRates struct {
	_             [0]func()
	Base          string                   `json:"base"`
	Amount        *float64                 `json:"amount,omitempty"`
	RequestedDate *string                  `json:"requested_date,omitempty"`
	Rates         map[string]CurrencyQuote `json:"rates"`
	Missing       []CurrencyRateMissing    `json:"missing"`
}
type CurrencyHistory struct {
	_        [0]func()
	Base     string                `json:"base"`
	Quote    string                `json:"quote"`
	From     string                `json:"from"`
	To       string                `json:"to"`
	Amount   *float64              `json:"amount,omitempty"`
	Data     []CurrencyObservation `json:"data"`
	Missing  []string              `json:"missing"`
	Interval string                `json:"interval,omitempty"`
}
type CurrencyHistoryMissing struct {
	CurrencyRateMissing
	RequestedDate string `json:"requested_date"`
}
type CurrencyRatesHistory struct {
	_        [0]func()
	Base     string                     `json:"base"`
	Quotes   []string                   `json:"quotes"`
	From     string                     `json:"from"`
	To       string                     `json:"to"`
	Amount   *float64                   `json:"amount,omitempty"`
	Data     []CurrencyQuoteObservation `json:"data"`
	Missing  []CurrencyHistoryMissing   `json:"missing"`
	Interval string                     `json:"interval,omitempty"`
}

type CurrencyResource struct {
	Currency
	Unit string `json:"unit,omitempty"`
}

type CurrencyList struct {
	_    [0]func()
	Data []CurrencyResource `json:"data"`
}

// CurrencyRateSummary counts original pair observation days without weighting carried dates again.
type CurrencyRateSummary struct {
	_                [0]func()
	Quote            string             `json:"quote"`
	From             string             `json:"from"`
	To               string             `json:"to"`
	Average          float64            `json:"average"`
	Min              float64            `json:"min"`
	Max              float64            `json:"max"`
	ObservationCount int                `json:"observation_count"`
	FirstDate        string             `json:"first_date"`
	LastDate         string             `json:"last_date"`
	MissingDayCount  int                `json:"missing_day_count"`
	Stale            *bool              `json:"stale,omitempty"`
	Units            *CurrencyRateUnits `json:"units,omitempty"`
}
type CurrencySummaryMissing struct {
	_               [0]func()
	Quote           string `json:"quote"`
	From            string `json:"from"`
	To              string `json:"to"`
	Reason          string `json:"reason"`
	MissingDayCount int    `json:"missing_day_count"`
}
type CurrencyRatesSummary struct {
	_        [0]func()
	Base     string                   `json:"base"`
	Quotes   []string                 `json:"quotes"`
	From     string                   `json:"from"`
	To       string                   `json:"to"`
	Interval string                   `json:"interval"`
	Data     []CurrencyRateSummary    `json:"data"`
	Missing  []CurrencySummaryMissing `json:"missing"`
}
