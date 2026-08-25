package validation

import (
	"meter-sync/internal/domain"
	"testing"
)

func TestValidateRecord(t *testing.T) {
	issues := ValidateRecord(domain.Record{})
	if len(issues) != 2 {
		t.Fatalf("issues=%v", issues)
	}
	if Explain(nil) != "valid" {
		t.Fatal("valid explanation")
	}
}
