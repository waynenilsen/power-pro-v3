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

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Lift routes
	mux.HandleFunc("GET /lifts", liftHandler.List)
	mux.HandleFunc("GET /lifts/{id}", liftHandler.Get)
	mux.HandleFunc("GET /lifts/by-slug/{slug}", liftHandler.GetBySlug)
	mux.HandleFunc("POST /lifts", liftHandler.Create)
	mux.HandleFunc("PUT /lifts/{id}", liftHandler.Update)
	mux.HandleFunc("DELETE /lifts/{id}", liftHandler.Delete)

	// LiftMax routes
	mux.HandleFunc("GET /users/{userId}/lift-maxes/current", liftMaxHandler.GetCurrent)
	mux.HandleFunc("GET /users/{userId}/lift-maxes", liftMaxHandler.List)
	mux.HandleFunc("GET /lift-maxes/{id}/convert", liftMaxHandler.Convert)
	mux.HandleFunc("GET /lift-maxes/{id}", liftMaxHandler.Get)
	mux.HandleFunc("POST /users/{userId}/lift-maxes", liftMaxHandler.Create)
	mux.HandleFunc("PUT /lift-maxes/{id}", liftMaxHandler.Update)
	mux.HandleFunc("DELETE /lift-maxes/{id}", liftMaxHandler.Delete)
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
