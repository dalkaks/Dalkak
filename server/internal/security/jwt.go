package appsecurity

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"dalkak/config"
	"dalkak/pkg/dtos"
	"dalkak/pkg/utils/generateutils"
	"dalkak/pkg/utils/timeutils"
	"encoding/asn1"
	"encoding/base64"
	"math/big"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAuthTokens(domain string, kmsSet *KmsSet, tokenDto *dtos.GenerateTokenDto) (*dtos.AuthTokens, int64, error) {
	nowTime := timeutils.GetTimestamp()
	accessToken, err := generateAccessToken(domain, kmsSet, nowTime, tokenDto)
	if err != nil {
		return nil, 0, err
	}

	refreshToken, err := generateRefreshToken(domain, kmsSet, nowTime, tokenDto)
	if err != nil {
		return nil, 0, err
	}

	return &dtos.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nowTime, nil
}

func generateAccessToken(domain string, kmsSet *KmsSet, nowTime int64, tokenDto *dtos.GenerateTokenDto) (string, error) {
	return createToken(jwt.MapClaims{
		"sub": tokenDto.WalletAddress,
		"iat": nowTime,
		"exp": nowTime + config.AccessTokenTTL,
		"iss": domain,
	}, kmsSet)
}

func generateRefreshToken(domain string, kmsSet *KmsSet, nowTime int64, tokenDto *dtos.GenerateTokenDto) (string, error) {
	tokenId := generateutils.GenerateUUID()
	return createToken(jwt.MapClaims{
		"sub": tokenDto.WalletAddress,
		"tid": tokenId,
		"exp": nowTime + config.RefreshTokenTTL,
	}, kmsSet)
}

func createToken(claims jwt.Claims, kmsSet *KmsSet) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	signedPart, err := token.SigningString()
	if err != nil {
		return "", &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to sign token",
		}
	}

	signInput := &kms.SignInput{
		KeyId:            &kmsSet.KeyId,
		MessageType:      types.MessageTypeRaw,
		Message:          []byte(signedPart),
		SigningAlgorithm: types.SigningAlgorithmSpecEcdsaSha256,
	}

	signOutput, err := kmsSet.Client.Sign(context.TODO(), signInput)
	if err != nil {
		return "", &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to sign token",
		}
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
		return "", &dtos.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Failed to parse token claims",
		}
	}

	if exp, ok := claims["exp"].(float64); ok {
		nowTime := timeutils.GetTimestamp()
		if int64(exp) < nowTime {
			return "", &dtos.AppError{
				Code:    http.StatusUnauthorized,
				Message: "Token is expired",
			}
		}
	} else {
		return "", &dtos.AppError{
			Code:    http.StatusUnauthorized,
			Message: "exp claim is missing or not a number",
		}
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", &dtos.AppError{
			Code:    http.StatusUnauthorized,
			Message: "sub claim is missing or not a string",
		}
	}
	return sub, nil
}

func verifyTokenSignature(tokenString string, kmsSet *KmsSet) error {
	jwtParts := strings.Split(tokenString, ".")
	hashInput := strings.Join(jwtParts[:2], ".")
	digest := sha256.Sum256([]byte(hashInput))

	sigDer, err := base64.RawURLEncoding.DecodeString(jwtParts[2])
	if err != nil {
		return &dtos.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Failed to decode signature",
		}
	}

	type ECDSASignature struct {
		R, S *big.Int
	}
	sigRS := &ECDSASignature{}
	_, err = asn1.Unmarshal(sigDer, sigRS)
	if err != nil {
		return &dtos.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Failed to unmarshal signature",
		}
	}

	ok := ecdsa.Verify(kmsSet.PublicKey, digest[:], sigRS.R, sigRS.S)
	if !ok {
		return &dtos.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Failed to verify signature",
		}
	}
	return nil
}
