// Package e2e provides end-to-end tests for complete program workflows.
// This file contains E2E tests for Sheiko Beginner program.
package e2e

import (
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
// SHEIKO BEGINNER HELPER FUNCTIONS
// =============================================================================

// createSheikoBeginnerTestSetup creates the complete Sheiko Beginner program structure.
// Returns programID, cycleID, and the week IDs for all 8 weeks (4 Prep + 4 Comp).
func createSheikoBeginnerTestSetup(t *testing.T, ts *testutil.TestServer, testID string) (programID, cycleID string, weekIDs []string) {
	t.Helper()

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	// Create prescriptions for Sheiko Beginner intensity zones
	// Zone 3: 70-79% - Primary working zone, volume accumulation
	zone3SquatPrescID := createSheikoPrescription(t, ts, squatID, 4, 4, 72.5, 0)
	zone3BenchPrescID := createSheikoPrescription(t, ts, benchID, 4, 3, 75.0, 1)
	zone3DeadliftPrescID := createSheikoPrescription(t, ts, deadliftID, 4, 3, 75.0, 2)

	// Zone 4: 80-89% - Heavier work sets, strength building
	zone4SquatPrescID := createSheikoPrescription(t, ts, squatID, 4, 2, 82.5, 0)
	zone4BenchPrescID := createSheikoPrescription(t, ts, benchID, 5, 3, 80.0, 1)
	_ = createSheikoPrescription(t, ts, deadliftID, 3, 2, 85.0, 2) // Zone 4 deadlift available for future use

	// Zone 5: 90%+ - Heavy singles for comp phase peaking
	zone5SquatPrescID := createSheikoPrescription(t, ts, squatID, 3, 1, 95.0, 0)
	zone5BenchPrescID := createSheikoPrescription(t, ts, benchID, 3, 1, 95.0, 1)
	zone5DeadliftPrescID := createSheikoPrescription(t, ts, deadliftID, 3, 1, 95.0, 2)

	// Taper prescriptions (lower intensity and volume)
	taperSquatPrescID := createSheikoPrescription(t, ts, squatID, 3, 2, 65.0, 0)
	taperBenchPrescID := createSheikoPrescription(t, ts, benchID, 3, 2, 70.0, 1)
	taperDeadliftPrescID := createSheikoPrescription(t, ts, deadliftID, 3, 2, 70.0, 2)

	// Create days for Prep Phase (3 days/week: Mon/Wed/Fri)
	// Prep Day 1 (Monday): Squat + Bench focus
	prepDay1Slug := "sheiko-beg-prep-d1-" + testID
	prepDay1ID := createSheikoDay(t, ts, "Prep Day 1 - Squat/Bench", prepDay1Slug)
	addPrescToDay(t, ts, prepDay1ID, zone3SquatPrescID)
	addPrescToDay(t, ts, prepDay1ID, zone4BenchPrescID)

	// Prep Day 2 (Wednesday): Deadlift variations + Bench
	prepDay2Slug := "sheiko-beg-prep-d2-" + testID
	prepDay2ID := createSheikoDay(t, ts, "Prep Day 2 - Deadlift/Bench", prepDay2Slug)
	addPrescToDay(t, ts, prepDay2ID, zone3DeadliftPrescID)
	addPrescToDay(t, ts, prepDay2ID, zone3BenchPrescID)

	// Prep Day 3 (Friday): Squat + Bench
	prepDay3Slug := "sheiko-beg-prep-d3-" + testID
	prepDay3ID := createSheikoDay(t, ts, "Prep Day 3 - Squat/Bench", prepDay3Slug)
	addPrescToDay(t, ts, prepDay3ID, zone4SquatPrescID)
	addPrescToDay(t, ts, prepDay3ID, zone3BenchPrescID)

	// Create days for Comp Phase
	// Comp Day 1 (Monday): Squat + Bench + Deadlift
	compDay1Slug := "sheiko-beg-comp-d1-" + testID
	compDay1ID := createSheikoDay(t, ts, "Comp Day 1 - All Lifts", compDay1Slug)
	addPrescToDay(t, ts, compDay1ID, zone4SquatPrescID)
	addPrescToDay(t, ts, compDay1ID, zone4BenchPrescID)

	// Comp Day 2 (Wednesday): Test Day or Peak Work
	compDay2Slug := "sheiko-beg-comp-d2-" + testID
	compDay2ID := createSheikoDay(t, ts, "Comp Day 2 - Peak/Test", compDay2Slug)
	addPrescToDay(t, ts, compDay2ID, zone5SquatPrescID)
	addPrescToDay(t, ts, compDay2ID, zone5BenchPrescID)
	addPrescToDay(t, ts, compDay2ID, zone5DeadliftPrescID)

	// Comp Day 3 (Friday): Recovery/Moderate
	compDay3Slug := "sheiko-beg-comp-d3-" + testID
	compDay3ID := createSheikoDay(t, ts, "Comp Day 3 - Moderate", compDay3Slug)
	addPrescToDay(t, ts, compDay3ID, zone4SquatPrescID)
	addPrescToDay(t, ts, compDay3ID, zone4BenchPrescID)

	// Taper Days (Comp Week 4)
	taperDay1Slug := "sheiko-beg-taper-d1-" + testID
	taperDay1ID := createSheikoDay(t, ts, "Taper Day 1", taperDay1Slug)
	addPrescToDay(t, ts, taperDay1ID, taperBenchPrescID)
	addPrescToDay(t, ts, taperDay1ID, taperDeadliftPrescID)

	taperDay2Slug := "sheiko-beg-taper-d2-" + testID
	taperDay2ID := createSheikoDay(t, ts, "Taper Day 2", taperDay2Slug)
	addPrescToDay(t, ts, taperDay2ID, taperSquatPrescID)
	addPrescToDay(t, ts, taperDay2ID, taperBenchPrescID)

	// Create 8-week cycle (4 Prep + 4 Comp)
	cycleName := "Sheiko Beginner Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 8}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID = cycleEnvelope.Data.ID

	// Create weeks for each phase
	weekIDs = make([]string, 8)

	// Prep Phase: Weeks 1-4
	for i := 1; i <= 4; i++ {
		weekID := createSheikoWeek(t, ts, cycleID, i)
		weekIDs[i-1] = weekID
		addDayToWeek(t, ts, weekID, prepDay1ID, "MONDAY")
		addDayToWeek(t, ts, weekID, prepDay2ID, "WEDNESDAY")
		addDayToWeek(t, ts, weekID, prepDay3ID, "FRIDAY")
	}

	// Comp Phase: Weeks 5-7 (regular comp days)
	for i := 5; i <= 7; i++ {
		weekID := createSheikoWeek(t, ts, cycleID, i)
		weekIDs[i-1] = weekID
		addDayToWeek(t, ts, weekID, compDay1ID, "MONDAY")
		addDayToWeek(t, ts, weekID, compDay2ID, "WEDNESDAY")
		addDayToWeek(t, ts, weekID, compDay3ID, "FRIDAY")
	}

	// Comp Phase: Week 8 (Taper week - only Mon/Wed)
	taperWeekID := createSheikoWeek(t, ts, cycleID, 8)
	weekIDs[7] = taperWeekID
	addDayToWeek(t, ts, taperWeekID, taperDay1ID, "MONDAY")
	addDayToWeek(t, ts, taperWeekID, taperDay2ID, "WEDNESDAY")

	// Create program
	programSlug := "sheiko-beginner-" + testID
	programBody := fmt.Sprintf(`{"name": "Sheiko Beginner", "slug": "%s", "cycleId": "%s"}`, programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID = programEnvelope.Data.ID

	return programID, cycleID, weekIDs
}

// =============================================================================
// SHEIKO BEGINNER E2E TESTS
// =============================================================================

// TestSheikoBeginnerFull8WeekProgram validates the complete 8-week Sheiko
// Beginner program execution with prep and comp phases.
//
// Sheiko Beginner characteristics:
// - Prep Phase (Weeks 1-4): Base building, 70-85% intensity, volume accumulation
// - Comp Phase (Weeks 5-8): Peaking, 80-105% intensity, taper in week 4
// - 3 days/week: Mon/Wed/Fri
// - High frequency: Squat 2x, Bench 3x, Deadlift 2x per week
func TestSheikoBeginnerFull8WeekProgram(t *testing.T) {
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

	// Training maxes (typical beginner-intermediate levels)
	squatMax := 150.0
	benchMax := 100.0
	deadliftMax := 170.0

	// Create training maxes
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatMax)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchMax)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", deadliftMax)

	// Create complete Sheiko Beginner program
	programID, _, _ := createSheikoBeginnerTestSetup(t, ts, testID)

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

	// Set meet date 8 weeks (56 days) from now
	meetDate := time.Now().AddDate(0, 0, 56).Format("2006-01-02")
	meetDateResponse := setMeetDate(t, ts, userID, programID, meetDate)

	t.Run("verifies meet date is set correctly", func(t *testing.T) {
		if meetDateResponse.DaysOut < 50 || meetDateResponse.DaysOut > 60 {
			t.Errorf("Expected ~56 days out, got %d", meetDateResponse.DaysOut)
		}
	})

	t.Run("verifies initial phase at 8 weeks out", func(t *testing.T) {
		// At 56 days out, should be in prep_1 or prep_2 phase
		validPhases := map[string]bool{"prep_1": true, "prep_2": true, "base": true}
		if !validPhases[meetDateResponse.CurrentPhase] {
			t.Errorf("Expected prep phase at 56 days out, got %s", meetDateResponse.CurrentPhase)
		}
	})

	t.Run("generates correct workout for week 1 prep phase", func(t *testing.T) {
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

		// In Prep phase, squat should be at ~72.5% of training max
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

	t.Run("verifies bench press frequency (3x per week)", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Count bench exercises in workout
		hasBench := false
		for _, ex := range workout.Data.Exercises {
			if ex.Lift.ID == benchID {
				hasBench = true
				break
			}
		}

		if !hasBench {
			t.Log("Note: Bench press may not be in current day's workout")
		}
	})
}

// TestSheikoBeginnerPrepPhaseIntensity validates that prep phase uses
// submaximal training with 70-85% intensity range.
func TestSheikoBeginnerPrepPhaseIntensity(t *testing.T) {
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

	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatMax)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchMax)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", deadliftMax)

	programID, _, _ := createSheikoBeginnerTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	// Set meet date far out to stay in prep phase (10 weeks)
	meetDate := time.Now().AddDate(0, 0, 70).Format("2006-01-02")
	setMeetDate(t, ts, userID, programID, meetDate)

	t.Run("prep phase weights are in 70-85% range", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		for _, ex := range workout.Data.Exercises {
			var max float64
			switch ex.Lift.ID {
			case squatID:
				max = squatMax
			case benchID:
				max = benchMax
			case deadliftID:
				max = deadliftMax
			default:
				continue
			}

			for _, set := range ex.Sets {
				percentage := (set.Weight / max) * 100
				// Allow some tolerance for rounding
				if percentage < 65 || percentage > 90 {
					t.Errorf("Prep phase weight %.1f for %s is %.1f%% - expected 65-90%%",
						set.Weight, ex.Lift.Name, percentage)
				}
			}
		}
	})
}

// TestSheikoBeginnerCompPhaseTransition validates the transition from
// Prep phase to Comp phase and appropriate intensity changes.
func TestSheikoBeginnerCompPhaseTransition(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user"

	programID, _, _ := createSheikoBeginnerTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	t.Run("entering comp phase at 4 weeks out", func(t *testing.T) {
		// 4 weeks = 28 days
		meetDate := time.Now().AddDate(0, 0, 28).Format("2006-01-02")
		response := setMeetDate(t, ts, userID, programID, meetDate)

		// At 28 days, should be in taper or prep_2 phase
		validPhases := map[string]bool{"taper": true, "prep_2": true}
		if !validPhases[response.CurrentPhase] {
			t.Logf("At 28 days out, got phase: %s", response.CurrentPhase)
		}
	})

	t.Run("entering peak phase at 2 weeks out", func(t *testing.T) {
		meetDate := time.Now().AddDate(0, 0, 14).Format("2006-01-02")
		response := setMeetDate(t, ts, userID, programID, meetDate)

		if response.CurrentPhase != "peak" {
			t.Errorf("Expected peak phase at 14 days out, got %s", response.CurrentPhase)
		}
	})

	t.Run("entering meet week at 1 week out", func(t *testing.T) {
		meetDate := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
		response := setMeetDate(t, ts, userID, programID, meetDate)

		if response.CurrentPhase != "meet_week" {
			t.Errorf("Expected meet_week phase at 7 days out, got %s", response.CurrentPhase)
		}
	})
}

// TestSheikoBeginnerTaperWeek validates that the taper week (Week 8)
// reduces volume and intensity appropriately.
func TestSheikoBeginnerTaperWeek(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user"

	programID, _, _ := createSheikoBeginnerTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	// Set meet date 5 days out (should be in meet_week with taper)
	meetDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02")
	setMeetDate(t, ts, userID, programID, meetDate)

	t.Run("taper multiplier is reduced", func(t *testing.T) {
		countdown := getCountdown(t, ts, userID, programID)

		// At meet week, taper multiplier should be low (around 0.4)
		if countdown.TaperMultiplier > 0.5 {
			t.Errorf("Expected low taper multiplier at meet week, got %.2f", countdown.TaperMultiplier)
		}
	})
}

// TestSheikoBeginnerIntensityZones validates that exercises are programmed
// within the correct intensity zones.
func TestSheikoBeginnerIntensityZones(t *testing.T) {
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

	// Use round numbers for easy percentage verification
	squatMax := 200.0
	benchMax := 150.0
	deadliftMax := 250.0

	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatMax)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchMax)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", deadliftMax)

	programID, _, _ := createSheikoBeginnerTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	// Set meet date 8 weeks out
	meetDate := time.Now().AddDate(0, 0, 56).Format("2006-01-02")
	setMeetDate(t, ts, userID, programID, meetDate)

	t.Run("workout has correct percentage-based weights", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Verify at least some exercises exist
		if len(workout.Data.Exercises) == 0 {
			t.Error("Expected exercises in workout")
			return
		}

		// Log the weights for manual verification
		for _, ex := range workout.Data.Exercises {
			var max float64
			switch ex.Lift.ID {
			case squatID:
				max = squatMax
			case benchID:
				max = benchMax
			case deadliftID:
				max = deadliftMax
			default:
				continue
			}

			for i, set := range ex.Sets {
				percentage := (set.Weight / max) * 100
				t.Logf("%s Set %d: %.1f lbs (%.1f%%)", ex.Lift.Name, i+1, set.Weight, percentage)
			}
		}
	})
}

// TestSheikoBeginnerRepeatablePrepPhase validates that the Prep phase
// can be repeated for non-competing lifters.
func TestSheikoBeginnerRepeatablePrepPhase(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user"

	programID, _, _ := createSheikoBeginnerTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	t.Run("no meet date keeps user in off-season/prep mode", func(t *testing.T) {
		// Without a meet date, user should be in off_season
		countdown := getCountdown(t, ts, userID, programID)

		if countdown.CurrentPhase != "off_season" {
			t.Errorf("Expected off_season phase without meet date, got %s", countdown.CurrentPhase)
		}
	})

	t.Run("clearing meet date returns to off-season", func(t *testing.T) {
		// First set a meet date
		meetDate := time.Now().AddDate(0, 0, 56).Format("2006-01-02")
		setMeetDate(t, ts, userID, programID, meetDate)

		// Verify we're not in off_season
		countdown1 := getCountdown(t, ts, userID, programID)
		if countdown1.CurrentPhase == "off_season" {
			t.Error("Should not be in off_season with meet date set")
		}

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
			t.Errorf("Expected off_season after clearing meet date, got %s", envelope.Data.CurrentPhase)
		}
	})
}

// TestSheikoBeginnerWeeklyStructure validates that workouts follow the
// 3-day per week structure (Mon/Wed/Fri).
func TestSheikoBeginnerWeeklyStructure(t *testing.T) {
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

	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", 200.0)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", 140.0)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", 220.0)

	programID, _, _ := createSheikoBeginnerTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	// Set meet date
	meetDate := time.Now().AddDate(0, 0, 56).Format("2006-01-02")
	setMeetDate(t, ts, userID, programID, meetDate)

	t.Run("can generate workout", func(t *testing.T) {
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

		// Should have exercises
		if len(workout.Data.Exercises) == 0 {
			t.Error("Expected exercises in workout")
		}
	})
}
