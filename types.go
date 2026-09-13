package parseapi

// Response types for the ParseAPI public API. Nullable fields are pointers.
// Deep objects follow the triad: nil when not requested, empty when
// requested but locked, populated when unlocked.

type IPDeep struct {
	_          [0]func()
	State      *string `json:"state"`
	City       *string `json:"city"`
	Registry   *string `json:"registry"`
	Datacenter *bool   `json:"datacenter"`
	Relay      *bool   `json:"relay"`
	Tor        *bool   `json:"tor"`
	VPN        *bool   `json:"vpn"`
	Provider   *string `json:"provider"`
}

type IP struct {
	_           [0]func()
	IP          string  `json:"ip"`
	Country     *string `json:"country"`
	CountryName *string `json:"country_name"`
	Continent   *string `json:"continent"`
	ASN         *string `json:"asn"`
	ASNName     *string `json:"asn_name"`
	Deep        *IPDeep `json:"deep,omitempty"`
}

type Continent struct {
	_          [0]func()
	Continent  string `json:"continent"`
	Name       string `json:"name"`
	Region     string `json:"region"`
	Subregion  string `json:"subregion"`
	Population *int64 `json:"population"`
	// Reporting year or period for population (YYYY or YYYY-YYYY). Null when unknown or unverifiable.
	PopulationPeriod *string  `json:"population_period"`
	Area             *float64 `json:"area"`
	Emoji            string   `json:"emoji"`
}

type ContinentCountryItem struct {
	_           [0]func()
	Country     string  `json:"country"`
	Name        string  `json:"name"`
	Emoji       *string `json:"emoji"`
	CallingCode *string `json:"calling_code"`
}

type ContinentCountries struct {
	_         [0]func()
	Continent string                 `json:"continent"`
	Countries []ContinentCountryItem `json:"countries"`
}

type Bloc struct {
	_       [0]func()
	Bloc    string `json:"bloc"`
	Name    string `json:"name"`
	Members int    `json:"members"`
}

type BlocCountryItem struct {
	_           [0]func()
	Country     string  `json:"country"`
	Name        string  `json:"name"`
	Emoji       *string `json:"emoji"`
	CallingCode *string `json:"calling_code"`
}

type BlocCountries struct {
	_         [0]func()
	Bloc      string            `json:"bloc"`
	Countries []BlocCountryItem `json:"countries"`
}

type Country struct {
	_              [0]func()
	Country        string       `json:"country"`
	Name           string       `json:"name"`
	NameLocal      *string      `json:"name_local"`
	Continent      string       `json:"continent"`
	Currency       *string      `json:"currency"`
	CurrencyName   *string      `json:"currency_name"`
	CurrencySymbol *string      `json:"currency_symbol"`
	CallingCode    *string      `json:"calling_code"`
	Emoji          *string      `json:"emoji"`
	Languages      []string     `json:"languages"`
	Deep           *CountryDeep `json:"deep,omitempty"`
	Timezones      []string     `json:"timezones"`
}

type CountryStateItem struct {
	_     [0]func()
	State string  `json:"state"`
	Name  string  `json:"name"`
	Type  *string `json:"type"`
}

type CountryStates struct {
	_       [0]func()
	Country string             `json:"country"`
	States  []CountryStateItem `json:"states"`
}

type State struct {
	_           [0]func()
	State       string     `json:"state"`
	Name        string     `json:"name"`
	NameLocal   *string    `json:"name_local"`
	Type        *string    `json:"type"`
	Country     string     `json:"country"`
	CountryName *string    `json:"country_name"`
	Latitude    *float64   `json:"latitude"`
	Longitude   *float64   `json:"longitude"`
	Timezone    *string    `json:"timezone"`
	Timezones   []string   `json:"timezones"`
	ISO3166_2   *string    `json:"iso_3166_2"`
	Deep        *StateDeep `json:"deep,omitempty"`
}

type StateDistrictItem struct {
	_        [0]func()
	District string                 `json:"district"`
	Name     string                 `json:"name"`
	Type     *string                `json:"type"`
	Deep     *StateDistrictItemDeep `json:"deep,omitempty"`
}

type StateDistricts struct {
	_           [0]func()
	State       string              `json:"state"`
	StateName   *string             `json:"state_name"`
	Country     string              `json:"country"`
	CountryName *string             `json:"country_name"`
	Districts   []StateDistrictItem `json:"districts"`
}

type District struct {
	_           [0]func()
	District    string        `json:"district"`
	Name        string        `json:"name"`
	Type        *string       `json:"type"`
	State       *string       `json:"state"`
	StateName   *string       `json:"state_name"`
	Country     string        `json:"country"`
	CountryName *string       `json:"country_name"`
	Latitude    *float64      `json:"latitude"`
	Longitude   *float64      `json:"longitude"`
	Timezone    *string       `json:"timezone"`
	Timezones   []string      `json:"timezones"`
	Deep        *DistrictDeep `json:"deep,omitempty"`
}

type City struct {
	_            [0]func()
	Name         string   `json:"name"`
	NameLocal    *string  `json:"name_local"`
	Type         *string  `json:"type"`
	State        *string  `json:"state"`
	StateName    *string  `json:"state_name"`
	District     *string  `json:"district"`
	DistrictName *string  `json:"district_name"`
	Country      string   `json:"country"`
	CountryName  *string  `json:"country_name"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	Timezone     *string  `json:"timezone"`
	// ID is the minted parse id (city_ + 12 chars). Stable pin via /city/id/{id}.
	ID   string    `json:"id"`
	Deep *CityDeep `json:"deep,omitempty"`
}

// CityNearest is a City plus the distance from the query point.
type CityNearest struct {
	_ [0]func()
	City
	Distance   float64 `json:"distance"`
	DistanceMi float64 `json:"distance_mi"`
}

type CitySearch struct {
	_       [0]func()
	Q       string `json:"q"`
	Country string `json:"country,omitempty"`
	State   string `json:"state,omitempty"`
	Cities  []City `json:"cities"`
}

type CityNearby struct {
	_       [0]func()
	City    string        `json:"city"`
	State   *string       `json:"state"`
	Country string        `json:"country"`
	Radius  float64       `json:"radius"`
	Unit    string        `json:"unit"`
	Nearby  []CityNearest `json:"nearby"`
}

// PostalMetro describes an area's share of ZIP addresses. Category shares are independent fractions.
type PostalMetro struct {
	_                [0]func()
	Code             string   `json:"code"`
	Name             string   `json:"name"`
	Type             string   `json:"type"`
	Share            *float64 `json:"share"`
	ResidentialShare *float64 `json:"residential_share"`
	BusinessShare    *float64 `json:"business_share"`
	OtherShare       *float64 `json:"other_share"`
}

type Postal struct {
	_                 [0]func()
	Postal            string      `json:"postal"`
	City              *string     `json:"city"`
	CityLocal         *string     `json:"city_local"`
	District          *string     `json:"district"`
	DistrictName      *string     `json:"district_name"`
	DistrictNameLocal *string     `json:"district_name_local"`
	State             *string     `json:"state"`
	StateName         *string     `json:"state_name"`
	StateNameLocal    *string     `json:"state_name_local"`
	Country           string      `json:"country"`
	CountryName       *string     `json:"country_name"`
	Latitude          *float64    `json:"latitude"`
	Longitude         *float64    `json:"longitude"`
	Timezone          *string     `json:"timezone"`
	Deep              *PostalDeep `json:"deep,omitempty"`
}

type PostalNearbyItem struct {
	_          [0]func()
	Postal     string           `json:"postal"`
	City       *string          `json:"city"`
	State      *string          `json:"state"`
	Country    string           `json:"country"`
	Distance   float64          `json:"distance"`
	DistanceMi float64          `json:"distance_mi"`
	Deep       *PostalMetroDeep `json:"deep,omitempty"`
}

type PostalNearby struct {
	_       [0]func()
	Postal  string             `json:"postal"`
	Country string             `json:"country"`
	Radius  float64            `json:"radius"`
	Unit    string             `json:"unit"`
	Nearby  []PostalNearbyItem `json:"nearby"`
	Deep    *PostalMetroDeep   `json:"deep,omitempty"`
}

type PostalDistanceEnd struct {
	_      [0]func()
	Postal string           `json:"postal"`
	City   *string          `json:"city"`
	Deep   *PostalMetroDeep `json:"deep,omitempty"`
}

type PostalDistance struct {
	_          [0]func()
	Country    string            `json:"country"`
	From       PostalDistanceEnd `json:"from"`
	To         PostalDistanceEnd `json:"to"`
	Distance   float64           `json:"distance"`
	DistanceMi float64           `json:"distance_mi"`
}

type EmailDeep struct {
	_           [0]func()
	Deliverable *bool `json:"deliverable"`
	Catchall    *bool `json:"catchall"`
}

type Email struct {
	_           [0]func()
	Email       string     `json:"email"`
	Valid       bool       `json:"valid"`
	Free        bool       `json:"free"`
	Domain      *string    `json:"domain"`
	DomainType  *string    `json:"domain_type"`
	DomainValid *bool      `json:"domain_valid"`
	Role        bool       `json:"role"`
	Disposable  bool       `json:"disposable"`
	Deep        *EmailDeep `json:"deep,omitempty"`
	DidYouMean  *string    `json:"didyoumean"`
}

type VATAddress struct {
	_       [0]func()
	Street  *string `json:"street"`
	City    *string `json:"city"`
	Postal  *string `json:"postal"`
	Country *string `json:"country"`
}

type VATDeep struct {
	_            [0]func()
	Registered   *bool       `json:"registered"`
	Name         *string     `json:"name"`
	Address      *VATAddress `json:"address"`
	Consultation *string     `json:"consultation"`
	// ConsultedAt is the registry-provided check time, or nil when unavailable.
	ConsultedAt *string `json:"consulted_at"`
}

type VAT struct {
	_       [0]func()
	VAT     *string  `json:"vat"`
	Valid   bool     `json:"valid"`
	Country *string  `json:"country"`
	From    *string  `json:"from,omitempty"`
	Deep    *VATDeep `json:"deep,omitempty"`
}

type IBAN struct {
	_       [0]func()
	IBAN    *string `json:"iban"`
	Valid   bool    `json:"valid"`
	Country *string `json:"country"`
	// Formatted is the print form in groups of four, for display. Nil when invalid.
	Formatted *string `json:"formatted"`
	// Bank is the identifier parsed from the number, not a name.
	Bank     *string   `json:"bank"`
	BankName *string   `json:"bank_name"`
	Bic      *string   `json:"bic"`
	Deep     *IBANDeep `json:"deep,omitempty"`
}

// NPI is a US healthcare provider record in the healthcare provider registry.
type NPI struct {
	_ [0]func()
	// NPI is the normalized 10-digit NPI. Invalid input still echoes the fold.
	NPI   *string `json:"npi"`
	Valid bool    `json:"valid"`
	// Registered reports whether the NPI exists in the registry.
	Registered *bool `json:"registered"`
	Active     *bool `json:"active"`
	// Excluded reports the OIG exclusion flag.
	Excluded *bool `json:"excluded"`
	// Type is individual or organization.
	Type       *string `json:"type"`
	Name       *string `json:"name"`
	First      *string `json:"first"`
	Last       *string `json:"last"`
	Credential *string `json:"credential"`
	Specialty  *string `json:"specialty"`
	// Taxonomy is the NUCC taxonomy code.
	Taxonomy  *string  `json:"taxonomy"`
	Address   *string  `json:"address"`
	City      *string  `json:"city"`
	State     *string  `json:"state"`
	StateName *string  `json:"state_name"`
	Postal    *string  `json:"postal"`
	Country   *string  `json:"country"`
	Phone     *string  `json:"phone"`
	Deep      *NPIDeep `json:"deep,omitempty"`
}

// NPIEnrollment is one Medicare FFS enrollment row.
type NPIEnrollment struct {
	_ [0]func()
	// Type is part_a, part_b, practitioner, dme, order_refer, or mdpp.
	Type      *string `json:"type"`
	Specialty *string `json:"specialty"`
	State     *string `json:"state"`
}

// NPIDeep is Medicare enrollment on paid plans.
type NPIDeep struct {
	_ [0]func()
	// Medicare is whether the NPI is in the published FFS enrollment extract.
	Medicare *bool `json:"medicare"`
	// OptOut is whether the NPI has a Medicare opt-out affidavit.
	OptOut *bool `json:"opt_out"`
	// Enrollments is type, specialty, and state. Empty when Medicare is false.
	Enrollments []NPIEnrollment `json:"enrollments"`
	// DeactivatedAt is the ISO date the NPI was deactivated.
	DeactivatedAt *string `json:"deactivated_at"`
}

type VINRecall struct {
	_ [0]func()
	// Campaign is the government campaign number.
	Campaign string `json:"campaign"`
	// Date is the report date, ISO YYYY-MM-DD.
	Date      *string `json:"date"`
	Component *string `json:"component"`
	// Summary is the filed summary verbatim.
	Summary *string `json:"summary"`
}

type TariffMeasure struct {
	_ [0]func()
	// Heading is the Chapter 99 heading, dotted (9903.01.24).
	Heading string `json:"heading"`
	// Description is the measure text verbatim.
	Description string `json:"description"`
	// Rate is the rate string verbatim.
	Rate *string `json:"rate"`
	// From is the effective date, ISO YYYY-MM-DD. Nil when the schedule states none.
	From *string `json:"from"`
	// Until is the expiry, ISO YYYY-MM-DD. Nil when open-ended.
	Until       *string `json:"until"`
	Conditional *bool   `json:"conditional"`
}

type TariffDeep struct {
	_ [0]func()
	// Origin is the country the measures were resolved for.
	Origin *string `json:"origin"`
	// EffectiveRate is the composed ad valorem percent. Nil when the
	// components do not compose cleanly.
	EffectiveRate *float64 `json:"effective_rate"`
	// Measures is every Chapter 99 tariff measure that applies to this code
	// from this origin.
	Measures []TariffMeasure `json:"measures"`
	// Units is the units of quantity (No., kg).
	Units []string `json:"units"`
	// Special is the column 1 special rate, verbatim.
	Special *string `json:"special"`
	// Other is the column 2 rate, verbatim.
	Other *string `json:"other"`
}

type Tariff struct {
	_ [0]func()
	// HTS is the normalized code with dots (8471.30.01.00).
	HTS string `json:"hts"`
	// Description is the schedule line verbatim.
	Description string `json:"description"`
	// Lineage is the parent descriptions from the schedule outline, outermost first.
	Lineage []string `json:"lineage"`
	// General is the column 1 general rate, verbatim.
	General *string `json:"general"`
	// Revision is the official release that answered (2026HTSRev17).
	Revision string      `json:"revision"`
	Deep     *TariffDeep `json:"deep,omitempty"`
}

type TariffSearchHit struct {
	_           [0]func()
	HTS         string  `json:"hts"`
	Description string  `json:"description"`
	General     *string `json:"general"`
}

type TariffSearch struct {
	_        [0]func()
	Q        string `json:"q"`
	Revision string `json:"revision"`
	// Lines is up to 20 tariff lines, best match first.
	Lines []TariffSearchHit `json:"lines"`
}

type VINDeep struct {
	_ [0]func()
	// Recalls is the open campaigns for the decoded vehicle. Empty when none,
	// nil when the registry did not answer.
	Recalls   []VINRecall `json:"recalls"`
	Series    *string     `json:"series"`
	Doors     *int        `json:"doors"`
	Cylinders *int        `json:"cylinders"`
	// Displacement is the engine displacement in liters.
	Displacement *float64 `json:"displacement"`
	Fuel         *string  `json:"fuel"`
	Horsepower   *float64 `json:"horsepower"`
	// Drive is fwd, rwd, awd, or 4wd.
	Drive *string `json:"drive"`
	// Transmission is automatic, manual, or cvt.
	Transmission *string `json:"transmission"`
	Manufacturer *string `json:"manufacturer"`
	PlantCity    *string `json:"plant_city"`
	PlantState   *string `json:"plant_state"`
	PlantCountry *string `json:"plant_country"`
	// Gvwr is the gross vehicle weight rating class as filed.
	Gvwr *string `json:"gvwr"`
}

type VIN struct {
	_ [0]func()
	// VIN is the normalized VIN, uppercase, no spaces. Invalid input still echoes the fold.
	VIN   *string `json:"vin"`
	Valid bool    `json:"valid"`
	Year  *int    `json:"year"`
	Make  *string `json:"make"`
	Model *string `json:"model"`
	Trim  *string `json:"trim"`
	// Body is the body style (sedan, coupe, suv, pickup).
	Body *string `json:"body"`
	// Type is the vehicle type (passenger car, truck, motorcycle, bus, trailer).
	Type *string  `json:"type"`
	Deep *VINDeep `json:"deep,omitempty"`
}

type Phone struct {
	_       [0]func()
	Phone   *string `json:"phone"`
	Valid   bool    `json:"valid"`
	Country *string `json:"country"`
	// Type is what the numbering plan can see: mobile, landline, toll_free, unknown. Never voip.
	Type          *string    `json:"type"`
	National      *string    `json:"national"`
	International *string    `json:"international"`
	Deep          *PhoneDeep `json:"deep,omitempty"`
}

type Carrier struct {
	_       [0]func()
	Phone   *string `json:"phone"`
	Valid   bool    `json:"valid"`
	Country *string `json:"country"`
	// Type is the network's word, including voip.
	Type *string `json:"type"`
	// Carrier is the current carrier display name. Nil when the probe had no answer.
	Carrier *string `json:"carrier"`
	// Burner reports whether the carrier is a known burner number app. Nil when carrier is unknown.
	Burner *bool        `json:"burner"`
	Deep   *CarrierDeep `json:"deep,omitempty"`
}

type Caller struct {
	_       [0]func()
	Phone   *string `json:"phone"`
	Valid   bool    `json:"valid"`
	Country *string `json:"country"`
	// Caller is the CNAM record verbatim (all-caps telco artifact). Nil when no record or outside NANP.
	Caller *string `json:"caller"`
}

type HLR struct {
	_       [0]func()
	Phone   *string `json:"phone"`
	Valid   bool    `json:"valid"`
	Country *string `json:"country"`
	// Live reports whether the number was assigned to a subscriber at the last check.
	Live *bool `json:"live"`
	// Connected reports whether the handset was reachable at the last check. Nil means unconfirmed, never no.
	Connected *bool    `json:"connected"`
	Deep      *HLRDeep `json:"deep,omitempty"`
}

type MXRecord struct {
	_        [0]func()
	Priority int    `json:"priority"`
	Host     string `json:"host"`
}

type DomainRegistration struct {
	_          [0]func()
	Registered bool     `json:"registered"`
	Created    *string  `json:"created"`
	Updated    *string  `json:"updated"`
	Expires    *string  `json:"expires"`
	Registrar  *string  `json:"registrar"`
	Status     []string `json:"status"`
	DNSSEC     bool     `json:"dnssec"`
}

type DomainDeep struct {
	_            [0]func()
	Registration *DomainRegistration `json:"registration"`
}

type Domain struct {
	_         [0]func()
	Domain    string      `json:"domain"`
	Available bool        `json:"available"`
	Deep      *DomainDeep `json:"deep,omitempty"`
}

type ASN struct {
	_           [0]func()
	ASN         uint32  `json:"asn"`
	Name        *string `json:"name"`
	Country     *string `json:"country"`
	CountryName *string `json:"country_name"`
}

type MAC struct {
	_         [0]func()
	MAC       string  `json:"mac"`
	Valid     bool    `json:"valid"`
	Vendor    *string `json:"vendor"`
	Local     *bool   `json:"local"`
	Multicast *bool   `json:"multicast"`
}

// BIN contains card-prefix reference data. Nil reference fields mean unknown.
type BIN struct {
	_   [0]func()
	BIN string `json:"bin"`
	// Prefix is the actual longest match and may be shorter than BIN.
	Prefix    *string        `json:"prefix"`
	Country   *string        `json:"country"`
	Issuer    *string        `json:"issuer"`
	Brand     *string        `json:"brand"`
	BrandName *string        `json:"brand_name"`
	Type      *string        `json:"type"`
	Prepaid   *bool          `json:"prepaid"`
	Deep    map[string]any `json:"deep,omitempty"`
}


// DNSRecord preserves DNS presentation text, including TXT quoting.
type DNSRecord struct {
	_     [0]func()
	Name  string `json:"name"`
	Type  string `json:"type"`
	TTL   uint32 `json:"ttl"`
	Value string `json:"value"`
}

type DNS struct {
	_       [0]func()
	Domain  string      `json:"domain"`
	Records []DNSRecord `json:"records"`
}

type MX struct {
	_      [0]func()
	Domain string     `json:"domain"`
	MX     []MXRecord `json:"mx"`
}

type UserAgentDeviceDeep struct {
	_           [0]func()
	Type        *string `json:"type"`
	Brand       *string `json:"brand"`
	Model       *string `json:"model"`
	CPU         *string `json:"cpu"`
	Touchscreen *bool   `json:"touchscreen"`
}

type UserAgentOSDeep struct {
	_        [0]func()
	Name     *string `json:"name"`
	Version  *string `json:"version"`
	Platform *string `json:"platform"`
}

type UserAgentBrowserBrand struct {
	_       [0]func()
	Brand   string `json:"brand"`
	Version string `json:"version"`
}

type UserAgentBrowserDeep struct {
	_       [0]func()
	Name    *string                 `json:"name"`
	Version *string                 `json:"version"`
	Type    *string                 `json:"type"`
	Brands  []UserAgentBrowserBrand `json:"brands,omitempty"`
}

type UserAgentEngineDeep struct {
	_       [0]func()
	Name    *string `json:"name"`
	Version *string `json:"version"`
}

type UserAgentDeep struct {
	_        [0]func()
	Device   *UserAgentDeviceDeep  `json:"device"`
	OS       *UserAgentOSDeep      `json:"os"`
	Browser  *UserAgentBrowserDeep `json:"browser"`
	Engine   *UserAgentEngineDeep  `json:"engine"`
	Headless *bool                 `json:"headless"`
	AI       *bool                 `json:"ai,omitempty"`
	Bot      map[string]any        `json:"bot,omitempty"`
}

type UserAgent struct {
	_         [0]func()
	UserAgent string         `json:"useragent"`
	Device    *string        `json:"device"`
	OS        *string        `json:"os"`
	Browser   *string        `json:"browser"`
	Bot       bool           `json:"bot"`
	Mobile    bool           `json:"mobile"`
	Deep      *UserAgentDeep `json:"deep,omitempty"`
}

type Currency struct {
	_            [0]func()
	Currency     string        `json:"currency"`
	Name         string        `json:"name"`
	Symbol       *string       `json:"symbol"`
	SymbolNative *string       `json:"symbol_native"`
	Digits       *int          `json:"digits"`
	Deep         *CurrencyDeep `json:"deep,omitempty"`
}

// Language is one language by BCP 47 shortest code or ISO 639-3. Codes are lowercase.
type Language struct {
	_         [0]func()
	Language  string        `json:"language"`
	Name      string        `json:"name"`
	NameLocal *string       `json:"name_local"`
	Script    *string       `json:"script"`
	Direction string        `json:"direction"`
	Deep      *LanguageDeep `json:"deep,omitempty"`
}

// Name is a parsed person name. Junk input returns Valid false, never an error.
// Gender comes from dictionary data and is nil when the data does not decide.
type Name struct {
	_      [0]func()
	Name   string    `json:"name"`
	Valid  bool      `json:"valid"`
	Prefix *string   `json:"prefix"`
	First  *string   `json:"first"`
	Middle *string   `json:"middle"`
	Last   *string   `json:"last"`
	Suffix *string   `json:"suffix"`
	Deep   *NameDeep `json:"deep,omitempty"`
}

type CurrencyRate struct {
	_         [0]func()
	Base      string   `json:"base"`
	Quote     string   `json:"quote"`
	Rate      float64  `json:"rate"`
	Date      string   `json:"date"`
	Amount    *float64 `json:"amount,omitempty"`
	Converted *float64 `json:"converted,omitempty"`
	Source    string   `json:"source,omitempty"`
}

type TimezoneNextDST struct {
	_            [0]func()
	At           string `json:"at"`
	DST          bool   `json:"dst"`
	Offset       string `json:"offset"`
	Abbreviation string `json:"abbreviation"`
}

// Time contains the local clock and timezone facts. Nil clock fields mean no resolved zone.
type Time = Timezone

type Timezone struct {
	_            [0]func()
	Timezone     *string                   `json:"timezone"`
	Abbreviation *string                   `json:"abbreviation"`
	Offset       *string                   `json:"offset"`
	DST          *bool                     `json:"dst"`
	Latitude     *float64                  `json:"latitude,omitempty"`
	Longitude    *float64                  `json:"longitude,omitempty"`
	At           *string                   `json:"at,omitempty"`
	Unix         *int64                    `json:"unix,omitempty"`
	To           *TimezoneConversionTarget `json:"to,omitempty"`
	Deep         *TimezoneDeep             `json:"deep,omitempty"`
}

// DateInfo contains calendar facts for a date. Calendar fields are nil
// when Valid is false. To and Days appear when a comparison was requested.
type DateInfo struct {
	_     [0]func()
	Date  string        `json:"date"`
	Valid bool          `json:"valid"`
	Unix  *int64        `json:"unix"`
	To    *string       `json:"to,omitempty"`
	Days  *int          `json:"days,omitempty"`
	Deep  *DateInfoDeep `json:"deep,omitempty"`
}

type Holiday struct {
	_          [0]func()
	Date       string   `json:"date"`
	Name       string   `json:"name"`
	NameLocal  *string  `json:"name_local"`
	Type       string   `json:"type"`
	Regions    []string `json:"regions"`
	Substitute bool     `json:"substitute"`
}

type HolidayYear struct {
	_        [0]func()
	Country  string    `json:"country"`
	Year     int       `json:"year"`
	Holidays []Holiday `json:"holidays"`
}

type HolidayDate struct {
	_       [0]func()
	Country string   `json:"country"`
	Date    string   `json:"date"`
	Holiday *Holiday `json:"holiday"`
}

type Elevation struct {
	_           [0]func()
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	Elevation   *float64 `json:"elevation"`
	ElevationFt *float64 `json:"elevation_ft"`
	Resolution  *float64 `json:"resolution"`
}

type PointDeep struct {
	_           [0]func()
	City        *PointCity `json:"city"`
	ElevationFt *float64   `json:"elevation_ft"`
	Resolution  *float64   `json:"resolution"`
}

type Point struct {
	_            [0]func()
	Latitude     float64    `json:"latitude"`
	Longitude    float64    `json:"longitude"`
	Country      *string    `json:"country"`
	CountryName  *string    `json:"country_name"`
	State        *string    `json:"state"`
	StateName    *string    `json:"state_name"`
	District     *string    `json:"district"`
	DistrictName *string    `json:"district_name"`
	Deep         *PointDeep `json:"deep,omitempty"`
	Timezone     *string    `json:"timezone"`
}

type WeatherForecastPeriod struct {
	_                   [0]func()
	Name                string   `json:"name"`
	Start               *string  `json:"start"`
	End                 *string  `json:"end"`
	Daytime             *bool    `json:"daytime"`
	Temperature         *float64 `json:"temperature"`
	TemperatureF        *float64 `json:"temperature_f"`
	PrecipitationChance *float64 `json:"precipitation_chance"`
	WindSpeed           *float64 `json:"wind_speed"`
	WindSpeedMph        *float64 `json:"wind_speed_mph"`
	WindDirection       *float64 `json:"wind_direction"`
	Condition           *string  `json:"condition"`
	ConditionName       *string  `json:"condition_name"`
	ConditionEmoji      *string  `json:"condition_emoji"`
}

type WeatherAlert struct {
	_        [0]func()
	Event    string  `json:"event"`
	Severity *string `json:"severity"`
	Urgency  *string `json:"urgency"`
	Headline *string `json:"headline"`
	Onset    *string `json:"onset"`
	Expires  *string `json:"expires"`
}

type WeatherHour struct {
	_                   [0]func()
	At                  *string  `json:"at"`
	Daytime             *bool    `json:"daytime"`
	Temperature         *float64 `json:"temperature"`
	TemperatureF        *float64 `json:"temperature_f"`
	Humidity            *float64 `json:"humidity"`
	PrecipitationChance *float64 `json:"precipitation_chance"`
	WindSpeed           *float64 `json:"wind_speed"`
	WindSpeedMph        *float64 `json:"wind_speed_mph"`
	WindDirection       *float64 `json:"wind_direction"`
	Condition           *string  `json:"condition"`
	ConditionName       *string  `json:"condition_name"`
	ConditionEmoji      *string  `json:"condition_emoji"`
	FeelsLike           *float64 `json:"feels_like"`
	FeelsLikeF          *float64 `json:"feels_like_f"`
	WindGust            *float64 `json:"wind_gust"`
	WindGustMph         *float64 `json:"wind_gust_mph"`
}

type WeatherMinute struct {
	_               [0]func()
	At              string   `json:"at"`
	Precipitation   *float64 `json:"precipitation"`
	PrecipitationIn *float64 `json:"precipitation_in"`
	Type            *string  `json:"type"`
}

type WeatherDay struct {
	_                   [0]func()
	Date                string   `json:"date"`
	High                *float64 `json:"high"`
	HighF               *float64 `json:"high_f"`
	Low                 *float64 `json:"low"`
	LowF                *float64 `json:"low_f"`
	PrecipitationChance *float64 `json:"precipitation_chance"`
	Condition           *string  `json:"condition"`
	ConditionName       *string  `json:"condition_name"`
	ConditionEmoji      *string  `json:"condition_emoji"`
	Sunrise             *string  `json:"sunrise"`
	Sunset              *string  `json:"sunset"`
	MoonPhase           *string  `json:"moon_phase"`
	MoonPhaseName       *string  `json:"moon_phase_name"`
	MoonPhaseEmoji      *string  `json:"moon_phase_emoji"`
}

type WeatherDeep struct {
	_        [0]func()
	Forecast []WeatherForecastPeriod `json:"forecast"`
	Alerts   []WeatherAlert          `json:"alerts"`
	Minutes  []WeatherMinute         `json:"minutes"`
	Hours    []WeatherHour           `json:"hours"`
	Days     []WeatherDay            `json:"days"`
	Air      *WeatherAir             `json:"air"`
	History  *WeatherHistory         `json:"history,omitempty"`
	Current  *WeatherCurrentDeep     `json:"current"`
}

type WeatherCurrent struct {
	_              [0]func()
	Temperature    *float64 `json:"temperature"`
	TemperatureF   *float64 `json:"temperature_f"`
	FeelsLike      *float64 `json:"feels_like"`
	FeelsLikeF     *float64 `json:"feels_like_f"`
	Humidity       *float64 `json:"humidity"`
	WindSpeed      *float64 `json:"wind_speed"`
	WindSpeedMph   *float64 `json:"wind_speed_mph"`
	WindDirection  *float64 `json:"wind_direction"`
	Condition      *string  `json:"condition"`
	ConditionName  *string  `json:"condition_name"`
	ConditionEmoji *string  `json:"condition_emoji"`
	ObservedAt     *string  `json:"observed_at"`
}

type WeatherStation struct {
	_          [0]func()
	ID         string   `json:"id"`
	Name       *string  `json:"name"`
	Distance   *float64 `json:"distance"`
	DistanceMi *float64 `json:"distance_mi"`
}

type Weather struct {
	_         [0]func()
	Latitude  float64         `json:"latitude"`
	Longitude float64         `json:"longitude"`
	Current   WeatherCurrent  `json:"current"`
	Station   *WeatherStation `json:"station"`
	Deep      *WeatherDeep    `json:"deep,omitempty"`
}

type EmojiSkin struct {
	_       [0]func()
	Emoji   string  `json:"emoji"`
	Tone    string  `json:"tone"`
	Unicode *string `json:"unicode"`
	Hex     *string `json:"hex"`
}

type Emoji struct {
	_          [0]func()
	Emoji      string     `json:"emoji"`
	Name       string     `json:"name"`
	Shortcodes []string   `json:"shortcodes"`
	Category   *string    `json:"category"`
	Deep       *EmojiDeep `json:"deep,omitempty"`
}

type EmojiSearch struct {
	_      [0]func()
	Q      string  `json:"q"`
	Emojis []Emoji `json:"emojis"`
}

type TimezoneConversionTarget struct {
	_            [0]func()
	Timezone     string                        `json:"timezone"`
	Abbreviation *string                       `json:"abbreviation"`
	Offset       string                        `json:"offset"`
	DST          bool                          `json:"dst"`
	At           string                        `json:"at"`
	Unix         *int64                        `json:"unix,omitempty"`
	Deep         *TimezoneConversionTargetDeep `json:"deep,omitempty"`
}

type WeatherAir struct {
	_       [0]func()
	AQI     *float64 `json:"aqi"`
	AQIName *string  `json:"aqi_name"`
	PM25    *float64 `json:"pm2_5"`
	Pm10    *float64 `json:"pm10"`
}

type WeatherHistory struct {
	_               [0]func()
	Date            string   `json:"date"`
	High            *float64 `json:"high"`
	HighF           *float64 `json:"high_f"`
	Low             *float64 `json:"low"`
	LowF            *float64 `json:"low_f"`
	Precipitation   *float64 `json:"precipitation"`
	PrecipitationIn *float64 `json:"precipitation_in"`
	WindMax         *float64 `json:"wind_max"`
	WindMaxMph      *float64 `json:"wind_max_mph"`
	Sunrise         *string  `json:"sunrise"`
	Sunset          *string  `json:"sunset"`
	MoonPhase       *string  `json:"moon_phase"`
	MoonPhaseName   *string  `json:"moon_phase_name"`
	MoonPhaseEmoji  *string  `json:"moon_phase_emoji"`
}

type Address struct {
	_            [0]func()
	Address      *string        `json:"address"`
	Valid        bool           `json:"valid"`
	Registered   *bool          `json:"registered"`
	Number       *string        `json:"number"`
	Street       *string        `json:"street"`
	Unit         *string        `json:"unit"`
	City         *string        `json:"city"`
	District     *string        `json:"district"`
	DistrictName *string        `json:"district_name"`
	State        *string        `json:"state"`
	StateName    *string        `json:"state_name"`
	Postal       *string        `json:"postal"`
	Country      *string        `json:"country"`
	CountryName  *string        `json:"country_name"`
	Latitude     *float64       `json:"latitude"`
	Longitude    *float64       `json:"longitude"`
	Deep         map[string]any `json:"deep,omitempty"`
}

type AddressSuggestion struct {
	_         [0]func()
	Address   string   `json:"address"`
	Number    *string  `json:"number"`
	Street    *string  `json:"street"`
	Unit      *string  `json:"unit"`
	City      *string  `json:"city"`
	State     *string  `json:"state"`
	Postal    *string  `json:"postal"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type AddressSearch struct {
	_         [0]func()
	Q         string              `json:"q"`
	Postal    *string             `json:"postal,omitempty"`
	City      *string             `json:"city,omitempty"`
	State     *string             `json:"state,omitempty"`
	Country   *string             `json:"country,omitempty"`
	Addresses []AddressSuggestion `json:"addresses"`
	// Why suggestions are empty: more_input, missing_context or no_matches. Null with suggestions. Open to future values. Operational failures are errors.
	Reason *string `json:"reason"`
}

type CompanyCountry struct {
	_     [0]func()
	Name  *string  `json:"name"`
	Blocs []string `json:"blocs"`
	// Levy name, such as VAT, GST or sales tax. Null when unknown or not applicable.
	Tax *string `json:"tax"`
}

type CompanyDeep struct {
	_           [0]func()
	Activity    *string `json:"activity"`
	StateName   *string `json:"state_name"`
	CountryName *string `json:"country_name"`
	VAT         *string `json:"vat"`
	GST         *bool   `json:"gst"`
	ACN         *string `json:"acn"`
	Siren       *string `json:"siren"`
	Siege       *bool   `json:"siege"`
	Kind        *string `json:"kind"`
	Invoice     *string `json:"invoice"`
}

type Company struct {
	_          [0]func()
	Company    *string      `json:"company"`
	Valid      bool         `json:"valid"`
	Registered *bool        `json:"registered"`
	Country    *string      `json:"country"`
	Type       *string      `json:"type"`
	Name       *string      `json:"name"`
	Active     *bool        `json:"active"`
	Address    *string      `json:"address"`
	City       *string      `json:"city"`
	State      *string      `json:"state"`
	Postal     *string      `json:"postal"`
	Deep       *CompanyDeep `json:"deep,omitempty"`
}

// MeasureChoice is one explicit interpretation of an ambiguous unit.
type MeasureChoice struct {
	Unit string `json:"unit"`
	Name string `json:"name"`
}

// Measure contains the parsed result. Amount stays a decimal string, including zero.
type Measure struct {
	Measure string          `json:"measure"`
	Valid   bool            `json:"valid"`
	Type    *string         `json:"type"`
	Amount  *string         `json:"amount"`
	Unit    *string         `json:"unit"`
	Reason  *string         `json:"reason"`
	Choices []MeasureChoice `json:"choices"`
}

type MeasureUnit struct {
	Unit    string   `json:"unit"`
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Aliases []string `json:"aliases"`
}

type MeasureUnits struct {
	Units []MeasureUnit `json:"units"`
}

// NAICSChild names a direct child industry code.
type NAICSChild struct {
	NAICS string `json:"naics"`
	Name  string `json:"name"`
}

// NAICSExclusion is a classification exclusion. Generic exclusions can have no linked codes.
type NAICSExclusion struct {
	Description string       `json:"description"`
	Codes       []NAICSChild `json:"codes"`
}

// NAICSCorrection is a query token corrected only during typo fallback.
type NAICSCorrection struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// NAICSMatch identifies the actual title, activity term or code matching a search.
type NAICSMatch struct {
	// Field is currently name, term or naics. It remains an open string.
	Field string `json:"field"`
	Text  string `json:"text"`
	// Corrections is empty for exact, plural and prefix matches.
	Corrections []NAICSCorrection `json:"corrections"`
}

// NAICS is a US NAICS 2022 definition and its hierarchy.
type NAICS struct {
	_          [0]func()
	NAICS      string  `json:"naics"`
	Name       string  `json:"name"`
	Level      int     `json:"level"`
	Parent     *string `json:"parent"`
	ParentName *string `json:"parent_name"`
	// Match is search evidence, absent on direct lookup and older responses.
	Match   *NAICSMatch `json:"match"`
	Year    int         `json:"year"`
	Country string      `json:"country"`
	Deep    *NAICSDeep  `json:"deep,omitempty"`
}

type NAICSSearch struct {
	_       [0]func()
	Q       string              `json:"q"`
	Year    int                 `json:"year"`
	Country string              `json:"country"`
	Results []NAICSSearchResult `json:"results"`
}

type CountryDeep struct {
	_          [0]func()
	ISO3       *string  `json:"iso3"`
	Numeric    *int     `json:"numeric"`
	FullName   *string  `json:"full_name"`
	Demonym    *string  `json:"demonym"`
	Capital    *string  `json:"capital"`
	CapitalLat *float64 `json:"capital_lat"`
	CapitalLon *float64 `json:"capital_lon"`
	Region     *string  `json:"region"`
	Subregion  *string  `json:"subregion"`
	Population *int64   `json:"population"`
	// Reporting year or period for population (YYYY or YYYY-YYYY). Null when unknown or unverifiable.
	PopulationPeriod *string  `json:"population_period"`
	Area             *float64 `json:"area"`
	TLD              *string  `json:"tld"`
	Borders          []string `json:"borders"`
	Blocs            []string `json:"blocs"`
	// Levy name, such as VAT, GST or sales tax. Null when unknown or not applicable.
	Tax *string `json:"tax"`
	// Standard country reference rate in percent (19 means 19%). Null is unknown, zero is known zero.
	TaxRate *float64 `json:"tax_rate"`
	// Tax registration number mask (9 is a digit, A is a letter). Describes format only.
	TaxIDFormat *string `json:"tax_id_format"`
	// Anchored tax registration number format regex. A match does not establish registration.
	TaxIDRegex   *string           `json:"tax_id_regex"`
	WeekStart    *string           `json:"week_start"`
	Units        *string           `json:"units"`
	DrivingSide  *string           `json:"driving_side"`
	Plugs        []string          `json:"plugs"`
	Voltage      *int              `json:"voltage"`
	Frequency    *int              `json:"frequency"`
	Emergency    *CountryEmergency `json:"emergency"`
	PostalFormat *string           `json:"postal_format"`
	PostalRegex  *string           `json:"postal_regex"`
	IOC          *string           `json:"ioc"`
	FIFA         *string           `json:"fifa"`
	Plate        *string           `json:"plate"`
}

type StateDeep struct {
	_          [0]func()
	Population *int64 `json:"population"`
	// Reporting year or period for population (YYYY or YYYY-YYYY). Null when unknown or unverifiable.
	PopulationPeriod *string  `json:"population_period"`
	Area             *float64 `json:"area"`
	FIPS             *string  `json:"fips"`
	Capital          *string  `json:"capital"`
	AreaCodes        []string `json:"area_codes"`
	// Levy name, such as VAT, GST or sales tax. Null when unknown or not applicable.
	Tax *string `json:"tax"`
	// State or province reference rate in percent. Country, state and postal rates are alternative references, not additive.
	TaxRate *float64 `json:"tax_rate"`
}

type DistrictDeep struct {
	_          [0]func()
	Population *int64 `json:"population"`
	// Reporting year or period for population (YYYY or YYYY-YYYY). Null when unknown or unverifiable.
	PopulationPeriod *string `json:"population_period"`
	// Area is the total in km2 (land + water, or the official total).
	Area *float64 `json:"area"`
	// LandArea / WaterArea are the km2 split, null when the source publishes total only.
	LandArea  *float64 `json:"land_area"`
	WaterArea *float64 `json:"water_area"`
	Seat      *string  `json:"seat"`
	// Median annual property tax payable on owner-occupied homes in this statistical area. Null when unsupported, missing or censored.
	PropertyTax *PropertyTax `json:"property_tax"`
}

type CityDeep struct {
	_ [0]func()
	// CapitalOf says what this city is the capital of: country, state, or null.
	CapitalOf   *string  `json:"capital_of"`
	Elevation   *float64 `json:"elevation"`
	ElevationFt *float64 `json:"elevation_ft"`
	Population  *int64   `json:"population"`
	// Reporting year or period for population (YYYY or YYYY-YYYY). Null when unknown or unverifiable.
	PopulationPeriod *string `json:"population_period"`
	// Area is the total in km2 (land + water, or the official total).
	Area *float64 `json:"area"`
	// LandArea / WaterArea are the km2 split, null when the source publishes total only.
	LandArea  *float64 `json:"land_area"`
	WaterArea *float64 `json:"water_area"`
}

type PostalDeep struct {
	_           [0]func()
	Elevation   *float64 `json:"elevation"`
	ElevationFt *float64 `json:"elevation_ft"`
	Population  *int64   `json:"population"`
	// Reporting year or period for population (YYYY or YYYY-YYYY). Null when unknown or unverifiable.
	PopulationPeriod *string `json:"population_period"`
	// Area is the total in km2, null when the source has no water split.
	Area *float64 `json:"area"`
	// LandArea / WaterArea are the km2 split where the source has them.
	LandArea  *float64 `json:"land_area"`
	WaterArea *float64 `json:"water_area"`
	Currency  *string  `json:"currency"`
	Neighbors []string `json:"neighbors"`
	// Metros is nil when unknown, and non-nil empty when observed outside all covered areas.
	Metros []PostalMetro `json:"metros"`
	// Levy name, such as VAT, GST or sales tax. Null when unknown or not applicable.
	Tax *string `json:"tax"`
	// Combined US ZIP reference rate in percent (7.9 means 7.9%). An exact address can differ. Null is unknown, zero is known zero.
	TaxRate *float64 `json:"tax_rate"`
	// State component of the ZIP reference rate, in percent. Null when unknown.
	TaxRateState *float64 `json:"tax_rate_state"`
	// County component of the ZIP reference rate, in percent. Null when unknown.
	TaxRateCounty *float64 `json:"tax_rate_county"`
	// City component of the ZIP reference rate, in percent. Null when unknown.
	TaxRateCity *float64 `json:"tax_rate_city"`
	// Special component of the ZIP reference rate, in percent. Null when unknown.
	TaxRateSpecial *float64 `json:"tax_rate_special"`
	// Median annual property tax payable on owner-occupied homes in this statistical area. Null when unsupported, missing or censored.
	PropertyTax *PropertyTax `json:"property_tax"`
}

type IBANDeep struct {
	_        [0]func()
	Checksum *string `json:"checksum"`
	// Branch is the identifier when that country has one.
	Branch  *string `json:"branch"`
	Account *string `json:"account"`
}

type PhoneDeep struct {
	_ [0]func()
	// State is the NPA-derived state code (US/CA).
	State     *string `json:"state"`
	StateName *string `json:"state_name"`
	// Timezone is the numbering-plan IANA id. Nil when the prefix covers more than one zone.
	Timezone *string `json:"timezone"`
}

type CarrierDeep struct {
	_ [0]func()
	// City is the issuing rate-center city.
	City      *string `json:"city"`
	State     *string `json:"state"`
	StateName *string `json:"state_name"`
}

type HLRDeep struct {
	_ [0]func()
	// Network diagnostics available from the last check. Nil when unconfirmed.
	Roaming        *bool   `json:"roaming"`
	RoamingNetwork *string `json:"roaming_network"`
	// RoamingCountry is ISO2, uppercase.
	RoamingCountry *string `json:"roaming_country"`
	// Network is the serving network name at the last check.
	Network         *string `json:"network"`
	OriginalNetwork *string `json:"original_network"`
	MCC             *string `json:"mcc"`
	MNC             *string `json:"mnc"`
}

type NAICSDeep struct {
	_           [0]func()
	Description *string      `json:"description"`
	Children    []NAICSChild `json:"children"`
	// Exclusions is nil for omitted/null older responses and non-nil when supplied.
	Exclusions []NAICSExclusion `json:"exclusions"`
}

type CurrencyDeep struct {
	_          [0]func()
	Numeric    *int     `json:"numeric"`
	NamePlural *string  `json:"name_plural"`
	Countries  []string `json:"countries"`
}

type LanguageDeep struct {
	_         [0]func()
	Iso3      *string  `json:"iso3"`
	Countries []string `json:"countries"`
}

type NameDeep struct {
	_ [0]func()
	// Known is name membership, independent of gender.
	Known      *bool   `json:"known"`
	Gender     *string `json:"gender"`
	Salutation *string `json:"salutation"`
}

type TimezoneDeep struct {
	_             [0]func()
	Name          *string          `json:"name"`
	OffsetSeconds *int             `json:"offset_seconds,omitempty"`
	OffsetMinutes *int             `json:"offset_minutes"`
	NextDST       *TimezoneNextDST `json:"next_dst"`
}

type TimezoneConversionTargetDeep struct {
	_             [0]func()
	Name          *string `json:"name"`
	OffsetSeconds *int    `json:"offset_seconds,omitempty"`
	OffsetMinutes *int    `json:"offset_minutes"`
}

type DateInfoDeep struct {
	_           [0]func()
	Year        *int    `json:"year"`
	Month       *int    `json:"month"`
	MonthName   *string `json:"month_name"`
	Day         *int    `json:"day"`
	Weekday     *int    `json:"weekday"`
	WeekdayName *string `json:"weekday_name"`
	Week        *int    `json:"week"`
	WeekYear    *int    `json:"week_year"`
	DayOfYear   *int    `json:"day_of_year"`
	Quarter     *int    `json:"quarter"`
	Leap        *bool   `json:"leap"`
	DaysInMonth *int    `json:"days_in_month"`
}

type EmojiDeep struct {
	_          [0]func()
	Codepoints []string    `json:"codepoints"`
	Hex        *string     `json:"hex"`
	Status     *string     `json:"status"`
	Version    *string     `json:"version"`
	Keywords   []string    `json:"keywords"`
	Skins      []EmojiSkin `json:"skins"`
}

type WeatherCurrentDeep struct {
	_            [0]func()
	Dewpoint     *float64 `json:"dewpoint"`
	DewpointF    *float64 `json:"dewpoint_f"`
	WindGust     *float64 `json:"wind_gust"`
	WindGustMph  *float64 `json:"wind_gust_mph"`
	Pressure     *float64 `json:"pressure"`
	PressureInhg *float64 `json:"pressure_inhg"`
	Visibility   *float64 `json:"visibility"`
	VisibilityMi *float64 `json:"visibility_mi"`
}

type PointCity struct {
	_           [0]func()
	Name        *string  `json:"name"`
	NameLocal   *string  `json:"name_local"`
	Type        *string  `json:"type"`
	State       *string  `json:"state"`
	StateName   *string  `json:"state_name"`
	Country     *string  `json:"country"`
	CountryName *string  `json:"country_name"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	// ID is the minted parse id (city_ + 12 chars). Stable pin via /city/id/{id}.
	ID         *string  `json:"id"`
	Distance   *float64 `json:"distance"`
	DistanceMi *float64 `json:"distance_mi"`
}

type PostalMetroDeep struct {
	_      [0]func()
	Metros []PostalMetro `json:"metros"`
}

type StateDistrictItemDeep struct {
	_          [0]func()
	Population *int64 `json:"population"`
	// Reporting year or period for population (YYYY or YYYY-YYYY). Null when unknown or unverifiable.
	PopulationPeriod *string `json:"population_period"`
}

type CountryEmergency struct {
	_         [0]func()
	Police    *string `json:"police"`
	Ambulance *string `json:"ambulance"`
	Fire      *string `json:"fire"`
}

type NAICSSearchResult struct {
	_          [0]func()
	NAICS      string  `json:"naics"`
	Name       string  `json:"name"`
	Level      int     `json:"level"`
	Parent     *string `json:"parent"`
	ParentName *string `json:"parent_name"`
	// Match is search evidence, absent on direct lookup and older responses.
	Match *NAICSMatch `json:"match"`
	Deep  *NAICSDeep  `json:"deep,omitempty"`
}

// PropertyTax is an area estimate, not a specific property bill.
type PropertyTax struct {
	_ [0]func()
	// Median annual tax payable, in currency units adjusted to the final year of period. Not a tax rate or an individual property bill.
	AnnualMedian float64 `json:"annual_median"`
	// ISO 4217 currency code, currently USD.
	Currency string `json:"currency"`
	// Reporting period, YYYY-YYYY. Monetary amounts use the final year of this period.
	Period string `json:"period"`
}
