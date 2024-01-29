package jwtutils

import (
	"context"
	"dalkak/pkg/dtos"
	"dalkak/pkg/utils/kmsutils"
	"encoding/base64"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const AccessTokenTTL = 30 * 60
const RefreshTokenTTL = 14 * 24 * 60 * 60

func GenerateAccessToken(domain string, kmsSet *kmsutils.KmsSet, tokenDto *dtos.GenerateTokenDto) (string, error) {
	return createToken(jwt.MapClaims{
		"sub": tokenDto.WalletAddress,
		"iat": tokenDto.NowTime,
		"exp": tokenDto.NowTime + AccessTokenTTL,
		"iss": domain,
	}, kmsSet)
}

func GenerateRefreshToken(domain string, kmsSet *kmsutils.KmsSet, tokenDto *dtos.GenerateTokenDto) (string, error) {
	tokenId := uuid.NewString()
	return createToken(jwt.MapClaims{
		"sub": tokenDto.WalletAddress,
		"tid": tokenId,
		"exp": tokenDto.NowTime + RefreshTokenTTL,
	}, kmsSet)
}

func createToken(claims jwt.Claims, kmsSet *kmsutils.KmsSet) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	signedPart, err := token.SigningString()
	if err != nil {
		return "", err
	}

	signInput := &kms.SignInput{
		KeyId:            &kmsSet.KeyId,
		MessageType:      types.MessageTypeRaw,
		Message:          []byte(signedPart),
		SigningAlgorithm: types.SigningAlgorithmSpecEcdsaSha256,
	}

	signOutput, err := kmsSet.Client.Sign(context.TODO(), signInput)
	if err != nil {
		return "", err
	}
	signature := base64.RawURLEncoding.EncodeToString(signOutput.Signature)

	signedToken := signedPart + "." + signature
	return signedToken, nil
}
