package format

import (
	"meter-sync/internal/domain"
	"testing"
)

func TestRecordJSON(t *testing.T) {
	data, err := RecordJSON(domain.Record{ID: "r"})
	if err != nil || len(data) == 0 {
		t.Fatalf("data=%s err=%v", data, err)
	}
}
