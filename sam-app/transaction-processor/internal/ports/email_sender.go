package ports

import (
	"gopkg.in/mail.v2"
	"transaction-processor/internal/domain/model"
)

// EmailSummary contains the data to be included in the summary email
type EmailSummary struct {
	TotalBalance             float64
	MonthlyTransactionCounts map[string]int
	AverageCreditAmount      float64
	AverageDebitAmount       float64
}

// EmailSender defines the interface for sending summary emails
type EmailSender interface {
	// SendSummaryEmail sends a summary email with account information
	SendSummaryEmail(recipient string, summary EmailSummary) error
}

// NewEmailSummaryFromAccount creates an EmailSummary from an Account
func NewEmailSummaryFromAccount(account *model.Account) EmailSummary {
	return EmailSummary{
		TotalBalance:             account.GetTotalBalance(),
		MonthlyTransactionCounts: account.GetMonthlyTransactionCounts(),
		AverageCreditAmount:      account.GetAverageCreditAmount(),
		AverageDebitAmount:       account.GetAverageDebitAmount(),
	}
}

// Mail-related interfaces for dependency injection and testing

// MailDialer interface for mail operations
type MailDialer interface {
	DialAndSend(m ...*mail.Message) error
}

// MailMessage interface for mail message operations
type MailMessage interface {
	SetHeader(field string, value ...string)
	SetBody(contentType, body string, settings ...mail.PartSetting)
}

// MailMessageFactory creates mail messages
type MailMessageFactory interface {
	NewMessage() MailMessage
}
