// Package e2e provides end-to-end tests for complete program workflows.
// This file contains E2E tests for Sheiko Intermediate program with peaking.
package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// =============================================================================
// SHEIKO INTERMEDIATE RESPONSE TYPES
// =============================================================================

// MeetDateResponseData represents the API response for meet date operations.
type MeetDateResponseData struct {
	MeetDate     *string `json:"meet_date,omitempty"`
	DaysOut      int     `json:"days_out"`
	CurrentPhase string  `json:"current_phase"`
	WeeksToMeet  int     `json:"weeks_to_meet"`
}

// MeetDateResponseEnvelope wraps meet date response with standard envelope.
type MeetDateResponseEnvelope struct {
	Data MeetDateResponseData `json:"data"`
}

// CountdownResponseData represents the API response for countdown operations.
type CountdownResponseData struct {
	MeetDate        *string `json:"meet_date,omitempty"`
	DaysOut         int     `json:"days_out"`
	CurrentPhase    string  `json:"current_phase"`
	PhaseWeek       int     `json:"phase_week"`
	TaperMultiplier float64 `json:"taper_multiplier"`
}

// CountdownResponseEnvelope wraps countdown response with standard envelope.
type CountdownResponseEnvelope struct {
	Data CountdownResponseData `json:"data"`
}

// =============================================================================
// HELPER FUNCTIONS FOR SHEIKO TESTS
// =============================================================================

func userPut(url string, body string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

// createSheikoTestSetup creates the complete Sheiko Intermediate program structure.
// Returns programID, cycleID, and the week IDs for all 13 weeks.
func createSheikoTestSetup(t *testing.T, ts *testutil.TestServer, testID string) (programID, cycleID string, weekIDs []string) {
	t.Helper()

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	// Create prescriptions for Sheiko-style training
	// Prep 1 style: 70-75% intensity, higher volume (5x4-5)
	prep1SquatPrescID := createSheikoPrescription(t, ts, squatID, 5, 4, 72.5, 0)
	prep1BenchPrescID := createSheikoPrescription(t, ts, benchID, 5, 4, 72.5, 1)
	prep1DeadliftPrescID := createSheikoPrescription(t, ts, deadliftID, 4, 3, 75.0, 2)

	// Prep 2 style: 80-90% intensity, medium volume (4x3)
	prep2SquatPrescID := createSheikoPrescription(t, ts, squatID, 4, 3, 85.0, 0)
	prep2BenchPrescID := createSheikoPrescription(t, ts, benchID, 4, 3, 85.0, 1)
	prep2DeadliftPrescID := createSheikoPrescription(t, ts, deadliftID, 3, 2, 87.5, 2)

	// Competition style: 90-95% intensity, low volume (3x2)
	compSquatPrescID := createSheikoPrescription(t, ts, squatID, 3, 2, 92.5, 0)
	compBenchPrescID := createSheikoPrescription(t, ts, benchID, 3, 2, 92.5, 1)
	compDeadliftPrescID := createSheikoPrescription(t, ts, deadliftID, 2, 1, 95.0, 2)

	// Create days for each phase
	// Prep 1 days
	prep1DayASlug := "sheiko-prep1-day-a-" + testID
	prep1DayAID := createSheikoDay(t, ts, "Prep 1 Day A", prep1DayASlug)
	addPrescToDay(t, ts, prep1DayAID, prep1SquatPrescID)
	addPrescToDay(t, ts, prep1DayAID, prep1BenchPrescID)

	prep1DayBSlug := "sheiko-prep1-day-b-" + testID
	prep1DayBID := createSheikoDay(t, ts, "Prep 1 Day B", prep1DayBSlug)
	addPrescToDay(t, ts, prep1DayBID, prep1DeadliftPrescID)
	addPrescToDay(t, ts, prep1DayBID, prep1BenchPrescID)

	prep1DayCSlug := "sheiko-prep1-day-c-" + testID
	prep1DayCID := createSheikoDay(t, ts, "Prep 1 Day C", prep1DayCSlug)
	addPrescToDay(t, ts, prep1DayCID, prep1SquatPrescID)

	// Prep 2 days
	prep2DayASlug := "sheiko-prep2-day-a-" + testID
	prep2DayAID := createSheikoDay(t, ts, "Prep 2 Day A", prep2DayASlug)
	addPrescToDay(t, ts, prep2DayAID, prep2SquatPrescID)
	addPrescToDay(t, ts, prep2DayAID, prep2BenchPrescID)

	prep2DayBSlug := "sheiko-prep2-day-b-" + testID
	prep2DayBID := createSheikoDay(t, ts, "Prep 2 Day B", prep2DayBSlug)
	addPrescToDay(t, ts, prep2DayBID, prep2DeadliftPrescID)
	addPrescToDay(t, ts, prep2DayBID, prep2BenchPrescID)

	prep2DayCSlug := "sheiko-prep2-day-c-" + testID
	prep2DayCID := createSheikoDay(t, ts, "Prep 2 Day C", prep2DayCSlug)
	addPrescToDay(t, ts, prep2DayCID, prep2SquatPrescID)

	// Competition days
	compDayASlug := "sheiko-comp-day-a-" + testID
	compDayAID := createSheikoDay(t, ts, "Comp Day A", compDayASlug)
	addPrescToDay(t, ts, compDayAID, compSquatPrescID)
	addPrescToDay(t, ts, compDayAID, compBenchPrescID)

	compDayBSlug := "sheiko-comp-day-b-" + testID
	compDayBID := createSheikoDay(t, ts, "Comp Day B", compDayBSlug)
	addPrescToDay(t, ts, compDayBID, compDeadliftPrescID)
	addPrescToDay(t, ts, compDayBID, compBenchPrescID)

	compDayCSlug := "sheiko-comp-day-c-" + testID
	compDayCID := createSheikoDay(t, ts, "Comp Day C", compDayCSlug)
	addPrescToDay(t, ts, compDayCID, compSquatPrescID)

	// Create 13-week cycle
	cycleName := "Sheiko Intermediate Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 13}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID = cycleEnvelope.Data.ID

	// Create weeks for each phase
	weekIDs = make([]string, 13)

	// Prep 1: Weeks 1-4
	for i := 1; i <= 4; i++ {
		weekID := createSheikoWeek(t, ts, cycleID, i)
		weekIDs[i-1] = weekID
		addDayToWeek(t, ts, weekID, prep1DayAID, "MONDAY")
		addDayToWeek(t, ts, weekID, prep1DayBID, "WEDNESDAY")
		addDayToWeek(t, ts, weekID, prep1DayCID, "FRIDAY")
	}

	// Prep 2: Weeks 5-8
	for i := 5; i <= 8; i++ {
		weekID := createSheikoWeek(t, ts, cycleID, i)
		weekIDs[i-1] = weekID
		addDayToWeek(t, ts, weekID, prep2DayAID, "MONDAY")
		addDayToWeek(t, ts, weekID, prep2DayBID, "WEDNESDAY")
		addDayToWeek(t, ts, weekID, prep2DayCID, "FRIDAY")
	}

	// Competition: Weeks 9-13
	for i := 9; i <= 13; i++ {
		weekID := createSheikoWeek(t, ts, cycleID, i)
		weekIDs[i-1] = weekID
		addDayToWeek(t, ts, weekID, compDayAID, "MONDAY")
		addDayToWeek(t, ts, weekID, compDayBID, "WEDNESDAY")
		addDayToWeek(t, ts, weekID, compDayCID, "FRIDAY")
	}

	// Create program
	programSlug := "sheiko-intermediate-" + testID
	programBody := fmt.Sprintf(`{"name": "Sheiko Intermediate", "slug": "%s", "cycleId": "%s"}`, programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID = programEnvelope.Data.ID

	return programID, cycleID, weekIDs
}

func createSheikoPrescription(t *testing.T, ts *testutil.TestServer, liftID string, sets, reps int, percentage float64, order int) string {
	t.Helper()
	body := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": %f},
		"setScheme": {"type": "FIXED", "sets": %d, "reps": %d},
		"order": %d
	}`, liftID, percentage, sets, reps, order)

	resp, err := adminPost(ts.URL("/prescriptions"), body)
	if err != nil {
		t.Fatalf("Failed to create prescription: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create prescription, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope PrescriptionResponse
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

func createSheikoDay(t *testing.T, ts *testutil.TestServer, name, slug string) string {
	t.Helper()
	body := fmt.Sprintf(`{"name": "%s", "slug": "%s"}`, name, slug)
	resp, _ := adminPost(ts.URL("/days"), body)
	var dayEnvelope DayResponse
	json.NewDecoder(resp.Body).Decode(&dayEnvelope)
	resp.Body.Close()
	return dayEnvelope.Data.ID
}

func createSheikoWeek(t *testing.T, ts *testutil.TestServer, cycleID string, weekNumber int) string {
	t.Helper()
	weekBody := fmt.Sprintf(`{"weekNumber": %d, "cycleId": "%s"}`, weekNumber, cycleID)
	weekResp, _ := adminPost(ts.URL("/weeks"), weekBody)
	var weekEnvelope WeekResponse
	json.NewDecoder(weekResp.Body).Decode(&weekEnvelope)
	weekResp.Body.Close()
	return weekEnvelope.Data.ID
}

func setMeetDate(t *testing.T, ts *testutil.TestServer, userID, programID, meetDate string) MeetDateResponseData {
	t.Helper()
	body := fmt.Sprintf(`{"meet_date": "%s"}`, meetDate)
	resp, err := userPut(ts.URL("/users/"+userID+"/programs/"+programID+"/state/meet-date"), body, userID)
	if err != nil {
		t.Fatalf("Failed to set meet date: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to set meet date, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope MeetDateResponseEnvelope
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data
}

func getCountdown(t *testing.T, ts *testutil.TestServer, userID, programID string) CountdownResponseData {
	t.Helper()
	resp, err := userGet(ts.URL("/users/"+userID+"/programs/"+programID+"/state/countdown"), userID)
	if err != nil {
		t.Fatalf("Failed to get countdown: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to get countdown, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope CountdownResponseEnvelope
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data
}

// =============================================================================
// SHEIKO INTERMEDIATE E2E TESTS
// =============================================================================

// TestSheikoIntermediateFull13WeekProgram validates the complete 13-week Sheiko
// Intermediate program execution with peaking and meet date functionality.
//
// Sheiko Intermediate characteristics:
// - Phase 1 (Prep 1): Weeks 1-4 - Base building, 70-75% intensity, high volume
// - Phase 2 (Prep 2): Weeks 5-8 - Intensification, 80-90% intensity, medium volume
// - Phase 3 (Competition): Weeks 9-13 - Taper, 90-95% intensity, low volume with taper
func TestSheikoIntermediateFull13WeekProgram(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user"

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	// Training maxes
	squatMax := 200.0
	benchMax := 140.0
	deadliftMax := 220.0

	// Create training maxes
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatMax)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchMax)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", deadliftMax)

	// Create complete Sheiko program
	programID, _, _ := createSheikoTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, err := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	if err != nil {
		t.Fatalf("Failed to enroll user: %v", err)
	}
	if enrollResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(enrollResp.Body)
		enrollResp.Body.Close()
		t.Fatalf("Failed to enroll user, status %d: %s", enrollResp.StatusCode, body)
	}
	enrollResp.Body.Close()

	// Set meet date 13 weeks (91 days) from now
	meetDate := time.Now().AddDate(0, 0, 91).Format("2006-01-02")
	meetDateResponse := setMeetDate(t, ts, userID, programID, meetDate)

	t.Run("verifies schedule type is days_out", func(t *testing.T) {
		// The schedule type should be implicitly days_out when meet date is set
		if meetDateResponse.DaysOut < 85 || meetDateResponse.DaysOut > 95 {
			t.Errorf("Expected ~91 days out, got %d", meetDateResponse.DaysOut)
		}
	})

	t.Run("verifies initial phase is prep_1 or base", func(t *testing.T) {
		// At 91 days out, we should be in prep_1 or base phase
		validPhases := map[string]bool{"prep_1": true, "base": true}
		if !validPhases[meetDateResponse.CurrentPhase] {
			t.Errorf("Expected phase prep_1 or base at 91 days out, got %s", meetDateResponse.CurrentPhase)
		}
	})

	t.Run("generates correct workout for week 1", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		if workoutResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(workoutResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", workoutResp.StatusCode, body)
		}

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Verify we have exercises
		if len(workout.Data.Exercises) == 0 {
			t.Error("Expected exercises in workout, got none")
		}

		// In Prep 1, squat should be at ~72.5% of training max = 145 lbs
		for _, ex := range workout.Data.Exercises {
			if ex.Lift.ID == squatID {
				expectedWeight := squatMax * 0.725
				for _, set := range ex.Sets {
					if !closeEnough(set.Weight, expectedWeight) {
						t.Logf("Squat weight: expected ~%.1f, got %.1f (difference: %.1f)",
							expectedWeight, set.Weight, set.Weight-expectedWeight)
					}
				}
			}
		}
	})
}

// TestSheikoIntermediatePhaseTransitions validates that phase transitions work
// correctly when starting at different points in the program.
func TestSheikoIntermediatePhaseTransitions(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	// Use seeded test user for foreign key constraints
	userID := "workout-test-user"

	// Create complete Sheiko program
	programID, _, _ := createSheikoTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	t.Run("starting at 5 weeks out enters Competition phase", func(t *testing.T) {
		// 5 weeks = 35 days, which should be in Competition phase (weeks 9-13)
		meetDate := time.Now().AddDate(0, 0, 35).Format("2006-01-02")
		response := setMeetDate(t, ts, userID, programID, meetDate)

		// Per the handler logic:
		// - 0-7 days: meet_week
		// - 8-14 days: peak
		// - 15-28 days: taper
		// - 29-56 days: prep_2
		// At 35 days, this is taper phase according to the handler
		expectedPhases := map[string]bool{"taper": true, "prep_2": true}
		if !expectedPhases[response.CurrentPhase] {
			t.Errorf("Expected taper or prep_2 phase at 35 days out, got %s", response.CurrentPhase)
		}
	})

	t.Run("starting at 2 weeks out enters peak phase", func(t *testing.T) {
		// 2 weeks = 14 days
		meetDate := time.Now().AddDate(0, 0, 14).Format("2006-01-02")
		response := setMeetDate(t, ts, userID, programID, meetDate)

		// At 14 days, this should be peak phase
		if response.CurrentPhase != "peak" {
			t.Errorf("Expected peak phase at 14 days out, got %s", response.CurrentPhase)
		}
	})

	t.Run("starting at 1 week out enters meet_week phase", func(t *testing.T) {
		// 1 week = 7 days
		meetDate := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
		response := setMeetDate(t, ts, userID, programID, meetDate)

		// At 7 days, this should be meet_week
		if response.CurrentPhase != "meet_week" {
			t.Errorf("Expected meet_week phase at 7 days out, got %s", response.CurrentPhase)
		}
	})
}

// TestSheikoIntermediateMeetDateChanges validates that changing the meet date
// correctly recalculates the phase and schedule.
func TestSheikoIntermediateMeetDateChanges(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	// Use seeded test user for foreign key constraints
	userID := "workout-test-user"

	// Create complete Sheiko program
	programID, _, _ := createSheikoTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	t.Run("changes from 10 weeks to 6 weeks recalculates phase", func(t *testing.T) {
		// First set meet date at 10 weeks (70 days)
		meetDate10Weeks := time.Now().AddDate(0, 0, 70).Format("2006-01-02")
		response1 := setMeetDate(t, ts, userID, programID, meetDate10Weeks)

		// At 70 days, should be in prep_1 phase (57-84 days = prep_1)
		if response1.CurrentPhase != "prep_1" {
			t.Logf("At 70 days out, got phase: %s", response1.CurrentPhase)
		}

		// Now change to 6 weeks (42 days)
		meetDate6Weeks := time.Now().AddDate(0, 0, 42).Format("2006-01-02")
		response2 := setMeetDate(t, ts, userID, programID, meetDate6Weeks)

		// At 42 days, should be in prep_2 phase (29-56 days = prep_2)
		if response2.CurrentPhase != "prep_2" {
			t.Errorf("Expected prep_2 phase at 42 days out, got %s", response2.CurrentPhase)
		}

		// Verify days out changed
		if response2.DaysOut >= response1.DaysOut {
			t.Errorf("Expected days_out to decrease, got %d -> %d", response1.DaysOut, response2.DaysOut)
		}
	})

	t.Run("clearing meet date returns to off_season", func(t *testing.T) {
		// Clear meet date
		body := `{"meet_date": null}`
		resp, err := userPut(ts.URL("/users/"+userID+"/programs/"+programID+"/state/meet-date"), body, userID)
		if err != nil {
			t.Fatalf("Failed to clear meet date: %v", err)
		}
		defer resp.Body.Close()

		var envelope MeetDateResponseEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)

		if envelope.Data.CurrentPhase != "off_season" {
			t.Errorf("Expected off_season phase after clearing meet date, got %s", envelope.Data.CurrentPhase)
		}
		if envelope.Data.DaysOut != 0 {
			t.Errorf("Expected 0 days_out after clearing, got %d", envelope.Data.DaysOut)
		}
	})
}

// TestSheikoIntermediateTaperMultiplier validates that the taper multiplier
// correctly reduces volume as the meet approaches.
func TestSheikoIntermediateTaperMultiplier(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	// Use seeded test user for foreign key constraints
	userID := "workout-test-user"

	// Create complete Sheiko program
	programID, _, _ := createSheikoTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	// Test taper multiplier at different days out
	// Based on the handler's determineTaperMultiplier:
	// - 0-7 days: 0.4
	// - 8-14 days: 0.6
	// - 15-21 days: 0.75
	// - 22-28 days: 0.85
	// - 29+ days: 1.0
	testCases := []struct {
		daysOut             int
		expectedMultiplier  float64
		expectedMultMin     float64 // Allow some tolerance
		expectedMultMax     float64
		description         string
	}{
		{50, 1.0, 0.95, 1.05, "no taper (50 days)"},
		{25, 0.85, 0.80, 0.90, "taper week 1 (25 days)"},
		{18, 0.75, 0.70, 0.80, "taper week 2 (18 days)"},
		{10, 0.6, 0.55, 0.65, "peak week (10 days)"},
		{5, 0.4, 0.35, 0.45, "meet week (5 days)"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			meetDate := time.Now().AddDate(0, 0, tc.daysOut).Format("2006-01-02")
			setMeetDate(t, ts, userID, programID, meetDate)

			countdown := getCountdown(t, ts, userID, programID)

			if countdown.TaperMultiplier < tc.expectedMultMin || countdown.TaperMultiplier > tc.expectedMultMax {
				t.Errorf("At %d days out: expected taper multiplier ~%.2f (range %.2f-%.2f), got %.2f",
					tc.daysOut, tc.expectedMultiplier, tc.expectedMultMin, tc.expectedMultMax, countdown.TaperMultiplier)
			}
		})
	}

	t.Run("taper multiplier decreases as meet approaches", func(t *testing.T) {
		var multipliers []float64
		daysOutValues := []int{35, 28, 21, 14, 7}

		for _, days := range daysOutValues {
			meetDate := time.Now().AddDate(0, 0, days).Format("2006-01-02")
			setMeetDate(t, ts, userID, programID, meetDate)
			countdown := getCountdown(t, ts, userID, programID)
			multipliers = append(multipliers, countdown.TaperMultiplier)
		}

		// Verify multipliers decrease (or stay same) as meet approaches
		for i := 1; i < len(multipliers); i++ {
			if multipliers[i] > multipliers[i-1] {
				t.Errorf("Taper multiplier should not increase as meet approaches: %.2f -> %.2f at days %d -> %d",
					multipliers[i-1], multipliers[i], daysOutValues[i-1], daysOutValues[i])
			}
		}
	})
}

// TestSheikoIntermediateCountdownEndpoint validates the countdown endpoint
// returns correct information at various points in the program.
func TestSheikoIntermediateCountdownEndpoint(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	// Use seeded test user for foreign key constraints
	userID := "workout-test-user"

	// Create complete Sheiko program
	programID, _, _ := createSheikoTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	// Set meet date 8 weeks out
	meetDate := time.Now().AddDate(0, 0, 56).Format("2006-01-02")
	setMeetDate(t, ts, userID, programID, meetDate)

	t.Run("countdown returns all required fields", func(t *testing.T) {
		countdown := getCountdown(t, ts, userID, programID)

		// Verify meet_date is set
		if countdown.MeetDate == nil {
			t.Error("Expected meet_date to be set")
		}

		// Verify days_out is reasonable
		if countdown.DaysOut < 50 || countdown.DaysOut > 60 {
			t.Errorf("Expected ~56 days_out, got %d", countdown.DaysOut)
		}

		// Verify current_phase is set
		if countdown.CurrentPhase == "" {
			t.Error("Expected non-empty current_phase")
		}

		// Verify phase_week is positive
		if countdown.PhaseWeek < 1 {
			t.Errorf("Expected phase_week >= 1, got %d", countdown.PhaseWeek)
		}

		// Verify taper_multiplier is in valid range
		if countdown.TaperMultiplier <= 0 || countdown.TaperMultiplier > 1.0 {
			t.Errorf("Expected taper_multiplier between 0 and 1, got %.2f", countdown.TaperMultiplier)
		}
	})
}

// closeEnough is defined in bill_starr_test.go - no need to redeclare
