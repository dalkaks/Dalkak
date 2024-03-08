package orderdto

import appdto "dalkak/pkg/dto/app"

type CreateOrderDto struct {
	UserInfo     *appdto.UserInfo
	CategoryType string
	CategoryId   string
	Name         string
	PaymentId    *string
}

func NewCreateOrderDto(userInfo *appdto.UserInfo, categoryType, categoryId, name string, paymentId *string) *CreateOrderDto {
	return &CreateOrderDto{
		UserInfo:     userInfo,
		CategoryType: categoryType,
		CategoryId:   categoryId,
		Name:         name,
		PaymentId:    paymentId,
	}
}
