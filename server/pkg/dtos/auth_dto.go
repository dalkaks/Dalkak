package dtos

import "dalkak/pkg/utils/timeutils"

type AuthToken struct {
	Token    string `json:"token"`
	TokenTTL int64  `json:"tokenTTL"`
}

type UserInfo struct {
	WalletAddress string `json:"walletAddress"`
}

type GenerateTokenDto struct {
	WalletAddress string `json:"walletAddress"`
	NowTime       int64  `json:"nowTime"`
}

func NewAuthToken(token string, tokenTTL int64) *AuthToken {
	return &AuthToken{
		Token:    token,
		TokenTTL: tokenTTL,
	}
}

func NewGenerateTokenDto(walletAddress string) *GenerateTokenDto {
	return &GenerateTokenDto{
		WalletAddress: walletAddress,
		NowTime:       timeutils.GetTimestamp(),
	}
}
