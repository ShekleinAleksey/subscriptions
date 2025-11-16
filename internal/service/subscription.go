package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ShekleinAleksey/subscriptions/internal/entity"
	"github.com/ShekleinAleksey/subscriptions/internal/repository"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, req *entity.CreateSubscriptionRequest) (*entity.Subscription, error)
	GetSubscription(ctx context.Context, id uuid.UUID) (*entity.Subscription, error)
	UpdateSubscription(ctx context.Context, id uuid.UUID, req *entity.UpdateSubscriptionRequest) error
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	ListSubscriptions(ctx context.Context, limit, offset int) ([]*entity.Subscription, error)
	GetSubscriptionSummary(ctx context.Context, req *entity.SubscriptionSummaryRequest) (*entity.SubscriptionSummary, error)
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) CreateSubscription(ctx context.Context, req *entity.CreateSubscriptionRequest) (*entity.Subscription, error) {
	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format: %w", err)
	}

	var endDate *time.Time
	if req.EndDate != nil {
		parsedEndDate, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date format: %w", err)
		}
		endDate = &parsedEndDate
	}

	subscription := &entity.Subscription{
		ID:          uuid.New(),
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *subscriptionService) GetSubscription(ctx context.Context, id uuid.UUID) (*entity.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *subscriptionService) UpdateSubscription(ctx context.Context, id uuid.UUID, req *entity.UpdateSubscriptionRequest) error {
	return s.repo.Update(ctx, id, req)
}

func (s *subscriptionService) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *subscriptionService) ListSubscriptions(ctx context.Context, limit, offset int) ([]*entity.Subscription, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.List(ctx, limit, offset)
}

func (s *subscriptionService) GetSubscriptionSummary(ctx context.Context, req *entity.SubscriptionSummaryRequest) (*entity.SubscriptionSummary, error) {
	return s.repo.GetSummary(ctx, req)
}
