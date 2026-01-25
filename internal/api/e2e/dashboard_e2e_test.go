// Package e2e provides end-to-end tests for complete API workflows.
// This file contains E2E tests for the Dashboard endpoint.
package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// =============================================================================
// DASHBOARD RESPONSE TYPES
// =============================================================================

// DashboardEnrollmentSummary represents the enrollment section of the dashboard.
type DashboardEnrollmentSummary struct {
	Status         string `json:"status"`
	ProgramName    string `json:"programName"`
	CycleIteration int    `json:"cycleIteration"`
	CycleStatus    string `json:"cycleStatus"`
	WeekNumber     int    `json:"weekNumber"`
	WeekStatus     string `json:"weekStatus"`
}

// DashboardNextWorkoutPreview represents the next workout section.
type DashboardNextWorkoutPreview struct {
	DayName       string `json:"dayName"`
	DaySlug       string `json:"daySlug"`
	ExerciseCount int    `json:"exerciseCount"`
	EstimatedSets int    `json:"estimatedSets"`
}

// DashboardSessionSummary represents the current session section.
type DashboardSessionSummary struct {
	SessionID     string `json:"sessionId"`
	DayName       string `json:"dayName"`
	StartedAt     string `json:"startedAt"`
	SetsCompleted int    `json:"setsCompleted"`
	TotalSets     int    `json:"totalSets"`
}

// DashboardWorkoutSummary represents a recent workout.
type DashboardWorkoutSummary struct {
	Date          string `json:"date"`
	DayName       string `json:"dayName"`
	SetsCompleted int    `json:"setsCompleted"`
}

// DashboardMaxSummary represents a current max.
type DashboardMaxSummary struct {
	Lift  string  `json:"lift"`
	Value float64 `json:"value"`
	Type  string  `json:"type"`
}

// DashboardResponseData represents the full dashboard response.
type DashboardResponseData struct {
	Enrollment     *DashboardEnrollmentSummary   `json:"enrollment"`
	NextWorkout    *DashboardNextWorkoutPreview  `json:"nextWorkout"`
	CurrentSession *DashboardSessionSummary      `json:"currentSession"`
	RecentWorkouts []DashboardWorkoutSummary     `json:"recentWorkouts"`
	CurrentMaxes   []DashboardMaxSummary         `json:"currentMaxes"`
}

// DashboardEnvelopeData wraps dashboard response in data envelope.
type DashboardEnvelopeData struct {
	Data DashboardResponseData `json:"data"`
}

// =============================================================================
// DASHBOARD HELPER FUNCTIONS
// =============================================================================

// bearerGetDashboard performs a GET request with Bearer token authentication.
func bearerGetDashboard(url string, token string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return http.DefaultClient.Do(req)
}

// userIDGetDashboard performs a GET request with X-User-ID header (test mode auth).
func userIDGetDashboard(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

// adminGetDashboard performs an admin-authenticated GET request.
func adminGetDashboard(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// =============================================================================
// E2E TESTS: GET /users/{id}/dashboard
// =============================================================================

func TestDashboardE2E_GetDashboard_Unenrolled(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]

	// Create test user (not enrolled in any program)
	email := fmt.Sprintf("dashboard-unenrolled-%s@example.com", testID)
	user := registerUser(t, ts, email, "password123", "Unenrolled User")
	loginResult := loginUser(t, ts, email, "password123")
	token := loginResult.Token

	t.Run("returns 200 with empty dashboard for unenrolled user", func(t *testing.T) {
		resp, err := bearerGetDashboard(ts.URL("/users/"+user.ID+"/dashboard"), token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result DashboardEnvelopeData
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Verify empty sections
		if result.Data.Enrollment != nil {
			t.Errorf("Expected enrollment to be null for unenrolled user, got %+v", result.Data.Enrollment)
		}
		if result.Data.NextWorkout != nil {
			t.Errorf("Expected nextWorkout to be null for unenrolled user, got %+v", result.Data.NextWorkout)
		}
		if result.Data.CurrentSession != nil {
			t.Errorf("Expected currentSession to be null for unenrolled user, got %+v", result.Data.CurrentSession)
		}
		if result.Data.RecentWorkouts == nil {
			t.Error("Expected recentWorkouts to be empty array, got nil")
		}
		if len(result.Data.RecentWorkouts) != 0 {
			t.Errorf("Expected recentWorkouts to be empty, got %d items", len(result.Data.RecentWorkouts))
		}
		if result.Data.CurrentMaxes == nil {
			t.Error("Expected currentMaxes to be empty array, got nil")
		}
		if len(result.Data.CurrentMaxes) != 0 {
			t.Errorf("Expected currentMaxes to be empty, got %d items", len(result.Data.CurrentMaxes))
		}
	})
}

func TestDashboardE2E_Authorization(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]

	// Create two test users
	userAEmail := fmt.Sprintf("dashboard-auth-a-%s@example.com", testID)
	userA := registerUser(t, ts, userAEmail, "password123", "User A")
	loginResultA := loginUser(t, ts, userAEmail, "password123")
	tokenA := loginResultA.Token

	userBEmail := fmt.Sprintf("dashboard-auth-b-%s@example.com", testID)
	userB := registerUser(t, ts, userBEmail, "password123", "User B")
	loginResultB := loginUser(t, ts, userBEmail, "password123")
	tokenB := loginResultB.Token

	t.Run("returns 401 without authentication", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/users/"+userA.ID+"/dashboard"), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 401, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("returns 403 for non-owner (user B accessing user A dashboard)", func(t *testing.T) {
		resp, err := bearerGetDashboard(ts.URL("/users/"+userA.ID+"/dashboard"), tokenB)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("returns 403 for admin accessing other user dashboard (owner-only)", func(t *testing.T) {
		resp, err := adminGetDashboard(ts.URL("/users/" + userA.ID + "/dashboard"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("returns 200 for owner using Bearer token", func(t *testing.T) {
		resp, err := bearerGetDashboard(ts.URL("/users/"+userA.ID+"/dashboard"), tokenA)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("returns 200 for owner using X-User-ID", func(t *testing.T) {
		resp, err := userIDGetDashboard(ts.URL("/users/"+userA.ID+"/dashboard"), userA.ID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	_ = userB
	_ = tokenB
}

func TestDashboardE2E_ResponseStructure(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]

	// Create test user
	email := fmt.Sprintf("dashboard-structure-%s@example.com", testID)
	user := registerUser(t, ts, email, "password123", "Structure Test User")
	loginResult := loginUser(t, ts, email, "password123")
	token := loginResult.Token

	t.Run("response has all expected fields with correct structure", func(t *testing.T) {
		resp, err := bearerGetDashboard(ts.URL("/users/"+user.ID+"/dashboard"), token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		bodyBytes, _ := io.ReadAll(resp.Body)

		// Parse as raw JSON to check structure
		var rawResult map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &rawResult); err != nil {
			t.Fatalf("Failed to parse response as JSON: %v", err)
		}

		// Check top-level data wrapper
		data, ok := rawResult["data"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected 'data' field in response")
		}

		// Check all required fields are present (even if null/empty)
		requiredFields := []string{"enrollment", "nextWorkout", "currentSession", "recentWorkouts", "currentMaxes"}
		for _, field := range requiredFields {
			if _, exists := data[field]; !exists {
				t.Errorf("Expected field '%s' to be present in response", field)
			}
		}

		// Verify recentWorkouts is an array (not null)
		if recentWorkouts, ok := data["recentWorkouts"].([]interface{}); !ok {
			t.Error("Expected recentWorkouts to be an array")
		} else if recentWorkouts == nil {
			t.Error("Expected recentWorkouts to be empty array, not null")
		}

		// Verify currentMaxes is an array (not null)
		if currentMaxes, ok := data["currentMaxes"].([]interface{}); !ok {
			t.Error("Expected currentMaxes to be an array")
		} else if currentMaxes == nil {
			t.Error("Expected currentMaxes to be empty array, not null")
		}
	})
}

func TestDashboardE2E_WithLiftMaxes(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]

	// Create test user
	email := fmt.Sprintf("dashboard-maxes-%s@example.com", testID)
	user := registerUser(t, ts, email, "password123", "Maxes Test User")
	loginResult := loginUser(t, ts, email, "password123")
	token := loginResult.Token

	// Create lifts (as admin) with unique slugs for this test
	squatID := dashboardCreateLift(t, ts, "Squat", "dash-squat-"+testID, false)
	benchID := dashboardCreateLift(t, ts, "Bench Press", "dash-bench-"+testID, false)
	deadliftID := dashboardCreateLift(t, ts, "Deadlift", "dash-deadlift-"+testID, false)

	// Create lift maxes for the user
	dashboardCreateLiftMax(t, ts, user.ID, squatID, "ONE_RM", 315.0)
	dashboardCreateLiftMax(t, ts, user.ID, benchID, "ONE_RM", 225.0)
	dashboardCreateLiftMax(t, ts, user.ID, deadliftID, "TRAINING_MAX", 365.0)

	t.Run("returns current maxes for user", func(t *testing.T) {
		resp, err := bearerGetDashboard(ts.URL("/users/"+user.ID+"/dashboard"), token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result DashboardEnvelopeData
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Should have 3 maxes
		if len(result.Data.CurrentMaxes) != 3 {
			t.Errorf("Expected 3 current maxes, got %d", len(result.Data.CurrentMaxes))
		}

		// Check that maxes are sorted alphabetically by lift name
		// Bench Press < Deadlift < Squat
		if len(result.Data.CurrentMaxes) >= 3 {
			if result.Data.CurrentMaxes[0].Lift != "Bench Press" {
				t.Errorf("Expected first max to be Bench Press, got %s", result.Data.CurrentMaxes[0].Lift)
			}
			if result.Data.CurrentMaxes[1].Lift != "Deadlift" {
				t.Errorf("Expected second max to be Deadlift, got %s", result.Data.CurrentMaxes[1].Lift)
			}
			if result.Data.CurrentMaxes[2].Lift != "Squat" {
				t.Errorf("Expected third max to be Squat, got %s", result.Data.CurrentMaxes[2].Lift)
			}
		}

		// Check values
		for _, max := range result.Data.CurrentMaxes {
			switch max.Lift {
			case "Squat":
				if max.Value != 315.0 {
					t.Errorf("Expected Squat max 315.0, got %.2f", max.Value)
				}
				if max.Type != "ONE_RM" {
					t.Errorf("Expected Squat type ONE_RM, got %s", max.Type)
				}
			case "Bench Press":
				if max.Value != 225.0 {
					t.Errorf("Expected Bench Press max 225.0, got %.2f", max.Value)
				}
			case "Deadlift":
				if max.Value != 365.0 {
					t.Errorf("Expected Deadlift max 365.0, got %.2f", max.Value)
				}
				if max.Type != "TRAINING_MAX" {
					t.Errorf("Expected Deadlift type TRAINING_MAX, got %s", max.Type)
				}
			}
		}
	})
}

// =============================================================================
// HELPER FUNCTIONS FOR CREATING TEST DATA
// =============================================================================

// dashboardCreateLift creates a lift via admin API and returns the ID.
func dashboardCreateLift(t *testing.T, ts *testutil.TestServer, name, slug string, isAccessory bool) string {
	t.Helper()

	body := fmt.Sprintf(`{"name": "%s", "slug": "%s", "isAccessory": %v}`, name, slug, isAccessory)

	req, err := http.NewRequest(http.MethodPost, ts.URL("/lifts"), bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to create lift: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create lift: %d - %s", resp.StatusCode, bodyBytes)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	data := result["data"].(map[string]interface{})
	return data["id"].(string)
}

// dashboardCreateLiftMax creates a lift max for a user and returns the ID.
func dashboardCreateLiftMax(t *testing.T, ts *testutil.TestServer, userID, liftID, maxType string, value float64) {
	t.Helper()

	body := fmt.Sprintf(`{"liftId": "%s", "type": "%s", "value": %f}`, liftID, maxType, value)

	req, err := http.NewRequest(http.MethodPost, ts.URL("/users/"+userID+"/lift-maxes"), bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-User-ID", userID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to create lift max: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create lift max: %d - %s", resp.StatusCode, bodyBytes)
	}
}

// Needed to mark strings import as used (used in test assertions)
var _ = strings.Contains
