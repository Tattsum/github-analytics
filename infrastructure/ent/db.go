package ent

import (
	"context"
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib" // register the "pgx" database/sql driver
)

// OpenPostgres opens a PostgreSQL-backed ent client using the pgx stdlib
// driver.
//
// dataSourceName is a standard libpq/pgx connection string, e.g.
// "postgres://user:pass@host:5432/db?sslmode=disable" (typically DATABASE_URL).
// The returned client owns the underlying *sql.DB and must be closed by the
// caller via client.Close.
//
// The generated ent.Open requires the caller to pass a registered dialect
// driver name; OpenPostgres wraps that wiring so callers only supply a DSN.
func OpenPostgres(dataSourceName string) (*Client, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	drv := entsql.OpenDB(dialect.Postgres, db)
	return NewClient(Driver(drv)), nil
}

// Migrate runs ent's schema auto-migration against the connected database,
// creating or altering tables to match the current schema. It is idempotent and
// safe to call on every startup.
func Migrate(ctx context.Context, client *Client) error {
	if err := client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("run schema migration: %w", err)
	}
	return nil
}
