package config

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type AppConfig struct {
	Origin string `json:"origin"`
}

func LoadConfig[T any](ctx context.Context, mod string, parameterName string) (*T, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	ssmClient := ssm.NewFromConfig(cfg)
	param, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(mod + "-" + parameterName),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	var config T
	err = json.Unmarshal([]byte(*param.Parameter.Value), &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
