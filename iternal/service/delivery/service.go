package delivery

import (
	"context"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/delivery"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/pkg/log"
)

type DeliveryService struct {
	l            log.Logger
	deliveryRepo deliveryRepository
	courierRepo  courierRepository
	timeCalc     DeliveryTimeFactory
}

func NewDeliveryService(deliveryRepo deliveryRepository, courierRepo courierRepository, deliveryTimeFactory DeliveryTimeFactory, logger log.Logger) *DeliveryService {
	return &DeliveryService{
		deliveryRepo: deliveryRepo,
		courierRepo:  courierRepo,
		timeCalc:     deliveryTimeFactory,
		l:            logger,
	}
}

func (s *DeliveryService) AssignDelivery(ctx context.Context, orderId string) (delivery.Delivery, courier.Courier, error) {
	tx, err := s.deliveryRepo.Begin(ctx)
	if err != nil {
		return delivery.Delivery{}, courier.Courier{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	availableCourier, err := s.courierRepo.GetAvailable(ctx)
	if err != nil {
		if err == courier.ErrNotFound {
			s.l.Warn("not found available courier")
			return delivery.Delivery{}, courier.Courier{}, delivery.ErrConflict
		}
		return delivery.Delivery{}, courier.Courier{}, err
	}

	t, err := s.timeCalc.CalculateDeliveryTime(availableCourier)
	if err != nil {
		return delivery.Delivery{}, courier.Courier{}, err
	}
	newDelivery := delivery.Delivery{
		CourierID: availableCourier.ID,
		OrderID:   &orderId,
		Deadline:  &t,
	}
	err = s.deliveryRepo.Create(ctx, newDelivery, tx)
	if err != nil {
		return delivery.Delivery{}, courier.Courier{}, err
	}

	status := courier.StatusBusy
	availableCourier.Status = &status
	err = s.courierRepo.Patch(ctx, availableCourier, tx)
	if err != nil {
		return delivery.Delivery{}, courier.Courier{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return delivery.Delivery{}, courier.Courier{}, err
	}
	s.l.Info("assign delivery", log.NewField("order_id", orderId), log.NewField("courier_id", availableCourier.ID))
	return newDelivery, availableCourier, nil
}

func (s *DeliveryService) UnassignDelivery(ctx context.Context, orderId string) (delivery.Delivery, error) {
	tx, err := s.deliveryRepo.Begin(ctx)
	if err != nil {
		return delivery.Delivery{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	del, err := s.deliveryRepo.GetByOrderID(ctx, orderId)
	if err != nil {
		return delivery.Delivery{}, err
	}

	err = s.deliveryRepo.Delete(ctx, *del.ID, tx)
	if err != nil {
		return delivery.Delivery{}, err
	}

	status := courier.StatusAvalible
	c := courier.Courier{
		ID:     del.CourierID,
		Status: &status,
	}
	err = s.courierRepo.Patch(ctx, c, tx)
	if err != nil {
		return delivery.Delivery{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return delivery.Delivery{}, err
	}
	s.l.Info("unassign delivery", log.NewField("order_id", orderId), log.NewField("courier_id", c.ID))
	return del, nil
}

func (s *DeliveryService) CompleteDelivery(ctx context.Context, orderId string) error {
	del, err := s.deliveryRepo.GetByOrderID(ctx, orderId)
	if err != nil {
		return err
	}

	status := courier.StatusAvalible
	c := courier.Courier{
		ID:     del.CourierID,
		Status: &status,
	}
	err = s.courierRepo.Patch(ctx, c, nil)
	if err != nil {
		return err
	}
	s.l.Info("complete delivery", log.NewField("order_id", orderId), log.NewField("courier_id", c.ID))
	return nil
}
