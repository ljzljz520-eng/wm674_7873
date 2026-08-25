package catalog

import (
	"fmt"
	"strings"
)

type Registry struct{ profiles map[string]Profile }

func NewRegistry() *Registry {
	registry := &Registry{profiles: make(map[string]Profile)}
	for _, profile := range All() {
		registry.profiles[profile.Code] = profile
	}
	return registry
}
func (r *Registry) Register(profile Profile) error {
	if strings.TrimSpace(profile.Code) == "" {
		return fmt.Errorf("profile code is required")
	}
	if _, exists := r.profiles[profile.Code]; exists {
		return fmt.Errorf("profile %s already registered", profile.Code)
	}
	r.profiles[profile.Code] = profile
	return nil
}
func (r *Registry) Get(code string) (Profile, bool) {
	profile, ok := r.profiles[strings.ToUpper(strings.TrimSpace(code))]
	return profile, ok
}
func (r *Registry) Codes() []string {
	codes := make([]string, 0, len(r.profiles))
	for code := range r.profiles {
		codes = append(codes, code)
	}
	sortStrings(codes)
	return codes
}
func sortStrings(values []string) {
	for i := 1; i < len(values); i++ {
		for j := i; j > 0 && values[j] < values[j-1]; j-- {
			values[j], values[j-1] = values[j-1], values[j]
		}
	}
}
func (r *Registry) ActiveProfiles() []Profile {
	result := make([]Profile, 0)
	for _, profile := range r.profiles {
		if profile.Active {
			result = append(result, profile)
		}
	}
	sortProfiles(result)
	return result
}
func sortProfiles(values []Profile) {
	for i := 1; i < len(values); i++ {
		for j := i; j > 0 && values[j].Code < values[j-1].Code; j-- {
			values[j], values[j-1] = values[j-1], values[j]
		}
	}
}
