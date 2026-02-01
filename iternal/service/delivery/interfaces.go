package delivery

import (
	"context"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/delivery"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/pkg/executor"
	"github.com/jackc/pgx/v5"
)

type courierRepository interface {
	Patch(ctx context.Context, courier courier.Courier, e executor.Executor) error
	GetAvailable(ctx context.Context) (courier.Courier, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

type deliveryRepository interface {
	Create(ctx context.Context, delivery delivery.Delivery, e executor.Executor) error
	Delete(ctx context.Context, id int64, e executor.Executor) error
	GetByOrderID(ctx context.Context, orderID string) (delivery.Delivery, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

type watcherRepository interface {
	ReleaseExpiredBusyCouriers(ctx context.Context) error
}
