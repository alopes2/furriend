module server/GetPetsLambda

go 1.21.4

require (
	github.com/aws/aws-lambda-go v1.41.0 // indirect
	github.com/aws/aws-sdk-go v1.48.3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	server/domain v0.0.0
)

replace server/domain => ../domain/
