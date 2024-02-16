package dtos

type AuthTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserInfo struct {
  WalletAddress string `json:"walletAddress"`
}

type GenerateTokenDto struct {
	WalletAddress string `json:"walletAddress"`
}

func NewGenerateTokenDto(walletAddress string) *GenerateTokenDto {
	return &GenerateTokenDto{
		WalletAddress: walletAddress,
	}
}