package processPayment

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type ProcessPaymentRequest struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var processPaymentRequest ProcessPaymentRequest
	err := json.Unmarshal([]byte(request.Body), &processPaymentRequest)
	if err != nil {
		log.Printf("Error al parsear el request body: %s", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Error al parsear el request body",
		}, nil
	}

	sqsURL := os.Getenv("SQS_URL")

	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(os.Getenv("LOCALSTACK_HOSTNAME")),
		Region:   aws.String("us-east-1"),
	})
	if err != nil {
		log.Printf("Error al crear la sesión de AWS: %s", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error al crear la sesión de AWS",
		}, nil
	}

	svc := sqs.New(sess)

	message := fmt.Sprintf(`{"order_id":"%s","status":"%s"}`, processPaymentRequest.OrderID, processPaymentRequest.Status)

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String(sqsURL),
	})
	if err != nil {
		log.Printf("Error al enviar el mensaje a la cola SQS: %s", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error al enviar el mensaje a la cola SQS",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Payment processed successfully",
	}, nil
}
