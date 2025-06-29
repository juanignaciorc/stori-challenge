package adapters

import (
	"context"
	"fmt"
	"strings"
	"transaction-processor/internal/ports"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

// SESEmailSender implements the EmailSender port using AWS SES
type SESEmailSender struct {
	sesClient *ses.Client
	sender    string
}

// NewSESEmailSender creates a new SESEmailSender
func NewSESEmailSender(sesClient *ses.Client, sender string) *SESEmailSender {
	return &SESEmailSender{
		sesClient: sesClient,
		sender:    sender,
	}
}

// SendSummaryEmail sends a summary email using AWS SES
func (s *SESEmailSender) SendSummaryEmail(recipient string, summary ports.EmailSummary) error {
	// Create the email content
	subject := "Your Account Transaction Summary"
	htmlBody := s.generateHTMLEmail(summary)
	textBody := s.generateTextEmail(summary)

	// Create the email message
	message := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{recipient},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(htmlBody),
				},
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(textBody),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(s.sender),
	}

	// Send the email
	_, err := s.sesClient.SendEmail(context.TODO(), message)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}

// generateHTMLEmail generates the HTML content for the summary email
func (s *SESEmailSender) generateHTMLEmail(summary ports.EmailSummary) string {
	var html strings.Builder

	html.WriteString("<!DOCTYPE html><html><head><style>")
	html.WriteString("body { font-family: Arial, sans-serif; line-height: 1.6; }")
	html.WriteString("h1 { color: #2a5885; }")
	html.WriteString("table { border-collapse: collapse; width: 100%; }")
	html.WriteString("th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }")
	html.WriteString("th { background-color: #f2f2f2; }")
	html.WriteString("</style></head><body>")

	html.WriteString("<h1>Account Transaction Summary</h1>")

	// Total balance
	html.WriteString(fmt.Sprintf("<p><strong>Total balance is:</strong> $%.2f</p>", summary.TotalBalance))

	// Transactions by month
	html.WriteString("<h2>Monthly Transaction Count</h2>")
	html.WriteString("<table><tr><th>Month</th><th>Number of Transactions</th></tr>")
	for month, count := range summary.MonthlyTransactionCounts {
		html.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%d</td></tr>", month, count))
	}
	html.WriteString("</table>")

	// Average amounts
	html.WriteString("<h2>Average Transaction Amounts</h2>")
	html.WriteString("<table>")
	html.WriteString(fmt.Sprintf("<tr><td>Average debit amount:</td><td>$%.2f</td></tr>", summary.AverageDebitAmount))
	html.WriteString(fmt.Sprintf("<tr><td>Average credit amount:</td><td>$%.2f</td></tr>", summary.AverageCreditAmount))
	html.WriteString("</table>")

	html.WriteString("</body></html>")

	return html.String()
}

// generateTextEmail generates the plain text content for the summary email
func (s *SESEmailSender) generateTextEmail(summary ports.EmailSummary) string {
	var text strings.Builder

	text.WriteString("Account Transaction Summary\n\n")

	// Total balance
	text.WriteString(fmt.Sprintf("Total balance is: $%.2f\n\n", summary.TotalBalance))

	// Transactions by month
	text.WriteString("Monthly Transaction Count:\n")
	for month, count := range summary.MonthlyTransactionCounts {
		text.WriteString(fmt.Sprintf("Number of transactions in %s: %d\n", month, count))
	}
	text.WriteString("\n")

	// Average amounts
	text.WriteString(fmt.Sprintf("Average debit amount: $%.2f\n", summary.AverageDebitAmount))
	text.WriteString(fmt.Sprintf("Average credit amount: $%.2f\n", summary.AverageCreditAmount))

	return text.String()
}
