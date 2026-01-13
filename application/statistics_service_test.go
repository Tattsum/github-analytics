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

// assertEmptyStats は空の統計のアサーションを行います.
func assertEmptyStats(t *testing.T, stats *domain.UserStatistics) {
	t.Helper()

	if stats.TotalCommits != 0 {
		t.Errorf("TotalCommits = %v, want 0", stats.TotalCommits)
	}

	if stats.TotalPRCreated != 0 {
		t.Errorf("TotalPRCreated = %v, want 0", stats.TotalPRCreated)
	}
}

// assertCommitStats はコミット統計のアサーションを行います.
func assertCommitStats(t *testing.T, stats *domain.UserStatistics, wantCommits, wantYear int) {
	t.Helper()

	if stats.TotalCommits != wantCommits {
		t.Errorf("TotalCommits = %v, want %v", stats.TotalCommits, wantCommits)
	}

	if stats.FirstActivityYear != wantYear {
		t.Errorf("FirstActivityYear = %v, want %v", stats.FirstActivityYear, wantYear)
	}
}

// assertPRStats はPR統計のアサーションを行います.
func assertPRStats(t *testing.T, stats *domain.UserStatistics, wantCreated, wantMerged int) {
	t.Helper()

	if stats.TotalPRCreated != wantCreated {
		t.Errorf("TotalPRCreated = %v, want %v", stats.TotalPRCreated, wantCreated)
	}

	if stats.TotalPRMerged != wantMerged {
		t.Errorf("TotalPRMerged = %v, want %v", stats.TotalPRMerged, wantMerged)
	}
}

// assertReviewStats はレビュー統計のアサーションを行います.
func assertReviewStats(t *testing.T, stats *domain.UserStatistics, wantReviews int) {
	t.Helper()

	if stats.TotalReviews != wantReviews {
		t.Errorf("TotalReviews = %v, want %v", stats.TotalReviews, wantReviews)
	}
}

// createPRsWithMerged はマージされたPRを含むPRリストを作成します.
func createPRsWithMerged() []*domain.Activity {
	pr1 := domain.NewActivity(domain.ActivityTypePR, "owner/repo", time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), 100, 50)
	pr1.IsMerged = true

	return []*domain.Activity{
		pr1,
		domain.NewActivity(domain.ActivityTypePR, "owner/repo", time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC), 200, 100),
	}
}

// createReviews はレビューリストを作成します.
func createReviews() []*domain.Activity {
	review1 := domain.NewActivity(domain.ActivityTypeReview, "owner/repo", time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), 0, 0)
	review1.IsReview = true
	review2 := domain.NewActivity(domain.ActivityTypeReview, "owner/repo", time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC), 0, 0)
	review2.IsReview = true

	return []*domain.Activity{review1, review2}
}

func TestStatisticsService_CalculateStatistics(t *testing.T) {
	t.Parallel()

	tests := createTestCases()

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

type testCase struct {
	name string
	data *infrastructure.UserActivityData
	want func(*testing.T, *domain.UserStatistics)
}

// createTestCases はテストケースを作成します.
func createTestCases() []testCase {
	return []testCase{
		createEmptyDataTestCase(),
		createCommitDataTestCase(),
		createPRDataTestCase(),
		createReviewDataTestCase(),
	}
}

// createEmptyDataTestCase は空のデータのテストケースを作成します.
func createEmptyDataTestCase() testCase {
	return testCase{
		name: "空のデータでも統計を計算できる",
		data: &infrastructure.UserActivityData{
			User:    domain.NewUser("testuser", "Test User", "2020-01-01T00:00:00Z"),
			Commits: []*domain.Activity{},
			PRs:     []*domain.Activity{},
			Issues:  []*domain.Activity{},
			Reviews: []*domain.Activity{},
		},
		want: assertEmptyStats,
	}
}

// createCommitDataTestCase はコミットデータのテストケースを作成します.
func createCommitDataTestCase() testCase {
	return testCase{
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
			assertCommitStats(t, stats, 2, 2020)
		},
	}
}

// createPRDataTestCase はPRデータのテストケースを作成します.
func createPRDataTestCase() testCase {
	return testCase{
		name: "PRデータがある場合",
		data: &infrastructure.UserActivityData{
			User:    domain.NewUser("testuser", "Test User", "2020-01-01T00:00:00Z"),
			Commits: []*domain.Activity{},
			PRs:     createPRsWithMerged(),
			Issues:  []*domain.Activity{},
			Reviews: []*domain.Activity{},
		},
		want: func(t *testing.T, stats *domain.UserStatistics) {
			assertPRStats(t, stats, 2, 1)
		},
	}
}

// createReviewDataTestCase はレビューデータのテストケースを作成します.
func createReviewDataTestCase() testCase {
	return testCase{
		name: "レビューデータがある場合",
		data: &infrastructure.UserActivityData{
			User:    domain.NewUser("testuser", "Test User", "2020-01-01T00:00:00Z"),
			Commits: []*domain.Activity{},
			PRs:     []*domain.Activity{},
			Issues:  []*domain.Activity{},
			Reviews: createReviews(),
		},
		want: func(t *testing.T, stats *domain.UserStatistics) {
			assertReviewStats(t, stats, 2)
		},
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
