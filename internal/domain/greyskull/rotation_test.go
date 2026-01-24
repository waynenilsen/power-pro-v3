package greyskull

import (
	"errors"
	"testing"
)

// TestGetVariantForDay tests the A/B rotation logic.
func TestGetVariantForDay(t *testing.T) {
	tests := []struct {
		name        string
		weekNumber  int
		dayPosition int
		want        Variant
		wantErr     error
	}{
		// Week 1 (odd): A, B, A
		{name: "Week 1, Day 1", weekNumber: 1, dayPosition: 1, want: VariantA},
		{name: "Week 1, Day 2", weekNumber: 1, dayPosition: 2, want: VariantB},
		{name: "Week 1, Day 3", weekNumber: 1, dayPosition: 3, want: VariantA},

		// Week 2 (even): B, A, B
		{name: "Week 2, Day 1", weekNumber: 2, dayPosition: 1, want: VariantB},
		{name: "Week 2, Day 2", weekNumber: 2, dayPosition: 2, want: VariantA},
		{name: "Week 2, Day 3", weekNumber: 2, dayPosition: 3, want: VariantB},

		// Week 3 (odd): A, B, A
		{name: "Week 3, Day 1", weekNumber: 3, dayPosition: 1, want: VariantA},
		{name: "Week 3, Day 2", weekNumber: 3, dayPosition: 2, want: VariantB},
		{name: "Week 3, Day 3", weekNumber: 3, dayPosition: 3, want: VariantA},

		// Week 4 (even): B, A, B
		{name: "Week 4, Day 1", weekNumber: 4, dayPosition: 1, want: VariantB},
		{name: "Week 4, Day 2", weekNumber: 4, dayPosition: 2, want: VariantA},
		{name: "Week 4, Day 3", weekNumber: 4, dayPosition: 3, want: VariantB},

		// Higher week numbers
		{name: "Week 10 (even), Day 2", weekNumber: 10, dayPosition: 2, want: VariantA},
		{name: "Week 11 (odd), Day 1", weekNumber: 11, dayPosition: 1, want: VariantA},
		{name: "Week 100 (even), Day 3", weekNumber: 100, dayPosition: 3, want: VariantB},
		{name: "Week 99 (odd), Day 2", weekNumber: 99, dayPosition: 2, want: VariantB},

		// Edge cases - errors
		{name: "Week 0 (invalid)", weekNumber: 0, dayPosition: 1, wantErr: ErrInvalidWeekNumber},
		{name: "Week -1 (invalid)", weekNumber: -1, dayPosition: 1, wantErr: ErrInvalidWeekNumber},
		{name: "Day 0 (invalid)", weekNumber: 1, dayPosition: 0, wantErr: ErrInvalidDayPosition},
		{name: "Day 4 (invalid)", weekNumber: 1, dayPosition: 4, wantErr: ErrInvalidDayPosition},
		{name: "Day -1 (invalid)", weekNumber: 1, dayPosition: -1, wantErr: ErrInvalidDayPosition},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetVariantForDay(tt.weekNumber, tt.dayPosition)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
				} else if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("GetVariantForDay(%d, %d) = %v, want %v", tt.weekNumber, tt.dayPosition, got, tt.want)
			}
		})
	}
}

// TestGetVariantForDay_FullCycle tests multiple weeks to verify the pattern holds.
func TestGetVariantForDay_FullCycle(t *testing.T) {
	// Verify the pattern over 8 weeks
	expected := [][]Variant{
		{VariantA, VariantB, VariantA}, // Week 1 (odd)
		{VariantB, VariantA, VariantB}, // Week 2 (even)
		{VariantA, VariantB, VariantA}, // Week 3 (odd)
		{VariantB, VariantA, VariantB}, // Week 4 (even)
		{VariantA, VariantB, VariantA}, // Week 5 (odd)
		{VariantB, VariantA, VariantB}, // Week 6 (even)
		{VariantA, VariantB, VariantA}, // Week 7 (odd)
		{VariantB, VariantA, VariantB}, // Week 8 (even)
	}

	for week := 1; week <= 8; week++ {
		for day := 1; day <= 3; day++ {
			got, err := GetVariantForDay(week, day)
			if err != nil {
				t.Errorf("Week %d, Day %d: unexpected error: %v", week, day, err)
				continue
			}
			want := expected[week-1][day-1]
			if got != want {
				t.Errorf("Week %d, Day %d: got %v, want %v", week, day, got, want)
			}
		}
	}
}

// TestGetVariantString tests the string pointer helper.
func TestGetVariantString(t *testing.T) {
	t.Run("valid variant A", func(t *testing.T) {
		got, err := GetVariantString(1, 1)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got == nil {
			t.Error("expected non-nil string pointer")
		}
		if *got != "A" {
			t.Errorf("expected 'A', got '%s'", *got)
		}
	})

	t.Run("valid variant B", func(t *testing.T) {
		got, err := GetVariantString(1, 2)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got == nil {
			t.Error("expected non-nil string pointer")
		}
		if *got != "B" {
			t.Errorf("expected 'B', got '%s'", *got)
		}
	})

	t.Run("invalid week number", func(t *testing.T) {
		got, err := GetVariantString(0, 1)
		if err == nil {
			t.Error("expected error for invalid week number")
		}
		if got != nil {
			t.Error("expected nil result for error case")
		}
	})

	t.Run("invalid day position", func(t *testing.T) {
		got, err := GetVariantString(1, 4)
		if err == nil {
			t.Error("expected error for invalid day position")
		}
		if got != nil {
			t.Error("expected nil result for error case")
		}
	})
}

// TestGetDayTemplate tests the day template retrieval.
func TestGetDayTemplate(t *testing.T) {
	t.Run("Variant A (Bench Day)", func(t *testing.T) {
		lifts, err := GetDayTemplate(VariantA)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(lifts) != 5 {
			t.Errorf("expected 5 lifts, got %d", len(lifts))
		}

		// Check main lifts
		expectedMain := []string{"bench-press", "barbell-row", "squat"}
		for i, expected := range expectedMain {
			if lifts[i].Slug != expected {
				t.Errorf("lift %d: expected slug '%s', got '%s'", i, expected, lifts[i].Slug)
			}
			if lifts[i].Sets != 3 {
				t.Errorf("lift %d (%s): expected 3 sets, got %d", i, expected, lifts[i].Sets)
			}
			if lifts[i].Reps != 5 {
				t.Errorf("lift %d (%s): expected 5 reps, got %d", i, expected, lifts[i].Reps)
			}
			if !lifts[i].IsAMRAP {
				t.Errorf("lift %d (%s): expected IsAMRAP=true", i, expected)
			}
		}

		// Check accessories
		expectedAccessory := []string{"tricep-extension", "ab-rollout"}
		for i, expected := range expectedAccessory {
			idx := i + 3
			if lifts[idx].Slug != expected {
				t.Errorf("lift %d: expected slug '%s', got '%s'", idx, expected, lifts[idx].Slug)
			}
			if lifts[idx].Sets != 3 {
				t.Errorf("lift %d (%s): expected 3 sets, got %d", idx, expected, lifts[idx].Sets)
			}
		}
	})

	t.Run("Variant B (OHP Day)", func(t *testing.T) {
		lifts, err := GetDayTemplate(VariantB)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(lifts) != 5 {
			t.Errorf("expected 5 lifts, got %d", len(lifts))
		}

		// Check main lifts
		expectedMain := []string{"overhead-press", "chin-up", "deadlift"}
		for i, expected := range expectedMain {
			if lifts[i].Slug != expected {
				t.Errorf("lift %d: expected slug '%s', got '%s'", i, expected, lifts[i].Slug)
			}
			if lifts[i].Sets != 3 {
				t.Errorf("lift %d (%s): expected 3 sets, got %d", i, expected, lifts[i].Sets)
			}
			if lifts[i].Reps != 5 {
				t.Errorf("lift %d (%s): expected 5 reps, got %d", i, expected, lifts[i].Reps)
			}
			if !lifts[i].IsAMRAP {
				t.Errorf("lift %d (%s): expected IsAMRAP=true", i, expected)
			}
		}

		// Check accessories
		expectedAccessory := []string{"bicep-curl", "shrug"}
		for i, expected := range expectedAccessory {
			idx := i + 3
			if lifts[idx].Slug != expected {
				t.Errorf("lift %d: expected slug '%s', got '%s'", idx, expected, lifts[idx].Slug)
			}
		}
	})

	t.Run("invalid variant", func(t *testing.T) {
		_, err := GetDayTemplate("C")
		if err == nil {
			t.Error("expected error for invalid variant")
		}
		if !errors.Is(err, ErrInvalidVariant) {
			t.Errorf("expected ErrInvalidVariant, got %v", err)
		}
	})

	t.Run("empty variant", func(t *testing.T) {
		_, err := GetDayTemplate("")
		if err == nil {
			t.Error("expected error for empty variant")
		}
	})
}

// TestGetDayTemplateForWeek tests the convenience function.
func TestGetDayTemplateForWeek(t *testing.T) {
	t.Run("Week 1, Day 1 (Variant A)", func(t *testing.T) {
		lifts, err := GetDayTemplateForWeek(1, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if lifts[0].Slug != "bench-press" {
			t.Errorf("expected bench-press, got %s", lifts[0].Slug)
		}
	})

	t.Run("Week 1, Day 2 (Variant B)", func(t *testing.T) {
		lifts, err := GetDayTemplateForWeek(1, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if lifts[0].Slug != "overhead-press" {
			t.Errorf("expected overhead-press, got %s", lifts[0].Slug)
		}
	})

	t.Run("invalid week number", func(t *testing.T) {
		_, err := GetDayTemplateForWeek(0, 1)
		if err == nil {
			t.Error("expected error for invalid week number")
		}
	})

	t.Run("invalid day position", func(t *testing.T) {
		_, err := GetDayTemplateForWeek(1, 5)
		if err == nil {
			t.Error("expected error for invalid day position")
		}
	})
}

// TestGetLiftSlugs tests the lift slug extraction.
func TestGetLiftSlugs(t *testing.T) {
	t.Run("Variant A", func(t *testing.T) {
		slugs, err := GetLiftSlugs(VariantA)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"bench-press", "barbell-row", "squat", "tricep-extension", "ab-rollout"}
		if len(slugs) != len(expected) {
			t.Errorf("expected %d slugs, got %d", len(expected), len(slugs))
		}
		for i, exp := range expected {
			if slugs[i] != exp {
				t.Errorf("slug %d: expected '%s', got '%s'", i, exp, slugs[i])
			}
		}
	})

	t.Run("Variant B", func(t *testing.T) {
		slugs, err := GetLiftSlugs(VariantB)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"overhead-press", "chin-up", "deadlift", "bicep-curl", "shrug"}
		if len(slugs) != len(expected) {
			t.Errorf("expected %d slugs, got %d", len(expected), len(slugs))
		}
		for i, exp := range expected {
			if slugs[i] != exp {
				t.Errorf("slug %d: expected '%s', got '%s'", i, exp, slugs[i])
			}
		}
	})

	t.Run("invalid variant", func(t *testing.T) {
		_, err := GetLiftSlugs("X")
		if err == nil {
			t.Error("expected error for invalid variant")
		}
	})
}

// TestGetMainLiftSlugs tests main lift extraction.
func TestGetMainLiftSlugs(t *testing.T) {
	t.Run("Variant A", func(t *testing.T) {
		slugs, err := GetMainLiftSlugs(VariantA)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"bench-press", "barbell-row", "squat"}
		if len(slugs) != 3 {
			t.Errorf("expected 3 main lifts, got %d", len(slugs))
		}
		for i, exp := range expected {
			if slugs[i] != exp {
				t.Errorf("slug %d: expected '%s', got '%s'", i, exp, slugs[i])
			}
		}
	})

	t.Run("Variant B", func(t *testing.T) {
		slugs, err := GetMainLiftSlugs(VariantB)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"overhead-press", "chin-up", "deadlift"}
		if len(slugs) != 3 {
			t.Errorf("expected 3 main lifts, got %d", len(slugs))
		}
		for i, exp := range expected {
			if slugs[i] != exp {
				t.Errorf("slug %d: expected '%s', got '%s'", i, exp, slugs[i])
			}
		}
	})

	t.Run("invalid variant", func(t *testing.T) {
		_, err := GetMainLiftSlugs("Z")
		if err == nil {
			t.Error("expected error for invalid variant")
		}
	})
}

// TestGetAccessoryLiftSlugs tests accessory lift extraction.
func TestGetAccessoryLiftSlugs(t *testing.T) {
	t.Run("Variant A", func(t *testing.T) {
		slugs, err := GetAccessoryLiftSlugs(VariantA)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"tricep-extension", "ab-rollout"}
		if len(slugs) != 2 {
			t.Errorf("expected 2 accessories, got %d", len(slugs))
		}
		for i, exp := range expected {
			if slugs[i] != exp {
				t.Errorf("slug %d: expected '%s', got '%s'", i, exp, slugs[i])
			}
		}
	})

	t.Run("Variant B", func(t *testing.T) {
		slugs, err := GetAccessoryLiftSlugs(VariantB)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"bicep-curl", "shrug"}
		if len(slugs) != 2 {
			t.Errorf("expected 2 accessories, got %d", len(slugs))
		}
		for i, exp := range expected {
			if slugs[i] != exp {
				t.Errorf("slug %d: expected '%s', got '%s'", i, exp, slugs[i])
			}
		}
	})

	t.Run("invalid variant", func(t *testing.T) {
		_, err := GetAccessoryLiftSlugs("Y")
		if err == nil {
			t.Error("expected error for invalid variant")
		}
	})
}

// TestLiftInfo_Structure tests the LiftInfo struct fields.
func TestLiftInfo_Structure(t *testing.T) {
	lifts, err := GetDayTemplate(VariantA)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test that we can access all fields
	for _, lift := range lifts {
		if lift.Slug == "" {
			t.Error("lift slug should not be empty")
		}
		if lift.Sets <= 0 {
			t.Errorf("lift %s: sets should be > 0, got %d", lift.Slug, lift.Sets)
		}
		if lift.Reps <= 0 {
			t.Errorf("lift %s: reps should be > 0, got %d", lift.Slug, lift.Reps)
		}
		// IsAMRAP is a bool, just verify it's accessible
		_ = lift.IsAMRAP
	}
}

// TestVariantConstants tests that variant constants have expected values.
func TestVariantConstants(t *testing.T) {
	if VariantA != "A" {
		t.Errorf("VariantA should be 'A', got '%s'", VariantA)
	}
	if VariantB != "B" {
		t.Errorf("VariantB should be 'B', got '%s'", VariantB)
	}
}

// TestWeekVariantAlternation verifies that each day's variant alternates correctly.
func TestWeekVariantAlternation(t *testing.T) {
	// For each day position, verify alternation between odd and even weeks
	for day := 1; day <= 3; day++ {
		oddWeek, _ := GetVariantForDay(1, day)
		evenWeek, _ := GetVariantForDay(2, day)

		if oddWeek == evenWeek {
			t.Errorf("Day %d: variant should alternate between weeks (both are %v)", day, oddWeek)
		}

		// Verify pattern continues
		oddWeek3, _ := GetVariantForDay(3, day)
		evenWeek4, _ := GetVariantForDay(4, day)

		if oddWeek != oddWeek3 {
			t.Errorf("Day %d: odd weeks should have same variant (week 1: %v, week 3: %v)", day, oddWeek, oddWeek3)
		}
		if evenWeek != evenWeek4 {
			t.Errorf("Day %d: even weeks should have same variant (week 2: %v, week 4: %v)", day, evenWeek, evenWeek4)
		}
	}
}

// TestABPattern_Documentation verifies the documented A/B pattern is correct.
func TestABPattern_Documentation(t *testing.T) {
	// From README.md:
	// Week 1: Day 1 (A), Day 2 (B), Day 3 (A)
	// Week 2: Day 1 (B), Day 2 (A), Day 3 (B)

	// Week 1
	w1d1, _ := GetVariantForDay(1, 1)
	w1d2, _ := GetVariantForDay(1, 2)
	w1d3, _ := GetVariantForDay(1, 3)
	if w1d1 != VariantA || w1d2 != VariantB || w1d3 != VariantA {
		t.Errorf("Week 1 pattern incorrect: got %v, %v, %v; want A, B, A", w1d1, w1d2, w1d3)
	}

	// Week 2
	w2d1, _ := GetVariantForDay(2, 1)
	w2d2, _ := GetVariantForDay(2, 2)
	w2d3, _ := GetVariantForDay(2, 3)
	if w2d1 != VariantB || w2d2 != VariantA || w2d3 != VariantB {
		t.Errorf("Week 2 pattern incorrect: got %v, %v, %v; want B, A, B", w2d1, w2d2, w2d3)
	}
}
