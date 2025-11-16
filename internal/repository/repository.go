package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	SubscriptionRepository SubscriptionRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		SubscriptionRepository: NewSubscriptionRepository(db),
	}
}
