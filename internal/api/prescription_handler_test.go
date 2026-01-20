package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// PaginatedPrescriptionsResponse is the paginated list response.
type PaginatedPrescriptionsResponse struct {
	Data       []PrescriptionResponse `json:"data"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"pageSize"`
	TotalItems int64                  `json:"totalItems"`
	TotalPages int64                  `json:"totalPages"`
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
		if result.TotalItems != 0 {
			t.Errorf("Expected totalItems 0, got %d", result.TotalItems)
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

		// Get page 1
		resp, err := authGet(ts.URL("/prescriptions?page=1&pageSize=2"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		var result PaginatedPrescriptionsResponse
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result.Data) != 2 {
			t.Errorf("Expected 2 prescriptions on page 1, got %d", len(result.Data))
		}
		if result.Page != 1 {
			t.Errorf("Expected page 1, got %d", result.Page)
		}
		if result.PageSize != 2 {
			t.Errorf("Expected pageSize 2, got %d", result.PageSize)
		}
		if result.TotalItems != 5 {
			t.Errorf("Expected totalItems 5, got %d", result.TotalItems)
		}
		if result.TotalPages != 3 {
			t.Errorf("Expected totalPages 3, got %d", result.TotalPages)
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
	var createdPrescription PrescriptionResponse
	json.NewDecoder(createResp.Body).Decode(&createdPrescription)
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

		var p PrescriptionResponse
		json.NewDecoder(resp.Body).Decode(&p)

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

		if errResp.Error != "Prescription not found" {
			t.Errorf("Expected error 'Prescription not found', got %s", errResp.Error)
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

		var p PrescriptionResponse
		json.NewDecoder(resp.Body).Decode(&p)

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

		var p PrescriptionResponse
		json.NewDecoder(resp.Body).Decode(&p)

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

		var p PrescriptionResponse
		json.NewDecoder(resp.Body).Decode(&p)

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
	var createdPrescription PrescriptionResponse
	json.NewDecoder(createResp.Body).Decode(&createdPrescription)
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

		var p PrescriptionResponse
		json.NewDecoder(resp.Body).Decode(&p)

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

		var p PrescriptionResponse
		json.NewDecoder(resp.Body).Decode(&p)

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

		var p PrescriptionResponse
		json.NewDecoder(resp.Body).Decode(&p)

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

		var p PrescriptionResponse
		json.NewDecoder(resp.Body).Decode(&p)

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

		var p PrescriptionResponse
		json.NewDecoder(resp.Body).Decode(&p)

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

		var p PrescriptionResponse
		json.NewDecoder(resp.Body).Decode(&p)

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

		var p PrescriptionResponse
		json.NewDecoder(resp.Body).Decode(&p)

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
		var createdPrescription PrescriptionResponse
		json.NewDecoder(createResp.Body).Decode(&createdPrescription)
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
	var createdPrescription PrescriptionResponse
	json.NewDecoder(createResp.Body).Decode(&createdPrescription)
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

		var p PrescriptionResponse
		json.NewDecoder(resp.Body).Decode(&p)

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

		var p PrescriptionResponse
		json.NewDecoder(resp.Body).Decode(&p)

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
