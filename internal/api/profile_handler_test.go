package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// ProfileData represents the profile object in responses.
type ProfileData struct {
	ID         string  `json:"id"`
	Email      string  `json:"email"`
	Name       *string `json:"name"`
	WeightUnit string  `json:"weightUnit"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

// ProfileResponseEnvelope wraps profile response in data envelope.
type ProfileResponseEnvelope struct {
	Data ProfileData `json:"data"`
}

// createTestUserForProfile creates a test user via the auth API and returns the user ID.
func createTestUserForProfile(t *testing.T, ts *testutil.TestServer, email, password, name string) string {
	t.Helper()
	body := `{"email": "` + email + `", "password": "` + password + `", "name": "` + name + `"}`
	resp, err := anonPost(ts.URL("/auth/register"), body)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create user, status %d: %s", resp.StatusCode, respBody)
	}

	var result AuthUserEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode user response: %v", err)
	}
	return result.Data.ID
}

// userGetProfile performs a GET request with X-User-ID header (test mode auth).
func userGetProfile(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

// userPutProfile performs a PUT request with X-User-ID header (test mode auth).
func userPutProfile(url string, userID string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func TestProfileGet(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a test user
	userID := createTestUserForProfile(t, ts, "profile-get@example.com", "password123", "Profile Get User")

	t.Run("gets own profile successfully", func(t *testing.T) {
		resp, err := userGetProfile(ts.URL("/users/"+userID+"/profile"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, respBody)
		}

		var result ProfileResponseEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.ID != userID {
			t.Errorf("Expected ID %s, got %s", userID, result.Data.ID)
		}
		if result.Data.Email != "profile-get@example.com" {
			t.Errorf("Expected email 'profile-get@example.com', got %s", result.Data.Email)
		}
		if result.Data.Name == nil || *result.Data.Name != "Profile Get User" {
			t.Errorf("Expected name 'Profile Get User', got %v", result.Data.Name)
		}
		if result.Data.WeightUnit != "lb" {
			t.Errorf("Expected default weightUnit 'lb', got %s", result.Data.WeightUnit)
		}
		if result.Data.CreatedAt == "" {
			t.Error("Expected createdAt to be set")
		}
		if result.Data.UpdatedAt == "" {
			t.Error("Expected updatedAt to be set")
		}
	})

	t.Run("admin can get other user profile", func(t *testing.T) {
		resp, err := adminGet(ts.URL("/users/" + userID + "/profile"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, respBody)
		}

		var result ProfileResponseEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.ID != userID {
			t.Errorf("Expected ID %s, got %s", userID, result.Data.ID)
		}
	})

	t.Run("returns 403 when accessing other user profile", func(t *testing.T) {
		otherUserID := "other-user-id"
		resp, err := userGetProfile(ts.URL("/users/"+otherUserID+"/profile"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("returns 401 without authentication", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/users/" + userID + "/profile"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 401, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("returns 404 for non-existent user", func(t *testing.T) {
		nonExistentID := "non-existent-user-id"
		resp, err := adminGet(ts.URL("/users/" + nonExistentID + "/profile"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 404, got %d: %s", resp.StatusCode, respBody)
		}
	})
}

func TestProfileUpdate(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a test user
	userID := createTestUserForProfile(t, ts, "profile-update@example.com", "password123", "Profile Update User")

	t.Run("updates own profile name successfully", func(t *testing.T) {
		body := `{"name": "Updated Name"}`
		resp, err := userPutProfile(ts.URL("/users/"+userID+"/profile"), userID, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, respBody)
		}

		var result ProfileResponseEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.Name == nil || *result.Data.Name != "Updated Name" {
			t.Errorf("Expected name 'Updated Name', got %v", result.Data.Name)
		}
	})

	t.Run("updates own profile weight unit successfully", func(t *testing.T) {
		body := `{"weightUnit": "kg"}`
		resp, err := userPutProfile(ts.URL("/users/"+userID+"/profile"), userID, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, respBody)
		}

		var result ProfileResponseEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.WeightUnit != "kg" {
			t.Errorf("Expected weightUnit 'kg', got %s", result.Data.WeightUnit)
		}
	})

	t.Run("updates both name and weight unit", func(t *testing.T) {
		body := `{"name": "Both Updated", "weightUnit": "lb"}`
		resp, err := userPutProfile(ts.URL("/users/"+userID+"/profile"), userID, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, respBody)
		}

		var result ProfileResponseEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.Name == nil || *result.Data.Name != "Both Updated" {
			t.Errorf("Expected name 'Both Updated', got %v", result.Data.Name)
		}
		if result.Data.WeightUnit != "lb" {
			t.Errorf("Expected weightUnit 'lb', got %s", result.Data.WeightUnit)
		}
	})

	t.Run("clears name with empty string", func(t *testing.T) {
		body := `{"name": ""}`
		resp, err := userPutProfile(ts.URL("/users/"+userID+"/profile"), userID, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, respBody)
		}

		var result ProfileResponseEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.Name != nil {
			t.Errorf("Expected name to be null, got %v", result.Data.Name)
		}
	})

	t.Run("empty body returns current profile", func(t *testing.T) {
		// First set a known state
		setupBody := `{"name": "Known Name", "weightUnit": "kg"}`
		_, err := userPutProfile(ts.URL("/users/"+userID+"/profile"), userID, setupBody)
		if err != nil {
			t.Fatalf("Failed to set up profile: %v", err)
		}

		// Then call with empty body
		body := `{}`
		resp, err := userPutProfile(ts.URL("/users/"+userID+"/profile"), userID, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, respBody)
		}

		var result ProfileResponseEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.Name == nil || *result.Data.Name != "Known Name" {
			t.Errorf("Expected name 'Known Name' to be preserved, got %v", result.Data.Name)
		}
		if result.Data.WeightUnit != "kg" {
			t.Errorf("Expected weightUnit 'kg' to be preserved, got %s", result.Data.WeightUnit)
		}
	})

	t.Run("admin cannot update other user profile", func(t *testing.T) {
		body := `{"name": "Admin Updated"}`
		resp, err := adminPut(ts.URL("/users/"+userID+"/profile"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("returns 400 for invalid weight unit", func(t *testing.T) {
		body := `{"weightUnit": "stone"}`
		resp, err := userPutProfile(ts.URL("/users/"+userID+"/profile"), userID, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 400, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("returns 400 for name too long", func(t *testing.T) {
		// Name over 100 characters
		longName := "This is a very long name that exceeds the maximum allowed length of one hundred characters for a user profile name"
		body := `{"name": "` + longName + `"}`
		resp, err := userPutProfile(ts.URL("/users/"+userID+"/profile"), userID, body)
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
		body := `{invalid json}`
		resp, err := userPutProfile(ts.URL("/users/"+userID+"/profile"), userID, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 400, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("returns 403 when updating other user profile", func(t *testing.T) {
		otherUserID := "other-user-id"
		body := `{"name": "Forbidden Update"}`
		resp, err := userPutProfile(ts.URL("/users/"+otherUserID+"/profile"), userID, body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("returns 401 without authentication", func(t *testing.T) {
		body := `{"name": "No Auth Update"}`
		req, err := http.NewRequest(http.MethodPut, ts.URL("/users/"+userID+"/profile"), bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 401, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("returns 403 for non-owner trying to update non-existent user", func(t *testing.T) {
		// Authorization check happens before database access, so even for non-existent users,
		// non-owners (including admins) get 403 Forbidden, not 404.
		nonExistentID := "non-existent-user-id"
		body := `{"name": "Non Existent"}`
		resp, err := adminPut(ts.URL("/users/"+nonExistentID+"/profile"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 403, got %d: %s", resp.StatusCode, respBody)
		}
	})
}
