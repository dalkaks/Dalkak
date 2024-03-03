package orderaggregate

import (
	orderentity "dalkak/internal/domain/order/object/entity"
	ordervalueobject "dalkak/internal/domain/order/object/valueobject"
)

type OrderAggregate struct {
	OrderEntity   orderentity.OrderEntity
	OrderCategory ordervalueobject.OrderCategory
	OrderPrice    ordervalueobject.OrderPrice
	OrderPayment  *interface{}
}

type OrderAggregateOption func(*OrderAggregate)

func NewOrderAggregate(entity *orderentity.OrderEntity, category *ordervalueobject.OrderCategory, price *ordervalueobject.OrderPrice, options ...OrderAggregateOption) *OrderAggregate {
	aggregate := &OrderAggregate{
		OrderEntity:   *entity,
		OrderCategory: *category,
		OrderPrice:    *price,
	}
	for _, option := range options {
		option(aggregate)
	}
	return aggregate
}
