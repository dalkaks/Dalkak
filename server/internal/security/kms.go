package appsecurity

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"dalkak/pkg/dtos"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"math/big"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

type KmsSet struct {
	Client    *kms.Client
	KeyId     string
	PublicKey *ecdsa.PublicKey
}

func GetKMSClient(ctx context.Context, mode string, keyId string) (*KmsSet, error) {
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

func (kmsSet *KmsSet) CreateSianature(sign string) (string, error) {
	signInput := &kms.SignInput{
		KeyId:            &kmsSet.KeyId,
		MessageType:      types.MessageTypeRaw,
		Message:          []byte(sign),
		SigningAlgorithm: types.SigningAlgorithmSpecEcdsaSha256,
	}

	signOutput, err := kmsSet.Client.Sign(context.TODO(), signInput)
	if err != nil {
		return "", dtos.NewAppError(dtos.ErrCodeInternal, dtos.ErrMsgTokenSignFailed, err)
	}
	signature := base64.RawURLEncoding.EncodeToString(signOutput.Signature)
	return signature, nil
}

func (kmsSet *KmsSet) VerifyTokenSignature(token string) error {
	jwtParts := strings.Split(token, ".")
	hashInput := strings.Join(jwtParts[:2], ".")
	digest := sha256.Sum256([]byte(hashInput))

	sigDer, err := base64.RawURLEncoding.DecodeString(jwtParts[2])
	if err != nil {
		return dtos.NewAppError(dtos.ErrCodeUnauthorized, dtos.ErrMsgTokenInvalidSignature, err)
	}

	type ECDSASignature struct {
		R, S *big.Int
	}
	sigRS := &ECDSASignature{}
	_, err = asn1.Unmarshal(sigDer, sigRS)
	if err != nil {
		return dtos.NewAppError(dtos.ErrCodeUnauthorized, dtos.ErrMsgTokenInvalidSignature, err)
	}

	ok := ecdsa.Verify(kmsSet.PublicKey, digest[:], sigRS.R, sigRS.S)
	if !ok {
		return dtos.NewAppError(dtos.ErrCodeUnauthorized, dtos.ErrMsgTokenInvalidSignature, errors.New("invalid token signature"))
	}
	return nil
}
