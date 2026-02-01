package delivery_test

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/delivery"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/pkg/executor"
)

type TxMock struct{ mock.Mock }

func (m *TxMock) Commit(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *TxMock) Rollback(ctx context.Context) error {
	return nil
}

func (m *TxMock) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(*TxMock), args.Error(1)
}

func (m *TxMock) CopyFrom(ctx context.Context, table pgx.Identifier, columns []string, rows pgx.CopyFromSource) (int64, error) {
	return 0, nil
}

func (m *TxMock) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return nil
}

func (m *TxMock) LargeObjects() pgx.LargeObjects {
	return pgx.LargeObjects{}
}

func (m *TxMock) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return nil, nil
}

func (m *TxMock) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (m *TxMock) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}

func (m *TxMock) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return nil
}

func (m *TxMock) Conn() *pgx.Conn {
	return nil
}

//--------

type DeliveryRepoMock struct{ mock.Mock }

func (m *DeliveryRepoMock) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(*TxMock), args.Error(1)
}

func (m *DeliveryRepoMock) Create(ctx context.Context, d delivery.Delivery, e executor.Executor) error {
	return m.Called(ctx, d, e).Error(0)
}

func (m *DeliveryRepoMock) GetByOrderID(ctx context.Context, id string) (delivery.Delivery, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(delivery.Delivery), args.Error(1)
}
func (m *DeliveryRepoMock) Delete(ctx context.Context, id int64, e executor.Executor) error {
	return m.Called(ctx, id, e).Error(0)
}

//--------

type CourierRepoMock struct{ mock.Mock }

func (m *CourierRepoMock) GetAvailable(ctx context.Context) (courier.Courier, error) {
	args := m.Called(ctx)
	return args.Get(0).(courier.Courier), args.Error(1)
}

func (m *CourierRepoMock) Patch(ctx context.Context, c courier.Courier, e executor.Executor) error {
	return m.Called(ctx, c, e).Error(0)
}

func (m *CourierRepoMock) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(*TxMock), args.Error(1)
}

//--------

type TimeCalcMock struct{ mock.Mock }

func (m *TimeCalcMock) CalculateDeliveryTime(c courier.Courier) (time.Time, error) {
	args := m.Called(c)
	return args.Get(0).(time.Time), args.Error(1)
}
