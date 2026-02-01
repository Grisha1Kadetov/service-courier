package common

import (
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/handler/common"
	"github.com/go-chi/chi/v5"
)

func NewRouter(router *chi.Mux) *chi.Mux {
	if router == nil {
		router = chi.NewRouter()
	}

	router.Head("/healthcheck", common.HealthCheck)
	router.Get("/ping", common.Ping)

	return router
}
