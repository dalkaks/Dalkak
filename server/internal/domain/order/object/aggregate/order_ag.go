package orderaggregate

import (
	orderentity "dalkak/internal/domain/order/object/entity"
	ordervalueobject "dalkak/internal/domain/order/object/valueobject"
)

type OrderAggregate struct {
	OrderEntity   orderentity.OrderEntity
	OrderPrice    ordervalueobject.OrderPrice
	OrderPayment  *interface{}
}

type OrderAggregateOption func(*OrderAggregate)

func NewOrderAggregate(entity *orderentity.OrderEntity, price *ordervalueobject.OrderPrice, options ...OrderAggregateOption) *OrderAggregate {
	aggregate := &OrderAggregate{
		OrderEntity:   *entity,
		OrderPrice:    *price,
	}
	for _, option := range options {
		option(aggregate)
	}
	return aggregate
}
