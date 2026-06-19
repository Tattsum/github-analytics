//go:build tools

// Package tools pins build-time code generation tools as module dependencies so
// that `go run entgo.io/ent/cmd/ent` and `go run github.com/99designs/gqlgen`
// resolve to versions tracked in go.mod. This file is never compiled into the
// application binary (guarded by the `tools` build tag) and only exists to keep
// the tool imports in the module graph.
package tools

import (
	_ "entgo.io/ent/cmd/ent"
	_ "github.com/99designs/gqlgen"

	// pgx is the PostgreSQL driver used at runtime by infrastructure/ (via the
	// database/sql stdlib adapter). It is pinned here until that code imports it
	// directly, so `go mod tidy` keeps it in the module graph.
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)
