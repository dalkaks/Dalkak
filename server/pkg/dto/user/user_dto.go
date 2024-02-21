package userdto

type CreateUserDto struct {
	WalletAddress string `json:"walletAddress"`
}

type FindUserDto struct {
	WalletAddress string `json:"walletAddress"`
}
