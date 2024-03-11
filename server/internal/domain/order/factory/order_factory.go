package orderfactory

import (
	orderaggregate "dalkak/internal/domain/order/object/aggregate"
	orderentity "dalkak/internal/domain/order/object/entity"
	ordervalueobject "dalkak/internal/domain/order/object/valueobject"
	"dalkak/internal/infrastructure/database/dao"
	orderdto "dalkak/pkg/dto/order"
)

type OrderAggregateFactory interface {
	CreateOrderAggregate(dto *orderdto.CreateOrderDto) (*orderaggregate.OrderAggregate, error)
}

type CreateOrderFactory struct {
}

func NewCreateOrderFactory() *CreateOrderFactory {
	return &CreateOrderFactory{}
}

func (factory *CreateOrderFactory) CreateOrderAggregate(dto *orderdto.CreateOrderDto) (*orderaggregate.OrderAggregate, error) {
	order, err := orderentity.NewOrderEntity(dto.UserInfo.GetUserId(), dto.Name, dto.CategoryType, dto.CategoryId)
	if err != nil {
		return nil, err
	}

	price, err := ordervalueobject.NewOrderPrice(order.CategoryType)
	if err != nil {
		return nil, err
	}

	orderAggregate := orderaggregate.NewOrderAggregate(
		order,
		price,
	)
	return orderAggregate, nil
}

func (factory *CreateOrderFactory) CreateOrderAggregateFromDao(dao *dao.BoardDao) (*orderaggregate.OrderAggregate, error) {
	order, err := orderentity.ConvertOrderEntity(dao.Id, dao.UserId, dao.Timestamp, dao.NftMetaName, dao.Status, orderentity.OrderCategoryTypeNft, dao.TypeId)
	if err != nil {
		return nil, err
	}

	price, err := ordervalueobject.NewOrderPrice(order.CategoryType)
	if err != nil {
		return nil, err
	}

	orderAggregate := orderaggregate.NewOrderAggregate(
		order,
		price,
	)
	return orderAggregate, nil
}
