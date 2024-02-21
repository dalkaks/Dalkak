package jwtutil

import (
	"dalkak/internal/core"
	responseutil "dalkak/pkg/utils/response"
	generateutil "dalkak/pkg/utils/generate"

	"github.com/golang-jwt/jwt/v5"
)

const AccessTokenTTL = 30 * 60
const RefreshTokenTTL = 14 * 24 * 60 * 60
const AdditonalTokenTTL = 10

type AuthToken struct {
	Token    string `json:"token"`
	TokenTTL int64  `json:"tokenTTL"`
}

type GenerateTokenDto struct {
	WalletAddress string `json:"walletAddress"`
	NowTime       int64  `json:"nowTime"`
}

func GenerateAuthToken(domain string, keyManager core.KeyManager, dto *GenerateTokenDto) (*AuthToken, *AuthToken, error) {
	accessToken, err := GenerateAccessToken(domain, keyManager, dto)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := GenerateRefreshToken(domain, keyManager, dto)
	if err != nil {
		return nil, nil, err
	}
	return accessToken, refreshToken, nil
}

func GenerateAccessToken(domain string, keyManager core.KeyManager, dto *GenerateTokenDto) (*AuthToken, error) {
	token, err := createToken(jwt.MapClaims{
		"sub": dto.WalletAddress,
		"iat": dto.NowTime,
		"exp": dto.NowTime + AccessTokenTTL + AdditonalTokenTTL,
		"iss": domain,
	}, keyManager)
	if err != nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgServerInternal, err)
	}
	return newAuthToken(token, AccessTokenTTL), nil
}

func GenerateRefreshToken(domain string, keyManager core.KeyManager, dto *GenerateTokenDto) (*AuthToken, error) {
	tokenId := generateutil.GenerateUUID()
	token, err := createToken(jwt.MapClaims{
		"sub": dto.WalletAddress,
		"tid": tokenId,
		"exp": dto.NowTime + RefreshTokenTTL + AdditonalTokenTTL,
	}, keyManager)
	if err != nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgServerInternal, err)
	}
	return newAuthToken(token, RefreshTokenTTL), nil
}

func createToken(claims jwt.Claims, keyManager core.KeyManager) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	signedPart, err := token.SigningString()
	if err != nil {
		return "", responseutil.NewAppError(responseutil.ErrCodeInternal, responseutil.ErrMsgServerInternal, err)
	}

	signature, err := keyManager.CreateSianature(signedPart)
	if err != nil {
		return "", err
	}

	signedToken := signedPart + "." + signature
	return signedToken, nil
}

func newAuthToken(token string, tokenTTL int64) *AuthToken {
	return &AuthToken{
		Token:    token,
		TokenTTL: tokenTTL,
	}
}
