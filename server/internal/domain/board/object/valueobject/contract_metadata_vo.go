package boardvalueobject

type ContractMetadata struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Image       string `json:"image" validate:"required"`
	ExternalUrl string `json:"externalUrl"`
}