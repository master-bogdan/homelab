package email

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"strings"
)

type SMTPConfig struct {
	From     string
	Host     string
	Port     int
	Username string
	Password string
}

type smtpClient struct {
	cfg SMTPConfig
}

func NewSMTPClient(cfg SMTPConfig) Client {
	if strings.TrimSpace(cfg.From) == "" || strings.TrimSpace(cfg.Host) == "" || cfg.Port <= 0 {
		return NewNoopClient()
	}

	cfg.From = strings.TrimSpace(cfg.From)
	cfg.Host = strings.TrimSpace(cfg.Host)
	cfg.Username = strings.TrimSpace(cfg.Username)

	return &smtpClient{cfg: cfg}
}

func (c *smtpClient) Send(_ context.Context, msg Message) error {
	recipients := normalizeRecipients(msg.To)
	if len(recipients) == 0 {
		return fmt.Errorf("email recipients are required")
	}

	subject := strings.TrimSpace(msg.Subject)
	if subject == "" {
		return fmt.Errorf("email subject is required")
	}

	body := strings.TrimSpace(msg.TextBody)
	if body == "" {
		return fmt.Errorf("email body is required")
	}

	var auth smtp.Auth
	if c.cfg.Username != "" || c.cfg.Password != "" {
		auth = smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.Host)
	}

	return smtp.SendMail(
		net.JoinHostPort(c.cfg.Host, strconv.Itoa(c.cfg.Port)),
		auth,
		c.cfg.From,
		recipients,
		buildMessage(c.cfg.From, recipients, subject, body),
	)
}

func normalizeRecipients(recipients []string) []string {
	normalized := make([]string, 0, len(recipients))
	for _, recipient := range recipients {
		trimmed := strings.TrimSpace(recipient)
		if trimmed != "" {
			normalized = append(normalized, trimmed)
		}
	}

	return normalized
}

func buildMessage(from string, to []string, subject, body string) []byte {
	var buffer bytes.Buffer

	buffer.WriteString("From: " + from + "\r\n")
	buffer.WriteString("To: " + strings.Join(to, ", ") + "\r\n")
	buffer.WriteString("Subject: " + subject + "\r\n")
	buffer.WriteString("MIME-Version: 1.0\r\n")
	buffer.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	buffer.WriteString("\r\n")
	buffer.WriteString(body)
	buffer.WriteString("\r\n")

	return buffer.Bytes()
}
