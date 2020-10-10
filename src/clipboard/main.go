package main

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func handlerCTX(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		return &events.APIGatewayProxyResponse{
			StatusCode: 503,
			Body:       "Something went wrong :(",
		}, nil
	}

	cc := lc.ClientContext

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Hello, " + cc.Client.AppTitle,
	}, nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	dbURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return nil, errors.New("DATABASE_URL is required")
	}
	/*
		if request.HTTPMethod == http.MethodGet {
			return &events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       dbURL,
			}, nil
		}
	*/
	clientOptions := options.Client().ApplyURI(dbURL)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	if request.HTTPMethod == http.MethodGet {
		id := request.QueryStringParameters["id"]
		if id == "" {
			id = "last"
		}

		return &events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       `{"id":"` + id + `"}`,
		}, nil
	}

	if request.HTTPMethod == http.MethodPost {
		return &events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       `{"status":"Post","data","` + request.Body + `"}`,
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"status":"ok"}`,
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
