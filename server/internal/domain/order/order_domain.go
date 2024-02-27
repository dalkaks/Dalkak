package orderdomain

import (
	"dalkak/internal/core"
	orderfactory "dalkak/internal/domain/order/factory"
	orderaggregate "dalkak/internal/domain/order/object/aggregate"
	orderdto "dalkak/pkg/dto/order"
)

type OrderDomainService interface {
	CreateOrder(dto *orderdto.CreateOrderDto) (*orderaggregate.OrderAggregate, error)
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
	factory := orderfactory.NewCreateOrderDtoFactory(dto)
	order, err := factory.CreateOrderAggregate()
	if err != nil {
		return nil, err
	}

	return order, nil
}
