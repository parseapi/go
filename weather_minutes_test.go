package parseapi

import (
	"encoding/json"
	"testing"
)

func TestMinuteRatesAndHistoricalAmounts(t *testing.T) {
	var weather Weather
	err := json.Unmarshal([]byte(`{"minutes":[{"at":"2026-10-27T12:00:00Z","precipitation_rate":1.8,"precipitation_rate_in":0.071,"type":"rain"},{"at":"2026-10-27T12:01:00Z","precipitation_rate":0,"precipitation_rate_in":0},{"at":"2026-10-27T12:02:00Z","precipitation_rate":null,"precipitation_rate_in":null}],"history":{"precipitation":12.3,"precipitation_in":0.484}}`), &weather)
	if err != nil {
		t.Fatal(err)
	}
	rows := weather.Minutes
	if *rows[0].PrecipitationRate != 1.8 || *rows[0].PrecipitationRateIn != 0.071 || *rows[1].PrecipitationRate != 0 || *rows[1].PrecipitationRateIn != 0 || rows[2].PrecipitationRate != nil || rows[2].PrecipitationRateIn != nil {
		t.Fatalf("minute rates: %#v", rows)
	}
	if *weather.History.Precipitation != 12.3 || *weather.History.PrecipitationIn != 0.484 {
		t.Fatalf("historical amounts: %#v", weather.History)
	}
}
