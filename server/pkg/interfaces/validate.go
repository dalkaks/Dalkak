package interfaces

type ValidatableRequest interface {
	IsValid() bool
}
