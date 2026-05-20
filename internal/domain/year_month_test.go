package domain

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParseYearMonth(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    YearMonth
		wantErr bool
	}{
		{
			name:  "valid",
			value: "07-2025",
			want: YearMonth{
				year:  2025,
				month: time.July,
			},
		},
		{
			name:    "single digit month",
			value:   "7-2025",
			wantErr: true,
		},
		{
			name:    "zero month",
			value:   "00-2025",
			wantErr: true,
		},
		{
			name:    "month greater than december",
			value:   "13-2025",
			wantErr: true,
		},
		{
			name:    "empty",
			value:   "",
			wantErr: true,
		},
		{
			name:    "wrong separator",
			value:   "07/2025",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseYearMonth(tt.value)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestYearMonthTime(t *testing.T) {
	ym, err := ParseYearMonth("07-2025")
	if err != nil {
		t.Fatalf("parse year month: %v", err)
	}

	want := time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC)
	if !ym.Time().Equal(want) {
		t.Fatalf("got %s, want %s", ym.Time(), want)
	}
}

func TestYearMonthString(t *testing.T) {
	ym, err := ParseYearMonth("07-2025")
	if err != nil {
		t.Fatalf("parse year month: %v", err)
	}

	if got, want := ym.String(), "07-2025"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestYearMonthJSON(t *testing.T) {
	type payload struct {
		StartDate YearMonth `json:"start_date"`
	}

	var got payload
	if err := json.Unmarshal([]byte(`{"start_date":"07-2025"}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.StartDate.Year() != 2025 || got.StartDate.Month() != time.July {
		t.Fatalf("got %+v", got.StartDate)
	}

	data, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if string(data) != `{"start_date":"07-2025"}` {
		t.Fatalf("got %s", data)
	}
}

func TestYearMonthJSONInvalid(t *testing.T) {
	type payload struct {
		StartDate YearMonth `json:"start_date"`
	}

	var got payload
	if err := json.Unmarshal([]byte(`{"start_date":"13-2025"}`), &got); err == nil {
		t.Fatal("expected error, got nil")
	}
}
