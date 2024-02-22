package userdto

type AuthAndSignUpRequest struct {
	WalletAddress string `json:"walletAddress" validate:"required"`
	Signature     string `json:"signature" validate:"required"`
}

type AuthAndSignUpResponse struct {
	AccessToken     string `json:"accessToken"`
	AccessTokenTTL  int64  `json:"accessTokenTTL"`
	RefreshToken    string `json:"refreshToken"`
	RefreshTokenTTL int64  `json:"refreshTokenTTL"`
}

func NewAuthAndSignUpResponse(accessToken string, accessTokenTTL int64, refreshToken string, refreshTokenTTL int64) *AuthAndSignUpResponse {
	return &AuthAndSignUpResponse{
		AccessToken:     accessToken,
		AccessTokenTTL:  accessTokenTTL,
		RefreshToken:    refreshToken,
		RefreshTokenTTL: refreshTokenTTL,
	}
}

// parse cookie
type ReissueAccessTokenRequest struct {
	WalletAddress string `json:"walletAddress" validate:"required"`
}

type ReissueAccessTokenResponse struct {
	AccessToken    string `json:"accessToken"`
	AccessTokenTTL int64  `json:"accessTokenTTL"`
}

func NewReissueAccessTokenResponse(accessToken string, accessTokenTTL int64) *ReissueAccessTokenResponse {
	return &ReissueAccessTokenResponse{
		AccessToken:    accessToken,
		AccessTokenTTL: accessTokenTTL,
	}
}
