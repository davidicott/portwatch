package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/api"
	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/metrics"
)

func newTestServer(t *testing.T) *api.Server {
	t.Helper()
	m := metrics.New()
	h := history.New(100)
	return api.New("127.0.0.1:0", m, h)
}

func TestHandleHealth(t *testing.T) {
	m := metrics.New()
	h := history.New(100)
	srv := api.New("127.0.0.1:0", m, h)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if body := rec.Body.String(); body != "ok" {
		t.Fatalf("unexpected body: %q", body)
	}
}

func TestHandleMetrics(t *testing.T) {
	m := metrics.New()
	m.RecordScan(3, time.Millisecond*12)
	h := history.New(100)
	srv := api.New("127.0.0.1:0", m, h)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var result map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if _, ok := result["total_scans"]; !ok {
		t.Error("missing total_scans key in metrics response")
	}
}

func TestHandleHistory(t *testing.T) {
	m := metrics.New()
	h := history.New(100)
	srv := api.New("127.0.0.1:0", m, h)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestServerStopsOnContextCancel(t *testing.T) {
	m := metrics.New()
	h := history.New(100)
	srv := api.New("127.0.0.1:19876", m, h)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	errCh := make(chan error, 1)
	go func() { errCh <- srv.Start(ctx) }()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("server did not stop after context cancel")
	}
}
