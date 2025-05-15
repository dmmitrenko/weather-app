package emailing

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

type SmtpConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type Sender struct {
	config SmtpConfig
}

func NewSender(cfg SmtpConfig) *Sender {
	return &Sender{
		config: cfg,
	}
}

func (s *Sender) Send(ctx context.Context, to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	msg := strings.Join([]string{
		fmt.Sprintf("From: %s", s.config.From),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"utf-8\"",
		"",
		body,
	}, "\r\n")

	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	return smtp.SendMail(addr, auth, s.config.From, []string{to}, []byte(msg))
}
