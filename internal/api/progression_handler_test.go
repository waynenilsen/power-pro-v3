package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// ProgressionResponse matches the API response format for a progression.
type ProgressionResponse struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	Parameters json.RawMessage `json:"parameters"`
	CreatedAt  string          `json:"createdAt"`
	UpdatedAt  string          `json:"updatedAt"`
}

// PaginatedMeta contains pagination metadata in the standard envelope.
type PaginatedMeta struct {
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasMore bool  `json:"hasMore"`
}

// PaginatedProgressionsResponse is the paginated list response.
type PaginatedProgressionsResponse struct {
	Data []ProgressionResponse `json:"data"`
	Meta PaginatedMeta         `json:"meta"`
}

// ProgressionEnvelope wraps single progression response with standard envelope.
type ProgressionEnvelope struct {
	Data ProgressionResponse `json:"data"`
}

func TestListProgressions(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("returns seeded progressions from canonical programs", func(t *testing.T) {
		resp, err := authGet(ts.URL("/progressions"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}

		var result PaginatedProgressionsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Canonical programs seed progressions - just verify we have some
		if result.Meta.Total == 0 {
			t.Errorf("Expected seeded progressions from canonical programs, got 0")
		}
	})

	t.Run("supports pagination", func(t *testing.T) {
		// Get baseline count of seeded progressions
		baselineResp, _ := authGet(ts.URL("/progressions"))
		var baselineResult PaginatedProgressionsResponse
		json.NewDecoder(baselineResp.Body).Decode(&baselineResult)
		baselineResp.Body.Close()
		baselineCount := baselineResult.Meta.Total

		// Create 3 additional progressions
		for i := 0; i < 3; i++ {
			body := `{"name": "Linear Prog Test ` + string(rune('A'+i)) + `", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
			resp, _ := adminPost(ts.URL("/progressions"), body)
			resp.Body.Close()
		}

		resp, err := authGet(ts.URL("/progressions?limit=2&offset=0"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedProgressionsResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) != 2 {
			t.Errorf("Expected 2 progressions with limit=2, got %d", len(result.Data))
		}
		if result.Meta.Offset != 0 {
			t.Errorf("Expected offset 0, got %d", result.Meta.Offset)
		}
		if result.Meta.Limit != 2 {
			t.Errorf("Expected limit 2, got %d", result.Meta.Limit)
		}
		expectedTotal := baselineCount + 3
		if result.Meta.Total != expectedTotal {
			t.Errorf("Expected total %d, got %d", expectedTotal, result.Meta.Total)
		}
		if !result.Meta.HasMore {
			t.Errorf("Expected hasMore to be true")
		}
	})

	t.Run("filters by type", func(t *testing.T) {
		// Get baseline counts for each type
		baselineLinearResp, _ := authGet(ts.URL("/progressions?type=LINEAR_PROGRESSION"))
		var baselineLinearResult PaginatedProgressionsResponse
		json.NewDecoder(baselineLinearResp.Body).Decode(&baselineLinearResult)
		baselineLinearResp.Body.Close()
		baselineLinearCount := baselineLinearResult.Meta.Total

		baselineCycleResp, _ := authGet(ts.URL("/progressions?type=CYCLE_PROGRESSION"))
		var baselineCycleResult PaginatedProgressionsResponse
		json.NewDecoder(baselineCycleResp.Body).Decode(&baselineCycleResult)
		baselineCycleResp.Body.Close()
		baselineCycleCount := baselineCycleResult.Meta.Total

		// Create additional LINEAR_PROGRESSION
		linearBody := `{"name": "Linear Test Filter", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
		resp, _ := adminPost(ts.URL("/progressions"), linearBody)
		resp.Body.Close()

		// Create CYCLE_PROGRESSION
		cycleBody := `{"name": "Cycle Test Filter", "type": "CYCLE_PROGRESSION", "parameters": {"increment": 10.0, "maxType": "TRAINING_MAX"}}`
		resp, _ = adminPost(ts.URL("/progressions"), cycleBody)
		resp.Body.Close()

		// Filter by LINEAR_PROGRESSION
		resp, err := authGet(ts.URL("/progressions?type=LINEAR_PROGRESSION"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedProgressionsResponse
		json.NewDecoder(resp.Body).Decode(&result)

		expectedLinearCount := baselineLinearCount + 1
		if result.Meta.Total != expectedLinearCount {
			t.Errorf("Expected %d LINEAR_PROGRESSION, got %d", expectedLinearCount, result.Meta.Total)
		}
		for _, prog := range result.Data {
			if prog.Type != "LINEAR_PROGRESSION" {
				t.Errorf("Expected type LINEAR_PROGRESSION, got %s", prog.Type)
			}
		}

		// Filter by CYCLE_PROGRESSION
		resp2, _ := authGet(ts.URL("/progressions?type=CYCLE_PROGRESSION"))
		defer resp2.Body.Close()

		var result2 PaginatedProgressionsResponse
		json.NewDecoder(resp2.Body).Decode(&result2)

		expectedCycleCount := baselineCycleCount + 1
		if result2.Meta.Total != expectedCycleCount {
			t.Errorf("Expected %d CYCLE_PROGRESSION, got %d", expectedCycleCount, result2.Meta.Total)
		}
		for _, prog := range result2.Data {
			if prog.Type != "CYCLE_PROGRESSION" {
				t.Errorf("Expected type CYCLE_PROGRESSION, got %s", prog.Type)
			}
		}
	})
}

func TestGetProgression(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a progression first
	createBody := `{"name": "Test Progression", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
	createResp, _ := adminPost(ts.URL("/progressions"), createBody)
	var createEnvelope ProgressionEnvelope
	json.NewDecoder(createResp.Body).Decode(&createEnvelope)
	created := createEnvelope.Data
	createResp.Body.Close()

	t.Run("returns progression by ID", func(t *testing.T) {
		resp, err := authGet(ts.URL("/progressions/" + created.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var envelope ProgressionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		prog := envelope.Data

		if prog.ID != created.ID {
			t.Errorf("Expected ID %s, got %s", created.ID, prog.ID)
		}
		if prog.Name != "Test Progression" {
			t.Errorf("Expected name 'Test Progression', got %s", prog.Name)
		}
		if prog.Type != "LINEAR_PROGRESSION" {
			t.Errorf("Expected type 'LINEAR_PROGRESSION', got %s", prog.Type)
		}
	})

	t.Run("returns 404 for non-existent ID", func(t *testing.T) {
		resp, _ := authGet(ts.URL("/progressions/non-existent-id"))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)

		if !strings.Contains(errResp.Error.Message, "progression not found") {
			t.Errorf("Expected error to contain 'progression not found', got %s", errResp.Error.Message)
		}
	})
}

func TestCreateProgression(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("creates LINEAR_PROGRESSION", func(t *testing.T) {
		body := `{"name": "Starting Strength Squat", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
		resp, err := adminPost(ts.URL("/progressions"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, body)
		}

		var envelope ProgressionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		prog := envelope.Data

		if prog.Name != "Starting Strength Squat" {
			t.Errorf("Expected name 'Starting Strength Squat', got %s", prog.Name)
		}
		if prog.Type != "LINEAR_PROGRESSION" {
			t.Errorf("Expected type 'LINEAR_PROGRESSION', got %s", prog.Type)
		}
		if prog.ID == "" {
			t.Errorf("Expected ID to be generated")
		}

		// Verify parameters
		var params struct {
			Increment   float64 `json:"increment"`
			MaxType     string  `json:"maxType"`
			TriggerType string  `json:"triggerType"`
		}
		json.Unmarshal(prog.Parameters, &params)
		if params.Increment != 5.0 {
			t.Errorf("Expected increment 5.0, got %f", params.Increment)
		}
		if params.MaxType != "TRAINING_MAX" {
			t.Errorf("Expected maxType TRAINING_MAX, got %s", params.MaxType)
		}
		if params.TriggerType != "AFTER_SESSION" {
			t.Errorf("Expected triggerType AFTER_SESSION, got %s", params.TriggerType)
		}
	})

	t.Run("creates CYCLE_PROGRESSION", func(t *testing.T) {
		body := `{"name": "5/3/1 Lower Body", "type": "CYCLE_PROGRESSION", "parameters": {"increment": 10.0, "maxType": "TRAINING_MAX"}}`
		resp, err := adminPost(ts.URL("/progressions"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, body)
		}

		var envelope ProgressionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		prog := envelope.Data

		if prog.Type != "CYCLE_PROGRESSION" {
			t.Errorf("Expected type 'CYCLE_PROGRESSION', got %s", prog.Type)
		}

		// Verify parameters
		var params struct {
			Increment float64 `json:"increment"`
			MaxType   string  `json:"maxType"`
		}
		json.Unmarshal(prog.Parameters, &params)
		if params.Increment != 10.0 {
			t.Errorf("Expected increment 10.0, got %f", params.Increment)
		}
	})

	t.Run("returns 400 for missing name", func(t *testing.T) {
		body := `{"type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
		resp, _ := adminPost(ts.URL("/progressions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid type", func(t *testing.T) {
		body := `{"name": "Invalid", "type": "INVALID_TYPE", "parameters": {}}`
		resp, _ := adminPost(ts.URL("/progressions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for negative increment in LINEAR_PROGRESSION", func(t *testing.T) {
		body := `{"name": "Invalid", "type": "LINEAR_PROGRESSION", "parameters": {"increment": -5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
		resp, _ := adminPost(ts.URL("/progressions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)

		// Should have validation error about increment - details is now in Error.Details
		detailsStr := fmt.Sprintf("%v", errResp.Error.Details)
		if !strings.Contains(detailsStr, "increment must be positive") {
			t.Errorf("Expected 'increment must be positive' in details, got %v", errResp.Error.Details)
		}
	})

	t.Run("returns 400 for zero increment in CYCLE_PROGRESSION", func(t *testing.T) {
		body := `{"name": "Invalid", "type": "CYCLE_PROGRESSION", "parameters": {"increment": 0.0, "maxType": "TRAINING_MAX"}}`
		resp, _ := adminPost(ts.URL("/progressions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid maxType", func(t *testing.T) {
		body := `{"name": "Invalid", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "INVALID", "triggerType": "AFTER_SESSION"}}`
		resp, _ := adminPost(ts.URL("/progressions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid triggerType in LINEAR_PROGRESSION", func(t *testing.T) {
		body := `{"name": "Invalid", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "INVALID"}}`
		resp, _ := adminPost(ts.URL("/progressions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for AFTER_CYCLE trigger in LINEAR_PROGRESSION", func(t *testing.T) {
		body := `{"name": "Invalid", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_CYCLE"}}`
		resp, _ := adminPost(ts.URL("/progressions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)

		// Should have validation error about trigger type - details is now in Error.Details
		detailsStr := fmt.Sprintf("%v", errResp.Error.Details)
		if !strings.Contains(detailsStr, "linear progression only supports AFTER_SESSION and AFTER_WEEK triggers") {
			t.Errorf("Expected trigger type error in details, got %v", errResp.Error.Details)
		}
	})

	t.Run("supports ONE_RM maxType", func(t *testing.T) {
		body := `{"name": "1RM Progression", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 2.5, "maxType": "ONE_RM", "triggerType": "AFTER_WEEK"}}`
		resp, err := adminPost(ts.URL("/progressions"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("requires auth", func(t *testing.T) {
		body := `{"name": "No Auth", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/progressions"), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		// No auth headers

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("requires admin", func(t *testing.T) {
		body := `{"name": "No Admin", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/progressions"), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", testutil.TestUserID)
		// No admin header

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

func TestUpdateProgression(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a progression to update
	createBody := `{"name": "Original Name", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
	createResp, _ := adminPost(ts.URL("/progressions"), createBody)
	var createEnvelope ProgressionEnvelope
	json.NewDecoder(createResp.Body).Decode(&createEnvelope)
	created := createEnvelope.Data
	createResp.Body.Close()

	t.Run("updates progression name", func(t *testing.T) {
		body := `{"name": "Updated Name"}`
		resp, err := adminPut(ts.URL("/progressions/"+created.ID), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope ProgressionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		prog := envelope.Data

		if prog.Name != "Updated Name" {
			t.Errorf("Expected name 'Updated Name', got %s", prog.Name)
		}
		// Type should remain unchanged
		if prog.Type != "LINEAR_PROGRESSION" {
			t.Errorf("Expected type to remain LINEAR_PROGRESSION, got %s", prog.Type)
		}
	})

	t.Run("updates progression type and parameters", func(t *testing.T) {
		body := `{"type": "CYCLE_PROGRESSION", "parameters": {"increment": 10.0, "maxType": "TRAINING_MAX"}}`
		resp, err := adminPut(ts.URL("/progressions/"+created.ID), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope ProgressionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		prog := envelope.Data

		if prog.Type != "CYCLE_PROGRESSION" {
			t.Errorf("Expected type 'CYCLE_PROGRESSION', got %s", prog.Type)
		}
	})

	t.Run("returns 404 for non-existent progression", func(t *testing.T) {
		body := `{"name": "Updated"}`
		resp, _ := adminPut(ts.URL("/progressions/non-existent-id"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for validation errors", func(t *testing.T) {
		body := `{"name": ""}`
		resp, _ := adminPut(ts.URL("/progressions/"+created.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid type", func(t *testing.T) {
		body := `{"type": "INVALID_TYPE"}`
		resp, _ := adminPut(ts.URL("/progressions/"+created.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestDeleteProgression(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("deletes progression successfully", func(t *testing.T) {
		// Create a progression to delete
		createBody := `{"name": "To Delete", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
		createResp, _ := adminPost(ts.URL("/progressions"), createBody)
		var createEnvelope ProgressionEnvelope
		json.NewDecoder(createResp.Body).Decode(&createEnvelope)
		created := createEnvelope.Data
		createResp.Body.Close()

		// Delete it
		resp, err := adminDelete(ts.URL("/progressions/" + created.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 204, got %d: %s", resp.StatusCode, bodyBytes)
		}

		// Verify it's deleted
		getResp, _ := authGet(ts.URL("/progressions/" + created.ID))
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected deleted progression to return 404, got %d", getResp.StatusCode)
		}
	})

	t.Run("returns 404 for non-existent progression", func(t *testing.T) {
		resp, _ := adminDelete(ts.URL("/progressions/non-existent-id"))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	// Note: The 409 conflict test for referenced progressions would require
	// setting up a program_progression reference, which depends on Program
	// entity. We'll add this test when ProgramProgression API is implemented.
}

func TestProgressionResponseFormat(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a progression
	createBody := `{"name": "Format Test", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
	createResp, _ := adminPost(ts.URL("/progressions"), createBody)
	var createEnvelope ProgressionEnvelope
	json.NewDecoder(createResp.Body).Decode(&createEnvelope)
	created := createEnvelope.Data
	createResp.Body.Close()

	t.Run("response has correct JSON field names", func(t *testing.T) {
		resp, _ := authGet(ts.URL("/progressions/" + created.ID))
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// Check camelCase field names per ERD spec (inside data envelope)
		expectedFields := []string{
			`"data"`,
			`"id"`,
			`"name"`,
			`"type"`,
			`"parameters"`,
			`"createdAt"`,
			`"updatedAt"`,
		}

		for _, field := range expectedFields {
			if !bytes.Contains(body, []byte(field)) {
				t.Errorf("Expected field %s in response, body: %s", field, bodyStr)
			}
		}
	})

	t.Run("parameters contains expected structure for LINEAR_PROGRESSION", func(t *testing.T) {
		resp, _ := authGet(ts.URL("/progressions/" + created.ID))
		defer resp.Body.Close()

		var envelope ProgressionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		prog := envelope.Data

		var params map[string]interface{}
		if err := json.Unmarshal(prog.Parameters, &params); err != nil {
			t.Fatalf("Failed to parse parameters: %v", err)
		}

		if _, ok := params["increment"]; !ok {
			t.Error("Expected 'increment' in parameters")
		}
		if _, ok := params["maxType"]; !ok {
			t.Error("Expected 'maxType' in parameters")
		}
		if _, ok := params["triggerType"]; !ok {
			t.Error("Expected 'triggerType' in parameters")
		}
	})
}
