package catalog

import (
	"fmt"
	"strings"
)

type Compatibility struct {
	Manufacturer string
	Model        string
	Currency     string
	Compatible   bool
	Reason       string
}

func CheckCompatibility(manufacturer, model, currency string) Compatibility {
	profile, err := Find(manufacturer)
	if err != nil {
		return Compatibility{Manufacturer: manufacturer, Model: model, Currency: currency, Reason: err.Error()}
	}
	if !profile.Active {
		return Compatibility{Manufacturer: profile.Name, Model: model, Currency: currency, Reason: "profile is inactive"}
	}
	if !SupportsModel(profile, model) {
		return Compatibility{Manufacturer: profile.Name, Model: model, Currency: currency, Reason: "model is unavailable"}
	}
	if !SupportsCurrency(profile, currency) {
		return Compatibility{Manufacturer: profile.Name, Model: model, Currency: currency, Reason: "currency is unavailable"}
	}
	return Compatibility{Manufacturer: profile.Name, Model: model, Currency: currency, Compatible: true, Reason: "compatible"}
}
func CompatibilityLabel(result Compatibility) string {
	if result.Compatible {
		return fmt.Sprintf("%s/%s %s", result.Manufacturer, result.Model, result.Currency)
	}
	return strings.TrimSpace(result.Reason)
}
func FilterCompatible(items []Compatibility) []Compatibility {
	result := make([]Compatibility, 0)
	for _, item := range items {
		if item.Compatible {
			result = append(result, item)
		}
	}
	return result
}
func GroupByCurrency(items []Compatibility) map[string][]Compatibility {
	groups := make(map[string][]Compatibility)
	for _, item := range items {
		if item.Currency == "" {
			continue
		}
		groups[item.Currency] = append(groups[item.Currency], item)
	}
	return groups
}
