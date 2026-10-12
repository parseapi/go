package parseapi_test

import (
	"context"
	parseapi "github.com/parseapi/go"
)

// Keep the old method reference callable alongside the optional date and amount options.
var _ func(*parseapi.Client, context.Context, string, string, ...parseapi.CurrencyRateOptions) (*parseapi.CurrencyRate, error) = (*parseapi.Client).CurrencyRate
var _ func(*parseapi.Client, context.Context, string, []string, ...parseapi.CurrencyRatesOptions) (*parseapi.CurrencyRates, error) = (*parseapi.Client).CurrencyRates
var _ func(*parseapi.Client, context.Context, string, string, string, string, ...parseapi.CurrencyHistoryOptions) (*parseapi.CurrencyHistory, error) = (*parseapi.Client).CurrencyHistory
var _ func(*parseapi.Client, context.Context, string, []string, string, string, ...parseapi.CurrencyRatesHistoryOptions) (*parseapi.CurrencyRatesHistory, error) = (*parseapi.Client).CurrencyRatesHistory

var _ func(*parseapi.Client, context.Context, ...parseapi.CurrencyListOptions) (*parseapi.CurrencyList, error) = (*parseapi.Client).CurrencyList
