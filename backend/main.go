package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SendEmailRequest represents the JSON request body for sending email
type SendEmailRequest struct {
	To      string `json:"to" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

// SendEmailResponse represents the JSON response after sending email
type SendEmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// EmailsResponse represents the JSON response for fetching emails
type EmailsResponse struct {
	Success bool        `json:"success"`
	Emails  []EmailInfo `json:"emails"`
	Error   string      `json:"error,omitempty"`
	Count   int         `json:"count"`
}

// Configuration for SMTP and IMAP servers
var (
	smtpConfig = SMTPConfig{
		Host:     "localhost",
		Port:     "587",
		Username: "test@mail.local",
		Password: "password123",
	}

	imapConfig = IMAPConfig{
		Host:     "localhost",
		Port:     "143",
		Username: "test@mail.local",
		Password: "password123",
	}
)

func main() {
	// Create Gin router with default middleware
	router := gin.Default()

	// Define routes
	router.POST("/send-email", handleSendEmail)
	router.GET("/emails", handleFetchEmails)
	router.GET("/health", handleHealth)

	// Start server
	router.Run(":8080")
}

// handleSendEmail processes POST /send-email requests
func handleSendEmail(c *gin.Context) {
	var req SendEmailRequest

	// Bind and validate JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, SendEmailResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Create email message
	emailMsg := EmailMessage{
		From:    smtpConfig.Username,
		To:      []string{req.To},
		Subject: req.Subject,
		Body:    req.Body,
	}

	// Send email
	if err := SendEmail(smtpConfig, emailMsg); err != nil {
		c.JSON(http.StatusInternalServerError, SendEmailResponse{
			Success: false,
			Message: "Failed to send email",
			Error:   err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, SendEmailResponse{
		Success: true,
		Message: "Email sent successfully",
	})
}

// handleFetchEmails processes GET /emails requests
func handleFetchEmails(c *gin.Context) {
	// Fetch emails from IMAP server
	emails, err := FetchEmails(imapConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, EmailsResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Success response with email list
	if emails == nil {
		emails = []EmailInfo{}
	}

	c.JSON(http.StatusOK, EmailsResponse{
		Success: true,
		Emails:  emails,
		Count:   len(emails),
	})
}

// handleHealth provides a simple health check endpoint
func handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "email-api",
	})
}
