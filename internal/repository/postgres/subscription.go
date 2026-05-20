package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gokost710/subscription-service/internal/domain"
	"github.com/gokost710/subscription-service/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultLimit = 50
	maxLimit     = 100
)

type SubscriptionRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, subscription domain.Subscription) (domain.Subscription, error) {
	const query = `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var endDate any
	if subscription.EndDate != nil {
		endDate = subscription.EndDate.Time()
	}

	err := r.db.QueryRow(
		ctx,
		query,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate.Time(),
		endDate,
	).Scan(&subscription.ID)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("create subscription: %w", err)
	}

	return subscription, nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id int64) (domain.Subscription, error) {
	const query = `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE id = $1
	`

	subscription, err := scanSubscription(r.db.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subscription{}, repository.ErrNotFound
		}

		return domain.Subscription{}, fmt.Errorf("get subscription by id: %w", err)
	}

	return subscription, nil
}

func (r *SubscriptionRepository) List(ctx context.Context, filter repository.SubscriptionFilter) ([]domain.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
	`

	args := make([]any, 0, 4)
	where := make([]string, 0, 2)

	if filter.UserID != nil {
		args = append(args, *filter.UserID)
		where = append(where, fmt.Sprintf("user_id = $%d", len(args)))
	}

	if filter.ServiceName != nil {
		args = append(args, *filter.ServiceName)
		where = append(where, fmt.Sprintf("service_name = $%d", len(args)))
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	limit := normalizeLimit(filter.Limit)
	offset := normalizeOffset(filter.Offset)

	args = append(args, limit)
	query += fmt.Sprintf(" ORDER BY id DESC LIMIT $%d", len(args))

	args = append(args, offset)
	query += fmt.Sprintf(" OFFSET $%d", len(args))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}
	defer rows.Close()

	subscriptions := make([]domain.Subscription, 0)
	for rows.Next() {
		subscription, err := scanSubscription(rows)
		if err != nil {
			return nil, fmt.Errorf("scan subscription: %w", err)
		}

		subscriptions = append(subscriptions, subscription)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate subscriptions: %w", err)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, subscription domain.Subscription) (domain.Subscription, error) {
	const query = `
		UPDATE subscriptions
		SET service_name = $2,
			price = $3,
			user_id = $4,
			start_date = $5,
			end_date = $6,
			updated_at = now()
		WHERE id = $1
		RETURNING id, service_name, price, user_id, start_date, end_date
	`

	var endDate any
	if subscription.EndDate != nil {
		endDate = subscription.EndDate.Time()
	}

	updated, err := scanSubscription(r.db.QueryRow(
		ctx,
		query,
		subscription.ID,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate.Time(),
		endDate,
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subscription{}, repository.ErrNotFound
		}

		return domain.Subscription{}, fmt.Errorf("update subscription: %w", err)
	}

	return updated, nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM subscriptions WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}

	if result.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}

type subscriptionScanner interface {
	Scan(dest ...any) error
}

func scanSubscription(scanner subscriptionScanner) (domain.Subscription, error) {
	var (
		subscription domain.Subscription
		userID       string
		startDate    pgtype.Date
		endDate      pgtype.Date
	)

	err := scanner.Scan(
		&subscription.ID,
		&subscription.ServiceName,
		&subscription.Price,
		&userID,
		&startDate,
		&endDate,
	)
	if err != nil {
		return domain.Subscription{}, err
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("parse user id: %w", err)
	}
	subscription.UserID = parsedUserID

	subscription.StartDate, err = yearMonthFromDate(startDate)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("parse start date: %w", err)
	}

	if endDate.Valid {
		parsedEndDate, err := yearMonthFromDate(endDate)
		if err != nil {
			return domain.Subscription{}, fmt.Errorf("parse end date: %w", err)
		}

		subscription.EndDate = &parsedEndDate
	}

	return subscription, nil
}

func yearMonthFromDate(date pgtype.Date) (domain.YearMonth, error) {
	if !date.Valid {
		return domain.YearMonth{}, fmt.Errorf("date is null")
	}

	return domain.NewYearMonth(date.Time.Year(), date.Time.Month())
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return defaultLimit
	}

	if limit > maxLimit {
		return maxLimit
	}

	return limit
}

func normalizeOffset(offset int) int {
	if offset < 0 {
		return 0
	}

	return offset
}

var _ repository.SubscriptionRepository = (*SubscriptionRepository)(nil)
