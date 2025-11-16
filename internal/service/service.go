package service

import "github.com/ShekleinAleksey/subscriptions/internal/repository"

type Service struct {
	SubscriptionService SubscriptionService
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		SubscriptionService: NewSubscriptionService(r.SubscriptionRepository),
	}
}
