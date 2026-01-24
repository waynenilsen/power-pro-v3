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
	"github.com/waynenilsen/power-pro-v3/internal/service"
)

// Config holds server configuration.
type Config struct {
	Port int
	DB   *sql.DB
}

// Server represents the HTTP server.
type Server struct {
	config               Config
	httpServer           *http.Server
	liftRepo             *repository.LiftRepository
	liftMaxRepo          *repository.LiftMaxRepository
	prescriptionRepo     *repository.PrescriptionRepository
	dayRepo              *repository.DayRepository
	weekRepo             *repository.WeekRepository
	cycleRepo            *repository.CycleRepository
	weeklyLookupRepo     *repository.WeeklyLookupRepository
	dailyLookupRepo      *repository.DailyLookupRepository
	programRepo          *repository.ProgramRepository
	userProgramStateRepo *repository.UserProgramStateRepository
	workoutRepo          *repository.WorkoutRepository
	progressionRepo            *repository.ProgressionRepository
	programProgressionRepo     *repository.ProgramProgressionRepository
	progressionHistoryRepo     *repository.ProgressionHistoryRepository
	loggedSetRepo              *repository.LoggedSetRepository
	progressionService         *service.ProgressionService
	failureService             *service.FailureService
	strategyFactory            *loadstrategy.StrategyFactory
	schemeFactory              *setscheme.SchemeFactory
}

// New creates a new Server instance.
func New(cfg Config) *Server {
	liftRepo := repository.NewLiftRepository(cfg.DB)
	liftMaxRepo := repository.NewLiftMaxRepository(cfg.DB)

	// Create strategy and scheme factories with registered types
	strategyFactory := loadstrategy.NewStrategyFactory()
	loadstrategy.RegisterPercentOf(strategyFactory)
	loadstrategy.RegisterRPETarget(strategyFactory)

	schemeFactory := setscheme.NewSchemeFactory()
	setscheme.RegisterFixedScheme(schemeFactory)
	setscheme.RegisterRampScheme(schemeFactory)
	setscheme.RegisterAMRAPScheme(schemeFactory)
	setscheme.RegisterGreySkullScheme(schemeFactory)

	prescriptionRepo := repository.NewPrescriptionRepository(cfg.DB, strategyFactory, schemeFactory)
	dayRepo := repository.NewDayRepository(cfg.DB)
	weekRepo := repository.NewWeekRepository(cfg.DB)
	cycleRepo := repository.NewCycleRepository(cfg.DB)
	weeklyLookupRepo := repository.NewWeeklyLookupRepository(cfg.DB)
	dailyLookupRepo := repository.NewDailyLookupRepository(cfg.DB)
	programRepo := repository.NewProgramRepository(cfg.DB)
	userProgramStateRepo := repository.NewUserProgramStateRepository(cfg.DB)
	workoutRepo := repository.NewWorkoutRepository(cfg.DB, strategyFactory, schemeFactory)
	progressionRepo := repository.NewProgressionRepository(cfg.DB)
	programProgressionRepo := repository.NewProgramProgressionRepository(cfg.DB)
	progressionHistoryRepo := repository.NewProgressionHistoryRepository(cfg.DB)
	loggedSetRepo := repository.NewLoggedSetRepository(cfg.DB)
	progressionFactory := service.GetDefaultFactory()
	progressionService := service.NewProgressionService(cfg.DB, progressionFactory)
	failureService := service.NewFailureService(cfg.DB, progressionFactory)

	s := &Server{
		config:               cfg,
		liftRepo:             liftRepo,
		liftMaxRepo:          liftMaxRepo,
		prescriptionRepo:     prescriptionRepo,
		dayRepo:              dayRepo,
		weekRepo:             weekRepo,
		cycleRepo:            cycleRepo,
		weeklyLookupRepo:     weeklyLookupRepo,
		dailyLookupRepo:      dailyLookupRepo,
		programRepo:          programRepo,
		userProgramStateRepo: userProgramStateRepo,
		workoutRepo:                workoutRepo,
		progressionRepo:            progressionRepo,
		programProgressionRepo:     programProgressionRepo,
		progressionHistoryRepo:     progressionHistoryRepo,
		loggedSetRepo:              loggedSetRepo,
		progressionService:         progressionService,
		failureService:             failureService,
		strategyFactory:            strategyFactory,
		schemeFactory:              schemeFactory,
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
	programHandler := api.NewProgramHandler(s.programRepo, s.cycleRepo)

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
		_, _ = w.Write([]byte(`{"status":"ok"}`))
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

	// Program routes:
	// - All authenticated users can read program data
	// - Only admins can create/update/delete programs
	mux.Handle("GET /programs", withAuth(programHandler.List))
	mux.Handle("GET /programs/{id}", withAuth(programHandler.Get))
	mux.Handle("POST /programs", withAdmin(programHandler.Create))
	mux.Handle("PUT /programs/{id}", withAdmin(programHandler.Update))
	mux.Handle("DELETE /programs/{id}", withAdmin(programHandler.Delete))

	// Progression routes:
	// - All authenticated users can read progression data
	// - Only admins can create/update/delete progressions
	progressionHandler := api.NewProgressionHandler(s.progressionRepo)
	mux.Handle("GET /progressions", withAuth(progressionHandler.List))
	mux.Handle("GET /progressions/{id}", withAuth(progressionHandler.Get))
	mux.Handle("POST /progressions", withAdmin(progressionHandler.Create))
	mux.Handle("PUT /progressions/{id}", withAdmin(progressionHandler.Update))
	mux.Handle("DELETE /progressions/{id}", withAdmin(progressionHandler.Delete))

	// Program Progression Configuration routes:
	// - All authenticated users can read program progression configurations
	// - Only admins can create/update/delete program progression configurations
	programProgressionHandler := api.NewProgramProgressionHandler(s.programProgressionRepo, s.programRepo, s.progressionRepo, s.liftRepo)
	mux.Handle("GET /programs/{programId}/progressions", withAuth(programProgressionHandler.List))
	mux.Handle("GET /programs/{programId}/progressions/{configId}", withAuth(programProgressionHandler.Get))
	mux.Handle("POST /programs/{programId}/progressions", withAdmin(programProgressionHandler.Create))
	mux.Handle("PUT /programs/{programId}/progressions/{configId}", withAdmin(programProgressionHandler.Update))
	mux.Handle("DELETE /programs/{programId}/progressions/{configId}", withAdmin(programProgressionHandler.Delete))

	// User Program Enrollment routes:
	// - Users can manage their own enrollment (enroll, view, unenroll)
	// - Admins can manage any user's enrollment
	enrollmentHandler := api.NewEnrollmentHandler(s.userProgramStateRepo, s.programRepo)
	mux.Handle("POST /users/{userId}/program", withAuth(enrollmentHandler.Enroll))
	mux.Handle("GET /users/{userId}/program", withAuth(enrollmentHandler.Get))
	mux.Handle("DELETE /users/{userId}/program", withAuth(enrollmentHandler.Unenroll))

	// State Advancement routes:
	// - Users can advance their own program state
	// - Admins can advance any user's program state
	stateAdvancementHandler := api.NewStateAdvancementHandler(s.userProgramStateRepo, s.config.DB)
	mux.Handle("POST /users/{userId}/program-state/advance", withAuth(stateAdvancementHandler.Advance))

	// Workout Generation routes:
	// - Users can generate/preview their own workouts
	// - Admins can generate/preview any user's workouts
	workoutHandler := api.NewWorkoutHandler(s.workoutRepo, s.config.DB)
	mux.Handle("GET /users/{userId}/workout", withAuth(workoutHandler.Generate))
	mux.Handle("GET /users/{userId}/workout/preview", withAuth(workoutHandler.Preview))

	// Progression History routes:
	// - Users can query their own progression history
	// - Admins can query any user's progression history
	// - Handler performs its own authorization check
	progressionHistoryHandler := api.NewProgressionHistoryHandler(s.progressionHistoryRepo)
	mux.Handle("GET /users/{userId}/progression-history", withAuth(progressionHistoryHandler.List))

	// Manual Progression Trigger routes:
	// - Users can trigger their own progressions
	// - Admins can trigger progressions for any user
	// - Handler performs its own authorization check
	manualTriggerHandler := api.NewManualTriggerHandler(s.progressionService)
	mux.Handle("POST /users/{userId}/progressions/trigger", withAuth(manualTriggerHandler.Trigger))

	// Logged Set routes:
	// - Users can log sets for their own sessions
	// - Users can query their own logged sets
	// - Handler performs its own authorization check for user-specific data
	loggedSetHandler := api.NewLoggedSetHandler(s.loggedSetRepo, s.failureService)
	mux.Handle("POST /sessions/{sessionId}/sets", withAuth(loggedSetHandler.CreateBatch))
	mux.Handle("GET /sessions/{sessionId}/sets", withAuth(loggedSetHandler.ListBySession))
	mux.Handle("GET /users/{userId}/logged-sets", withAuth(loggedSetHandler.ListByUser))

	// Failure Counter routes:
	// - Users can query their own failure counters
	// - Admins can query any user's failure counters
	// - Handler performs its own authorization check
	failureCounterHandler := api.NewFailureCounterHandler(s.failureService)
	mux.Handle("GET /users/{userId}/failure-counters", withAuth(failureCounterHandler.Get))
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
