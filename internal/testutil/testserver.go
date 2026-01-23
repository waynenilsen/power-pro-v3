// Package testutil provides testing utilities for E2E tests.
package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

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

// NewTestServer creates and starts a new test server with an isolated database.
func NewTestServer() (*TestServer, error) {
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

	// Find available port
	port, err := server.FindAvailablePort()
	if err != nil {
		db.Close()
		os.Remove(dbPath)
		return nil, fmt.Errorf("failed to find available port: %w", err)
	}

	// Create server
	srv := server.New(server.Config{
		Port: port,
		DB:   db,
	})

	// Start server in goroutine
	errCh := make(chan error, 1)
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Wait for server to be ready
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	if err := waitForServer(baseURL, 5*time.Second); err != nil {
		db.Close()
		os.Remove(dbPath)
		return nil, fmt.Errorf("server failed to start: %w", err)
	}

	ts := &TestServer{
		Server:  srv,
		BaseURL: baseURL,
		port:    port,
		dbPath:  dbPath,
		db:      db,
		cleanups: []func(){
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

// waitForServer waits for the server to be ready by polling the health endpoint.
func waitForServer(baseURL string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 1 * time.Second}

	for time.Now().Before(deadline) {
		resp, err := client.Get(baseURL + "/health")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(50 * time.Millisecond)
	}

	return fmt.Errorf("server not ready after %v", timeout)
}

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
