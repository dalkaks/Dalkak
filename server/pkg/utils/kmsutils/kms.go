package kmsutils

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

type KmsSet struct {
	Client *kms.Client
	KeyId  string
}

func GetKMSClient(ctx context.Context, keyId string) (*KmsSet, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &KmsSet{
		Client: kms.NewFromConfig(cfg),
		KeyId:  keyId,
	}, nil
}
