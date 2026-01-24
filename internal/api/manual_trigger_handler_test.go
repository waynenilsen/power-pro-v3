package api_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// ManualTriggerRequest matches the API request body.
type ManualTriggerRequest struct {
	ProgressionID string `json:"progressionId"`
	LiftID        string `json:"liftId,omitempty"`
	Force         bool   `json:"force"`
}

// ManualTriggerResultResponse represents a single progression result in the API response.
type ManualTriggerResultResponse struct {
	ProgressionID string                         `json:"progressionId"`
	LiftID        string                         `json:"liftId"`
	Applied       bool                           `json:"applied"`
	Skipped       bool                           `json:"skipped,omitempty"`
	SkipReason    string                         `json:"skipReason,omitempty"`
	Result        *ManualTriggerResultDetail     `json:"result,omitempty"`
	Error         string                         `json:"error,omitempty"`
}

// ManualTriggerResultDetail contains the details of an applied progression.
type ManualTriggerResultDetail struct {
	PreviousValue float64   `json:"previousValue"`
	NewValue      float64   `json:"newValue"`
	Delta         float64   `json:"delta"`
	MaxType       string    `json:"maxType"`
	AppliedAt     time.Time `json:"appliedAt"`
}

// ManualTriggerResponse represents the response for manual progression trigger.
type ManualTriggerResponse struct {
	Results      []ManualTriggerResultResponse `json:"results"`
	TotalApplied int                           `json:"totalApplied"`
	TotalSkipped int                           `json:"totalSkipped"`
	TotalErrors  int                           `json:"totalErrors"`
}

// ProgressionResponseEnvelope wraps the API response with standard envelope.
type ProgressionResponseEnvelope struct {
	Data ProgressionResponse `json:"data"`
}

// CycleResponseData represents cycle data in the API response.
type CycleResponseData struct {
	ID string `json:"id"`
}

// CycleResponseEnvelope wraps cycle response with standard envelope.
type CycleResponseEnvelope struct {
	Data CycleResponseData `json:"data"`
}

// ProgramResponseData represents program data in the API response.
type ProgramResponseData struct {
	ID string `json:"id"`
}

// ProgramResponseEnvelope wraps program response with standard envelope.
type ProgramResponseEnvelope struct {
	Data ProgramResponseData `json:"data"`
}

// ManualTriggerResponseEnvelope wraps manual trigger response with standard envelope.
type ManualTriggerResponseEnvelope struct {
	Data ManualTriggerResponse `json:"data"`
}

// Helper functions for manual trigger tests

func authPostTrigger(url string, body interface{}, userID string) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func adminPostTrigger(url string, body interface{}) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// createMTTestLift creates a test lift and returns its ID
func createMTTestLift(t *testing.T, ts *testutil.TestServer, name, slug string) string {
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

// createMTTestProgression creates a test progression and returns its ID
func createMTTestProgression(t *testing.T, ts *testutil.TestServer, name string, progType string, triggerType string) string {
	t.Helper()
	var body string
	if progType == "LINEAR_PROGRESSION" {
		body = `{"name": "` + name + `", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "` + triggerType + `"}}`
	} else {
		body = `{"name": "` + name + `", "type": "CYCLE_PROGRESSION", "parameters": {"increment": 10.0, "maxType": "TRAINING_MAX"}}`
	}
	resp, err := adminPost(ts.URL("/progressions"), body)
	if err != nil {
		t.Fatalf("Failed to create test progression: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create test progression (status %d): %s", resp.StatusCode, bodyBytes)
	}

	var envelope ProgressionResponseEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		t.Fatalf("Failed to decode progression response: %v", err)
	}
	return envelope.Data.ID
}

// setupManualTriggerTestData creates all necessary test data for manual trigger tests
func setupManualTriggerTestData(t *testing.T, ts *testutil.TestServer) (userID, liftID, progressionID, programID string) {
	t.Helper()

	// Use the test user ID
	userID = testutil.TestUserID

	// Create a cycle
	cycleResp, err := adminPost(ts.URL("/cycles"), `{"name": "MT Test Cycle", "lengthWeeks": 4}`)
	if err != nil {
		t.Fatalf("Failed to create cycle: %v", err)
	}
	defer cycleResp.Body.Close()
	var cycleEnvelope CycleResponseEnvelope
	if err := json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope); err != nil {
		t.Fatalf("Failed to decode cycle response: %v", err)
	}

	// Create a program
	progResp, err := adminPost(ts.URL("/programs"), `{"name": "MT Test Program", "slug": "mt-test-program", "cycleId": "`+cycleEnvelope.Data.ID+`"}`)
	if err != nil {
		t.Fatalf("Failed to create program: %v", err)
	}
	defer progResp.Body.Close()
	var programEnvelope ProgramResponseEnvelope
	if err := json.NewDecoder(progResp.Body).Decode(&programEnvelope); err != nil {
		t.Fatalf("Failed to decode program response: %v", err)
	}
	programID = programEnvelope.Data.ID

	// Create a lift
	liftID = createMTTestLift(t, ts, "MT Test Squat", "mt-test-squat-"+uuid.New().String()[:8])

	// Create a progression
	progressionID = createMTTestProgression(t, ts, "MT Test Linear", "LINEAR_PROGRESSION", "AFTER_SESSION")

	// Link progression to program
	ppBody := `{"progressionId": "` + progressionID + `", "liftId": "` + liftID + `", "priority": 1, "enabled": true}`
	ppResp, err := adminPost(ts.URL("/programs/"+programID+"/progressions"), ppBody)
	if err != nil {
		t.Fatalf("Failed to create program progression: %v", err)
	}
	ppResp.Body.Close()

	// Enroll user in program
	enrollResp, err := authPostUser(ts.URL("/users/"+userID+"/program"), `{"programId": "`+programID+`"}`, userID)
	if err != nil {
		t.Fatalf("Failed to enroll user: %v", err)
	}
	enrollResp.Body.Close()

	// Create initial lift max
	maxBody := `{"liftId": "` + liftID + `", "type": "TRAINING_MAX", "value": 300, "effectiveDate": "2024-01-01T00:00:00Z"}`
	maxResp, err := authPostUser(ts.URL("/users/"+userID+"/lift-maxes"), maxBody, userID)
	if err != nil {
		t.Fatalf("Failed to create lift max: %v", err)
	}
	maxResp.Body.Close()

	return userID, liftID, progressionID, programID
}

func TestManualTriggerAuthorization(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("returns 401 for unauthenticated request", func(t *testing.T) {
		body := ManualTriggerRequest{ProgressionID: "some-id"}
		bodyBytes, _ := json.Marshal(body)

		resp, err := http.Post(ts.URL("/users/"+testutil.TestUserID+"/progressions/trigger"), "application/json", bytes.NewReader(bodyBytes))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 403 when user tries to trigger for another user", func(t *testing.T) {
		otherUserID := "other-user-id"
		body := ManualTriggerRequest{ProgressionID: "some-id"}

		resp, err := authPostTrigger(ts.URL("/users/"+otherUserID+"/progressions/trigger"), body, testutil.TestUserID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("admin can trigger for any user", func(t *testing.T) {
		// This will fail with 404 for progression not found, but not 403
		body := ManualTriggerRequest{ProgressionID: "some-id"}

		resp, err := adminPostTrigger(ts.URL("/users/some-other-user/progressions/trigger"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Should not be 403 (admin allowed), but will be 404 or 400
		if resp.StatusCode == http.StatusForbidden {
			t.Error("Expected admin to be allowed, got 403")
		}
	})
}

func TestManualTriggerValidation(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("returns 400 for missing progressionId", func(t *testing.T) {
		body := ManualTriggerRequest{} // Empty progressionId

		resp, err := authPostTrigger(ts.URL("/users/"+testutil.TestUserID+"/progressions/trigger"), body, testutil.TestUserID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 400, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/users/"+testutil.TestUserID+"/progressions/trigger"), bytes.NewReader([]byte("not json")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", testutil.TestUserID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestManualTriggerErrorHandling(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("returns 404 for non-existent progressionId", func(t *testing.T) {
		body := ManualTriggerRequest{ProgressionID: "non-existent-progression"}

		resp, err := authPostTrigger(ts.URL("/users/"+testutil.TestUserID+"/progressions/trigger"), body, testutil.TestUserID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 404, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("returns 404 for non-existent liftId", func(t *testing.T) {
		// Set up test data (user must be enrolled first)
		userID, _, progressionID, _ := setupManualTriggerTestData(t, ts)

		body := ManualTriggerRequest{
			ProgressionID: progressionID,
			LiftID:        "non-existent-lift",
		}

		resp, err := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 404, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("returns 400 when user not enrolled", func(t *testing.T) {
		// Create a new user that is not enrolled
		newUserID := uuid.New().String()
		progressionID := createMTTestProgression(t, ts, "Unenrolled Test Prog", "LINEAR_PROGRESSION", "AFTER_SESSION")

		// Create the user in the database
		// Note: This may fail if the user doesn't exist - the API will return appropriate error
		body := ManualTriggerRequest{ProgressionID: progressionID}

		resp, err := authPostTrigger(ts.URL("/users/"+newUserID+"/progressions/trigger"), body, newUserID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Should return 400 for user not enrolled
		if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusNotFound {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 400 or 404, got %d: %s", resp.StatusCode, respBody)
		}
	})
}

func TestManualTriggerResponseFormat(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Set up complete test data
	userID, liftID, progressionID, _ := setupManualTriggerTestData(t, ts)

	t.Run("response has correct JSON structure", func(t *testing.T) {
		body := ManualTriggerRequest{
			ProgressionID: progressionID,
			LiftID:        liftID,
			Force:         false,
		}

		resp, err := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, respBody)
		}

		var envelope ManualTriggerResponseEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		result := envelope.Data

		// Verify response structure
		if result.Results == nil {
			t.Error("Expected results array to be present")
		}

		// Verify counts
		total := result.TotalApplied + result.TotalSkipped + result.TotalErrors
		if total != len(result.Results) {
			t.Errorf("Expected totals to match results length, got applied=%d + skipped=%d + errors=%d = %d, but len=%d",
				result.TotalApplied, result.TotalSkipped, result.TotalErrors, total, len(result.Results))
		}
	})

	t.Run("successful trigger returns applied result", func(t *testing.T) {
		body := ManualTriggerRequest{
			ProgressionID: progressionID,
			LiftID:        liftID,
			Force:         true, // Use force to ensure it applies
		}

		resp, err := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, respBody)
		}

		var envelope ManualTriggerResponseEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		result := envelope.Data

		if result.TotalApplied != 1 {
			t.Errorf("Expected TotalApplied=1, got %d", result.TotalApplied)
		}

		if len(result.Results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(result.Results))
		}

		r := result.Results[0]
		if !r.Applied {
			t.Errorf("Expected Applied=true, got false (skipReason: %s, error: %s)", r.SkipReason, r.Error)
		}

		if r.LiftID != liftID {
			t.Errorf("Expected liftId=%s, got %s", liftID, r.LiftID)
		}

		if r.Result == nil {
			t.Error("Expected result detail to be present")
		} else {
			if r.Result.Delta != 5.0 {
				t.Errorf("Expected delta=5.0, got %f", r.Result.Delta)
			}
			if r.Result.NewValue != r.Result.PreviousValue+5.0 {
				t.Errorf("Expected newValue=%f, got %f", r.Result.PreviousValue+5.0, r.Result.NewValue)
			}
		}
	})
}

func TestManualTriggerIdempotency(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Set up test data
	userID, liftID, progressionID, _ := setupManualTriggerTestData(t, ts)

	t.Run("force=true allows repeated applications", func(t *testing.T) {
		body := ManualTriggerRequest{
			ProgressionID: progressionID,
			LiftID:        liftID,
			Force:         true,
		}

		// First application
		resp1, err := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make first request: %v", err)
		}
		defer resp1.Body.Close()

		var envelope1 ManualTriggerResponseEnvelope
		json.NewDecoder(resp1.Body).Decode(&envelope1)
		result1 := envelope1.Data

		if result1.TotalApplied != 1 {
			t.Fatalf("Expected first application to succeed, got TotalApplied=%d", result1.TotalApplied)
		}

		firstValue := result1.Results[0].Result.NewValue

		// Second application with force=true should also succeed
		resp2, err := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make second request: %v", err)
		}
		defer resp2.Body.Close()

		var envelope2 ManualTriggerResponseEnvelope
		json.NewDecoder(resp2.Body).Decode(&envelope2)
		result2 := envelope2.Data

		if result2.TotalApplied != 1 {
			t.Errorf("Expected second force application to succeed, got TotalApplied=%d (skipped=%d, errors=%d)",
				result2.TotalApplied, result2.TotalSkipped, result2.TotalErrors)
		}

		if result2.Results[0].Result != nil {
			secondValue := result2.Results[0].Result.NewValue
			if secondValue != firstValue+5.0 {
				t.Errorf("Expected second value to be %f (first + increment), got %f", firstValue+5.0, secondValue)
			}
		}
	})
}

func TestManualTriggerMultipleLifts(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Set up test data with multiple lifts
	userID := testutil.TestUserID

	// Create a cycle
	cycleResp, err := adminPost(ts.URL("/cycles"), `{"name": "Multi Lift Cycle", "lengthWeeks": 4}`)
	if err != nil {
		t.Fatalf("Failed to create cycle: %v", err)
	}
	var cycleEnvelope CycleResponseEnvelope
	if err := json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope); err != nil {
		t.Fatalf("Failed to decode cycle response: %v", err)
	}
	cycleResp.Body.Close()

	// Create a program
	progResp, err := adminPost(ts.URL("/programs"), `{"name": "Multi Lift Program", "slug": "multi-lift-program-`+uuid.New().String()[:8]+`", "cycleId": "`+cycleEnvelope.Data.ID+`"}`)
	if err != nil {
		t.Fatalf("Failed to create program: %v", err)
	}
	var programEnvelope ProgramResponseEnvelope
	if err := json.NewDecoder(progResp.Body).Decode(&programEnvelope); err != nil {
		t.Fatalf("Failed to decode program response: %v", err)
	}
	progResp.Body.Close()

	// Create multiple lifts
	lift1ID := createMTTestLift(t, ts, "Multi Squat", "multi-squat-"+uuid.New().String()[:8])
	lift2ID := createMTTestLift(t, ts, "Multi Bench", "multi-bench-"+uuid.New().String()[:8])

	// Create a progression
	progressionID := createMTTestProgression(t, ts, "Multi Lift Prog", "LINEAR_PROGRESSION", "AFTER_SESSION")

	// Link progression to program for both lifts
	ppBody1 := `{"progressionId": "` + progressionID + `", "liftId": "` + lift1ID + `", "priority": 1, "enabled": true}`
	pp1Resp, _ := adminPost(ts.URL("/programs/"+programEnvelope.Data.ID+"/progressions"), ppBody1)
	pp1Resp.Body.Close()

	ppBody2 := `{"progressionId": "` + progressionID + `", "liftId": "` + lift2ID + `", "priority": 2, "enabled": true}`
	pp2Resp, _ := adminPost(ts.URL("/programs/"+programEnvelope.Data.ID+"/progressions"), ppBody2)
	pp2Resp.Body.Close()

	// Enroll user in program (need to unenroll first if already enrolled)
	unenrollResp, _ := authDelete(ts.URL("/users/"+userID+"/program"), userID)
	unenrollResp.Body.Close()

	enrollResp, _ := authPostUser(ts.URL("/users/"+userID+"/program"), `{"programId": "`+programEnvelope.Data.ID+`"}`, userID)
	enrollResp.Body.Close()

	// Create initial lift maxes
	max1Body := `{"liftId": "` + lift1ID + `", "type": "TRAINING_MAX", "value": 300, "effectiveDate": "2024-01-01T00:00:00Z"}`
	max1Resp, _ := authPostUser(ts.URL("/users/"+userID+"/lift-maxes"), max1Body, userID)
	max1Resp.Body.Close()

	max2Body := `{"liftId": "` + lift2ID + `", "type": "TRAINING_MAX", "value": 200, "effectiveDate": "2024-01-01T00:00:00Z"}`
	max2Resp, _ := authPostUser(ts.URL("/users/"+userID+"/lift-maxes"), max2Body, userID)
	max2Resp.Body.Close()

	t.Run("applies to all configured lifts when liftId is empty", func(t *testing.T) {
		body := ManualTriggerRequest{
			ProgressionID: progressionID,
			LiftID:        "", // Empty - should apply to all
			Force:         true,
		}

		resp, err := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, respBody)
		}

		var envelope ManualTriggerResponseEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		result := envelope.Data

		if result.TotalApplied != 2 {
			t.Errorf("Expected TotalApplied=2 (both lifts), got %d (skipped=%d, errors=%d)",
				result.TotalApplied, result.TotalSkipped, result.TotalErrors)
		}

		if len(result.Results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(result.Results))
		}

		// Verify both lifts are in results
		foundLift1 := false
		foundLift2 := false
		for _, r := range result.Results {
			if r.LiftID == lift1ID {
				foundLift1 = true
			}
			if r.LiftID == lift2ID {
				foundLift2 = true
			}
		}
		if !foundLift1 {
			t.Error("Expected lift1 to be in results")
		}
		if !foundLift2 {
			t.Error("Expected lift2 to be in results")
		}
	})

	t.Run("applies to specific lift when liftId is provided", func(t *testing.T) {
		body := ManualTriggerRequest{
			ProgressionID: progressionID,
			LiftID:        lift1ID, // Specific lift
			Force:         true,
		}

		resp, err := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var envelope ManualTriggerResponseEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		result := envelope.Data

		if result.TotalApplied != 1 {
			t.Errorf("Expected TotalApplied=1, got %d", result.TotalApplied)
		}

		if len(result.Results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(result.Results))
		}

		if result.Results[0].LiftID != lift1ID {
			t.Errorf("Expected liftId=%s, got %s", lift1ID, result.Results[0].LiftID)
		}
	})
}

// Helper function
func authDelete(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

// TestManualTriggerAuditLogging tests that force=true creates proper audit logs
func TestManualTriggerAuditLogging(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Set up test data
	userID, liftID, progressionID, _ := setupManualTriggerTestData(t, ts)

	t.Run("force=true trigger creates log with manual marker", func(t *testing.T) {
		body := ManualTriggerRequest{
			ProgressionID: progressionID,
			LiftID:        liftID,
			Force:         true,
		}

		resp, err := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		resp.Body.Close()

		// Verify by checking progression history
		histResp, err := authGetHistory(ts.URL("/users/"+userID+"/progression-history?limit=1"), userID)
		if err != nil {
			t.Fatalf("Failed to get history: %v", err)
		}
		defer histResp.Body.Close()

		var history ProgressionHistoryTestListResponse
		json.NewDecoder(histResp.Body).Decode(&history)

		if len(history.Data) > 0 {
			// Check the trigger context contains manual and force markers
			contextStr := string(history.Data[0].TriggerContext)
			if !containsSubstring(contextStr, `"manual"`) {
				t.Errorf("Expected trigger context to contain 'manual', got: %s", contextStr)
			}
		}
	})
}

// Helper to check substring (simple implementation)
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestManualTriggerWithTestServer tests the complete flow using TestServer's DB access
func TestManualTriggerWithTestServer(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Get direct DB access
	sqlDB := ts.DB()
	if sqlDB == nil {
		t.Skip("TestServer does not expose DB connection")
	}

	queries := db.New(sqlDB)

	// Create test data directly in DB
	userID := testutil.TestUserID
	now := time.Now().Format(time.RFC3339)
	pastDate := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)

	// Create cycle
	cycleID := uuid.New().String()
	err = queries.CreateCycle(t.Context(), db.CreateCycleParams{
		ID:          cycleID,
		Name:        "Direct DB Cycle",
		LengthWeeks: 4,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		t.Fatalf("Failed to create cycle: %v", err)
	}

	// Create program
	programID := uuid.New().String()
	err = queries.CreateProgram(t.Context(), db.CreateProgramParams{
		ID:        programID,
		Name:      "Direct DB Program",
		Slug:      "direct-db-program-" + uuid.New().String()[:8],
		CycleID:   cycleID,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		t.Fatalf("Failed to create program: %v", err)
	}

	// Use seeded lifts
	liftID := "00000000-0000-0000-0000-000000000001" // Seeded squat

	// Create progression
	progressionID := uuid.New().String()
	err = queries.CreateProgression(t.Context(), db.CreateProgressionParams{
		ID:   progressionID,
		Name: "Direct DB Prog",
		Type: string(progression.TypeLinear),
		Parameters: `{
			"id": "` + progressionID + `",
			"name": "Direct DB Prog",
			"increment": 5.0,
			"maxType": "TRAINING_MAX",
			"triggerType": "AFTER_SESSION"
		}`,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		t.Fatalf("Failed to create progression: %v", err)
	}

	// Create program progression
	ppID := uuid.New().String()
	err = queries.CreateProgramProgression(t.Context(), db.CreateProgramProgressionParams{
		ID:            ppID,
		ProgramID:     programID,
		ProgressionID: progressionID,
		LiftID:        sql.NullString{String: liftID, Valid: true},
		Priority:      1,
		Enabled:       1,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		t.Fatalf("Failed to create program progression: %v", err)
	}

	// Unenroll user first (if enrolled)
	queries.DeleteUserProgramStateByUserID(t.Context(), userID)

	// Enroll user
	enrollmentID := uuid.New().String()
	err = queries.CreateUserProgramState(t.Context(), db.CreateUserProgramStateParams{
		ID:                    enrollmentID,
		UserID:                userID,
		ProgramID:             programID,
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		EnrolledAt:            now,
		UpdatedAt:             now,
	})
	if err != nil {
		t.Fatalf("Failed to enroll user: %v", err)
	}

	// Create lift max
	maxID := uuid.New().String()
	err = queries.CreateLiftMax(t.Context(), db.CreateLiftMaxParams{
		ID:            maxID,
		UserID:        userID,
		LiftID:        liftID,
		Type:          "TRAINING_MAX",
		Value:         300,
		EffectiveDate: pastDate,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		t.Fatalf("Failed to create lift max: %v", err)
	}

	t.Run("trigger via API creates expected database entries", func(t *testing.T) {
		body := ManualTriggerRequest{
			ProgressionID: progressionID,
			LiftID:        liftID,
			Force:         true,
		}

		resp, err := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var envelope ManualTriggerResponseEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		result := envelope.Data

		if result.TotalApplied != 1 {
			respBody, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected TotalApplied=1, got %d. Response: %s", result.TotalApplied, respBody)
		}

		// Verify LiftMax was created
		currentMax, err := queries.GetCurrentMax(t.Context(), db.GetCurrentMaxParams{
			UserID: userID,
			LiftID: liftID,
			Type:   "TRAINING_MAX",
		})
		if err != nil {
			t.Fatalf("Failed to get current max: %v", err)
		}

		if currentMax.Value != 305 {
			t.Errorf("Expected new max value 305, got %f", currentMax.Value)
		}

		// Verify ProgressionLog was created with manual marker
		logs, err := queries.ListProgressionLogsByUser(t.Context(), db.ListProgressionLogsByUserParams{
			UserID: userID,
			Limit:  1,
			Offset: 0,
		})
		if err != nil {
			t.Fatalf("Failed to get logs: %v", err)
		}

		if len(logs) == 0 {
			t.Fatal("Expected at least one log entry")
		}

		// Verify trigger context contains manual marker
		triggerCtx := logs[0].TriggerContext.String
		if !containsSubstring(triggerCtx, `"manual":true`) && !containsSubstring(triggerCtx, `"manual": true`) {
			t.Errorf("Expected trigger_context to contain 'manual:true', got: %s", triggerCtx)
		}
	})
}
