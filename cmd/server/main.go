// Command server runs the GitHub team analytics web application: it serves the
// gqlgen GraphQL API at POST /query and the embedded React SPA at /.
//
// Required environment:
//
//	DATABASE_URL  PostgreSQL connection string (pgx/libpq format).
//
// Optional environment:
//
//	PORT          HTTP listen port (default 8080).
//	ENV           When "development"/"dev", the GraphQL playground is mounted at
//	              GET /playground. It is omitted in production.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/Tattsum/github-analytics/graph"
	"github.com/Tattsum/github-analytics/infrastructure"
	"github.com/Tattsum/github-analytics/infrastructure/snapshotdb"
)

const (
	defaultPort        = "8080"
	migrateTimeout     = 30 * time.Second
	readHeaderTimeout  = 10 * time.Second
	shutdownTimeout    = 10 * time.Second
	graphQLEndpoint    = "/query"
	playgroundEndpoint = "/playground"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("server: %v", err)
	}
}

// run wires the dependency graph and serves until interrupted. It returns an
// error instead of calling log.Fatal so that deferred cleanup (the DB Close)
// always runs.
func run() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return errors.New("DATABASE_URL environment variable is not set")
	}

	client, err := infrastructure.OpenPostgres(databaseURL)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			log.Printf("server: failed to close PostgreSQL connection: %v", cerr)
		}
	}()

	migrateCtx, cancel := context.WithTimeout(context.Background(), migrateTimeout)
	defer cancel()
	if err := infrastructure.Migrate(migrateCtx, client); err != nil {
		return err
	}

	reader := snapshotdb.NewSnapshotReader(client)
	resolver := graph.NewResolver(reader)

	mux := http.NewServeMux()
	mountGraphQL(mux, resolver)
	mountSPA(mux)

	srv := &http.Server{
		Addr:              ":" + port(),
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return serve(srv)
}

// mountGraphQL registers the gqlgen handler at POST /query and, in development,
// the GraphQL playground at GET /playground.
func mountGraphQL(mux *http.ServeMux, resolver *graph.Resolver) {
	gql := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	mux.Handle(graphQLEndpoint, gql)

	if isDevelopment() {
		mux.Handle(playgroundEndpoint, playground.Handler("GitHub Analytics", graphQLEndpoint))
		log.Printf("server: GraphQL playground mounted at GET %s", playgroundEndpoint)
	}
}

// serve starts the HTTP server and blocks until SIGINT/SIGTERM, then performs a
// graceful shutdown.
func serve(srv *http.Server) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		log.Printf("server: listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Println("server: shutdown signal received")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}

func port() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return defaultPort
}

func isDevelopment() bool {
	switch os.Getenv("ENV") {
	case "development", "dev":
		return true
	default:
		return false
	}
}
