package snapshotdb

import (
	"reflect"
	"testing"

	"github.com/Tattsum/github-analytics/application"
	"github.com/Tattsum/github-analytics/domain"
)

func TestBuildRepoDayStats(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   func(t *testing.T) *application.Snapshot
		want []memberRepoDayStatInput
	}{
		{
			name: "single member maps every field with renamed columns",
			in: func(t *testing.T) *application.Snapshot {
				t.Helper()
				m := newMember(t, "Tattsum")
				m.RepoDailyStats = []*domain.RepoDailyStatistics{
					{
						Repository:     "Tattsum/foo",
						Date:           "2024-03-14",
						CommitCount:    42,
						PRCreated:      7,
						PRMerged:       6,
						IssueCount:     3,
						ReviewCount:    9,
						TotalAdditions: 1200,
						TotalDeletions: 350,
					},
				}
				return &application.Snapshot{Members: []*domain.UserStatistics{m}}
			},
			want: []memberRepoDayStatInput{
				{
					login:         "Tattsum",
					nameWithOwner: "Tattsum/foo",
					day:           "2024-03-14",
					commitCount:   42,
					prCreated:     7,
					prMerged:      6,
					issueCount:    3,
					reviewCount:   9,
					additions:     1200,
					deletions:     350,
				},
			},
		},
		{
			name: "nil member, nil user and nil stat entries are skipped",
			in: func(t *testing.T) *application.Snapshot {
				t.Helper()
				valid := newMember(t, "Tattsum")
				valid.RepoDailyStats = []*domain.RepoDailyStatistics{
					nil, // nil entry inside RepoDailyStats must be skipped
					{
						Repository:     "Tattsum/foo",
						Date:           "2024-03-14",
						CommitCount:    42,
						PRCreated:      7,
						PRMerged:       6,
						IssueCount:     3,
						ReviewCount:    9,
						TotalAdditions: 1200,
						TotalDeletions: 350,
					},
				}
				nilUser := &domain.UserStatistics{
					RepoDailyStats: []*domain.RepoDailyStatistics{
						{Repository: "Tattsum/ignored", Date: "2024-03-14", CommitCount: 99},
					},
				}
				return &application.Snapshot{Members: []*domain.UserStatistics{nil, nilUser, valid}}
			},
			want: []memberRepoDayStatInput{
				{
					login:         "Tattsum",
					nameWithOwner: "Tattsum/foo",
					day:           "2024-03-14",
					commitCount:   42,
					prCreated:     7,
					prMerged:      6,
					issueCount:    3,
					reviewCount:   9,
					additions:     1200,
					deletions:     350,
				},
			},
		},
		{
			name: "multiple members produce rows for each preserving order",
			in: func(t *testing.T) *application.Snapshot {
				t.Helper()
				first := newMember(t, "Tattsum")
				first.RepoDailyStats = []*domain.RepoDailyStatistics{
					{
						Repository:     "Tattsum/foo",
						Date:           "2024-03-14",
						CommitCount:    42,
						PRCreated:      7,
						PRMerged:       6,
						IssueCount:     3,
						ReviewCount:    9,
						TotalAdditions: 1200,
						TotalDeletions: 350,
					},
				}
				second := newMember(t, "octocat")
				second.RepoDailyStats = []*domain.RepoDailyStatistics{
					{
						Repository:     "octocat/bar",
						Date:           "2024-03-15",
						CommitCount:    21,
						PRCreated:      4,
						PRMerged:       2,
						IssueCount:     5,
						ReviewCount:    8,
						TotalAdditions: 640,
						TotalDeletions: 180,
					},
				}
				return &application.Snapshot{Members: []*domain.UserStatistics{first, second}}
			},
			want: []memberRepoDayStatInput{
				{
					login:         "Tattsum",
					nameWithOwner: "Tattsum/foo",
					day:           "2024-03-14",
					commitCount:   42,
					prCreated:     7,
					prMerged:      6,
					issueCount:    3,
					reviewCount:   9,
					additions:     1200,
					deletions:     350,
				},
				{
					login:         "octocat",
					nameWithOwner: "octocat/bar",
					day:           "2024-03-15",
					commitCount:   21,
					prCreated:     4,
					prMerged:      2,
					issueCount:    5,
					reviewCount:   8,
					additions:     640,
					deletions:     180,
				},
			},
		},
		{
			name: "no members yields empty slice",
			in: func(t *testing.T) *application.Snapshot {
				t.Helper()
				return &application.Snapshot{}
			},
			want: []memberRepoDayStatInput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := buildRepoDayStats(tt.in(t))

			if len(got) != len(tt.want) {
				t.Fatalf("repo day stat count = %d, want %d", len(got), len(tt.want))
			}

			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("repoDay[%d] = %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestBuildRepoMetas(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   func(t *testing.T) *application.Snapshot
		want []repoMetaInput
	}{
		{
			name: "owner is derived from the nameWithOwner prefix ignoring collected owner",
			in: func(t *testing.T) *application.Snapshot {
				t.Helper()
				m := newMember(t, "Tattsum")
				m.AllRepositories = []*domain.RepositoryActivity{
					// Owner deliberately differs from the prefix; the prefix must win.
					{Repository: "Tattsum/foo", Owner: "someone-else", OwnerType: "Organization"},
				}
				return &application.Snapshot{Members: []*domain.UserStatistics{m}}
			},
			want: []repoMetaInput{
				{nameWithOwner: "Tattsum/foo", owner: "Tattsum", ownerType: "Organization"},
			},
		},
		{
			name: "ownerType comes from the first member with a non-empty value",
			in: func(t *testing.T) *application.Snapshot {
				t.Helper()
				memberA := newMember(t, "Tattsum")
				memberA.AllRepositories = []*domain.RepositoryActivity{
					{Repository: "Tattsum/foo", Owner: "Tattsum", OwnerType: ""},
				}
				memberB := newMember(t, "octocat")
				memberB.AllRepositories = []*domain.RepositoryActivity{
					{Repository: "Tattsum/foo", Owner: "Tattsum", OwnerType: "Organization"},
				}
				return &application.Snapshot{Members: []*domain.UserStatistics{memberA, memberB}}
			},
			want: []repoMetaInput{
				{nameWithOwner: "Tattsum/foo", owner: "Tattsum", ownerType: "Organization"},
			},
		},
		{
			name: "output keeps first-seen order across members and repositories",
			in: func(t *testing.T) *application.Snapshot {
				t.Helper()
				memberA := newMember(t, "Tattsum")
				memberA.AllRepositories = []*domain.RepositoryActivity{
					{Repository: "Tattsum/foo", Owner: "Tattsum", OwnerType: "User"},
					{Repository: "Tattsum/bar", Owner: "Tattsum", OwnerType: "Organization"},
				}
				memberB := newMember(t, "octocat")
				memberB.AllRepositories = []*domain.RepositoryActivity{
					{Repository: "octocat/baz", Owner: "octocat", OwnerType: "User"},
				}
				return &application.Snapshot{Members: []*domain.UserStatistics{memberA, memberB}}
			},
			want: []repoMetaInput{
				{nameWithOwner: "Tattsum/foo", owner: "Tattsum", ownerType: "User"},
				{nameWithOwner: "Tattsum/bar", owner: "Tattsum", ownerType: "Organization"},
				{nameWithOwner: "octocat/baz", owner: "octocat", ownerType: "User"},
			},
		},
		{
			name: "a repo seen across multiple members yields exactly one row",
			in: func(t *testing.T) *application.Snapshot {
				t.Helper()
				memberA := newMember(t, "Tattsum")
				memberA.AllRepositories = []*domain.RepositoryActivity{
					{Repository: "Tattsum/foo", Owner: "Tattsum", OwnerType: "Organization"},
				}
				memberB := newMember(t, "octocat")
				memberB.AllRepositories = []*domain.RepositoryActivity{
					{Repository: "Tattsum/foo", Owner: "Tattsum", OwnerType: "Organization"},
				}
				return &application.Snapshot{Members: []*domain.UserStatistics{memberA, memberB}}
			},
			want: []repoMetaInput{
				{nameWithOwner: "Tattsum/foo", owner: "Tattsum", ownerType: "Organization"},
			},
		},
		{
			name: "empty repository strings and nil members are skipped",
			in: func(t *testing.T) *application.Snapshot {
				t.Helper()
				valid := newMember(t, "Tattsum")
				valid.AllRepositories = []*domain.RepositoryActivity{
					nil, // nil entry must be skipped
					{Repository: "", Owner: "Tattsum", OwnerType: "User"}, // empty repository must be skipped
					{Repository: "Tattsum/foo", Owner: "Tattsum", OwnerType: "Organization"},
				}
				return &application.Snapshot{Members: []*domain.UserStatistics{nil, valid}}
			},
			want: []repoMetaInput{
				{nameWithOwner: "Tattsum/foo", owner: "Tattsum", ownerType: "Organization"},
			},
		},
		{
			name: "no members yields empty slice",
			in: func(t *testing.T) *application.Snapshot {
				t.Helper()
				return &application.Snapshot{}
			},
			want: []repoMetaInput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := buildRepoMetas(tt.in(t))

			if len(got) != len(tt.want) {
				t.Fatalf("repo meta count = %d, want %d", len(got), len(tt.want))
			}

			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("repoMeta[%d] = %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestOwnerOf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		nameWithOwner  string
		collectedOwner string
		want           string
	}{
		{
			name:           "prefix wins even when collected owner differs",
			nameWithOwner:  "Tattsum/foo",
			collectedOwner: "someone-else",
			want:           "Tattsum",
		},
		{
			name:           "no slash falls back to collected owner",
			nameWithOwner:  "weirdname",
			collectedOwner: "Organization",
			want:           "Organization",
		},
		{
			name:           "empty owner prefix falls back to collected owner",
			nameWithOwner:  "/foo",
			collectedOwner: "Tattsum",
			want:           "Tattsum",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ownerOf(tt.nameWithOwner, tt.collectedOwner)

			if got != tt.want {
				t.Errorf("ownerOf(%q, %q) = %q, want %q", tt.nameWithOwner, tt.collectedOwner, got, tt.want)
			}
		})
	}
}

func TestChunkRows(t *testing.T) {
	t.Parallel()

	// seq builds [0, 1, ..., n-1] so chunk boundaries and contents are verifiable.
	seq := func(n int) []int {
		out := make([]int, n)
		for i := range out {
			out[i] = i
		}

		return out
	}

	tests := []struct {
		name      string
		n         int
		wantSizes []int
	}{
		{name: "empty yields no chunks", n: 0, wantSizes: []int{}},
		{name: "fewer than the cap is one chunk", n: 3, wantSizes: []int{3}},
		{name: "exactly the cap is one full chunk", n: maxBulkRows, wantSizes: []int{maxBulkRows}},
		{name: "one over the cap splits with a remainder", n: maxBulkRows + 1, wantSizes: []int{maxBulkRows, 1}},
		{name: "multiple caps plus remainder", n: 2*maxBulkRows + 7, wantSizes: []int{maxBulkRows, maxBulkRows, 7}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chunks := chunkRows(seq(tt.n))

			sizes := make([]int, len(chunks))
			for i, c := range chunks {
				sizes[i] = len(c)
			}
			if !reflect.DeepEqual(sizes, tt.wantSizes) {
				t.Fatalf("chunk sizes = %v, want %v", sizes, tt.wantSizes)
			}

			// Reassembling the chunks must reproduce the original order exactly,
			// so no row is dropped or duplicated across chunk boundaries.
			var flat []int
			for _, c := range chunks {
				flat = append(flat, c...)
			}
			if tt.n > 0 && !reflect.DeepEqual(flat, seq(tt.n)) {
				t.Fatalf("reassembled = %v, want %v", flat, seq(tt.n))
			}
		})
	}
}
