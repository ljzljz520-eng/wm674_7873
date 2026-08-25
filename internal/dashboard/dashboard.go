package dashboard

import (
	"fmt"
	"meter-sync/internal/domain"
	"meter-sync/internal/service"
	"sort"
	"strings"
)

type Dashboard struct {
	Summary         service.Summary
	Metrics         []service.ManufacturerMetric
	OpenStatuses    map[domain.Status]int
	TopManufacturer string
	TopAmount       int64
}

func Build(records []domain.Record) Dashboard {
	summary := service.Summarize(records)
	metrics := service.Metrics(records)
	top, _ := service.LargestMetric(metrics)
	return Dashboard{Summary: summary, Metrics: metrics, OpenStatuses: service.StatusCounts(records), TopManufacturer: top.Manufacturer, TopAmount: top.Amount}
}
func (d Dashboard) Headline() string {
	return fmt.Sprintf("%d records / %d CNY / leader %s", d.Summary.Total, d.Summary.Amount, d.TopManufacturer)
}
func (d Dashboard) IsEmpty() bool                   { return d.Summary.Total == 0 }
func (d Dashboard) Status(status domain.Status) int { return d.OpenStatuses[status] }
func (d Dashboard) Metric(manufacturer string) (service.ManufacturerMetric, bool) {
	for _, metric := range d.Metrics {
		if strings.EqualFold(metric.Manufacturer, manufacturer) {
			return metric, true
		}
	}
	return service.ManufacturerMetric{}, false
}
func (d Dashboard) SortedMetrics() []service.ManufacturerMetric {
	metrics := append([]service.ManufacturerMetric(nil), d.Metrics...)
	sort.Slice(metrics, func(i, j int) bool {
		if metrics[i].Amount == metrics[j].Amount {
			return metrics[i].Manufacturer < metrics[j].Manufacturer
		}
		return metrics[i].Amount > metrics[j].Amount
	})
	return metrics
}
func (d Dashboard) ExportRows() [][]string {
	rows := [][]string{{"manufacturer", "count", "amount", "latest_version"}}
	for _, metric := range d.SortedMetrics() {
		rows = append(rows, []string{metric.Manufacturer, fmt.Sprint(metric.Count), fmt.Sprint(metric.Amount), fmt.Sprint(metric.LatestVersion)})
	}
	return rows
}
func (d Dashboard) ExportCSV() string {
	rows := d.ExportRows()
	var builder strings.Builder
	for _, row := range rows {
		builder.WriteString(strings.Join(row, ","))
		builder.WriteByte('\n')
	}
	return builder.String()
}
