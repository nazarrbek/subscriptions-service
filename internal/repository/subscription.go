package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/nazarrbek/subscriptions-service/internal/models"
)

type SubscriptionRepository struct {
	db *pgx.Conn
}

func NewSubscriptionRepository(db *pgx.Conn) *SubscriptionRepository {
	return &SubscriptionRepository{
		db: db,
	}
}

func (r *SubscriptionRepository) Create(
	ctx context.Context,
	subscription *models.Subscription,
) error {
	const query = `INSERT INTO subscriptions (
    id,
    service_name,
    price,
    user_id,
    start_date,
    end_date
)
VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(ctx,
		query,
		subscription.ID,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
	)
	if err != nil {
		return fmt.Errorf("create subscription: %w", err)
	}
	return nil
}
