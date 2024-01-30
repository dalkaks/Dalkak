package dtos

type AuthTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type GenerateTokenDto struct {
	WalletAddress string `json:"walletAddress"`
	NowTime       int64  `json:"nowTime"`
}
