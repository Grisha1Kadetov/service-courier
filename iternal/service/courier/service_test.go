package courier_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	model "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
	repo "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/repository/courier"
	service "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/service/courier"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/test/testutils"
)

func setup(t *testing.T) (*service.CourierService, *repo.CourierRepositoryPG, func()) {
	pool, cleanup := testutils.SetupDB(t)
	r := repo.NewCourierRepository(pool)
	s := service.NewCourierService(r)
	return s, r, cleanup
}
func TestCourierService_CreateGetUpdate(t *testing.T) {
	t.Parallel()
	service, _, cleanup := setup(t)
	defer cleanup()
	ctx := context.Background()

	m := testCreate(ctx, t, service)
	id := testGet(ctx, t, service, m)
	testUpdate(ctx, t, service, id)
}

func testCreate(ctx context.Context, t *testing.T, service *service.CourierService) model.Courier {
	// Create
	m := model.Courier{
		Name:          stringPtr("Иван"),
		Phone:         stringPtr("1234567890"),
		Status:        statusPtr(model.StatusAvalible),
		TransportType: transportPtr(model.TransportCar),
	}
	err := service.CreateCourier(ctx, m)
	require.NoError(t, err)

	// Create duplicate
	err = service.CreateCourier(ctx, m)
	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrConflict)

	return m
}

func testGet(ctx context.Context, t *testing.T, service *service.CourierService, m model.Courier) *int64 {
	equal := func(expected, actual model.Courier) {
		assert.Equal(t, expected.Name, actual.Name)
		assert.Equal(t, expected.Phone, actual.Phone)
		assert.Equal(t, expected.Status, actual.Status)
		assert.Equal(t, expected.TransportType, actual.TransportType)
	}

	// Get all
	couriers, err := service.GetCouriers(ctx)
	require.NoError(t, err)
	require.Len(t, couriers, 1)
	equal(m, couriers[0])

	// Get by ID
	id := couriers[0].ID
	courier, err := service.GetCourier(ctx, *id)
	require.NoError(t, err)
	equal(m, courier)

	// Get by invalid ID
	_, err = service.GetCourier(ctx, 100)
	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrNotFound)
	return id
}

func testUpdate(ctx context.Context, t *testing.T, service *service.CourierService, id *int64) {
	tests := []struct {
		name      string
		model     model.Courier
		actual    func(c model.Courier) any
		expect    any
		expectErr error
	}{
		{
			name:   "update name",
			model:  model.Courier{ID: id, Name: stringPtr("Ivan")},
			actual: func(c model.Courier) any { return c.Name },
			expect: stringPtr("Ivan"),
		},
		{
			name:   "update status",
			model:  model.Courier{ID: id, Status: statusPtr(model.StatusPaused)},
			actual: func(c model.Courier) any { return c.Status },
			expect: statusPtr(model.StatusPaused),
		},
		{
			name:   "update phone",
			model:  model.Courier{ID: id, Phone: stringPtr("+1111111111")},
			actual: func(c model.Courier) any { return c.Phone },
			expect: stringPtr("+1111111111"),
		},
		{
			name:   "update transport",
			model:  model.Courier{ID: id, TransportType: transportPtr(model.TransportCar)},
			actual: func(c model.Courier) any { return c.TransportType },
			expect: transportPtr(model.TransportCar),
		},
		{
			name:      "not found",
			model:     model.Courier{ID: int64Ptr(100), Name: stringPtr("Ivan")},
			expectErr: model.ErrNotFound,
		},
		{
			name:      "nothing to update",
			model:     model.Courier{ID: id},
			expectErr: model.ErrNothingToUpdate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.PatchCourier(ctx, tt.model)

			if tt.expectErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)

			got, err := service.GetCourier(ctx, *tt.model.ID)
			require.NoError(t, err)
			value := tt.actual(got)
			assert.Equal(t, tt.expect, value)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}

func statusPtr(s model.Status) *model.Status {
	return &s
}

func transportPtr(t model.TransportType) *model.TransportType {
	return &t
}
