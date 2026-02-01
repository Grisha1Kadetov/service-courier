package delivery_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
	deliveryModel "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/delivery"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/service/delivery"
)

func ptr[T any](v T) *T { return &v }

func TestDeliveryTimeFactory_CalculateDeliveryTime(t *testing.T) {
	t.Parallel()
	f := &delivery.DefaultDeliveryTimeFactory{}

	tests := []struct {
		name       string
		courier    courier.Courier
		addMinutes time.Duration
		expectErr  error
	}{
		{
			name:       "car adds 5 minutes",
			courier:    courier.Courier{TransportType: ptr(courier.TransportCar)},
			addMinutes: 5,
		},
		{
			name:       "scooter adds 15 minutes",
			courier:    courier.Courier{TransportType: ptr(courier.TransportScooter)},
			addMinutes: 15,
		},
		{
			name:       "on foot adds 30 minutes",
			courier:    courier.Courier{TransportType: ptr(courier.TransportOnFoot)},
			addMinutes: 30,
		},
		{
			name:      "nil transport returns error",
			courier:   courier.Courier{},
			expectErr: deliveryModel.ErrCannotCalculateDeliveryTime,
		},
	}

	allowedDrift := time.Minute

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()

			result, err := f.CalculateDeliveryTime(tt.courier)

			if tt.expectErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
				return
			}

			require.NoError(t, err)

			expected := now.Add(tt.addMinutes * time.Minute)
			assert.WithinDuration(t, expected, result, allowedDrift)
		})
	}
}
