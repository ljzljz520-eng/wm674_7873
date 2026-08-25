package approval

import (
	"errors"
	"fmt"
	"meter-sync/internal/domain"
	"sort"
	"strings"
)

type Decision string

const (
	Approved Decision = "approved"
	Rejected Decision = "rejected"
	Pending  Decision = "pending"
)

type Request struct {
	ID        string
	RecordID  string
	Reviewer  string
	Reason    string
	Decision  Decision
	CreatedAt string
	DecidedAt string
}
type Queue struct{ requests map[string]Request }

func New() *Queue { return &Queue{requests: make(map[string]Request)} }
func (q *Queue) Submit(record domain.Record, reviewer, reason, stamp string) (Request, error) {
	if record.ID == "" {
		return Request{}, errors.New("record id is required")
	}
	if strings.TrimSpace(reviewer) == "" {
		return Request{}, errors.New("reviewer is required")
	}
	id := fmt.Sprintf("approval:%s:%d", record.ID, record.Version)
	request := Request{ID: id, RecordID: record.ID, Reviewer: reviewer, Reason: reason, Decision: Pending, CreatedAt: stamp}
	q.requests[id] = request
	return request, nil
}
func (q *Queue) Decide(id string, decision Decision, stamp string) error {
	request, ok := q.requests[id]
	if !ok {
		return errors.New("approval request not found")
	}
	if decision != Approved && decision != Rejected {
		return errors.New("decision must be approved or rejected")
	}
	request.Decision = decision
	request.DecidedAt = stamp
	q.requests[id] = request
	return nil
}
func (q *Queue) Get(id string) (Request, bool) { request, ok := q.requests[id]; return request, ok }
func (q *Queue) Pending(recordID string) []Request {
	result := make([]Request, 0)
	for _, request := range q.requests {
		if request.Decision == Pending && (recordID == "" || request.RecordID == recordID) {
			result = append(result, request)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}
func (q *Queue) HasApproved(recordID string) bool {
	for _, request := range q.requests {
		if request.RecordID == recordID && request.Decision == Approved {
			return true
		}
	}
	return false
}
func (q *Queue) Count(decision Decision) int {
	count := 0
	for _, request := range q.requests {
		if request.Decision == decision {
			count++
		}
	}
	return count
}
