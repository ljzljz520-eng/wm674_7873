package policy

import (
	"errors"
	"meter-sync/internal/domain"
)

type Rule struct {
	Name    string
	Check   func(domain.Record) bool
	Failure string
}
type Engine struct{ Rules []Rule }

func Default() Engine {
	return Engine{Rules: []Rule{{"positive-version", func(r domain.Record) bool { return r.Version > 0 }, "version must be positive"}, {"known-currency", func(r domain.Record) bool { return r.Currency == "CNY" }, "currency must be CNY"}, {"reference-present", func(r domain.Record) bool { return r.SyncReference != "" }, "reference is required"}}}
}
func (e Engine) Check(record domain.Record) error {
	for _, rule := range e.Rules {
		if rule.Check == nil {
			continue
		}
		if !rule.Check(record) {
			return errors.New(rule.Failure)
		}
	}
	return nil
}
func (e Engine) CheckBatch(records []domain.Record) map[string]error {
	failures := make(map[string]error)
	for _, record := range records {
		if err := e.Check(record); err != nil {
			failures[record.ID] = err
		}
	}
	return failures
}
