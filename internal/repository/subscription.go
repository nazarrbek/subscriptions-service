package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nazarrbek/subscriptions-service/internal/apperror"
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, fmt.Errorf("get subscription by id: %w", err)
	}

	return &sub, nil
}

func (r *SubscriptionRepository) List(
	ctx context.Context,
	limit, offset int,
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
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
`

	rows, err := r.db.Query(ctx, query, limit, offset)
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

// Count returns the total number of subscriptions.
func (r *SubscriptionRepository) Count(ctx context.Context) (int, error) {
	const query = `SELECT COUNT(*) FROM subscriptions`

	var total int
	if err := r.db.QueryRow(ctx, query).Scan(&total); err != nil {
		return 0, fmt.Errorf("count subscriptions: %w", err)
	}
	return total, nil
}

func (r *SubscriptionRepository) Update(
	ctx context.Context,
	sub *models.Subscription,
) (*models.Subscription, error) {

	const query = `
UPDATE subscriptions
SET
	service_name = $1,
	price = $2,
	start_date = $3,
	end_date = $4,
	updated_at = CURRENT_TIMESTAMP
WHERE id = $5
RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at;
`

	var updated models.Subscription

	err := r.db.QueryRow(
		ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	).Scan(
		&updated.ID,
		&updated.ServiceName,
		&updated.Price,
		&updated.UserID,
		&updated.StartDate,
		&updated.EndDate,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, fmt.Errorf("update subscription: %w", err)
	}

	return &updated, nil
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
		return apperror.ErrNotFound
	}

	return nil
}

func (r *SubscriptionRepository) CalculateTotal(
	ctx context.Context,
	userID *uuid.UUID,
	serviceName *string,
	from time.Time,
	to time.Time,
) (int, error) {

	const query = `
SELECT
	price,
	start_date,
	end_date
FROM subscriptions
WHERE start_date <= $1
AND (end_date IS NULL OR end_date >= $2)
AND ($3::uuid IS NULL OR user_id = $3)
AND ($4::text IS NULL OR service_name = $4)
ORDER BY start_date;
`

	var userIDArg any
	if userID != nil {
		userIDArg = *userID
	}

	var serviceNameArg any
	if serviceName != nil {
		serviceNameArg = *serviceName
	}

	rows, err := r.db.Query(ctx, query, to, from, userIDArg, serviceNameArg)
	if err != nil {
		return 0, fmt.Errorf("calculate total: %w", err)
	}
	defer rows.Close()

	var total int
	queryFrom := monthStart(from)
	queryTo := monthStart(to)

	for rows.Next() {
		var price int
		var startDate time.Time
		var endDate *time.Time

		if err := rows.Scan(&price, &startDate, &endDate); err != nil {
			return 0, fmt.Errorf("scan subscription for total: %w", err)
		}

		overlapStart := maxMonth(monthStart(startDate), queryFrom)
		overlapEnd := queryTo
		if endDate != nil {
			candidateEnd := monthStart(*endDate)
			if candidateEnd.Before(overlapEnd) {
				overlapEnd = candidateEnd
			}
		}

		if overlapEnd.Before(overlapStart) {
			continue
		}

		months := monthsInclusive(overlapStart, overlapEnd)
		total += price * months
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("calculate total rows: %w", err)
	}

	return total, nil
}

func monthStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func maxMonth(first time.Time, second time.Time) time.Time {
	if second.After(first) {
		return second
	}

	return first
}

func monthsInclusive(start time.Time, end time.Time) int {
	if end.Before(start) {
		return 0
	}

	return (end.Year()-start.Year())*12 + int(end.Month()-start.Month()) + 1
}
