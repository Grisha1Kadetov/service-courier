package delivery

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/delivery"
	"github.com/go-chi/render"
)

type DeliveryHandler struct {
	deliveryService deliveryService
}

func NewDeliveryHandler(deliveryService deliveryService) *DeliveryHandler {
	return &DeliveryHandler{deliveryService: deliveryService}
}

func (d *DeliveryHandler) UnassignDelivery(w http.ResponseWriter, r *http.Request) {
	var request UnassignRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if request.OrderID == nil || *request.OrderID == "" {
		http.Error(w, "missing order_id", http.StatusBadRequest)
		return
	}
	requstModel := request.ToModel()
	dm, err := d.deliveryService.UnassignDelivery(r.Context(), *requstModel.OrderID)
	if err != nil {
		handlError(w, err)
		return
	}
	response := FromModelToUnassignDTO(dm)
	render.JSON(w, r, response)
}

func (d *DeliveryHandler) AssignDelivery(w http.ResponseWriter, r *http.Request) {
	var request AssignRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if request.OrderID == nil || *request.OrderID == "" {
		http.Error(w, "missing order_id", http.StatusBadRequest)
		return
	}
	requstModel := request.ToModel()
	dm, cm, err := d.deliveryService.AssignDelivery(r.Context(), *requstModel.OrderID)
	if err != nil {
		handlError(w, err)
		return
	}
	response := FromModelToAssignDTO(dm, cm)
	render.JSON(w, r, response)
}

func handlError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, delivery.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, delivery.ErrConflict):
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
