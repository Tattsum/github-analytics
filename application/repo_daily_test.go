package application

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Tattsum/github-analytics/domain"
	"github.com/Tattsum/github-analytics/infrastructure"
)

func TestAggregateRepositoryDaily(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		stats []*MemberRepoDayStat
		metas []*RepoMeta
		want  []*RepositoryDailyStats
	}{
		{
			name: "rows from different logins for same repo/day are summed into one DailyStatistics",
			stats: []*MemberRepoDayStat{
				{
					Login:         "alice",
					NameWithOwner: "Tattsum/foo",
					Day:           "2026-03-01",
					CommitCount:   7,
					PRCreated:     3,
					PRMerged:      2,
					IssueCount:    4,
					ReviewCount:   5,
					Additions:     120,
					Deletions:     45,
				},
				{
					Login:         "bob",
					NameWithOwner: "Tattsum/foo",
					Day:           "2026-03-01",
					CommitCount:   11,
					PRCreated:     6,
					PRMerged:      1,
					IssueCount:    2,
					ReviewCount:   8,
					Additions:     300,
					Deletions:     90,
				},
			},
			metas: []*RepoMeta{
				{NameWithOwner: "Tattsum/foo", Owner: "Tattsum", OwnerType: "Organization"},
			},
			want: []*RepositoryDailyStats{
				{
					NameWithOwner: "Tattsum/foo",
					Owner:         "Tattsum",
					OwnerType:     "Organization",
					DailyStats: []*domain.DailyStatistics{
						{
							Date:           "2026-03-01",
							CommitCount:    18,
							PRCreated:      9,
							PRMerged:       3,
							IssueCount:     6,
							ReviewCount:    13,
							TotalAdditions: 420,
							TotalDeletions: 135,
						},
					},
				},
			},
		},
		{
			name: "repos sorted ascending by NameWithOwner and each repo's DailyStats sorted ascending by Date",
			stats: []*MemberRepoDayStat{
				{
					Login:         "alice",
					NameWithOwner: "Tattsum/zeta",
					Day:           "2026-03-02",
					CommitCount:   9,
					PRCreated:     2,
					PRMerged:      1,
					IssueCount:    3,
					ReviewCount:   4,
					Additions:     50,
					Deletions:     10,
				},
				{
					Login:         "alice",
					NameWithOwner: "Tattsum/alpha",
					Day:           "2026-03-05",
					CommitCount:   13,
					PRCreated:     5,
					PRMerged:      4,
					IssueCount:    6,
					ReviewCount:   7,
					Additions:     200,
					Deletions:     80,
				},
				{
					Login:         "alice",
					NameWithOwner: "Tattsum/alpha",
					Day:           "2026-03-01",
					CommitCount:   8,
					PRCreated:     1,
					PRMerged:      1,
					IssueCount:    2,
					ReviewCount:   3,
					Additions:     60,
					Deletions:     20,
				},
			},
			metas: []*RepoMeta{
				{NameWithOwner: "Tattsum/alpha", Owner: "Tattsum", OwnerType: "Organization"},
				{NameWithOwner: "Tattsum/zeta", Owner: "Tattsum", OwnerType: "User"},
			},
			want: []*RepositoryDailyStats{
				{
					NameWithOwner: "Tattsum/alpha",
					Owner:         "Tattsum",
					OwnerType:     "Organization",
					DailyStats: []*domain.DailyStatistics{
						{
							Date:           "2026-03-01",
							CommitCount:    8,
							PRCreated:      1,
							PRMerged:       1,
							IssueCount:     2,
							ReviewCount:    3,
							TotalAdditions: 60,
							TotalDeletions: 20,
						},
						{
							Date:           "2026-03-05",
							CommitCount:    13,
							PRCreated:      5,
							PRMerged:       4,
							IssueCount:     6,
							ReviewCount:    7,
							TotalAdditions: 200,
							TotalDeletions: 80,
						},
					},
				},
				{
					NameWithOwner: "Tattsum/zeta",
					Owner:         "Tattsum",
					OwnerType:     "User",
					DailyStats: []*domain.DailyStatistics{
						{
							Date:           "2026-03-02",
							CommitCount:    9,
							PRCreated:      2,
							PRMerged:       1,
							IssueCount:     3,
							ReviewCount:    4,
							TotalAdditions: 50,
							TotalDeletions: 10,
						},
					},
				},
			},
		},
		{
			name: "repo with no matching RepoMeta gets empty Owner and OwnerType",
			stats: []*MemberRepoDayStat{
				{
					Login:         "carol",
					NameWithOwner: "Tattsum/orphan",
					Day:           "2026-03-03",
					CommitCount:   14,
					PRCreated:     6,
					PRMerged:      3,
					IssueCount:    2,
					ReviewCount:   9,
					Additions:     77,
					Deletions:     33,
				},
			},
			metas: []*RepoMeta{
				{NameWithOwner: "Tattsum/other", Owner: "Tattsum", OwnerType: "Organization"},
			},
			want: []*RepositoryDailyStats{
				{
					NameWithOwner: "Tattsum/orphan",
					Owner:         "",
					OwnerType:     "",
					DailyStats: []*domain.DailyStatistics{
						{
							Date:           "2026-03-03",
							CommitCount:    14,
							PRCreated:      6,
							PRMerged:       3,
							IssueCount:     2,
							ReviewCount:    9,
							TotalAdditions: 77,
							TotalDeletions: 33,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := AggregateRepositoryDaily(tt.stats, tt.metas)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAggregateRepositoryContributorDaily(t *testing.T) {
	t.Parallel()

	stats := []*MemberRepoDayStat{
		{
			Login:         "alice",
			NameWithOwner: "Tattsum/foo",
			Day:           "2026-04-02",
			CommitCount:   12,
			PRCreated:     4,
			PRMerged:      3,
			IssueCount:    5,
			ReviewCount:   6,
			Additions:     210,
			Deletions:     70,
		},
		{
			Login:         "alice",
			NameWithOwner: "Tattsum/foo",
			Day:           "2026-04-01",
			CommitCount:   8,
			PRCreated:     2,
			PRMerged:      1,
			IssueCount:    3,
			ReviewCount:   4,
			Additions:     90,
			Deletions:     30,
		},
		{
			Login:         "bob",
			NameWithOwner: "Tattsum/foo",
			Day:           "2026-04-01",
			CommitCount:   15,
			PRCreated:     7,
			PRMerged:      5,
			IssueCount:    2,
			ReviewCount:   9,
			Additions:     400,
			Deletions:     120,
		},
		{
			Login:         "alice",
			NameWithOwner: "Tattsum/other",
			Day:           "2026-04-01",
			CommitCount:   99,
			PRCreated:     99,
			PRMerged:      99,
			IssueCount:    99,
			ReviewCount:   99,
			Additions:     999,
			Deletions:     999,
		},
	}

	want := map[string][]*domain.DailyStatistics{
		"alice": {
			{
				Date:           "2026-04-01",
				CommitCount:    8,
				PRCreated:      2,
				PRMerged:       1,
				IssueCount:     3,
				ReviewCount:    4,
				TotalAdditions: 90,
				TotalDeletions: 30,
			},
			{
				Date:           "2026-04-02",
				CommitCount:    12,
				PRCreated:      4,
				PRMerged:       3,
				IssueCount:     5,
				ReviewCount:    6,
				TotalAdditions: 210,
				TotalDeletions: 70,
			},
		},
		"bob": {
			{
				Date:           "2026-04-01",
				CommitCount:    15,
				PRCreated:      7,
				PRMerged:       5,
				IssueCount:     2,
				ReviewCount:    9,
				TotalAdditions: 400,
				TotalDeletions: 120,
			},
		},
	}

	got := AggregateRepositoryContributorDaily(stats, "Tattsum/foo")

	assert.Equal(t, want, got, "rows for other repos must be excluded, grouped by login, sorted by Date")
}

func TestStatisticsServiceCalculateStatisticsRepoDailyStats(t *testing.T) {
	t.Parallel()

	const (
		repoFoo = "Tattsum/foo"
		repoBar = "Tattsum/bar"
	)

	dayOne := time.Date(2026, 5, 1, 9, 30, 0, 0, time.UTC)
	dayTwo := time.Date(2026, 5, 2, 18, 15, 0, 0, time.UTC)

	commitFooDay1 := domain.NewActivity(domain.ActivityTypeCommit, repoFoo, dayOne, 50, 20)
	commitFooDay1.RepositoryOwner = "Tattsum"
	commitFooDay1.RepositoryOwnerType = "Organization"

	commitFooDay2 := domain.NewActivity(domain.ActivityTypeCommit, repoFoo, dayTwo, 30, 10)
	commitFooDay2.RepositoryOwner = "Tattsum"
	commitFooDay2.RepositoryOwnerType = "Organization"

	commitBarDay1 := domain.NewActivity(domain.ActivityTypeCommit, repoBar, dayOne, 70, 25)
	commitBarDay1.RepositoryOwner = "octocat"
	commitBarDay1.RepositoryOwnerType = "User"

	prFooDay1Merged := domain.NewActivity(domain.ActivityTypePR, repoFoo, dayOne, 12, 4)
	prFooDay1Merged.RepositoryOwner = "Tattsum"
	prFooDay1Merged.RepositoryOwnerType = "Organization"
	prFooDay1Merged.IsMerged = true

	prFooDay1Open := domain.NewActivity(domain.ActivityTypePR, repoFoo, dayOne, 8, 3)
	prFooDay1Open.RepositoryOwner = "Tattsum"
	prFooDay1Open.RepositoryOwnerType = "Organization"
	prFooDay1Open.IsMerged = false

	prBarDay2Merged := domain.NewActivity(domain.ActivityTypePR, repoBar, dayTwo, 6, 2)
	prBarDay2Merged.RepositoryOwner = "octocat"
	prBarDay2Merged.RepositoryOwnerType = "User"
	prBarDay2Merged.IsMerged = true

	issueBarDay1 := domain.NewActivity(domain.ActivityTypeIssue, repoBar, dayOne, 0, 0)
	issueBarDay1.RepositoryOwner = "octocat"
	issueBarDay1.RepositoryOwnerType = "User"

	reviewFooDay2 := domain.NewActivity(domain.ActivityTypeReview, repoFoo, dayTwo, 0, 0)
	reviewFooDay2.RepositoryOwner = "Tattsum"
	reviewFooDay2.RepositoryOwnerType = "Organization"

	data := &infrastructure.UserActivityData{
		User:    domain.NewUser("alice", "Alice Example", "2018-01-02T03:04:05Z"),
		Commits: []*domain.Activity{commitFooDay1, commitFooDay2, commitBarDay1},
		PRs:     []*domain.Activity{prFooDay1Merged, prFooDay1Open, prBarDay2Merged},
		Issues:  []*domain.Activity{issueBarDay1},
		Reviews: []*domain.Activity{reviewFooDay2},
	}

	service := NewStatisticsService()

	stats, err := service.CalculateStatistics(data)
	require.NoError(t, err)
	require.NotNil(t, stats)

	wantRepoDaily := []*domain.RepoDailyStatistics{
		{
			Repository:     repoBar,
			Date:           "2026-05-01",
			CommitCount:    1,
			PRCreated:      0,
			PRMerged:       0,
			IssueCount:     1,
			ReviewCount:    0,
			TotalAdditions: 70,
			TotalDeletions: 25,
		},
		{
			Repository:     repoBar,
			Date:           "2026-05-02",
			CommitCount:    0,
			PRCreated:      1,
			PRMerged:       1,
			IssueCount:     0,
			ReviewCount:    0,
			TotalAdditions: 6,
			TotalDeletions: 2,
		},
		{
			Repository:     repoFoo,
			Date:           "2026-05-01",
			CommitCount:    1,
			PRCreated:      2,
			PRMerged:       1,
			IssueCount:     0,
			ReviewCount:    0,
			TotalAdditions: 70,
			TotalDeletions: 27,
		},
		{
			Repository:     repoFoo,
			Date:           "2026-05-02",
			CommitCount:    1,
			PRCreated:      0,
			PRMerged:       0,
			IssueCount:     0,
			ReviewCount:    1,
			TotalAdditions: 30,
			TotalDeletions: 10,
		},
	}

	assert.Equal(t, wantRepoDaily, stats.RepoDailyStats, "RepoDailyStats must be bucketed per (repository, date) and sorted by (Repository, Date)")

	ownerByRepo := make(map[string]string, len(stats.AllRepositories))
	ownerTypeByRepo := make(map[string]string, len(stats.AllRepositories))
	for _, repo := range stats.AllRepositories {
		ownerByRepo[repo.Repository] = repo.Owner
		ownerTypeByRepo[repo.Repository] = repo.OwnerType
	}

	assert.Equal(t, "Tattsum", ownerByRepo[repoFoo], "AllRepositories Owner must come from activity RepositoryOwner")
	assert.Equal(t, "Organization", ownerTypeByRepo[repoFoo], "AllRepositories OwnerType must come from activity RepositoryOwnerType")
	assert.Equal(t, "octocat", ownerByRepo[repoBar], "AllRepositories Owner must come from activity RepositoryOwner")
	assert.Equal(t, "User", ownerTypeByRepo[repoBar], "AllRepositories OwnerType must come from activity RepositoryOwnerType")
}
