package orderdto

import appdto "dalkak/pkg/dto/app"

type CreateOrderDto struct {
	UserInfo      *appdto.UserInfo
	CategoryType  string
	CatetoryId    string
	Name          string
	PaymentId     *string
	OriginPrice   int64
	DiscountPrice int64
	PaymentPrice  int64
}

func NewCreateOrderDto(userInfo *appdto.UserInfo, categoryType, categoryId, name string, paymentId *string, originPrice, discountPrice, paymentPrice int64) *CreateOrderDto {
	return &CreateOrderDto{
		UserInfo:      userInfo,
		CategoryType:  categoryType,
		CatetoryId:    categoryId,
		Name:          name,
		PaymentId:     paymentId,
		OriginPrice:   originPrice,
		DiscountPrice: discountPrice,
		PaymentPrice:  paymentPrice,
	}
}
