package infrastructure

import (
	"context"

	"github.com/Tattsum/github-analytics/infrastructure/ent"
)

// OpenPostgres opens a PostgreSQL-backed ent client using the pgx stdlib
// driver. dataSourceName is a standard pgx/libpq connection string (typically
// DATABASE_URL), e.g. "postgres://user:pass@host:5432/db?sslmode=disable".
//
// The returned client owns the underlying *sql.DB; the caller must close it via
// EntClient.Close once finished. This is a thin wrapper over the generated
// ent.OpenPostgres so that callers depend on the infrastructure package rather
// than reaching into the generated ent package directly.
func OpenPostgres(dataSourceName string) (*EntClient, error) {
	return ent.OpenPostgres(dataSourceName)
}

// Migrate runs ent's schema auto-migration against the connected database. It
// is idempotent and safe to call on every batch run before writing a snapshot.
func Migrate(ctx context.Context, client *EntClient) error {
	return ent.Migrate(ctx, client)
}

// EntClient is the persistence client used by the write side. It is an alias of
// the generated ent client and is also accepted by NewSnapshotWriter.
type EntClient = ent.Client
