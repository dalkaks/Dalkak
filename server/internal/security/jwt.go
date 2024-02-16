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

func GenerateAuthTokens(domain string, KMS interfaces.KMS, tokenDto *dtos.GenerateTokenDto) (*dtos.AuthTokens, int64, error) {
	nowTime := timeutils.GetTimestamp()
	accessToken, err := generateAccessToken(domain, KMS, nowTime, tokenDto)
	if err != nil {
		return nil, 0, err
	}

	refreshToken, err := generateRefreshToken(domain, KMS, nowTime, tokenDto)
	if err != nil {
		return nil, 0, err
	}

	return &dtos.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nowTime, nil
}

func generateAccessToken(domain string, KMS interfaces.KMS, nowTime int64, tokenDto *dtos.GenerateTokenDto) (string, error) {
	return createToken(jwt.MapClaims{
		"sub": tokenDto.WalletAddress,
		"iat": nowTime,
		"exp": nowTime + config.AccessTokenTTL,
		"iss": domain,
	}, KMS)
}

func generateRefreshToken(domain string, KMS interfaces.KMS, nowTime int64, tokenDto *dtos.GenerateTokenDto) (string, error) {
	tokenId := generateutils.GenerateUUID()
	return createToken(jwt.MapClaims{
		"sub": tokenDto.WalletAddress,
		"tid": tokenId,
		"exp": nowTime + config.RefreshTokenTTL,
	}, KMS)
}

func createToken(claims jwt.Claims, KMS interfaces.KMS) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	signedPart, err := token.SigningString()
	if err != nil {
		return "", dtos.NewAppError(dtos.ErrCodeInternal, dtos.ErrMsgTokenSignFailed, err)
	}

	signature, err := KMS.CreateSianature(signedPart)
	if err != nil {
		return "", nil
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
