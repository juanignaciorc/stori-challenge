package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"os"
	"strconv"
	"transaction-processor/internal/adapters"
	"transaction-processor/internal/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// Configuration holds the Lambda function configuration
type Configuration struct {
	EmailSender       string `json:"emailSender"`
	EmailPassword     string `json:"emailPassword"`
	SmtpServer        string `json:"smtpServer"`
	SmtpPort          int    `json:"smtpPort"`
	TransactionsTable string `json:"transactionsTable"`
	AccountsTable     string `json:"accountsTable"`
	AccountID         string `json:"accountID"`
}

// RequestBody represents the expected structure of the POST request body
type RequestBody struct {
	Email string `json:"email" validate:"required,email"`
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

	// Validate email format using validator
	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
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

	dynamoClient := dynamodb.NewFromConfig(awsConfig)

	// Create adapters
	fileReader := adapters.NewCSVFileReader()

	// Create SMTP email sender
	smtpConfig := adapters.SMTPConfiguration{
		Sender:     config.EmailSender,
		Password:   config.EmailPassword,
		SmtpServer: config.SmtpServer,
		SmtpPort:   config.SmtpPort,
	}
	emailSender := adapters.NewSMTPEmailSender(smtpConfig)

	log.Printf("Using SMTP email sender with server: %s, port: %d", config.SmtpServer, config.SmtpPort)

	var repository *adapters.DynamoDBRepository
	if config.TransactionsTable != "" && config.AccountsTable != "" {
		repository = adapters.NewDynamoDBRepository(dynamoClient, config.TransactionsTable, config.AccountsTable)
	}

	// Create trx service
	service := services.NewTransactionService(fileReader, emailSender, repository)

	// Process transactions and send summary
	filePath := "transactions.csv"
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
	// Get SMTP port from environment variable, default to 587 if not set
	smtpPort := 587
	if portStr := os.Getenv("SMTP_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			smtpPort = port
		}
	}

	config := Configuration{
		EmailSender:       os.Getenv("EMAIL_SENDER"),
		EmailPassword:     os.Getenv("EMAIL_PASSWORD"),
		SmtpServer:        os.Getenv("SMTP_SERVER"),
		SmtpPort:          smtpPort,
		TransactionsTable: os.Getenv("TRANSACTIONS_TABLE"),
		AccountsTable:     os.Getenv("ACCOUNTS_TABLE"),
		AccountID:         os.Getenv("ACCOUNT_ID"),
	}

	// If ACCOUNT_ID is not set, use a default value
	if config.AccountID == "" {
		config.AccountID = "default"
	}

	// Set default SMTP server if not provided
	if config.SmtpServer == "" {
		config.SmtpServer = "smtp.gmail.com"
	}

	return config
}

func main() {
	lambda.Start(handler)
}
