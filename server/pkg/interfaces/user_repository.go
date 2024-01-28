package interfaces

type UserRepository interface {
	FindUser(walletAddress string) (*UserDto, error)
	CreateUser(walletAddress string) (string, error)
}
