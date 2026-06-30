package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/nazarrbek/subscriptions-service/internal/dto"
	"github.com/nazarrbek/subscriptions-service/internal/models"
	"github.com/nazarrbek/subscriptions-service/internal/repository"
)

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
	}
}

func (s *SubscriptionService) Create(
	ctx context.Context,
	req *dto.CreateSubscriptionRequest,
) error {

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start date: %w", err)
	}

	var endDate *time.Time
	if req.EndDate != "" {
		parsedEndDate, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			return fmt.Errorf("invalid end date: %w", err)
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

	return s.repo.Create(ctx, subscription)
}

func (s *SubscriptionService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Subscription, error) {

	return s.repo.GetByID(ctx, id)
}
