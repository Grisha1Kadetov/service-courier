package delivery

import (
	"context"
	"errors"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/delivery"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/pkg/executor"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var psq = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type DeliveryRepositoryPG struct {
	pool *pgxpool.Pool
}

func NewDeliveryRepository(pool *pgxpool.Pool) *DeliveryRepositoryPG {
	return &DeliveryRepositoryPG{pool: pool}
}

func (r *DeliveryRepositoryPG) Create(ctx context.Context, delivery delivery.Delivery, e executor.Executor) error {
	if e == nil {
		e = r.pool
	}
	sql, args, _ := psq.Insert("delivery").
		Columns("courier_id", "order_id", "deadline").
		Values(delivery.CourierID, delivery.OrderID, delivery.Deadline.UTC()).ToSql()

	_, err := e.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *DeliveryRepositoryPG) Delete(ctx context.Context, id int64, e executor.Executor) error {
	if e == nil {
		e = r.pool
	}

	sql, args, _ := psq.
		Delete("delivery").
		Where(sq.Eq{"id": id}).
		ToSql()

	tag, err := e.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return delivery.ErrNotFound
	}

	return nil
}

func (r *DeliveryRepositoryPG) GetByOrderID(ctx context.Context, orderID string) (delivery.Delivery, error) {
	sql, args, _ := psq.
		Select("id", "courier_id", "order_id", "assigned_at", "deadline").
		From("delivery").
		Where(sq.Eq{"order_id": orderID}).
		Limit(1).
		ToSql()

	var d delivery.Delivery

	err := r.pool.QueryRow(ctx, sql, args...).Scan(
		&d.ID, &d.CourierID, &d.OrderID, &d.AssignedAt, &d.Deadline,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return d, delivery.ErrNotFound
		}
		return d, err
	}

	return d, nil
}

func (r *DeliveryRepositoryPG) Begin(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	return tx, err
}

func (r *DeliveryRepositoryPG) ReleaseExpiredBusyCouriers(ctx context.Context) error {
	sql, args, _ := psq.
		Update("couriers c").
		Set("status", "available").
		Where(sq.Eq{"c.status": "busy"}).
		Where(sq.Expr(`
        NOT EXISTS (
            SELECT 1
            FROM delivery d
            WHERE d.courier_id = c.id
              AND d.deadline >= NOW()
        )
    `)).
		ToSql()

	_, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}
