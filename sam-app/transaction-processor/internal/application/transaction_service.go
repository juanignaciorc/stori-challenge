package application

import (
	"transaction-processor/internal/domain/model"
	"transaction-processor/internal/ports"
)

// TransactionService orchestrates the transaction processing use case
type TransactionService struct {
	fileReader            ports.FileReader
	emailSender           ports.EmailSender
	transactionRepository ports.TransactionRepository
}

// NewTransactionService creates a new TransactionService
func NewTransactionService(
	fileReader ports.FileReader,
	emailSender ports.EmailSender,
	transactionRepository ports.TransactionRepository,
) *TransactionService {
	return &TransactionService{
		fileReader:            fileReader,
		emailSender:           emailSender,
		transactionRepository: transactionRepository,
	}
}

// ProcessTransactionsAndSendSummary processes a transaction file and sends a summary email
func (s *TransactionService) ProcessTransactionsAndSendSummary(filePath, emailRecipient, accountID string) error {
	// Read transactions from file
	transactions, err := s.fileReader.ReadTransactions(filePath)
	if err != nil {
		return err
	}

	// Create account and add transactions
	account := model.NewAccount()
	for _, tx := range transactions {
		account.AddTransaction(tx)

		// Save transaction to database if repository is provided
		if s.transactionRepository != nil {
			if err := s.transactionRepository.SaveTransaction(tx); err != nil {
				return err
			}
		}
	}

	// Create email summary
	summary := ports.NewEmailSummaryFromAccount(account)

	// Save account summary to database if repository is provided
	if s.transactionRepository != nil {
		if err := s.transactionRepository.SaveAccount(accountID, summary); err != nil {
			return err
		}
	}

	// Send summary email
	return s.emailSender.SendSummaryEmail(emailRecipient, summary)
}
