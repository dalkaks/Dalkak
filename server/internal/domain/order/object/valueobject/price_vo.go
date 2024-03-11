package ordervalueobject

import (
	orderentity "dalkak/internal/domain/order/object/entity"
	responseutil "dalkak/pkg/utils/response"
)

type OrderPrice struct {
	OriginPrice   int64
	DiscountPrice int64
	PaymentPrice  int64
}

func NewOrderPrice(categoryType orderentity.OrderCategoryType) (*OrderPrice, error) {
	// todo : get price from categoryType
	if categoryType == orderentity.OrderCategoryTypeNft {
		return &OrderPrice{
			OriginPrice:   2000,
			DiscountPrice: 0,
			PaymentPrice:  2000,
		}, nil
	} else {
		return nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}
}
