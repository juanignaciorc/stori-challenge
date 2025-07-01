package factory

import (
	"context"
	"log"

	"transaction-processor/internal/adapters"
	"transaction-processor/internal/config"
	"transaction-processor/internal/services"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// ServiceFactory creates and configures application services
type ServiceFactory struct {
	config config.Configuration
}

// NewServiceFactory creates a new ServiceFactory
func NewServiceFactory(cfg config.Configuration) *ServiceFactory {
	return &ServiceFactory{
		config: cfg,
	}
}

// CreateTransactionService creates a fully configured TransactionService
func (f *ServiceFactory) CreateTransactionService() (*services.TransactionService, error) {
	// Initialize AWS SDK clients
	awsConfig, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Error loading AWS config: %v", err)
		return nil, err
	}

	dynamoClient := dynamodb.NewFromConfig(awsConfig)

	// Create adapters
	fileReader := adapters.NewCSVFileReader()

	// Create SMTP email sender
	smtpConfig := adapters.SMTPConfiguration{
		Sender:     f.config.EmailSender,
		Password:   f.config.EmailPassword,
		SmtpServer: f.config.SmtpServer,
		SmtpPort:   f.config.SmtpPort,
	}
	emailSender := adapters.NewSMTPEmailSender(smtpConfig)

	log.Printf("Using SMTP email sender with server: %s, port: %d", f.config.SmtpServer, f.config.SmtpPort)

	var repository *adapters.DynamoDBRepository
	if f.config.TransactionsTable != "" && f.config.AccountsTable != "" {
		repository = adapters.NewDynamoDBRepository(dynamoClient, f.config.TransactionsTable, f.config.AccountsTable)
	}

	// Create and return transaction service
	return services.NewTransactionService(fileReader, emailSender, repository), nil
}