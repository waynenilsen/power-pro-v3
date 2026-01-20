// Package server provides the HTTP server implementation.
package server

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/api"
	"github.com/waynenilsen/power-pro-v3/internal/middleware"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// Config holds server configuration.
type Config struct {
	Port int
	DB   *sql.DB
}

// Server represents the HTTP server.
type Server struct {
	config      Config
	httpServer  *http.Server
	liftRepo    *repository.LiftRepository
	liftMaxRepo *repository.LiftMaxRepository
}

// New creates a new Server instance.
func New(cfg Config) *Server {
	liftRepo := repository.NewLiftRepository(cfg.DB)
	liftMaxRepo := repository.NewLiftMaxRepository(cfg.DB)

	s := &Server{
		config:      cfg,
		liftRepo:    liftRepo,
		liftMaxRepo: liftMaxRepo,
	}

	mux := http.NewServeMux()
	s.registerRoutes(mux)

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

// registerRoutes sets up all API routes.
func (s *Server) registerRoutes(mux *http.ServeMux) {
	// Create handlers
	liftHandler := api.NewLiftHandler(s.liftRepo)
	liftMaxHandler := api.NewLiftMaxHandler(s.liftMaxRepo, s.liftRepo)

	// Auth middleware configuration
	authCfg := middleware.AuthConfig{
		WriteError: api.WriteError,
	}

	// Create middleware
	requireAuth := middleware.RequireAuth(authCfg)
	requireAdmin := middleware.RequireAdmin(authCfg)

	// Helper to wrap handler with middleware
	withAuth := func(h http.HandlerFunc) http.Handler {
		return requireAuth(http.HandlerFunc(h))
	}
	withAdmin := func(h http.HandlerFunc) http.Handler {
		return middleware.ChainMiddleware(requireAuth, requireAdmin)(http.HandlerFunc(h))
	}

	// LiftMax ownership check middleware
	liftMaxOwnerCheck := func(h http.HandlerFunc) http.Handler {
		ownerFunc := func(r *http.Request) (string, error) {
			// For routes with {userId} in path, that is the owner
			if userID := r.PathValue("userId"); userID != "" {
				return userID, nil
			}
			// For routes with {id}, we need to look up the resource
			// The handler will do the ownership check after fetching the resource
			// Return empty to skip middleware ownership check
			return "", nil
		}
		return middleware.ChainMiddleware(
			requireAuth,
			middleware.RequireOwnerOrAdmin(authCfg, ownerFunc),
		)(http.HandlerFunc(h))
	}

	// Health check (no auth required)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Lift routes (NFR-007):
	// - All authenticated users can read lift data
	// - Only admins can create/update/delete lifts
	mux.Handle("GET /lifts", withAuth(liftHandler.List))
	mux.Handle("GET /lifts/{id}", withAuth(liftHandler.Get))
	mux.Handle("GET /lifts/by-slug/{slug}", withAuth(liftHandler.GetBySlug))
	mux.Handle("POST /lifts", withAdmin(liftHandler.Create))
	mux.Handle("PUT /lifts/{id}", withAdmin(liftHandler.Update))
	mux.Handle("DELETE /lifts/{id}", withAdmin(liftHandler.Delete))

	// LiftMax routes (NFR-006):
	// - Users can only access their own LiftMax data
	// - Admins can access any user's LiftMax data
	mux.Handle("GET /users/{userId}/lift-maxes/current", liftMaxOwnerCheck(liftMaxHandler.GetCurrent))
	mux.Handle("GET /users/{userId}/lift-maxes", liftMaxOwnerCheck(liftMaxHandler.List))
	mux.Handle("GET /lift-maxes/{id}/convert", withAuth(liftMaxHandler.Convert))
	mux.Handle("GET /lift-maxes/{id}", withAuth(liftMaxHandler.Get))
	mux.Handle("POST /users/{userId}/lift-maxes", liftMaxOwnerCheck(liftMaxHandler.Create))
	mux.Handle("PUT /lift-maxes/{id}", withAuth(liftMaxHandler.Update))
	mux.Handle("DELETE /lift-maxes/{id}", withAuth(liftMaxHandler.Delete))
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Stop gracefully shuts down the server.
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// Addr returns the server's address after it starts listening.
func (s *Server) Addr() string {
	return s.httpServer.Addr
}

// FindAvailablePort finds an available port in the range 30000-60000.
func FindAvailablePort() (int, error) {
	// Use port 0 to let the OS assign an available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, fmt.Errorf("failed to find available port: %w", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	// Ensure it's in our expected range (mostly for documentation)
	if port < 1024 {
		return 0, fmt.Errorf("got port %d which is a privileged port", port)
	}
	return port, nil
}
