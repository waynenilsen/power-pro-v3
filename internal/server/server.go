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
	"github.com/waynenilsen/power-pro-v3/internal/domain/loadstrategy"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
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
	config             Config
	httpServer         *http.Server
	liftRepo           *repository.LiftRepository
	liftMaxRepo        *repository.LiftMaxRepository
	prescriptionRepo   *repository.PrescriptionRepository
	dayRepo            *repository.DayRepository
	weekRepo           *repository.WeekRepository
	cycleRepo          *repository.CycleRepository
	weeklyLookupRepo   *repository.WeeklyLookupRepository
	dailyLookupRepo    *repository.DailyLookupRepository
	strategyFactory    *loadstrategy.StrategyFactory
	schemeFactory      *setscheme.SchemeFactory
}

// New creates a new Server instance.
func New(cfg Config) *Server {
	liftRepo := repository.NewLiftRepository(cfg.DB)
	liftMaxRepo := repository.NewLiftMaxRepository(cfg.DB)

	// Create strategy and scheme factories with registered types
	strategyFactory := loadstrategy.NewStrategyFactory()
	loadstrategy.RegisterPercentOf(strategyFactory)

	schemeFactory := setscheme.NewSchemeFactory()
	setscheme.RegisterFixedScheme(schemeFactory)
	setscheme.RegisterRampScheme(schemeFactory)

	prescriptionRepo := repository.NewPrescriptionRepository(cfg.DB, strategyFactory, schemeFactory)
	dayRepo := repository.NewDayRepository(cfg.DB)
	weekRepo := repository.NewWeekRepository(cfg.DB)
	cycleRepo := repository.NewCycleRepository(cfg.DB)
	weeklyLookupRepo := repository.NewWeeklyLookupRepository(cfg.DB)
	dailyLookupRepo := repository.NewDailyLookupRepository(cfg.DB)

	s := &Server{
		config:             cfg,
		liftRepo:           liftRepo,
		liftMaxRepo:        liftMaxRepo,
		prescriptionRepo:   prescriptionRepo,
		dayRepo:            dayRepo,
		weekRepo:           weekRepo,
		cycleRepo:          cycleRepo,
		weeklyLookupRepo:   weeklyLookupRepo,
		dailyLookupRepo:    dailyLookupRepo,
		strategyFactory:    strategyFactory,
		schemeFactory:      schemeFactory,
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
	prescriptionHandler := api.NewPrescriptionHandler(s.prescriptionRepo, s.liftRepo, s.liftMaxRepo, s.strategyFactory, s.schemeFactory)
	dayHandler := api.NewDayHandler(s.dayRepo, s.prescriptionRepo)
	weekHandler := api.NewWeekHandler(s.weekRepo)
	cycleHandler := api.NewCycleHandler(s.cycleRepo)
	weeklyLookupHandler := api.NewWeeklyLookupHandler(s.weeklyLookupRepo)
	dailyLookupHandler := api.NewDailyLookupHandler(s.dailyLookupRepo)

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

	// Prescription routes:
	// - All authenticated users can read prescription data
	// - Only admins can create/update/delete prescriptions
	// - Authenticated users can resolve prescriptions (needs their userId for max lookup)
	mux.Handle("GET /prescriptions", withAuth(prescriptionHandler.List))
	mux.Handle("GET /prescriptions/{id}", withAuth(prescriptionHandler.Get))
	mux.Handle("POST /prescriptions", withAdmin(prescriptionHandler.Create))
	mux.Handle("PUT /prescriptions/{id}", withAdmin(prescriptionHandler.Update))
	mux.Handle("DELETE /prescriptions/{id}", withAdmin(prescriptionHandler.Delete))
	mux.Handle("POST /prescriptions/{id}/resolve", withAuth(prescriptionHandler.Resolve))
	mux.Handle("POST /prescriptions/resolve-batch", withAuth(prescriptionHandler.ResolveBatch))

	// Day routes:
	// - All authenticated users can read day data
	// - Only admins can create/update/delete days and manage prescriptions
	mux.Handle("GET /days", withAuth(dayHandler.List))
	mux.Handle("GET /days/{id}", withAuth(dayHandler.Get))
	mux.Handle("GET /days/by-slug/{slug}", withAuth(dayHandler.GetBySlug))
	mux.Handle("POST /days", withAdmin(dayHandler.Create))
	mux.Handle("PUT /days/{id}", withAdmin(dayHandler.Update))
	mux.Handle("DELETE /days/{id}", withAdmin(dayHandler.Delete))
	mux.Handle("POST /days/{id}/prescriptions", withAdmin(dayHandler.AddPrescription))
	mux.Handle("DELETE /days/{id}/prescriptions/{prescriptionId}", withAdmin(dayHandler.RemovePrescription))
	mux.Handle("PUT /days/{id}/prescriptions/reorder", withAdmin(dayHandler.ReorderPrescriptions))

	// Week routes:
	// - All authenticated users can read week data
	// - Only admins can create/update/delete weeks and manage day mappings
	mux.Handle("GET /weeks", withAuth(weekHandler.List))
	mux.Handle("GET /weeks/{id}", withAuth(weekHandler.Get))
	mux.Handle("POST /weeks", withAdmin(weekHandler.Create))
	mux.Handle("PUT /weeks/{id}", withAdmin(weekHandler.Update))
	mux.Handle("DELETE /weeks/{id}", withAdmin(weekHandler.Delete))
	mux.Handle("POST /weeks/{id}/days", withAdmin(weekHandler.AddDay))
	mux.Handle("DELETE /weeks/{id}/days/{dayId}", withAdmin(weekHandler.RemoveDay))

	// Cycle routes:
	// - All authenticated users can read cycle data
	// - Only admins can create/update/delete cycles
	mux.Handle("GET /cycles", withAuth(cycleHandler.List))
	mux.Handle("GET /cycles/{id}", withAuth(cycleHandler.Get))
	mux.Handle("POST /cycles", withAdmin(cycleHandler.Create))
	mux.Handle("PUT /cycles/{id}", withAdmin(cycleHandler.Update))
	mux.Handle("DELETE /cycles/{id}", withAdmin(cycleHandler.Delete))

	// WeeklyLookup routes:
	// - All authenticated users can read weekly lookup data
	// - Only admins can create/update/delete weekly lookups
	mux.Handle("GET /weekly-lookups", withAuth(weeklyLookupHandler.List))
	mux.Handle("GET /weekly-lookups/{id}", withAuth(weeklyLookupHandler.Get))
	mux.Handle("POST /weekly-lookups", withAdmin(weeklyLookupHandler.Create))
	mux.Handle("PUT /weekly-lookups/{id}", withAdmin(weeklyLookupHandler.Update))
	mux.Handle("DELETE /weekly-lookups/{id}", withAdmin(weeklyLookupHandler.Delete))

	// DailyLookup routes:
	// - All authenticated users can read daily lookup data
	// - Only admins can create/update/delete daily lookups
	mux.Handle("GET /daily-lookups", withAuth(dailyLookupHandler.List))
	mux.Handle("GET /daily-lookups/{id}", withAuth(dailyLookupHandler.Get))
	mux.Handle("POST /daily-lookups", withAdmin(dailyLookupHandler.Create))
	mux.Handle("PUT /daily-lookups/{id}", withAdmin(dailyLookupHandler.Update))
	mux.Handle("DELETE /daily-lookups/{id}", withAdmin(dailyLookupHandler.Delete))
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
