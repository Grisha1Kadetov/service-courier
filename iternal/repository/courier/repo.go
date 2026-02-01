package courier

import (
	"context"
	"errors"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/pkg/executor"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var psq = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
var UniqueViolationCode string = "23505"

type CourierRepositoryPG struct {
	pool *pgxpool.Pool
}

func NewCourierRepository(pool *pgxpool.Pool) *CourierRepositoryPG {
	return &CourierRepositoryPG{pool: pool}
}

func (r *CourierRepositoryPG) Create(ctx context.Context, courier courier.Courier, e executor.Executor) error {
	if e == nil {
		e = r.pool
	}
	sql, args, _ := psq.Insert("couriers").
		Columns("name", "phone", "status", "transport_type").
		Values(courier.Name, courier.Phone, courier.Status, courier.TransportType).ToSql()

	_, err := e.Exec(ctx, sql, args...)
	if err != nil {
		return checkUnicError(err)
	}
	return nil
}

func (r *CourierRepositoryPG) Patch(ctx context.Context, c courier.Courier, e executor.Executor) error {
	if e == nil {
		e = r.pool
	}
	updates := make(map[string]any)

	if c.Phone != nil {
		updates["phone"] = c.Phone
	}
	if c.Status != nil {
		updates["status"] = c.Status
	}
	if c.Name != nil {
		updates["name"] = c.Name
	}
	if c.TransportType != nil {
		updates["transport_type"] = c.TransportType
	}

	if len(updates) == 0 {
		return courier.ErrNothingToUpdate
	} else {
		updates["updated_at"] = sq.Expr("NOW()")
	}

	sql, args, _ := psq.Update("couriers").SetMap(updates).Where(sq.Eq{"id": c.ID}).ToSql()

	tag, err := e.Exec(ctx, sql, args...)
	if err != nil {
		return checkUnicError(err)
	}
	if tag.RowsAffected() == 0 {
		return courier.ErrNotFound
	}
	return nil
}

func (r *CourierRepositoryPG) GetByID(ctx context.Context, id int64) (courier.Courier, error) {
	sql, args, _ := psq.Select("id", "name", "phone", "status", "transport_type").From("couriers").Where(sq.Eq{"id": id}).ToSql()

	var res courier.Courier
	err := r.pool.QueryRow(ctx, sql, args...).Scan(&res.ID, &res.Name, &res.Phone, &res.Status, &res.TransportType)
	if err != nil {
		if err == pgx.ErrNoRows {
			return res, courier.ErrNotFound
		}
		return res, checkUnicError(err)
	}

	return res, nil
}

func (r *CourierRepositoryPG) GetAll(ctx context.Context) ([]courier.Courier, error) {
	sql, args, _ := psq.Select("id", "name", "phone", "status", "transport_type").From("couriers").ToSql()
	var res []courier.Courier
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return res, checkUnicError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var courier courier.Courier
		err := rows.Scan(&courier.ID, &courier.Name, &courier.Phone, &courier.Status, &courier.TransportType)
		if err != nil {
			return res, checkUnicError(err)
		}
		res = append(res, courier)
	}
	return res, nil
}

func (r *CourierRepositoryPG) GetAvailable(ctx context.Context) (courier.Courier, error) {
	sql, args, _ := psq.
		Select(
			"c.id",
			"c.name",
			"c.phone",
			"c.status",
			"c.transport_type",
		).
		From("couriers c").
		LeftJoin("delivery d ON d.courier_id = c.id AND d.deadline < NOW()").
		Where(sq.Eq{"c.status": courier.StatusAvalible}).
		GroupBy("c.id", "c.name", "c.phone", "c.status", "c.transport_type").
		OrderBy("COUNT(d.courier_id) ASC").
		Limit(1).
		ToSql()

	var res courier.Courier
	err := r.pool.QueryRow(ctx, sql, args...).Scan(
		&res.ID, &res.Name, &res.Phone, &res.Status, &res.TransportType,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return res, courier.ErrNotFound
		}
		return res, err
	}

	return res, nil
}

func (r *CourierRepositoryPG) Begin(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	return tx, err
}

func checkUnicError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == UniqueViolationCode {
		return courier.ErrConflict
	}
	return err
}
