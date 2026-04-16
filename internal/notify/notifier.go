// Package notify provides pluggable notification backends for portwatch alerts.
package notify

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// Notifier sends alert events to some destination.
type Notifier interface {
	Notify(ctx context.Context, events []alert.Event) error
}

// StdoutNotifier writes human-readable alerts to anoutNotifier struct {
	w io.Writer
}

// NewStdoutNotifier returns a StdoutNotifier writing to w.
// If w is nil, os.Stdout is used.
func NewStdoutNotifier(w io.Writer) *StdoutNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &StdoutNotifier{w: w}
}

// Notify writes each event as a formatted line to the writer.
// It stops and returns an error on the first write failure.
func (s *StdoutNotifier) Notify(_ context.Context, events []alert.Event) error {
	for _, e := range events {
		line := fmt.Sprintf("[portwatch] %s %s\n", strings.ToUpper(string(e.Kind)), e.Port)
		if _, err := fmt.Fprint(s.w, line); err != nil {
			return fmt.Errorf("stdout notifier: write failed for port %s: %w", e.Port, err)
		}
	}
	return nil
}
