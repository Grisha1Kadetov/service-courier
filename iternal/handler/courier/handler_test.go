package courier_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	handler "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/handler/courier"
	model "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
	server "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/server/courier"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type CourierServiceMock struct {
	mock.Mock
}

func (m *CourierServiceMock) CreateCourier(ctx context.Context, c model.Courier) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *CourierServiceMock) PatchCourier(ctx context.Context, c model.Courier) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *CourierServiceMock) GetCourier(ctx context.Context, id int64) (model.Courier, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Courier), args.Error(1)
}

func (m *CourierServiceMock) GetCouriers(ctx context.Context) ([]model.Courier, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Courier), args.Error(1)
}

func setupRouter(h *handler.CourierHandler) http.Handler {
	return server.NewRouter(h, nil)
}

func TestCreateCourier_Success(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	reqBody := handler.RequestCreate{
		Name:          stringPtr("John Doe"),
		Phone:         stringPtr("1234567890"),
		Status:        stringPtr(string(model.StatusAvalible)),
		TransportType: stringPtr(string(model.TransportOnFoot)),
	}
	body, _ := json.Marshal(reqBody)

	svc.On("CreateCourier", mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest("POST", "/courier", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	svc.AssertCalled(t, "CreateCourier", mock.Anything, mock.Anything)
}

func TestCreateCourier_InvalidJSON(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	req := httptest.NewRequest("POST", "/courier", bytes.NewBufferString("{invalid"))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "invalid JSON")
	svc.AssertNotCalled(t, "CreateCourier", mock.Anything, mock.Anything)
}

func TestCreateCourier_InvalidStatus(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	reqBody := handler.RequestCreate{
		Name:          stringPtr("John Doe"),
		Phone:         stringPtr("1234567890"),
		Status:        stringPtr("wrong_status"),
		TransportType: stringPtr(string(model.TransportOnFoot)),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/courier", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "invalid status")
	svc.AssertNotCalled(t, "CreateCourier", mock.Anything, mock.Anything)
}

func TestCreateCourier_InvalidTransportType(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	reqBody := handler.RequestCreate{
		Name:          stringPtr("John Doe"),
		Phone:         stringPtr("1234567890"),
		Status:        stringPtr(string(model.StatusAvalible)),
		TransportType: stringPtr("wrong_transport"),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/courier", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "invalid transport type")
	svc.AssertNotCalled(t, "CreateCourier", mock.Anything, mock.Anything)
}

func TestCreateCourier_MissingName(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	reqBody := handler.RequestCreate{
		Phone:         stringPtr("1234567890"),
		Status:        stringPtr(string(model.StatusAvalible)),
		TransportType: stringPtr(string(model.TransportOnFoot)),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/courier", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "missing name")
	svc.AssertNotCalled(t, "CreateCourier", mock.Anything, mock.Anything)
}

func TestCreateCourier_MissingPhone(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	reqBody := handler.RequestCreate{
		Name:          stringPtr("John Doe"),
		Status:        stringPtr(string(model.StatusAvalible)),
		TransportType: stringPtr(string(model.TransportOnFoot)),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/courier", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "missing phone")
	svc.AssertNotCalled(t, "CreateCourier", mock.Anything, mock.Anything)
}

func TestCreateCourier_MissingStatus(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	reqBody := handler.RequestCreate{
		Name:          stringPtr("John Doe"),
		Phone:         stringPtr("1234567890"),
		TransportType: stringPtr(string(model.TransportOnFoot)),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/courier", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "missing status")
	svc.AssertNotCalled(t, "CreateCourier", mock.Anything, mock.Anything)
}

func TestCreateCourier_MissingTransportType(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	reqBody := handler.RequestCreate{
		Name:  stringPtr("John Doe"),
		Phone: stringPtr("1234567890"),
		Status: stringPtr(
			string(model.StatusAvalible),
		),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/courier", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "missing transport type")
	svc.AssertNotCalled(t, "CreateCourier", mock.Anything, mock.Anything)
}

func TestCreateCourier_Conflict(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	reqBody := handler.RequestCreate{
		Name:          stringPtr("John Doe"),
		Phone:         stringPtr("1234567890"),
		Status:        stringPtr(string(model.StatusAvalible)),
		TransportType: stringPtr(string(model.TransportOnFoot)),
	}
	body, _ := json.Marshal(reqBody)

	svc.On("CreateCourier", mock.Anything, mock.Anything).Return(model.ErrConflict)

	req := httptest.NewRequest("POST", "/courier", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusConflict, rr.Code)
}

func TestCreateCourier_InternalError(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	reqBody := handler.RequestCreate{
		Name:          stringPtr("John Doe"),
		Phone:         stringPtr("1234567890"),
		Status:        stringPtr(string(model.StatusAvalible)),
		TransportType: stringPtr(string(model.TransportOnFoot)),
	}
	body, _ := json.Marshal(reqBody)

	svc.On("CreateCourier", mock.Anything, mock.Anything).Return(errors.New("boom"))

	req := httptest.NewRequest("POST", "/courier", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestPatchCourier_Success(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	reqBody := handler.RequestUpdate{
		ID:            int64Ptr(1),
		Name:          stringPtr("John Doe"),
		Phone:         stringPtr("1234567890"),
		Status:        stringPtr(string(model.StatusBusy)),
		TransportType: stringPtr(string(model.TransportCar)),
	}
	body, _ := json.Marshal(reqBody)

	svc.On("PatchCourier", mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	svc.AssertCalled(t, "PatchCourier", mock.Anything, mock.Anything)
}

func TestPatchCourier_InvalidJSON(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	req := httptest.NewRequest("PUT", "/courier", bytes.NewBufferString("{invalid"))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	svc.AssertNotCalled(t, "PatchCourier", mock.Anything, mock.Anything)
}

func TestPatchCourier_InvalidStatus(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	reqBody := handler.RequestUpdate{
		ID:            int64Ptr(1),
		Status:        stringPtr("wrong_status"),
		TransportType: stringPtr(string(model.TransportCar)),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "invalid status")
	svc.AssertNotCalled(t, "PatchCourier", mock.Anything, mock.Anything)
}

func TestPatchCourier_InvalidTransportType(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	reqBody := handler.RequestUpdate{
		ID:            int64Ptr(1),
		Status:        stringPtr(string(model.StatusBusy)),
		TransportType: stringPtr("wrong_transport"),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "invalid transport type")
	svc.AssertNotCalled(t, "PatchCourier", mock.Anything, mock.Anything)
}

func TestPatchCourier_NotFound(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	reqBody := handler.RequestUpdate{
		ID:            int64Ptr(1),
		Status:        stringPtr(string(model.StatusBusy)),
		TransportType: stringPtr(string(model.TransportCar)),
	}
	body, _ := json.Marshal(reqBody)

	svc.On("PatchCourier", mock.Anything, mock.Anything).Return(model.ErrNotFound)

	req := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestPatchCourier_Conflict(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	reqBody := handler.RequestUpdate{
		ID:            int64Ptr(1),
		Status:        stringPtr(string(model.StatusBusy)),
		TransportType: stringPtr(string(model.TransportCar)),
	}
	body, _ := json.Marshal(reqBody)

	svc.On("PatchCourier", mock.Anything, mock.Anything).Return(model.ErrConflict)

	req := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusConflict, rr.Code)
}

func TestPatchCourier_NothingToUpdate(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	reqBody := handler.RequestUpdate{
		ID: int64Ptr(1),
	}
	body, _ := json.Marshal(reqBody)

	svc.On("PatchCourier", mock.Anything, mock.Anything).Return(model.ErrNothingToUpdate)

	req := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestPatchCourier_InternalError(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	reqBody := handler.RequestUpdate{
		ID:            int64Ptr(1),
		Status:        stringPtr(string(model.StatusBusy)),
		TransportType: stringPtr(string(model.TransportCar)),
	}
	body, _ := json.Marshal(reqBody)

	svc.On("PatchCourier", mock.Anything, mock.Anything).Return(errors.New("boom"))

	req := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetCourier_Success(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	id := int64(10)
	name := "John Doe"
	phone := "1234567890"
	status := model.StatusAvalible
	transport := model.TransportOnFoot
	now := time.Now()

	courierModel := model.Courier{
		ID:            &id,
		Name:          &name,
		Phone:         &phone,
		Status:        &status,
		CreatedAt:     &now,
		UpdatedAt:     &now,
		TransportType: &transport,
	}

	svc.On("GetCourier", mock.Anything, id).Return(courierModel, nil)

	req := httptest.NewRequest("GET", "/courier/"+strconv.FormatInt(id, 10), nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp handler.CourierResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.NotNil(t, resp.ID)
	require.Equal(t, id, *resp.ID)
	require.Equal(t, courierModel.Name, resp.Name)
	require.Equal(t, courierModel.Phone, resp.Phone)
	require.Equal(t, string(*courierModel.Status), *resp.Status)
	require.Equal(t, string(*courierModel.TransportType), *resp.TransportType)
}

func TestGetCourier_WrongID(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	req := httptest.NewRequest("GET", "/courier/abc", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "wrong id")
	svc.AssertNotCalled(t, "GetCourier", mock.Anything, mock.Anything)
}

func TestGetCourier_NotFound(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	id := int64(99)

	svc.On("GetCourier", mock.Anything, id).Return(model.Courier{}, model.ErrNotFound)

	req := httptest.NewRequest("GET", "/courier/"+strconv.FormatInt(id, 10), nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetCourier_InternalError(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	id := int64(99)

	svc.On("GetCourier", mock.Anything, id).Return(model.Courier{}, errors.New("boom"))

	req := httptest.NewRequest("GET", "/courier/"+strconv.FormatInt(id, 10), nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetCouriers_Success(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	id1 := int64(1)
	id2 := int64(2)
	name1 := "John"
	name2 := "Jane"
	phone1 := "111"
	phone2 := "222"
	status1 := model.StatusAvalible
	status2 := model.StatusBusy
	tr1 := model.TransportOnFoot
	tr2 := model.TransportCar

	list := []model.Courier{
		{
			ID:            &id1,
			Name:          &name1,
			Phone:         &phone1,
			Status:        &status1,
			TransportType: &tr1,
		},
		{
			ID:            &id2,
			Name:          &name2,
			Phone:         &phone2,
			Status:        &status2,
			TransportType: &tr2,
		},
	}

	svc.On("GetCouriers", mock.Anything).Return(list, nil)

	req := httptest.NewRequest("GET", "/couriers", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp []handler.CourierResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.Len(t, resp, 2)
	require.NotNil(t, resp[0].ID)
	require.Equal(t, id1, *resp[0].ID)
	require.Equal(t, name1, *resp[0].Name)
	require.Equal(t, phone1, *resp[0].Phone)
	require.Equal(t, string(status1), *resp[0].Status)
	require.Equal(t, string(tr1), *resp[0].TransportType)
}

func TestGetCouriers_InternalError(t *testing.T) {
	t.Parallel()

	svc := new(CourierServiceMock)
	h := handler.NewCourierHandler(svc)
	router := setupRouter(h)

	svc.On("GetCouriers", mock.Anything).Return([]model.Courier{}, errors.New("boom"))

	req := httptest.NewRequest("GET", "/couriers", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
