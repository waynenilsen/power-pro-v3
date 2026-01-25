// Package e2e provides end-to-end tests for complete API workflows.
// This file contains E2E tests for the Profile endpoints.
package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// =============================================================================
// PROFILE RESPONSE TYPES
// =============================================================================

// ProfileResponseData represents the profile object in responses.
type ProfileResponseData struct {
	ID         string  `json:"id"`
	Email      string  `json:"email"`
	Name       *string `json:"name"`
	WeightUnit string  `json:"weightUnit"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

// ProfileEnvelopeData wraps profile response in data envelope.
type ProfileEnvelopeData struct {
	Data ProfileResponseData `json:"data"`
}

// ErrorResponseProfile represents the API error response structure.
type ErrorResponseProfile struct {
	Error struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Details interface{} `json:"details,omitempty"`
	} `json:"error"`
}

// =============================================================================
// PROFILE HELPER FUNCTIONS
// =============================================================================

// registerUserForProfile registers a new user and returns the user data and token.
func registerUserForProfile(t *testing.T, ts *testutil.TestServer, email, password, name string) (AuthUserResponseData, string) {
	t.Helper()

	// Register
	user := registerUser(t, ts, email, password, name)

	// Login to get token
	loginResult := loginUser(t, ts, email, password)

	return user, loginResult.Token
}

// bearerGetProfile performs a GET request with Bearer token authentication.
func bearerGetProfile(url string, token string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return http.DefaultClient.Do(req)
}

// bearerPutProfile performs a PUT request with Bearer token authentication.
func bearerPutProfile(url string, token string, body string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}
	req, err := http.NewRequest(http.MethodPut, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return http.DefaultClient.Do(req)
}

// userIDGetProfile performs a GET request with X-User-ID header (test mode auth).
func userIDGetProfile(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

// userIDPutProfile performs a PUT request with X-User-ID header (test mode auth).
func userIDPutProfile(url string, userID string, body string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}
	req, err := http.NewRequest(http.MethodPut, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return http.DefaultClient.Do(req)
}

// adminGetProfile performs an admin-authenticated GET request.
func adminGetProfile(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// adminPutProfile performs an admin-authenticated PUT request.
func adminPutProfile(url string, body string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}
	req, err := http.NewRequest(http.MethodPut, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return http.DefaultClient.Do(req)
}

// =============================================================================
// E2E TESTS: GET /users/{id}/profile
// =============================================================================

func TestProfileE2E_GetProfile(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]

	// Create test users
	userAEmail := fmt.Sprintf("profile-get-a-%s@example.com", testID)
	userA, tokenA := registerUserForProfile(t, ts, userAEmail, "password123", "User A")

	userBEmail := fmt.Sprintf("profile-get-b-%s@example.com", testID)
	userB, tokenB := registerUserForProfile(t, ts, userBEmail, "password123", "User B")

	t.Run("returns 200 with correct profile data for owner using Bearer token", func(t *testing.T) {
		resp, err := bearerGetProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenA)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result ProfileEnvelopeData
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Verify profile data
		if result.Data.ID != userA.ID {
			t.Errorf("Expected ID '%s', got '%s'", userA.ID, result.Data.ID)
		}
		if result.Data.Email != userAEmail {
			t.Errorf("Expected email '%s', got '%s'", userAEmail, result.Data.Email)
		}
		if result.Data.Name == nil || *result.Data.Name != "User A" {
			t.Errorf("Expected name 'User A', got %v", result.Data.Name)
		}
		if result.Data.WeightUnit != "lb" {
			t.Errorf("Expected default weightUnit 'lb', got '%s'", result.Data.WeightUnit)
		}
	})

	t.Run("returns 200 with correct profile data for owner using X-User-ID", func(t *testing.T) {
		resp, err := userIDGetProfile(ts.URL("/users/"+userA.ID+"/profile"), userA.ID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result ProfileEnvelopeData
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.ID != userA.ID {
			t.Errorf("Expected ID '%s', got '%s'", userA.ID, result.Data.ID)
		}
	})

	t.Run("returns 200 with correct profile data for admin viewing other user", func(t *testing.T) {
		resp, err := adminGetProfile(ts.URL("/users/" + userA.ID + "/profile"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result ProfileEnvelopeData
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.ID != userA.ID {
			t.Errorf("Expected ID '%s', got '%s'", userA.ID, result.Data.ID)
		}
		if result.Data.Email != userAEmail {
			t.Errorf("Expected email '%s', got '%s'", userAEmail, result.Data.Email)
		}
	})

	t.Run("returns 401 without authentication", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/users/"+userA.ID+"/profile"), nil)
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

	t.Run("returns 403 for non-owner non-admin", func(t *testing.T) {
		// User B trying to access User A's profile
		resp, err := bearerGetProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenB)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("returns 404 for non-existent user", func(t *testing.T) {
		nonExistentID := uuid.New().String()
		resp, err := adminGetProfile(ts.URL("/users/" + nonExistentID + "/profile"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 404, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("response has correct JSON structure with camelCase fields", func(t *testing.T) {
		resp, err := bearerGetProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenA)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyStr := string(bodyBytes)

		// Check for camelCase fields
		expectedFields := []string{"id", "email", "name", "weightUnit", "createdAt", "updatedAt"}
		for _, field := range expectedFields {
			if !strings.Contains(bodyStr, `"`+field+`"`) {
				t.Errorf("Expected camelCase field '%s' in response, got: %s", field, bodyStr)
			}
		}

		// Check that snake_case is NOT present
		unexpectedFields := []string{"weight_unit", "created_at", "updated_at"}
		for _, field := range unexpectedFields {
			if strings.Contains(bodyStr, `"`+field+`"`) {
				t.Errorf("Unexpected snake_case field '%s' in response, got: %s", field, bodyStr)
			}
		}
	})

	_ = userB
	_ = tokenB
}

// =============================================================================
// E2E TESTS: PUT /users/{id}/profile
// =============================================================================

func TestProfileE2E_UpdateProfile(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]

	// Create test users
	userAEmail := fmt.Sprintf("profile-update-a-%s@example.com", testID)
	userA, tokenA := registerUserForProfile(t, ts, userAEmail, "password123", "User A")

	userBEmail := fmt.Sprintf("profile-update-b-%s@example.com", testID)
	userB, tokenB := registerUserForProfile(t, ts, userBEmail, "password123", "User B")

	t.Run("returns 200 and updates name", func(t *testing.T) {
		body := `{"name": "Updated Name A"}`
		resp, err := bearerPutProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenA, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result ProfileEnvelopeData
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.Name == nil || *result.Data.Name != "Updated Name A" {
			t.Errorf("Expected name 'Updated Name A', got %v", result.Data.Name)
		}
	})

	t.Run("returns 200 and updates weightUnit to kg", func(t *testing.T) {
		body := `{"weightUnit": "kg"}`
		resp, err := bearerPutProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenA, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result ProfileEnvelopeData
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.WeightUnit != "kg" {
			t.Errorf("Expected weightUnit 'kg', got '%s'", result.Data.WeightUnit)
		}
	})

	t.Run("returns 200 and updates weightUnit to lb", func(t *testing.T) {
		// First set to kg
		body := `{"weightUnit": "kg"}`
		_, _ = bearerPutProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenA, body)

		// Then set back to lb
		body = `{"weightUnit": "lb"}`
		resp, err := bearerPutProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenA, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result ProfileEnvelopeData
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.WeightUnit != "lb" {
			t.Errorf("Expected weightUnit 'lb', got '%s'", result.Data.WeightUnit)
		}
	})

	t.Run("returns 200 with partial update (only name)", func(t *testing.T) {
		// Set a known state first
		setupBody := `{"name": "Before Partial", "weightUnit": "kg"}`
		_, _ = bearerPutProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenA, setupBody)

		// Update only name
		body := `{"name": "After Partial"}`
		resp, err := bearerPutProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenA, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result ProfileEnvelopeData
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.Name == nil || *result.Data.Name != "After Partial" {
			t.Errorf("Expected name 'After Partial', got %v", result.Data.Name)
		}
		// weightUnit should remain unchanged
		if result.Data.WeightUnit != "kg" {
			t.Errorf("Expected weightUnit 'kg' (unchanged), got '%s'", result.Data.WeightUnit)
		}
	})

	t.Run("returns 200 with partial update (only weightUnit)", func(t *testing.T) {
		// Set a known state first
		setupBody := `{"name": "Preserved Name", "weightUnit": "lb"}`
		_, _ = bearerPutProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenA, setupBody)

		// Update only weightUnit
		body := `{"weightUnit": "kg"}`
		resp, err := bearerPutProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenA, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result ProfileEnvelopeData
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// name should remain unchanged
		if result.Data.Name == nil || *result.Data.Name != "Preserved Name" {
			t.Errorf("Expected name 'Preserved Name' (unchanged), got %v", result.Data.Name)
		}
		if result.Data.WeightUnit != "kg" {
			t.Errorf("Expected weightUnit 'kg', got '%s'", result.Data.WeightUnit)
		}
	})

	t.Run("returns 200 with empty body (no-op)", func(t *testing.T) {
		// Set a known state first
		setupBody := `{"name": "No-Op Name", "weightUnit": "lb"}`
		_, _ = bearerPutProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenA, setupBody)

		// Send empty body
		body := `{}`
		resp, err := bearerPutProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenA, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result ProfileEnvelopeData
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// All fields should remain unchanged
		if result.Data.Name == nil || *result.Data.Name != "No-Op Name" {
			t.Errorf("Expected name 'No-Op Name' (unchanged), got %v", result.Data.Name)
		}
		if result.Data.WeightUnit != "lb" {
			t.Errorf("Expected weightUnit 'lb' (unchanged), got '%s'", result.Data.WeightUnit)
		}
	})

	t.Run("clears name when empty string provided", func(t *testing.T) {
		// Set a name first
		setupBody := `{"name": "Name To Clear"}`
		_, _ = bearerPutProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenA, setupBody)

		// Clear with empty string
		body := `{"name": ""}`
		resp, err := bearerPutProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenA, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result ProfileEnvelopeData
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.Name != nil {
			t.Errorf("Expected name to be null, got %v", result.Data.Name)
		}
	})

	t.Run("returns 400 for invalid weightUnit", func(t *testing.T) {
		body := `{"weightUnit": "stone"}`
		resp, err := bearerPutProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenA, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 400, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("returns 400 for name exceeding 100 characters", func(t *testing.T) {
		// Generate a name over 100 characters
		longName := strings.Repeat("a", 101)
		body := fmt.Sprintf(`{"name": "%s"}`, longName)
		resp, err := bearerPutProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenA, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 400, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("returns 401 without authentication", func(t *testing.T) {
		body := `{"name": "No Auth Update"}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/users/"+userA.ID+"/profile"), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

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

	t.Run("returns 403 for non-owner", func(t *testing.T) {
		// User B trying to update User A's profile
		body := `{"name": "Forbidden Update"}`
		resp, err := bearerPutProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenB, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("returns 403 for admin trying to update other user", func(t *testing.T) {
		body := `{"name": "Admin Forbidden Update"}`
		resp, err := adminPutProfile(ts.URL("/users/"+userA.ID+"/profile"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("returns 403 for non-owner even if user does not exist", func(t *testing.T) {
		// Authorization check happens before database lookup
		nonExistentID := uuid.New().String()
		body := `{"name": "Non Existent"}`
		resp, err := bearerPutProfile(ts.URL("/users/"+nonExistentID+"/profile"), tokenA, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// User A is not the owner of the non-existent user's profile
		if resp.StatusCode != http.StatusForbidden {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("verify updated_at changes on update", func(t *testing.T) {
		// Get current profile
		resp1, err := bearerGetProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenA)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp1.Body.Close()

		var result1 ProfileEnvelopeData
		if err := json.NewDecoder(resp1.Body).Decode(&result1); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		originalUpdatedAt := result1.Data.UpdatedAt

		// Wait to ensure timestamp difference (RFC3339 has second precision)
		time.Sleep(1100 * time.Millisecond)

		// Update profile
		body := `{"name": "Timestamp Test"}`
		resp2, err := bearerPutProfile(ts.URL("/users/"+userA.ID+"/profile"), tokenA, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp2.Body.Close()

		var result2 ProfileEnvelopeData
		if err := json.NewDecoder(resp2.Body).Decode(&result2); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result2.Data.UpdatedAt == originalUpdatedAt {
			t.Error("Expected updated_at to change after update")
		}
	})

	t.Run("works with X-User-ID header", func(t *testing.T) {
		body := `{"name": "X-User-ID Update"}`
		resp, err := userIDPutProfile(ts.URL("/users/"+userA.ID+"/profile"), userA.ID, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result ProfileEnvelopeData
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.Name == nil || *result.Data.Name != "X-User-ID Update" {
			t.Errorf("Expected name 'X-User-ID Update', got %v", result.Data.Name)
		}
	})

	_ = userB
	_ = tokenB
}

// =============================================================================
// E2E TESTS: WEIGHT UNIT DEFAULT
// =============================================================================

func TestProfileE2E_WeightUnitDefault(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]

	t.Run("new users have weight_unit defaulted to lb", func(t *testing.T) {
		email := fmt.Sprintf("default-weight-%s@example.com", testID)
		user, token := registerUserForProfile(t, ts, email, "password123", "Default Weight User")

		// Get profile immediately after registration
		resp, err := bearerGetProfile(ts.URL("/users/"+user.ID+"/profile"), token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result ProfileEnvelopeData
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.WeightUnit != "lb" {
			t.Errorf("Expected default weightUnit 'lb', got '%s'", result.Data.WeightUnit)
		}
	})
}

// =============================================================================
// E2E TESTS: FULL FLOW INTEGRATION
// =============================================================================

func TestProfileE2E_FullFlow(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]

	t.Run("register -> get profile -> update profile -> verify changes", func(t *testing.T) {
		email := fmt.Sprintf("fullflow-profile-%s@example.com", testID)

		// Step 1: Register
		user, token := registerUserForProfile(t, ts, email, "password123", "Full Flow User")

		// Step 2: Get profile
		resp1, err := bearerGetProfile(ts.URL("/users/"+user.ID+"/profile"), token)
		if err != nil {
			t.Fatalf("Failed to get profile: %v", err)
		}
		defer resp1.Body.Close()

		var profile1 ProfileEnvelopeData
		json.NewDecoder(resp1.Body).Decode(&profile1)

		if profile1.Data.WeightUnit != "lb" {
			t.Errorf("Expected default weightUnit 'lb', got '%s'", profile1.Data.WeightUnit)
		}
		if profile1.Data.Name == nil || *profile1.Data.Name != "Full Flow User" {
			t.Errorf("Expected name 'Full Flow User', got %v", profile1.Data.Name)
		}

		// Step 3: Update profile
		updateBody := `{"name": "Updated Full Flow", "weightUnit": "kg"}`
		resp2, err := bearerPutProfile(ts.URL("/users/"+user.ID+"/profile"), token, updateBody)
		if err != nil {
			t.Fatalf("Failed to update profile: %v", err)
		}
		defer resp2.Body.Close()

		if resp2.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp2.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp2.StatusCode, bodyBytes)
		}

		// Step 4: Get profile again to verify changes persisted
		resp3, err := bearerGetProfile(ts.URL("/users/"+user.ID+"/profile"), token)
		if err != nil {
			t.Fatalf("Failed to get profile: %v", err)
		}
		defer resp3.Body.Close()

		var profile2 ProfileEnvelopeData
		json.NewDecoder(resp3.Body).Decode(&profile2)

		if profile2.Data.Name == nil || *profile2.Data.Name != "Updated Full Flow" {
			t.Errorf("Expected updated name 'Updated Full Flow', got %v", profile2.Data.Name)
		}
		if profile2.Data.WeightUnit != "kg" {
			t.Errorf("Expected updated weightUnit 'kg', got '%s'", profile2.Data.WeightUnit)
		}
	})
}
