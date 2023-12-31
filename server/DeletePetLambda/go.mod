module server/DeletePetLambda

go 1.21.4

require server/domain v0.0.0

require (
	github.com/aws/aws-lambda-go v1.41.0 // indirect
	github.com/aws/aws-sdk-go v1.48.3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

replace server/domain => ../domain/
