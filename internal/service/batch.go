package service

import (
	"errors"
	"fmt"
	"meter-sync/internal/domain"
)

type BatchOutcome struct {
	Row    int
	Record domain.Record
	Error  string
}
type BatchSummary struct {
	Outcomes    []BatchOutcome
	Accepted    int
	Rejected    int
	TotalAmount int64
}

func (s *Service) ProcessBatch(rows []domain.ImportRow, actor string) BatchSummary {
	summary := BatchSummary{Outcomes: make([]BatchOutcome, 0, len(rows))}
	for index, row := range rows {
		record, err := s.processRow(row, actor)
		outcome := BatchOutcome{Row: index + 1, Record: record}
		if err != nil {
			outcome.Error = err.Error()
			summary.Rejected++
		} else {
			summary.Accepted++
			summary.TotalAmount += record.Amount
		}
		summary.Outcomes = append(summary.Outcomes, outcome)
	}
	return summary
}
func (s *Service) processRow(row domain.ImportRow, actor string) (domain.Record, error) {
	if row.ID == "" {
		return domain.Record{}, errors.New("row identifier is required")
	}
	if row.Amount < 0 {
		return domain.Record{}, errors.New("row amount cannot be negative")
	}
	record, err := s.store.GetRecord(row.ID)
	if err != nil {
		return s.Register(row.ID, row.Manufacturer, row.Model, row.Reference, actor, row.Amount)
	}
	if record.SyncReference != row.Reference {
		return domain.Record{}, fmt.Errorf("reference mismatch for %s", row.ID)
	}
	return s.UpdateAmount(row.ID, row.Amount, actor)
}
func (s *Service) ValidateBatch(rows []domain.ImportRow) []domain.ImportError {
	errorsFound := make([]domain.ImportError, 0)
	seen := make(map[string]bool)
	for index, row := range rows {
		if row.ID == "" {
			errorsFound = append(errorsFound, domain.ImportError{Row: index + 1, Cause: "identifier is required"})
			continue
		}
		if seen[row.ID] {
			errorsFound = append(errorsFound, domain.ImportError{Row: index + 1, Cause: "duplicate identifier"})
		}
		seen[row.ID] = true
		if row.Amount < 0 {
			errorsFound = append(errorsFound, domain.ImportError{Row: index + 1, Cause: "amount cannot be negative"})
		}
		if row.Manufacturer == "" {
			errorsFound = append(errorsFound, domain.ImportError{Row: index + 1, Cause: "manufacturer is required"})
		}
		if row.Model == "" {
			errorsFound = append(errorsFound, domain.ImportError{Row: index + 1, Cause: "model is required"})
		}
		if row.Reference == "" {
			errorsFound = append(errorsFound, domain.ImportError{Row: index + 1, Cause: "reference is required"})
		}
	}
	return errorsFound
}
func (s *Service) DryRun(rows []domain.ImportRow) BatchSummary {
	summary := BatchSummary{Outcomes: make([]BatchOutcome, 0, len(rows))}
	issues := s.ValidateBatch(rows)
	for index, row := range rows {
		outcome := BatchOutcome{Row: index + 1}
		for _, issue := range issues {
			if issue.Row == index+1 {
				outcome.Error = outcome.Error + issue.Cause + "; "
			}
		}
		if outcome.Error == "" {
			outcome.Record = domain.Record{ID: row.ID, Manufacturer: row.Manufacturer, MeterModel: row.Model, SyncReference: row.Reference, Amount: row.Amount}
			summary.Accepted++
			summary.TotalAmount += row.Amount
		} else {
			summary.Rejected++
		}
		summary.Outcomes = append(summary.Outcomes, outcome)
	}
	return summary
}
func (s *Service) ReplaceReference(id, reference, actor string) (domain.Record, error) {
	if reference == "" {
		return domain.Record{}, errors.New("reference is required")
	}
	record, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	record.SyncReference = reference
	record.Version++
	record.UpdatedAt = s.clock.Now()
	if err := s.store.SaveRecord(record); err != nil {
		return domain.Record{}, err
	}
	if err := s.audit(record, "reference_update", actor, reference); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}
func (s *Service) ReviewWithPolicy(id, actor, note string, check func(domain.Record) error) (domain.Record, error) {
	record, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	if check != nil {
		if err := check(record); err != nil {
			return domain.Record{}, err
		}
	}
	return s.Review(id, actor, note)
}
