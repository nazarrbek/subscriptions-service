package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/nazarrbek/subscriptions-service/internal/dto"
	"github.com/nazarrbek/subscriptions-service/internal/models"
)

// SubscriptionRepo defines the repository contract used by the service layer.
// This interface enables dependency injection and unit testing with mocks.
type SubscriptionRepo interface {
	Create(ctx context.Context, subscription *models.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	List(ctx context.Context, limit, offset int) ([]models.Subscription, error)
	Count(ctx context.Context) (int, error)
	Update(ctx context.Context, sub *models.Subscription) (*models.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CalculateTotal(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to time.Time) (int, error)
}

type SubscriptionService struct {
	repo SubscriptionRepo
}

func NewSubscriptionService(repo SubscriptionRepo) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
	}
}

func (s *SubscriptionService) Create(
	ctx context.Context,
	req *dto.CreateSubscriptionRequest,
) (*models.Subscription, error) {

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	var endDate *time.Time
	if req.EndDate != "" {
		parsedEndDate, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date: %w", err)
		}
		endDate = &parsedEndDate
	}

	subscription := &models.Subscription{
		ID:          uuid.New(),
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := s.repo.Create(ctx, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *SubscriptionService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Subscription, error) {

	return s.repo.GetByID(ctx, id)
}

func (s *SubscriptionService) List(
	ctx context.Context,
	params dto.ListParams,
) ([]models.Subscription, int, error) {

	subscriptions, err := s.repo.List(ctx, params.Limit, params.Offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return subscriptions, total, nil
}

func (s *SubscriptionService) Update(
	ctx context.Context,
	id uuid.UUID,
	req *dto.UpdateSubscriptionRequest,
) (*models.Subscription, error) {

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	var endDate *time.Time

	if req.EndDate != "" {
		t, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date: %w", err)
		}
		endDate = &t
	}

	sub := &models.Subscription{
		ID:          id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	return s.repo.Update(ctx, sub)
}

func (s *SubscriptionService) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {

	return s.repo.Delete(ctx, id)
}

func (s *SubscriptionService) CalculateTotal(
	ctx context.Context,
	userID *uuid.UUID,
	serviceName *string,
	from time.Time,
	to time.Time,
) (int, error) {

	return s.repo.CalculateTotal(
		ctx,
		userID,
		serviceName,
		from,
		to,
	)
}
