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

// EnrollmentProgramTestResponse represents program info in an enrollment response.
type EnrollmentProgramTestResponse struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Slug             string  `json:"slug"`
	Description      *string `json:"description,omitempty"`
	CycleLengthWeeks int     `json:"cycleLengthWeeks"`
}

// EnrollmentStateTestResponse represents the state portion of an enrollment response.
type EnrollmentStateTestResponse struct {
	CurrentWeek           int  `json:"currentWeek"`
	CurrentCycleIteration int  `json:"currentCycleIteration"`
	CurrentDayIndex       *int `json:"currentDayIndex,omitempty"`
}

// CurrentWorkoutSessionTestResponse represents the current workout session in an enrollment response.
type CurrentWorkoutSessionTestResponse struct {
	ID         string     `json:"id"`
	WeekNumber int        `json:"weekNumber"`
	DayIndex   int        `json:"dayIndex"`
	Status     string     `json:"status"`
	StartedAt  time.Time  `json:"startedAt"`
	FinishedAt *time.Time `json:"finishedAt,omitempty"`
}

// EnrollmentTestResponse represents the API response format for a user's program enrollment.
type EnrollmentTestResponse struct {
	ID                    string                             `json:"id"`
	UserID                string                             `json:"userId"`
	Program               EnrollmentProgramTestResponse      `json:"program"`
	State                 EnrollmentStateTestResponse        `json:"state"`
	EnrollmentStatus      string                             `json:"enrollmentStatus"`
	CycleStatus           string                             `json:"cycleStatus"`
	WeekStatus            string                             `json:"weekStatus"`
	CurrentWorkoutSession *CurrentWorkoutSessionTestResponse `json:"currentWorkoutSession"`
	EnrolledAt            time.Time                          `json:"enrolledAt"`
	UpdatedAt             time.Time                          `json:"updatedAt"`
}

// EnrollmentEnvelope wraps single enrollment response with standard envelope.
type EnrollmentEnvelope struct {
	Data EnrollmentTestResponse `json:"data"`
}

// CycleResponseForEnrollment represents cycle data in the API response.
type CycleResponseForEnrollment struct {
	ID string `json:"id"`
}

// CycleEnvelopeForEnrollment wraps cycle response with standard envelope.
type CycleEnvelopeForEnrollment struct {
	Data CycleResponseForEnrollment `json:"data"`
}

// ProgramResponseForEnrollment represents program data in the API response.
type ProgramResponseForEnrollment struct {
	ID string `json:"id"`
}

// ProgramEnvelopeForEnrollment wraps program response with standard envelope.
type ProgramEnvelopeForEnrollment struct {
	Data ProgramResponseForEnrollment `json:"data"`
}

// Helper functions for enrollment tests

func userPostEnrollment(url string, body string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func userGetEnrollment(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func userDeleteEnrollment(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func adminPostEnrollment(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

func adminGetEnrollment(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

func adminDeleteEnrollment(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

// createEnrollmentTestCycle creates a test cycle and returns its ID
func createEnrollmentTestCycle(t *testing.T, ts *testutil.TestServer, name string) string {
	body := `{"name": "` + name + `", "lengthWeeks": 4}`
	resp, err := adminPostCycle(ts.URL("/cycles"), body)
	if err != nil {
		t.Fatalf("Failed to create test cycle: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		t.Fatalf("Expected status 201 for cycle creation, got %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope CycleEnvelopeForEnrollment
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		resp.Body.Close()
		t.Fatalf("Failed to decode cycle response: %v", err)
	}
	resp.Body.Close()
	if envelope.Data.ID == "" {
		t.Fatal("Cycle ID is empty")
	}
	return envelope.Data.ID
}

// createEnrollmentTestProgram creates a test program and returns its ID
func createEnrollmentTestProgram(t *testing.T, ts *testutil.TestServer, name, slug, cycleID string) string {
	body := `{"name": "` + name + `", "slug": "` + slug + `", "cycleId": "` + cycleID + `"}`
	resp, err := adminPostProgram(ts.URL("/programs"), body)
	if err != nil {
		t.Fatalf("Failed to create test program: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 201 for program creation, got %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope ProgramEnvelopeForEnrollment
	json.NewDecoder(resp.Body).Decode(&envelope)
	if envelope.Data.ID == "" {
		t.Fatal("Program ID is empty")
	}
	return envelope.Data.ID
}

func TestEnrollmentCRUD(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a cycle and program for enrollment
	cycleID := createEnrollmentTestCycle(t, ts, "Enrollment Test Cycle")
	programID := createEnrollmentTestProgram(t, ts, "Test Program", "test-program-enrollment", cycleID)

	userID := testutil.TestUserID

	t.Run("enrolls user in program", func(t *testing.T) {
		body := `{"programId": "` + programID + `"}`
		resp, err := userPostEnrollment(ts.URL("/users/"+userID+"/program"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope EnrollmentEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		enrollment := envelope.Data

		if enrollment.ID == "" {
			t.Error("Expected non-empty ID")
		}
		if enrollment.UserID != userID {
			t.Errorf("Expected userId %s, got %s", userID, enrollment.UserID)
		}
		if enrollment.Program.ID != programID {
			t.Errorf("Expected programId %s, got %s", programID, enrollment.Program.ID)
		}
		if enrollment.State.CurrentWeek != 1 {
			t.Errorf("Expected currentWeek 1, got %d", enrollment.State.CurrentWeek)
		}
		if enrollment.State.CurrentCycleIteration != 1 {
			t.Errorf("Expected currentCycleIteration 1, got %d", enrollment.State.CurrentCycleIteration)
		}
	})

	t.Run("gets user's current enrollment", func(t *testing.T) {
		resp, err := userGetEnrollment(ts.URL("/users/"+userID+"/program"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope EnrollmentEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		enrollment := envelope.Data

		if enrollment.UserID != userID {
			t.Errorf("Expected userId %s, got %s", userID, enrollment.UserID)
		}
		if enrollment.Program.Name != "Test Program" {
			t.Errorf("Expected program name 'Test Program', got %s", enrollment.Program.Name)
		}
	})

	t.Run("re-enrollment replaces existing enrollment", func(t *testing.T) {
		// Create another program
		program2ID := createEnrollmentTestProgram(t, ts, "Second Program", "second-program", cycleID)

		body := `{"programId": "` + program2ID + `"}`
		resp, err := userPostEnrollment(ts.URL("/users/"+userID+"/program"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope EnrollmentEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		enrollment := envelope.Data

		if enrollment.Program.ID != program2ID {
			t.Errorf("Expected programId %s after re-enrollment, got %s", program2ID, enrollment.Program.ID)
		}
		if enrollment.Program.Name != "Second Program" {
			t.Errorf("Expected program name 'Second Program', got %s", enrollment.Program.Name)
		}
	})

	t.Run("unenrolls user from program", func(t *testing.T) {
		resp, err := userDeleteEnrollment(ts.URL("/users/"+userID+"/program"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 204, got %d: %s", resp.StatusCode, bodyBytes)
		}

		// Verify user is unenrolled
		getResp, _ := userGetEnrollment(ts.URL("/users/"+userID+"/program"), userID)
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404 after unenroll, got %d", getResp.StatusCode)
		}
	})
}

func TestEnrollmentValidation(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	userID := testutil.TestUserID

	t.Run("rejects empty programId", func(t *testing.T) {
		body := `{"programId": ""}`
		resp, err := userPostEnrollment(ts.URL("/users/"+userID+"/program"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects non-existent programId", func(t *testing.T) {
		body := `{"programId": "non-existent-program"}`
		resp, err := userPostEnrollment(ts.URL("/users/"+userID+"/program"), body, userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 404 when getting non-enrolled user", func(t *testing.T) {
		// Use a different user ID that hasn't been enrolled
		resp, err := userGetEnrollment(ts.URL("/users/non-enrolled-user/program"), "non-enrolled-user")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 404 when deleting non-enrolled user", func(t *testing.T) {
		resp, err := userDeleteEnrollment(ts.URL("/users/non-enrolled-user-delete/program"), "non-enrolled-user-delete")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestEnrollmentAuthorization(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a cycle and program for enrollment
	cycleID := createEnrollmentTestCycle(t, ts, "Auth Test Cycle")
	programID := createEnrollmentTestProgram(t, ts, "Auth Test Program", "auth-test-program", cycleID)

	userID := "auth-test-user"
	otherUserID := "other-user"

	// First enroll the user
	body := `{"programId": "` + programID + `"}`
	resp, _ := userPostEnrollment(ts.URL("/users/"+userID+"/program"), body, userID)
	resp.Body.Close()

	t.Run("unauthenticated user gets 401 on POST", func(t *testing.T) {
		body := `{"programId": "` + programID + `"}`
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/users/"+userID+"/program"), bytes.NewBufferString(body))
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
		req, _ := http.NewRequest(http.MethodGet, ts.URL("/users/"+userID+"/program"), nil)
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

	t.Run("unauthenticated user gets 401 on DELETE", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL("/users/"+userID+"/program"), nil)
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

	t.Run("user cannot enroll another user", func(t *testing.T) {
		body := `{"programId": "` + programID + `"}`
		resp, err := userPostEnrollment(ts.URL("/users/"+userID+"/program"), body, otherUserID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("user cannot view another user's enrollment", func(t *testing.T) {
		resp, err := userGetEnrollment(ts.URL("/users/"+userID+"/program"), otherUserID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("user cannot unenroll another user", func(t *testing.T) {
		resp, err := userDeleteEnrollment(ts.URL("/users/"+userID+"/program"), otherUserID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("admin can view any user's enrollment", func(t *testing.T) {
		resp, err := adminGetEnrollment(ts.URL("/users/" + userID + "/program"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("admin can enroll any user", func(t *testing.T) {
		body := `{"programId": "` + programID + `"}`
		resp, err := adminPostEnrollment(ts.URL("/users/admin-enrolled-user/program"), body)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 201, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})

	t.Run("admin can unenroll any user", func(t *testing.T) {
		resp, err := adminDeleteEnrollment(ts.URL("/users/admin-enrolled-user/program"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 204, got %d: %s", resp.StatusCode, bodyBytes)
		}
	})
}

func TestEnrollmentResponseFormat(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a cycle and program for enrollment
	cycleID := createEnrollmentTestCycle(t, ts, "Response Format Test Cycle")
	programID := createEnrollmentTestProgram(t, ts, "Format Test Program", "format-test-program", cycleID)

	userID := "format-test-user"

	// Enroll user
	body := `{"programId": "` + programID + `"}`
	resp, _ := userPostEnrollment(ts.URL("/users/"+userID+"/program"), body, userID)
	resp.Body.Close()

	t.Run("response has correct JSON field names", func(t *testing.T) {
		resp, _ := userGetEnrollment(ts.URL("/users/"+userID+"/program"), userID)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		expectedFields := []string{
			`"id"`,
			`"userId"`,
			`"program"`,
			`"state"`,
			`"enrolledAt"`,
			`"updatedAt"`,
			`"currentWeek"`,
			`"currentCycleIteration"`,
			`"cycleLengthWeeks"`,
		}

		for _, field := range expectedFields {
			if !bytes.Contains(body, []byte(field)) {
				t.Errorf("Expected field %s in response, body: %s", field, bodyStr)
			}
		}
	})

	t.Run("initial state is correct", func(t *testing.T) {
		resp, _ := userGetEnrollment(ts.URL("/users/"+userID+"/program"), userID)
		defer resp.Body.Close()

		var envelope EnrollmentEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		enrollment := envelope.Data

		if enrollment.State.CurrentWeek != 1 {
			t.Errorf("Expected initial currentWeek 1, got %d", enrollment.State.CurrentWeek)
		}
		if enrollment.State.CurrentCycleIteration != 1 {
			t.Errorf("Expected initial currentCycleIteration 1, got %d", enrollment.State.CurrentCycleIteration)
		}
		if enrollment.State.CurrentDayIndex != nil {
			t.Errorf("Expected nil currentDayIndex, got %d", *enrollment.State.CurrentDayIndex)
		}
	})

	t.Run("program info is embedded correctly", func(t *testing.T) {
		resp, _ := userGetEnrollment(ts.URL("/users/"+userID+"/program"), userID)
		defer resp.Body.Close()

		var envelope EnrollmentEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		enrollment := envelope.Data

		if enrollment.Program.ID != programID {
			t.Errorf("Expected program.id %s, got %s", programID, enrollment.Program.ID)
		}
		if enrollment.Program.Name != "Format Test Program" {
			t.Errorf("Expected program.name 'Format Test Program', got %s", enrollment.Program.Name)
		}
		if enrollment.Program.Slug != "format-test-program" {
			t.Errorf("Expected program.slug 'format-test-program', got %s", enrollment.Program.Slug)
		}
		if enrollment.Program.CycleLengthWeeks != 4 {
			t.Errorf("Expected program.cycleLengthWeeks 4, got %d", enrollment.Program.CycleLengthWeeks)
		}
	})

	t.Run("includes enrollment status fields", func(t *testing.T) {
		resp, _ := userGetEnrollment(ts.URL("/users/"+userID+"/program"), userID)
		defer resp.Body.Close()

		var envelope EnrollmentEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		enrollment := envelope.Data

		if enrollment.EnrollmentStatus != "ACTIVE" {
			t.Errorf("Expected enrollmentStatus 'ACTIVE', got %s", enrollment.EnrollmentStatus)
		}
		if enrollment.CycleStatus != "PENDING" {
			t.Errorf("Expected cycleStatus 'PENDING', got %s", enrollment.CycleStatus)
		}
		if enrollment.WeekStatus != "PENDING" {
			t.Errorf("Expected weekStatus 'PENDING', got %s", enrollment.WeekStatus)
		}
		if enrollment.CurrentWorkoutSession != nil {
			t.Error("Expected currentWorkoutSession to be nil for new enrollment")
		}
	})
}

// Helper functions for next-cycle and advance-week endpoints

func userPostNextCycle(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func userPostAdvanceWeek(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func adminPostNextCycle(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

func adminPostAdvanceWeek(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

func TestEnrollmentAdvanceWeek(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a cycle with 4 weeks and program for enrollment
	cycleID := createEnrollmentTestCycle(t, ts, "Advance Week Test Cycle")
	programID := createEnrollmentTestProgram(t, ts, "Advance Week Test Program", "advance-week-test-program", cycleID)

	userID := testutil.TestUserID

	// Enroll user
	body := `{"programId": "` + programID + `"}`
	resp, err := userPostEnrollment(ts.URL("/users/"+userID+"/program"), body, userID)
	if err != nil {
		t.Fatalf("Failed to enroll user: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 201 for enrollment, got %d: %s", resp.StatusCode, bodyBytes)
	}
	resp.Body.Close()

	t.Run("advances from week 1 to week 2", func(t *testing.T) {
		resp, err := userPostAdvanceWeek(ts.URL("/users/"+userID+"/enrollment/advance-week"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope EnrollmentEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		enrollment := envelope.Data

		if enrollment.State.CurrentWeek != 2 {
			t.Errorf("Expected currentWeek 2, got %d", enrollment.State.CurrentWeek)
		}
		if enrollment.EnrollmentStatus != "ACTIVE" {
			t.Errorf("Expected enrollmentStatus 'ACTIVE', got %s", enrollment.EnrollmentStatus)
		}
		if enrollment.WeekStatus != "PENDING" {
			t.Errorf("Expected weekStatus 'PENDING', got %s", enrollment.WeekStatus)
		}
	})

	t.Run("advances to week 3", func(t *testing.T) {
		resp, _ := userPostAdvanceWeek(ts.URL("/users/"+userID+"/enrollment/advance-week"), userID)
		defer resp.Body.Close()

		var envelope EnrollmentEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)

		if envelope.Data.State.CurrentWeek != 3 {
			t.Errorf("Expected currentWeek 3, got %d", envelope.Data.State.CurrentWeek)
		}
	})

	t.Run("advances to week 4 (final week)", func(t *testing.T) {
		resp, _ := userPostAdvanceWeek(ts.URL("/users/"+userID+"/enrollment/advance-week"), userID)
		defer resp.Body.Close()

		var envelope EnrollmentEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)

		if envelope.Data.State.CurrentWeek != 4 {
			t.Errorf("Expected currentWeek 4, got %d", envelope.Data.State.CurrentWeek)
		}
	})

	t.Run("advancing from final week transitions to BETWEEN_CYCLES", func(t *testing.T) {
		resp, err := userPostAdvanceWeek(ts.URL("/users/"+userID+"/enrollment/advance-week"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope EnrollmentEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)
		enrollment := envelope.Data

		if enrollment.EnrollmentStatus != "BETWEEN_CYCLES" {
			t.Errorf("Expected enrollmentStatus 'BETWEEN_CYCLES', got %s", enrollment.EnrollmentStatus)
		}
		if enrollment.CycleStatus != "COMPLETED" {
			t.Errorf("Expected cycleStatus 'COMPLETED', got %s", enrollment.CycleStatus)
		}
		if enrollment.WeekStatus != "COMPLETED" {
			t.Errorf("Expected weekStatus 'COMPLETED', got %s", enrollment.WeekStatus)
		}
	})

	t.Run("cannot advance week when BETWEEN_CYCLES", func(t *testing.T) {
		resp, err := userPostAdvanceWeek(ts.URL("/users/"+userID+"/enrollment/advance-week"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestEnrollmentNextCycle(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a cycle with 4 weeks and program for enrollment
	cycleID := createEnrollmentTestCycle(t, ts, "Next Cycle Test Cycle")
	programID := createEnrollmentTestProgram(t, ts, "Next Cycle Test Program", "next-cycle-test-program", cycleID)

	userID := testutil.TestUserID

	// Enroll user
	body := `{"programId": "` + programID + `"}`
	resp, err := userPostEnrollment(ts.URL("/users/"+userID+"/program"), body, userID)
	if err != nil {
		t.Fatalf("Failed to enroll user: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 201 for enrollment, got %d: %s", resp.StatusCode, bodyBytes)
	}
	resp.Body.Close()

	// Advance through all 4 weeks to reach BETWEEN_CYCLES
	for i := 0; i < 4; i++ {
		resp, err := userPostAdvanceWeek(ts.URL("/users/"+userID+"/enrollment/advance-week"), userID)
		if err != nil {
			t.Fatalf("Failed to advance week: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200 for advance-week, got %d: %s", resp.StatusCode, bodyBytes)
		}
		resp.Body.Close()
	}

	// Verify we're in BETWEEN_CYCLES
	getResp, _ := userGetEnrollment(ts.URL("/users/"+userID+"/program"), userID)
	var checkEnv EnrollmentEnvelope
	json.NewDecoder(getResp.Body).Decode(&checkEnv)
	getResp.Body.Close()
	if checkEnv.Data.EnrollmentStatus != "BETWEEN_CYCLES" {
		t.Fatalf("Expected to be in BETWEEN_CYCLES state, got %s", checkEnv.Data.EnrollmentStatus)
	}

	t.Run("starts new cycle from BETWEEN_CYCLES", func(t *testing.T) {
		resp, err := userPostNextCycle(ts.URL("/users/"+userID+"/enrollment/next-cycle"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, bodyBytes)
		}

		var envelope EnrollmentEnvelope
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		enrollment := envelope.Data

		if enrollment.EnrollmentStatus != "ACTIVE" {
			t.Errorf("Expected enrollmentStatus 'ACTIVE', got %s", enrollment.EnrollmentStatus)
		}
		if enrollment.State.CurrentCycleIteration != 2 {
			t.Errorf("Expected currentCycleIteration 2, got %d", enrollment.State.CurrentCycleIteration)
		}
		if enrollment.State.CurrentWeek != 1 {
			t.Errorf("Expected currentWeek 1, got %d", enrollment.State.CurrentWeek)
		}
		if enrollment.CycleStatus != "PENDING" {
			t.Errorf("Expected cycleStatus 'PENDING', got %s", enrollment.CycleStatus)
		}
		if enrollment.WeekStatus != "PENDING" {
			t.Errorf("Expected weekStatus 'PENDING', got %s", enrollment.WeekStatus)
		}
	})

	t.Run("cannot start next cycle when ACTIVE", func(t *testing.T) {
		resp, err := userPostNextCycle(ts.URL("/users/"+userID+"/enrollment/next-cycle"), userID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestEnrollmentStateEndpointsValidation(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	t.Run("returns 404 for non-enrolled user on next-cycle", func(t *testing.T) {
		resp, err := userPostNextCycle(ts.URL("/users/not-enrolled-user/enrollment/next-cycle"), "not-enrolled-user")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 404 for non-enrolled user on advance-week", func(t *testing.T) {
		resp, err := userPostAdvanceWeek(ts.URL("/users/not-enrolled-user/enrollment/advance-week"), "not-enrolled-user")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestEnrollmentStateEndpointsAuthorization(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Create a cycle and program for enrollment
	cycleID := createEnrollmentTestCycle(t, ts, "Auth State Test Cycle")
	programID := createEnrollmentTestProgram(t, ts, "Auth State Test Program", "auth-state-test-program", cycleID)

	userID := "auth-test-user"
	otherUserID := "other-user"

	// Enroll user
	body := `{"programId": "` + programID + `"}`
	resp, err := userPostEnrollment(ts.URL("/users/"+userID+"/program"), body, userID)
	if err != nil {
		t.Fatalf("Failed to enroll user: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 201 for enrollment, got %d: %s", resp.StatusCode, bodyBytes)
	}
	resp.Body.Close()

	t.Run("unauthenticated user gets 401 on advance-week", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/users/"+userID+"/enrollment/advance-week"), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("unauthenticated user gets 401 on next-cycle", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, ts.URL("/users/"+userID+"/enrollment/next-cycle"), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("user cannot advance another user's week", func(t *testing.T) {
		resp, err := userPostAdvanceWeek(ts.URL("/users/"+userID+"/enrollment/advance-week"), otherUserID)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})

	t.Run("admin can advance any user's week", func(t *testing.T) {
		resp, err := adminPostAdvanceWeek(ts.URL("/users/" + userID + "/enrollment/advance-week"))
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
