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

func (detail *CountryDeep) UnmarshalJSON(data []byte) error {
	type wire CountryDeep
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	isoWireAlias(raw, "iso_alpha3", "iso3")
	isoWireAlias(raw, "iso_numeric", "numeric")
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

// Decode older capital objects while emitting only the canonical coordinate keys.
func decodeCapitalCoordinates(data []byte, target any) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	isoWireAlias(raw, "latitude", "lat")
	isoWireAlias(raw, "longitude", "lon")
	delete(raw, "lat")
	delete(raw, "lon")
	normalized, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	return json.Unmarshal(normalized, target)
}
func (value *CountryCapital) UnmarshalJSON(data []byte) error {
	type wire CountryCapital
	var decoded wire
	if err := decodeCapitalCoordinates(data, &decoded); err != nil {
		return err
	}
	*value = CountryCapital(decoded)
	return nil
}
func (value *StateCapital) UnmarshalJSON(data []byte) error {
	type wire StateCapital
	var decoded wire
	if err := decodeCapitalCoordinates(data, &decoded); err != nil {
		return err
	}
	*value = StateCapital(decoded)
	return nil
}

