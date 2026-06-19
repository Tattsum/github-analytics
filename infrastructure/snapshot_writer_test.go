package infrastructure

import (
	"sort"
	"testing"

	"github.com/Tattsum/github-analytics/domain"
)

// newMember builds a UserStatistics fixture for the mapping tests.
func newMember(t *testing.T, login string) *domain.UserStatistics {
	t.Helper()

	return domain.NewUserStatistics(domain.NewUser(login, login+" Name", "2020-01-01T00:00:00Z"))
}

func TestBuildStatCreates_MemberScalars(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   func(t *testing.T) *SnapshotInput
		want []memberStatInput
	}{
		{
			name: "single member maps every scalar field",
			in: func(t *testing.T) *SnapshotInput {
				t.Helper()
				m := newMember(t, "octocat")
				m.TotalCommits = 12
				m.TotalPRCreated = 7
				m.TotalPRMerged = 5
				m.TotalIssues = 3
				m.TotalReviews = 9
				m.TotalAdditions = 1000
				m.TotalDeletions = 400
				m.FirstActivityYear = 2018
				m.PeakActivityYear = 2021
				m.PeakActivityCommits = 8
				m.PRToReviewRatio = 1.2857
				return &SnapshotInput{Members: []*domain.UserStatistics{m}}
			},
			want: []memberStatInput{
				{
					login: "octocat", totalCommits: 12, totalPRCreated: 7, totalPRMerged: 5,
					totalIssues: 3, totalReviews: 9, totalAdditions: 1000, totalDeletions: 400,
					firstActivityYear: 2018, peakActivityYear: 2021, peakActivityCommits: 8,
					prToReviewRatio: 1.2857,
				},
			},
		},
		{
			name: "nil member and nil user are skipped",
			in: func(t *testing.T) *SnapshotInput {
				t.Helper()
				valid := newMember(t, "valid")
				valid.TotalCommits = 1
				nilUser := &domain.UserStatistics{TotalCommits: 99}
				return &SnapshotInput{Members: []*domain.UserStatistics{nil, nilUser, valid}}
			},
			want: []memberStatInput{
				{login: "valid", totalCommits: 1},
			},
		},
		{
			name: "no members yields empty slice",
			in: func(t *testing.T) *SnapshotInput {
				t.Helper()
				return &SnapshotInput{}
			},
			want: []memberStatInput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, _, _ := buildStatCreates(tt.in(t))

			if len(got) != len(tt.want) {
				t.Fatalf("member stat count = %d, want %d", len(got), len(tt.want))
			}

			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("member[%d] = %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestBuildStatCreates_YearStats(t *testing.T) {
	t.Parallel()

	m := newMember(t, "dev")
	m.YearlyStats = map[int]*domain.YearlyStatistics{
		2020: {Year: 2020, CommitCount: 10, PRCreated: 4, PRMerged: 3, IssueCount: 2, ReviewCount: 6, TotalAdditions: 500, TotalDeletions: 200},
		2021: {Year: 2021, CommitCount: 20, PRCreated: 8, PRMerged: 7, IssueCount: 1, ReviewCount: 12, TotalAdditions: 800, TotalDeletions: 300},
	}
	// A nil yearly entry must be skipped without panicking.
	m.YearlyStats[2019] = nil

	_, yearStats, _ := buildStatCreates(&SnapshotInput{Members: []*domain.UserStatistics{m}})

	if len(yearStats) != 2 {
		t.Fatalf("year stat count = %d, want 2", len(yearStats))
	}

	// Map iteration order is non-deterministic; sort by year for comparison.
	sort.Slice(yearStats, func(i, j int) bool { return yearStats[i].year < yearStats[j].year })

	want := []memberYearStatInput{
		{login: "dev", year: 2020, commitCount: 10, prCreated: 4, prMerged: 3, issueCount: 2, reviewCount: 6, additions: 500, deletions: 200},
		{login: "dev", year: 2021, commitCount: 20, prCreated: 8, prMerged: 7, issueCount: 1, reviewCount: 12, additions: 800, deletions: 300},
	}

	for i := range want {
		if yearStats[i] != want[i] {
			t.Errorf("year[%d] = %+v, want %+v", i, yearStats[i], want[i])
		}
	}
}

func TestBuildStatCreates_RepoStats(t *testing.T) {
	t.Parallel()

	m := newMember(t, "dev")
	m.AllRepositories = []*domain.RepositoryActivity{
		{Repository: "org/repo-a", CommitCount: 5, PRCount: 3, IssueCount: 1, ReviewCount: 4, TotalAdditions: 100, TotalDeletions: 40},
		nil, // must be skipped
		{Repository: "org/repo-b", CommitCount: 2, PRCount: 1, IssueCount: 0, ReviewCount: 0, TotalAdditions: 10, TotalDeletions: 5},
	}

	_, _, repoStats := buildStatCreates(&SnapshotInput{Members: []*domain.UserStatistics{m}})

	if len(repoStats) != 2 {
		t.Fatalf("repo stat count = %d, want 2", len(repoStats))
	}

	want := []memberRepoStatInput{
		{login: "dev", nameWithOwner: "org/repo-a", commitCount: 5, prCreated: 3, issueCount: 1, reviewCount: 4, additions: 100, deletions: 40},
		{login: "dev", nameWithOwner: "org/repo-b", commitCount: 2, prCreated: 1, additions: 10, deletions: 5},
	}

	for i := range want {
		if repoStats[i] != want[i] {
			t.Errorf("repo[%d] = %+v, want %+v", i, repoStats[i], want[i])
		}
	}
}

func TestBuildStatCreates_PersistsAllRepositories(t *testing.T) {
	t.Parallel()

	// AllRepositories must be persisted in full, not just the top 3.
	m := newMember(t, "prolific")
	const repoCount = 7
	m.AllRepositories = make([]*domain.RepositoryActivity, 0, repoCount)
	for i := 0; i < repoCount; i++ {
		m.AllRepositories = append(m.AllRepositories, &domain.RepositoryActivity{
			Repository:  "org/repo-" + string(rune('a'+i)),
			CommitCount: repoCount - i,
		})
	}

	_, _, repoStats := buildStatCreates(&SnapshotInput{Members: []*domain.UserStatistics{m}})

	if len(repoStats) != repoCount {
		t.Fatalf("repo stat count = %d, want %d (all repositories must be stored)", len(repoStats), repoCount)
	}
}
