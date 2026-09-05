package parseapi

// Response types for the parseAPI public API. Fields are appended as the API grows. Nullable fields are pointers.
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
	Continent  string   `json:"continent"`
	Name       string   `json:"name"`
	Region     string   `json:"region"`
	Subregion  string   `json:"subregion"`
	Population *int64   `json:"population"`
	Area       *float64 `json:"area"`
	Emoji      string   `json:"emoji"`
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
	Country        string   `json:"country"`
	ISO3           string   `json:"iso3"`
	Numeric        int      `json:"numeric"`
	Name           string   `json:"name"`
	FullName       *string  `json:"full_name"`
	LocalName      *string  `json:"local_name"`
	Demonym        *string  `json:"demonym"`
	Capital        *string  `json:"capital"`
	CapitalLat     *float64 `json:"capital_lat"`
	CapitalLon     *float64 `json:"capital_lon"`
	Continent      string   `json:"continent"`
	Region         *string  `json:"region"`
	Subregion      *string  `json:"subregion"`
	Population     *int64   `json:"population"`
	Area           *float64 `json:"area"`
	Currency       *string  `json:"currency"`
	CurrencyName   *string  `json:"currency_name"`
	CurrencySymbol *string  `json:"currency_symbol"`
	TLD            *string  `json:"tld"`
	CallingCode    *string  `json:"calling_code"`
	Emoji          *string  `json:"emoji"`
	Languages      []string `json:"languages"`
	Borders        []string `json:"borders"`
	Blocs          []string `json:"blocs"`
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
	State       string   `json:"state"`
	Name        string   `json:"name"`
	LocalName   *string  `json:"local_name"`
	Type        *string  `json:"type"`
	Country     string   `json:"country"`
	CountryName *string  `json:"country_name"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	Population  *int64   `json:"population"`
	Area        *float64 `json:"area"`
	Timezone    *string  `json:"timezone"`
	Timezones   []string `json:"timezones"`
	ISO3166_2   *string  `json:"iso_3166_2"`
	FIPS        *string  `json:"fips"`
	Capital     *string  `json:"capital"`
	AreaCodes   []string `json:"area_codes"`
	Tax         *string  `json:"tax"`
	TaxRate     *float64 `json:"tax_rate"`
}

type StateDistrictItem struct {
	_        [0]func()
	District string  `json:"district"`
	Name     string  `json:"name"`
	Type     *string `json:"type"`
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
	District    string   `json:"district"`
	Name        string   `json:"name"`
	Type        *string  `json:"type"`
	State       *string  `json:"state"`
	StateName   *string  `json:"state_name"`
	Country     string   `json:"country"`
	CountryName *string  `json:"country_name"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	Population  *int64   `json:"population"`
	// Area is the total in km2 (land + water, or the official total).
	Area *float64 `json:"area"`
	// LandArea / WaterArea are the km2 split, null when the source publishes total only.
	LandArea  *float64 `json:"land_area"`
	WaterArea *float64 `json:"water_area"`
	Seat      *string  `json:"seat"`
	Timezone  *string  `json:"timezone"`
	Timezones []string `json:"timezones"`
}

type City struct {
	_         [0]func()
	Name      string  `json:"name"`
	LocalName *string `json:"local_name"`
	Type      *string `json:"type"`
	// CapitalOf says what this city is the capital of: country, state, or null.
	CapitalOf    *string  `json:"capital_of"`
	State        *string  `json:"state"`
	StateName    *string  `json:"state_name"`
	District     *string  `json:"district"`
	DistrictName *string  `json:"district_name"`
	Country      string   `json:"country"`
	CountryName  *string  `json:"country_name"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	Elevation    *float64 `json:"elevation"`
	ElevationFt  *float64 `json:"elevation_ft"`
	Population   *int64   `json:"population"`
	// Area is the total in km2 (land + water, or the official total).
	Area *float64 `json:"area"`
	// LandArea / WaterArea are the km2 split, null when the source publishes total only.
	LandArea  *float64 `json:"land_area"`
	WaterArea *float64 `json:"water_area"`
	Timezone  *string  `json:"timezone"`
	// ID is the minted parse id (city_ + 12 chars). Stable pin via /city/id/{id}.
	ID string `json:"id"`
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

type Postal struct {
	_                 [0]func()
	Postal            string   `json:"postal"`
	City              *string  `json:"city"`
	CityLocal         *string  `json:"city_local"`
	District          *string  `json:"district"`
	DistrictName      *string  `json:"district_name"`
	DistrictNameLocal *string  `json:"district_name_local"`
	State             *string  `json:"state"`
	StateName         *string  `json:"state_name"`
	StateNameLocal    *string  `json:"state_name_local"`
	Country           string   `json:"country"`
	CountryName       *string  `json:"country_name"`
	Latitude          *float64 `json:"latitude"`
	Longitude         *float64 `json:"longitude"`
	Elevation         *float64 `json:"elevation"`
	ElevationFt       *float64 `json:"elevation_ft"`
	Population        *int64   `json:"population"`
	// Area is the total in km2, null when the source has no water split.
	Area *float64 `json:"area"`
	// LandArea / WaterArea are the km2 split where the source has them.
	LandArea  *float64 `json:"land_area"`
	WaterArea *float64 `json:"water_area"`
	Timezone  *string  `json:"timezone"`
	Currency  *string  `json:"currency"`
	Neighbors []string `json:"neighbors"`
}

type PostalNearbyItem struct {
	_          [0]func()
	Postal     string  `json:"postal"`
	City       *string `json:"city"`
	State      *string `json:"state"`
	Country    string  `json:"country"`
	Distance   float64 `json:"distance"`
	DistanceMi float64 `json:"distance_mi"`
}

type PostalNearby struct {
	_       [0]func()
	Postal  string             `json:"postal"`
	Country string             `json:"country"`
	Radius  float64            `json:"radius"`
	Unit    string             `json:"unit"`
	Nearby  []PostalNearbyItem `json:"nearby"`
}

type PostalDistanceEnd struct {
	_      [0]func()
	Postal string  `json:"postal"`
	City   *string `json:"city"`
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
	Domain      *string    `json:"domain"`
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
	// ConsultedAt is the registry timestamp of this check, ISO.
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
	Checksum  *string `json:"checksum"`
	// Bank is the identifier parsed from the number, not a name.
	Bank     *string `json:"bank"`
	BankName *string `json:"bank_name"`
	Bic      *string `json:"bic"`
	// Branch is the identifier when that country has one.
	Branch  *string `json:"branch"`
	Account *string `json:"account"`
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
	// DeactivatedAt is the ISO date the NPI was deactivated.
	DeactivatedAt *string `json:"deactivated_at"`
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
	Until *string `json:"until"`
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
}

type Tariff struct {
	_ [0]func()
	// HTS is the normalized code with dots (8471.30.01.00).
	HTS string `json:"hts"`
	// Description is the schedule line verbatim.
	Description string `json:"description"`
	// Lineage is the parent descriptions from the schedule outline, outermost first.
	Lineage []string `json:"lineage"`
	// Units is the units of quantity (No., kg).
	Units []string `json:"units"`
	// General is the column 1 general rate, verbatim.
	General *string `json:"general"`
	// Special is the column 1 special rate, verbatim.
	Special *string `json:"special"`
	// Other is the column 2 rate, verbatim.
	Other *string `json:"other"`
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
	Recalls []VINRecall `json:"recalls"`
}

type VIN struct {
	_ [0]func()
	// VIN is the normalized VIN, uppercase, no spaces. Invalid input still echoes the fold.
	VIN    *string `json:"vin"`
	Valid  bool    `json:"valid"`
	Year   *int    `json:"year"`
	Make   *string `json:"make"`
	Model  *string `json:"model"`
	Trim   *string `json:"trim"`
	Series *string `json:"series"`
	// Body is the body style (sedan, coupe, suv, pickup).
	Body *string `json:"body"`
	// Type is the vehicle type (passenger car, truck, motorcycle, bus, trailer).
	Type      *string `json:"type"`
	Doors     *int    `json:"doors"`
	Cylinders *int    `json:"cylinders"`
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
	Gvwr *string  `json:"gvwr"`
	Deep *VINDeep `json:"deep,omitempty"`
}

type Phone struct {
	_       [0]func()
	Phone   *string `json:"phone"`
	Valid   bool    `json:"valid"`
	Country *string `json:"country"`
	// Type is what the numbering plan can see: mobile, landline, toll_free, unknown. Never voip.
	Type *string `json:"type"`
	// State is the NPA-derived state code (US/CA).
	State     *string `json:"state"`
	StateName *string `json:"state_name"`
	// Timezone is the numbering-plan IANA id. Nil when the prefix covers more than one zone.
	Timezone      *string `json:"timezone"`
	National      *string `json:"national"`
	International *string `json:"international"`
	// Deep is always empty. The metered proves are their own endpoints: Carrier, Caller, HLR.
	Deep map[string]any `json:"deep,omitempty"`
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
	Burner *bool `json:"burner"`
	// City is the issuing rate-center city.
	City      *string `json:"city"`
	State     *string `json:"state"`
	StateName *string `json:"state_name"`
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
	// Live reports whether the number is assigned to a subscriber.
	Live *bool `json:"live"`
	// Connected reports whether the handset is reachable right now. Nil means unconfirmed, never no.
	Connected *bool `json:"connected"`
	// The six network extras fill on live HLR dips only. Nil elsewhere (NANP, failover).
	Roaming        *bool   `json:"roaming"`
	RoamingNetwork *string `json:"roaming_network"`
	// RoamingCountry is ISO2, uppercase.
	RoamingCountry *string `json:"roaming_country"`
	// Network is the current serving network name.
	Network         *string `json:"network"`
	OriginalNetwork *string `json:"original_network"`
	MCC             *string `json:"mcc"`
	MNC             *string `json:"mnc"`
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
	A            []string            `json:"a"`
	AAAA         []string            `json:"aaaa"`
	NS           []string            `json:"ns"`
	MX           []MXRecord          `json:"mx"`
	TXT          []string            `json:"txt"`
	Mailhost     *string             `json:"mailhost"`
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
	Currency     string   `json:"currency"`
	Numeric      *int     `json:"numeric"`
	Name         string   `json:"name"`
	NamePlural   *string  `json:"name_plural"`
	Symbol       *string  `json:"symbol"`
	SymbolNative *string  `json:"symbol_native"`
	Digits       *int     `json:"digits"`
	Countries    []string `json:"countries"`
}

// Language is one language by BCP 47 shortest code or ISO 639-3. Codes are lowercase.
type Language struct {
	_         [0]func()
	Language  string   `json:"language"`
	Iso3      *string  `json:"iso3"`
	Name      string   `json:"name"`
	LocalName *string  `json:"local_name"`
	Script    *string  `json:"script"`
	Direction string   `json:"direction"`
	Countries []string `json:"countries"`
}

// Name is a parsed person name. Junk input returns Valid false, never an error.
// Gender comes from dictionary data and is nil when the data does not decide.
type Name struct {
	_          [0]func()
	Name       string  `json:"name"`
	Valid      bool    `json:"valid"`
	Prefix     *string `json:"prefix"`
	First      *string `json:"first"`
	Middle     *string `json:"middle"`
	Last       *string `json:"last"`
	Suffix     *string `json:"suffix"`
	Gender     *string `json:"gender"`
	Salutation *string `json:"salutation"`
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

type Timezone struct {
	_             [0]func()
	Timezone      *string                   `json:"timezone"`
	Name          *string                   `json:"name"`
	Abbreviation  *string                   `json:"abbreviation"`
	Offset        *string                   `json:"offset"`
	OffsetMinutes *int                      `json:"offset_minutes"`
	DST           *bool                     `json:"dst"`
	NextDST       *TimezoneNextDST          `json:"next_dst"`
	Latitude      *float64                  `json:"latitude,omitempty"`
	Longitude     *float64                  `json:"longitude,omitempty"`
	At            *string                   `json:"at,omitempty"`
	To            *TimezoneConversionTarget `json:"to,omitempty"`
}

// DateInfo contains calendar facts for a date. Calendar fields are nil
// when Valid is false. To and Days appear when a comparison was requested.
type DateInfo struct {
	_           [0]func()
	Date        string  `json:"date"`
	Valid       bool    `json:"valid"`
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
	Unix        *int64  `json:"unix"`
	To          *string `json:"to,omitempty"`
	Days        *int    `json:"days,omitempty"`
}

type Holiday struct {
	_          [0]func()
	Date       string   `json:"date"`
	Name       string   `json:"name"`
	LocalName  *string  `json:"local_name"`
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
	_        [0]func()
	City     *CityNearest `json:"city"`
	Timezone *Timezone    `json:"timezone"`
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
	Elevation    *float64   `json:"elevation"`
	ElevationFt  *float64   `json:"elevation_ft"`
	Resolution   *float64   `json:"resolution"`
	Deep         *PointDeep `json:"deep,omitempty"`
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
}

type WeatherCurrent struct {
	_              [0]func()
	Temperature    *float64 `json:"temperature"`
	TemperatureF   *float64 `json:"temperature_f"`
	FeelsLike      *float64 `json:"feels_like"`
	FeelsLikeF     *float64 `json:"feels_like_f"`
	Dewpoint       *float64 `json:"dewpoint"`
	DewpointF      *float64 `json:"dewpoint_f"`
	Humidity       *float64 `json:"humidity"`
	WindSpeed      *float64 `json:"wind_speed"`
	WindSpeedMph   *float64 `json:"wind_speed_mph"`
	WindGust       *float64 `json:"wind_gust"`
	WindGustMph    *float64 `json:"wind_gust_mph"`
	WindDirection  *float64 `json:"wind_direction"`
	Pressure       *float64 `json:"pressure"`
	PressureInhg   *float64 `json:"pressure_inhg"`
	Visibility     *float64 `json:"visibility"`
	VisibilityMi   *float64 `json:"visibility_mi"`
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

type WeatherSource struct {
	_    [0]func()
	ID   string  `json:"id"`
	Name *string `json:"name"`
}

type Weather struct {
	_         [0]func()
	Latitude  float64         `json:"latitude"`
	Longitude float64         `json:"longitude"`
	Current   WeatherCurrent  `json:"current"`
	Station   *WeatherStation `json:"station"`
	Source    WeatherSource   `json:"source"`
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
	Emoji      string      `json:"emoji"`
	Name       string      `json:"name"`
	Shortcodes []string    `json:"shortcodes"`
	Codepoints []string    `json:"codepoints"`
	Hex        string      `json:"hex"`
	Category   *string     `json:"category"`
	Status     *string     `json:"status"`
	Version    *string     `json:"version"`
	Keywords   []string    `json:"keywords"`
	Skins      []EmojiSkin `json:"skins"`
}

type EmojiSearch struct {
	_      [0]func()
	Q      string  `json:"q"`
	Emojis []Emoji `json:"emojis"`
}

type TimezoneConversionTarget struct {
	_             [0]func()
	Timezone      string  `json:"timezone"`
	Name          *string `json:"name"`
	Abbreviation  *string `json:"abbreviation"`
	Offset        string  `json:"offset"`
	OffsetMinutes int     `json:"offset_minutes"`
	DST           bool    `json:"dst"`
	At            string  `json:"at"`
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
}

type CompanyCountry struct {
	_     [0]func()
	Name  *string  `json:"name"`
	Blocs []string `json:"blocs"`
	Tax   *string  `json:"tax"`
}

type CompanyDeep struct {
	_       [0]func()
	Country *CompanyCountry `json:"country"`
	Postal  *Postal         `json:"postal"`
	City    *City           `json:"city"`
}

type Company struct {
	_           [0]func()
	Company     *string      `json:"company"`
	Valid       bool         `json:"valid"`
	Registered  *bool        `json:"registered"`
	Country     *string      `json:"country"`
	Type        *string      `json:"type"`
	Name        *string      `json:"name"`
	Active      *bool        `json:"active"`
	Activity    *string      `json:"activity"`
	Address     *string      `json:"address"`
	City        *string      `json:"city"`
	State       *string      `json:"state"`
	StateName   *string      `json:"state_name"`
	Postal      *string      `json:"postal"`
	CountryName *string      `json:"country_name"`
	VAT         *string      `json:"vat"`
	GST         *bool        `json:"gst"`
	ACN         *string      `json:"acn"`
	Siren       *string      `json:"siren"`
	Siege       *bool        `json:"siege"`
	Kind        *string      `json:"kind"`
	Invoice     *string      `json:"invoice"`
	Deep        *CompanyDeep `json:"deep,omitempty"`
}
