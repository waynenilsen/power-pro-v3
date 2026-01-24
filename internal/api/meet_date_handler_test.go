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

// MeetDateTestResponse represents the API response for meet date operations.
type MeetDateTestResponse struct {
	MeetDate     *string `json:"meet_date,omitempty"`
	DaysOut      int     `json:"days_out"`
	CurrentPhase string  `json:"current_phase"`
	WeeksToMeet  int     `json:"weeks_to_meet"`
}

// MeetDateEnvelope wraps meet date response with standard envelope.
type MeetDateEnvelope struct {
	Data MeetDateTestResponse `json:"data"`
}

// CountdownTestResponse represents the API response for countdown operations.
type CountdownTestResponse struct {
	MeetDate        *string `json:"meet_date,omitempty"`
	DaysOut         int     `json:"days_out"`
	CurrentPhase    string  `json:"current_phase"`
	PhaseWeek       int     `json:"phase_week"`
	TaperMultiplier float64 `json:"taper_multiplier"`
}

// CountdownEnvelope wraps countdown response with standard envelope.
type CountdownEnvelope struct {
	Data CountdownTestResponse `json:"data"`
}

// Helper functions for meet date tests

func userPutMeetDate(url string, body string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func userGetCountdown(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func adminPutMeetDate(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

func adminGetCountdown(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// createMeetDateTestCycle creates a test cycle and returns its ID
func createMeetDateTestCycle(t *testing.T, ts *testutil.TestServer, name string) string {
	body := `{"name": "` + name + `", "lengthWeeks": 4}`
	resp, err := adminPostCycle(ts.URL("/cycles"), body)
	if err != nil {
		t.Fatalf("Failed to create test cycle: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create test cycle, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope CycleEnvelopeForEnrollment
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// createMeetDateTestProgram creates a test program and returns its ID
func createMeetDateTestProgram(t *testing.T, ts *testutil.TestServer, name, slug, cycleID string) string {
	body := `{"name": "` + name + `", "slug": "` + slug + `", "cycleId": "` + cycleID + `"}`
	resp, err := adminPostProgram(ts.URL("/programs"), body)
	if err != nil {
		t.Fatalf("Failed to create test program: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create test program, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope ProgramEnvelopeForEnrollment
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// enrollUserForMeetDateTest enrolls a user in a program for testing
func enrollUserForMeetDateTest(t *testing.T, ts *testutil.TestServer, userID, programID string) {
	body := `{"programId": "` + programID + `"}`
	resp, err := userPostEnrollment(ts.URL("/users/"+userID+"/program"), body, userID)
	if err != nil {
		t.Fatalf("Failed to enroll user: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to enroll user, status %d: %s", resp.StatusCode, bodyBytes)
	}
}

func TestMeetDateSetAndGet(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a cycle and program for enrollment
	cycleID := createMeetDateTestCycle(t, ts, "Meet Date Test Cycle")
	programID := createMeetDateTestProgram(t, ts, "Meet Date Test Program", "meet-date-test-program", cycleID)

	userID := "meet-date-test-user"
	enrollUserForMeetDateTest(t, ts, userID, programID)

	// Use a future date for testing
	futureDate := time.Now().AddDate(0, 3, 0).Format("2006-01-02") // 3 months from now

	t.Run("sets meet date via API", func(t *testing.T) {
		body := `{"meet_date": "` + futureDate + `"}`
		resp, err := userPutMeetDate(ts.URL("/users/"+userID+"/programs/"+programID+"/state/meet-date"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope MeetDateEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		response := envelope.Data

		if response.MeetDate == nil {
			t.Error("Expected non-nil meet_date")
		} else if *response.MeetDate != futureDate {
			t.Errorf("Expected meet_date %s, got %s", futureDate, *response.MeetDate)
		}
		if response.DaysOut <= 0 {
			t.Errorf("Expected positive days_out, got %d", response.DaysOut)
		}
		if response.CurrentPhase == "" {
			t.Error("Expected non-empty current_phase")
		}
		if response.WeeksToMeet <= 0 {
			t.Errorf("Expected positive weeks_to_meet, got %d", response.WeeksToMeet)
		}
	})

	t.Run("gets countdown information", func(t *testing.T) {
		resp, err := userGetCountdown(ts.URL("/users/"+userID+"/programs/"+programID+"/state/countdown"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope CountdownEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		response := envelope.Data

		if response.MeetDate == nil {
			t.Error("Expected non-nil meet_date")
		} else if *response.MeetDate != futureDate {
			t.Errorf("Expected meet_date %s, got %s", futureDate, *response.MeetDate)
		}
		if response.DaysOut <= 0 {
			t.Errorf("Expected positive days_out, got %d", response.DaysOut)
		}
		if response.CurrentPhase == "" {
			t.Error("Expected non-empty current_phase")
		}
		if response.PhaseWeek < 1 {
			t.Errorf("Expected phase_week >= 1, got %d", response.PhaseWeek)
		}
		if response.TaperMultiplier <= 0 || response.TaperMultiplier > 1 {
			t.Errorf("Expected taper_multiplier between 0 and 1, got %f", response.TaperMultiplier)
		}
	})

	t.Run("clears meet date via API", func(t *testing.T) {
		body := `{"meet_date": null}`
		resp, err := userPutMeetDate(ts.URL("/users/"+userID+"/programs/"+programID+"/state/meet-date"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope MeetDateEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		response := envelope.Data

		if response.MeetDate != nil {
			t.Errorf("Expected nil meet_date, got %s", *response.MeetDate)
		}
		if response.DaysOut != 0 {
			t.Errorf("Expected days_out 0, got %d", response.DaysOut)
		}
		if response.CurrentPhase != "off_season" {
			t.Errorf("Expected current_phase 'off_season', got %s", response.CurrentPhase)
		}
	})
}

func TestMeetDateValidation(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a cycle and program for enrollment
	cycleID := createMeetDateTestCycle(t, ts, "Meet Date Validation Cycle")
	programID := createMeetDateTestProgram(t, ts, "Meet Date Validation Program", "meet-date-validation-program", cycleID)

	userID := "meet-date-validation-user"
	enrollUserForMeetDateTest(t, ts, userID, programID)

	t.Run("rejects past meet date", func(t *testing.T) {
		pastDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02") // yesterday
		body := `{"meet_date": "` + pastDate + `"}`
		resp, err := userPutMeetDate(ts.URL("/users/"+userID+"/programs/"+programID+"/state/meet-date"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 400, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("rejects invalid date format", func(t *testing.T) {
		body := `{"meet_date": "not-a-date"}`
		resp, err := userPutMeetDate(ts.URL("/users/"+userID+"/programs/"+programID+"/state/meet-date"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 400, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("accepts RFC3339 format", func(t *testing.T) {
		futureDate := time.Now().AddDate(0, 3, 0).Format(time.RFC3339) // 3 months from now
		body := `{"meet_date": "` + futureDate + `"}`
		resp, err := userPutMeetDate(ts.URL("/users/"+userID+"/programs/"+programID+"/state/meet-date"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("returns 404 for non-enrolled user", func(t *testing.T) {
		body := `{"meet_date": "2025-06-15"}`
		resp, err := userPutMeetDate(ts.URL("/users/non-enrolled-user/programs/some-program/state/meet-date"), body, "non-enrolled-user")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestMeetDateAuthorization(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a cycle and program for enrollment
	cycleID := createMeetDateTestCycle(t, ts, "Meet Date Auth Cycle")
	programID := createMeetDateTestProgram(t, ts, "Meet Date Auth Program", "meet-date-auth-program", cycleID)

	userID := "meet-date-auth-user"
	otherUserID := "other-user"
	enrollUserForMeetDateTest(t, ts, userID, programID)

	futureDate := time.Now().AddDate(0, 3, 0).Format("2006-01-02")

	t.Run("unauthenticated user gets 401 on PUT", func(t *testing.T) {
		body := `{"meet_date": "` + futureDate + `"}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL("/users/"+userID+"/programs/"+programID+"/state/meet-date"), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		// No X-User-ID header

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("unauthenticated user gets 401 on GET", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/users/"+userID+"/programs/"+programID+"/state/countdown"), nil)
		// No X-User-ID header

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("user cannot set another user's meet date", func(t *testing.T) {
		body := `{"meet_date": "` + futureDate + `"}`
		resp, err := userPutMeetDate(ts.URL("/users/"+userID+"/programs/"+programID+"/state/meet-date"), body, otherUserID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("user cannot view another user's countdown", func(t *testing.T) {
		resp, err := userGetCountdown(ts.URL("/users/"+userID+"/programs/"+programID+"/state/countdown"), otherUserID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("admin can set any user's meet date", func(t *testing.T) {
		body := `{"meet_date": "` + futureDate + `"}`
		resp, err := adminPutMeetDate(ts.URL("/users/"+userID+"/programs/"+programID+"/state/meet-date"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("admin can view any user's countdown", func(t *testing.T) {
		resp, err := adminGetCountdown(ts.URL("/users/" + userID + "/programs/" + programID + "/state/countdown"))
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

func TestMeetDatePhaseCalculation(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a cycle and program for enrollment
	cycleID := createMeetDateTestCycle(t, ts, "Phase Calc Cycle")
	programID := createMeetDateTestProgram(t, ts, "Phase Calc Program", "phase-calc-program", cycleID)

	userID := "phase-calc-user"
	enrollUserForMeetDateTest(t, ts, userID, programID)

	testCases := []struct {
		name          string
		daysFromNow   int
		expectedPhase string
	}{
		{"base phase (>84 days)", 100, "base"},
		{"prep_1 phase (57-84 days)", 70, "prep_1"},
		{"prep_2 phase (29-56 days)", 45, "prep_2"},
		{"taper phase (15-28 days)", 21, "taper"},
		{"peak phase (8-14 days)", 10, "peak"},
		{"meet_week phase (1-7 days)", 5, "meet_week"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			meetDate := time.Now().AddDate(0, 0, tc.daysFromNow).Format("2006-01-02")
			body := `{"meet_date": "` + meetDate + `"}`
			resp, err := userPutMeetDate(ts.URL("/users/"+userID+"/programs/"+programID+"/state/meet-date"), body, userID)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				bodyBytes, _ := io.ReadAll(resp.Body)
				t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
			}

			var envelope MeetDateEnvelope
			json.NewDecoder(resp.Body).Decode(&envelope)
			response := envelope.Data

			if response.CurrentPhase != tc.expectedPhase {
				t.Errorf("Expected phase %s for %d days out, got %s", tc.expectedPhase, tc.daysFromNow, response.CurrentPhase)
			}
		})
	}
}

func TestMeetDateResponseFormat(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a cycle and program for enrollment
	cycleID := createMeetDateTestCycle(t, ts, "Response Format Cycle")
	programID := createMeetDateTestProgram(t, ts, "Response Format Program", "response-format-program", cycleID)

	userID := "response-format-user"
	enrollUserForMeetDateTest(t, ts, userID, programID)

	futureDate := time.Now().AddDate(0, 3, 0).Format("2006-01-02")

	// Set meet date first
	body := `{"meet_date": "` + futureDate + `"}`
	resp, _ := userPutMeetDate(ts.URL("/users/"+userID+"/programs/"+programID+"/state/meet-date"), body, userID)
	resp.Body.Close()

	t.Run("PUT response has correct JSON field names", func(t *testing.T) {
		body := `{"meet_date": "` + futureDate + `"}`
		resp, _ := userPutMeetDate(ts.URL("/users/"+userID+"/programs/"+programID+"/state/meet-date"), body, userID)
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		bodyStr := string(respBody)

		expectedFields := []string{
			`"meet_date"`,
			`"days_out"`,
			`"current_phase"`,
			`"weeks_to_meet"`,
		}

		for _, field := range expectedFields {
			if !bytes.Contains(respBody, []byte(field)) {
				t.Errorf("Expected field %s in response, body: %s", field, bodyStr)
			}
		}
	})

	t.Run("GET countdown response has correct JSON field names", func(t *testing.T) {
		resp, _ := userGetCountdown(ts.URL("/users/"+userID+"/programs/"+programID+"/state/countdown"), userID)
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		bodyStr := string(respBody)

		expectedFields := []string{
			`"meet_date"`,
			`"days_out"`,
			`"current_phase"`,
			`"phase_week"`,
			`"taper_multiplier"`,
		}

		for _, field := range expectedFields {
			if !bytes.Contains(respBody, []byte(field)) {
				t.Errorf("Expected field %s in response, body: %s", field, bodyStr)
			}
		}
	})
}
