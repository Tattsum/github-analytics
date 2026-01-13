package domain

import "time"

// ActivityType は活動の種類を表します.
type ActivityType string

const (
	ActivityTypeCommit  ActivityType = "commit"
	ActivityTypePR      ActivityType = "pull_request"
	ActivityTypeIssue   ActivityType = "issue"
	ActivityTypeReview  ActivityType = "review"
	ActivityTypePRMerge ActivityType = "pr_merge"
)

// Activity は1つの活動を表す値オブジェクトです.
type Activity struct {
	Type       ActivityType
	Repository string
	Date       time.Time
	Additions  int
	Deletions  int
	IsMerged   bool // PRの場合のみ有効
	IsReview   bool // Reviewの場合のみ有効
}

// NewActivity は新しいActivity値オブジェクトを作成します.
func NewActivity(activityType ActivityType, repo string, date time.Time, additions, deletions int) *Activity {
	return &Activity{
		Type:       activityType,
		Repository: repo,
		Date:       date,
		Additions:  additions,
		Deletions:  deletions,
	}
}

// RepositoryActivity はリポジトリごとの活動を集計した値オブジェクトです.
type RepositoryActivity struct {
	Repository     string
	CommitCount    int
	PRCount        int
	IssueCount     int
	ReviewCount    int
	TotalAdditions int
	TotalDeletions int
	FirstActivity  time.Time
	LastActivity   time.Time
}

// NewRepositoryActivity は新しいRepositoryActivity値オブジェクトを作成します.
func NewRepositoryActivity(repo string) *RepositoryActivity {
	return &RepositoryActivity{
		Repository: repo,
	}
}
