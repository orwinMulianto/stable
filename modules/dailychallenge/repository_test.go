package dailychallenge

import (
	"testing"
	"time"
)

func TestCalculateNextStreak(t *testing.T) {
	today := time.Date(2026, time.June, 22, 0, 0, 0, 0, jakartaLocation)
	yesterday := today.AddDate(0, 0, -1)
	twoDaysAgo := today.AddDate(0, 0, -2)

	tests := []struct {
		name     string
		current  int
		lastDate *time.Time
		expected int
	}{
		{name: "first completion", current: 0, lastDate: nil, expected: 1},
		{name: "consecutive day", current: 4, lastDate: &yesterday, expected: 5},
		{name: "missed day", current: 4, lastDate: &twoDaysAgo, expected: 1},
		{name: "same day does not increase", current: 4, lastDate: &today, expected: 4},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := calculateNextStreak(test.current, test.lastDate, today)
			if result != test.expected {
				t.Fatalf("expected %d, got %d", test.expected, result)
			}
		})
	}
}
