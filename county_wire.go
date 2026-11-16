package parseapi

import "encoding/json"

// decodeCountyWire accepts historical District keys while canonical presence (even null) wins.
func decodeCountyWire(data []byte, target any) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for _, pair := range [][2]string{{"county", "district"}, {"county_name", "district_name"}, {"county_name_local", "district_name_local"}, {"counties", "districts"}} {
		isoWireAlias(raw, pair[0], pair[1])
	}
	normalized, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	return json.Unmarshal(normalized, target)
}

func (value *IPDeep) UnmarshalJSON(data []byte) error {
	type wire IPDeep
	var result wire
	if err := decodeCountyWire(data, &result); err != nil {
		return err
	}
	*value = IPDeep(result)
	value.District = value.County
	value.DistrictName = value.CountyName
	return nil
}

func (value *IP) UnmarshalJSON(data []byte) error {
	type wire IP
	var result wire
	if err := decodeCountyWire(data, &result); err != nil {
		return err
	}
	*value = IP(result)
	value.District = value.County
	value.DistrictName = value.CountyName
	return nil
}

func (value *StateCountyItem) UnmarshalJSON(data []byte) error {
	type wire StateCountyItem
	var result wire
	if err := decodeCountyWire(data, &result); err != nil {
		return err
	}
	*value = StateCountyItem(result)
	value.District = value.County
	return nil
}

func (value *StateCounties) UnmarshalJSON(data []byte) error {
	type wire StateCounties
	var result wire
	if err := decodeCountyWire(data, &result); err != nil {
		return err
	}
	*value = StateCounties(result)
	value.Districts = value.Counties
	return nil
}

func (value *County) UnmarshalJSON(data []byte) error {
	type wire County
	var result wire
	if err := decodeCountyWire(data, &result); err != nil {
		return err
	}
	*value = County(result)
	value.District = value.County
	return nil
}

func (value *City) UnmarshalJSON(data []byte) error {
	type wire City
	var result wire
	if err := decodeCountyWire(data, &result); err != nil {
		return err
	}
	*value = City(result)
	value.District = value.County
	value.DistrictName = value.CountyName
	return nil
}

func (value *Postal) UnmarshalJSON(data []byte) error {
	type wire Postal
	var result wire
	if err := decodeCountyWire(data, &result); err != nil {
		return err
	}
	*value = Postal(result)
	value.District = value.County
	value.DistrictName = value.CountyName
	value.DistrictNameLocal = value.CountyNameLocal
	return nil
}

func (value *Point) UnmarshalJSON(data []byte) error {
	type wire Point
	var result wire
	if err := decodeCountyWire(data, &result); err != nil {
		return err
	}
	*value = Point(result)
	value.District = value.County
	value.DistrictName = value.CountyName
	return nil
}

func (value *Address) UnmarshalJSON(data []byte) error {
	type wire Address
	var result wire
	if err := decodeCountyWire(data, &result); err != nil {
		return err
	}
	*value = Address(result)
	value.District = value.County
	value.DistrictName = value.CountyName
	return nil
}

func (value *CityResource) UnmarshalJSON(data []byte) error {
	type wire CityResource
	var result wire
	if err := decodeCountyWire(data, &result); err != nil {
		return err
	}
	*value = CityResource(result)
	value.District = value.County
	value.DistrictName = value.CountyName
	return nil
}

// CityNearest embeds City; decode the distance fields explicitly rather than inheriting City.UnmarshalJSON.
func (value *CityNearest) UnmarshalJSON(data []byte) error {
	type cityWire City
	var wire struct {
		cityWire
		Distance   float64 `json:"distance"`
		DistanceMi float64 `json:"distance_mi"`
	}
	if err := decodeCountyWire(data, &wire); err != nil {
		return err
	}
	value.City = City(wire.cityWire)
	value.District = value.County
	value.DistrictName = value.CountyName
	value.Distance = wire.Distance
	value.DistanceMi = wire.DistanceMi
	return nil
}
