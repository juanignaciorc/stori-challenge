package adapters

import (
	"errors"
	"strings"
	"testing"
	"transaction-processor/internal/mocks"
	"transaction-processor/internal/ports"

	"go.uber.org/mock/gomock"
)

func TestSMTPClient_SendSummaryEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		recipient     string
		summary       ports.EmailSummary
		setupMocks    func(*mocks.MockMailDialer, *mocks.MockMailMessageFactory, *mocks.MockMailMessage)
		expectedError string
		wantErr       bool
	}{
		{
			name:      "successful email send",
			recipient: "test@example.com",
			summary: ports.EmailSummary{
				TotalBalance:             1500.75,
				MonthlyTransactionCounts: map[string]int{"January": 5, "February": 3},
				AverageCreditAmount:      200.50,
				AverageDebitAmount:       150.25,
			},
			setupMocks: func(mockDialer *mocks.MockMailDialer, mockFactory *mocks.MockMailMessageFactory, mockMsg *mocks.MockMailMessage) {
				mockFactory.EXPECT().
					NewMessage().
					Return(mockMsg).
					Times(1)

				mockMsg.EXPECT().
					SetHeader("From", "sender@example.com").
					Times(1)

				mockMsg.EXPECT().
					SetHeader("To", "test@example.com").
					Times(1)

				mockMsg.EXPECT().
					SetHeader("Subject", "Transaction Summary").
					Times(1)

				mockMsg.EXPECT().
					SetBody("text/html", gomock.Any()).
					Times(1)

				mockDialer.EXPECT().
					DialAndSend(gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:      "empty recipient",
			recipient: "",
			summary: ports.EmailSummary{
				TotalBalance:             1000.0,
				MonthlyTransactionCounts: map[string]int{"January": 2},
				AverageCreditAmount:      100.0,
				AverageDebitAmount:       50.0,
			},
			setupMocks: func(mockDialer *mocks.MockMailDialer, mockFactory *mocks.MockMailMessageFactory, mockMsg *mocks.MockMailMessage) {
				// No expectations because generateEmailBody should fail before any mail operations
			},
			expectedError: "error generating email body: recipient cannot be empty",
			wantErr:       true,
		},
		{
			name:      "dialer fails",
			recipient: "test@example.com",
			summary: ports.EmailSummary{
				TotalBalance:             500.0,
				MonthlyTransactionCounts: map[string]int{"March": 1},
				AverageCreditAmount:      250.0,
				AverageDebitAmount:       250.0,
			},
			setupMocks: func(mockDialer *mocks.MockMailDialer, mockFactory *mocks.MockMailMessageFactory, mockMsg *mocks.MockMailMessage) {
				mockFactory.EXPECT().
					NewMessage().
					Return(mockMsg).
					Times(1)

				mockMsg.EXPECT().
					SetHeader(gomock.Any(), gomock.Any()).
					AnyTimes()

				mockMsg.EXPECT().
					SetBody(gomock.Any(), gomock.Any()).
					Times(1)

				mockDialer.EXPECT().
					DialAndSend(gomock.Any()).
					Return(errors.New("SMTP connection failed")).
					Times(1)
			},
			expectedError: "error sending email: SMTP connection failed",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDialer := mocks.NewMockMailDialer(ctrl)
			mockFactory := mocks.NewMockMailMessageFactory(ctrl)
			mockMsg := mocks.NewMockMailMessage(ctrl)

			tt.setupMocks(mockDialer, mockFactory, mockMsg)

			smtpClient := NewSMTPEmailSenderWithDependencies(
				mockDialer,
				mockFactory,
				"sender@example.com",
			)

			err := smtpClient.SendSummaryEmail(tt.recipient, tt.summary)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if tt.expectedError != "" && err.Error() != tt.expectedError {
					t.Errorf("expected error '%s', got: %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			}
		})
	}
}

func TestSMTPClient_generateEmailBody(t *testing.T) {
	tests := []struct {
		name              string
		recipient         string
		summary           ports.EmailSummary
		wantErr           bool
		expectedInBody    []string
		notExpectedInBody []string
	}{
		{
			name:      "valid email generation with multiple months",
			recipient: "test@example.com",
			summary: ports.EmailSummary{
				TotalBalance:             1234.56,
				MonthlyTransactionCounts: map[string]int{"January": 10, "February": 5},
				AverageCreditAmount:      100.25,
				AverageDebitAmount:       75.50,
			},
			wantErr: false,
			expectedInBody: []string{
				"$1234.56",
				"January",
				"10",
				"February",
				"5",
				"$100.25",
				"$-75.50",
				"Transaction Summary",
				"Monthly Transaction Count",
				"Average Transaction Amounts",
			},
		},
		{
			name:      "empty recipient should fail",
			recipient: "",
			summary: ports.EmailSummary{
				TotalBalance: 100.0,
			},
			wantErr: true,
		},
		{
			name:      "zero values should format correctly",
			recipient: "test@example.com",
			summary: ports.EmailSummary{
				TotalBalance:             0.0,
				MonthlyTransactionCounts: map[string]int{},
				AverageCreditAmount:      0.0,
				AverageDebitAmount:       0.0,
			},
			wantErr: false,
			expectedInBody: []string{
				"$0.00",
				"$-0.00",
				"Transaction Summary",
			},
		},
		{
			name:      "negative balance should format correctly",
			recipient: "user@test.com",
			summary: ports.EmailSummary{
				TotalBalance:             -500.75,
				MonthlyTransactionCounts: map[string]int{"March": 3, "April": 7},
				AverageCreditAmount:      200.00,
				AverageDebitAmount:       150.25,
			},
			wantErr: false,
			expectedInBody: []string{
				"$-500.75",
				"March",
				"3",
				"April",
				"7",
				"$200.00",
				"$-150.25",
			},
		},
		{
			name:      "single month transaction",
			recipient: "single@example.com",
			summary: ports.EmailSummary{
				TotalBalance:             999.99,
				MonthlyTransactionCounts: map[string]int{"December": 1},
				AverageCreditAmount:      999.99,
				AverageDebitAmount:       0.0,
			},
			wantErr: false,
			expectedInBody: []string{
				"$999.99",
				"December",
				"1",
				"$999.99",
				"$-0.00",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.wantErr {
				// Test error cases by attempting to send and expecting failure
				mockDialer := mocks.NewMockMailDialer(ctrl)
				mockFactory := mocks.NewMockMailMessageFactory(ctrl)

				client := NewSMTPEmailSenderWithDependencies(
					mockDialer,
					mockFactory,
					"sender@example.com",
				)

				err := client.SendSummaryEmail(tt.recipient, tt.summary)
				if err == nil {
					t.Errorf("expected error for empty recipient, got nil")
				}
				if !strings.Contains(err.Error(), "recipient cannot be empty") {
					t.Errorf("expected error message about empty recipient, got: %v", err)
				}
				return
			}

			// Test successful email body generation by capturing the body content
			mockDialer := mocks.NewMockMailDialer(ctrl)
			mockFactory := mocks.NewMockMailMessageFactory(ctrl)
			mockMsg := mocks.NewMockMailMessage(ctrl)

			var capturedBody string

			// Set up expectations for successful email generation
			mockFactory.EXPECT().NewMessage().Return(mockMsg).Times(1)

			// Capture headers
			mockMsg.EXPECT().SetHeader("From", "sender@example.com").Times(1)
			mockMsg.EXPECT().SetHeader("To", tt.recipient).Times(1)
			mockMsg.EXPECT().SetHeader("Subject", "Transaction Summary").Times(1)

			// Capture the email body
			mockMsg.EXPECT().SetBody("text/html", gomock.Any()).
				Do(func(contentType, body string, settings ...interface{}) {
					capturedBody = body
				}).Times(1)

			// Mock successful sending
			mockDialer.EXPECT().DialAndSend(gomock.Any()).Return(nil).Times(1)

			client := NewSMTPEmailSenderWithDependencies(
				mockDialer,
				mockFactory,
				"sender@example.com",
			)

			err := client.SendSummaryEmail(tt.recipient, tt.summary)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify the captured body contains expected content
			for _, expected := range tt.expectedInBody {
				if !strings.Contains(capturedBody, expected) {
					t.Errorf("expected '%s' to be in email body, but it wasn't found.\nActual body:\n%s", expected, capturedBody)
				}
			}

			// Verify the captured body doesn't contain unexpected content
			for _, notExpected := range tt.notExpectedInBody {
				if strings.Contains(capturedBody, notExpected) {
					t.Errorf("did not expect '%s' to be in email body, but it was found.\nActual body:\n%s", notExpected, capturedBody)
				}
			}

			// Additional validations for HTML structure
			if !strings.Contains(capturedBody, "<!DOCTYPE html>") {
				t.Error("email body should contain HTML DOCTYPE declaration")
			}
			if !strings.Contains(capturedBody, "<html>") || !strings.Contains(capturedBody, "</html>") {
				t.Error("email body should be properly formatted HTML")
			}
			if !strings.Contains(capturedBody, "<table>") {
				t.Error("email body should contain HTML tables for transaction data")
			}
		})
	}
}

func TestNewSMTPEmailSender(t *testing.T) {
	conf := SMTPConfiguration{
		Sender:     "test@example.com",
		Password:   "password",
		SmtpServer: "smtp.example.com",
		SmtpPort:   587,
	}

	client := NewSMTPEmailSender(conf)

	if client == nil {
		t.Error("expected client to be created, got nil")
	}
}
