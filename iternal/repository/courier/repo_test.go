package courier_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	model "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
	repo "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/repository/courier"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/test/testutils"
)

func setup(t *testing.T) (*repo.CourierRepositoryPG, func()) {
	pool, cleanup := testutils.SetupDB(t)
	r := repo.NewCourierRepository(pool)
	return r, cleanup
}

func TestCourierService_Begin(t *testing.T) {
	t.Parallel()
	repo, cleanup := setup(t)
	defer cleanup()
	ctx := context.Background()

	tx, err := repo.Begin(ctx)
	require.NoError(t, err)
	require.NotNil(t, tx)
	_ = tx.Rollback(ctx)
}

func TestCourierRepo_CreateGetUpdate(t *testing.T) {
	t.Parallel()
	repo, cleanup := setup(t)
	defer cleanup()
	ctx := context.Background()

	m := testCreate(ctx, t, repo)
	id := testGet(ctx, t, repo, m)
	testUpdate(ctx, t, repo, id)
}

func testCreate(ctx context.Context, t *testing.T, repo *repo.CourierRepositoryPG) model.Courier {
	// Create
	m := model.Courier{
		Name:          stringPtr("Иван"),
		Phone:         stringPtr("1234567890"),
		Status:        statusPtr(model.StatusAvalible),
		TransportType: transportPtr(model.TransportCar),
	}
	err := repo.Create(ctx, m, nil)
	require.NoError(t, err)

	// Create duplicate
	err = repo.Create(ctx, m, nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrConflict)

	return m
}

func testGet(ctx context.Context, t *testing.T, repo *repo.CourierRepositoryPG, m model.Courier) *int64 {
	equal := func(expected, actual model.Courier) {
		assert.Equal(t, expected.Name, actual.Name)
		assert.Equal(t, expected.Phone, actual.Phone)
		assert.Equal(t, expected.Status, actual.Status)
		assert.Equal(t, expected.TransportType, actual.TransportType)
	}

	// Get all
	couriers, err := repo.GetAll(ctx)
	require.NoError(t, err)
	require.Len(t, couriers, 1)
	equal(m, couriers[0])

	// Get by ID
	id := couriers[0].ID
	courier, err := repo.GetByID(ctx, *id)
	require.NoError(t, err)
	equal(m, courier)

	// Get by invalid ID
	_, err = repo.GetByID(ctx, 100)
	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrNotFound)
	return id
}

func testUpdate(ctx context.Context, t *testing.T, repo *repo.CourierRepositoryPG, id *int64) {
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
			err := repo.Patch(ctx, tt.model, nil)

			if tt.expectErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)

			got, err := repo.GetByID(ctx, *tt.model.ID)
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
