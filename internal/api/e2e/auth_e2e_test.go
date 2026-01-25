// Package e2e provides end-to-end tests for complete API workflows.
package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// =============================================================================
// AUTH RESPONSE TYPES
// =============================================================================

// AuthUserResponseData represents the user object in auth responses.
type AuthUserResponseData struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// AuthUserEnvelopeData wraps user response in data envelope.
type AuthUserEnvelopeData struct {
	Data AuthUserResponseData `json:"data"`
}

// AuthLoginResponseData represents the login response.
type AuthLoginResponseData struct {
	Token     string               `json:"token"`
	ExpiresAt string               `json:"expiresAt"`
	User      AuthUserResponseData `json:"user"`
}

// AuthLoginEnvelopeData wraps login response in data envelope.
type AuthLoginEnvelopeData struct {
	Data AuthLoginResponseData `json:"data"`
}

// =============================================================================
// AUTH HELPER FUNCTIONS
// =============================================================================

// anonPostE2E performs an unauthenticated POST request.
func anonPostE2E(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}

// bearerGetE2E performs a GET request with Bearer token authentication.
func bearerGetE2E(url string, token string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return http.DefaultClient.Do(req)
}

// bearerPostE2E performs a POST request with Bearer token authentication.
func bearerPostE2E(url string, token string, body string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return http.DefaultClient.Do(req)
}

// registerUser registers a new user and returns the user data.
func registerUser(t *testing.T, ts *testutil.TestServer, email, password, name string) AuthUserResponseData {
	t.Helper()
	var body string
	if name != "" {
		body = fmt.Sprintf(`{"email": "%s", "password": "%s", "name": "%s"}`, email, password, name)
	} else {
		body = fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)
	}
	resp, err := anonPostE2E(ts.URL("/auth/register"), body)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to register user, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var result AuthUserEnvelopeData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode register response: %v", err)
	}
	return result.Data
}

// loginUser logs in a user and returns the login response.
func loginUser(t *testing.T, ts *testutil.TestServer, email, password string) AuthLoginResponseData {
	t.Helper()
	body := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)
	resp, err := anonPostE2E(ts.URL("/auth/login"), body)
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to login user, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var result AuthLoginEnvelopeData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}
	return result.Data
}

// logoutUser logs out a user.
func logoutUser(t *testing.T, ts *testutil.TestServer, token string) {
	t.Helper()
	resp, err := bearerPostE2E(ts.URL("/auth/logout"), token, "")
	if err != nil {
		t.Fatalf("Failed to logout user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to logout user, status %d: %s", resp.StatusCode, bodyBytes)
	}
}

// =============================================================================
// E2E TESTS: REGISTRATION
// =============================================================================

func TestAuthE2E_Registration(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]

	t.Run("successful registration with email, password, and name", func(t *testing.T) {
		email := fmt.Sprintf("user1-%s@example.com", testID)
		user := registerUser(t, ts, email, "password123", "Test User")

		if user.ID == "" {
			t.Error("Expected user ID to be set")
		}
		if user.Email != email {
			t.Errorf("Expected email '%s', got '%s'", email, user.Email)
		}
		if user.Name != "Test User" {
			t.Errorf("Expected name 'Test User', got '%s'", user.Name)
		}
	})

	t.Run("successful registration with email and password only", func(t *testing.T) {
		email := fmt.Sprintf("user2-%s@example.com", testID)
		user := registerUser(t, ts, email, "password123", "")

		if user.ID == "" {
			t.Error("Expected user ID to be set")
		}
		if user.Email != email {
			t.Errorf("Expected email '%s', got '%s'", email, user.Email)
		}
	})

	t.Run("fails with missing email", func(t *testing.T) {
		body := `{"password": "password123", "name": "No Email"}`
		resp, err := anonPostE2E(ts.URL("/auth/register"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("fails with invalid email format", func(t *testing.T) {
		body := `{"email": "notanemail", "password": "password123"}`
		resp, err := anonPostE2E(ts.URL("/auth/register"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("fails with missing password", func(t *testing.T) {
		body := `{"email": "missingpass@example.com", "name": "Missing Pass"}`
		resp, err := anonPostE2E(ts.URL("/auth/register"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("fails with short password", func(t *testing.T) {
		body := `{"email": "shortpass@example.com", "password": "short"}`
		resp, err := anonPostE2E(ts.URL("/auth/register"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("fails with duplicate email", func(t *testing.T) {
		email := fmt.Sprintf("duplicate-%s@example.com", testID)

		// First registration
		registerUser(t, ts, email, "password123", "First User")

		// Second registration with same email
		body := fmt.Sprintf(`{"email": "%s", "password": "password123", "name": "Second User"}`, email)
		resp, err := anonPostE2E(ts.URL("/auth/register"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409, got %d", resp.StatusCode)
		}
	})

	t.Run("email is case-insensitive", func(t *testing.T) {
		email := fmt.Sprintf("CaseTest-%s@example.com", testID)

		// First registration with mixed case
		registerUser(t, ts, email, "password123", "Case User")

		// Second registration with lowercase
		lowerEmail := fmt.Sprintf("casetest-%s@example.com", testID)
		body := fmt.Sprintf(`{"email": "%s", "password": "password123"}`, lowerEmail)
		resp, err := anonPostE2E(ts.URL("/auth/register"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409 (case-insensitive), got %d", resp.StatusCode)
		}
	})
}

// =============================================================================
// E2E TESTS: LOGIN
// =============================================================================

func TestAuthE2E_Login(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	email := fmt.Sprintf("login-%s@example.com", testID)

	// Register a user first
	registerUser(t, ts, email, "password123", "Login User")

	t.Run("successful login returns token and user", func(t *testing.T) {
		result := loginUser(t, ts, email, "password123")

		if result.Token == "" {
			t.Error("Expected token to be set")
		}
		if result.ExpiresAt == "" {
			t.Error("Expected expiresAt to be set")
		}
		if result.User.Email != email {
			t.Errorf("Expected email '%s', got '%s'", email, result.User.Email)
		}
		if result.User.Name != "Login User" {
			t.Errorf("Expected name 'Login User', got '%s'", result.User.Name)
		}
	})

	t.Run("fails with wrong email", func(t *testing.T) {
		body := `{"email": "nonexistent@example.com", "password": "password123"}`
		resp, err := anonPostE2E(ts.URL("/auth/login"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("fails with wrong password", func(t *testing.T) {
		body := fmt.Sprintf(`{"email": "%s", "password": "wrongpassword"}`, email)
		resp, err := anonPostE2E(ts.URL("/auth/login"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("fails with missing credentials", func(t *testing.T) {
		resp, err := anonPostE2E(ts.URL("/auth/login"), `{}`)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("token works for authenticated requests", func(t *testing.T) {
		result := loginUser(t, ts, email, "password123")

		// Use token to access /auth/me
		resp, err := bearerGetE2E(ts.URL("/auth/me"), result.Token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var meResult AuthUserEnvelopeData
		if err := json.NewDecoder(resp.Body).Decode(&meResult); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if meResult.Data.Email != email {
			t.Errorf("Expected email '%s', got '%s'", email, meResult.Data.Email)
		}
	})
}

// =============================================================================
// E2E TESTS: LOGOUT
// =============================================================================

func TestAuthE2E_Logout(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	email := fmt.Sprintf("logout-%s@example.com", testID)

	// Register and login
	registerUser(t, ts, email, "password123", "Logout User")
	loginResult := loginUser(t, ts, email, "password123")

	t.Run("returns 204 on successful logout", func(t *testing.T) {
		// Get a fresh token
		freshLogin := loginUser(t, ts, email, "password123")

		resp, err := bearerPostE2E(ts.URL("/auth/logout"), freshLogin.Token, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 204, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("invalidates session after logout", func(t *testing.T) {
		// Get a fresh token
		freshLogin := loginUser(t, ts, email, "password123")

		// Logout
		logoutUser(t, ts, freshLogin.Token)

		// Try to use the token
		resp, err := bearerGetE2E(ts.URL("/auth/me"), freshLogin.Token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 after logout, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 401 without session", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/auth/logout"), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	// Keep the original token for later tests
	_ = loginResult
}

// =============================================================================
// E2E TESTS: CURRENT USER
// =============================================================================

func TestAuthE2E_CurrentUser(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	email := fmt.Sprintf("me-%s@example.com", testID)

	// Register and login
	registerUser(t, ts, email, "password123", "Me User")
	loginResult := loginUser(t, ts, email, "password123")

	t.Run("returns user with valid session", func(t *testing.T) {
		resp, err := bearerGetE2E(ts.URL("/auth/me"), loginResult.Token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result AuthUserEnvelopeData
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.Email != email {
			t.Errorf("Expected email '%s', got '%s'", email, result.Data.Email)
		}
		if result.Data.Name != "Me User" {
			t.Errorf("Expected name 'Me User', got '%s'", result.Data.Name)
		}
	})

	t.Run("returns 401 without session", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/auth/me"), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 401 with invalid token", func(t *testing.T) {
		resp, err := bearerGetE2E(ts.URL("/auth/me"), "invalid-token-12345")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})
}

// =============================================================================
// E2E TESTS: BACKWARDS COMPATIBILITY
// =============================================================================

func TestAuthE2E_BackwardsCompatibility(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]

	// Use a seeded test user ID that exists in the test database
	userID := "workout-test-user"

	t.Run("X-User-ID header still works", func(t *testing.T) {
		// Access a protected endpoint using X-User-ID header
		req, err := http.NewRequest(http.MethodGet, ts.URL("/users/"+userID+"/lift-maxes"), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("X-User-ID", userID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Should succeed (200) - the endpoint exists and user is authenticated
		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200 with X-User-ID, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("X-Admin header still works", func(t *testing.T) {
		// Create a new lift (admin-only operation)
		liftSlug := "admin-test-lift-" + testID
		body := fmt.Sprintf(`{"name": "Admin Test Lift", "slug": "%s", "isCompetitionLift": false}`, liftSlug)

		req, err := http.NewRequest(http.MethodPost, ts.URL("/lifts"), bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", testutil.TestAdminID)
		req.Header.Set("X-Admin", "true")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 201 with X-Admin, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("Authorization Bearer takes precedence over X-User-ID", func(t *testing.T) {
		// Register and login to get a real token
		email := fmt.Sprintf("precedence-%s@example.com", testID)
		user := registerUser(t, ts, email, "password123", "Precedence User")
		loginResult := loginUser(t, ts, email, "password123")

		// Make request with both Bearer token and X-User-ID header
		// The Bearer token should take precedence
		req, err := http.NewRequest(http.MethodGet, ts.URL("/auth/me"), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+loginResult.Token)
		req.Header.Set("X-User-ID", "some-other-user-id") // This should be ignored

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result AuthUserEnvelopeData
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Should return the user from the Bearer token, not from X-User-ID
		if result.Data.ID != user.ID {
			t.Errorf("Expected user ID from Bearer token '%s', got '%s'", user.ID, result.Data.ID)
		}
		if result.Data.Email != email {
			t.Errorf("Expected email from Bearer token '%s', got '%s'", email, result.Data.Email)
		}
	})
}

// =============================================================================
// E2E TESTS: FULL FLOW INTEGRATION
// =============================================================================

func TestAuthE2E_FullFlow(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]

	t.Run("register -> login -> access protected endpoint -> logout", func(t *testing.T) {
		email := fmt.Sprintf("fullflow-%s@example.com", testID)

		// Step 1: Register
		user := registerUser(t, ts, email, "password123", "Full Flow User")
		if user.ID == "" {
			t.Fatal("Registration should return user ID")
		}

		// Step 2: Login
		loginResult := loginUser(t, ts, email, "password123")
		if loginResult.Token == "" {
			t.Fatal("Login should return token")
		}
		if loginResult.User.ID != user.ID {
			t.Errorf("Login user ID should match registered user ID")
		}

		// Step 3: Access protected endpoint (/auth/me)
		resp, err := bearerGetE2E(ts.URL("/auth/me"), loginResult.Token)
		if err != nil {
			t.Fatalf("Failed to access protected endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Protected endpoint should be accessible, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var meResult AuthUserEnvelopeData
		json.NewDecoder(resp.Body).Decode(&meResult)
		if meResult.Data.ID != user.ID {
			t.Errorf("Protected endpoint should return correct user")
		}

		// Step 4: Logout
		logoutUser(t, ts, loginResult.Token)

		// Step 5: Verify token is invalidated
		resp2, err := bearerGetE2E(ts.URL("/auth/me"), loginResult.Token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp2.Body.Close()

		if resp2.StatusCode != http.StatusUnauthorized {
			t.Errorf("Token should be invalidated after logout, got %d", resp2.StatusCode)
		}
	})

	t.Run("user can access own resources with session auth", func(t *testing.T) {
		email := fmt.Sprintf("ownresource-%s@example.com", testID)

		// Register and login
		user := registerUser(t, ts, email, "password123", "Own Resource User")
		loginResult := loginUser(t, ts, email, "password123")

		// Access user-specific endpoint (lift maxes for the logged-in user)
		req, err := http.NewRequest(http.MethodGet, ts.URL("/users/"+user.ID+"/lift-maxes"), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+loginResult.Token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Should succeed - user accessing their own resources
		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("User should access own resources, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("user cannot access other users resources", func(t *testing.T) {
		email := fmt.Sprintf("otherresource-%s@example.com", testID)

		// Register and login user A
		registerUser(t, ts, email, "password123", "User A")
		loginResult := loginUser(t, ts, email, "password123")

		// Try to access another user's resources
		otherUserID := "some-other-user-id-12345"
		req, err := http.NewRequest(http.MethodGet, ts.URL("/users/"+otherUserID+"/lift-maxes"), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+loginResult.Token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Should be forbidden (403) - user accessing other user's resources
		if resp.StatusCode != http.StatusForbidden {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("User should NOT access other user's resources, expected 403, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})
}

// =============================================================================
// E2E TESTS: SESSION INTEGRATION WITH PROTECTED ENDPOINTS
// =============================================================================

func TestAuthE2E_ProtectedEndpointAccess(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]

	t.Run("authenticated user can create and access lift maxes", func(t *testing.T) {
		email := fmt.Sprintf("liftmax-%s@example.com", testID)

		// Register and login
		user := registerUser(t, ts, email, "password123", "LiftMax User")
		loginResult := loginUser(t, ts, email, "password123")

		// Use seeded squat lift ID
		squatID := "00000000-0000-0000-0000-000000000001"

		// Create a lift max using session auth
		createBody := fmt.Sprintf(`{"liftId": "%s", "type": "TRAINING_MAX", "value": 225.0}`, squatID)
		req, err := http.NewRequest(http.MethodPost, ts.URL("/users/"+user.ID+"/lift-maxes"), bytes.NewBufferString(createBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+loginResult.Token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		// Retrieve lift maxes using session auth
		req2, err := http.NewRequest(http.MethodGet, ts.URL("/users/"+user.ID+"/lift-maxes"), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req2.Header.Set("Authorization", "Bearer "+loginResult.Token)

		resp2, err := http.DefaultClient.Do(req2)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp2.Body.Close()

		if resp2.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp2.Body)
			t.Errorf("Expected status 200, got %d: %s", resp2.StatusCode, bodyBytes)
		}
	})

	t.Run("unauthenticated requests to protected endpoints fail", func(t *testing.T) {
		// Try to access a protected endpoint without authentication
		req, err := http.NewRequest(http.MethodGet, ts.URL("/users/some-user/lift-maxes"), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
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
}
