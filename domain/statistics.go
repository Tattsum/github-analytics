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

// DailyStatistics は日別の統計情報を表す値オブジェクトです.
// Date は "2006-01-02" 形式（UTC基準で丸めた日）のISO日付文字列です.
// 任意の日付範囲での絞り込みと時系列推移の描画に用います.
type DailyStatistics struct {
	Date           string
	CommitCount    int
	PRCreated      int
	PRMerged       int
	IssueCount     int
	ReviewCount    int
	TotalAdditions int
	TotalDeletions int
}

// NewDailyStatistics は新しいDailyStatistics値オブジェクトを作成します.
func NewDailyStatistics(date string) *DailyStatistics {
	return &DailyStatistics{
		Date: date,
	}
}

// UserStatistics はユーザーの統計情報を集約するドメインモデルです.
type UserStatistics struct {
	User                *User
	TotalCommits        int
	TotalPRCreated      int
	TotalPRMerged       int
	TotalIssues         int
	TotalReviews        int
	TotalAdditions      int
	TotalDeletions      int
	FirstActivityYear   int
	PeakActivityYear    int
	PeakActivityCommits int
	YearlyStats         map[int]*YearlyStatistics
	// DailyStats は日別の統計情報です（キーは "2006-01-02" 形式のISO日付文字列）.
	DailyStats           map[string]*DailyStatistics
	TopRepositories      []*RepositoryActivity
	LongTermRepositories []*RepositoryActivity
	// AllRepositories は関与した全リポジトリの活動内訳です（TopRepositories/LongTermRepositoriesとは別に全件を保持する）.
	AllRepositories []*RepositoryActivity
	PRToReviewRatio float64 // PR作成数に対するレビュー数の比率
	RoleTransition  []RoleTransitionPoint
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
		DailyStats:           make(map[string]*DailyStatistics),
		TopRepositories:      make([]*RepositoryActivity, 0),
		LongTermRepositories: make([]*RepositoryActivity, 0),
		AllRepositories:      make([]*RepositoryActivity, 0),
		RoleTransition:       make([]RoleTransitionPoint, 0),
	}
}

// CalculatePRToReviewRatio はPR作成数に対するレビュー数の比率を計算します.
func (us *UserStatistics) CalculatePRToReviewRatio() {
	if us.TotalPRCreated > 0 {
		us.PRToReviewRatio = float64(us.TotalReviews) / float64(us.TotalPRCreated)
	}
}
