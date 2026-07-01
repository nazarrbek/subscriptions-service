package repository

import (
	"context"
	"fmt"
	"time"

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

func (r *SubscriptionRepository) Update(
	ctx context.Context,
	sub *models.Subscription,
) error {

	const query = `
UPDATE subscriptions
SET
	service_name = $1,
	price = $2,
	start_date = $3,
	end_date = $4,
	updated_at = CURRENT_TIMESTAMP
WHERE id = $5;
`

	_, err := r.db.Exec(
		ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)

	if err != nil {
		return fmt.Errorf("update subscription: %w", err)
	}

	return nil
}

func (r *SubscriptionRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {

	const query = `
DELETE FROM subscriptions
WHERE id = $1;
`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

func (r *SubscriptionRepository) CalculateTotal(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
	from time.Time,
	to time.Time,
) (int, error) {

	const query = `
SELECT COALESCE(SUM(price), 0)
FROM subscriptions
WHERE user_id = $1
AND service_name = $2
AND start_date BETWEEN $3 AND $4;
`

	var total int

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
		serviceName,
		from,
		to,
	).Scan(&total)

	if err != nil {
		return 0, fmt.Errorf("calculate total: %w", err)
	}

	return total, nil
}
