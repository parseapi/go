package parseapi

import "context"

// DomainRegistrationResult contains partial registration detail from the same observation.
// Nil means unknown. A non-nil empty Status is a known empty list.
type DomainRegistrationResult struct {
	_          [0]func()
	Registered *bool    `json:"registered"`
	Created    *string  `json:"created"`
	Updated    *string  `json:"updated"`
	Expires    *string  `json:"expires"`
	Registrar  *string  `json:"registrar"`
	Status     []string `json:"status"`
	DNSSEC     *bool    `json:"dnssec"`
}

// DomainResultDeep is omitted when not requested and empty when detail is unavailable.
type DomainResultDeep struct {
	_            [0]func()
	Registration *DomainRegistrationResult `json:"registration,omitempty"`
}

// DomainResult preserves unknown availability instead of reducing it to a boolean.
type DomainResult struct {
	_      [0]func()
	Domain string `json:"domain"`
	// True means unregistered, false means registered, nil means unknown.
	Available *bool             `json:"available"`
	Deep      *DomainResultDeep `json:"deep,omitempty"`
}

// DomainResult checks registration with explicit unknown availability and partial detail.
// It makes the same lookup as Domain and accepts the same options.
func (c *Client) DomainResult(ctx context.Context, domain string, options ...DomainOptions) (*DomainResult, error) {
	opts, err := oneOption(options)
	if err != nil {
		return nil, err
	}
	deep := ""
	if opts.Deep {
		deep = "true"
	}
	out := &DomainResult{}
	if err := c.get(ctx, "/domain/"+seg(domain), values("deep", deep), nil, out); err != nil {
		return nil, err
	}
	return out, nil
}
