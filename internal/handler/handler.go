package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nazarrbek/subscriptions-service/internal/apperror"
	"github.com/nazarrbek/subscriptions-service/internal/dto"
	"github.com/nazarrbek/subscriptions-service/internal/models"
)

const (
	defaultLimit = 10
	maxLimit     = 100
)

// SubscriptionServiceI defines the service contract used by the handler layer.
type SubscriptionServiceI interface {
	Create(ctx context.Context, req *dto.CreateSubscriptionRequest) (*models.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	List(ctx context.Context, params dto.ListParams) ([]models.Subscription, int, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateSubscriptionRequest) (*models.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CalculateTotal(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to time.Time) (int, error)
}

type SubscriptionHandler struct {
	service SubscriptionServiceI
}

func NewSubscriptionHandler(service SubscriptionServiceI) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
	}
}

// errorResponse writes a structured JSON error response.
func errorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Create godoc
// @Summary Create subscription
// @Description Create a new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body dto.CreateSubscriptionRequest true "Subscription"
// @Success 201 {object} dto.CreateSubscriptionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	createdSubscription, err := h.service.Create(r.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "validation:") {
			errorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		errorResponse(w, http.StatusInternalServerError, err.Error())
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
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	subscription, err := h.service.GetByID(r.Context(), parsedID)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			errorResponse(w, http.StatusNotFound, "subscription not found")
			return
		}
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)
}

// List godoc
// @Summary List subscriptions
// @Description Get subscriptions with pagination
// @Tags subscriptions
// @Produce json
// @Param limit query int false "Number of records per page (default 10, max 100)"
// @Param offset query int false "Number of records to skip (default 0)"
// @Success 200 {object} dto.ListResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {

	limit := defaultLimit
	offset := 0

	if rawLimit := r.URL.Query().Get("limit"); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil || parsed < 1 {
			errorResponse(w, http.StatusBadRequest, "limit must be a positive integer")
			return
		}
		if parsed > maxLimit {
			parsed = maxLimit
		}
		limit = parsed
	}

	if rawOffset := r.URL.Query().Get("offset"); rawOffset != "" {
		parsed, err := strconv.Atoi(rawOffset)
		if err != nil || parsed < 0 {
			errorResponse(w, http.StatusBadRequest, "offset must be a non-negative integer")
			return
		}
		offset = parsed
	}

	params := dto.ListParams{
		Limit:  limit,
		Offset: offset,
	}

	subscriptions, total, err := h.service.List(r.Context(), params)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if subscriptions == nil {
		subscriptions = []models.Subscription{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ListResponse{
		Data:   subscriptions,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// Update godoc
// @Summary Update subscription
// @Description Update subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param request body dto.UpdateSubscriptionRequest true "Subscription"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := chi.URLParam(r, "id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.UpdateSubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updated, err := h.service.Update(r.Context(), parsedID, &req)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			errorResponse(w, http.StatusNotFound, "subscription not found")
			return
		}
		if strings.Contains(err.Error(), "validation:") {
			errorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// Delete godoc
// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := chi.URLParam(r, "id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.service.Delete(r.Context(), parsedID); err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			errorResponse(w, http.StatusNotFound, "subscription not found")
			return
		}
		errorResponse(w, http.StatusInternalServerError, err.Error())
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
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/total [get]
func (h *SubscriptionHandler) CalculateTotal(
	w http.ResponseWriter,
	r *http.Request,
) {
	var userID *uuid.UUID
	if rawUserID := r.URL.Query().Get("user_id"); rawUserID != "" {
		parsedUserID, err := uuid.Parse(rawUserID)
		if err != nil {
			errorResponse(w, http.StatusBadRequest, "invalid user_id")
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
		errorResponse(w, http.StatusBadRequest, "invalid from parameter, expected MM-YYYY format")
		return
	}

	to, err := time.Parse("01-2006", r.URL.Query().Get("to"))
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid to parameter, expected MM-YYYY format")
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
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.TotalResponse{
		Total: total,
	})
}
