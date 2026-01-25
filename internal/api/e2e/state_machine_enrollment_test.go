package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// TestStateMachineEnrollmentTransitions tests enrollment state machine transitions.
// This test focuses on the high-level enrollment status transitions:
// - NONE -> ACTIVE (enroll)
// - ACTIVE -> BETWEEN_CYCLES (cycle completes via advanceWeek at final week)
// - BETWEEN_CYCLES -> ACTIVE (start new cycle)
// - ACTIVE -> QUIT (unenroll)
// - BETWEEN_CYCLES -> QUIT (unenroll while deciding)
// - Error: Can't start workout when BETWEEN_CYCLES
// - Error: Can't start new cycle when ACTIVE
func TestStateMachineEnrollmentTransitions(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Test-unique identifier
	testID := uuid.New().String()[:8]

	// Use a seeded test user
	userID := "workout-test-user"

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"

	// Create training maxes
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", 225.0)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", 135.0)

	// Create a simple 2-week program for testing state transitions
	// Week 1: Day A only (simple)
	// Week 2: Day A only (simple)

	// Create prescriptions
	squatPrescID := createPrescription(t, ts, squatID, 3, 5, 100.0, 0)
	benchPrescID := createPrescription(t, ts, benchID, 3, 5, 100.0, 1)

	// Create Day A
	dayASlug := "sm-day-a-" + testID
	dayABody := fmt.Sprintf(`{"name": "State Machine Day A", "slug": "%s"}`, dayASlug)
	dayAResp, err := adminPost(ts.URL("/days"), dayABody)
	if err != nil {
		t.Fatalf("Failed to create day: %v", err)
	}
	var dayAEnvelope DayResponse
	json.NewDecoder(dayAResp.Body).Decode(&dayAEnvelope)
	dayAResp.Body.Close()
	dayAID := dayAEnvelope.Data.ID

	// Add prescriptions to Day A
	addPrescToDay(t, ts, dayAID, squatPrescID)
	addPrescToDay(t, ts, dayAID, benchPrescID)

	// Create 2-week cycle
	cycleName := "SM Test Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 2}`, cycleName)
	cycleResp, err := adminPost(ts.URL("/cycles"), cycleBody)
	if err != nil {
		t.Fatalf("Failed to create cycle: %v", err)
	}
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID := cycleEnvelope.Data.ID

	// Create week 1
	week1Body := fmt.Sprintf(`{"weekNumber": 1, "cycleId": "%s"}`, cycleID)
	week1Resp, _ := adminPost(ts.URL("/weeks"), week1Body)
	var week1Envelope WeekResponse
	json.NewDecoder(week1Resp.Body).Decode(&week1Envelope)
	week1Resp.Body.Close()
	week1ID := week1Envelope.Data.ID

	// Create week 2
	week2Body := fmt.Sprintf(`{"weekNumber": 2, "cycleId": "%s"}`, cycleID)
	week2Resp, _ := adminPost(ts.URL("/weeks"), week2Body)
	var week2Envelope WeekResponse
	json.NewDecoder(week2Resp.Body).Decode(&week2Envelope)
	week2Resp.Body.Close()
	week2ID := week2Envelope.Data.ID

	// Add day to weeks
	addDayToWeek(t, ts, week1ID, dayAID, "MONDAY")
	addDayToWeek(t, ts, week2ID, dayAID, "MONDAY")

	// Create program
	programSlug := "sm-test-program-" + testID
	programBody := fmt.Sprintf(`{"name": "State Machine Test Program", "slug": "%s", "cycleId": "%s"}`, programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	// =============================================================================
	// TEST 1: NONE -> ACTIVE (enroll)
	// =============================================================================
	t.Run("Test NONE to ACTIVE (enroll)", func(t *testing.T) {
		enrollment := enrollUser(t, ts, userID, programID)

		assertEnrollmentState(t, enrollment, ExpectedEnrollmentState{
			EnrollmentStatus: "ACTIVE",
			CycleStatus:      "PENDING",
			WeekStatus:       "PENDING",
			CurrentWeek:      1,
			CycleIteration:   1,
			HasActiveSession: false,
		})
	})

	// =============================================================================
	// TEST 2: ACTIVE -> BETWEEN_CYCLES (complete all weeks in cycle)
	// =============================================================================
	t.Run("Test ACTIVE to BETWEEN_CYCLES (cycle completes)", func(t *testing.T) {
		// We need to complete all workouts in the cycle
		// First, start and complete a workout for week 1
		sessionID := startWorkoutSession(t, ts, userID)
		finishWorkoutSession(t, ts, sessionID, userID)

		// Verify still ACTIVE after week 1 workout
		enrollment := getEnrollment(t, ts, userID)
		if enrollment.EnrollmentStatus != "ACTIVE" {
			t.Errorf("Expected ACTIVE after week 1 workout, got %s", enrollment.EnrollmentStatus)
		}

		// Advance to week 2
		enrollment = advanceWeek(t, ts, userID)
		assertEnrollmentState(t, enrollment, ExpectedEnrollmentState{
			EnrollmentStatus: "ACTIVE",
			CycleStatus:      "PENDING",
			WeekStatus:       "PENDING",
			CurrentWeek:      2,
			CycleIteration:   1,
			HasActiveSession: false,
		})

		// Start and complete week 2 workout
		sessionID = startWorkoutSession(t, ts, userID)
		finishWorkoutSession(t, ts, sessionID, userID)

		// Advance past final week -> should transition to BETWEEN_CYCLES
		enrollment = advanceWeek(t, ts, userID)
		assertEnrollmentState(t, enrollment, ExpectedEnrollmentState{
			EnrollmentStatus: "BETWEEN_CYCLES",
			CycleStatus:      "COMPLETED",
			WeekStatus:       "COMPLETED",
			CurrentWeek:      2,
			CycleIteration:   1,
			HasActiveSession: false,
		})
	})

	// =============================================================================
	// TEST 3: BETWEEN_CYCLES -> ACTIVE (start new cycle)
	// =============================================================================
	t.Run("Test BETWEEN_CYCLES to ACTIVE (start new cycle)", func(t *testing.T) {
		// User is already in BETWEEN_CYCLES from previous test
		enrollment := getEnrollment(t, ts, userID)
		if enrollment.EnrollmentStatus != "BETWEEN_CYCLES" {
			t.Fatalf("Expected BETWEEN_CYCLES state to start test, got %s", enrollment.EnrollmentStatus)
		}

		// Start new cycle
		enrollment = startNextCycle(t, ts, userID)

		assertEnrollmentState(t, enrollment, ExpectedEnrollmentState{
			EnrollmentStatus: "ACTIVE",
			CycleStatus:      "PENDING",
			WeekStatus:       "PENDING",
			CurrentWeek:      1,
			CycleIteration:   2, // Incremented!
			HasActiveSession: false,
		})
	})

	// =============================================================================
	// TEST 4: ACTIVE -> QUIT (unenroll)
	// =============================================================================
	t.Run("Test ACTIVE to QUIT (unenroll)", func(t *testing.T) {
		// User is ACTIVE from previous test (cycle 2, week 1)
		enrollment := getEnrollment(t, ts, userID)
		if enrollment.EnrollmentStatus != "ACTIVE" {
			t.Fatalf("Expected ACTIVE state to start test, got %s", enrollment.EnrollmentStatus)
		}

		// Unenroll
		unenrollUser(t, ts, userID)

		// Verify no longer enrolled - GET should return 404
		req, err := http.NewRequest(http.MethodGet, ts.URL("/users/"+userID+"/program"), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("X-User-ID", userID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to get enrollment: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 404 after unenroll, got %d", resp.StatusCode)
		}
	})

	// =============================================================================
	// TEST 5: BETWEEN_CYCLES -> QUIT (unenroll while deciding)
	// =============================================================================
	t.Run("Test BETWEEN_CYCLES to QUIT (unenroll while deciding)", func(t *testing.T) {
		// Re-enroll for this test
		enrollment := enrollUser(t, ts, userID, programID)
		if enrollment.EnrollmentStatus != "ACTIVE" {
			t.Fatalf("Expected ACTIVE after re-enrollment, got %s", enrollment.EnrollmentStatus)
		}

		// Complete both weeks to reach BETWEEN_CYCLES
		sessionID := startWorkoutSession(t, ts, userID)
		finishWorkoutSession(t, ts, sessionID, userID)
		advanceWeek(t, ts, userID)

		sessionID = startWorkoutSession(t, ts, userID)
		finishWorkoutSession(t, ts, sessionID, userID)
		enrollment = advanceWeek(t, ts, userID)

		if enrollment.EnrollmentStatus != "BETWEEN_CYCLES" {
			t.Fatalf("Expected BETWEEN_CYCLES, got %s", enrollment.EnrollmentStatus)
		}

		// Unenroll while in BETWEEN_CYCLES
		unenrollUser(t, ts, userID)

		// Verify no longer enrolled
		req, err := http.NewRequest(http.MethodGet, ts.URL("/users/"+userID+"/program"), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("X-User-ID", userID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to get enrollment: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 404 after unenroll from BETWEEN_CYCLES, got %d", resp.StatusCode)
		}
	})

	// =============================================================================
	// TEST 6: Can't start workout when BETWEEN_CYCLES
	// =============================================================================
	t.Run("Test cannot start workout when BETWEEN_CYCLES", func(t *testing.T) {
		// Re-enroll for this test
		enrollment := enrollUser(t, ts, userID, programID)
		if enrollment.EnrollmentStatus != "ACTIVE" {
			t.Fatalf("Expected ACTIVE after re-enrollment, got %s", enrollment.EnrollmentStatus)
		}

		// Complete both weeks to reach BETWEEN_CYCLES
		sessionID := startWorkoutSession(t, ts, userID)
		finishWorkoutSession(t, ts, sessionID, userID)
		advanceWeek(t, ts, userID)

		sessionID = startWorkoutSession(t, ts, userID)
		finishWorkoutSession(t, ts, sessionID, userID)
		enrollment = advanceWeek(t, ts, userID)

		if enrollment.EnrollmentStatus != "BETWEEN_CYCLES" {
			t.Fatalf("Expected BETWEEN_CYCLES, got %s", enrollment.EnrollmentStatus)
		}

		// Attempt to start workout - should fail
		req, err := http.NewRequest(http.MethodPost, ts.URL("/workouts/start"), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("X-User-ID", userID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to start workout: %v", err)
		}
		defer resp.Body.Close()

		// Should get 400 Bad Request
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400 Bad Request when starting workout in BETWEEN_CYCLES, got %d", resp.StatusCode)
		}

		// Parse and verify error
		body, _ := io.ReadAll(resp.Body)
		var errResp ErrorResponseData
		if err := json.Unmarshal(body, &errResp); err != nil {
			t.Fatalf("Failed to parse error response: %v", err)
		}

		// Verify error code
		if errResp.Error.Code != "invalid_enrollment_state" {
			t.Errorf("Expected error code 'invalid_enrollment_state', got '%s'", errResp.Error.Code)
		}

		// Clean up - unenroll
		unenrollUser(t, ts, userID)
	})

	// =============================================================================
	// TEST 7: Can't start new cycle when ACTIVE
	// =============================================================================
	t.Run("Test cannot start new cycle when ACTIVE", func(t *testing.T) {
		// Enroll fresh
		enrollment := enrollUser(t, ts, userID, programID)
		if enrollment.EnrollmentStatus != "ACTIVE" {
			t.Fatalf("Expected ACTIVE after enrollment, got %s", enrollment.EnrollmentStatus)
		}

		// Attempt to start next cycle - should fail
		req, err := http.NewRequest(http.MethodPost, ts.URL("/users/"+userID+"/enrollment/next-cycle"), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("X-User-ID", userID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to call next-cycle: %v", err)
		}
		defer resp.Body.Close()

		// Should get 400 Bad Request
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400 Bad Request, got %d", resp.StatusCode)
		}

		// Parse and verify error
		body, _ := io.ReadAll(resp.Body)
		var errResp ErrorResponseData
		if err := json.Unmarshal(body, &errResp); err != nil {
			t.Fatalf("Failed to parse error response: %v", err)
		}

		// Verify error code
		if errResp.Error.Code != "invalid_enrollment_state" {
			t.Errorf("Expected error code 'invalid_enrollment_state', got '%s'", errResp.Error.Code)
		}

		// Clean up
		unenrollUser(t, ts, userID)
	})
}

// ErrorResponseData represents the API error response structure.
type ErrorResponseData struct {
	Error struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Details interface{} `json:"details,omitempty"`
	} `json:"error"`
}
