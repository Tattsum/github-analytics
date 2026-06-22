// Package snapshotdb provides the ent/PostgreSQL-backed read side for
// aggregated analytics snapshots. It implements application.SnapshotReader.
//
// It lives in its own package (not the infrastructure root) because the
// application package already depends on the infrastructure root package, so a
// reader placed there would create an import cycle.
package snapshotdb

import (
	"context"
	"fmt"
	"sort"

	"entgo.io/ent/dialect/sql"

	"github.com/Tattsum/github-analytics/application"
	"github.com/Tattsum/github-analytics/domain"
	"github.com/Tattsum/github-analytics/infrastructure/ent"
	"github.com/Tattsum/github-analytics/infrastructure/ent/snapshot"
)

// SnapshotReader は ent クライアントを用いて最新スナップショットを読み取り、
// application.SnapshotReader を満たす実装です.
type SnapshotReader struct {
	client *ent.Client
}

// SnapshotReader が application.SnapshotReader を満たすことをコンパイル時に保証します.
var _ application.SnapshotReader = (*SnapshotReader)(nil)

// NewSnapshotReader は ent クライアントを背後に持つ SnapshotReader を作成します.
func NewSnapshotReader(client *ent.Client) *SnapshotReader {
	return &SnapshotReader{client: client}
}

// latest は captured_at が最も新しいスナップショットを、要求された各 stat エッジを eager-load して返します.
// スナップショットが1件も存在しない場合は (nil, nil) を返し、呼び出し元が「未取得」として扱えるようにします.
func (r *SnapshotReader) latest(
	ctx context.Context,
	with func(*ent.SnapshotQuery) *ent.SnapshotQuery,
) (*ent.Snapshot, error) {
	query := r.client.Snapshot.
		Query().
		Order(snapshot.ByCapturedAt(sql.OrderDesc()))

	if with != nil {
		query = with(query)
	}

	snap, err := query.First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("query latest snapshot: %w", err)
	}

	return snap, nil
}

// LatestMembers は最新スナップショットのメンバー横断スカラー指標を返します.
// スナップショットが無い場合は空スライスを返します（エラーにしません）.
func (r *SnapshotReader) LatestMembers(ctx context.Context) ([]*application.MemberStats, error) {
	snap, err := r.latest(ctx, func(q *ent.SnapshotQuery) *ent.SnapshotQuery {
		return q.WithMemberStats()
	})
	if err != nil {
		return nil, err
	}

	if snap == nil {
		return []*application.MemberStats{}, nil
	}

	members := make([]*application.MemberStats, 0, len(snap.Edges.MemberStats))
	for _, ms := range snap.Edges.MemberStats {
		members = append(members, toMemberStats(ms))
	}

	return members, nil
}

// Member は指定ログインのドリルダウン統計（年次推移・全リポジトリ内訳）を返します.
// 該当メンバーが最新スナップショットに存在しない場合は (nil, nil) を返します.
func (r *SnapshotReader) Member(ctx context.Context, login string) (*domain.UserStatistics, error) {
	snap, err := r.latest(ctx, func(q *ent.SnapshotQuery) *ent.SnapshotQuery {
		return q.
			WithMemberStats().
			WithMemberYearStats().
			WithMemberDayStats().
			WithMemberRepoStats()
	})
	if err != nil {
		return nil, err
	}

	if snap == nil {
		return nil, nil
	}

	var member *ent.MemberStat
	for _, ms := range snap.Edges.MemberStats {
		if ms.Login == login {
			member = ms

			break
		}
	}

	if member == nil {
		return nil, nil
	}

	return buildUserStatistics(member, snap.Edges.MemberYearStats, snap.Edges.MemberDayStats, snap.Edges.MemberRepoStats), nil
}

// TeamDailyStats は最新スナップショットのメンバー日別統計をメンバー横断で同一日に合算し、
// チーム全体の日別合計を日付昇順の時系列で返します.
// スナップショットが無い場合は空スライスを返します（エラーにしません）.
func (r *SnapshotReader) TeamDailyStats(ctx context.Context) ([]*domain.DailyStatistics, error) {
	snap, err := r.latest(ctx, func(q *ent.SnapshotQuery) *ent.SnapshotQuery {
		return q.WithMemberDayStats()
	})
	if err != nil {
		return nil, err
	}

	if snap == nil {
		return []*domain.DailyStatistics{}, nil
	}

	return application.AggregateTeamDaily(toDailyStatistics(snap.Edges.MemberDayStats)), nil
}

// toDailyStatistics は ent の MemberDayStat 群を domain.DailyStatistics へマッピングします.
// 各行は1メンバーの1日分で、チーム合算は呼び出し元（AggregateTeamDaily）が行います.
func toDailyStatistics(stats []*ent.MemberDayStat) []*domain.DailyStatistics {
	out := make([]*domain.DailyStatistics, 0, len(stats))
	for _, mds := range stats {
		daily := domain.NewDailyStatistics(mds.Day)
		daily.CommitCount = mds.CommitCount
		daily.PRCreated = mds.PrCreated
		daily.PRMerged = mds.PrMerged
		daily.IssueCount = mds.IssueCount
		daily.ReviewCount = mds.ReviewCount
		daily.TotalAdditions = mds.Additions
		daily.TotalDeletions = mds.Deletions
		out = append(out, daily)
	}

	return out
}

// TeamSummary はチーム全体の合計・集計値を返します.
// RepositoryCount は最新スナップショット内のユニークな nameWithOwner 数です.
func (r *SnapshotReader) TeamSummary(ctx context.Context) (*application.TeamSummary, error) {
	snap, err := r.latest(ctx, func(q *ent.SnapshotQuery) *ent.SnapshotQuery {
		return q.
			WithMemberStats().
			WithMemberRepoStats()
	})
	if err != nil {
		return nil, err
	}

	if snap == nil {
		return &application.TeamSummary{}, nil
	}

	members := make([]*application.MemberStats, 0, len(snap.Edges.MemberStats))
	for _, ms := range snap.Edges.MemberStats {
		members = append(members, toMemberStats(ms))
	}

	summary := application.SummarizeTeam(members)
	summary.RepositoryCount = countRepositories(snap.Edges.MemberRepoStats)

	return summary, nil
}

// Repositories はリポジトリ軸の横断集計を返します.
// MemberRepoStat を nameWithOwner でグルーピングし、リポジトリごとの合計・貢献者一覧へ再集計します.
func (r *SnapshotReader) Repositories(ctx context.Context) ([]*application.RepositoryStats, error) {
	snap, err := r.latest(ctx, func(q *ent.SnapshotQuery) *ent.SnapshotQuery {
		return q.WithMemberRepoStats()
	})
	if err != nil {
		return nil, err
	}

	if snap == nil {
		return []*application.RepositoryStats{}, nil
	}

	return application.AggregateRepositories(toRepoStatInputs(snap.Edges.MemberRepoStats)), nil
}

// Repository は指定リポジトリの集計を返します.
// 該当リポジトリが最新スナップショットに存在しない場合は (nil, nil) を返します.
func (r *SnapshotReader) Repository(
	ctx context.Context,
	nameWithOwner string,
) (*application.RepositoryStats, error) {
	repos, err := r.Repositories(ctx)
	if err != nil {
		return nil, err
	}

	for _, repo := range repos {
		if repo.NameWithOwner == nameWithOwner {
			return repo, nil
		}
	}

	return nil, nil
}

// toMemberStats は ent の MemberStat を application.MemberStats へマッピングします.
// MemberStat スキーマは name を保持しないため、Name は login をそのまま用います.
func toMemberStats(ms *ent.MemberStat) *application.MemberStats {
	return &application.MemberStats{
		Login:           ms.Login,
		Name:            ms.Login,
		TotalCommits:    ms.TotalCommits,
		TotalPRCreated:  ms.TotalPrCreated,
		TotalPRMerged:   ms.TotalPrMerged,
		TotalIssues:     ms.TotalIssues,
		TotalReviews:    ms.TotalReviews,
		TotalAdditions:  ms.TotalAdditions,
		TotalDeletions:  ms.TotalDeletions,
		PRToReviewRatio: ms.PrToReviewRatio,
	}
}

// toRepoStatInputs は ent の MemberRepoStat 群を集計関数の入力構造体へマッピングします.
func toRepoStatInputs(stats []*ent.MemberRepoStat) []*application.MemberRepoStat {
	inputs := make([]*application.MemberRepoStat, 0, len(stats))
	for _, mrs := range stats {
		inputs = append(inputs, &application.MemberRepoStat{
			Login:         mrs.Login,
			NameWithOwner: mrs.NameWithOwner,
			CommitCount:   mrs.CommitCount,
			PRCreated:     mrs.PrCreated,
			PRMerged:      mrs.PrMerged,
			IssueCount:    mrs.IssueCount,
			ReviewCount:   mrs.ReviewCount,
			Additions:     mrs.Additions,
			Deletions:     mrs.Deletions,
		})
	}

	return inputs
}

// countRepositories は MemberRepoStat 群に含まれるユニークな nameWithOwner 数を返します.
func countRepositories(stats []*ent.MemberRepoStat) int {
	seen := make(map[string]struct{}, len(stats))
	for _, mrs := range stats {
		seen[mrs.NameWithOwner] = struct{}{}
	}

	return len(seen)
}

// buildUserStatistics は ent の各 stat 行から、指定メンバーの UserStatistics を組み立てます.
// 年次推移は MemberYearStat、全リポジトリ内訳は MemberRepoStat（当該 login 分のみ）から構築します.
func buildUserStatistics(
	member *ent.MemberStat,
	yearStats []*ent.MemberYearStat,
	dayStats []*ent.MemberDayStat,
	repoStats []*ent.MemberRepoStat,
) *domain.UserStatistics {
	stats := domain.NewUserStatistics(domain.NewUser(member.Login, member.Login, ""))
	stats.TotalCommits = member.TotalCommits
	stats.TotalPRCreated = member.TotalPrCreated
	stats.TotalPRMerged = member.TotalPrMerged
	stats.TotalIssues = member.TotalIssues
	stats.TotalReviews = member.TotalReviews
	stats.TotalAdditions = member.TotalAdditions
	stats.TotalDeletions = member.TotalDeletions
	stats.FirstActivityYear = member.FirstActivityYear
	stats.PeakActivityYear = member.PeakActivityYear
	stats.PeakActivityCommits = member.PeakActivityCommits
	stats.PRToReviewRatio = member.PrToReviewRatio

	for _, ys := range yearStats {
		if ys.Login != member.Login {
			continue
		}

		yearly := domain.NewYearlyStatistics(ys.Year)
		yearly.CommitCount = ys.CommitCount
		yearly.PRCreated = ys.PrCreated
		yearly.PRMerged = ys.PrMerged
		yearly.IssueCount = ys.IssueCount
		yearly.ReviewCount = ys.ReviewCount
		yearly.TotalAdditions = ys.Additions
		yearly.TotalDeletions = ys.Deletions
		stats.YearlyStats[ys.Year] = yearly
	}

	for _, ds := range dayStats {
		if ds.Login != member.Login {
			continue
		}

		daily := domain.NewDailyStatistics(ds.Day)
		daily.CommitCount = ds.CommitCount
		daily.PRCreated = ds.PrCreated
		daily.PRMerged = ds.PrMerged
		daily.IssueCount = ds.IssueCount
		daily.ReviewCount = ds.ReviewCount
		daily.TotalAdditions = ds.Additions
		daily.TotalDeletions = ds.Deletions
		stats.DailyStats[ds.Day] = daily
	}

	stats.AllRepositories = buildMemberRepositories(member.Login, repoStats)
	stats.TopRepositories = topRepositoriesByCommits(stats.AllRepositories, topRepositoryCount)

	return stats
}

// topRepositoryCount はメンバー詳細で表示する上位リポジトリ数です.
const topRepositoryCount = 3

// topRepositoriesByCommits はコミット数の降順で上位n件のリポジトリを返します.
// スナップショットには活動日時が無いため LongTermRepositories は算出できません.
func topRepositoriesByCommits(repos []*domain.RepositoryActivity, n int) []*domain.RepositoryActivity {
	sorted := make([]*domain.RepositoryActivity, len(repos))
	copy(sorted, repos)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].CommitCount > sorted[j].CommitCount
	})

	if len(sorted) > n {
		sorted = sorted[:n]
	}

	return sorted
}

// buildMemberRepositories は当該メンバーのリポジトリ活動内訳を、コミット数の降順で組み立てます.
func buildMemberRepositories(login string, repoStats []*ent.MemberRepoStat) []*domain.RepositoryActivity {
	repos := make([]*domain.RepositoryActivity, 0)
	for _, mrs := range repoStats {
		if mrs.Login != login {
			continue
		}

		repo := domain.NewRepositoryActivity(mrs.NameWithOwner)
		repo.CommitCount = mrs.CommitCount
		repo.PRCount = mrs.PrCreated
		repo.IssueCount = mrs.IssueCount
		repo.ReviewCount = mrs.ReviewCount
		repo.TotalAdditions = mrs.Additions
		repo.TotalDeletions = mrs.Deletions
		repos = append(repos, repo)
	}

	return repos
}
