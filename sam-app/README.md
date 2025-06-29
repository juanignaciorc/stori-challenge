# Stori Transaction Processor

This application processes transaction files from an S3 bucket, calculates account statistics, and sends a summary email to the user. It also stores transaction and account data in DynamoDB tables.

## Architecture

The application follows a hexagonal architecture pattern:

```
.
├── README.md                                <-- This instructions file
├── transactions.csv                         <-- Sample transaction file
├── transaction-processor                   <-- Source code for the lambda function
│   ├── main.go                              <-- Lambda handler
│   ├── main_test.go                         <-- Unit tests
│   ├── Dockerfile                           <-- Dockerfile
│   └── internal/                            <-- Internal packages
│       ├── domain/                          <-- Domain models
│       │   └── model/                       <-- Core domain models
│       │       ├── transaction.go           <-- Transaction model
│       │       └── account.go               <-- Account model
│       ├── ports/                           <-- Interface definitions
│       │   ├── file_reader.go               <-- File reading interface
│       │   ├── email_sender.go              <-- Email sending interface
│       │   └── transaction_repository.go    <-- Database interface
│       ├── adapters/                        <-- Interface implementations
│       │   ├── csv_file_reader.go           <-- CSV file reader
│       │   ├── s3_file_reader.go            <-- S3 file reader
│       │   ├── ses_email_sender.go          <-- SES email sender
│       │   └── dynamodb_repository.go       <-- DynamoDB repository
│       └── application/                     <-- Application services
│           └── transaction_service.go       <-- Transaction processing service
└── template.yaml                            <-- SAM template
```

## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
* [Golang](https://golang.org) (version 1.19 or later)
* An AWS account with permissions to create S3 buckets, DynamoDB tables, Lambda functions, and send emails via SES

## Setup process

### Installing dependencies & building the target 

This project uses the built-in `sam build` to build a Docker image from a Dockerfile and then copy the source of your application inside the Docker image.

```bash
sam build
```

### Local development

Before testing locally, you need to create a local transactions.csv file:

```bash
# The file is already created in the project root
cat transactions.csv
```

**Invoking function locally through local API Gateway**

```bash
sam local start-api
```

If the previous command ran successfully you should now be able to hit the following local endpoint to invoke your function:

```
http://localhost:3000/process-transactions
```

**Testing with local environment variables**

You can also test the function locally with environment variables:

```bash
sam local invoke HelloWorldFunction --env-vars env.json
```

Where env.json contains:

```json
{
  "HelloWorldFunction": {
    "S3_BUCKET": "local-bucket",
    "S3_KEY": "transactions.csv",
    "EMAIL_RECIPIENT": "your-email@example.com",
    "EMAIL_SENDER": "noreply@example.com",
    "TRANSACTIONS_TABLE": "Transactions",
    "ACCOUNTS_TABLE": "Accounts",
    "ACCOUNT_ID": "default"
  }
}
```

Note: When testing locally, the S3 file reader will not work. You can modify the code to use the CSV file reader instead for local testing.

## Packaging and deployment

This application uses Docker for packaging. The SAM template is configured to build a Docker image from the Dockerfile in the transaction-processor directory.

To deploy your application for the first time, run the following in your shell:

```bash
sam deploy --guided
```

The command will package and deploy your application to AWS, with a series of prompts:

* **Stack Name**: The name of the stack to deploy to CloudFormation. This should be unique to your account and region, and a good starting point would be "stori-transaction-processor".
* **AWS Region**: The AWS region you want to deploy your app to.
* **Parameter EmailRecipient**: The email address to send transaction summaries to.
* **Parameter EmailSender**: The email address to send transaction summaries from. This must be verified in SES.
* **Confirm changes before deploy**: If set to yes, any change sets will be shown to you before execution for manual review.
* **Allow SAM CLI IAM role creation**: Set to yes to allow SAM to create the necessary IAM roles.
* **Save arguments to samconfig.toml**: If set to yes, your choices will be saved to a configuration file inside the project.

You can find your API Gateway Endpoint URL in the output values displayed after deployment.

### Verifying SES Email Addresses

Before you can send emails, you need to verify the email addresses in SES:

1. Go to the AWS SES console
2. Click on "Email Addresses" under "Identity Management"
3. Click "Verify a New Email Address"
4. Enter the email address you specified as the EmailSender parameter
5. Click "Verify This Email Address"
6. Check your email and click the verification link

### Uploading a Transaction File

After deployment, you can upload a transaction file to the S3 bucket:

```bash
# Get the bucket name from the CloudFormation outputs
BUCKET_NAME=$(aws cloudformation describe-stacks --stack-name stori-transaction-processor --query "Stacks[0].Outputs[?OutputKey=='TransactionsBucketName'].OutputValue" --output text)

# Upload the sample transactions.csv file
aws s3 cp transactions.csv s3://$BUCKET_NAME/

# Alternatively, you can use the AWS Management Console to upload the file
```

The Lambda function will be triggered automatically when the file is uploaded, and a summary email will be sent to the specified recipient.

### Testing

You can run the unit tests for the Lambda function using the Go testing package:

```shell
cd ./transaction-processor/
go test -v .
```

To test the deployed application, you can:

1. Upload a transaction file to the S3 bucket as described above
2. Check your email for the summary
3. Verify that the transactions and account information are stored in the DynamoDB tables:

```bash
# List transactions in the DynamoDB table
aws dynamodb scan --table-name Transactions

# List account information in the DynamoDB table
aws dynamodb scan --table-name Accounts
```

You can also invoke the Lambda function directly through the API Gateway endpoint:

```bash
# Get the API Gateway endpoint URL
API_URL=$(aws cloudformation describe-stacks --stack-name stori-transaction-processor --query "Stacks[0].Outputs[?OutputKey=='TransactionProcessorAPI'].OutputValue" --output text)

# Invoke the API
curl -X GET $API_URL
```

## Features

This application implements the following features:

1. **Transaction Processing**: Reads transaction data from a CSV file in an S3 bucket
2. **Account Summary**: Calculates total balance, transactions per month, and average credit/debit amounts
3. **Email Notification**: Sends a summary email to the user with the account information
4. **Database Storage**: Stores transaction and account information in DynamoDB tables
5. **Styled Email**: Includes HTML formatting for better readability

## Hexagonal Architecture

The application follows a hexagonal architecture pattern, which provides several benefits:

1. **Separation of Concerns**: The core domain logic is separated from external dependencies
2. **Testability**: The core domain logic can be tested without external dependencies
3. **Flexibility**: External dependencies can be swapped out without changing the core domain logic
4. **Maintainability**: The codebase is organized in a way that makes it easy to understand and maintain

## Future Improvements

Here are some ideas for future improvements:

1. Add more comprehensive error handling and logging
2. Implement pagination for large transaction files
3. Add authentication and authorization for the API
4. Create a web interface for uploading transaction files and viewing account information
5. Implement real-time notifications using WebSockets or SNS
6. Add support for different file formats (JSON, XML, etc.)
7. Implement data validation and sanitization
8. Add support for multiple accounts and users
