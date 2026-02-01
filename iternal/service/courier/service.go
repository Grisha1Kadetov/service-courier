package courier

import (
	"context"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
)

type CourierService struct {
	courierRepo courierRepository
}

func NewCourierService(courierRepo courierRepository) *CourierService {
	return &CourierService{courierRepo: courierRepo}
}

func (s *CourierService) CreateCourier(ctx context.Context, courier courier.Courier) error {
	err := s.courierRepo.Create(ctx, courier, nil)
	return err
}

func (s *CourierService) PatchCourier(ctx context.Context, courier courier.Courier) error {
	err := s.courierRepo.Patch(ctx, courier, nil)
	return err
}

func (s *CourierService) GetCourier(ctx context.Context, id int64) (courier.Courier, error) {
	courier, err := s.courierRepo.GetByID(ctx, id)
	return courier, err
}

func (s *CourierService) GetCouriers(ctx context.Context) ([]courier.Courier, error) {
	couriers, err := s.courierRepo.GetAll(ctx)
	return couriers, err
}
