package main

import (
	"fmt"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// IMAPConfig holds IMAP server configuration
type IMAPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

// EmailInfo represents a simplified email message
type EmailInfo struct {
	Subject string `json:"subject"`
	From    string `json:"from"`
	Date    string `json:"date"`
}

// FetchEmails retrieves emails from IMAP server and returns a list
func FetchEmails(config IMAPConfig) ([]EmailInfo, error) {
	var emails []EmailInfo

	// Connect to IMAP server
	c, err := client.Dial(config.Host + ":" + config.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to IMAP server: %v", err)
	}
	defer c.Logout()

	// Login
	if err := c.Login(config.Username, config.Password); err != nil {
		return nil, fmt.Errorf("failed to login to IMAP: %v", err)
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("failed to select INBOX: %v", err)
	}

	// If no messages, return empty list
	if mbox.Messages == 0 {
		return emails, nil
	}

	// Fetch all messages
	seqset := new(imap.SeqSet)
	seqset.AddRange(1, mbox.Messages)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	// Collect email information
	for msg := range messages {
		if msg.Envelope != nil {
			from := ""
			if len(msg.Envelope.From) > 0 {
				from = msg.Envelope.From[0].Address()
			}

			emailInfo := EmailInfo{
				Subject: msg.Envelope.Subject,
				From:    from,
				Date:    msg.Envelope.Date.String(),
			}
			emails = append(emails, emailInfo)
		}
	}

	// Wait for fetch to complete
	if err := <-done; err != nil {
		return nil, fmt.Errorf("failed to fetch emails: %v", err)
	}

	return emails, nil
}
