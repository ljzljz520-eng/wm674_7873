package dashboard

import (
	"meter-sync/internal/domain"
	"strings"
	"testing"
)

func TestDashboard(t *testing.T) {
	dashboard := Build([]domain.Record{{ID: "r", Manufacturer: "A", Amount: 3, Status: domain.StatusDraft, Version: 1}})
	if dashboard.IsEmpty() || dashboard.TopManufacturer != "A" {
		t.Fatalf("dashboard=%+v", dashboard)
	}
	if !strings.Contains(dashboard.ExportCSV(), "manufacturer") {
		t.Fatal("csv")
	}
}
