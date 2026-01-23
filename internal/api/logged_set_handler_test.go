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

func TestLoggedSetHandler_CreateBatch(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a test lift
	liftID := createLSTestLift(t, ts, "Squat", "squat-ls-test")
	sessionID := uuid.New().String()
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

		resp, err := authPostLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), body, testutil.TestUserID)
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
		resp, err := authPostLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), body, testutil.TestUserID)
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
		resp, err := authPostLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), body, testutil.TestUserID)
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

		resp, err := authPostLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), body, testutil.TestUserID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestLoggedSetHandler_ListBySession(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a test lift
	liftID := createLSTestLift(t, ts, "Bench Press", "bench-ls-test")
	sessionID := uuid.New().String()
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
	resp, _ := authPostLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), body, testutil.TestUserID)
	resp.Body.Close()

	t.Run("lists sets by session", func(t *testing.T) {
		resp, err := authGetLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), testutil.TestUserID)
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
		resp, err := authGetLoggedSets(ts.URL("/sessions/"+nonExistentSessionID+"/sets"), testutil.TestUserID)
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

	// Create a test lift
	liftID := createLSTestLift(t, ts, "Deadlift", "deadlift-ls-test")
	sessionID := uuid.New().String()
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
	resp, _ := authPostLoggedSets(ts.URL("/sessions/"+sessionID+"/sets"), body, testutil.TestUserID)
	resp.Body.Close()

	t.Run("user can list their own logged sets", func(t *testing.T) {
		resp, err := authGetLoggedSets(ts.URL("/users/"+testutil.TestUserID+"/logged-sets"), testutil.TestUserID)
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
		resp, err := authGetLoggedSets(ts.URL("/users/"+otherUserID+"/logged-sets"), testutil.TestUserID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("admin can list any user's logged sets", func(t *testing.T) {
		resp, err := adminGetLoggedSets(ts.URL("/users/" + testutil.TestUserID + "/logged-sets"))
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
		resp, err := http.Get(ts.URL("/users/" + testutil.TestUserID + "/logged-sets"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})
}
