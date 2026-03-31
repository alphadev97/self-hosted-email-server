package main

import (
	"fmt"
	"net/smtp"
	"time"
)

// SMTPConfig holds SMTP server configuration
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

// EmailMessage represents an email to send
type EmailMessage struct {
	From    string
	To      []string
	Subject string
	Body    string
}

// SendEmail sends an email via SMTP with RFC-compliant headers
func SendEmail(config SMTPConfig, msg EmailMessage) error {
	if msg.From == "" {
		msg.From = config.Username
	}

	// Generate proper Message-ID (required for Postfix/Amavis)
	messageID := generateMessageID(msg.From)

	// Get current date in RFC 5322 format
	currentDate := time.Now().Format(time.RFC1123Z)

	// Build RFC-compliant email message
	messageBytes := []byte(
		"From: " + msg.From + "\r\n" +
			"To: " + msg.To[0] + "\r\n" +
			"Subject: " + msg.Subject + "\r\n" +
			"Message-ID: " + messageID + "\r\n" +
			"Date: " + currentDate + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
			"Content-Transfer-Encoding: 8bit\r\n" +
			"\r\n" +
			msg.Body,
	)

	// Create SMTP auth
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	// Send email
	err := smtp.SendMail(
		config.Host+":"+config.Port,
		auth,
		msg.From,
		msg.To,
		messageBytes,
	)

	return err
}

// generateMessageID creates an RFC 5322 compliant Message-ID
// Format: <timestamp.randomID@domain>
func generateMessageID(from string) string {
	// Extract domain from email address
	var domain string
	for i := len(from) - 1; i >= 0; i-- {
		if from[i] == '@' {
			domain = from[i+1:]
			break
		}
	}
	if domain == "" {
		domain = "mail.local"
	}

	// Create Message-ID with timestamp
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("<%d.go@%s>", timestamp, domain)
}
