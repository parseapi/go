package parseapi

import "encoding/json"

func countryFlag(raw map[string]json.RawMessage) {
	if _, ok := raw["flag"]; !ok {
		if value, found := raw["emoji"]; found {
			raw["flag"] = value
		}
	}
	delete(raw, "emoji")
}

func (item *ContinentCountryItem) UnmarshalJSON(data []byte) error {
	type wire ContinentCountryItem
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	countryFlag(raw)
	normalized, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	var value wire
	if err = json.Unmarshal(normalized, &value); err != nil {
		return err
	}
	*item = ContinentCountryItem(value)
	item.Emoji = item.Flag
	return nil
}

// UnmarshalJSON accepts historical scalar capital/name spellings and current objects.
func (country *Country) UnmarshalJSON(data []byte) error {
	type wire Country
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	_, present := raw["name_full"]
	if !present {
		if legacy, ok := raw["full_name"]; ok {
			raw["name_full"] = legacy
		}
	}
	countryFlag(raw)
	if value, ok := raw["capital"]; ok && len(value) > 0 && value[0] != '{' && string(value) != "null" {
		capital := map[string]json.RawMessage{"name": value, "lat": json.RawMessage("null"), "lon": json.RawMessage("null")}
		if value, ok := raw["capital_lat"]; ok {
			capital["lat"] = value
		}
		if value, ok := raw["capital_lon"]; ok {
			capital["lon"] = value
		}
		encoded, err := json.Marshal(capital)
		if err != nil {
			return err
		}
		raw["capital"] = encoded
	} else if ok && string(value) == "null" {
		_, lat := raw["capital_lat"]
		_, lon := raw["capital_lon"]
		if lat || lon {
			capital := map[string]json.RawMessage{"name": value, "lat": json.RawMessage("null"), "lon": json.RawMessage("null")}
			if lat {
				capital["lat"] = raw["capital_lat"]
			}
			if lon {
				capital["lon"] = raw["capital_lon"]
			}
			encoded, err := json.Marshal(capital)
			if err != nil {
				return err
			}
			raw["capital"] = encoded
		}
	}
	normalized, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	var result wire
	if err = json.Unmarshal(normalized, &result); err != nil {
		return err
	}
	*country = Country(result)
	country.nameFullPresent = present
	country.FullName = country.NameFull
	country.Emoji = country.Flag
	if country.Capital != nil {
		country.CapitalLat = country.Capital.Lat
		country.CapitalLon = country.Capital.Lon
	}
	return nil
}

func (country Country) MarshalJSON() ([]byte, error) {
	type wire Country
	name := country.NameFull
	if !country.nameFullPresent && name == nil {
		name = country.FullName
	}
	return json.Marshal(struct {
		wire
		NameFull *string `json:"name_full"`
	}{wire: wire(country), NameFull: name})
}

func (detail *CountryDeep) UnmarshalJSON(data []byte) error {
	type wire CountryDeep
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	value, present := raw["tax"]
	_, rate := raw["tax_rate"]
	_, format := raw["tax_id_format"]
	_, regex := raw["tax_id_regex"]
	legacy := (present && len(value) > 0 && value[0] != '{' && string(value) != "null") || ((rate || format || regex) && (!present || string(value) == "null"))
	if legacy {
		tax := map[string]json.RawMessage{"name": json.RawMessage("null"), "rate": json.RawMessage("null")}
		if present {
			tax["name"] = value
		}
		if rate {
			tax["rate"] = raw["tax_rate"]
		}
		encoded, err := json.Marshal(tax)
		if err != nil {
			return err
		}
		raw["tax"] = encoded
	}
	normalized, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	var result wire
	if err = json.Unmarshal(normalized, &result); err != nil {
		return err
	}
	*detail = CountryDeep(result)
	if detail.Tax != nil {
		detail.TaxRate = detail.Tax.Rate
	}
	// These published read aliases describe only historical flat payloads.
	if legacy {
		if format {
			if err := json.Unmarshal(raw["tax_id_format"], &detail.TaxIDFormat); err != nil {
				return err
			}
		}
		if regex {
			if err := json.Unmarshal(raw["tax_id_regex"], &detail.TaxIDRegex); err != nil {
				return err
			}
		}
	}
	return nil
}
