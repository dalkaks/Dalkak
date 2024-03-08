package orderentity

import (
	responseutil "dalkak/pkg/utils/response"
	timeutil "dalkak/pkg/utils/time"
)

type OrderEntity struct {
	Id           string            `json:"id"`
	UserId       string            `json:"userId"`
	Timestamp    int64             `json:"timestamp"`
	Name         string            `json:"name"`
	Status       OrderStatus       `json:"status"`
	CategoryType OrderCategoryType `json:"categoryType"`
	CategoryId   string            `json:"categoryId"`
}

type OrderStatus string

const (
	OrderCreated          OrderStatus = "created"
	PaymentStatusPaid     OrderStatus = "paid"
	PaymentStatusFailed   OrderStatus = "payFailed"
	PaymentStatusCanceled OrderStatus = "payCanceled"
)

type OrderCategoryType string

const (
	OrderCategoryTypeNft OrderCategoryType = "NFT"
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

func NewOrderEntity(userId, name, categoryTypeStr, categoryId string, options ...OrderEntityOption) (*OrderEntity, error) {
	orderId := newOrderId(categoryTypeStr, categoryId)
	categoryType, err := newOrderCategoryType(categoryTypeStr)
	if err != nil {
		return nil, err
	}

	oe := &OrderEntity{
		Id:           orderId,
		UserId:       userId,
		Timestamp:    timeutil.GetTimestamp(),
		Name:         name,
		Status:       OrderCreated,
		CategoryType: categoryType,
		CategoryId:   categoryId,
	}
	for _, option := range options {
		option(oe)
	}
	return oe, nil
}

func newOrderId(categoryTypeStr, categoryId string) string {
	return categoryTypeStr + "-" + categoryId
}

func newOrderCategoryType(categoryTypeStr string) (OrderCategoryType, error) {
	switch categoryTypeStr {
	case string(OrderCategoryTypeNft):
		return OrderCategoryTypeNft, nil
	default:
		return "", responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}
}

func (os OrderStatus) String() string {
	return string(os)
}

func (oct OrderCategoryType) String() string {
	return string(oct)
}
