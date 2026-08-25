package service

import (
	"meter-sync/internal/domain"
	"sort"
)

type ManufacturerMetric struct {
	Manufacturer  string
	Count         int
	Amount        int64
	LatestVersion int
}

func Metrics(records []domain.Record) []ManufacturerMetric {
	grouped := make(map[string]ManufacturerMetric)
	for _, record := range records {
		metric := grouped[record.Manufacturer]
		metric.Manufacturer = record.Manufacturer
		metric.Count++
		metric.Amount += record.Amount
		if record.Version > metric.LatestVersion {
			metric.LatestVersion = record.Version
		}
		grouped[record.Manufacturer] = metric
	}
	result := make([]ManufacturerMetric, 0, len(grouped))
	for _, metric := range grouped {
		result = append(result, metric)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Manufacturer < result[j].Manufacturer })
	return result
}
func LargestMetric(metrics []ManufacturerMetric) (ManufacturerMetric, bool) {
	if len(metrics) == 0 {
		return ManufacturerMetric{}, false
	}
	largest := metrics[0]
	for _, metric := range metrics[1:] {
		if metric.Amount > largest.Amount {
			largest = metric
		} else if metric.Amount == largest.Amount && metric.Manufacturer < largest.Manufacturer {
			largest = metric
		}
	}
	return largest, true
}
func StatusCounts(records []domain.Record) map[domain.Status]int {
	counts := make(map[domain.Status]int)
	for _, record := range records {
		counts[record.Status]++
	}
	return counts
}
func AmountByStatus(records []domain.Record) map[domain.Status]int64 {
	amounts := make(map[domain.Status]int64)
	for _, record := range records {
		amounts[record.Status] += record.Amount
	}
	return amounts
}
