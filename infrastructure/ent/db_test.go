package ent

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestOpenPostgres(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		dsn     string
		wantErr bool
	}{
		{
			name: "valid url dsn",
			dsn:  "postgres://user:pass@localhost:5432/db?sslmode=disable",
		},
		{
			name: "valid keyword dsn",
			dsn:  "host=localhost port=5432 user=u dbname=d sslmode=disable",
		},
		{
			// sql.Open is lazy, so even an empty DSN constructs a client; it
			// only fails when a query actually dials the database.
			name: "empty dsn still constructs lazily",
			dsn:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, err := OpenPostgres(tt.dsn)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("OpenPostgres(%q) = nil error, want error", tt.dsn)
				}
				return
			}
			if err != nil {
				t.Fatalf("OpenPostgres(%q) returned unexpected error: %v", tt.dsn, err)
			}
			if client == nil {
				t.Fatalf("OpenPostgres(%q) returned nil client", tt.dsn)
			}
			t.Cleanup(func() {
				if cerr := client.Close(); cerr != nil {
					t.Errorf("client.Close() error: %v", cerr)
				}
			})
		})
	}
}

// TestMigrate_UnreachableDBWrapsError verifies Migrate surfaces a wrapped error
// (rather than panicking or swallowing it) when the database cannot be reached.
func TestMigrate_UnreachableDBWrapsError(t *testing.T) {
	t.Parallel()

	// Point at a closed port on localhost so the dial fails fast.
	client, err := OpenPostgres("postgres://user:pass@127.0.0.1:1/db?sslmode=disable&connect_timeout=1")
	if err != nil {
		t.Fatalf("OpenPostgres returned unexpected error: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	err = Migrate(context.Background(), client)
	if err == nil {
		t.Fatal("Migrate against unreachable DB = nil error, want error")
	}
	if !strings.Contains(err.Error(), "run schema migration") {
		t.Errorf("Migrate error %q does not include expected context %q", err, "run schema migration")
	}
	// The wrapped underlying error must remain unwrappable via errors machinery.
	if errors.Unwrap(err) == nil {
		t.Errorf("Migrate error %q does not wrap an underlying error", err)
	}
}
