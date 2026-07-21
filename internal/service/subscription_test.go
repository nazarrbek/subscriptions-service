package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nazarrbek/subscriptions-service/internal/apperror"
	"github.com/nazarrbek/subscriptions-service/internal/dto"
	"github.com/nazarrbek/subscriptions-service/internal/models"
)

// mockRepo is a test double implementing SubscriptionRepo.
type mockRepo struct {
	createFn         func(ctx context.Context, sub *models.Subscription) error
	getByIDFn        func(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	listFn           func(ctx context.Context, limit, offset int) ([]models.Subscription, error)
	countFn          func(ctx context.Context) (int, error)
	updateFn         func(ctx context.Context, sub *models.Subscription) (*models.Subscription, error)
	deleteFn         func(ctx context.Context, id uuid.UUID) error
	calculateTotalFn func(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to time.Time) (int, error)
}

func (m *mockRepo) Create(ctx context.Context, sub *models.Subscription) error {
	if m.createFn != nil {
		return m.createFn(ctx, sub)
	}
	return nil
}

func (m *mockRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockRepo) List(ctx context.Context, limit, offset int) ([]models.Subscription, error) {
	if m.listFn != nil {
		return m.listFn(ctx, limit, offset)
	}
	return nil, nil
}

func (m *mockRepo) Count(ctx context.Context) (int, error) {
	if m.countFn != nil {
		return m.countFn(ctx)
	}
	return 0, nil
}

func (m *mockRepo) Update(ctx context.Context, sub *models.Subscription) (*models.Subscription, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, sub)
	}
	return nil, nil
}

func (m *mockRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockRepo) CalculateTotal(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to time.Time) (int, error) {
	if m.calculateTotalFn != nil {
		return m.calculateTotalFn(ctx, userID, serviceName, from, to)
	}
	return 0, nil
}

func TestCreate_Success(t *testing.T) {
	var capturedSub *models.Subscription

	repo := &mockRepo{
		createFn: func(_ context.Context, sub *models.Subscription) error {
			capturedSub = sub
			return nil
		},
	}

	svc := NewSubscriptionService(repo)

	req := &dto.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       999,
		UserID:      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		StartDate:   "01-2025",
	}

	result, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.ServiceName != "Netflix" {
		t.Errorf("expected ServiceName=Netflix, got %s", result.ServiceName)
	}

	if result.Price != 999 {
		t.Errorf("expected Price=999, got %d", result.Price)
	}

	if capturedSub == nil {
		t.Fatal("expected repo.Create to be called")
	}

	if capturedSub.ID == uuid.Nil {
		t.Error("expected non-nil UUID")
	}
}

func TestCreate_ValidationError(t *testing.T) {
	repo := &mockRepo{}
	svc := NewSubscriptionService(repo)

	req := &dto.CreateSubscriptionRequest{
		ServiceName: "",
		Price:       999,
		UserID:      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		StartDate:   "01-2025",
	}

	_, err := svc.Create(context.Background(), req)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestCreate_InvalidUserID(t *testing.T) {
	repo := &mockRepo{}
	svc := NewSubscriptionService(repo)

	req := &dto.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       999,
		UserID:      "not-a-uuid",
		StartDate:   "01-2025",
	}

	_, err := svc.Create(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for invalid user_id, got nil")
	}
}

func TestUpdate_Success(t *testing.T) {
	expectedID := uuid.New()
	expectedSub := &models.Subscription{
		ID:          expectedID,
		ServiceName: "Updated",
		Price:       1500,
		UserID:      uuid.New(),
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo := &mockRepo{
		updateFn: func(_ context.Context, _ *models.Subscription) (*models.Subscription, error) {
			return expectedSub, nil
		},
	}

	svc := NewSubscriptionService(repo)

	req := &dto.UpdateSubscriptionRequest{
		ServiceName: "Updated",
		Price:       1500,
		StartDate:   "01-2025",
	}

	result, err := svc.Update(context.Background(), expectedID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ServiceName != "Updated" {
		t.Errorf("expected ServiceName=Updated, got %s", result.ServiceName)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	repo := &mockRepo{
		updateFn: func(_ context.Context, _ *models.Subscription) (*models.Subscription, error) {
			return nil, apperror.ErrNotFound
		},
	}

	svc := NewSubscriptionService(repo)

	req := &dto.UpdateSubscriptionRequest{
		ServiceName: "Updated",
		Price:       1500,
		StartDate:   "01-2025",
	}

	_, err := svc.Update(context.Background(), uuid.New(), req)
	if err == nil {
		t.Fatal("expected ErrNotFound, got nil")
	}

	if err != apperror.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdate_ValidationError(t *testing.T) {
	repo := &mockRepo{}
	svc := NewSubscriptionService(repo)

	req := &dto.UpdateSubscriptionRequest{
		ServiceName: "",
		Price:       999,
		StartDate:   "01-2025",
	}

	_, err := svc.Update(context.Background(), uuid.New(), req)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestList_PassesPagination(t *testing.T) {
	var capturedLimit, capturedOffset int

	repo := &mockRepo{
		listFn: func(_ context.Context, limit, offset int) ([]models.Subscription, error) {
			capturedLimit = limit
			capturedOffset = offset
			return []models.Subscription{}, nil
		},
		countFn: func(_ context.Context) (int, error) {
			return 42, nil
		},
	}

	svc := NewSubscriptionService(repo)

	params := dto.ListParams{Limit: 20, Offset: 10}
	subs, total, err := svc.List(context.Background(), params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if capturedLimit != 20 {
		t.Errorf("expected limit=20, got %d", capturedLimit)
	}
	if capturedOffset != 10 {
		t.Errorf("expected offset=10, got %d", capturedOffset)
	}
	if total != 42 {
		t.Errorf("expected total=42, got %d", total)
	}
	if subs == nil {
		t.Error("expected non-nil subscriptions slice")
	}
}

func TestGetByID_NotFound(t *testing.T) {
	repo := &mockRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*models.Subscription, error) {
			return nil, apperror.ErrNotFound
		},
	}

	svc := NewSubscriptionService(repo)

	_, err := svc.GetByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected ErrNotFound, got nil")
	}

	if err != apperror.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDelete_Success(t *testing.T) {
	repo := &mockRepo{
		deleteFn: func(_ context.Context, _ uuid.UUID) error {
			return nil
		},
	}

	svc := NewSubscriptionService(repo)

	err := svc.Delete(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDelete_NotFound(t *testing.T) {
	repo := &mockRepo{
		deleteFn: func(_ context.Context, _ uuid.UUID) error {
			return apperror.ErrNotFound
		},
	}

	svc := NewSubscriptionService(repo)

	err := svc.Delete(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected ErrNotFound, got nil")
	}

	if err != apperror.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
