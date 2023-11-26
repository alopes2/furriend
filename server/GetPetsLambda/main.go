package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"server/domain"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	log "github.com/sirupsen/logrus"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  petID := request.PathParameters["petID"]
  
  logger := log.New().WithContext(ctx).WithFields(log.Fields{
    domain.ReqKey: request.RequestContext.RequestID,
    domain.PetID: petID,
  });

	logger.Info("Processing request data for request")

  if strings.TrimSpace(petID) == "" {
    return events.APIGatewayProxyResponse{StatusCode: 400}, nil
  }

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(domain.PetsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"PetID": {
				S: aws.String(petID),
			},
		},
	})
  
	if err != nil {
		logger.Error("Got error calling GetItem: ", err)
	}

	if result.Item == nil {
		logger.Warn(fmt.Sprintf("Could not find pet with ID {petID}, %s", petID));
		return events.APIGatewayProxyResponse{StatusCode: 404}, nil
	}

	pet := domain.Pet{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &pet)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

  jsonResponse, jsonErr := json.Marshal(pet)

  if jsonErr != nil {
		logger.Error("Failed to unmarshal Record", jsonErr)
  }


	return events.APIGatewayProxyResponse{Body: string(jsonResponse), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
