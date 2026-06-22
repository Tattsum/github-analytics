package application

import (
	"testing"

	"github.com/Tattsum/github-analytics/domain"
)

func TestSummarizeTeam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		members []*MemberStats
		want    *TeamSummary
	}{
		{
			name:    "no members yields zero totals and zero member count",
			members: []*MemberStats{},
			want: &TeamSummary{
				MemberCount:    0,
				TotalCommits:   0,
				TotalPRCreated: 0,
				TotalPRMerged:  0,
				TotalIssues:    0,
				TotalReviews:   0,
				TotalAdditions: 0,
				TotalDeletions: 0,
			},
		},
		{
			name: "single member totals equal that member",
			members: []*MemberStats{
				{
					Login:          "alice",
					TotalCommits:   42,
					TotalPRCreated: 7,
					TotalPRMerged:  5,
					TotalIssues:    3,
					TotalReviews:   11,
					TotalAdditions: 1200,
					TotalDeletions: 340,
				},
			},
			want: &TeamSummary{
				MemberCount:    1,
				TotalCommits:   42,
				TotalPRCreated: 7,
				TotalPRMerged:  5,
				TotalIssues:    3,
				TotalReviews:   11,
				TotalAdditions: 1200,
				TotalDeletions: 340,
			},
		},
		{
			name: "multiple members are summed across every metric",
			members: []*MemberStats{
				{
					Login:          "alice",
					TotalCommits:   42,
					TotalPRCreated: 7,
					TotalPRMerged:  5,
					TotalIssues:    3,
					TotalReviews:   11,
					TotalAdditions: 1200,
					TotalDeletions: 340,
				},
				{
					Login:          "bob",
					TotalCommits:   8,
					TotalPRCreated: 4,
					TotalPRMerged:  2,
					TotalIssues:    6,
					TotalReviews:   9,
					TotalAdditions: 300,
					TotalDeletions: 60,
				},
			},
			want: &TeamSummary{
				MemberCount:    2,
				TotalCommits:   50,
				TotalPRCreated: 11,
				TotalPRMerged:  7,
				TotalIssues:    9,
				TotalReviews:   20,
				TotalAdditions: 1500,
				TotalDeletions: 400,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := SummarizeTeam(tt.members)

			// RepositoryCount is populated by the caller from the repository
			// axis, so SummarizeTeam must leave it at zero.
			if got.RepositoryCount != 0 {
				t.Errorf("RepositoryCount = %d, want 0 (caller-populated)", got.RepositoryCount)
			}

			assertTeamSummary(t, got, tt.want)
		})
	}
}

func TestAggregateRepositories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		stats []*MemberRepoStat
		want  []*RepositoryStats
	}{
		{
			name:  "no stats yields empty result",
			stats: []*MemberRepoStat{},
			want:  []*RepositoryStats{},
		},
		{
			name: "single contributor in single repository",
			stats: []*MemberRepoStat{
				{
					Login:         "alice",
					NameWithOwner: "acme/api",
					CommitCount:   10,
					PRCreated:     4,
					PRMerged:      3,
					IssueCount:    2,
					ReviewCount:   5,
					Additions:     500,
					Deletions:     120,
				},
			},
			want: []*RepositoryStats{
				{
					NameWithOwner:    "acme/api",
					TotalCommits:     10,
					TotalPRCreated:   4,
					TotalPRMerged:    3,
					TotalIssues:      2,
					TotalReviews:     5,
					TotalAdditions:   500,
					TotalDeletions:   120,
					ContributorCount: 1,
					Contributors: []*RepositoryContributor{
						{
							Login:       "alice",
							CommitCount: 10,
							PRCreated:   4,
							ReviewCount: 5,
							Additions:   500,
							Deletions:   120,
						},
					},
				},
			},
		},
		{
			name: "multiple members in the same repository are summed and contributors sorted by login",
			stats: []*MemberRepoStat{
				{
					Login:         "bob",
					NameWithOwner: "acme/api",
					CommitCount:   6,
					PRCreated:     2,
					PRMerged:      1,
					IssueCount:    4,
					ReviewCount:   8,
					Additions:     200,
					Deletions:     30,
				},
				{
					Login:         "alice",
					NameWithOwner: "acme/api",
					CommitCount:   10,
					PRCreated:     4,
					PRMerged:      3,
					IssueCount:    2,
					ReviewCount:   5,
					Additions:     500,
					Deletions:     120,
				},
			},
			want: []*RepositoryStats{
				{
					NameWithOwner:    "acme/api",
					TotalCommits:     16,
					TotalPRCreated:   6,
					TotalPRMerged:    4,
					TotalIssues:      6,
					TotalReviews:     13,
					TotalAdditions:   700,
					TotalDeletions:   150,
					ContributorCount: 2,
					Contributors: []*RepositoryContributor{
						{
							Login:       "alice",
							CommitCount: 10,
							PRCreated:   4,
							ReviewCount: 5,
							Additions:   500,
							Deletions:   120,
						},
						{
							Login:       "bob",
							CommitCount: 6,
							PRCreated:   2,
							ReviewCount: 8,
							Additions:   200,
							Deletions:   30,
						},
					},
				},
			},
		},
		{
			name: "multiple repositories are returned sorted by nameWithOwner",
			stats: []*MemberRepoStat{
				{
					Login:         "alice",
					NameWithOwner: "acme/web",
					CommitCount:   3,
					PRCreated:     1,
					ReviewCount:   2,
					Additions:     90,
					Deletions:     10,
				},
				{
					Login:         "alice",
					NameWithOwner: "acme/api",
					CommitCount:   10,
					PRCreated:     4,
					ReviewCount:   5,
					Additions:     500,
					Deletions:     120,
				},
			},
			want: []*RepositoryStats{
				{
					NameWithOwner:    "acme/api",
					TotalCommits:     10,
					TotalPRCreated:   4,
					TotalReviews:     5,
					TotalAdditions:   500,
					TotalDeletions:   120,
					ContributorCount: 1,
					Contributors: []*RepositoryContributor{
						{Login: "alice", CommitCount: 10, PRCreated: 4, ReviewCount: 5, Additions: 500, Deletions: 120},
					},
				},
				{
					NameWithOwner:    "acme/web",
					TotalCommits:     3,
					TotalPRCreated:   1,
					TotalReviews:     2,
					TotalAdditions:   90,
					TotalDeletions:   10,
					ContributorCount: 1,
					Contributors: []*RepositoryContributor{
						{Login: "alice", CommitCount: 3, PRCreated: 1, ReviewCount: 2, Additions: 90, Deletions: 10},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := AggregateRepositories(tt.stats)

			assertRepositoryStats(t, got, tt.want)
		})
	}
}

func assertTeamSummary(t *testing.T, got, want *TeamSummary) {
	t.Helper()

	if got.MemberCount != want.MemberCount {
		t.Errorf("MemberCount = %d, want %d", got.MemberCount, want.MemberCount)
	}

	if got.TotalCommits != want.TotalCommits {
		t.Errorf("TotalCommits = %d, want %d", got.TotalCommits, want.TotalCommits)
	}

	if got.TotalPRCreated != want.TotalPRCreated {
		t.Errorf("TotalPRCreated = %d, want %d", got.TotalPRCreated, want.TotalPRCreated)
	}

	if got.TotalPRMerged != want.TotalPRMerged {
		t.Errorf("TotalPRMerged = %d, want %d", got.TotalPRMerged, want.TotalPRMerged)
	}

	if got.TotalIssues != want.TotalIssues {
		t.Errorf("TotalIssues = %d, want %d", got.TotalIssues, want.TotalIssues)
	}

	if got.TotalReviews != want.TotalReviews {
		t.Errorf("TotalReviews = %d, want %d", got.TotalReviews, want.TotalReviews)
	}

	if got.TotalAdditions != want.TotalAdditions {
		t.Errorf("TotalAdditions = %d, want %d", got.TotalAdditions, want.TotalAdditions)
	}

	if got.TotalDeletions != want.TotalDeletions {
		t.Errorf("TotalDeletions = %d, want %d", got.TotalDeletions, want.TotalDeletions)
	}
}

func assertRepositoryStats(t *testing.T, got, want []*RepositoryStats) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len(repositories) = %d, want %d", len(got), len(want))
	}

	for i := range want {
		gotRepo, wantRepo := got[i], want[i]

		if gotRepo.NameWithOwner != wantRepo.NameWithOwner {
			t.Errorf("repo[%d].NameWithOwner = %q, want %q", i, gotRepo.NameWithOwner, wantRepo.NameWithOwner)
		}

		if gotRepo.TotalCommits != wantRepo.TotalCommits {
			t.Errorf("repo[%d].TotalCommits = %d, want %d", i, gotRepo.TotalCommits, wantRepo.TotalCommits)
		}

		if gotRepo.TotalPRCreated != wantRepo.TotalPRCreated {
			t.Errorf("repo[%d].TotalPRCreated = %d, want %d", i, gotRepo.TotalPRCreated, wantRepo.TotalPRCreated)
		}

		if gotRepo.TotalPRMerged != wantRepo.TotalPRMerged {
			t.Errorf("repo[%d].TotalPRMerged = %d, want %d", i, gotRepo.TotalPRMerged, wantRepo.TotalPRMerged)
		}

		if gotRepo.TotalIssues != wantRepo.TotalIssues {
			t.Errorf("repo[%d].TotalIssues = %d, want %d", i, gotRepo.TotalIssues, wantRepo.TotalIssues)
		}

		if gotRepo.TotalReviews != wantRepo.TotalReviews {
			t.Errorf("repo[%d].TotalReviews = %d, want %d", i, gotRepo.TotalReviews, wantRepo.TotalReviews)
		}

		if gotRepo.TotalAdditions != wantRepo.TotalAdditions {
			t.Errorf("repo[%d].TotalAdditions = %d, want %d", i, gotRepo.TotalAdditions, wantRepo.TotalAdditions)
		}

		if gotRepo.TotalDeletions != wantRepo.TotalDeletions {
			t.Errorf("repo[%d].TotalDeletions = %d, want %d", i, gotRepo.TotalDeletions, wantRepo.TotalDeletions)
		}

		if gotRepo.ContributorCount != wantRepo.ContributorCount {
			t.Errorf("repo[%d].ContributorCount = %d, want %d", i, gotRepo.ContributorCount, wantRepo.ContributorCount)
		}

		assertContributors(t, i, gotRepo.Contributors, wantRepo.Contributors)
	}
}

func assertContributors(t *testing.T, repoIdx int, got, want []*RepositoryContributor) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("repo[%d] len(contributors) = %d, want %d", repoIdx, len(got), len(want))
	}

	for j := range want {
		gotC, wantC := got[j], want[j]

		if gotC.Login != wantC.Login {
			t.Errorf("repo[%d].contributor[%d].Login = %q, want %q", repoIdx, j, gotC.Login, wantC.Login)
		}

		if gotC.CommitCount != wantC.CommitCount {
			t.Errorf("repo[%d].contributor[%d].CommitCount = %d, want %d", repoIdx, j, gotC.CommitCount, wantC.CommitCount)
		}

		if gotC.PRCreated != wantC.PRCreated {
			t.Errorf("repo[%d].contributor[%d].PRCreated = %d, want %d", repoIdx, j, gotC.PRCreated, wantC.PRCreated)
		}

		if gotC.ReviewCount != wantC.ReviewCount {
			t.Errorf("repo[%d].contributor[%d].ReviewCount = %d, want %d", repoIdx, j, gotC.ReviewCount, wantC.ReviewCount)
		}

		if gotC.Additions != wantC.Additions {
			t.Errorf("repo[%d].contributor[%d].Additions = %d, want %d", repoIdx, j, gotC.Additions, wantC.Additions)
		}

		if gotC.Deletions != wantC.Deletions {
			t.Errorf("repo[%d].contributor[%d].Deletions = %d, want %d", repoIdx, j, gotC.Deletions, wantC.Deletions)
		}
	}
}

func TestAggregateTeamDaily(t *testing.T) {
	t.Parallel()

	rows := []*domain.DailyStatistics{
		{Date: "2024-01-09", CommitCount: 4, PRCreated: 1, ReviewCount: 2, TotalAdditions: 40},
		{Date: "2024-01-08", CommitCount: 3, PRCreated: 2, ReviewCount: 1, TotalAdditions: 30},
		{Date: "2024-01-08", CommitCount: 5, PRCreated: 1, ReviewCount: 4, TotalAdditions: 50},
	}

	got := AggregateTeamDaily(rows)

	if len(got) != 2 {
		t.Fatalf("AggregateTeamDaily() len = %d, want 2", len(got))
	}

	if got[0].Date != "2024-01-08" || got[1].Date != "2024-01-09" {
		t.Fatalf("AggregateTeamDaily() not sorted ascending by date: %q, %q", got[0].Date, got[1].Date)
	}

	if got[0].CommitCount != 8 || got[0].PRCreated != 3 || got[0].ReviewCount != 5 || got[0].TotalAdditions != 80 {
		t.Errorf("AggregateTeamDaily() 2024-01-08 = %+v, want commits 8 / prCreated 3 / reviews 5 / additions 80", got[0])
	}

	if got[1].CommitCount != 4 {
		t.Errorf("AggregateTeamDaily() 2024-01-09 commits = %d, want 4", got[1].CommitCount)
	}
}

func TestAggregateTeamDaily_Empty(t *testing.T) {
	t.Parallel()

	if got := AggregateTeamDaily(nil); len(got) != 0 {
		t.Errorf("AggregateTeamDaily(nil) = %+v, want empty", got)
	}
}
