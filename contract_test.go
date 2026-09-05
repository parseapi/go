// This is the frozen consumer contract after the final pre-launch cleanup.
// Keep these method types working when adding future features. Do not regenerate
// this fixture to hide a public signature change.
package parseapi_test

import (
	"context"
	parseapi "github.com/parseapi/go"
)

var _ func(*parseapi.Client, context.Context, string, ...parseapi.IPOptions) (*parseapi.IP, error) = (*parseapi.Client).IP
var _ func(*parseapi.Client, context.Context, ...parseapi.IPSelfOptions) (*parseapi.IP, error) = (*parseapi.Client).IPSelf
var _ func(*parseapi.Client, context.Context, string, ...parseapi.ContinentOptions) (*parseapi.Continent, error) = (*parseapi.Client).Continent
var _ func(*parseapi.Client, context.Context, string, ...parseapi.ContinentCountriesOptions) (*parseapi.ContinentCountries, error) = (*parseapi.Client).ContinentCountries
var _ func(*parseapi.Client, context.Context, string, ...parseapi.BlocOptions) (*parseapi.Bloc, error) = (*parseapi.Client).Bloc
var _ func(*parseapi.Client, context.Context, string, ...parseapi.BlocCountriesOptions) (*parseapi.BlocCountries, error) = (*parseapi.Client).BlocCountries
var _ func(*parseapi.Client, context.Context, string, ...parseapi.CountryOptions) (*parseapi.Country, error) = (*parseapi.Client).Country
var _ func(*parseapi.Client, context.Context, string, ...parseapi.CountryStatesOptions) (*parseapi.CountryStates, error) = (*parseapi.Client).CountryStates
var _ func(*parseapi.Client, context.Context, string, ...parseapi.StateOptions) (*parseapi.State, error) = (*parseapi.Client).State
var _ func(*parseapi.Client, context.Context, string, ...parseapi.StateDistrictsOptions) (*parseapi.StateDistricts, error) = (*parseapi.Client).StateDistricts
var _ func(*parseapi.Client, context.Context, string, ...parseapi.DistrictOptions) (*parseapi.District, error) = (*parseapi.Client).District
var _ func(*parseapi.Client, context.Context, string, ...parseapi.CityOptions) (*parseapi.City, error) = (*parseapi.Client).City
var _ func(*parseapi.Client, context.Context, string, ...parseapi.CityIDOptions) (*parseapi.City, error) = (*parseapi.Client).CityID
var _ func(*parseapi.Client, context.Context, string, ...parseapi.CitySearchOptions) (*parseapi.CitySearch, error) = (*parseapi.Client).CitySearch
var _ func(*parseapi.Client, context.Context, float64, float64, ...parseapi.CityNearestOptions) (*parseapi.CityNearest, error) = (*parseapi.Client).CityNearest
var _ func(*parseapi.Client, context.Context, string, ...parseapi.CityNearbyOptions) (*parseapi.CityNearby, error) = (*parseapi.Client).CityNearby
var _ func(*parseapi.Client, context.Context, string, ...parseapi.PostalOptions) (*parseapi.Postal, error) = (*parseapi.Client).Postal
var _ func(*parseapi.Client, context.Context, string, ...parseapi.PostalNearbyOptions) (*parseapi.PostalNearby, error) = (*parseapi.Client).PostalNearby
var _ func(*parseapi.Client, context.Context, string, string, ...parseapi.PostalDistanceOptions) (*parseapi.PostalDistance, error) = (*parseapi.Client).PostalDistance
var _ func(*parseapi.Client, context.Context, string, ...parseapi.EmailOptions) (*parseapi.Email, error) = (*parseapi.Client).Email
var _ func(*parseapi.Client, context.Context, string, ...parseapi.VATOptions) (*parseapi.VAT, error) = (*parseapi.Client).VAT
var _ func(*parseapi.Client, context.Context, string, ...parseapi.IBANOptions) (*parseapi.IBAN, error) = (*parseapi.Client).IBAN
var _ func(*parseapi.Client, context.Context, string, ...parseapi.NPIOptions) (*parseapi.NPI, error) = (*parseapi.Client).NPI
var _ func(*parseapi.Client, context.Context, string, ...parseapi.PhoneOptions) (*parseapi.Phone, error) = (*parseapi.Client).Phone
var _ func(*parseapi.Client, context.Context, string, ...parseapi.CarrierOptions) (*parseapi.Carrier, error) = (*parseapi.Client).Carrier
var _ func(*parseapi.Client, context.Context, string, ...parseapi.CallerOptions) (*parseapi.Caller, error) = (*parseapi.Client).Caller
var _ func(*parseapi.Client, context.Context, string, ...parseapi.HLROptions) (*parseapi.HLR, error) = (*parseapi.Client).HLR
var _ func(*parseapi.Client, context.Context, string, ...parseapi.DomainOptions) (*parseapi.Domain, error) = (*parseapi.Client).Domain
var _ func(*parseapi.Client, context.Context, string, ...parseapi.ASNOptions) (*parseapi.ASN, error) = (*parseapi.Client).ASN
var _ func(*parseapi.Client, context.Context, string, ...parseapi.MACOptions) (*parseapi.MAC, error) = (*parseapi.Client).MAC
var _ func(*parseapi.Client, context.Context, string, ...parseapi.MXOptions) (*parseapi.MX, error) = (*parseapi.Client).MX
var _ func(*parseapi.Client, context.Context, string, ...parseapi.UserAgentOptions) (*parseapi.UserAgent, error) = (*parseapi.Client).UserAgent
var _ func(*parseapi.Client, context.Context, string, ...parseapi.VINOptions) (*parseapi.VIN, error) = (*parseapi.Client).VIN
var _ func(*parseapi.Client, context.Context, string, ...parseapi.TariffOptions) (*parseapi.Tariff, error) = (*parseapi.Client).Tariff
var _ func(*parseapi.Client, context.Context, string, ...parseapi.TariffSearchOptions) (*parseapi.TariffSearch, error) = (*parseapi.Client).TariffSearch
var _ func(*parseapi.Client, context.Context, string, ...parseapi.CurrencyOptions) (*parseapi.Currency, error) = (*parseapi.Client).Currency
var _ func(*parseapi.Client, context.Context, string, ...parseapi.LanguageOptions) (*parseapi.Language, error) = (*parseapi.Client).Language
var _ func(*parseapi.Client, context.Context, string, ...parseapi.NameOptions) (*parseapi.Name, error) = (*parseapi.Client).Name
var _ func(*parseapi.Client, context.Context, string, string, ...parseapi.CurrencyRateOptions) (*parseapi.CurrencyRate, error) = (*parseapi.Client).CurrencyRate
var _ func(*parseapi.Client, context.Context, string, ...parseapi.TimezoneOptions) (*parseapi.Timezone, error) = (*parseapi.Client).Timezone
var _ func(*parseapi.Client, context.Context, float64, float64, ...parseapi.TimezoneAtOptions) (*parseapi.Timezone, error) = (*parseapi.Client).TimezoneAt
var _ func(*parseapi.Client, context.Context, string, ...parseapi.DateOptions) (*parseapi.DateInfo, error) = (*parseapi.Client).Date
var _ func(*parseapi.Client, context.Context, ...parseapi.DateTodayOptions) (*parseapi.DateInfo, error) = (*parseapi.Client).DateToday
var _ func(*parseapi.Client, context.Context, string, ...parseapi.HolidayOptions) (*parseapi.HolidayYear, error) = (*parseapi.Client).Holiday
var _ func(*parseapi.Client, context.Context, string, string, ...parseapi.HolidayDateOptions) (*parseapi.HolidayDate, error) = (*parseapi.Client).HolidayDate
var _ func(*parseapi.Client, context.Context, float64, float64, ...parseapi.ElevationOptions) (*parseapi.Elevation, error) = (*parseapi.Client).Elevation
var _ func(*parseapi.Client, context.Context, float64, float64, ...parseapi.PointOptions) (*parseapi.Point, error) = (*parseapi.Client).Point
var _ func(*parseapi.Client, context.Context, float64, float64, ...parseapi.WeatherOptions) (*parseapi.Weather, error) = (*parseapi.Client).Weather
var _ func(*parseapi.Client, context.Context, string, ...parseapi.EmojiOptions) (*parseapi.Emoji, error) = (*parseapi.Client).Emoji
var _ func(*parseapi.Client, context.Context, string, ...parseapi.EmojiSearchOptions) (*parseapi.EmojiSearch, error) = (*parseapi.Client).EmojiSearch
var _ func(*parseapi.Client, context.Context, string, ...parseapi.AddressOptions) (*parseapi.Address, error) = (*parseapi.Client).Address
var _ func(*parseapi.Client, context.Context, string, ...parseapi.AddressSearchOptions) (*parseapi.AddressSearch, error) = (*parseapi.Client).AddressSearch
var _ func(*parseapi.Client, context.Context, string, ...parseapi.CompanyOptions) (*parseapi.Company, error) = (*parseapi.Client).Company
