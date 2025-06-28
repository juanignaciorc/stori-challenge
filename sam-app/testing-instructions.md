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

## 2. Testing via AWS CLI

### Invoke the Lambda function directly

```bash
aws lambda invoke \
  --function-name arn:aws:lambda:us-east-1:458093367779:function:sam-app-HelloWorldFunction-VyalFgXnhbtG \
  --payload '{"requestContext": {"identity": {"sourceIP": "test-ip"}}}' \
  response.json

# View the response
cat response.json
```

### Test the API Gateway endpoint

```bash
aws apigateway test-invoke-method \
  --rest-api-id 7a96wdh36m \
  --resource-id <resource-id> \
  --http-method GET \
  --path-with-query-string "/hello-world"
```

Note: You'll need to find the resource-id from the API Gateway console.

## 3. Testing via AWS Console

### Lambda Console
1. Go to the [AWS Lambda Console](https://console.aws.amazon.com/lambda)
2. Find and select your function: `sam-app-HelloWorldFunction-VyalFgXnhbtG`
3. Click on the "Test" tab
4. Create a new test event with the following JSON:
   ```json
   {
     "requestContext": {
       "identity": {
         "sourceIP": "test-ip"
       }
     }
   }
   ```
5. Click "Test" to invoke the function with this event

### API Gateway Console
1. Go to the [API Gateway Console](https://console.aws.amazon.com/apigateway)
2. Find and select your API: It should be associated with the ID `7a96wdh36m`
3. Click on "Resources" in the left navigation
4. Select the GET method under the /hello-world resource
5. Click the "Test" tab in the Method Execution pane
6. Click the "Test" button to send a test request

## 4. Local Testing

### Using SAM CLI

You can test your function locally using the SAM CLI:

```bash
# Start a local API Gateway
sam local start-api

# Then in another terminal, make a request to the local endpoint
curl http://localhost:3000/hello-world
```

### Using Go test

You can run the unit tests for your function:

```bash
cd hello-world
go test -v .
```

## 5. Monitoring and Logs

After testing your function, you can check the logs:

```bash
# Get the recent logs for your function
aws logs get-log-events \
  --log-group-name /aws/lambda/sam-app-HelloWorldFunction-VyalFgXnhbtG \
  --log-stream-name <log-stream-name>
```

You can also view logs in the CloudWatch Logs console.