package adapters

import (
	"gopkg.in/mail.v2"
	"strings"
	"testing"
	"transaction-processor/internal/ports"
)

func TestNewSMTPEmailSender(t *testing.T) {
	type args struct {
		conf SMTPConfiguration
	}

	// Create test configuration
	testConf := SMTPConfiguration{
		Sender:     "test@example.com",
		Password:   "password123",
		SmtpServer: "smtp.example.com",
		SmtpPort:   587,
	}

	// Create expected dialer
	expectedDialer := mail.NewDialer(testConf.SmtpServer, testConf.SmtpPort, testConf.Sender, testConf.Password)
	expectedDialer.StartTLSPolicy = mail.MandatoryStartTLS

	tests := []struct {
		name string
		args args
		want *SMTPClient
	}{
		{
			name: "Creates SMTP client with correct configuration",
			args: args{
				conf: testConf,
			},
			want: &SMTPClient{
				dialer: expectedDialer,
				sender: testConf.Sender,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSMTPEmailSender(tt.args.conf)

			// Check sender field
			if got.sender != tt.want.sender {
				t.Errorf("NewSMTPEmailSender().sender = %v, want %v", got.sender, tt.want.sender)
			}

			// Check dialer fields individually since we can't directly compare dialers
			if got.dialer.Host != tt.want.dialer.Host {
				t.Errorf("NewSMTPEmailSender().dialer.Host = %v, want %v", got.dialer.Host, tt.want.dialer.Host)
			}
			if got.dialer.Port != tt.want.dialer.Port {
				t.Errorf("NewSMTPEmailSender().dialer.Port = %v, want %v", got.dialer.Port, tt.want.dialer.Port)
			}
			if got.dialer.Username != tt.want.dialer.Username {
				t.Errorf("NewSMTPEmailSender().dialer.Username = %v, want %v", got.dialer.Username, tt.want.dialer.Username)
			}
			if got.dialer.Password != tt.want.dialer.Password {
				t.Errorf("NewSMTPEmailSender().dialer.Password = %v, want %v", got.dialer.Password, tt.want.dialer.Password)
			}
			if got.dialer.StartTLSPolicy != tt.want.dialer.StartTLSPolicy {
				t.Errorf("NewSMTPEmailSender().dialer.StartTLSPolicy = %v, want %v", got.dialer.StartTLSPolicy, tt.want.dialer.StartTLSPolicy)
			}
		})
	}
}

// TestSMTPClient_SendSummaryEmail_Integration tests the SendSummaryEmail method
// Note: This is an integration test that would actually send an email if run.
// It's skipped by default to avoid sending actual emails during testing.
func TestSMTPClient_SendSummaryEmail_Integration(t *testing.T) {
	// Skip this test by default to avoid sending actual emails
	t.Skip("Skipping integration test that would send an actual email")

	// Create test data
	testSummary := ports.EmailSummary{
		TotalBalance: 1234.56,
		MonthlyTransactionCounts: map[string]int{
			"January":  5,
			"February": 8,
		},
		AverageCreditAmount: 100.25,
		AverageDebitAmount:  50.75,
	}

	testRecipient := "test@example.com"

	// Create a real SMTP client with test configuration
	// Note: These credentials won't work, they're just for the test
	conf := SMTPConfiguration{
		Sender:     "test@example.com",
		Password:   "password123",
		SmtpServer: "smtp.example.com",
		SmtpPort:   587,
	}

	client := NewSMTPEmailSender(conf)

	// Call the method under test
	err := client.SendSummaryEmail(testRecipient, testSummary)

	// Check there was no error
	if err != nil {
		t.Errorf("SendSummaryEmail() error = %v, want nil", err)
	}
}

func TestSMTPClient_generateEmailBody(t *testing.T) {
	type fields struct {
		dialer *mail.Dialer
		sender string
	}
	type args struct {
		recipient string
		summary   ports.EmailSummary
	}

	// Create test data
	testSummary := ports.EmailSummary{
		TotalBalance: 1234.56,
		MonthlyTransactionCounts: map[string]int{
			"January":  5,
			"February": 8,
		},
		AverageCreditAmount: 100.25,
		AverageDebitAmount:  50.75,
	}

	// We don't need to test the exact HTML output, just that it contains the expected data
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantContains []string
		wantErr bool
	}{
		{
			name: "Generates HTML email with correct data",
			fields: fields{
				dialer: nil, // Not used in generateEmailBody
				sender: "test@example.com",
			},
			args: args{
				recipient: "recipient@example.com",
				summary:   testSummary,
			},
			wantContains: []string{
				"$1234.56", // Total balance
				"January", "5", // Monthly transaction count
				"February", "8", // Monthly transaction count
				"$100.25", // Average credit amount
				"$50.75", // Average debit amount
			},
			wantErr: false,
		},
		{
			name: "Handles empty monthly transaction counts",
			fields: fields{
				dialer: nil,
				sender: "test@example.com",
			},
			args: args{
				recipient: "recipient@example.com",
				summary: ports.EmailSummary{
					TotalBalance:             500.00,
					MonthlyTransactionCounts: map[string]int{},
					AverageCreditAmount:      75.50,
					AverageDebitAmount:       25.25,
				},
			},
			wantContains: []string{
				"$500.00", // Total balance
				"$75.50",  // Average credit amount
				"$25.25",  // Average debit amount
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SMTPClient{
				dialer: tt.fields.dialer,
				sender: tt.fields.sender,
			}
			got, err := s.generateEmailBody(tt.args.recipient, tt.args.summary)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateEmailBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check that the generated HTML contains all expected strings
			for _, wantStr := range tt.wantContains {
				if !strings.Contains(got, wantStr) {
					t.Errorf("generateEmailBody() output doesn't contain expected string: %s", wantStr)
				}
			}
		})
	}
}
