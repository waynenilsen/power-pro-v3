package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// AuthUserResponse represents the user object in auth responses.
type AuthUserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// AuthUserEnvelope wraps user response in data envelope.
type AuthUserEnvelope struct {
	Data AuthUserResponse `json:"data"`
}

// AuthLoginResponse represents the login response.
type AuthLoginResponse struct {
	Token     string           `json:"token"`
	ExpiresAt string           `json:"expiresAt"`
	User      AuthUserResponse `json:"user"`
}

// AuthLoginEnvelope wraps login response in data envelope.
type AuthLoginEnvelope struct {
	Data AuthLoginResponse `json:"data"`
}

// anonPost performs an unauthenticated POST request.
func anonPost(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}

// bearerGet performs a GET request with Bearer token authentication.
func bearerGet(url string, token string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return http.DefaultClient.Do(req)
}

// bearerPost performs a POST request with Bearer token authentication.
func bearerPost(url string, token string, body string) (*http.Response, error) {
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

func TestAuthRegister(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("registers new user successfully", func(t *testing.T) {
		body := `{"email": "newuser@example.com", "password": "password123", "name": "New User"}`
		resp, err := anonPost(ts.URL("/auth/register"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			respBody, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, respBody)
		}

		var result AuthUserEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.ID == "" {
			t.Error("Expected user ID to be set")
		}
		if result.Data.Email != "newuser@example.com" {
			t.Errorf("Expected email 'newuser@example.com', got '%s'", result.Data.Email)
		}
		if result.Data.Name != "New User" {
			t.Errorf("Expected name 'New User', got '%s'", result.Data.Name)
		}
		if result.Data.CreatedAt == "" {
			t.Error("Expected createdAt to be set")
		}
		if result.Data.UpdatedAt == "" {
			t.Error("Expected updatedAt to be set")
		}
	})

	t.Run("returns 400 for missing email", func(t *testing.T) {
		body := `{"password": "password123", "name": "No Email User"}`
		resp, err := anonPost(ts.URL("/auth/register"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid email format", func(t *testing.T) {
		body := `{"email": "notanemail", "password": "password123"}`
		resp, err := anonPost(ts.URL("/auth/register"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for missing password", func(t *testing.T) {
		body := `{"email": "nopassword@example.com", "name": "No Password"}`
		resp, err := anonPost(ts.URL("/auth/register"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for password too short", func(t *testing.T) {
		body := `{"email": "shortpass@example.com", "password": "short"}`
		resp, err := anonPost(ts.URL("/auth/register"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 409 for duplicate email", func(t *testing.T) {
		// First registration
		body := `{"email": "duplicate@example.com", "password": "password123"}`
		resp1, err := anonPost(ts.URL("/auth/register"), body)
		if err != nil {
			t.Fatalf("Failed to make first request: %v", err)
		}
		resp1.Body.Close()

		// Second registration with same email
		resp2, err := anonPost(ts.URL("/auth/register"), body)
		if err != nil {
			t.Fatalf("Failed to make second request: %v", err)
		}
		defer resp2.Body.Close()

		if resp2.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409, got %d", resp2.StatusCode)
		}
	})

	t.Run("email is case-insensitive for uniqueness", func(t *testing.T) {
		// First registration
		body1 := `{"email": "CaseSensitive@example.com", "password": "password123"}`
		resp1, err := anonPost(ts.URL("/auth/register"), body1)
		if err != nil {
			t.Fatalf("Failed to make first request: %v", err)
		}
		resp1.Body.Close()

		// Second registration with different case
		body2 := `{"email": "casesensitive@example.com", "password": "password123"}`
		resp2, err := anonPost(ts.URL("/auth/register"), body2)
		if err != nil {
			t.Fatalf("Failed to make second request: %v", err)
		}
		defer resp2.Body.Close()

		if resp2.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409 (email should be case-insensitive), got %d", resp2.StatusCode)
		}
	})

	t.Run("name is optional", func(t *testing.T) {
		body := `{"email": "noname@example.com", "password": "password123"}`
		resp, err := anonPost(ts.URL("/auth/register"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 201, got %d: %s", resp.StatusCode, respBody)
		}
	})
}

func TestAuthLogin(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Register a user first
	registerBody := `{"email": "login@example.com", "password": "password123", "name": "Login User"}`
	regResp, err := anonPost(ts.URL("/auth/register"), registerBody)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}
	regResp.Body.Close()

	t.Run("logs in with valid credentials", func(t *testing.T) {
		body := `{"email": "login@example.com", "password": "password123"}`
		resp, err := anonPost(ts.URL("/auth/login"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, respBody)
		}

		var result AuthLoginEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.Token == "" {
			t.Error("Expected token to be set")
		}
		if result.Data.ExpiresAt == "" {
			t.Error("Expected expiresAt to be set")
		}
		if result.Data.User.Email != "login@example.com" {
			t.Errorf("Expected email 'login@example.com', got '%s'", result.Data.User.Email)
		}
		if result.Data.User.Name != "Login User" {
			t.Errorf("Expected name 'Login User', got '%s'", result.Data.User.Name)
		}
	})

	t.Run("returns 400 for missing email", func(t *testing.T) {
		body := `{"password": "password123"}`
		resp, err := anonPost(ts.URL("/auth/login"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for missing password", func(t *testing.T) {
		body := `{"email": "login@example.com"}`
		resp, err := anonPost(ts.URL("/auth/login"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 401 for wrong password", func(t *testing.T) {
		body := `{"email": "login@example.com", "password": "wrongpassword"}`
		resp, err := anonPost(ts.URL("/auth/login"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 401 for non-existent email", func(t *testing.T) {
		body := `{"email": "nonexistent@example.com", "password": "password123"}`
		resp, err := anonPost(ts.URL("/auth/login"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("login is case-insensitive for email", func(t *testing.T) {
		body := `{"email": "LOGIN@example.com", "password": "password123"}`
		resp, err := anonPost(ts.URL("/auth/login"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 (case-insensitive login), got %d", resp.StatusCode)
		}
	})
}

func TestAuthLogout(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Register and login to get a token
	registerBody := `{"email": "logout@example.com", "password": "password123"}`
	regResp, _ := anonPost(ts.URL("/auth/register"), registerBody)
	regResp.Body.Close()

	loginBody := `{"email": "logout@example.com", "password": "password123"}`
	loginResp, _ := anonPost(ts.URL("/auth/login"), loginBody)
	var loginResult AuthLoginEnvelope
	json.NewDecoder(loginResp.Body).Decode(&loginResult)
	loginResp.Body.Close()
	token := loginResult.Data.Token

	t.Run("logs out with valid token", func(t *testing.T) {
		resp, err := bearerPost(ts.URL("/auth/logout"), token, "")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 204, got %d: %s", resp.StatusCode, respBody)
		}
	})

	t.Run("logout is idempotent", func(t *testing.T) {
		// Get a new token first
		loginResp2, _ := anonPost(ts.URL("/auth/login"), loginBody)
		var loginResult2 AuthLoginEnvelope
		json.NewDecoder(loginResp2.Body).Decode(&loginResult2)
		loginResp2.Body.Close()
		token2 := loginResult2.Data.Token

		// First logout
		resp1, err := bearerPost(ts.URL("/auth/logout"), token2, "")
		if err != nil {
			t.Fatalf("Failed to make first request: %v", err)
		}
		resp1.Body.Close()

		if resp1.StatusCode != http.StatusNoContent {
			t.Errorf("Expected status 204 for first logout, got %d", resp1.StatusCode)
		}

		// Second logout with same token should also succeed (or return 401)
		resp2, err := bearerPost(ts.URL("/auth/logout"), token2, "")
		if err != nil {
			t.Fatalf("Failed to make second request: %v", err)
		}
		resp2.Body.Close()

		// After the token is invalidated, the second logout should return 401
		// because the middleware validates the session first
		if resp2.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for second logout with invalidated token, got %d", resp2.StatusCode)
		}
	})

	t.Run("returns 401 without authentication", func(t *testing.T) {
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
}

func TestAuthMe(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Register and login to get a token
	registerBody := `{"email": "me@example.com", "password": "password123", "name": "Me User"}`
	regResp, _ := anonPost(ts.URL("/auth/register"), registerBody)
	regResp.Body.Close()

	loginBody := `{"email": "me@example.com", "password": "password123"}`
	loginResp, _ := anonPost(ts.URL("/auth/login"), loginBody)
	var loginResult AuthLoginEnvelope
	json.NewDecoder(loginResp.Body).Decode(&loginResult)
	loginResp.Body.Close()
	token := loginResult.Data.Token

	t.Run("returns current user with valid token", func(t *testing.T) {
		resp, err := bearerGet(ts.URL("/auth/me"), token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, respBody)
		}

		var result AuthUserEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Data.Email != "me@example.com" {
			t.Errorf("Expected email 'me@example.com', got '%s'", result.Data.Email)
		}
		if result.Data.Name != "Me User" {
			t.Errorf("Expected name 'Me User', got '%s'", result.Data.Name)
		}
	})

	t.Run("returns 401 without authentication", func(t *testing.T) {
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
		resp, err := bearerGet(ts.URL("/auth/me"), "invalid-token")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 401 after logout", func(t *testing.T) {
		// Get a fresh token
		loginResp2, _ := anonPost(ts.URL("/auth/login"), loginBody)
		var loginResult2 AuthLoginEnvelope
		json.NewDecoder(loginResp2.Body).Decode(&loginResult2)
		loginResp2.Body.Close()
		freshToken := loginResult2.Data.Token

		// Logout
		logoutResp, _ := bearerPost(ts.URL("/auth/logout"), freshToken, "")
		logoutResp.Body.Close()

		// Try to access /auth/me with the invalidated token
		resp, err := bearerGet(ts.URL("/auth/me"), freshToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 after logout, got %d", resp.StatusCode)
		}
	})
}

func TestAuthFullFlow(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("complete auth flow: register -> login -> me -> logout", func(t *testing.T) {
		// Step 1: Register
		registerBody := `{"email": "fullflow@example.com", "password": "password123", "name": "Full Flow User"}`
		regResp, err := anonPost(ts.URL("/auth/register"), registerBody)
		if err != nil {
			t.Fatalf("Failed to register: %v", err)
		}
		if regResp.StatusCode != http.StatusCreated {
			respBody, _ := io.ReadAll(regResp.Body)
			regResp.Body.Close()
			t.Fatalf("Expected register status 201, got %d: %s", regResp.StatusCode, respBody)
		}
		var regResult AuthUserEnvelope
		json.NewDecoder(regResp.Body).Decode(&regResult)
		regResp.Body.Close()
		userID := regResult.Data.ID

		// Step 2: Login
		loginBody := `{"email": "fullflow@example.com", "password": "password123"}`
		loginResp, err := anonPost(ts.URL("/auth/login"), loginBody)
		if err != nil {
			t.Fatalf("Failed to login: %v", err)
		}
		if loginResp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(loginResp.Body)
			loginResp.Body.Close()
			t.Fatalf("Expected login status 200, got %d: %s", loginResp.StatusCode, respBody)
		}
		var loginResult AuthLoginEnvelope
		json.NewDecoder(loginResp.Body).Decode(&loginResult)
		loginResp.Body.Close()
		token := loginResult.Data.Token

		// Step 3: Get current user
		meResp, err := bearerGet(ts.URL("/auth/me"), token)
		if err != nil {
			t.Fatalf("Failed to get current user: %v", err)
		}
		if meResp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(meResp.Body)
			meResp.Body.Close()
			t.Fatalf("Expected /auth/me status 200, got %d: %s", meResp.StatusCode, respBody)
		}
		var meResult AuthUserEnvelope
		json.NewDecoder(meResp.Body).Decode(&meResult)
		meResp.Body.Close()

		if meResult.Data.ID != userID {
			t.Errorf("Expected user ID %s, got %s", userID, meResult.Data.ID)
		}
		if meResult.Data.Email != "fullflow@example.com" {
			t.Errorf("Expected email 'fullflow@example.com', got '%s'", meResult.Data.Email)
		}

		// Step 4: Logout
		logoutResp, err := bearerPost(ts.URL("/auth/logout"), token, "")
		if err != nil {
			t.Fatalf("Failed to logout: %v", err)
		}
		if logoutResp.StatusCode != http.StatusNoContent {
			respBody, _ := io.ReadAll(logoutResp.Body)
			logoutResp.Body.Close()
			t.Fatalf("Expected logout status 204, got %d: %s", logoutResp.StatusCode, respBody)
		}
		logoutResp.Body.Close()

		// Step 5: Verify token is invalidated
		meResp2, err := bearerGet(ts.URL("/auth/me"), token)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer meResp2.Body.Close()

		if meResp2.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 after logout, got %d", meResp2.StatusCode)
		}
	})
}

func TestAuthPasswordNeverExposed(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Register
	registerBody := `{"email": "nohash@example.com", "password": "password123", "name": "No Hash User"}`
	regResp, _ := anonPost(ts.URL("/auth/register"), registerBody)
	regRespBody, _ := io.ReadAll(regResp.Body)
	regResp.Body.Close()

	// Check register response doesn't contain password_hash
	if bytes.Contains(regRespBody, []byte("password_hash")) || bytes.Contains(regRespBody, []byte("passwordHash")) {
		t.Error("Register response should not contain password_hash")
	}

	// Login
	loginBody := `{"email": "nohash@example.com", "password": "password123"}`
	loginResp, _ := anonPost(ts.URL("/auth/login"), loginBody)
	loginRespBody, _ := io.ReadAll(loginResp.Body)
	loginResp.Body.Close()

	// Check login response doesn't contain password_hash
	if bytes.Contains(loginRespBody, []byte("password_hash")) || bytes.Contains(loginRespBody, []byte("passwordHash")) {
		t.Error("Login response should not contain password_hash")
	}

	// Get token for /auth/me
	var loginResult AuthLoginEnvelope
	json.Unmarshal(loginRespBody, &loginResult)
	token := loginResult.Data.Token

	// Get current user
	meResp, _ := bearerGet(ts.URL("/auth/me"), token)
	meRespBody, _ := io.ReadAll(meResp.Body)
	meResp.Body.Close()

	// Check /auth/me response doesn't contain password_hash
	if bytes.Contains(meRespBody, []byte("password_hash")) || bytes.Contains(meRespBody, []byte("passwordHash")) {
		t.Error("/auth/me response should not contain password_hash")
	}
}
