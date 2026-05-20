package domain

import "github.com/google/uuid"

type Subscription struct {
	ID          int64
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   YearMonth
	EndDate     *YearMonth
}
