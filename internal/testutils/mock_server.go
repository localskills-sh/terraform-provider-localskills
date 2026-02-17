package testutils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// NewMockLocalskillsServer creates a test HTTP server with basic health check endpoint.
// Callers can add additional handlers to the returned mux before starting tests.
func NewMockLocalskillsServer() (*httptest.Server, *http.ServeMux) {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    map[string]string{"status": "ok"},
		})
	})

	server := httptest.NewServer(mux)
	return server, mux
}
