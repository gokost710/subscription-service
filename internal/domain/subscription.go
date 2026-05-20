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

func (s Subscription) ActiveMonthsInPeriod(from YearMonth, to YearMonth) int {
	if from.IsZero() || to.IsZero() || to.Before(from) {
		return 0
	}

	activeFrom := maxYearMonth(s.StartDate, from)
	activeTo := to
	if s.EndDate != nil {
		activeTo = minYearMonth(*s.EndDate, to)
	}

	if activeTo.Before(activeFrom) {
		return 0
	}

	return MonthsBetweenInclusive(activeFrom, activeTo)
}
