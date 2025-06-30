package adapters

import (
	"bytes"
	"fmt"
	"gopkg.in/mail.v2"
	"html/template"
	"transaction-processor/internal/ports"
)

// SMTPConfiguration holds the SMTP configuration
type SMTPConfiguration struct {
	Sender     string
	Password   string
	SmtpServer string
	SmtpPort   int
}

// RealMailMessage wraps mail.Message to implement MailMessage interface
type RealMailMessage struct {
	*mail.Message
}

func (r *RealMailMessage) SetHeader(field string, value ...string) {
	r.Message.SetHeader(field, value...)
}

func (r *RealMailMessage) SetBody(contentType, body string, settings ...mail.PartSetting) {
	r.Message.SetBody(contentType, body, settings...)
}

// RealMailMessageFactory implements MailMessageFactory
type RealMailMessageFactory struct{}

func (r *RealMailMessageFactory) NewMessage() ports.MailMessage {
	return &RealMailMessage{mail.NewMessage()}
}

// RealMailDialer wraps mail.Dialer to implement MailDialer interface
type RealMailDialer struct {
	*mail.Dialer
}

func (r *RealMailDialer) DialAndSend(m ...*mail.Message) error {
	return r.Dialer.DialAndSend(m...)
}

// SMTPClient implements the EmailSender port using SMTP
type SMTPClient struct {
	dialer         ports.MailDialer
	messageFactory ports.MailMessageFactory
	sender         string
}

// NewSMTPEmailSender creates a new SMTPClient with real dependencies
func NewSMTPEmailSender(conf SMTPConfiguration) *SMTPClient {
	dialer := mail.NewDialer(conf.SmtpServer, conf.SmtpPort, conf.Sender, conf.Password)
	dialer.StartTLSPolicy = mail.MandatoryStartTLS

	return &SMTPClient{
		dialer:         &RealMailDialer{dialer},
		messageFactory: &RealMailMessageFactory{},
		sender:         conf.Sender,
	}
}

// NewSMTPEmailSenderWithDependencies creates a new SMTPClient with injected dependencies (for testing)
func NewSMTPEmailSenderWithDependencies(
	dialer ports.MailDialer,
	messageFactory ports.MailMessageFactory,
	sender string,
) *SMTPClient {
	return &SMTPClient{
		dialer:         dialer,
		messageFactory: messageFactory,
		sender:         sender,
	}
}

// SendSummaryEmail sends a summary email using SMTP with gomail
func (s *SMTPClient) SendSummaryEmail(recipient string, summary ports.EmailSummary) error {
	// Generate the email content using the HTML template
	emailBody, err := s.generateEmailBody(recipient, summary)
	if err != nil {
		return fmt.Errorf("error generating email body: %w", err)
	}

	// Create the email message using the factory
	msg := s.messageFactory.NewMessage()
	msg.SetHeader("From", s.sender)
	msg.SetHeader("To", recipient)
	msg.SetHeader("Subject", "Transaction Summary")
	msg.SetBody("text/html", emailBody)

	// We need to create a concrete mail.Message for DialAndSend
	// This is necessary because the mail library expects *mail.Message
	realMsg := mail.NewMessage()
	realMsg.SetHeader("From", s.sender)
	realMsg.SetHeader("To", recipient)
	realMsg.SetHeader("Subject", "Transaction Summary")
	realMsg.SetBody("text/html", emailBody)

	// Send the email
	if err := s.dialer.DialAndSend(realMsg); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}

// generateEmailBody generates the email content using the HTML template
func (s *SMTPClient) generateEmailBody(recipient string, summary ports.EmailSummary) (string, error) {
	if recipient == "" {
		return "", fmt.Errorf("recipient cannot be empty")
	}

	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>Transaction Summary</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        h1 { color: #2a5885; }
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <h1>Transaction Summary</h1>
    <p><strong>Total balance is:</strong> ${{printf "%.2f" .TotalBalance}}</p>

    <h2>Monthly Transaction Count</h2>
    <table>
        <tr><th>Month</th><th>Number of Transactions</th></tr>
        {{range $month, $count := .MonthlyTransactionCounts}}
        <tr><td>{{$month}}</td><td>{{$count}}</td></tr>
        {{end}}
    </table>
    
    <h2>Average Transaction Amounts</h2>
    <table>
        <tr><td>Average debit amount:</td><td>${{printf "%.2f" .AverageDebitAmount}}</td></tr>
        <tr><td>Average credit amount:</td><td>${{printf "%.2f" .AverageCreditAmount}}</td></tr>
    </table>
</body>
</html>`

	tmplData := struct {
		TotalBalance             float64
		MonthlyTransactionCounts map[string]int
		AverageCreditAmount      float64
		AverageDebitAmount       float64
	}{
		TotalBalance:             summary.TotalBalance,
		MonthlyTransactionCounts: summary.MonthlyTransactionCounts,
		AverageCreditAmount:      summary.AverageCreditAmount,
		AverageDebitAmount:       summary.AverageDebitAmount,
	}

	var emailBody bytes.Buffer
	tmplObj := template.Must(template.New("emailTemplate").Parse(tmpl))
	err := tmplObj.Execute(&emailBody, tmplData)
	if err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return emailBody.String(), nil
}
