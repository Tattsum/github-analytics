package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// assertActivity はActivityのアサーションを行います.
func assertActivity(t *testing.T, got, want *Activity) {
	t.Helper()

	assert.Equal(t, want.Type, got.Type, "Type should match")
	assert.Equal(t, want.Repository, got.Repository, "Repository should match")
	assert.True(t, got.Date.Equal(want.Date), "Date should match")
	assert.Equal(t, want.Additions, got.Additions, "Additions should match")
	assert.Equal(t, want.Deletions, got.Deletions, "Deletions should match")
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
			assert.Equal(t, tt.want.Repository, got.Repository, "Repository should match")
			assert.Equal(t, 0, got.CommitCount, "CommitCount should be 0")
		})
	}
}
