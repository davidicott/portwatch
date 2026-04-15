package metrics

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Reporter writes human-readable metric summaries to a writer.
type Reporter struct {
	w io.Writer
}

// NewReporter creates a Reporter that writes to w.
func NewReporter(w io.Writer) *Reporter {
	return &Reporter{w: w}
}

// Report formats and writes the current counters snapshot to the writer.
func (rp *Reporter) Report(c Counters) error {
	tw := tabwriter.NewWriter(rp.w, 0, 0, 2, ' ', 0)

	uptime := time.Since(c.StartTime).Truncate(time.Second)
	lastScan := "never"
	if !c.LastScanTime.IsZero() {
		lastScan = c.LastScanTime.Format(time.RFC3339)
	}

	lines := []struct{ k, v string }{
		{"uptime", uptime.String()},
		{"scans_total", fmt.Sprintf("%d", c.ScansTotal)},
		{"alerts_total", fmt.Sprintf("%d", c.AlertsTotal)},
		{"ports_opened", fmt.Sprintf("%d", c.OpenedPorts)},
		{"ports_closed", fmt.Sprintf("%d", c.ClosedPorts)},
		{"last_scan", lastScan},
	}

	for _, l := range lines {
		if _, err := fmt.Fprintf(tw, "%s\t%s\n", l.k, l.v); err != nil {
			return err
		}
	}
	return tw.Flush()
}
