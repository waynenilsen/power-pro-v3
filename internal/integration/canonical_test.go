// Package integration provides integration tests for cross-component behavior.
// This file contains verification tests that validate seeded canonical programs
// have correct structure and accurate prescriptions.
package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/waynenilsen/power-pro-v3/internal/database"
	"github.com/waynenilsen/power-pro-v3/internal/db"
)

// =============================================================================
// TEST INFRASTRUCTURE
// =============================================================================

// findCanonicalMigrationsPath finds the migrations directory path relative to this file.
func findCanonicalMigrationsPath() (string, error) {
	// Get the directory of this source file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get caller info")
	}

	// Navigate up to project root
	dir := filepath.Dir(filename)
	for i := 0; i < 5; i++ {
		migrationsPath := filepath.Join(dir, "migrations")
		if _, err := os.Stat(migrationsPath); err == nil {
			return migrationsPath, nil
		}
		dir = filepath.Dir(dir)
	}

	return "", fmt.Errorf("migrations directory not found")
}

// setupCanonicalTestDB creates an in-memory database with all migrations applied.
func setupCanonicalTestDB(t *testing.T) *sql.DB {
	t.Helper()

	migrationsPath, err := findCanonicalMigrationsPath()
	if err != nil {
		t.Fatalf("failed to find migrations path: %v", err)
	}

	sqlDB, err := database.OpenInMemory(migrationsPath)
	if err != nil {
		t.Fatalf("failed to setup test database: %v", err)
	}

	return sqlDB
}

// =============================================================================
// HELPER TYPES AND FUNCTIONS
// =============================================================================

// programSpec defines expected program structure.
type programSpec struct {
	Slug        string
	Name        string
	Days        int
	CycleLengthWeeks int
}

// dayPrescriptionSpec defines expected prescription details for a day.
type dayPrescriptionSpec struct {
	LiftSlug   string
	Sets       int
	Reps       int
	Percentage float64
	IsAmrap    bool
}

// loadStrategyJSON represents the JSON structure for load strategy.
type loadStrategyJSON struct {
	Type          string  `json:"type"`
	ReferenceType string  `json:"referenceType"`
	Percentage    float64 `json:"percentage"`
}

// setSchemeJSON represents the JSON structure for set scheme.
type setSchemeJSON struct {
	Type    string `json:"type"`
	Sets    int    `json:"sets"`
	Reps    int    `json:"reps"`
	IsAmrap bool   `json:"isAmrap"`
	Tier    string `json:"tier,omitempty"`
	Stage   int    `json:"stage,omitempty"`
}

// getProgramBySlug retrieves a program by its slug.
func getProgramBySlug(t *testing.T, queries *db.Queries, slug string) db.Program {
	t.Helper()
	prog, err := queries.GetProgramBySlug(context.Background(), slug)
	if err != nil {
		t.Fatalf("failed to get program by slug %s: %v", slug, err)
	}
	return prog
}

// getDaysForProgram retrieves all days for a program.
func getDaysForProgram(t *testing.T, sqlDB *sql.DB, programID string) []db.Day {
	t.Helper()
	rows, err := sqlDB.Query(`SELECT id, name, slug, metadata, program_id, created_at, updated_at FROM days WHERE program_id = ? ORDER BY slug`, programID)
	if err != nil {
		t.Fatalf("failed to query days: %v", err)
	}
	defer rows.Close()

	var days []db.Day
	for rows.Next() {
		var d db.Day
		if err := rows.Scan(&d.ID, &d.Name, &d.Slug, &d.Metadata, &d.ProgramID, &d.CreatedAt, &d.UpdatedAt); err != nil {
			t.Fatalf("failed to scan day: %v", err)
		}
		days = append(days, d)
	}
	return days
}

// getWeeksForCycle retrieves all weeks for a cycle.
func getWeeksForCycle(t *testing.T, sqlDB *sql.DB, cycleID string) []db.Week {
	t.Helper()
	rows, err := sqlDB.Query(`SELECT id, week_number, variant, cycle_id, created_at, updated_at FROM weeks WHERE cycle_id = ? ORDER BY week_number`, cycleID)
	if err != nil {
		t.Fatalf("failed to query weeks: %v", err)
	}
	defer rows.Close()

	var weeks []db.Week
	for rows.Next() {
		var w db.Week
		if err := rows.Scan(&w.ID, &w.WeekNumber, &w.Variant, &w.CycleID, &w.CreatedAt, &w.UpdatedAt); err != nil {
			t.Fatalf("failed to scan week: %v", err)
		}
		weeks = append(weeks, w)
	}
	return weeks
}

// getPrescriptionsForDay retrieves all prescriptions for a day with lift info.
type prescriptionWithLift struct {
	PrescriptionID string
	LiftID         string
	LiftSlug       string
	LoadStrategy   string
	SetScheme      string
	Order          int64
}

func getPrescriptionsForDay(t *testing.T, sqlDB *sql.DB, dayID string) []prescriptionWithLift {
	t.Helper()
	query := `
		SELECT p.id, p.lift_id, l.slug, p.load_strategy, p.set_scheme, dp."order"
		FROM day_prescriptions dp
		JOIN prescriptions p ON dp.prescription_id = p.id
		JOIN lifts l ON p.lift_id = l.id
		WHERE dp.day_id = ?
		ORDER BY dp."order"
	`
	rows, err := sqlDB.Query(query, dayID)
	if err != nil {
		t.Fatalf("failed to query prescriptions for day: %v", err)
	}
	defer rows.Close()

	var prescriptions []prescriptionWithLift
	for rows.Next() {
		var p prescriptionWithLift
		if err := rows.Scan(&p.PrescriptionID, &p.LiftID, &p.LiftSlug, &p.LoadStrategy, &p.SetScheme, &p.Order); err != nil {
			t.Fatalf("failed to scan prescription: %v", err)
		}
		prescriptions = append(prescriptions, p)
	}
	return prescriptions
}

// getProgressionsForProgram retrieves all progression rules for a program.
type programProgressionInfo struct {
	LiftSlug  string
	ProgName  string
	ProgType  string
	Params    string
	Increment *float64
}

func getProgressionsForProgram(t *testing.T, sqlDB *sql.DB, programID string) []programProgressionInfo {
	t.Helper()
	query := `
		SELECT l.slug, pr.name, pr.type, pr.parameters, pp.override_increment
		FROM program_progressions pp
		JOIN progressions pr ON pp.progression_id = pr.id
		JOIN lifts l ON pp.lift_id = l.id
		WHERE pp.program_id = ?
		ORDER BY pp.priority
	`
	rows, err := sqlDB.Query(query, programID)
	if err != nil {
		t.Fatalf("failed to query progressions: %v", err)
	}
	defer rows.Close()

	var progressions []programProgressionInfo
	for rows.Next() {
		var p programProgressionInfo
		var overrideIncrement sql.NullFloat64
		if err := rows.Scan(&p.LiftSlug, &p.ProgName, &p.ProgType, &p.Params, &overrideIncrement); err != nil {
			t.Fatalf("failed to scan progression: %v", err)
		}
		if overrideIncrement.Valid {
			p.Increment = &overrideIncrement.Float64
		}
		progressions = append(progressions, p)
	}
	return progressions
}

// =============================================================================
// CANONICAL PROGRAM VERIFICATION TESTS
// =============================================================================

func TestCanonicalPrograms_Existence(t *testing.T) {
	sqlDB := setupCanonicalTestDB(t)
	defer sqlDB.Close()
	queries := db.New(sqlDB)

	programs := []programSpec{
		{Slug: "starting-strength", Name: "Starting Strength", Days: 2, CycleLengthWeeks: 1},
		{Slug: "texas-method", Name: "Texas Method", Days: 3, CycleLengthWeeks: 1},
		{Slug: "531", Name: "Wendler 5/3/1", Days: 16, CycleLengthWeeks: 4}, // 4 weeks x 4 days = 16 day instances
		{Slug: "gzclp", Name: "GZCLP", Days: 4, CycleLengthWeeks: 1},
	}

	for _, spec := range programs {
		t.Run(spec.Slug, func(t *testing.T) {
			prog := getProgramBySlug(t, queries, spec.Slug)

			if prog.Name != spec.Name {
				t.Errorf("program name: expected %q, got %q", spec.Name, prog.Name)
			}

			// Verify cycle length
			var lengthWeeks int64
			err := sqlDB.QueryRow(`SELECT length_weeks FROM cycles WHERE id = ?`, prog.CycleID).Scan(&lengthWeeks)
			if err != nil {
				t.Fatalf("failed to get cycle: %v", err)
			}
			if int(lengthWeeks) != spec.CycleLengthWeeks {
				t.Errorf("cycle length: expected %d weeks, got %d", spec.CycleLengthWeeks, lengthWeeks)
			}

			// Verify day count
			days := getDaysForProgram(t, sqlDB, prog.ID)
			if len(days) != spec.Days {
				t.Errorf("day count: expected %d, got %d", spec.Days, len(days))
			}
		})
	}
}

// =============================================================================
// STARTING STRENGTH VERIFICATION
// =============================================================================

func TestCanonicalPrograms_StartingStrength(t *testing.T) {
	sqlDB := setupCanonicalTestDB(t)
	defer sqlDB.Close()
	queries := db.New(sqlDB)

	prog := getProgramBySlug(t, queries, "starting-strength")
	days := getDaysForProgram(t, sqlDB, prog.ID)

	t.Run("has exactly 2 workout days", func(t *testing.T) {
		if len(days) != 2 {
			t.Errorf("expected 2 days, got %d", len(days))
		}
	})

	// Find Day A and Day B
	var dayA, dayB db.Day
	for _, d := range days {
		if d.Slug == "workout-a" {
			dayA = d
		} else if d.Slug == "workout-b" {
			dayB = d
		}
	}

	t.Run("Day A prescriptions", func(t *testing.T) {
		if dayA.ID == "" {
			t.Fatal("Day A (workout-a) not found")
		}

		prescriptions := getPrescriptionsForDay(t, sqlDB, dayA.ID)
		if len(prescriptions) != 3 {
			t.Fatalf("expected 3 prescriptions for Day A, got %d", len(prescriptions))
		}

		// Expected: Squat 3x5, Bench 3x5, Deadlift 1x5
		expected := []struct {
			liftSlug   string
			sets, reps int
		}{
			{"squat", 3, 5},
			{"bench-press", 3, 5},
			{"deadlift", 1, 5},
		}

		for i, exp := range expected {
			p := prescriptions[i]
			if p.LiftSlug != exp.liftSlug {
				t.Errorf("prescription %d: expected lift %s, got %s", i, exp.liftSlug, p.LiftSlug)
			}

			var scheme setSchemeJSON
			if err := json.Unmarshal([]byte(p.SetScheme), &scheme); err != nil {
				t.Fatalf("failed to parse set scheme: %v", err)
			}
			if scheme.Sets != exp.sets {
				t.Errorf("prescription %d (%s): expected %d sets, got %d", i, exp.liftSlug, exp.sets, scheme.Sets)
			}
			if scheme.Reps != exp.reps {
				t.Errorf("prescription %d (%s): expected %d reps, got %d", i, exp.liftSlug, exp.reps, scheme.Reps)
			}
		}
	})

	t.Run("Day B prescriptions", func(t *testing.T) {
		if dayB.ID == "" {
			t.Fatal("Day B (workout-b) not found")
		}

		prescriptions := getPrescriptionsForDay(t, sqlDB, dayB.ID)
		if len(prescriptions) != 3 {
			t.Fatalf("expected 3 prescriptions for Day B, got %d", len(prescriptions))
		}

		// Expected: Squat 3x5, Press 3x5, Power Clean 5x3
		expected := []struct {
			liftSlug   string
			sets, reps int
		}{
			{"squat", 3, 5},
			{"overhead-press", 3, 5},
			{"power-clean", 5, 3},
		}

		for i, exp := range expected {
			p := prescriptions[i]
			if p.LiftSlug != exp.liftSlug {
				t.Errorf("prescription %d: expected lift %s, got %s", i, exp.liftSlug, p.LiftSlug)
			}

			var scheme setSchemeJSON
			if err := json.Unmarshal([]byte(p.SetScheme), &scheme); err != nil {
				t.Fatalf("failed to parse set scheme: %v", err)
			}
			if scheme.Sets != exp.sets {
				t.Errorf("prescription %d (%s): expected %d sets, got %d", i, exp.liftSlug, exp.sets, scheme.Sets)
			}
			if scheme.Reps != exp.reps {
				t.Errorf("prescription %d (%s): expected %d reps, got %d", i, exp.liftSlug, exp.reps, scheme.Reps)
			}
		}
	})

	t.Run("progression rules", func(t *testing.T) {
		progressions := getProgressionsForProgram(t, sqlDB, prog.ID)

		// Expected: Squat +10, Bench +5, Press +5, Deadlift +10, Power Clean +5
		expectedIncrements := map[string]float64{
			"squat":          10.0,
			"bench-press":    5.0,
			"overhead-press": 5.0,
			"deadlift":       10.0,
			"power-clean":    5.0,
		}

		for _, p := range progressions {
			expected, ok := expectedIncrements[p.LiftSlug]
			if !ok {
				t.Errorf("unexpected lift %s in progressions", p.LiftSlug)
				continue
			}

			// Parse the increment from parameters
			var params struct {
				Increment float64 `json:"increment"`
			}
			if err := json.Unmarshal([]byte(p.Params), &params); err != nil {
				t.Fatalf("failed to parse progression params for %s: %v", p.LiftSlug, err)
			}

			if params.Increment != expected {
				t.Errorf("%s progression: expected +%.1f, got +%.1f", p.LiftSlug, expected, params.Increment)
			}
		}
	})
}

// =============================================================================
// TEXAS METHOD VERIFICATION
// =============================================================================

func TestCanonicalPrograms_TexasMethod(t *testing.T) {
	sqlDB := setupCanonicalTestDB(t)
	defer sqlDB.Close()
	queries := db.New(sqlDB)

	prog := getProgramBySlug(t, queries, "texas-method")
	days := getDaysForProgram(t, sqlDB, prog.ID)

	t.Run("has exactly 3 workout days", func(t *testing.T) {
		if len(days) != 3 {
			t.Errorf("expected 3 days, got %d", len(days))
		}
	})

	// Find Volume, Recovery, Intensity days
	daysBySlug := make(map[string]db.Day)
	for _, d := range days {
		daysBySlug[d.Slug] = d
	}

	t.Run("Volume Day prescriptions at 90%", func(t *testing.T) {
		day, ok := daysBySlug["volume-day"]
		if !ok {
			t.Fatal("Volume Day not found")
		}

		prescriptions := getPrescriptionsForDay(t, sqlDB, day.ID)

		// Check that volume day prescriptions are at 90%
		for _, p := range prescriptions {
			var strategy loadStrategyJSON
			if err := json.Unmarshal([]byte(p.LoadStrategy), &strategy); err != nil {
				t.Fatalf("failed to parse load strategy: %v", err)
			}

			if strategy.Percentage != 90.0 {
				t.Errorf("Volume Day %s: expected 90%%, got %.1f%%", p.LiftSlug, strategy.Percentage)
			}
		}
	})

	t.Run("Recovery Day squat at 72%", func(t *testing.T) {
		day, ok := daysBySlug["recovery-day"]
		if !ok {
			t.Fatal("Recovery Day not found")
		}

		prescriptions := getPrescriptionsForDay(t, sqlDB, day.ID)

		// Find squat prescription
		for _, p := range prescriptions {
			if p.LiftSlug == "squat" {
				var strategy loadStrategyJSON
				if err := json.Unmarshal([]byte(p.LoadStrategy), &strategy); err != nil {
					t.Fatalf("failed to parse load strategy: %v", err)
				}

				if strategy.Percentage != 72.0 {
					t.Errorf("Recovery Day squat: expected 72%%, got %.1f%%", strategy.Percentage)
				}
				break
			}
		}
	})

	t.Run("Intensity Day prescriptions at 100%", func(t *testing.T) {
		day, ok := daysBySlug["intensity-day"]
		if !ok {
			t.Fatal("Intensity Day not found")
		}

		prescriptions := getPrescriptionsForDay(t, sqlDB, day.ID)

		// Check that intensity day prescriptions are at 100%
		for _, p := range prescriptions {
			var strategy loadStrategyJSON
			if err := json.Unmarshal([]byte(p.LoadStrategy), &strategy); err != nil {
				t.Fatalf("failed to parse load strategy: %v", err)
			}

			if strategy.Percentage != 100.0 {
				t.Errorf("Intensity Day %s: expected 100%%, got %.1f%%", p.LiftSlug, strategy.Percentage)
			}
		}
	})

	t.Run("weekly progression model", func(t *testing.T) {
		progressions := getProgressionsForProgram(t, sqlDB, prog.ID)

		for _, p := range progressions {
			if p.ProgType != "LINEAR_PROGRESSION" {
				t.Errorf("%s: expected LINEAR_PROGRESSION type, got %s", p.LiftSlug, p.ProgType)
			}

			// Check trigger type is AFTER_WEEK
			var params struct {
				TriggerType string `json:"triggerType"`
			}
			if err := json.Unmarshal([]byte(p.Params), &params); err != nil {
				t.Fatalf("failed to parse progression params: %v", err)
			}

			if params.TriggerType != "AFTER_WEEK" {
				t.Errorf("%s: expected AFTER_WEEK trigger, got %s", p.LiftSlug, params.TriggerType)
			}
		}
	})
}

// =============================================================================
// WENDLER 5/3/1 VERIFICATION
// =============================================================================

func TestCanonicalPrograms_531(t *testing.T) {
	sqlDB := setupCanonicalTestDB(t)
	defer sqlDB.Close()
	queries := db.New(sqlDB)

	prog := getProgramBySlug(t, queries, "531")

	t.Run("has 4 weeks per cycle", func(t *testing.T) {
		weeks := getWeeksForCycle(t, sqlDB, prog.CycleID)
		if len(weeks) != 4 {
			t.Errorf("expected 4 weeks, got %d", len(weeks))
		}
	})

	t.Run("has 4 workout days per week (16 total day instances)", func(t *testing.T) {
		days := getDaysForProgram(t, sqlDB, prog.ID)
		if len(days) != 16 {
			t.Errorf("expected 16 day instances (4 weeks x 4 days), got %d", len(days))
		}
	})

	t.Run("Week 1 percentages (65%, 75%, 85%)", func(t *testing.T) {
		// Find a Week 1 day
		var week1DayID string
		err := sqlDB.QueryRow(`
			SELECT d.id FROM days d
			WHERE d.program_id = ? AND d.slug LIKE '%-w1'
			LIMIT 1
		`, prog.ID).Scan(&week1DayID)
		if err != nil {
			t.Fatalf("failed to find Week 1 day: %v", err)
		}

		prescriptions := getPrescriptionsForDay(t, sqlDB, week1DayID)
		expectedPcts := []float64{65.0, 75.0, 85.0}

		if len(prescriptions) != 3 {
			t.Fatalf("expected 3 prescriptions for Week 1 day, got %d", len(prescriptions))
		}

		for i, p := range prescriptions {
			var strategy loadStrategyJSON
			if err := json.Unmarshal([]byte(p.LoadStrategy), &strategy); err != nil {
				t.Fatalf("failed to parse load strategy: %v", err)
			}

			if strategy.Percentage != expectedPcts[i] {
				t.Errorf("Week 1 set %d: expected %.1f%%, got %.1f%%", i+1, expectedPcts[i], strategy.Percentage)
			}
		}
	})

	t.Run("Week 2 percentages (70%, 80%, 90%)", func(t *testing.T) {
		var week2DayID string
		err := sqlDB.QueryRow(`
			SELECT d.id FROM days d
			WHERE d.program_id = ? AND d.slug LIKE '%-w2'
			LIMIT 1
		`, prog.ID).Scan(&week2DayID)
		if err != nil {
			t.Fatalf("failed to find Week 2 day: %v", err)
		}

		prescriptions := getPrescriptionsForDay(t, sqlDB, week2DayID)
		expectedPcts := []float64{70.0, 80.0, 90.0}

		for i, p := range prescriptions {
			var strategy loadStrategyJSON
			if err := json.Unmarshal([]byte(p.LoadStrategy), &strategy); err != nil {
				t.Fatalf("failed to parse load strategy: %v", err)
			}

			if strategy.Percentage != expectedPcts[i] {
				t.Errorf("Week 2 set %d: expected %.1f%%, got %.1f%%", i+1, expectedPcts[i], strategy.Percentage)
			}
		}
	})

	t.Run("Week 3 percentages (75%, 85%, 95%)", func(t *testing.T) {
		var week3DayID string
		err := sqlDB.QueryRow(`
			SELECT d.id FROM days d
			WHERE d.program_id = ? AND d.slug LIKE '%-w3'
			LIMIT 1
		`, prog.ID).Scan(&week3DayID)
		if err != nil {
			t.Fatalf("failed to find Week 3 day: %v", err)
		}

		prescriptions := getPrescriptionsForDay(t, sqlDB, week3DayID)
		expectedPcts := []float64{75.0, 85.0, 95.0}

		for i, p := range prescriptions {
			var strategy loadStrategyJSON
			if err := json.Unmarshal([]byte(p.LoadStrategy), &strategy); err != nil {
				t.Fatalf("failed to parse load strategy: %v", err)
			}

			if strategy.Percentage != expectedPcts[i] {
				t.Errorf("Week 3 set %d: expected %.1f%%, got %.1f%%", i+1, expectedPcts[i], strategy.Percentage)
			}
		}
	})

	t.Run("Week 4 (Deload) percentages (40%, 50%, 60%)", func(t *testing.T) {
		var week4DayID string
		err := sqlDB.QueryRow(`
			SELECT d.id FROM days d
			WHERE d.program_id = ? AND d.slug LIKE '%-w4'
			LIMIT 1
		`, prog.ID).Scan(&week4DayID)
		if err != nil {
			t.Fatalf("failed to find Week 4 day: %v", err)
		}

		prescriptions := getPrescriptionsForDay(t, sqlDB, week4DayID)
		expectedPcts := []float64{40.0, 50.0, 60.0}

		for i, p := range prescriptions {
			var strategy loadStrategyJSON
			if err := json.Unmarshal([]byte(p.LoadStrategy), &strategy); err != nil {
				t.Fatalf("failed to parse load strategy: %v", err)
			}

			if strategy.Percentage != expectedPcts[i] {
				t.Errorf("Week 4 set %d: expected %.1f%%, got %.1f%%", i+1, expectedPcts[i], strategy.Percentage)
			}
		}
	})

	t.Run("AMRAP sets marked correctly", func(t *testing.T) {
		// Find a Week 1 day and check the last set is AMRAP
		var week1DayID string
		err := sqlDB.QueryRow(`
			SELECT d.id FROM days d
			WHERE d.program_id = ? AND d.slug LIKE '%-w1'
			LIMIT 1
		`, prog.ID).Scan(&week1DayID)
		if err != nil {
			t.Fatalf("failed to find Week 1 day: %v", err)
		}

		prescriptions := getPrescriptionsForDay(t, sqlDB, week1DayID)
		if len(prescriptions) != 3 {
			t.Fatalf("expected 3 prescriptions, got %d", len(prescriptions))
		}

		// First two sets should not be AMRAP
		for i := 0; i < 2; i++ {
			var scheme setSchemeJSON
			if err := json.Unmarshal([]byte(prescriptions[i].SetScheme), &scheme); err != nil {
				t.Fatalf("failed to parse set scheme: %v", err)
			}
			if scheme.IsAmrap {
				t.Errorf("Set %d should not be AMRAP", i+1)
			}
		}

		// Last set should be AMRAP
		var lastScheme setSchemeJSON
		if err := json.Unmarshal([]byte(prescriptions[2].SetScheme), &lastScheme); err != nil {
			t.Fatalf("failed to parse set scheme: %v", err)
		}
		if !lastScheme.IsAmrap {
			t.Error("Last set (set 3) should be AMRAP")
		}
	})

	t.Run("cycle progression rules", func(t *testing.T) {
		progressions := getProgressionsForProgram(t, sqlDB, prog.ID)

		// Expected: Squat +10, Bench +5, Deadlift +10, OHP +5
		expectedIncrements := map[string]float64{
			"squat":          10.0,
			"bench-press":    5.0,
			"deadlift":       10.0,
			"overhead-press": 5.0,
		}

		for _, p := range progressions {
			expected, ok := expectedIncrements[p.LiftSlug]
			if !ok {
				continue // Skip unexpected lifts
			}

			var params struct {
				Increment   float64 `json:"increment"`
				TriggerType string  `json:"triggerType"`
			}
			if err := json.Unmarshal([]byte(p.Params), &params); err != nil {
				t.Fatalf("failed to parse progression params: %v", err)
			}

			if params.Increment != expected {
				t.Errorf("%s progression: expected +%.1f, got +%.1f", p.LiftSlug, expected, params.Increment)
			}

			if params.TriggerType != "AFTER_CYCLE" {
				t.Errorf("%s progression: expected AFTER_CYCLE trigger, got %s", p.LiftSlug, params.TriggerType)
			}
		}
	})
}

// =============================================================================
// GZCLP VERIFICATION
// =============================================================================

func TestCanonicalPrograms_GZCLP(t *testing.T) {
	sqlDB := setupCanonicalTestDB(t)
	defer sqlDB.Close()
	queries := db.New(sqlDB)

	prog := getProgramBySlug(t, queries, "gzclp")
	days := getDaysForProgram(t, sqlDB, prog.ID)

	t.Run("has exactly 4 workout days", func(t *testing.T) {
		if len(days) != 4 {
			t.Errorf("expected 4 days, got %d", len(days))
		}
	})

	t.Run("T1/T2 pairings are correct", func(t *testing.T) {
		// Expected pairings:
		// Day 1: T1 Squat, T2 Bench
		// Day 2: T1 OHP, T2 Deadlift
		// Day 3: T1 Bench, T2 Squat
		// Day 4: T1 Deadlift, T2 OHP
		expectedPairings := map[string]struct{ t1, t2 string }{
			"gzclp-day-1": {"squat", "bench-press"},
			"gzclp-day-2": {"overhead-press", "deadlift"},
			"gzclp-day-3": {"bench-press", "squat"},
			"gzclp-day-4": {"deadlift", "overhead-press"},
		}

		for _, day := range days {
			expected, ok := expectedPairings[day.Slug]
			if !ok {
				continue
			}

			prescriptions := getPrescriptionsForDay(t, sqlDB, day.ID)
			if len(prescriptions) != 2 {
				t.Errorf("%s: expected 2 prescriptions (T1 + T2), got %d", day.Slug, len(prescriptions))
				continue
			}

			// First prescription should be T1
			if prescriptions[0].LiftSlug != expected.t1 {
				t.Errorf("%s T1: expected %s, got %s", day.Slug, expected.t1, prescriptions[0].LiftSlug)
			}

			// Second prescription should be T2
			if prescriptions[1].LiftSlug != expected.t2 {
				t.Errorf("%s T2: expected %s, got %s", day.Slug, expected.t2, prescriptions[1].LiftSlug)
			}
		}
	})

	t.Run("T1 default scheme is 5x3+", func(t *testing.T) {
		// Check a T1 prescription
		var day1ID string
		for _, d := range days {
			if d.Slug == "gzclp-day-1" {
				day1ID = d.ID
				break
			}
		}

		prescriptions := getPrescriptionsForDay(t, sqlDB, day1ID)
		if len(prescriptions) < 1 {
			t.Fatal("no prescriptions found for Day 1")
		}

		var scheme setSchemeJSON
		if err := json.Unmarshal([]byte(prescriptions[0].SetScheme), &scheme); err != nil {
			t.Fatalf("failed to parse set scheme: %v", err)
		}

		if scheme.Sets != 5 {
			t.Errorf("T1 sets: expected 5, got %d", scheme.Sets)
		}
		if scheme.Reps != 3 {
			t.Errorf("T1 reps: expected 3, got %d", scheme.Reps)
		}
		if !scheme.IsAmrap {
			t.Error("T1 should be AMRAP (5x3+)")
		}
		if scheme.Tier != "T1" {
			t.Errorf("T1 tier: expected T1, got %s", scheme.Tier)
		}
	})

	t.Run("T2 default scheme is 3x10", func(t *testing.T) {
		// Check a T2 prescription
		var day1ID string
		for _, d := range days {
			if d.Slug == "gzclp-day-1" {
				day1ID = d.ID
				break
			}
		}

		prescriptions := getPrescriptionsForDay(t, sqlDB, day1ID)
		if len(prescriptions) < 2 {
			t.Fatal("not enough prescriptions found for Day 1")
		}

		var scheme setSchemeJSON
		if err := json.Unmarshal([]byte(prescriptions[1].SetScheme), &scheme); err != nil {
			t.Fatalf("failed to parse set scheme: %v", err)
		}

		if scheme.Sets != 3 {
			t.Errorf("T2 sets: expected 3, got %d", scheme.Sets)
		}
		if scheme.Reps != 10 {
			t.Errorf("T2 reps: expected 10, got %d", scheme.Reps)
		}
		if scheme.Tier != "T2" {
			t.Errorf("T2 tier: expected T2, got %s", scheme.Tier)
		}
	})

	t.Run("progression increments", func(t *testing.T) {
		progressions := getProgressionsForProgram(t, sqlDB, prog.ID)

		// Build a map of lift -> progressions
		// GZCLP has multiple progressions per lift (T1 and T2)
		t1Lower := 5.0  // T1 lower body
		t1Upper := 2.5  // T1 upper body
		t2All := 2.5    // T2 all lifts

		for _, p := range progressions {
			var params struct {
				Increment    float64 `json:"increment"`
				Tier         string  `json:"tier,omitempty"`
				LiftCategory string  `json:"liftCategory,omitempty"`
			}
			if err := json.Unmarshal([]byte(p.Params), &params); err != nil {
				t.Fatalf("failed to parse progression params: %v", err)
			}

			// Check expected increments based on lift and tier
			switch {
			case params.Tier == "T1" && params.LiftCategory == "LOWER":
				if params.Increment != t1Lower {
					t.Errorf("T1 Lower progression: expected +%.1f, got +%.1f", t1Lower, params.Increment)
				}
			case params.Tier == "T1" && params.LiftCategory == "UPPER":
				if params.Increment != t1Upper {
					t.Errorf("T1 Upper progression: expected +%.1f, got +%.1f", t1Upper, params.Increment)
				}
			case params.Tier == "T2":
				if params.Increment != t2All {
					t.Errorf("T2 progression: expected +%.1f, got +%.1f", t2All, params.Increment)
				}
			}
		}
	})
}
