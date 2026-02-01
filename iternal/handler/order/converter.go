package order

import (
	"fmt"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/order"
)

func (o Order) ToModel() (order.Order, error) {
	status, ok := order.ValidStatuses[o.Status]
	if !ok {
		return order.Order{}, fmt.Errorf("invalid order status: %s", o.Status)
	}
	return order.Order{
		ID:        o.ID,
		Status:    status,
		CreatedAt: o.CreatedAt,
	}, nil
}
