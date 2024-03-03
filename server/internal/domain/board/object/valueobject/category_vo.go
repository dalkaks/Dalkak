package boardvalueobject

import responseutil "dalkak/pkg/utils/response"

type BoardCategory struct {
	CategoryType BoardCategoryType
	CategoryId   string
	Network      NetworkType
}

type BoardCategoryType string

const (
	BoardDefaultNft BoardCategoryType = "defaultNft"
	BoardCustomNft  BoardCategoryType = "customNft"
)

type NetworkType string

const (
	EthereumMainnet NetworkType = "ethereumMainnet"
)

// todo remove hardcoding
func NewBoardCategory(categoryTypeStr, networkStr string) (*BoardCategory, error) {
	var categoryType BoardCategoryType
	switch categoryTypeStr {
	case string(BoardDefaultNft):
		categoryType = BoardDefaultNft
	case string(BoardCustomNft):
		categoryType = BoardCustomNft
	default:
		return nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}

	var network NetworkType
	switch networkStr {
	case string(EthereumMainnet):
		network = EthereumMainnet
	default:
		return nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}

	return &BoardCategory{
		CategoryType: categoryType,
		CategoryId:   "default",
		Network:      network,
	}, nil
}

func (bc BoardCategory) GetCategoryType() string {
	return string(bc.CategoryType)
}

func (bc BoardCategory) GetCategoryId() string {
	return bc.CategoryId
}

func (bc BoardCategory) GetNetwork() string {
	return string(bc.Network)
}
