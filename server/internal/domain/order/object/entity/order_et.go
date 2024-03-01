package orderentity

import (
	generateutil "dalkak/pkg/utils/generate"
	timeutil "dalkak/pkg/utils/time"
)

type OrderEntity struct {
	Id        string      `json:"id"`
	UserId    string      `json:"userId"`
	Timestamp int64       `json:"timestamp"`
	Name      string      `json:"name"`
	Status    OrderStatus `json:"status"`
}

type OrderStatus string

const (
	OrderCreated          OrderStatus = "created"
	PaymentStatusPaid     OrderStatus = "paid"
	PaymentStatusFailed   OrderStatus = "payFailed"
	PaymentStatusCanceled OrderStatus = "payCanceled"
)

type OrderEntityOption func(*OrderEntity)

func WithID(id string) OrderEntityOption {
	return func(oe *OrderEntity) {
		oe.Id = id
	}
}

func WithStatus(status OrderStatus) OrderEntityOption {
	return func(oe *OrderEntity) {
		oe.Status = status
	}
}

func NewOrderEntity(userId, name string, options ...OrderEntityOption) *OrderEntity {
	oe := &OrderEntity{
		Id:        generateutil.GenerateUUID(),
		UserId:    userId,
		Timestamp: timeutil.GetTimestamp(),
		Name:      name,
		Status:    OrderCreated,
	}
	for _, option := range options {
		option(oe)
	}
	return oe
}

func (os OrderStatus) String() string {
	return string(os)
}
