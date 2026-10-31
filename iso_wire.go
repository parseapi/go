package parseapi

import "encoding/json"

// isoWireAlias preserves an explicit canonical null and accepts historical wire keys.
func isoWireAlias(raw map[string]json.RawMessage, canonical string, legacy ...string) {
	if _, present := raw[canonical]; !present {
		for _, key := range legacy {
			if value, exists := raw[key]; exists {
				raw[canonical] = value
				break
			}
		}
	}
	for _, key := range legacy {
		delete(raw, key)
	}
}

func decodeISOWire(data []byte, target any, canonical string, legacy ...string) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	isoWireAlias(raw, canonical, legacy...)
	normalized, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	return json.Unmarshal(normalized, target)
}

func (value *Language) UnmarshalJSON(data []byte) error {
	type wire Language
	decoded := wire(*value)
	if err := decodeISOWire(data, &decoded, "iso_alpha3", "iso3"); err != nil {
		return err
	}
	*value = Language(decoded)
	return nil
}

func (value *State) UnmarshalJSON(data []byte) error {
	type wire State
	decoded := wire(*value)
	if err := decodeISOWire(data, &decoded, "iso_subdivision", "iso", "iso_3166_2"); err != nil {
		return err
	}
	decoded.ISO31662 = decoded.ISO
	*value = State(decoded)
	return nil
}

func (value *LanguageDeep) UnmarshalJSON(data []byte) error {
	var wire struct {
		ISO3      *string  `json:"iso_alpha3"`
		Countries []string `json:"countries"`
	}
	if err := decodeISOWire(data, &wire, "iso_alpha3", "iso3"); err != nil {
		return err
	}
	value.ISO3, value.Countries = wire.ISO3, wire.Countries
	return nil
}

func (value *StateDeep) UnmarshalJSON(data []byte) error {
	type wire StateDeep
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	tax, present := raw["tax"]
	_, ratePresent := raw["tax_rate"]
	if (present && len(tax) > 0 && tax[0] != '{' && string(tax) != "null") || (ratePresent && (!present || string(tax) == "null")) {
		name := tax
		if !present {
			name = json.RawMessage("null")
		}
		rate := raw["tax_rate"]
		if !ratePresent {
			rate = json.RawMessage("null")
		}
		nested, err := json.Marshal(map[string]json.RawMessage{"name": name, "rate": rate})
		if err != nil {
			return err
		}
		raw["tax"] = nested
	}
	normalized, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	var decoded wire
	if err = json.Unmarshal(normalized, &decoded); err != nil {
		return err
	}
	for key, target := range map[string]any{"fips": &decoded.FIPS, "capital": &decoded.Capital, "area_codes": &decoded.AreaCodes, "tax_rate": &decoded.TaxRate} {
		if data, present := raw[key]; present {
			if err = json.Unmarshal(data, target); err != nil {
				return err
			}
		}
	}
	if decoded.Tax != nil {
		decoded.TaxRate = decoded.Tax.Rate
	}
	*value = StateDeep(decoded)
	return nil
}
