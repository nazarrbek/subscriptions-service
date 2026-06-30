package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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
func (r *SubscriptionRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Subscription, error) {

	const query = `
SELECT
	id,
	service_name,
	price,
	user_id,
	start_date,
	end_date,
	created_at,
	updated_at
FROM subscriptions
WHERE id = $1`

	var sub models.Subscription

	err := r.db.QueryRow(ctx, query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get subscription by id: %w", err)
	}

	return &sub, nil
}

func (r *SubscriptionRepository) List(
	ctx context.Context,
) ([]models.Subscription, error) {

	const query = `
SELECT
	id,
	service_name,
	price,
	user_id,
	start_date,
	end_date,
	created_at,
	updated_at
FROM subscriptions
ORDER BY created_at DESC;
`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []models.Subscription

	for rows.Next() {
		var sub models.Subscription

		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan subscription: %w", err)
		}

		subscriptions = append(subscriptions, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return subscriptions, nil
}
