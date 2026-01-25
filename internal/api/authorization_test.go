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

// ========================================
// Lift Authorization Tests (NFR-007)
// ========================================

func TestLiftAuthorizationRead(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	squatID := "00000000-0000-0000-0000-000000000001"

	t.Run("unauthenticated user gets 401 on GET /lifts", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/lifts"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 401, got %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("unauthenticated user gets 401 on GET /lifts/{id}", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/lifts/" + squatID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("unauthenticated user gets 401 on GET /lifts/by-slug/{slug}", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/lifts/by-slug/squat"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("authenticated user can GET /lifts", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/lifts"), nil)
		req.Header.Set("X-User-ID", testutil.TestUserID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("authenticated user can GET /lifts/{id}", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/lifts/"+squatID), nil)
		req.Header.Set("X-User-ID", testutil.TestUserID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("authenticated user can GET /lifts/by-slug/{slug}", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/lifts/by-slug/squat"), nil)
		req.Header.Set("X-User-ID", testutil.TestUserID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}
	})
}

func TestLiftAuthorizationWrite(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("unauthenticated user gets 401 on POST /lifts", func(t *testing.T) {
		body := `{"name": "Test Lift"}`
		resp, err := http.Post(ts.URL("/lifts"), "application/json", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("non-admin user gets 403 on POST /lifts", func(t *testing.T) {
		body := `{"name": "Test Lift"}`
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/lifts"), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", testutil.TestUserID)
		// Not setting X-Admin

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("admin user can POST /lifts", func(t *testing.T) {
		body := `{"name": "Admin Created Lift", "slug": "admin-created-lift"}`
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/lifts"), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", testutil.TestAdminID)
		req.Header.Set("X-Admin", "true")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 201, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("unauthenticated user gets 401 on PUT /lifts/{id}", func(t *testing.T) {
		squatID := "00000000-0000-0000-0000-000000000001"
		body := `{"name": "Updated Squat"}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/lifts/"+squatID), bytes.NewBufferString(body))
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

	t.Run("non-admin user gets 403 on PUT /lifts/{id}", func(t *testing.T) {
		squatID := "00000000-0000-0000-0000-000000000001"
		body := `{"name": "Updated Squat"}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/lifts/"+squatID), bytes.NewBufferString(body))
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

	t.Run("admin user can PUT /lifts/{id}", func(t *testing.T) {
		squatID := "00000000-0000-0000-0000-000000000001"
		body := `{"name": "Admin Updated Squat"}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/lifts/"+squatID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", testutil.TestAdminID)
		req.Header.Set("X-Admin", "true")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("unauthenticated user gets 401 on DELETE /lifts/{id}", func(t *testing.T) {
		// First create a lift as admin
		createBody := `{"name": "To Delete 1", "slug": "to-delete-1"}`
		createReq, _ := http.NewRequest(http.MethodPost, ts.URL("/lifts"), bytes.NewBufferString(createBody))
		createReq.Header.Set("Content-Type", "application/json")
		createReq.Header.Set("X-User-ID", testutil.TestAdminID)
		createReq.Header.Set("X-Admin", "true")
		createResp, _ := http.DefaultClient.Do(createReq)
		var createdLift LiftResponse
		json.NewDecoder(createResp.Body).Decode(&createdLift)
		createResp.Body.Close()

		// Try to delete without auth
		req, _ := http.NewRequest(http.MethodDelete, ts.URL("/lifts/"+createdLift.Data.ID), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("non-admin user gets 403 on DELETE /lifts/{id}", func(t *testing.T) {
		// First create a lift as admin
		createBody := `{"name": "To Delete 2", "slug": "to-delete-2"}`
		createReq, _ := http.NewRequest(http.MethodPost, ts.URL("/lifts"), bytes.NewBufferString(createBody))
		createReq.Header.Set("Content-Type", "application/json")
		createReq.Header.Set("X-User-ID", testutil.TestAdminID)
		createReq.Header.Set("X-Admin", "true")
		createResp, _ := http.DefaultClient.Do(createReq)
		var createdLift LiftResponse
		json.NewDecoder(createResp.Body).Decode(&createdLift)
		createResp.Body.Close()

		// Try to delete as non-admin
		req, _ := http.NewRequest(http.MethodDelete, ts.URL("/lifts/"+createdLift.Data.ID), nil)
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

	t.Run("admin user can DELETE /lifts/{id}", func(t *testing.T) {
		// First create a lift as admin
		createBody := `{"name": "To Delete 3", "slug": "to-delete-3"}`
		createReq, _ := http.NewRequest(http.MethodPost, ts.URL("/lifts"), bytes.NewBufferString(createBody))
		createReq.Header.Set("Content-Type", "application/json")
		createReq.Header.Set("X-User-ID", testutil.TestAdminID)
		createReq.Header.Set("X-Admin", "true")
		createResp, _ := http.DefaultClient.Do(createReq)
		var createdLift LiftResponse
		json.NewDecoder(createResp.Body).Decode(&createdLift)
		createResp.Body.Close()

		// Delete as admin
		req, _ := http.NewRequest(http.MethodDelete, ts.URL("/lifts/"+createdLift.Data.ID), nil)
		req.Header.Set("X-User-ID", testutil.TestAdminID)
		req.Header.Set("X-Admin", "true")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 204, got %d: %s", resp.StatusCode, respBody)
		}
	})
}

// ========================================
// LiftMax Authorization Tests (NFR-006)
// ========================================

func TestLiftMaxAuthorizationOwnership(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	userAID := "user-a"
	userBID := "user-b"
	squatID := "00000000-0000-0000-0000-000000000001"

	// Create a lift max for user A
	createMax := func(userID string) string {
		body := fmt.Sprintf(`{"liftId": "%s", "type": "ONE_RM", "value": 315.0}`, squatID)
		req, _ := http.NewRequest(http.MethodPost, ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userID)), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", userID)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to create lift max: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			respBody, _ := io.ReadAll(resp.Body)
			t.Fatalf("Failed to create lift max: %d: %s", resp.StatusCode, respBody)
		}
		var max LiftMaxResponse
		json.NewDecoder(resp.Body).Decode(&max)
		return max.Data.ID
	}

	userAMaxID := createMax(userAID)

	t.Run("unauthenticated user gets 401 on GET /users/{userId}/lift-maxes", func(t *testing.T) {
		resp, err := http.Get(ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userAID)))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("user can access their own lift max list", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userAID)), nil)
		req.Header.Set("X-User-ID", userAID)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("user cannot access another user's lift max list", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userAID)), nil)
		req.Header.Set("X-User-ID", userBID)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("admin can access any user's lift max list", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userAID)), nil)
		req.Header.Set("X-User-ID", testutil.TestAdminID)
		req.Header.Set("X-Admin", "true")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("unauthenticated user gets 401 on GET /lift-maxes/{id}", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/lift-maxes/" + userAMaxID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("owner can access their own lift max by ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/lift-maxes/"+userAMaxID), nil)
		req.Header.Set("X-User-ID", userAID)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("non-owner cannot access another user's lift max by ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/lift-maxes/"+userAMaxID), nil)
		req.Header.Set("X-User-ID", userBID)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("admin can access any user's lift max by ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/lift-maxes/"+userAMaxID), nil)
		req.Header.Set("X-User-ID", testutil.TestAdminID)
		req.Header.Set("X-Admin", "true")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("user can create lift max for themselves", func(t *testing.T) {
		body := fmt.Sprintf(`{"liftId": "%s", "type": "TRAINING_MAX", "value": 285.0}`, squatID)
		req, _ := http.NewRequest(http.MethodPost, ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userBID)), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", userBID)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 201, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("user cannot create lift max for another user", func(t *testing.T) {
		body := fmt.Sprintf(`{"liftId": "%s", "type": "TRAINING_MAX", "value": 300.0}`, squatID)
		req, _ := http.NewRequest(http.MethodPost, ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userAID)), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", userBID)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("admin can create lift max for any user", func(t *testing.T) {
		effectiveDate := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
		body := fmt.Sprintf(`{"liftId": "%s", "type": "ONE_RM", "value": 350.0, "effectiveDate": "%s"}`, squatID, effectiveDate)
		req, _ := http.NewRequest(http.MethodPost, ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userBID)), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", testutil.TestAdminID)
		req.Header.Set("X-Admin", "true")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 201, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("owner can update their own lift max", func(t *testing.T) {
		body := `{"value": 320.0}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/lift-maxes/"+userAMaxID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", userAID)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("non-owner cannot update another user's lift max", func(t *testing.T) {
		body := `{"value": 400.0}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/lift-maxes/"+userAMaxID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", userBID)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("admin can update any user's lift max", func(t *testing.T) {
		body := `{"value": 325.0}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/lift-maxes/"+userAMaxID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", testutil.TestAdminID)
		req.Header.Set("X-Admin", "true")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("non-owner cannot delete another user's lift max", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL("/lift-maxes/"+userAMaxID), nil)
		req.Header.Set("X-User-ID", userBID)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("owner can delete their own lift max", func(t *testing.T) {
		// Create a new max to delete
		body := fmt.Sprintf(`{"liftId": "%s", "type": "ONE_RM", "value": 315.0, "effectiveDate": "2020-01-01T00:00:00Z"}`, squatID)
		createReq, _ := http.NewRequest(http.MethodPost, ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userAID)), bytes.NewBufferString(body))
		createReq.Header.Set("Content-Type", "application/json")
		createReq.Header.Set("X-User-ID", userAID)
		createResp, _ := http.DefaultClient.Do(createReq)
		var newMax LiftMaxResponse
		json.NewDecoder(createResp.Body).Decode(&newMax)
		createResp.Body.Close()

		req, _ := http.NewRequest(http.MethodDelete, ts.URL("/lift-maxes/"+newMax.Data.ID), nil)
		req.Header.Set("X-User-ID", userAID)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 204, got %d: %s", resp.StatusCode, respBody)
		}
	})
}

func TestLiftMaxCurrentMaxAuthorization(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	userAID := "current-user-a"
	userBID := "current-user-b"
	squatID := "00000000-0000-0000-0000-000000000001"

	// Create a lift max for user A
	body := fmt.Sprintf(`{"liftId": "%s", "type": "ONE_RM", "value": 315.0}`, squatID)
	createReq, _ := http.NewRequest(http.MethodPost, ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userAID)), bytes.NewBufferString(body))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("X-User-ID", userAID)
	createResp, _ := http.DefaultClient.Do(createReq)
	createResp.Body.Close()

	t.Run("unauthenticated user gets 401 on GET current", func(t *testing.T) {
		url := fmt.Sprintf("/users/%s/lift-maxes/current?lift=%s&type=ONE_RM", userAID, squatID)
		resp, err := http.Get(ts.URL(url))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("owner can access their current max", func(t *testing.T) {
		url := fmt.Sprintf("/users/%s/lift-maxes/current?lift=%s&type=ONE_RM", userAID, squatID)
		req, _ := http.NewRequest(http.MethodGet, ts.URL(url), nil)
		req.Header.Set("X-User-ID", userAID)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("non-owner cannot access another user's current max", func(t *testing.T) {
		url := fmt.Sprintf("/users/%s/lift-maxes/current?lift=%s&type=ONE_RM", userAID, squatID)
		req, _ := http.NewRequest(http.MethodGet, ts.URL(url), nil)
		req.Header.Set("X-User-ID", userBID)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("admin can access any user's current max", func(t *testing.T) {
		url := fmt.Sprintf("/users/%s/lift-maxes/current?lift=%s&type=ONE_RM", userAID, squatID)
		req, _ := http.NewRequest(http.MethodGet, ts.URL(url), nil)
		req.Header.Set("X-User-ID", testutil.TestAdminID)
		req.Header.Set("X-Admin", "true")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}
	})
}

func TestLiftMaxConvertAuthorization(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	userAID := "convert-user-a"
	userBID := "convert-user-b"
	squatID := "00000000-0000-0000-0000-000000000001"

	// Create a lift max for user A
	body := fmt.Sprintf(`{"liftId": "%s", "type": "ONE_RM", "value": 400.0}`, squatID)
	createReq, _ := http.NewRequest(http.MethodPost, ts.URL(fmt.Sprintf("/users/%s/lift-maxes", userAID)), bytes.NewBufferString(body))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("X-User-ID", userAID)
	createResp, _ := http.DefaultClient.Do(createReq)
	var max LiftMaxResponse
	json.NewDecoder(createResp.Body).Decode(&max)
	createResp.Body.Close()

	t.Run("unauthenticated user gets 401 on convert", func(t *testing.T) {
		url := fmt.Sprintf("/lift-maxes/%s/convert?to_type=TRAINING_MAX", max.Data.ID)
		resp, err := http.Get(ts.URL(url))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("owner can convert their lift max", func(t *testing.T) {
		url := fmt.Sprintf("/lift-maxes/%s/convert?to_type=TRAINING_MAX", max.Data.ID)
		req, _ := http.NewRequest(http.MethodGet, ts.URL(url), nil)
		req.Header.Set("X-User-ID", userAID)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("non-owner cannot convert another user's lift max", func(t *testing.T) {
		url := fmt.Sprintf("/lift-maxes/%s/convert?to_type=TRAINING_MAX", max.Data.ID)
		req, _ := http.NewRequest(http.MethodGet, ts.URL(url), nil)
		req.Header.Set("X-User-ID", userBID)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("admin can convert any user's lift max", func(t *testing.T) {
		url := fmt.Sprintf("/lift-maxes/%s/convert?to_type=TRAINING_MAX", max.Data.ID)
		req, _ := http.NewRequest(http.MethodGet, ts.URL(url), nil)
		req.Header.Set("X-User-ID", testutil.TestAdminID)
		req.Header.Set("X-Admin", "true")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, respBody)
		}
	})
}

func TestBearerTokenAuthentication(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Register and login to get a real token for Bearer auth tests
	registerBody := `{"email": "bearer@example.com", "password": "password123"}`
	regReq, _ := http.NewRequest(http.MethodPost, ts.URL("/auth/register"), bytes.NewBufferString(registerBody))
	regReq.Header.Set("Content-Type", "application/json")
	regResp, _ := http.DefaultClient.Do(regReq)
	regResp.Body.Close()

	loginBody := `{"email": "bearer@example.com", "password": "password123"}`
	loginReq, _ := http.NewRequest(http.MethodPost, ts.URL("/auth/login"), bytes.NewBufferString(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp, _ := http.DefaultClient.Do(loginReq)

	var loginResult struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	json.NewDecoder(loginResp.Body).Decode(&loginResult)
	loginResp.Body.Close()
	realToken := loginResult.Data.Token

	t.Run("can authenticate with Bearer token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/lifts"), nil)
		req.Header.Set("Authorization", "Bearer "+realToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("empty Bearer token gets 401", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/lifts"), nil)
		req.Header.Set("Authorization", "Bearer ")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("non-Bearer auth header falls back to X-User-ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/lifts"), nil)
		req.Header.Set("Authorization", "Basic something")
		req.Header.Set("X-User-ID", testutil.TestUserID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}
	})
}

func TestAuthorizationErrorMessages(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("unauthenticated returns clear error", func(t *testing.T) {
		resp, _ := http.Get(ts.URL("/lifts"))
		defer resp.Body.Close()

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)

		if errResp.Error.Message != "Authentication required" {
			t.Errorf("Expected error 'Authentication required', got '%s'", errResp.Error.Message)
		}
	})

	t.Run("forbidden returns clear error for non-admin", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/lifts"), bytes.NewBufferString(`{"name": "Test"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", testutil.TestUserID)

		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)

		if errResp.Error.Message != "Admin privileges required" {
			t.Errorf("Expected error 'Admin privileges required', got '%s'", errResp.Error.Message)
		}
	})
}
