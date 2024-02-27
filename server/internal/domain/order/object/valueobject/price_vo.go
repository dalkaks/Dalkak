package ordervalueobject

import responseutil "dalkak/pkg/utils/response"

type OrderPrice struct {
	OriginPrice   int64
	DiscountPrice int64
	PaymentPrice  int64
}

func NewOrderPrice(categoryType *OrderCategory) (*OrderPrice, error) {
	// todo : get price from categoryType
	if categoryType.CategoryType == OrderCategoryTypeNft {
		return &OrderPrice{
			OriginPrice:   2000,
			DiscountPrice: 0,
			PaymentPrice:  2000,
		}, nil
	} else {
		return nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}
}
