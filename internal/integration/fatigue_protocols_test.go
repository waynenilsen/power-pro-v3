// Package integration provides end-to-end integration tests for fatigue protocols.
// This file tests MRS (GZCL-style) and FatigueDrop (RTS-style) set schemes
// demonstrating complete workout flows including termination conditions.
package integration

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/loggedset"
	"github.com/waynenilsen/power-pro-v3/internal/domain/prescription"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
	"github.com/waynenilsen/power-pro-v3/internal/service"
)

// =============================================================================
// MOCK REPOSITORIES FOR INTEGRATION TESTS
// =============================================================================

// mockPrescriptionRepo implements service.PrescriptionRepository for testing.
type mockPrescriptionRepo struct {
	prescriptions map[string]*prescription.Prescription
}

func newMockPrescriptionRepo() *mockPrescriptionRepo {
	return &mockPrescriptionRepo{
		prescriptions: make(map[string]*prescription.Prescription),
	}
}

func (m *mockPrescriptionRepo) Add(p *prescription.Prescription) {
	m.prescriptions[p.ID] = p
}

func (m *mockPrescriptionRepo) GetByID(id string) (*prescription.Prescription, error) {
	return m.prescriptions[id], nil
}

// mockLoggedSetLister implements service.LoggedSetLister for testing.
type mockLoggedSetLister struct {
	sets []loggedset.LoggedSet
}

func newMockLoggedSetLister() *mockLoggedSetLister {
	return &mockLoggedSetLister{
		sets: make([]loggedset.LoggedSet, 0),
	}
}

func (m *mockLoggedSetLister) Add(ls loggedset.LoggedSet) {
	m.sets = append(m.sets, ls)
}

func (m *mockLoggedSetLister) Clear() {
	m.sets = make([]loggedset.LoggedSet, 0)
}

func (m *mockLoggedSetLister) ListBySessionAndPrescription(sessionID, prescriptionID string) ([]loggedset.LoggedSet, error) {
	var result []loggedset.LoggedSet
	for _, ls := range m.sets {
		if ls.SessionID == sessionID && ls.PrescriptionID == prescriptionID {
			result = append(result, ls)
		}
	}
	return result, nil
}

// =============================================================================
// GZCL VDIP (MRS) INTEGRATION TESTS
// =============================================================================

// TestGZCL_VDIP_T1_ThreeMRS simulates a GZCL T1 workout using 3 MRS blocks.
//
// GZCL T1 Protocol:
// - 3 Max Rep Sets (MRS)
// - Target: 25 total reps
// - Minimum: 3 reps per set (technical failure threshold)
// - Same weight for all sets
//
// Example session:
// - Set 1: 225 lbs x 10 (total: 10)
// - Set 2: 225 lbs x 8 (total: 18)
// - Set 3: 225 lbs x 6 (total: 24)
// - Set 4: 225 lbs x 4 (total: 28) → STOP (exceeded target 25)
func TestGZCL_VDIP_T1_ThreeMRS(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	// Setup test IDs
	userID := uuid.New().String()
	sessionID := uuid.New().String()
	prescriptionID := uuid.New().String()
	liftID := uuid.New().String()

	// Create MRS scheme for T1: 3 MRS, target 25 total reps, minimum 3 per set
	mrsScheme, err := setscheme.NewMRS(25, 3, 10, 3) // target=25, min=3, max=10, numMRS=3
	if err != nil {
		t.Fatalf("failed to create MRS scheme: %v", err)
	}

	// Create prescription with MRS scheme
	presc := &prescription.Prescription{
		ID:        prescriptionID,
		LiftID:    liftID,
		SetScheme: mrsScheme,
	}

	// Setup mock repositories
	prescRepo := newMockPrescriptionRepo()
	prescRepo.Add(presc)
	loggedSetLister := newMockLoggedSetLister()

	// Create session service
	sessionService := service.NewSessionService(prescRepo, loggedSetLister)

	// === Workout Simulation ===
	baseWeight := 225.0

	// Generate first set
	genCtx := setscheme.DefaultSetGenerationContext()
	initialSets, err := mrsScheme.GenerateSets(baseWeight, genCtx)
	if err != nil {
		t.Fatalf("failed to generate initial sets: %v", err)
	}
	if len(initialSets) != 1 {
		t.Fatalf("expected 1 initial set, got %d", len(initialSets))
	}

	// Verify first set is provisional
	if !initialSets[0].IsProvisional {
		t.Error("expected first set to be provisional")
	}

	// Simulate workout sets (realistic T1 performance)
	workoutSets := []struct {
		repsPerformed int
		description   string
	}{
		{10, "Fresh, strong set"},
		{8, "Still strong, slight fatigue"},
		{6, "Noticeable fatigue"},
		{4, "Final set, significant fatigue"},
	}

	totalReps := 0
	for i, setData := range workoutSets {
		setNumber := i + 1
		totalReps += setData.repsPerformed

		// Create logged set
		input := loggedset.CreateLoggedSetInput{
			UserID:         userID,
			SessionID:      sessionID,
			PrescriptionID: prescriptionID,
			LiftID:         liftID,
			SetNumber:      setNumber,
			Weight:         baseWeight,
			TargetReps:     3, // MinRepsPerSet
			RepsPerformed:  setData.repsPerformed,
			IsAMRAP:        false,
		}

		logged, valResult := loggedset.NewLoggedSet(input, uuid.New().String())
		if !valResult.Valid {
			t.Fatalf("set %d: failed to create logged set: %v", setNumber, valResult.Errors)
		}
		loggedSetLister.Add(*logged)

		// Check if we should continue (except after final set)
		if i < len(workoutSets)-1 {
			req := service.NextSetRequest{
				SessionID:      sessionID,
				PrescriptionID: prescriptionID,
				UserID:         userID,
			}
			result, err := sessionService.GetNextSet(ctx, req)
			if err != nil {
				t.Fatalf("set %d: GetNextSet failed: %v", setNumber, err)
			}

			if result.IsComplete {
				// Check if we hit the target
				if totalReps >= 25 {
					// Expected termination
					if result.NextSet != nil {
						t.Errorf("set %d: expected no next set after completion", setNumber)
					}
					break
				}
				t.Errorf("set %d: unexpected completion at %d total reps", setNumber, totalReps)
				break
			}

			// Verify next set uses same weight (MRS keeps weight constant)
			if result.NextSet == nil {
				t.Fatalf("set %d: expected next set, got nil", setNumber)
			}
			if result.NextSet.Weight != baseWeight {
				t.Errorf("set %d: expected weight %.1f, got %.1f", setNumber, baseWeight, result.NextSet.Weight)
			}
		}
	}

	// Final verification: check completion state
	req := service.NextSetRequest{
		SessionID:      sessionID,
		PrescriptionID: prescriptionID,
		UserID:         userID,
	}
	result, err := sessionService.GetNextSet(ctx, req)
	if err != nil {
		t.Fatalf("final GetNextSet failed: %v", err)
	}

	if !result.IsComplete {
		t.Error("expected workout to be complete")
	}
	if result.TotalRepsCompleted != 28 {
		t.Errorf("expected 28 total reps, got %d", result.TotalRepsCompleted)
	}
	if result.TotalSetsCompleted != 4 {
		t.Errorf("expected 4 sets completed, got %d", result.TotalSetsCompleted)
	}
	if result.TerminationReason == "" {
		t.Error("expected termination reason to be set")
	}
	// Should indicate target reached
	expectedReason := "Target total reps reached (28/25)"
	if result.TerminationReason != expectedReason {
		t.Errorf("expected termination reason '%s', got '%s'", expectedReason, result.TerminationReason)
	}
}

// TestGZCL_VDIP_T3_FourMRS simulates a GZCL T3 workout using 4 MRS blocks.
//
// GZCL T3 Protocol:
// - 4 Max Rep Sets (MRS)
// - Target: 40 total reps (higher volume accessory work)
// - Minimum: 5 reps per set
// - Lighter weight, higher rep focus
func TestGZCL_VDIP_T3_FourMRS(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	userID := uuid.New().String()
	sessionID := uuid.New().String()
	prescriptionID := uuid.New().String()
	liftID := uuid.New().String()

	// Create MRS scheme for T3: 4 MRS, target 40 total reps, minimum 5 per set
	mrsScheme, err := setscheme.NewMRS(40, 5, 10, 4) // target=40, min=5, max=10, numMRS=4
	if err != nil {
		t.Fatalf("failed to create MRS scheme: %v", err)
	}

	presc := &prescription.Prescription{
		ID:        prescriptionID,
		LiftID:    liftID,
		SetScheme: mrsScheme,
	}

	prescRepo := newMockPrescriptionRepo()
	prescRepo.Add(presc)
	loggedSetLister := newMockLoggedSetLister()
	sessionService := service.NewSessionService(prescRepo, loggedSetLister)

	// T3 workout simulation (lat pulldowns, curls, etc.)
	baseWeight := 100.0
	workoutSets := []int{15, 12, 10, 8} // Typical T3 performance (higher reps)

	totalReps := 0
	for i, reps := range workoutSets {
		setNumber := i + 1
		totalReps += reps

		input := loggedset.CreateLoggedSetInput{
			UserID:         userID,
			SessionID:      sessionID,
			PrescriptionID: prescriptionID,
			LiftID:         liftID,
			SetNumber:      setNumber,
			Weight:         baseWeight,
			TargetReps:     5,
			RepsPerformed:  reps,
			IsAMRAP:        false,
		}

		logged, valResult := loggedset.NewLoggedSet(input, uuid.New().String())
		if !valResult.Valid {
			t.Fatalf("set %d: failed to create logged set: %v", setNumber, valResult.Errors)
		}
		loggedSetLister.Add(*logged)
	}

	// Verify completion
	req := service.NextSetRequest{
		SessionID:      sessionID,
		PrescriptionID: prescriptionID,
		UserID:         userID,
	}
	result, err := sessionService.GetNextSet(ctx, req)
	if err != nil {
		t.Fatalf("GetNextSet failed: %v", err)
	}

	if !result.IsComplete {
		t.Error("expected workout to be complete")
	}
	// Total: 15 + 12 + 10 + 8 = 45 reps (exceeds 40 target)
	if result.TotalRepsCompleted != 45 {
		t.Errorf("expected 45 total reps, got %d", result.TotalRepsCompleted)
	}
	if result.TotalSetsCompleted != 4 {
		t.Errorf("expected 4 sets, got %d", result.TotalSetsCompleted)
	}
}

// TestGZCL_MRS_FailureTermination tests MRS termination when lifter fails
// to hit minimum reps (technical failure).
//
// Scenario:
// - Heavy weight causes premature fatigue
// - Lifter fails to hit 3 reps on a set
// - Workout terminates due to failure
func TestGZCL_MRS_FailureTermination(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	userID := uuid.New().String()
	sessionID := uuid.New().String()
	prescriptionID := uuid.New().String()
	liftID := uuid.New().String()

	// MRS scheme: target 25, minimum 3 per set
	mrsScheme, err := setscheme.NewMRS(25, 3, 10, 3)
	if err != nil {
		t.Fatalf("failed to create MRS scheme: %v", err)
	}

	presc := &prescription.Prescription{
		ID:        prescriptionID,
		LiftID:    liftID,
		SetScheme: mrsScheme,
	}

	prescRepo := newMockPrescriptionRepo()
	prescRepo.Add(presc)
	loggedSetLister := newMockLoggedSetLister()
	sessionService := service.NewSessionService(prescRepo, loggedSetLister)

	// Workout where lifter fails on set 3
	baseWeight := 250.0 // Heavier weight causes early failure
	workoutSets := []int{6, 4, 2} // Set 3 fails to hit minimum of 3

	totalReps := 0
	for i, reps := range workoutSets {
		setNumber := i + 1
		totalReps += reps

		input := loggedset.CreateLoggedSetInput{
			UserID:         userID,
			SessionID:      sessionID,
			PrescriptionID: prescriptionID,
			LiftID:         liftID,
			SetNumber:      setNumber,
			Weight:         baseWeight,
			TargetReps:     3,
			RepsPerformed:  reps,
			IsAMRAP:        false,
		}

		logged, valResult := loggedset.NewLoggedSet(input, uuid.New().String())
		if !valResult.Valid {
			t.Fatalf("set %d: failed to create logged set: %v", setNumber, valResult.Errors)
		}
		loggedSetLister.Add(*logged)
	}

	// Verify failure termination
	req := service.NextSetRequest{
		SessionID:      sessionID,
		PrescriptionID: prescriptionID,
		UserID:         userID,
	}
	result, err := sessionService.GetNextSet(ctx, req)
	if err != nil {
		t.Fatalf("GetNextSet failed: %v", err)
	}

	if !result.IsComplete {
		t.Error("expected workout to be complete due to failure")
	}
	// Total: 6 + 4 + 2 = 12 (below target of 25 due to failure)
	if result.TotalRepsCompleted != 12 {
		t.Errorf("expected 12 total reps, got %d", result.TotalRepsCompleted)
	}
	// Termination reason should indicate failure
	expectedReason := "Failed to hit minimum reps (2/3)"
	if result.TerminationReason != expectedReason {
		t.Errorf("expected termination reason '%s', got '%s'", expectedReason, result.TerminationReason)
	}
}

// =============================================================================
// RTS INTERMEDIATE (FATIGUE DROP) INTEGRATION TESTS
// =============================================================================

// TestRTS_LoadDrop_RPETermination simulates an RTS-style load drop workout
// where weight decreases each set until target RPE is reached.
//
// RTS Load Drop Protocol:
// - Start at target RPE (e.g., 8)
// - Drop weight by 5% each set
// - Stop when RPE reaches 10 (max effort)
//
// Example session:
// - Set 1: 315 lbs x 3 @ RPE 8.0
// - Set 2: 299 lbs x 3 @ RPE 8.5
// - Set 3: 284 lbs x 3 @ RPE 9.0
// - Set 4: 270 lbs x 3 @ RPE 9.5
// - Set 5: 256 lbs x 3 @ RPE 10.0 → STOP
func TestRTS_LoadDrop_RPETermination(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	userID := uuid.New().String()
	sessionID := uuid.New().String()
	prescriptionID := uuid.New().String()
	liftID := uuid.New().String()

	// Create FatigueDrop scheme: 3 reps, start RPE 8, stop RPE 10, drop 5%
	fdScheme, err := setscheme.NewFatigueDrop(3, 8.0, 10.0, 0.05, 10)
	if err != nil {
		t.Fatalf("failed to create FatigueDrop scheme: %v", err)
	}

	presc := &prescription.Prescription{
		ID:        prescriptionID,
		LiftID:    liftID,
		SetScheme: fdScheme,
	}

	prescRepo := newMockPrescriptionRepo()
	prescRepo.Add(presc)
	loggedSetLister := newMockLoggedSetLister()
	sessionService := service.NewSessionService(prescRepo, loggedSetLister)

	// Generate first set
	baseWeight := 315.0
	genCtx := setscheme.DefaultSetGenerationContext()
	initialSets, err := fdScheme.GenerateSets(baseWeight, genCtx)
	if err != nil {
		t.Fatalf("failed to generate initial sets: %v", err)
	}
	if len(initialSets) != 1 {
		t.Fatalf("expected 1 initial set, got %d", len(initialSets))
	}

	// Workout simulation with progressive RPE increase
	workoutSets := []struct {
		weight float64
		rpe    float64
	}{
		{315.0, 8.0},  // Set 1: Starting weight
		{299.0, 8.5},  // Set 2: 5% drop, RPE increases
		{284.0, 9.0},  // Set 3: 5% drop
		{270.0, 9.5},  // Set 4: 5% drop
		{256.0, 10.0}, // Set 5: Stop RPE reached
	}

	var lastResult *service.NextSetResult
	for i, setData := range workoutSets {
		setNumber := i + 1
		rpe := setData.rpe

		input := loggedset.CreateLoggedSetInput{
			UserID:         userID,
			SessionID:      sessionID,
			PrescriptionID: prescriptionID,
			LiftID:         liftID,
			SetNumber:      setNumber,
			Weight:         setData.weight,
			TargetReps:     3,
			RepsPerformed:  3, // Always hit target reps
			IsAMRAP:        false,
			RPE:            &rpe,
		}

		logged, valResult := loggedset.NewLoggedSet(input, uuid.New().String())
		if !valResult.Valid {
			t.Fatalf("set %d: failed to create logged set: %v", setNumber, valResult.Errors)
		}
		loggedSetLister.Add(*logged)

		// Check for next set
		req := service.NextSetRequest{
			SessionID:      sessionID,
			PrescriptionID: prescriptionID,
			UserID:         userID,
		}
		result, err := sessionService.GetNextSet(ctx, req)
		if err != nil {
			t.Fatalf("set %d: GetNextSet failed: %v", setNumber, err)
		}
		lastResult = result

		if result.IsComplete {
			// Should complete after RPE 10
			if rpe < 10.0 {
				t.Errorf("set %d: unexpected completion at RPE %.1f", setNumber, rpe)
			}
			break
		}

		// Verify next set weight drops by ~5%
		if result.NextSet != nil && i < len(workoutSets)-1 {
			expectedNextWeight := workoutSets[i+1].weight
			// Allow for rounding differences (5 lb increment)
			if math.Abs(result.NextSet.Weight-expectedNextWeight) > 5.0 {
				t.Errorf("set %d: expected next weight ~%.1f, got %.1f",
					setNumber, expectedNextWeight, result.NextSet.Weight)
			}
		}
	}

	// Verify final state
	if lastResult == nil {
		t.Fatal("expected final result")
	}
	if !lastResult.IsComplete {
		t.Error("expected workout to be complete")
	}
	if lastResult.TotalSetsCompleted != 5 {
		t.Errorf("expected 5 sets completed, got %d", lastResult.TotalSetsCompleted)
	}
	// Termination reason should indicate RPE target
	expectedReason := "Target RPE reached (10.0/10.0)"
	if lastResult.TerminationReason != expectedReason {
		t.Errorf("expected termination reason '%s', got '%s'", expectedReason, lastResult.TerminationReason)
	}
}

// TestRTS_LoadDrop_MaxSetsLimit tests the safety limit on maximum sets.
//
// Scenario:
// - RPE doesn't increase as expected (lifter adapts well to fatigue)
// - Safety limit of max sets triggers termination
func TestRTS_LoadDrop_MaxSetsLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	userID := uuid.New().String()
	sessionID := uuid.New().String()
	prescriptionID := uuid.New().String()
	liftID := uuid.New().String()

	// Create FatigueDrop with low max sets for testing
	// StopRPE 10 but max sets 3 - should hit max sets first
	fdScheme, err := setscheme.NewFatigueDrop(3, 8.0, 10.0, 0.05, 3) // max 3 sets
	if err != nil {
		t.Fatalf("failed to create FatigueDrop scheme: %v", err)
	}

	presc := &prescription.Prescription{
		ID:        prescriptionID,
		LiftID:    liftID,
		SetScheme: fdScheme,
	}

	prescRepo := newMockPrescriptionRepo()
	prescRepo.Add(presc)
	loggedSetLister := newMockLoggedSetLister()
	sessionService := service.NewSessionService(prescRepo, loggedSetLister)

	// Workout where RPE stays low (lifter adapts well)
	workoutSets := []struct {
		weight float64
		rpe    float64
	}{
		{315.0, 8.0}, // Set 1
		{299.0, 8.5}, // Set 2
		{284.0, 9.0}, // Set 3 - max sets reached before RPE 10
	}

	var lastResult *service.NextSetResult
	for i, setData := range workoutSets {
		setNumber := i + 1
		rpe := setData.rpe

		input := loggedset.CreateLoggedSetInput{
			UserID:         userID,
			SessionID:      sessionID,
			PrescriptionID: prescriptionID,
			LiftID:         liftID,
			SetNumber:      setNumber,
			Weight:         setData.weight,
			TargetReps:     3,
			RepsPerformed:  3,
			IsAMRAP:        false,
			RPE:            &rpe,
		}

		logged, valResult := loggedset.NewLoggedSet(input, uuid.New().String())
		if !valResult.Valid {
			t.Fatalf("set %d: failed to create logged set: %v", setNumber, valResult.Errors)
		}
		loggedSetLister.Add(*logged)

		req := service.NextSetRequest{
			SessionID:      sessionID,
			PrescriptionID: prescriptionID,
			UserID:         userID,
		}
		result, err := sessionService.GetNextSet(ctx, req)
		if err != nil {
			t.Fatalf("set %d: GetNextSet failed: %v", setNumber, err)
		}
		lastResult = result
	}

	// Verify max sets termination
	if !lastResult.IsComplete {
		t.Error("expected workout to be complete")
	}
	if lastResult.TotalSetsCompleted != 3 {
		t.Errorf("expected 3 sets completed, got %d", lastResult.TotalSetsCompleted)
	}
	expectedReason := "Maximum sets reached (safety limit)"
	if lastResult.TerminationReason != expectedReason {
		t.Errorf("expected termination reason '%s', got '%s'", expectedReason, lastResult.TerminationReason)
	}
}

// TestRTS_RepeatSets_RPEIncrease tests repeated sets at the same weight
// where RPE increases until termination.
//
// This is a variation of RTS training where weight stays constant
// but RPE tracks fatigue accumulation.
func TestRTS_RepeatSets_RPEIncrease(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	userID := uuid.New().String()
	sessionID := uuid.New().String()
	prescriptionID := uuid.New().String()
	liftID := uuid.New().String()

	// FatigueDrop with 0% drop = repeat sets at same weight
	// Stop when RPE reaches 9.5 (stop before true max)
	fdScheme, err := setscheme.NewFatigueDrop(5, 7.0, 9.5, 0.0, 10) // 0% drop
	if err != nil {
		t.Fatalf("failed to create FatigueDrop scheme: %v", err)
	}

	presc := &prescription.Prescription{
		ID:        prescriptionID,
		LiftID:    liftID,
		SetScheme: fdScheme,
	}

	prescRepo := newMockPrescriptionRepo()
	prescRepo.Add(presc)
	loggedSetLister := newMockLoggedSetLister()
	sessionService := service.NewSessionService(prescRepo, loggedSetLister)

	// Same weight, RPE increases each set
	baseWeight := 200.0
	rpeProgression := []float64{7.0, 7.5, 8.0, 8.5, 9.0, 9.5}

	var lastResult *service.NextSetResult
	for i, rpe := range rpeProgression {
		setNumber := i + 1
		currentRPE := rpe

		input := loggedset.CreateLoggedSetInput{
			UserID:         userID,
			SessionID:      sessionID,
			PrescriptionID: prescriptionID,
			LiftID:         liftID,
			SetNumber:      setNumber,
			Weight:         baseWeight, // Same weight every set
			TargetReps:     5,
			RepsPerformed:  5,
			IsAMRAP:        false,
			RPE:            &currentRPE,
		}

		logged, valResult := loggedset.NewLoggedSet(input, uuid.New().String())
		if !valResult.Valid {
			t.Fatalf("set %d: failed to create logged set: %v", setNumber, valResult.Errors)
		}
		loggedSetLister.Add(*logged)

		req := service.NextSetRequest{
			SessionID:      sessionID,
			PrescriptionID: prescriptionID,
			UserID:         userID,
		}
		result, err := sessionService.GetNextSet(ctx, req)
		if err != nil {
			t.Fatalf("set %d: GetNextSet failed: %v", setNumber, err)
		}
		lastResult = result

		if result.IsComplete {
			break
		}

		// Verify weight stays the same when drop is 0%
		if result.NextSet != nil {
			// With 0% drop, weight should stay at baseWeight
			// However, due to rounding, it might vary slightly
			if math.Abs(result.NextSet.Weight-baseWeight) > 0.001 {
				t.Errorf("set %d: expected next weight %.1f (no drop), got %.1f",
					setNumber, baseWeight, result.NextSet.Weight)
			}
		}
	}

	// Verify completion at RPE 9.5
	if !lastResult.IsComplete {
		t.Error("expected workout to be complete")
	}
	if lastResult.TotalSetsCompleted != 6 {
		t.Errorf("expected 6 sets completed, got %d", lastResult.TotalSetsCompleted)
	}
	expectedReason := "Target RPE reached (9.5/9.5)"
	if lastResult.TerminationReason != expectedReason {
		t.Errorf("expected termination reason '%s', got '%s'", expectedReason, lastResult.TerminationReason)
	}
}

// =============================================================================
// ADDITIONAL EDGE CASE TESTS
// =============================================================================

// TestMRS_ExactTargetReps tests MRS termination when exactly hitting the target.
func TestMRS_ExactTargetReps(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	userID := uuid.New().String()
	sessionID := uuid.New().String()
	prescriptionID := uuid.New().String()
	liftID := uuid.New().String()

	// Target exactly 15 reps
	mrsScheme, err := setscheme.NewMRS(15, 3, 10, 3)
	if err != nil {
		t.Fatalf("failed to create MRS scheme: %v", err)
	}

	presc := &prescription.Prescription{
		ID:        prescriptionID,
		LiftID:    liftID,
		SetScheme: mrsScheme,
	}

	prescRepo := newMockPrescriptionRepo()
	prescRepo.Add(presc)
	loggedSetLister := newMockLoggedSetLister()
	sessionService := service.NewSessionService(prescRepo, loggedSetLister)

	// Sets that exactly hit 15 reps
	workoutSets := []int{5, 5, 5} // Exactly 15 total

	baseWeight := 225.0
	for i, reps := range workoutSets {
		setNumber := i + 1

		input := loggedset.CreateLoggedSetInput{
			UserID:         userID,
			SessionID:      sessionID,
			PrescriptionID: prescriptionID,
			LiftID:         liftID,
			SetNumber:      setNumber,
			Weight:         baseWeight,
			TargetReps:     3,
			RepsPerformed:  reps,
			IsAMRAP:        false,
		}

		logged, valResult := loggedset.NewLoggedSet(input, uuid.New().String())
		if !valResult.Valid {
			t.Fatalf("set %d: failed to create logged set: %v", setNumber, valResult.Errors)
		}
		loggedSetLister.Add(*logged)
	}

	req := service.NextSetRequest{
		SessionID:      sessionID,
		PrescriptionID: prescriptionID,
		UserID:         userID,
	}
	result, err := sessionService.GetNextSet(ctx, req)
	if err != nil {
		t.Fatalf("GetNextSet failed: %v", err)
	}

	if !result.IsComplete {
		t.Error("expected workout to be complete")
	}
	if result.TotalRepsCompleted != 15 {
		t.Errorf("expected exactly 15 total reps, got %d", result.TotalRepsCompleted)
	}
	expectedReason := "Target total reps reached (15/15)"
	if result.TerminationReason != expectedReason {
		t.Errorf("expected termination reason '%s', got '%s'", expectedReason, result.TerminationReason)
	}
}

// TestFatigueDrop_WeightRounding tests that FatigueDrop correctly rounds weights.
func TestFatigueDrop_WeightRounding(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Create FatigueDrop scheme
	fdScheme, err := setscheme.NewFatigueDrop(3, 8.0, 10.0, 0.05, 10)
	if err != nil {
		t.Fatalf("failed to create FatigueDrop scheme: %v", err)
	}

	genCtx := setscheme.DefaultSetGenerationContext()

	// Test starting weight that results in odd numbers after 5% drop
	baseWeight := 315.0 // 315 * 0.95 = 299.25 → should round to 295 or 300

	initialSets, err := fdScheme.GenerateSets(baseWeight, genCtx)
	if err != nil {
		t.Fatalf("failed to generate sets: %v", err)
	}

	// Build history for next set generation
	history := []setscheme.GeneratedSet{initialSets[0]}

	rpe := 8.0
	termCtx := setscheme.TerminationContext{
		SetNumber:  1,
		TotalSets:  1,
		LastRPE:    &rpe,
		LastReps:   3,
		TotalReps:  3,
		TargetReps: 3,
	}

	nextSet, shouldContinue := fdScheme.GenerateNextSet(genCtx, history, termCtx)
	if !shouldContinue {
		t.Fatal("expected to continue")
	}
	if nextSet == nil {
		t.Fatal("expected next set")
	}

	// 315 * 0.95 = 299.25, should round down to 295 (conservative for fatigue)
	expectedWeight := 295.0
	if nextSet.Weight != expectedWeight {
		t.Errorf("expected weight %.1f after rounding, got %.1f", expectedWeight, nextSet.Weight)
	}
}

// TestFatigueDrop_LargeDropPercent tests FatigueDrop with a larger drop percentage.
func TestFatigueDrop_LargeDropPercent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// 10% drop per set
	fdScheme, err := setscheme.NewFatigueDrop(5, 7.0, 9.0, 0.10, 10)
	if err != nil {
		t.Fatalf("failed to create FatigueDrop scheme: %v", err)
	}

	genCtx := setscheme.DefaultSetGenerationContext()
	baseWeight := 200.0

	initialSets, err := fdScheme.GenerateSets(baseWeight, genCtx)
	if err != nil {
		t.Fatalf("failed to generate sets: %v", err)
	}

	history := []setscheme.GeneratedSet{initialSets[0]}

	rpe := 7.5
	termCtx := setscheme.TerminationContext{
		SetNumber:  1,
		TotalSets:  1,
		LastRPE:    &rpe,
		LastReps:   5,
		TotalReps:  5,
		TargetReps: 5,
	}

	nextSet, shouldContinue := fdScheme.GenerateNextSet(genCtx, history, termCtx)
	if !shouldContinue {
		t.Fatal("expected to continue")
	}
	if nextSet == nil {
		t.Fatal("expected next set")
	}

	// 200 * 0.90 = 180
	expectedWeight := 180.0
	if nextSet.Weight != expectedWeight {
		t.Errorf("expected weight %.1f after 10%% drop, got %.1f", expectedWeight, nextSet.Weight)
	}
}

// =============================================================================
// ACCEPTANCE CRITERIA TESTS
// =============================================================================

// TestFatigueProtocols_AcceptanceCriteria validates all business requirements.
func TestFatigueProtocols_AcceptanceCriteria(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test")
	}

	// AC1: MRS scheme generates provisional first set
	t.Run("AC1: MRS generates provisional first set", func(t *testing.T) {
		scheme, _ := setscheme.NewMRS(25, 3, 10, 3)
		sets, err := scheme.GenerateSets(225.0, setscheme.DefaultSetGenerationContext())
		if err != nil {
			t.Fatalf("failed to generate sets: %v", err)
		}
		if len(sets) != 1 {
			t.Errorf("expected 1 set, got %d", len(sets))
		}
		if !sets[0].IsProvisional {
			t.Error("expected set to be provisional")
		}
	})

	// AC2: MRS terminates at total reps target
	t.Run("AC2: MRS terminates at total reps target", func(t *testing.T) {
		scheme, _ := setscheme.NewMRS(10, 3, 10, 3)
		ctx := setscheme.DefaultSetGenerationContext()
		history := []setscheme.GeneratedSet{{SetNumber: 1, Weight: 225.0}}
		termCtx := setscheme.TerminationContext{
			TotalSets:  2,
			TotalReps:  12, // Exceeds target of 10
			LastReps:   5,
			TargetReps: 3,
		}
		_, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		if shouldContinue {
			t.Error("expected termination at total reps target")
		}
	})

	// AC3: MRS terminates at rep failure
	t.Run("AC3: MRS terminates at rep failure", func(t *testing.T) {
		scheme, _ := setscheme.NewMRS(25, 3, 10, 3)
		ctx := setscheme.DefaultSetGenerationContext()
		history := []setscheme.GeneratedSet{{SetNumber: 1, Weight: 225.0}}
		termCtx := setscheme.TerminationContext{
			TotalSets:  2,
			TotalReps:  10,
			LastReps:   2, // Below minimum of 3
			TargetReps: 3,
		}
		_, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		if shouldContinue {
			t.Error("expected termination at rep failure")
		}
	})

	// AC4: FatigueDrop generates provisional first set
	t.Run("AC4: FatigueDrop generates provisional first set", func(t *testing.T) {
		scheme, _ := setscheme.NewFatigueDrop(3, 8.0, 10.0, 0.05, 10)
		sets, err := scheme.GenerateSets(315.0, setscheme.DefaultSetGenerationContext())
		if err != nil {
			t.Fatalf("failed to generate sets: %v", err)
		}
		if len(sets) != 1 {
			t.Errorf("expected 1 set, got %d", len(sets))
		}
		if !sets[0].IsProvisional {
			t.Error("expected set to be provisional")
		}
	})

	// AC5: FatigueDrop terminates at RPE threshold
	t.Run("AC5: FatigueDrop terminates at RPE threshold", func(t *testing.T) {
		scheme, _ := setscheme.NewFatigueDrop(3, 8.0, 10.0, 0.05, 10)
		ctx := setscheme.DefaultSetGenerationContext()
		history := []setscheme.GeneratedSet{{SetNumber: 1, Weight: 315.0}}
		rpe := 10.0
		termCtx := setscheme.TerminationContext{
			TotalSets: 3,
			LastRPE:   &rpe, // Reached stop RPE
			LastReps:  3,
		}
		_, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		if shouldContinue {
			t.Error("expected termination at RPE threshold")
		}
	})

	// AC6: FatigueDrop drops weight correctly
	t.Run("AC6: FatigueDrop drops weight correctly", func(t *testing.T) {
		scheme, _ := setscheme.NewFatigueDrop(3, 8.0, 10.0, 0.05, 10)
		ctx := setscheme.DefaultSetGenerationContext()
		history := []setscheme.GeneratedSet{{SetNumber: 1, Weight: 300.0}} // Even number for clean math
		rpe := 8.0
		termCtx := setscheme.TerminationContext{
			TotalSets: 1,
			LastRPE:   &rpe,
			LastReps:  3,
		}
		nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		if !shouldContinue {
			t.Fatal("expected to continue")
		}
		// 300 * 0.95 = 285
		if nextSet.Weight != 285.0 {
			t.Errorf("expected weight 285.0, got %.1f", nextSet.Weight)
		}
	})

	// AC7: SessionService provides correct termination reasons
	t.Run("AC7: SessionService provides correct termination reasons", func(t *testing.T) {
		// Setup MRS
		mrsScheme, _ := setscheme.NewMRS(10, 3, 10, 3)
		presc := &prescription.Prescription{ID: "p1", SetScheme: mrsScheme}
		prescRepo := newMockPrescriptionRepo()
		prescRepo.Add(presc)
		loggedSetLister := newMockLoggedSetLister()
		sessionService := service.NewSessionService(prescRepo, loggedSetLister)

		// Log sets to hit target
		for i := 1; i <= 2; i++ {
			input := loggedset.CreateLoggedSetInput{
				UserID:         "u1",
				SessionID:      "s1",
				PrescriptionID: "p1",
				LiftID:         "l1",
				SetNumber:      i,
				Weight:         225.0,
				TargetReps:     3,
				RepsPerformed:  6,
			}
			logged, _ := loggedset.NewLoggedSet(input, uuid.New().String())
			loggedSetLister.Add(*logged)
		}

		req := service.NextSetRequest{SessionID: "s1", PrescriptionID: "p1", UserID: "u1"}
		result, _ := sessionService.GetNextSet(context.Background(), req)

		if !result.IsComplete {
			t.Error("expected completion")
		}
		if result.TerminationReason == "" {
			t.Error("expected termination reason")
		}
	})

	// AC8: LoggedSet correctly stores RPE
	t.Run("AC8: LoggedSet correctly stores RPE", func(t *testing.T) {
		rpe := 8.5
		input := loggedset.CreateLoggedSetInput{
			UserID:         "u1",
			SessionID:      "s1",
			PrescriptionID: "p1",
			LiftID:         "l1",
			SetNumber:      1,
			Weight:         225.0,
			TargetReps:     3,
			RepsPerformed:  3,
			RPE:            &rpe,
		}
		logged, valResult := loggedset.NewLoggedSet(input, uuid.New().String())
		if !valResult.Valid {
			t.Fatalf("failed to create logged set: %v", valResult.Errors)
		}
		if logged.RPE == nil || *logged.RPE != 8.5 {
			t.Error("expected RPE 8.5 to be stored")
		}
	})
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// createTestLoggedSet is a helper to create logged sets for tests.
func createTestLoggedSet(
	userID, sessionID, prescriptionID, liftID string,
	setNumber int, weight float64, targetReps, repsPerformed int,
	rpe *float64,
) loggedset.LoggedSet {
	input := loggedset.CreateLoggedSetInput{
		UserID:         userID,
		SessionID:      sessionID,
		PrescriptionID: prescriptionID,
		LiftID:         liftID,
		SetNumber:      setNumber,
		Weight:         weight,
		TargetReps:     targetReps,
		RepsPerformed:  repsPerformed,
		IsAMRAP:        false,
		RPE:            rpe,
	}
	logged, _ := loggedset.NewLoggedSet(input, uuid.New().String())
	return *logged
}

// Ensure time import is used
var _ = time.Now
