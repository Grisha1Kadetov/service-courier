package order

import (
	"context"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/order"
)

type orderService interface {
	ProcessOrder(ctx context.Context, o order.Order) error
}
