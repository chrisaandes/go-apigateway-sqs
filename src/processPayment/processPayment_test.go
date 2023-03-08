package processPayment

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	requestBody := `{"order_id": "test-order-id", "status": "paid"}`
	req, err := http.NewRequest("POST", "/process-payment", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	Handler(context.Background(), events.APIGatewayProxyRequest{
		Body: requestBody,
	}, rr)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, rr.Code)
	}

	expectedBody := "Payment processed successfully"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected response body %s but got %s", expectedBody, rr.Body.String())
	}

	if !messageExistsInQueue(sqsURL, requestBody) {
		t.Errorf("Message not found in SQS queue or has incorrect content")
	}
}

func messageExistsInQueue(queueURL string, message string) bool {
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(os.Getenv("LOCALSTACK_HOSTNAME")),
		Region:   aws.String("us-east-1"),
	})
	if err != nil {
		log.Printf("Error creating new AWS session: %v", err)
		return false
	}

	svc := sqs.New(sess)

	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queueURL),
	})
	if err != nil {
		log.Printf("Error receiving messages from SQS: %v", err)
		return false
	}
	for _, message := range result.Messages {
		if *message.Body == message {
			return true
		}
	}
	return false
}
