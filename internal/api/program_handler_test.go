package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// ProgramTestResponse matches the API response format for a program (list view).
type ProgramTestResponse struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Slug            string    `json:"slug"`
	Description     *string   `json:"description,omitempty"`
	CycleID         string    `json:"cycleId"`
	WeeklyLookupID  *string   `json:"weeklyLookupId,omitempty"`
	DailyLookupID   *string   `json:"dailyLookupId,omitempty"`
	DefaultRounding *float64  `json:"defaultRounding,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// ProgramCycleWeekTestResponse represents a week within a cycle.
type ProgramCycleWeekTestResponse struct {
	ID         string `json:"id"`
	WeekNumber int    `json:"weekNumber"`
}

// ProgramCycleTestResponse represents embedded cycle info in a program response.
type ProgramCycleTestResponse struct {
	ID          string                         `json:"id"`
	Name        string                         `json:"name"`
	LengthWeeks int                            `json:"lengthWeeks"`
	Weeks       []ProgramCycleWeekTestResponse `json:"weeks"`
}

// LookupReferenceTestResponse represents a lookup table reference.
type LookupReferenceTestResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ProgramDetailTestResponse represents the API response format for a program (detail view with embedded cycle).
type ProgramDetailTestResponse struct {
	ID              string                       `json:"id"`
	Name            string                       `json:"name"`
	Slug            string                       `json:"slug"`
	Description     *string                      `json:"description,omitempty"`
	Cycle           *ProgramCycleTestResponse    `json:"cycle"`
	WeeklyLookup    *LookupReferenceTestResponse `json:"weeklyLookup,omitempty"`
	DailyLookup     *LookupReferenceTestResponse `json:"dailyLookup,omitempty"`
	DefaultRounding *float64                     `json:"defaultRounding,omitempty"`
	CreatedAt       time.Time                    `json:"createdAt"`
	UpdatedAt       time.Time                    `json:"updatedAt"`
}

// ProgramPaginationMeta contains pagination metadata.
type ProgramPaginationMeta struct {
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasMore bool  `json:"hasMore"`
}

// PaginatedProgramsResponse is the paginated list response.
type PaginatedProgramsResponse struct {
	Data []ProgramTestResponse  `json:"data"`
	Meta *ProgramPaginationMeta `json:"meta"`
}

// authGetProgram performs an authenticated GET request
func authGetProgram(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestUserID)
	return http.DefaultClient.Do(req)
}

// adminPostProgram performs an admin-authenticated POST request
func adminPostProgram(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// adminPutProgram performs an admin-authenticated PUT request
func adminPutProgram(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// adminDeleteProgram performs an admin-authenticated DELETE request
func adminDeleteProgram(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// createProgramTestCycle creates a test cycle and returns its ID
func createProgramTestCycle(t *testing.T, ts *testutil.TestServer, name string) string {
	body := `{"name": "` + name + `", "lengthWeeks": 4}`
	resp, err := adminPostCycle(ts.URL("/cycles"), body)
	if err != nil {
		t.Fatalf("Failed to create test cycle: %v", err)
	}
	defer resp.Body.Close()

	var cycle CycleTestResponse
	json.NewDecoder(resp.Body).Decode(&cycle)
	return cycle.ID
}

func TestProgramCRUD(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a cycle first (required for programs)
	cycleID := createProgramTestCycle(t, ts, "Test Cycle for Programs")

	var createdProgram ProgramTestResponse

	t.Run("creates program with required fields", func(t *testing.T) {
		body := `{"name": "5/3/1 BBB", "slug": "531-bbb", "cycleId": "` + cycleID + `"}`
		resp, err := adminPostProgram(ts.URL("/programs"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		if err := json.NewDecoder(resp.Body).Decode(&createdProgram); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if createdProgram.ID == "" {
			t.Error("Expected non-empty ID")
		}
		if createdProgram.Name != "5/3/1 BBB" {
			t.Errorf("Expected name '5/3/1 BBB', got %s", createdProgram.Name)
		}
		if createdProgram.Slug != "531-bbb" {
			t.Errorf("Expected slug '531-bbb', got %s", createdProgram.Slug)
		}
		if createdProgram.CycleID != cycleID {
			t.Errorf("Expected cycleId %s, got %s", cycleID, createdProgram.CycleID)
		}
	})

	t.Run("creates program with all fields", func(t *testing.T) {
		body := `{
			"name": "Starting Strength",
			"slug": "starting-strength",
			"description": "Classic beginner program",
			"cycleId": "` + cycleID + `",
			"defaultRounding": 2.5
		}`
		resp, err := adminPostProgram(ts.URL("/programs"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var program ProgramTestResponse
		if err := json.NewDecoder(resp.Body).Decode(&program); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if program.Description == nil || *program.Description != "Classic beginner program" {
			t.Errorf("Expected description 'Classic beginner program', got %v", program.Description)
		}
		if program.DefaultRounding == nil || *program.DefaultRounding != 2.5 {
			t.Errorf("Expected defaultRounding 2.5, got %v", program.DefaultRounding)
		}
	})

	t.Run("gets program by ID with embedded cycle", func(t *testing.T) {
		resp, err := authGetProgram(ts.URL("/programs/" + createdProgram.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var program ProgramDetailTestResponse
		if err := json.NewDecoder(resp.Body).Decode(&program); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if program.ID != createdProgram.ID {
			t.Errorf("Expected ID %s, got %s", createdProgram.ID, program.ID)
		}
		if program.Name != "5/3/1 BBB" {
			t.Errorf("Expected name '5/3/1 BBB', got %s", program.Name)
		}
		// Verify cycle is embedded
		if program.Cycle == nil {
			t.Error("Expected embedded cycle, got nil")
		} else {
			if program.Cycle.ID != cycleID {
				t.Errorf("Expected cycle ID %s, got %s", cycleID, program.Cycle.ID)
			}
			if program.Cycle.Name != "Test Cycle for Programs" {
				t.Errorf("Expected cycle name 'Test Cycle for Programs', got %s", program.Cycle.Name)
			}
		}
	})

	t.Run("returns 404 for non-existent program", func(t *testing.T) {
		resp, err := authGetProgram(ts.URL("/programs/non-existent-id"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("lists programs with pagination", func(t *testing.T) {
		resp, err := authGetProgram(ts.URL("/programs"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result PaginatedProgramsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Meta == nil || result.Meta.Total < 2 {
			total := int64(0)
			if result.Meta != nil {
				total = result.Meta.Total
			}
			t.Errorf("Expected at least 2 programs, got %d", total)
		}
		if result.Meta == nil || result.Meta.Offset != 0 {
			offset := 0
			if result.Meta != nil {
				offset = result.Meta.Offset
			}
			t.Errorf("Expected offset 0, got %d", offset)
		}
	})

	t.Run("updates program name", func(t *testing.T) {
		body := `{"name": "Modified 5/3/1 BBB"}`
		resp, err := adminPutProgram(ts.URL("/programs/"+createdProgram.ID), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var updated ProgramTestResponse
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if updated.Name != "Modified 5/3/1 BBB" {
			t.Errorf("Expected name 'Modified 5/3/1 BBB', got %s", updated.Name)
		}
	})

	t.Run("updates program slug", func(t *testing.T) {
		body := `{"slug": "modified-531-bbb"}`
		resp, err := adminPutProgram(ts.URL("/programs/"+createdProgram.ID), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var updated ProgramTestResponse
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if updated.Slug != "modified-531-bbb" {
			t.Errorf("Expected slug 'modified-531-bbb', got %s", updated.Slug)
		}
	})

	t.Run("deletes program", func(t *testing.T) {
		// Create a program to delete
		body := `{"name": "Program To Delete", "slug": "program-to-delete", "cycleId": "` + cycleID + `"}`
		createResp, _ := adminPostProgram(ts.URL("/programs"), body)
		var toDelete ProgramTestResponse
		json.NewDecoder(createResp.Body).Decode(&toDelete)
		createResp.Body.Close()

		resp, err := adminDeleteProgram(ts.URL("/programs/" + toDelete.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 204, got %d: %s", resp.StatusCode, bodyBytes)
		}

		// Verify it's deleted
		getResp, _ := authGetProgram(ts.URL("/programs/" + toDelete.ID))
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404 after delete, got %d", getResp.StatusCode)
		}
	})
}

func TestProgramValidation(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	cycleID := createProgramTestCycle(t, ts, "Validation Test Cycle")

	t.Run("rejects empty name", func(t *testing.T) {
		body := `{"name": "", "slug": "test-slug", "cycleId": "` + cycleID + `"}`
		resp, err := adminPostProgram(ts.URL("/programs"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects empty slug", func(t *testing.T) {
		body := `{"name": "Test Program", "slug": "", "cycleId": "` + cycleID + `"}`
		resp, err := adminPostProgram(ts.URL("/programs"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects invalid slug format", func(t *testing.T) {
		body := `{"name": "Test Program", "slug": "Invalid Slug!", "cycleId": "` + cycleID + `"}`
		resp, err := adminPostProgram(ts.URL("/programs"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects empty cycleId", func(t *testing.T) {
		body := `{"name": "Test Program", "slug": "test-slug", "cycleId": ""}`
		resp, err := adminPostProgram(ts.URL("/programs"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects non-existent cycleId", func(t *testing.T) {
		body := `{"name": "Test Program", "slug": "test-slug", "cycleId": "non-existent-cycle"}`
		resp, err := adminPostProgram(ts.URL("/programs"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects invalid defaultRounding (negative)", func(t *testing.T) {
		body := `{"name": "Test Program", "slug": "test-slug-neg", "cycleId": "` + cycleID + `", "defaultRounding": -1}`
		resp, err := adminPostProgram(ts.URL("/programs"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects invalid defaultRounding (zero)", func(t *testing.T) {
		body := `{"name": "Test Program", "slug": "test-slug-zero", "cycleId": "` + cycleID + `", "defaultRounding": 0}`
		resp, err := adminPostProgram(ts.URL("/programs"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 404 when updating non-existent program", func(t *testing.T) {
		body := `{"name": "New Name"}`
		resp, err := adminPutProgram(ts.URL("/programs/non-existent-id"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 404 when deleting non-existent program", func(t *testing.T) {
		resp, err := adminDeleteProgram(ts.URL("/programs/non-existent-id"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestProgramSlugUniqueness(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	cycleID := createProgramTestCycle(t, ts, "Slug Uniqueness Test Cycle")

	// Create first program
	body := `{"name": "First Program", "slug": "unique-slug", "cycleId": "` + cycleID + `"}`
	resp, _ := adminPostProgram(ts.URL("/programs"), body)
	var firstProgram ProgramTestResponse
	json.NewDecoder(resp.Body).Decode(&firstProgram)
	resp.Body.Close()

	t.Run("rejects duplicate slug on create", func(t *testing.T) {
		body := `{"name": "Second Program", "slug": "unique-slug", "cycleId": "` + cycleID + `"}`
		resp, err := adminPostProgram(ts.URL("/programs"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 409, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("rejects duplicate slug on update", func(t *testing.T) {
		// Create another program with different slug
		body := `{"name": "Third Program", "slug": "another-slug", "cycleId": "` + cycleID + `"}`
		createResp, _ := adminPostProgram(ts.URL("/programs"), body)
		var secondProgram ProgramTestResponse
		json.NewDecoder(createResp.Body).Decode(&secondProgram)
		createResp.Body.Close()

		// Try to update to use the first program's slug
		updateBody := `{"slug": "unique-slug"}`
		resp, err := adminPutProgram(ts.URL("/programs/"+secondProgram.ID), updateBody)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 409, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("allows updating to same slug on same program", func(t *testing.T) {
		body := `{"slug": "unique-slug"}`
		resp, err := adminPutProgram(ts.URL("/programs/"+firstProgram.ID), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})
}

func TestProgramAuthorization(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	cycleID := createProgramTestCycle(t, ts, "Auth Test Cycle")

	// Create program as admin
	body := `{"name": "Auth Test Program", "slug": "auth-test-program", "cycleId": "` + cycleID + `"}`
	createResp, _ := adminPostProgram(ts.URL("/programs"), body)
	var createdProgram ProgramTestResponse
	json.NewDecoder(createResp.Body).Decode(&createdProgram)
	createResp.Body.Close()

	t.Run("unauthenticated user gets 401 on GET /programs", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/programs"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("unauthenticated user gets 401 on GET /programs/{id}", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/programs/" + createdProgram.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("authenticated user can GET /programs", func(t *testing.T) {
		resp, err := authGetProgram(ts.URL("/programs"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("non-admin user gets 403 on POST /programs", func(t *testing.T) {
		body := `{"name": "Unauthorized Program", "slug": "unauthorized-program", "cycleId": "` + cycleID + `"}`
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/programs"), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", testutil.TestUserID)
		// Not setting X-Admin

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("non-admin user gets 403 on PUT /programs/{id}", func(t *testing.T) {
		body := `{"name": "Modified Name"}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/programs/"+createdProgram.ID), bytes.NewBufferString(body))
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

	t.Run("non-admin user gets 403 on DELETE /programs/{id}", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL("/programs/"+createdProgram.ID), nil)
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
}

func TestProgramResponseFormat(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	cycleID := createProgramTestCycle(t, ts, "Response Format Test Cycle")

	body := `{"name": "Format Test Program", "slug": "format-test-program", "cycleId": "` + cycleID + `"}`
	createResp, _ := adminPostProgram(ts.URL("/programs"), body)
	var createdProgram ProgramTestResponse
	json.NewDecoder(createResp.Body).Decode(&createdProgram)
	createResp.Body.Close()

	t.Run("response has correct JSON field names", func(t *testing.T) {
		resp, _ := authGetProgram(ts.URL("/programs/" + createdProgram.ID))
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// Check camelCase field names per ERD spec
		expectedFields := []string{
			`"id"`,
			`"name"`,
			`"slug"`,
			`"cycle"`,
			`"createdAt"`,
			`"updatedAt"`,
		}

		for _, field := range expectedFields {
			if !bytes.Contains(body, []byte(field)) {
				t.Errorf("Expected field %s in response, body: %s", field, bodyStr)
			}
		}
	})

	t.Run("list response has correct pagination fields", func(t *testing.T) {
		resp, _ := authGetProgram(ts.URL("/programs"))
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		expectedFields := []string{
			`"data"`,
			`"meta"`,
			`"total"`,
			`"limit"`,
			`"offset"`,
			`"hasMore"`,
		}

		for _, field := range expectedFields {
			if !bytes.Contains(body, []byte(field)) {
				t.Errorf("Expected field %s in response, body: %s", field, bodyStr)
			}
		}
	})
}

func TestProgramSorting(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	cycleID := createProgramTestCycle(t, ts, "Sorting Test Cycle")

	// Create programs with different names
	testPrograms := []string{"zeta-program", "alpha-program", "beta-program"}
	for _, slug := range testPrograms {
		body := `{"name": "` + slug + `", "slug": "` + slug + `", "cycleId": "` + cycleID + `"}`
		resp, _ := adminPostProgram(ts.URL("/programs"), body)
		resp.Body.Close()
	}

	t.Run("sorts by name ascending by default", func(t *testing.T) {
		resp, err := authGetProgram(ts.URL("/programs"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedProgramsResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) < 2 {
			t.Skip("Not enough programs to test sorting")
		}

		for i := 1; i < len(result.Data); i++ {
			if result.Data[i].Name < result.Data[i-1].Name {
				t.Errorf("Programs not sorted correctly by name: %s before %s",
					result.Data[i-1].Name, result.Data[i].Name)
			}
		}
	})

	t.Run("sorts by name descending", func(t *testing.T) {
		resp, err := authGetProgram(ts.URL("/programs?sortBy=name&sortOrder=desc"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedProgramsResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) < 2 {
			t.Skip("Not enough programs to test sorting")
		}

		for i := 1; i < len(result.Data); i++ {
			if result.Data[i].Name > result.Data[i-1].Name {
				t.Errorf("Programs not sorted correctly (desc) by name: %s before %s",
					result.Data[i-1].Name, result.Data[i].Name)
			}
		}
	})
}

func TestProgramPagination(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	cycleID := createProgramTestCycle(t, ts, "Pagination Test Cycle")

	// Create multiple programs
	for i := 1; i <= 15; i++ {
		slug := "pagination-program-" + string(rune('a'+i-1))
		body := `{"name": "Pagination Program ` + string(rune('A'+i-1)) + `", "slug": "` + slug + `", "cycleId": "` + cycleID + `"}`
		resp, _ := adminPostProgram(ts.URL("/programs"), body)
		resp.Body.Close()
	}

	t.Run("respects limit parameter", func(t *testing.T) {
		resp, err := authGetProgram(ts.URL("/programs?limit=5"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedProgramsResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) > 5 {
			t.Errorf("Expected at most 5 items, got %d", len(result.Data))
		}
		if result.Meta == nil || result.Meta.Limit != 5 {
			limit := 0
			if result.Meta != nil {
				limit = result.Meta.Limit
			}
			t.Errorf("Expected limit 5, got %d", limit)
		}
	})

	t.Run("respects offset parameter", func(t *testing.T) {
		resp, err := authGetProgram(ts.URL("/programs?limit=5&offset=5"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedProgramsResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if result.Meta == nil || result.Meta.Offset != 5 {
			offset := 0
			if result.Meta != nil {
				offset = result.Meta.Offset
			}
			t.Errorf("Expected offset 5, got %d", offset)
		}
	})

	t.Run("returns hasMore correctly", func(t *testing.T) {
		resp, err := authGetProgram(ts.URL("/programs?limit=10"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedProgramsResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if result.Meta == nil {
			t.Fatal("Expected meta to be present")
		}
		// hasMore should be true if there are more items than offset + limit
		expectedHasMore := result.Meta.Total > int64(result.Meta.Offset+len(result.Data))
		if result.Meta.HasMore != expectedHasMore {
			t.Errorf("Expected hasMore %v, got %v", expectedHasMore, result.Meta.HasMore)
		}
	})
}
