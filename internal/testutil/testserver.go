// Package testutil provides testing utilities for E2E tests.
package testutil

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/waynenilsen/power-pro-v3/internal/database"
	"github.com/waynenilsen/power-pro-v3/internal/server"
)

// TestServer wraps a server instance for testing.
type TestServer struct {
	Server   *server.Server
	BaseURL  string
	port     int
	dbPath   string
	db       *sql.DB
	cleanups []func()
}

type inProcessRoundTripper struct {
	baseHost string
	handler  http.Handler
	fallback http.RoundTripper
}

func (rt *inProcessRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req == nil || req.URL == nil {
		return nil, errors.New("nil request")
	}

	if !strings.EqualFold(req.URL.Host, rt.baseHost) {
		if rt.fallback != nil {
			return rt.fallback.RoundTrip(req)
		}
		return nil, fmt.Errorf("unexpected host %q (expected %q)", req.URL.Host, rt.baseHost)
	}

	if req.Body != nil {
		defer req.Body.Close()
	}

	r := req.Clone(req.Context())
	r.RequestURI = req.URL.RequestURI()
	if r.Host == "" {
		r.Host = rt.baseHost
	}
	if r.RemoteAddr == "" {
		r.RemoteAddr = "127.0.0.1:0"
	}

	rec := httptest.NewRecorder()
	rt.handler.ServeHTTP(rec, r)

	res := rec.Result()
	res.Request = req
	return res, nil
}

// NewTestServer creates and starts a new test server with an isolated database.
// It automatically enables test mode (POWERPRO_TEST_MODE=true) to allow X-User-ID
// and X-Admin headers to work for authentication in tests.
func NewTestServer() (*TestServer, error) {
	originalTestMode, hadOriginalTestMode := os.LookupEnv("POWERPRO_TEST_MODE")
	// Enable test mode for X-User-ID header authentication
	os.Setenv("POWERPRO_TEST_MODE", "true")

	// Find migrations path relative to project root
	migrationsPath, err := findMigrationsPath()
	if err != nil {
		return nil, fmt.Errorf("failed to find migrations path: %w", err)
	}

	// Create temporary database
	tmpFile, err := os.CreateTemp("", "powerpro-test-*.db")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp db file: %w", err)
	}
	dbPath := tmpFile.Name()
	tmpFile.Close()

	// Open database with migrations
	db, err := database.Open(database.Config{
		Path:           dbPath,
		MigrationsPath: migrationsPath,
	})
	if err != nil {
		os.Remove(dbPath)
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create server without binding to an actual TCP port.
	//
	// The test suite routes HTTP requests directly into the handler via a custom
	// http.RoundTripper. This makes tests work in sandboxed environments where
	// binding/listening on TCP ports is not allowed.
	srv := server.New(server.Config{
		Port: 0,
		DB:   db,
	})

	baseURL := "http://powerpro.test"
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		db.Close()
		os.Remove(dbPath)
		return nil, fmt.Errorf("failed to parse base url: %w", err)
	}

	originalTransport := http.DefaultClient.Transport
	fallbackTransport := originalTransport
	if fallbackTransport == nil {
		fallbackTransport = http.DefaultTransport
	}

	http.DefaultClient.Transport = &inProcessRoundTripper{
		baseHost: parsedBaseURL.Host,
		handler:  srv.Handler(),
		fallback: fallbackTransport,
	}

	// Basic readiness check against /health.
	req, err := http.NewRequest(http.MethodGet, baseURL+"/health", nil)
	if err != nil {
		http.DefaultClient.Transport = originalTransport
		db.Close()
		os.Remove(dbPath)
		return nil, fmt.Errorf("failed to create health request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.DefaultClient.Transport = originalTransport
		db.Close()
		os.Remove(dbPath)
		return nil, fmt.Errorf("server failed health check: %w", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		http.DefaultClient.Transport = originalTransport
		db.Close()
		os.Remove(dbPath)
		return nil, fmt.Errorf("server health check failed: status %d", resp.StatusCode)
	}

	ts := &TestServer{
		Server:  srv,
		BaseURL: baseURL,
		port:    0,
		dbPath:  dbPath,
		db:      db,
		cleanups: []func(){
			func() { http.DefaultClient.Transport = originalTransport },
			func() {
				if hadOriginalTestMode {
					_ = os.Setenv("POWERPRO_TEST_MODE", originalTestMode)
					return
				}
				_ = os.Unsetenv("POWERPRO_TEST_MODE")
			},
			func() { _ = srv.Stop(context.Background()) },
			func() { db.Close() },
			func() { _ = os.Remove(dbPath) },
		},
	}

	return ts, nil
}

// Close shuts down the test server and cleans up resources.
func (ts *TestServer) Close() {
	for _, cleanup := range ts.cleanups {
		cleanup()
	}
}

// URL returns a full URL for the given path.
func (ts *TestServer) URL(path string) string {
	return ts.BaseURL + path
}

// DB returns the underlying database connection for direct access in tests.
func (ts *TestServer) DB() *sql.DB {
	return ts.db
}

// AuthHeaders returns HTTP headers for an authenticated user.
func AuthHeaders(userID string, isAdmin bool) map[string]string {
	headers := map[string]string{
		"X-User-ID": userID,
	}
	if isAdmin {
		headers["X-Admin"] = "true"
	}
	return headers
}

// TestUserID is the standard test user ID used in tests.
const TestUserID = "test-user-001"

// TestAdminID is the standard test admin ID used in tests.
const TestAdminID = "test-admin-001"

// findMigrationsPath finds the migrations directory path.
func findMigrationsPath() (string, error) {
	// Get the directory of this source file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get caller info")
	}

	// Navigate up to project root
	dir := filepath.Dir(filename)
	for i := 0; i < 5; i++ {
		migrationsPath := filepath.Join(dir, "migrations")
		if _, err := os.Stat(migrationsPath); err == nil {
			return migrationsPath, nil
		}
		dir = filepath.Dir(dir)
	}

	return "", fmt.Errorf("migrations directory not found")
}
