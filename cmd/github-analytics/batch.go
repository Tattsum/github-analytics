package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Tattsum/github-analytics/application"
	"github.com/Tattsum/github-analytics/domain"
	"github.com/Tattsum/github-analytics/infrastructure"
	"github.com/Tattsum/github-analytics/infrastructure/snapshotdb"
)

var (
	// errMissingDatabaseURL is returned when DATABASE_URL is unset in batch mode.
	errMissingDatabaseURL = errors.New("DATABASE_URL environment variable is not set")
	// errNoMemberStatistics is returned when no member could be aggregated.
	errNoMemberStatistics = errors.New("no member statistics were computed; aborting snapshot write")
)

// runBatch fetches activity for the given users, aggregates per-member
// statistics, and writes exactly one snapshot to PostgreSQL. Fatal exit is kept
// at the top level so deferred cleanup runs before the process terminates.
func runBatch(users []string, includePrivate bool, token string) {
	if err := executeBatch(users, includePrivate, token); err != nil {
		log.Fatalf("batch: %v", err)
	}
}

// executeBatch performs the batch run and returns an error instead of exiting,
// so that the deferred context cancel and DB Close always run.
//
// DATABASE_URL must point at the target PostgreSQL instance. Migrations are run
// before writing so the batch is safe to run against a fresh database.
func executeBatch(users []string, includePrivate bool, token string) error {
	const timeoutMinutes = 30

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return errMissingDatabaseURL
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutMinutes*time.Minute)
	defer cancel()

	client, err := infrastructure.OpenPostgres(databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}

	defer func() {
		if cerr := client.Close(); cerr != nil {
			log.Printf("Failed to close PostgreSQL connection: %v", cerr)
		}
	}()

	if err := infrastructure.Migrate(ctx, client); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	members := computeMemberStatistics(ctx, users, includePrivate, token)
	if len(members) == 0 {
		return errNoMemberStatistics
	}

	snapshot := &application.Snapshot{
		CapturedAt: time.Now(),
		Members:    members,
	}

	writer := snapshotdb.NewSnapshotWriter(client)
	if err := writer.Save(ctx, snapshot); err != nil {
		return fmt.Errorf("failed to save snapshot: %w", err)
	}

	fmt.Printf("\n=== バッチ完了 ===\nスナップショットを保存しました（メンバー数: %d, captured_at: %s）\n",
		len(members), snapshot.CapturedAt.Format(time.RFC3339))

	return nil
}

// computeMemberStatistics fetches and aggregates statistics for each user
// sequentially. Per-user failures are logged and skipped so one unreachable
// account does not abort the whole snapshot.
func computeMemberStatistics(ctx context.Context, users []string, includePrivate bool, token string) []*domain.UserStatistics {
	client := infrastructure.NewGitHubClient(token)
	repo := infrastructure.NewGitHubRepository(client)
	fetcher := infrastructure.NewGitHubDataFetcher(repo)
	statsService := application.NewStatisticsService()

	members := make([]*domain.UserStatistics, 0, len(users))

	for _, user := range users {
		stats, err := processUser(ctx, user, includePrivate, fetcher, statsService)
		if err != nil {
			log.Printf("Error processing user %s: %v", user, err)
			continue
		}

		members = append(members, stats)

		fmt.Printf("Completed processing user: %s\n", user)
	}

	return members
}
