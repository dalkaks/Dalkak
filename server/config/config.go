package config

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type AppConfig struct {
	Origin     string `json:"origin"`
	Domain     string `json:"domain"`
	StaticLink string `json:"staticLink"`
	KmsKeyId   string `json:"kmsKeyId"`
	QueryKey   string `json:"queryKey"`
}

type LambdaConfig struct {
	DevDomain  string `json:"devDomain"`
	ProdDomain string `json:"prodDomain"`
}

func LoadConfig[T any](ctx context.Context, mode string, parameterName string) (*T, error) {
	var cfg aws.Config
	var err error

	if mode == "LOCAL" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile("dalkak"))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}

	if err != nil {
		return nil, err
	}

	ssmClient := ssm.NewFromConfig(cfg)
	param, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(mode + "-" + parameterName),
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
