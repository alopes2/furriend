package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	domain "Furriend.Domain"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func handleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(request.Body))

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	tableName := "furriend_pets"
	petID := "123"

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PetID": {
				N: aws.String(petID),
			},
		},
	})
  
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	if result.Item == nil {
		msg := "Could not find '" + petID + "'"
		return events.APIGatewayProxyResponse{}, errors.New(msg)
	}

	pet :=domain.Pet{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &pet)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

  jsonResponse, jsonErr := json.Marshal(pet)

  if jsonErr != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", jsonErr))
  }


	return events.APIGatewayProxyResponse{Body: string(jsonResponse), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
