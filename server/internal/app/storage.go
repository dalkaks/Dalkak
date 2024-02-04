package app

import (
	"context"
	"dalkak/pkg/interfaces"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage struct {
	client *s3.Client
	bucket string
}

var _ interfaces.Storage = (*Storage)(nil)

func NewStorage(ctx context.Context, mode string) (*Storage, error) {
	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	storageClient := s3.NewFromConfig(cfg)
	var bucket string
	if mode == "PROD" {
		bucket = "dalkak-prod"
	} else {
		bucket = "dalkak-dev"
	}

  return &Storage{client: storageClient, bucket: bucket}, nil
}
