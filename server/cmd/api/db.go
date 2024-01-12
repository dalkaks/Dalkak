package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DB struct {
	client *dynamodb.Client
}

func (app *application) connectToDB(ctx context.Context) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	dbClient := dynamodb.NewFromConfig(cfg)
	return dbClient, nil
}
