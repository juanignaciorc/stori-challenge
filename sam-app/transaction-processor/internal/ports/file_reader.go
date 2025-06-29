package ports

import (
	"transaction-processor/internal/domain/model"
)

// FileReader defines the interface for reading transaction data from a file
type FileReader interface {
	// ReadTransactions reads transactions from a file and returns them
	ReadTransactions(filePath string) ([]*model.Transaction, error)
}
