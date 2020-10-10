package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == http.MethodGet {
		id := uuid.New()
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       `{"uuid":"` + id.String() + `"}`,
		}, nil
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       `{"error":"bad request"}`,
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
