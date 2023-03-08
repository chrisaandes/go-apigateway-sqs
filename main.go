package main

import (
    "context"
    "fmt"
    "os"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "your-app/src/createOrder"
    "your-app/src/processPayment"
)

func main() {
    lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
        switch request.Path {
        case "/create-order":
            return createOrder.Handler(ctx, request)
        case "/process-payment":
            return processPayment.Handler(ctx, request)
        default:
            return events.APIGatewayProxyResponse{
                StatusCode: 404,
                Body:       fmt.Sprintf("Path %s not found", request.Path),
            }, nil
        }
    })
}
