package domain

// YearlyStatistics は年別の統計情報を表す値オブジェクトです.
type YearlyStatistics struct {
	Year           int
	CommitCount    int
	PRCreated      int
	PRMerged       int
	IssueCount     int
	ReviewCount    int
	TotalAdditions int
	TotalDeletions int
}

// NewYearlyStatistics は新しいYearlyStatistics値オブジェクトを作成します.
func NewYearlyStatistics(year int) *YearlyStatistics {
	return &YearlyStatistics{
		Year: year,
	}
}

// UserStatistics はユーザーの統計情報を集約するドメインモデルです.
type UserStatistics struct {
	User                 *User
	TotalCommits         int
	TotalPRCreated       int
	TotalPRMerged        int
	TotalIssues          int
	TotalReviews         int
	TotalAdditions       int
	TotalDeletions       int
	FirstActivityYear    int
	PeakActivityYear     int
	PeakActivityCommits  int
	YearlyStats          map[int]*YearlyStatistics
	TopRepositories      []*RepositoryActivity
	LongTermRepositories []*RepositoryActivity
	PRToReviewRatio      float64 // PR作成数に対するレビュー数の比率
	RoleTransition       []RoleTransitionPoint
}

// RoleTransitionPoint はロール変化のポイントを表します.
type RoleTransitionPoint struct {
	Year        int
	PRCreated   int
	ReviewCount int
	Ratio       float64
	Description string
}

// NewUserStatistics は新しいUserStatisticsドメインモデルを作成します.
func NewUserStatistics(user *User) *UserStatistics {
	return &UserStatistics{
		User:                 user,
		YearlyStats:          make(map[int]*YearlyStatistics),
		TopRepositories:      make([]*RepositoryActivity, 0),
		LongTermRepositories: make([]*RepositoryActivity, 0),
		RoleTransition:       make([]RoleTransitionPoint, 0),
	}
}

// CalculatePRToReviewRatio はPR作成数に対するレビュー数の比率を計算します.
func (us *UserStatistics) CalculatePRToReviewRatio() {
	if us.TotalPRCreated > 0 {
		us.PRToReviewRatio = float64(us.TotalReviews) / float64(us.TotalPRCreated)
	}
}
