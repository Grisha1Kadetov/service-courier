package delivery

import (
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/delivery"
)

func FromModelToAssignDTO(d delivery.Delivery, c courier.Courier) AssignResponse {
	return AssignResponse{
		CourierID:     d.CourierID,
		OrderID:       d.OrderID,
		TransportType: (*string)(c.TransportType),
		Deadline:      d.Deadline,
	}
}

func FromModelToUnassignDTO(d delivery.Delivery) UnassignResponse {
	s := "unassigned"
	return UnassignResponse{
		OrderID:   d.OrderID,
		Status:    &s,
		CourierID: d.CourierID,
	}
}

func (r AssignRequest) ToModel() delivery.Delivery {
	return delivery.Delivery{
		OrderID: r.OrderID,
	}
}

func (r UnassignRequest) ToModel() delivery.Delivery {
	return delivery.Delivery{
		OrderID: r.OrderID,
	}
}
