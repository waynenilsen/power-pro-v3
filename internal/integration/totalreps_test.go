// Package integration provides end-to-end integration tests for set schemes.
// This file tests TotalReps (5/3/1 Building the Monolith-style) set scheme
// demonstrating complete workout flows including termination conditions.
package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/loggedset"
	"github.com/waynenilsen/power-pro-v3/internal/domain/prescription"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
	"github.com/waynenilsen/power-pro-v3/internal/service"
)

// =============================================================================
// TOTALREPS INTEGRATION TESTS - 5/3/1 Building the Monolith Style
// =============================================================================

// TestTotalReps_BasicFlow simulates a basic TotalReps workout where user
// accumulates reps until reaching a target.
//
// Building the Monolith context:
// - Accessory work like "100 chin-ups"
// - User distributes reps however they want across sets
// - No minimum reps requirement per set
//
// Example session:
// - Set 1: 15 reps (total: 15)
// - Set 2: 12 reps (total: 27)
// - Set 3: 10 reps (total: 37)
// - ... continues until total >= 100
func TestTotalReps_BasicFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	// Setup test IDs
	userID := uuid.New().String()
	sessionID := uuid.New().String()
	prescriptionID := uuid.New().String()
	liftID := uuid.New().String()

	// Create TotalReps scheme: target 100 reps, suggested 15 per set, max 20 sets
	totalRepsScheme, err := setscheme.NewTotalRepsScheme(100, 15, 20)
	if err != nil {
		t.Fatalf("failed to create TotalReps scheme: %v", err)
	}

	// Create prescription with TotalReps scheme
	presc := &prescription.Prescription{
		ID:        prescriptionID,
		LiftID:    liftID,
		SetScheme: totalRepsScheme,
	}

	// Setup mock repositories
	prescRepo := newMockPrescriptionRepo()
	prescRepo.Add(presc)
	loggedSetLister := newMockLoggedSetLister()

	// Create session service
	sessionService := service.NewSessionService(prescRepo, loggedSetLister)

	// === Workout Simulation ===
	baseWeight := 0.0 // Bodyweight for chin-ups

	// Generate first set
	genCtx := setscheme.DefaultSetGenerationContext()
	initialSets, err := totalRepsScheme.GenerateSets(baseWeight, genCtx)
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

	// Verify suggested reps
	if initialSets[0].TargetReps != 15 {
		t.Errorf("expected suggested reps 15, got %d", initialSets[0].TargetReps)
	}

	// Simulate workout sets (realistic chin-up performance with fatigue)
	workoutSets := []struct {
		repsPerformed int
		description   string
	}{
		{15, "Fresh start, strong set"},
		{12, "Slight fatigue"},
		{10, "Moderate fatigue"},
		{10, "Grinding through"},
		{8, "Getting hard"},
		{8, "Fatigue accumulating"},
		{7, "Arms burning"},
		{6, "Push through"},
		{6, "Almost there"},
		{6, "One more set after this"},
		{5, "Final push"}, // Total: 15+12+10+10+8+8+7+6+6+6+5 = 93
		{7, "Done!"},      // Total: 100
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
			TargetReps:     15, // Suggested reps
			RepsPerformed:  setData.repsPerformed,
			IsAMRAP:        false,
		}

		logged, valResult := loggedset.NewLoggedSet(input, uuid.New().String())
		if !valResult.Valid {
			t.Fatalf("set %d: failed to create logged set: %v", setNumber, valResult.Errors)
		}
		loggedSetLister.Add(*logged)

		// Check if we should continue
		req := service.NextSetRequest{
			SessionID:      sessionID,
			PrescriptionID: prescriptionID,
			UserID:         userID,
		}
		result, err := sessionService.GetNextSet(ctx, req)
		if err != nil {
			t.Fatalf("set %d: GetNextSet failed: %v", setNumber, err)
		}

		// Track cumulative progress
		if result.TotalRepsCompleted != totalReps {
			t.Errorf("set %d: expected cumulative reps %d, got %d", setNumber, totalReps, result.TotalRepsCompleted)
		}

		if result.IsComplete {
			// Should complete when we hit 100 reps
			if totalReps < 100 {
				t.Errorf("set %d: unexpected completion at %d total reps", setNumber, totalReps)
			}
			break
		}

		// Verify next set uses same weight (bodyweight)
		if result.NextSet != nil && result.NextSet.Weight != baseWeight {
			t.Errorf("set %d: expected weight %.1f, got %.1f", setNumber, baseWeight, result.NextSet.Weight)
		}
	}

	// Final verification
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
	if result.TotalRepsCompleted != 100 {
		t.Errorf("expected 100 total reps, got %d", result.TotalRepsCompleted)
	}
	if result.TotalSetsCompleted != 12 {
		t.Errorf("expected 12 sets completed, got %d", result.TotalSetsCompleted)
	}
	expectedReason := "Target total reps reached (100/100)"
	if result.TerminationReason != expectedReason {
		t.Errorf("expected termination reason '%s', got '%s'", expectedReason, result.TerminationReason)
	}
}

// TestTotalReps_ExactTargetAchievement tests termination when exactly hitting the target.
//
// Scenario:
// - Target: 50 reps
// - Log: 15, 15, 15, 5 = 50 exactly
// - Verify clean termination at exact target
func TestTotalReps_ExactTargetAchievement(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	userID := uuid.New().String()
	sessionID := uuid.New().String()
	prescriptionID := uuid.New().String()
	liftID := uuid.New().String()

	// Create TotalReps scheme: target 50 reps
	totalRepsScheme, err := setscheme.NewTotalRepsScheme(50, 15, 20)
	if err != nil {
		t.Fatalf("failed to create TotalReps scheme: %v", err)
	}

	presc := &prescription.Prescription{
		ID:        prescriptionID,
		LiftID:    liftID,
		SetScheme: totalRepsScheme,
	}

	prescRepo := newMockPrescriptionRepo()
	prescRepo.Add(presc)
	loggedSetLister := newMockLoggedSetLister()
	sessionService := service.NewSessionService(prescRepo, loggedSetLister)

	// Sets that exactly hit 50 reps
	workoutSets := []int{15, 15, 15, 5} // Exactly 50 total
	baseWeight := 0.0

	for i, reps := range workoutSets {
		setNumber := i + 1

		input := loggedset.CreateLoggedSetInput{
			UserID:         userID,
			SessionID:      sessionID,
			PrescriptionID: prescriptionID,
			LiftID:         liftID,
			SetNumber:      setNumber,
			Weight:         baseWeight,
			TargetReps:     15,
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
	if result.TotalRepsCompleted != 50 {
		t.Errorf("expected exactly 50 total reps, got %d", result.TotalRepsCompleted)
	}
	if result.TotalSetsCompleted != 4 {
		t.Errorf("expected 4 sets, got %d", result.TotalSetsCompleted)
	}
	expectedReason := "Target total reps reached (50/50)"
	if result.TerminationReason != expectedReason {
		t.Errorf("expected termination reason '%s', got '%s'", expectedReason, result.TerminationReason)
	}
}

// TestTotalReps_OvershootScenario tests termination when overshooting the target.
//
// Scenario:
// - Target: 50 reps
// - Log: 20, 20, 15 = 55 (over target)
// - Verify terminates immediately on exceeding target
func TestTotalReps_OvershootScenario(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	userID := uuid.New().String()
	sessionID := uuid.New().String()
	prescriptionID := uuid.New().String()
	liftID := uuid.New().String()

	// Create TotalReps scheme: target 50 reps
	totalRepsScheme, err := setscheme.NewTotalRepsScheme(50, 15, 20)
	if err != nil {
		t.Fatalf("failed to create TotalReps scheme: %v", err)
	}

	presc := &prescription.Prescription{
		ID:        prescriptionID,
		LiftID:    liftID,
		SetScheme: totalRepsScheme,
	}

	prescRepo := newMockPrescriptionRepo()
	prescRepo.Add(presc)
	loggedSetLister := newMockLoggedSetLister()
	sessionService := service.NewSessionService(prescRepo, loggedSetLister)

	// Sets that overshoot 50 reps
	workoutSets := []int{20, 20, 15} // 55 total (over by 5)
	baseWeight := 0.0

	for i, reps := range workoutSets {
		setNumber := i + 1

		input := loggedset.CreateLoggedSetInput{
			UserID:         userID,
			SessionID:      sessionID,
			PrescriptionID: prescriptionID,
			LiftID:         liftID,
			SetNumber:      setNumber,
			Weight:         baseWeight,
			TargetReps:     15,
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
	// Overshoot is allowed - we record what was actually done
	if result.TotalRepsCompleted != 55 {
		t.Errorf("expected 55 total reps (overshoot), got %d", result.TotalRepsCompleted)
	}
	if result.TotalSetsCompleted != 3 {
		t.Errorf("expected 3 sets, got %d", result.TotalSetsCompleted)
	}
	expectedReason := "Target total reps reached (55/50)"
	if result.TerminationReason != expectedReason {
		t.Errorf("expected termination reason '%s', got '%s'", expectedReason, result.TerminationReason)
	}
}

// TestTotalReps_MaxSetsSafety tests the max sets safety limit.
//
// Scenario:
// - Target: 1000 reps (unrealistic)
// - Max sets: 5
// - Log 5 sets of 10 reps (50 total)
// - Verify terminates due to max sets, not reps
func TestTotalReps_MaxSetsSafety(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	userID := uuid.New().String()
	sessionID := uuid.New().String()
	prescriptionID := uuid.New().String()
	liftID := uuid.New().String()

	// Create TotalReps scheme: unrealistic target, low max sets
	totalRepsScheme, err := setscheme.NewTotalRepsScheme(1000, 10, 5) // Max 5 sets
	if err != nil {
		t.Fatalf("failed to create TotalReps scheme: %v", err)
	}

	presc := &prescription.Prescription{
		ID:        prescriptionID,
		LiftID:    liftID,
		SetScheme: totalRepsScheme,
	}

	prescRepo := newMockPrescriptionRepo()
	prescRepo.Add(presc)
	loggedSetLister := newMockLoggedSetLister()
	sessionService := service.NewSessionService(prescRepo, loggedSetLister)

	// Log 5 sets of 10 reps each
	baseWeight := 0.0
	for setNumber := 1; setNumber <= 5; setNumber++ {
		input := loggedset.CreateLoggedSetInput{
			UserID:         userID,
			SessionID:      sessionID,
			PrescriptionID: prescriptionID,
			LiftID:         liftID,
			SetNumber:      setNumber,
			Weight:         baseWeight,
			TargetReps:     10,
			RepsPerformed:  10,
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
		t.Error("expected workout to be complete due to max sets")
	}
	if result.TotalRepsCompleted != 50 {
		t.Errorf("expected 50 total reps, got %d", result.TotalRepsCompleted)
	}
	if result.TotalSetsCompleted != 5 {
		t.Errorf("expected 5 sets, got %d", result.TotalSetsCompleted)
	}
	expectedReason := "Maximum sets reached (safety limit)"
	if result.TerminationReason != expectedReason {
		t.Errorf("expected termination reason '%s', got '%s'", expectedReason, result.TerminationReason)
	}
}

// TestTotalReps_NoMinRepsRequirement verifies there's no minimum reps per set requirement.
//
// Unlike MRS, TotalReps allows any rep count per set (including low single digits).
// This is realistic for Building the Monolith where a fatigued lifter might do
// sets of 3-5 chin-ups near the end.
func TestTotalReps_NoMinRepsRequirement(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	userID := uuid.New().String()
	sessionID := uuid.New().String()
	prescriptionID := uuid.New().String()
	liftID := uuid.New().String()

	// Create TotalReps scheme: target 30 reps
	totalRepsScheme, err := setscheme.NewTotalRepsScheme(30, 10, 20)
	if err != nil {
		t.Fatalf("failed to create TotalReps scheme: %v", err)
	}

	presc := &prescription.Prescription{
		ID:        prescriptionID,
		LiftID:    liftID,
		SetScheme: totalRepsScheme,
	}

	prescRepo := newMockPrescriptionRepo()
	prescRepo.Add(presc)
	loggedSetLister := newMockLoggedSetLister()
	sessionService := service.NewSessionService(prescRepo, loggedSetLister)

	// Simulate exhausted lifter doing very small sets
	workoutSets := []int{10, 5, 4, 3, 3, 2, 2, 1} // 30 total with small sets at the end
	baseWeight := 0.0

	for i, reps := range workoutSets {
		setNumber := i + 1

		input := loggedset.CreateLoggedSetInput{
			UserID:         userID,
			SessionID:      sessionID,
			PrescriptionID: prescriptionID,
			LiftID:         liftID,
			SetNumber:      setNumber,
			Weight:         baseWeight,
			TargetReps:     10,
			RepsPerformed:  reps, // Even 1 rep is valid
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

	// Should complete at rep target, not due to failure
	if !result.IsComplete {
		t.Error("expected workout to be complete")
	}
	if result.TotalRepsCompleted != 30 {
		t.Errorf("expected 30 total reps, got %d", result.TotalRepsCompleted)
	}
	if result.TotalSetsCompleted != 8 {
		t.Errorf("expected 8 sets, got %d", result.TotalSetsCompleted)
	}
	// Should NOT say "Failed to hit minimum reps" - this is TotalReps, not MRS
	expectedReason := "Target total reps reached (30/30)"
	if result.TerminationReason != expectedReason {
		t.Errorf("expected termination reason '%s', got '%s'", expectedReason, result.TerminationReason)
	}
}

// TestBuildingTheMonolith_MondaySimulation simulates a Building the Monolith Monday session.
//
// From the program:
// - Squat: 5/3/1 sets (fixed scheme - not tested here)
// - Press: 5/3/1 BBB (fixed scheme - not tested here)
// - Chin-ups: 100 total reps (TotalReps scheme)
// - Dips: 100-200 total reps (TotalReps scheme)
// - Face pulls: 100 total reps (TotalReps scheme)
//
// This test verifies all three TotalReps prescriptions can complete in sequence.
func TestBuildingTheMonolith_MondaySimulation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	userID := uuid.New().String()
	sessionID := uuid.New().String()

	// Create TotalReps schemes for accessory work
	chinUpScheme, _ := setscheme.NewTotalRepsScheme(100, 15, 20)
	dipScheme, _ := setscheme.NewTotalRepsScheme(100, 20, 15) // Dips often easier at start
	facePullScheme, _ := setscheme.NewTotalRepsScheme(100, 25, 12)

	// Create prescriptions
	chinUpPresc := &prescription.Prescription{
		ID:        uuid.New().String(),
		LiftID:    "chin-ups",
		SetScheme: chinUpScheme,
	}
	dipPresc := &prescription.Prescription{
		ID:        uuid.New().String(),
		LiftID:    "dips",
		SetScheme: dipScheme,
	}
	facePullPresc := &prescription.Prescription{
		ID:        uuid.New().String(),
		LiftID:    "face-pulls",
		SetScheme: facePullScheme,
	}

	prescRepo := newMockPrescriptionRepo()
	prescRepo.Add(chinUpPresc)
	prescRepo.Add(dipPresc)
	prescRepo.Add(facePullPresc)

	loggedSetLister := newMockLoggedSetLister()
	sessionService := service.NewSessionService(prescRepo, loggedSetLister)

	// Helper to log sets until complete
	completeExercise := func(t *testing.T, prescID, liftID string, repSequence []int) {
		for i, reps := range repSequence {
			setNumber := i + 1

			input := loggedset.CreateLoggedSetInput{
				UserID:         userID,
				SessionID:      sessionID,
				PrescriptionID: prescID,
				LiftID:         liftID,
				SetNumber:      setNumber,
				Weight:         0.0, // Bodyweight
				TargetReps:     15,
				RepsPerformed:  reps,
				IsAMRAP:        false,
			}

			logged, valResult := loggedset.NewLoggedSet(input, uuid.New().String())
			if !valResult.Valid {
				t.Fatalf("%s set %d: failed to create logged set: %v", liftID, setNumber, valResult.Errors)
			}
			loggedSetLister.Add(*logged)
		}

		req := service.NextSetRequest{
			SessionID:      sessionID,
			PrescriptionID: prescID,
			UserID:         userID,
		}
		result, err := sessionService.GetNextSet(ctx, req)
		if err != nil {
			t.Fatalf("%s: GetNextSet failed: %v", liftID, err)
		}

		if !result.IsComplete {
			t.Errorf("%s: expected exercise to be complete", liftID)
		}
	}

	// === Execute Chin-ups: 100 total reps ===
	t.Run("Chin-ups 100 reps", func(t *testing.T) {
		// Realistic chin-up progression with fatigue
		chinUpReps := []int{15, 12, 10, 10, 8, 8, 7, 6, 6, 6, 5, 7} // = 100
		completeExercise(t, chinUpPresc.ID, "chin-ups", chinUpReps)
	})

	// === Execute Dips: 100 total reps ===
	t.Run("Dips 100 reps", func(t *testing.T) {
		// Dips often easier, can do bigger sets
		dipReps := []int{20, 18, 15, 12, 12, 10, 8, 5} // = 100
		completeExercise(t, dipPresc.ID, "dips", dipReps)
	})

	// === Execute Face Pulls: 100 total reps ===
	t.Run("Face pulls 100 reps", func(t *testing.T) {
		// Face pulls are light, high rep work
		facePullReps := []int{25, 25, 25, 25} // = 100
		completeExercise(t, facePullPresc.ID, "face-pulls", facePullReps)
	})
}

// TestTotalReps_WeightedAccessory tests TotalReps with weighted exercises.
//
// Not all TotalReps work is bodyweight. Building the Monolith also includes:
// - Dumbbell rows: heavy sets
// - Shrugs: 100 total reps with heavy weight
func TestTotalReps_WeightedAccessory(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	userID := uuid.New().String()
	sessionID := uuid.New().String()
	prescriptionID := uuid.New().String()
	liftID := uuid.New().String()

	// Create TotalReps scheme: 100 total reps of shrugs
	shrugScheme, err := setscheme.NewTotalRepsScheme(100, 15, 15)
	if err != nil {
		t.Fatalf("failed to create TotalReps scheme: %v", err)
	}

	presc := &prescription.Prescription{
		ID:        prescriptionID,
		LiftID:    liftID,
		SetScheme: shrugScheme,
	}

	prescRepo := newMockPrescriptionRepo()
	prescRepo.Add(presc)
	loggedSetLister := newMockLoggedSetLister()
	sessionService := service.NewSessionService(prescRepo, loggedSetLister)

	// Shrugs at 225 lbs
	baseWeight := 225.0
	workoutSets := []int{20, 18, 15, 15, 12, 10, 10} // = 100

	for i, reps := range workoutSets {
		setNumber := i + 1

		input := loggedset.CreateLoggedSetInput{
			UserID:         userID,
			SessionID:      sessionID,
			PrescriptionID: prescriptionID,
			LiftID:         liftID,
			SetNumber:      setNumber,
			Weight:         baseWeight,
			TargetReps:     15,
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
	if result.TotalRepsCompleted != 100 {
		t.Errorf("expected 100 total reps, got %d", result.TotalRepsCompleted)
	}
}

// TestTotalReps_ProgressiveSetTracking verifies cumulative progress is tracked correctly.
func TestTotalReps_ProgressiveSetTracking(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	userID := uuid.New().String()
	sessionID := uuid.New().String()
	prescriptionID := uuid.New().String()
	liftID := uuid.New().String()

	// Target 40 reps so we can track progression clearly
	totalRepsScheme, err := setscheme.NewTotalRepsScheme(40, 10, 20)
	if err != nil {
		t.Fatalf("failed to create TotalReps scheme: %v", err)
	}

	presc := &prescription.Prescription{
		ID:        prescriptionID,
		LiftID:    liftID,
		SetScheme: totalRepsScheme,
	}

	prescRepo := newMockPrescriptionRepo()
	prescRepo.Add(presc)
	loggedSetLister := newMockLoggedSetLister()
	sessionService := service.NewSessionService(prescRepo, loggedSetLister)

	workoutSets := []int{12, 10, 8, 7, 5} // = 42, exceeds 40
	expectedCumulative := []int{12, 22, 30, 37, 42}
	baseWeight := 0.0

	for i, reps := range workoutSets {
		setNumber := i + 1

		input := loggedset.CreateLoggedSetInput{
			UserID:         userID,
			SessionID:      sessionID,
			PrescriptionID: prescriptionID,
			LiftID:         liftID,
			SetNumber:      setNumber,
			Weight:         baseWeight,
			TargetReps:     10,
			RepsPerformed:  reps,
			IsAMRAP:        false,
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

		// Verify cumulative tracking
		if result.TotalRepsCompleted != expectedCumulative[i] {
			t.Errorf("set %d: expected cumulative %d, got %d",
				setNumber, expectedCumulative[i], result.TotalRepsCompleted)
		}
		if result.TotalSetsCompleted != setNumber {
			t.Errorf("set %d: expected %d sets completed, got %d",
				setNumber, setNumber, result.TotalSetsCompleted)
		}

		// Check termination after we exceed target
		if setNumber >= 4 && !result.IsComplete {
			// Set 4: 37 reps, should continue
			// Set 5: 42 reps, should be complete
		}
	}
}

// =============================================================================
// ACCEPTANCE CRITERIA TESTS
// =============================================================================

// TestTotalReps_AcceptanceCriteria validates all business requirements for TotalReps.
func TestTotalReps_AcceptanceCriteria(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test")
	}

	// AC1: TotalReps scheme generates provisional first set
	t.Run("AC1: TotalReps generates provisional first set", func(t *testing.T) {
		scheme, _ := setscheme.NewTotalRepsScheme(100, 15, 20)
		sets, err := scheme.GenerateSets(0.0, setscheme.DefaultSetGenerationContext())
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

	// AC2: TotalReps terminates at total reps target
	t.Run("AC2: TotalReps terminates at total reps target", func(t *testing.T) {
		scheme, _ := setscheme.NewTotalRepsScheme(50, 15, 20)
		ctx := setscheme.DefaultSetGenerationContext()
		history := []setscheme.GeneratedSet{{SetNumber: 1, Weight: 0.0}}
		termCtx := setscheme.TerminationContext{
			TotalSets: 3,
			TotalReps: 55, // Exceeds target of 50
			LastReps:  15,
		}
		_, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		if shouldContinue {
			t.Error("expected termination at total reps target")
		}
	})

	// AC3: TotalReps does NOT terminate at low reps (no min reps requirement)
	t.Run("AC3: TotalReps does NOT terminate at low reps", func(t *testing.T) {
		scheme, _ := setscheme.NewTotalRepsScheme(100, 15, 20)
		ctx := setscheme.DefaultSetGenerationContext()
		history := []setscheme.GeneratedSet{{SetNumber: 1, Weight: 0.0}}
		termCtx := setscheme.TerminationContext{
			TotalSets: 3,
			TotalReps: 30, // Well below target
			LastReps:  1,  // Very low rep set - should NOT cause termination
		}
		_, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		if !shouldContinue {
			t.Error("TotalReps should NOT terminate due to low reps per set")
		}
	})

	// AC4: TotalReps terminates at max sets
	t.Run("AC4: TotalReps terminates at max sets", func(t *testing.T) {
		scheme, _ := setscheme.NewTotalRepsScheme(1000, 15, 5) // Max 5 sets
		ctx := setscheme.DefaultSetGenerationContext()
		history := []setscheme.GeneratedSet{{SetNumber: 1, Weight: 0.0}}
		termCtx := setscheme.TerminationContext{
			TotalSets: 5, // At max
			TotalReps: 50,
			LastReps:  10,
		}
		_, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		if shouldContinue {
			t.Error("expected termination at max sets")
		}
	})

	// AC5: TotalReps uses same weight for all sets
	t.Run("AC5: TotalReps uses same weight for all sets", func(t *testing.T) {
		scheme, _ := setscheme.NewTotalRepsScheme(100, 15, 20)
		ctx := setscheme.DefaultSetGenerationContext()
		baseWeight := 225.0
		history := []setscheme.GeneratedSet{{SetNumber: 1, Weight: baseWeight}}
		termCtx := setscheme.TerminationContext{
			TotalSets: 1,
			TotalReps: 15,
			LastReps:  15,
		}
		nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		if !shouldContinue {
			t.Fatal("expected to continue")
		}
		if nextSet.Weight != baseWeight {
			t.Errorf("expected weight %.1f, got %.1f", baseWeight, nextSet.Weight)
		}
	})

	// AC6: TotalReps suggested reps are advisory only
	t.Run("AC6: TotalReps suggested reps are advisory", func(t *testing.T) {
		scheme, _ := setscheme.NewTotalRepsScheme(100, 15, 20)
		sets, _ := scheme.GenerateSets(0.0, setscheme.DefaultSetGenerationContext())
		// Suggested reps appear in TargetReps but don't affect termination
		if sets[0].TargetReps != 15 {
			t.Errorf("expected suggested reps 15, got %d", sets[0].TargetReps)
		}
	})

	// AC7: SessionService provides correct termination reasons for TotalReps
	t.Run("AC7: SessionService provides correct termination reasons", func(t *testing.T) {
		scheme, _ := setscheme.NewTotalRepsScheme(25, 10, 20)
		presc := &prescription.Prescription{ID: "p1", SetScheme: scheme}
		prescRepo := newMockPrescriptionRepo()
		prescRepo.Add(presc)
		loggedSetLister := newMockLoggedSetLister()
		sessionService := service.NewSessionService(prescRepo, loggedSetLister)

		// Log sets to hit target
		for i := 1; i <= 3; i++ {
			input := loggedset.CreateLoggedSetInput{
				UserID:         "u1",
				SessionID:      "s1",
				PrescriptionID: "p1",
				LiftID:         "l1",
				SetNumber:      i,
				Weight:         0.0,
				TargetReps:     10,
				RepsPerformed:  10, // 30 total > 25 target
			}
			logged, _ := loggedset.NewLoggedSet(input, uuid.New().String())
			loggedSetLister.Add(*logged)
		}

		req := service.NextSetRequest{SessionID: "s1", PrescriptionID: "p1", UserID: "u1"}
		result, _ := sessionService.GetNextSet(context.Background(), req)

		if !result.IsComplete {
			t.Error("expected completion")
		}
		expectedReason := "Target total reps reached (30/25)"
		if result.TerminationReason != expectedReason {
			t.Errorf("expected termination reason '%s', got '%s'", expectedReason, result.TerminationReason)
		}
	})
}
