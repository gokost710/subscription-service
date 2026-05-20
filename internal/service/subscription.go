package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/gokost710/subscription-service/internal/domain"
	"github.com/gokost710/subscription-service/internal/repository"
	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) Create(ctx context.Context, subscription domain.Subscription) (domain.Subscription, error) {
	subscription = normalizeSubscription(subscription)
	if err := validateSubscription(subscription); err != nil {
		return domain.Subscription{}, err
	}

	return s.repo.Create(ctx, subscription)
}

func (s *SubscriptionService) GetByID(ctx context.Context, id int64) (domain.Subscription, error) {
	if id <= 0 {
		return domain.Subscription{}, fmt.Errorf("%w: id must be positive", ErrInvalidSubscription)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *SubscriptionService) List(ctx context.Context, filter repository.SubscriptionFilter) ([]domain.Subscription, error) {
	filter = normalizeFilter(filter)

	return s.repo.List(ctx, filter)
}

func (s *SubscriptionService) Update(ctx context.Context, subscription domain.Subscription) (domain.Subscription, error) {
	if subscription.ID <= 0 {
		return domain.Subscription{}, fmt.Errorf("%w: id must be positive", ErrInvalidSubscription)
	}

	subscription = normalizeSubscription(subscription)
	if err := validateSubscription(subscription); err != nil {
		return domain.Subscription{}, err
	}

	return s.repo.Update(ctx, subscription)
}

func (s *SubscriptionService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidSubscription)
	}

	return s.repo.Delete(ctx, id)
}

func normalizeSubscription(subscription domain.Subscription) domain.Subscription {
	subscription.ServiceName = strings.TrimSpace(subscription.ServiceName)

	return subscription
}

func normalizeFilter(filter repository.SubscriptionFilter) repository.SubscriptionFilter {
	if filter.ServiceName != nil {
		serviceName := strings.TrimSpace(*filter.ServiceName)
		if serviceName == "" {
			filter.ServiceName = nil
		} else {
			filter.ServiceName = &serviceName
		}
	}

	return filter
}

func validateSubscription(subscription domain.Subscription) error {
	if subscription.ServiceName == "" {
		return fmt.Errorf("%w: service_name is required", ErrInvalidSubscription)
	}

	if subscription.Price < 0 {
		return fmt.Errorf("%w: price must be greater than or equal to zero", ErrInvalidSubscription)
	}

	if subscription.UserID == uuid.Nil {
		return fmt.Errorf("%w: user_id is required", ErrInvalidSubscription)
	}

	if subscription.StartDate.IsZero() {
		return fmt.Errorf("%w: start_date is required", ErrInvalidSubscription)
	}

	if subscription.EndDate != nil && subscription.EndDate.Time().Before(subscription.StartDate.Time()) {
		return fmt.Errorf("%w: end_date must be greater than or equal to start_date", ErrInvalidSubscription)
	}

	return nil
}
