package interfaces

type KMS interface {
	CreateSianature(sign string) (string, error)
	VerifyTokenSignature(sign string) error
}
