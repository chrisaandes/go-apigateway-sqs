package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/chrisaandes/go-apigateway-sqs/src/createOrder"
)

type CreateOrderEvent struct {
	OrderID    string `json:"order_id"`
	TotalPrice int64  `json:"total_price"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req createOrder.CreateOrderRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Printf("Error unmarshaling request: %s", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Error unmarshaling request",
		}, nil
	}

	orderID, err := createOrder.CreateOrderHandler(ctx, req)
	if err != nil {
		log.Printf("Error creating order: %s", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error creating order",
		}, nil
	}

	// Send a payment event via SQS or EventBridge
	err = sendPaymentEvent(orderID, req.TotalPrice)
	if err != nil {
		log.Printf("Error sending payment event: %s", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error sending payment event",
		}, nil
	}

	// Return a successful response with the order ID
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       orderID,
	}, nil
}

func sendPaymentEvent(orderID string, totalPrice int64) error {
	// Get the SQS URL from an environment variable
	sqsURL := os.Getenv("SQS_URL")

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(os.Getenv("LOCALSTACK_HOSTNAME")),
		Region:   aws.String("us-east-1"),
	})
	if err != nil {
		return err
	}

	// Create a new SQS client
	svc := sqs.New(sess)

	// Create a new payment event
	paymentEvent := CreateOrderEvent{
		OrderID:    orderID,
		TotalPrice: totalPrice,
	}

	// Convert the payment event to JSON
	paymentEventJSON, err := json.Marshal(paymentEvent)
	if err != nil {
		return err
	}

	// Send the payment event to the SQS queue
	_, err = svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(paymentEventJSON)),
		QueueUrl:    aws.String(sqsURL),
	})
	if err != nil {
		return err
	}

	log.Printf("Payment event sent for order ID: %s", orderID)
	return nil
}

func main() {
	lambda.Start(handler)
}
