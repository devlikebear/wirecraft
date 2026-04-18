package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHandlesHealthz(t *testing.T) {
	handler := New()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /healthz status = %d, want %d", rec.Code, http.StatusOK)
	}

	if got := rec.Body.String(); got != "ok\n" {
		t.Fatalf("GET /healthz body = %q, want %q", got, "ok\n")
	}
}
