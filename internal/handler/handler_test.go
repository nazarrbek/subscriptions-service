package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nazarrbek/subscriptions-service/internal/apperror"
	"github.com/nazarrbek/subscriptions-service/internal/dto"
	"github.com/nazarrbek/subscriptions-service/internal/models"
)

// mockService is a test double implementing SubscriptionServiceI.
type mockService struct {
	createFn         func(ctx context.Context, req *dto.CreateSubscriptionRequest) (*models.Subscription, error)
	getByIDFn        func(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	listFn           func(ctx context.Context, params dto.ListParams) ([]models.Subscription, int, error)
	updateFn         func(ctx context.Context, id uuid.UUID, req *dto.UpdateSubscriptionRequest) (*models.Subscription, error)
	deleteFn         func(ctx context.Context, id uuid.UUID) error
	calculateTotalFn func(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to time.Time) (int, error)
}

func (m *mockService) Create(ctx context.Context, req *dto.CreateSubscriptionRequest) (*models.Subscription, error) {
	if m.createFn != nil {
		return m.createFn(ctx, req)
	}
	return nil, nil
}

func (m *mockService) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockService) List(ctx context.Context, params dto.ListParams) ([]models.Subscription, int, error) {
	if m.listFn != nil {
		return m.listFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *mockService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateSubscriptionRequest) (*models.Subscription, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, req)
	}
	return nil, nil
}

func (m *mockService) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockService) CalculateTotal(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to time.Time) (int, error) {
	if m.calculateTotalFn != nil {
		return m.calculateTotalFn(ctx, userID, serviceName, from, to)
	}
	return 0, nil
}

// newChiRequest creates an HTTP request with chi URL params set.
func newChiRequest(method, target string, body []byte, params map[string]string) *http.Request {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, target, bytes.NewReader(body))
	} else {
		req = httptest.NewRequest(method, target, nil)
	}

	if len(params) > 0 {
		rctx := chi.NewRouteContext()
		for k, v := range params {
			rctx.URLParams.Add(k, v)
		}
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	}

	return req
}

func TestCreate_Success(t *testing.T) {
	subID := uuid.New()
	svc := &mockService{
		createFn: func(_ context.Context, _ *dto.CreateSubscriptionRequest) (*models.Subscription, error) {
			return &models.Subscription{ID: subID}, nil
		},
	}

	h := NewSubscriptionHandler(svc)
	body, _ := json.Marshal(dto.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       999,
		UserID:      uuid.New().String(),
		StartDate:   "01-2025",
	})

	req := newChiRequest(http.MethodPost, "/subscriptions", body, nil)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Create(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}

	var resp dto.CreateSubscriptionResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.ID != subID.String() {
		t.Errorf("expected ID=%s, got %s", subID.String(), resp.ID)
	}
}

func TestCreate_InvalidJSON(t *testing.T) {
	svc := &mockService{}
	h := NewSubscriptionHandler(svc)

	req := newChiRequest(http.MethodPost, "/subscriptions", []byte("{invalid"), nil)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Create(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestGetByID_Success(t *testing.T) {
	subID := uuid.New()
	expected := &models.Subscription{
		ID:          subID,
		ServiceName: "Netflix",
		Price:       999,
		UserID:      uuid.New(),
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	svc := &mockService{
		getByIDFn: func(_ context.Context, id uuid.UUID) (*models.Subscription, error) {
			if id != subID {
				t.Errorf("expected id=%s, got %s", subID, id)
			}
			return expected, nil
		},
	}

	h := NewSubscriptionHandler(svc)

	req := newChiRequest(http.MethodGet, "/subscriptions/"+subID.String(), nil, map[string]string{"id": subID.String()})
	rr := httptest.NewRecorder()

	h.GetByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var resp models.Subscription
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.ServiceName != "Netflix" {
		t.Errorf("expected ServiceName=Netflix, got %s", resp.ServiceName)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	svc := &mockService{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*models.Subscription, error) {
			return nil, apperror.ErrNotFound
		},
	}

	h := NewSubscriptionHandler(svc)
	subID := uuid.New()

	req := newChiRequest(http.MethodGet, "/subscriptions/"+subID.String(), nil, map[string]string{"id": subID.String()})
	rr := httptest.NewRecorder()

	h.GetByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestGetByID_InvalidID(t *testing.T) {
	svc := &mockService{}
	h := NewSubscriptionHandler(svc)

	req := newChiRequest(http.MethodGet, "/subscriptions/not-a-uuid", nil, map[string]string{"id": "not-a-uuid"})
	rr := httptest.NewRecorder()

	h.GetByID(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestUpdate_Success(t *testing.T) {
	subID := uuid.New()
	updated := &models.Subscription{
		ID:          subID,
		ServiceName: "Updated",
		Price:       1500,
		UserID:      uuid.New(),
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	svc := &mockService{
		updateFn: func(_ context.Context, _ uuid.UUID, _ *dto.UpdateSubscriptionRequest) (*models.Subscription, error) {
			return updated, nil
		},
	}

	h := NewSubscriptionHandler(svc)
	body, _ := json.Marshal(dto.UpdateSubscriptionRequest{
		ServiceName: "Updated",
		Price:       1500,
		StartDate:   "01-2025",
	})

	req := newChiRequest(http.MethodPut, "/subscriptions/"+subID.String(), body, map[string]string{"id": subID.String()})
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Update(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var resp models.Subscription
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.ServiceName != "Updated" {
		t.Errorf("expected ServiceName=Updated, got %s", resp.ServiceName)
	}

	if resp.ID != subID {
		t.Errorf("expected ID=%s, got %s", subID, resp.ID)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	svc := &mockService{
		updateFn: func(_ context.Context, _ uuid.UUID, _ *dto.UpdateSubscriptionRequest) (*models.Subscription, error) {
			return nil, apperror.ErrNotFound
		},
	}

	h := NewSubscriptionHandler(svc)
	subID := uuid.New()
	body, _ := json.Marshal(dto.UpdateSubscriptionRequest{
		ServiceName: "Updated",
		Price:       1500,
		StartDate:   "01-2025",
	})

	req := newChiRequest(http.MethodPut, "/subscriptions/"+subID.String(), body, map[string]string{"id": subID.String()})
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Update(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestUpdate_InvalidJSON(t *testing.T) {
	svc := &mockService{}
	h := NewSubscriptionHandler(svc)
	subID := uuid.New()

	req := newChiRequest(http.MethodPut, "/subscriptions/"+subID.String(), []byte("{bad"), map[string]string{"id": subID.String()})
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Update(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestDelete_Success(t *testing.T) {
	svc := &mockService{
		deleteFn: func(_ context.Context, _ uuid.UUID) error {
			return nil
		},
	}

	h := NewSubscriptionHandler(svc)
	subID := uuid.New()

	req := newChiRequest(http.MethodDelete, "/subscriptions/"+subID.String(), nil, map[string]string{"id": subID.String()})
	rr := httptest.NewRecorder()

	h.Delete(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rr.Code)
	}
}

func TestDelete_NotFound(t *testing.T) {
	svc := &mockService{
		deleteFn: func(_ context.Context, _ uuid.UUID) error {
			return apperror.ErrNotFound
		},
	}

	h := NewSubscriptionHandler(svc)
	subID := uuid.New()

	req := newChiRequest(http.MethodDelete, "/subscriptions/"+subID.String(), nil, map[string]string{"id": subID.String()})
	rr := httptest.NewRecorder()

	h.Delete(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestList_Success(t *testing.T) {
	svc := &mockService{
		listFn: func(_ context.Context, params dto.ListParams) ([]models.Subscription, int, error) {
			if params.Limit != 10 {
				t.Errorf("expected default limit=10, got %d", params.Limit)
			}
			if params.Offset != 0 {
				t.Errorf("expected default offset=0, got %d", params.Offset)
			}
			return []models.Subscription{
				{ID: uuid.New(), ServiceName: "Netflix", Price: 999},
			}, 1, nil
		},
	}

	h := NewSubscriptionHandler(svc)

	req := newChiRequest(http.MethodGet, "/subscriptions", nil, nil)
	rr := httptest.NewRecorder()

	h.List(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var resp dto.ListResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Total != 1 {
		t.Errorf("expected total=1, got %d", resp.Total)
	}
	if resp.Limit != 10 {
		t.Errorf("expected limit=10, got %d", resp.Limit)
	}
	if resp.Offset != 0 {
		t.Errorf("expected offset=0, got %d", resp.Offset)
	}
	if len(resp.Data) != 1 {
		t.Errorf("expected 1 item, got %d", len(resp.Data))
	}
}

func TestList_CustomPagination(t *testing.T) {
	svc := &mockService{
		listFn: func(_ context.Context, params dto.ListParams) ([]models.Subscription, int, error) {
			if params.Limit != 25 {
				t.Errorf("expected limit=25, got %d", params.Limit)
			}
			if params.Offset != 50 {
				t.Errorf("expected offset=50, got %d", params.Offset)
			}
			return []models.Subscription{}, 100, nil
		},
	}

	h := NewSubscriptionHandler(svc)

	req := newChiRequest(http.MethodGet, "/subscriptions?limit=25&offset=50", nil, nil)
	rr := httptest.NewRecorder()

	h.List(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestList_InvalidLimit(t *testing.T) {
	svc := &mockService{}
	h := NewSubscriptionHandler(svc)

	req := newChiRequest(http.MethodGet, "/subscriptions?limit=-5", nil, nil)
	rr := httptest.NewRecorder()

	h.List(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestList_InvalidOffset(t *testing.T) {
	svc := &mockService{}
	h := NewSubscriptionHandler(svc)

	req := newChiRequest(http.MethodGet, "/subscriptions?offset=-1", nil, nil)
	rr := httptest.NewRecorder()

	h.List(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestList_LimitCappedAtMax(t *testing.T) {
	svc := &mockService{
		listFn: func(_ context.Context, params dto.ListParams) ([]models.Subscription, int, error) {
			if params.Limit != maxLimit {
				t.Errorf("expected limit capped at %d, got %d", maxLimit, params.Limit)
			}
			return []models.Subscription{}, 0, nil
		},
	}

	h := NewSubscriptionHandler(svc)

	req := newChiRequest(http.MethodGet, "/subscriptions?limit=999", nil, nil)
	rr := httptest.NewRecorder()

	h.List(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestErrorResponseFormat(t *testing.T) {
	svc := &mockService{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*models.Subscription, error) {
			return nil, apperror.ErrNotFound
		},
	}

	h := NewSubscriptionHandler(svc)
	subID := uuid.New()

	req := newChiRequest(http.MethodGet, "/subscriptions/"+subID.String(), nil, map[string]string{"id": subID.String()})
	rr := httptest.NewRecorder()

	h.GetByID(rr, req)

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type=application/json, got %s", contentType)
	}

	var errResp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("expected JSON error response, got parse error: %v", err)
	}

	if _, ok := errResp["error"]; !ok {
		t.Error("expected 'error' key in JSON response")
	}
}
