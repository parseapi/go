package parseapi

import "encoding/json"

// UnmarshalJSON accepts older minute payloads only when the canonical key is absent.
func (minute *WeatherMinute) UnmarshalJSON(data []byte) error {
	type plain WeatherMinute
	var decoded WeatherMinute
	wire := struct {
		*plain
		Rate       json.RawMessage `json:"precipitation_rate"`
		RateIn     json.RawMessage `json:"precipitation_rate_in"`
		LegacyRate *float64        `json:"precipitation"`
		LegacyIn   *float64        `json:"precipitation_in"`
	}{plain: (*plain)(&decoded)}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	decoded.PrecipitationRate = wire.LegacyRate
	if wire.Rate != nil {
		if err := json.Unmarshal(wire.Rate, &decoded.PrecipitationRate); err != nil {
			return err
		}
	}
	decoded.PrecipitationRateIn = wire.LegacyIn
	if wire.RateIn != nil {
		if err := json.Unmarshal(wire.RateIn, &decoded.PrecipitationRateIn); err != nil {
			return err
		}
	}
	*minute = decoded
	return nil
}
