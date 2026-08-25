package store

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"meter-sync/internal/domain"
)

type Snapshot struct {
	Records     []domain.Record     `json:"records"`
	Audits      []domain.AuditEvent `json:"audits"`
	Workflows   []domain.Workflow   `json:"workflows"`
	Attachments []domain.Attachment `json:"attachments"`
}

func (s *Store) Snapshot() (Snapshot, error) {
	snapshot := Snapshot{Records: make([]domain.Record, 0), Audits: make([]domain.AuditEvent, 0), Workflows: make([]domain.Workflow, 0), Attachments: make([]domain.Attachment, 0)}
	s.mu.RLock()
	defer s.mu.RUnlock()
	err := s.db.View(func(tx *bbolt.Tx) error {
		if err := readRecords(tx, snapshotPointer(&snapshot)); err != nil {
			return err
		}
		if err := readAudits(tx, &snapshot); err != nil {
			return err
		}
		if err := readWorkflows(tx, &snapshot); err != nil {
			return err
		}
		return readAttachments(tx, &snapshot)
	})
	return snapshot, err
}
func snapshotPointer(snapshot *Snapshot) *Snapshot { return snapshot }
func readRecords(tx *bbolt.Tx, snapshot *Snapshot) error {
	return tx.Bucket(bucketRecords).ForEach(func(_, value []byte) error {
		var record domain.Record
		if err := json.Unmarshal(value, &record); err != nil {
			return err
		}
		snapshot.Records = append(snapshot.Records, record)
		return nil
	})
}
func readAudits(tx *bbolt.Tx, snapshot *Snapshot) error {
	return tx.Bucket(bucketAudits).ForEach(func(_, value []byte) error {
		var event domain.AuditEvent
		if err := json.Unmarshal(value, &event); err != nil {
			return err
		}
		snapshot.Audits = append(snapshot.Audits, event)
		return nil
	})
}
func readWorkflows(tx *bbolt.Tx, snapshot *Snapshot) error {
	return tx.Bucket(bucketFlows).ForEach(func(_, value []byte) error {
		var flow domain.Workflow
		if err := json.Unmarshal(value, &flow); err != nil {
			return err
		}
		snapshot.Workflows = append(snapshot.Workflows, flow)
		return nil
	})
}
func readAttachments(tx *bbolt.Tx, snapshot *Snapshot) error {
	return tx.Bucket(bucketAttachments).ForEach(func(_, value []byte) error {
		var attachment domain.Attachment
		if err := json.Unmarshal(value, &attachment); err != nil {
			return err
		}
		snapshot.Attachments = append(snapshot.Attachments, attachment)
		return nil
	})
}
func EncodeSnapshot(snapshot Snapshot) ([]byte, error) { return json.MarshalIndent(snapshot, "", "  ") }
func SnapshotSummary(snapshot Snapshot) string {
	return fmt.Sprintf("records=%d audits=%d workflows=%d attachments=%d", len(snapshot.Records), len(snapshot.Audits), len(snapshot.Workflows), len(snapshot.Attachments))
}
