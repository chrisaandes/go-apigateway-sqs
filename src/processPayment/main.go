package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chrisaandes/go-apigateway-sqs/src/processPayment"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req processPayment.ProcessPaymentRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Printf("Error unmarshaling request: %s", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Error unmarshaling request",
		}, nil
	}

	err = processPayment.ProcessPaymentHandler(ctx, req)
	if err != nil {
		log.Printf("Error processing payment: %s", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error processing payment",
		}, nil
	}

	// Return a successful response
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Payment processed successfully",
	}, nil
}

func main() {
	lambda.Start(handler)
}
