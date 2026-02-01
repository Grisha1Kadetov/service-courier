package order

import (
	"context"
	"log"
	"time"
)

type OrderFetcher struct {
	gw      orderGateway
	ds      deliveryService
	coursor time.Time
}

func NewOrderFetcher(gw orderGateway, ds deliveryService) *OrderFetcher {
	return &OrderFetcher{
		gw:      gw,
		ds:      ds,
		coursor: time.Time{},
	}
}

func (f *OrderFetcher) Start(ctx context.Context, tick time.Duration) {
	go func() {
		ticker := time.NewTicker(tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := f.fetchOrders(ctx)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
}

func (f *OrderFetcher) fetchOrders(ctx context.Context) error {
	orders, err := f.gw.GetOrders(ctx, f.coursor)
	if err != nil {
		return err
	}
	m := time.Time{}
	for _, order := range orders {
		_, _, err := f.ds.AssignDelivery(ctx, order.ID)
		if err != nil {
			return err
		}
		if m.Before(order.CreatedAt) {
			m = order.CreatedAt
		}
	}
	f.coursor = m
	return nil
}
