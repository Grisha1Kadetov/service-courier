package courier_test

import (
	"testing"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/handler/courier"
	model "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
	"github.com/stretchr/testify/require"
)

func TestConverter(t *testing.T) {
	t.Parallel()
	rCreate := courier.RequestCreate{
		Name:          stringPtr("John Doe"),
		Phone:         stringPtr("1234567890"),
		Status:        stringPtr(string(model.StatusAvalible)),
		TransportType: stringPtr(string(model.TransportOnFoot)),
	}
	m, err := rCreate.ToModel()
	require.NoError(t, err)

	testCourierStatus(t, rCreate)
	testCourierTransportType(t, rCreate)

	require.Equal(t, rCreate.Name, m.Name)
	require.Equal(t, rCreate.Phone, m.Phone)

	rUpdate := courier.RequestUpdate{
		ID:            int64Ptr(1),
		Name:          stringPtr("Jane Doe"),
		Phone:         stringPtr("1234567890"),
		Status:        stringPtr(string(model.StatusBusy)),
		TransportType: stringPtr(string(model.TransportCar)),
	}
	m, err = rUpdate.ToModel()
	require.NoError(t, err)
	require.Equal(t, rUpdate.ID, m.ID)
	require.Equal(t, rUpdate.Name, m.Name)
	require.Equal(t, rUpdate.Phone, m.Phone)
}

func testCourierStatus(t *testing.T, request courier.RequestCreate) {
	start := request.Status
	request.Status = stringPtr(string(model.StatusAvalible))
	_, err := request.ToModel()
	require.NoError(t, err)

	request.Status = stringPtr(string(model.StatusBusy))
	_, err = request.ToModel()
	require.NoError(t, err)

	request.Status = stringPtr(string(model.StatusPaused))
	_, err = request.ToModel()
	require.NoError(t, err)

	request.Status = stringPtr("invalid_status")
	_, err = request.ToModel()
	require.ErrorIs(t, err, courier.ErrInvalidStatus)

	request.Status = start
}

func testCourierTransportType(t *testing.T, request courier.RequestCreate) {
	start := request.TransportType
	request.TransportType = stringPtr(string(model.TransportOnFoot))
	_, err := request.ToModel()
	require.NoError(t, err)

	request.TransportType = stringPtr(string(model.TransportCar))
	_, err = request.ToModel()
	require.NoError(t, err)

	request.TransportType = stringPtr(string(model.TransportScooter))
	_, err = request.ToModel()
	require.NoError(t, err)

	request.TransportType = stringPtr("invalid_transport")
	_, err = request.ToModel()
	require.ErrorIs(t, err, courier.ErrInvalidTransportType)

	request.TransportType = start
}

func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}
