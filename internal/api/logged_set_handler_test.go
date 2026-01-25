package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// LoggedSetTestResponse represents a logged set in the API response.
type LoggedSetTestResponse struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	SessionID      string    `json:"sessionId"`
	PrescriptionID string    `json:"prescriptionId"`
	LiftID         string    `json:"liftId"`
	SetNumber      int       `json:"setNumber"`
	Weight         float64   `json:"weight"`
	TargetReps     int       `json:"targetReps"`
	RepsPerformed  int       `json:"repsPerformed"`
	IsAMRAP        bool      `json:"isAmrap"`
	CreatedAt      time.Time `json:"createdAt"`
}

// LoggedSetTestListResponse wraps an array of logged sets.
type LoggedSetTestListResponse struct {
	Data []LoggedSetTestResponse `json:"data"`
}

// LoggedSetTestPaginatedResponse wraps paginated logged set responses.
type LoggedSetTestPaginatedResponse struct {
	Data []LoggedSetTestResponse `json:"data"`
	Meta *struct {
		Total   int64 `json:"total"`
		Limit   int   `json:"limit"`
		Offset  int   `json:"offset"`
		HasMore bool  `json:"hasMore"`
	} `json:"meta"`
}

// LSWorkoutSessionEnvelope wraps workout session response.
type LSWorkoutSessionEnvelope struct {
	Data struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"data"`
}

// Helper functions for logged set tests

func authPostLoggedSets(url string, body string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func authGetLoggedSets(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func adminGetLoggedSets(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// createLSTestUser creates a test user in the database
func createLSTestUser(t *testing.T, ts *testutil.TestServer, userID string) {
	t.Helper()
	now := time.Now().Format(time.RFC3339)
	_, err := ts.DB().Exec("INSERT OR IGNORE INTO users (id, created_at, updated_at) VALUES (?, ?, ?)", userID, now, now)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
}

// createLSTestLift creates a test lift and returns its ID
func createLSTestLift(t *testing.T, ts *testutil.TestServer, name, slug string) string {
	t.Helper()
	body := `{"name": "` + name + `", "slug": "` + slug + `", "isCompetitionLift": true}`
	resp, err := adminPost(ts.URL("/lifts"), body)
	if err != nil {
		t.Fatalf("Failed to create test lift: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create test lift (status %d): %s", resp.StatusCode, bodyBytes)
	}

	var lift LiftResponse
	json.NewDecoder(resp.Body).Decode(&lift)
	return lift.Data.ID
}

// createLSTestCycle creates a test cycle and returns its ID
func createLSTestCycle(t *testing.T, ts *testutil.TestServer, name string) string {
	t.Helper()
	body := `{"name": "` + name + `", "lengthWeeks": 4}`
	req, err := http.NewRequest(http.MethodPost, ts.URL("/cycles"), bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to create test cycle: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create cycle (status %d): %s", resp.StatusCode, bodyBytes)
	}

	var envelope struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// createLSTestProgram creates a test program and returns its ID
func createLSTestProgram(t *testing.T, ts *testutil.TestServer, name, slug, cycleID string) string {
	t.Helper()
	body := `{"name": "` + name + `", "slug": "` + slug + `", "cycleId": "` + cycleID + `"}`
	req, err := http.NewRequest(http.MethodPost, ts.URL("/programs"), bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to create test program: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create program (status %d): %s", resp.StatusCode, bodyBytes)
	}

	var envelope struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// enrollLSTestUser enrolls a user in a program
func enrollLSTestUser(t *testing.T, ts *testutil.TestServer, userID, programID string) {
	t.Helper()
	body := `{"programId": "` + programID + `"}`
	req, err := http.NewRequest(http.MethodPost, ts.URL("/users/"+userID+"/program"), bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to enroll user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 201 for enrollment, got %d: %s", resp.StatusCode, bodyBytes)
	}
}

// startLSWorkoutSession starts a workout session and returns its ID
func startLSWorkoutSession(t *testing.T, ts *testutil.TestServer, userID string) string {
	t.Helper()
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

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope LSWorkoutSessionEnvelope
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// finishLSWorkoutSession finishes a workout session
func finishLSWorkoutSession(t *testing.T, ts *testutil.TestServer, sessionID, userID string) {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, ts.URL("/workouts/"+sessionID+"/finish"), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to finish workout: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
	}
}

func TestLoggedSetHandler_CreateBatch(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create test fixtures: user, lift, cycle, program, enrollment, and workout session
	userID := "ls-test-user"
	createLSTestUser(t, ts, userID)
	liftID := createLSTestLift(t, ts, "Squat", "squat-ls-test")
	cycleID := createLSTestCycle(t, ts, "LS Test Cycle")
	programID := createLSTestProgram(t, ts, "LS Test Program", "ls-test-program", cycleID)
	enrollLSTestUser(t, ts, userID, programID)
	sessionID := startLSWorkoutSession(t, ts, userID)
	prescriptionID := uuid.New().String()

	t.Run("creates logged sets successfully", func(t *testing.T) {
		body := `{
			"sets": [
				{
					"prescriptionId": "` + prescriptionID + `",
					"liftId": "` + liftID + `",
					"setNumber": 1,
					"weight": 225.0,
					"targetReps": 5,
					"repsPerformed": 7,
					"isAmrap": true
				},
				{
					"prescriptionId": "` + prescriptionID + `",
					"liftId": "` + liftID + `",
					"setNumber": 2,
					"weight": 225.0,
					"targetReps": 5,
					"repsPerformed": 5,
					"isAmrap": false
				}
			]
		}`

		resp, err := authPostLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result LoggedSetTestListResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(result.Data) != 2 {
			t.Errorf("Expected 2 logged sets, got %d", len(result.Data))
		}

		// Verify first set (AMRAP)
		if result.Data[0].SetNumber != 1 {
			t.Errorf("Expected set number 1, got %d", result.Data[0].SetNumber)
		}
		if result.Data[0].RepsPerformed != 7 {
			t.Errorf("Expected 7 reps performed, got %d", result.Data[0].RepsPerformed)
		}
		if !result.Data[0].IsAMRAP {
			t.Error("Expected first set to be AMRAP")
		}

		// Verify second set
		if result.Data[1].SetNumber != 2 {
			t.Errorf("Expected set number 2, got %d", result.Data[1].SetNumber)
		}
		if result.Data[1].IsAMRAP {
			t.Error("Expected second set to not be AMRAP")
		}
	})

	t.Run("returns 401 for unauthenticated request", func(t *testing.T) {
		body := `{"sets": []}`
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/sessions/"+sessionID+"/sets"), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for empty sets array", func(t *testing.T) {
		body := `{"sets": []}`
		resp, err := authPostLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		body := `{invalid json}`
		resp, err := authPostLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns validation error for invalid set data", func(t *testing.T) {
		body := `{
			"sets": [
				{
					"prescriptionId": "",
					"liftId": "",
					"setNumber": 0,
					"weight": -10,
					"targetReps": 0,
					"repsPerformed": -1,
					"isAmrap": false
				}
			]
		}`

		resp, err := authPostLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 404 for non-existent session", func(t *testing.T) {
		nonExistentSessionID := uuid.New().String()
		body := `{
			"sets": [
				{
					"prescriptionId": "` + prescriptionID + `",
					"liftId": "` + liftID + `",
					"setNumber": 1,
					"weight": 225.0,
					"targetReps": 5,
					"repsPerformed": 5,
					"isAmrap": false
				}
			]
		}`

		resp, err := authPostLoggedSets(ts.URL("/sessions/"+nonExistentSessionID+"/sets"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for completed session", func(t *testing.T) {
		// First finish the session
		finishLSWorkoutSession(t, ts, sessionID, userID)

		body := `{
			"sets": [
				{
					"prescriptionId": "` + prescriptionID + `",
					"liftId": "` + liftID + `",
					"setNumber": 3,
					"weight": 225.0,
					"targetReps": 5,
					"repsPerformed": 5,
					"isAmrap": false
				}
			]
		}`

		resp, err := authPostLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 400 for completed session, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})
}

func TestLoggedSetHandler_ListBySession(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create test fixtures
	userID := "ls-list-session-user"
	createLSTestUser(t, ts, userID)
	liftID := createLSTestLift(t, ts, "Bench Press", "bench-ls-test")
	cycleID := createLSTestCycle(t, ts, "LS List Session Cycle")
	programID := createLSTestProgram(t, ts, "LS List Session Program", "ls-list-session-program", cycleID)
	enrollLSTestUser(t, ts, userID, programID)
	sessionID := startLSWorkoutSession(t, ts, userID)
	prescriptionID := uuid.New().String()

	// Create some logged sets
	body := `{
		"sets": [
			{
				"prescriptionId": "` + prescriptionID + `",
				"liftId": "` + liftID + `",
				"setNumber": 1,
				"weight": 185.0,
				"targetReps": 5,
				"repsPerformed": 5,
				"isAmrap": false
			}
		]
	}`
	resp, _ := authPostLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), body, userID)
	resp.Body.Close()

	t.Run("lists sets by session", func(t *testing.T) {
		resp, err := authGetLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result LoggedSetTestListResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(result.Data) != 1 {
			t.Errorf("Expected 1 logged set, got %d", len(result.Data))
		}
		if result.Data[0].SessionID != sessionID {
			t.Errorf("Expected session ID %s, got %s", sessionID, result.Data[0].SessionID)
		}
	})

	t.Run("returns empty list for non-existent session", func(t *testing.T) {
		nonExistentSessionID := uuid.New().String()
		resp, err := authGetLoggedSets(ts.URL("/sessions/"+nonExistentSessionID+"/sets"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result LoggedSetTestListResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(result.Data) != 0 {
			t.Errorf("Expected 0 logged sets for non-existent session, got %d", len(result.Data))
		}
	})
}

func TestLoggedSetHandler_ListByUser(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create test fixtures
	userID := "ls-list-user-test"
	createLSTestUser(t, ts, userID)
	liftID := createLSTestLift(t, ts, "Deadlift", "deadlift-ls-test")
	cycleID := createLSTestCycle(t, ts, "LS List User Cycle")
	programID := createLSTestProgram(t, ts, "LS List User Program", "ls-list-user-program", cycleID)
	enrollLSTestUser(t, ts, userID, programID)
	sessionID := startLSWorkoutSession(t, ts, userID)
	prescriptionID := uuid.New().String()

	// Create some logged sets
	body := `{
		"sets": [
			{
				"prescriptionId": "` + prescriptionID + `",
				"liftId": "` + liftID + `",
				"setNumber": 1,
				"weight": 315.0,
				"targetReps": 3,
				"repsPerformed": 5,
				"isAmrap": true
			}
		]
	}`
	resp, _ := authPostLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), body, userID)
	resp.Body.Close()

	t.Run("user can list their own logged sets", func(t *testing.T) {
		resp, err := authGetLoggedSets(ts.URL("/users/"+userID+"/logged-sets"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result LoggedSetTestPaginatedResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(result.Data) == 0 {
			t.Error("Expected at least one logged set")
		}
		if result.Meta == nil {
			t.Error("Expected pagination metadata")
		}
	})

	t.Run("user cannot list other user's logged sets", func(t *testing.T) {
		otherUserID := uuid.New().String()
		resp, err := authGetLoggedSets(ts.URL("/users/"+otherUserID+"/logged-sets"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("admin can list any user's logged sets", func(t *testing.T) {
		resp, err := adminGetLoggedSets(ts.URL("/users/" + userID + "/logged-sets"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("returns 401 for unauthenticated request", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/users/" + userID + "/logged-sets"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})
}
