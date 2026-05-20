package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gokost710/subscription-service/internal/domain"
	"github.com/gokost710/subscription-service/internal/repository"
	"github.com/google/uuid"
)

func TestSubscriptionServiceCreate(t *testing.T) {
	startDate := mustYearMonth(t, 2025, time.July)
	repo := &subscriptionRepositoryMock{}
	svc := NewSubscriptionService(repo)

	got, err := svc.Create(context.Background(), domain.Subscription{
		ServiceName: " Yandex Plus ",
		Price:       400,
		UserID:      uuid.New(),
		StartDate:   startDate,
	})
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}

	if got.ServiceName != "Yandex Plus" {
		t.Fatalf("got service name %q", got.ServiceName)
	}

	if !repo.createCalled {
		t.Fatal("expected repository Create call")
	}
}

func TestSubscriptionServiceCreateValidation(t *testing.T) {
	startDate := mustYearMonth(t, 2025, time.July)
	userID := uuid.New()
	endDateBeforeStart := mustYearMonth(t, 2025, time.June)

	tests := []struct {
		name         string
		subscription domain.Subscription
	}{
		{
			name: "empty service name",
			subscription: domain.Subscription{
				ServiceName: " ",
				Price:       400,
				UserID:      userID,
				StartDate:   startDate,
			},
		},
		{
			name: "negative price",
			subscription: domain.Subscription{
				ServiceName: "Yandex Plus",
				Price:       -1,
				UserID:      userID,
				StartDate:   startDate,
			},
		},
		{
			name: "empty user id",
			subscription: domain.Subscription{
				ServiceName: "Yandex Plus",
				Price:       400,
				StartDate:   startDate,
			},
		},
		{
			name: "empty start date",
			subscription: domain.Subscription{
				ServiceName: "Yandex Plus",
				Price:       400,
				UserID:      userID,
			},
		},
		{
			name: "end date before start date",
			subscription: domain.Subscription{
				ServiceName: "Yandex Plus",
				Price:       400,
				UserID:      userID,
				StartDate:   startDate,
				EndDate:     &endDateBeforeStart,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &subscriptionRepositoryMock{}
			svc := NewSubscriptionService(repo)

			_, err := svc.Create(context.Background(), tt.subscription)
			if !errors.Is(err, ErrInvalidSubscription) {
				t.Fatalf("got error %v, want ErrInvalidSubscription", err)
			}

			if repo.createCalled {
				t.Fatal("repository Create must not be called")
			}
		})
	}
}

func TestSubscriptionServiceIDValidation(t *testing.T) {
	svc := NewSubscriptionService(&subscriptionRepositoryMock{})

	if _, err := svc.GetByID(context.Background(), 0); !errors.Is(err, ErrInvalidSubscription) {
		t.Fatalf("got error %v, want ErrInvalidSubscription", err)
	}

	if err := svc.Delete(context.Background(), -1); !errors.Is(err, ErrInvalidSubscription) {
		t.Fatalf("got error %v, want ErrInvalidSubscription", err)
	}
}

func TestSubscriptionServiceTotalPriceValidation(t *testing.T) {
	svc := NewSubscriptionService(&subscriptionRepositoryMock{})

	_, err := svc.TotalPrice(context.Background(), repository.SubscriptionSummaryFilter{
		From: mustYearMonth(t, 2025, time.March),
		To:   mustYearMonth(t, 2025, time.January),
	})
	if !errors.Is(err, ErrInvalidSubscription) {
		t.Fatalf("got error %v, want ErrInvalidSubscription", err)
	}
}

func mustYearMonth(t *testing.T, year int, month time.Month) domain.YearMonth {
	t.Helper()

	value, err := domain.NewYearMonth(year, month)
	if err != nil {
		t.Fatalf("create year month: %v", err)
	}

	return value
}

type subscriptionRepositoryMock struct {
	createCalled bool
}

func (m *subscriptionRepositoryMock) Create(_ context.Context, subscription domain.Subscription) (domain.Subscription, error) {
	m.createCalled = true
	subscription.ID = 1

	return subscription, nil
}

func (m *subscriptionRepositoryMock) GetByID(_ context.Context, id int64) (domain.Subscription, error) {
	return domain.Subscription{ID: id}, nil
}

func (m *subscriptionRepositoryMock) List(_ context.Context, _ repository.SubscriptionFilter) ([]domain.Subscription, error) {
	return nil, nil
}

func (m *subscriptionRepositoryMock) Update(_ context.Context, subscription domain.Subscription) (domain.Subscription, error) {
	return subscription, nil
}

func (m *subscriptionRepositoryMock) Delete(_ context.Context, _ int64) error {
	return nil
}

func (m *subscriptionRepositoryMock) TotalPrice(_ context.Context, _ repository.SubscriptionSummaryFilter) (int, error) {
	return 0, nil
}
