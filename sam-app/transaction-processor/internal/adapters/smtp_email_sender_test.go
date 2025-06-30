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
	// Test the email body generation separately
	tests := []struct {
		name              string
		recipient         string
		summary           ports.EmailSummary
		wantErr           bool
		expectedInBody    []string
		notExpectedInBody []string
	}{
		{
			name:      "valid email generation",
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
				"February",
				"$100.25",
				"$75.50",
				"Transaction Summary",
			},
		},
		{
			name:      "empty recipient",
			recipient: "",
			summary: ports.EmailSummary{
				TotalBalance: 100.0,
			},
			wantErr: true,
		},
		{
			name:      "zero values",
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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For now, we'll test via the public method but intercept before sending
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.wantErr {
				// Test by calling SendSummaryEmail and expecting the body generation error
				mockDialer := mocks.NewMockMailDialer(ctrl)
				mockFactory := mocks.NewMockMailMessageFactory(ctrl)

				client := NewSMTPEmailSenderWithDependencies(
					mockDialer,
					mockFactory,
					"sender@example.com",
				)

				err := client.SendSummaryEmail(tt.recipient, tt.summary)
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			// For successful cases, we'll test by mocking and capturing the body
			mockDialer := mocks.NewMockMailDialer(ctrl)
			mockFactory := mocks.NewMockMailMessageFactory(ctrl)
			mockMsg := mocks.NewMockMailMessage(ctrl)

			var capturedBody string

			mockFactory.EXPECT().NewMessage().Return(mockMsg).Times(1)
			mockMsg.EXPECT().SetHeader(gomock.Any(), gomock.Any()).AnyTimes()
			mockMsg.EXPECT().SetBody("text/html", gomock.Any()).
				Do(func(contentType, body string, settings ...interface{}) {
					capturedBody = body
				}).Times(1)
			mockDialer.EXPECT().DialAndSend(gomock.Any()).Return(nil).Times(1)

			client := NewSMTPEmailSenderWithDependencies(
				mockDialer,
				mockFactory,
				"sender@example.com",
			)

			err := client.SendSummaryEmail(tt.recipient, tt.summary)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Check expected content in body
			for _, expected := range tt.expectedInBody {
				if !strings.Contains(capturedBody, expected) {
					t.Errorf("expected '%s' to be in email body, but it wasn't found", expected)
				}
			}

			// Check content that should not be in body
			for _, notExpected := range tt.notExpectedInBody {
				if strings.Contains(capturedBody, notExpected) {
					t.Errorf("did not expect '%s' to be in email body, but it was found", notExpected)
				}
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
