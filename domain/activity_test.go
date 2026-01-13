package domain

import (
	"testing"
	"time"
)

// assertActivity はActivityのアサーションを行います.
func assertActivity(t *testing.T, got, want *Activity) {
	t.Helper()

	if got.Type != want.Type {
		t.Errorf("NewActivity().Type = %v, want %v", got.Type, want.Type)
	}

	if got.Repository != want.Repository {
		t.Errorf("NewActivity().Repository = %v, want %v", got.Repository, want.Repository)
	}

	if !got.Date.Equal(want.Date) {
		t.Errorf("NewActivity().Date = %v, want %v", got.Date, want.Date)
	}

	if got.Additions != want.Additions {
		t.Errorf("NewActivity().Additions = %v, want %v", got.Additions, want.Additions)
	}

	if got.Deletions != want.Deletions {
		t.Errorf("NewActivity().Deletions = %v, want %v", got.Deletions, want.Deletions)
	}
}

func TestNewActivity(t *testing.T) {
	t.Parallel()

	now := time.Now()
	tests := []struct {
		name         string
		activityType ActivityType
		repo         string
		date         time.Time
		additions    int
		deletions    int
		want         *Activity
	}{
		{
			name:         "コミット活動の作成",
			activityType: ActivityTypeCommit,
			repo:         "owner/repo",
			date:         now,
			additions:    100,
			deletions:    50,
			want: &Activity{
				Type:       ActivityTypeCommit,
				Repository: "owner/repo",
				Date:       now,
				Additions:  100,
				Deletions:  50,
			},
		},
		{
			name:         "PR活動の作成",
			activityType: ActivityTypePR,
			repo:         "owner/repo",
			date:         now,
			additions:    200,
			deletions:    100,
			want: &Activity{
				Type:       ActivityTypePR,
				Repository: "owner/repo",
				Date:       now,
				Additions:  200,
				Deletions:  100,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NewActivity(tt.activityType, tt.repo, tt.date, tt.additions, tt.deletions)
			assertActivity(t, got, tt.want)
		})
	}
}

func TestNewRepositoryActivity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		repo string
		want *RepositoryActivity
	}{
		{
			name: "リポジトリ活動の作成",
			repo: "owner/repo",
			want: &RepositoryActivity{
				Repository: "owner/repo",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NewRepositoryActivity(tt.repo)
			if got.Repository != tt.want.Repository {
				t.Errorf("NewRepositoryActivity().Repository = %v, want %v", got.Repository, tt.want.Repository)
			}

			if got.CommitCount != 0 {
				t.Errorf("NewRepositoryActivity().CommitCount = %v, want 0", got.CommitCount)
			}
		})
	}
}
