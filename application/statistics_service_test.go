package application

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Tattsum/github-analytics/domain"
	"github.com/Tattsum/github-analytics/infrastructure"
)

func TestNewStatisticsService(t *testing.T) {
	t.Parallel()

	service := NewStatisticsService()
	assert.NotNil(t, service, "NewStatisticsService() should not return nil")
}

// assertEmptyStats は空の統計のアサーションを行います.
func assertEmptyStats(t *testing.T, stats *domain.UserStatistics) {
	t.Helper()

	assert.Equal(t, 0, stats.TotalCommits, "TotalCommits should be 0")
	assert.Equal(t, 0, stats.TotalPRCreated, "TotalPRCreated should be 0")
}

// assertCommitStats はコミット統計のアサーションを行います.
func assertCommitStats(t *testing.T, stats *domain.UserStatistics, wantCommits, wantYear int) {
	t.Helper()

	assert.Equal(t, wantCommits, stats.TotalCommits, "TotalCommits should match")
	assert.Equal(t, wantYear, stats.FirstActivityYear, "FirstActivityYear should match")
}

// assertPRStats はPR統計のアサーションを行います.
func assertPRStats(t *testing.T, stats *domain.UserStatistics, wantCreated, wantMerged int) {
	t.Helper()

	assert.Equal(t, wantCreated, stats.TotalPRCreated, "TotalPRCreated should match")
	assert.Equal(t, wantMerged, stats.TotalPRMerged, "TotalPRMerged should match")
}

// assertReviewStats はレビュー統計のアサーションを行います.
func assertReviewStats(t *testing.T, stats *domain.UserStatistics, wantReviews int) {
	t.Helper()

	assert.Equal(t, wantReviews, stats.TotalReviews, "TotalReviews should match")
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
			require.NoError(t, err, "CalculateStatistics() should not return error")
			require.NotNil(t, stats, "CalculateStatistics() should not return nil")

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
	require.NoError(t, err, "CalculateStatistics() should not return error")

	assert.Equal(t, 2, len(stats.YearlyStats), "YearlyStats should have 2 entries")
	assert.Equal(t, 1, stats.YearlyStats[2020].CommitCount, "YearlyStats[2020].CommitCount should be 1")
	assert.Equal(t, 1, stats.YearlyStats[2021].CommitCount, "YearlyStats[2021].CommitCount should be 1")
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
	require.NoError(t, err, "CalculateStatistics() should not return error")

	assert.LessOrEqual(t, len(stats.TopRepositories), 3, "TopRepositories length should be <= 3")

	if len(stats.TopRepositories) > 0 {
		// コミット数が多い順にソートされていることを確認
		prevCount := stats.TopRepositories[0].CommitCount
		for i := 1; i < len(stats.TopRepositories); i++ {
			assert.LessOrEqual(t, stats.TopRepositories[i].CommitCount, prevCount,
				"TopRepositories should be sorted by commit count (descending)")
			prevCount = stats.TopRepositories[i].CommitCount
		}
	}
}
