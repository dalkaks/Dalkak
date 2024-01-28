package interfaces

type UserRepository interface {
	FindOrCreateUser(walletAddress string) (string, error)
}
