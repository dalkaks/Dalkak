package interfaces

type UserService interface {
	AuthAndSignUp(walletAddress string, signature string) (string, error)
}
