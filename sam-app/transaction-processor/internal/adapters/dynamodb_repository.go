package adapters

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"transaction-processor/internal/domain/model"
	"transaction-processor/internal/ports"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDBRepository implements the TransactionRepository port using DynamoDB
type DynamoDBRepository struct {
	dynamoClient      ports.DynamoDBClient
	transactionsTable string
	accountsTable     string
}

// NewDynamoDBRepository creates a new DynamoDBRepository
func NewDynamoDBRepository(
	dynamoClient *dynamodb.Client,
	transactionsTable string,
	accountsTable string,
) *DynamoDBRepository {
	return &DynamoDBRepository{
		dynamoClient:      dynamoClient,
		transactionsTable: transactionsTable,
		accountsTable:     accountsTable,
	}
}

// SaveTransaction saves a transaction to DynamoDB
func (r *DynamoDBRepository) SaveTransaction(tx *model.Transaction) error {
	// Create the item
	item := map[string]types.AttributeValue{
		"ID":        &types.AttributeValueMemberS{Value: tx.ID},
		"Date":      &types.AttributeValueMemberS{Value: tx.Date.Format(time.RFC3339)},
		"Amount":    &types.AttributeValueMemberN{Value: strconv.FormatFloat(tx.Amount, 'f', 2, 64)},
		"IsCredit":  &types.AttributeValueMemberBOOL{Value: tx.IsCredit},
		"Timestamp": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
	}

	// Put the item in the table
	_, err := r.dynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.transactionsTable),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("error saving transaction to DynamoDB: %w", err)
	}

	return nil
}

// SaveAccount saves account information to DynamoDB
func (r *DynamoDBRepository) SaveAccount(accountID string, summary ports.EmailSummary) error {
	// Create the monthly transaction counts attribute
	monthlyCountsMap := make(map[string]types.AttributeValue)
	for month, count := range summary.MonthlyTransactionCounts {
		monthlyCountsMap[month] = &types.AttributeValueMemberN{Value: strconv.Itoa(count)}
	}

	// Create the item
	item := map[string]types.AttributeValue{
		"AccountID":           &types.AttributeValueMemberS{Value: accountID},
		"TotalBalance":        &types.AttributeValueMemberN{Value: strconv.FormatFloat(summary.TotalBalance, 'f', 2, 64)},
		"MonthlyTransactions": &types.AttributeValueMemberM{Value: monthlyCountsMap},
		"AverageCreditAmount": &types.AttributeValueMemberN{Value: strconv.FormatFloat(summary.AverageCreditAmount, 'f', 2, 64)},
		"AverageDebitAmount":  &types.AttributeValueMemberN{Value: strconv.FormatFloat(summary.AverageDebitAmount, 'f', 2, 64)},
		"Timestamp":           &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
	}

	// Put the item in the table
	_, err := r.dynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.accountsTable),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("error saving account to DynamoDB: %w", err)
	}

	return nil
}

// GetTransactions retrieves all transactions for an account from DynamoDB
func (r *DynamoDBRepository) GetTransactions(accountID string) ([]*model.Transaction, error) {
	// This is a simplified implementation that doesn't filter by account ID
	// In a real services, you would use a GSI or a query with a filter expression

	// Scan the table
	result, err := r.dynamoClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(r.transactionsTable),
	})
	if err != nil {
		return nil, fmt.Errorf("error scanning transactions table: %w", err)
	}

	// Convert the items to transactions
	var transactions []*model.Transaction
	for _, item := range result.Items {
		// Extract the values
		id := item["ID"].(*types.AttributeValueMemberS).Value
		dateStr := item["Date"].(*types.AttributeValueMemberS).Value
		amountStr := item["Amount"].(*types.AttributeValueMemberN).Value
		isCredit := item["IsCredit"].(*types.AttributeValueMemberBOOL).Value

		// Parse the date
		date, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing date: %w", err)
		}

		// Parse the amount
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing amount: %w", err)
		}

		// Create the transaction
		tx := &model.Transaction{
			ID:       id,
			Date:     date,
			Amount:   amount,
			IsCredit: isCredit,
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}
