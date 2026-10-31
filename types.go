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
	// Present only when supported. District is qualified by the response country.
	District       *string `json:"district,omitempty"`
	DistrictName   *string `json:"district_name,omitempty"`
	CensusDivision *string `json:"census_division,omitempty"`
	CensusRegion   *string `json:"census_region,omitempty"`
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
	_       [0]func()
	Country string  `json:"country"`
	Name    string  `json:"name"`
	Flag    *string `json:"flag"`
	// Deprecated: use Flag.
	Emoji       *string `json:"-"`
	CallingCode *string `json:"calling_code"`
}

type ContinentCountries struct {
	_         [0]func()
	Continent string                 `json:"continent"`
	Countries []ContinentCountryItem `json:"countries"`
}







type Country struct {
	_              [0]func()
	Country        string  `json:"country"`
	Name           string  `json:"name"`
	NameLocal      *string `json:"name_local"`
	Continent      string  `json:"continent"`
	Currency       *string `json:"currency"`
	CurrencyName   *string `json:"currency_name"`
	CurrencySymbol *string `json:"currency_symbol"`
	CallingCode    *string `json:"calling_code"`
	Flag           *string `json:"flag"`
	// Deprecated: use Flag. Historical source adapter, not a second JSON key.
	Emoji     *string      `json:"-"`
	Languages []string     `json:"languages"`
	Deep      *CountryDeep `json:"deep,omitempty"`
	Timezones []string     `json:"timezones"`
	ISO3      *string      `json:"iso_alpha3,omitempty"`
	Numeric   *int         `json:"iso_numeric,omitempty"`

	NameFull *string `json:"name_full,omitempty"`
	// Deprecated: use NameFull. Retained as a source-compatible adapter.
	FullName        *string         `json:"-"`
	Region          *string         `json:"region,omitempty"`
	Subregion       *string         `json:"subregion,omitempty"`
	Capital         *CountryCapital `json:"capital,omitempty"`
	CapitalLat      *float64        `json:"-"`
	CapitalLon      *float64        `json:"-"`
	Borders         []string        `json:"borders"`
	Blocs           []string        `json:"blocs"`
	Demonym         *string         `json:"demonym,omitempty"`
	nameFullPresent bool
}

type CountryCapital struct {
	_         [0]func()
	Name      *string  `json:"name"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type CountryTax struct {
	_    [0]func()
	Name *string  `json:"name"`
	Rate *float64 `json:"rate"`
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

type StateCapital struct {
	_         [0]func()
	Name      *string  `json:"name"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type State struct {
	_           [0]func()
	State       string        `json:"state"`
	Name        string        `json:"name"`
	NameLocal   *string       `json:"name_local"`
	Type        *string       `json:"type"`
	Country     string        `json:"country"`
	CountryName *string       `json:"country_name"`
	Latitude    *float64      `json:"latitude"`
	Longitude   *float64      `json:"longitude"`
	Timezones   []string      `json:"timezones"`
	ISO         *string       `json:"iso_subdivision"`
	FIPS        *string       `json:"fips"`
	AreaCodes   []string      `json:"area_codes"`
	Capital     *StateCapital `json:"capital"`
	Deep        *StateDeep    `json:"deep,omitempty"`
	// Deprecated: retained native member for the released ISO subdivision field.
	ISO31662 *string `json:"-"`
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
	Timezones   []string      `json:"timezones"`
	Deep        *DistrictDeep `json:"deep,omitempty"`
}

type City struct {
	Match        *CityMatch `json:"match"`
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
	// Grain and Parent are present when hierarchy is requested.
	Grain                *string     `json:"grain,omitempty"`
	Parent               *CityParent `json:"parent,omitempty"`
	NameRomanized        *string     `json:"name_romanized"`
	RomanizationLanguage *string     `json:"romanization_language"`
}

// CityNearest is a City plus the distance from the query point.
type CityNearest struct {
	_ [0]func()
	City
	Distance   float64 `json:"distance"`
	DistanceMi float64 `json:"distance_mi"`
}

type CitySearch struct {
	_          [0]func()
	Q          string  `json:"q"`
	Country    string  `json:"country,omitempty"`
	State      string  `json:"state,omitempty"`
	Cities     []City  `json:"cities"`
	NextCursor *string `json:"next_cursor,omitempty"`
}

type CityNearby struct {
	_          [0]func()
	City       string        `json:"city"`
	State      *string       `json:"state"`
	Country    string        `json:"country"`
	Radius     float64       `json:"radius"`
	Unit       string        `json:"unit"`
	Nearby     []CityNearest `json:"nearby"`
	NextCursor *string       `json:"next_cursor,omitempty"`
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
	// Percentage of ZIP addresses (0-100). Nil is unknown or absent.
	SharePercent *float64 `json:"share_percent"`
}

// PostalLocality is a supported Australian postcode suburb choice, independent of the scalar city.
type PostalLocality struct {
	_         [0]func()
	City      string `json:"city"`
	State     string `json:"state"`
	StateName string `json:"state_name"`
}

type Postal struct {
	_         [0]func()
	Postal    string  `json:"postal"`
	City      *string `json:"city"`
	CityLocal *string `json:"city_local"`
	// Nil is unknown. An empty slice has no eligible choices. One choice does not imply City.
	Localities        []PostalLocality `json:"localities"`
	District          *string          `json:"district"`
	DistrictName      *string          `json:"district_name"`
	DistrictNameLocal *string          `json:"district_name_local"`
	State             *string          `json:"state"`
	StateName         *string          `json:"state_name"`
	StateNameLocal    *string          `json:"state_name_local"`
	Country           string           `json:"country"`
	CountryName       *string          `json:"country_name"`
	Latitude          *float64         `json:"latitude"`
	Longitude         *float64         `json:"longitude"`
	Timezone          *string          `json:"timezone"`
	Deep              *PostalDeep      `json:"deep,omitempty"`
	// Open ZIP type, including future values. Nil is unknown or absent.
	Type *string `json:"type"`
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
	// FirstName is a suggested first name, not a verified identity.
	FirstName    *string `json:"first_name"`
	NoReply      *bool   `json:"no_reply"`
	Tag          *string `json:"tag"`
	MailProvider *string `json:"mail_provider"`
	_            [0]func()
	Deliverable  *bool   `json:"deliverable"`
	Catchall     *bool   `json:"catchall"`
	Status       *string `json:"status"` // Mailbox status: deliverable, undeliverable or risky. Null when unavailable.
	Reason       *string `json:"reason"` // Why the address received this result, such as mailbox_full or mailbox_not_found. Null when unavailable.
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

type StackTechnology struct {
	_          [0]func()
	Technology string  `json:"technology"`
	Name       string  `json:"name"`
	Version    *string `json:"version"`
}

type Stack struct {
	_         [0]func()
	Domain    string  `json:"domain"`
	URL       string  `json:"url"`
	CheckedAt *string `json:"checked_at"`
	// Nil collections mean the check did not complete. Empty slices mean no matches.
	CMS        []StackTechnology `json:"cms"`
	Servers    []StackTechnology `json:"servers"`
	Frameworks []StackTechnology `json:"frameworks"`
	Ecommerce  []StackTechnology `json:"ecommerce"`
	Analytics  []StackTechnology `json:"analytics"`
	Chat       []StackTechnology `json:"chat"`
	Payments   []StackTechnology `json:"payments"`
	Hosting    []StackTechnology `json:"hosting"`
	Deep       map[string]any    `json:"deep,omitempty"`
}

type DomainRegistration struct {
	_          [0]func()
	Registered bool     `json:"registered"`
	Created    *string  `json:"created"`
	Updated    *string  `json:"updated"`
	Expires    *string  `json:"expires"`
	Registrar  *string  `json:"registrar"`
	Status     []string `json:"status"`
	DNSSEC     *bool    `json:"dnssec"`
}

type DomainDeep struct {
	_            [0]func()
	Registration *DomainRegistration `json:"registration"`
}

type Domain struct {
	_      [0]func()
	Domain string `json:"domain"`
	// True means unregistered, false means registered, nil means unknown.
	Available *bool       `json:"available"`
	Deep      *DomainDeep `json:"deep,omitempty"`
}

type ASN struct {
	_           [0]func()
	ASN         uint32  `json:"asn"`
	Name        *string `json:"name"`
	Country     *string `json:"country"`
	CountryName *string `json:"country_name"`
}


// Card contains network identity; nil brand means unknown or ambiguous.
type Card struct {
	_         [0]func()
	BIN       string    `json:"bin"`
	Brand     *string   `json:"brand"`
	BrandName *string   `json:"brand_name"`
	Logo      string    `json:"logo"`
	Deep      *CardDeep `json:"deep,omitempty"`
}

// CardDeep is optional, partial recorded issuer data. Nil fields mean unknown.
type CardDeep struct {
	_       [0]func()
	Prefix  *string `json:"prefix"`
	Issuer  *string `json:"issuer"`
	Country *string `json:"country"`
	Type    *string `json:"type"`
	Prepaid *bool   `json:"prepaid"`
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


// Language is one language by BCP 47 shortest code or ISO 639-3. Codes are lowercase.
type Language struct {
	_         [0]func()
	Language  string        `json:"language"`
	Name      string        `json:"name"`
	NameLocal *string       `json:"name_local"`
	Script    *string       `json:"script"`
	Direction string        `json:"direction"`
	Deep      *LanguageDeep `json:"deep,omitempty"`
	ISO3      *string       `json:"iso_alpha3,omitempty"`
	Countries []string      `json:"countries"`
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


type TimezoneNextDST struct {
	_            [0]func()
	At           string `json:"at"`
	DST          bool   `json:"dst"`
	Offset       string `json:"offset"`
	Abbreviation string `json:"abbreviation"`
}

// Time contains the local clock and timezone facts. Nil clock fields mean no resolved zone.
type Time = Timezone

// TimeZones lists the serving timezone IDs and pinned rule edition.
type TimeZones struct {
	_                       [0]func()
	TimezoneDatabaseVersion string          `json:"timezone_database_version"`
	Timezones               []string        `json:"timezones"`
	At                      *string         `json:"at,omitempty"`
	Zones                   []TimeZoneEntry `json:"zones,omitempty"`
}

type Timezone struct {
	Location     *TimeLocation `json:"location,omitempty"`
	_            [0]func()
	Timezone     *string                    `json:"timezone"`
	Abbreviation *string                    `json:"abbreviation"`
	Offset       *string                    `json:"offset"`
	DST          *bool                      `json:"dst"`
	Latitude     *float64                   `json:"latitude,omitempty"`
	Longitude    *float64                   `json:"longitude,omitempty"`
	At           *string                    `json:"at,omitempty"`
	Unix         *int64                     `json:"unix,omitempty"`
	To           *TimezoneConversionTarget  `json:"to,omitempty"`
	Targets      []TimezoneConversionTarget `json:"targets,omitempty"`
	Deep         *TimezoneDeep              `json:"deep,omitempty"`
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

type EmojiBrowse struct {
	_          [0]func()
	Emojis     []Emoji `json:"emojis"`
	Total      int     `json:"total"`
	NextOffset *int    `json:"next_offset"`
}
type EmojiVariant struct {
	_     [0]func()
	Emoji string   `json:"emoji"`
	Tones []string `json:"tones"`
}
type EmojiCategory struct {
	_             [0]func()
	Category      string   `json:"category"`
	Subcategories []string `json:"subcategories"`
}
type EmojiCategories struct {
	_          [0]func()
	Categories []EmojiCategory `json:"categories"`
}

// EmojiPosition uses UTF-16 code units, with an exclusive end.
type EmojiPosition struct {
	_     [0]func()
	Start int `json:"start"`
	End   int `json:"end"`
}
type EmojiParseItem struct {
	Emoji
	Positions []EmojiPosition `json:"positions"`
}
type EmojiParse struct {
	_      [0]func()
	Emojis []EmojiParseItem `json:"emojis"`
}

type EmojiSearch struct {
	Total      *int `json:"total"`
	NextOffset *int `json:"next_offset"`
	_          [0]func()
	Q          string  `json:"q"`
	Emojis     []Emoji `json:"emojis"`
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
	// Latest available GDP and GDP per person in current US dollars.
	GDP          *float64 `json:"gdp,omitempty"`
	GDPPerCapita *float64 `json:"gdp_per_capita,omitempty"`
	// Annual percent change; negative values indicate decline.
	PopulationGrowth *float64 `json:"population_growth,omitempty"`
	// Shares of population, from 0 to 100.
	UrbanPopulation *float64 `json:"urban_population,omitempty"`
	InternetUsers   *float64 `json:"internet_users,omitempty"`

	_          [0]func()
	ISO3       *string  `json:"iso_alpha3,omitempty"`
	Numeric    *int     `json:"iso_numeric,omitempty"`
	FullName   *string  `json:"full_name,omitempty"`
	Demonym    *string  `json:"demonym,omitempty"`
	Capital    *string  `json:"capital,omitempty"`
	CapitalLat *float64 `json:"capital_lat,omitempty"`
	CapitalLon *float64 `json:"capital_lon,omitempty"`
	Region     *string  `json:"region,omitempty"`
	Subregion  *string  `json:"subregion,omitempty"`
	Population *int64   `json:"population"`
	// Deprecated: no longer returned; retained for source compatibility and omitted from JSON.
	PopulationPeriod *string  `json:"-"`
	Area             *float64 `json:"area"`
	// Land area in km2.
	LandArea *float64 `json:"land_area"`
	// Water area in km2.
	WaterArea *float64 `json:"water_area"`
	// Coastline length in km. Zero is a known landlocked coastline.
	Coastline *float64 `json:"coastline"`
	// Mean elevation in metres above sea level.
	Elevation    *float64               `json:"elevation,omitempty"`
	LowestPoint  *CountryElevationPoint `json:"lowest_point"`
	HighestPoint *CountryElevationPoint `json:"highest_point"`
	TLD          *string                `json:"tld"`
	Borders      []string               `json:"borders,omitempty"`
	Blocs        []string               `json:"blocs,omitempty"`
	// Levy name, such as VAT, GST or sales tax. Null when unknown or not applicable.
	Tax *CountryTax `json:"tax,omitempty"`
	// Standard country reference rate in percent (19 means 19%). Null is unknown, zero is known zero.
	TaxRate *float64 `json:"-"`
	// Tax registration number mask (9 is a digit, A is a letter). Describes format only.
	TaxIDFormat *string `json:"-"`
	// Anchored tax registration number format regex. A match does not establish registration.
	TaxIDRegex              *string           `json:"-"`
	WeekStart               *string           `json:"week_start"`
	Units                   *string           `json:"units"`
	DrivingSide             *string           `json:"driving_side"`
	Plugs                   []string          `json:"plugs"`
	Voltage                 *int              `json:"voltage"`
	Frequency               *int              `json:"frequency"`
	Emergency               *CountryEmergency `json:"emergency"`
	PostalFormat            *string           `json:"postal_format"`
	PostalRegex             *string           `json:"postal_regex"`
	IOC                     *string           `json:"ioc"`
	FIFA                    *string           `json:"fifa"`
	Plate                   *string           `json:"plate"`
	GDPUSD                  *float64          `json:"gdp_usd,omitempty"`
	GDPPerCapitaUSD         *float64          `json:"gdp_per_capita_usd,omitempty"`
	PopulationGrowthPercent *float64          `json:"population_growth_percent,omitempty"`
	UrbanPopulationPercent  *float64          `json:"urban_population_percent,omitempty"`
	InternetUsersPercent    *float64          `json:"internet_users_percent,omitempty"`
	// Average terrain elevation in meters above sea level.
	ElevationAverage *float64 `json:"elevation_average,omitempty"`
}

type CountryElevationPoint struct {
	_    [0]func()
	Name *string `json:"name"`
	// Elevation in metres above sea level. Values below sea level are negative.
	Elevation float64 `json:"elevation"`
}

type StateDeep struct {
	_          [0]func()
	Population *int64 `json:"population"`
	// Deprecated: no longer returned; retained for source compatibility and omitted from JSON.
	PopulationPeriod *string   `json:"-"`
	Area             *float64  `json:"area"`
	Tax              *StateTax `json:"tax"`
	// Historical Deep members remain readable for existing source consumers.
	FIPS      *string  `json:"-"`
	Capital   *string  `json:"-"`
	AreaCodes []string `json:"-"`
	TaxRate   *float64 `json:"-"`
}

type DistrictDeep struct {
	_          [0]func()
	Population *int64 `json:"population"`
	// Deprecated: no longer returned; retained for source compatibility and omitted from JSON.
	PopulationPeriod *string `json:"-"`
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
	// Deprecated: no longer returned; retained for source compatibility and omitted from JSON.
	PopulationPeriod *string `json:"-"`
	// Official observation date and population basis; null when unproven.
	PopulationObservedOn *string `json:"population_observed_on"`
	PopulationBasis      *string `json:"population_basis"`
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
	// Deprecated: no longer returned; retained for source compatibility and omitted from JSON.
	PopulationPeriod *string `json:"-"`
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


type LanguageDeep struct {
	_ [0]func()
	// Historical Deep read adapters; canonical responses keep these facts in Core.
	ISO3      *string  `json:"-"`
	Countries []string `json:"-"`
}

type NameDeep struct {
	_          [0]func()
	Gender     *string `json:"gender"`
	Salutation *string `json:"salutation"`
	Short      *string `json:"short"`
	Directory  *string `json:"directory"`
	Initials   *string `json:"initials"`
}

type TimeResolutionAlternative struct {
	_      [0]func()
	At     *string `json:"at"`
	Unix   *int64  `json:"unix"`
	Offset *string `json:"offset"`
}

type TimeResolution struct {
	_                 [0]func()
	Kind              *string                     `json:"kind"`
	Policy            *string                     `json:"policy"`
	AdjustmentSeconds *int                        `json:"adjustment_seconds"`
	Alternatives      []TimeResolutionAlternative `json:"alternatives"`
}

type TimezoneDeep struct {
	StandardOffset          *string         `json:"standard_offset"`
	StandardOffsetSeconds   *int            `json:"standard_offset_seconds"`
	DSTOffsetSeconds        *int            `json:"dst_offset_seconds"`
	Season                  *TimeSeason     `json:"season"`
	TimezoneDatabaseVersion *string         `json:"timezone_database_version"`
	Resolution              *TimeResolution `json:"resolution"`
	_                       [0]func()
	Name                    *string          `json:"name"`
	OffsetSeconds           *int             `json:"offset_seconds,omitempty"`
	OffsetMinutes           *int             `json:"offset_minutes"`
	NextDST                 *TimezoneNextDST `json:"next_dst"`
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
	Artwork     map[string]EmojiArtwork `json:"artwork"`
	Subcategory *string                 `json:"subcategory"`
	Base        *string                 `json:"base"`
	Variants    []EmojiVariant          `json:"variants"`
	_           [0]func()
	Codepoints  []string    `json:"codepoints"`
	Hex         *string     `json:"hex"`
	Status      *string     `json:"status"`
	Version     *string     `json:"version"`
	Keywords    []string    `json:"keywords"`
	Skins       []EmojiSkin `json:"skins"`
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
	// Deprecated: no longer returned; retained for source compatibility and omitted from JSON.
	PopulationPeriod *string `json:"-"`
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

// TimeZoneEntry is one supported timezone evaluated at the catalog instant.
type TimeZoneEntry struct {
	_             [0]func()
	Timezone      string   `json:"timezone"`
	Countries     []string `json:"countries"`
	Area          *string  `json:"area"`
	Abbreviation  string   `json:"abbreviation"`
	Offset        string   `json:"offset"`
	OffsetSeconds int      `json:"offset_seconds"`
	DST           bool     `json:"dst"`
	ObservesDST   bool     `json:"observes_dst"`
}
type TimeTransitionState struct {
	_             [0]func()
	At            *string `json:"at"`
	Offset        *string `json:"offset"`
	OffsetSeconds *int    `json:"offset_seconds"`
	Abbreviation  *string `json:"abbreviation"`
	DST           *bool   `json:"dst"`
}
type TimeTransition struct {
	_             [0]func()
	At            *string              `json:"at"`
	Before        *TimeTransitionState `json:"before"`
	After         *TimeTransitionState `json:"after"`
	ChangeSeconds *int                 `json:"change_seconds"`
}
type TimeSeason struct {
	_     [0]func()
	Start *TimeTransition `json:"start"`
	End   *TimeTransition `json:"end"`
}

type TimeLocationInput struct {
	_     [0]func()
	Type  string `json:"type"`
	Value string `json:"value"`
}
type TimeLocationCandidate struct {
	_         [0]func()
	ID        *string  `json:"id"`
	Name      *string  `json:"name"`
	Country   *string  `json:"country"`
	State     *string  `json:"state"`
	Timezone  *string  `json:"timezone"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}
type TimeLocation struct {
	_          [0]func()
	Input      TimeLocationInput       `json:"input"`
	Status     string                  `json:"status"`
	Candidates []TimeLocationCandidate `json:"candidates"`
	Truncated  bool                    `json:"truncated"`
	Source     string                  `json:"source"`
}
















// Industry names for the existing US NAICS response contract.
type Industry = NAICS
type IndustryChild = NAICSChild
type IndustryCorrection = NAICSCorrection
type IndustryDeep = NAICSDeep
type IndustryExclusion = NAICSExclusion
type IndustryMatch = NAICSMatch
type IndustrySearch = NAICSSearch
type IndustrySearchResult = NAICSSearchResult










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
	Deep      map[string]any `json:"deep,omitempty"`
}



// Artwork URLs exist only for files present in a pinned source manifest.
type EmojiArtwork struct {
	_           [0]func()
	License     *string           `json:"license"`
	LicenseURL  *string           `json:"license_url"`
	Attribution *string           `json:"attribution"`
	SVG         *string           `json:"svg"`
	PNG         map[string]string `json:"png"`
}

type EmojiCompleteness struct {
	_               [0]func()
	MissingCount    int `json:"missing_count"`
	UnexpectedCount int `json:"unexpected_count"`
}
type EmojiSource struct {
	_          [0]func()
	ID         string  `json:"id"`
	Version    *string `json:"version"`
	SourceURL  string  `json:"source_url"`
	License    string  `json:"license"`
	LicenseURL string  `json:"license_url"`
}
type EmojiNotice struct {
	_       [0]func()
	ID      string `json:"id"`
	License string `json:"license"`
	Notice  string `json:"notice"`
}
type EmojiFlagLicense struct {
	_           [0]func()
	License     string `json:"license"`
	LicenseURL  string `json:"license_url"`
	Description string `json:"description"`
}

type EmojiArtworkCoverage struct {
	_   [0]func()
	SVG int `json:"svg"`
	PNG int `json:"png"`
}
type EmojiArtworkSource struct {
	_              [0]func()
	Style          string                `json:"style"`
	Version        *string               `json:"version"`
	UnicodeVersion *string               `json:"unicode_version"`
	SourceURL      *string               `json:"source_url"`
	License        *string               `json:"license"`
	LicenseURL     *string               `json:"license_url"`
	Attribution    *string               `json:"attribution"`
	FlagLicense    *EmojiFlagLicense     `json:"flag_license"`
	LoadedCoverage *EmojiArtworkCoverage `json:"loaded_coverage"`
	Notices        []string              `json:"notices"`
}
type EmojiMetadata struct {
	ShortcodeSources map[string]EmojiShortcodeSourceMetadata `json:"shortcode_sources"`
	_                [0]func()
	Edition          string               `json:"edition"`
	MaxVersion       *string              `json:"max_version"`
	SourceVersion    string               `json:"source_version"`
	Total            uint32               `json:"total"`
	Completeness     EmojiCompleteness    `json:"completeness"`
	Sources          []EmojiSource        `json:"sources"`
	ArtworkSources   []EmojiArtworkSource `json:"artwork_sources"`
	Notices          []EmojiNotice        `json:"notices"`
}
type EmojiCatalogScope struct {
	_                 [0]func()
	Filters           map[string]string `json:"filters"`
	CompleteInventory bool              `json:"complete_inventory"`
}
type EmojiCatalog struct {
	Localization *EmojiCatalogLocalization `json:"localization"`
	Profile      *string                   `json:"profile"`
	Scope        *EmojiCatalogScope        `json:"scope"`
	Search       *EmojiCatalogSearchFacts  `json:"search"`
	_            [0]func()
	Metadata     EmojiMetadata                  `json:"metadata"`
	Total        int                            `json:"total"`
	Emojis       []Emoji                        `json:"emojis"`
	Forms        map[string]string              `json:"forms"`
	Shortcodes   map[string]map[string][]string `json:"shortcodes"`
	Legacy       map[string]string              `json:"legacy"`
}

type EmojiNativeSearchFacts struct {
	Order    *int `json:"order"`
	_        [0]func()
	Name     *string  `json:"name"`
	Keywords []string `json:"keywords"`
}
type EmojiCatalogSearchFacts struct {
	LocalizedKeywords []string `json:"localized_keywords"`
	_                 [0]func()
	Native            map[string]EmojiNativeSearchFacts `json:"native"`
}

// Emoji discovery and conversion preserve input UTF-16 positions.
type EmojiUnresolved struct {
	Emoji     *string `json:"emoji"`
	Shortcode *string `json:"shortcode"`
	Start     int     `json:"start"`
	End       int     `json:"end"`
	Reason    string  `json:"reason"`
}
type EmojiConversion struct {
	Text       string            `json:"text"`
	Unresolved []EmojiUnresolved `json:"unresolved"`
}

// Verified localization and shortcode provenance are optional on older servers.
type EmojiCatalogLocalization struct {
	RequestedLocale    *string `json:"requested_locale"`
	ResolvedLocale     *string `json:"resolved_locale"`
	TranslatedNames    int     `json:"translated_names"`
	TranslatedKeywords int     `json:"translated_keywords"`
	FallbackNames      int     `json:"fallback_names"`
	FallbackKeywords   int     `json:"fallback_keywords"`
}
type EmojiShortcodeSourceMetadata struct {
	Records          uint32 `json:"records"`
	MissingRecords   uint32 `json:"missing_records"`
	AmbiguousAliases uint32 `json:"ambiguous_aliases"`
}

// CityParent is a reviewed immediate parent. Nullable names remain unknown.
type CityParent struct {
	_         [0]func()
	ID        string  `json:"id"`
	Name      *string `json:"name"`
	NameLocal *string `json:"name_local"`
	Country   string  `json:"country"`
	State     *string `json:"state"`
	Type      *string `json:"type"`
	Grain     *string `json:"grain"`
}

// CityList is a page from the curated City collection.
type CityList struct {
	_          [0]func()
	Object     string         `json:"object"`
	HasMore    *bool          `json:"has_more,omitempty"`
	URL        *string        `json:"url,omitempty"`
	Q          *string        `json:"q,omitempty"`
	Country    *string        `json:"country,omitempty"`
	State      *string        `json:"state,omitempty"`
	Type       *string        `json:"type,omitempty"`
	Grain      *string        `json:"grain,omitempty"`
	Data       []CityResource `json:"data"`
	NextCursor *string        `json:"next_cursor"`
}

// CityNearbyCoords is a page of cities around the supplied coordinate.
type CityNearbyCoords struct {
	_          [0]func()
	Object     string         `json:"object"`
	HasMore    *bool          `json:"has_more,omitempty"`
	URL        *string        `json:"url,omitempty"`
	Latitude   float64        `json:"latitude"`
	Longitude  float64        `json:"longitude"`
	Radius     float64        `json:"radius"`
	Unit       string         `json:"unit"`
	Data       []CityResource `json:"data"`
	NextCursor *string        `json:"next_cursor"`
}



type StateTax struct {
	_    [0]func()
	Name *string `json:"name"`
	// Reference rate in percent.
	Rate *float64 `json:"rate"`
}



type Currency struct {
	_ [0]func()
	Currency string `json:"currency"`
	ISONumeric *int `json:"iso_numeric"`
	Name string `json:"name"`
	NamePlural *string `json:"name_plural"`
	Symbol *string `json:"symbol"`
	SymbolNative *string `json:"symbol_native"`
	Digits *int `json:"digits"`
	Countries []string `json:"countries"`
	Deep *CurrencyDeep `json:"deep,omitempty"`
}

type CurrencyRate struct {
	_         [0]func()
	Base      string   `json:"base"`
	Quote     string   `json:"quote"`
	Rate      float64  `json:"rate"`
	Date      string   `json:"date"`
	Amount    *float64 `json:"amount,omitempty"`
	Converted *float64 `json:"converted,omitempty"`
}

type CurrencyDeep struct { _ [0]func() }
