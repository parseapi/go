package parseapi

// Company directory responses. National Company and CompanyDeep remain separate.

// CompanyProfileAddress is an address with its recorded role; role does not imply mailing validity or headquarters.
type CompanyProfileAddress struct {
	_       [0]func()
	Type    string  `json:"type"`
	Street  *string `json:"street,omitempty"`
	City    *string `json:"city,omitempty"`
	State   *string `json:"state,omitempty"`
	Postal  *string `json:"postal,omitempty"`
	Country *string `json:"country,omitempty"`
}

// CompanyProfileListing is a reported exchange/symbol pair. No listings does not establish private ownership.
type CompanyProfileListing struct {
	_        [0]func()
	Exchange string `json:"exchange"`
	Symbol   string `json:"symbol"`
}

// CompanyProfileJurisdiction is the recorded registration jurisdiction, separate from address or operating location.
type CompanyProfileJurisdiction struct {
	_       [0]func()
	Country string  `json:"country"`
	State   *string `json:"state,omitempty"`
}

// CompanyProfileWebsite is another associated hostname and its recorded URL, when known.
type CompanyProfileWebsite struct {
	_      [0]func()
	Domain string  `json:"domain"`
	URL    *string `json:"url,omitempty"`
}

// CompanyProfileIdentifier is an authority-scoped identifier; values preserve leading zeros.
type CompanyProfileIdentifier struct {
	_         [0]func()
	Type      string `json:"type"`
	Authority string `json:"authority"`
	Value     string `json:"value"`
}

// CompanyProfileIndustry is a reported classification; type is an open scheme string.
type CompanyProfileIndustry struct {
	_    [0]func()
	Type string  `json:"type"`
	Code string  `json:"code"`
	Name *string `json:"name,omitempty"`
}

// CompanyProfileFounding holds the reported founding value and precision (year, month or day), distinct from incorporation.
type CompanyProfileFounding struct {
	_     [0]func()
	Value string `json:"value"`
	// Open string; currently year, month or day. Preserve the source value without padding.
	Precision string `json:"precision"`
}

// CompanyProfileEmployees is a reported total headcount at its explicit measurement date.
type CompanyProfileEmployees struct {
	_     [0]func()
	Count int64  `json:"count"`
	AsOf  string `json:"as_of"`
	// Open string, currently legal_entity or consolidated_group.
	Scope string `json:"scope"`
	// Open string, currently reported.
	Method      string `json:"method"`
	Approximate bool   `json:"approximate"`
}

// CompanyProfileRegistrationLegalForm holds a register's open legal-form code and label.
type CompanyProfileRegistrationLegalForm struct {
	_    [0]func()
	Code string `json:"code"`
	Name string `json:"name"`
}

// CompanyProfileRegistrationAddress preserves recorded principal-address components, not inferred ISO codes.
type CompanyProfileRegistrationAddress struct {
	_          [0]func()
	Kind       string  `json:"kind"`
	Line1      *string `json:"line1"`
	Line2      *string `json:"line2"`
	City       *string `json:"city"`
	State      *string `json:"state"`
	Postal     *string `json:"postal"`
	CountryRaw *string `json:"country_raw"`
}

// CompanyProfileRegistration holds registry-scoped facts, not an operation or tax-exemption verdict.
type CompanyProfileRegistration struct {
	_            [0]func()
	Authority    string                              `json:"authority"`
	Number       string                              `json:"number"`
	Jurisdiction CompanyProfileJurisdiction          `json:"jurisdiction"`
	Role         string                              `json:"role"`
	LegalForm    CompanyProfileRegistrationLegalForm `json:"legal_form"`
	Status       string                              `json:"status"`
	// This register's reported entity-form date, not universal incorporation or founding.
	FormationDate *string                            `json:"formation_date"`
	Address       *CompanyProfileRegistrationAddress `json:"address"`
}

// CompanyProfileSource provides attribution only for the named selected enrichment fields; observation is not a source update.
type CompanyProfileSource struct {
	_      [0]func()
	Type   string   `json:"type"`
	URL    string   `json:"url"`
	Fields []string `json:"fields"`
	// Artifact observation timestamp.
	ObservedAt string `json:"observed_at"`
	// Explicit source update timestamp, or null. Measurement dates belong to employees.as_of.
	UpdatedAt *string `json:"updated_at,omitempty"`
}

// CompanyProfileDeep holds optional directory detail. Every member may be missing or null; existing releases may omit enrichment fields.
type CompanyProfileDeep struct {
	_            [0]func()
	LegalName    *string                     `json:"legal_name,omitempty"`
	Aliases      []string                    `json:"aliases,omitempty"`
	Jurisdiction *CompanyProfileJurisdiction `json:"jurisdiction,omitempty"`
	// Recorded legal status; not an operating or compliance verdict.
	Status       *string                    `json:"status,omitempty"`
	Websites     []CompanyProfileWebsite    `json:"websites,omitempty"`
	Identifiers  []CompanyProfileIdentifier `json:"identifiers,omitempty"`
	Incorporated *string                    `json:"incorporated,omitempty"`
	Addresses    []CompanyProfileAddress    `json:"addresses,omitempty"`
	Industries   []CompanyProfileIndustry   `json:"industries,omitempty"`
	Parent       *string                    `json:"parent,omitempty"`
	Description  *string                    `json:"description,omitempty"`
	// Reported asset URL; the client does not fetch or license the asset.
	Logo *string `json:"logo,omitempty"`
	// Selected company account URLs; an empty array does not prove no accounts exist.
	Socials       []string                     `json:"socials,omitempty"`
	Founded       *CompanyProfileFounding      `json:"founded,omitempty"`
	Employees     *CompanyProfileEmployees     `json:"employees,omitempty"`
	Registrations []CompanyProfileRegistration `json:"registrations,omitempty"`
	// Attribution for projected enrichment fields only, not the entire legal profile.
	Sources []CompanyProfileSource `json:"sources,omitempty"`
}

// CompanyMatch holds search match evidence. Open strings permit future fields and identifier/listing namespaces.
type CompanyMatch struct {
	_         [0]func()
	Field     *string `json:"field,omitempty"`
	Value     *string `json:"value,omitempty"`
	Type      *string `json:"type,omitempty"`
	Authority *string `json:"authority,omitempty"`
	Exchange  *string `json:"exchange,omitempty"`
}

// CompanyProfile is a directory profile, distinct from national company-number validation.
type CompanyProfile struct {
	_        [0]func()
	ID       string                  `json:"id"`
	Name     string                  `json:"name"`
	Country  *string                 `json:"country,omitempty"`
	Website  *string                 `json:"website,omitempty"`
	Listings []CompanyProfileListing `json:"listings"`
	Address  *CompanyProfileAddress  `json:"address,omitempty"`
	Deep     *CompanyProfileDeep     `json:"deep,omitempty"`
}

// CompanyCandidate is a directory search profile with match evidence; a match is not proof of legal identity.
type CompanyCandidate struct {
	_        [0]func()
	ID       string                  `json:"id"`
	Name     string                  `json:"name"`
	Country  *string                 `json:"country,omitempty"`
	Website  *string                 `json:"website,omitempty"`
	Listings []CompanyProfileListing `json:"listings"`
	Address  *CompanyProfileAddress  `json:"address,omitempty"`
	Deep     *CompanyProfileDeep     `json:"deep,omitempty"`
	Match    CompanyMatch            `json:"match"`
}

// CompanySearch is one page of company candidates. Reuse next with the same selector, filters and limit.
type CompanySearch struct {
	_         [0]func()
	Companies []CompanyCandidate `json:"companies"`
	Next      *string            `json:"next,omitempty"`
}

// CompanyCoverage contains counts for this directory edition, not complete country or worldwide coverage.
type CompanyCoverage struct {
	_            [0]func()
	Scope        string   `json:"scope"`
	Label        string   `json:"label"`
	Description  string   `json:"description"`
	SnapshotAt   string   `json:"snapshot_at"`
	Companies    int64    `json:"companies"`
	Countries    []string `json:"countries"`
	WithWebsite  int64    `json:"with_website"`
	WithListings int64    `json:"with_listings"`
	WithAddress  int64    `json:"with_address"`
}
