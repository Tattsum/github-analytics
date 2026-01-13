package application

import (
	"sort"
	"time"

	"github.com/Tattsum/github-analytics/domain"
	"github.com/Tattsum/github-analytics/infrastructure"
)

// StatisticsService は統計情報を計算するサービスです.
type StatisticsService struct {
}

// NewStatisticsService は新しいStatisticsServiceを作成します.
func NewStatisticsService() *StatisticsService {
	return &StatisticsService{}
}

// CalculateStatistics は活動データから統計情報を計算します.
func (s *StatisticsService) CalculateStatistics(data *infrastructure.UserActivityData) (*domain.UserStatistics, error) {
	stats := domain.NewUserStatistics(data.User)

	// 全活動を統合
	allActivities := make([]*domain.Activity, 0)
	allActivities = append(allActivities, data.Commits...)
	allActivities = append(allActivities, data.PRs...)
	allActivities = append(allActivities, data.Issues...)
	allActivities = append(allActivities, data.Reviews...)

	// 基本統計を計算
	s.calculateBasicStatistics(stats, allActivities, data)

	// 年別統計を計算
	s.calculateYearlyStatistics(stats, allActivities, data)

	// リポジトリ統計を計算
	s.calculateRepositoryStatistics(stats, allActivities)

	// 継続性・キャリア変遷を分析
	s.analyzeContinuityAndCareer(stats)

	return stats, nil
}

// calculateBasicStatistics は基本統計を計算します.
func (s *StatisticsService) calculateBasicStatistics(
	stats *domain.UserStatistics,
	allActivities []*domain.Activity,
	data *infrastructure.UserActivityData,
) {
	stats.TotalCommits = len(data.Commits)
	stats.TotalPRCreated = len(data.PRs)
	stats.TotalIssues = len(data.Issues)
	stats.TotalReviews = len(data.Reviews)

	// PRマージ数を計算
	for _, pr := range data.PRs {
		if pr.IsMerged {
			stats.TotalPRMerged++
		}

		stats.TotalAdditions += pr.Additions
		stats.TotalDeletions += pr.Deletions
	}

	// コミットの変更行数を計算
	for _, commit := range data.Commits {
		stats.TotalAdditions += commit.Additions
		stats.TotalDeletions += commit.Deletions
	}

	// 最初の活動年を取得
	if len(allActivities) > 0 {
		firstActivity := allActivities[0]
		for _, activity := range allActivities {
			if activity.Date.Before(firstActivity.Date) {
				firstActivity = activity
			}
		}

		stats.FirstActivityYear = firstActivity.Date.Year()
	}

	stats.CalculatePRToReviewRatio()
}

// calculateYearlyStatistics は年別統計を計算します.
func (s *StatisticsService) calculateYearlyStatistics(
	stats *domain.UserStatistics,
	allActivities []*domain.Activity,
	data *infrastructure.UserActivityData,
) {
	// 年別にコミットを集計
	yearlyCommits := make(map[int]int)

	for _, commit := range data.Commits {
		year := commit.Date.Year()
		yearlyCommits[year]++
	}

	// 年別にPRを集計
	yearlyPRCreated := make(map[int]int)
	yearlyPRMerged := make(map[int]int)

	for _, pr := range data.PRs {
		year := pr.Date.Year()

		yearlyPRCreated[year]++
		if pr.IsMerged {
			yearlyPRMerged[year]++
		}
	}

	// 年別にIssueを集計
	yearlyIssues := make(map[int]int)

	for _, issue := range data.Issues {
		year := issue.Date.Year()
		yearlyIssues[year]++
	}

	// 年別にReviewを集計
	yearlyReviews := make(map[int]int)

	for _, review := range data.Reviews {
		year := review.Date.Year()
		yearlyReviews[year]++
	}

	// 年別に変更行数を集計
	yearlyAdditions := make(map[int]int)
	yearlyDeletions := make(map[int]int)

	for _, activity := range allActivities {
		year := activity.Date.Year()
		yearlyAdditions[year] += activity.Additions
		yearlyDeletions[year] += activity.Deletions
	}

	// 全ての年を取得
	years := make(map[int]bool)
	for year := range yearlyCommits {
		years[year] = true
	}

	for year := range yearlyPRCreated {
		years[year] = true
	}

	for year := range yearlyIssues {
		years[year] = true
	}

	for year := range yearlyReviews {
		years[year] = true
	}

	// 年別統計を作成
	for year := range years {
		yearlyStat := domain.NewYearlyStatistics(year)
		yearlyStat.CommitCount = yearlyCommits[year]
		yearlyStat.PRCreated = yearlyPRCreated[year]
		yearlyStat.PRMerged = yearlyPRMerged[year]
		yearlyStat.IssueCount = yearlyIssues[year]
		yearlyStat.ReviewCount = yearlyReviews[year]
		yearlyStat.TotalAdditions = yearlyAdditions[year]
		yearlyStat.TotalDeletions = yearlyDeletions[year]
		stats.YearlyStats[year] = yearlyStat
	}

	// ピーク年を計算
	maxCommits := 0
	for year, yearlyStat := range stats.YearlyStats {
		if yearlyStat.CommitCount > maxCommits {
			maxCommits = yearlyStat.CommitCount
			stats.PeakActivityYear = year
			stats.PeakActivityCommits = maxCommits
		}
	}
}

// calculateRepositoryStatistics はリポジトリ統計を計算します.
func (s *StatisticsService) calculateRepositoryStatistics(
	stats *domain.UserStatistics,
	allActivities []*domain.Activity,
) {
	repoMap := make(map[string]*domain.RepositoryActivity)

	// リポジトリごとに活動を集計
	for _, activity := range allActivities {
		repo, exists := repoMap[activity.Repository]
		if !exists {
			repo = domain.NewRepositoryActivity(activity.Repository)
			repoMap[activity.Repository] = repo
			repo.FirstActivity = activity.Date
			repo.LastActivity = activity.Date
		}

		switch activity.Type {
		case domain.ActivityTypeCommit:
			repo.CommitCount++
		case domain.ActivityTypePR:
			repo.PRCount++
		case domain.ActivityTypeIssue:
			repo.IssueCount++
		case domain.ActivityTypeReview:
			repo.ReviewCount++
		}

		repo.TotalAdditions += activity.Additions
		repo.TotalDeletions += activity.Deletions

		if activity.Date.Before(repo.FirstActivity) {
			repo.FirstActivity = activity.Date
		}

		if activity.Date.After(repo.LastActivity) {
			repo.LastActivity = activity.Date
		}
	}

	// TOP3リポジトリを取得（コミット数順）
	repos := make([]*domain.RepositoryActivity, 0, len(repoMap))
	for _, repo := range repoMap {
		repos = append(repos, repo)
	}

	sort.Slice(repos, func(i, j int) bool {
		return repos[i].CommitCount > repos[j].CommitCount
	})

	if len(repos) > 3 {
		stats.TopRepositories = repos[:3]
	} else {
		stats.TopRepositories = repos
	}

	// 長期間関与しているリポジトリを取得（1年以上）
	for _, repo := range repos {
		duration := repo.LastActivity.Sub(repo.FirstActivity)
		if duration >= 365*24*time.Hour {
			stats.LongTermRepositories = append(stats.LongTermRepositories, repo)
		}
	}

	// 長期間関与リポジトリを期間順にソート
	sort.Slice(stats.LongTermRepositories, func(i, j int) bool {
		durationI := stats.LongTermRepositories[i].LastActivity.Sub(stats.LongTermRepositories[i].FirstActivity)
		durationJ := stats.LongTermRepositories[j].LastActivity.Sub(stats.LongTermRepositories[j].FirstActivity)

		return durationI > durationJ
	})
}

// analyzeContinuityAndCareer は継続性・キャリア変遷を分析します.
func (s *StatisticsService) analyzeContinuityAndCareer(stats *domain.UserStatistics) {
	// 年ごとのPR作成数とレビュー数の比率を計算
	years := make([]int, 0, len(stats.YearlyStats))
	for year := range stats.YearlyStats {
		years = append(years, year)
	}

	sort.Ints(years)

	for _, year := range years {
		yearlyStat := stats.YearlyStats[year]

		var ratio float64
		if yearlyStat.PRCreated > 0 {
			ratio = float64(yearlyStat.ReviewCount) / float64(yearlyStat.PRCreated)
		}

		description := s.generateRoleDescription(yearlyStat.PRCreated, yearlyStat.ReviewCount, ratio)

		stats.RoleTransition = append(stats.RoleTransition, domain.RoleTransitionPoint{
			Year:        year,
			PRCreated:   yearlyStat.PRCreated,
			ReviewCount: yearlyStat.ReviewCount,
			Ratio:       ratio,
			Description: description,
		})
	}
}

// generateRoleDescription はロール変化の説明を生成します.
func (s *StatisticsService) generateRoleDescription(prCreated, reviewCount int, ratio float64) string {
	if prCreated == 0 && reviewCount == 0 {
		return "活動なし"
	}

	if prCreated == 0 && reviewCount > 0 {
		return "レビュー中心の活動"
	}

	if prCreated > 0 && reviewCount == 0 {
		return "開発中心の活動"
	}

	if ratio < 0.5 {
		return "開発中心、レビューも実施"
	} else if ratio < 1.0 {
		return "開発とレビューのバランス"
	} else if ratio < 2.0 {
		return "レビュー重視の活動"
	} else {
		return "レビュー中心、チーム品質向上に貢献"
	}
}
