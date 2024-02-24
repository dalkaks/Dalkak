package userdto

type CheckAndCreateUserDto struct {
	WalletAddress string
}

func NewCheckAndCreateUserDto(walletAddress string) *CheckAndCreateUserDto {
	return &CheckAndCreateUserDto{
		WalletAddress: walletAddress,
	}
}
