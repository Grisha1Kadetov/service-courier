package delivery

import (
	"context"
	"log"
	"time"
)

type Watcher struct {
	watcherRepository watcherRepository
}

func NewWatcher(watcherRepository watcherRepository) *Watcher {
	return &Watcher{
		watcherRepository: watcherRepository,
	}
}

func (w *Watcher) RunWatcherDelivery(ctx context.Context, tick time.Duration) {
	go func() {
		ticker := time.NewTicker(tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := w.watcherRepository.ReleaseExpiredBusyCouriers(ctx)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
	}()
}
