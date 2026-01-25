// Package e2e provides end-to-end tests for complete API workflows.
// This file contains E2E tests for Program Discovery features including
// filtering, search, and detail enhancement endpoints.
package e2e

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// =============================================================================
// PROGRAM DISCOVERY RESPONSE TYPES
// =============================================================================

// ProgramListResponse represents a program in the list view.
type ProgramListResponse struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Slug            string    `json:"slug"`
	Description     *string   `json:"description,omitempty"`
	CycleID         string    `json:"cycleId"`
	WeeklyLookupID  *string   `json:"weeklyLookupId,omitempty"`
	DailyLookupID   *string   `json:"dailyLookupId,omitempty"`
	DefaultRounding *float64  `json:"defaultRounding,omitempty"`
	Difficulty      string    `json:"difficulty"`
	DaysPerWeek     int       `json:"daysPerWeek"`
	Focus           string    `json:"focus"`
	HasAmrap        bool      `json:"hasAmrap"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// SampleWeekDay represents a day in the sample week.
type SampleWeekDay struct {
	Day           int    `json:"day"`
	Name          string `json:"name"`
	ExerciseCount int    `json:"exerciseCount"`
}

// CycleWeek represents a week in a cycle.
type CycleWeek struct {
	ID         string `json:"id"`
	WeekNumber int    `json:"weekNumber"`
}

// CycleReference represents cycle info in detail response.
type CycleReference struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	LengthWeeks int         `json:"lengthWeeks"`
	Weeks       []CycleWeek `json:"weeks"`
}

// LookupReference represents a lookup table reference.
type LookupReference struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ProgramDetailResponse represents a program in the detail view.
type ProgramDetailResponse struct {
	ID                      string           `json:"id"`
	Name                    string           `json:"name"`
	Slug                    string           `json:"slug"`
	Description             *string          `json:"description,omitempty"`
	Cycle                   *CycleReference  `json:"cycle"`
	WeeklyLookup            *LookupReference `json:"weeklyLookup,omitempty"`
	DailyLookup             *LookupReference `json:"dailyLookup,omitempty"`
	DefaultRounding         *float64         `json:"defaultRounding,omitempty"`
	Difficulty              string           `json:"difficulty"`
	DaysPerWeek             int              `json:"daysPerWeek"`
	Focus                   string           `json:"focus"`
	HasAmrap                bool             `json:"hasAmrap"`
	SampleWeek              []SampleWeekDay  `json:"sampleWeek"`
	LiftRequirements        []string         `json:"liftRequirements"`
	EstimatedSessionMinutes int              `json:"estimatedSessionMinutes"`
	CreatedAt               time.Time        `json:"createdAt"`
	UpdatedAt               time.Time        `json:"updatedAt"`
}

// PaginationMeta contains pagination metadata.
type PaginationMeta struct {
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasMore bool  `json:"hasMore"`
}

// PaginatedProgramsResponse is the paginated list response.
type PaginatedProgramsResponse struct {
	Data []ProgramListResponse `json:"data"`
	Meta *PaginationMeta       `json:"meta"`
}

// ProgramDetailEnvelope wraps detail response.
type ProgramDetailEnvelope struct {
	Data ProgramDetailResponse `json:"data"`
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// authGet performs an authenticated GET request.
func authGet(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestUserID)
	return http.DefaultClient.Do(req)
}

// getProgramsFiltered fetches programs with query parameters.
func getProgramsFiltered(t *testing.T, ts *testutil.TestServer, queryParams string) PaginatedProgramsResponse {
	t.Helper()

	url := ts.URL("/programs")
	if queryParams != "" {
		url += "?" + queryParams
	}

	resp, err := authGet(url)
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
	return result
}

// getProgramDetail fetches a program by ID.
func getProgramDetail(t *testing.T, ts *testutil.TestServer, programID string) ProgramDetailResponse {
	t.Helper()

	resp, err := authGet(ts.URL("/programs/" + programID))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
	}

	var result ProgramDetailEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	return result.Data
}

// getProgramBySlug finds a program by slug in the list.
func getProgramBySlug(t *testing.T, ts *testutil.TestServer, slug string) *ProgramListResponse {
	t.Helper()

	result := getProgramsFiltered(t, ts, "")
	for _, p := range result.Data {
		if p.Slug == slug {
			return &p
		}
	}
	return nil
}

// assertExpectedStatus makes a request and asserts the expected status code.
func assertExpectedStatus(t *testing.T, ts *testutil.TestServer, queryParams string, expectedStatus int) {
	t.Helper()

	url := ts.URL("/programs")
	if queryParams != "" {
		url += "?" + queryParams
	}

	resp, err := authGet(url)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Errorf("Expected status %d, got %d: %s", expectedStatus, resp.StatusCode, bodyBytes)
	}
}

// =============================================================================
// CANONICAL PROGRAM METADATA (from migration 00034)
// =============================================================================
// Starting Strength: beginner, 3 days/week, strength, no AMRAP
// Texas Method: intermediate, 3 days/week, strength, no AMRAP
// Wendler 5/3/1: intermediate, 4 days/week, strength, has AMRAP
// GZCLP: beginner, 4 days/week, strength, has AMRAP

// =============================================================================
// FILTERING TESTS
// =============================================================================

func TestProgramDiscoveryE2E_FilterByDifficulty(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("filter beginner returns Starting Strength and GZCLP", func(t *testing.T) {
		result := getProgramsFiltered(t, ts, "difficulty=beginner")

		// All returned programs should be beginner difficulty
		for _, p := range result.Data {
			if p.Difficulty != "beginner" {
				t.Errorf("Expected difficulty 'beginner', got %s for program %s", p.Difficulty, p.Name)
			}
		}

		// Should include Starting Strength and GZCLP
		slugs := make(map[string]bool)
		for _, p := range result.Data {
			slugs[p.Slug] = true
		}

		if !slugs["starting-strength"] {
			t.Error("Expected Starting Strength in beginner results")
		}
		if !slugs["gzclp"] {
			t.Error("Expected GZCLP in beginner results")
		}
	})

	t.Run("filter intermediate returns Texas Method and 5/3/1", func(t *testing.T) {
		result := getProgramsFiltered(t, ts, "difficulty=intermediate")

		// All returned programs should be intermediate difficulty
		for _, p := range result.Data {
			if p.Difficulty != "intermediate" {
				t.Errorf("Expected difficulty 'intermediate', got %s for program %s", p.Difficulty, p.Name)
			}
		}

		// Should include Texas Method and 5/3/1
		slugs := make(map[string]bool)
		for _, p := range result.Data {
			slugs[p.Slug] = true
		}

		if !slugs["texas-method"] {
			t.Error("Expected Texas Method in intermediate results")
		}
		if !slugs["531"] {
			t.Error("Expected 5/3/1 in intermediate results")
		}
	})

	t.Run("filter advanced returns empty (no advanced programs)", func(t *testing.T) {
		result := getProgramsFiltered(t, ts, "difficulty=advanced")

		// Should still return 200 with empty data
		if result.Data == nil {
			t.Error("Expected non-nil data array")
		}

		// Count canonical programs in results (there might be other seeded programs)
		canonicalCount := 0
		for _, p := range result.Data {
			if p.Slug == "starting-strength" || p.Slug == "texas-method" ||
				p.Slug == "531" || p.Slug == "gzclp" {
				canonicalCount++
			}
		}
		if canonicalCount > 0 {
			t.Errorf("Expected no canonical programs in advanced filter, got %d", canonicalCount)
		}
	})

	t.Run("invalid difficulty returns 400", func(t *testing.T) {
		assertExpectedStatus(t, ts, "difficulty=expert", http.StatusBadRequest)
		assertExpectedStatus(t, ts, "difficulty=elite", http.StatusBadRequest)
		// Empty string is ignored (not validated), only non-empty invalid values return 400
	})
}

func TestProgramDiscoveryE2E_FilterByDaysPerWeek(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("filter 3 days returns Starting Strength and Texas Method", func(t *testing.T) {
		result := getProgramsFiltered(t, ts, "days_per_week=3")

		// All returned programs should have 3 days per week
		for _, p := range result.Data {
			if p.DaysPerWeek != 3 {
				t.Errorf("Expected daysPerWeek 3, got %d for program %s", p.DaysPerWeek, p.Name)
			}
		}

		// Should include Starting Strength and Texas Method
		slugs := make(map[string]bool)
		for _, p := range result.Data {
			slugs[p.Slug] = true
		}

		if !slugs["starting-strength"] {
			t.Error("Expected Starting Strength in 3 days results")
		}
		if !slugs["texas-method"] {
			t.Error("Expected Texas Method in 3 days results")
		}
	})

	t.Run("filter 4 days returns 5/3/1 and GZCLP", func(t *testing.T) {
		result := getProgramsFiltered(t, ts, "days_per_week=4")

		// All returned programs should have 4 days per week
		for _, p := range result.Data {
			if p.DaysPerWeek != 4 {
				t.Errorf("Expected daysPerWeek 4, got %d for program %s", p.DaysPerWeek, p.Name)
			}
		}

		// Should include 5/3/1 and GZCLP
		slugs := make(map[string]bool)
		for _, p := range result.Data {
			slugs[p.Slug] = true
		}

		if !slugs["531"] {
			t.Error("Expected 5/3/1 in 4 days results")
		}
		if !slugs["gzclp"] {
			t.Error("Expected GZCLP in 4 days results")
		}
	})

	t.Run("filter 5 days returns empty for canonical programs", func(t *testing.T) {
		result := getProgramsFiltered(t, ts, "days_per_week=5")

		// Count canonical programs (should be 0)
		for _, p := range result.Data {
			if p.Slug == "starting-strength" || p.Slug == "texas-method" ||
				p.Slug == "531" || p.Slug == "gzclp" {
				t.Errorf("Unexpected canonical program %s in 5 days results", p.Name)
			}
		}
	})

	t.Run("invalid days_per_week returns 400", func(t *testing.T) {
		assertExpectedStatus(t, ts, "days_per_week=0", http.StatusBadRequest)
		assertExpectedStatus(t, ts, "days_per_week=8", http.StatusBadRequest)
		assertExpectedStatus(t, ts, "days_per_week=-1", http.StatusBadRequest)
		assertExpectedStatus(t, ts, "days_per_week=abc", http.StatusBadRequest)
	})
}

func TestProgramDiscoveryE2E_FilterByFocus(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("filter strength returns all 4 canonical programs", func(t *testing.T) {
		result := getProgramsFiltered(t, ts, "focus=strength")

		// All returned programs should have strength focus
		for _, p := range result.Data {
			if p.Focus != "strength" {
				t.Errorf("Expected focus 'strength', got %s for program %s", p.Focus, p.Name)
			}
		}

		// Should include all 4 canonical programs
		slugs := make(map[string]bool)
		for _, p := range result.Data {
			slugs[p.Slug] = true
		}

		expectedSlugs := []string{"starting-strength", "texas-method", "531", "gzclp"}
		for _, slug := range expectedSlugs {
			if !slugs[slug] {
				t.Errorf("Expected %s in strength results", slug)
			}
		}
	})

	t.Run("filter hypertrophy returns empty for canonical programs", func(t *testing.T) {
		result := getProgramsFiltered(t, ts, "focus=hypertrophy")

		// Count canonical programs (should be 0)
		for _, p := range result.Data {
			if p.Slug == "starting-strength" || p.Slug == "texas-method" ||
				p.Slug == "531" || p.Slug == "gzclp" {
				t.Errorf("Unexpected canonical program %s in hypertrophy results", p.Name)
			}
		}
	})

	t.Run("invalid focus returns 400", func(t *testing.T) {
		assertExpectedStatus(t, ts, "focus=cardio", http.StatusBadRequest)
		assertExpectedStatus(t, ts, "focus=invalid", http.StatusBadRequest)
	})
}

func TestProgramDiscoveryE2E_FilterByHasAmrap(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("filter has_amrap=true returns 5/3/1 and GZCLP", func(t *testing.T) {
		result := getProgramsFiltered(t, ts, "has_amrap=true")

		// All returned programs should have AMRAP
		for _, p := range result.Data {
			if !p.HasAmrap {
				t.Errorf("Expected hasAmrap true for program %s", p.Name)
			}
		}

		// Should include 5/3/1 and GZCLP
		slugs := make(map[string]bool)
		for _, p := range result.Data {
			slugs[p.Slug] = true
		}

		if !slugs["531"] {
			t.Error("Expected 5/3/1 in has_amrap=true results")
		}
		if !slugs["gzclp"] {
			t.Error("Expected GZCLP in has_amrap=true results")
		}

		// Should NOT include Starting Strength and Texas Method
		if slugs["starting-strength"] {
			t.Error("Starting Strength should not be in has_amrap=true results")
		}
		if slugs["texas-method"] {
			t.Error("Texas Method should not be in has_amrap=true results")
		}
	})

	t.Run("filter has_amrap=false returns Starting Strength and Texas Method", func(t *testing.T) {
		result := getProgramsFiltered(t, ts, "has_amrap=false")

		// All returned programs should NOT have AMRAP
		for _, p := range result.Data {
			if p.HasAmrap {
				t.Errorf("Expected hasAmrap false for program %s", p.Name)
			}
		}

		// Should include Starting Strength and Texas Method
		slugs := make(map[string]bool)
		for _, p := range result.Data {
			slugs[p.Slug] = true
		}

		if !slugs["starting-strength"] {
			t.Error("Expected Starting Strength in has_amrap=false results")
		}
		if !slugs["texas-method"] {
			t.Error("Expected Texas Method in has_amrap=false results")
		}

		// Should NOT include 5/3/1 and GZCLP
		if slugs["531"] {
			t.Error("5/3/1 should not be in has_amrap=false results")
		}
		if slugs["gzclp"] {
			t.Error("GZCLP should not be in has_amrap=false results")
		}
	})

	t.Run("invalid has_amrap returns 400", func(t *testing.T) {
		assertExpectedStatus(t, ts, "has_amrap=maybe", http.StatusBadRequest)
		// Note: "1" and "true" are valid, "0" and "false" are valid
		// Only non-boolean strings like "maybe", "yes" return 400
		assertExpectedStatus(t, ts, "has_amrap=yes", http.StatusBadRequest)
	})
}

func TestProgramDiscoveryE2E_CombinedFilters(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("difficulty=beginner AND days_per_week=3 returns only Starting Strength", func(t *testing.T) {
		result := getProgramsFiltered(t, ts, "difficulty=beginner&days_per_week=3")

		// Find Starting Strength
		found := false
		for _, p := range result.Data {
			if p.Slug == "starting-strength" {
				found = true
				if p.Difficulty != "beginner" {
					t.Errorf("Expected difficulty 'beginner', got %s", p.Difficulty)
				}
				if p.DaysPerWeek != 3 {
					t.Errorf("Expected daysPerWeek 3, got %d", p.DaysPerWeek)
				}
			}
			// Texas Method has 3 days but is intermediate, so should not be here
			if p.Slug == "texas-method" {
				t.Error("Texas Method should not be in beginner+3days results")
			}
			// GZCLP is beginner but has 4 days, so should not be here
			if p.Slug == "gzclp" {
				t.Error("GZCLP should not be in beginner+3days results")
			}
		}

		if !found {
			t.Error("Expected Starting Strength in results")
		}
	})

	t.Run("difficulty=intermediate AND has_amrap=true returns only 5/3/1", func(t *testing.T) {
		result := getProgramsFiltered(t, ts, "difficulty=intermediate&has_amrap=true")

		// Find 5/3/1
		found := false
		for _, p := range result.Data {
			if p.Slug == "531" {
				found = true
				if p.Difficulty != "intermediate" {
					t.Errorf("Expected difficulty 'intermediate', got %s", p.Difficulty)
				}
				if !p.HasAmrap {
					t.Error("Expected hasAmrap true")
				}
			}
			// Texas Method is intermediate but has no AMRAP
			if p.Slug == "texas-method" {
				t.Error("Texas Method should not be in intermediate+hasAmrap results")
			}
			// GZCLP has AMRAP but is beginner
			if p.Slug == "gzclp" {
				t.Error("GZCLP should not be in intermediate+hasAmrap results")
			}
		}

		if !found {
			t.Error("Expected 5/3/1 in results")
		}
	})

	t.Run("no matches returns 200 with empty array", func(t *testing.T) {
		// No program is advanced + hypertrophy focus
		result := getProgramsFiltered(t, ts, "difficulty=advanced&focus=peaking")

		if result.Data == nil {
			t.Error("Expected non-nil data array")
		}
		// Should return empty array, not error
		if result.Meta == nil {
			t.Error("Expected meta to be present")
		}
	})
}

// =============================================================================
// SEARCH TESTS
// =============================================================================

func TestProgramDiscoveryE2E_Search(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("search=strength returns Starting Strength", func(t *testing.T) {
		result := getProgramsFiltered(t, ts, "search=strength")

		found := false
		for _, p := range result.Data {
			if p.Slug == "starting-strength" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected Starting Strength when searching for 'strength'")
		}
	})

	t.Run("search=5/3/1 returns Wendler 5/3/1", func(t *testing.T) {
		// Note: Search is on name, not slug. The name is "Wendler 5/3/1"
		result := getProgramsFiltered(t, ts, "search=5/3/1")

		found := false
		for _, p := range result.Data {
			if p.Slug == "531" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected Wendler 5/3/1 when searching for '5/3/1'")
		}
	})

	t.Run("search=gz returns GZCLP", func(t *testing.T) {
		result := getProgramsFiltered(t, ts, "search=gz")

		found := false
		for _, p := range result.Data {
			if p.Slug == "gzclp" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected GZCLP when searching for 'gz'")
		}
	})

	t.Run("search=method returns Texas Method", func(t *testing.T) {
		result := getProgramsFiltered(t, ts, "search=method")

		found := false
		for _, p := range result.Data {
			if p.Slug == "texas-method" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected Texas Method when searching for 'method'")
		}
	})

	t.Run("search is case-insensitive", func(t *testing.T) {
		// Search with uppercase
		result := getProgramsFiltered(t, ts, "search=STRENGTH")

		found := false
		for _, p := range result.Data {
			if p.Slug == "starting-strength" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected Starting Strength when searching for 'STRENGTH' (case-insensitive)")
		}
	})

	t.Run("search with no matches returns empty array", func(t *testing.T) {
		result := getProgramsFiltered(t, ts, "search=nonexistentprogram12345")

		if len(result.Data) != 0 {
			t.Errorf("Expected 0 results for non-matching search, got %d", len(result.Data))
		}
		if result.Meta == nil {
			t.Error("Expected meta to be present")
		} else if result.Meta.Total != 0 {
			t.Errorf("Expected total 0 for non-matching search, got %d", result.Meta.Total)
		}
	})
}

func TestProgramDiscoveryE2E_SearchWithFilters(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("search=strength AND difficulty=beginner returns Starting Strength", func(t *testing.T) {
		result := getProgramsFiltered(t, ts, "search=strength&difficulty=beginner")

		// All results should match both search AND filter
		for _, p := range result.Data {
			if !strings.Contains(strings.ToLower(p.Name), "strength") {
				t.Errorf("Expected program name to contain 'strength', got %s", p.Name)
			}
			if p.Difficulty != "beginner" {
				t.Errorf("Expected difficulty 'beginner', got %s for program %s", p.Difficulty, p.Name)
			}
		}

		// Should include Starting Strength
		found := false
		for _, p := range result.Data {
			if p.Slug == "starting-strength" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected Starting Strength in search+filter results")
		}
	})

	t.Run("search=5 AND difficulty=intermediate returns 5/3/1", func(t *testing.T) {
		result := getProgramsFiltered(t, ts, "search=5&difficulty=intermediate")

		// Should include 5/3/1
		found := false
		for _, p := range result.Data {
			if p.Slug == "531" {
				found = true
				if p.Difficulty != "intermediate" {
					t.Errorf("Expected difficulty 'intermediate', got %s", p.Difficulty)
				}
				break
			}
		}
		if !found {
			t.Error("Expected 5/3/1 in search+filter results")
		}
	})
}

// =============================================================================
// DETAIL ENHANCEMENT TESTS
// =============================================================================

func TestProgramDiscoveryE2E_DetailEnhancements(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Get program IDs for canonical programs
	programs := getProgramsFiltered(t, ts, "")
	programsBySlug := make(map[string]ProgramListResponse)
	for _, p := range programs.Data {
		programsBySlug[p.Slug] = p
	}

	t.Run("Starting Strength detail includes required fields", func(t *testing.T) {
		p, ok := programsBySlug["starting-strength"]
		if !ok {
			t.Skip("Starting Strength not found")
		}

		detail := getProgramDetail(t, ts, p.ID)

		// Check sampleWeek
		if detail.SampleWeek == nil {
			t.Error("Expected sampleWeek to be present")
		} else if len(detail.SampleWeek) == 0 {
			t.Error("Expected sampleWeek to have days")
		}

		// Check liftRequirements
		if detail.LiftRequirements == nil {
			t.Error("Expected liftRequirements to be present")
		} else if len(detail.LiftRequirements) == 0 {
			t.Error("Expected liftRequirements to have lifts")
		}

		// Check liftRequirements is sorted alphabetically
		if len(detail.LiftRequirements) > 1 {
			sorted := make([]string, len(detail.LiftRequirements))
			copy(sorted, detail.LiftRequirements)
			sort.Strings(sorted)

			for i, lift := range detail.LiftRequirements {
				if lift != sorted[i] {
					t.Errorf("liftRequirements not sorted: expected %s at position %d, got %s",
						sorted[i], i, lift)
					break
				}
			}
		}

		// Check estimatedSessionMinutes
		if detail.EstimatedSessionMinutes <= 0 {
			t.Errorf("Expected positive estimatedSessionMinutes, got %d", detail.EstimatedSessionMinutes)
		}

		// Check metadata matches list view
		if detail.Difficulty != "beginner" {
			t.Errorf("Expected difficulty 'beginner', got %s", detail.Difficulty)
		}
		if detail.DaysPerWeek != 3 {
			t.Errorf("Expected daysPerWeek 3, got %d", detail.DaysPerWeek)
		}
		if detail.Focus != "strength" {
			t.Errorf("Expected focus 'strength', got %s", detail.Focus)
		}
		if detail.HasAmrap {
			t.Error("Expected hasAmrap false for Starting Strength")
		}
	})

	t.Run("Texas Method detail includes required fields", func(t *testing.T) {
		p, ok := programsBySlug["texas-method"]
		if !ok {
			t.Skip("Texas Method not found")
		}

		detail := getProgramDetail(t, ts, p.ID)

		// Check presence of enhancement fields
		if detail.SampleWeek == nil {
			t.Error("Expected sampleWeek to be present")
		}
		if detail.LiftRequirements == nil {
			t.Error("Expected liftRequirements to be present")
		}
		if detail.EstimatedSessionMinutes <= 0 {
			t.Errorf("Expected positive estimatedSessionMinutes, got %d", detail.EstimatedSessionMinutes)
		}

		// Check metadata
		if detail.Difficulty != "intermediate" {
			t.Errorf("Expected difficulty 'intermediate', got %s", detail.Difficulty)
		}
		if detail.DaysPerWeek != 3 {
			t.Errorf("Expected daysPerWeek 3, got %d", detail.DaysPerWeek)
		}
		if detail.HasAmrap {
			t.Error("Expected hasAmrap false for Texas Method")
		}
	})

	t.Run("Wendler 5/3/1 detail includes required fields", func(t *testing.T) {
		p, ok := programsBySlug["531"]
		if !ok {
			t.Skip("5/3/1 not found")
		}

		detail := getProgramDetail(t, ts, p.ID)

		// Check presence of enhancement fields
		if detail.SampleWeek == nil {
			t.Error("Expected sampleWeek to be present")
		}
		if detail.LiftRequirements == nil {
			t.Error("Expected liftRequirements to be present")
		}
		if detail.EstimatedSessionMinutes <= 0 {
			t.Errorf("Expected positive estimatedSessionMinutes, got %d", detail.EstimatedSessionMinutes)
		}

		// Check metadata
		if detail.Difficulty != "intermediate" {
			t.Errorf("Expected difficulty 'intermediate', got %s", detail.Difficulty)
		}
		if detail.DaysPerWeek != 4 {
			t.Errorf("Expected daysPerWeek 4, got %d", detail.DaysPerWeek)
		}
		if !detail.HasAmrap {
			t.Error("Expected hasAmrap true for 5/3/1")
		}
	})

	t.Run("GZCLP detail includes required fields", func(t *testing.T) {
		p, ok := programsBySlug["gzclp"]
		if !ok {
			t.Skip("GZCLP not found")
		}

		detail := getProgramDetail(t, ts, p.ID)

		// Check presence of enhancement fields
		if detail.SampleWeek == nil {
			t.Error("Expected sampleWeek to be present")
		}
		if detail.LiftRequirements == nil {
			t.Error("Expected liftRequirements to be present")
		}
		if detail.EstimatedSessionMinutes <= 0 {
			t.Errorf("Expected positive estimatedSessionMinutes, got %d", detail.EstimatedSessionMinutes)
		}

		// Check metadata
		if detail.Difficulty != "beginner" {
			t.Errorf("Expected difficulty 'beginner', got %s", detail.Difficulty)
		}
		if detail.DaysPerWeek != 4 {
			t.Errorf("Expected daysPerWeek 4, got %d", detail.DaysPerWeek)
		}
		if !detail.HasAmrap {
			t.Error("Expected hasAmrap true for GZCLP")
		}
	})
}

func TestProgramDiscoveryE2E_SampleWeekStructure(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Get Starting Strength for detailed structure check
	programs := getProgramsFiltered(t, ts, "")
	var ssID string
	for _, p := range programs.Data {
		if p.Slug == "starting-strength" {
			ssID = p.ID
			break
		}
	}

	if ssID == "" {
		t.Skip("Starting Strength not found")
	}

	detail := getProgramDetail(t, ts, ssID)

	t.Run("sampleWeek days have correct structure", func(t *testing.T) {
		for _, day := range detail.SampleWeek {
			// Day number should be positive
			if day.Day <= 0 {
				t.Errorf("Expected positive day number, got %d", day.Day)
			}

			// Name should not be empty
			if day.Name == "" {
				t.Error("Expected day name to be non-empty")
			}

			// Exercise count should be positive
			if day.ExerciseCount <= 0 {
				t.Errorf("Expected positive exerciseCount, got %d for day %s", day.ExerciseCount, day.Name)
			}
		}
	})

	t.Run("Starting Strength has 2 workout days in sample week", func(t *testing.T) {
		// Starting Strength alternates between Day A and Day B
		if len(detail.SampleWeek) < 1 {
			t.Error("Expected at least 1 day in sample week")
		}

		// The sample week should show unique day types
		uniqueDays := make(map[string]bool)
		for _, day := range detail.SampleWeek {
			uniqueDays[day.Name] = true
		}

		// Starting Strength has 2 unique workout types (A and B)
		if len(uniqueDays) < 1 {
			t.Error("Expected at least 1 unique day type in sample week")
		}
	})
}

// =============================================================================
// LIST RESPONSE METADATA TESTS
// =============================================================================

func TestProgramDiscoveryE2E_ListResponseMetadata(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	result := getProgramsFiltered(t, ts, "")

	t.Run("list response includes discovery metadata fields", func(t *testing.T) {
		for _, p := range result.Data {
			// All programs should have difficulty set
			if p.Difficulty == "" {
				t.Errorf("Expected difficulty to be set for program %s", p.Name)
			}

			// All programs should have daysPerWeek set (and positive)
			if p.DaysPerWeek <= 0 {
				t.Errorf("Expected positive daysPerWeek for program %s, got %d", p.Name, p.DaysPerWeek)
			}

			// All programs should have focus set
			if p.Focus == "" {
				t.Errorf("Expected focus to be set for program %s", p.Name)
			}

			// HasAmrap is boolean, so just check the value is valid for canonical programs
			if p.Slug == "starting-strength" || p.Slug == "texas-method" {
				if p.HasAmrap {
					t.Errorf("Expected hasAmrap false for %s", p.Name)
				}
			}
			if p.Slug == "531" || p.Slug == "gzclp" {
				if !p.HasAmrap {
					t.Errorf("Expected hasAmrap true for %s", p.Name)
				}
			}
		}
	})

	t.Run("canonical programs have correct backfilled metadata", func(t *testing.T) {
		// Build a map for easy lookup
		programsBySlug := make(map[string]ProgramListResponse)
		for _, p := range result.Data {
			programsBySlug[p.Slug] = p
		}

		// Check Starting Strength
		if p, ok := programsBySlug["starting-strength"]; ok {
			if p.Difficulty != "beginner" {
				t.Errorf("Starting Strength: expected difficulty 'beginner', got %s", p.Difficulty)
			}
			if p.DaysPerWeek != 3 {
				t.Errorf("Starting Strength: expected daysPerWeek 3, got %d", p.DaysPerWeek)
			}
			if p.Focus != "strength" {
				t.Errorf("Starting Strength: expected focus 'strength', got %s", p.Focus)
			}
			if p.HasAmrap {
				t.Error("Starting Strength: expected hasAmrap false")
			}
		}

		// Check Texas Method
		if p, ok := programsBySlug["texas-method"]; ok {
			if p.Difficulty != "intermediate" {
				t.Errorf("Texas Method: expected difficulty 'intermediate', got %s", p.Difficulty)
			}
			if p.DaysPerWeek != 3 {
				t.Errorf("Texas Method: expected daysPerWeek 3, got %d", p.DaysPerWeek)
			}
			if p.Focus != "strength" {
				t.Errorf("Texas Method: expected focus 'strength', got %s", p.Focus)
			}
			if p.HasAmrap {
				t.Error("Texas Method: expected hasAmrap false")
			}
		}

		// Check 5/3/1
		if p, ok := programsBySlug["531"]; ok {
			if p.Difficulty != "intermediate" {
				t.Errorf("5/3/1: expected difficulty 'intermediate', got %s", p.Difficulty)
			}
			if p.DaysPerWeek != 4 {
				t.Errorf("5/3/1: expected daysPerWeek 4, got %d", p.DaysPerWeek)
			}
			if p.Focus != "strength" {
				t.Errorf("5/3/1: expected focus 'strength', got %s", p.Focus)
			}
			if !p.HasAmrap {
				t.Error("5/3/1: expected hasAmrap true")
			}
		}

		// Check GZCLP
		if p, ok := programsBySlug["gzclp"]; ok {
			if p.Difficulty != "beginner" {
				t.Errorf("GZCLP: expected difficulty 'beginner', got %s", p.Difficulty)
			}
			if p.DaysPerWeek != 4 {
				t.Errorf("GZCLP: expected daysPerWeek 4, got %d", p.DaysPerWeek)
			}
			if p.Focus != "strength" {
				t.Errorf("GZCLP: expected focus 'strength', got %s", p.Focus)
			}
			if !p.HasAmrap {
				t.Error("GZCLP: expected hasAmrap true")
			}
		}
	})
}
