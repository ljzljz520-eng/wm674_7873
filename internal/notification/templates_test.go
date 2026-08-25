package notification

import "testing"

func TestTemplates(t *testing.T) {
	template := FindTemplate("review-request")
	if !template.Enabled {
		t.Fatal("template disabled")
	}
	if len(EnabledTemplates()) != 4 {
		t.Fatal("template count")
	}
}
