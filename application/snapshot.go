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
	// DailyStats はこのリポジトリ内での当該メンバーの日別活動の時系列です（日付昇順）.
	// リポジトリ内メンバー間の時系列比較（多系列の重ね合わせ）に用います.
	// 単一リポジトリのドリルダウン（Repository）でのみ設定され、一覧（Repositories）では nil です.
	DailyStats []*domain.DailyStatistics
}

// MemberRepoDayStat はメンバー×リポジトリ×日1件分の集計済みメトリクスです.
// リポジトリ間・リポジトリ内メンバー間の時系列比較の集計入力になります.
type MemberRepoDayStat struct {
	Login         string
	NameWithOwner string
	Day           string
	CommitCount   int
	PRCreated     int
	PRMerged      int
	IssueCount    int
	ReviewCount   int
	Additions     int
	Deletions     int
}

// RepoMeta はリポジトリの所有者メタ情報です（スナップショット内で1リポジトリ1件）.
// 組織内リポジトリへの絞り込み判定の権威的な情報源です.
type RepoMeta struct {
	NameWithOwner string
	Owner         string
	// OwnerType は所有者の種別（"Organization" / "User"）です（不明な場合は空文字）.
	OwnerType string
}

// RepositoryDailyStats はリポジトリ軸での日別時系列（メンバー横断で合算済み）です.
// 複数リポジトリの活動推移を重ね合わせて比較するためのデータ源で、所有者メタを同梱します.
type RepositoryDailyStats struct {
	NameWithOwner string
	Owner         string
	OwnerType     string
	// DailyStats はこのリポジトリの日別合計の時系列です（日付昇順）.
	DailyStats []*domain.DailyStatistics
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
	// Repository は指定リポジトリの集計を返します（貢献者ごとの日別時系列を含む）.
	Repository(ctx context.Context, nameWithOwner string) (*RepositoryStats, error)
	// RepositoryDailyStats は各リポジトリの日別合計を、所有者メタ付きで返します.
	// 複数リポジトリの活動推移を重ね合わせて比較するためのデータ源です.
	RepositoryDailyStats(ctx context.Context) ([]*RepositoryDailyStats, error)
}

// SnapshotWriter はバッチが集計済みスナップショットを永続化するための契約です.
// 実装は infrastructure 層（ent/Postgres）が提供し、冪等に1スナップショットを書き込みます.
type SnapshotWriter interface {
	// Save は1回分の集計済みスナップショットを保存します.
	Save(ctx context.Context, snapshot *Snapshot) error
}
