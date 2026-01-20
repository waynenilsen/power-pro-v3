package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// LiftResponse matches the API response format.
type LiftResponse struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Slug              string  `json:"slug"`
	IsCompetitionLift bool    `json:"isCompetitionLift"`
	ParentLiftID      *string `json:"parentLiftId"`
	CreatedAt         string  `json:"createdAt"`
	UpdatedAt         string  `json:"updatedAt"`
}

// PaginatedLiftsResponse is the paginated list response.
type PaginatedLiftsResponse struct {
	Data       []LiftResponse `json:"data"`
	Page       int            `json:"page"`
	PageSize   int            `json:"pageSize"`
	TotalItems int64          `json:"totalItems"`
	TotalPages int64          `json:"totalPages"`
}

// ErrorResponse is the error response format.
type ErrorResponse struct {
	Error   string   `json:"error"`
	Details []string `json:"details,omitempty"`
}

func TestListLifts(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("returns seeded lifts", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/lifts"))
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

		// Should have 3 seeded lifts
		if len(result.Data) != 3 {
			t.Errorf("Expected 3 lifts, got %d", len(result.Data))
		}
		if result.TotalItems != 3 {
			t.Errorf("Expected totalItems 3, got %d", result.TotalItems)
		}
	})

	t.Run("supports pagination", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/lifts?page=1&pageSize=2"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedLiftsResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) != 2 {
			t.Errorf("Expected 2 lifts on page 1, got %d", len(result.Data))
		}
		if result.Page != 1 {
			t.Errorf("Expected page 1, got %d", result.Page)
		}
		if result.PageSize != 2 {
			t.Errorf("Expected pageSize 2, got %d", result.PageSize)
		}

		// Get page 2
		resp2, _ := http.Get(ts.URL("/lifts?page=2&pageSize=2"))
		defer resp2.Body.Close()

		var result2 PaginatedLiftsResponse
		json.NewDecoder(resp2.Body).Decode(&result2)

		if len(result2.Data) != 1 {
			t.Errorf("Expected 1 lift on page 2, got %d", len(result2.Data))
		}
	})

	t.Run("filters by is_competition_lift", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/lifts?is_competition_lift=true"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedLiftsResponse
		json.NewDecoder(resp.Body).Decode(&result)

		// All seeded lifts are competition lifts
		if len(result.Data) != 3 {
			t.Errorf("Expected 3 competition lifts, got %d", len(result.Data))
		}

		for _, lift := range result.Data {
			if !lift.IsCompetitionLift {
				t.Errorf("Expected all lifts to be competition lifts")
			}
		}

		// Filter for non-competition lifts
		resp2, _ := http.Get(ts.URL("/lifts?is_competition_lift=false"))
		defer resp2.Body.Close()

		var result2 PaginatedLiftsResponse
		json.NewDecoder(resp2.Body).Decode(&result2)

		if len(result2.Data) != 0 {
			t.Errorf("Expected 0 non-competition lifts, got %d", len(result2.Data))
		}
	})

	t.Run("supports sorting by name", func(t *testing.T) {
		// Ascending (default)
		resp, _ := http.Get(ts.URL("/lifts?sortBy=name&sortOrder=asc"))
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
		resp2, _ := http.Get(ts.URL("/lifts?sortBy=name&sortOrder=desc"))
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

		resp, err := http.Get(ts.URL("/lifts/" + squatID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var lift LiftResponse
		json.NewDecoder(resp.Body).Decode(&lift)

		if lift.ID != squatID {
			t.Errorf("Expected ID %s, got %s", squatID, lift.ID)
		}
		if lift.Name != "Squat" {
			t.Errorf("Expected name 'Squat', got %s", lift.Name)
		}
		if lift.Slug != "squat" {
			t.Errorf("Expected slug 'squat', got %s", lift.Slug)
		}
		if !lift.IsCompetitionLift {
			t.Errorf("Expected isCompetitionLift to be true")
		}
	})

	t.Run("returns 404 for non-existent ID", func(t *testing.T) {
		resp, _ := http.Get(ts.URL("/lifts/non-existent-id"))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)

		if errResp.Error != "Lift not found" {
			t.Errorf("Expected error 'Lift not found', got %s", errResp.Error)
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
		resp, err := http.Get(ts.URL("/lifts/by-slug/bench-press"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var lift LiftResponse
		json.NewDecoder(resp.Body).Decode(&lift)

		if lift.Slug != "bench-press" {
			t.Errorf("Expected slug 'bench-press', got %s", lift.Slug)
		}
		if lift.Name != "Bench Press" {
			t.Errorf("Expected name 'Bench Press', got %s", lift.Name)
		}
	})

	t.Run("returns 404 for non-existent slug", func(t *testing.T) {
		resp, _ := http.Get(ts.URL("/lifts/by-slug/non-existent"))
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
		resp, err := http.Post(ts.URL("/lifts"), "application/json", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, body)
		}

		var lift LiftResponse
		json.NewDecoder(resp.Body).Decode(&lift)

		if lift.Name != "Pause Squat" {
			t.Errorf("Expected name 'Pause Squat', got %s", lift.Name)
		}
		if lift.Slug != "pause-squat" {
			t.Errorf("Expected slug 'pause-squat', got %s", lift.Slug)
		}
		if lift.IsCompetitionLift {
			t.Errorf("Expected isCompetitionLift to be false")
		}
		if lift.ID == "" {
			t.Errorf("Expected ID to be generated")
		}
	})

	t.Run("auto-generates slug from name", func(t *testing.T) {
		body := `{"name": "Front Squat"}`
		resp, _ := http.Post(ts.URL("/lifts"), "application/json", bytes.NewBufferString(body))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var lift LiftResponse
		json.NewDecoder(resp.Body).Decode(&lift)

		if lift.Slug != "front-squat" {
			t.Errorf("Expected auto-generated slug 'front-squat', got %s", lift.Slug)
		}
	})

	t.Run("creates lift with parent", func(t *testing.T) {
		squatID := "00000000-0000-0000-0000-000000000001"
		body := fmt.Sprintf(`{"name": "Box Squat", "parentLiftId": "%s"}`, squatID)
		resp, _ := http.Post(ts.URL("/lifts"), "application/json", bytes.NewBufferString(body))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var lift LiftResponse
		json.NewDecoder(resp.Body).Decode(&lift)

		if lift.ParentLiftID == nil || *lift.ParentLiftID != squatID {
			t.Errorf("Expected parentLiftId %s, got %v", squatID, lift.ParentLiftID)
		}
	})

	t.Run("returns 400 for missing name", func(t *testing.T) {
		body := `{"slug": "no-name"}`
		resp, _ := http.Post(ts.URL("/lifts"), "application/json", bytes.NewBufferString(body))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid slug format", func(t *testing.T) {
		body := `{"name": "Invalid Slug", "slug": "INVALID_SLUG!"}`
		resp, _ := http.Post(ts.URL("/lifts"), "application/json", bytes.NewBufferString(body))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 409 for duplicate slug", func(t *testing.T) {
		body := `{"name": "Another Squat", "slug": "squat"}`
		resp, _ := http.Post(ts.URL("/lifts"), "application/json", bytes.NewBufferString(body))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for non-existent parent lift", func(t *testing.T) {
		body := `{"name": "Orphan Lift", "parentLiftId": "non-existent-id"}`
		resp, _ := http.Post(ts.URL("/lifts"), "application/json", bytes.NewBufferString(body))
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
	createResp, _ := http.Post(ts.URL("/lifts"), "application/json", bytes.NewBufferString(createBody))
	var createdLift LiftResponse
	json.NewDecoder(createResp.Body).Decode(&createdLift)
	createResp.Body.Close()

	t.Run("updates lift name", func(t *testing.T) {
		body := `{"name": "Updated Lift"}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/lifts/"+createdLift.ID), bytes.NewBufferString(body))
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

		var lift LiftResponse
		json.NewDecoder(resp.Body).Decode(&lift)

		if lift.Name != "Updated Lift" {
			t.Errorf("Expected name 'Updated Lift', got %s", lift.Name)
		}
		// Slug should remain unchanged
		if lift.Slug != "test-lift" {
			t.Errorf("Expected slug 'test-lift', got %s", lift.Slug)
		}
	})

	t.Run("updates lift slug", func(t *testing.T) {
		body := `{"slug": "updated-slug"}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/lifts/"+createdLift.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var lift LiftResponse
		json.NewDecoder(resp.Body).Decode(&lift)

		if lift.Slug != "updated-slug" {
			t.Errorf("Expected slug 'updated-slug', got %s", lift.Slug)
		}
	})

	t.Run("returns 404 for non-existent lift", func(t *testing.T) {
		body := `{"name": "Updated"}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/lifts/non-existent-id"), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 409 for duplicate slug", func(t *testing.T) {
		body := `{"slug": "squat"}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/lifts/"+createdLift.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for validation errors", func(t *testing.T) {
		body := `{"name": ""}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/lifts/"+createdLift.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
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
		createResp, _ := http.Post(ts.URL("/lifts"), "application/json", bytes.NewBufferString(createBody))
		var createdLift LiftResponse
		json.NewDecoder(createResp.Body).Decode(&createdLift)
		createResp.Body.Close()

		// Delete it
		req, _ := http.NewRequest(http.MethodDelete, ts.URL("/lifts/"+createdLift.ID), nil)
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
		getResp, _ := http.Get(ts.URL("/lifts/" + createdLift.ID))
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected deleted lift to return 404, got %d", getResp.StatusCode)
		}
	})

	t.Run("returns 404 for non-existent lift", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL("/lifts/non-existent-id"), nil)
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 409 when lift has child references", func(t *testing.T) {
		// Create a parent lift
		createParent := `{"name": "Parent Lift", "slug": "parent-lift"}`
		parentResp, _ := http.Post(ts.URL("/lifts"), "application/json", bytes.NewBufferString(createParent))
		var parentLift LiftResponse
		json.NewDecoder(parentResp.Body).Decode(&parentLift)
		parentResp.Body.Close()

		// Create a child lift referencing the parent
		createChild := fmt.Sprintf(`{"name": "Child Lift", "slug": "child-lift", "parentLiftId": "%s"}`, parentLift.ID)
		childResp, _ := http.Post(ts.URL("/lifts"), "application/json", bytes.NewBufferString(createChild))
		childResp.Body.Close()

		// Try to delete the parent
		req, _ := http.NewRequest(http.MethodDelete, ts.URL("/lifts/"+parentLift.ID), nil)
		resp, _ := http.DefaultClient.Do(req)
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
	createResp, _ := http.Post(ts.URL("/lifts"), "application/json", bytes.NewBufferString(createBody))
	var createdLift LiftResponse
	json.NewDecoder(createResp.Body).Decode(&createdLift)
	createResp.Body.Close()

	t.Run("rejects self-reference on update", func(t *testing.T) {
		body := fmt.Sprintf(`{"parentLiftId": "%s"}`, createdLift.ID)
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/lifts/"+createdLift.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
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
	respA, _ := http.Post(ts.URL("/lifts"), "application/json", bytes.NewBufferString(createA))
	var liftA LiftResponse
	json.NewDecoder(respA.Body).Decode(&liftA)
	respA.Body.Close()

	// Create lift B with parent A
	createB := fmt.Sprintf(`{"name": "Lift B", "slug": "lift-b", "parentLiftId": "%s"}`, liftA.ID)
	respB, _ := http.Post(ts.URL("/lifts"), "application/json", bytes.NewBufferString(createB))
	var liftB LiftResponse
	json.NewDecoder(respB.Body).Decode(&liftB)
	respB.Body.Close()

	// Create lift C with parent B
	createC := fmt.Sprintf(`{"name": "Lift C", "slug": "lift-c", "parentLiftId": "%s"}`, liftB.ID)
	respC, _ := http.Post(ts.URL("/lifts"), "application/json", bytes.NewBufferString(createC))
	var liftC LiftResponse
	json.NewDecoder(respC.Body).Decode(&liftC)
	respC.Body.Close()

	t.Run("rejects circular reference A->C", func(t *testing.T) {
		// Try to set A's parent to C (would create A->C->B->A cycle)
		body := fmt.Sprintf(`{"parentLiftId": "%s"}`, liftC.ID)
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/lifts/"+liftA.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 for circular reference, got %d", resp.StatusCode)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)

		found := false
		for _, detail := range errResp.Details {
			if detail == "circular reference detected: lift cannot be its own ancestor" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected circular reference error in details, got %v", errResp.Details)
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
		resp, _ := http.Get(ts.URL("/lifts/00000000-0000-0000-0000-000000000001"))
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
