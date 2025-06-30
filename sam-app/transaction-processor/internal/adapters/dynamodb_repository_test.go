package adapters

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
	"transaction-processor/internal/domain/model"
	"transaction-processor/internal/mocks"
	"transaction-processor/internal/ports"
)

func TestDynamoDBRepository_SaveTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamo := mocks.NewMockDynamoDBClient(ctrl)

	repo := &DynamoDBRepository{
		dynamoClient:      mockDynamo,
		transactionsTable: "TransactionsTable",
		accountsTable:     "AccountsTable",
	}

	tx := &model.Transaction{
		ID:       "tx123",
		Date:     time.Date(2025, 6, 30, 12, 0, 0, 0, time.UTC),
		Amount:   123.45,
		IsCredit: true,
	}

	mockDynamo.
		EXPECT().
		PutItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.PutItemOutput{}, nil)

	err := ports.TransactionRepository.SaveTransaction(repo, tx)

	assert.NoError(t, err)
}

func TestDynamoDBRepository_SaveAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamo := mocks.NewMockDynamoDBClient(ctrl)

	repo := &DynamoDBRepository{
		dynamoClient:      mockDynamo,
		transactionsTable: "TransactionsTable",
		accountsTable:     "AccountsTable",
	}

	accountID := "account123"
	summary := ports.EmailSummary{
		TotalBalance:        1000.50,
		AverageCreditAmount: 300.75,
		AverageDebitAmount:  200.25,
		MonthlyTransactionCounts: map[string]int{
			"2025-01": 5,
			"2025-02": 8,
		},
	}

	mockDynamo.
		EXPECT().
		PutItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.PutItemOutput{}, nil)

	err := repo.SaveAccount(accountID, summary)

	assert.NoError(t, err)
}

func TestDynamoDBRepository_GetTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamo := mocks.NewMockDynamoDBClient(ctrl)

	repo := &DynamoDBRepository{
		dynamoClient:      mockDynamo,
		transactionsTable: "TransactionsTable",
		accountsTable:     "AccountsTable",
	}

	accountID := "account123"

	// Mock DynamoDB response
	now := time.Now().UTC().Truncate(time.Second)
	mockDynamo.EXPECT().
		Scan(gomock.Any(), gomock.Any()).
		Return(&dynamodb.ScanOutput{
			Items: []map[string]types.AttributeValue{
				{
					"ID":       &types.AttributeValueMemberS{Value: "tx1"},
					"Date":     &types.AttributeValueMemberS{Value: now.Format(time.RFC3339)},
					"Amount":   &types.AttributeValueMemberN{Value: "123.45"},
					"IsCredit": &types.AttributeValueMemberBOOL{Value: true},
				},
			},
		}, nil)

	txs, err := repo.GetTransactions(accountID)

	assert.NoError(t, err)
	assert.Len(t, txs, 1)
	assert.Equal(t, "tx1", txs[0].ID)
	assert.Equal(t, 123.45, txs[0].Amount)
	assert.Equal(t, true, txs[0].IsCredit)
	assert.Equal(t, now, txs[0].Date)
}
