# My Project

This is a sample project that demonstrates how to use AWS Lambda with API Gateway, DynamoDB, and SQS using the Serverless Framework and LocalStack. The project is written in Go and includes two Lambda functions:

- `createOrder` - creates a new order in a DynamoDB table
- `processPayment` - simulates processing a payment using an SQS queue

## Project Structure

The project is organized into the following directories and files:

```.
├── bin/
├── src/
│ ├── createOrder/
│ │ ├── createOrder.go
│ │ └── main.go
│ └── processPayment/
│ ├── processPayment.go
│ └── main.go
├── build.sh
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
└── serverless.yml
```

- `bin/` - contains the compiled Go binaries for each Lambda function
- `src/` - contains the source code for each Lambda function, organized into subdirectories
- `build.sh` - a script for compiling the Go binaries for each Lambda function
- `Dockerfile` - a Dockerfile for running the project using LocalStack
- `docker-compose.yml` - a Docker Compose file for starting LocalStack and running the project
- `go.mod` - the Go module file that specifies the project's dependencies
- `go.sum` - a checksum file for the project's dependencies
- `main.go` - the main file that contains the Lambda handler function that routes requests to the appropriate Lambda function
- `serverless.yml` - the Serverless Framework configuration file that specifies the API Gateway, DynamoDB, and SQS resources for the project



## Running the Project Locally

To run the project locally using LocalStack, follow these steps:

1. Install Docker and Docker Compose.
2. Clone the repository and navigate to the project directory.
3. Run `docker-compose up` to start LocalStack and run the project.
4. Navigate to `http://localhost:4566/health` to ensure that LocalStack is running.
5. Use an HTTP client such as `curl` or `Postman` to send requests to the API Gateway endpoints. The endpoints are defined in `serverless.yml` and are printed to the console when you run `docker-compose up`.



## Production (AWS)

To deploy the project to production on AWS, follow these steps:

1. Ensure that you have an AWS account and have set up your AWS credentials on your local machine.
2. Install the Serverless Framework CLI.
3. Set up your Serverless Framework credentials by running serverless config credentials --provider aws --key YOUR_ACCESS_KEY --secret YOUR_SECRET_KEY.
4. Navigate to the project directory.
5. Run ./deploy.sh prod to deploy the project to AWS.

This project is licensed under the MIT License. See `LICENSE` for more information.

