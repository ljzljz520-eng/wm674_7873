package catalog

import (
	"errors"
	"sort"
	"strings"
)

type Profile struct {
	Code       string
	Name       string
	Region     string
	Models     []string
	Currencies []string
	Active     bool
}

var profiles = []Profile{
	{Code: "HD", Name: "Huadian Meter Systems", Region: "East", Models: []string{"HX-100", "HX-200", "HX-300"}, Currencies: []string{"CNY"}, Active: true},
	{Code: "NARI", Name: "Nari Technology Meters", Region: "East", Models: []string{"NR-01", "NR-02", "NR-03"}, Currencies: []string{"CNY"}, Active: true},
	{Code: "WASION", Name: "Wasion Group", Region: "Central", Models: []string{"WS-10", "WS-20", "WS-30"}, Currencies: []string{"CNY", "USD"}, Active: true},
	{Code: "LANDIS", Name: "Landis+Gyr China", Region: "North", Models: []string{"LG-100", "LG-200"}, Currencies: []string{"CNY"}, Active: true},
	{Code: "SIEMENS", Name: "Siemens Metering", Region: "North", Models: []string{"S7-120", "S7-220"}, Currencies: []string{"CNY", "EUR"}, Active: true},
	{Code: "CHINT", Name: "Chint Instrument", Region: "East", Models: []string{"DTS-100", "DTS-200"}, Currencies: []string{"CNY"}, Active: true},
	{Code: "ABB", Name: "ABB Measurement", Region: "South", Models: []string{"A1", "A2"}, Currencies: []string{"CNY", "USD"}, Active: true},
	{Code: "SCHNEIDER", Name: "Schneider Electric Meters", Region: "South", Models: []string{"SE-1", "SE-2"}, Currencies: []string{"CNY", "EUR"}, Active: true},
}

func All() []Profile {
	result := make([]Profile, len(profiles))
	copy(result, profiles)
	for i := range result {
		result[i].Models = append([]string(nil), result[i].Models...)
		result[i].Currencies = append([]string(nil), result[i].Currencies...)
	}
	return result
}
func Find(code string) (Profile, error) {
	for _, profile := range profiles {
		if strings.EqualFold(profile.Code, strings.TrimSpace(code)) {
			return profile, nil
		}
	}
	return Profile{}, errors.New("manufacturer profile not found")
}
func FindByName(name string) []Profile {
	result := make([]Profile, 0)
	needle := strings.ToLower(strings.TrimSpace(name))
	for _, profile := range profiles {
		if strings.Contains(strings.ToLower(profile.Name), needle) {
			result = append(result, profile)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Code < result[j].Code })
	return result
}
func SupportsModel(profile Profile, model string) bool {
	for _, candidate := range profile.Models {
		if candidate == model {
			return true
		}
	}
	return false
}
func SupportsCurrency(profile Profile, currency string) bool {
	for _, candidate := range profile.Currencies {
		if strings.EqualFold(candidate, currency) {
			return true
		}
	}
	return false
}
func Validate(profile Profile, model, currency string) error {
	if !profile.Active {
		return errors.New("manufacturer profile is inactive")
	}
	if !SupportsModel(profile, model) {
		return errors.New("model is not registered for manufacturer")
	}
	if !SupportsCurrency(profile, currency) {
		return errors.New("currency is not supported")
	}
	return nil
}
