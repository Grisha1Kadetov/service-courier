package server

import (
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/handler/delivery"
	"github.com/go-chi/chi/v5"
)

func NewRouter(deliveryHandler *delivery.DeliveryHandler, router *chi.Mux) *chi.Mux {
	if router == nil {
		router = chi.NewRouter()
	}

	router.Post("/delivery/assign", deliveryHandler.AssignDelivery)
	router.Post("/delivery/unassign", deliveryHandler.UnassignDelivery)

	return router
}
