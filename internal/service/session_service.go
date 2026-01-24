// Package service provides business logic services.
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/waynenilsen/power-pro-v3/internal/domain/loggedset"
	"github.com/waynenilsen/power-pro-v3/internal/domain/prescription"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
)

// ErrPrescriptionNotFound indicates the prescription was not found.
var ErrPrescriptionNotFound = errors.New("prescription not found")

// ErrNotVariableScheme indicates the prescription doesn't use a variable set scheme.
var ErrNotVariableScheme = errors.New("prescription does not use a variable set scheme")

// ErrNoSetsLogged indicates no sets have been logged yet for the prescription.
var ErrNoSetsLogged = errors.New("no sets logged for this prescription")

// PrescriptionRepository defines the interface for prescription lookup.
type PrescriptionRepository interface {
	GetByID(id string) (*prescription.Prescription, error)
}

// LoggedSetLister defines the interface for listing logged sets.
type LoggedSetLister interface {
	ListBySessionAndPrescription(sessionID, prescriptionID string) ([]loggedset.LoggedSet, error)
}

// SessionService provides business logic for workout session operations.
type SessionService struct {
	prescriptionRepo PrescriptionRepository
	loggedSetLister  LoggedSetLister
}

// NewSessionService creates a new SessionService.
func NewSessionService(prescriptionRepo PrescriptionRepository, loggedSetLister LoggedSetLister) *SessionService {
	return &SessionService{
		prescriptionRepo: prescriptionRepo,
		loggedSetLister:  loggedSetLister,
	}
}

// NextSetRequest contains the parameters for requesting the next set.
type NextSetRequest struct {
	SessionID      string
	PrescriptionID string
	UserID         string // For authorization check
}

// NextSetResult contains the result of a next set request.
type NextSetResult struct {
	// NextSet is the next set to perform. Nil if exercise is complete.
	NextSet *setscheme.GeneratedSet
	// IsComplete is true if the variable scheme is done (termination condition met).
	IsComplete bool
	// TotalSetsCompleted is the number of sets logged so far.
	TotalSetsCompleted int
	// TotalRepsCompleted is the total reps performed so far.
	TotalRepsCompleted int
	// TerminationReason explains why the exercise is complete (if applicable).
	TerminationReason string
}

// GetNextSet generates the next set for a variable scheme based on session performance.
// Returns the next set to perform, or indicates completion if termination conditions are met.
func (s *SessionService) GetNextSet(ctx context.Context, req NextSetRequest) (*NextSetResult, error) {
	// Get the prescription to access the set scheme
	presc, err := s.prescriptionRepo.GetByID(req.PrescriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get prescription: %w", err)
	}
	if presc == nil {
		return nil, ErrPrescriptionNotFound
	}

	// Check if the scheme is a variable scheme
	variableScheme, ok := presc.SetScheme.(setscheme.VariableSetScheme)
	if !ok || !variableScheme.IsVariableCount() {
		return nil, ErrNotVariableScheme
	}

	// Get logged sets for this session and prescription
	loggedSets, err := s.loggedSetLister.ListBySessionAndPrescription(req.SessionID, req.PrescriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get logged sets: %w", err)
	}

	// Calculate session stats from logged sets
	totalSets := len(loggedSets)
	totalReps := 0
	var lastReps int
	var lastRPE *float64
	var lastWeight float64

	for _, ls := range loggedSets {
		totalReps += ls.RepsPerformed
		lastReps = ls.RepsPerformed
		lastRPE = ls.RPE
		lastWeight = ls.Weight
	}

	// Build termination context from logged performance
	termCtx := setscheme.TerminationContext{
		SetNumber:  totalSets + 1, // Next set number
		LastRPE:    lastRPE,
		LastReps:   lastReps,
		TotalReps:  totalReps,
		TotalSets:  totalSets,
		TargetReps: 0, // Will be set by the scheme if needed
	}

	// Build history of generated sets from logged data
	history := make([]setscheme.GeneratedSet, len(loggedSets))
	for i, ls := range loggedSets {
		history[i] = setscheme.GeneratedSet{
			SetNumber:  ls.SetNumber,
			Weight:     ls.Weight,
			TargetReps: ls.TargetReps,
			IsWorkSet:  true,
		}
	}

	// If no sets logged yet, we need to return the first provisional set
	if totalSets == 0 {
		return nil, ErrNoSetsLogged
	}

	// Set target reps from the scheme's configuration for termination checking
	// This depends on the scheme type
	if mrs, ok := variableScheme.(*setscheme.MRS); ok {
		termCtx.TargetReps = mrs.MinRepsPerSet
	} else if fd, ok := variableScheme.(*setscheme.FatigueDrop); ok {
		termCtx.TargetReps = fd.TargetReps
	} else if tr, ok := variableScheme.(*setscheme.TotalRepsScheme); ok {
		termCtx.TargetReps = tr.SuggestedRepsPerSet
	}

	// Generate next set using the variable scheme
	genCtx := setscheme.DefaultSetGenerationContext()
	nextSet, shouldContinue := variableScheme.GenerateNextSet(genCtx, history, termCtx)

	result := &NextSetResult{
		TotalSetsCompleted: totalSets,
		TotalRepsCompleted: totalReps,
	}

	if !shouldContinue {
		result.IsComplete = true
		result.TerminationReason = determineTerminationReason(variableScheme, termCtx, lastWeight)
		return result, nil
	}

	result.NextSet = nextSet
	return result, nil
}

// determineTerminationReason provides a human-readable explanation for why the exercise ended.
func determineTerminationReason(scheme setscheme.VariableSetScheme, termCtx setscheme.TerminationContext, lastWeight float64) string {
	switch v := scheme.(type) {
	case *setscheme.MRS:
		if termCtx.TotalReps >= v.TargetTotalReps {
			return fmt.Sprintf("Target total reps reached (%d/%d)", termCtx.TotalReps, v.TargetTotalReps)
		}
		if termCtx.LastReps < v.MinRepsPerSet {
			return fmt.Sprintf("Failed to hit minimum reps (%d/%d)", termCtx.LastReps, v.MinRepsPerSet)
		}
		if termCtx.TotalSets >= v.MaxSets || (v.MaxSets == 0 && termCtx.TotalSets >= setscheme.DefaultMRSMaxSets) {
			return "Maximum sets reached (safety limit)"
		}
	case *setscheme.FatigueDrop:
		if termCtx.LastRPE != nil && *termCtx.LastRPE >= v.StopRPE {
			return fmt.Sprintf("Target RPE reached (%.1f/%.1f)", *termCtx.LastRPE, v.StopRPE)
		}
		if lastWeight <= 0 {
			return "Weight dropped to zero"
		}
		maxSets := v.MaxSets
		if maxSets == 0 {
			maxSets = setscheme.DefaultMaxSets
		}
		if termCtx.TotalSets >= maxSets {
			return "Maximum sets reached (safety limit)"
		}
	case *setscheme.TotalRepsScheme:
		if termCtx.TotalReps >= v.TargetTotalReps {
			return fmt.Sprintf("Target total reps reached (%d/%d)", termCtx.TotalReps, v.TargetTotalReps)
		}
		maxSets := v.MaxSets
		if maxSets == 0 {
			maxSets = setscheme.DefaultTotalRepsMaxSets
		}
		if termCtx.TotalSets >= maxSets {
			return "Maximum sets reached (safety limit)"
		}
	}
	return "Exercise complete"
}
