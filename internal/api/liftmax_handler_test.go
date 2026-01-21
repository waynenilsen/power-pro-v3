package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// LiftMaxData matches the lift max data within a response envelope.
type LiftMaxData struct {
	ID            string    `json:"id"`
	UserID        string    `json:"userId"`
	LiftID        string    `json:"liftId"`
	Type          string    `json:"type"`
	Value         float64   `json:"value"`
	EffectiveDate time.Time `json:"effectiveDate"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// LiftMaxResponse is the standard envelope for single lift max responses.
type LiftMaxResponse struct {
	Data     LiftMaxData `json:"data"`
	Warnings []string    `json:"warnings,omitempty"`
}

// LiftMaxPaginationMeta contains pagination metadata.
type LiftMaxPaginationMeta struct {
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasMore bool  `json:"hasMore"`
}

// PaginatedLiftMaxesResponse is the paginated list response with standard envelope.
type PaginatedLiftMaxesResponse struct {
	Data []LiftMaxData          `json:"data"`
	Meta *LiftMaxPaginationMeta `json:"meta"`
}

// LiftMaxWithWarningsResponse wraps response with warnings.
// Deprecated: use LiftMaxResponse which already includes warnings
type LiftMaxWithWarningsResponse struct {
	Data     LiftMaxData `json:"data"`
	Warnings []string    `json:"warnings,omitempty"`
}

// Test constants
const (
	testSquatID = "00000000-0000-0000-0000-000000000001"
	testBenchID = "00000000-0000-0000-0000-000000000002"
)

func TestListLiftMaxes(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	userID := testutil.TestUserID

	// Create some test maxes first
	createMax(t, ts, userID, testSquatID, "ONE_RM", 315.0, nil)
	createMax(t, ts, userID, testSquatID, "TRAINING_MAX", 285.0, nil)
	createMax(t, ts, userID, testBenchID, "ONE_RM", 225.0, nil)

	t.Run("returns user's lift maxes", func(t *testing.T) {
		resp, err := authGetUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userID)), userID)
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
		if result.Meta == nil || result.Meta.Total != 3 {
			total := int64(0)
			if result.Meta != nil {
				total = result.Meta.Total
			}
			t.Errorf("Expected total 3, got %d", total)
		}
	})

	t.Run("supports pagination", func(t *testing.T) {
		resp, err := authGetUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes?limit=2&offset=0", userID)), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedLiftMaxesResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) != 2 {
			t.Errorf("Expected 2 lift maxes on page 1, got %d", len(result.Data))
		}
		if result.Meta == nil || result.Meta.Limit != 2 {
			limit := 0
			if result.Meta != nil {
				limit = result.Meta.Limit
			}
			t.Errorf("Expected limit 2, got %d", limit)
		}
		if result.Meta == nil || result.Meta.HasMore != true {
			t.Errorf("Expected hasMore to be true")
		}

		// Get page 2 (offset=2)
		resp2, _ := authGetUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes?limit=2&offset=2", userID)), userID)
		defer resp2.Body.Close()

		var result2 PaginatedLiftMaxesResponse
		json.NewDecoder(resp2.Body).Decode(&result2)

		if len(result2.Data) != 1 {
			t.Errorf("Expected 1 lift max on page 2, got %d", len(result2.Data))
		}
	})

	t.Run("filters by lift_id", func(t *testing.T) {
		resp, err := authGetUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes?lift_id=%s", userID, testSquatID)), userID)
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
		resp, err := authGetUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes?type=ONE_RM", userID)), userID)
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
		resp, err := authGetUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes?lift_id=%s&type=ONE_RM", userID, testSquatID)), userID)
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
		// Use admin to access non-existent user's list (which should be empty)
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/users/non-existent-user/lift-maxes"), nil)
		req.Header.Set("X-User-ID", testutil.TestAdminID)
		req.Header.Set("X-Admin", "true")
		resp, err := http.DefaultClient.Do(req)
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
		resp, err := authGetUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes?type=INVALID", userID)), userID)
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

		resp, err := authGetUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", newUserID)), newUserID)
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

	userID := testutil.TestUserID

	// Create a lift max
	created := createMax(t, ts, userID, testSquatID, "ONE_RM", 315.0, nil)

	t.Run("returns lift max by ID", func(t *testing.T) {
		resp, err := authGetUser(ts.URL("/lift-maxes/"+created.Data.ID), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var max LiftMaxResponse
		json.NewDecoder(resp.Body).Decode(&max)

		if max.Data.ID != created.Data.ID {
			t.Errorf("Expected ID %s, got %s", created.Data.ID, max.Data.ID)
		}
		if max.Data.UserID != userID {
			t.Errorf("Expected userId %s, got %s", userID, max.Data.UserID)
		}
		if max.Data.LiftID != testSquatID {
			t.Errorf("Expected liftId %s, got %s", testSquatID, max.Data.LiftID)
		}
		if max.Data.Type != "ONE_RM" {
			t.Errorf("Expected type ONE_RM, got %s", max.Data.Type)
		}
		if max.Data.Value != 315.0 {
			t.Errorf("Expected value 315.0, got %f", max.Data.Value)
		}
	})

	t.Run("returns 404 for non-existent ID", func(t *testing.T) {
		resp, _ := authGetUser(ts.URL("/lift-maxes/non-existent-id"), userID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)

		if !strings.Contains(errResp.Error.Message, "lift max not found") {
			t.Errorf("Expected error to contain 'lift max not found', got %s", errResp.Error.Message)
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
		userID := "create-test-user"
		body := fmt.Sprintf(`{"liftId": "%s", "type": "ONE_RM", "value": 405.0}`, testSquatID)
		resp, err := authPostUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userID)), body, userID)
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

		if max.Data.LiftID != testSquatID {
			t.Errorf("Expected liftId %s, got %s", testSquatID, max.Data.LiftID)
		}
		if max.Data.Type != "ONE_RM" {
			t.Errorf("Expected type ONE_RM, got %s", max.Data.Type)
		}
		if max.Data.Value != 405.0 {
			t.Errorf("Expected value 405.0, got %f", max.Data.Value)
		}
		if max.Data.ID == "" {
			t.Errorf("Expected ID to be generated")
		}
	})

	t.Run("creates lift max with custom effective date", func(t *testing.T) {
		userID := "date-test-user"
		effectiveDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
		body := fmt.Sprintf(`{"liftId": "%s", "type": "ONE_RM", "value": 410.0, "effectiveDate": "%s"}`, testSquatID, effectiveDate.Format(time.RFC3339))
		resp, _ := authPostUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userID)), body, userID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var max LiftMaxResponse
		json.NewDecoder(resp.Body).Decode(&max)

		if max.Data.EffectiveDate.Year() != 2024 || max.Data.EffectiveDate.Month() != 1 || max.Data.EffectiveDate.Day() != 15 {
			t.Errorf("Expected effectiveDate 2024-01-15, got %s", max.Data.EffectiveDate)
		}
	})

	t.Run("returns 400 for missing liftId", func(t *testing.T) {
		userID := "missing-lift-user"
		body := `{"type": "ONE_RM", "value": 315.0}`
		resp, _ := authPostUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userID)), body, userID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for missing type", func(t *testing.T) {
		userID := "missing-type-user"
		body := fmt.Sprintf(`{"liftId": "%s", "value": 315.0}`, testSquatID)
		resp, _ := authPostUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userID)), body, userID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid type", func(t *testing.T) {
		userID := "invalid-type-user"
		body := fmt.Sprintf(`{"liftId": "%s", "type": "INVALID", "value": 315.0}`, testSquatID)
		resp, _ := authPostUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userID)), body, userID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid value precision", func(t *testing.T) {
		userID := "invalid-precision-user"
		body := fmt.Sprintf(`{"liftId": "%s", "type": "ONE_RM", "value": 315.33}`, testSquatID)
		resp, _ := authPostUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userID)), body, userID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for non-positive value", func(t *testing.T) {
		userID := "zero-value-user"
		body := fmt.Sprintf(`{"liftId": "%s", "type": "ONE_RM", "value": 0}`, testSquatID)
		resp, _ := authPostUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userID)), body, userID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for non-existent lift", func(t *testing.T) {
		userID := "bad-lift-user"
		body := `{"liftId": "non-existent-lift", "type": "ONE_RM", "value": 315.0}`
		resp, _ := authPostUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userID)), body, userID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 409 for duplicate unique constraint", func(t *testing.T) {
		userID := "duplicate-test-user"
		// Create first max
		effectiveDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
		body := fmt.Sprintf(`{"liftId": "%s", "type": "ONE_RM", "value": 315.0, "effectiveDate": "%s"}`, testSquatID, effectiveDate.Format(time.RFC3339))
		resp1, _ := authPostUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userID)), body, userID)
		resp1.Body.Close()

		// Try to create duplicate
		resp2, _ := authPostUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userID)), body, userID)
		defer resp2.Body.Close()

		if resp2.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409, got %d", resp2.StatusCode)
		}
	})

	t.Run("includes warning for TM outside expected range", func(t *testing.T) {
		userID := "tm-warning-user"
		// First create a 1RM
		body1 := fmt.Sprintf(`{"liftId": "%s", "type": "ONE_RM", "value": 400.0}`, testSquatID)
		resp1, _ := authPostUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userID)), body1, userID)
		resp1.Body.Close()

		// Now create a TM that's too low (70% = 280)
		body2 := fmt.Sprintf(`{"liftId": "%s", "type": "TRAINING_MAX", "value": 280.0}`, testSquatID)
		resp2, _ := authPostUser(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userID)), body2, userID)
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
	userID := "update-user"
	created := createMax(t, ts, userID, testSquatID, "ONE_RM", 315.0, nil)

	t.Run("updates lift max value", func(t *testing.T) {
		body := `{"value": 325.0}`
		resp, err := authPutUser(ts.URL("/lift-maxes/"+created.Data.ID), body, userID)
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

		if max.Data.Value != 325.0 {
			t.Errorf("Expected value 325.0, got %f", max.Data.Value)
		}
		// Type should remain unchanged
		if max.Data.Type != "ONE_RM" {
			t.Errorf("Expected type to remain ONE_RM, got %s", max.Data.Type)
		}
	})

	t.Run("updates lift max effective date", func(t *testing.T) {
		newDate := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
		body := fmt.Sprintf(`{"effectiveDate": "%s"}`, newDate.Format(time.RFC3339))
		resp, _ := authPutUser(ts.URL("/lift-maxes/"+created.Data.ID), body, userID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var max LiftMaxResponse
		json.NewDecoder(resp.Body).Decode(&max)

		if max.Data.EffectiveDate.Year() != 2024 || max.Data.EffectiveDate.Month() != 7 {
			t.Errorf("Expected effective date 2024-07, got %s", max.Data.EffectiveDate)
		}
	})

	t.Run("returns 404 for non-existent lift max", func(t *testing.T) {
		body := `{"value": 400.0}`
		resp, _ := authPutUser(ts.URL("/lift-maxes/non-existent-id"), body, userID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid value precision", func(t *testing.T) {
		body := `{"value": 325.33}`
		resp, _ := authPutUser(ts.URL("/lift-maxes/"+created.Data.ID), body, userID)
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
		userID := "delete-user"
		// Create a lift max to delete
		created := createMax(t, ts, userID, testSquatID, "ONE_RM", 315.0, nil)

		// Delete it
		resp, err := authDeleteUser(ts.URL("/lift-maxes/"+created.Data.ID), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 204, got %d: %s", resp.StatusCode, bodyBytes)
		}

		// Verify it's deleted
		getResp, _ := authGetUser(ts.URL("/lift-maxes/"+created.Data.ID), userID)
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected deleted lift max to return 404, got %d", getResp.StatusCode)
		}
	})

	t.Run("returns 404 for non-existent lift max", func(t *testing.T) {
		userID := testutil.TestUserID
		resp, _ := authDeleteUser(ts.URL("/lift-maxes/non-existent-id"), userID)
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
	userID := "format-user"
	created := createMax(t, ts, userID, testSquatID, "ONE_RM", 315.0, nil)

	t.Run("response has correct JSON field names", func(t *testing.T) {
		resp, _ := authGetUser(ts.URL("/lift-maxes/"+created.Data.ID), userID)
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

// ConversionResponse matches the API response format for conversions.
type ConversionResponse struct {
	OriginalValue  float64 `json:"originalValue"`
	OriginalType   string  `json:"originalType"`
	ConvertedValue float64 `json:"convertedValue"`
	ConvertedType  string  `json:"convertedType"`
	Percentage     float64 `json:"percentage"`
}

func TestConvertLiftMax(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create test data
	oneRMUserID := "convert-test-user"
	oneRM := createMax(t, ts, oneRMUserID, testSquatID, "ONE_RM", 400.0, nil)

	tmUserID := "convert-tm-user"
	tm := createMax(t, ts, tmUserID, testSquatID, "TRAINING_MAX", 360.0, nil)

	t.Run("converts 1RM to Training Max with default percentage", func(t *testing.T) {
		resp, err := authGetUser(ts.URL("/lift-maxes/"+oneRM.Data.ID+"/convert?to_type=TRAINING_MAX"), oneRMUserID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}

		var result ConversionResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.OriginalValue != 400.0 {
			t.Errorf("Expected originalValue 400.0, got %f", result.OriginalValue)
		}
		if result.OriginalType != "ONE_RM" {
			t.Errorf("Expected originalType ONE_RM, got %s", result.OriginalType)
		}
		// 400 * 0.90 = 360
		if result.ConvertedValue != 360.0 {
			t.Errorf("Expected convertedValue 360.0, got %f", result.ConvertedValue)
		}
		if result.ConvertedType != "TRAINING_MAX" {
			t.Errorf("Expected convertedType TRAINING_MAX, got %s", result.ConvertedType)
		}
		if result.Percentage != 90.0 {
			t.Errorf("Expected percentage 90.0, got %f", result.Percentage)
		}
	})

	t.Run("converts Training Max to 1RM with default percentage", func(t *testing.T) {
		resp, err := authGetUser(ts.URL("/lift-maxes/"+tm.Data.ID+"/convert?to_type=ONE_RM"), tmUserID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}

		var result ConversionResponse
		json.NewDecoder(resp.Body).Decode(&result)

		// 360 / 0.90 = 400
		if result.ConvertedValue != 400.0 {
			t.Errorf("Expected convertedValue 400.0, got %f", result.ConvertedValue)
		}
		if result.ConvertedType != "ONE_RM" {
			t.Errorf("Expected convertedType ONE_RM, got %s", result.ConvertedType)
		}
	})

	t.Run("converts with custom percentage", func(t *testing.T) {
		resp, _ := authGetUser(ts.URL("/lift-maxes/"+oneRM.Data.ID+"/convert?to_type=TRAINING_MAX&percentage=85"), oneRMUserID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}

		var result ConversionResponse
		json.NewDecoder(resp.Body).Decode(&result)

		// 400 * 0.85 = 340
		if result.ConvertedValue != 340.0 {
			t.Errorf("Expected convertedValue 340.0, got %f", result.ConvertedValue)
		}
		if result.Percentage != 85.0 {
			t.Errorf("Expected percentage 85.0, got %f", result.Percentage)
		}
	})

	t.Run("rounds converted value to nearest 0.25", func(t *testing.T) {
		// Create a max that will produce a non-quarter result
		oddUserID := "round-test-user"
		oddMax := createMax(t, ts, oddUserID, testSquatID, "ONE_RM", 315.0, nil)

		resp, _ := authGetUser(ts.URL("/lift-maxes/"+oddMax.Data.ID+"/convert?to_type=TRAINING_MAX&percentage=87"), oddUserID)
		defer resp.Body.Close()

		var result ConversionResponse
		json.NewDecoder(resp.Body).Decode(&result)

		// 315 * 0.87 = 274.05 -> rounds to 274.0
		// Verify it's a multiple of 0.25
		scaled := result.ConvertedValue * 4
		if scaled != float64(int(scaled)) {
			t.Errorf("Expected converted value to be multiple of 0.25, got %f", result.ConvertedValue)
		}
	})

	t.Run("returns 400 for missing to_type parameter", func(t *testing.T) {
		resp, _ := authGetUser(ts.URL("/lift-maxes/"+oneRM.Data.ID+"/convert"), oneRMUserID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		if errResp.Error.Message != "Missing required query parameter: to_type" {
			t.Errorf("Expected error about missing to_type, got: %s", errResp.Error.Message)
		}
	})

	t.Run("returns 400 for invalid to_type", func(t *testing.T) {
		resp, _ := authGetUser(ts.URL("/lift-maxes/"+oneRM.Data.ID+"/convert?to_type=INVALID"), oneRMUserID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 when converting to same type", func(t *testing.T) {
		resp, _ := authGetUser(ts.URL("/lift-maxes/"+oneRM.Data.ID+"/convert?to_type=ONE_RM"), oneRMUserID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		if errResp.Error.Message != "Cannot convert to same type: lift max is already ONE_RM" {
			t.Errorf("Expected error about same type, got: %s", errResp.Error.Message)
		}
	})

	t.Run("returns 400 for percentage below 1", func(t *testing.T) {
		resp, _ := authGetUser(ts.URL("/lift-maxes/"+oneRM.Data.ID+"/convert?to_type=TRAINING_MAX&percentage=0"), oneRMUserID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for percentage above 100", func(t *testing.T) {
		resp, _ := authGetUser(ts.URL("/lift-maxes/"+oneRM.Data.ID+"/convert?to_type=TRAINING_MAX&percentage=101"), oneRMUserID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for non-numeric percentage", func(t *testing.T) {
		resp, _ := authGetUser(ts.URL("/lift-maxes/"+oneRM.Data.ID+"/convert?to_type=TRAINING_MAX&percentage=abc"), oneRMUserID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 404 for non-existent lift max", func(t *testing.T) {
		resp, _ := authGetUser(ts.URL("/lift-maxes/non-existent-id/convert?to_type=TRAINING_MAX"), "any-user")
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("accepts lowercase to_type", func(t *testing.T) {
		resp, _ := authGetUser(ts.URL("/lift-maxes/"+oneRM.Data.ID+"/convert?to_type=training_max"), oneRMUserID)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200 for lowercase to_type, got %d: %s", resp.StatusCode, body)
		}
	})
}

func TestGetCurrentLiftMax(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testUser := "current-max-user"
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	lastWeek := now.Add(-7 * 24 * time.Hour)

	// Create multiple maxes with different dates to test "most recent" logic
	createMax(t, ts, testUser, testSquatID, "ONE_RM", 300.0, &lastWeek)
	createMax(t, ts, testUser, testSquatID, "ONE_RM", 315.0, &yesterday)
	mostRecent := createMax(t, ts, testUser, testSquatID, "ONE_RM", 320.0, &now)

	// Also create a training max
	createMax(t, ts, testUser, testSquatID, "TRAINING_MAX", 285.0, &yesterday)
	mostRecentTM := createMax(t, ts, testUser, testSquatID, "TRAINING_MAX", 290.0, &now)

	// Create max for a different lift
	createMax(t, ts, testUser, testBenchID, "ONE_RM", 225.0, &now)

	t.Run("returns most recent max for user/lift/type", func(t *testing.T) {
		url := fmt.Sprintf("/users/%s/lift-maxes/current?lift=%s&type=ONE_RM", testUser, testSquatID)
		resp, err := authGetUser(ts.URL(url), testUser)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}

		var max LiftMaxResponse
		if err := json.NewDecoder(resp.Body).Decode(&max); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if max.Data.ID != mostRecent.Data.ID {
			t.Errorf("Expected most recent max ID %s, got %s", mostRecent.Data.ID, max.Data.ID)
		}
		if max.Data.Value != 320.0 {
			t.Errorf("Expected value 320.0, got %f", max.Data.Value)
		}
		if max.Data.Type != "ONE_RM" {
			t.Errorf("Expected type ONE_RM, got %s", max.Data.Type)
		}
	})

	t.Run("returns most recent training max", func(t *testing.T) {
		url := fmt.Sprintf("/users/%s/lift-maxes/current?lift=%s&type=TRAINING_MAX", testUser, testSquatID)
		resp, err := authGetUser(ts.URL(url), testUser)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var max LiftMaxResponse
		json.NewDecoder(resp.Body).Decode(&max)

		if max.Data.ID != mostRecentTM.Data.ID {
			t.Errorf("Expected most recent TM ID %s, got %s", mostRecentTM.Data.ID, max.Data.ID)
		}
		if max.Data.Value != 290.0 {
			t.Errorf("Expected value 290.0, got %f", max.Data.Value)
		}
	})

	t.Run("returns correct max for different lift", func(t *testing.T) {
		url := fmt.Sprintf("/users/%s/lift-maxes/current?lift=%s&type=ONE_RM", testUser, testBenchID)
		resp, _ := authGetUser(ts.URL(url), testUser)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var max LiftMaxResponse
		json.NewDecoder(resp.Body).Decode(&max)

		if max.Data.LiftID != testBenchID {
			t.Errorf("Expected lift %s, got %s", testBenchID, max.Data.LiftID)
		}
		if max.Data.Value != 225.0 {
			t.Errorf("Expected value 225.0, got %f", max.Data.Value)
		}
	})

	t.Run("returns 404 when no max exists for combination", func(t *testing.T) {
		// Query for a lift that doesn't have any maxes
		nonExistentLift := "11111111-1111-1111-1111-111111111111"
		url := fmt.Sprintf("/users/%s/lift-maxes/current?lift=%s&type=ONE_RM", testUser, nonExistentLift)
		resp, _ := authGetUser(ts.URL(url), testUser)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		if errResp.Error.Message != "No lift max found for the specified user, lift, and type" {
			t.Errorf("Expected specific error message, got: %s", errResp.Error.Message)
		}
	})

	t.Run("returns 404 when user has no maxes for this type", func(t *testing.T) {
		// User exists but doesn't have a TRAINING_MAX for bench
		url := fmt.Sprintf("/users/%s/lift-maxes/current?lift=%s&type=TRAINING_MAX", testUser, testBenchID)
		resp, _ := authGetUser(ts.URL(url), testUser)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 when lift param is missing", func(t *testing.T) {
		url := fmt.Sprintf("/users/%s/lift-maxes/current?type=ONE_RM", testUser)
		resp, _ := authGetUser(ts.URL(url), testUser)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		if errResp.Error.Message != "Missing required query parameter: lift" {
			t.Errorf("Expected error about missing lift, got: %s", errResp.Error.Message)
		}
	})

	t.Run("returns 400 when lift param is invalid UUID", func(t *testing.T) {
		url := fmt.Sprintf("/users/%s/lift-maxes/current?lift=not-a-uuid&type=ONE_RM", testUser)
		resp, _ := authGetUser(ts.URL(url), testUser)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		if errResp.Error.Message != "Invalid lift parameter: must be a valid UUID" {
			t.Errorf("Expected error about invalid UUID, got: %s", errResp.Error.Message)
		}
	})

	t.Run("returns 400 when type param is missing", func(t *testing.T) {
		url := fmt.Sprintf("/users/%s/lift-maxes/current?lift=%s", testUser, testSquatID)
		resp, _ := authGetUser(ts.URL(url), testUser)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		if errResp.Error.Message != "Missing required query parameter: type" {
			t.Errorf("Expected error about missing type, got: %s", errResp.Error.Message)
		}
	})

	t.Run("returns 400 when type param is invalid", func(t *testing.T) {
		url := fmt.Sprintf("/users/%s/lift-maxes/current?lift=%s&type=INVALID", testUser, testSquatID)
		resp, _ := authGetUser(ts.URL(url), testUser)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		if errResp.Error.Message != "Invalid type parameter: must be ONE_RM or TRAINING_MAX" {
			t.Errorf("Expected error about invalid type, got: %s", errResp.Error.Message)
		}
	})

	t.Run("accepts lowercase type parameter", func(t *testing.T) {
		url := fmt.Sprintf("/users/%s/lift-maxes/current?lift=%s&type=one_rm", testUser, testSquatID)
		resp, _ := authGetUser(ts.URL(url), testUser)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200 for lowercase type, got %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("works with single max record", func(t *testing.T) {
		singleUser := "single-max-user"
		single := createMax(t, ts, singleUser, testSquatID, "ONE_RM", 405.0, nil)

		url := fmt.Sprintf("/users/%s/lift-maxes/current?lift=%s&type=ONE_RM", singleUser, testSquatID)
		resp, _ := authGetUser(ts.URL(url), singleUser)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var max LiftMaxResponse
		json.NewDecoder(resp.Body).Decode(&max)

		if max.Data.ID != single.Data.ID {
			t.Errorf("Expected max ID %s, got %s", single.Data.ID, max.Data.ID)
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

	req, err := http.NewRequest(http.MethodPost, ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userID)), bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to create lift max: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create lift max: status %d, body: %s", resp.StatusCode, bodyBytes)
	}

	// Decode as standard response envelope
	bodyBytes, _ := io.ReadAll(resp.Body)

	var max LiftMaxResponse
	if err := json.Unmarshal(bodyBytes, &max); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	return max
}

// Helper function to make authenticated GET request for user's own resources
func authGetUser(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

// Helper function to make authenticated PUT request for user's own resources
func authPutUser(url string, body string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

// Helper function to make authenticated DELETE request for user's own resources
func authDeleteUser(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

// Helper function to make authenticated POST request for user's own resources
func authPostUser(url string, body string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}
