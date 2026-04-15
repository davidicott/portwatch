// Package api provides a lightweight HTTP server exposing portwatch
// runtime data such as metrics and recent alert history.
package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/metrics"
)

// Server is a minimal HTTP API server.
type Server struct {
	addr     string
	metrics  *metrics.Metrics
	history  *history.History
	httpSrv  *http.Server
}

// New creates a new Server bound to addr.
func New(addr string, m *metrics.Metrics, h *history.History) *Server {
	s := &Server{
		addr:    addr,
		metrics: m,
		history: h,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", s.handleMetrics)
	mux.HandleFunc("/history", s.handleHistory)
	mux.HandleFunc("/healthz", s.handleHealth)

	s.httpSrv = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return s
}

// Start begins serving HTTP requests. It returns when ctx is cancelled.
func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.httpSrv.Shutdown(shutCtx)
	}
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	snap := s.metrics.Snapshot()
	jsonResponse(w, snap)
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	events := s.history.Latest(50)
	jsonResponse(w, events)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func jsonResponse(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
	}
}
