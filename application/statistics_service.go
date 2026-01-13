// Package application はユースケース層を提供します.
// このパッケージはドメイン層の協調を行い、ビジネスロジックを実装します.
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

// aggregateYearlyData は年別データを集計します.
func (s *StatisticsService) aggregateYearlyData(
	data *infrastructure.UserActivityData,
	allActivities []*domain.Activity,
) (map[int]int, map[int]int, map[int]int, map[int]int, map[int]int, map[int]int, map[int]int) {
	yearlyCommits := make(map[int]int)
	yearlyPRCreated := make(map[int]int)
	yearlyPRMerged := make(map[int]int)
	yearlyIssues := make(map[int]int)
	yearlyReviews := make(map[int]int)
	yearlyAdditions := make(map[int]int)
	yearlyDeletions := make(map[int]int)

	for _, commit := range data.Commits {
		year := commit.Date.Year()
		yearlyCommits[year]++
	}

	for _, pr := range data.PRs {
		year := pr.Date.Year()
		yearlyPRCreated[year]++

		if pr.IsMerged {
			yearlyPRMerged[year]++
		}
	}

	for _, issue := range data.Issues {
		year := issue.Date.Year()
		yearlyIssues[year]++
	}

	for _, review := range data.Reviews {
		year := review.Date.Year()
		yearlyReviews[year]++
	}

	for _, activity := range allActivities {
		year := activity.Date.Year()
		yearlyAdditions[year] += activity.Additions
		yearlyDeletions[year] += activity.Deletions
	}

	return yearlyCommits, yearlyPRCreated, yearlyPRMerged, yearlyIssues, yearlyReviews, yearlyAdditions, yearlyDeletions
}

// collectAllYears は全ての年を収集します.
func (s *StatisticsService) collectAllYears(
	yearlyCommits, yearlyPRCreated, yearlyIssues, yearlyReviews map[int]int,
) map[int]bool {
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

	return years
}

// buildYearlyStats は年別統計を作成します.
func (s *StatisticsService) buildYearlyStats(
	stats *domain.UserStatistics,
	years map[int]bool,
	yearlyCommits, yearlyPRCreated, yearlyPRMerged, yearlyIssues, yearlyReviews, yearlyAdditions, yearlyDeletions map[int]int,
) {
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
}

// calculatePeakYear はピーク年を計算します.
func (s *StatisticsService) calculatePeakYear(stats *domain.UserStatistics) {
	maxCommits := 0
	for year, yearlyStat := range stats.YearlyStats {
		if yearlyStat.CommitCount > maxCommits {
			maxCommits = yearlyStat.CommitCount
			stats.PeakActivityYear = year
			stats.PeakActivityCommits = maxCommits
		}
	}
}

// calculateYearlyStatistics は年別統計を計算します.
func (s *StatisticsService) calculateYearlyStatistics(
	stats *domain.UserStatistics,
	allActivities []*domain.Activity,
	data *infrastructure.UserActivityData,
) {
	yearlyCommits, yearlyPRCreated, yearlyPRMerged, yearlyIssues, yearlyReviews, yearlyAdditions, yearlyDeletions :=
		s.aggregateYearlyData(data, allActivities)

	years := s.collectAllYears(yearlyCommits, yearlyPRCreated, yearlyIssues, yearlyReviews)

	s.buildYearlyStats(stats, years, yearlyCommits, yearlyPRCreated, yearlyPRMerged, yearlyIssues, yearlyReviews, yearlyAdditions, yearlyDeletions)

	s.calculatePeakYear(stats)
}

// aggregateRepositoryActivities はリポジトリごとに活動を集計します.
func (s *StatisticsService) aggregateRepositoryActivities(allActivities []*domain.Activity) map[string]*domain.RepositoryActivity {
	repoMap := make(map[string]*domain.RepositoryActivity)

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

	return repoMap
}

// selectTopRepositories はTOP3リポジトリを選択します.
func (s *StatisticsService) selectTopRepositories(repoMap map[string]*domain.RepositoryActivity) []*domain.RepositoryActivity {
	repos := make([]*domain.RepositoryActivity, 0, len(repoMap))
	for _, repo := range repoMap {
		repos = append(repos, repo)
	}

	sort.Slice(repos, func(i, j int) bool {
		return repos[i].CommitCount > repos[j].CommitCount
	})

	const topRepoLimit = 3
	if len(repos) > topRepoLimit {
		return repos[:topRepoLimit]
	}

	return repos
}

// findLongTermRepositories は長期間関与しているリポジトリを取得します.
func (s *StatisticsService) findLongTermRepositories(repos []*domain.RepositoryActivity) []*domain.RepositoryActivity {
	const oneYear = 365 * 24 * time.Hour

	longTermRepos := make([]*domain.RepositoryActivity, 0)

	for _, repo := range repos {
		duration := repo.LastActivity.Sub(repo.FirstActivity)
		if duration >= oneYear {
			longTermRepos = append(longTermRepos, repo)
		}
	}

	sort.Slice(longTermRepos, func(i, j int) bool {
		durationI := longTermRepos[i].LastActivity.Sub(longTermRepos[i].FirstActivity)
		durationJ := longTermRepos[j].LastActivity.Sub(longTermRepos[j].FirstActivity)

		return durationI > durationJ
	})

	return longTermRepos
}

// calculateRepositoryStatistics はリポジトリ統計を計算します.
func (s *StatisticsService) calculateRepositoryStatistics(
	stats *domain.UserStatistics,
	allActivities []*domain.Activity,
) {
	repoMap := s.aggregateRepositoryActivities(allActivities)
	stats.TopRepositories = s.selectTopRepositories(repoMap)
	stats.LongTermRepositories = s.findLongTermRepositories(stats.TopRepositories)
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

// generateRoleDescriptionForNoActivity は活動がない場合の説明を生成します.
func (s *StatisticsService) generateRoleDescriptionForNoActivity(prCreated, reviewCount int) (string, bool) {
	if prCreated == 0 && reviewCount == 0 {
		return "活動なし", true
	}

	if prCreated == 0 && reviewCount > 0 {
		return "レビュー中心の活動", true
	}

	if prCreated > 0 && reviewCount == 0 {
		return "開発中心の活動", true
	}

	return "", false
}

// generateRoleDescriptionByRatio は比率に基づいてロール説明を生成します.
func (s *StatisticsService) generateRoleDescriptionByRatio(ratio float64) string {
	const (
		ratioThresholdLow  = 0.5
		ratioThresholdMid  = 1.0
		ratioThresholdHigh = 2.0
	)

	switch {
	case ratio < ratioThresholdLow:
		return "開発中心、レビューも実施"
	case ratio < ratioThresholdMid:
		return "開発とレビューのバランス"
	case ratio < ratioThresholdHigh:
		return "レビュー重視の活動"
	default:
		return "レビュー中心、チーム品質向上に貢献"
	}
}

// generateRoleDescription はロール変化の説明を生成します.
func (s *StatisticsService) generateRoleDescription(prCreated, reviewCount int, ratio float64) string {
	if desc, ok := s.generateRoleDescriptionForNoActivity(prCreated, reviewCount); ok {
		return desc
	}

	return s.generateRoleDescriptionByRatio(ratio)
}
