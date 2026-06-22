package application

import (
	"context"
	"time"

	"github.com/Tattsum/github-analytics/domain"
)

// MemberStats はメンバー横断比較・ランキングに用いる、比較可能なスカラー指標の集合です.
// ランキング・ソート・比較はフロントエンドで計算するため、ここでは並び順を持ちません.
type MemberStats struct {
	Login          string
	Name           string
	TotalCommits   int
	TotalPRCreated int
	TotalPRMerged  int
	TotalIssues    int
	TotalReviews   int
	TotalAdditions int
	TotalDeletions int
	// PRToReviewRatio はPR作成数に対するレビュー数の比率です.
	PRToReviewRatio float64
}

// TeamSummary はチーム全体の合計・集計値を表します.
type TeamSummary struct {
	MemberCount     int
	RepositoryCount int
	TotalCommits    int
	TotalPRCreated  int
	TotalPRMerged   int
	TotalIssues     int
	TotalReviews    int
	TotalAdditions  int
	TotalDeletions  int
}

// RepositoryContributor はリポジトリ軸でのメンバーごとの貢献内訳です.
type RepositoryContributor struct {
	Login       string
	CommitCount int
	PRCreated   int
	ReviewCount int
	Additions   int
	Deletions   int
}

// RepositoryStats はリポジトリ軸での横断集計を表します.
type RepositoryStats struct {
	NameWithOwner    string
	TotalCommits     int
	TotalPRCreated   int
	TotalPRMerged    int
	TotalIssues      int
	TotalReviews     int
	TotalAdditions   int
	TotalDeletions   int
	ContributorCount int
	Contributors     []*RepositoryContributor
}

// Snapshot はバッチ実行1回分の集計済みスナップショットです.
// captured_at をキーに蓄積され、Web はデフォルトで最新スナップショットを参照します.
type Snapshot struct {
	CapturedAt time.Time
	// Members はメンバーごとの集計済み統計です（member-levelスカラー・member×year・member×repositoryを含む）.
	Members []*domain.UserStatistics
}

// SnapshotReader は最新スナップショットを読み取るための契約です.
// 実装は infrastructure 層（ent/Postgres）が提供します.
type SnapshotReader interface {
	// LatestMembers は最新スナップショットのメンバー横断スカラー指標を返します.
	LatestMembers(ctx context.Context) ([]*MemberStats, error)
	// Member は指定ログインのドリルダウン統計（年次推移・トップリポジトリ等）を返します.
	Member(ctx context.Context, login string) (*domain.UserStatistics, error)
	// TeamSummary はチーム全体の合計・集計値を返します.
	TeamSummary(ctx context.Context) (*TeamSummary, error)
	// TeamDailyStats はチーム全体の日別合計を、日付昇順の時系列で返します.
	// メンバー横断で同一日の指標を合算したもので、期間絞り込み・推移グラフのデータ源です.
	TeamDailyStats(ctx context.Context) ([]*domain.DailyStatistics, error)
	// Repositories はリポジトリ軸の横断集計を返します.
	Repositories(ctx context.Context) ([]*RepositoryStats, error)
	// Repository は指定リポジトリの集計を返します.
	Repository(ctx context.Context, nameWithOwner string) (*RepositoryStats, error)
}

// SnapshotWriter はバッチが集計済みスナップショットを永続化するための契約です.
// 実装は infrastructure 層（ent/Postgres）が提供し、冪等に1スナップショットを書き込みます.
type SnapshotWriter interface {
	// Save は1回分の集計済みスナップショットを保存します.
	Save(ctx context.Context, snapshot *Snapshot) error
}
