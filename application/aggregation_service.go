package application

import (
	"sort"

	"github.com/Tattsum/github-analytics/domain"
)

// MemberRepoStat はメンバー×リポジトリ1件分の集計済みメトリクスです.
// リポジトリ軸の横断集計（リポジトリごとの合計・貢献者一覧）を組み立てる入力として用います.
// infrastructure 層が ent の行をこの純粋な構造体へマッピングし、集計ロジック本体は DB に依存せずテスト可能にします.
type MemberRepoStat struct {
	Login         string
	NameWithOwner string
	CommitCount   int
	PRCreated     int
	PRMerged      int
	IssueCount    int
	ReviewCount   int
	Additions     int
	Deletions     int
}

// SummarizeTeam はメンバー横断スカラー指標を合計し、チーム全体の集計値を返します.
// RepositoryCount は別途リポジトリ軸の集計から求めるため、ここでは設定しません（呼び出し元が補完します）.
func SummarizeTeam(members []*MemberStats) *TeamSummary {
	summary := &TeamSummary{
		MemberCount: len(members),
	}

	for _, member := range members {
		summary.TotalCommits += member.TotalCommits
		summary.TotalPRCreated += member.TotalPRCreated
		summary.TotalPRMerged += member.TotalPRMerged
		summary.TotalIssues += member.TotalIssues
		summary.TotalReviews += member.TotalReviews
		summary.TotalAdditions += member.TotalAdditions
		summary.TotalDeletions += member.TotalDeletions
	}

	return summary
}

// AggregateRepositories はメンバー×リポジトリの行を nameWithOwner でグルーピングし、
// リポジトリ軸の横断集計（合計・貢献者一覧・貢献者数）へ再集計します.
// 戻り値は nameWithOwner の昇順で安定ソートされ、各リポジトリの Contributors は login の昇順です.
// ランキング・ソート・比較はフロントエンドで行うため、ここでは決定的な順序のみ保証します.
func AggregateRepositories(stats []*MemberRepoStat) []*RepositoryStats {
	repoIndex := make(map[string]*RepositoryStats)

	for _, stat := range stats {
		repo, exists := repoIndex[stat.NameWithOwner]
		if !exists {
			repo = &RepositoryStats{
				NameWithOwner: stat.NameWithOwner,
				Contributors:  make([]*RepositoryContributor, 0, 1),
			}
			repoIndex[stat.NameWithOwner] = repo
		}

		repo.TotalCommits += stat.CommitCount
		repo.TotalPRCreated += stat.PRCreated
		repo.TotalPRMerged += stat.PRMerged
		repo.TotalIssues += stat.IssueCount
		repo.TotalReviews += stat.ReviewCount
		repo.TotalAdditions += stat.Additions
		repo.TotalDeletions += stat.Deletions

		repo.Contributors = append(repo.Contributors, &RepositoryContributor{
			Login:       stat.Login,
			CommitCount: stat.CommitCount,
			PRCreated:   stat.PRCreated,
			ReviewCount: stat.ReviewCount,
			Additions:   stat.Additions,
			Deletions:   stat.Deletions,
		})
	}

	repos := make([]*RepositoryStats, 0, len(repoIndex))
	for _, repo := range repoIndex {
		repo.ContributorCount = len(repo.Contributors)

		sort.Slice(repo.Contributors, func(i, j int) bool {
			return repo.Contributors[i].Login < repo.Contributors[j].Login
		})

		repos = append(repos, repo)
	}

	sort.Slice(repos, func(i, j int) bool {
		return repos[i].NameWithOwner < repos[j].NameWithOwner
	})

	return repos
}

// AggregateTeamDaily はメンバー単位の日別統計行を同一日付で合算し、
// チーム全体の日別合計を日付(YYYY-MM-DD)昇順の時系列で返します.
// 期間絞り込み・推移グラフはフロントエンドで行うため、ここでは決定的な昇順のみ保証します.
func AggregateTeamDaily(rows []*domain.DailyStatistics) []*domain.DailyStatistics {
	byDay := make(map[string]*domain.DailyStatistics)

	for _, row := range rows {
		if row == nil {
			continue
		}

		day, exists := byDay[row.Date]
		if !exists {
			day = domain.NewDailyStatistics(row.Date)
			byDay[row.Date] = day
		}

		day.CommitCount += row.CommitCount
		day.PRCreated += row.PRCreated
		day.PRMerged += row.PRMerged
		day.IssueCount += row.IssueCount
		day.ReviewCount += row.ReviewCount
		day.TotalAdditions += row.TotalAdditions
		day.TotalDeletions += row.TotalDeletions
	}

	days := make([]*domain.DailyStatistics, 0, len(byDay))
	for _, day := range byDay {
		days = append(days, day)
	}

	sort.Slice(days, func(i, j int) bool {
		return days[i].Date < days[j].Date
	})

	return days
}
