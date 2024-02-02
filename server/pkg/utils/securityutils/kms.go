package securityutils

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

type KmsSet struct {
	Client    *kms.Client
	KeyId     string
	PublicKey *ecdsa.PublicKey
}

func GetKMSClient(ctx context.Context, keyId string) (*KmsSet, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := kms.NewFromConfig(cfg)
	publicKey, err := client.GetPublicKey(ctx, &kms.GetPublicKeyInput{
		KeyId: &keyId,
	})
	if err != nil {
		return nil, err
	}

	pemEncodedPublicKey := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKey.PublicKey,
	})
	block, _ := pem.Decode(pemEncodedPublicKey)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode public key")
	}
	pubEdcsaKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &KmsSet{
		Client:    client,
		KeyId:     keyId,
		PublicKey: pubEdcsaKey.(*ecdsa.PublicKey),
	}, nil
}
