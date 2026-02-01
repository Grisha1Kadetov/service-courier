package delivery_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
	delivery "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/delivery"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/pkg/log"
	service "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/service/delivery"
)

func TestDeliveryService_AssignDelivery_Success(t *testing.T) {
	ctx := context.Background()

	delRepo := new(DeliveryRepoMock)
	courRepo := new(CourierRepoMock)
	timeCalc := new(TimeCalcMock)
	l := log.NewStdLogger()
	tx := new(TxMock)

	svc := service.NewDeliveryService(delRepo, courRepo, timeCalc, l)

	orderId := "A1"

	c := courier.Courier{
		ID:            ptr((int64)(10)),
		TransportType: ptr(courier.TransportCar),
	}

	expectedDeadline := time.Now().Add(5 * time.Minute)

	delRepo.On("Begin", ctx).Return(tx, nil)
	courRepo.On("GetAvailable", ctx).Return(c, nil)
	timeCalc.On("CalculateDeliveryTime", c).Return(expectedDeadline, nil)

	delRepo.On(
		"Create",
		ctx,
		mock.MatchedBy(func(d delivery.Delivery) bool {
			return *d.CourierID == *c.ID && *d.OrderID == orderId
		}),
		tx,
	).Return(nil)

	updated := c
	statusBusy := courier.StatusBusy
	updated.Status = &statusBusy

	courRepo.On("Patch", ctx, updated, tx).Return(nil)
	tx.On("Commit", ctx).Return(nil)

	delivery, resultCourier, err := svc.AssignDelivery(ctx, orderId)

	require.NoError(t, err)
	assert.Equal(t, *c.ID, *delivery.CourierID)
	assert.Equal(t, orderId, *delivery.OrderID)
	assert.WithinDuration(t, expectedDeadline, *delivery.Deadline, time.Minute)
	assert.Equal(t, courier.StatusBusy, *resultCourier.Status)

	delRepo.AssertExpectations(t)
	courRepo.AssertExpectations(t)
	timeCalc.AssertExpectations(t)
}

func TestDeliveryService_AssignDelivery_NoCourier(t *testing.T) {
	ctx := context.Background()

	delRepo := new(DeliveryRepoMock)
	courRepo := new(CourierRepoMock)
	timeCalc := new(TimeCalcMock)
	tx := &TxMock{}

	l := log.NewStdLogger()
	svc := service.NewDeliveryService(delRepo, courRepo, timeCalc, l)

	delRepo.On("Begin", ctx).Return(tx, nil)
	courRepo.On("GetAvailable", ctx).Return(courier.Courier{}, courier.ErrNotFound)

	d, c, err := svc.AssignDelivery(ctx, "ORDER")

	assert.ErrorIs(t, err, delivery.ErrConflict)
	assert.Equal(t, delivery.Delivery{}, d)
	assert.Equal(t, courier.Courier{}, c)

	delRepo.AssertExpectations(t)
	courRepo.AssertExpectations(t)
}

func TestDeliveryService_UnassignDelivery_Success(t *testing.T) {
	ctx := context.Background()

	delRepo := new(DeliveryRepoMock)
	courRepo := new(CourierRepoMock)
	timeCalc := new(TimeCalcMock)
	tx := &TxMock{}

	l := log.NewStdLogger()
	svc := service.NewDeliveryService(delRepo, courRepo, timeCalc, l)

	orderId := "123"

	del := delivery.Delivery{
		ID:        ptr((int64)(10)),
		CourierID: ptr((int64)(1)),
		OrderID:   &orderId,
	}

	delRepo.On("Begin", ctx).Return(tx, nil)
	delRepo.On("GetByOrderID", ctx, orderId).Return(del, nil)
	delRepo.On("Delete", ctx, int64(10), tx).Return(nil)

	status := courier.StatusAvalible
	courierUpdated := courier.Courier{
		ID:     del.CourierID,
		Status: &status,
	}

	courRepo.On("Patch", ctx, courierUpdated, tx).Return(nil)

	tx.On("Commit", ctx).Return(nil)

	res, err := svc.UnassignDelivery(ctx, orderId)

	require.NoError(t, err)
	assert.Equal(t, del, res)

	delRepo.AssertExpectations(t)
	courRepo.AssertExpectations(t)
}
