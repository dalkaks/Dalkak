package dtos

type AuthTokens struct {
  AccessToken  string `json:"accessToken"`
  RefreshToken string `json:"refreshToken"`
}