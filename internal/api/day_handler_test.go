package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// DayResponse matches the API response format.
type DayResponse struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Slug      string                 `json:"slug"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	ProgramID *string                `json:"programId,omitempty"`
	CreatedAt string                 `json:"createdAt"`
	UpdatedAt string                 `json:"updatedAt"`
}

// DayWithPrescriptionsResponse is the detailed day response with prescriptions.
type DayWithPrescriptionsResponse struct {
	ID            string                    `json:"id"`
	Name          string                    `json:"name"`
	Slug          string                    `json:"slug"`
	Metadata      map[string]interface{}    `json:"metadata,omitempty"`
	ProgramID     *string                   `json:"programId,omitempty"`
	Prescriptions []DayPrescriptionResponse `json:"prescriptions"`
	CreatedAt     string                    `json:"createdAt"`
	UpdatedAt     string                    `json:"updatedAt"`
}

// DayPrescriptionResponse represents a prescription within a day.
type DayPrescriptionResponse struct {
	ID             string `json:"id"`
	PrescriptionID string `json:"prescriptionId"`
	Order          int    `json:"order"`
	CreatedAt      string `json:"createdAt"`
}

// DayPaginationMeta contains pagination metadata.
type DayPaginationMeta struct {
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasMore bool  `json:"hasMore"`
}

// PaginatedDaysResponse is the paginated list response.
type PaginatedDaysResponse struct {
	Data []DayResponse      `json:"data"`
	Meta *DayPaginationMeta `json:"meta"`
}

// DayEnvelope wraps single day response with standard envelope.
type DayEnvelope struct {
	Data DayResponse `json:"data"`
}

// DayWithPrescriptionsEnvelope wraps day with prescriptions response.
type DayWithPrescriptionsEnvelope struct {
	Data DayWithPrescriptionsResponse `json:"data"`
}

func TestListDays(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create some days for testing
	day1Body := `{"name": "Day A", "slug": "day-a"}`
	day1Resp, _ := adminPost(ts.URL("/days"), day1Body)
	day1Resp.Body.Close()

	day2Body := `{"name": "Day B", "slug": "day-b"}`
	day2Resp, _ := adminPost(ts.URL("/days"), day2Body)
	day2Resp.Body.Close()

	day3Body := `{"name": "Heavy Day", "slug": "heavy-day"}`
	day3Resp, _ := adminPost(ts.URL("/days"), day3Body)
	day3Resp.Body.Close()

	t.Run("returns created days", func(t *testing.T) {
		resp, err := authGet(ts.URL("/days"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}

		var result PaginatedDaysResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(result.Data) != 3 {
			t.Errorf("Expected 3 days, got %d", len(result.Data))
		}
		if result.Meta == nil || result.Meta.Total != 3 {
			total := int64(0)
			if result.Meta != nil {
				total = result.Meta.Total
			}
			t.Errorf("Expected total 3, got %d", total)
		}
	})

	t.Run("supports pagination", func(t *testing.T) {
		resp, err := authGet(ts.URL("/days?limit=2&offset=0"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedDaysResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) != 2 {
			t.Errorf("Expected 2 days on page 1, got %d", len(result.Data))
		}
		if result.Meta == nil {
			t.Fatal("Expected meta to be present")
		}
		if result.Meta.Offset != 0 {
			t.Errorf("Expected offset 0, got %d", result.Meta.Offset)
		}
		if result.Meta.Limit != 2 {
			t.Errorf("Expected limit 2, got %d", result.Meta.Limit)
		}

		// Get page 2 (offset=2)
		resp2, _ := authGet(ts.URL("/days?limit=2&offset=2"))
		defer resp2.Body.Close()

		var result2 PaginatedDaysResponse
		json.NewDecoder(resp2.Body).Decode(&result2)

		if len(result2.Data) != 1 {
			t.Errorf("Expected 1 day on page 2, got %d", len(result2.Data))
		}
	})

	t.Run("supports sorting by name", func(t *testing.T) {
		// Ascending (default)
		resp, _ := authGet(ts.URL("/days?sortBy=name&sortOrder=asc"))
		defer resp.Body.Close()

		var result PaginatedDaysResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) < 2 {
			t.Fatal("Need at least 2 days for sort test")
		}

		// First should be "Day A" (alphabetically first)
		if result.Data[0].Name != "Day A" {
			t.Errorf("Expected first day to be 'Day A', got %s", result.Data[0].Name)
		}

		// Descending
		resp2, _ := authGet(ts.URL("/days?sortBy=name&sortOrder=desc"))
		defer resp2.Body.Close()

		var result2 PaginatedDaysResponse
		json.NewDecoder(resp2.Body).Decode(&result2)

		// First should be "Heavy Day" (alphabetically last)
		if result2.Data[0].Name != "Heavy Day" {
			t.Errorf("Expected first day to be 'Heavy Day', got %s", result2.Data[0].Name)
		}
	})
}

func TestGetDay(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a day
	createBody := `{"name": "Test Day", "slug": "test-day"}`
	createResp, _ := adminPost(ts.URL("/days"), createBody)
	var createEnvelope DayEnvelope
	json.NewDecoder(createResp.Body).Decode(&createEnvelope)
	createdDay := createEnvelope.Data
	createResp.Body.Close()

	t.Run("returns day by ID with prescriptions", func(t *testing.T) {
		resp, err := authGet(ts.URL("/days/" + createdDay.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var envelope DayWithPrescriptionsEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		day := envelope.Data

		if day.ID != createdDay.ID {
			t.Errorf("Expected ID %s, got %s", createdDay.ID, day.ID)
		}
		if day.Name != "Test Day" {
			t.Errorf("Expected name 'Test Day', got %s", day.Name)
		}
		if day.Slug != "test-day" {
			t.Errorf("Expected slug 'test-day', got %s", day.Slug)
		}
		if day.Prescriptions == nil {
			t.Errorf("Expected prescriptions to be an empty array, got nil")
		}
	})

	t.Run("returns 404 for non-existent ID", func(t *testing.T) {
		resp, _ := authGet(ts.URL("/days/non-existent-id"))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)

		if !strings.Contains(errResp.Error.Message, "day not found") {
			t.Errorf("Expected error to contain 'day not found', got %s", errResp.Error.Message)
		}
	})
}

func TestGetDayBySlug(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a day
	createBody := `{"name": "Heavy Day", "slug": "heavy-day"}`
	createResp, _ := adminPost(ts.URL("/days"), createBody)
	createResp.Body.Close()

	t.Run("returns day by slug", func(t *testing.T) {
		resp, err := authGet(ts.URL("/days/by-slug/heavy-day"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var envelope DayWithPrescriptionsEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		day := envelope.Data

		if day.Slug != "heavy-day" {
			t.Errorf("Expected slug 'heavy-day', got %s", day.Slug)
		}
		if day.Name != "Heavy Day" {
			t.Errorf("Expected name 'Heavy Day', got %s", day.Name)
		}
	})

	t.Run("returns 404 for non-existent slug", func(t *testing.T) {
		resp, _ := authGet(ts.URL("/days/by-slug/non-existent"))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestCreateDay(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("creates day with all fields", func(t *testing.T) {
		body := `{"name": "Day A", "slug": "day-a", "metadata": {"intensityLevel": "HEAVY"}}`
		resp, err := adminPost(ts.URL("/days"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, body)
		}

		var envelope DayEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		day := envelope.Data

		if day.Name != "Day A" {
			t.Errorf("Expected name 'Day A', got %s", day.Name)
		}
		if day.Slug != "day-a" {
			t.Errorf("Expected slug 'day-a', got %s", day.Slug)
		}
		if day.Metadata == nil || day.Metadata["intensityLevel"] != "HEAVY" {
			t.Errorf("Expected metadata with intensityLevel HEAVY, got %v", day.Metadata)
		}
		if day.ID == "" {
			t.Errorf("Expected ID to be generated")
		}
	})

	t.Run("auto-generates slug from name", func(t *testing.T) {
		body := `{"name": "Light Day"}`
		resp, _ := adminPost(ts.URL("/days"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope DayEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		day := envelope.Data

		if day.Slug != "light-day" {
			t.Errorf("Expected auto-generated slug 'light-day', got %s", day.Slug)
		}
	})

	t.Run("returns 400 for missing name", func(t *testing.T) {
		body := `{"slug": "no-name"}`
		resp, _ := adminPost(ts.URL("/days"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid slug format", func(t *testing.T) {
		body := `{"name": "Invalid Slug", "slug": "INVALID_SLUG!"}`
		resp, _ := adminPost(ts.URL("/days"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 409 for duplicate slug", func(t *testing.T) {
		// Create first day
		firstBody := `{"name": "First Day", "slug": "unique-slug"}`
		firstResp, _ := adminPost(ts.URL("/days"), firstBody)
		firstResp.Body.Close()

		// Try to create second day with same slug
		body := `{"name": "Another Day", "slug": "unique-slug"}`
		resp, _ := adminPost(ts.URL("/days"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409, got %d", resp.StatusCode)
		}
	})
}

func TestUpdateDay(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a day to update
	createBody := `{"name": "Update Test", "slug": "update-test"}`
	createResp, _ := adminPost(ts.URL("/days"), createBody)
	var createEnvelope DayEnvelope
	json.NewDecoder(createResp.Body).Decode(&createEnvelope)
	createdDay := createEnvelope.Data
	createResp.Body.Close()

	t.Run("updates day name", func(t *testing.T) {
		body := `{"name": "Updated Name"}`
		resp, err := adminPut(ts.URL("/days/"+createdDay.ID), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope DayEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		day := envelope.Data

		if day.Name != "Updated Name" {
			t.Errorf("Expected name 'Updated Name', got %s", day.Name)
		}
		// Slug should remain unchanged
		if day.Slug != "update-test" {
			t.Errorf("Expected slug 'update-test', got %s", day.Slug)
		}
	})

	t.Run("updates day slug", func(t *testing.T) {
		body := `{"slug": "updated-slug"}`
		resp, _ := adminPut(ts.URL("/days/"+createdDay.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var envelope DayEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		day := envelope.Data

		if day.Slug != "updated-slug" {
			t.Errorf("Expected slug 'updated-slug', got %s", day.Slug)
		}
	})

	t.Run("updates day metadata", func(t *testing.T) {
		body := `{"metadata": {"intensityLevel": "MEDIUM", "focus": "recovery"}}`
		resp, _ := adminPut(ts.URL("/days/"+createdDay.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var envelope DayEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		day := envelope.Data

		if day.Metadata == nil || day.Metadata["intensityLevel"] != "MEDIUM" {
			t.Errorf("Expected metadata with intensityLevel MEDIUM, got %v", day.Metadata)
		}
		if day.Metadata["focus"] != "recovery" {
			t.Errorf("Expected metadata with focus 'recovery', got %v", day.Metadata["focus"])
		}
	})

	t.Run("clears metadata", func(t *testing.T) {
		body := `{"clearMetadata": true}`
		resp, _ := adminPut(ts.URL("/days/"+createdDay.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var envelope DayEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		day := envelope.Data

		if day.Metadata != nil && len(day.Metadata) > 0 {
			t.Errorf("Expected metadata to be nil/empty, got %v", day.Metadata)
		}
	})

	t.Run("returns 404 for non-existent day", func(t *testing.T) {
		body := `{"name": "Updated"}`
		resp, _ := adminPut(ts.URL("/days/non-existent-id"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 409 for duplicate slug", func(t *testing.T) {
		// Create another day
		otherBody := `{"name": "Other Day", "slug": "other-day"}`
		otherResp, _ := adminPost(ts.URL("/days"), otherBody)
		otherResp.Body.Close()

		// Try to update our day to use that slug
		body := `{"slug": "other-day"}`
		resp, _ := adminPut(ts.URL("/days/"+createdDay.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for validation errors", func(t *testing.T) {
		body := `{"name": ""}`
		resp, _ := adminPut(ts.URL("/days/"+createdDay.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestDeleteDay(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("deletes day successfully", func(t *testing.T) {
		// Create a day to delete
		createBody := `{"name": "To Delete", "slug": "to-delete"}`
		createResp, _ := adminPost(ts.URL("/days"), createBody)
		var createEnvelope DayEnvelope
		json.NewDecoder(createResp.Body).Decode(&createEnvelope)
		createdDay := createEnvelope.Data
		createResp.Body.Close()

		// Delete it
		resp, err := adminDelete(ts.URL("/days/" + createdDay.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 204, got %d: %s", resp.StatusCode, bodyBytes)
		}

		// Verify it's deleted
		getResp, _ := authGet(ts.URL("/days/" + createdDay.ID))
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected deleted day to return 404, got %d", getResp.StatusCode)
		}
	})

	t.Run("returns 404 for non-existent day", func(t *testing.T) {
		resp, _ := adminDelete(ts.URL("/days/non-existent-id"))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestDayPrescriptionManagement(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a day
	createDayBody := `{"name": "Prescription Test Day", "slug": "prescription-test-day"}`
	dayResp, _ := adminPost(ts.URL("/days"), createDayBody)
	var dayEnvelope DayEnvelope
	json.NewDecoder(dayResp.Body).Decode(&dayEnvelope)
	createdDay := dayEnvelope.Data
	dayResp.Body.Close()

	if createdDay.ID == "" {
		t.Fatal("Failed to create test day")
	}

	// Create some prescriptions
	squatID := "00000000-0000-0000-0000-000000000001" // seeded squat ID
	createPrescription1 := fmt.Sprintf(`{"liftId": "%s", "loadStrategy": {"type": "PERCENT_OF", "percentage": 80, "referenceType": "TRAINING_MAX"}, "setScheme": {"type": "FIXED", "sets": 3, "reps": 5}, "order": 0}`, squatID)
	p1Resp, _ := adminPost(ts.URL("/prescriptions"), createPrescription1)
	p1Body, _ := io.ReadAll(p1Resp.Body)
	p1Resp.Body.Close()

	var p1Envelope struct {
		Data PrescriptionResponse `json:"data"`
	}
	if err := json.Unmarshal(p1Body, &p1Envelope); err != nil {
		t.Fatalf("Failed to decode prescription 1: %v, body: %s", err, string(p1Body))
	}
	p1 := p1Envelope.Data
	if p1.ID == "" {
		t.Fatalf("Failed to create prescription 1, response: %s", string(p1Body))
	}

	createPrescription2 := fmt.Sprintf(`{"liftId": "%s", "loadStrategy": {"type": "PERCENT_OF", "percentage": 85, "referenceType": "TRAINING_MAX"}, "setScheme": {"type": "FIXED", "sets": 3, "reps": 3}, "order": 1}`, squatID)
	p2Resp, _ := adminPost(ts.URL("/prescriptions"), createPrescription2)
	p2Body, _ := io.ReadAll(p2Resp.Body)
	p2Resp.Body.Close()

	var p2Envelope struct {
		Data PrescriptionResponse `json:"data"`
	}
	if err := json.Unmarshal(p2Body, &p2Envelope); err != nil {
		t.Fatalf("Failed to decode prescription 2: %v, body: %s", err, string(p2Body))
	}
	p2 := p2Envelope.Data
	if p2.ID == "" {
		t.Fatalf("Failed to create prescription 2, response: %s", string(p2Body))
	}

	t.Run("adds prescription to day with auto-order", func(t *testing.T) {
		body := fmt.Sprintf(`{"prescriptionId": "%s"}`, p1.ID)
		resp, err := adminPost(ts.URL("/days/"+createdDay.ID+"/prescriptions"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var dpEnvelope struct {
			Data DayPrescriptionResponse `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&dpEnvelope)
		dp := dpEnvelope.Data

		if dp.PrescriptionID != p1.ID {
			t.Errorf("Expected prescriptionId %s, got %s", p1.ID, dp.PrescriptionID)
		}
		if dp.Order != 0 {
			t.Errorf("Expected order 0, got %d", dp.Order)
		}
	})

	t.Run("adds second prescription with auto-incremented order", func(t *testing.T) {
		body := fmt.Sprintf(`{"prescriptionId": "%s"}`, p2.ID)
		resp, _ := adminPost(ts.URL("/days/"+createdDay.ID+"/prescriptions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var dpEnvelope struct {
			Data DayPrescriptionResponse `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&dpEnvelope)
		dp := dpEnvelope.Data

		if dp.Order != 1 {
			t.Errorf("Expected order 1, got %d", dp.Order)
		}
	})

	t.Run("returns 409 for duplicate prescription", func(t *testing.T) {
		body := fmt.Sprintf(`{"prescriptionId": "%s"}`, p1.ID)
		resp, _ := adminPost(ts.URL("/days/"+createdDay.ID+"/prescriptions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409, got %d", resp.StatusCode)
		}
	})

	t.Run("returns day with prescriptions", func(t *testing.T) {
		resp, _ := authGet(ts.URL("/days/" + createdDay.ID))
		defer resp.Body.Close()

		var envelope DayWithPrescriptionsEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		day := envelope.Data

		if len(day.Prescriptions) != 2 {
			t.Fatalf("Expected 2 prescriptions, got %d", len(day.Prescriptions))
		}

		// Should be ordered
		if day.Prescriptions[0].Order != 0 {
			t.Errorf("Expected first prescription order 0, got %d", day.Prescriptions[0].Order)
		}
		if day.Prescriptions[1].Order != 1 {
			t.Errorf("Expected second prescription order 1, got %d", day.Prescriptions[1].Order)
		}
	})

	t.Run("reorders prescriptions", func(t *testing.T) {
		// Reorder: swap p2 to be first
		body := fmt.Sprintf(`{"prescriptionIds": ["%s", "%s"]}`, p2.ID, p1.ID)
		resp, _ := adminPut(ts.URL("/days/"+createdDay.ID+"/prescriptions/reorder"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope DayWithPrescriptionsEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		day := envelope.Data

		if len(day.Prescriptions) != 2 {
			t.Fatalf("Expected 2 prescriptions, got %d", len(day.Prescriptions))
		}

		// p2 should be first now (order 0)
		if day.Prescriptions[0].PrescriptionID != p2.ID {
			t.Errorf("Expected first prescription to be %s, got %s", p2.ID, day.Prescriptions[0].PrescriptionID)
		}
		if day.Prescriptions[0].Order != 0 {
			t.Errorf("Expected first prescription order 0, got %d", day.Prescriptions[0].Order)
		}

		// p1 should be second now (order 1)
		if day.Prescriptions[1].PrescriptionID != p1.ID {
			t.Errorf("Expected second prescription to be %s, got %s", p1.ID, day.Prescriptions[1].PrescriptionID)
		}
		if day.Prescriptions[1].Order != 1 {
			t.Errorf("Expected second prescription order 1, got %d", day.Prescriptions[1].Order)
		}
	})

	t.Run("removes prescription from day", func(t *testing.T) {
		resp, err := adminDelete(ts.URL("/days/" + createdDay.ID + "/prescriptions/" + p1.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 204, got %d: %s", resp.StatusCode, bodyBytes)
		}

		// Verify it's removed
		getResp, _ := authGet(ts.URL("/days/" + createdDay.ID))
		defer getResp.Body.Close()

		var envelope DayWithPrescriptionsEnvelope
		json.NewDecoder(getResp.Body).Decode(&envelope)
		day := envelope.Data

		if len(day.Prescriptions) != 1 {
			t.Errorf("Expected 1 prescription after removal, got %d", len(day.Prescriptions))
		}
	})

	t.Run("returns 404 for non-existent prescription in day", func(t *testing.T) {
		resp, _ := adminDelete(ts.URL("/days/" + createdDay.ID + "/prescriptions/non-existent-id"))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestDayResponseFormat(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a day with metadata
	createBody := `{"name": "Format Test", "slug": "format-test", "metadata": {"intensityLevel": "HEAVY", "focus": "squats"}}`
	createResp, _ := adminPost(ts.URL("/days"), createBody)
	var createEnvelope DayEnvelope
	json.NewDecoder(createResp.Body).Decode(&createEnvelope)
	createdDay := createEnvelope.Data
	createResp.Body.Close()

	t.Run("response has correct JSON field names", func(t *testing.T) {
		resp, _ := authGet(ts.URL("/days/" + createdDay.ID))
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// Check camelCase field names per ERD spec
		expectedFields := []string{
			`"id"`,
			`"name"`,
			`"slug"`,
			`"metadata"`,
			`"prescriptions"`,
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
