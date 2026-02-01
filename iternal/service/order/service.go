package order

import (
	"context"
	"errors"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/delivery"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/order"
)

type OrderService struct {
	orderProcessFactory orderProcessFactory
	orderGateway        orderGateway
}

func NewOrderService(orderGateway orderGateway, orderProcessFactory orderProcessFactory) *OrderService {
	return &OrderService{
		orderProcessFactory: orderProcessFactory,
		orderGateway:        orderGateway,
	}
}

func (s *OrderService) ProcessOrder(ctx context.Context, o order.Order) error {
	actual, err := s.orderGateway.GetOrderById(ctx, o.ID)
	if err != nil {
		if errors.Is(err, order.ErrNotFound) {
			return nil // skip
		}
		return err
	}
	if actual.Status != o.Status {
		return nil // skip
	}

	processer, err := s.orderProcessFactory.CreateByStatus(o.Status)
	if err != nil {
		if errors.Is(err, order.ErrNotHandledStatus) {
			return nil // skip
		}
		return err
	}
	err = processer.ProcessOrder(ctx, o.ID)
	if err != nil {
		if errors.Is(err, delivery.ErrNotFound) {
			return nil
		}
		return err
	}
	return nil
}
