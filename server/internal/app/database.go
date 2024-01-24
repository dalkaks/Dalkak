package app

import (
	"context"
	"dalkak/internal/interfaces"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DB struct {
	client *dynamodb.Client
}

var _ interfaces.Database = (*DB)(nil)

func NewDB(ctx context.Context) (*DB, error) {
	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	dbClient := dynamodb.NewFromConfig(cfg)
	return &DB{client: dbClient}, nil
}
