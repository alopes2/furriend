module server/CreatePetLambda

go 1.21.4

require (
	github.com/aws/aws-lambda-go v1.41.0
	github.com/aws/aws-sdk-go v1.48.3
	github.com/sirupsen/logrus v1.9.3
	server/domain v0.0.0
)

require (
	github.com/google/uuid v1.4.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

replace server/domain => ../domain/
