package domain

import (
	"errors"
	"fmt"
	"strings"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusReviewed  Status = "reviewed"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

type Record struct {
	ID            string `json:"id"`
	Manufacturer  string `json:"manufacturer"`
	MeterModel    string `json:"meter_model"`
	SyncReference string `json:"sync_reference"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	Status        Status `json:"status"`
	Version       int    `json:"version"`
	CreatedBy     string `json:"created_by"`
	ReviewedBy    string `json:"reviewed_by"`
	ReviewNote    string `json:"review_note"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	ArchivedAt    string `json:"archived_at"`
}

func NewRecord(id, manufacturer, model, reference, actor, stamp string, amount int64) (Record, error) {
	record := Record{ID: strings.TrimSpace(id), Manufacturer: strings.TrimSpace(manufacturer), MeterModel: strings.TrimSpace(model), SyncReference: strings.TrimSpace(reference), Amount: amount, Currency: "CNY", Status: StatusDraft, Version: 1, CreatedBy: strings.TrimSpace(actor), CreatedAt: stamp, UpdatedAt: stamp}
	if err := record.Validate(); err != nil {
		return Record{}, err
	}
	return record, nil
}

func (r Record) Validate() error {
	if r.ID == "" || r.Manufacturer == "" || r.MeterModel == "" || r.SyncReference == "" {
		return errors.New("record identity fields are required")
	}
	if r.Amount < 0 {
		return errors.New("amount cannot be negative")
	}
	if r.Currency == "" {
		return errors.New("currency is required")
	}
	if r.Version < 1 {
		return errors.New("version must be positive")
	}
	if r.Status != StatusDraft && r.Status != StatusReviewed && r.Status != StatusPublished && r.Status != StatusArchived {
		return fmt.Errorf("unsupported status %q", r.Status)
	}
	return nil
}

func (r Record) CanReview() bool  { return r.Status == StatusDraft }
func (r Record) CanPublish() bool { return r.Status == StatusReviewed }
func (r Record) CanArchive() bool { return r.Status == StatusPublished }

func (r Record) WithAmount(amount int64, actor, stamp string) (Record, error) {
	if amount < 0 {
		return Record{}, errors.New("amount cannot be negative")
	}
	if strings.TrimSpace(actor) == "" {
		return Record{}, errors.New("actor is required")
	}
	r.Amount = amount
	r.Version++
	r.UpdatedAt = stamp
	return r, nil
}

func (r Record) MarkReviewed(actor, note, stamp string) (Record, error) {
	if !r.CanReview() {
		return Record{}, errors.New("record is not reviewable")
	}
	if strings.TrimSpace(actor) == "" {
		return Record{}, errors.New("reviewer is required")
	}
	r.Status = StatusReviewed
	r.ReviewedBy = strings.TrimSpace(actor)
	r.ReviewNote = strings.TrimSpace(note)
	r.Version++
	r.UpdatedAt = stamp
	return r, nil
}

func (r Record) MarkPublished(stamp string) (Record, error) {
	if !r.CanPublish() {
		return Record{}, errors.New("record is not publishable")
	}
	r.Status = StatusPublished
	r.Version++
	r.UpdatedAt = stamp
	return r, nil
}

func (r Record) MarkArchived(stamp string) (Record, error) {
	if !r.CanArchive() {
		return Record{}, errors.New("record is not archivable")
	}
	r.Status = StatusArchived
	r.ArchivedAt = stamp
	r.Version++
	r.UpdatedAt = stamp
	return r, nil
}
