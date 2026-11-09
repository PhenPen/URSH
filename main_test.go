package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthHandler checks the handler's own response: status, headers and body.
func TestHealthHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	healthHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if got, want := res.StatusCode, http.StatusOK; got != want {
		t.Errorf("status = %d, want %d", got, want)
	}
	if got, want := res.Header.Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Errorf("Content-Type = %q, want %q", got, want)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}
	if got, want := string(body), "ok\n"; got != want {
		t.Errorf("body = %q, want %q", got, want)
	}
}

// TestNewMux checks that routes are registered for the expected methods and
// that unregistered methods and paths are rejected by the router itself.
func TestNewMux(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		target     string
		wantStatus int
		wantAllow  string
	}{
		{name: "get health", method: http.MethodGet, target: "/health", wantStatus: http.StatusOK},
		{name: "head health", method: http.MethodHead, target: "/health", wantStatus: http.StatusOK},
		{name: "post health", method: http.MethodPost, target: "/health", wantStatus: http.StatusMethodNotAllowed, wantAllow: "GET, HEAD"},
		{name: "unknown path", method: http.MethodGet, target: "/missing", wantStatus: http.StatusNotFound},
	}

	mux := newMux()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.target, nil)

			mux.ServeHTTP(rec, req)

			if got := rec.Code; got != tt.wantStatus {
				t.Errorf("status = %d, want %d", got, tt.wantStatus)
			}
			if tt.wantAllow != "" {
				if got := rec.Header().Get("Allow"); got != tt.wantAllow {
					t.Errorf("Allow = %q, want %q", got, tt.wantAllow)
				}
			}
		})
	}
}
