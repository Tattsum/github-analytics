package main

import (
	"io/fs"
	"net/http"

	"github.com/Tattsum/github-analytics/frontend"
)

// mountSPA serves the embedded React build at "/". Static assets are served
// directly; any unknown path that is not an asset falls back to index.html so
// client-side routes (e.g. /members/:login) resolve in the browser.
//
// It is registered last (on the root pattern) so the more specific GraphQL
// routes take precedence in the ServeMux.
func mountSPA(mux *http.ServeMux) {
	dist, err := fs.Sub(frontend.DistFS, "dist")
	if err != nil {
		// fs.Sub on a static embed.FS only fails for a malformed path, which is
		// a programming error, so panicking at startup is appropriate.
		panic("server: invalid embedded frontend dist: " + err.Error())
	}

	mux.Handle("/", spaHandler(dist))
}

// spaHandler serves static assets from dist and falls back to index.html for
// any path that does not correspond to an existing file, so client-side routes
// resolve in the browser.
func spaHandler(dist fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(dist))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path != "/" {
			if _, statErr := fs.Stat(dist, trimLeadingSlash(path)); statErr != nil {
				serveIndex(w, dist)
				return
			}
		}
		fileServer.ServeHTTP(w, r)
	})
}

// serveIndex writes the SPA entry document (index.html) for client-routed paths.
func serveIndex(w http.ResponseWriter, dist fs.FS) {
	index, err := fs.ReadFile(dist, "index.html")
	if err != nil {
		http.Error(w, "frontend not built", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write(index); err != nil {
		// The client disconnected mid-write; nothing useful to do but log via
		// the default server error path is unavailable here, so we drop it.
		_ = err
	}
}

// trimLeadingSlash converts an HTTP path ("/assets/x.js") into an fs path
// ("assets/x.js"). fs.FS paths must not start with a slash.
func trimLeadingSlash(p string) string {
	if p != "" && p[0] == '/' {
		return p[1:]
	}
	return p
}
