package lambdautils

import (
	"context"
	"dalkak/config"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
)

func GetLambdaSdkConfig(ctx context.Context) (*aws.Config, *config.LambdaConfig, error) {
	sdkConfig, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, nil, err
	}

	lambdaConfig, err := config.LoadConfig[config.LambdaConfig](ctx, "Common", "LambdaConfig")
	if err != nil {
		return nil, nil, err
	}

	return &sdkConfig, lambdaConfig, nil
}

func HandleError(err error) {
	log.Printf("Error: %v", err)
}
