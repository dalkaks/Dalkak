package models

type UserAuthAndSignUpRequest struct {
	WalletAddress string `json:"walletAddress" binding:"required"`
	Signature     string `json:"signature" binding:"required"`
}
