package orderfactory

import (
	orderaggregate "dalkak/internal/domain/order/object/aggregate"
	orderentity "dalkak/internal/domain/order/object/entity"
	ordervalueobject "dalkak/internal/domain/order/object/valueobject"
	orderdto "dalkak/pkg/dto/order"
)

type OrderAggregateFactory interface {
	CreateOrderAggregate() (*orderaggregate.OrderAggregate, error)
}

type CreateOrderDtoFactory struct {
	dto *orderdto.CreateOrderDto
}

func NewCreateOrderDtoFactory(dto *orderdto.CreateOrderDto) *CreateOrderDtoFactory {
	return &CreateOrderDtoFactory{
		dto: dto,
	}
}

func (factory *CreateOrderDtoFactory) CreateOrderAggregate() (*orderaggregate.OrderAggregate, error) {
	order := orderentity.NewOrderEntity(factory.dto.UserInfo.GetUserId(), factory.dto.Name)
	category, err := ordervalueobject.NewOrderCategory(factory.dto.CategoryType, factory.dto.CatetoryId)
	if err != nil {
		return nil, err
	}
	price, err := ordervalueobject.NewOrderPrice(factory.dto.OriginPrice, factory.dto.DiscountPrice, factory.dto.PaymentPrice)
	if err != nil {
		return nil, err
	}

	orderAggregate := orderaggregate.NewOrderAggregate(
		order,
		category,
		price,
	)
	return orderAggregate, nil
}
