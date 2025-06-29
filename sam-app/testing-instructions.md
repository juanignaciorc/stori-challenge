# Testing Your Transaction Processor Lambda Function

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
