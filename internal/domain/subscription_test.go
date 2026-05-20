package domain

import (
	"testing"
	"time"
)

func TestSubscriptionActiveMonthsInPeriod(t *testing.T) {
	tests := []struct {
		name         string
		subscription Subscription
		from         YearMonth
		to           YearMonth
		want         int
	}{
		{
			name: "subscription fully inside period",
			subscription: Subscription{
				StartDate: mustYearMonth(t, 2025, time.February),
				EndDate:   yearMonthPtr(t, 2025, time.April),
			},
			from: mustYearMonth(t, 2025, time.January),
			to:   mustYearMonth(t, 2025, time.December),
			want: 3,
		},
		{
			name: "subscription starts before period and has no end",
			subscription: Subscription{
				StartDate: mustYearMonth(t, 2024, time.December),
			},
			from: mustYearMonth(t, 2025, time.January),
			to:   mustYearMonth(t, 2025, time.March),
			want: 3,
		},
		{
			name: "subscription overlaps period end",
			subscription: Subscription{
				StartDate: mustYearMonth(t, 2025, time.March),
				EndDate:   yearMonthPtr(t, 2025, time.July),
			},
			from: mustYearMonth(t, 2025, time.January),
			to:   mustYearMonth(t, 2025, time.May),
			want: 3,
		},
		{
			name: "subscription before period",
			subscription: Subscription{
				StartDate: mustYearMonth(t, 2024, time.January),
				EndDate:   yearMonthPtr(t, 2024, time.December),
			},
			from: mustYearMonth(t, 2025, time.January),
			to:   mustYearMonth(t, 2025, time.May),
			want: 0,
		},
		{
			name: "subscription after period",
			subscription: Subscription{
				StartDate: mustYearMonth(t, 2025, time.June),
			},
			from: mustYearMonth(t, 2025, time.January),
			to:   mustYearMonth(t, 2025, time.May),
			want: 0,
		},
		{
			name: "same month overlap",
			subscription: Subscription{
				StartDate: mustYearMonth(t, 2025, time.May),
				EndDate:   yearMonthPtr(t, 2025, time.May),
			},
			from: mustYearMonth(t, 2025, time.May),
			to:   mustYearMonth(t, 2025, time.May),
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.subscription.ActiveMonthsInPeriod(tt.from, tt.to)
			if got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestMonthsBetweenInclusive(t *testing.T) {
	tests := []struct {
		name string
		from YearMonth
		to   YearMonth
		want int
	}{
		{
			name: "same month",
			from: mustYearMonth(t, 2025, time.January),
			to:   mustYearMonth(t, 2025, time.January),
			want: 1,
		},
		{
			name: "same year",
			from: mustYearMonth(t, 2025, time.January),
			to:   mustYearMonth(t, 2025, time.March),
			want: 3,
		},
		{
			name: "different years",
			from: mustYearMonth(t, 2024, time.November),
			to:   mustYearMonth(t, 2025, time.February),
			want: 4,
		},
		{
			name: "invalid range",
			from: mustYearMonth(t, 2025, time.March),
			to:   mustYearMonth(t, 2025, time.January),
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MonthsBetweenInclusive(tt.from, tt.to)
			if got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func mustYearMonth(t *testing.T, year int, month time.Month) YearMonth {
	t.Helper()

	value, err := NewYearMonth(year, month)
	if err != nil {
		t.Fatalf("create year month: %v", err)
	}

	return value
}

func yearMonthPtr(t *testing.T, year int, month time.Month) *YearMonth {
	t.Helper()

	value := mustYearMonth(t, year, month)
	return &value
}
