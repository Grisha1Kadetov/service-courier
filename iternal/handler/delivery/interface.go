package delivery

import (
	"context"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/delivery"
)

type deliveryService interface {
	UnassignDelivery(ctx context.Context, orderId string) (delivery.Delivery, error)
	AssignDelivery(ctx context.Context, orderId string) (delivery.Delivery, courier.Courier, error)
}
