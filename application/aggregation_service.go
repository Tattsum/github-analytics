package application

import "sort"

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
