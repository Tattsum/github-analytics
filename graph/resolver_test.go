package graph

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Tattsum/github-analytics/application"
	"github.com/Tattsum/github-analytics/domain"
	"github.com/Tattsum/github-analytics/graph/model"
)

// fakeSnapshotReader is a configurable test double for application.SnapshotReader.
type fakeSnapshotReader struct {
	members     []*application.MemberStats
	member      *domain.UserStatistics
	teamSummary *application.TeamSummary
	teamDaily   []*domain.DailyStatistics
	repos       []*application.RepositoryStats
	repo        *application.RepositoryStats
	repoDaily   []*application.RepositoryDailyStats
	err         error
}

func (f *fakeSnapshotReader) LatestMembers(_ context.Context) ([]*application.MemberStats, error) {
	return f.members, f.err
}

func (f *fakeSnapshotReader) Member(_ context.Context, _ string) (*domain.UserStatistics, error) {
	return f.member, f.err
}

func (f *fakeSnapshotReader) TeamSummary(_ context.Context) (*application.TeamSummary, error) {
	return f.teamSummary, f.err
}

func (f *fakeSnapshotReader) TeamDailyStats(_ context.Context) ([]*domain.DailyStatistics, error) {
	return f.teamDaily, f.err
}

func (f *fakeSnapshotReader) Repositories(_ context.Context) ([]*application.RepositoryStats, error) {
	return f.repos, f.err
}

func (f *fakeSnapshotReader) Repository(_ context.Context, _ string) (*application.RepositoryStats, error) {
	return f.repo, f.err
}

func (f *fakeSnapshotReader) RepositoryDailyStats(_ context.Context) ([]*application.RepositoryDailyStats, error) {
	return f.repoDaily, f.err
}

func newTestQueryResolver(t *testing.T, reader application.SnapshotReader) QueryResolver {
	t.Helper()
	return NewResolver(reader).Query()
}

func TestQueryResolver_Members(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("boom")

	tests := []struct {
		name    string
		reader  *fakeSnapshotReader
		want    []*model.MemberStats
		wantErr bool
	}{
		{
			name: "maps member scalars",
			reader: &fakeSnapshotReader{
				members: []*application.MemberStats{
					{
						Login:           "octocat",
						Name:            "The Octocat",
						TotalCommits:    42,
						TotalPRCreated:  7,
						TotalPRMerged:   5,
						TotalIssues:     3,
						TotalReviews:    11,
						TotalAdditions:  120,
						TotalDeletions:  30,
						PRToReviewRatio: 1.57,
					},
				},
			},
			want: []*model.MemberStats{
				{
					Login:           "octocat",
					Name:            "The Octocat",
					TotalCommits:    42,
					TotalPRCreated:  7,
					TotalPRMerged:   5,
					TotalIssues:     3,
					TotalReviews:    11,
					TotalAdditions:  120,
					TotalDeletions:  30,
					PrToReviewRatio: 1.57,
				},
			},
		},
		{
			name:   "empty members yields empty slice",
			reader: &fakeSnapshotReader{members: nil},
			want:   []*model.MemberStats{},
		},
		{
			name:    "reader error is wrapped",
			reader:  &fakeSnapshotReader{err: sentinel},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newTestQueryResolver(t, tt.reader)

			got, err := r.Members(context.Background())
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, sentinel)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestQueryResolver_Member(t *testing.T) {
	t.Parallel()

	first := time.Date(2021, time.March, 2, 8, 0, 0, 0, time.UTC)
	last := time.Date(2023, time.July, 9, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		reader  *fakeSnapshotReader
		assert  func(t *testing.T, got *model.UserStatistics)
		wantNil bool
		wantErr bool
	}{
		{
			name: "maps drill-down and sorts yearly stats ascending",
			reader: &fakeSnapshotReader{
				member: &domain.UserStatistics{
					User:                &domain.User{Login: "octocat", Name: "The Octocat"},
					TotalCommits:        42,
					TotalReviews:        11,
					TotalPRCreated:      7,
					PRToReviewRatio:     1.57,
					FirstActivityYear:   2021,
					PeakActivityYear:    2022,
					PeakActivityCommits: 30,
					// Insertion order deliberately non-chronological to prove sorting.
					YearlyStats: map[int]*domain.YearlyStatistics{
						2023: {Year: 2023, CommitCount: 12, PRCreated: 4},
						2021: {Year: 2021, CommitCount: 5, PRCreated: 1},
						2022: {Year: 2022, CommitCount: 30, PRCreated: 2},
					},
					TopRepositories: []*domain.RepositoryActivity{
						{
							Repository:     "Tattsum/github-analytics",
							CommitCount:    25,
							PRCount:        6,
							IssueCount:     2,
							ReviewCount:    9,
							TotalAdditions: 100,
							TotalDeletions: 20,
							FirstActivity:  first,
							LastActivity:   last,
						},
					},
					LongTermRepositories: []*domain.RepositoryActivity{
						{Repository: "Tattsum/dotfiles", FirstActivity: first, LastActivity: last},
					},
					RoleTransition: []domain.RoleTransitionPoint{
						{Year: 2022, PRCreated: 2, ReviewCount: 9, Ratio: 4.5, Description: "shift to reviewer"},
					},
				},
			},
			assert: func(t *testing.T, got *model.UserStatistics) {
				t.Helper()
				assert.Equal(t, "octocat", got.Login)
				assert.Equal(t, "The Octocat", got.Name)
				assert.Equal(t, 42, got.TotalCommits)
				assert.InEpsilon(t, 1.57, got.PrToReviewRatio, 1e-9)

				require.Len(t, got.YearlyStats, 3)
				assert.Equal(t, []int{2021, 2022, 2023},
					[]int{got.YearlyStats[0].Year, got.YearlyStats[1].Year, got.YearlyStats[2].Year})
				assert.Equal(t, 5, got.YearlyStats[0].CommitCount)
				assert.Equal(t, 30, got.YearlyStats[1].CommitCount)

				require.Len(t, got.TopRepositories, 1)
				top := got.TopRepositories[0]
				assert.Equal(t, "Tattsum/github-analytics", top.Repository)
				assert.Equal(t, 6, top.PrCount)
				assert.Equal(t, "2021-03-02T08:00:00Z", top.FirstActivity)
				assert.Equal(t, "2023-07-09T12:00:00Z", top.LastActivity)

				require.Len(t, got.LongTermRepositories, 1)
				assert.Equal(t, "Tattsum/dotfiles", got.LongTermRepositories[0].Repository)

				require.Len(t, got.RoleTransition, 1)
				assert.Equal(t, "shift to reviewer", got.RoleTransition[0].Description)
				assert.InEpsilon(t, 4.5, got.RoleTransition[0].Ratio, 1e-9)
			},
		},
		{
			name:    "not found returns nil without error",
			reader:  &fakeSnapshotReader{member: nil},
			wantNil: true,
		},
		{
			name:    "reader error is wrapped",
			reader:  &fakeSnapshotReader{err: errors.New("boom")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newTestQueryResolver(t, tt.reader)

			got, err := r.Member(context.Background(), "octocat")
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.wantNil {
				assert.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			tt.assert(t, got)
		})
	}
}

func TestQueryResolver_TeamSummary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		reader  *fakeSnapshotReader
		want    *model.TeamSummary
		wantErr bool
	}{
		{
			name: "maps team aggregates",
			reader: &fakeSnapshotReader{
				teamSummary: &application.TeamSummary{
					MemberCount:     8,
					RepositoryCount: 4,
					TotalCommits:    123,
					TotalPRCreated:  45,
					TotalPRMerged:   40,
					TotalIssues:     17,
					TotalReviews:    60,
					TotalAdditions:  9000,
					TotalDeletions:  3000,
				},
			},
			want: &model.TeamSummary{
				MemberCount:     8,
				RepositoryCount: 4,
				TotalCommits:    123,
				TotalPRCreated:  45,
				TotalPRMerged:   40,
				TotalIssues:     17,
				TotalReviews:    60,
				TotalAdditions:  9000,
				TotalDeletions:  3000,
			},
		},
		{
			name:    "reader error is wrapped",
			reader:  &fakeSnapshotReader{err: errors.New("boom")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newTestQueryResolver(t, tt.reader)

			got, err := r.TeamSummary(context.Background())
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestQueryResolver_TeamDailyStats(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		reader  *fakeSnapshotReader
		want    []*model.DailyStatistics
		wantErr bool
	}{
		{
			name: "maps daily series preserving order",
			reader: &fakeSnapshotReader{
				teamDaily: []*domain.DailyStatistics{
					{Date: "2024-01-08", CommitCount: 8, PRCreated: 3, PRMerged: 2, IssueCount: 1, ReviewCount: 5, TotalAdditions: 80, TotalDeletions: 20},
					{Date: "2024-01-09", CommitCount: 4, PRCreated: 1, PRMerged: 1, IssueCount: 0, ReviewCount: 2, TotalAdditions: 40, TotalDeletions: 10},
				},
			},
			want: []*model.DailyStatistics{
				{Date: "2024-01-08", CommitCount: 8, PrCreated: 3, PrMerged: 2, IssueCount: 1, ReviewCount: 5, TotalAdditions: 80, TotalDeletions: 20},
				{Date: "2024-01-09", CommitCount: 4, PrCreated: 1, PrMerged: 1, IssueCount: 0, ReviewCount: 2, TotalAdditions: 40, TotalDeletions: 10},
			},
		},
		{
			name:   "no data yields empty slice",
			reader: &fakeSnapshotReader{teamDaily: nil},
			want:   []*model.DailyStatistics{},
		},
		{
			name:    "reader error is wrapped",
			reader:  &fakeSnapshotReader{err: errors.New("boom")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newTestQueryResolver(t, tt.reader)

			got, err := r.TeamDailyStats(context.Background())
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestQueryResolver_Repositories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		reader  *fakeSnapshotReader
		want    []*model.RepositoryStats
		wantErr bool
	}{
		{
			name: "maps repository totals and contributors",
			reader: &fakeSnapshotReader{
				repos: []*application.RepositoryStats{
					{
						NameWithOwner:    "Tattsum/github-analytics",
						TotalCommits:     50,
						TotalPRCreated:   20,
						TotalPRMerged:    18,
						TotalIssues:      9,
						TotalReviews:     33,
						TotalAdditions:   4000,
						TotalDeletions:   1500,
						ContributorCount: 2,
						Contributors: []*application.RepositoryContributor{
							{Login: "octocat", CommitCount: 30, PRCreated: 12, ReviewCount: 20, Additions: 2500, Deletions: 900},
							{Login: "hubot", CommitCount: 20, PRCreated: 8, ReviewCount: 13, Additions: 1500, Deletions: 600},
						},
					},
				},
			},
			want: []*model.RepositoryStats{
				{
					NameWithOwner: "Tattsum/github-analytics",
					Total: &model.RepositoryTotals{
						Commits:   50,
						PrCreated: 20,
						PrMerged:  18,
						Issues:    9,
						Reviews:   33,
						Additions: 4000,
						Deletions: 1500,
					},
					ContributorCount: 2,
					Contributors: []*model.RepositoryContributor{
						{Login: "octocat", CommitCount: 30, PrCreated: 12, ReviewCount: 20, Additions: 2500, Deletions: 900},
						{Login: "hubot", CommitCount: 20, PrCreated: 8, ReviewCount: 13, Additions: 1500, Deletions: 600},
					},
				},
			},
		},
		{
			name:   "empty repositories yields empty slice",
			reader: &fakeSnapshotReader{repos: nil},
			want:   []*model.RepositoryStats{},
		},
		{
			name:    "reader error is wrapped",
			reader:  &fakeSnapshotReader{err: errors.New("boom")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newTestQueryResolver(t, tt.reader)

			got, err := r.Repositories(context.Background())
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestQueryResolver_Repository(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		reader  *fakeSnapshotReader
		want    *model.RepositoryStats
		wantNil bool
		wantErr bool
	}{
		{
			name: "maps single repository",
			reader: &fakeSnapshotReader{
				repo: &application.RepositoryStats{
					NameWithOwner:    "Tattsum/dotfiles",
					TotalCommits:     12,
					ContributorCount: 1,
					Contributors: []*application.RepositoryContributor{
						{Login: "octocat", CommitCount: 12},
					},
				},
			},
			want: &model.RepositoryStats{
				NameWithOwner: "Tattsum/dotfiles",
				Total: &model.RepositoryTotals{
					Commits: 12,
				},
				ContributorCount: 1,
				Contributors: []*model.RepositoryContributor{
					{Login: "octocat", CommitCount: 12},
				},
			},
		},
		{
			name:    "not found returns nil without error",
			reader:  &fakeSnapshotReader{repo: nil},
			wantNil: true,
		},
		{
			name:    "reader error is wrapped",
			reader:  &fakeSnapshotReader{err: errors.New("boom")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newTestQueryResolver(t, tt.reader)

			got, err := r.Repository(context.Background(), "Tattsum/dotfiles")
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.wantNil {
				assert.Nil(t, got)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
