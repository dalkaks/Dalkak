package database

import (
	orderaggregate "dalkak/internal/domain/order/object/aggregate"
)

const OrderDataType = "Order"

func GenerateOrderDataPk(orderId string) string {
	return OrderDataType + `#` + orderId
}

type OrderData struct {
	Pk         string
	Sk         string
	EntityType string
	Timestamp  int64

	Id     string
	Name   string
	Status string
	UserId string

	Type   string
	TypeId string

	OriginPrice   int64
	DiscountPrice int64
	PaymentPrice  int64
}

func CreateOrderData(order *orderaggregate.OrderAggregate) *OrderData {
	return &OrderData{
		Pk:         GenerateOrderDataPk(order.OrderEntity.Id),
		Sk:         GenerateOrderDataPk(order.OrderEntity.Id),
		EntityType: OrderDataType,
		Timestamp:  order.OrderEntity.Timestamp,

		Id:     order.OrderEntity.Id,
		Name:   order.OrderEntity.Name,
		Status: order.OrderEntity.Status.String(),
		UserId: order.OrderEntity.UserId,
		Type:   order.OrderEntity.CategoryType.String(),
		TypeId: order.OrderEntity.CategoryId,

		OriginPrice:   order.OrderPrice.OriginPrice,
		DiscountPrice: order.OrderPrice.DiscountPrice,
		PaymentPrice:  order.OrderPrice.PaymentPrice,
	}
}

// func CreateOrderKey(orderId string) map[string]types.AttributeValue {
// 	return map[string]types.AttributeValue{
// 		"Pk": &types.AttributeValueMemberS{Value: GenerateOrderDataPk(orderId)},
// 		"Sk": &types.AttributeValueMemberS{Value: GenerateOrderDataPk(orderId)},
// 	}
// }
