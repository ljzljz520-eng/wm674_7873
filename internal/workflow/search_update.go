package workflow

import (
	"fmt"
	"meter-sync/internal/domain"
	"meter-sync/internal/service"
)

type SearchUpdate struct{ Service *service.Service }

func (w SearchUpdate) FindAndUpdate(manufacturer string, amount int64, actor string) (domain.Record, error) {
	records, err := w.Service.Search(domain.Query{Manufacturer: manufacturer, Limit: 1})
	if err != nil {
		return domain.Record{}, err
	}
	if len(records) == 0 {
		return domain.Record{}, fmt.Errorf("manufacturer not found: %s", manufacturer)
	}
	return w.Service.UpdateAmount(records[0].ID, amount, actor)
}
func (w SearchUpdate) PublishUpdated(id, actor string) (domain.Record, error) {
	record, err := w.Service.Review(id, actor, "change review")
	if err != nil {
		return domain.Record{}, err
	}
	return w.Service.Publish(record.ID, actor)
}
