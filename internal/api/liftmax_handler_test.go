package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// LiftMaxResponse matches the API response format.
type LiftMaxResponse struct {
	ID            string    `json:"id"`
	UserID        string    `json:"userId"`
	LiftID        string    `json:"liftId"`
	Type          string    `json:"type"`
	Value         float64   `json:"value"`
	EffectiveDate time.Time `json:"effectiveDate"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// PaginatedLiftMaxesResponse is the paginated list response.
type PaginatedLiftMaxesResponse struct {
	Data       []LiftMaxResponse `json:"data"`
	Page       int               `json:"page"`
	PageSize   int               `json:"pageSize"`
	TotalItems int64             `json:"totalItems"`
	TotalPages int64             `json:"totalPages"`
}

// LiftMaxWithWarningsResponse wraps response with warnings.
type LiftMaxWithWarningsResponse struct {
	Data     LiftMaxResponse `json:"data"`
	Warnings []string        `json:"warnings,omitempty"`
}

// Test constants
const (
	testUserID  = "test-user-001"
	testSquatID = "00000000-0000-0000-0000-000000000001"
	testBenchID = "00000000-0000-0000-0000-000000000002"
)

func TestListLiftMaxes(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create some test maxes first
	createMax(t, ts, testUserID, testSquatID, "ONE_RM", 315.0, nil)
	createMax(t, ts, testUserID, testSquatID, "TRAINING_MAX", 285.0, nil)
	createMax(t, ts, testUserID, testBenchID, "ONE_RM", 225.0, nil)

	t.Run("returns user's lift maxes", func(t *testing.T) {
		resp, err := http.Get(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", testUserID)))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}

		var result PaginatedLiftMaxesResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(result.Data) != 3 {
			t.Errorf("Expected 3 lift maxes, got %d", len(result.Data))
		}
		if result.TotalItems != 3 {
			t.Errorf("Expected totalItems 3, got %d", result.TotalItems)
		}
	})

	t.Run("supports pagination", func(t *testing.T) {
		resp, err := http.Get(ts.URL(fmt.Sprintf("/users/%s/lift-maxes?page=1&pageSize=2", testUserID)))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedLiftMaxesResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) != 2 {
			t.Errorf("Expected 2 lift maxes on page 1, got %d", len(result.Data))
		}
		if result.Page != 1 {
			t.Errorf("Expected page 1, got %d", result.Page)
		}
		if result.PageSize != 2 {
			t.Errorf("Expected pageSize 2, got %d", result.PageSize)
		}
		if result.TotalPages != 2 {
			t.Errorf("Expected totalPages 2, got %d", result.TotalPages)
		}

		// Get page 2
		resp2, _ := http.Get(ts.URL(fmt.Sprintf("/users/%s/lift-maxes?page=2&pageSize=2", testUserID)))
		defer resp2.Body.Close()

		var result2 PaginatedLiftMaxesResponse
		json.NewDecoder(resp2.Body).Decode(&result2)

		if len(result2.Data) != 1 {
			t.Errorf("Expected 1 lift max on page 2, got %d", len(result2.Data))
		}
	})

	t.Run("filters by lift_id", func(t *testing.T) {
		resp, err := http.Get(ts.URL(fmt.Sprintf("/users/%s/lift-maxes?lift_id=%s", testUserID, testSquatID)))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedLiftMaxesResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) != 2 {
			t.Errorf("Expected 2 squat maxes, got %d", len(result.Data))
		}

		for _, max := range result.Data {
			if max.LiftID != testSquatID {
				t.Errorf("Expected all maxes to be for squat, got lift %s", max.LiftID)
			}
		}
	})

	t.Run("filters by type", func(t *testing.T) {
		resp, err := http.Get(ts.URL(fmt.Sprintf("/users/%s/lift-maxes?type=ONE_RM", testUserID)))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedLiftMaxesResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) != 2 {
			t.Errorf("Expected 2 ONE_RM maxes, got %d", len(result.Data))
		}

		for _, max := range result.Data {
			if max.Type != "ONE_RM" {
				t.Errorf("Expected all maxes to be ONE_RM, got %s", max.Type)
			}
		}
	})

	t.Run("filters by lift_id and type combined", func(t *testing.T) {
		resp, err := http.Get(ts.URL(fmt.Sprintf("/users/%s/lift-maxes?lift_id=%s&type=ONE_RM", testUserID, testSquatID)))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedLiftMaxesResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) != 1 {
			t.Errorf("Expected 1 squat ONE_RM max, got %d", len(result.Data))
		}

		if len(result.Data) > 0 {
			if result.Data[0].LiftID != testSquatID {
				t.Errorf("Expected squat lift, got %s", result.Data[0].LiftID)
			}
			if result.Data[0].Type != "ONE_RM" {
				t.Errorf("Expected ONE_RM type, got %s", result.Data[0].Type)
			}
		}
	})

	t.Run("returns empty for non-existent user", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/users/non-existent-user/lift-maxes"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedLiftMaxesResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) != 0 {
			t.Errorf("Expected 0 lift maxes for non-existent user, got %d", len(result.Data))
		}
	})

	t.Run("returns 400 for invalid type filter", func(t *testing.T) {
		resp, err := http.Get(ts.URL(fmt.Sprintf("/users/%s/lift-maxes?type=INVALID", testUserID)))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("sorts by effective_date descending by default", func(t *testing.T) {
		// Create a new user with specific dates to test ordering
		newUserID := "sort-test-user"
		now := time.Now()
		yesterday := now.Add(-24 * time.Hour)
		lastWeek := now.Add(-7 * 24 * time.Hour)

		createMax(t, ts, newUserID, testSquatID, "ONE_RM", 300.0, &lastWeek)
		createMax(t, ts, newUserID, testSquatID, "ONE_RM", 310.0, &yesterday)
		createMax(t, ts, newUserID, testSquatID, "ONE_RM", 320.0, &now)

		resp, err := http.Get(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", newUserID)))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedLiftMaxesResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) < 3 {
			t.Fatalf("Expected at least 3 maxes, got %d", len(result.Data))
		}

		// Most recent should be first (highest value 320)
		if result.Data[0].Value != 320.0 {
			t.Errorf("Expected first max value to be 320.0, got %f", result.Data[0].Value)
		}
	})
}

func TestGetLiftMax(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a lift max
	created := createMax(t, ts, testUserID, testSquatID, "ONE_RM", 315.0, nil)

	t.Run("returns lift max by ID", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/lift-maxes/" + created.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var max LiftMaxResponse
		json.NewDecoder(resp.Body).Decode(&max)

		if max.ID != created.ID {
			t.Errorf("Expected ID %s, got %s", created.ID, max.ID)
		}
		if max.UserID != testUserID {
			t.Errorf("Expected userId %s, got %s", testUserID, max.UserID)
		}
		if max.LiftID != testSquatID {
			t.Errorf("Expected liftId %s, got %s", testSquatID, max.LiftID)
		}
		if max.Type != "ONE_RM" {
			t.Errorf("Expected type ONE_RM, got %s", max.Type)
		}
		if max.Value != 315.0 {
			t.Errorf("Expected value 315.0, got %f", max.Value)
		}
	})

	t.Run("returns 404 for non-existent ID", func(t *testing.T) {
		resp, _ := http.Get(ts.URL("/lift-maxes/non-existent-id"))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)

		if errResp.Error != "Lift max not found" {
			t.Errorf("Expected error 'Lift max not found', got %s", errResp.Error)
		}
	})
}

func TestCreateLiftMax(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("creates lift max with all fields", func(t *testing.T) {
		body := fmt.Sprintf(`{"liftId": "%s", "type": "ONE_RM", "value": 405.0}`, testSquatID)
		resp, err := http.Post(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", "create-test-user")), "application/json", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, body)
		}

		var max LiftMaxResponse
		json.NewDecoder(resp.Body).Decode(&max)

		if max.LiftID != testSquatID {
			t.Errorf("Expected liftId %s, got %s", testSquatID, max.LiftID)
		}
		if max.Type != "ONE_RM" {
			t.Errorf("Expected type ONE_RM, got %s", max.Type)
		}
		if max.Value != 405.0 {
			t.Errorf("Expected value 405.0, got %f", max.Value)
		}
		if max.ID == "" {
			t.Errorf("Expected ID to be generated")
		}
	})

	t.Run("creates lift max with custom effective date", func(t *testing.T) {
		effectiveDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
		body := fmt.Sprintf(`{"liftId": "%s", "type": "ONE_RM", "value": 410.0, "effectiveDate": "%s"}`, testSquatID, effectiveDate.Format(time.RFC3339))
		resp, _ := http.Post(ts.URL("/users/date-test-user/lift-maxes"), "application/json", bytes.NewBufferString(body))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var max LiftMaxResponse
		json.NewDecoder(resp.Body).Decode(&max)

		if max.EffectiveDate.Year() != 2024 || max.EffectiveDate.Month() != 1 || max.EffectiveDate.Day() != 15 {
			t.Errorf("Expected effectiveDate 2024-01-15, got %s", max.EffectiveDate)
		}
	})

	t.Run("returns 400 for missing liftId", func(t *testing.T) {
		body := `{"type": "ONE_RM", "value": 315.0}`
		resp, _ := http.Post(ts.URL("/users/missing-lift-user/lift-maxes"), "application/json", bytes.NewBufferString(body))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for missing type", func(t *testing.T) {
		body := fmt.Sprintf(`{"liftId": "%s", "value": 315.0}`, testSquatID)
		resp, _ := http.Post(ts.URL("/users/missing-type-user/lift-maxes"), "application/json", bytes.NewBufferString(body))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid type", func(t *testing.T) {
		body := fmt.Sprintf(`{"liftId": "%s", "type": "INVALID", "value": 315.0}`, testSquatID)
		resp, _ := http.Post(ts.URL("/users/invalid-type-user/lift-maxes"), "application/json", bytes.NewBufferString(body))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid value precision", func(t *testing.T) {
		body := fmt.Sprintf(`{"liftId": "%s", "type": "ONE_RM", "value": 315.33}`, testSquatID)
		resp, _ := http.Post(ts.URL("/users/invalid-precision-user/lift-maxes"), "application/json", bytes.NewBufferString(body))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for non-positive value", func(t *testing.T) {
		body := fmt.Sprintf(`{"liftId": "%s", "type": "ONE_RM", "value": 0}`, testSquatID)
		resp, _ := http.Post(ts.URL("/users/zero-value-user/lift-maxes"), "application/json", bytes.NewBufferString(body))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for non-existent lift", func(t *testing.T) {
		body := `{"liftId": "non-existent-lift", "type": "ONE_RM", "value": 315.0}`
		resp, _ := http.Post(ts.URL("/users/bad-lift-user/lift-maxes"), "application/json", bytes.NewBufferString(body))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 409 for duplicate unique constraint", func(t *testing.T) {
		// Create first max
		effectiveDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
		body := fmt.Sprintf(`{"liftId": "%s", "type": "ONE_RM", "value": 315.0, "effectiveDate": "%s"}`, testSquatID, effectiveDate.Format(time.RFC3339))
		resp1, _ := http.Post(ts.URL("/users/duplicate-test-user/lift-maxes"), "application/json", bytes.NewBufferString(body))
		resp1.Body.Close()

		// Try to create duplicate
		resp2, _ := http.Post(ts.URL("/users/duplicate-test-user/lift-maxes"), "application/json", bytes.NewBufferString(body))
		defer resp2.Body.Close()

		if resp2.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409, got %d", resp2.StatusCode)
		}
	})

	t.Run("includes warning for TM outside expected range", func(t *testing.T) {
		// First create a 1RM
		body1 := fmt.Sprintf(`{"liftId": "%s", "type": "ONE_RM", "value": 400.0}`, testSquatID)
		resp1, _ := http.Post(ts.URL("/users/tm-warning-user/lift-maxes"), "application/json", bytes.NewBufferString(body1))
		resp1.Body.Close()

		// Now create a TM that's too low (70% = 280)
		body2 := fmt.Sprintf(`{"liftId": "%s", "type": "TRAINING_MAX", "value": 280.0}`, testSquatID)
		resp2, _ := http.Post(ts.URL("/users/tm-warning-user/lift-maxes"), "application/json", bytes.NewBufferString(body2))
		defer resp2.Body.Close()

		if resp2.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp2.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp2.StatusCode, bodyBytes)
		}

		var result LiftMaxWithWarningsResponse
		json.NewDecoder(resp2.Body).Decode(&result)

		if len(result.Warnings) == 0 {
			t.Errorf("Expected warnings for TM below 80%% of 1RM")
		}
	})
}

func TestUpdateLiftMax(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a lift max to update
	created := createMax(t, ts, "update-user", testSquatID, "ONE_RM", 315.0, nil)

	t.Run("updates lift max value", func(t *testing.T) {
		body := `{"value": 325.0}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/lift-maxes/"+created.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var max LiftMaxResponse
		json.NewDecoder(resp.Body).Decode(&max)

		if max.Value != 325.0 {
			t.Errorf("Expected value 325.0, got %f", max.Value)
		}
		// Type should remain unchanged
		if max.Type != "ONE_RM" {
			t.Errorf("Expected type to remain ONE_RM, got %s", max.Type)
		}
	})

	t.Run("updates lift max effective date", func(t *testing.T) {
		newDate := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
		body := fmt.Sprintf(`{"effectiveDate": "%s"}`, newDate.Format(time.RFC3339))
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/lift-maxes/"+created.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var max LiftMaxResponse
		json.NewDecoder(resp.Body).Decode(&max)

		if max.EffectiveDate.Year() != 2024 || max.EffectiveDate.Month() != 7 {
			t.Errorf("Expected effective date 2024-07, got %s", max.EffectiveDate)
		}
	})

	t.Run("returns 404 for non-existent lift max", func(t *testing.T) {
		body := `{"value": 400.0}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/lift-maxes/non-existent-id"), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid value precision", func(t *testing.T) {
		body := `{"value": 325.33}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/lift-maxes/"+created.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestDeleteLiftMax(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("deletes lift max successfully", func(t *testing.T) {
		// Create a lift max to delete
		created := createMax(t, ts, "delete-user", testSquatID, "ONE_RM", 315.0, nil)

		// Delete it
		req, _ := http.NewRequest(http.MethodDelete, ts.URL("/lift-maxes/"+created.ID), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 204, got %d: %s", resp.StatusCode, bodyBytes)
		}

		// Verify it's deleted
		getResp, _ := http.Get(ts.URL("/lift-maxes/" + created.ID))
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected deleted lift max to return 404, got %d", getResp.StatusCode)
		}
	})

	t.Run("returns 404 for non-existent lift max", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL("/lift-maxes/non-existent-id"), nil)
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestLiftMaxResponseFormat(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a lift max
	created := createMax(t, ts, "format-user", testSquatID, "ONE_RM", 315.0, nil)

	t.Run("response has correct JSON field names", func(t *testing.T) {
		resp, _ := http.Get(ts.URL("/lift-maxes/" + created.ID))
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// Check camelCase field names per ERD spec
		expectedFields := []string{
			`"id"`,
			`"userId"`,
			`"liftId"`,
			`"type"`,
			`"value"`,
			`"effectiveDate"`,
			`"createdAt"`,
			`"updatedAt"`,
		}

		for _, field := range expectedFields {
			if !bytes.Contains(body, []byte(field)) {
				t.Errorf("Expected field %s in response, body: %s", field, bodyStr)
			}
		}
	})
}

// Helper function to create a lift max for testing
func createMax(t *testing.T, ts *testutil.TestServer, userID, liftID, maxType string, value float64, effectiveDate *time.Time) LiftMaxResponse {
	t.Helper()

	body := fmt.Sprintf(`{"liftId": "%s", "type": "%s", "value": %f`, liftID, maxType, value)
	if effectiveDate != nil {
		body += fmt.Sprintf(`, "effectiveDate": "%s"`, effectiveDate.Format(time.RFC3339))
	}
	body += "}"

	resp, err := http.Post(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userID)), "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("Failed to create lift max: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create lift max: status %d, body: %s", resp.StatusCode, bodyBytes)
	}

	// Check if response has warnings wrapper
	bodyBytes, _ := io.ReadAll(resp.Body)

	// Try to decode as response with warnings first
	var withWarnings LiftMaxWithWarningsResponse
	if err := json.Unmarshal(bodyBytes, &withWarnings); err == nil && withWarnings.Data.ID != "" {
		return withWarnings.Data
	}

	// Otherwise decode as plain response
	var max LiftMaxResponse
	if err := json.Unmarshal(bodyBytes, &max); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	return max
}
