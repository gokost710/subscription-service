package domain

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const yearMonthLayout = "01-2006"

type YearMonth struct {
	year  int
	month time.Month
}

func NewYearMonth(year int, month time.Month) (YearMonth, error) {
	if year < 1 {
		return YearMonth{}, fmt.Errorf("invalid year: %d", year)
	}

	if month < time.January || month > time.December {
		return YearMonth{}, fmt.Errorf("invalid month: %d", month)
	}

	return YearMonth{
		year:  year,
		month: month,
	}, nil
}

func ParseYearMonth(value string) (YearMonth, error) {
	if len(value) != len(yearMonthLayout) {
		return YearMonth{}, fmt.Errorf("invalid year-month format: %q", value)
	}

	if value[2] != '-' {
		return YearMonth{}, fmt.Errorf("invalid year-month format: %q", value)
	}

	month, err := strconv.Atoi(value[:2])
	if err != nil {
		return YearMonth{}, fmt.Errorf("parse month: %w", err)
	}

	year, err := strconv.Atoi(value[3:])
	if err != nil {
		return YearMonth{}, fmt.Errorf("parse year: %w", err)
	}

	return NewYearMonth(year, time.Month(month))
}

func (ym YearMonth) Year() int {
	return ym.year
}

func (ym YearMonth) Month() time.Month {
	return ym.month
}

func (ym YearMonth) IsZero() bool {
	return ym.year == 0 && ym.month == 0
}

func (ym YearMonth) Time() time.Time {
	return time.Date(ym.year, ym.month, 1, 0, 0, 0, 0, time.UTC)
}

func (ym YearMonth) String() string {
	return ym.Time().Format(yearMonthLayout)
}

func (ym YearMonth) MarshalJSON() ([]byte, error) {
	return json.Marshal(ym.String())
}

func (ym *YearMonth) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	parsed, err := ParseYearMonth(value)
	if err != nil {
		return err
	}

	*ym = parsed

	return nil
}
