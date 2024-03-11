package orderentity

import (
	boardentity "dalkak/internal/domain/board/object/entity"
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

func NewOrderEntity(userId, name, categoryTypeStr, categoryId string) (*OrderEntity, error) {
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
	return oe, nil
}

func ConvertOrderEntity(id, userId string, timestamp int64, name, statusStr string, categoryType OrderCategoryType, categoryId string) (*OrderEntity, error) {
	status, err := newOrderStatus(statusStr)
	if err != nil {
		return nil, err
	}

	return &OrderEntity{
		Id:           id,
		UserId:       userId,
		Timestamp:    timestamp,
		Name:         name,
		Status:       status,
		CategoryType: categoryType,
		CategoryId:   categoryId,
	}, nil
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

func newOrderStatus(statusStr string) (OrderStatus, error) {
	switch statusStr {
	case string(OrderCreated):
		return OrderCreated, nil
	case string(PaymentStatusPaid):
		return PaymentStatusPaid, nil
	case string(PaymentStatusFailed):
		return PaymentStatusFailed, nil
	case string(PaymentStatusCanceled):
		return PaymentStatusCanceled, nil

	// board
	// case string(boardentity.BoardCreated):
	// case string(boardentity.PaymentStatusPaid):
	// case string(boardentity.PaymentStatusFailed):
	case string(boardentity.ContractUploaded),
		string(boardentity.ContractUploadFailed),
		string(boardentity.NFTUploaded),
		string(boardentity.NFTUploadFailed),
		string(boardentity.BoardCancelled),
		string(boardentity.BoardPosted):
		return PaymentStatusPaid, nil

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
