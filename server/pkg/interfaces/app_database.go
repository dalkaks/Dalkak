package interfaces

type Database interface {
	FindOrCreateUser(walletAddress string) (string, error)
}
