package notify

import (
	"fmt"
	"log/syslog"

	"github.com/user/portwatch/internal/alert"
)

// SyslogNotifier sends alerts to the local syslog daemon.
type SyslogNotifier struct {
	writer *syslog.Writer
	tag    string
}

// NewSyslogNotifier creates a SyslogNotifier that writes to syslog with the
// given tag. Priority defaults to LOG_ALERT | LOG_DAEMON.
func NewSyslogNotifier(tag string) (*SyslogNotifier, error) {
	if tag == "" {
		tag = "portwatch"
	}
	w, err := syslog.New(syslog.LOG_ALERT|syslog.LOG_DAEMON, tag)
	if err != nil {
		return nil, fmt.Errorf("syslog: open writer: %w", err)
	}
	return &SyslogNotifier{writer: w, tag: tag}, nil
}

// Notify writes each event as a syslog alert message.
func (s *SyslogNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}
	for _, e := range events {
		msg := fmt.Sprintf("portwatch: %s %s/%d", e.Kind, e.Port.Protocol, e.Port.Port)
		if err := s.writer.Alert(msg); err != nil {
			return fmt.Errorf("syslog: write alert: %w", err)
		}
	}
	return nil
}

// Close releases the underlying syslog connection.
func (s *SyslogNotifier) Close() error {
	return s.writer.Close()
}
