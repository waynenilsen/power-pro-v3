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

// LiftData matches the lift data format within the response envelope.
type LiftData struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Slug              string  `json:"slug"`
	IsCompetitionLift bool    `json:"isCompetitionLift"`
	ParentLiftID      *string `json:"parentLiftId"`
	CreatedAt         string  `json:"createdAt"`
	UpdatedAt         string  `json:"updatedAt"`
}

// LiftResponse is the standard envelope for single lift responses.
type LiftResponse struct {
	Data LiftData `json:"data"`
}

// PaginationMeta contains pagination metadata.
type PaginationMeta struct {
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasMore bool  `json:"hasMore"`
}

// PaginatedLiftsResponse is the paginated list response with standard envelope.
type PaginatedLiftsResponse struct {
	Data []LiftData      `json:"data"`
	Meta *PaginationMeta `json:"meta"`
}

// ErrorDetail contains structured error information.
type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ErrorResponse is the error response format.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// authGet performs an authenticated GET request
func authGet(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestUserID)
	return http.DefaultClient.Do(req)
}

// adminGet performs an admin-authenticated GET request
func adminGet(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// adminPost performs an admin-authenticated POST request
func adminPost(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// adminPut performs an admin-authenticated PUT request
func adminPut(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// adminDelete performs an admin-authenticated DELETE request
func adminDelete(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

func TestListLifts(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("returns seeded lifts", func(t *testing.T) {
		resp, err := authGet(ts.URL("/lifts"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}

		var result PaginatedLiftsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Should have 5 seeded lifts (squat, bench, deadlift, overhead press, power clean)
		if len(result.Data) != 5 {
			t.Errorf("Expected 5 lifts, got %d", len(result.Data))
		}
		if result.Meta == nil {
			t.Fatal("Expected meta to be present")
		}
		if result.Meta.Total != 5 {
			t.Errorf("Expected total 5, got %d", result.Meta.Total)
		}
	})

	t.Run("supports pagination", func(t *testing.T) {
		resp, err := authGet(ts.URL("/lifts?limit=2&offset=0"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedLiftsResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) != 2 {
			t.Errorf("Expected 2 lifts on page 1, got %d", len(result.Data))
		}
		if result.Meta == nil {
			t.Fatal("Expected meta to be present")
		}
		if result.Meta.Offset != 0 {
			t.Errorf("Expected offset 0, got %d", result.Meta.Offset)
		}
		if result.Meta.Limit != 2 {
			t.Errorf("Expected limit 2, got %d", result.Meta.Limit)
		}
		if !result.Meta.HasMore {
			t.Error("Expected hasMore to be true for page 1")
		}

		// Get page 3 (offset=4) - should have 1 lift with 5 total
		resp2, _ := authGet(ts.URL("/lifts?limit=2&offset=4"))
		defer resp2.Body.Close()

		var result2 PaginatedLiftsResponse
		json.NewDecoder(resp2.Body).Decode(&result2)

		if len(result2.Data) != 1 {
			t.Errorf("Expected 1 lift on page 3 (offset=4), got %d", len(result2.Data))
		}
		if result2.Meta.HasMore {
			t.Error("Expected hasMore to be false for last page")
		}
	})

	t.Run("filters by is_competition_lift", func(t *testing.T) {
		resp, err := authGet(ts.URL("/lifts?is_competition_lift=true"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedLiftsResponse
		json.NewDecoder(resp.Body).Decode(&result)

		// Only 3 lifts are competition lifts (squat, bench, deadlift)
		// Overhead press and power clean are not competition lifts
		if len(result.Data) != 3 {
			t.Errorf("Expected 3 competition lifts, got %d", len(result.Data))
		}

		for _, lift := range result.Data {
			if !lift.IsCompetitionLift {
				t.Errorf("Expected all lifts to be competition lifts")
			}
		}

		// Filter for non-competition lifts
		resp2, _ := authGet(ts.URL("/lifts?is_competition_lift=false"))
		defer resp2.Body.Close()

		var result2 PaginatedLiftsResponse
		json.NewDecoder(resp2.Body).Decode(&result2)

		// 2 non-competition lifts (overhead press, power clean)
		if len(result2.Data) != 2 {
			t.Errorf("Expected 2 non-competition lifts, got %d", len(result2.Data))
		}
	})

	t.Run("supports sorting by name", func(t *testing.T) {
		// Ascending (default)
		resp, _ := authGet(ts.URL("/lifts?sortBy=name&sortOrder=asc"))
		defer resp.Body.Close()

		var result PaginatedLiftsResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) < 2 {
			t.Fatal("Need at least 2 lifts for sort test")
		}

		// First should be "Bench Press" (alphabetically first)
		if result.Data[0].Name != "Bench Press" {
			t.Errorf("Expected first lift to be 'Bench Press', got %s", result.Data[0].Name)
		}

		// Descending
		resp2, _ := authGet(ts.URL("/lifts?sortBy=name&sortOrder=desc"))
		defer resp2.Body.Close()

		var result2 PaginatedLiftsResponse
		json.NewDecoder(resp2.Body).Decode(&result2)

		// First should be "Squat" (alphabetically last)
		if result2.Data[0].Name != "Squat" {
			t.Errorf("Expected first lift to be 'Squat', got %s", result2.Data[0].Name)
		}
	})
}

func TestGetLift(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("returns lift by ID", func(t *testing.T) {
		// Use seeded squat ID
		squatID := "00000000-0000-0000-0000-000000000001"

		resp, err := authGet(ts.URL("/lifts/" + squatID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var result LiftResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if result.Data.ID != squatID {
			t.Errorf("Expected ID %s, got %s", squatID, result.Data.ID)
		}
		if result.Data.Name != "Squat" {
			t.Errorf("Expected name 'Squat', got %s", result.Data.Name)
		}
		if result.Data.Slug != "squat" {
			t.Errorf("Expected slug 'squat', got %s", result.Data.Slug)
		}
		if !result.Data.IsCompetitionLift {
			t.Errorf("Expected isCompetitionLift to be true")
		}
	})

	t.Run("returns 404 for non-existent ID", func(t *testing.T) {
		resp, _ := authGet(ts.URL("/lifts/non-existent-id"))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)

		if !strings.Contains(errResp.Error.Message, "lift not found") {
			t.Errorf("Expected error to contain 'lift not found', got %s", errResp.Error.Message)
		}
	})
}

func TestGetLiftBySlug(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("returns lift by slug", func(t *testing.T) {
		resp, err := authGet(ts.URL("/lifts/by-slug/bench-press"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var result LiftResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if result.Data.Slug != "bench-press" {
			t.Errorf("Expected slug 'bench-press', got %s", result.Data.Slug)
		}
		if result.Data.Name != "Bench Press" {
			t.Errorf("Expected name 'Bench Press', got %s", result.Data.Name)
		}
	})

	t.Run("returns 404 for non-existent slug", func(t *testing.T) {
		resp, _ := authGet(ts.URL("/lifts/by-slug/non-existent"))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestCreateLift(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("creates lift with all fields", func(t *testing.T) {
		body := `{"name": "Pause Squat", "slug": "pause-squat", "isCompetitionLift": false}`
		resp, err := adminPost(ts.URL("/lifts"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, body)
		}

		var result LiftResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if result.Data.Name != "Pause Squat" {
			t.Errorf("Expected name 'Pause Squat', got %s", result.Data.Name)
		}
		if result.Data.Slug != "pause-squat" {
			t.Errorf("Expected slug 'pause-squat', got %s", result.Data.Slug)
		}
		if result.Data.IsCompetitionLift {
			t.Errorf("Expected isCompetitionLift to be false")
		}
		if result.Data.ID == "" {
			t.Errorf("Expected ID to be generated")
		}
	})

	t.Run("auto-generates slug from name", func(t *testing.T) {
		body := `{"name": "Front Squat"}`
		resp, _ := adminPost(ts.URL("/lifts"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result LiftResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if result.Data.Slug != "front-squat" {
			t.Errorf("Expected auto-generated slug 'front-squat', got %s", result.Data.Slug)
		}
	})

	t.Run("creates lift with parent", func(t *testing.T) {
		squatID := "00000000-0000-0000-0000-000000000001"
		body := fmt.Sprintf(`{"name": "Box Squat", "parentLiftId": "%s"}`, squatID)
		resp, _ := adminPost(ts.URL("/lifts"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result LiftResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if result.Data.ParentLiftID == nil || *result.Data.ParentLiftID != squatID {
			t.Errorf("Expected parentLiftId %s, got %v", squatID, result.Data.ParentLiftID)
		}
	})

	t.Run("returns 400 for missing name", func(t *testing.T) {
		body := `{"slug": "no-name"}`
		resp, _ := adminPost(ts.URL("/lifts"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid slug format", func(t *testing.T) {
		body := `{"name": "Invalid Slug", "slug": "INVALID_SLUG!"}`
		resp, _ := adminPost(ts.URL("/lifts"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 409 for duplicate slug", func(t *testing.T) {
		body := `{"name": "Another Squat", "slug": "squat"}`
		resp, _ := adminPost(ts.URL("/lifts"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for non-existent parent lift", func(t *testing.T) {
		body := `{"name": "Orphan Lift", "parentLiftId": "non-existent-id"}`
		resp, _ := adminPost(ts.URL("/lifts"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestUpdateLift(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a lift to update
	createBody := `{"name": "Test Lift", "slug": "test-lift"}`
	createResp, _ := adminPost(ts.URL("/lifts"), createBody)
	var createdResult LiftResponse
	json.NewDecoder(createResp.Body).Decode(&createdResult)
	createResp.Body.Close()

	t.Run("updates lift name", func(t *testing.T) {
		body := `{"name": "Updated Lift"}`
		resp, err := adminPut(ts.URL("/lifts/"+createdResult.Data.ID), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result LiftResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if result.Data.Name != "Updated Lift" {
			t.Errorf("Expected name 'Updated Lift', got %s", result.Data.Name)
		}
		// Slug should remain unchanged
		if result.Data.Slug != "test-lift" {
			t.Errorf("Expected slug 'test-lift', got %s", result.Data.Slug)
		}
	})

	t.Run("updates lift slug", func(t *testing.T) {
		body := `{"slug": "updated-slug"}`
		resp, _ := adminPut(ts.URL("/lifts/"+createdResult.Data.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var result LiftResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if result.Data.Slug != "updated-slug" {
			t.Errorf("Expected slug 'updated-slug', got %s", result.Data.Slug)
		}
	})

	t.Run("returns 404 for non-existent lift", func(t *testing.T) {
		body := `{"name": "Updated"}`
		resp, _ := adminPut(ts.URL("/lifts/non-existent-id"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 409 for duplicate slug", func(t *testing.T) {
		body := `{"slug": "squat"}`
		resp, _ := adminPut(ts.URL("/lifts/"+createdResult.Data.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for validation errors", func(t *testing.T) {
		body := `{"name": ""}`
		resp, _ := adminPut(ts.URL("/lifts/"+createdResult.Data.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestDeleteLift(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("deletes lift successfully", func(t *testing.T) {
		// Create a lift to delete
		createBody := `{"name": "To Delete", "slug": "to-delete"}`
		createResp, _ := adminPost(ts.URL("/lifts"), createBody)
		var createdResult LiftResponse
		json.NewDecoder(createResp.Body).Decode(&createdResult)
		createResp.Body.Close()

		// Delete it
		resp, err := adminDelete(ts.URL("/lifts/" + createdResult.Data.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 204, got %d: %s", resp.StatusCode, bodyBytes)
		}

		// Verify it's deleted
		getResp, _ := authGet(ts.URL("/lifts/" + createdResult.Data.ID))
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected deleted lift to return 404, got %d", getResp.StatusCode)
		}
	})

	t.Run("returns 404 for non-existent lift", func(t *testing.T) {
		resp, _ := adminDelete(ts.URL("/lifts/non-existent-id"))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 409 when lift has child references", func(t *testing.T) {
		// Create a parent lift
		createParent := `{"name": "Parent Lift", "slug": "parent-lift"}`
		parentResp, _ := adminPost(ts.URL("/lifts"), createParent)
		var parentResult LiftResponse
		json.NewDecoder(parentResp.Body).Decode(&parentResult)
		parentResp.Body.Close()

		// Create a child lift referencing the parent
		createChild := fmt.Sprintf(`{"name": "Child Lift", "slug": "child-lift", "parentLiftId": "%s"}`, parentResult.Data.ID)
		childResp, _ := adminPost(ts.URL("/lifts"), createChild)
		childResp.Body.Close()

		// Try to delete the parent
		resp, _ := adminDelete(ts.URL("/lifts/" + parentResult.Data.ID))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409 when deleting lift with children, got %d", resp.StatusCode)
		}
	})
}

func TestSelfReferenceRejection(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a lift
	createBody := `{"name": "Self Ref Test", "slug": "self-ref-test"}`
	createResp, _ := adminPost(ts.URL("/lifts"), createBody)
	var createdResult LiftResponse
	json.NewDecoder(createResp.Body).Decode(&createdResult)
	createResp.Body.Close()

	t.Run("rejects self-reference on update", func(t *testing.T) {
		body := fmt.Sprintf(`{"parentLiftId": "%s"}`, createdResult.Data.ID)
		resp, _ := adminPut(ts.URL("/lifts/"+createdResult.Data.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 for self-reference, got %d", resp.StatusCode)
		}
	})
}

func TestCircularReferenceDetection(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create lift A
	createA := `{"name": "Lift A", "slug": "lift-a"}`
	respA, _ := adminPost(ts.URL("/lifts"), createA)
	var liftAResult LiftResponse
	json.NewDecoder(respA.Body).Decode(&liftAResult)
	respA.Body.Close()

	// Create lift B with parent A
	createB := fmt.Sprintf(`{"name": "Lift B", "slug": "lift-b", "parentLiftId": "%s"}`, liftAResult.Data.ID)
	respB, _ := adminPost(ts.URL("/lifts"), createB)
	var liftBResult LiftResponse
	json.NewDecoder(respB.Body).Decode(&liftBResult)
	respB.Body.Close()

	// Create lift C with parent B
	createC := fmt.Sprintf(`{"name": "Lift C", "slug": "lift-c", "parentLiftId": "%s"}`, liftBResult.Data.ID)
	respC, _ := adminPost(ts.URL("/lifts"), createC)
	var liftCResult LiftResponse
	json.NewDecoder(respC.Body).Decode(&liftCResult)
	respC.Body.Close()

	t.Run("rejects circular reference A->C", func(t *testing.T) {
		// Try to set A's parent to C (would create A->C->B->A cycle)
		body := fmt.Sprintf(`{"parentLiftId": "%s"}`, liftCResult.Data.ID)
		resp, _ := adminPut(ts.URL("/lifts/"+liftAResult.Data.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 for circular reference, got %d", resp.StatusCode)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)

		// The circular reference error is in the details, check various formats
		found := false
		switch details := errResp.Error.Details.(type) {
		case []interface{}:
			for _, detail := range details {
				if detailStr, ok := detail.(string); ok && strings.Contains(detailStr, "circular reference") {
					found = true
					break
				}
			}
		case map[string]interface{}:
			if validationErrors, ok := details["validationErrors"].([]interface{}); ok {
				for _, err := range validationErrors {
					if errStr, ok := err.(string); ok && strings.Contains(errStr, "circular reference") {
						found = true
						break
					}
				}
			}
		}
		if !found && !strings.Contains(errResp.Error.Message, "circular reference") {
			t.Errorf("Expected circular reference error message, got %v (details: %v)", errResp.Error.Message, errResp.Error.Details)
		}
	})
}

func TestResponseFormat(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("response has correct JSON field names", func(t *testing.T) {
		resp, _ := authGet(ts.URL("/lifts/00000000-0000-0000-0000-000000000001"))
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// Check camelCase field names per ERD spec
		expectedFields := []string{
			`"id"`,
			`"name"`,
			`"slug"`,
			`"isCompetitionLift"`,
			`"parentLiftId"`,
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
