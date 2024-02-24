package blockchaintype

type NftAttribute struct {
	TraitType   string      `json:"traitType,omitempty"`
	Value       interface{} `json:"value" validate:"required"`
	DisplayType string      `json:"displayType,omitempty"`
	MaxValue    int         `json:"maxValue,omitempty"`
}
