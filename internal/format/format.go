package format

import (
	"encoding/json"
	"meter-sync/internal/domain"
)

func RecordJSON(record domain.Record) ([]byte, error) { return json.MarshalIndent(record, "", "  ") }
func RecordsJSON(records []domain.Record) ([]byte, error) {
	return json.MarshalIndent(records, "", "  ")
}
func AuditJSON(events []domain.AuditEvent) ([]byte, error) {
	return json.MarshalIndent(events, "", "  ")
}
