package ports

import (
	"transaction-processor/internal/domain/model"
)

// EmailSummary contains the data to be included in the summary email
type EmailSummary struct {
	TotalBalance           float64
	MonthlyTransactionCounts map[string]int
	AverageCreditAmount    float64
	AverageDebitAmount     float64
}

// EmailSender defines the interface for sending summary emails
type EmailSender interface {
	// SendSummaryEmail sends a summary email with account information
	SendSummaryEmail(recipient string, summary EmailSummary) error
}

// NewEmailSummaryFromAccount creates an EmailSummary from an Account
func NewEmailSummaryFromAccount(account *model.Account) EmailSummary {
	return EmailSummary{
		TotalBalance:           account.GetTotalBalance(),
		MonthlyTransactionCounts: account.GetMonthlyTransactionCounts(),
		AverageCreditAmount:    account.GetAverageCreditAmount(),
		AverageDebitAmount:    account.GetAverageDebitAmount(),
	}
}
