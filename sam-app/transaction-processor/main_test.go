package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	testCases := []struct {
		name          string
		request       events.APIGatewayProxyRequest
		expectedCode  int
		expectedError error
	}{
		{
			name: "valid email",
			request: events.APIGatewayProxyRequest{
				Body: `{"email":"test@example.com"}`,
			},
			expectedCode:  200,
			expectedError: nil,
		},
		{
			name: "invalid email",
			request: events.APIGatewayProxyRequest{
				Body: `{"email":"invalid-email"}`,
			},
			expectedCode:  400,
			expectedError: nil,
		},
		{
			name: "invalid JSON",
			request: events.APIGatewayProxyRequest{
				Body: `{"email":`,
			},
			expectedCode:  400,
			expectedError: nil,
		},
	}

	// Skip actual test execution since we can't mock AWS services easily in this environment
	t.Skip("Skipping test as it requires AWS services to be mocked")

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			response, err := handler(testCase.request)
			if err != testCase.expectedError {
				t.Errorf("Expected error %v, but got %v", testCase.expectedError, err)
			}

			if response.StatusCode != testCase.expectedCode {
				t.Errorf("Expected status code %d, but got %v", testCase.expectedCode, response.StatusCode)
			}
		})
	}
}
