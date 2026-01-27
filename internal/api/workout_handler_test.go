package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// Test response types for workout tests

// LiftTestResponse matches the API response format for a lift.
type LiftTestResponse struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Slug              string  `json:"slug"`
	IsCompetitionLift bool    `json:"isCompetitionLift"`
	ParentLiftID      *string `json:"parentLiftId"`
	CreatedAt         string  `json:"createdAt"`
	UpdatedAt         string  `json:"updatedAt"`
}

// LiftMaxTestResponse matches the API response format for a lift max.
type LiftMaxTestResponse struct {
	ID            string    `json:"id"`
	UserID        string    `json:"userId"`
	LiftID        string    `json:"liftId"`
	Type          string    `json:"type"`
	Value         float64   `json:"value"`
	EffectiveDate time.Time `json:"effectiveDate"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// PrescriptionTestResponse matches the API response format for a prescription.
type PrescriptionTestResponse struct {
	ID           string                 `json:"id"`
	LiftID       string                 `json:"liftId"`
	LoadStrategy map[string]interface{} `json:"loadStrategy"`
	SetScheme    map[string]interface{} `json:"setScheme"`
	Order        int                    `json:"order"`
	Notes        string                 `json:"notes,omitempty"`
	RestSeconds  *int                   `json:"restSeconds,omitempty"`
	CreatedAt    string                 `json:"createdAt"`
	UpdatedAt    string                 `json:"updatedAt"`
}

// DayTestResponse matches the API response format for a day.
type DayTestResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	ProgramID *string `json:"programId,omitempty"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

// WeekTestResponse matches the API response format for a week.
type WeekTestResponse struct {
	ID         string  `json:"id"`
	WeekNumber int     `json:"weekNumber"`
	Variant    *string `json:"variant,omitempty"`
	CycleID    string  `json:"cycleId"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

// Envelope types for decoding API responses in workout tests.
// All API responses are wrapped in {"data": ...} envelopes.

// LiftTestEnvelope wraps lift responses.
type LiftTestEnvelope struct {
	Data LiftTestResponse `json:"data"`
}

// LiftMaxTestEnvelope wraps lift max responses.
type LiftMaxTestEnvelope struct {
	Data LiftMaxTestResponse `json:"data"`
}

// PrescriptionTestEnvelope wraps prescription responses.
type PrescriptionTestEnvelope struct {
	Data PrescriptionTestResponse `json:"data"`
}

// DayTestEnvelope wraps day responses.
type DayTestEnvelope struct {
	Data DayTestResponse `json:"data"`
}

// WeekTestEnvelope wraps week responses.
type WeekTestEnvelope struct {
	Data WeekTestResponse `json:"data"`
}

// WorkoutTestEnvelope wraps workout responses.
type WorkoutTestEnvelope struct {
	Data WorkoutTestResponse `json:"data"`
}

// Helper functions for workout tests

func userPostLiftMax(url string, body string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

// WorkoutLiftTestResponse represents lift info in a workout response.
type WorkoutLiftTestResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// WorkoutSetTestResponse represents a set in a workout response.
type WorkoutSetTestResponse struct {
	SetNumber  int     `json:"setNumber"`
	Weight     float64 `json:"weight"`
	TargetReps int     `json:"targetReps"`
	IsWorkSet  bool    `json:"isWorkSet"`
}

// WorkoutExerciseTestResponse represents an exercise in a workout response.
type WorkoutExerciseTestResponse struct {
	PrescriptionID string                   `json:"prescriptionId"`
	Lift           WorkoutLiftTestResponse  `json:"lift"`
	Sets           []WorkoutSetTestResponse `json:"sets"`
	Notes          string                   `json:"notes,omitempty"`
	RestSeconds    *int                     `json:"restSeconds,omitempty"`
}

// WorkoutTestResponse represents the API response for a generated workout.
type WorkoutTestResponse struct {
	UserID         string                        `json:"userId"`
	ProgramID      string                        `json:"programId"`
	CycleIteration int                           `json:"cycleIteration"`
	WeekNumber     int                           `json:"weekNumber"`
	DaySlug        string                        `json:"daySlug"`
	Date           string                        `json:"date"`
	Exercises      []WorkoutExerciseTestResponse `json:"exercises"`
}

// Helper functions for workout tests

func userGetWorkout(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func adminGetWorkout(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// workoutTestSetup creates all the necessary entities for workout testing
type workoutTestSetup struct {
	LiftID         string
	MaxID          string
	PrescriptionID string
	DayID          string
	WeekID         string
	CycleID        string
	ProgramID      string
}

func setupWorkoutTest(t *testing.T, ts *testutil.TestServer, userID string) *workoutTestSetup {
	t.Helper()
	setup := &workoutTestSetup{}

	// Use unique slugs per test to avoid conflicts
	liftSlug := "back-squat-" + userID
	daySlug := "squat-day-" + userID
	programSlug := "workout-test-program-" + userID
	cycleName := "Test Cycle " + userID

	// 1. Create a lift
	liftBody := `{"name": "Back Squat", "slug": "` + liftSlug + `"}`
	resp, err := adminPost(ts.URL("/lifts"), liftBody)
	if err != nil {
		t.Fatalf("Failed to create lift: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create lift, status %d: %s", resp.StatusCode, body)
	}
	var liftEnvelope LiftTestEnvelope
	json.Unmarshal(body, &liftEnvelope)
	setup.LiftID = liftEnvelope.Data.ID

	// 2. Create a lift max for the user (1RM that results in TM=300)
	// Backend auto-calculates TM as 90% of 1RM, so 1RM = 300/0.9 = 333.33, rounded to 333.25
	maxBody := `{"liftId": "` + setup.LiftID + `", "type": "ONE_RM", "value": 333.25}`
	resp, err = userPostLiftMax(ts.URL("/users/"+userID+"/lift-maxes"), maxBody, userID)
	if err != nil {
		t.Fatalf("Failed to create lift max: %v", err)
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create lift max, status %d: %s", resp.StatusCode, body)
	}
	var liftMaxEnvelope LiftMaxTestEnvelope
	json.Unmarshal(body, &liftMaxEnvelope)
	setup.MaxID = liftMaxEnvelope.Data.ID

	// 3. Create a prescription with PERCENT_OF load strategy and FIXED set scheme
	prescriptionBody := `{
		"liftId": "` + setup.LiftID + `",
		"loadStrategy": {
			"type": "PERCENT_OF",
			"referenceType": "TRAINING_MAX",
			"percentage": 75.0
		},
		"setScheme": {
			"type": "FIXED",
			"sets": 5,
			"reps": 5
		},
		"order": 0,
		"notes": "Focus on form",
		"restSeconds": 180
	}`
	resp, err = adminPost(ts.URL("/prescriptions"), prescriptionBody)
	if err != nil {
		t.Fatalf("Failed to create prescription: %v", err)
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create prescription, status %d: %s", resp.StatusCode, body)
	}
	var prescriptionEnvelope PrescriptionTestEnvelope
	json.Unmarshal(body, &prescriptionEnvelope)
	setup.PrescriptionID = prescriptionEnvelope.Data.ID

	// 4. Create a day
	dayBody := `{"name": "Squat Day", "slug": "` + daySlug + `"}`
	resp, err = adminPost(ts.URL("/days"), dayBody)
	if err != nil {
		t.Fatalf("Failed to create day: %v", err)
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create day, status %d: %s", resp.StatusCode, body)
	}
	var dayEnvelope DayTestEnvelope
	json.Unmarshal(body, &dayEnvelope)
	setup.DayID = dayEnvelope.Data.ID

	// 5. Add prescription to day
	addPrescriptionBody := `{"prescriptionId": "` + setup.PrescriptionID + `"}`
	resp, err = adminPost(ts.URL("/days/"+setup.DayID+"/prescriptions"), addPrescriptionBody)
	if err != nil {
		t.Fatalf("Failed to add prescription to day: %v", err)
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to add prescription to day, status %d: %s", resp.StatusCode, body)
	}

	// 6. Create a cycle
	cycleBody := `{"name": "` + cycleName + `", "lengthWeeks": 4}`
	resp, err = adminPostCycle(ts.URL("/cycles"), cycleBody)
	if err != nil {
		t.Fatalf("Failed to create cycle: %v", err)
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create cycle, status %d: %s", resp.StatusCode, body)
	}
	var cycleEnvelope CycleEnvelopeInteg
	json.Unmarshal(body, &cycleEnvelope)
	setup.CycleID = cycleEnvelope.Data.ID

	// 7. Create a week in the cycle
	weekBody := `{"weekNumber": 1, "cycleId": "` + setup.CycleID + `"}`
	resp, err = adminPost(ts.URL("/weeks"), weekBody)
	if err != nil {
		t.Fatalf("Failed to create week: %v", err)
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create week, status %d: %s", resp.StatusCode, body)
	}
	var weekEnvelope WeekTestEnvelope
	json.Unmarshal(body, &weekEnvelope)
	setup.WeekID = weekEnvelope.Data.ID

	// 8. Add day to week
	addDayBody := `{"dayId": "` + setup.DayID + `", "dayOfWeek": "MONDAY"}`
	resp, err = adminPost(ts.URL("/weeks/"+setup.WeekID+"/days"), addDayBody)
	if err != nil {
		t.Fatalf("Failed to add day to week: %v", err)
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to add day to week, status %d: %s", resp.StatusCode, body)
	}

	// 9. Create a program
	programBody := `{"name": "Workout Test Program", "slug": "` + programSlug + `", "cycleId": "` + setup.CycleID + `"}`
	resp, err = adminPostProgram(ts.URL("/programs"), programBody)
	if err != nil {
		t.Fatalf("Failed to create program: %v", err)
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create program, status %d: %s", resp.StatusCode, body)
	}
	var programEnvelope ProgramEnvelopeInteg
	json.Unmarshal(body, &programEnvelope)
	setup.ProgramID = programEnvelope.Data.ID

	// 10. Enroll user in program
	enrollBody := `{"programId": "` + setup.ProgramID + `"}`
	resp, err = userPostEnrollment(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	if err != nil {
		t.Fatalf("Failed to enroll user: %v", err)
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to enroll user, status %d: %s", resp.StatusCode, body)
	}

	return setup
}

func TestWorkoutGenerate(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	userID := "workout-test-user"
	setup := setupWorkoutTest(t, ts, userID)

	daySlug := "squat-day-" + userID
	liftSlug := "back-squat-" + userID

	t.Run("generates workout for enrolled user", func(t *testing.T) {
		resp, err := userGetWorkout(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var workoutEnvelope WorkoutTestEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&workoutEnvelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		workout := workoutEnvelope.Data

		if workout.UserID != userID {
			t.Errorf("Expected userId %s, got %s", userID, workout.UserID)
		}
		if workout.ProgramID != setup.ProgramID {
			t.Errorf("Expected programId %s, got %s", setup.ProgramID, workout.ProgramID)
		}
		if workout.CycleIteration != 1 {
			t.Errorf("Expected cycleIteration 1, got %d", workout.CycleIteration)
		}
		if workout.WeekNumber != 1 {
			t.Errorf("Expected weekNumber 1, got %d", workout.WeekNumber)
		}
		if workout.DaySlug != daySlug {
			t.Errorf("Expected daySlug '%s', got %s", daySlug, workout.DaySlug)
		}
		if workout.Date == "" {
			t.Error("Expected non-empty date")
		}
		if len(workout.Exercises) != 1 {
			t.Errorf("Expected 1 exercise, got %d", len(workout.Exercises))
		}
	})

	t.Run("exercise contains correct information", func(t *testing.T) {
		resp, err := userGetWorkout(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var workoutEnvelope WorkoutTestEnvelope
		json.NewDecoder(resp.Body).Decode(&workoutEnvelope)
		workout := workoutEnvelope.Data

		if len(workout.Exercises) != 1 {
			t.Fatalf("Expected 1 exercise, got %d", len(workout.Exercises))
		}

		exercise := workout.Exercises[0]
		if exercise.PrescriptionID != setup.PrescriptionID {
			t.Errorf("Expected prescriptionId %s, got %s", setup.PrescriptionID, exercise.PrescriptionID)
		}
		if exercise.Lift.ID != setup.LiftID {
			t.Errorf("Expected lift.id %s, got %s", setup.LiftID, exercise.Lift.ID)
		}
		if exercise.Lift.Name != "Back Squat" {
			t.Errorf("Expected lift.name 'Back Squat', got %s", exercise.Lift.Name)
		}
		if exercise.Lift.Slug != liftSlug {
			t.Errorf("Expected lift.slug '%s', got %s", liftSlug, exercise.Lift.Slug)
		}
		if exercise.Notes != "Focus on form" {
			t.Errorf("Expected notes 'Focus on form', got %s", exercise.Notes)
		}
		if exercise.RestSeconds == nil || *exercise.RestSeconds != 180 {
			t.Errorf("Expected restSeconds 180, got %v", exercise.RestSeconds)
		}
	})

	t.Run("sets are generated correctly with 75% of 300 TM = 225", func(t *testing.T) {
		resp, err := userGetWorkout(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var workoutEnvelope WorkoutTestEnvelope
		json.NewDecoder(resp.Body).Decode(&workoutEnvelope)
		workout := workoutEnvelope.Data

		exercise := workout.Exercises[0]
		if len(exercise.Sets) != 5 {
			t.Fatalf("Expected 5 sets, got %d", len(exercise.Sets))
		}

		// 75% of 300 = 225
		expectedWeight := 225.0
		for i, set := range exercise.Sets {
			if set.SetNumber != i+1 {
				t.Errorf("Set %d: expected setNumber %d, got %d", i, i+1, set.SetNumber)
			}
			if set.Weight != expectedWeight {
				t.Errorf("Set %d: expected weight %f, got %f", i, expectedWeight, set.Weight)
			}
			if set.TargetReps != 5 {
				t.Errorf("Set %d: expected targetReps 5, got %d", i, set.TargetReps)
			}
			if !set.IsWorkSet {
				t.Errorf("Set %d: expected IsWorkSet true", i)
			}
		}
	})

	t.Run("accepts optional weekNumber parameter", func(t *testing.T) {
		// Note: weekNumber=1 is the only week we created
		resp, err := userGetWorkout(ts.URL("/users/"+userID+"/workout?weekNumber=1"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var workoutEnvelope WorkoutTestEnvelope
		json.NewDecoder(resp.Body).Decode(&workoutEnvelope)
		workout := workoutEnvelope.Data

		if workout.WeekNumber != 1 {
			t.Errorf("Expected weekNumber 1, got %d", workout.WeekNumber)
		}
	})

	t.Run("accepts optional daySlug parameter", func(t *testing.T) {
		resp, err := userGetWorkout(ts.URL("/users/"+userID+"/workout?daySlug="+daySlug), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var workoutEnvelope WorkoutTestEnvelope
		json.NewDecoder(resp.Body).Decode(&workoutEnvelope)
		workout := workoutEnvelope.Data

		if workout.DaySlug != daySlug {
			t.Errorf("Expected daySlug '%s', got %s", daySlug, workout.DaySlug)
		}
	})

	t.Run("accepts optional date parameter", func(t *testing.T) {
		resp, err := userGetWorkout(ts.URL("/users/"+userID+"/workout?date=2024-01-15"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var workoutEnvelope WorkoutTestEnvelope
		json.NewDecoder(resp.Body).Decode(&workoutEnvelope)
		workout := workoutEnvelope.Data

		if workout.Date != "2024-01-15" {
			t.Errorf("Expected date '2024-01-15', got %s", workout.Date)
		}
	})
}

func TestWorkoutGenerateErrors(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("returns 404 for non-enrolled user", func(t *testing.T) {
		resp, err := userGetWorkout(ts.URL("/users/non-enrolled-user/workout"), "non-enrolled-user")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid weekNumber", func(t *testing.T) {
		// Set up a user for this test
		userID := "workout-error-test-user"
		setupWorkoutTest(t, ts, userID)

		resp, err := userGetWorkout(ts.URL("/users/"+userID+"/workout?weekNumber=0"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for non-existent week", func(t *testing.T) {
		userID := "workout-error-test-user-2"
		setupWorkoutTest(t, ts, userID)

		// Week 99 doesn't exist in our setup
		resp, err := userGetWorkout(ts.URL("/users/"+userID+"/workout?weekNumber=99"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for non-existent day slug", func(t *testing.T) {
		userID := "workout-error-test-user-3"
		setupWorkoutTest(t, ts, userID)

		resp, err := userGetWorkout(ts.URL("/users/"+userID+"/workout?daySlug=non-existent-day"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestWorkoutPreview(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	userID := "workout-preview-test-user"
	setup := setupWorkoutTest(t, ts, userID)
	daySlug := "squat-day-" + userID

	t.Run("previews workout with required parameters", func(t *testing.T) {
		resp, err := userGetWorkout(ts.URL("/users/"+userID+"/workout/preview?week=1&day="+daySlug), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var workoutEnvelope WorkoutTestEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&workoutEnvelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		workout := workoutEnvelope.Data

		if workout.UserID != userID {
			t.Errorf("Expected userId %s, got %s", userID, workout.UserID)
		}
		if workout.ProgramID != setup.ProgramID {
			t.Errorf("Expected programId %s, got %s", setup.ProgramID, workout.ProgramID)
		}
		if workout.WeekNumber != 1 {
			t.Errorf("Expected weekNumber 1, got %d", workout.WeekNumber)
		}
		if workout.DaySlug != daySlug {
			t.Errorf("Expected daySlug '%s', got %s", daySlug, workout.DaySlug)
		}
		if len(workout.Exercises) != 1 {
			t.Errorf("Expected 1 exercise, got %d", len(workout.Exercises))
		}
	})

	t.Run("returns 400 when week is missing", func(t *testing.T) {
		resp, err := userGetWorkout(ts.URL("/users/"+userID+"/workout/preview?day="+daySlug), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 when day is missing", func(t *testing.T) {
		resp, err := userGetWorkout(ts.URL("/users/"+userID+"/workout/preview?week=1"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid week", func(t *testing.T) {
		resp, err := userGetWorkout(ts.URL("/users/"+userID+"/workout/preview?week=0&day=squat-day"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 404 for non-enrolled user", func(t *testing.T) {
		resp, err := userGetWorkout(ts.URL("/users/non-enrolled-preview/workout/preview?week=1&day=squat-day"), "non-enrolled-preview")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestWorkoutAuthorization(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	userID := "workout-auth-test-user"
	setupWorkoutTest(t, ts, userID)
	daySlug := "squat-day-" + userID

	t.Run("unauthenticated user gets 401 on generate", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/users/"+userID+"/workout"), nil)
		// No X-User-ID header

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("unauthenticated user gets 401 on preview", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/users/"+userID+"/workout/preview?week=1&day="+daySlug), nil)
		// No X-User-ID header

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("user cannot generate another user's workout", func(t *testing.T) {
		resp, err := userGetWorkout(ts.URL("/users/"+userID+"/workout"), "other-user")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("user cannot preview another user's workout", func(t *testing.T) {
		resp, err := userGetWorkout(ts.URL("/users/"+userID+"/workout/preview?week=1&day="+daySlug), "other-user")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("admin can generate any user's workout", func(t *testing.T) {
		resp, err := adminGetWorkout(ts.URL("/users/" + userID + "/workout"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("admin can preview any user's workout", func(t *testing.T) {
		resp, err := adminGetWorkout(ts.URL("/users/" + userID + "/workout/preview?week=1&day=" + daySlug))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})
}

func TestWorkoutResponseFormat(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	userID := "workout-format-test-user"
	setupWorkoutTest(t, ts, userID)

	t.Run("response has correct JSON field names", func(t *testing.T) {
		resp, err := userGetWorkout(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		expectedFields := []string{
			`"userId"`,
			`"programId"`,
			`"cycleIteration"`,
			`"weekNumber"`,
			`"daySlug"`,
			`"date"`,
			`"exercises"`,
			`"prescriptionId"`,
			`"lift"`,
			`"sets"`,
			`"setNumber"`,
			`"weight"`,
			`"targetReps"`,
			`"isWorkSet"`,
		}

		for _, field := range expectedFields {
			if !bytes.Contains(body, []byte(field)) {
				t.Errorf("Expected field %s in response, body: %s", field, string(body))
			}
		}
	})
}
