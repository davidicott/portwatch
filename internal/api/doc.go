// Package api exposes a minimal HTTP interface for the portwatch daemon.
//
// Endpoints:
//
//	 GET /healthz  — liveness probe; returns 200 OK with body "ok"
//	 GET /metrics  — JSON snapshot of runtime metrics (scan counts, uptime, etc.)
//	 GET /history  — JSON array of the 50 most recent alert events
//
// The server is created with [New] and started by calling [Server.Start] with
// a context. Cancelling the context triggers a graceful shutdown with a
// five-second deadline.
//
// Example:
//
//	srv := api.New(":8080", m, h)
//	if err := srv.Start(ctx); err != nil {
//	    log.Fatal(err)
//	}
package api
