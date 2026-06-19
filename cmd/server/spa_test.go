package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

func newTestDist(t *testing.T) fstest.MapFS {
	t.Helper()
	return fstest.MapFS{
		"index.html":    {Data: []byte("<!doctype html><div id=root></div>")},
		"assets/app.js": {Data: []byte("console.log('app')")},
		"favicon.svg":   {Data: []byte("<svg/>")},
	}
}

func TestSPAHandler(t *testing.T) {
	t.Parallel()

	const indexBody = "<!doctype html><div id=root></div>"

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string // exact body for fallback/index cases; "" means don't assert body
	}{
		{
			name:       "root serves index",
			path:       "/",
			wantStatus: http.StatusOK,
			wantBody:   indexBody,
		},
		{
			name:       "existing asset served directly",
			path:       "/assets/app.js",
			wantStatus: http.StatusOK,
			wantBody:   "console.log('app')",
		},
		{
			name:       "existing top-level file served directly",
			path:       "/favicon.svg",
			wantStatus: http.StatusOK,
			wantBody:   "<svg/>",
		},
		{
			name:       "unknown client route falls back to index",
			path:       "/members/octocat",
			wantStatus: http.StatusOK,
			wantBody:   indexBody,
		},
		{
			name:       "unknown nested path falls back to index",
			path:       "/repositories/owner%2Frepo",
			wantStatus: http.StatusOK,
			wantBody:   indexBody,
		},
	}

	handler := spaHandler(newTestDist(t))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status: got %d, want %d", rec.Code, tt.wantStatus)
			}

			body, err := io.ReadAll(rec.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}
			if got := string(body); got != tt.wantBody {
				t.Fatalf("body: got %q, want %q", got, tt.wantBody)
			}
		})
	}
}

func TestTrimLeadingSlash(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "leading slash removed", in: "/assets/app.js", want: "assets/app.js"},
		{name: "root slash becomes empty", in: "/", want: ""},
		{name: "no leading slash unchanged", in: "assets/app.js", want: "assets/app.js"},
		{name: "empty unchanged", in: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := trimLeadingSlash(tt.in); got != tt.want {
				t.Fatalf("trimLeadingSlash(%q): got %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
