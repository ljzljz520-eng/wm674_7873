package domain

import "strings"

type Query struct {
	Manufacturer string
	Model        string
	Status       Status
	Reference    string
	Limit        int
}

func (q Query) Matches(r Record) bool {
	if q.Manufacturer != "" && !strings.Contains(strings.ToLower(r.Manufacturer), strings.ToLower(q.Manufacturer)) {
		return false
	}
	if q.Model != "" && !strings.Contains(strings.ToLower(r.MeterModel), strings.ToLower(q.Model)) {
		return false
	}
	if q.Status != "" && r.Status != q.Status {
		return false
	}
	if q.Reference != "" && r.SyncReference != q.Reference {
		return false
	}
	return true
}

type ImportRow struct {
	ID           string
	Manufacturer string
	Model        string
	Reference    string
	Amount       int64
}

type ImportResult struct {
	Accepted []Record
	Rejected []ImportError
}

type ImportError struct {
	Row   int
	Cause string
}
