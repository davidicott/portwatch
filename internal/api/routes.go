package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/user/portwatch/internal/history"
)

// handleHealth returns a simple liveness check response.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// handleMetrics returns current runtime metrics as JSON.
func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	snap := s.reporter.Snapshot()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(snap)
}

// handleHistory returns recent alert history, optionally limited by ?n=<count>.
func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	n := 50
	if raw := r.URL.Query().Get("n"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 {
			n = v
		}
	}

	events := s.history.Latest(n)

	format := r.URL.Query().Get("format")
	switch format {
	case "table":
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_ = history.ExportTable(w, events)
	default:
		w.Header().Set("Content-Type", "application/json")
		_ = history.ExportJSON(w, events)
	}
}
