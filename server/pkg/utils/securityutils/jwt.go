package securityutils

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"dalkak/pkg/dtos"
	"dalkak/pkg/utils/timeutils"
	"encoding/asn1"
	"encoding/base64"
	"errors"
	"math/big"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const AccessTokenTTL = 30 * 60
const RefreshTokenTTL = 14 * 24 * 60 * 60

func GenerateAccessToken(domain string, kmsSet *KmsSet, tokenDto *dtos.GenerateTokenDto) (string, error) {
	return createToken(jwt.MapClaims{
		"sub": tokenDto.WalletAddress,
		"iat": tokenDto.NowTime,
		"exp": tokenDto.NowTime + AccessTokenTTL,
		"iss": domain,
	}, kmsSet)
}

func GenerateRefreshToken(domain string, kmsSet *KmsSet, tokenDto *dtos.GenerateTokenDto) (string, error) {
	tokenId := uuid.NewString()
	return createToken(jwt.MapClaims{
		"sub": tokenDto.WalletAddress,
		"tid": tokenId,
		"exp": tokenDto.NowTime + RefreshTokenTTL,
	}, kmsSet)
}

func createToken(claims jwt.Claims, kmsSet *KmsSet) (string, error) {
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

func ParseTokenWithPublicKey(tokenString string, kmsSet *KmsSet) (string, error) {
	err := verifyTokenSignature(tokenString, kmsSet)
	if err != nil {
		return "", err
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("unable to parse claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		nowTime := timeutils.GetTimestamp()
		if int64(exp) < nowTime {
			return "", errors.New("token is expired")
		}
	} else {
		return "", errors.New("exp claim is missing")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("sub claim is missing or not a string")
	}
	return sub, nil
}

func verifyTokenSignature(tokenString string, kmsSet *KmsSet) error {
	jwtParts := strings.Split(tokenString, ".")
	hashInput := strings.Join(jwtParts[:2], ".")
	digest := sha256.Sum256([]byte(hashInput))

	sigDer, err := base64.RawURLEncoding.DecodeString(jwtParts[2])
	if err != nil {
		return err
	}

	type ECDSASignature struct {
		R, S *big.Int
	}
	sigRS := &ECDSASignature{}
	_, err = asn1.Unmarshal(sigDer, sigRS)
	if err != nil {
		return err
	}

	ok := ecdsa.Verify(kmsSet.PublicKey, digest[:], sigRS.R, sigRS.S)
	if !ok {
		return errors.New("invalid signature")
	}
	return nil
}
