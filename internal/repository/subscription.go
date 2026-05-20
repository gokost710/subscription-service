package repository

import (
	"context"

	"github.com/gokost710/subscription-service/internal/domain"
	"github.com/google/uuid"
)

type SubscriptionFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	Limit       int
	Offset      int
}

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription domain.Subscription) (domain.Subscription, error)
	GetByID(ctx context.Context, id int64) (domain.Subscription, error)
	List(ctx context.Context, filter SubscriptionFilter) ([]domain.Subscription, error)
	Update(ctx context.Context, subscription domain.Subscription) (domain.Subscription, error)
	Delete(ctx context.Context, id int64) error
}
