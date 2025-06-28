# Testing Your Lambda Function

Based on your deployed Lambda function, here are several ways to test it:

## 1. Testing via API Gateway Endpoint

The simplest way to test your Lambda function is by making a GET request to the API Gateway endpoint that was provided in the deployment outputs:

```bash
curl https://7a96wdh36m.execute-api.us-east-1.amazonaws.com/Prod/hello-world/
```

You can also open this URL in a web browser:
https://7a96wdh36m.execute-api.us-east-1.amazonaws.com/Prod/hello-world/

The function should return a greeting message that includes your IP address.