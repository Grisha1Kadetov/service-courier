package order

import (
	"context"
	"time"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/delivery"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/order"
)

type orderGateway interface {
	GetOrders(context.Context, time.Time) ([]order.Order, error)
	GetOrderById(context.Context, string) (order.Order, error)
}

type orderProcessFactory interface {
	CreateByStatus(status order.Status) (orderProcesser, error)
}

type deliveryService interface {
	AssignDelivery(ctx context.Context, orderId string) (delivery.Delivery, courier.Courier, error)
	UnassignDelivery(ctx context.Context, orderId string) (delivery.Delivery, error)
	CompleteDelivery(ctx context.Context, orderId string) error
}

type orderProcesser interface {
	ProcessOrder(ctx context.Context, orderID string) error
}
