package delivery

import (
	"time"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/delivery"
)

type DeliveryTimeFactory interface {
	CalculateDeliveryTime(courier.Courier) (time.Time, error)
}

type DefaultDeliveryTimeFactory struct{}

func (f *DefaultDeliveryTimeFactory) CalculateDeliveryTime(c courier.Courier) (time.Time, error) {
	if c.TransportType == nil {
		return time.Time{}, delivery.ErrCannotCalculateDeliveryTime
	}

	t := *c.TransportType
	switch t {
	case courier.TransportCar:
		return time.Now().Add(5 * time.Minute), nil
	case courier.TransportScooter:
		return time.Now().Add(15 * time.Minute), nil
	case courier.TransportOnFoot:
		return time.Now().Add(30 * time.Minute), nil
	}
	return time.Time{}, delivery.ErrCannotCalculateDeliveryTime
}
