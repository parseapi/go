package parseapi

import (
	"context"
	"strconv"
)

// CompanyIDOptions adds directory detail in the same pooled request on every plan.
type CompanyIDOptions struct {
	_    [0]func()
	Deep bool
}

// CompanySearchOptions uses at most one identity selector, or discovers by Country, exact Industry or selected registration.
// The API validates selector combinations. Exchange narrows Ticker; Authority
// narrows Identifier. Reuse Cursor with the same selector, filters and Limit.
// Deep adds detail inside each profile without an additional lookup.
type CompanySearchOptions struct {
	_          [0]func()
	Query      string
	Domain     string
	Ticker     string
	Identifier string
	Country    string
	// Industry is an exact four-digit SIC code; pair it with IndustryType.
	Industry string
	// IndustryType is an open namespace string, currently sic.
	IndustryType string
	// RegistrationAuthority scopes selected registration evidence; form/status require it.
	RegistrationAuthority string
	// RegistrationForm is an exact source legal-form code, not ownership or tax status.
	RegistrationForm string
	// RegistrationStatus is an exact administrative status, not operating activity.
	RegistrationStatus string
	Exchange           string
	Authority          string
	Limit              int
	Cursor             string
	Deep               bool
}

// CompanyCoverageOptions configures the edition coverage request.
type CompanyCoverageOptions struct{ _ [0]func() }

// CompanyID retrieves an existing directory profile by its stable co_ ID.
// National company-number validation remains Company.
func (c *Client) CompanyID(ctx context.Context, id string, options ...CompanyIDOptions) (*CompanyProfile, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &CompanyProfile{}
	if err := c.get(ctx, "/company/id/"+seg(id), values("deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CompanySearch returns one bounded page. No match returns an empty companies
// array; service errors remain errors. One matching domain is not ownership proof.
func (c *Client) CompanySearch(ctx context.Context, options ...CompanySearchOptions) (*CompanySearch, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep, limit := "", ""
	if opts.Deep {
		deep = "true"
	}
	if opts.Limit != 0 {
		limit = strconv.Itoa(opts.Limit)
	}
	out := &CompanySearch{}
	if err := c.get(ctx, "/company", values("q", opts.Query, "domain", opts.Domain, "ticker", opts.Ticker,
		"identifier", opts.Identifier, "country", opts.Country, "industry", opts.Industry, "industry_type", opts.IndustryType, "registration_authority", opts.RegistrationAuthority,
		"registration_form", opts.RegistrationForm, "registration_status", opts.RegistrationStatus, "exchange", opts.Exchange,
		"authority", opts.Authority, "limit", limit, "cursor", opts.Cursor, "deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CompanyCoverage returns counts for the directory edition, not global completeness.
func (c *Client) CompanyCoverage(ctx context.Context, options ...CompanyCoverageOptions) (*CompanyCoverage, error) {
	if _, err := oneOption(options); err != nil {
		return nil, err
	}
	out := &CompanyCoverage{}
	if err := c.get(ctx, "/company/directory/coverage", nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
