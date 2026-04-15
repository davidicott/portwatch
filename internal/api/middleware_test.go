package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggingMiddleware_PassesThrough(t *testing.T) {
	handler := loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Errorf("expected status 202, got %d", rec.Code)
	}
}

func TestRecoveryMiddleware_HandlesPanic(t *testing.T) {
	handler := recoveryMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	}))

	rec := httptest.NewRec := httptest.NewRequest(http.MethodGet, "/panic", nil)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}
}

func TestResponseRecorder_DefaultStatus(t *testing.T) {
	w := httptest.NewRecorder()
	rr := newResponseRecorder(w)

	if rr.statusCode != http.StatusOK {
		t.Errorf("expected default status 200, got %d", rr.statusCode)
	}
}

func TestResponseRecorder_CapturesStatus(t *testing.T) {
	w := httptest.NewRecorder()
	rr := newResponseRecorder(w)
	rr.WriteHeader(http.StatusNotFound)

	if rr.statusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.statusCode)
	}
}
