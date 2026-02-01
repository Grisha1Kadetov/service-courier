package courier

import (
	"errors"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
)

var ErrInvalidStatus error = errors.New("invalid status")
var ErrInvalidTransportType error = errors.New("invalid transport type")

type requestMapper interface {
	ToModel() (courier.Courier, error)
}

func FromModelToDTO(c courier.Courier) CourierResponse {
	return CourierResponse{
		ID:            c.ID,
		Name:          c.Name,
		Phone:         c.Phone,
		Status:        (*string)(c.Status),
		TransportType: (*string)(c.TransportType),
	}
}

func FromModelSliceToDTO(couriers []courier.Courier) []CourierResponse {
	result := make([]CourierResponse, len(couriers))
	for i, c := range couriers {
		result[i] = FromModelToDTO(c)
	}
	return result
}

func (c RequestCreate) ToModel() (courier.Courier, error) {
	status, err := convertStringToStatus(c.Status)
	if err != nil {
		return courier.Courier{}, err
	}
	tt, err := convertStringToTransportType(c.TransportType)
	if err != nil {
		return courier.Courier{}, err
	}

	return courier.Courier{
		Name:          c.Name,
		Phone:         c.Phone,
		Status:        status,
		TransportType: tt,
	}, nil
}

func (c RequestUpdate) ToModel() (courier.Courier, error) {
	status, err := convertStringToStatus(c.Status)
	if err != nil {
		return courier.Courier{}, err
	}
	tt, err := convertStringToTransportType(c.TransportType)
	if err != nil {
		return courier.Courier{}, err
	}

	return courier.Courier{
		ID:            c.ID,
		Name:          c.Name,
		Phone:         c.Phone,
		Status:        status,
		TransportType: tt,
	}, nil
}

func convertStringToStatus(s *string) (*courier.Status, error) {
	var status *courier.Status
	if s != nil {
		switch *s {
		case string(courier.StatusAvalible), string(courier.StatusBusy), string(courier.StatusPaused):
			s := courier.Status(*s)
			status = &s
		default:
			return nil, ErrInvalidStatus
		}
	} else {
		status = nil
	}
	return status, nil
}

func convertStringToTransportType(s *string) (*courier.TransportType, error) {
	var tType *courier.TransportType
	if s != nil {
		switch *s {
		case string(courier.TransportCar), string(courier.TransportScooter), string(courier.TransportOnFoot):
			s := courier.TransportType(*s)
			tType = &s
		default:
			return nil, ErrInvalidTransportType
		}
	} else {
		tType = nil
	}
	return tType, nil
}
