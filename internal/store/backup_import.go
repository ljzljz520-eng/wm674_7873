package store

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"meter-sync/internal/domain"
)

func (s *Store) Restore(snapshot Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Update(func(tx *bbolt.Tx) error {
		if err := clearBucket(tx, bucketRecords); err != nil {
			return err
		}
		if err := clearBucket(tx, bucketAudits); err != nil {
			return err
		}
		if err := clearBucket(tx, bucketFlows); err != nil {
			return err
		}
		if err := clearBucket(tx, bucketAttachments); err != nil {
			return err
		}
		for _, record := range snapshot.Records {
			if err := put(tx, bucketRecords, record.ID, record); err != nil {
				return err
			}
		}
		for _, event := range snapshot.Audits {
			if err := put(tx, bucketAudits, event.ID, event); err != nil {
				return err
			}
		}
		for _, flow := range snapshot.Workflows {
			if err := put(tx, bucketFlows, flow.ID, flow); err != nil {
				return err
			}
		}
		for _, attachment := range snapshot.Attachments {
			if err := put(tx, bucketAttachments, attachment.ID, attachment); err != nil {
				return err
			}
		}
		return nil
	})
}
func clearBucket(tx *bbolt.Tx, name []byte) error {
	if err := tx.DeleteBucket(name); err != nil && err != bbolt.ErrBucketNotFound {
		return err
	}
	_, err := tx.CreateBucket(name)
	return err
}
func DecodeSnapshot(data []byte) (Snapshot, error) {
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return Snapshot{}, fmt.Errorf("decode snapshot: %w", err)
	}
	return snapshot, nil
}
func ValidateSnapshot(snapshot Snapshot) error {
	seen := make(map[string]struct{})
	for _, record := range snapshot.Records {
		if err := record.Validate(); err != nil {
			return err
		}
		if _, ok := seen[record.ID]; ok {
			return fmt.Errorf("duplicate record %s", record.ID)
		}
		seen[record.ID] = struct{}{}
	}
	for _, event := range snapshot.Audits {
		if event.ID == "" || event.RecordID == "" {
			return fmt.Errorf("invalid audit event")
		}
	}
	return nil
}
func MergeSnapshots(left, right Snapshot) Snapshot {
	merged := left
	index := make(map[string]int)
	for i, record := range merged.Records {
		index[record.ID] = i
	}
	for _, record := range right.Records {
		if i, ok := index[record.ID]; ok {
			if record.Version > merged.Records[i].Version {
				merged.Records[i] = record
			}
		} else {
			index[record.ID] = len(merged.Records)
			merged.Records = append(merged.Records, record)
		}
	}
	merged.Audits = append(merged.Audits, right.Audits...)
	merged.Workflows = append(merged.Workflows, right.Workflows...)
	merged.Attachments = append(merged.Attachments, right.Attachments...)
	return merged
}

var _ = domain.StatusDraft
