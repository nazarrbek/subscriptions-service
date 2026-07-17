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

// Create godoc
// @Summary Create subscription
// @Description Create a new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body dto.CreateSubscriptionRequest true "Subscription"
// @Success 201 {object} dto.CreateSubscriptionResponse
// @Failure 400
// @Failure 500
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdSubscription, err := h.service.Create(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.CreateSubscriptionResponse{ID: createdSubscription.ID.String()})
}

// GetByID godoc
// @Summary Get subscription by ID
// @Description Get subscription by UUID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} models.Subscription
// @Failure 400
// @Failure 404
// @Router /subscriptions/{id} [get]
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

// List godoc
// @Summary List subscriptions
// @Description Get all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} models.Subscription
// @Failure 500
// @Router /subscriptions [get]
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

// Update godoc
// @Summary Update subscription
// @Description Update subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param request body dto.UpdateSubscriptionRequest true "Subscription"
// @Success 200
// @Failure 400
// @Failure 404
// @Router /subscriptions/{id} [put]
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

// Delete godoc
// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 204
// @Failure 400
// @Failure 404
// @Router /subscriptions/{id} [delete]
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

// CalculateTotal godoc
// @Summary Calculate total subscription cost
// @Description Calculate total subscription cost for selected period
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service name"
// @Param from query string true "Start period (MM-YYYY)"
// @Param to query string true "End period (MM-YYYY)"
// @Success 200 {object} dto.TotalResponse
// @Failure 400
// @Failure 500
// @Router /subscriptions/total [get]
func (h *SubscriptionHandler) CalculateTotal(
	w http.ResponseWriter,
	r *http.Request,
) {
	var userID *uuid.UUID
	if rawUserID := r.URL.Query().Get("user_id"); rawUserID != "" {
		parsedUserID, err := uuid.Parse(rawUserID)
		if err != nil {
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}
		userID = &parsedUserID
	}

	var serviceName *string
	if rawServiceName := r.URL.Query().Get("service_name"); rawServiceName != "" {
		serviceName = &rawServiceName
	}

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.TotalResponse{
		Total: total,
	})
}
