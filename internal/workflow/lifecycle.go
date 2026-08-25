package workflow

import (
	"fmt"
	"meter-sync/internal/domain"
	"meter-sync/internal/service"
)

type Lifecycle struct{ Service *service.Service }

func (l Lifecycle) CreateReviewPublishArchive(id, manufacturer, model, reference, actor string, amount int64) (domain.Record, error) {
	record, err := l.Service.Register(id, manufacturer, model, reference, actor, amount)
	if err != nil {
		return domain.Record{}, err
	}
	if record.Status != domain.StatusDraft {
		return domain.Record{}, fmt.Errorf("unexpected initial state")
	}
	record, err = l.Service.Review(id, actor, "validated manufacturer payload")
	if err != nil {
		return domain.Record{}, err
	}
	record, err = l.Service.Publish(id, actor)
	if err != nil {
		return domain.Record{}, err
	}
	return l.Service.Archive(id, actor)
}

func (l Lifecycle) ReopenAndDescribe(id string) (string, error) {
	record, err := l.Service.Search(domain.Query{Reference: id})
	if err != nil {
		return "", err
	}
	if len(record) == 0 {
		return "", fmt.Errorf("no records for %s", id)
	}
	return fmt.Sprintf("%s:%s:%d", record[0].ID, record[0].Status, record[0].Amount), nil
}
