package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nazarrbek/subscriptions-service/internal/dto"
	"github.com/nazarrbek/subscriptions-service/internal/service"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

func NewSubscriptionHandler(service *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
	}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	subscription, err := h.service.GetByID(r.Context(), parsedID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(subscription)
}

func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {

	subscriptions, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(subscriptions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *SubscriptionHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := chi.URLParam(r, "id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateSubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Update(r.Context(), parsedID, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SubscriptionHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := chi.URLParam(r, "id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), parsedID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) CalculateTotal(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, err := uuid.Parse(r.URL.Query().Get("user_id"))
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	serviceName := r.URL.Query().Get("service_name")

	from, err := time.Parse("01-2006", r.URL.Query().Get("from"))
	if err != nil {
		http.Error(w, "invalid from", http.StatusBadRequest)
		return
	}

	to, err := time.Parse("01-2006", r.URL.Query().Get("to"))
	if err != nil {
		http.Error(w, "invalid to", http.StatusBadRequest)
		return
	}

	total, err := h.service.CalculateTotal(
		r.Context(),
		userID,
		serviceName,
		from,
		to,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(dto.TotalResponse{
		Total: total,
	})
}
