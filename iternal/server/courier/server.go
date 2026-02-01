package server

import (
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/handler/courier"
	"github.com/go-chi/chi/v5"
)

func NewRouter(courierHandler *courier.CourierHandler, router *chi.Mux) *chi.Mux {
	if router == nil {
		router = chi.NewRouter()
	}

	router.Get("/courier/{id}", courierHandler.GetCourier)
	router.Get("/couriers", courierHandler.GetCouriers)
	router.Post("/courier", courierHandler.CreateCourier)
	router.Put("/courier", courierHandler.PatchCourier)

	return router
}
