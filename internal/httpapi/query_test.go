package httpapi

import (
	"net/http/httptest"
	"testing"
)

func TestQueryFromRequest(t *testing.T) {
	request := httptest.NewRequest("GET", "/?manufacturer=M&limit=4", nil)
	query := QueryFromRequest(request)
	if query.Manufacturer != "M" || query.Limit != 4 {
		t.Fatal(query)
	}
}
