package orderdto

import appdto "dalkak/pkg/dto/app"

type CreateOrderDto struct {
	UserInfo      *appdto.UserInfo
	CategoryType  string
	CatetoryId    string
	Name          string
	PaymentId     *string
}

func NewCreateOrderDto(userInfo *appdto.UserInfo, categoryType, categoryId, name string, paymentId *string) *CreateOrderDto {
	return &CreateOrderDto{
		UserInfo:      userInfo,
		CategoryType:  categoryType,
		CatetoryId:    categoryId,
		Name:          name,
		PaymentId:     paymentId,
	}
}
