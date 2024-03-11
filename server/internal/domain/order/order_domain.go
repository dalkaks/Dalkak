package orderdomain

import (
	"dalkak/internal/core"
	orderfactory "dalkak/internal/domain/order/factory"
	orderaggregate "dalkak/internal/domain/order/object/aggregate"
	"dalkak/internal/infrastructure/database/dao"
	orderdto "dalkak/pkg/dto/order"
)

type OrderDomainService interface {
	CreateOrder(dto *orderdto.CreateOrderDto) (*orderaggregate.OrderAggregate, error)
	ConvertBoardDaoToOrder(dao *dao.BoardDao) (*orderaggregate.OrderAggregate, error)
}

type OrderDomainServiceImpl struct {
	Database     OrderRepository
	EventManager core.EventManager
}

func NewOrderDomainService(database OrderRepository, eventManager core.EventManager) OrderDomainService {
	return &OrderDomainServiceImpl{
		Database:     database,
		EventManager: eventManager,
	}
}

func (service *OrderDomainServiceImpl) CreateOrder(dto *orderdto.CreateOrderDto) (*orderaggregate.OrderAggregate, error) {
	factory := orderfactory.NewCreateOrderFactory()
	order, err := factory.CreateOrderAggregate(dto)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (service *OrderDomainServiceImpl) ConvertBoardDaoToOrder(dao *dao.BoardDao) (*orderaggregate.OrderAggregate, error) {
	factory := orderfactory.NewCreateOrderFactory()
	order, err := factory.CreateOrderAggregateFromDao(dao)
	if err != nil {
		return nil, err
	}

	return order, nil
}
