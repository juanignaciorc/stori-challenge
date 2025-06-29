package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"transaction-processor/internal/adapters"
	"transaction-processor/internal/application"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ses"
)

// Configuration holds the Lambda function configuration
type Configuration struct {
	EmailSender       string `json:"emailSender"`
	TransactionsTable string `json:"transactionsTable"`
	AccountsTable     string `json:"accountsTable"`
	AccountID         string `json:"accountID"`
}

// RequestBody represents the expected structure of the POST request body
type RequestBody struct {
	Email string `json:"email"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse and validate email from request body
	var requestBody RequestBody
	if err := json.Unmarshal([]byte(request.Body), &requestBody); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Invalid request body: %v", err),
		}, nil
	}

	// Validate email format
	if _, err := mail.ParseAddress(requestBody.Email); err != nil {
		log.Printf("Invalid email format: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid email format. Please provide a valid email address.",
		}, nil
	}

	// Load configuration from environment variables
	config := loadConfiguration()

	// Initialize AWS SDK clients
	awsConfig, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Error loading AWS config: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error loading AWS config: %v", err),
		}, nil
	}

	sesClient := ses.NewFromConfig(awsConfig)
	dynamoClient := dynamodb.NewFromConfig(awsConfig)

	// Create adapters
	fileReader := adapters.NewCSVFileReader()
	emailSender := adapters.NewSESEmailSender(sesClient, config.EmailSender)
	var repository *adapters.DynamoDBRepository
	if config.TransactionsTable != "" && config.AccountsTable != "" {
		repository = adapters.NewDynamoDBRepository(dynamoClient, config.TransactionsTable, config.AccountsTable)
	}

	// Create application service
	service := application.NewTransactionService(fileReader, emailSender, repository)

	// Process transactions and send summary
	filePath := filepath.Join(".", "transactions.csv")
	err = service.ProcessTransactionsAndSendSummary(filePath, requestBody.Email, config.AccountID)
	if err != nil {
		log.Printf("Error processing transactions: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error processing transactions: %v", err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Transactions processed successfully. Summary email sent.",
	}, nil
}

func loadConfiguration() Configuration {
	config := Configuration{
		EmailSender:       os.Getenv("EMAIL_SENDER"),
		TransactionsTable: os.Getenv("TRANSACTIONS_TABLE"),
		AccountsTable:     os.Getenv("ACCOUNTS_TABLE"),
		AccountID:         os.Getenv("ACCOUNT_ID"),
	}

	// If ACCOUNT_ID is not set, use a default value
	if config.AccountID == "" {
		config.AccountID = "default"
	}

	return config
}

func main() {
	lambda.Start(handler)
}
