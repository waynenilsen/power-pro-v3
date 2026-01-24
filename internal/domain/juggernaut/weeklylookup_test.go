package juggernaut

import (
	"testing"
)

func TestCreate531WeeklyLookup_ReturnsValidLookup(t *testing.T) {
	lookup := Create531WeeklyLookup("test-id", nil)

	if lookup == nil {
		t.Fatal("expected non-nil WeeklyLookup")
	}

	if lookup.ID != "test-id" {
		t.Errorf("expected ID %q, got %q", "test-id", lookup.ID)
	}

	if lookup.Name != "Inverted Juggernaut 5/3/1" {
		t.Errorf("expected Name %q, got %q", "Inverted Juggernaut 5/3/1", lookup.Name)
	}

	if lookup.ProgramID != nil {
		t.Errorf("expected nil ProgramID, got %v", lookup.ProgramID)
	}

	if len(lookup.Entries) != 4 {
		t.Errorf("expected 4 entries, got %d", len(lookup.Entries))
	}
}

func TestCreate531WeeklyLookup_WithProgramID(t *testing.T) {
	programID := "program-123"
	lookup := Create531WeeklyLookup("test-id", &programID)

	if lookup.ProgramID == nil {
		t.Fatal("expected non-nil ProgramID")
	}

	if *lookup.ProgramID != programID {
		t.Errorf("expected ProgramID %q, got %q", programID, *lookup.ProgramID)
	}
}

func TestCreate531WeeklyLookup_Week1Accumulation(t *testing.T) {
	lookup := Create531WeeklyLookup("test-id", nil)
	entry := lookup.GetByWeekNumber(1)

	if entry == nil {
		t.Fatal("expected non-nil entry for week 1")
	}

	expectedPercentages := []float64{65.0, 75.0, 85.0, 75.0, 65.0}
	if len(entry.Percentages) != len(expectedPercentages) {
		t.Fatalf("expected %d percentages, got %d", len(expectedPercentages), len(entry.Percentages))
	}
	for i, pct := range expectedPercentages {
		if entry.Percentages[i] != pct {
			t.Errorf("week 1 percentages[%d]: expected %.1f, got %.1f", i, pct, entry.Percentages[i])
		}
	}

	expectedReps := []int{5, 5, -5, 5, -5}
	if len(entry.Reps) != len(expectedReps) {
		t.Fatalf("expected %d reps, got %d", len(expectedReps), len(entry.Reps))
	}
	for i, rep := range expectedReps {
		if entry.Reps[i] != rep {
			t.Errorf("week 1 reps[%d]: expected %d, got %d", i, rep, entry.Reps[i])
		}
	}
}

func TestCreate531WeeklyLookup_Week2Intensification(t *testing.T) {
	lookup := Create531WeeklyLookup("test-id", nil)
	entry := lookup.GetByWeekNumber(2)

	if entry == nil {
		t.Fatal("expected non-nil entry for week 2")
	}

	expectedPercentages := []float64{70.0, 80.0, 90.0, 80.0, 70.0}
	if len(entry.Percentages) != len(expectedPercentages) {
		t.Fatalf("expected %d percentages, got %d", len(expectedPercentages), len(entry.Percentages))
	}
	for i, pct := range expectedPercentages {
		if entry.Percentages[i] != pct {
			t.Errorf("week 2 percentages[%d]: expected %.1f, got %.1f", i, pct, entry.Percentages[i])
		}
	}

	expectedReps := []int{3, 3, -3, 3, -3}
	if len(entry.Reps) != len(expectedReps) {
		t.Fatalf("expected %d reps, got %d", len(expectedReps), len(entry.Reps))
	}
	for i, rep := range expectedReps {
		if entry.Reps[i] != rep {
			t.Errorf("week 2 reps[%d]: expected %d, got %d", i, rep, entry.Reps[i])
		}
	}
}

func TestCreate531WeeklyLookup_Week3Realization(t *testing.T) {
	lookup := Create531WeeklyLookup("test-id", nil)
	entry := lookup.GetByWeekNumber(3)

	if entry == nil {
		t.Fatal("expected non-nil entry for week 3")
	}

	expectedPercentages := []float64{75.0, 85.0, 95.0, 85.0, 75.0}
	if len(entry.Percentages) != len(expectedPercentages) {
		t.Fatalf("expected %d percentages, got %d", len(expectedPercentages), len(entry.Percentages))
	}
	for i, pct := range expectedPercentages {
		if entry.Percentages[i] != pct {
			t.Errorf("week 3 percentages[%d]: expected %.1f, got %.1f", i, pct, entry.Percentages[i])
		}
	}

	expectedReps := []int{5, 3, -1, 3, -5}
	if len(entry.Reps) != len(expectedReps) {
		t.Fatalf("expected %d reps, got %d", len(expectedReps), len(entry.Reps))
	}
	for i, rep := range expectedReps {
		if entry.Reps[i] != rep {
			t.Errorf("week 3 reps[%d]: expected %d, got %d", i, rep, entry.Reps[i])
		}
	}
}

func TestCreate531WeeklyLookup_Week4Deload(t *testing.T) {
	lookup := Create531WeeklyLookup("test-id", nil)
	entry := lookup.GetByWeekNumber(4)

	if entry == nil {
		t.Fatal("expected non-nil entry for week 4")
	}

	// Deload has only 3 sets
	expectedPercentages := []float64{40.0, 50.0, 60.0}
	if len(entry.Percentages) != len(expectedPercentages) {
		t.Fatalf("week 4 (deload) expected %d percentages (sets), got %d", len(expectedPercentages), len(entry.Percentages))
	}
	for i, pct := range expectedPercentages {
		if entry.Percentages[i] != pct {
			t.Errorf("week 4 percentages[%d]: expected %.1f, got %.1f", i, pct, entry.Percentages[i])
		}
	}

	expectedReps := []int{5, 5, 5}
	if len(entry.Reps) != len(expectedReps) {
		t.Fatalf("week 4 (deload) expected %d reps, got %d", len(expectedReps), len(entry.Reps))
	}
	for i, rep := range expectedReps {
		if entry.Reps[i] != rep {
			t.Errorf("week 4 reps[%d]: expected %d, got %d", i, rep, entry.Reps[i])
		}
	}
}

func TestCreate531WeeklyLookup_DeloadHasOnlyThreeSets(t *testing.T) {
	lookup := Create531WeeklyLookup("test-id", nil)
	entry := lookup.GetByWeekNumber(4)

	if entry == nil {
		t.Fatal("expected non-nil entry for week 4")
	}

	if len(entry.Percentages) != 3 {
		t.Errorf("deload week should have 3 sets, got %d", len(entry.Percentages))
	}

	if len(entry.Reps) != 3 {
		t.Errorf("deload week should have 3 reps, got %d", len(entry.Reps))
	}
}

func TestCreate531WeeklyLookup_GetByWeekNumber(t *testing.T) {
	lookup := Create531WeeklyLookup("test-id", nil)

	// Valid weeks
	for week := 1; week <= 4; week++ {
		entry := lookup.GetByWeekNumber(week)
		if entry == nil {
			t.Errorf("expected non-nil entry for week %d", week)
		}
		if entry.WeekNumber != week {
			t.Errorf("week %d: expected WeekNumber %d, got %d", week, week, entry.WeekNumber)
		}
	}

	// Invalid week returns nil
	entry := lookup.GetByWeekNumber(5)
	if entry != nil {
		t.Error("expected nil entry for week 5")
	}

	entry = lookup.GetByWeekNumber(0)
	if entry != nil {
		t.Error("expected nil entry for week 0")
	}
}

func TestCreate531WeeklyLookup_PercentagesAndRepsMatchLength(t *testing.T) {
	lookup := Create531WeeklyLookup("test-id", nil)

	for _, entry := range lookup.Entries {
		if len(entry.Percentages) != len(entry.Reps) {
			t.Errorf("week %d: percentages length (%d) != reps length (%d)",
				entry.WeekNumber, len(entry.Percentages), len(entry.Reps))
		}
	}
}

func TestCreate531WeeklyLookup_NonDeloadWeeksHaveFiveSets(t *testing.T) {
	lookup := Create531WeeklyLookup("test-id", nil)

	for week := 1; week <= 3; week++ {
		entry := lookup.GetByWeekNumber(week)
		if entry == nil {
			t.Fatalf("expected non-nil entry for week %d", week)
		}
		if len(entry.Percentages) != 5 {
			t.Errorf("week %d: expected 5 sets, got %d", week, len(entry.Percentages))
		}
	}
}

func TestCreate531WeeklyLookup_TimestampsAreSet(t *testing.T) {
	lookup := Create531WeeklyLookup("test-id", nil)

	if lookup.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}

	if lookup.UpdatedAt.IsZero() {
		t.Error("expected non-zero UpdatedAt")
	}
}
