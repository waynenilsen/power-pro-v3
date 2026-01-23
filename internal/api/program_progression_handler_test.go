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

// ProgramProgressionResponse matches the API response format for a program progression config.
type ProgramProgressionTestResponse struct {
	ID                string                  `json:"id"`
	ProgramID         string                  `json:"programId"`
	ProgressionID     string                  `json:"progressionId"`
	LiftID            *string                 `json:"liftId"`
	Priority          int64                   `json:"priority"`
	Enabled           bool                    `json:"enabled"`
	OverrideIncrement *float64                `json:"overrideIncrement,omitempty"`
	CreatedAt         time.Time               `json:"createdAt"`
	UpdatedAt         time.Time               `json:"updatedAt"`
	Progression       *ProgressionRefResponse `json:"progression,omitempty"`
}

// ProgressionRefResponse is a reference to progression details.
type ProgressionRefResponse struct {
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	Parameters json.RawMessage `json:"parameters"`
}

// ProgramProgressionEnvelope wraps a single program progression response.
type ProgramProgressionEnvelope struct {
	Data ProgramProgressionTestResponse `json:"data"`
}

// PPProgressionEnvelope wraps a single progression response.
type PPProgressionEnvelope struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

// PPCycleEnvelope wraps a single cycle response.
type PPCycleEnvelope struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

// PPProgramEnvelope wraps a single program response.
type PPProgramEnvelope struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

// PPPaginatedProgressionsResponse wraps a paginated list of program progressions.
type PPPaginatedProgressionsResponse struct {
	Data []ProgramProgressionTestResponse `json:"data"`
	Meta *struct {
		Total   int64 `json:"total"`
		Limit   int   `json:"limit"`
		Offset  int   `json:"offset"`
		HasMore bool  `json:"hasMore"`
	} `json:"meta"`
}

// Helper functions for this test file

func authGetProgramProgression(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestUserID)
	return http.DefaultClient.Do(req)
}

func adminPostProgramProgression(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

func adminPutProgramProgression(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

func adminDeleteProgramProgression(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// createPPTestProgression creates a test progression and returns its ID
func createPPTestProgression(t *testing.T, ts *testutil.TestServer, name string, progType string) string {
	t.Helper()
	var body string
	if progType == "LINEAR_PROGRESSION" {
		body = `{"name": "` + name + `", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
	} else {
		body = `{"name": "` + name + `", "type": "CYCLE_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX"}}`
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

	var envelope PPProgressionEnvelope
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// createPPTestCycle creates a test cycle and returns its ID
func createPPTestCycle(t *testing.T, ts *testutil.TestServer, name string) string {
	t.Helper()
	body := `{"name": "` + name + `", "lengthWeeks": 4}`
	resp, err := adminPostCycle(ts.URL("/cycles"), body)
	if err != nil {
		t.Fatalf("Failed to create test cycle: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create test cycle (status %d): %s", resp.StatusCode, bodyBytes)
	}

	var envelope PPCycleEnvelope
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// createPPTestProgram creates a test program and returns its ID
func createPPTestProgram(t *testing.T, ts *testutil.TestServer, name, slug, cycleID string) string {
	t.Helper()
	body := `{"name": "` + name + `", "slug": "` + slug + `", "cycleId": "` + cycleID + `"}`
	resp, err := adminPostProgram(ts.URL("/programs"), body)
	if err != nil {
		t.Fatalf("Failed to create test program: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create test program (status %d): %s", resp.StatusCode, bodyBytes)
	}

	var envelope PPProgramEnvelope
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// createPPTestLift creates a test lift and returns its ID
func createPPTestLift(t *testing.T, ts *testutil.TestServer, name, slug string) string {
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

func TestProgramProgressionCRUD(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Setup: create a program and progression
	cycleID := createPPTestCycle(t, ts, "PP Test Cycle")
	programID := createPPTestProgram(t, ts, "PP Test Program", "pp-test-program", cycleID)
	progressionID := createPPTestProgression(t, ts, "PP Test Progression", "LINEAR_PROGRESSION")

	var createdConfig ProgramProgressionTestResponse

	t.Run("creates program progression with required fields only", func(t *testing.T) {
		body := `{"progressionId": "` + progressionID + `"}`
		resp, err := adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope ProgramProgressionEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		createdConfig = envelope.Data

		if createdConfig.ID == "" {
			t.Error("Expected non-empty ID")
		}
		if createdConfig.ProgramID != programID {
			t.Errorf("Expected programId %s, got %s", programID, createdConfig.ProgramID)
		}
		if createdConfig.ProgressionID != progressionID {
			t.Errorf("Expected progressionId %s, got %s", progressionID, createdConfig.ProgressionID)
		}
		if createdConfig.LiftID != nil {
			t.Errorf("Expected liftId nil, got %s", *createdConfig.LiftID)
		}
		if createdConfig.Priority != 0 {
			t.Errorf("Expected default priority 0, got %d", createdConfig.Priority)
		}
		if !createdConfig.Enabled {
			t.Error("Expected default enabled true, got false")
		}
		if createdConfig.OverrideIncrement != nil {
			t.Errorf("Expected overrideIncrement nil, got %v", createdConfig.OverrideIncrement)
		}
	})

	t.Run("creates program progression with all fields", func(t *testing.T) {
		prog2ID := createPPTestProgression(t, ts, "Another Progression", "CYCLE_PROGRESSION")
		liftID := createPPTestLift(t, ts, "Test Squat", "test-squat-pp")

		body := `{
			"progressionId": "` + prog2ID + `",
			"liftId": "` + liftID + `",
			"priority": 10,
			"enabled": false,
			"overrideIncrement": 10.0
		}`
		resp, err := adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope ProgramProgressionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		config := envelope.Data

		if config.LiftID == nil || *config.LiftID != liftID {
			t.Errorf("Expected liftId %s, got %v", liftID, config.LiftID)
		}
		if config.Priority != 10 {
			t.Errorf("Expected priority 10, got %d", config.Priority)
		}
		if config.Enabled {
			t.Error("Expected enabled false, got true")
		}
		if config.OverrideIncrement == nil || *config.OverrideIncrement != 10.0 {
			t.Errorf("Expected overrideIncrement 10.0, got %v", config.OverrideIncrement)
		}
	})

	t.Run("lists program progressions with progression details", func(t *testing.T) {
		resp, err := authGetProgramProgression(ts.URL("/programs/" + programID + "/progressions"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result PPPaginatedProgressionsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		configs := result.Data

		if len(configs) < 2 {
			t.Errorf("Expected at least 2 configs, got %d", len(configs))
		}

		// Check that progression details are included
		for _, config := range configs {
			if config.Progression == nil {
				t.Error("Expected progression details to be included")
			} else {
				if config.Progression.Name == "" {
					t.Error("Expected progression name to be populated")
				}
				if config.Progression.Type == "" {
					t.Error("Expected progression type to be populated")
				}
			}
		}

		// Verify ordering by priority (lower first)
		for i := 1; i < len(configs); i++ {
			if configs[i].Priority < configs[i-1].Priority {
				t.Errorf("Expected configs ordered by priority: %d before %d",
					configs[i-1].Priority, configs[i].Priority)
			}
		}
	})

	t.Run("gets single program progression by ID", func(t *testing.T) {
		resp, err := authGetProgramProgression(ts.URL("/programs/" + programID + "/progressions/" + createdConfig.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope ProgramProgressionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		config := envelope.Data

		if config.ID != createdConfig.ID {
			t.Errorf("Expected ID %s, got %s", createdConfig.ID, config.ID)
		}
	})

	t.Run("updates program progression", func(t *testing.T) {
		body := `{"priority": 5, "enabled": false, "overrideIncrement": 7.5}`
		resp, err := adminPutProgramProgression(ts.URL("/programs/"+programID+"/progressions/"+createdConfig.ID), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope ProgramProgressionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		updated := envelope.Data

		if updated.Priority != 5 {
			t.Errorf("Expected priority 5, got %d", updated.Priority)
		}
		if updated.Enabled {
			t.Error("Expected enabled false, got true")
		}
		if updated.OverrideIncrement == nil || *updated.OverrideIncrement != 7.5 {
			t.Errorf("Expected overrideIncrement 7.5, got %v", updated.OverrideIncrement)
		}
		// Verify programId, progressionId, liftId unchanged
		if updated.ProgramID != createdConfig.ProgramID {
			t.Errorf("Expected programId unchanged")
		}
		if updated.ProgressionID != createdConfig.ProgressionID {
			t.Errorf("Expected progressionId unchanged")
		}
	})

	t.Run("deletes program progression", func(t *testing.T) {
		// Create one to delete
		prog3ID := createPPTestProgression(t, ts, "To Delete Progression", "LINEAR_PROGRESSION")
		body := `{"progressionId": "` + prog3ID + `"}`
		createResp, _ := adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
		var createEnvelope ProgramProgressionEnvelope
		json.NewDecoder(createResp.Body).Decode(&createEnvelope)
		toDelete := createEnvelope.Data
		createResp.Body.Close()

		resp, err := adminDeleteProgramProgression(ts.URL("/programs/" + programID + "/progressions/" + toDelete.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 204, got %d: %s", resp.StatusCode, bodyBytes)
		}

		// Verify it's deleted
		getResp, _ := authGetProgramProgression(ts.URL("/programs/" + programID + "/progressions/" + toDelete.ID))
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404 after delete, got %d", getResp.StatusCode)
		}
	})
}

func TestProgramProgressionValidation(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	cycleID := createPPTestCycle(t, ts, "Validation Test Cycle")
	programID := createPPTestProgram(t, ts, "Validation Test Program", "validation-test-program", cycleID)
	progressionID := createPPTestProgression(t, ts, "Validation Test Progression", "LINEAR_PROGRESSION")

	t.Run("returns 404 for non-existent program", func(t *testing.T) {
		resp, err := authGetProgramProgression(ts.URL("/programs/non-existent-program/progressions"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 404 for non-existent config", func(t *testing.T) {
		resp, err := authGetProgramProgression(ts.URL("/programs/" + programID + "/progressions/non-existent-config"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for missing progressionId", func(t *testing.T) {
		body := `{}`
		resp, err := adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for non-existent progressionId", func(t *testing.T) {
		body := `{"progressionId": "non-existent-progression"}`
		resp, err := adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 400, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("returns 400 for non-existent liftId", func(t *testing.T) {
		body := `{"progressionId": "` + progressionID + `", "liftId": "non-existent-lift"}`
		resp, err := adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 400, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("returns 400 for negative priority", func(t *testing.T) {
		body := `{"progressionId": "` + progressionID + `", "priority": -1}`
		resp, err := adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for non-positive overrideIncrement", func(t *testing.T) {
		body := `{"progressionId": "` + progressionID + `", "overrideIncrement": 0}`
		resp, err := adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestProgramProgressionDuplicate(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	cycleID := createPPTestCycle(t, ts, "Duplicate Test Cycle")
	programID := createPPTestProgram(t, ts, "Duplicate Test Program", "duplicate-test-program", cycleID)
	progressionID := createPPTestProgression(t, ts, "Duplicate Test Progression", "LINEAR_PROGRESSION")
	liftID := createPPTestLift(t, ts, "Duplicate Test Lift", "duplicate-test-lift")

	t.Run("returns 409 for duplicate program-progression combination", func(t *testing.T) {
		// Create first config (program-wide)
		body := `{"progressionId": "` + progressionID + `"}`
		resp, _ := adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
		resp.Body.Close()

		// Try to create duplicate
		resp, err := adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 409, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("returns 409 for duplicate program-progression-lift combination", func(t *testing.T) {
		// Create config with liftId
		body := `{"progressionId": "` + progressionID + `", "liftId": "` + liftID + `"}`
		resp, _ := adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
		resp.Body.Close()

		// Try to create duplicate
		resp, err := adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 409, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("allows same progression with different liftId", func(t *testing.T) {
		lift2ID := createPPTestLift(t, ts, "Different Lift", "different-lift-dup")

		body := `{"progressionId": "` + progressionID + `", "liftId": "` + lift2ID + `"}`
		resp, err := adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
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

func TestProgramProgressionAuthorization(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	cycleID := createPPTestCycle(t, ts, "Auth Test Cycle")
	programID := createPPTestProgram(t, ts, "Auth Test Program", "auth-test-program-pp", cycleID)
	progressionID := createPPTestProgression(t, ts, "Auth Test Progression", "LINEAR_PROGRESSION")

	// Create a config as admin
	body := `{"progressionId": "` + progressionID + `"}`
	createResp, _ := adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
	var createEnvelope ProgramProgressionEnvelope
	json.NewDecoder(createResp.Body).Decode(&createEnvelope)
	createdConfig := createEnvelope.Data
	createResp.Body.Close()

	t.Run("unauthenticated user gets 401 on GET list", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/programs/" + programID + "/progressions"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("unauthenticated user gets 401 on GET single", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/programs/" + programID + "/progressions/" + createdConfig.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("authenticated user can GET list", func(t *testing.T) {
		resp, err := authGetProgramProgression(ts.URL("/programs/" + programID + "/progressions"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("non-admin user gets 403 on POST", func(t *testing.T) {
		prog2ID := createPPTestProgression(t, ts, "Non-admin Test Prog", "LINEAR_PROGRESSION")
		body := `{"progressionId": "` + prog2ID + `"}`
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/programs/"+programID+"/progressions"), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", testutil.TestUserID)
		// Not setting X-Admin

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("non-admin user gets 403 on PUT", func(t *testing.T) {
		body := `{"priority": 100}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/programs/"+programID+"/progressions/"+createdConfig.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", testutil.TestUserID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("non-admin user gets 403 on DELETE", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL("/programs/"+programID+"/progressions/"+createdConfig.ID), nil)
		req.Header.Set("X-User-ID", testutil.TestUserID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})
}

func TestProgramProgressionResponseFormat(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	cycleID := createPPTestCycle(t, ts, "Format Test Cycle")
	programID := createPPTestProgram(t, ts, "Format Test Program", "format-test-program", cycleID)
	progressionID := createPPTestProgression(t, ts, "Format Test Progression", "LINEAR_PROGRESSION")

	body := `{"progressionId": "` + progressionID + `"}`
	createResp, _ := adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
	var createEnvelope ProgramProgressionEnvelope
	json.NewDecoder(createResp.Body).Decode(&createEnvelope)
	createdConfig := createEnvelope.Data
	createResp.Body.Close()

	t.Run("single config response has correct JSON field names", func(t *testing.T) {
		resp, _ := authGetProgramProgression(ts.URL("/programs/" + programID + "/progressions/" + createdConfig.ID))
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// Check camelCase field names per ticket spec
		expectedFields := []string{
			`"id"`,
			`"programId"`,
			`"progressionId"`,
			`"priority"`,
			`"enabled"`,
			`"createdAt"`,
			`"updatedAt"`,
		}

		for _, field := range expectedFields {
			if !bytes.Contains(body, []byte(field)) {
				t.Errorf("Expected field %s in response, body: %s", field, bodyStr)
			}
		}
	})

	t.Run("list response includes progression details", func(t *testing.T) {
		resp, _ := authGetProgramProgression(ts.URL("/programs/" + programID + "/progressions"))
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// Check for progression sub-object
		expectedFields := []string{
			`"progression"`,
			`"name"`,
			`"type"`,
			`"parameters"`,
		}

		for _, field := range expectedFields {
			if !bytes.Contains(body, []byte(field)) {
				t.Errorf("Expected field %s in response, body: %s", field, bodyStr)
			}
		}
	})
}

func TestProgramProgressionPriorityOrdering(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	cycleID := createPPTestCycle(t, ts, "Priority Test Cycle")
	programID := createPPTestProgram(t, ts, "Priority Test Program", "priority-test-program", cycleID)

	// Create progressions with different priorities
	prog1ID := createPPTestProgression(t, ts, "Priority 10 Prog", "LINEAR_PROGRESSION")
	prog2ID := createPPTestProgression(t, ts, "Priority 0 Prog", "LINEAR_PROGRESSION")
	prog3ID := createPPTestProgression(t, ts, "Priority 5 Prog", "LINEAR_PROGRESSION")

	// Create configs out of order
	body := `{"progressionId": "` + prog1ID + `", "priority": 10}`
	resp, _ := adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
	resp.Body.Close()

	body = `{"progressionId": "` + prog2ID + `", "priority": 0}`
	resp, _ = adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
	resp.Body.Close()

	body = `{"progressionId": "` + prog3ID + `", "priority": 5}`
	resp, _ = adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
	resp.Body.Close()

	t.Run("list returns configs ordered by priority ascending", func(t *testing.T) {
		resp, err := authGetProgramProgression(ts.URL("/programs/" + programID + "/progressions"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PPPaginatedProgressionsResponse
		json.NewDecoder(resp.Body).Decode(&result)
		configs := result.Data

		if len(configs) != 3 {
			t.Fatalf("Expected 3 configs, got %d", len(configs))
		}

		// Should be ordered: priority 0, priority 5, priority 10
		expectedPriorities := []int64{0, 5, 10}
		for i, expected := range expectedPriorities {
			if configs[i].Priority != expected {
				t.Errorf("Expected priority %d at index %d, got %d", expected, i, configs[i].Priority)
			}
		}
	})
}

func TestProgramProgressionCrossProgram(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	cycleID := createPPTestCycle(t, ts, "Cross Program Test Cycle")
	program1ID := createPPTestProgram(t, ts, "Program 1", "cross-program-1", cycleID)
	program2ID := createPPTestProgram(t, ts, "Program 2", "cross-program-2", cycleID)
	progressionID := createPPTestProgression(t, ts, "Cross Program Progression", "LINEAR_PROGRESSION")

	// Create config for program 1
	body := `{"progressionId": "` + progressionID + `"}`
	resp, _ := adminPostProgramProgression(ts.URL("/programs/"+program1ID+"/progressions"), body)
	var config1Envelope ProgramProgressionEnvelope
	json.NewDecoder(resp.Body).Decode(&config1Envelope)
	config1 := config1Envelope.Data
	resp.Body.Close()

	t.Run("config cannot be accessed from wrong program", func(t *testing.T) {
		// Try to get program1's config via program2's URL
		resp, err := authGetProgramProgression(ts.URL("/programs/" + program2ID + "/progressions/" + config1.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("config cannot be updated from wrong program", func(t *testing.T) {
		body := `{"priority": 99}`
		resp, err := adminPutProgramProgression(ts.URL("/programs/"+program2ID+"/progressions/"+config1.ID), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("config cannot be deleted from wrong program", func(t *testing.T) {
		resp, err := adminDeleteProgramProgression(ts.URL("/programs/" + program2ID + "/progressions/" + config1.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestDeleteProgressionReferencedByProgramProgression(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	cycleID := createPPTestCycle(t, ts, "Delete Ref Test Cycle")
	programID := createPPTestProgram(t, ts, "Delete Ref Test Program", "delete-ref-test-program", cycleID)
	progressionID := createPPTestProgression(t, ts, "Referenced Progression", "LINEAR_PROGRESSION")

	// Create a program progression that references this progression
	body := `{"progressionId": "` + progressionID + `"}`
	resp, _ := adminPostProgramProgression(ts.URL("/programs/"+programID+"/progressions"), body)
	resp.Body.Close()

	t.Run("cannot delete progression that is referenced by program progressions", func(t *testing.T) {
		resp, err := adminDelete(ts.URL("/progressions/" + progressionID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 409, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})
}
