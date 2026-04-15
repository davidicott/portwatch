package history

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// ExportJSON writes all events from the history as a JSON array to w.
func (h *History) ExportJSON(w io.Writer) error {
	events := h.Latest(h.Len())
	return json.NewEncoder(w).Encode(events)
}

// ExportTable writes all events from the history as a human-readable table to w.
func (h *History) ExportTable(w io.Writer) error {
	events := h.Latest(h.Len())
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIME\tEVENT\tPROTO\tADDR\tPORT\tPID")
	for _, e := range events {
		pid := "-"
		if e.Port.PID != 0 {
			pid = fmt.Sprintf("%d", e.Port.PID)
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%d\t%s\n",
			e.Time.Format(time.RFC3339),
			labelFor(e.Kind),
			e.Port.Proto,
			e.Port.Addr,
			e.Port.Port,
			pid,
		)
	}
	return tw.Flush()
}

func labelFor(kind alert.EventKind) string {
	switch kind {
	case alert.Opened:
		return "OPENED"
	case alert.Closed:
		return "CLOSED"
	default:
		return "UNKNOWN"
	}
}
