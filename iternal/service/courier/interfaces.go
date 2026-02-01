package courier

import (
	"context"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/pkg/executor"
	"github.com/jackc/pgx/v5"
)

type courierRepository interface {
	Create(ctx context.Context, courier courier.Courier, e executor.Executor) error
	Patch(ctx context.Context, courier courier.Courier, e executor.Executor) error
	GetByID(ctx context.Context, id int64) (courier.Courier, error)
	GetAll(ctx context.Context) ([]courier.Courier, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}
