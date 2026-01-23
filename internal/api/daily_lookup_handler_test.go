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

// DailyLookupEntryTestResponse matches the API response format for a daily lookup entry.
type DailyLookupEntryTestResponse struct {
	DayIdentifier      string  `json:"dayIdentifier"`
	PercentageModifier float64 `json:"percentageModifier"`
	IntensityLevel     *string `json:"intensityLevel,omitempty"`
}

// DailyLookupTestResponse matches the API response format for a daily lookup.
type DailyLookupTestResponse struct {
	ID        string                         `json:"id"`
	Name      string                         `json:"name"`
	Entries   []DailyLookupEntryTestResponse `json:"entries"`
	ProgramID *string                        `json:"programId,omitempty"`
	CreatedAt time.Time                      `json:"createdAt"`
	UpdatedAt time.Time                      `json:"updatedAt"`
}

// DailyLookupPaginationMeta contains pagination metadata.
type DailyLookupPaginationMeta struct {
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasMore bool  `json:"hasMore"`
}

// PaginatedDailyLookupsResponse is the paginated list response.
type PaginatedDailyLookupsResponse struct {
	Data []DailyLookupTestResponse  `json:"data"`
	Meta *DailyLookupPaginationMeta `json:"meta"`
}

// DailyLookupEnvelope wraps single daily lookup response with standard envelope.
type DailyLookupEnvelope struct {
	Data DailyLookupTestResponse `json:"data"`
}

// authGetDailyLookup performs an authenticated GET request
func authGetDailyLookup(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestUserID)
	return http.DefaultClient.Do(req)
}

// adminPostDailyLookup performs an admin-authenticated POST request
func adminPostDailyLookup(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// adminPutDailyLookup performs an admin-authenticated PUT request
func adminPutDailyLookup(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// adminDeleteDailyLookup performs an admin-authenticated DELETE request
func adminDeleteDailyLookup(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

func TestDailyLookupCRUD(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	var createdLookup DailyLookupTestResponse

	t.Run("creates daily lookup with required fields", func(t *testing.T) {
		body := `{
			"name": "Bill Starr Intensities",
			"entries": [
				{"dayIdentifier": "heavy", "percentageModifier": 100, "intensityLevel": "HEAVY"},
				{"dayIdentifier": "light", "percentageModifier": 70, "intensityLevel": "LIGHT"},
				{"dayIdentifier": "medium", "percentageModifier": 80, "intensityLevel": "MEDIUM"}
			]
		}`
		resp, err := adminPostDailyLookup(ts.URL("/daily-lookups"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope DailyLookupEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		createdLookup = envelope.Data

		if createdLookup.ID == "" {
			t.Error("Expected non-empty ID")
		}
		if createdLookup.Name != "Bill Starr Intensities" {
			t.Errorf("Expected name 'Bill Starr Intensities', got %s", createdLookup.Name)
		}
		if len(createdLookup.Entries) != 3 {
			t.Errorf("Expected 3 entries, got %d", len(createdLookup.Entries))
		}
	})

	t.Run("creates daily lookup without intensity level", func(t *testing.T) {
		body := `{
			"name": "Simple Day Lookup",
			"entries": [
				{"dayIdentifier": "day-a", "percentageModifier": 100},
				{"dayIdentifier": "day-b", "percentageModifier": 95}
			]
		}`
		resp, err := adminPostDailyLookup(ts.URL("/daily-lookups"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope DailyLookupEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		lookup := envelope.Data

		// Verify intensity level is nil
		for _, entry := range lookup.Entries {
			if entry.IntensityLevel != nil {
				t.Errorf("Expected nil intensityLevel, got %v", *entry.IntensityLevel)
			}
		}
	})

	t.Run("gets daily lookup by ID", func(t *testing.T) {
		resp, err := authGetDailyLookup(ts.URL("/daily-lookups/" + createdLookup.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope DailyLookupEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		lookup := envelope.Data

		if lookup.ID != createdLookup.ID {
			t.Errorf("Expected ID %s, got %s", createdLookup.ID, lookup.ID)
		}
		if lookup.Name != "Bill Starr Intensities" {
			t.Errorf("Expected name 'Bill Starr Intensities', got %s", lookup.Name)
		}
	})

	t.Run("returns 404 for non-existent daily lookup", func(t *testing.T) {
		resp, err := authGetDailyLookup(ts.URL("/daily-lookups/non-existent-id"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("lists daily lookups with pagination", func(t *testing.T) {
		resp, err := authGetDailyLookup(ts.URL("/daily-lookups"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var result PaginatedDailyLookupsResponse
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

	t.Run("updates daily lookup name", func(t *testing.T) {
		body := `{"name": "Modified Bill Starr Intensities"}`
		resp, err := adminPutDailyLookup(ts.URL("/daily-lookups/"+createdLookup.ID), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope DailyLookupEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		updated := envelope.Data

		if updated.Name != "Modified Bill Starr Intensities" {
			t.Errorf("Expected name 'Modified Bill Starr Intensities', got %s", updated.Name)
		}
	})

	t.Run("updates daily lookup entries", func(t *testing.T) {
		body := `{
			"entries": [
				{"dayIdentifier": "heavy", "percentageModifier": 100},
				{"dayIdentifier": "light", "percentageModifier": 60}
			]
		}`
		resp, err := adminPutDailyLookup(ts.URL("/daily-lookups/"+createdLookup.ID), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope DailyLookupEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		updated := envelope.Data

		if len(updated.Entries) != 2 {
			t.Errorf("Expected 2 entries, got %d", len(updated.Entries))
		}
	})

	t.Run("deletes daily lookup", func(t *testing.T) {
		// Create a lookup to delete
		body := `{
			"name": "Lookup To Delete",
			"entries": [{"dayIdentifier": "test", "percentageModifier": 100}]
		}`
		createResp, _ := adminPostDailyLookup(ts.URL("/daily-lookups"), body)
		var createEnvelope DailyLookupEnvelope
		json.NewDecoder(createResp.Body).Decode(&createEnvelope)
		toDelete := createEnvelope.Data
		createResp.Body.Close()

		resp, err := adminDeleteDailyLookup(ts.URL("/daily-lookups/" + toDelete.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 204, got %d: %s", resp.StatusCode, bodyBytes)
		}

		// Verify it's deleted
		getResp, _ := authGetDailyLookup(ts.URL("/daily-lookups/" + toDelete.ID))
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404 after delete, got %d", getResp.StatusCode)
		}
	})
}

func TestDailyLookupValidation(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("rejects empty name", func(t *testing.T) {
		body := `{
			"name": "",
			"entries": [{"dayIdentifier": "heavy", "percentageModifier": 100}]
		}`
		resp, err := adminPostDailyLookup(ts.URL("/daily-lookups"), body)
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
		resp, err := adminPostDailyLookup(ts.URL("/daily-lookups"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects empty day identifier", func(t *testing.T) {
		body := `{
			"name": "Test Lookup",
			"entries": [{"dayIdentifier": "", "percentageModifier": 100}]
		}`
		resp, err := adminPostDailyLookup(ts.URL("/daily-lookups"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects duplicate day identifiers", func(t *testing.T) {
		body := `{
			"name": "Test Lookup",
			"entries": [
				{"dayIdentifier": "heavy", "percentageModifier": 100},
				{"dayIdentifier": "heavy", "percentageModifier": 90}
			]
		}`
		resp, err := adminPostDailyLookup(ts.URL("/daily-lookups"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects invalid intensity level", func(t *testing.T) {
		body := `{
			"name": "Test Lookup",
			"entries": [{"dayIdentifier": "monday", "percentageModifier": 100, "intensityLevel": "INVALID"}]
		}`
		resp, err := adminPostDailyLookup(ts.URL("/daily-lookups"), body)
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
		resp, err := adminPutDailyLookup(ts.URL("/daily-lookups/non-existent-id"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 404 when deleting non-existent lookup", func(t *testing.T) {
		resp, err := adminDeleteDailyLookup(ts.URL("/daily-lookups/non-existent-id"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestDailyLookupAuthorization(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create lookup as admin
	lookupBody := `{
		"name": "Auth Test Lookup",
		"entries": [{"dayIdentifier": "heavy", "percentageModifier": 100}]
	}`
	lookupResp, _ := adminPostDailyLookup(ts.URL("/daily-lookups"), lookupBody)
	var lookupEnvelope DailyLookupEnvelope
	json.NewDecoder(lookupResp.Body).Decode(&lookupEnvelope)
	createdLookup := lookupEnvelope.Data
	lookupResp.Body.Close()

	t.Run("unauthenticated user gets 401 on GET /daily-lookups", func(t *testing.T) {
		resp, err := http.Get(ts.URL("/daily-lookups"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("authenticated user can GET /daily-lookups", func(t *testing.T) {
		resp, err := authGetDailyLookup(ts.URL("/daily-lookups"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("non-admin user gets 403 on POST /daily-lookups", func(t *testing.T) {
		body := `{
			"name": "Unauthorized Lookup",
			"entries": [{"dayIdentifier": "heavy", "percentageModifier": 100}]
		}`
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/daily-lookups"), bytes.NewBufferString(body))
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

	t.Run("non-admin user gets 403 on PUT /daily-lookups/{id}", func(t *testing.T) {
		body := `{"name": "Modified Name"}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/daily-lookups/"+createdLookup.ID), bytes.NewBufferString(body))
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

	t.Run("non-admin user gets 403 on DELETE /daily-lookups/{id}", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL("/daily-lookups/"+createdLookup.ID), nil)
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
