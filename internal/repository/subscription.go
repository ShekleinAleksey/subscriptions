package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ShekleinAleksey/subscriptions/internal/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
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
        INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	_, err := r.db.ExecContext(ctx, query,
		subscription.ID,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
	)

	if err != nil {
		logrus.Fatalf("Error creating subscription: %v", err)
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	logrus.Info("Subscription created successfully: %s", subscription.ID)
	return nil
}

func (r *subscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error) {
	query := `
        SELECT *
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
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subscription not found")
	}
	if err != nil {
		logrus.Fatalf("Error getting subscription by ID: %v", err)
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &subscription, nil
}

func (r *subscriptionRepo) Update(ctx context.Context, id uuid.UUID, req *entity.UpdateSubscriptionRequest) error {
	query := "UPDATE subscriptions SET"
	params := []interface{}{}
	paramCount := 1

	if req.ServiceName != nil {
		query += fmt.Sprintf(" service_name = $%d", paramCount)
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
		logrus.Fatalf("Error updating subscription: %v", err)
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	logrus.Info("Subscription updated successfully: %s", id)
	return nil
}

func (r *subscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM subscriptions WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		logrus.Fatalf("Error deleting subscription: %v", err)
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	logrus.Info("Subscription deleted successfully: %s", id)
	return nil
}

func (r *subscriptionRepo) List(ctx context.Context, limit, offset int) ([]*entity.Subscription, error) {
	query := `
        SELECT *
        FROM subscriptions 
        LIMIT $1 OFFSET $2
    `

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		logrus.Fatalf("Error listing subscriptions: %v", err)
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
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, &subscription)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating subscriptions: %w", err)
	}

	logrus.Info("Listed %d subscriptions", len(subscriptions))
	return subscriptions, nil
}

func (r *subscriptionRepo) GetSummary(ctx context.Context, req *entity.SubscriptionSummaryRequest) (*entity.SubscriptionSummary, error) {
	query := `SELECT COALESCE(SUM(price), 0), COUNT(*) FROM subscriptions WHERE 1=1`
	params := []interface{}{}
	paramCount := 1

	// Если периоды не переданы, используем все время
	if req.StartPeriod != nil && req.EndPeriod != nil {
		startPeriod, err := time.Parse("01-2006", *req.StartPeriod)
		if err != nil {
			return nil, fmt.Errorf("invalid start_period format: %w", err)
		}

		endPeriod, err := time.Parse("01-2006", *req.EndPeriod)
		if err != nil {
			return nil, fmt.Errorf("invalid end_period format: %w", err)
		}

		query += fmt.Sprintf(" AND start_date <= $%d AND (end_date IS NULL OR end_date >= $%d)", paramCount, paramCount+1)
		params = append(params, endPeriod, startPeriod)
		paramCount += 2
	}

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
	err := r.db.QueryRowContext(ctx, query, params...).Scan(&summary.TotalCost, &summary.Count)
	if err != nil {
		logrus.WithError(err).Error("failed to get subscription summary")
		return nil, fmt.Errorf("failed to get subscription summary: %w", err)
	}

	return &summary, nil
}
