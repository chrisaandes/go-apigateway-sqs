package createOrder

import (
	"context"
	"testing"
)

func TestCreateOrderHandler(t *testing.T) {
	// Create a new CreateOrderRequest for testing
	request := CreateOrderRequest{
		UserID:     "test-user-id",
		Item:       "test-item",
		Quantity:   2,
		TotalPrice: 100,
	}

	// Call the CreateOrderHandler function with the test request
	orderID, err := CreateOrderHandler(context.Background(), request)

	// Verify that there were no errors returned by the function
	if err != nil {
		t.Errorf("Error creating new order: %v", err)
	}

	expectedAttributes := map[string]interface{}{
		"order_id":    orderID,
		"user_id":     request.UserID,
		"item":        request.Item,
		"quantity":    request.Quantity,
		"total_price": request.TotalPrice,
	}
	if !orderExistsWithAttributes(orderID, expectedAttributes) {
		t.Errorf("New order not found in DynamoDB or has incorrect attributes")
	}
}

func orderExistsWithAttributes(orderID string, expectedAttributes map[string]interface{}) bool {
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String("http://localhost:4566"),
		Region:   aws.String("us-east-1"),
	})
	if err != nil {
		log.Printf("Error creating new AWS session: %v", err)
		return false
	}

	// Create a new DynamoDB client
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"order_id": {
				S: aws.String(orderID),
			},
		},
		TableName: aws.String("orders"),
	})
	if err != nil {
		log.Printf("Error querying DynamoDB: %v", err)
		return false
	}

	if len(result.Item) == 0 {
		return false
	}
	for key, value := range expectedAttributes {
		attributeValue, ok := result.Item[key]
		if !ok {
			return false
		}
		switch value.(type) {
		case string:
			if *attributeValue.S != value.(string) {
				return false
			}
		case int:
			if int(*attributeValue.N) != value.(int) {
				return false
			}
		case int64:
			if int64(*attributeValue.N) != value.(int64) {
				return false
			}
		default:
			return false
		}
	}
	return true
}
