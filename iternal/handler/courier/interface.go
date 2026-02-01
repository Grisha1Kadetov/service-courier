package courier

import (
	"context"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
)

type courierService interface {
	CreateCourier(ctx context.Context, courier courier.Courier) error
	PatchCourier(ctx context.Context, courier courier.Courier) error
	GetCourier(ctx context.Context, id int64) (courier.Courier, error)
	GetCouriers(ctx context.Context) ([]courier.Courier, error)
}
