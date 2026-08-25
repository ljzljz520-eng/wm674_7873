package httpapi

import (
	"encoding/json"
	"meter-sync/internal/domain"
	"meter-sync/internal/service"
	"net/http"
)

type Handler struct{ Service *service.Service }

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	query := domain.Query{Manufacturer: r.URL.Query().Get("manufacturer"), Model: r.URL.Query().Get("model"), Reference: r.URL.Query().Get("reference"), Status: domain.Status(r.URL.Query().Get("status"))}
	records, err := h.Service.Search(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(records)
}
func Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
