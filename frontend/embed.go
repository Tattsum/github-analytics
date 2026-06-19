// Package frontend embeds the built React SPA (frontend/dist) so the Go server
// can serve it same-origin alongside the GraphQL endpoint.
//
// The embed pattern requires frontend/dist to exist at build time. Run
// `pnpm --dir frontend build` (or `make frontend`) before building the server
// binary; the production Dockerfile does this for you.
package frontend

import "embed"

// DistFS holds the built SPA assets under the "dist" directory. Use Dist() to
// obtain the same files rooted at "dist" (so the FS looks like the web root).
//
//go:embed all:dist
var DistFS embed.FS
