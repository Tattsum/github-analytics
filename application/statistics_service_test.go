package application

import (
	"testing"
	"time"

	"github.com/Tattsum/github-analytics/domain"
	"github.com/Tattsum/github-analytics/infrastructure"
)

func TestNewStatisticsService(t *testing.T) {
	t.Parallel()

	service := NewStatisticsService()
	if service == nil {
		t.Error("NewStatisticsService() should not return nil")
	}
}

func TestStatisticsService_CalculateStatistics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data *infrastructure.UserActivityData
		want func(*testing.T, *domain.UserStatistics)
	}{
		{
			name: "空のデータでも統計を計算できる",
			data: &infrastructure.UserActivityData{
				User:    domain.NewUser("testuser", "Test User", "2020-01-01T00:00:00Z"),
				Commits: []*domain.Activity{},
				PRs:     []*domain.Activity{},
				Issues:  []*domain.Activity{},
				Reviews: []*domain.Activity{},
			},
			want: func(t *testing.T, stats *domain.UserStatistics) {
				if stats.TotalCommits != 0 {
					t.Errorf("TotalCommits = %v, want 0", stats.TotalCommits)
				}

				if stats.TotalPRCreated != 0 {
					t.Errorf("TotalPRCreated = %v, want 0", stats.TotalPRCreated)
				}
			},
		},
		{
			name: "コミットデータがある場合",
			data: &infrastructure.UserActivityData{
				User: domain.NewUser("testuser", "Test User", "2020-01-01T00:00:00Z"),
				Commits: []*domain.Activity{
					domain.NewActivity(domain.ActivityTypeCommit, "owner/repo", time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), 100, 50),
					domain.NewActivity(domain.ActivityTypeCommit, "owner/repo", time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC), 200, 100),
				},
				PRs:     []*domain.Activity{},
				Issues:  []*domain.Activity{},
				Reviews: []*domain.Activity{},
			},
			want: func(t *testing.T, stats *domain.UserStatistics) {
				if stats.TotalCommits != 2 {
					t.Errorf("TotalCommits = %v, want 2", stats.TotalCommits)
				}

				if stats.FirstActivityYear != 2020 {
					t.Errorf("FirstActivityYear = %v, want 2020", stats.FirstActivityYear)
				}
			},
		},
		{
			name: "PRデータがある場合",
			data: &infrastructure.UserActivityData{
				User:    domain.NewUser("testuser", "Test User", "2020-01-01T00:00:00Z"),
				Commits: []*domain.Activity{},
				PRs: []*domain.Activity{
					func() *domain.Activity {
						pr := domain.NewActivity(domain.ActivityTypePR, "owner/repo", time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), 100, 50)
						pr.IsMerged = true

						return pr
					}(),
					domain.NewActivity(domain.ActivityTypePR, "owner/repo", time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC), 200, 100),
				},
				Issues:  []*domain.Activity{},
				Reviews: []*domain.Activity{},
			},
			want: func(t *testing.T, stats *domain.UserStatistics) {
				if stats.TotalPRCreated != 2 {
					t.Errorf("TotalPRCreated = %v, want 2", stats.TotalPRCreated)
				}

				if stats.TotalPRMerged != 1 {
					t.Errorf("TotalPRMerged = %v, want 1", stats.TotalPRMerged)
				}
			},
		},
		{
			name: "レビューデータがある場合",
			data: &infrastructure.UserActivityData{
				User:    domain.NewUser("testuser", "Test User", "2020-01-01T00:00:00Z"),
				Commits: []*domain.Activity{},
				PRs:     []*domain.Activity{},
				Issues:  []*domain.Activity{},
				Reviews: []*domain.Activity{
					func() *domain.Activity {
						review := domain.NewActivity(domain.ActivityTypeReview, "owner/repo", time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), 0, 0)
						review.IsReview = true

						return review
					}(),
					func() *domain.Activity {
						review := domain.NewActivity(domain.ActivityTypeReview, "owner/repo", time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC), 0, 0)
						review.IsReview = true

						return review
					}(),
				},
			},
			want: func(t *testing.T, stats *domain.UserStatistics) {
				if stats.TotalReviews != 2 {
					t.Errorf("TotalReviews = %v, want 2", stats.TotalReviews)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := NewStatisticsService()

			stats, err := service.CalculateStatistics(tt.data)
			if err != nil {
				t.Fatalf("CalculateStatistics() error = %v", err)
			}

			if stats == nil {
				t.Fatal("CalculateStatistics() returned nil")
			}

			tt.want(t, stats)
		})
	}
}

func TestStatisticsService_CalculateStatistics_YearlyStats(t *testing.T) {
	t.Parallel()

	service := NewStatisticsService()
	data := &infrastructure.UserActivityData{
		User: domain.NewUser("testuser", "Test User", "2020-01-01T00:00:00Z"),
		Commits: []*domain.Activity{
			domain.NewActivity(domain.ActivityTypeCommit, "owner/repo", time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), 100, 50),
			domain.NewActivity(domain.ActivityTypeCommit, "owner/repo", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), 200, 100),
		},
		PRs:     []*domain.Activity{},
		Issues:  []*domain.Activity{},
		Reviews: []*domain.Activity{},
	}

	stats, err := service.CalculateStatistics(data)
	if err != nil {
		t.Fatalf("CalculateStatistics() error = %v", err)
	}

	if len(stats.YearlyStats) != 2 {
		t.Errorf("YearlyStats length = %v, want 2", len(stats.YearlyStats))
	}

	if stats.YearlyStats[2020].CommitCount != 1 {
		t.Errorf("YearlyStats[2020].CommitCount = %v, want 1", stats.YearlyStats[2020].CommitCount)
	}

	if stats.YearlyStats[2021].CommitCount != 1 {
		t.Errorf("YearlyStats[2021].CommitCount = %v, want 1", stats.YearlyStats[2021].CommitCount)
	}
}

func TestStatisticsService_CalculateStatistics_TopRepositories(t *testing.T) {
	t.Parallel()

	service := NewStatisticsService()
	data := &infrastructure.UserActivityData{
		User: domain.NewUser("testuser", "Test User", "2020-01-01T00:00:00Z"),
		Commits: []*domain.Activity{
			domain.NewActivity(domain.ActivityTypeCommit, "owner/repo1", time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), 100, 50),
			domain.NewActivity(domain.ActivityTypeCommit, "owner/repo1", time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC), 200, 100),
			domain.NewActivity(domain.ActivityTypeCommit, "owner/repo2", time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC), 300, 150),
			domain.NewActivity(domain.ActivityTypeCommit, "owner/repo3", time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC), 400, 200),
		},
		PRs:     []*domain.Activity{},
		Issues:  []*domain.Activity{},
		Reviews: []*domain.Activity{},
	}

	stats, err := service.CalculateStatistics(data)
	if err != nil {
		t.Fatalf("CalculateStatistics() error = %v", err)
	}

	if len(stats.TopRepositories) > 3 {
		t.Errorf("TopRepositories length = %v, want <= 3", len(stats.TopRepositories))
	}

	if len(stats.TopRepositories) > 0 {
		// コミット数が多い順にソートされていることを確認
		prevCount := stats.TopRepositories[0].CommitCount
		for i := 1; i < len(stats.TopRepositories); i++ {
			if stats.TopRepositories[i].CommitCount > prevCount {
				t.Error("TopRepositories should be sorted by commit count (descending)")
			}

			prevCount = stats.TopRepositories[i].CommitCount
		}
	}
}
