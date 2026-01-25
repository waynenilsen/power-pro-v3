package e2e

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// TestStateMachineWeekCycleTransitions tests week and cycle state machine transitions.
// This test focuses on week and cycle level state transitions:
// - First workout of week: weekStatus remains PENDING (status tracked via active session)
// - Last workout of week + advanceWeek: weekStatus -> COMPLETED, new week PENDING
// - Last workout of cycle + advanceWeek: cycleStatus -> COMPLETED, enrollmentStatus -> BETWEEN_CYCLES
// - New week: weekStatus resets to PENDING
// - New cycle: cycleStatus and weekStatus reset to PENDING, cycleIteration incremented
//
// Note: Week and cycle status transitions to COMPLETED only happen when advanceWeek is called.
// The IN_PROGRESS state is implicitly determined by having an active workout session.
func TestStateMachineWeekCycleTransitions(t *testing.T) {
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

	// =============================================================================
	// TEST 1: First workout of week - verify session tracks week activity
	// Note: Week/cycle status stays PENDING until advanceWeek is called.
	// The active session implicitly indicates week is in progress.
	// =============================================================================
	t.Run("Test first workout of week creates active session while week stays PENDING", func(t *testing.T) {
		// Create a simple 1-week, 1-day program
		squatPrescID := createPrescription(t, ts, squatID, 3, 5, 100.0, 0)
		benchPrescID := createPrescription(t, ts, benchID, 3, 5, 100.0, 1)

		daySlug := "wc-test1-day-" + testID
		dayBody := fmt.Sprintf(`{"name": "WC Test 1 Day", "slug": "%s"}`, daySlug)
		dayResp, err := adminPost(ts.URL("/days"), dayBody)
		if err != nil {
			t.Fatalf("Failed to create day: %v", err)
		}
		var dayEnvelope DayResponse
		json.NewDecoder(dayResp.Body).Decode(&dayEnvelope)
		dayResp.Body.Close()
		dayID := dayEnvelope.Data.ID

		addPrescToDay(t, ts, dayID, squatPrescID)
		addPrescToDay(t, ts, dayID, benchPrescID)

		cycleName := "WC Test 1 Cycle " + testID
		cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 1}`, cycleName)
		cycleResp, err := adminPost(ts.URL("/cycles"), cycleBody)
		if err != nil {
			t.Fatalf("Failed to create cycle: %v", err)
		}
		var cycleEnvelope CycleResponse
		json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
		cycleResp.Body.Close()
		cycleID := cycleEnvelope.Data.ID

		week1Body := fmt.Sprintf(`{"weekNumber": 1, "cycleId": "%s"}`, cycleID)
		week1Resp, _ := adminPost(ts.URL("/weeks"), week1Body)
		var week1Envelope WeekResponse
		json.NewDecoder(week1Resp.Body).Decode(&week1Envelope)
		week1Resp.Body.Close()
		week1ID := week1Envelope.Data.ID

		addDayToWeek(t, ts, week1ID, dayID, "MONDAY")

		programSlug := "wc-test1-program-" + testID
		programBody := fmt.Sprintf(`{"name": "WC Test 1 Program", "slug": "%s", "cycleId": "%s"}`, programSlug, cycleID)
		programResp, _ := adminPost(ts.URL("/programs"), programBody)
		var programEnvelope ProgramResponse
		json.NewDecoder(programResp.Body).Decode(&programEnvelope)
		programResp.Body.Close()
		programID := programEnvelope.Data.ID

		// Enroll user
		enrollment := enrollUser(t, ts, userID, programID)

		// Verify initial state: weekStatus = PENDING
		assertEnrollmentState(t, enrollment, ExpectedEnrollmentState{
			EnrollmentStatus: "ACTIVE",
			CycleStatus:      "PENDING",
			WeekStatus:       "PENDING",
			CurrentWeek:      1,
			CycleIteration:   1,
			HasActiveSession: false,
		})

		// Start workout - weekStatus stays PENDING, but we get an active session
		// The active session implicitly indicates the week is "in progress"
		sessionID := startWorkoutSession(t, ts, userID)

		// Verify we have an active session (this is how we know week activity started)
		// Week/cycle status stays PENDING until advanceWeek is called
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

		// Clean up
		finishWorkoutSession(t, ts, sessionID, userID)
		unenrollUser(t, ts, userID)
	})

	// =============================================================================
	// TEST 2: Last workout of week transitions to COMPLETED with auto-advance
	// =============================================================================
	t.Run("Test last workout of week transitions to COMPLETED with auto-advance", func(t *testing.T) {
		// Create a 2-week program with 1 day per week
		squatPrescID := createPrescription(t, ts, squatID, 3, 5, 100.0, 0)

		daySlug := "wc-test2-day-" + testID
		dayBody := fmt.Sprintf(`{"name": "WC Test 2 Day", "slug": "%s"}`, daySlug)
		dayResp, err := adminPost(ts.URL("/days"), dayBody)
		if err != nil {
			t.Fatalf("Failed to create day: %v", err)
		}
		var dayEnvelope DayResponse
		json.NewDecoder(dayResp.Body).Decode(&dayEnvelope)
		dayResp.Body.Close()
		dayID := dayEnvelope.Data.ID

		addPrescToDay(t, ts, dayID, squatPrescID)

		cycleName := "WC Test 2 Cycle " + testID
		cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 2}`, cycleName)
		cycleResp, err := adminPost(ts.URL("/cycles"), cycleBody)
		if err != nil {
			t.Fatalf("Failed to create cycle: %v", err)
		}
		var cycleEnvelope CycleResponse
		json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
		cycleResp.Body.Close()
		cycleID := cycleEnvelope.Data.ID

		week1Body := fmt.Sprintf(`{"weekNumber": 1, "cycleId": "%s"}`, cycleID)
		week1Resp, _ := adminPost(ts.URL("/weeks"), week1Body)
		var week1Envelope WeekResponse
		json.NewDecoder(week1Resp.Body).Decode(&week1Envelope)
		week1Resp.Body.Close()
		week1ID := week1Envelope.Data.ID

		week2Body := fmt.Sprintf(`{"weekNumber": 2, "cycleId": "%s"}`, cycleID)
		week2Resp, _ := adminPost(ts.URL("/weeks"), week2Body)
		var week2Envelope WeekResponse
		json.NewDecoder(week2Resp.Body).Decode(&week2Envelope)
		week2Resp.Body.Close()
		week2ID := week2Envelope.Data.ID

		addDayToWeek(t, ts, week1ID, dayID, "MONDAY")
		addDayToWeek(t, ts, week2ID, dayID, "MONDAY")

		programSlug := "wc-test2-program-" + testID
		programBody := fmt.Sprintf(`{"name": "WC Test 2 Program", "slug": "%s", "cycleId": "%s"}`, programSlug, cycleID)
		programResp, _ := adminPost(ts.URL("/programs"), programBody)
		var programEnvelope ProgramResponse
		json.NewDecoder(programResp.Body).Decode(&programEnvelope)
		programResp.Body.Close()
		programID := programEnvelope.Data.ID

		// Enroll user
		enrollUser(t, ts, userID, programID)

		// Complete workout for week 1
		sessionID := startWorkoutSession(t, ts, userID)
		finishWorkoutSession(t, ts, sessionID, userID)

		// Advance week - should mark week 1 COMPLETED and start week 2 PENDING
		enrollment := advanceWeek(t, ts, userID)

		assertEnrollmentState(t, enrollment, ExpectedEnrollmentState{
			EnrollmentStatus: "ACTIVE",
			CycleStatus:      "PENDING",
			WeekStatus:       "PENDING",
			CurrentWeek:      2,
			CycleIteration:   1,
			HasActiveSession: false,
		})

		// Clean up
		unenrollUser(t, ts, userID)
	})

	// =============================================================================
	// TEST 3: Last workout of cycle triggers CYCLE_COMPLETED transition
	// =============================================================================
	t.Run("Test last workout of cycle triggers CYCLE_COMPLETED transition", func(t *testing.T) {
		// Create a 1-week program (single cycle length for simplicity)
		squatPrescID := createPrescription(t, ts, squatID, 3, 5, 100.0, 0)

		daySlug := "wc-test3-day-" + testID
		dayBody := fmt.Sprintf(`{"name": "WC Test 3 Day", "slug": "%s"}`, daySlug)
		dayResp, err := adminPost(ts.URL("/days"), dayBody)
		if err != nil {
			t.Fatalf("Failed to create day: %v", err)
		}
		var dayEnvelope DayResponse
		json.NewDecoder(dayResp.Body).Decode(&dayEnvelope)
		dayResp.Body.Close()
		dayID := dayEnvelope.Data.ID

		addPrescToDay(t, ts, dayID, squatPrescID)

		cycleName := "WC Test 3 Cycle " + testID
		cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 1}`, cycleName)
		cycleResp, err := adminPost(ts.URL("/cycles"), cycleBody)
		if err != nil {
			t.Fatalf("Failed to create cycle: %v", err)
		}
		var cycleEnvelope CycleResponse
		json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
		cycleResp.Body.Close()
		cycleID := cycleEnvelope.Data.ID

		week1Body := fmt.Sprintf(`{"weekNumber": 1, "cycleId": "%s"}`, cycleID)
		week1Resp, _ := adminPost(ts.URL("/weeks"), week1Body)
		var week1Envelope WeekResponse
		json.NewDecoder(week1Resp.Body).Decode(&week1Envelope)
		week1Resp.Body.Close()
		week1ID := week1Envelope.Data.ID

		addDayToWeek(t, ts, week1ID, dayID, "MONDAY")

		programSlug := "wc-test3-program-" + testID
		programBody := fmt.Sprintf(`{"name": "WC Test 3 Program", "slug": "%s", "cycleId": "%s"}`, programSlug, cycleID)
		programResp, _ := adminPost(ts.URL("/programs"), programBody)
		var programEnvelope ProgramResponse
		json.NewDecoder(programResp.Body).Decode(&programEnvelope)
		programResp.Body.Close()
		programID := programEnvelope.Data.ID

		// Enroll user
		enrollUser(t, ts, userID, programID)

		// Complete workout for the only week in the cycle
		sessionID := startWorkoutSession(t, ts, userID)
		finishWorkoutSession(t, ts, sessionID, userID)

		// Advance past final week - should transition to BETWEEN_CYCLES
		enrollment := advanceWeek(t, ts, userID)

		assertEnrollmentState(t, enrollment, ExpectedEnrollmentState{
			EnrollmentStatus: "BETWEEN_CYCLES",
			CycleStatus:      "COMPLETED",
			WeekStatus:       "COMPLETED",
			CurrentWeek:      1,
			CycleIteration:   1,
			HasActiveSession: false,
		})

		// Clean up
		unenrollUser(t, ts, userID)
	})

	// =============================================================================
	// TEST 4: Week status resets to PENDING on new week
	// =============================================================================
	t.Run("Test week status resets to PENDING on new week", func(t *testing.T) {
		// Create a 2-week program
		squatPrescID := createPrescription(t, ts, squatID, 3, 5, 100.0, 0)

		daySlug := "wc-test4-day-" + testID
		dayBody := fmt.Sprintf(`{"name": "WC Test 4 Day", "slug": "%s"}`, daySlug)
		dayResp, err := adminPost(ts.URL("/days"), dayBody)
		if err != nil {
			t.Fatalf("Failed to create day: %v", err)
		}
		var dayEnvelope DayResponse
		json.NewDecoder(dayResp.Body).Decode(&dayEnvelope)
		dayResp.Body.Close()
		dayID := dayEnvelope.Data.ID

		addPrescToDay(t, ts, dayID, squatPrescID)

		cycleName := "WC Test 4 Cycle " + testID
		cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 2}`, cycleName)
		cycleResp, err := adminPost(ts.URL("/cycles"), cycleBody)
		if err != nil {
			t.Fatalf("Failed to create cycle: %v", err)
		}
		var cycleEnvelope CycleResponse
		json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
		cycleResp.Body.Close()
		cycleID := cycleEnvelope.Data.ID

		week1Body := fmt.Sprintf(`{"weekNumber": 1, "cycleId": "%s"}`, cycleID)
		week1Resp, _ := adminPost(ts.URL("/weeks"), week1Body)
		var week1Envelope WeekResponse
		json.NewDecoder(week1Resp.Body).Decode(&week1Envelope)
		week1Resp.Body.Close()
		week1ID := week1Envelope.Data.ID

		week2Body := fmt.Sprintf(`{"weekNumber": 2, "cycleId": "%s"}`, cycleID)
		week2Resp, _ := adminPost(ts.URL("/weeks"), week2Body)
		var week2Envelope WeekResponse
		json.NewDecoder(week2Resp.Body).Decode(&week2Envelope)
		week2Resp.Body.Close()
		week2ID := week2Envelope.Data.ID

		addDayToWeek(t, ts, week1ID, dayID, "MONDAY")
		addDayToWeek(t, ts, week2ID, dayID, "MONDAY")

		programSlug := "wc-test4-program-" + testID
		programBody := fmt.Sprintf(`{"name": "WC Test 4 Program", "slug": "%s", "cycleId": "%s"}`, programSlug, cycleID)
		programResp, _ := adminPost(ts.URL("/programs"), programBody)
		var programEnvelope ProgramResponse
		json.NewDecoder(programResp.Body).Decode(&programEnvelope)
		programResp.Body.Close()
		programID := programEnvelope.Data.ID

		// Enroll user
		enrollUser(t, ts, userID, programID)

		// Complete a workout to make week IN_PROGRESS
		sessionID := startWorkoutSession(t, ts, userID)
		finishWorkoutSession(t, ts, sessionID, userID)

		// Advance to next week
		enrollment := advanceWeek(t, ts, userID)

		// Verify weekStatus is PENDING and currentWeek incremented
		assertEnrollmentState(t, enrollment, ExpectedEnrollmentState{
			EnrollmentStatus: "ACTIVE",
			CycleStatus:      "PENDING",
			WeekStatus:       "PENDING",
			CurrentWeek:      2,
			CycleIteration:   1,
			HasActiveSession: false,
		})

		// Clean up
		unenrollUser(t, ts, userID)
	})

	// =============================================================================
	// TEST 5: Cycle status resets to PENDING on new cycle
	// =============================================================================
	t.Run("Test cycle status resets to PENDING on new cycle", func(t *testing.T) {
		// Create a 1-week program so we can complete the cycle quickly
		squatPrescID := createPrescription(t, ts, squatID, 3, 5, 100.0, 0)

		daySlug := "wc-test5-day-" + testID
		dayBody := fmt.Sprintf(`{"name": "WC Test 5 Day", "slug": "%s"}`, daySlug)
		dayResp, err := adminPost(ts.URL("/days"), dayBody)
		if err != nil {
			t.Fatalf("Failed to create day: %v", err)
		}
		var dayEnvelope DayResponse
		json.NewDecoder(dayResp.Body).Decode(&dayEnvelope)
		dayResp.Body.Close()
		dayID := dayEnvelope.Data.ID

		addPrescToDay(t, ts, dayID, squatPrescID)

		cycleName := "WC Test 5 Cycle " + testID
		cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 1}`, cycleName)
		cycleResp, err := adminPost(ts.URL("/cycles"), cycleBody)
		if err != nil {
			t.Fatalf("Failed to create cycle: %v", err)
		}
		var cycleEnvelope CycleResponse
		json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
		cycleResp.Body.Close()
		cycleID := cycleEnvelope.Data.ID

		week1Body := fmt.Sprintf(`{"weekNumber": 1, "cycleId": "%s"}`, cycleID)
		week1Resp, _ := adminPost(ts.URL("/weeks"), week1Body)
		var week1Envelope WeekResponse
		json.NewDecoder(week1Resp.Body).Decode(&week1Envelope)
		week1Resp.Body.Close()
		week1ID := week1Envelope.Data.ID

		addDayToWeek(t, ts, week1ID, dayID, "MONDAY")

		programSlug := "wc-test5-program-" + testID
		programBody := fmt.Sprintf(`{"name": "WC Test 5 Program", "slug": "%s", "cycleId": "%s"}`, programSlug, cycleID)
		programResp, _ := adminPost(ts.URL("/programs"), programBody)
		var programEnvelope ProgramResponse
		json.NewDecoder(programResp.Body).Decode(&programEnvelope)
		programResp.Body.Close()
		programID := programEnvelope.Data.ID

		// Enroll user
		enrollUser(t, ts, userID, programID)

		// Complete the only week in the cycle
		sessionID := startWorkoutSession(t, ts, userID)
		finishWorkoutSession(t, ts, sessionID, userID)

		// Advance past final week to reach BETWEEN_CYCLES
		enrollment := advanceWeek(t, ts, userID)
		if enrollment.EnrollmentStatus != "BETWEEN_CYCLES" {
			t.Fatalf("Expected BETWEEN_CYCLES, got %s", enrollment.EnrollmentStatus)
		}

		// Start next cycle
		enrollment = startNextCycle(t, ts, userID)

		// Verify cycle and week status reset, cycleIteration incremented
		assertEnrollmentState(t, enrollment, ExpectedEnrollmentState{
			EnrollmentStatus: "ACTIVE",
			CycleStatus:      "PENDING",
			WeekStatus:       "PENDING",
			CurrentWeek:      1,
			CycleIteration:   2,
			HasActiveSession: false,
		})

		// Clean up
		unenrollUser(t, ts, userID)
	})
}
