package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"transaction-processor/internal/config"
	"transaction-processor/internal/factory"
	"transaction-processor/internal/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
)

// TransactionHandler handles transaction processing requests
type TransactionHandler struct {
	config         config.Configuration
	serviceFactory *factory.ServiceFactory
}

// NewTransactionHandler creates a new TransactionHandler
func NewTransactionHandler(cfg config.Configuration) *TransactionHandler {
	return &TransactionHandler{
		config:         cfg,
		serviceFactory: factory.NewServiceFactory(cfg),
	}
}

// Handle processes the Lambda request for transaction processing
func (h *TransactionHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse and validate email from request body
	var requestBody models.RequestBody
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

	// Create transaction service using factory
	service, err := h.serviceFactory.CreateTransactionService()
	if err != nil {
		log.Printf("Error creating transaction service: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error creating transaction service: %v", err),
		}, nil
	}

	// Process transactions and send summary
	filePath := "transactions.csv"
	err = service.ProcessTransactionsAndSendSummary(filePath, requestBody.Email, h.config.AccountID)
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
