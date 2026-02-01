package order

import (
	"context"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/order"
)

type OrderProcesserCreated struct {
	deliveryService deliveryService
}

func (p *OrderProcesserCreated) ProcessOrder(ctx context.Context, orderID string) error {
	_, _, err := p.deliveryService.AssignDelivery(ctx, orderID)
	return err
}

type OrderProcesserCancelled struct {
	deliveryService deliveryService
}

func (p *OrderProcesserCancelled) ProcessOrder(ctx context.Context, orderID string) error {
	_, err := p.deliveryService.UnassignDelivery(ctx, orderID)
	return err
}

type OrderProcesserCompleted struct {
	deliveryService deliveryService
}

func (p *OrderProcesserCompleted) ProcessOrder(ctx context.Context, orderID string) error {
	return p.deliveryService.CompleteDelivery(ctx, orderID)
}

type OrderProcesserFactory struct {
	deliveryService deliveryService
}

func NewOrderProcessorFactory(deliveryService deliveryService) *OrderProcesserFactory {
	return &OrderProcesserFactory{
		deliveryService: deliveryService,
	}
}

func (f *OrderProcesserFactory) CreateByStatus(status order.Status) (orderProcesser, error) {
	switch status {
	case order.OrderStatusCreated:
		return &OrderProcesserCreated{deliveryService: f.deliveryService}, nil
	case order.OrderStatusDeleted:
		return &OrderProcesserCancelled{deliveryService: f.deliveryService}, nil
	case order.OrderStatusCompleted:
		return &OrderProcesserCompleted{deliveryService: f.deliveryService}, nil
	default:
		return nil, order.ErrNotHandledStatus
	}
}
