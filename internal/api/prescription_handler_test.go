package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// PrescriptionResponse matches the API response format.
type PrescriptionResponse struct {
	ID           string          `json:"id"`
	LiftID       string          `json:"liftId"`
	LoadStrategy json.RawMessage `json:"loadStrategy"`
	SetScheme    json.RawMessage `json:"setScheme"`
	Order        int             `json:"order"`
	Notes        string          `json:"notes,omitempty"`
	RestSeconds  *int            `json:"restSeconds,omitempty"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

// PrescriptionPaginationMeta contains pagination metadata.
type PrescriptionPaginationMeta struct {
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasMore bool  `json:"hasMore"`
}

// PaginatedPrescriptionsResponse is the paginated list response.
type PaginatedPrescriptionsResponse struct {
	Data []PrescriptionResponse      `json:"data"`
	Meta *PrescriptionPaginationMeta `json:"meta"`
}

// PrescriptionEnvelope wraps single prescription response with standard envelope.
type PrescriptionEnvelope struct {
	Data PrescriptionResponse `json:"data"`
}

// LoadStrategyResponse represents the load strategy in responses.
type LoadStrategyResponse struct {
	Type              string  `json:"type"`
	ReferenceType     string  `json:"referenceType"`
	Percentage        float64 `json:"percentage"`
	RoundingIncrement float64 `json:"roundingIncrement,omitempty"`
	RoundingDirection string  `json:"roundingDirection,omitempty"`
}

// SetSchemeResponse represents the set scheme in responses.
type SetSchemeResponse struct {
	Type string `json:"type"`
	Sets int    `json:"sets,omitempty"`
	Reps int    `json:"reps,omitempty"`
}

// Test data: use seeded squat ID
const seededSquatID = "00000000-0000-0000-0000-000000000001"

func TestListPrescriptions(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("returns empty list initially", func(t *testing.T) {
		resp, err := authGet(ts.URL("/prescriptions"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, body)
		}

		var result PaginatedPrescriptionsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(result.Data) != 0 {
			t.Errorf("Expected 0 prescriptions, got %d", len(result.Data))
		}
		if result.Meta == nil || result.Meta.Total != 0 {
			total := int64(0)
			if result.Meta != nil {
				total = result.Meta.Total
			}
			t.Errorf("Expected total 0, got %d", total)
		}
	})

	t.Run("returns prescriptions with pagination", func(t *testing.T) {
		// Create multiple prescriptions
		for i := 0; i < 5; i++ {
			body := fmt.Sprintf(`{
				"liftId": "%s",
				"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 75},
				"setScheme": {"type": "FIXED", "sets": 5, "reps": 5},
				"order": %d
			}`, seededSquatID, i)
			resp, _ := adminPost(ts.URL("/prescriptions"), body)
			resp.Body.Close()
		}

		// Get page 1 (offset=0, limit=2)
		resp, err := authGet(ts.URL("/prescriptions?limit=2&offset=0"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedPrescriptionsResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) != 2 {
			t.Errorf("Expected 2 prescriptions on page 1, got %d", len(result.Data))
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
		if result.Meta.Total != 5 {
			t.Errorf("Expected total 5, got %d", result.Meta.Total)
		}
		if !result.Meta.HasMore {
			t.Error("Expected hasMore to be true")
		}
	})
}

func TestListPrescriptionsFilterByLift(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	benchID := "00000000-0000-0000-0000-000000000002"

	// Create prescriptions for squat
	for i := 0; i < 3; i++ {
		body := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 75},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5},
			"order": %d
		}`, seededSquatID, i)
		resp, _ := adminPost(ts.URL("/prescriptions"), body)
		resp.Body.Close()
	}

	// Create prescriptions for bench
	for i := 0; i < 2; i++ {
		body := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "ONE_RM", "percentage": 80},
			"setScheme": {"type": "FIXED", "sets": 3, "reps": 8},
			"order": %d
		}`, benchID, i)
		resp, _ := adminPost(ts.URL("/prescriptions"), body)
		resp.Body.Close()
	}

	t.Run("filters by lift_id", func(t *testing.T) {
		resp, err := authGet(ts.URL("/prescriptions?lift_id=" + seededSquatID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedPrescriptionsResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) != 3 {
			t.Errorf("Expected 3 squat prescriptions, got %d", len(result.Data))
		}

		for _, p := range result.Data {
			if p.LiftID != seededSquatID {
				t.Errorf("Expected all prescriptions to have liftId %s, got %s", seededSquatID, p.LiftID)
			}
		}
	})
}

func TestListPrescriptionsSorting(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create prescriptions with different orders
	orders := []int{3, 1, 2}
	for _, order := range orders {
		body := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 75},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5},
			"order": %d
		}`, seededSquatID, order)
		resp, _ := adminPost(ts.URL("/prescriptions"), body)
		resp.Body.Close()
	}

	t.Run("sorts by order ascending by default", func(t *testing.T) {
		resp, _ := authGet(ts.URL("/prescriptions"))
		defer resp.Body.Close()

		var result PaginatedPrescriptionsResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) < 2 {
			t.Fatal("Need at least 2 prescriptions for sort test")
		}

		// Check sorted order ascending
		for i := 1; i < len(result.Data); i++ {
			if result.Data[i].Order < result.Data[i-1].Order {
				t.Errorf("Prescriptions not sorted by order ascending")
			}
		}
	})

	t.Run("sorts by order descending", func(t *testing.T) {
		resp, _ := authGet(ts.URL("/prescriptions?sortBy=order&sortOrder=desc"))
		defer resp.Body.Close()

		var result PaginatedPrescriptionsResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) < 2 {
			t.Fatal("Need at least 2 prescriptions for sort test")
		}

		// Check sorted order descending
		for i := 1; i < len(result.Data); i++ {
			if result.Data[i].Order > result.Data[i-1].Order {
				t.Errorf("Prescriptions not sorted by order descending")
			}
		}
	})
}

func TestGetPrescription(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a prescription
	createBody := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85, "roundingIncrement": 5, "roundingDirection": "NEAREST"},
		"setScheme": {"type": "FIXED", "sets": 5, "reps": 5},
		"order": 1,
		"notes": "Focus on depth",
		"restSeconds": 180
	}`, seededSquatID)
	createResp, _ := adminPost(ts.URL("/prescriptions"), createBody)
	var createEnvelope PrescriptionEnvelope
	json.NewDecoder(createResp.Body).Decode(&createEnvelope)
	createdPrescription := createEnvelope.Data
	createResp.Body.Close()

	t.Run("returns prescription by ID", func(t *testing.T) {
		resp, err := authGet(ts.URL("/prescriptions/" + createdPrescription.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var envelope PrescriptionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		p := envelope.Data

		if p.ID != createdPrescription.ID {
			t.Errorf("Expected ID %s, got %s", createdPrescription.ID, p.ID)
		}
		if p.LiftID != seededSquatID {
			t.Errorf("Expected liftId %s, got %s", seededSquatID, p.LiftID)
		}
		if p.Order != 1 {
			t.Errorf("Expected order 1, got %d", p.Order)
		}
		if p.Notes != "Focus on depth" {
			t.Errorf("Expected notes 'Focus on depth', got %s", p.Notes)
		}
		if p.RestSeconds == nil || *p.RestSeconds != 180 {
			t.Errorf("Expected restSeconds 180, got %v", p.RestSeconds)
		}

		// Check load strategy
		var ls LoadStrategyResponse
		json.Unmarshal(p.LoadStrategy, &ls)
		if ls.Type != "PERCENT_OF" {
			t.Errorf("Expected loadStrategy type PERCENT_OF, got %s", ls.Type)
		}
		if ls.ReferenceType != "TRAINING_MAX" {
			t.Errorf("Expected referenceType TRAINING_MAX, got %s", ls.ReferenceType)
		}
		if ls.Percentage != 85 {
			t.Errorf("Expected percentage 85, got %f", ls.Percentage)
		}

		// Check set scheme
		var ss SetSchemeResponse
		json.Unmarshal(p.SetScheme, &ss)
		if ss.Type != "FIXED" {
			t.Errorf("Expected setScheme type FIXED, got %s", ss.Type)
		}
		if ss.Sets != 5 {
			t.Errorf("Expected sets 5, got %d", ss.Sets)
		}
		if ss.Reps != 5 {
			t.Errorf("Expected reps 5, got %d", ss.Reps)
		}
	})

	t.Run("returns 404 for non-existent ID", func(t *testing.T) {
		resp, _ := authGet(ts.URL("/prescriptions/non-existent-id"))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)

		if !strings.Contains(errResp.Error.Message, "prescription not found") {
			t.Errorf("Expected error to contain 'prescription not found', got %s", errResp.Error.Message)
		}
	})
}

func TestCreatePrescription(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("creates prescription with all fields", func(t *testing.T) {
		body := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85, "roundingIncrement": 5.0, "roundingDirection": "NEAREST"},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5},
			"order": 1,
			"notes": "Main work",
			"restSeconds": 180
		}`, seededSquatID)
		resp, err := adminPost(ts.URL("/prescriptions"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, body)
		}

		var envelope PrescriptionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		p := envelope.Data

		if p.ID == "" {
			t.Errorf("Expected ID to be generated")
		}
		if p.LiftID != seededSquatID {
			t.Errorf("Expected liftId %s, got %s", seededSquatID, p.LiftID)
		}
		if p.Order != 1 {
			t.Errorf("Expected order 1, got %d", p.Order)
		}
		if p.Notes != "Main work" {
			t.Errorf("Expected notes 'Main work', got %s", p.Notes)
		}
		if p.RestSeconds == nil || *p.RestSeconds != 180 {
			t.Errorf("Expected restSeconds 180, got %v", p.RestSeconds)
		}
	})

	t.Run("creates prescription with minimal fields", func(t *testing.T) {
		body := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "ONE_RM", "percentage": 75},
			"setScheme": {"type": "FIXED", "sets": 3, "reps": 8}
		}`, seededSquatID)
		resp, _ := adminPost(ts.URL("/prescriptions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope PrescriptionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		p := envelope.Data

		// Order defaults to 0
		if p.Order != 0 {
			t.Errorf("Expected order 0 (default), got %d", p.Order)
		}
		// Notes should be empty
		if p.Notes != "" {
			t.Errorf("Expected empty notes, got %s", p.Notes)
		}
		// RestSeconds should be nil
		if p.RestSeconds != nil {
			t.Errorf("Expected nil restSeconds, got %v", p.RestSeconds)
		}
	})

	t.Run("creates prescription with RAMP set scheme", func(t *testing.T) {
		body := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100},
			"setScheme": {"type": "RAMP", "steps": [{"percentage": 50, "reps": 5}, {"percentage": 60, "reps": 5}, {"percentage": 70, "reps": 5}, {"percentage": 80, "reps": 3}, {"percentage": 90, "reps": 1}], "workSetThreshold": 80},
			"order": 0
		}`, seededSquatID)
		resp, _ := adminPost(ts.URL("/prescriptions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope PrescriptionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		p := envelope.Data

		// Verify RAMP scheme was stored
		var ss map[string]interface{}
		json.Unmarshal(p.SetScheme, &ss)
		if ss["type"] != "RAMP" {
			t.Errorf("Expected setScheme type RAMP, got %v", ss["type"])
		}
	})

	t.Run("returns 400 for missing liftId", func(t *testing.T) {
		body := `{
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5}
		}`
		resp, _ := adminPost(ts.URL("/prescriptions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for non-existent lift", func(t *testing.T) {
		body := `{
			"liftId": "non-existent-lift-id",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5}
		}`
		resp, _ := adminPost(ts.URL("/prescriptions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid loadStrategy", func(t *testing.T) {
		body := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "INVALID_TYPE"},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5}
		}`, seededSquatID)
		resp, _ := adminPost(ts.URL("/prescriptions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid setScheme", func(t *testing.T) {
		body := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85},
			"setScheme": {"type": "INVALID_SCHEME"}
		}`, seededSquatID)
		resp, _ := adminPost(ts.URL("/prescriptions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid percentage", func(t *testing.T) {
		body := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": -5},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5}
		}`, seededSquatID)
		resp, _ := adminPost(ts.URL("/prescriptions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid sets", func(t *testing.T) {
		body := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85},
			"setScheme": {"type": "FIXED", "sets": 0, "reps": 5}
		}`, seededSquatID)
		resp, _ := adminPost(ts.URL("/prescriptions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for negative order", func(t *testing.T) {
		body := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5},
			"order": -1
		}`, seededSquatID)
		resp, _ := adminPost(ts.URL("/prescriptions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for notes too long", func(t *testing.T) {
		longNotes := make([]byte, 501)
		for i := range longNotes {
			longNotes[i] = 'a'
		}
		body := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5},
			"notes": "%s"
		}`, seededSquatID, string(longNotes))
		resp, _ := adminPost(ts.URL("/prescriptions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for negative restSeconds", func(t *testing.T) {
		body := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5},
			"restSeconds": -30
		}`, seededSquatID)
		resp, _ := adminPost(ts.URL("/prescriptions"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		resp, _ := adminPost(ts.URL("/prescriptions"), "{invalid json}")
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestUpdatePrescription(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a prescription to update
	createBody := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85},
		"setScheme": {"type": "FIXED", "sets": 5, "reps": 5},
		"order": 1,
		"notes": "Original notes",
		"restSeconds": 180
	}`, seededSquatID)
	createResp, _ := adminPost(ts.URL("/prescriptions"), createBody)
	var createEnvelope PrescriptionEnvelope
	json.NewDecoder(createResp.Body).Decode(&createEnvelope)
	createdPrescription := createEnvelope.Data
	createResp.Body.Close()

	t.Run("updates prescription notes", func(t *testing.T) {
		body := `{"notes": "Updated notes"}`
		resp, err := adminPut(ts.URL("/prescriptions/"+createdPrescription.ID), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope PrescriptionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		p := envelope.Data

		if p.Notes != "Updated notes" {
			t.Errorf("Expected notes 'Updated notes', got %s", p.Notes)
		}
		// Other fields should remain unchanged
		if p.Order != 1 {
			t.Errorf("Expected order 1, got %d", p.Order)
		}
	})

	t.Run("updates prescription order", func(t *testing.T) {
		body := `{"order": 5}`
		resp, _ := adminPut(ts.URL("/prescriptions/"+createdPrescription.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var envelope PrescriptionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		p := envelope.Data

		if p.Order != 5 {
			t.Errorf("Expected order 5, got %d", p.Order)
		}
	})

	t.Run("updates prescription restSeconds", func(t *testing.T) {
		body := `{"restSeconds": 240}`
		resp, _ := adminPut(ts.URL("/prescriptions/"+createdPrescription.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var envelope PrescriptionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		p := envelope.Data

		if p.RestSeconds == nil || *p.RestSeconds != 240 {
			t.Errorf("Expected restSeconds 240, got %v", p.RestSeconds)
		}
	})

	t.Run("clears prescription restSeconds", func(t *testing.T) {
		body := `{"clearRestSeconds": true}`
		resp, _ := adminPut(ts.URL("/prescriptions/"+createdPrescription.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var envelope PrescriptionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		p := envelope.Data

		if p.RestSeconds != nil {
			t.Errorf("Expected restSeconds to be cleared, got %v", p.RestSeconds)
		}
	})

	t.Run("updates prescription loadStrategy", func(t *testing.T) {
		body := `{"loadStrategy": {"type": "PERCENT_OF", "referenceType": "ONE_RM", "percentage": 90}}`
		resp, _ := adminPut(ts.URL("/prescriptions/"+createdPrescription.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope PrescriptionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		p := envelope.Data

		var ls LoadStrategyResponse
		json.Unmarshal(p.LoadStrategy, &ls)
		if ls.ReferenceType != "ONE_RM" {
			t.Errorf("Expected referenceType ONE_RM, got %s", ls.ReferenceType)
		}
		if ls.Percentage != 90 {
			t.Errorf("Expected percentage 90, got %f", ls.Percentage)
		}
	})

	t.Run("updates prescription setScheme", func(t *testing.T) {
		body := `{"setScheme": {"type": "FIXED", "sets": 3, "reps": 10}}`
		resp, _ := adminPut(ts.URL("/prescriptions/"+createdPrescription.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope PrescriptionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		p := envelope.Data

		var ss SetSchemeResponse
		json.Unmarshal(p.SetScheme, &ss)
		if ss.Sets != 3 {
			t.Errorf("Expected sets 3, got %d", ss.Sets)
		}
		if ss.Reps != 10 {
			t.Errorf("Expected reps 10, got %d", ss.Reps)
		}
	})

	t.Run("updates prescription liftId", func(t *testing.T) {
		benchID := "00000000-0000-0000-0000-000000000002"
		body := fmt.Sprintf(`{"liftId": "%s"}`, benchID)
		resp, _ := adminPut(ts.URL("/prescriptions/"+createdPrescription.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope PrescriptionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		p := envelope.Data

		if p.LiftID != benchID {
			t.Errorf("Expected liftId %s, got %s", benchID, p.LiftID)
		}
	})

	t.Run("returns 404 for non-existent prescription", func(t *testing.T) {
		body := `{"notes": "Updated"}`
		resp, _ := adminPut(ts.URL("/prescriptions/non-existent-id"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for non-existent liftId", func(t *testing.T) {
		body := `{"liftId": "non-existent-lift-id"}`
		resp, _ := adminPut(ts.URL("/prescriptions/"+createdPrescription.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid loadStrategy", func(t *testing.T) {
		body := `{"loadStrategy": {"type": "INVALID"}}`
		resp, _ := adminPut(ts.URL("/prescriptions/"+createdPrescription.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for invalid setScheme", func(t *testing.T) {
		body := `{"setScheme": {"type": "INVALID"}}`
		resp, _ := adminPut(ts.URL("/prescriptions/"+createdPrescription.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for validation errors", func(t *testing.T) {
		body := `{"order": -1}`
		resp, _ := adminPut(ts.URL("/prescriptions/"+createdPrescription.ID), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestDeletePrescription(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("deletes prescription successfully", func(t *testing.T) {
		// Create a prescription to delete
		createBody := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5}
		}`, seededSquatID)
		createResp, _ := adminPost(ts.URL("/prescriptions"), createBody)
		var createEnvelope PrescriptionEnvelope
		json.NewDecoder(createResp.Body).Decode(&createEnvelope)
		createdPrescription := createEnvelope.Data
		createResp.Body.Close()

		// Delete it
		resp, err := adminDelete(ts.URL("/prescriptions/" + createdPrescription.ID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 204, got %d: %s", resp.StatusCode, bodyBytes)
		}

		// Verify it's deleted
		getResp, _ := authGet(ts.URL("/prescriptions/" + createdPrescription.ID))
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected deleted prescription to return 404, got %d", getResp.StatusCode)
		}
	})

	t.Run("returns 404 for non-existent prescription", func(t *testing.T) {
		resp, _ := adminDelete(ts.URL("/prescriptions/non-existent-id"))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestPrescriptionAuth(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("returns 401 for unauthenticated requests", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/prescriptions"), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 403 when non-admin tries to create", func(t *testing.T) {
		body := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5}
		}`, seededSquatID)
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/prescriptions"), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", testutil.TestUserID) // Non-admin user
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("allows authenticated user to read prescriptions", func(t *testing.T) {
		resp, err := authGet(ts.URL("/prescriptions"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}

func TestPrescriptionResponseFormat(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a prescription
	createBody := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85, "roundingIncrement": 5.0, "roundingDirection": "NEAREST"},
		"setScheme": {"type": "FIXED", "sets": 5, "reps": 5},
		"order": 1,
		"notes": "Test notes",
		"restSeconds": 180
	}`, seededSquatID)
	createResp, _ := adminPost(ts.URL("/prescriptions"), createBody)
	var createEnvelope PrescriptionEnvelope
	json.NewDecoder(createResp.Body).Decode(&createEnvelope)
	createdPrescription := createEnvelope.Data
	createResp.Body.Close()

	t.Run("response has correct JSON field names", func(t *testing.T) {
		resp, _ := authGet(ts.URL("/prescriptions/" + createdPrescription.ID))
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// Check camelCase field names per API spec
		expectedFields := []string{
			`"id"`,
			`"liftId"`,
			`"loadStrategy"`,
			`"setScheme"`,
			`"order"`,
			`"notes"`,
			`"restSeconds"`,
			`"createdAt"`,
			`"updatedAt"`,
		}

		for _, field := range expectedFields {
			if !bytes.Contains(body, []byte(field)) {
				t.Errorf("Expected field %s in response, body: %s", field, bodyStr)
			}
		}
	})

	t.Run("loadStrategy has correct structure", func(t *testing.T) {
		resp, _ := authGet(ts.URL("/prescriptions/" + createdPrescription.ID))
		defer resp.Body.Close()

		var envelope PrescriptionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		p := envelope.Data

		var ls map[string]interface{}
		json.Unmarshal(p.LoadStrategy, &ls)

		if ls["type"] != "PERCENT_OF" {
			t.Errorf("Expected loadStrategy.type = PERCENT_OF, got %v", ls["type"])
		}
		if ls["referenceType"] != "TRAINING_MAX" {
			t.Errorf("Expected loadStrategy.referenceType = TRAINING_MAX, got %v", ls["referenceType"])
		}
		if ls["percentage"].(float64) != 85 {
			t.Errorf("Expected loadStrategy.percentage = 85, got %v", ls["percentage"])
		}
	})

	t.Run("setScheme has correct structure", func(t *testing.T) {
		resp, _ := authGet(ts.URL("/prescriptions/" + createdPrescription.ID))
		defer resp.Body.Close()

		var envelope PrescriptionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		p := envelope.Data

		var ss map[string]interface{}
		json.Unmarshal(p.SetScheme, &ss)

		if ss["type"] != "FIXED" {
			t.Errorf("Expected setScheme.type = FIXED, got %v", ss["type"])
		}
		if int(ss["sets"].(float64)) != 5 {
			t.Errorf("Expected setScheme.sets = 5, got %v", ss["sets"])
		}
		if int(ss["reps"].(float64)) != 5 {
			t.Errorf("Expected setScheme.reps = 5, got %v", ss["reps"])
		}
	})
}

// ResolvedPrescriptionTestResponse matches the resolved prescription API response format.
type ResolvedPrescriptionTestResponse struct {
	PrescriptionID string `json:"prescriptionId"`
	Lift           struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"lift"`
	Sets []struct {
		SetNumber  int     `json:"setNumber"`
		Weight     float64 `json:"weight"`
		TargetReps int     `json:"targetReps"`
		IsWorkSet  bool    `json:"isWorkSet"`
	} `json:"sets"`
	Notes       string `json:"notes,omitempty"`
	RestSeconds *int   `json:"restSeconds,omitempty"`
}

// ResolvedPrescriptionEnvelope wraps single resolved prescription response with standard envelope.
type ResolvedPrescriptionEnvelope struct {
	Data ResolvedPrescriptionTestResponse `json:"data"`
}

// BatchResolveResultItem matches a single item in the batch resolution response.
type BatchResolveResultItem struct {
	PrescriptionID string                            `json:"prescriptionId"`
	Status         string                            `json:"status"`
	Resolved       *ResolvedPrescriptionTestResponse `json:"resolved,omitempty"`
	Error          string                            `json:"error,omitempty"`
}

// BatchResolveTestResponse matches the batch resolution API response format.
type BatchResolveTestResponse struct {
	Results []BatchResolveResultItem `json:"results"`
}

// BatchResolveEnvelope wraps batch resolution response with standard envelope.
type BatchResolveEnvelope struct {
	Data BatchResolveTestResponse `json:"data"`
}

// authPost performs an authenticated POST request
func authPost(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestUserID)
	return http.DefaultClient.Do(req)
}

func TestResolvePrescription(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a prescription
	createBody := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85, "roundingIncrement": 5, "roundingDirection": "NEAREST"},
		"setScheme": {"type": "FIXED", "sets": 5, "reps": 5},
		"order": 1,
		"notes": "Focus on depth",
		"restSeconds": 180
	}`, seededSquatID)
	createResp, _ := adminPost(ts.URL("/prescriptions"), createBody)
	var createEnvelope PrescriptionEnvelope
	json.NewDecoder(createResp.Body).Decode(&createEnvelope)
	createdPrescription := createEnvelope.Data
	createResp.Body.Close()

	// Create a training max for the user
	maxBody := fmt.Sprintf(`{
		"liftId": "%s",
		"type": "TRAINING_MAX",
		"value": 300,
		"effectiveDate": "2025-01-15T00:00:00Z"
	}`, seededSquatID)
	maxResp, _ := adminPost(ts.URL("/users/"+testutil.TestUserID+"/lift-maxes"), maxBody)
	maxResp.Body.Close()

	t.Run("resolves prescription successfully", func(t *testing.T) {
		body := fmt.Sprintf(`{"userId": "%s"}`, testutil.TestUserID)
		resp, err := authPost(ts.URL("/prescriptions/"+createdPrescription.ID+"/resolve"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope ResolvedPrescriptionEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		resolved := envelope.Data

		if resolved.PrescriptionID != createdPrescription.ID {
			t.Errorf("Expected prescriptionId %s, got %s", createdPrescription.ID, resolved.PrescriptionID)
		}

		if resolved.Lift.ID != seededSquatID {
			t.Errorf("Expected lift id %s, got %s", seededSquatID, resolved.Lift.ID)
		}

		if resolved.Lift.Name != "Squat" {
			t.Errorf("Expected lift name 'Squat', got %s", resolved.Lift.Name)
		}

		if len(resolved.Sets) != 5 {
			t.Errorf("Expected 5 sets, got %d", len(resolved.Sets))
		}

		// Check calculated weight (85% of 300 = 255, rounded to 5lb increment)
		expectedWeight := 255.0
		if len(resolved.Sets) > 0 && resolved.Sets[0].Weight != expectedWeight {
			t.Errorf("Expected weight %f, got %f", expectedWeight, resolved.Sets[0].Weight)
		}

		if resolved.Notes != "Focus on depth" {
			t.Errorf("Expected notes 'Focus on depth', got %s", resolved.Notes)
		}

		if resolved.RestSeconds == nil || *resolved.RestSeconds != 180 {
			t.Errorf("Expected restSeconds 180, got %v", resolved.RestSeconds)
		}
	})

	t.Run("returns 404 for non-existent prescription", func(t *testing.T) {
		body := fmt.Sprintf(`{"userId": "%s"}`, testutil.TestUserID)
		resp, _ := authPost(ts.URL("/prescriptions/non-existent-id/resolve"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for missing userId", func(t *testing.T) {
		resp, _ := authPost(ts.URL("/prescriptions/"+createdPrescription.ID+"/resolve"), `{}`)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 422 when max not found", func(t *testing.T) {
		// Create a prescription with a different lift that has no max
		benchID := "00000000-0000-0000-0000-000000000002"
		createBody := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 75},
			"setScheme": {"type": "FIXED", "sets": 3, "reps": 8}
		}`, benchID)
		createResp, _ := adminPost(ts.URL("/prescriptions"), createBody)
		var benchEnvelope PrescriptionEnvelope
		json.NewDecoder(createResp.Body).Decode(&benchEnvelope)
		benchPrescription := benchEnvelope.Data
		createResp.Body.Close()

		// Try to resolve without a max
		body := fmt.Sprintf(`{"userId": "%s"}`, testutil.TestUserID)
		resp, _ := authPost(ts.URL("/prescriptions/"+benchPrescription.ID+"/resolve"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 422, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)

		if !strings.Contains(errResp.Error.Message, "max not found") && !strings.Contains(errResp.Error.Message, "No TRAINING_MAX found") {
			t.Errorf("Expected error about max not found, got: %s", errResp.Error.Message)
		}
	})

	t.Run("returns 401 for unauthenticated request", func(t *testing.T) {
		body := fmt.Sprintf(`{"userId": "%s"}`, testutil.TestUserID)
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/prescriptions/"+createdPrescription.ID+"/resolve"), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})
}

func TestResolvePrescriptionBatch(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	benchID := "00000000-0000-0000-0000-000000000002"

	// Create training max for squat
	maxBody := fmt.Sprintf(`{
		"liftId": "%s",
		"type": "TRAINING_MAX",
		"value": 300,
		"effectiveDate": "2025-01-15T00:00:00Z"
	}`, seededSquatID)
	maxResp, _ := adminPost(ts.URL("/users/"+testutil.TestUserID+"/lift-maxes"), maxBody)
	maxResp.Body.Close()

	// Create training max for bench
	maxBody2 := fmt.Sprintf(`{
		"liftId": "%s",
		"type": "TRAINING_MAX",
		"value": 200,
		"effectiveDate": "2025-01-15T00:00:00Z"
	}`, benchID)
	maxResp2, _ := adminPost(ts.URL("/users/"+testutil.TestUserID+"/lift-maxes"), maxBody2)
	maxResp2.Body.Close()

	// Create multiple prescriptions
	var prescriptionIDs []string

	// Squat prescription 1
	createBody := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 75},
		"setScheme": {"type": "FIXED", "sets": 5, "reps": 5},
		"order": 0,
		"notes": "Warmup"
	}`, seededSquatID)
	createResp, _ := adminPost(ts.URL("/prescriptions"), createBody)
	var p1Envelope PrescriptionEnvelope
	json.NewDecoder(createResp.Body).Decode(&p1Envelope)
	createResp.Body.Close()
	prescriptionIDs = append(prescriptionIDs, p1Envelope.Data.ID)

	// Squat prescription 2 (using same lift/max - should hit cache)
	createBody2 := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85},
		"setScheme": {"type": "FIXED", "sets": 3, "reps": 3},
		"order": 1,
		"notes": "Work sets"
	}`, seededSquatID)
	createResp2, _ := adminPost(ts.URL("/prescriptions"), createBody2)
	var p2Envelope PrescriptionEnvelope
	json.NewDecoder(createResp2.Body).Decode(&p2Envelope)
	createResp2.Body.Close()
	prescriptionIDs = append(prescriptionIDs, p2Envelope.Data.ID)

	// Bench prescription
	createBody3 := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 80},
		"setScheme": {"type": "FIXED", "sets": 4, "reps": 8},
		"order": 2
	}`, benchID)
	createResp3, _ := adminPost(ts.URL("/prescriptions"), createBody3)
	var p3Envelope PrescriptionEnvelope
	json.NewDecoder(createResp3.Body).Decode(&p3Envelope)
	createResp3.Body.Close()
	prescriptionIDs = append(prescriptionIDs, p3Envelope.Data.ID)

	t.Run("batch resolves all prescriptions successfully", func(t *testing.T) {
		body := fmt.Sprintf(`{"prescriptionIds": ["%s", "%s", "%s"], "userId": "%s"}`,
			prescriptionIDs[0], prescriptionIDs[1], prescriptionIDs[2], testutil.TestUserID)
		resp, err := authPost(ts.URL("/prescriptions/resolve-batch"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var batchEnvelope BatchResolveEnvelope
		json.NewDecoder(resp.Body).Decode(&batchEnvelope)
		batchResp := batchEnvelope.Data

		if len(batchResp.Results) != 3 {
			t.Fatalf("Expected 3 results, got %d", len(batchResp.Results))
		}

		// Check all results are successful
		for i, result := range batchResp.Results {
			if result.Status != "success" {
				t.Errorf("Expected result %d to be 'success', got '%s' (error: %s)", i, result.Status, result.Error)
			}
			if result.Resolved == nil {
				t.Errorf("Expected result %d to have resolved data", i)
				continue
			}
		}

		// Check specific weights
		// First: 75% of 300 = 225
		if batchResp.Results[0].Resolved != nil && batchResp.Results[0].Resolved.Sets[0].Weight != 225 {
			t.Errorf("Expected first prescription weight 225, got %f", batchResp.Results[0].Resolved.Sets[0].Weight)
		}

		// Second: 85% of 300 = 255
		if batchResp.Results[1].Resolved != nil && batchResp.Results[1].Resolved.Sets[0].Weight != 255 {
			t.Errorf("Expected second prescription weight 255, got %f", batchResp.Results[1].Resolved.Sets[0].Weight)
		}

		// Third: 80% of 200 = 160
		if batchResp.Results[2].Resolved != nil && batchResp.Results[2].Resolved.Sets[0].Weight != 160 {
			t.Errorf("Expected third prescription weight 160, got %f", batchResp.Results[2].Resolved.Sets[0].Weight)
		}
	})

	t.Run("batch handles partial failures", func(t *testing.T) {
		// Create a prescription without a corresponding max
		deadliftID := "00000000-0000-0000-0000-000000000003"
		createBody := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 80},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5}
		}`, deadliftID)
		createResp, _ := adminPost(ts.URL("/prescriptions"), createBody)
		var noMaxEnvelope PrescriptionEnvelope
		json.NewDecoder(createResp.Body).Decode(&noMaxEnvelope)
		noMaxPrescription := noMaxEnvelope.Data
		createResp.Body.Close()

		// Include one valid and one invalid prescription
		body := fmt.Sprintf(`{"prescriptionIds": ["%s", "%s", "non-existent-id"], "userId": "%s"}`,
			prescriptionIDs[0], noMaxPrescription.ID, testutil.TestUserID)
		resp, err := authPost(ts.URL("/prescriptions/resolve-batch"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var batchEnvelope BatchResolveEnvelope
		json.NewDecoder(resp.Body).Decode(&batchEnvelope)
		batchResp := batchEnvelope.Data

		if len(batchResp.Results) != 3 {
			t.Fatalf("Expected 3 results, got %d", len(batchResp.Results))
		}

		// First should succeed
		if batchResp.Results[0].Status != "success" {
			t.Errorf("Expected first result to be 'success', got '%s'", batchResp.Results[0].Status)
		}

		// Second should fail (no max)
		if batchResp.Results[1].Status != "error" {
			t.Errorf("Expected second result to be 'error', got '%s'", batchResp.Results[1].Status)
		}
		if !strings.Contains(batchResp.Results[1].Error, "max not found") && !strings.Contains(batchResp.Results[1].Error, "No TRAINING_MAX found") {
			t.Errorf("Expected error about max not found, got: %s", batchResp.Results[1].Error)
		}

		// Third should fail (not found)
		if batchResp.Results[2].Status != "error" {
			t.Errorf("Expected third result to be 'error', got '%s'", batchResp.Results[2].Status)
		}
		if !strings.Contains(batchResp.Results[2].Error, "prescription not found") {
			t.Errorf("Expected error to contain 'prescription not found', got: %s", batchResp.Results[2].Error)
		}
	})

	t.Run("returns 400 for missing userId", func(t *testing.T) {
		body := fmt.Sprintf(`{"prescriptionIds": ["%s"]}`, prescriptionIDs[0])
		resp, _ := authPost(ts.URL("/prescriptions/resolve-batch"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 400 for empty prescriptionIds", func(t *testing.T) {
		body := fmt.Sprintf(`{"prescriptionIds": [], "userId": "%s"}`, testutil.TestUserID)
		resp, _ := authPost(ts.URL("/prescriptions/resolve-batch"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 401 for unauthenticated request", func(t *testing.T) {
		body := fmt.Sprintf(`{"prescriptionIds": ["%s"], "userId": "%s"}`, prescriptionIDs[0], testutil.TestUserID)
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/prescriptions/resolve-batch"), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})
}

func TestResolveWithRampScheme(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a training max
	maxBody := fmt.Sprintf(`{
		"liftId": "%s",
		"type": "TRAINING_MAX",
		"value": 400,
		"effectiveDate": "2025-01-15T00:00:00Z"
	}`, seededSquatID)
	maxResp, _ := adminPost(ts.URL("/users/"+testutil.TestUserID+"/lift-maxes"), maxBody)
	maxResp.Body.Close()

	// Create a prescription with RAMP scheme
	createBody := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100, "roundingIncrement": 5},
		"setScheme": {"type": "RAMP", "steps": [{"percentage": 50, "reps": 5}, {"percentage": 60, "reps": 5}, {"percentage": 70, "reps": 3}, {"percentage": 80, "reps": 1}], "workSetThreshold": 70},
		"order": 0
	}`, seededSquatID)
	createResp, _ := adminPost(ts.URL("/prescriptions"), createBody)
	var prescriptionEnvelope PrescriptionEnvelope
	json.NewDecoder(createResp.Body).Decode(&prescriptionEnvelope)
	prescription := prescriptionEnvelope.Data
	createResp.Body.Close()

	t.Run("resolves RAMP scheme correctly", func(t *testing.T) {
		body := fmt.Sprintf(`{"userId": "%s"}`, testutil.TestUserID)
		resp, err := authPost(ts.URL("/prescriptions/"+prescription.ID+"/resolve"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var resolvedEnvelope ResolvedPrescriptionEnvelope
		json.NewDecoder(resp.Body).Decode(&resolvedEnvelope)
		resolved := resolvedEnvelope.Data

		if len(resolved.Sets) != 4 {
			t.Fatalf("Expected 4 sets, got %d", len(resolved.Sets))
		}

		// Check progressive weights (100% of TM = 400)
		// Step 1: 50% = 200
		// Step 2: 60% = 240
		// Step 3: 70% = 280
		// Step 4: 80% = 320
		expectedWeights := []float64{200, 240, 280, 320}
		for i, set := range resolved.Sets {
			if set.Weight != expectedWeights[i] {
				t.Errorf("Set %d: expected weight %f, got %f", i+1, expectedWeights[i], set.Weight)
			}
			if set.SetNumber != i+1 {
				t.Errorf("Set %d: expected setNumber %d, got %d", i+1, i+1, set.SetNumber)
			}
		}

		// Check work set classification (threshold is 70%, so sets 3 and 4 are work sets)
		expectedWorkSets := []bool{false, false, true, true}
		for i, set := range resolved.Sets {
			if set.IsWorkSet != expectedWorkSets[i] {
				t.Errorf("Set %d: expected isWorkSet %v, got %v", i+1, expectedWorkSets[i], set.IsWorkSet)
			}
		}
	})
}

func TestResolveResponseFormat(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a training max
	maxBody := fmt.Sprintf(`{
		"liftId": "%s",
		"type": "TRAINING_MAX",
		"value": 315,
		"effectiveDate": "2025-01-15T00:00:00Z"
	}`, seededSquatID)
	maxResp, _ := adminPost(ts.URL("/users/"+testutil.TestUserID+"/lift-maxes"), maxBody)
	maxResp.Body.Close()

	// Create a prescription
	restSeconds := 180
	createBody := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85},
		"setScheme": {"type": "FIXED", "sets": 3, "reps": 5},
		"order": 1,
		"notes": "Test notes",
		"restSeconds": %d
	}`, seededSquatID, restSeconds)
	createResp, _ := adminPost(ts.URL("/prescriptions"), createBody)
	var prescriptionEnvelope PrescriptionEnvelope
	json.NewDecoder(createResp.Body).Decode(&prescriptionEnvelope)
	prescription := prescriptionEnvelope.Data
	createResp.Body.Close()

	t.Run("response has correct JSON field names", func(t *testing.T) {
		body := fmt.Sprintf(`{"userId": "%s"}`, testutil.TestUserID)
		resp, _ := authPost(ts.URL("/prescriptions/"+prescription.ID+"/resolve"), body)
		defer resp.Body.Close()

		bodyBytes, _ := io.ReadAll(resp.Body)

		// Check camelCase field names per ticket spec
		expectedFields := []string{
			`"prescriptionId"`,
			`"lift"`,
			`"sets"`,
			`"setNumber"`,
			`"weight"`,
			`"targetReps"`,
			`"isWorkSet"`,
			`"notes"`,
			`"restSeconds"`,
		}

		for _, field := range expectedFields {
			if !bytes.Contains(bodyBytes, []byte(field)) {
				t.Errorf("Expected field %s in response, body: %s", field, string(bodyBytes))
			}
		}
	})

	t.Run("lift object has correct structure", func(t *testing.T) {
		body := fmt.Sprintf(`{"userId": "%s"}`, testutil.TestUserID)
		resp, _ := authPost(ts.URL("/prescriptions/"+prescription.ID+"/resolve"), body)
		defer resp.Body.Close()

		var resolvedEnvelope ResolvedPrescriptionEnvelope
		json.NewDecoder(resp.Body).Decode(&resolvedEnvelope)
		resolved := resolvedEnvelope.Data

		if resolved.Lift.ID != seededSquatID {
			t.Errorf("Expected lift.id = %s, got %s", seededSquatID, resolved.Lift.ID)
		}
		if resolved.Lift.Name != "Squat" {
			t.Errorf("Expected lift.name = 'Squat', got %s", resolved.Lift.Name)
		}
		if resolved.Lift.Slug != "squat" {
			t.Errorf("Expected lift.slug = 'squat', got %s", resolved.Lift.Slug)
		}
	})
}

func TestResolveAdditionalCases(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("returns 400 for invalid JSON body", func(t *testing.T) {
		// Create a prescription first
		createBody := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5}
		}`, seededSquatID)
		createResp, _ := adminPost(ts.URL("/prescriptions"), createBody)
		var prescriptionEnvelope PrescriptionEnvelope
		json.NewDecoder(createResp.Body).Decode(&prescriptionEnvelope)
		prescription := prescriptionEnvelope.Data
		createResp.Body.Close()

		resp, _ := authPost(ts.URL("/prescriptions/"+prescription.ID+"/resolve"), "{invalid json}")
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("batch returns 400 for invalid JSON body", func(t *testing.T) {
		resp, _ := authPost(ts.URL("/prescriptions/resolve-batch"), "{invalid json}")
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("resolve with empty userId string", func(t *testing.T) {
		// Create a prescription
		createBody := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5}
		}`, seededSquatID)
		createResp, _ := adminPost(ts.URL("/prescriptions"), createBody)
		var prescriptionEnvelope PrescriptionEnvelope
		json.NewDecoder(createResp.Body).Decode(&prescriptionEnvelope)
		prescription := prescriptionEnvelope.Data
		createResp.Body.Close()

		resp, _ := authPost(ts.URL("/prescriptions/"+prescription.ID+"/resolve"), `{"userId": "   "}`)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("batch resolve with empty userId string", func(t *testing.T) {
		// Create a prescription
		createBody := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5}
		}`, seededSquatID)
		createResp, _ := adminPost(ts.URL("/prescriptions"), createBody)
		var prescriptionEnvelope PrescriptionEnvelope
		json.NewDecoder(createResp.Body).Decode(&prescriptionEnvelope)
		prescription := prescriptionEnvelope.Data
		createResp.Body.Close()

		body := fmt.Sprintf(`{"prescriptionIds": ["%s"], "userId": "   "}`, prescription.ID)
		resp, _ := authPost(ts.URL("/prescriptions/resolve-batch"), body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestCachingBehavior(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create training max for squat
	maxBody := fmt.Sprintf(`{
		"liftId": "%s",
		"type": "TRAINING_MAX",
		"value": 300,
		"effectiveDate": "2025-01-15T00:00:00Z"
	}`, seededSquatID)
	maxResp, _ := adminPost(ts.URL("/users/"+testutil.TestUserID+"/lift-maxes"), maxBody)
	maxResp.Body.Close()

	// Create multiple prescriptions using the same lift (to test cache hits)
	var prescriptionIDs []string
	for i := 0; i < 5; i++ {
		createBody := fmt.Sprintf(`{
			"liftId": "%s",
			"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": %d},
			"setScheme": {"type": "FIXED", "sets": 5, "reps": 5},
			"order": %d
		}`, seededSquatID, 70+i*5, i)
		createResp, _ := adminPost(ts.URL("/prescriptions"), createBody)
		var pEnvelope PrescriptionEnvelope
		json.NewDecoder(createResp.Body).Decode(&pEnvelope)
		createResp.Body.Close()
		prescriptionIDs = append(prescriptionIDs, pEnvelope.Data.ID)
	}

	t.Run("batch uses cached max for same lift", func(t *testing.T) {
		// All 5 prescriptions use the same lift, so the max should only be looked up once
		body := fmt.Sprintf(`{"prescriptionIds": ["%s", "%s", "%s", "%s", "%s"], "userId": "%s"}`,
			prescriptionIDs[0], prescriptionIDs[1], prescriptionIDs[2], prescriptionIDs[3], prescriptionIDs[4], testutil.TestUserID)
		resp, err := authPost(ts.URL("/prescriptions/resolve-batch"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var batchEnvelope BatchResolveEnvelope
		json.NewDecoder(resp.Body).Decode(&batchEnvelope)
		batchResp := batchEnvelope.Data

		if len(batchResp.Results) != 5 {
			t.Fatalf("Expected 5 results, got %d", len(batchResp.Results))
		}

		// All should succeed
		for i, result := range batchResp.Results {
			if result.Status != "success" {
				t.Errorf("Expected result %d to be 'success', got '%s' (error: %s)", i, result.Status, result.Error)
			}
		}

		// Verify weights are calculated correctly (cache should work)
		// 70% of 300 = 210
		// 75% of 300 = 225
		// 80% of 300 = 240
		// 85% of 300 = 255
		// 90% of 300 = 270
		expectedWeights := []float64{210, 225, 240, 255, 270}
		for i, result := range batchResp.Results {
			if result.Resolved != nil && len(result.Resolved.Sets) > 0 {
				if result.Resolved.Sets[0].Weight != expectedWeights[i] {
					t.Errorf("Result %d: expected weight %f, got %f", i, expectedWeights[i], result.Resolved.Sets[0].Weight)
				}
			}
		}
	})
}
