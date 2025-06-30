# Stori Transaction Processor

A serverless application that processes transaction files, calculates account statistics, and sends summary emails. Built with AWS Lambda, API Gateway, S3, DynamoDB, and SMTP.

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  API Gateway │────▶│    Lambda   │────▶│     S3      │     │    SMTP     │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
                          │                                        ▲
                          │                                        │
                          ▼                                        │
                    ┌─────────────┐                         ┌─────────────┐
                    │   DynamoDB  │                         │  Email      │
                    └─────────────┘                         │  Summary    │
                                                            └─────────────┘
```

This project implements **Hexagonal Architecture** (ports and adapters) to separate core business logic from external dependencies, making the code more maintainable, testable, and flexible.

## Quick Start

### Test the Deployed Lambda

```bash
curl -X POST https://7a96wdh36m.execute-api.us-east-1.amazonaws.com/Prod/process-transactions/ \
  -H "Content-Type: application/json" \
  -d '{"email": "your-email-here@gmail.com"}'
```

### Local Development

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
