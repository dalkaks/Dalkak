package appsecurity

import (
	"dalkak/config"
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/utils/generateutils"
	"dalkak/pkg/utils/timeutils"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(domain string, KMS interfaces.KMS, dto *dtos.GenerateTokenDto) (*dtos.AuthToken, error) {
	token, err := createToken(jwt.MapClaims{
		"sub": dto.WalletAddress,
		"iat": dto.NowTime,
		"exp": dto.NowTime + config.AccessTokenTTL + config.AdditonalTokenTTL,
		"iss": domain,
	}, KMS)
	if err != nil {
		return nil, err
	}

	return dtos.NewAuthToken(token, config.AccessTokenTTL), nil
}

func GenerateRefreshToken(domain string, KMS interfaces.KMS, dto *dtos.GenerateTokenDto) (*dtos.AuthToken, error) {
	tokenId := generateutils.GenerateUUID()
	token, err := createToken(jwt.MapClaims{
		"sub": dto.WalletAddress,
		"tid": tokenId,
		"exp": dto.NowTime + config.RefreshTokenTTL + config.AdditonalTokenTTL,
	}, KMS)
	if err != nil {
		return nil, err
	}

	return dtos.NewAuthToken(token, config.RefreshTokenTTL), nil
}

func createToken(claims jwt.Claims, KMS interfaces.KMS) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	signedPart, err := token.SigningString()
	if err != nil {
		return "", dtos.NewAppError(dtos.ErrCodeInternal, dtos.ErrMsgTokenSignFailed, err)
	}

	signature, err := KMS.CreateSianature(signedPart)
	if err != nil {
		return "", err
	}

	signedToken := signedPart + "." + signature
	return signedToken, nil
}

func ParseTokenWithPublicKey(tokenString string, KMS interfaces.KMS) (string, error) {
	err := KMS.VerifyTokenSignature(tokenString)
	if err != nil {
		return "", err
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", dtos.NewAppError(dtos.ErrCodeUnauthorized, dtos.ErrMsgTokenParseFailed, errors.New("failed to parse token"))
	}

	if exp, ok := claims["exp"].(float64); ok {
		nowTime := timeutils.GetTimestamp()
		if int64(exp) < nowTime {
			return "", dtos.NewAppError(dtos.ErrCodeUnauthorized, dtos.ErrMsgTokenExpired, errors.New("token is expired"))
		}
	} else {
		return "", dtos.NewAppError(dtos.ErrCodeUnauthorized, dtos.ErrMsgTokenInvalidClaim, errors.New("invalid token claim"))
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", dtos.NewAppError(dtos.ErrCodeUnauthorized, dtos.ErrMsgTokenInvalidClaim, errors.New("invalid token claim"))
	}
	return sub, nil
}
