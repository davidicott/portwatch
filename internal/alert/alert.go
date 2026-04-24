package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event describes a port change event.
type Event struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Port      scanner.Port
}

// Notifier sends alert events to a destination.
type Notifier interface {
	Notify(e Event) error
}

// LogNotifier writes alerts to an io.Writer (default: os.Stdout).
type LogNotifier struct {
	Out io.Writer
}

// NewLogNotifier creates a LogNotifier writing to stdout.
func NewLogNotifier() *LogNotifier {
	return &LogNotifier{Out: os.Stdout}
}

// Notify formats and writes the event to the configured writer.
func (l *LogNotifier) Notify(e Event) error {
	_, err := fmt.Fprintf(
		l.Out,
		"[%s] %s — %s\n",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.Message,
	)
	return err
}

// NotifyAll sends all events through the notifier, returning the first error
// encountered. Events after the failing one are not sent.
func NotifyAll(n Notifier, events []Event) error {
	for _, e := range events {
		if err := n.Notify(e); err != nil {
			return fmt.Errorf("alert: failed to notify event %q: %w", e.Message, err)
		}
	}
	return nil
}

// BuildEvents converts a diff result into a slice of alert Events.
func BuildEvents(opened, closed []scanner.Port) []Event {
	now := time.Now()
	events := make([]Event, 0, len(opened)+len(closed))

	for _, p := range opened {
		events = append(events, Event{
			Timestamp: now,
			Level:     LevelAlert,
			Message:   fmt.Sprintf("port opened: %s", p),
			Port:      p,
		})
	}

	for _, p := range closed {
		events = append(events, Event{
			Timestamp: now,
			Level:     LevelWarn,
			Message:   fmt.Sprintf("port closed: %s", p),
			Port:      p,
		})
	}

	return events
}
