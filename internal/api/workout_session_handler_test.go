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

// WorkoutSessionTestResponse represents the API response format for a workout session.
type WorkoutSessionTestResponse struct {
	ID                 string     `json:"id"`
	UserProgramStateID string     `json:"userProgramStateId"`
	WeekNumber         int        `json:"weekNumber"`
	DayIndex           int        `json:"dayIndex"`
	Status             string     `json:"status"`
	StartedAt          time.Time  `json:"startedAt"`
	FinishedAt         *time.Time `json:"finishedAt,omitempty"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
}

// WorkoutSessionEnvelope wraps single workout session response with standard envelope.
type WorkoutSessionEnvelope struct {
	Data WorkoutSessionTestResponse `json:"data"`
}

// WorkoutSessionListEnvelope wraps list workout session response with standard envelope.
type WorkoutSessionListEnvelope struct {
	Data []WorkoutSessionTestResponse `json:"data"`
	Meta struct {
		Total   int64 `json:"total"`
		Limit   int   `json:"limit"`
		Offset  int   `json:"offset"`
		HasMore bool  `json:"hasMore"`
	} `json:"meta"`
}

// Helper functions for workout session tests

func userPostWorkoutStart(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func userGetWorkoutSession(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func userPostWorkoutFinish(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func userPostWorkoutAbandon(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func adminGetWorkoutSession(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

func adminPostWorkoutFinish(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// createWorkoutSessionTestCycle creates a test cycle and returns its ID
func createWorkoutSessionTestCycle(t *testing.T, ts *testutil.TestServer, name string) string {
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

	var envelope struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// createWorkoutSessionTestProgram creates a test program and returns its ID
func createWorkoutSessionTestProgram(t *testing.T, ts *testutil.TestServer, name, slug, cycleID string) string {
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

	var envelope struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// enrollUserForWorkoutSession enrolls a user in a program
func enrollUserForWorkoutSession(t *testing.T, ts *testutil.TestServer, userID, programID string) {
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

func TestWorkoutSessionCRUD(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Set up test data
	cycleID := createWorkoutSessionTestCycle(t, ts, "Workout Session Test Cycle")
	programID := createWorkoutSessionTestProgram(t, ts, "Workout Session Test Program", "ws-test-program", cycleID)
	userID := "ws-test-user"

	// Enroll user in program
	enrollUserForWorkoutSession(t, ts, userID, programID)

	var sessionID string

	t.Run("starts a new workout session", func(t *testing.T) {
		resp, err := userPostWorkoutStart(ts.URL("/workouts/start"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope WorkoutSessionEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		session := envelope.Data

		if session.ID == "" {
			t.Error("Expected non-empty ID")
		}
		if session.Status != "IN_PROGRESS" {
			t.Errorf("Expected status IN_PROGRESS, got %s", session.Status)
		}
		if session.WeekNumber != 1 {
			t.Errorf("Expected weekNumber 1, got %d", session.WeekNumber)
		}
		if session.DayIndex != 0 {
			t.Errorf("Expected dayIndex 0, got %d", session.DayIndex)
		}

		sessionID = session.ID
	})

	t.Run("cannot start workout when one already in progress", func(t *testing.T) {
		resp, err := userPostWorkoutStart(ts.URL("/workouts/start"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409, got %d", resp.StatusCode)
		}
	})

	t.Run("gets workout session by ID", func(t *testing.T) {
		resp, err := userGetWorkoutSession(ts.URL("/workouts/"+sessionID), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope WorkoutSessionEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if envelope.Data.ID != sessionID {
			t.Errorf("Expected ID %s, got %s", sessionID, envelope.Data.ID)
		}
	})

	t.Run("gets current workout for user", func(t *testing.T) {
		resp, err := userGetWorkoutSession(ts.URL("/users/"+userID+"/workouts/current"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope WorkoutSessionEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if envelope.Data.ID != sessionID {
			t.Errorf("Expected current session ID %s, got %s", sessionID, envelope.Data.ID)
		}
	})

	t.Run("lists user workouts", func(t *testing.T) {
		resp, err := userGetWorkoutSession(ts.URL("/users/"+userID+"/workouts"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope WorkoutSessionListEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(envelope.Data) != 1 {
			t.Errorf("Expected 1 session, got %d", len(envelope.Data))
		}
		if envelope.Meta.Total != 1 {
			t.Errorf("Expected total 1, got %d", envelope.Meta.Total)
		}
	})

	t.Run("finishes workout session", func(t *testing.T) {
		resp, err := userPostWorkoutFinish(ts.URL("/workouts/"+sessionID+"/finish"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope WorkoutSessionEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if envelope.Data.Status != "COMPLETED" {
			t.Errorf("Expected status COMPLETED, got %s", envelope.Data.Status)
		}
		if envelope.Data.FinishedAt == nil {
			t.Error("Expected finishedAt to be set")
		}
	})

	t.Run("cannot finish already completed session", func(t *testing.T) {
		resp, err := userPostWorkoutFinish(ts.URL("/workouts/"+sessionID+"/finish"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409, got %d", resp.StatusCode)
		}
	})

	t.Run("no current workout after completion", func(t *testing.T) {
		resp, err := userGetWorkoutSession(ts.URL("/users/"+userID+"/workouts/current"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestWorkoutSessionAbandon(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Set up test data
	cycleID := createWorkoutSessionTestCycle(t, ts, "Abandon Test Cycle")
	programID := createWorkoutSessionTestProgram(t, ts, "Abandon Test Program", "abandon-test-program", cycleID)
	userID := "abandon-test-user"

	enrollUserForWorkoutSession(t, ts, userID, programID)

	// Start a workout
	resp, _ := userPostWorkoutStart(ts.URL("/workouts/start"), userID)
	var startEnvelope WorkoutSessionEnvelope
	json.NewDecoder(resp.Body).Decode(&startEnvelope)
	resp.Body.Close()
	sessionID := startEnvelope.Data.ID

	t.Run("abandons workout session", func(t *testing.T) {
		resp, err := userPostWorkoutAbandon(ts.URL("/workouts/"+sessionID+"/abandon"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope WorkoutSessionEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if envelope.Data.Status != "ABANDONED" {
			t.Errorf("Expected status ABANDONED, got %s", envelope.Data.Status)
		}
		if envelope.Data.FinishedAt == nil {
			t.Error("Expected finishedAt to be set")
		}
	})

	t.Run("cannot abandon already abandoned session", func(t *testing.T) {
		resp, err := userPostWorkoutAbandon(ts.URL("/workouts/"+sessionID+"/abandon"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409, got %d", resp.StatusCode)
		}
	})

	t.Run("can start new workout after abandoning", func(t *testing.T) {
		resp, err := userPostWorkoutStart(ts.URL("/workouts/start"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})
}

func TestWorkoutSessionValidation(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	userID := "validation-test-user"

	t.Run("cannot start workout without enrollment", func(t *testing.T) {
		resp, err := userPostWorkoutStart(ts.URL("/workouts/start"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 404 for non-existent session", func(t *testing.T) {
		resp, err := userGetWorkoutSession(ts.URL("/workouts/non-existent-id"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestWorkoutSessionAuthorization(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Set up test data
	cycleID := createWorkoutSessionTestCycle(t, ts, "Auth Test Cycle")
	programID := createWorkoutSessionTestProgram(t, ts, "Auth Test Program", "auth-ws-test-program", cycleID)
	userID := "auth-ws-test-user"
	otherUserID := "other-ws-user"

	enrollUserForWorkoutSession(t, ts, userID, programID)

	// Start a workout
	resp, _ := userPostWorkoutStart(ts.URL("/workouts/start"), userID)
	var startEnvelope WorkoutSessionEnvelope
	json.NewDecoder(resp.Body).Decode(&startEnvelope)
	resp.Body.Close()
	sessionID := startEnvelope.Data.ID

	t.Run("unauthenticated user gets 401 on start", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/workouts/start"), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("unauthenticated user gets 401 on get", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/workouts/"+sessionID), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("user cannot view another user's session", func(t *testing.T) {
		resp, err := userGetWorkoutSession(ts.URL("/workouts/"+sessionID), otherUserID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("user cannot finish another user's session", func(t *testing.T) {
		resp, err := userPostWorkoutFinish(ts.URL("/workouts/"+sessionID+"/finish"), otherUserID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("user cannot view another user's workout list", func(t *testing.T) {
		resp, err := userGetWorkoutSession(ts.URL("/users/"+userID+"/workouts"), otherUserID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("admin can view any user's session", func(t *testing.T) {
		resp, err := adminGetWorkoutSession(ts.URL("/workouts/" + sessionID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("admin can finish any user's session", func(t *testing.T) {
		resp, err := adminPostWorkoutFinish(ts.URL("/workouts/" + sessionID + "/finish"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("admin can view any user's workout list", func(t *testing.T) {
		resp, err := adminGetWorkoutSession(ts.URL("/users/" + userID + "/workouts"))
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

func TestWorkoutSessionStatusFilter(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Set up test data
	cycleID := createWorkoutSessionTestCycle(t, ts, "Filter Test Cycle")
	programID := createWorkoutSessionTestProgram(t, ts, "Filter Test Program", "filter-test-program", cycleID)
	userID := "filter-test-user"

	enrollUserForWorkoutSession(t, ts, userID, programID)

	// Create sessions with different statuses
	// Session 1: Start and complete
	resp, _ := userPostWorkoutStart(ts.URL("/workouts/start"), userID)
	var env WorkoutSessionEnvelope
	json.NewDecoder(resp.Body).Decode(&env)
	resp.Body.Close()
	resp, _ = userPostWorkoutFinish(ts.URL("/workouts/"+env.Data.ID+"/finish"), userID)
	resp.Body.Close()

	// Session 2: Start and abandon
	resp, _ = userPostWorkoutStart(ts.URL("/workouts/start"), userID)
	json.NewDecoder(resp.Body).Decode(&env)
	resp.Body.Close()
	resp, _ = userPostWorkoutAbandon(ts.URL("/workouts/"+env.Data.ID+"/abandon"), userID)
	resp.Body.Close()

	// Session 3: Start (leave in progress)
	resp, _ = userPostWorkoutStart(ts.URL("/workouts/start"), userID)
	resp.Body.Close()

	t.Run("filters by COMPLETED status", func(t *testing.T) {
		resp, err := userGetWorkoutSession(ts.URL("/users/"+userID+"/workouts?status=COMPLETED"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var envelope WorkoutSessionListEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)

		if envelope.Meta.Total != 1 {
			t.Errorf("Expected 1 completed session, got %d", envelope.Meta.Total)
		}
		if len(envelope.Data) > 0 && envelope.Data[0].Status != "COMPLETED" {
			t.Errorf("Expected COMPLETED status, got %s", envelope.Data[0].Status)
		}
	})

	t.Run("filters by ABANDONED status", func(t *testing.T) {
		resp, err := userGetWorkoutSession(ts.URL("/users/"+userID+"/workouts?status=ABANDONED"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var envelope WorkoutSessionListEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)

		if envelope.Meta.Total != 1 {
			t.Errorf("Expected 1 abandoned session, got %d", envelope.Meta.Total)
		}
	})

	t.Run("filters by IN_PROGRESS status", func(t *testing.T) {
		resp, err := userGetWorkoutSession(ts.URL("/users/"+userID+"/workouts?status=IN_PROGRESS"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var envelope WorkoutSessionListEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)

		if envelope.Meta.Total != 1 {
			t.Errorf("Expected 1 in-progress session, got %d", envelope.Meta.Total)
		}
	})

	t.Run("returns all sessions without filter", func(t *testing.T) {
		resp, err := userGetWorkoutSession(ts.URL("/users/"+userID+"/workouts"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var envelope WorkoutSessionListEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)

		if envelope.Meta.Total != 3 {
			t.Errorf("Expected 3 total sessions, got %d", envelope.Meta.Total)
		}
	})

	t.Run("rejects invalid status filter", func(t *testing.T) {
		resp, err := userGetWorkoutSession(ts.URL("/users/"+userID+"/workouts?status=INVALID"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}
