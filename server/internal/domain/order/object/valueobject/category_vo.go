package ordervalueobject

import responseutil "dalkak/pkg/utils/response"

type OrderCategory struct {
	CategoryType OrderCategoryType
	CategoryId   string
}

type OrderCategoryType string

const (
	OrderCategoryTypeNft OrderCategoryType = "NFT"
)

func NewOrderCategory(categoryTypeStr string, categoryId string) (*OrderCategory, error) {
	categoryType, err := NewOrderCategoryType(categoryTypeStr)
	if err != nil {
		return nil, err
	}

	return &OrderCategory{
		CategoryType: categoryType,
		CategoryId:   categoryId,
	}, nil
}

func NewOrderCategoryType(categoryTypeStr string) (OrderCategoryType, error) {
	switch categoryTypeStr {
	case string(OrderCategoryTypeNft):
		return OrderCategoryTypeNft, nil
	default:
		return "", responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}
}

func (oct OrderCategoryType) String() string {
	return string(oct)
}
