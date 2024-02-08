package app

import (
	"context"
	"dalkak/pkg/interfaces"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DB struct {
	client *dynamodb.Client
	table  string
}

var _ interfaces.Database = (*DB)(nil)

func NewDB(ctx context.Context, mode string) (*DB, error) {
	var cfg aws.Config
	var err error

	if mode == "LOCAL" {
		cfg, err = awsConfig.LoadDefaultConfig(ctx, awsConfig.WithSharedConfigProfile("dalkak"))
	} else {
		cfg, err = awsConfig.LoadDefaultConfig(ctx)
	}

	if err != nil {
		return nil, err
	}

	dbClient := dynamodb.NewFromConfig(cfg)
	var table string
	if mode == "PROD" {
		table = "dalkak_prod"
	} else {
		table = "dalkak_dev"
	}

	return &DB{client: dbClient, table: table}, nil
}

func (db *DB) GetClient() *dynamodb.Client {
	return db.client
}

func (db *DB) GetTable() string {
	return db.table
}
