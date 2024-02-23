package key

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	responseutil "dalkak/pkg/utils/response"
	timeutil "dalkak/pkg/utils/time"
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
	"github.com/golang-jwt/jwt/v5"
)

type KmsSet struct {
	Client    *kms.Client
	KeyId     string
	PublicKey *ecdsa.PublicKey
	Mode			string
	Domain		string
}

func NewKeyManager(ctx context.Context, mode string, keyId string, domain string) (*KmsSet, error) {
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
		Mode:			mode,
		Domain:		domain,
	}, nil
}

func (kmsSet *KmsSet) GetDomain() string {
	return kmsSet.Domain
}

func (kmsSet *KmsSet) GetMode() string {
	return kmsSet.Mode
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
		return "", responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgServerInternal, err)
	}
	signature := base64.RawURLEncoding.EncodeToString(signOutput.Signature)
	return signature, nil
}

func (kmsSet *KmsSet) ParseTokenWithPublicKey(tokenString string) (string, error) {
	err := kmsSet.verifyTokenSignature(tokenString)
	if err != nil {
		return "", err
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgTokenParseFailed)
	}

	if exp, ok := claims["exp"].(float64); ok {
		nowTime := timeutil.GetTimestamp()
		if int64(exp) < nowTime {
			return "", responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgTokenExpired)
		}
	} else {
		return "", responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgTokenParseFailed)
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgTokenParseFailed)
	}
	return sub, nil
}

func (kmsSet *KmsSet) verifyTokenSignature(token string) error {
	jwtParts := strings.Split(token, ".")
	hashInput := strings.Join(jwtParts[:2], ".")
	digest := sha256.Sum256([]byte(hashInput))

	sigDer, err := base64.RawURLEncoding.DecodeString(jwtParts[2])
	if err != nil {
		return responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgTokenInvalidSignature, err)
	}

	type ECDSASignature struct {
		R, S *big.Int
	}
	sigRS := &ECDSASignature{}
	_, err = asn1.Unmarshal(sigDer, sigRS)
	if err != nil {
		return responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgTokenInvalidSignature, err)
	}

	ok := ecdsa.Verify(kmsSet.PublicKey, digest[:], sigRS.R, sigRS.S)
	if !ok {
		return responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgTokenInvalidSignature)
	}
	return nil
}
