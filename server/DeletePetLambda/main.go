package main

import (
	"context"
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

	result, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(domain.PetsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"PetID": {
				S: aws.String(petID),
			},
		},
	})

  logger.WithFields(logrus.Fields{
    "DeleteResult": result,
  }).Info("Deletion result")
  
	if err != nil {
		logger.Error("Got error calling DeleteItem: ", err)
	  return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}


	return events.APIGatewayProxyResponse{StatusCode: 204}, nil
}

func main() {
	lambda.Start(handleRequest)
}
