package courier

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/courier"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type CourierHandler struct {
	courierService courierService
}

func NewCourierHandler(courierService courierService) *CourierHandler {
	return &CourierHandler{courierService: courierService}
}

func (h *CourierHandler) CreateCourier(w http.ResponseWriter, r *http.Request) {
	var request RequestCreate
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	courier, ok := HandleRequest(w, request)
	if !ok {
		return
	}

	if courier.Name == nil || *courier.Name == "" {
		http.Error(w, "missing name", http.StatusBadRequest)
		return
	}
	if courier.Phone == nil || *courier.Phone == "" {
		http.Error(w, "missing phone", http.StatusBadRequest)
		return
	}
	if courier.Status == nil || *courier.Status == "" {
		http.Error(w, "missing status", http.StatusBadRequest)
		return
	}
	if courier.TransportType == nil || *courier.TransportType == "" {
		http.Error(w, "missing transport type", http.StatusBadRequest)
		return
	}

	err := h.courierService.CreateCourier(r.Context(), courier)
	if err != nil {
		HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *CourierHandler) PatchCourier(w http.ResponseWriter, r *http.Request) {
	var request RequestUpdate
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	courier, ok := HandleRequest(w, request)
	if !ok {
		return
	}

	err := h.courierService.PatchCourier(r.Context(), courier)
	if err != nil {
		HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CourierHandler) GetCourier(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "wrong id", http.StatusBadRequest)
		return
	}

	courier, err := h.courierService.GetCourier(r.Context(), id)
	if err != nil {
		HandleError(w, err)
		return
	}

	render.JSON(w, r, FromModelToDTO(courier))
}

func (h *CourierHandler) GetCouriers(w http.ResponseWriter, r *http.Request) {
	couriers, err := h.courierService.GetCouriers(r.Context())
	if err != nil {
		HandleError(w, err)
		return
	}

	render.JSON(w, r, FromModelSliceToDTO(couriers))
}

func HandleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, courier.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, courier.ErrConflict):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, courier.ErrNothingToUpdate):
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleRequest(w http.ResponseWriter, req requestMapper) (courier.Courier, bool) {
	model, err := req.ToModel()
	if err != nil {
		if errors.Is(err, ErrInvalidStatus) {
			http.Error(w, "invalid status", http.StatusBadRequest)
			return courier.Courier{}, false
		}
		if errors.Is(err, ErrInvalidTransportType) {
			http.Error(w, "invalid transport type", http.StatusBadRequest)
			return courier.Courier{}, false
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return courier.Courier{}, false
	}
	return model, true
}
