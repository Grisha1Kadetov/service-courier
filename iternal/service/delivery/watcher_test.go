package delivery_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/service/delivery"
)

type WatcherRepositoryMock struct {
	mock.Mock
}

func (m *WatcherRepositoryMock) ReleaseExpiredBusyCouriers(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestWatcher_StartWatcherDelivery(t *testing.T) {
	t.Parallel()

	repo := new(WatcherRepositoryMock)
	repo.On("ReleaseExpiredBusyCouriers", mock.Anything).Return(nil)

	watcher := delivery.NewWatcher(repo)

	ctx, cancel := context.WithCancel(context.Background())
	tick := 10 * time.Millisecond
	watcher.RunWatcherDelivery(ctx, tick)

	time.Sleep(35 * time.Millisecond)
	cancel()
	time.Sleep(10 * time.Millisecond)

	repo.AssertCalled(t, "ReleaseExpiredBusyCouriers", mock.Anything)
}
