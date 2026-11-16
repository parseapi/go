package parseapi

// CityMatch describes the explicit search suggestion evidence.
type CityMatch struct {
	_            [0]func()
	Type         string  `json:"type"`
	Source       string  `json:"source,omitempty"`
	Value        string  `json:"value"`
	Language     *string `json:"language"`
	EditDistance *int    `json:"edit_distance,omitempty"`
}

// CityResource is canonical JSON. Unknown facts are nil; ID/Object remain.
type CityResource struct {
	Match       *CityMatch `json:"match"`
	_           [0]func()
	ID          string      `json:"id"`
	Object      string      `json:"object"`
	Name        *string     `json:"name"`
	NameLocal   *string     `json:"name_local"`
	Country     *string     `json:"country"`
	CountryName *string     `json:"country_name"`
	State       *string     `json:"state"`
	StateName   *string     `json:"state_name"`
	County      *string     `json:"county"`
	CountyName  *string     `json:"county_name"`
	Type        *string     `json:"type"`
	Grain       *string     `json:"grain"`
	Latitude    *float64    `json:"latitude"`
	Longitude   *float64    `json:"longitude"`
	Timezone    *string     `json:"timezone"`
	Distance    *float64    `json:"distance"`
	DistanceMi  *float64    `json:"distance_mi"`
	Parent      *CityParent `json:"parent"`

	NameRomanized        *string  `json:"name_romanized"`
	RomanizationLanguage *string  `json:"romanization_language"`
	CapitalOf            *string  `json:"capital_of"`
	Elevation            *float64 `json:"elevation"`
	ElevationFt          *float64 `json:"elevation_ft"`
	Population           *int64   `json:"population"`
	PopulationObservedOn *string  `json:"population_observed_on"`
	PopulationBasis      *string  `json:"population_basis"`
	Area                 *float64 `json:"area"`
	LandArea             *float64 `json:"land_area"`
	WaterArea            *float64 `json:"water_area"`

	// Deprecated: use County. Retained as a read adapter, never a second JSON field.
	District *string `json:"-"`
	// Deprecated: use CountyName. Retained as a read adapter, never a second JSON field.
	DistrictName *string `json:"-"`
}

// CityAmbiguityCandidate is a bounded safe identity, never Deep data.
type CityAmbiguityCandidate struct {
	_           [0]func()
	ID          string  `json:"id"`
	Object      string  `json:"object"`
	Name        *string `json:"name"`
	NameLocal   *string `json:"name_local"`
	Country     *string `json:"country"`
	CountryName *string `json:"country_name"`
	State       *string `json:"state"`
	StateName   *string `json:"state_name"`
	Type        *string `json:"type"`
}

// CityRelations holds reviewed relationships. An empty page does not prove no relatives.
type CityRelations struct {
	_            [0]func()
	Object       string         `json:"object"`
	Data         []CityResource `json:"data"`
	HasMore      *bool          `json:"has_more,omitempty"`
	URL          *string        `json:"url,omitempty"`
	NextCursor   *string        `json:"next_cursor"`
	CityID       string         `json:"city_id"`
	Relationship string         `json:"relationship"`
	Coverage     string         `json:"coverage"`
	ParentID     *string        `json:"parent_id"`
}
type CityDistance struct {
	_          [0]func()
	Object     string   `json:"object"`
	From       string   `json:"from"`
	To         string   `json:"to"`
	Distance   *float64 `json:"distance"`
	DistanceMi *float64 `json:"distance_mi"`
	Method     string   `json:"method"`
}
