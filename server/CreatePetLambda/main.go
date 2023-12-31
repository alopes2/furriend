package main

import (
	"context"
	"encoding/json"
	"strings"

	"server/domain"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger := log.New().WithContext(ctx).WithFields(log.Fields{
		domain.ReqKey: request.RequestContext.RequestID,
	})

	logger.Info("Processing request data for request")

	pet := domain.Pet{
    PetID: uuid.New().String(),
  }

	err := json.Unmarshal([]byte(request.Body), &pet)
  pet.Specie = strings.ToUpper(pet.Specie)

	if err != nil {
		logger.WithField("Body", request.Body).Error("Failed to marshal Pet ", err)
    errorResponse, _ := json.Marshal([]string{"Body incorrect"});
		return events.APIGatewayProxyResponse{Body: string(errorResponse), StatusCode: 400}, nil
	}

	validationResult := domain.PetRequestValidator(pet)

	if !validationResult.IsValid {
		errorResponse, _ := json.Marshal(validationResult.Errors)
		return events.APIGatewayProxyResponse{Body: string(errorResponse), StatusCode: 400}, nil
	}

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(pet)

	if err != nil {
		logger.WithField("Pet", pet).Error("Failed to marshal Pet", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	_, err = svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(domain.PetsTable),
		Item:      av,
	})

	if err != nil {
		logger.Error("Got error calling PutItem: ", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	jsonResponse, jsonErr := json.Marshal(pet)

	if jsonErr != nil {
		logger.Error("Failed to unmarshal Record ", jsonErr)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{Body: string(jsonResponse), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
