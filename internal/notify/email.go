package notify

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// EmailConfig holds SMTP configuration for the email notifier.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

type emailNotifier struct {
	cfg  EmailConfig
	auth smtp.Auth
}

// NewEmailNotifier creates a Notifier that sends alerts via SMTP email.
func NewEmailNotifier(cfg EmailConfig) Notifier {
	var auth smtp.Auth
	if cfg.Username != "" {
		auth = smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	}
	return &emailNotifier{cfg: cfg, auth: auth}
}

func (e *emailNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	subject := fmt.Sprintf("portwatch: %d port change(s) detected", len(events))
	body := buildEmailBody(events)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		e.cfg.From,
		strings.Join(e.cfg.To, ", "),
		subject,
		body,
	))

	addr := fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port)
	return smtp.SendMail(addr, e.auth, e.cfg.From, e.cfg.To, msg)
}

func buildEmailBody(events []alert.Event) string {
	var sb strings.Builder
	sb.WriteString("Port change summary:\n\n")
	for _, ev := range events {
		sb.WriteString(fmt.Sprintf("  [%s] %s://%s:%d\n",
			strings.ToUpper(ev.Kind),
			ev.Port.Proto,
			ev.Port.Addr,
			ev.Port.Number,
		))
	}
	return sb.String()
}
