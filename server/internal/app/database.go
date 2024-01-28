package app

import (
	"context"
	"dalkak/pkg/interfaces"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DB struct {
	client *dynamodb.Client
	prefix string
}

var _ interfaces.Database = (*DB)(nil)

func NewDB(ctx context.Context, mode string) (*DB, error) {
	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	dbClient := dynamodb.NewFromConfig(cfg)
	prefix := "dalkak_" + mode + "_"

	return &DB{client: dbClient, prefix: prefix}, nil
}
