package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ShekleinAleksey/subscriptions/internal/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *entity.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, req *entity.UpdateSubscriptionRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entity.Subscription, error)
	GetSummary(ctx context.Context, req *entity.SubscriptionSummaryRequest) (*entity.SubscriptionSummary, error)
}

type subscriptionRepo struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) SubscriptionRepository {
	return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) Create(ctx context.Context, subscription *entity.Subscription) error {
	query := `
        INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

	_, err := r.db.ExecContext(ctx, query,
		subscription.ID,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
		subscription.CreatedAt,
		subscription.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error creating subscription: %v", err)
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	log.Printf("Subscription created successfully: %s", subscription.ID)
	return nil
}

func (r *subscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error) {
	query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions WHERE id = $1
    `

	var subscription entity.Subscription
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subscription.ID,
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.UserID,
		&subscription.StartDate,
		&subscription.EndDate,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subscription not found")
	}
	if err != nil {
		log.Printf("Error getting subscription by ID: %v", err)
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &subscription, nil
}

func (r *subscriptionRepo) Update(ctx context.Context, id uuid.UUID, req *entity.UpdateSubscriptionRequest) error {
	query := "UPDATE subscriptions SET updated_at = $1"
	params := []interface{}{time.Now()}
	paramCount := 2

	if req.ServiceName != nil {
		query += fmt.Sprintf(", service_name = $%d", paramCount)
		params = append(params, *req.ServiceName)
		paramCount++
	}

	if req.Price != nil {
		query += fmt.Sprintf(", price = $%d", paramCount)
		params = append(params, *req.Price)
		paramCount++
	}

	if req.StartDate != nil {
		startDate, err := time.Parse("01-2006", *req.StartDate)
		if err != nil {
			return fmt.Errorf("invalid start_date format: %w", err)
		}
		query += fmt.Sprintf(", start_date = $%d", paramCount)
		params = append(params, startDate)
		paramCount++
	}

	if req.EndDate != nil {
		if *req.EndDate == "" {
			query += fmt.Sprintf(", end_date = $%d", paramCount)
			params = append(params, nil)
		} else {
			endDate, err := time.Parse("01-2006", *req.EndDate)
			if err != nil {
				return fmt.Errorf("invalid end_date format: %w", err)
			}
			query += fmt.Sprintf(", end_date = $%d", paramCount)
			params = append(params, endDate)
		}
		paramCount++
	}

	query += " WHERE id = $" + fmt.Sprint(paramCount)
	params = append(params, id)

	result, err := r.db.ExecContext(ctx, query, params...)
	if err != nil {
		log.Printf("Error updating subscription: %v", err)
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	log.Printf("Subscription updated successfully: %s", id)
	return nil
}

func (r *subscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM subscriptions WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("Error deleting subscription: %v", err)
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	log.Printf("Subscription deleted successfully: %s", id)
	return nil
}

func (r *subscriptionRepo) List(ctx context.Context, limit, offset int) ([]*entity.Subscription, error) {
	query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions 
        ORDER BY created_at DESC 
        LIMIT $1 OFFSET $2
    `

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		log.Printf("Error listing subscriptions: %v", err)
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*entity.Subscription
	for rows.Next() {
		var subscription entity.Subscription
		err := rows.Scan(
			&subscription.ID,
			&subscription.ServiceName,
			&subscription.Price,
			&subscription.UserID,
			&subscription.StartDate,
			&subscription.EndDate,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, &subscription)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating subscriptions: %w", err)
	}

	log.Printf("Listed %d subscriptions", len(subscriptions))
	return subscriptions, nil
}

func (r *subscriptionRepo) GetSummary(ctx context.Context, req *entity.SubscriptionSummaryRequest) (*entity.SubscriptionSummary, error) {
	query := `
        SELECT COALESCE(SUM(price), 0), COUNT(*)
        FROM subscriptions 
        WHERE start_date <= $1 AND (end_date IS NULL OR end_date >= $2)
    `

	params := []interface{}{}

	startPeriod, err := time.Parse("01-2006", req.StartPeriod)
	if err != nil {
		return nil, fmt.Errorf("invalid start_period format: %w", err)
	}

	endPeriod, err := time.Parse("01-2006", req.EndPeriod)
	if err != nil {
		return nil, fmt.Errorf("invalid end_period format: %w", err)
	}

	// Adjust end period to end of month
	endPeriod = time.Date(endPeriod.Year(), endPeriod.Month()+1, 0, 0, 0, 0, 0, time.UTC)

	params = append(params, endPeriod, startPeriod)
	paramCount := 3

	if req.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", paramCount)
		params = append(params, *req.UserID)
		paramCount++
	}

	if req.ServiceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", paramCount)
		params = append(params, *req.ServiceName)
		paramCount++
	}

	var summary entity.SubscriptionSummary
	err = r.db.QueryRowContext(ctx, query, params...).Scan(&summary.TotalCost, &summary.Count)
	if err != nil {
		log.Printf("Error getting subscription summary: %v", err)
		return nil, fmt.Errorf("failed to get subscription summary: %w", err)
	}

	log.Printf("Subscription summary calculated: total cost %d, count %d", summary.TotalCost, summary.Count)
	return &summary, nil
}
