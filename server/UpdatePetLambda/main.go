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

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	petID := request.PathParameters["petID"]

	logger := log.New().WithContext(ctx).WithFields(log.Fields{
		domain.ReqKey: request.RequestContext.RequestID,
		domain.PetID:  petID,
	})

	logger.Info("Processing request data for request")

	if strings.TrimSpace(petID) == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400}, nil
	}

  pet := domain.Pet{}
  
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

	result, err := svc.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(domain.PetsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"PetID": {
				S: aws.String(petID),
			},
		},
    ExpressionAttributeNames: map[string]*string{
			"#name": aws.String("Name"),
			"#specie": aws.String("Specie"),
		},
    ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {
				S: aws.String(pet.Name),
			},
			":specie": {
				S: aws.String(pet.Specie),
			},
		},
    UpdateExpression: aws.String("set #name = :name, #specie = :specie"),
	})

	logger.WithFields(logrus.Fields{
		"UpdateResult": result,
	}).Info("Update result")

	if err != nil {
		logger.Error("Got error calling UpdateItem ", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
