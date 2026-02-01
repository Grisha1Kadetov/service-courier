package delivery_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/handler/delivery"
	modelCourier "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
	modelDelivery "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/delivery"
	server "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/server/delivery"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type DeliveryServiceMock struct {
	mock.Mock
}

func (m *DeliveryServiceMock) AssignDelivery(ctx context.Context, orderId string) (modelDelivery.Delivery, modelCourier.Courier, error) {
	args := m.Called(ctx, orderId)
	return args.Get(0).(modelDelivery.Delivery), args.Get(1).(modelCourier.Courier), args.Error(2)
}

func (m *DeliveryServiceMock) UnassignDelivery(ctx context.Context, orderId string) (modelDelivery.Delivery, error) {
	args := m.Called(ctx, orderId)
	return args.Get(0).(modelDelivery.Delivery), args.Error(1)
}

func setupRouter(h *delivery.DeliveryHandler) *chi.Mux {
	return server.NewRouter(h, nil)
}

func TestAssignDelivery_Success(t *testing.T) {
	t.Parallel()

	svc := new(DeliveryServiceMock)
	h := delivery.NewDeliveryHandler(svc)
	r := setupRouter(h)

	orderID := "order123"
	courierID := int64(10)
	now := time.Now()

	deliveryModel := modelDelivery.Delivery{
		OrderID:   &orderID,
		CourierID: &courierID,
		Deadline:  &now,
	}

	tr := modelCourier.TransportCar
	courierModel := modelCourier.Courier{
		ID:            &courierID,
		TransportType: &tr,
	}

	svc.On("AssignDelivery", mock.Anything, orderID).Return(
		deliveryModel,
		courierModel,
		nil,
	)

	body, _ := json.Marshal(delivery.AssignRequest{OrderID: &orderID})
	req := httptest.NewRequest("POST", "/delivery/assign", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp delivery.AssignResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

	require.Equal(t, deliveryModel.OrderID, resp.OrderID)
	require.Equal(t, deliveryModel.CourierID, resp.CourierID)
	require.Equal(t, string(*courierModel.TransportType), *resp.TransportType)
	require.Equal(t, deliveryModel.Deadline.Format(time.RFC3339), resp.Deadline.Format(time.RFC3339))
}

func TestAssignDelivery_InvalidJSON(t *testing.T) {
	t.Parallel()

	svc := new(DeliveryServiceMock)
	h := delivery.NewDeliveryHandler(svc)
	r := setupRouter(h)

	req := httptest.NewRequest("POST", "/delivery/assign", bytes.NewBufferString("{invalid}"))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAssignDelivery_MissingOrderID(t *testing.T) {
	t.Parallel()

	svc := new(DeliveryServiceMock)
	h := delivery.NewDeliveryHandler(svc)
	r := setupRouter(h)

	req := httptest.NewRequest("POST", "/delivery/assign", bytes.NewBufferString("{}"))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAssignDelivery_NotFound(t *testing.T) {
	t.Parallel()

	svc := new(DeliveryServiceMock)
	h := delivery.NewDeliveryHandler(svc)
	r := setupRouter(h)

	orderID := "abc"

	svc.On("AssignDelivery", mock.Anything, orderID).
		Return(modelDelivery.Delivery{}, modelCourier.Courier{}, modelDelivery.ErrNotFound)

	body, _ := json.Marshal(delivery.AssignRequest{OrderID: &orderID})
	req := httptest.NewRequest("POST", "/delivery/assign", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestAssignDelivery_Conflict(t *testing.T) {
	t.Parallel()

	svc := new(DeliveryServiceMock)
	h := delivery.NewDeliveryHandler(svc)
	r := setupRouter(h)

	orderID := "abc"

	svc.On("AssignDelivery", mock.Anything, orderID).
		Return(modelDelivery.Delivery{}, modelCourier.Courier{}, modelDelivery.ErrConflict)

	body, _ := json.Marshal(delivery.AssignRequest{OrderID: &orderID})
	req := httptest.NewRequest("POST", "/delivery/assign", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusConflict, rr.Code)
}

func TestUnassignDelivery_Success(t *testing.T) {
	t.Parallel()

	svc := new(DeliveryServiceMock)
	h := delivery.NewDeliveryHandler(svc)
	r := setupRouter(h)

	orderID := "o1"
	courierID := int64(20)

	deliveryModel := modelDelivery.Delivery{
		OrderID:   &orderID,
		CourierID: &courierID,
	}

	svc.On("UnassignDelivery", mock.Anything, orderID).
		Return(deliveryModel, nil)

	body, _ := json.Marshal(delivery.UnassignRequest{OrderID: &orderID})
	req := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp delivery.UnassignResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

	require.Equal(t, deliveryModel.OrderID, resp.OrderID)
	require.Equal(t, deliveryModel.CourierID, resp.CourierID)
	require.Equal(t, "unassigned", *resp.Status)
}

func TestUnassignDelivery_InvalidJSON(t *testing.T) {
	t.Parallel()

	svc := new(DeliveryServiceMock)
	h := delivery.NewDeliveryHandler(svc)
	r := setupRouter(h)

	req := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewBufferString("{invalid}"))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUnassignDelivery_MissingOrderID(t *testing.T) {
	t.Parallel()

	svc := new(DeliveryServiceMock)
	h := delivery.NewDeliveryHandler(svc)
	r := setupRouter(h)

	req := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewBufferString("{}"))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUnassignDelivery_NotFound(t *testing.T) {
	t.Parallel()

	svc := new(DeliveryServiceMock)
	h := delivery.NewDeliveryHandler(svc)
	r := setupRouter(h)

	orderID := "abc"

	svc.On("UnassignDelivery", mock.Anything, orderID).
		Return(modelDelivery.Delivery{}, modelDelivery.ErrNotFound)

	body, _ := json.Marshal(delivery.UnassignRequest{OrderID: &orderID})
	req := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestUnassignDelivery_Conflict(t *testing.T) {
	t.Parallel()

	svc := new(DeliveryServiceMock)
	h := delivery.NewDeliveryHandler(svc)
	r := setupRouter(h)

	orderID := "abc"

	svc.On("UnassignDelivery", mock.Anything, orderID).
		Return(modelDelivery.Delivery{}, modelDelivery.ErrConflict)

	body, _ := json.Marshal(delivery.UnassignRequest{OrderID: &orderID})
	req := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusConflict, rr.Code)
}
