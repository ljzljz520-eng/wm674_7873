package dashboard

import (
	"meter-sync/internal/domain"
	"sort"
	"strings"
)

type Filter struct {
	Manufacturer  string
	Status        domain.Status
	MinimumAmount int64
	MaximumAmount int64
}

func (f Filter) Matches(record domain.Record) bool {
	if f.Manufacturer != "" && !strings.Contains(strings.ToLower(record.Manufacturer), strings.ToLower(f.Manufacturer)) {
		return false
	}
	if f.Status != "" && record.Status != f.Status {
		return false
	}
	if record.Amount < f.MinimumAmount {
		return false
	}
	if f.MaximumAmount > 0 && record.Amount > f.MaximumAmount {
		return false
	}
	return true
}
func Apply(records []domain.Record, filter Filter) []domain.Record {
	result := make([]domain.Record, 0)
	for _, record := range records {
		if filter.Matches(record) {
			result = append(result, record)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}
func GroupByManufacturer(records []domain.Record) map[string][]domain.Record {
	groups := make(map[string][]domain.Record)
	for _, record := range records {
		groups[record.Manufacturer] = append(groups[record.Manufacturer], record)
	}
	for _, items := range groups {
		sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	}
	return groups
}
func GroupAmount(records []domain.Record) map[string]int64 {
	groups := make(map[string]int64)
	for _, record := range records {
		groups[record.Manufacturer] += record.Amount
	}
	return groups
}
func IDs(records []domain.Record) []string {
	ids := make([]string, 0, len(records))
	for _, record := range records {
		ids = append(ids, record.ID)
	}
	sort.Strings(ids)
	return ids
}
func Latest(records []domain.Record) (domain.Record, bool) {
	if len(records) == 0 {
		return domain.Record{}, false
	}
	latest := records[0]
	for _, record := range records[1:] {
		if record.Version > latest.Version || (record.Version == latest.Version && record.UpdatedAt > latest.UpdatedAt) {
			latest = record
		}
	}
	return latest, true
}
func Count(records []domain.Record, filter Filter) int { return len(Apply(records, filter)) }
func Amount(records []domain.Record) int64 {
	var total int64
	for _, record := range records {
		total += record.Amount
	}
	return total
}
func Statuses(records []domain.Record) []domain.Status {
	set := make(map[domain.Status]struct{})
	for _, record := range records {
		set[record.Status] = struct{}{}
	}
	statuses := make([]domain.Status, 0, len(set))
	for status := range set {
		statuses = append(statuses, status)
	}
	sort.Slice(statuses, func(i, j int) bool { return statuses[i] < statuses[j] })
	return statuses
}
func Empty(records []domain.Record) bool { return len(records) == 0 }
