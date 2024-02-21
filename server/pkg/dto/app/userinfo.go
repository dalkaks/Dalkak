package appdto

type UserInfo struct {
	WalletAddress string `json:"walletAddress"`
}

func (u *UserInfo) GetUserId() string {
	return u.WalletAddress
}
