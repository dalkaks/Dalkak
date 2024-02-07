package dtos

type UserDto struct {
	WalletAddress string `json:"walletAddress"`
	Timestamp     int64  `json:"timestamp"`
}
