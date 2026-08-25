package ledger

import (
	"errors"
	"fmt"
	"meter-sync/internal/domain"
	"sort"
	"strings"
)

type Entry struct {
	ID         string
	RecordID   string
	Kind       string
	Amount     int64
	Currency   string
	Reference  string
	CreatedAt  string
	Reconciled bool
}
type Ledger struct{ entries map[string]Entry }

func New() *Ledger { return &Ledger{entries: make(map[string]Entry)} }
func (l *Ledger) Add(record domain.Record, kind, stamp string) (Entry, error) {
	if record.ID == "" {
		return Entry{}, errors.New("record id is required")
	}
	if kind == "" {
		return Entry{}, errors.New("entry kind is required")
	}
	id := fmt.Sprintf("%s:%s:%d", record.ID, kind, record.Version)
	if _, ok := l.entries[id]; ok {
		return Entry{}, fmt.Errorf("entry %s already exists", id)
	}
	entry := Entry{ID: id, RecordID: record.ID, Kind: kind, Amount: record.Amount, Currency: record.Currency, Reference: record.SyncReference, CreatedAt: stamp}
	l.entries[id] = entry
	return entry, nil
}
func (l *Ledger) Get(id string) (Entry, bool) { entry, ok := l.entries[id]; return entry, ok }
func (l *Ledger) MarkReconciled(id string) error {
	entry, ok := l.entries[id]
	if !ok {
		return fmt.Errorf("entry %s not found", id)
	}
	entry.Reconciled = true
	l.entries[id] = entry
	return nil
}
func (l *Ledger) Entries(recordID string) []Entry {
	result := make([]Entry, 0)
	for _, entry := range l.entries {
		if recordID == "" || entry.RecordID == recordID {
			result = append(result, entry)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}
func (l *Ledger) Total(recordID string) int64 {
	var total int64
	for _, entry := range l.Entries(recordID) {
		total += entry.Amount
	}
	return total
}
func (l *Ledger) Unreconciled() []Entry {
	result := make([]Entry, 0)
	for _, entry := range l.entries {
		if !entry.Reconciled {
			result = append(result, entry)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}
func (l *Ledger) ReconcileAll() int {
	count := 0
	for _, entry := range l.Unreconciled() {
		if l.MarkReconciled(entry.ID) == nil {
			count++
		}
	}
	return count
}
func (l *Ledger) CurrencyTotals() map[string]int64 {
	totals := make(map[string]int64)
	for _, entry := range l.entries {
		totals[strings.ToUpper(entry.Currency)] += entry.Amount
	}
	return totals
}
