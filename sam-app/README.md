# Stori Transaction Processor

A serverless application that processes transaction files, calculates account statistics, and sends summary emails. Built with AWS Lambda, API Gateway, S3, DynamoDB, and SMTP.

## Architecture

```
┌─────────────┐      ┌───────────────────────── ┐     ┌─────────────┐     
│  API Gateway │────▶│    Lambda -              │────▶│      S3      │   
└─────────────┘      └──────────────────────────┘     └─────────────┘     
                          │                  │                                    
                          │                  │                                    
                          ▼                  ▼                                    
                       ┌─────────────┐┌─────────────┐                          
                       │   DynamoDB  ││ SMTP Server │                    
                       └─────────────┘└─────────────┘                         

```

This project implements **Hexagonal Architecture** (ports and adapters) to separate core business logic from external dependencies, making the code more maintainable, testable, and flexible.

## Quick Start

# Testing Transaction Processor Lambda Function

Based on your deployed Lambda function, here are several ways to test it:

## 1. Testing via API Gateway Endpoint

The simplest way to test your Lambda function is by making a POST request to the API Gateway endpoint that was provided in the deployment outputs:

```bash
curl -X POST https://7a96wdh36m.execute-api.us-east-1.amazonaws.com/Prod/process-transactions/ \
  -H "Content-Type: application/json" \
  -d '{"email": "juanignacioroldan01@gmail.com"}'
```

Replace `juanignacioroldan01@gmail.com` with your actual email address where you want to receive the transaction summary.

The function will process the transactions from the CSV file and send a summary email to the provided email address. If successful, it will return a message indicating that the transactions were processed successfully.

### Local Development

**⚠️ DISCLAIMER**: The application will fail when running locally because it requires a local DynamoDB instance to be running. The app is designed to work with AWS DynamoDB and does not include local DynamoDB setup. For full functionality testing, please deploy to AWS or set up DynamoDB Local separately.

1. Build the application:
   ```bash
   sam build
   ```

2. Start local API:
   ```bash
   sam local start-api
   ```

3. Test locally:
   ```bash
   curl -X POST http://localhost:3000/process-transactions \
     -H "Content-Type: application/json" \
     -d '{"email": "your-email@example.com"}'
   ```
### Deployment

```bash
sam deploy --guided
```

## Project Structure

The codebase follows hexagonal architecture with:
- `domain/model`: Core business entities
- `ports`: Interface definitions
- `adapters`: External implementations
- `services`: Application business logic

## Features

- Transaction processing from CSV files
- Account summary calculation
- Email notifications with summary
- DynamoDB storage for transactions and accounts
