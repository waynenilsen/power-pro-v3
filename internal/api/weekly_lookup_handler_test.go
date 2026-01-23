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

// WeeklyLookupEntryTestResponse matches the API response format for a weekly lookup entry.
type WeeklyLookupEntryTestResponse struct {
	WeekNumber         int       `json:"weekNumber"`
	Percentages        []float64 `json:"percentages"`
	Reps               []int     `json:"reps"`
	PercentageModifier *float64  `json:"percentageModifier,omitempty"`
}

// WeeklyLookupTestResponse matches the API response format for a weekly lookup.
type WeeklyLookupTestResponse struct {
	ID        string                          `json:"id"`
	Name      string                          `json:"name"`
	Entries   []WeeklyLookupEntryTestResponse `json:"entries"`
	ProgramID *string                         `json:"programId,omitempty"`
	CreatedAt time.Time                       `json:"createdAt"`
	UpdatedAt time.Time                       `json:"updatedAt"`
}

// WeeklyLookupPaginationMeta contains pagination metadata.
type WeeklyLookupPaginationMeta struct {
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasMore bool  `json:"hasMore"`
}

// PaginatedWeeklyLookupsResponse is the paginated list response.
type PaginatedWeeklyLookupsResponse struct {
	Data []WeeklyLookupTestResponse  `json:"data"`
	Meta *WeeklyLookupPaginationMeta `json:"meta"`
}

// WeeklyLookupEnvelope wraps single weekly lookup response with standard envelope.
type WeeklyLookupEnvelope struct {
	Data WeeklyLookupTestResponse `json:"data"`
}

// authGetWeeklyLookup performs an authenticated GET request
func authGetWeeklyLookup(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestUserID)
	return http.DefaultClient.Do(req)
}

// adminPostWeeklyLookup performs an admin-authenticated POST request
func adminPostWeeklyLookup(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// adminPutWeeklyLookup performs an admin-authenticated PUT request
func adminPutWeeklyLookup(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// adminDeleteWeeklyLookup performs an admin-authenticated DELETE request
func adminDeleteWeeklyLookup(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

func TestWeeklyLookupCRUD(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	var createdLookup WeeklyLookupTestResponse

	t.Run("creates weekly lookup with required fields", func(t *testing.T) {
		body := `{
			"name": "5/3/1 Percentages",
			"entries": [
				{"weekNumber": 1, "percentages": [65, 75, 85], "reps": [5, 5, 5]},
				{"weekNumber": 2, "percentages": [70, 80, 90], "reps": [3, 3, 3]},
				{"weekNumber": 3, "percentages": [75, 85, 95], "reps": [5, 3, 1]},
				{"weekNumber": 4, "percentages": [40, 50, 60], "reps": [5, 5, 5]}
			]
		}`
		resp, err := adminPostWeeklyLookup(ts.URL("/weekly-lookups"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope WeeklyLookupEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		createdLookup = envelope.Data

		if createdLookup.ID == "" {
			t.Error("Expected non-empty ID")
		}
		if createdLookup.Name != "5/3/1 Percentages" {
			t.Errorf("Expected name '5/3/1 Percentages', got %s", createdLookup.Name)
		}
		if len(createdLookup.Entries) != 4 {
			t.Errorf("Expected 4 entries, got %d", len(createdLookup.Entries))
		}
	})

	t.Run("creates weekly lookup with single entry", func(t *testing.T) {
		body := `{
			"name": "Simple Lookup",
			"entries": [
				{"weekNumber": 1, "percentages": [80], "reps": [5]}
			]
		}`
		resp, err := adminPostWeeklyLookup(ts.URL("/weekly-lookups"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("creates weekly lookup with percentage modifier", func(t *testing.T) {
		body := `{
			"name": "Lookup with Modifier",
			"entries": [
				{"weekNumber": 1, "percentages": [80], "reps": [5], "percentageModifier": 0.9}
			]
		}`
		resp, err := adminPostWeeklyLookup(ts.URL("/weekly-lookups"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope WeeklyLookupEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		lookup := envelope.Data

		if lookup.Entries[0].PercentageModifier == nil {
			t.Error("Expected percentageModifier to be set")
		} else if *lookup.Entries[0].PercentageModifier != 0.9 {
			t.Errorf("Expected percentageModifier 0.9, got %f", *lookup.Entries[0].PercentageModifier)
		}
	})

	t.Run("gets weekly lookup by ID", func(t *testing.T) {
		resp, err := authGetWeeklyLookup(ts.URL("/weekly-lookups/" + createdLookup.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope WeeklyLookupEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		lookup := envelope.Data

		if lookup.ID != createdLookup.ID {
			t.Errorf("Expected ID %s, got %s", createdLookup.ID, lookup.ID)
		}
		if lookup.Name != "5/3/1 Percentages" {
			t.Errorf("Expected name '5/3/1 Percentages', got %s", lookup.Name)
		}
	})

	t.Run("returns 404 for non-existent weekly lookup", func(t *testing.T) {
		resp, err := authGetWeeklyLookup(ts.URL("/weekly-lookups/non-existent-id"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("lists weekly lookups with pagination", func(t *testing.T) {
		resp, err := authGetWeeklyLookup(ts.URL("/weekly-lookups"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result PaginatedWeeklyLookupsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result.Meta == nil || result.Meta.Total < 1 {
			total := int64(0)
			if result.Meta != nil {
				total = result.Meta.Total
			}
			t.Errorf("Expected at least 1 lookup, got %d", total)
		}
		if result.Meta == nil || result.Meta.Offset != 0 {
			offset := 0
			if result.Meta != nil {
				offset = result.Meta.Offset
			}
			t.Errorf("Expected offset 0, got %d", offset)
		}
	})

	t.Run("updates weekly lookup name", func(t *testing.T) {
		body := `{"name": "Modified 5/3/1 Percentages"}`
		resp, err := adminPutWeeklyLookup(ts.URL("/weekly-lookups/"+createdLookup.ID), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope WeeklyLookupEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		updated := envelope.Data

		if updated.Name != "Modified 5/3/1 Percentages" {
			t.Errorf("Expected name 'Modified 5/3/1 Percentages', got %s", updated.Name)
		}
	})

	t.Run("updates weekly lookup entries", func(t *testing.T) {
		body := `{
			"entries": [
				{"weekNumber": 1, "percentages": [70, 80, 90], "reps": [3, 3, 3]}
			]
		}`
		resp, err := adminPutWeeklyLookup(ts.URL("/weekly-lookups/"+createdLookup.ID), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope WeeklyLookupEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		updated := envelope.Data

		if len(updated.Entries) != 1 {
			t.Errorf("Expected 1 entry, got %d", len(updated.Entries))
		}
	})

	t.Run("deletes weekly lookup", func(t *testing.T) {
		// Create a lookup to delete
		body := `{
			"name": "Lookup To Delete",
			"entries": [{"weekNumber": 1, "percentages": [80], "reps": [5]}]
		}`
		createResp, _ := adminPostWeeklyLookup(ts.URL("/weekly-lookups"), body)
		var createEnvelope WeeklyLookupEnvelope
		json.NewDecoder(createResp.Body).Decode(&createEnvelope)
		toDelete := createEnvelope.Data
		createResp.Body.Close()

		resp, err := adminDeleteWeeklyLookup(ts.URL("/weekly-lookups/" + toDelete.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 204, got %d: %s", resp.StatusCode, bodyBytes)
		}

		// Verify it's deleted
		getResp, _ := authGetWeeklyLookup(ts.URL("/weekly-lookups/" + toDelete.ID))
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404 after delete, got %d", getResp.StatusCode)
		}
	})
}

func TestWeeklyLookupValidation(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("rejects empty name", func(t *testing.T) {
		body := `{
			"name": "",
			"entries": [{"weekNumber": 1, "percentages": [80], "reps": [5]}]
		}`
		resp, err := adminPostWeeklyLookup(ts.URL("/weekly-lookups"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects empty entries", func(t *testing.T) {
		body := `{"name": "Test Lookup", "entries": []}`
		resp, err := adminPostWeeklyLookup(ts.URL("/weekly-lookups"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects invalid week number (zero)", func(t *testing.T) {
		body := `{
			"name": "Test Lookup",
			"entries": [{"weekNumber": 0, "percentages": [80], "reps": [5]}]
		}`
		resp, err := adminPostWeeklyLookup(ts.URL("/weekly-lookups"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects duplicate week numbers", func(t *testing.T) {
		body := `{
			"name": "Test Lookup",
			"entries": [
				{"weekNumber": 1, "percentages": [80], "reps": [5]},
				{"weekNumber": 1, "percentages": [85], "reps": [3]}
			]
		}`
		resp, err := adminPostWeeklyLookup(ts.URL("/weekly-lookups"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects mismatched percentages and reps lengths", func(t *testing.T) {
		body := `{
			"name": "Test Lookup",
			"entries": [{"weekNumber": 1, "percentages": [80, 85, 90], "reps": [5, 3]}]
		}`
		resp, err := adminPostWeeklyLookup(ts.URL("/weekly-lookups"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 404 when updating non-existent lookup", func(t *testing.T) {
		body := `{"name": "New Name"}`
		resp, err := adminPutWeeklyLookup(ts.URL("/weekly-lookups/non-existent-id"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 404 when deleting non-existent lookup", func(t *testing.T) {
		resp, err := adminDeleteWeeklyLookup(ts.URL("/weekly-lookups/non-existent-id"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestWeeklyLookupAuthorization(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create lookup as admin
	lookupBody := `{
		"name": "Auth Test Lookup",
		"entries": [{"weekNumber": 1, "percentages": [80], "reps": [5]}]
	}`
	lookupResp, _ := adminPostWeeklyLookup(ts.URL("/weekly-lookups"), lookupBody)
	var lookupEnvelope WeeklyLookupEnvelope
	json.NewDecoder(lookupResp.Body).Decode(&lookupEnvelope)
	createdLookup := lookupEnvelope.Data
	lookupResp.Body.Close()

	t.Run("unauthenticated user gets 401 on GET /weekly-lookups", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/weekly-lookups"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("authenticated user can GET /weekly-lookups", func(t *testing.T) {
		resp, err := authGetWeeklyLookup(ts.URL("/weekly-lookups"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("non-admin user gets 403 on POST /weekly-lookups", func(t *testing.T) {
		body := `{
			"name": "Unauthorized Lookup",
			"entries": [{"weekNumber": 1, "percentages": [80], "reps": [5]}]
		}`
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/weekly-lookups"), bytes.NewBufferString(body))
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

	t.Run("non-admin user gets 403 on PUT /weekly-lookups/{id}", func(t *testing.T) {
		body := `{"name": "Modified Name"}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/weekly-lookups/"+createdLookup.ID), bytes.NewBufferString(body))
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

	t.Run("non-admin user gets 403 on DELETE /weekly-lookups/{id}", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL("/weekly-lookups/"+createdLookup.ID), nil)
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
