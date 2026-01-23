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

// StateAdvancementTestResponse represents the API response for state advancement.
type StateAdvancementTestResponse struct {
	CurrentWeek           int       `json:"currentWeek"`
	CurrentCycleIteration int       `json:"currentCycleIteration"`
	CurrentDayIndex       *int      `json:"currentDayIndex,omitempty"`
	CycleCompleted        bool      `json:"cycleCompleted"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

// StateAdvancementEnvelope wraps state advancement response with standard envelope.
type StateAdvancementEnvelope struct {
	Data StateAdvancementTestResponse `json:"data"`
}

// CycleEnvelopeForAdvancement wraps cycle response with standard envelope.
type CycleEnvelopeForAdvancement struct {
	Data CycleTestResponse `json:"data"`
}

// ProgramEnvelopeForAdvancement wraps program response with standard envelope.
type ProgramEnvelopeForAdvancement struct {
	Data ProgramTestResponse `json:"data"`
}

// WeekEnvelopeForAdvancement wraps week response with standard envelope.
type WeekEnvelopeForAdvancement struct {
	Data WeekTestResponse `json:"data"`
}

// DayEnvelopeForAdvancement wraps day response with standard envelope.
type DayEnvelopeForAdvancement struct {
	Data DayTestResponse `json:"data"`
}

// Helper functions for state advancement tests

func userPostAdvance(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func adminPostAdvance(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// createAdvancementTestCycle creates a test cycle and returns its ID
func createAdvancementTestCycle(t *testing.T, ts *testutil.TestServer, name string, lengthWeeks int) string {
	body := `{"name": "` + name + `", "lengthWeeks": ` + intToStr(lengthWeeks) + `}`
	resp, err := adminPostCycle(ts.URL("/cycles"), body)
	if err != nil {
		t.Fatalf("Failed to create test cycle: %v", err)
	}
	defer resp.Body.Close()

	var envelope CycleEnvelopeForAdvancement
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// createAdvancementTestProgram creates a test program and returns its ID
func createAdvancementTestProgram(t *testing.T, ts *testutil.TestServer, name, slug, cycleID string) string {
	body := `{"name": "` + name + `", "slug": "` + slug + `", "cycleId": "` + cycleID + `"}`
	resp, err := adminPostProgram(ts.URL("/programs"), body)
	if err != nil {
		t.Fatalf("Failed to create test program: %v", err)
	}
	defer resp.Body.Close()

	var envelope ProgramEnvelopeForAdvancement
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// createAdvancementTestWeek creates a week for the cycle and returns its ID
func createAdvancementTestWeek(t *testing.T, ts *testutil.TestServer, cycleID string, weekNumber int) string {
	body := `{"weekNumber": ` + intToStr(weekNumber) + `, "cycleId": "` + cycleID + `"}`
	resp, err := adminPostWeek(ts.URL("/weeks"), body)
	if err != nil {
		t.Fatalf("Failed to create test week: %v", err)
	}
	defer resp.Body.Close()

	var envelope WeekEnvelopeForAdvancement
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// createAdvancementTestDay creates a day and returns its ID
func createAdvancementTestDay(t *testing.T, ts *testutil.TestServer, name, slug string) string {
	body := `{"name": "` + name + `", "slug": "` + slug + `"}`
	resp, err := adminPostDay(ts.URL("/days"), body)
	if err != nil {
		t.Fatalf("Failed to create test day: %v", err)
	}
	defer resp.Body.Close()

	var envelope DayEnvelopeForAdvancement
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// addDayToWeek adds a day to a week
func addDayToWeek(t *testing.T, ts *testutil.TestServer, weekID, dayID, dayOfWeek string) {
	body := `{"dayId": "` + dayID + `", "dayOfWeek": "` + dayOfWeek + `"}`
	req, _ := http.NewRequest(http.MethodPost, ts.URL("/weeks/"+weekID+"/days"), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to add day to week: %v", err)
	}
	resp.Body.Close()
}

// enrollUserInProgram enrolls a user in a program
func enrollUserInProgram(t *testing.T, ts *testutil.TestServer, userID, programID string) {
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

// intToStr converts int to string
func intToStr(i int) string {
	return fmt.Sprintf("%d", i)
}

// Note: adminPostWeek is defined in week_handler_test.go

func adminPostDay(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

func TestStateAdvancementBasic(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create cycle with 2 weeks
	cycleID := createAdvancementTestCycle(t, ts, "Adv Test Cycle", 2)
	programID := createAdvancementTestProgram(t, ts, "Adv Test Program", "adv-test-program", cycleID)

	// Create weeks with days
	week1ID := createAdvancementTestWeek(t, ts, cycleID, 1)
	week2ID := createAdvancementTestWeek(t, ts, cycleID, 2)

	// Create days
	dayAID := createAdvancementTestDay(t, ts, "Day A", "adv-day-a")
	dayBID := createAdvancementTestDay(t, ts, "Day B", "adv-day-b")

	// Add 2 days to each week
	addDayToWeek(t, ts, week1ID, dayAID, "MONDAY")
	addDayToWeek(t, ts, week1ID, dayBID, "WEDNESDAY")
	addDayToWeek(t, ts, week2ID, dayAID, "MONDAY")
	addDayToWeek(t, ts, week2ID, dayBID, "WEDNESDAY")

	// Use standard test user that is seeded in migrations
	userID := testutil.TestUserID
	enrollUserInProgram(t, ts, userID, programID)

	t.Run("advances state from initial position", func(t *testing.T) {
		resp, err := userPostAdvance(ts.URL("/users/"+userID+"/program-state/advance"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope StateAdvancementEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		advResp := envelope.Data

		// First advancement: nil -> day 1, week 1
		if advResp.CurrentDayIndex == nil || *advResp.CurrentDayIndex != 1 {
			t.Errorf("Expected currentDayIndex 1, got %v", advResp.CurrentDayIndex)
		}
		if advResp.CurrentWeek != 1 {
			t.Errorf("Expected currentWeek 1, got %d", advResp.CurrentWeek)
		}
		if advResp.CurrentCycleIteration != 1 {
			t.Errorf("Expected currentCycleIteration 1, got %d", advResp.CurrentCycleIteration)
		}
		if advResp.CycleCompleted {
			t.Error("Expected cycleCompleted false")
		}
	})

	t.Run("advances to next week", func(t *testing.T) {
		// Advance once more to complete week 1, start week 2
		resp, _ := userPostAdvance(ts.URL("/users/"+userID+"/program-state/advance"), userID)
		defer resp.Body.Close()

		var envelope StateAdvancementEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		advResp := envelope.Data

		// Second advancement: day 1 -> day 0, week 2
		if advResp.CurrentDayIndex == nil || *advResp.CurrentDayIndex != 0 {
			t.Errorf("Expected currentDayIndex 0 (new week), got %v", advResp.CurrentDayIndex)
		}
		if advResp.CurrentWeek != 2 {
			t.Errorf("Expected currentWeek 2, got %d", advResp.CurrentWeek)
		}
		if advResp.CycleCompleted {
			t.Error("Expected cycleCompleted false")
		}
	})

	t.Run("advances through week 2", func(t *testing.T) {
		// Advance within week 2
		resp, _ := userPostAdvance(ts.URL("/users/"+userID+"/program-state/advance"), userID)
		defer resp.Body.Close()

		var envelope StateAdvancementEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		advResp := envelope.Data

		// Third advancement: day 0, week 2 -> day 1, week 2
		if advResp.CurrentDayIndex == nil || *advResp.CurrentDayIndex != 1 {
			t.Errorf("Expected currentDayIndex 1, got %v", advResp.CurrentDayIndex)
		}
		if advResp.CurrentWeek != 2 {
			t.Errorf("Expected currentWeek 2, got %d", advResp.CurrentWeek)
		}
		if advResp.CycleCompleted {
			t.Error("Expected cycleCompleted false")
		}
	})

	t.Run("completes cycle", func(t *testing.T) {
		// Advance to complete cycle
		resp, _ := userPostAdvance(ts.URL("/users/"+userID+"/program-state/advance"), userID)
		defer resp.Body.Close()

		var envelope StateAdvancementEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		advResp := envelope.Data

		// Fourth advancement: day 1, week 2 -> day 0, week 1, cycle 2 (CYCLE COMPLETE!)
		if advResp.CurrentDayIndex == nil || *advResp.CurrentDayIndex != 0 {
			t.Errorf("Expected currentDayIndex 0 (new cycle), got %v", advResp.CurrentDayIndex)
		}
		if advResp.CurrentWeek != 1 {
			t.Errorf("Expected currentWeek 1 (new cycle), got %d", advResp.CurrentWeek)
		}
		if advResp.CurrentCycleIteration != 2 {
			t.Errorf("Expected currentCycleIteration 2, got %d", advResp.CurrentCycleIteration)
		}
		if !advResp.CycleCompleted {
			t.Error("Expected cycleCompleted true")
		}
	})
}

func TestStateAdvancementNotEnrolled(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Use a seeded user who is not enrolled
	userID := "non-enrolled-user"

	t.Run("returns 404 when not enrolled", func(t *testing.T) {
		resp, err := userPostAdvance(ts.URL("/users/"+userID+"/program-state/advance"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 404, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})
}

func TestStateAdvancementAuthorization(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create cycle and program
	cycleID := createAdvancementTestCycle(t, ts, "Auth Test Cycle", 4)
	programID := createAdvancementTestProgram(t, ts, "Auth Test Program", "auth-test-program-adv", cycleID)

	// Create week and day
	weekID := createAdvancementTestWeek(t, ts, cycleID, 1)
	dayID := createAdvancementTestDay(t, ts, "Auth Day", "auth-day-adv")
	addDayToWeek(t, ts, weekID, dayID, "MONDAY")

	userID := "auth-test-user-adv"
	enrollUserInProgram(t, ts, userID, programID)

	t.Run("unauthenticated user gets 401", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/users/"+userID+"/program-state/advance"), nil)
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

	t.Run("user cannot advance another user's state", func(t *testing.T) {
		resp, err := userPostAdvance(ts.URL("/users/"+userID+"/program-state/advance"), "other-user")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("admin can advance any user's state", func(t *testing.T) {
		resp, err := adminPostAdvance(ts.URL("/users/" + userID + "/program-state/advance"))
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

func TestStateAdvancementResponseFormat(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create cycle and program
	cycleID := createAdvancementTestCycle(t, ts, "Format Test Cycle", 4)
	programID := createAdvancementTestProgram(t, ts, "Format Test Program", "format-test-program-adv", cycleID)

	// Create week and day
	weekID := createAdvancementTestWeek(t, ts, cycleID, 1)
	dayID := createAdvancementTestDay(t, ts, "Format Day", "format-day-adv")
	addDayToWeek(t, ts, weekID, dayID, "MONDAY")

	userID := "format-test-user-adv"
	enrollUserInProgram(t, ts, userID, programID)

	t.Run("response has correct JSON field names", func(t *testing.T) {
		resp, _ := userPostAdvance(ts.URL("/users/"+userID+"/program-state/advance"), userID)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		expectedFields := []string{
			`"currentWeek"`,
			`"currentCycleIteration"`,
			`"currentDayIndex"`,
			`"cycleCompleted"`,
			`"updatedAt"`,
		}

		for _, field := range expectedFields {
			if !bytes.Contains(body, []byte(field)) {
				t.Errorf("Expected field %s in response, body: %s", field, bodyStr)
			}
		}
	})
}

func TestStateAdvancementNoDaysConfigured(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create cycle and program with NO days configured (edge case)
	cycleID := createAdvancementTestCycle(t, ts, "No Days Cycle", 2)
	programID := createAdvancementTestProgram(t, ts, "No Days Program", "no-days-program", cycleID)

	// Create weeks but don't add any days
	createAdvancementTestWeek(t, ts, cycleID, 1)
	createAdvancementTestWeek(t, ts, cycleID, 2)

	userID := "no-days-user"
	enrollUserInProgram(t, ts, userID, programID)

	t.Run("handles week with no days configured", func(t *testing.T) {
		// Should still work - defaults to 1 day per week
		resp, err := userPostAdvance(ts.URL("/users/"+userID+"/program-state/advance"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope StateAdvancementEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		advResp := envelope.Data

		// With 0 days, we default to 1, so first advancement completes day 0 -> moves to next week
		if advResp.CurrentWeek != 2 {
			t.Errorf("Expected currentWeek 2 (0 days = 1 day default), got %d", advResp.CurrentWeek)
		}
	})
}
