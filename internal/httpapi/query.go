package httpapi

import (
	"meter-sync/internal/domain"
	"net/http"
	"strconv"
)

func QueryFromRequest(request *http.Request) domain.Query {
	values := request.URL.Query()
	limit, _ := strconv.Atoi(values.Get("limit"))
	return domain.Query{Manufacturer: values.Get("manufacturer"), Model: values.Get("model"), Reference: values.Get("reference"), Status: domain.Status(values.Get("status")), Limit: limit}
}
func MethodAllowed(request *http.Request, methods ...string) bool {
	for _, method := range methods {
		if request.Method == method {
			return true
		}
	}
	return false
}
func HeaderValue(request *http.Request, name string) string { return request.Header.Get(name) }
func IsJSON(request *http.Request) bool {
	return request.Header.Get("Accept") == "application/json" || request.Header.Get("Content-Type") == "application/json"
}
