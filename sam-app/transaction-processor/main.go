package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"transaction-processor/internal/config"
	"transaction-processor/internal/handlers"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Load configuration from environment variables
	cfg := config.Load()

	// Create transaction handler
	transactionHandler := handlers.NewTransactionHandler(cfg)

	// Handle the request
	return transactionHandler.Handle(request)
}

func main() {
	lambda.Start(handler)
}
