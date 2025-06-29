package ports

import (
	"transaction-processor/internal/domain/model"
)

// TransactionRepository defines the interface for storing and retrieving transactions
type TransactionRepository interface {
	// SaveTransaction saves a transaction to the database
	SaveTransaction(tx *model.Transaction) error

	// SaveAccount saves account information to the database
	SaveAccount(accountID string, summary EmailSummary) error

	// GetTransactions retrieves all transactions for an account
	GetTransactions(accountID string) ([]*model.Transaction, error)
}
