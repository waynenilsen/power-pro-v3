package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// TestStateMachineWorkoutSessionTransitions tests workout session state transitions.
// This test focuses on the session-level state machine:
// - Start workout -> IN_PROGRESS session
// - Can't start second workout while one IN_PROGRESS
// - Finish workout -> COMPLETED session
// - Abandon workout -> ABANDONED session
// - Can start new workout after COMPLETED
// - Can start new workout after ABANDONED
// - Logging sets requires active session
// - Can't log sets to COMPLETED session
func TestStateMachineWorkoutSessionTransitions(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Test-unique identifier
	testID := uuid.New().String()[:8]

	// Use the seeded workout test user
	userID := "workout-test-user"

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"

	// Create training maxes
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", 225.0)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", 135.0)

	// Create a simple 1-week program for testing session transitions
	// Create prescriptions
	squatPrescID := createPrescription(t, ts, squatID, 3, 5, 100.0, 0)
	benchPrescID := createPrescription(t, ts, benchID, 3, 5, 100.0, 1)

	// Create Day A
	dayASlug := "ws-day-a-" + testID
	dayABody := fmt.Sprintf(`{"name": "WS Test Day A", "slug": "%s"}`, dayASlug)
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

	// Create 1-week cycle
	cycleName := "WS Test Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 1}`, cycleName)
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

	// Add day to week
	addDayToWeek(t, ts, week1ID, dayAID, "MONDAY")

	// Create program
	programSlug := "ws-test-program-" + testID
	programBody := fmt.Sprintf(`{"name": "WS State Test Program", "slug": "%s", "cycleId": "%s"}`, programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	// =============================================================================
	// TEST 1: Start workout creates IN_PROGRESS session
	// =============================================================================
	t.Run("Test start workout creates IN_PROGRESS session", func(t *testing.T) {
		// Enroll user fresh
		enrollment := enrollUser(t, ts, userID, programID)

		// Verify initial state
		assertEnrollmentState(t, enrollment, ExpectedEnrollmentState{
			EnrollmentStatus: "ACTIVE",
			CycleStatus:      "PENDING",
			WeekStatus:       "PENDING",
			CurrentWeek:      1,
			CycleIteration:   1,
			HasActiveSession: false,
		})

		// Start workout session
		sessionID := startWorkoutSession(t, ts, userID)
		if sessionID == "" {
			t.Fatal("Expected session ID from start workout")
		}

		// Verify enrollment has active session
		// Note: CycleStatus and WeekStatus remain PENDING when starting a workout
		// They only transition to COMPLETED after advancing week/cycle
		enrollment = getEnrollment(t, ts, userID)
		assertEnrollmentState(t, enrollment, ExpectedEnrollmentState{
			EnrollmentStatus: "ACTIVE",
			CycleStatus:      "PENDING",
			WeekStatus:       "PENDING",
			CurrentWeek:      1,
			CycleIteration:   1,
			HasActiveSession: true,
			SessionStatus:    "IN_PROGRESS",
		})

		// Clean up for next test
		finishWorkoutSession(t, ts, sessionID, userID)
		unenrollUser(t, ts, userID)
	})

	// =============================================================================
	// TEST 2: Can't start second workout while one IN_PROGRESS
	// =============================================================================
	t.Run("Test cannot start second workout while one IN_PROGRESS", func(t *testing.T) {
		// Enroll and start first workout
		enrollUser(t, ts, userID, programID)
		startWorkoutSession(t, ts, userID)

		// Attempt to start second workout
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

		// Should get 400 Bad Request (or 409 Conflict)
		if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected 400/409 when starting second workout, got %d", resp.StatusCode)
		}

		// Parse and verify error
		body, _ := io.ReadAll(resp.Body)
		var errResp ErrorResponseData
		if err := json.Unmarshal(body, &errResp); err != nil {
			t.Fatalf("Failed to parse error response: %v", err)
		}

		// Verify we got an appropriate error
		if errResp.Error.Code == "" {
			t.Error("Expected error code in response")
		}

		// Clean up
		unenrollUser(t, ts, userID)
	})

	// =============================================================================
	// TEST 3: Finish workout transitions to COMPLETED
	// =============================================================================
	t.Run("Test finish workout transitions to COMPLETED", func(t *testing.T) {
		// Enroll and start workout
		enrollUser(t, ts, userID, programID)
		sessionID := startWorkoutSession(t, ts, userID)

		// Log at least one set (required before finishing)
		logTestSets(t, ts, userID, sessionID, squatPrescID, squatID)

		// Finish workout
		finishWorkoutSession(t, ts, sessionID, userID)

		// Verify enrollment no longer has active session
		// Note: CycleStatus and WeekStatus remain PENDING after finishing a workout
		// They only transition to COMPLETED after calling advance-week
		enrollment := getEnrollment(t, ts, userID)
		assertEnrollmentState(t, enrollment, ExpectedEnrollmentState{
			EnrollmentStatus: "ACTIVE",
			CycleStatus:      "PENDING",
			WeekStatus:       "PENDING",
			CurrentWeek:      1,
			CycleIteration:   1,
			HasActiveSession: false,
		})

		// Verify session is COMPLETED by checking it exists
		sessionResp := getWorkoutSession(t, ts, sessionID, userID)
		if sessionResp.Status != "COMPLETED" {
			t.Errorf("Expected session status COMPLETED, got %s", sessionResp.Status)
		}

		// Clean up
		unenrollUser(t, ts, userID)
	})

	// =============================================================================
	// TEST 4: Abandon workout transitions to ABANDONED
	// =============================================================================
	t.Run("Test abandon workout transitions to ABANDONED", func(t *testing.T) {
		// Enroll and start workout
		enrollUser(t, ts, userID, programID)
		sessionID := startWorkoutSession(t, ts, userID)

		// Abandon the workout
		abandonWorkoutSession(t, ts, sessionID, userID)

		// Verify enrollment no longer has active session
		enrollment := getEnrollment(t, ts, userID)
		if enrollment.CurrentWorkoutSession != nil {
			t.Error("Expected no active session after abandon")
		}

		// Verify session is ABANDONED
		sessionResp := getWorkoutSession(t, ts, sessionID, userID)
		if sessionResp.Status != "ABANDONED" {
			t.Errorf("Expected session status ABANDONED, got %s", sessionResp.Status)
		}

		// Clean up
		unenrollUser(t, ts, userID)
	})

	// =============================================================================
	// TEST 5: Can start new workout after COMPLETED
	// =============================================================================
	t.Run("Test can start new workout after COMPLETED", func(t *testing.T) {
		// Enroll
		enrollUser(t, ts, userID, programID)

		// Start and complete first workout
		sessionID1 := startWorkoutSession(t, ts, userID)
		logTestSets(t, ts, userID, sessionID1, squatPrescID, squatID)
		finishWorkoutSession(t, ts, sessionID1, userID)

		// Verify first session is completed
		session1 := getWorkoutSession(t, ts, sessionID1, userID)
		if session1.Status != "COMPLETED" {
			t.Errorf("Expected first session status COMPLETED, got %s", session1.Status)
		}

		// Start new workout
		sessionID2 := startWorkoutSession(t, ts, userID)
		if sessionID2 == "" {
			t.Fatal("Expected to start new workout after completing previous")
		}

		// Verify new session is different
		if sessionID1 == sessionID2 {
			t.Error("New session should have different ID than completed session")
		}

		// Verify enrollment has new active session
		enrollment := getEnrollment(t, ts, userID)
		if enrollment.CurrentWorkoutSession == nil {
			t.Error("Expected active session after starting new workout")
		} else if enrollment.CurrentWorkoutSession.Status != "IN_PROGRESS" {
			t.Errorf("Expected new session IN_PROGRESS, got %s", enrollment.CurrentWorkoutSession.Status)
		}

		// Clean up
		unenrollUser(t, ts, userID)
	})

	// =============================================================================
	// TEST 6: Can start new workout after ABANDONED
	// =============================================================================
	t.Run("Test can start new workout after ABANDONED", func(t *testing.T) {
		// Enroll
		enrollUser(t, ts, userID, programID)

		// Start and abandon first workout
		sessionID1 := startWorkoutSession(t, ts, userID)
		abandonWorkoutSession(t, ts, sessionID1, userID)

		// Verify first session is abandoned
		session1 := getWorkoutSession(t, ts, sessionID1, userID)
		if session1.Status != "ABANDONED" {
			t.Errorf("Expected first session status ABANDONED, got %s", session1.Status)
		}

		// Start new workout
		sessionID2 := startWorkoutSession(t, ts, userID)
		if sessionID2 == "" {
			t.Fatal("Expected to start new workout after abandoning previous")
		}

		// Verify new session is different
		if sessionID1 == sessionID2 {
			t.Error("New session should have different ID than abandoned session")
		}

		// Verify enrollment has new active session
		enrollment := getEnrollment(t, ts, userID)
		if enrollment.CurrentWorkoutSession == nil {
			t.Error("Expected active session after starting new workout")
		} else if enrollment.CurrentWorkoutSession.Status != "IN_PROGRESS" {
			t.Errorf("Expected new session IN_PROGRESS, got %s", enrollment.CurrentWorkoutSession.Status)
		}

		// Clean up
		unenrollUser(t, ts, userID)
	})

	// =============================================================================
	// TEST 7: Logging sets requires active session
	// =============================================================================
	t.Run("Test logging sets requires active session", func(t *testing.T) {
		// Enroll but don't start a session
		enrollUser(t, ts, userID, programID)

		// Attempt to log sets with a fake session ID
		fakeSessionID := uuid.New().String()
		err := tryLogSets(t, ts, userID, fakeSessionID, squatPrescID, squatID)
		if err == nil {
			t.Error("Expected error when logging sets without active session")
		}

		// Clean up
		unenrollUser(t, ts, userID)
	})

	// =============================================================================
	// TEST 8: Can't log sets to COMPLETED session
	// =============================================================================
	t.Run("Test cannot log sets to COMPLETED session", func(t *testing.T) {
		// Enroll and complete a workout
		enrollUser(t, ts, userID, programID)
		sessionID := startWorkoutSession(t, ts, userID)
		logTestSets(t, ts, userID, sessionID, squatPrescID, squatID)
		finishWorkoutSession(t, ts, sessionID, userID)

		// Verify session is completed
		session := getWorkoutSession(t, ts, sessionID, userID)
		if session.Status != "COMPLETED" {
			t.Fatalf("Expected COMPLETED session, got %s", session.Status)
		}

		// Attempt to log more sets to the completed session
		err := tryLogSets(t, ts, userID, sessionID, benchPrescID, benchID)
		if err == nil {
			t.Error("Expected error when logging sets to COMPLETED session")
		}

		// Clean up
		unenrollUser(t, ts, userID)
	})
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// abandonWorkoutSession abandons a workout session.
func abandonWorkoutSession(t *testing.T, ts *testutil.TestServer, sessionID, userID string) {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, ts.URL("/workouts/"+sessionID+"/abandon"), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to abandon workout: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to abandon workout session, status %d: %s", resp.StatusCode, body)
	}
}

// getWorkoutSession retrieves a workout session by ID.
func getWorkoutSession(t *testing.T, ts *testutil.TestServer, sessionID, userID string) WorkoutSessionData {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, ts.URL("/workouts/"+sessionID), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to get workout session: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to get workout session, status %d: %s", resp.StatusCode, body)
	}

	var envelope WorkoutSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		t.Fatalf("Failed to decode workout session response: %v", err)
	}
	return envelope.Data
}

// logTestSets logs a minimal set of sets to satisfy session completion requirements.
func logTestSets(t *testing.T, ts *testutil.TestServer, userID, sessionID, prescriptionID, liftID string) {
	t.Helper()

	type setRequest struct {
		PrescriptionID string  `json:"prescriptionId"`
		LiftID         string  `json:"liftId"`
		SetNumber      int     `json:"setNumber"`
		Weight         float64 `json:"weight"`
		TargetReps     int     `json:"targetReps"`
		RepsPerformed  int     `json:"repsPerformed"`
		IsAMRAP        bool    `json:"isAmrap"`
	}

	// Log 3 sets of 5 reps
	setsReq := []setRequest{
		{PrescriptionID: prescriptionID, LiftID: liftID, SetNumber: 1, Weight: 225.0, TargetReps: 5, RepsPerformed: 5},
		{PrescriptionID: prescriptionID, LiftID: liftID, SetNumber: 2, Weight: 225.0, TargetReps: 5, RepsPerformed: 5},
		{PrescriptionID: prescriptionID, LiftID: liftID, SetNumber: 3, Weight: 225.0, TargetReps: 5, RepsPerformed: 5},
	}

	body, _ := json.Marshal(map[string]interface{}{"sets": setsReq})
	req, _ := http.NewRequest(http.MethodPost, ts.URL("/sessions/"+sessionID+"/sets"), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to log sets: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to log sets, status %d: %s", resp.StatusCode, respBody)
	}
}

// tryLogSets attempts to log sets and returns an error if it fails.
func tryLogSets(t *testing.T, ts *testutil.TestServer, userID, sessionID, prescriptionID, liftID string) error {
	t.Helper()

	type setRequest struct {
		PrescriptionID string  `json:"prescriptionId"`
		LiftID         string  `json:"liftId"`
		SetNumber      int     `json:"setNumber"`
		Weight         float64 `json:"weight"`
		TargetReps     int     `json:"targetReps"`
		RepsPerformed  int     `json:"repsPerformed"`
		IsAMRAP        bool    `json:"isAmrap"`
	}

	setsReq := []setRequest{
		{PrescriptionID: prescriptionID, LiftID: liftID, SetNumber: 1, Weight: 225.0, TargetReps: 5, RepsPerformed: 5},
	}

	body, _ := json.Marshal(map[string]interface{}{"sets": setsReq})
	req, _ := http.NewRequest(http.MethodPost, ts.URL("/sessions/"+sessionID+"/sets"), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to log sets, status %d: %s", resp.StatusCode, respBody)
	}

	return nil
}
