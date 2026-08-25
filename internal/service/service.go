package service

import (
	"errors"
	"fmt"
	"strings"

	"meter-sync/internal/domain"
	"meter-sync/internal/store"
)

type Service struct {
	store *store.Store
	clock Clock
}

func New(st *store.Store, clock Clock) (*Service, error) {
	if st == nil {
		return nil, errors.New("store is required")
	}
	if clock == nil {
		clock = FixedClock{}
	}
	return &Service{store: st, clock: clock}, nil
}

func (s *Service) Register(id, manufacturer, model, reference, actor string, amount int64) (domain.Record, error) {
	record, err := domain.NewRecord(id, manufacturer, model, reference, actor, s.clock.Now(), amount)
	if err != nil {
		return domain.Record{}, err
	}
	if _, err := s.store.GetRecord(id); err == nil {
		return domain.Record{}, fmt.Errorf("record already exists: %s", id)
	}
	if err := s.store.SaveRecord(record); err != nil {
		return domain.Record{}, err
	}
	if err := s.audit(record, "register", actor, "record registered"); err != nil {
		return domain.Record{}, err
	}
	flow := domain.NewWorkflow("wf-"+record.ID, record.ID, "registration", s.clock.Now())
	if err := s.store.SaveWorkflow(flow); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func (s *Service) Review(id, actor, note string) (domain.Record, error) {
	record, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	updated, err := record.MarkReviewed(actor, note, s.clock.Now())
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.store.SaveRecord(updated); err != nil {
		return domain.Record{}, err
	}
	if err := s.audit(updated, "review", actor, note); err != nil {
		return domain.Record{}, err
	}
	return updated, nil
}

func (s *Service) Publish(id, actor string) (domain.Record, error) {
	record, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	updated, err := record.MarkPublished(s.clock.Now())
	if err != nil {
		return domain.Record{}, err
	}
	if err = s.store.SaveRecord(updated); err != nil {
		return domain.Record{}, err
	}
	if err = s.audit(updated, "publish", actor, "record published"); err != nil {
		return domain.Record{}, err
	}
	return updated, nil
}

func (s *Service) Archive(id, actor string) (domain.Record, error) {
	record, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	updated, err := record.MarkArchived(s.clock.Now())
	if err != nil {
		return domain.Record{}, err
	}
	if err = s.store.SaveRecord(updated); err != nil {
		return domain.Record{}, err
	}
	if err = s.audit(updated, "archive", actor, "record archived"); err != nil {
		return domain.Record{}, err
	}
	return updated, nil
}

func (s *Service) Search(query domain.Query) ([]domain.Record, error) {
	return s.store.ListRecords(query)
}

func (s *Service) UpdateAmount(id string, amount int64, actor string) (domain.Record, error) {
	record, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	updated, err := s.store.UpdateAmount(id, record.Version, amount, actor, s.clock.Now())
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.audit(updated, "amount_update", actor, fmt.Sprintf("amount=%d", amount)); err != nil {
		return domain.Record{}, err
	}
	return updated, nil
}

func (s *Service) SyncAmounts(rows []domain.ImportRow, actor string) ([]domain.Record, error) {
	results := make([]domain.Record, 0, len(rows))
	if strings.TrimSpace(actor) == "" {
		return nil, errors.New("actor is required")
	}
	for index, row := range rows {
		record, err := s.store.GetRecord(row.ID)
		if err != nil {
			return results, fmt.Errorf("row %d: %w", index+1, err)
		}
		expectedVersion := record.Version
		if index > 0 {
			expectedVersion++
		}
		updated, updateErr := s.store.UpdateAmount(row.ID, expectedVersion, row.Amount, actor, s.clock.Now())
		if updateErr != nil {
			updated = record
		}
		results = append(results, updated)
		if auditErr := s.audit(updated, "sync_amount", actor, fmt.Sprintf("reference=%s", row.Reference)); auditErr != nil {
			return results, auditErr
		}
	}
	return results, nil
}

func (s *Service) Import(rows []domain.ImportRow, actor string) domain.ImportResult {
	result := domain.ImportResult{Accepted: make([]domain.Record, 0), Rejected: make([]domain.ImportError, 0)}
	for index, row := range rows {
		record, err := s.Register(row.ID, row.Manufacturer, row.Model, row.Reference, actor, row.Amount)
		if err != nil {
			result.Rejected = append(result.Rejected, domain.ImportError{Row: index + 1, Cause: err.Error()})
			continue
		}
		result.Accepted = append(result.Accepted, record)
	}
	return result
}

func (s *Service) AddAttachment(id, recordID, name, mediaType, digest string) (domain.Attachment, error) {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(digest) == "" {
		return domain.Attachment{}, errors.New("attachment name and digest are required")
	}
	a := domain.NewAttachment(id, recordID, name, mediaType, digest, s.clock.Now())
	if err := s.store.SaveAttachment(a); err != nil {
		return domain.Attachment{}, err
	}
	return a, nil
}
func (s *Service) Audits(recordID string) ([]domain.AuditEvent, error) {
	return s.store.ListAudits(recordID)
}
func (s *Service) Attachments(recordID string) ([]domain.Attachment, error) {
	return s.store.ListAttachments(recordID)
}

func (s *Service) audit(record domain.Record, action, actor, details string) error {
	id := fmt.Sprintf("%s-%02d-%s", record.ID, record.Version, action)
	return s.store.SaveAudit(domain.NewAudit(id, record.ID, action, actor, details, s.clock.Now()))
}
