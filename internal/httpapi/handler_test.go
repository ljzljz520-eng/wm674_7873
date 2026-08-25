package httpapi

import (
	"meter-sync/internal/service"
	"meter-sync/internal/store"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestHandlerSearch(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc, _ := service.New(st, service.FixedClock{Value: "t"})
	_, _ = svc.Register("r", "Maker", "M", "S", "u", 1)
	request := httptest.NewRequest("GET", "/?manufacturer=Maker", nil)
	response := httptest.NewRecorder()
	(Handler{Service: svc}).ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatalf("code=%d", response.Code)
	}
}
