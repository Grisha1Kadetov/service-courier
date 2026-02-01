package common

import (
	"net/http"

	"github.com/go-chi/render"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	render.NoContent(w, r)
}
