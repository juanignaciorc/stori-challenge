package config

import (
	"os"
	"strconv"
)

// Configuration holds the Lambda function configuration
type Configuration struct {
	EmailSender       string `json:"emailSender"`
	EmailPassword     string `json:"emailPassword"`
	SmtpServer        string `json:"smtpServer"`
	SmtpPort          int    `json:"smtpPort"`
	TransactionsTable string `json:"transactionsTable"`
	AccountsTable     string `json:"accountsTable"`
	AccountID         string `json:"accountID"`
}

// Load loads configuration from environment variables
func Load() Configuration {
	// Get SMTP port from environment variable, default to 587 if not set
	smtpPort := 587
	if portStr := os.Getenv("SMTP_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			smtpPort = port
		}
	}

	config := Configuration{
		EmailSender:       os.Getenv("EMAIL_SENDER"),
		EmailPassword:     os.Getenv("EMAIL_PASSWORD"),
		SmtpServer:        os.Getenv("SMTP_SERVER"),
		SmtpPort:          smtpPort,
		TransactionsTable: os.Getenv("TRANSACTIONS_TABLE"),
		AccountsTable:     os.Getenv("ACCOUNTS_TABLE"),
		AccountID:         os.Getenv("ACCOUNT_ID"),
	}

	// If ACCOUNT_ID is not set, use a default value
	if config.AccountID == "" {
		config.AccountID = "default"
	}

	// Set default SMTP server if not provided
	if config.SmtpServer == "" {
		config.SmtpServer = "smtp.gmail.com"
	}

	return config
}