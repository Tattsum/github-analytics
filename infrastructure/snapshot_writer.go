package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Tattsum/github-analytics/domain"
	"github.com/Tattsum/github-analytics/infrastructure/ent"
)

// ErrNilSnapshot is returned when Save is called with a nil snapshot.
var ErrNilSnapshot = errors.New("snapshot must not be nil")

// SnapshotInput is the aggregated data the writer persists for one batch run.
//
// It is defined here, in terms of domain types only, so the infrastructure
// write side does not depend on the application package (application already
// depends on infrastructure for UserActivityData, so the reverse import would
// create a cycle). The batch entrypoint adapts application.Snapshot into this.
type SnapshotInput struct {
	CapturedAt time.Time
	Members    []*domain.UserStatistics
}

// SnapshotWriter persists an aggregated snapshot into PostgreSQL via ent.
//
// Each Save creates exactly one new Snapshot row (idempotent at the snapshot
// level: a new run is a new snapshot), together with all of its MemberStat /
// MemberYearStat / MemberRepoStat rows, in a single transaction.
type SnapshotWriter struct {
	client *ent.Client
}

// NewSnapshotWriter constructs a SnapshotWriter backed by the given ent client.
func NewSnapshotWriter(client *ent.Client) *SnapshotWriter {
	return &SnapshotWriter{client: client}
}

// Save writes one aggregated snapshot and all of its member-level rows in a
// single transaction. On any failure the transaction is rolled back so a
// snapshot is never persisted partially.
func (w *SnapshotWriter) Save(ctx context.Context, snapshot *SnapshotInput) error {
	if snapshot == nil {
		return fmt.Errorf("save snapshot: %w", ErrNilSnapshot)
	}

	tx, err := w.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("begin snapshot transaction: %w", err)
	}

	if err := w.saveTx(ctx, tx, snapshot); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("%w (rollback failed: %v)", err, rbErr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit snapshot transaction: %w", err)
	}

	return nil
}

// saveTx performs the actual writes inside the given transaction.
func (w *SnapshotWriter) saveTx(ctx context.Context, tx *ent.Tx, snapshot *SnapshotInput) error {
	snapRow, err := tx.Snapshot.Create().
		SetCapturedAt(snapshot.CapturedAt).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("create snapshot row: %w", err)
	}

	memberStats, yearStats, repoStats := buildStatCreates(snapshot)

	if err := applyStatCreates(ctx, tx, snapRow.ID, memberStats, yearStats, repoStats); err != nil {
		return err
	}

	return nil
}

// memberStatInput captures the scalar fields of one MemberStat row.
type memberStatInput struct {
	login               string
	totalCommits        int
	totalPRCreated      int
	totalPRMerged       int
	totalIssues         int
	totalReviews        int
	totalAdditions      int
	totalDeletions      int
	firstActivityYear   int
	peakActivityYear    int
	peakActivityCommits int
	prToReviewRatio     float64
}

// memberYearStatInput captures the fields of one MemberYearStat row.
type memberYearStatInput struct {
	login       string
	year        int
	commitCount int
	prCreated   int
	prMerged    int
	issueCount  int
	reviewCount int
	additions   int
	deletions   int
}

// memberRepoStatInput captures the fields of one MemberRepoStat row.
type memberRepoStatInput struct {
	login         string
	nameWithOwner string
	commitCount   int
	prCreated     int
	prMerged      int
	issueCount    int
	reviewCount   int
	additions     int
	deletions     int
}

// buildStatCreates maps the aggregated per-member statistics into the flat
// row inputs persisted for a snapshot. It is pure (no I/O) so it can be unit
// tested without a database.
func buildStatCreates(snapshot *SnapshotInput) ([]memberStatInput, []memberYearStatInput, []memberRepoStatInput) {
	memberStats := make([]memberStatInput, 0, len(snapshot.Members))
	yearStats := make([]memberYearStatInput, 0)
	repoStats := make([]memberRepoStatInput, 0)

	for _, member := range snapshot.Members {
		if member == nil || member.User == nil {
			continue
		}

		login := member.User.Login

		memberStats = append(memberStats, memberStatInput{
			login:               login,
			totalCommits:        member.TotalCommits,
			totalPRCreated:      member.TotalPRCreated,
			totalPRMerged:       member.TotalPRMerged,
			totalIssues:         member.TotalIssues,
			totalReviews:        member.TotalReviews,
			totalAdditions:      member.TotalAdditions,
			totalDeletions:      member.TotalDeletions,
			firstActivityYear:   member.FirstActivityYear,
			peakActivityYear:    member.PeakActivityYear,
			peakActivityCommits: member.PeakActivityCommits,
			prToReviewRatio:     member.PRToReviewRatio,
		})

		yearStats = append(yearStats, buildYearStats(login, member.YearlyStats)...)
		repoStats = append(repoStats, buildRepoStats(login, member.AllRepositories)...)
	}

	return memberStats, yearStats, repoStats
}

// buildYearStats maps a member's yearly statistics into row inputs.
func buildYearStats(login string, yearly map[int]*domain.YearlyStatistics) []memberYearStatInput {
	out := make([]memberYearStatInput, 0, len(yearly))

	for year, stat := range yearly {
		if stat == nil {
			continue
		}

		out = append(out, memberYearStatInput{
			login:       login,
			year:        year,
			commitCount: stat.CommitCount,
			prCreated:   stat.PRCreated,
			prMerged:    stat.PRMerged,
			issueCount:  stat.IssueCount,
			reviewCount: stat.ReviewCount,
			additions:   stat.TotalAdditions,
			deletions:   stat.TotalDeletions,
		})
	}

	return out
}

// buildRepoStats maps a member's per-repository activity into row inputs.
// All repositories are persisted, not just the top ones.
func buildRepoStats(login string, repos []*domain.RepositoryActivity) []memberRepoStatInput {
	out := make([]memberRepoStatInput, 0, len(repos))

	for _, repo := range repos {
		if repo == nil {
			continue
		}

		out = append(out, memberRepoStatInput{
			login:         login,
			nameWithOwner: repo.Repository,
			commitCount:   repo.CommitCount,
			prCreated:     repo.PRCount,
			// PRMerged is not tracked per repository in RepositoryActivity, so it
			// stays at the column default (0) until v2.
			issueCount:  repo.IssueCount,
			reviewCount: repo.ReviewCount,
			additions:   repo.TotalAdditions,
			deletions:   repo.TotalDeletions,
		})
	}

	return out
}

// applyStatCreates persists the mapped row inputs against the given snapshot ID
// using ent bulk creates.
func applyStatCreates(
	ctx context.Context,
	tx *ent.Tx,
	snapshotID int,
	memberStats []memberStatInput,
	yearStats []memberYearStatInput,
	repoStats []memberRepoStatInput,
) error {
	if len(memberStats) > 0 {
		_, err := tx.MemberStat.MapCreateBulk(memberStats, func(c *ent.MemberStatCreate, i int) {
			m := memberStats[i]
			c.SetSnapshotID(snapshotID).
				SetLogin(m.login).
				SetTotalCommits(m.totalCommits).
				SetTotalPrCreated(m.totalPRCreated).
				SetTotalPrMerged(m.totalPRMerged).
				SetTotalIssues(m.totalIssues).
				SetTotalReviews(m.totalReviews).
				SetTotalAdditions(m.totalAdditions).
				SetTotalDeletions(m.totalDeletions).
				SetFirstActivityYear(m.firstActivityYear).
				SetPeakActivityYear(m.peakActivityYear).
				SetPeakActivityCommits(m.peakActivityCommits).
				SetPrToReviewRatio(m.prToReviewRatio)
		}).Save(ctx)
		if err != nil {
			return fmt.Errorf("create member stats: %w", err)
		}
	}

	if len(yearStats) > 0 {
		_, err := tx.MemberYearStat.MapCreateBulk(yearStats, func(c *ent.MemberYearStatCreate, i int) {
			y := yearStats[i]
			c.SetSnapshotID(snapshotID).
				SetLogin(y.login).
				SetYear(y.year).
				SetCommitCount(y.commitCount).
				SetPrCreated(y.prCreated).
				SetPrMerged(y.prMerged).
				SetIssueCount(y.issueCount).
				SetReviewCount(y.reviewCount).
				SetAdditions(y.additions).
				SetDeletions(y.deletions)
		}).Save(ctx)
		if err != nil {
			return fmt.Errorf("create member year stats: %w", err)
		}
	}

	if len(repoStats) > 0 {
		_, err := tx.MemberRepoStat.MapCreateBulk(repoStats, func(c *ent.MemberRepoStatCreate, i int) {
			r := repoStats[i]
			c.SetSnapshotID(snapshotID).
				SetLogin(r.login).
				SetNameWithOwner(r.nameWithOwner).
				SetCommitCount(r.commitCount).
				SetPrCreated(r.prCreated).
				SetPrMerged(r.prMerged).
				SetIssueCount(r.issueCount).
				SetReviewCount(r.reviewCount).
				SetAdditions(r.additions).
				SetDeletions(r.deletions)
		}).Save(ctx)
		if err != nil {
			return fmt.Errorf("create member repo stats: %w", err)
		}
	}

	return nil
}
