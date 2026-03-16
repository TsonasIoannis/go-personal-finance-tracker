package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestRunHealthcheck(t *testing.T) {
	t.Run("returns nil when health endpoint responds with 200", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/health" {
				t.Fatalf("expected /health path, got %s", r.URL.Path)
			}

			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		previousPort, hadPort := os.LookupEnv("PORT")
		t.Cleanup(func() {
			if hadPort {
				_ = os.Setenv("PORT", previousPort)
				return
			}

			_ = os.Unsetenv("PORT")
		})

		parsedURL, err := url.Parse(server.URL)
		if err != nil {
			t.Fatalf("failed to parse test server url: %v", err)
		}

		_ = os.Setenv("PORT", parsedURL.Port())

		if err := runHealthcheck(); err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	})
}
