package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"sync"

	"go.etcd.io/bbolt"
	"meter-sync/internal/domain"
)

var (
	bucketRecords     = []byte("records")
	bucketAudits      = []byte("audits")
	bucketFlows       = []byte("workflows")
	bucketAttachments = []byte("attachments")
)

type Store struct {
	db *bbolt.DB
	mu sync.RWMutex
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	db, err := bbolt.Open(filepath.Clean(path), 0o600, nil)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, name := range [][]byte{bucketRecords, bucketAudits, bucketFlows, bucketAttachments} {
			if _, e := tx.CreateBucketIfNotExists(name); e != nil {
				return e
			}
		}
		return nil
	})
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("initialize database: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func encode(value any) ([]byte, error)     { return json.Marshal(value) }
func decode(data []byte, target any) error { return json.Unmarshal(data, target) }

func put(tx *bbolt.Tx, bucket []byte, key string, value any) error {
	data, err := encode(value)
	if err != nil {
		return err
	}
	return tx.Bucket(bucket).Put([]byte(key), data)
}

func get(bucket *bbolt.Bucket, key string, target any) error {
	data := bucket.Get([]byte(key))
	if data == nil {
		return fmt.Errorf("not found: %s", key)
	}
	copyData := append([]byte(nil), data...)
	return decode(copyData, target)
}

func (s *Store) SaveRecord(record domain.Record) error {
	if err := record.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Update(func(tx *bbolt.Tx) error { return put(tx, bucketRecords, record.ID, record) })
}

func (s *Store) GetRecord(id string) (domain.Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var record domain.Record
	err := s.db.View(func(tx *bbolt.Tx) error { return get(tx.Bucket(bucketRecords), id, &record) })
	return record, err
}

func (s *Store) ListRecords(query domain.Query) ([]domain.Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	records := make([]domain.Record, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketRecords).ForEach(func(_, value []byte) error {
			var record domain.Record
			if err := decode(value, &record); err != nil {
				return err
			}
			if query.Matches(record) {
				records = append(records, record)
			}
			return nil
		})
	})
	sort.Slice(records, func(i, j int) bool { return records[i].ID < records[j].ID })
	if query.Limit > 0 && len(records) > query.Limit {
		records = records[:query.Limit]
	}
	return records, err
}

func (s *Store) UpdateAmount(id string, expectedVersion int, amount int64, actor, stamp string) (domain.Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var updated domain.Record
	err := s.db.Update(func(tx *bbolt.Tx) error {
		var current domain.Record
		if err := get(tx.Bucket(bucketRecords), id, &current); err != nil {
			return err
		}
		if current.Version != expectedVersion {
			return fmt.Errorf("version conflict: expected %d got %d", expectedVersion, current.Version)
		}
		candidate, err := current.WithAmount(amount, actor, stamp)
		if err != nil {
			return err
		}
		updated = candidate
		return put(tx, bucketRecords, id, candidate)
	})
	return updated, err
}

func (s *Store) SaveAudit(event domain.AuditEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Update(func(tx *bbolt.Tx) error { return put(tx, bucketAudits, event.ID, event) })
}
func (s *Store) ListAudits(recordID string) ([]domain.AuditEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]domain.AuditEvent, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketAudits).ForEach(func(_, value []byte) error {
			var event domain.AuditEvent
			if e := decode(value, &event); e != nil {
				return e
			}
			if recordID == "" || event.RecordID == recordID {
				events = append(events, event)
			}
			return nil
		})
	})
	sort.Slice(events, func(i, j int) bool { return events[i].ID < events[j].ID })
	return events, err
}

func (s *Store) SaveWorkflow(flow domain.Workflow) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Update(func(tx *bbolt.Tx) error { return put(tx, bucketFlows, flow.ID, flow) })
}
func (s *Store) GetWorkflow(id string) (domain.Workflow, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var flow domain.Workflow
	err := s.db.View(func(tx *bbolt.Tx) error { return get(tx.Bucket(bucketFlows), id, &flow) })
	return flow, err
}
func (s *Store) SaveAttachment(attachment domain.Attachment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Update(func(tx *bbolt.Tx) error { return put(tx, bucketAttachments, attachment.ID, attachment) })
}
func (s *Store) ListAttachments(recordID string) ([]domain.Attachment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]domain.Attachment, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketAttachments).ForEach(func(_, v []byte) error {
			var a domain.Attachment
			if e := decode(v, &a); e != nil {
				return e
			}
			if recordID == "" || a.RecordID == recordID {
				items = append(items, a)
			}
			return nil
		})
	})
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items, err
}
