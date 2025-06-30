package ports

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"transaction-processor/internal/domain/model"
)

type DynamoDBClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

// TransactionRepository defines the interface for storing and retrieving transactions
type TransactionRepository interface {
	// SaveTransaction saves a transaction to the database
	SaveTransaction(tx *model.Transaction) error

	// SaveAccount saves account information to the database
	SaveAccount(accountID string, summary EmailSummary) error

	// GetTransactions retrieves all transactions for an account
	GetTransactions(accountID string) ([]*model.Transaction, error)
}
