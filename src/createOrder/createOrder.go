package createOrder

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	UserID     string `json:"user_id"`
	Item       string `json:"item"`
	Quantity   int    `json:"quantity"`
	TotalPrice int64  `json:"total_price"`
}

func CreateOrderHandler(ctx context.Context, request CreateOrderRequest) (string, error) {
	// Generate a random order ID
	orderID := uuid.New().String()

	// Create a new order item for DynamoDB
	orderItem := map[string]*dynamodb.AttributeValue{
		"order_id": {
			S: aws.String(orderID),
		},
		"user_id": {
			S: aws.String(request.UserID),
		},
		"item": {
			S: aws.String(request.Item),
		},
		"quantity": {
			N: aws.String(string(request.Quantity)),
		},
		"total_price": {
			N: aws.String(string(request.TotalPrice)),
		},
		"created_at": {
			S: aws.String(time.Now().Format(time.RFC3339)),
		},
	}

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String("http://localhost:4566"),
		Region:   aws.String("us-east-1"),
	})
	if err != nil {
		return "", err
	}

	// Create a new DynamoDB client
	svc := dynamodb.New(sess)

	// Put the new order item in DynamoDB
	_, err = svc.PutItem(&dynamodb.PutItemInput{
		Item:      orderItem,
		TableName: aws.String("orders"),
	})
	if err != nil {
		log.Printf("Error putting item in DynamoDB: %s", err.Error())
		return "", err
	}

	log.Printf("New order created with ID: %s", orderID)
	return orderID, nil
}
