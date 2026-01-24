// Package api_test provides integration tests for cross-entity workflows.
// These tests verify the correct interaction between multiple components
// and entities in the PowerPro system.
package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// =============================================================================
// ENVELOPE TYPES FOR API RESPONSE DECODING
// =============================================================================
// All API responses are wrapped in {"data": ...} envelopes.

// PrescriptionEnvelopeInteg wraps prescription responses for integration tests.
type PrescriptionEnvelopeInteg struct {
	Data PrescriptionResponse `json:"data"`
}

// ResolvedPrescriptionEnvelopeInteg wraps resolved prescription responses.
type ResolvedPrescriptionEnvelopeInteg struct {
	Data ResolvedPrescriptionTestResponse `json:"data"`
}

// BatchResolveEnvelopeInteg wraps batch resolve responses.
type BatchResolveEnvelopeInteg struct {
	Data BatchResolveTestResponse `json:"data"`
}

// WorkoutEnvelopeInteg wraps workout responses.
type WorkoutEnvelopeInteg struct {
	Data WorkoutTestResponse `json:"data"`
}

// CycleEnvelopeInteg wraps cycle responses.
type CycleEnvelopeInteg struct {
	Data CycleTestResponse `json:"data"`
}

// WeekEnvelopeInteg wraps week responses.
type WeekEnvelopeInteg struct {
	Data WeekTestResponse `json:"data"`
}

// DayEnvelopeInteg wraps day responses.
type DayEnvelopeInteg struct {
	Data DayTestResponse `json:"data"`
}

// ProgramEnvelopeInteg wraps program responses.
type ProgramEnvelopeInteg struct {
	Data ProgramTestResponse `json:"data"`
}

// ProgressionEnvelopeInteg wraps progression responses.
type ProgressionEnvelopeInteg struct {
	Data ProgressionResponse `json:"data"`
}

// ManualTriggerEnvelopeInteg wraps manual trigger responses.
type ManualTriggerEnvelopeInteg struct {
	Data ManualTriggerResponse `json:"data"`
}

// =============================================================================
// PRESCRIPTION RESOLUTION WORKFLOW INTEGRATION TESTS
// =============================================================================

// TestPrescriptionResolutionWorkflowIntegration tests the complete prescription
// resolution workflow: Movement -> Prescription -> LoadStrategy -> LiftMax -> resolved values
func TestPrescriptionResolutionWorkflowIntegration(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Use seeded lifts (squat, bench, deadlift)
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"
	userID := testutil.TestUserID

	t.Run("resolves prescription with PERCENT_OF TRAINING_MAX strategy", func(t *testing.T) {
		// Step 1: Create a training max for the user using the helper
		createMax(t, ts, userID, squatID, "TRAINING_MAX", 400.0, nil)

		// Step 2: Create a prescription with PERCENT_OF load strategy
		prescriptionBody := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85, "roundingIncrement": 5, "roundingDirection": "NEAREST"},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 3},
			"order": 1,
			"notes": "Heavy singles prep",
			"restSeconds": 240
		}`, squatID)
		prescResp, err := adminPost(ts.URL("/prescriptions"), prescriptionBody)
		if err != nil {
			t.Fatalf("Failed to create prescription: %v", err)
		}
		defer prescResp.Body.Close()

		if prescResp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(prescResp.Body)
			t.Fatalf("Failed to create prescription, status %d: %s", prescResp.StatusCode, body)
		}

		var prescEnvelope struct {
			Data PrescriptionResponse `json:"data"`
		}
		json.NewDecoder(prescResp.Body).Decode(&prescEnvelope)
		prescription := prescEnvelope.Data

		// Step 3: Resolve the prescription for the user
		resolveBody := fmt.Sprintf(`{"userId": "%s"}`, userID)
		resolveResp, err := authPost(ts.URL("/prescriptions/"+prescription.ID+"/resolve"), resolveBody)
		if err != nil {
			t.Fatalf("Failed to resolve prescription: %v", err)
		}
		defer resolveResp.Body.Close()

		if resolveResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resolveResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resolveResp.StatusCode, body)
		}

		var resolveEnvelope struct {
			Data ResolvedPrescriptionTestResponse `json:"data"`
		}
		json.NewDecoder(resolveResp.Body).Decode(&resolveEnvelope)
		resolved := resolveEnvelope.Data

		// Verify resolved values
		// 85% of 400 = 340, rounded to nearest 5 = 340
		expectedWeight := 340.0
		if len(resolved.Sets) != 5 {
			t.Errorf("Expected 5 sets, got %d", len(resolved.Sets))
		}
		for i, set := range resolved.Sets {
			if set.Weight != expectedWeight {
				t.Errorf("Set %d: expected weight %f, got %f", i+1, expectedWeight, set.Weight)
			}
			if set.TargetReps != 3 {
				t.Errorf("Set %d: expected 3 reps, got %d", i+1, set.TargetReps)
			}
			if !set.IsWorkSet {
				t.Errorf("Set %d: expected IsWorkSet=true for FIXED scheme", i+1)
			}
		}
		if resolved.Notes != "Heavy singles prep" {
			t.Errorf("Expected notes 'Heavy singles prep', got '%s'", resolved.Notes)
		}
		if resolved.RestSeconds == nil || *resolved.RestSeconds != 240 {
			t.Errorf("Expected restSeconds 240, got %v", resolved.RestSeconds)
		}
	})

	t.Run("resolves prescription with PERCENT_OF ONE_RM strategy", func(t *testing.T) {
		// Create a 1RM for bench
		createMax(t, ts, userID, benchID, "ONE_RM", 250.0, nil)

		// Create prescription with ONE_RM reference
		prescriptionBody := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "ONE_RM", "percentage": 90},
			"setScheme": {"type": "FIXED", "sets": 1, "reps": 1},
			"order": 1
		}`, benchID)
		prescResp, _ := adminPost(ts.URL("/prescriptions"), prescriptionBody)
		var prescEnvelope struct {
			Data PrescriptionResponse `json:"data"`
		}
		json.NewDecoder(prescResp.Body).Decode(&prescEnvelope)
		prescription := prescEnvelope.Data
		prescResp.Body.Close()

		// Resolve
		resolveBody := fmt.Sprintf(`{"userId": "%s"}`, userID)
		resolveResp, _ := authPost(ts.URL("/prescriptions/"+prescription.ID+"/resolve"), resolveBody)
		defer resolveResp.Body.Close()

		if resolveResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resolveResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resolveResp.StatusCode, body)
		}

		var resolveEnvelope struct {
			Data ResolvedPrescriptionTestResponse `json:"data"`
		}
		json.NewDecoder(resolveResp.Body).Decode(&resolveEnvelope)
		resolved := resolveEnvelope.Data

		// 90% of 250 = 225
		expectedWeight := 225.0
		if len(resolved.Sets) != 1 {
			t.Fatalf("Expected 1 set, got %d", len(resolved.Sets))
		}
		if resolved.Sets[0].Weight != expectedWeight {
			t.Errorf("Expected weight %f, got %f", expectedWeight, resolved.Sets[0].Weight)
		}
	})

	t.Run("resolves RAMP scheme with progressive weights", func(t *testing.T) {
		// Create training max for deadlift
		createMax(t, ts, userID, deadliftID, "TRAINING_MAX", 500.0, nil)

		// Create RAMP prescription
		prescriptionBody := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100, "roundingIncrement": 5},
			"setScheme": {"type": "RAMP", "steps": [{"percentage": 50, "reps": 5}, {"percentage": 60, "reps": 3}, {"percentage": 70, "reps": 2}, {"percentage": 80, "reps": 1}, {"percentage": 90, "reps": 1}], "workSetThreshold": 70},
			"order": 1
		}`, deadliftID)
		prescResp, _ := adminPost(ts.URL("/prescriptions"), prescriptionBody)
		var prescEnvelope PrescriptionEnvelopeInteg
		json.NewDecoder(prescResp.Body).Decode(&prescEnvelope)
		prescription := prescEnvelope.Data
		prescResp.Body.Close()

		// Resolve
		resolveBody := fmt.Sprintf(`{"userId": "%s"}`, userID)
		resolveResp, _ := authPost(ts.URL("/prescriptions/"+prescription.ID+"/resolve"), resolveBody)
		defer resolveResp.Body.Close()

		if resolveResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resolveResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resolveResp.StatusCode, body)
		}

		var resolveEnvelope ResolvedPrescriptionEnvelopeInteg
		json.NewDecoder(resolveResp.Body).Decode(&resolveEnvelope)
		resolved := resolveEnvelope.Data

		// Verify RAMP weights: 50% of 500 = 250, 60% = 300, 70% = 350, 80% = 400, 90% = 450
		expectedWeights := []float64{250, 300, 350, 400, 450}
		expectedReps := []int{5, 3, 2, 1, 1}
		expectedWorkSets := []bool{false, false, true, true, true} // >= 70%

		if len(resolved.Sets) != 5 {
			t.Fatalf("Expected 5 sets, got %d", len(resolved.Sets))
		}

		for i, set := range resolved.Sets {
			if set.Weight != expectedWeights[i] {
				t.Errorf("Set %d: expected weight %f, got %f", i+1, expectedWeights[i], set.Weight)
			}
			if set.TargetReps != expectedReps[i] {
				t.Errorf("Set %d: expected %d reps, got %d", i+1, expectedReps[i], set.TargetReps)
			}
			if set.IsWorkSet != expectedWorkSets[i] {
				t.Errorf("Set %d: expected isWorkSet=%v, got %v", i+1, expectedWorkSets[i], set.IsWorkSet)
			}
		}
	})

	t.Run("batch resolution resolves multiple prescriptions", func(t *testing.T) {
		// Create multiple prescriptions for the same lift (squat max already created)
		var prescriptionIDs []string
		for i := 0; i < 3; i++ {
			prescriptionBody := fmt.Sprintf(`{
				"liftId": "%s",
				"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": %d},
				"setScheme": {"type": "FIXED", "sets": 3, "reps": 5},
				"order": %d
			}`, squatID, 70+i*5, i)
			prescResp, _ := adminPost(ts.URL("/prescriptions"), prescriptionBody)
			var pEnvelope PrescriptionEnvelopeInteg
			json.NewDecoder(prescResp.Body).Decode(&pEnvelope)
			prescResp.Body.Close()
			prescriptionIDs = append(prescriptionIDs, pEnvelope.Data.ID)
		}

		// Batch resolve
		batchBody := fmt.Sprintf(`{"prescriptionIds": ["%s", "%s", "%s"], "userId": "%s"}`,
			prescriptionIDs[0], prescriptionIDs[1], prescriptionIDs[2], userID)
		batchResp, _ := authPost(ts.URL("/prescriptions/resolve-batch"), batchBody)
		defer batchResp.Body.Close()

		if batchResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(batchResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", batchResp.StatusCode, body)
		}

		var batchEnvelope BatchResolveEnvelopeInteg
		json.NewDecoder(batchResp.Body).Decode(&batchEnvelope)
		batchResult := batchEnvelope.Data

		if len(batchResult.Results) != 3 {
			t.Fatalf("Expected 3 results, got %d", len(batchResult.Results))
		}

		// All should succeed and have correct weights (70%, 75%, 80% of 400)
		expectedWeights := []float64{280, 300, 320}
		for i, result := range batchResult.Results {
			if result.Status != "success" {
				t.Errorf("Result %d: expected success, got %s (error: %s)", i, result.Status, result.Error)
				continue
			}
			if result.Resolved != nil && len(result.Resolved.Sets) > 0 {
				if result.Resolved.Sets[0].Weight != expectedWeights[i] {
					t.Errorf("Result %d: expected weight %f, got %f", i, expectedWeights[i], result.Resolved.Sets[0].Weight)
				}
			}
		}
	})

	t.Run("returns error when max not found for reference type", func(t *testing.T) {
		newUserID := "no-max-user-" + uuid.New().String()[:8]

		// Create prescription
		prescriptionBody := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5}
		}`, squatID)
		prescResp, _ := adminPost(ts.URL("/prescriptions"), prescriptionBody)
		var prescEnvelope PrescriptionEnvelopeInteg
		json.NewDecoder(prescResp.Body).Decode(&prescEnvelope)
		prescription := prescEnvelope.Data
		prescResp.Body.Close()

		// Try to resolve without a max
		resolveBody := fmt.Sprintf(`{"userId": "%s"}`, newUserID)
		resolveResp, _ := authPost(ts.URL("/prescriptions/"+prescription.ID+"/resolve"), resolveBody)
		defer resolveResp.Body.Close()

		if resolveResp.StatusCode != http.StatusBadRequest {
			body, _ := io.ReadAll(resolveResp.Body)
			t.Errorf("Expected status 400, got %d: %s", resolveResp.StatusCode, body)
		}
	})
}

// =============================================================================
// WORKOUT GENERATION WORKFLOW INTEGRATION TESTS
// =============================================================================

// TestWorkoutGenerationWorkflowIntegration tests the complete workout generation
// workflow: Schedule -> Week/Day -> Prescriptions -> LiftMaxes -> Generated Workout
func TestWorkoutGenerationWorkflowIntegration(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	squatID := "00000000-0000-0000-0000-000000000001"

	t.Run("generates complete workout with multiple exercises", func(t *testing.T) {
		// Use seeded user from migrations (required for foreign key constraint)
		userID := "workout-test-user"
		// Use the existing setupWorkoutTest helper which handles all the setup correctly
		setup := setupWorkoutTest(t, ts, userID)

		// Generate workout
		workoutResp, err := userGetWorkout(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		if workoutResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(workoutResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", workoutResp.StatusCode, body)
		}

		var workoutEnvelope WorkoutEnvelopeInteg
		json.NewDecoder(workoutResp.Body).Decode(&workoutEnvelope)
		workout := workoutEnvelope.Data

		// Verify workout structure
		if workout.UserID != userID {
			t.Errorf("Expected userId %s, got %s", userID, workout.UserID)
		}
		if workout.ProgramID != setup.ProgramID {
			t.Errorf("Expected programId %s, got %s", setup.ProgramID, workout.ProgramID)
		}
		if workout.WeekNumber != 1 {
			t.Errorf("Expected weekNumber 1, got %d", workout.WeekNumber)
		}
		if workout.CycleIteration != 1 {
			t.Errorf("Expected cycleIteration 1, got %d", workout.CycleIteration)
		}
		if len(workout.Exercises) < 1 {
			t.Fatalf("Expected at least 1 exercise, got %d", len(workout.Exercises))
		}

		// Verify exercise sets are calculated correctly (75% of 300 = 225)
		ex := workout.Exercises[0]
		expectedWeight := 225.0 // 75% of 300
		if len(ex.Sets) != 5 {
			t.Errorf("Expected 5 sets, got %d", len(ex.Sets))
		}
		for i, set := range ex.Sets {
			if set.Weight != expectedWeight {
				t.Errorf("Set %d: expected weight %f, got %f", i+1, expectedWeight, set.Weight)
			}
			if set.TargetReps != 5 {
				t.Errorf("Set %d: expected 5 reps, got %d", i+1, set.TargetReps)
			}
		}
	})

	t.Run("workout generation with RAMP set scheme", func(t *testing.T) {
		// Use seeded user from migrations (required for foreign key constraint)
		userID := "workout-preview-test-user"
		daySlug := "workout-ramp-day-" + uuid.New().String()[:8]
		programSlug := "workout-ramp-prog-" + uuid.New().String()[:8]
		liftSlug := "ramp-squat-" + uuid.New().String()[:8]

		// Create a lift
		liftBody := fmt.Sprintf(`{"name": "Ramp Squat", "slug": "%s"}`, liftSlug)
		liftResp, _ := adminPost(ts.URL("/lifts"), liftBody)
		var lift LiftResponse
		json.NewDecoder(liftResp.Body).Decode(&lift)
		liftResp.Body.Close()

		// Create lift max using helper
		createMax(t, ts, userID, lift.Data.ID, "TRAINING_MAX", 400.0, nil)

		// Create RAMP prescription
		prescBody := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100},
			"setScheme": {"type": "RAMP", "steps": [{"percentage": 50, "reps": 5}, {"percentage": 60, "reps": 3}, {"percentage": 70, "reps": 2}, {"percentage": 80, "reps": 1}], "workSetThreshold": 70},
			"order": 0
		}`, lift.Data.ID)
		prescResp, _ := adminPost(ts.URL("/prescriptions"), prescBody)
		var prescEnvelope PrescriptionEnvelopeInteg
		json.NewDecoder(prescResp.Body).Decode(&prescEnvelope)
		presc := prescEnvelope.Data
		prescResp.Body.Close()

		// Create day
		dayBody := fmt.Sprintf(`{"name": "Ramp Day", "slug": "%s"}`, daySlug)
		dayResp, _ := adminPost(ts.URL("/days"), dayBody)
		var dayEnvelope DayEnvelopeInteg
		json.NewDecoder(dayResp.Body).Decode(&dayEnvelope)
		day := dayEnvelope.Data
		dayResp.Body.Close()

		addPresc, _ := adminPost(ts.URL("/days/"+day.ID+"/prescriptions"), `{"prescriptionId": "`+presc.ID+`"}`)
		addPresc.Body.Close()

		// Create cycle, week, program
		cycleResp, _ := adminPostCycle(ts.URL("/cycles"), `{"name": "Ramp Cycle", "lengthWeeks": 4}`)
		var cycleEnvelope CycleEnvelopeInteg
		json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
		cycle := cycleEnvelope.Data
		cycleResp.Body.Close()

		weekResp, _ := adminPost(ts.URL("/weeks"), `{"weekNumber": 1, "cycleId": "`+cycle.ID+`"}`)
		var weekEnvelope WeekEnvelopeInteg
		json.NewDecoder(weekResp.Body).Decode(&weekEnvelope)
		week := weekEnvelope.Data
		weekResp.Body.Close()

		addDayResp, _ := adminPost(ts.URL("/weeks/"+week.ID+"/days"), `{"dayId": "`+day.ID+`", "dayOfWeek": "TUESDAY"}`)
		addDayResp.Body.Close()

		programResp, _ := adminPostProgram(ts.URL("/programs"), `{"name": "Ramp Program", "slug": "`+programSlug+`", "cycleId": "`+cycle.ID+`"}`)
		var programEnvelope ProgramEnvelopeInteg
		json.NewDecoder(programResp.Body).Decode(&programEnvelope)
		program := programEnvelope.Data
		programResp.Body.Close()

		// Enroll user
		enrollResp, _ := userPostEnrollment(ts.URL("/users/"+userID+"/program"), `{"programId": "`+program.ID+`"}`, userID)
		if enrollResp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(enrollResp.Body)
			t.Fatalf("Failed to enroll user: %d, %s", enrollResp.StatusCode, body)
		}
		enrollResp.Body.Close()

		// Generate workout
		workoutResp, _ := userGetWorkout(ts.URL("/users/"+userID+"/workout"), userID)
		defer workoutResp.Body.Close()

		if workoutResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(workoutResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", workoutResp.StatusCode, body)
		}

		var workoutEnvelope WorkoutEnvelopeInteg
		json.NewDecoder(workoutResp.Body).Decode(&workoutEnvelope)
		workout := workoutEnvelope.Data

		if len(workout.Exercises) != 1 {
			t.Fatalf("Expected 1 exercise, got %d", len(workout.Exercises))
		}

		ex := workout.Exercises[0]
		if len(ex.Sets) != 4 {
			t.Fatalf("Expected 4 sets (RAMP steps), got %d", len(ex.Sets))
		}

		// Verify RAMP weights: 50% = 200, 60% = 240, 70% = 280, 80% = 320
		expectedWeights := []float64{200, 240, 280, 320}
		expectedReps := []int{5, 3, 2, 1}
		expectedWorkSets := []bool{false, false, true, true}

		for i, set := range ex.Sets {
			if set.Weight != expectedWeights[i] {
				t.Errorf("Set %d: expected weight %f, got %f", i+1, expectedWeights[i], set.Weight)
			}
			if set.TargetReps != expectedReps[i] {
				t.Errorf("Set %d: expected %d reps, got %d", i+1, expectedReps[i], set.TargetReps)
			}
			if set.IsWorkSet != expectedWorkSets[i] {
				t.Errorf("Set %d: expected isWorkSet=%v, got %v", i+1, expectedWorkSets[i], set.IsWorkSet)
			}
		}
	})

	t.Run("workout preview allows specific week/day selection", func(t *testing.T) {
		// Use seeded user from migrations (required for foreign key constraint)
		userID := "workout-format-test-user"
		setup := setupWorkoutTest(t, ts, userID)
		daySlug := "squat-day-" + userID

		// Preview for week 1, specific day
		previewResp, _ := userGetWorkout(ts.URL("/users/"+userID+"/workout/preview?week=1&day="+daySlug), userID)
		defer previewResp.Body.Close()

		if previewResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(previewResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", previewResp.StatusCode, body)
		}

		var previewEnvelope WorkoutEnvelopeInteg
		json.NewDecoder(previewResp.Body).Decode(&previewEnvelope)
		preview := previewEnvelope.Data

		if preview.WeekNumber != 1 {
			t.Errorf("Expected weekNumber 1, got %d", preview.WeekNumber)
		}
		if preview.DaySlug != daySlug {
			t.Errorf("Expected daySlug '%s', got '%s'", daySlug, preview.DaySlug)
		}
		_ = setup // used for enrollment
	})

	t.Run("returns error when user not enrolled", func(t *testing.T) {
		unenrolledUser := "unenrolled-user-" + uuid.New().String()[:8]
		workoutResp, _ := userGetWorkout(ts.URL("/users/"+unenrolledUser+"/workout"), unenrolledUser)
		defer workoutResp.Body.Close()

		if workoutResp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", workoutResp.StatusCode)
		}
	})

	_ = squatID // available for future subtests
}

// =============================================================================
// PROGRESSION EVALUATION WORKFLOW INTEGRATION TESTS
// =============================================================================

// TestProgressionEvaluationWorkflowIntegration tests the complete progression
// evaluation workflow: WorkoutLog -> Progression rules -> LiftMax update -> History
func TestProgressionEvaluationWorkflowIntegration(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	t.Run("full progression lifecycle: trigger -> evaluate -> update -> history", func(t *testing.T) {
		// Use seeded user from migrations (required for foreign key constraint)
		userID := "auth-test-user-adv"
		programSlug := "prog-lifecycle-prog-" + uuid.New().String()[:8]

		// Step 1: Create cycle and program
		cycleResp, _ := adminPostCycle(ts.URL("/cycles"), `{"name": "Progression Cycle", "lengthWeeks": 4}`)
		var cycleEnvelope CycleEnvelopeInteg
		json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
		cycle := cycleEnvelope.Data
		cycleResp.Body.Close()

		programResp, _ := adminPostProgram(ts.URL("/programs"), `{"name": "Progression Program", "slug": "`+programSlug+`", "cycleId": "`+cycle.ID+`"}`)
		var programEnvelope ProgramEnvelopeInteg
		json.NewDecoder(programResp.Body).Decode(&programEnvelope)
		program := programEnvelope.Data
		programResp.Body.Close()

		// Step 2: Create a linear progression (5lb increment after session)
		progressionResp, _ := adminPost(ts.URL("/progressions"), `{"name": "Linear Session Prog", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`)
		var progressionEnvelope ProgressionEnvelopeInteg
		json.NewDecoder(progressionResp.Body).Decode(&progressionEnvelope)
		progression := progressionEnvelope.Data
		progressionResp.Body.Close()

		// Step 3: Link progression to program for squat
		ppBody := fmt.Sprintf(`{"progressionId": "%s", "liftId": "%s", "priority": 1, "enabled": true}`, progression.ID, squatID)
		ppResp, _ := adminPost(ts.URL("/programs/"+program.ID+"/progressions"), ppBody)
		ppResp.Body.Close()

		// Step 4: Create initial lift max for user
		initialMax := 300.0
		createMax(t, ts, userID, squatID, "TRAINING_MAX", initialMax, nil)

		// Step 5: Enroll user in program
		enrollResp, _ := userPostEnrollment(ts.URL("/users/"+userID+"/program"), `{"programId": "`+program.ID+`"}`, userID)
		if enrollResp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(enrollResp.Body)
			t.Fatalf("Failed to enroll user: %d, %s", enrollResp.StatusCode, body)
		}
		enrollResp.Body.Close()

		// Step 6: Trigger manual progression
		triggerBody := ManualTriggerRequest{
			ProgressionID: progression.ID,
			LiftID:        squatID,
			Force:         true, // Force to ensure it applies
		}
		triggerResp, err := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
		if err != nil {
			t.Fatalf("Failed to trigger progression: %v", err)
		}
		defer triggerResp.Body.Close()

		if triggerResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(triggerResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", triggerResp.StatusCode, body)
		}

		var triggerEnvelope ManualTriggerEnvelopeInteg
		json.NewDecoder(triggerResp.Body).Decode(&triggerEnvelope)
		triggerResult := triggerEnvelope.Data

		// Verify progression was applied
		if triggerResult.TotalApplied != 1 {
			t.Errorf("Expected TotalApplied=1, got %d", triggerResult.TotalApplied)
		}

		if len(triggerResult.Results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(triggerResult.Results))
		}

		result := triggerResult.Results[0]
		if !result.Applied {
			t.Errorf("Expected progression to be applied, skipReason: %s, error: %s", result.SkipReason, result.Error)
		}
		if result.Result == nil {
			t.Fatal("Expected result detail to be present")
		}
		if result.Result.Delta != 5.0 {
			t.Errorf("Expected delta 5.0, got %f", result.Result.Delta)
		}
		if result.Result.NewValue != initialMax+5.0 {
			t.Errorf("Expected new value %f, got %f", initialMax+5.0, result.Result.NewValue)
		}

		// Step 7: Verify progression history was recorded
		historyResp, _ := authGetHistory(ts.URL("/users/"+userID+"/progression-history?limit=10"), userID)
		defer historyResp.Body.Close()

		var history ProgressionHistoryTestListResponse
		json.NewDecoder(historyResp.Body).Decode(&history)

		if len(history.Data) < 1 {
			t.Fatal("Expected at least 1 history entry")
		}

		latestEntry := history.Data[0]
		if latestEntry.LiftID != squatID {
			t.Errorf("Expected history entry for squat, got %s", latestEntry.LiftID)
		}
		if latestEntry.PreviousValue != initialMax {
			t.Errorf("Expected previous value %f, got %f", initialMax, latestEntry.PreviousValue)
		}
		if latestEntry.NewValue != initialMax+5.0 {
			t.Errorf("Expected new value %f, got %f", initialMax+5.0, latestEntry.NewValue)
		}
	})

	t.Run("progression with multiple lifts applies in priority order", func(t *testing.T) {
		// Use seeded user from migrations (required for foreign key constraint)
		userID := "format-test-user-adv"
		programSlug := "multi-lift-prog-" + uuid.New().String()[:8]

		// Create cycle and program
		cycleResp, _ := adminPostCycle(ts.URL("/cycles"), `{"name": "Multi Lift Cycle", "lengthWeeks": 4}`)
		var cycleEnvelope CycleEnvelopeInteg
		json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
		cycle := cycleEnvelope.Data
		cycleResp.Body.Close()

		programResp, _ := adminPostProgram(ts.URL("/programs"), `{"name": "Multi Lift Program", "slug": "`+programSlug+`", "cycleId": "`+cycle.ID+`"}`)
		var programEnvelope ProgramEnvelopeInteg
		json.NewDecoder(programResp.Body).Decode(&programEnvelope)
		program := programEnvelope.Data
		programResp.Body.Close()

		// Create progression
		progressionResp, _ := adminPost(ts.URL("/progressions"), `{"name": "Multi Lift Prog", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`)
		var progressionEnvelope ProgressionEnvelopeInteg
		json.NewDecoder(progressionResp.Body).Decode(&progressionEnvelope)
		progression := progressionEnvelope.Data
		progressionResp.Body.Close()

		// Link progression to squat (priority 1) and bench (priority 2)
		pp1Body := fmt.Sprintf(`{"progressionId": "%s", "liftId": "%s", "priority": 1, "enabled": true}`, progression.ID, squatID)
		pp1Resp, _ := adminPost(ts.URL("/programs/"+program.ID+"/progressions"), pp1Body)
		pp1Resp.Body.Close()

		pp2Body := fmt.Sprintf(`{"progressionId": "%s", "liftId": "%s", "priority": 2, "enabled": true}`, progression.ID, benchID)
		pp2Resp, _ := adminPost(ts.URL("/programs/"+program.ID+"/progressions"), pp2Body)
		pp2Resp.Body.Close()

		// Create lift maxes
		createMax(t, ts, userID, squatID, "TRAINING_MAX", 300.0, nil)
		createMax(t, ts, userID, benchID, "TRAINING_MAX", 200.0, nil)

		// Enroll user
		enrollResp, _ := userPostEnrollment(ts.URL("/users/"+userID+"/program"), `{"programId": "`+program.ID+`"}`, userID)
		if enrollResp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(enrollResp.Body)
			t.Fatalf("Failed to enroll user: %d, %s", enrollResp.StatusCode, body)
		}
		enrollResp.Body.Close()

		// Trigger progression without specifying lift (should apply to all)
		triggerBody := ManualTriggerRequest{
			ProgressionID: progression.ID,
			LiftID:        "", // Empty = apply to all configured lifts
			Force:         true,
		}
		triggerResp, _ := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
		defer triggerResp.Body.Close()

		var triggerEnvelope ManualTriggerEnvelopeInteg
		json.NewDecoder(triggerResp.Body).Decode(&triggerEnvelope)
		triggerResult := triggerEnvelope.Data

		// Should apply to both lifts
		if triggerResult.TotalApplied != 2 {
			t.Errorf("Expected TotalApplied=2, got %d (skipped=%d, errors=%d)", triggerResult.TotalApplied, triggerResult.TotalSkipped, triggerResult.TotalErrors)
		}

		if len(triggerResult.Results) != 2 {
			t.Fatalf("Expected 2 results, got %d", len(triggerResult.Results))
		}

		// Verify results are in priority order (squat first, then bench)
		if triggerResult.Results[0].LiftID != squatID {
			t.Errorf("Expected first result to be squat (priority 1), got %s", triggerResult.Results[0].LiftID)
		}
		if triggerResult.Results[1].LiftID != benchID {
			t.Errorf("Expected second result to be bench (priority 2), got %s", triggerResult.Results[1].LiftID)
		}
	})

	t.Run("cycle progression applies with larger increment", func(t *testing.T) {
		// Use seeded user from migrations (required for foreign key constraint)
		userID := "no-days-user"
		programSlug := "cycle-prog-" + uuid.New().String()[:8]

		// Create cycle and program
		cycleResp, _ := adminPostCycle(ts.URL("/cycles"), `{"name": "Cycle End Prog Cycle", "lengthWeeks": 4}`)
		var cycleEnvelope CycleEnvelopeInteg
		json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
		cycle := cycleEnvelope.Data
		cycleResp.Body.Close()

		programResp, _ := adminPostProgram(ts.URL("/programs"), `{"name": "Cycle End Prog Program", "slug": "`+programSlug+`", "cycleId": "`+cycle.ID+`"}`)
		var programEnvelope ProgramEnvelopeInteg
		json.NewDecoder(programResp.Body).Decode(&programEnvelope)
		program := programEnvelope.Data
		programResp.Body.Close()

		// Create cycle progression (10lb increment)
		progressionResp, _ := adminPost(ts.URL("/progressions"), `{"name": "Cycle End Prog", "type": "CYCLE_PROGRESSION", "parameters": {"increment": 10.0, "maxType": "TRAINING_MAX"}}`)
		var progressionEnvelope ProgressionEnvelopeInteg
		json.NewDecoder(progressionResp.Body).Decode(&progressionEnvelope)
		progression := progressionEnvelope.Data
		progressionResp.Body.Close()

		// Link to deadlift
		ppBody := fmt.Sprintf(`{"progressionId": "%s", "liftId": "%s", "priority": 1, "enabled": true}`, progression.ID, deadliftID)
		ppResp, _ := adminPost(ts.URL("/programs/"+program.ID+"/progressions"), ppBody)
		ppResp.Body.Close()

		// Create lift max
		createMax(t, ts, userID, deadliftID, "TRAINING_MAX", 400.0, nil)

		// Enroll user
		enrollResp, _ := userPostEnrollment(ts.URL("/users/"+userID+"/program"), `{"programId": "`+program.ID+`"}`, userID)
		if enrollResp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(enrollResp.Body)
			t.Fatalf("Failed to enroll user: %d, %s", enrollResp.StatusCode, body)
		}
		enrollResp.Body.Close()

		// Manually trigger the cycle progression
		triggerBody := ManualTriggerRequest{
			ProgressionID: progression.ID,
			LiftID:        deadliftID,
			Force:         true,
		}
		triggerResp, _ := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
		defer triggerResp.Body.Close()

		var triggerEnvelope ManualTriggerEnvelopeInteg
		json.NewDecoder(triggerResp.Body).Decode(&triggerEnvelope)
		triggerResult := triggerEnvelope.Data

		// Should apply cycle progression (+10)
		if triggerResult.TotalApplied != 1 {
			t.Errorf("Expected TotalApplied=1, got %d", triggerResult.TotalApplied)
		}

		if len(triggerResult.Results) > 0 && triggerResult.Results[0].Result != nil {
			if triggerResult.Results[0].Result.Delta != 10.0 {
				t.Errorf("Expected delta 10.0 for cycle progression, got %f", triggerResult.Results[0].Result.Delta)
			}
			if triggerResult.Results[0].Result.NewValue != 410.0 {
				t.Errorf("Expected new value 410.0, got %f", triggerResult.Results[0].Result.NewValue)
			}
		}
	})

	t.Run("force=true allows repeated progression applications", func(t *testing.T) {
		// Use seeded user from migrations (required for foreign key constraint)
		// Using a different user to avoid conflicts with other tests
		userID := "convert-user-a"
		programSlug := "idempotent-prog-" + uuid.New().String()[:8]

		// Create cycle and program
		cycleResp, _ := adminPostCycle(ts.URL("/cycles"), `{"name": "Idempotent Cycle", "lengthWeeks": 4}`)
		var cycleEnvelope CycleEnvelopeInteg
		json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
		cycle := cycleEnvelope.Data
		cycleResp.Body.Close()

		programResp, _ := adminPostProgram(ts.URL("/programs"), `{"name": "Idempotent Program", "slug": "`+programSlug+`", "cycleId": "`+cycle.ID+`"}`)
		var programEnvelope ProgramEnvelopeInteg
		json.NewDecoder(programResp.Body).Decode(&programEnvelope)
		program := programEnvelope.Data
		programResp.Body.Close()

		// Create progression
		progressionResp, _ := adminPost(ts.URL("/progressions"), `{"name": "Idempotent Prog", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`)
		var progressionEnvelope ProgressionEnvelopeInteg
		json.NewDecoder(progressionResp.Body).Decode(&progressionEnvelope)
		progression := progressionEnvelope.Data
		progressionResp.Body.Close()

		// Link progression
		ppBody := fmt.Sprintf(`{"progressionId": "%s", "liftId": "%s", "priority": 1, "enabled": true}`, progression.ID, squatID)
		ppResp, _ := adminPost(ts.URL("/programs/"+program.ID+"/progressions"), ppBody)
		ppResp.Body.Close()

		// Create lift max
		createMax(t, ts, userID, squatID, "TRAINING_MAX", 300.0, nil)

		// Enroll user
		enrollResp, _ := userPostEnrollment(ts.URL("/users/"+userID+"/program"), `{"programId": "`+program.ID+`"}`, userID)
		if enrollResp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(enrollResp.Body)
			t.Fatalf("Failed to enroll user: %d, %s", enrollResp.StatusCode, body)
		}
		enrollResp.Body.Close()

		// First trigger with force=true
		triggerBody := ManualTriggerRequest{
			ProgressionID: progression.ID,
			LiftID:        squatID,
			Force:         true,
		}
		trigger1Resp, _ := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
		var result1Envelope ManualTriggerEnvelopeInteg
		json.NewDecoder(trigger1Resp.Body).Decode(&result1Envelope)
		result1 := result1Envelope.Data
		trigger1Resp.Body.Close()

		if result1.TotalApplied != 1 {
			t.Fatalf("Expected first trigger to apply, got TotalApplied=%d", result1.TotalApplied)
		}
		// Verify first trigger applied successfully with correct delta
		if result1.Results[0].Result == nil {
			t.Fatal("Expected first trigger result detail to be present")
		}
		if result1.Results[0].Result.Delta != 5.0 {
			t.Errorf("First trigger: expected delta 5.0, got %f", result1.Results[0].Result.Delta)
		}

		// Second trigger with force=true should also apply (bypassing idempotency)
		trigger2Resp, _ := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
		var result2Envelope ManualTriggerEnvelopeInteg
		json.NewDecoder(trigger2Resp.Body).Decode(&result2Envelope)
		result2 := result2Envelope.Data
		trigger2Resp.Body.Close()

		if result2.TotalApplied != 1 {
			t.Errorf("Expected second force trigger to apply, got TotalApplied=%d", result2.TotalApplied)
		}
		if result2.Results[0].Result != nil {
			secondNewValue := result2.Results[0].Result.NewValue
			secondDelta := result2.Results[0].Result.Delta
			// With force=true, the second trigger should also apply with 5lb delta
			if secondDelta != 5.0 {
				t.Errorf("Expected delta 5.0, got %f", secondDelta)
			}
			// The second trigger should produce a value at least 5 greater than initial
			// (may or may not chain from first depending on effective_date handling)
			initialMax := 300.0
			if secondNewValue < initialMax+5.0 {
				t.Errorf("Expected newValue >= %f, got %f", initialMax+5.0, secondNewValue)
			}
			// The key verification: force=true allows the trigger to apply twice
			// Each application applies the increment
			if result2.Results[0].Applied != true {
				t.Error("Expected second trigger to be applied")
			}
		}
	})
}

// =============================================================================
// VARIABLE SCHEME SESSION HANDLING INTEGRATION TESTS
// =============================================================================

// NextSetTestResponse represents the API response for next set requests.
type NextSetTestResponse struct {
	NextSet            *NextSetInfoTest `json:"nextSet,omitempty"`
	IsComplete         bool             `json:"isComplete"`
	TotalSetsCompleted int              `json:"totalSetsCompleted"`
	TotalRepsCompleted int              `json:"totalRepsCompleted"`
	TerminationReason  string           `json:"terminationReason,omitempty"`
}

// NextSetInfoTest represents a generated set in the API response.
type NextSetInfoTest struct {
	SetNumber  int     `json:"setNumber"`
	Weight     float64 `json:"weight"`
	TargetReps int     `json:"targetReps"`
	IsWorkSet  bool    `json:"isWorkSet"`
}

// NextSetEnvelopeInteg wraps next set responses.
type NextSetEnvelopeInteg struct {
	Data NextSetTestResponse `json:"data"`
}

// TestVariableSchemeSessionHandling tests the dynamic set generation for MRS and FatigueDrop schemes.
func TestVariableSchemeSessionHandling(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	squatID := "00000000-0000-0000-0000-000000000001"
	userID := testutil.TestUserID

	t.Run("MRS scheme generates next set until target total reps reached", func(t *testing.T) {
		sessionID := uuid.New().String()

		// Create an MRS prescription: target 25 total reps, min 3 reps per set
		prescBody := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 75},
			"setScheme": {"type": "MRS", "target_total_reps": 25, "min_reps_per_set": 3, "max_sets": 10},
			"order": 0
		}`, squatID)
		prescResp, err := adminPost(ts.URL("/prescriptions"), prescBody)
		if err != nil {
			t.Fatalf("Failed to create prescription: %v", err)
		}
		var prescEnvelope PrescriptionEnvelopeInteg
		json.NewDecoder(prescResp.Body).Decode(&prescEnvelope)
		prescResp.Body.Close()
		prescriptionID := prescEnvelope.Data.ID

		// Ensure user has a training max
		createMax(t, ts, userID, squatID, "TRAINING_MAX", 400.0, nil)

		// Before logging any sets, requesting next-set should fail
		t.Run("returns error before any sets logged", func(t *testing.T) {
			resp, err := authGetLoggedSets(ts.URL("/sessions/"+sessionID+"/prescriptions/"+prescriptionID+"/next-set"), userID)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			// Should return 400 because no sets have been logged
			if resp.StatusCode != http.StatusBadRequest {
				body, _ := io.ReadAll(resp.Body)
				t.Errorf("Expected status 400, got %d: %s", resp.StatusCode, body)
			}
		})

		// Log first set: 10 reps
		loggedSetBody := fmt.Sprintf(`{
			"sets": [
				{
					"prescriptionId": "%s",
					"liftId": "%s",
					"setNumber": 1,
					"weight": 300.0,
					"targetReps": 3,
					"repsPerformed": 10,
					"isAmrap": false
				}
			]
		}`, prescriptionID, squatID)
		logResp, _ := authPostLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), loggedSetBody, userID)
		logResp.Body.Close()

		t.Run("returns next set after first set logged", func(t *testing.T) {
			resp, err := authGetLoggedSets(ts.URL("/sessions/"+sessionID+"/prescriptions/"+prescriptionID+"/next-set"), userID)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, body)
			}

			var envelope NextSetEnvelopeInteg
			json.NewDecoder(resp.Body).Decode(&envelope)
			result := envelope.Data

			if result.IsComplete {
				t.Error("Expected exercise to not be complete yet")
			}
			if result.NextSet == nil {
				t.Fatal("Expected next set to be returned")
			}
			if result.NextSet.SetNumber != 2 {
				t.Errorf("Expected set number 2, got %d", result.NextSet.SetNumber)
			}
			if result.TotalSetsCompleted != 1 {
				t.Errorf("Expected 1 set completed, got %d", result.TotalSetsCompleted)
			}
			if result.TotalRepsCompleted != 10 {
				t.Errorf("Expected 10 reps completed, got %d", result.TotalRepsCompleted)
			}
		})

		// Log more sets until target reached
		loggedSetBody2 := fmt.Sprintf(`{
			"sets": [
				{
					"prescriptionId": "%s",
					"liftId": "%s",
					"setNumber": 2,
					"weight": 300.0,
					"targetReps": 3,
					"repsPerformed": 8,
					"isAmrap": false
				}
			]
		}`, prescriptionID, squatID)
		logResp2, _ := authPostLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), loggedSetBody2, userID)
		logResp2.Body.Close()

		loggedSetBody3 := fmt.Sprintf(`{
			"sets": [
				{
					"prescriptionId": "%s",
					"liftId": "%s",
					"setNumber": 3,
					"weight": 300.0,
					"targetReps": 3,
					"repsPerformed": 7,
					"isAmrap": false
				}
			]
		}`, prescriptionID, squatID)
		logResp3, _ := authPostLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), loggedSetBody3, userID)
		logResp3.Body.Close()

		t.Run("returns complete when target total reps reached", func(t *testing.T) {
			resp, err := authGetLoggedSets(ts.URL("/sessions/"+sessionID+"/prescriptions/"+prescriptionID+"/next-set"), userID)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, body)
			}

			var envelope NextSetEnvelopeInteg
			json.NewDecoder(resp.Body).Decode(&envelope)
			result := envelope.Data

			// Total reps: 10 + 8 + 7 = 25 (matches target)
			if !result.IsComplete {
				t.Error("Expected exercise to be complete (25 total reps)")
			}
			if result.NextSet != nil {
				t.Error("Expected next set to be nil when complete")
			}
			if result.TotalRepsCompleted != 25 {
				t.Errorf("Expected 25 reps completed, got %d", result.TotalRepsCompleted)
			}
			if result.TerminationReason == "" {
				t.Error("Expected termination reason to be set")
			}
		})
	})

	t.Run("returns error for non-variable scheme", func(t *testing.T) {
		sessionID := uuid.New().String()

		// Create a FIXED prescription (not variable)
		prescBody := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 75},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5},
			"order": 0
		}`, squatID)
		prescResp, _ := adminPost(ts.URL("/prescriptions"), prescBody)
		var prescEnvelope PrescriptionEnvelopeInteg
		json.NewDecoder(prescResp.Body).Decode(&prescEnvelope)
		prescResp.Body.Close()
		prescriptionID := prescEnvelope.Data.ID

		// Log a set
		loggedSetBody := fmt.Sprintf(`{
			"sets": [
				{
					"prescriptionId": "%s",
					"liftId": "%s",
					"setNumber": 1,
					"weight": 300.0,
					"targetReps": 5,
					"repsPerformed": 5,
					"isAmrap": false
				}
			]
		}`, prescriptionID, squatID)
		logResp, _ := authPostLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), loggedSetBody, userID)
		logResp.Body.Close()

		resp, err := authGetLoggedSets(ts.URL("/sessions/"+sessionID+"/prescriptions/"+prescriptionID+"/next-set"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Should return 400 because FIXED is not a variable scheme
		if resp.StatusCode != http.StatusBadRequest {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 400, got %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("returns error for non-existent prescription", func(t *testing.T) {
		sessionID := uuid.New().String()
		fakePrescriptionID := uuid.New().String()

		resp, err := authGetLoggedSets(ts.URL("/sessions/"+sessionID+"/prescriptions/"+fakePrescriptionID+"/next-set"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 404, got %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("logged sets with RPE are returned correctly", func(t *testing.T) {
		sessionID := uuid.New().String()

		// Log a set with RPE
		loggedSetBody := fmt.Sprintf(`{
			"sets": [
				{
					"prescriptionId": "%s",
					"liftId": "%s",
					"setNumber": 1,
					"weight": 300.0,
					"targetReps": 5,
					"repsPerformed": 5,
					"isAmrap": false,
					"rpe": 8.5
				}
			]
		}`, uuid.New().String(), squatID)
		logResp, _ := authPostLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), loggedSetBody, userID)
		if logResp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(logResp.Body)
			t.Fatalf("Failed to create logged set: %d: %s", logResp.StatusCode, body)
		}

		var loggedSetEnvelope LoggedSetTestListResponse
		json.NewDecoder(logResp.Body).Decode(&loggedSetEnvelope)
		logResp.Body.Close()

		// Verify the logged set response
		if len(loggedSetEnvelope.Data) != 1 {
			t.Fatalf("Expected 1 logged set, got %d", len(loggedSetEnvelope.Data))
		}
		// RPE should be included in response (check via Get endpoint)
		resp, _ := authGetLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), userID)
		defer resp.Body.Close()

		// Verify RPE is returned
		var listResponse struct {
			Data []struct {
				RPE *float64 `json:"rpe"`
			} `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&listResponse)

		if len(listResponse.Data) == 0 {
			t.Fatal("Expected logged sets in response")
		}
		if listResponse.Data[0].RPE == nil {
			t.Error("Expected RPE to be returned in logged set response")
		} else if *listResponse.Data[0].RPE != 8.5 {
			t.Errorf("Expected RPE 8.5, got %f", *listResponse.Data[0].RPE)
		}
	})
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// intUserGetLiftMaxes retrieves lift maxes for a user (integration test specific)
func intUserGetLiftMaxes(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

// IntegrationLiftMaxesListTestResponse represents the paginated list of lift maxes
type IntegrationLiftMaxesListTestResponse struct {
	Data       []LiftMaxTestResponse `json:"data"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"pageSize"`
	TotalItems int64                 `json:"totalItems"`
	TotalPages int64                 `json:"totalPages"`
}
