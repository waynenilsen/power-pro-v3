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

// CycleTestResponse matches the API response format for a cycle.
type CycleTestResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	LengthWeeks int       `json:"lengthWeeks"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CycleWeekTestResponse matches the API response format for a week within a cycle.
type CycleWeekTestResponse struct {
	ID         string `json:"id"`
	WeekNumber int    `json:"weekNumber"`
}

// CycleWithWeeksTestResponse matches the API response format for a cycle with its weeks.
type CycleWithWeeksTestResponse struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	LengthWeeks int                     `json:"lengthWeeks"`
	Weeks       []CycleWeekTestResponse `json:"weeks"`
	CreatedAt   time.Time               `json:"createdAt"`
	UpdatedAt   time.Time               `json:"updatedAt"`
}

// CyclePaginationMeta contains pagination metadata.
type CyclePaginationMeta struct {
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasMore bool  `json:"hasMore"`
}

// PaginatedCyclesResponse is the paginated list response.
type PaginatedCyclesResponse struct {
	Data []CycleTestResponse  `json:"data"`
	Meta *CyclePaginationMeta `json:"meta"`
}

// authGetCycle performs an authenticated GET request
func authGetCycle(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestUserID)
	return http.DefaultClient.Do(req)
}

// adminPostCycle performs an admin-authenticated POST request
func adminPostCycle(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// adminPutCycle performs an admin-authenticated PUT request
func adminPutCycle(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// adminDeleteCycle performs an admin-authenticated DELETE request
func adminDeleteCycle(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

func TestCycleCRUD(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	var createdCycle CycleTestResponse

	t.Run("creates cycle with required fields", func(t *testing.T) {
		body := `{"name": "5/3/1 Cycle", "lengthWeeks": 4}`
		resp, err := adminPostCycle(ts.URL("/cycles"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		if err := json.NewDecoder(resp.Body).Decode(&createdCycle); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if createdCycle.ID == "" {
			t.Error("Expected non-empty ID")
		}
		if createdCycle.Name != "5/3/1 Cycle" {
			t.Errorf("Expected name '5/3/1 Cycle', got %s", createdCycle.Name)
		}
		if createdCycle.LengthWeeks != 4 {
			t.Errorf("Expected length_weeks 4, got %d", createdCycle.LengthWeeks)
		}
	})

	t.Run("creates cycle with 1 week", func(t *testing.T) {
		body := `{"name": "Starting Strength Cycle", "lengthWeeks": 1}`
		resp, err := adminPostCycle(ts.URL("/cycles"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var cycle CycleTestResponse
		if err := json.NewDecoder(resp.Body).Decode(&cycle); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if cycle.LengthWeeks != 1 {
			t.Errorf("Expected length_weeks 1, got %d", cycle.LengthWeeks)
		}
	})

	t.Run("gets cycle by ID", func(t *testing.T) {
		resp, err := authGetCycle(ts.URL("/cycles/" + createdCycle.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var cycle CycleWithWeeksTestResponse
		if err := json.NewDecoder(resp.Body).Decode(&cycle); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if cycle.ID != createdCycle.ID {
			t.Errorf("Expected ID %s, got %s", createdCycle.ID, cycle.ID)
		}
		if cycle.Name != "5/3/1 Cycle" {
			t.Errorf("Expected name '5/3/1 Cycle', got %s", cycle.Name)
		}
		// Initially no weeks
		if len(cycle.Weeks) != 0 {
			t.Errorf("Expected 0 weeks initially, got %d", len(cycle.Weeks))
		}
	})

	t.Run("returns 404 for non-existent cycle", func(t *testing.T) {
		resp, err := authGetCycle(ts.URL("/cycles/non-existent-id"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("lists cycles with pagination", func(t *testing.T) {
		resp, err := authGetCycle(ts.URL("/cycles"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result PaginatedCyclesResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Meta == nil || result.Meta.Total < 2 {
			total := int64(0)
			if result.Meta != nil {
				total = result.Meta.Total
			}
			t.Errorf("Expected at least 2 cycles, got %d", total)
		}
		if result.Meta == nil || result.Meta.Offset != 0 {
			offset := 0
			if result.Meta != nil {
				offset = result.Meta.Offset
			}
			t.Errorf("Expected offset 0, got %d", offset)
		}
	})

	t.Run("updates cycle name", func(t *testing.T) {
		body := `{"name": "Modified 5/3/1 Cycle"}`
		resp, err := adminPutCycle(ts.URL("/cycles/"+createdCycle.ID), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var updated CycleTestResponse
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if updated.Name != "Modified 5/3/1 Cycle" {
			t.Errorf("Expected name 'Modified 5/3/1 Cycle', got %s", updated.Name)
		}
	})

	t.Run("updates cycle length_weeks", func(t *testing.T) {
		body := `{"lengthWeeks": 3}`
		resp, err := adminPutCycle(ts.URL("/cycles/"+createdCycle.ID), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var updated CycleTestResponse
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if updated.LengthWeeks != 3 {
			t.Errorf("Expected length_weeks 3, got %d", updated.LengthWeeks)
		}
	})

	t.Run("deletes cycle", func(t *testing.T) {
		// Create a cycle to delete
		body := `{"name": "Cycle To Delete", "lengthWeeks": 2}`
		createResp, _ := adminPostCycle(ts.URL("/cycles"), body)
		var toDelete CycleTestResponse
		json.NewDecoder(createResp.Body).Decode(&toDelete)
		createResp.Body.Close()

		resp, err := adminDeleteCycle(ts.URL("/cycles/" + toDelete.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 204, got %d: %s", resp.StatusCode, bodyBytes)
		}

		// Verify it's deleted
		getResp, _ := authGetCycle(ts.URL("/cycles/" + toDelete.ID))
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404 after delete, got %d", getResp.StatusCode)
		}
	})
}

func TestCycleValidation(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("rejects empty name", func(t *testing.T) {
		body := `{"name": "", "lengthWeeks": 4}`
		resp, err := adminPostCycle(ts.URL("/cycles"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects whitespace-only name", func(t *testing.T) {
		body := `{"name": "   ", "lengthWeeks": 4}`
		resp, err := adminPostCycle(ts.URL("/cycles"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects invalid length_weeks (zero)", func(t *testing.T) {
		body := `{"name": "Test Cycle", "lengthWeeks": 0}`
		resp, err := adminPostCycle(ts.URL("/cycles"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects invalid length_weeks (negative)", func(t *testing.T) {
		body := `{"name": "Test Cycle", "lengthWeeks": -1}`
		resp, err := adminPostCycle(ts.URL("/cycles"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects missing length_weeks", func(t *testing.T) {
		body := `{"name": "Test Cycle"}`
		resp, err := adminPostCycle(ts.URL("/cycles"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Missing lengthWeeks defaults to 0, which is invalid
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 404 when updating non-existent cycle", func(t *testing.T) {
		body := `{"name": "New Name"}`
		resp, err := adminPutCycle(ts.URL("/cycles/non-existent-id"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 404 when deleting non-existent cycle", func(t *testing.T) {
		resp, err := adminDeleteCycle(ts.URL("/cycles/non-existent-id"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestCycleWithWeeks(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a cycle
	cycleBody := `{"name": "Week Test Cycle", "lengthWeeks": 4}`
	cycleResp, _ := adminPostCycle(ts.URL("/cycles"), cycleBody)
	var createdCycle CycleTestResponse
	json.NewDecoder(cycleResp.Body).Decode(&createdCycle)
	cycleResp.Body.Close()

	// Create weeks for the cycle
	for i := 1; i <= 4; i++ {
		weekBody := `{"weekNumber": ` + string(rune('0'+i)) + `, "cycleId": "` + createdCycle.ID + `"}`
		resp, _ := adminPostWeek(ts.URL("/weeks"), weekBody)
		resp.Body.Close()
	}

	t.Run("get cycle includes associated weeks", func(t *testing.T) {
		resp, err := authGetCycle(ts.URL("/cycles/" + createdCycle.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var cycle CycleWithWeeksTestResponse
		if err := json.NewDecoder(resp.Body).Decode(&cycle); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(cycle.Weeks) != 4 {
			t.Errorf("Expected 4 weeks, got %d", len(cycle.Weeks))
		}

		// Verify weeks are ordered by week number
		for i, w := range cycle.Weeks {
			expectedWeekNumber := i + 1
			if w.WeekNumber != expectedWeekNumber {
				t.Errorf("Expected week %d at position %d, got week %d", expectedWeekNumber, i, w.WeekNumber)
			}
		}
	})
}

func TestCycleAuthorization(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create cycle as admin
	cycleBody := `{"name": "Auth Test Cycle", "lengthWeeks": 4}`
	cycleResp, _ := adminPostCycle(ts.URL("/cycles"), cycleBody)
	var createdCycle CycleTestResponse
	json.NewDecoder(cycleResp.Body).Decode(&createdCycle)
	cycleResp.Body.Close()

	t.Run("unauthenticated user gets 401 on GET /cycles", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/cycles"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("unauthenticated user gets 401 on GET /cycles/{id}", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/cycles/" + createdCycle.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("authenticated user can GET /cycles", func(t *testing.T) {
		resp, err := authGetCycle(ts.URL("/cycles"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("non-admin user gets 403 on POST /cycles", func(t *testing.T) {
		body := `{"name": "Unauthorized Cycle", "lengthWeeks": 4}`
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/cycles"), bytes.NewBufferString(body))
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

	t.Run("non-admin user gets 403 on PUT /cycles/{id}", func(t *testing.T) {
		body := `{"name": "Modified Name"}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/cycles/"+createdCycle.ID), bytes.NewBufferString(body))
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

	t.Run("non-admin user gets 403 on DELETE /cycles/{id}", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL("/cycles/"+createdCycle.ID), nil)
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

func TestCycleResponseFormat(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create cycle
	cycleBody := `{"name": "Format Test Cycle", "lengthWeeks": 4}`
	cycleResp, _ := adminPostCycle(ts.URL("/cycles"), cycleBody)
	var createdCycle CycleTestResponse
	json.NewDecoder(cycleResp.Body).Decode(&createdCycle)
	cycleResp.Body.Close()

	t.Run("response has correct JSON field names", func(t *testing.T) {
		resp, _ := authGetCycle(ts.URL("/cycles/" + createdCycle.ID))
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// Check camelCase field names per ERD spec
		expectedFields := []string{
			`"id"`,
			`"name"`,
			`"lengthWeeks"`,
			`"weeks"`,
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

func TestCycleSorting(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create cycles with different names and lengths
	testCycles := []struct {
		name        string
		lengthWeeks int
	}{
		{"Charlie", 3},
		{"Alpha", 1},
		{"Bravo", 4},
	}

	for _, tc := range testCycles {
		body := `{"name": "` + tc.name + `", "lengthWeeks": ` + string(rune('0'+tc.lengthWeeks)) + `}`
		resp, _ := adminPostCycle(ts.URL("/cycles"), body)
		resp.Body.Close()
	}

	t.Run("sorts by name ascending by default", func(t *testing.T) {
		resp, err := authGetCycle(ts.URL("/cycles"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedCyclesResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) < 2 {
			t.Skip("Not enough cycles to test sorting")
		}

		for i := 1; i < len(result.Data); i++ {
			if result.Data[i].Name < result.Data[i-1].Name {
				t.Errorf("Cycles not sorted correctly by name: %s before %s",
					result.Data[i-1].Name, result.Data[i].Name)
			}
		}
	})

	t.Run("sorts by name descending", func(t *testing.T) {
		resp, err := authGetCycle(ts.URL("/cycles?sortBy=name&sortOrder=desc"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedCyclesResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) < 2 {
			t.Skip("Not enough cycles to test sorting")
		}

		for i := 1; i < len(result.Data); i++ {
			if result.Data[i].Name > result.Data[i-1].Name {
				t.Errorf("Cycles not sorted correctly (desc) by name: %s before %s",
					result.Data[i-1].Name, result.Data[i].Name)
			}
		}
	})

	t.Run("sorts by length_weeks ascending", func(t *testing.T) {
		resp, err := authGetCycle(ts.URL("/cycles?sortBy=length_weeks&sortOrder=asc"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedCyclesResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) < 2 {
			t.Skip("Not enough cycles to test sorting")
		}

		for i := 1; i < len(result.Data); i++ {
			if result.Data[i].LengthWeeks < result.Data[i-1].LengthWeeks {
				t.Errorf("Cycles not sorted correctly by length_weeks: %d before %d",
					result.Data[i-1].LengthWeeks, result.Data[i].LengthWeeks)
			}
		}
	})
}

func TestCyclePagination(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create multiple cycles
	for i := 1; i <= 25; i++ {
		body := `{"name": "Pagination Cycle ` + string(rune('A'+i-1)) + `", "lengthWeeks": 4}`
		resp, _ := adminPostCycle(ts.URL("/cycles"), body)
		resp.Body.Close()
	}

	t.Run("respects limit parameter", func(t *testing.T) {
		resp, err := authGetCycle(ts.URL("/cycles?limit=5"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedCyclesResponse
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
		resp, err := authGetCycle(ts.URL("/cycles?limit=5&offset=5"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedCyclesResponse
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
		resp, err := authGetCycle(ts.URL("/cycles?limit=10"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedCyclesResponse
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
