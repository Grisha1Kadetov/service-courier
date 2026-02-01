package delivery_test

import (
	"testing"
	"time"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/handler/delivery"
	modelCourier "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
	modelDelivery "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/delivery"

	"github.com/stretchr/testify/require"
)

func TestFromModelToAssignDTO(t *testing.T) {
	t.Parallel()

	now := time.Now()

	courierID := int64(10)
	orderID := "order123"
	transport := modelCourier.TransportCar

	mDelivery := modelDelivery.Delivery{
		CourierID: &courierID,
		OrderID:   &orderID,
		Deadline:  &now,
	}

	mCourier := modelCourier.Courier{
		ID:            &courierID,
		TransportType: &transport,
	}

	dtoResp := delivery.FromModelToAssignDTO(mDelivery, mCourier)

	require.Equal(t, mDelivery.CourierID, dtoResp.CourierID)
	require.Equal(t, mDelivery.OrderID, dtoResp.OrderID)
	require.Equal(t, string(*mCourier.TransportType), *dtoResp.TransportType)
	require.Equal(t, mDelivery.Deadline, dtoResp.Deadline)
}

func TestFromModelToUnassignDTO(t *testing.T) {
	t.Parallel()

	orderID := "order999"
	courierID := int64(55)

	mDelivery := modelDelivery.Delivery{
		OrderID:   &orderID,
		CourierID: &courierID,
	}

	dtoResp := delivery.FromModelToUnassignDTO(mDelivery)

	require.Equal(t, mDelivery.OrderID, dtoResp.OrderID)
	require.Equal(t, mDelivery.CourierID, dtoResp.CourierID)
	require.NotNil(t, dtoResp.Status)
	require.Equal(t, "unassigned", *dtoResp.Status)
}

func TestAssignRequestToModel(t *testing.T) {
	t.Parallel()

	order := "abc"
	req := delivery.AssignRequest{OrderID: &order}

	m := req.ToModel()

	require.Equal(t, req.OrderID, m.OrderID)
}

func TestUnassignRequestToModel(t *testing.T) {
	t.Parallel()

	order := "xyz"
	req := delivery.UnassignRequest{OrderID: &order}

	m := req.ToModel()

	require.Equal(t, req.OrderID, m.OrderID)
}
