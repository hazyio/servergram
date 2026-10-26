package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultCORSAllowsNextDevPreflight(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "")

	req := httptest.NewRequest(http.MethodOptions, "/account/validate", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	req.Header.Set("Access-Control-Request-Headers", "content-type")
	rec := httptest.NewRecorder()

	NewRouter(nil).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "http://localhost:3000")
	}
}
