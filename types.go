package parseapi

// Response types for the parseAPI public API. Shapes are append-only
// upstream, so these only ever grow. Nullable fields are pointers.
// Deep objects follow the triad: nil when not requested, empty when
// requested but locked, populated when unlocked.

type IPDeep struct {
	State      *string `json:"state"`
	City       *string `json:"city"`
	Registry   *string `json:"registry"`
	Datacenter *bool   `json:"datacenter"`
	Relay      *bool   `json:"relay"`
	Tor        *bool   `json:"tor"`
	Provider   *string `json:"provider"`
}

type IP struct {
	IP          string  `json:"ip"`
	Country     *string `json:"country"`
	CountryName *string `json:"country_name"`
	Continent   *string `json:"continent"`
	ASN         *string `json:"asn"`
	ASNName     *string `json:"asn_name"`
	Deep        *IPDeep `json:"deep,omitempty"`
}

type Continent struct {
	Continent  string   `json:"continent"`
	Name       string   `json:"name"`
	Region     string   `json:"region"`
	Subregion  string   `json:"subregion"`
	Population *int64   `json:"population"`
	Area       *float64 `json:"area"`
	Emoji      string   `json:"emoji"`
}

type ContinentCountryItem struct {
	Country     string  `json:"country"`
	Name        string  `json:"name"`
	Emoji       *string `json:"emoji"`
	CallingCode *string `json:"calling_code"`
}

type ContinentCountries struct {
	Continent string                 `json:"continent"`
	Countries []ContinentCountryItem `json:"countries"`
}

type Country struct {
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
}

type CountryStateItem struct {
	State string  `json:"state"`
	Name  string  `json:"name"`
	Type  *string `json:"type"`
}

type CountryStates struct {
	Country string             `json:"country"`
	States  []CountryStateItem `json:"states"`
}

type State struct {
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
}

type StateDistrictItem struct {
	District string  `json:"district"`
	Name     string  `json:"name"`
	Type     *string `json:"type"`
}

type StateDistricts struct {
	State       string              `json:"state"`
	StateName   *string             `json:"state_name"`
	Country     string              `json:"country"`
	CountryName *string             `json:"country_name"`
	Districts   []StateDistrictItem `json:"districts"`
}

type District struct {
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
	LandArea    *float64 `json:"land_area"`
	WaterArea   *float64 `json:"water_area"`
}

type City struct {
	Name       string   `json:"name"`
	LocalName  *string  `json:"local_name"`
	State      *string  `json:"state"`
	StateName  *string  `json:"state_name"`
	Country    string   `json:"country"`
	Latitude   *float64 `json:"latitude"`
	Longitude  *float64 `json:"longitude"`
	Population *int64   `json:"population"`
	Timezone   *string  `json:"timezone"`
	// ID is the minted parse id (city_ + 12 chars). Stable pin via /city/id/{id}.
	ID string `json:"id"`
}

// CityNearest is a City plus the distance from the query point.
type CityNearest struct {
	City
	Distance   float64 `json:"distance"`
	DistanceMi float64 `json:"distance_mi"`
}

type CitySearch struct {
	Q       string `json:"q"`
	Country string `json:"country,omitempty"`
	State   string `json:"state,omitempty"`
	Cities  []City `json:"cities"`
}

type Postal struct {
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
	Latitude          *float64 `json:"latitude"`
	Longitude         *float64 `json:"longitude"`
	Elevation         *float64 `json:"elevation"`
	ElevationFt       *float64 `json:"elevation_ft"`
	Population        *int64   `json:"population"`
	Timezone          *string  `json:"timezone"`
	Currency          *string  `json:"currency"`
	Neighbors         []string `json:"neighbors"`
}

type PostalNearbyItem struct {
	Postal     string  `json:"postal"`
	City       *string `json:"city"`
	State      *string `json:"state"`
	Country    string  `json:"country"`
	Distance   float64 `json:"distance"`
	DistanceMi float64 `json:"distance_mi"`
}

type PostalNearby struct {
	Postal  string             `json:"postal"`
	Country string             `json:"country"`
	Radius  float64            `json:"radius"`
	Unit    string             `json:"unit"`
	Nearby  []PostalNearbyItem `json:"nearby"`
}

type PostalDistanceEnd struct {
	Postal string  `json:"postal"`
	City   *string `json:"city"`
}

type PostalDistance struct {
	Country    string            `json:"country"`
	From       PostalDistanceEnd `json:"from"`
	To         PostalDistanceEnd `json:"to"`
	Distance   float64           `json:"distance"`
	DistanceMi float64           `json:"distance_mi"`
}

type EmailDeep struct {
	Deliverable *bool `json:"deliverable"`
	Catchall    *bool `json:"catchall"`
}

type Email struct {
	Email       string     `json:"email"`
	Valid       bool       `json:"valid"`
	Domain      *string    `json:"domain"`
	DomainValid *bool      `json:"domain_valid"`
	Role        bool       `json:"role"`
	Disposable  bool       `json:"disposable"`
	Deep        *EmailDeep `json:"deep,omitempty"`
}

type PhoneDeep struct {
	Type    *string `json:"type"`
	Carrier *string `json:"carrier"`
	// Burner reports whether the carrier is a known burner number app. Nil when carrier is unknown.
	Burner    *bool   `json:"burner"`
	City      *string `json:"city"`
	State     *string `json:"state"`
	StateName *string `json:"state_name"`
}

type Phone struct {
	Phone         *string    `json:"phone"`
	Valid         bool       `json:"valid"`
	Country       *string    `json:"country"`
	National      *string    `json:"national"`
	International *string    `json:"international"`
	Deep          *PhoneDeep `json:"deep,omitempty"`
}

type MXRecord struct {
	Priority int    `json:"priority"`
	Host     string `json:"host"`
}

type DomainRegistration struct {
	Registered bool     `json:"registered"`
	Created    *string  `json:"created"`
	Updated    *string  `json:"updated"`
	Expires    *string  `json:"expires"`
	Registrar  *string  `json:"registrar"`
	Status     []string `json:"status"`
	DNSSEC     bool     `json:"dnssec"`
}

type DomainDeep struct {
	A            []string            `json:"a"`
	AAAA         []string            `json:"aaaa"`
	NS           []string            `json:"ns"`
	MX           []MXRecord          `json:"mx"`
	TXT          []string            `json:"txt"`
	Provider     *string             `json:"provider"`
	Registration *DomainRegistration `json:"registration"`
}

type Domain struct {
	Domain    string      `json:"domain"`
	Available bool        `json:"available"`
	Deep      *DomainDeep `json:"deep,omitempty"`
}

type MX struct {
	Domain string     `json:"domain"`
	MX     []MXRecord `json:"mx"`
}

type UseragentDeviceDeep struct {
	Type        *string `json:"type"`
	Brand       *string `json:"brand"`
	Model       *string `json:"model"`
	CPU         *string `json:"cpu"`
	Touchscreen *bool   `json:"touchscreen"`
}

type UseragentOSDeep struct {
	Name     *string `json:"name"`
	Version  *string `json:"version"`
	Platform *string `json:"platform"`
}

type UseragentBrowserBrand struct {
	Brand   string `json:"brand"`
	Version string `json:"version"`
}

type UseragentBrowserDeep struct {
	Name    *string                 `json:"name"`
	Version *string                 `json:"version"`
	Type    *string                 `json:"type"`
	Brands  []UseragentBrowserBrand `json:"brands,omitempty"`
}

type UseragentEngineDeep struct {
	Name    *string `json:"name"`
	Version *string `json:"version"`
}

type UseragentDeep struct {
	Device   *UseragentDeviceDeep  `json:"device"`
	OS       *UseragentOSDeep      `json:"os"`
	Browser  *UseragentBrowserDeep `json:"browser"`
	Engine   *UseragentEngineDeep  `json:"engine"`
	Headless *bool                 `json:"headless"`
	AI       *bool                 `json:"ai,omitempty"`
}

type Useragent struct {
	Useragent string         `json:"useragent"`
	Device    *string        `json:"device"`
	OS        *string        `json:"os"`
	Browser   *string        `json:"browser"`
	Bot       bool           `json:"bot"`
	Mobile    bool           `json:"mobile"`
	Deep      *UseragentDeep `json:"deep,omitempty"`
}

type Currency struct {
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
	Base   string  `json:"base"`
	Quote  string  `json:"quote"`
	Rate   float64 `json:"rate"`
	Date   string  `json:"date"`
	Source string  `json:"source,omitempty"`
}

type TimezoneNextDST struct {
	At           string `json:"at"`
	DST          bool   `json:"dst"`
	Offset       string `json:"offset"`
	Abbreviation string `json:"abbreviation"`
}

type Timezone struct {
	Timezone      string           `json:"timezone"`
	Name          *string          `json:"name"`
	Abbreviation  string           `json:"abbreviation"`
	Offset        string           `json:"offset"`
	OffsetMinutes int              `json:"offset_minutes"`
	DST           bool             `json:"dst"`
	NextDST       *TimezoneNextDST `json:"next_dst"`
}

type Holiday struct {
	Date       string   `json:"date"`
	Name       string   `json:"name"`
	LocalName  *string  `json:"local_name"`
	Type       string   `json:"type"`
	Regions    []string `json:"regions"`
	Substitute bool     `json:"substitute"`
}

type HolidayYear struct {
	Country  string    `json:"country"`
	Year     int       `json:"year"`
	Holidays []Holiday `json:"holidays"`
}

type HolidayDate struct {
	Country string   `json:"country"`
	Date    string   `json:"date"`
	Holiday *Holiday `json:"holiday"`
}

type Elevation struct {
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	Elevation   *float64 `json:"elevation"`
	ElevationFt *float64 `json:"elevation_ft"`
	Resolution  *float64 `json:"resolution"`
}

type PointDeep struct {
	City     *CityNearest `json:"city"`
	Timezone *Timezone    `json:"timezone"`
}

type Point struct {
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
	Event    string  `json:"event"`
	Severity *string `json:"severity"`
	Urgency  *string `json:"urgency"`
	Headline *string `json:"headline"`
	Onset    *string `json:"onset"`
	Expires  *string `json:"expires"`
}

type WeatherDeep struct {
	Forecast []WeatherForecastPeriod `json:"forecast"`
	Alerts   []WeatherAlert          `json:"alerts"`
}

type Weather struct {
	Latitude          float64      `json:"latitude"`
	Longitude         float64      `json:"longitude"`
	Temperature       *float64     `json:"temperature"`
	TemperatureF      *float64     `json:"temperature_f"`
	FeelsLike         *float64     `json:"feels_like"`
	FeelsLikeF        *float64     `json:"feels_like_f"`
	Dewpoint          *float64     `json:"dewpoint"`
	DewpointF         *float64     `json:"dewpoint_f"`
	Humidity          *float64     `json:"humidity"`
	WindSpeed         *float64     `json:"wind_speed"`
	WindSpeedMph      *float64     `json:"wind_speed_mph"`
	WindGust          *float64     `json:"wind_gust"`
	WindGustMph       *float64     `json:"wind_gust_mph"`
	WindDirection     *float64     `json:"wind_direction"`
	Pressure          *float64     `json:"pressure"`
	PressureInhg      *float64     `json:"pressure_inhg"`
	Visibility        *float64     `json:"visibility"`
	VisibilityMi      *float64     `json:"visibility_mi"`
	Condition         *string      `json:"condition"`
	ConditionName     *string      `json:"condition_name"`
	ConditionEmoji    *string      `json:"condition_emoji"`
	ObservedAt        *string      `json:"observed_at"`
	Station           string       `json:"station"`
	StationName       *string      `json:"station_name"`
	StationDistance   float64      `json:"station_distance"`
	StationDistanceMi float64      `json:"station_distance_mi"`
	Source            string       `json:"source"`
	SourceName        *string      `json:"source_name"`
	Deep              *WeatherDeep `json:"deep,omitempty"`
}

type EmojiSkin struct {
	Emoji   string  `json:"emoji"`
	Tone    string  `json:"tone"`
	Unicode *string `json:"unicode"`
	Hex     *string `json:"hex"`
}

type Emoji struct {
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
	Q      string  `json:"q"`
	Emojis []Emoji `json:"emojis"`
}
