package domain

import (
	"testing"
)

func TestNewYearlyStatistics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		year int
		want *YearlyStatistics
	}{
		{
			name: "2020年の統計作成",
			year: 2020,
			want: &YearlyStatistics{
				Year: 2020,
			},
		},
		{
			name: "2024年の統計作成",
			year: 2024,
			want: &YearlyStatistics{
				Year: 2024,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NewYearlyStatistics(tt.year)
			if got.Year != tt.want.Year {
				t.Errorf("NewYearlyStatistics().Year = %v, want %v", got.Year, tt.want.Year)
			}
		})
	}
}

func TestNewUserStatistics(t *testing.T) {
	t.Parallel()

	user := NewUser("testuser", "Test User", "2020-01-01T00:00:00Z")
	got := NewUserStatistics(user)

	if got.User != user {
		t.Errorf("NewUserStatistics().User = %v, want %v", got.User, user)
	}

	if got.TotalCommits != 0 {
		t.Errorf("NewUserStatistics().TotalCommits = %v, want 0", got.TotalCommits)
	}

	if got.YearlyStats == nil {
		t.Error("NewUserStatistics().YearlyStats should not be nil")
	}

	if got.TopRepositories == nil {
		t.Error("NewUserStatistics().TopRepositories should not be nil")
	}
}

func TestUserStatistics_CalculatePRToReviewRatio(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		totalPRCreated int
		totalReviews   int
		want           float64
	}{
		{
			name:           "PR作成数がレビュー数より多い場合",
			totalPRCreated: 10,
			totalReviews:   5,
			want:           0.5,
		},
		{
			name:           "レビュー数がPR作成数より多い場合",
			totalPRCreated: 5,
			totalReviews:   10,
			want:           2.0,
		},
		{
			name:           "PR作成数が0の場合、比率は0",
			totalPRCreated: 0,
			totalReviews:   5,
			want:           0.0,
		},
		{
			name:           "同じ数の場合",
			totalPRCreated: 10,
			totalReviews:   10,
			want:           1.0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			user := NewUser("testuser", "Test User", "2020-01-01T00:00:00Z")
			stats := NewUserStatistics(user)
			stats.TotalPRCreated = tt.totalPRCreated
			stats.TotalReviews = tt.totalReviews
			stats.CalculatePRToReviewRatio()

			if stats.PRToReviewRatio != tt.want {
				t.Errorf("UserStatistics.CalculatePRToReviewRatio() = %v, want %v", stats.PRToReviewRatio, tt.want)
			}
		})
	}
}
