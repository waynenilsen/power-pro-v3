package loadstrategy

import (
	"math"
	"testing"

	"github.com/waynenilsen/power-pro-v3/internal/domain/dailylookup"
	"github.com/waynenilsen/power-pro-v3/internal/domain/weeklylookup"
)

func TestLookupContext_HasWeeklyLookup(t *testing.T) {
	tests := []struct {
		name     string
		context  *LookupContext
		expected bool
	}{
		{
			name:     "nil context",
			context:  nil,
			expected: false,
		},
		{
			name:     "nil weekly lookup",
			context:  &LookupContext{WeekNumber: 1},
			expected: false,
		},
		{
			name: "zero week number",
			context: &LookupContext{
				WeeklyLookup: &weeklylookup.WeeklyLookup{},
				WeekNumber:   0,
			},
			expected: false,
		},
		{
			name: "valid weekly lookup",
			context: &LookupContext{
				WeeklyLookup: &weeklylookup.WeeklyLookup{},
				WeekNumber:   1,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.context.HasWeeklyLookup()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestLookupContext_HasDailyLookup(t *testing.T) {
	tests := []struct {
		name     string
		context  *LookupContext
		expected bool
	}{
		{
			name:     "nil context",
			context:  nil,
			expected: false,
		},
		{
			name:     "nil daily lookup",
			context:  &LookupContext{DaySlug: "heavy"},
			expected: false,
		},
		{
			name: "empty day slug",
			context: &LookupContext{
				DailyLookup: &dailylookup.DailyLookup{},
				DaySlug:     "",
			},
			expected: false,
		},
		{
			name: "valid daily lookup",
			context: &LookupContext{
				DailyLookup: &dailylookup.DailyLookup{},
				DaySlug:     "heavy",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.context.HasDailyLookup()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestLookupContext_GetWeeklyEntry(t *testing.T) {
	weeklyLookup := &weeklylookup.WeeklyLookup{
		Entries: []weeklylookup.WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65, 75, 85}, Reps: []int{5, 5, 5}},
			{WeekNumber: 2, Percentages: []float64{70, 80, 90}, Reps: []int{3, 3, 3}},
			{WeekNumber: 3, Percentages: []float64{75, 85, 95}, Reps: []int{5, 3, 1}},
		},
	}

	tests := []struct {
		name        string
		context     *LookupContext
		expectNil   bool
		expectWeek  int
	}{
		{
			name:      "nil context",
			context:   nil,
			expectNil: true,
		},
		{
			name: "week 1",
			context: &LookupContext{
				WeeklyLookup: weeklyLookup,
				WeekNumber:   1,
			},
			expectNil:  false,
			expectWeek: 1,
		},
		{
			name: "week 2",
			context: &LookupContext{
				WeeklyLookup: weeklyLookup,
				WeekNumber:   2,
			},
			expectNil:  false,
			expectWeek: 2,
		},
		{
			name: "week not found",
			context: &LookupContext{
				WeeklyLookup: weeklyLookup,
				WeekNumber:   5,
			},
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := tt.context.GetWeeklyEntry()
			if tt.expectNil {
				if entry != nil {
					t.Error("expected nil entry")
				}
				return
			}
			if entry == nil {
				t.Fatal("expected non-nil entry")
			}
			if entry.WeekNumber != tt.expectWeek {
				t.Errorf("expected week %d, got %d", tt.expectWeek, entry.WeekNumber)
			}
		})
	}
}

func TestLookupContext_GetDailyEntry(t *testing.T) {
	dailyLookup := &dailylookup.DailyLookup{
		Entries: []dailylookup.DailyLookupEntry{
			{DayIdentifier: "heavy", PercentageModifier: 100},
			{DayIdentifier: "light", PercentageModifier: 70},
			{DayIdentifier: "medium", PercentageModifier: 80},
		},
	}

	tests := []struct {
		name            string
		context         *LookupContext
		expectNil       bool
		expectModifier  float64
	}{
		{
			name:      "nil context",
			context:   nil,
			expectNil: true,
		},
		{
			name: "heavy day",
			context: &LookupContext{
				DailyLookup: dailyLookup,
				DaySlug:     "heavy",
			},
			expectNil:      false,
			expectModifier: 100,
		},
		{
			name: "light day case insensitive",
			context: &LookupContext{
				DailyLookup: dailyLookup,
				DaySlug:     "LIGHT",
			},
			expectNil:      false,
			expectModifier: 70,
		},
		{
			name: "day not found",
			context: &LookupContext{
				DailyLookup: dailyLookup,
				DaySlug:     "extra-heavy",
			},
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := tt.context.GetDailyEntry()
			if tt.expectNil {
				if entry != nil {
					t.Error("expected nil entry")
				}
				return
			}
			if entry == nil {
				t.Fatal("expected non-nil entry")
			}
			if entry.PercentageModifier != tt.expectModifier {
				t.Errorf("expected modifier %f, got %f", tt.expectModifier, entry.PercentageModifier)
			}
		})
	}
}

func TestLookupContext_ApplyModifiers_NoLookups(t *testing.T) {
	// Test with nil context
	var nilCtx *LookupContext
	result := nilCtx.ApplyModifiers(85.0)
	if result != 85.0 {
		t.Errorf("nil context: expected 85.0, got %f", result)
	}

	// Test with empty context
	emptyCtx := &LookupContext{}
	result = emptyCtx.ApplyModifiers(85.0)
	if result != 85.0 {
		t.Errorf("empty context: expected 85.0, got %f", result)
	}
}

func TestLookupContext_ApplyModifiers_WeeklySetSpecificPercentage(t *testing.T) {
	// 5/3/1 style weekly lookup with set-specific percentages
	weeklyLookup := &weeklylookup.WeeklyLookup{
		Entries: []weeklylookup.WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65, 75, 85}, Reps: []int{5, 5, 5}},
			{WeekNumber: 2, Percentages: []float64{70, 80, 90}, Reps: []int{3, 3, 3}},
			{WeekNumber: 3, Percentages: []float64{75, 85, 95}, Reps: []int{5, 3, 1}},
		},
	}

	tests := []struct {
		name           string
		weekNumber     int
		setNumber      int
		basePercentage float64
		expected       float64
	}{
		{
			name:           "week 1 set 1",
			weekNumber:     1,
			setNumber:      1,
			basePercentage: 100, // base doesn't matter when set-specific is used
			expected:       65,
		},
		{
			name:           "week 1 set 2",
			weekNumber:     1,
			setNumber:      2,
			basePercentage: 100,
			expected:       75,
		},
		{
			name:           "week 1 set 3",
			weekNumber:     1,
			setNumber:      3,
			basePercentage: 100,
			expected:       85,
		},
		{
			name:           "week 2 set 3",
			weekNumber:     2,
			setNumber:      3,
			basePercentage: 100,
			expected:       90,
		},
		{
			name:           "week 3 set 3 (top set)",
			weekNumber:     3,
			setNumber:      3,
			basePercentage: 100,
			expected:       95,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &LookupContext{
				WeeklyLookup: weeklyLookup,
				WeekNumber:   tt.weekNumber,
				SetNumber:    tt.setNumber,
			}

			result := ctx.ApplyModifiers(tt.basePercentage)
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestLookupContext_ApplyModifiers_WeeklyPercentageModifier(t *testing.T) {
	// Deload week style with percentage modifier
	modifier90 := 90.0
	modifier80 := 80.0
	weeklyLookup := &weeklylookup.WeeklyLookup{
		Entries: []weeklylookup.WeeklyLookupEntry{
			{WeekNumber: 1, PercentageModifier: nil},                  // No modification
			{WeekNumber: 2, PercentageModifier: &modifier90},         // 90% of prescribed
			{WeekNumber: 4, PercentageModifier: &modifier80},         // Deload: 80% of prescribed
		},
	}

	tests := []struct {
		name           string
		weekNumber     int
		basePercentage float64
		expected       float64
	}{
		{
			name:           "week 1 no modifier",
			weekNumber:     1,
			basePercentage: 85,
			expected:       85, // No modifier, base percentage unchanged
		},
		{
			name:           "week 2 90% modifier",
			weekNumber:     2,
			basePercentage: 100, // 100 * (90/100) = 90
			expected:       90,
		},
		{
			name:           "week 4 deload 80% modifier",
			weekNumber:     4,
			basePercentage: 85, // 85 * (80/100) = 68
			expected:       68,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &LookupContext{
				WeeklyLookup: weeklyLookup,
				WeekNumber:   tt.weekNumber,
				SetNumber:    0, // No set-specific percentage
			}

			result := ctx.ApplyModifiers(tt.basePercentage)
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestLookupContext_ApplyModifiers_DailyOnly(t *testing.T) {
	// Bill Starr Heavy/Light/Medium style
	dailyLookup := &dailylookup.DailyLookup{
		Entries: []dailylookup.DailyLookupEntry{
			{DayIdentifier: "heavy", PercentageModifier: 100},
			{DayIdentifier: "light", PercentageModifier: 70},
			{DayIdentifier: "medium", PercentageModifier: 80},
		},
	}

	tests := []struct {
		name           string
		daySlug        string
		basePercentage float64
		expected       float64
	}{
		{
			name:           "heavy day",
			daySlug:        "heavy",
			basePercentage: 85, // 85 * (100/100) = 85
			expected:       85,
		},
		{
			name:           "light day",
			daySlug:        "light",
			basePercentage: 85, // 85 * (70/100) = 59.5
			expected:       59.5,
		},
		{
			name:           "medium day",
			daySlug:        "medium",
			basePercentage: 85, // 85 * (80/100) = 68
			expected:       68,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &LookupContext{
				DailyLookup: dailyLookup,
				DaySlug:     tt.daySlug,
			}

			result := ctx.ApplyModifiers(tt.basePercentage)
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestLookupContext_ApplyModifiers_WeeklyAndDaily(t *testing.T) {
	// Combined: weekly set-specific percentages + daily intensity modifier
	weeklyLookup := &weeklylookup.WeeklyLookup{
		Entries: []weeklylookup.WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65, 75, 85}, Reps: []int{5, 5, 5}},
		},
	}

	dailyLookup := &dailylookup.DailyLookup{
		Entries: []dailylookup.DailyLookupEntry{
			{DayIdentifier: "heavy", PercentageModifier: 100},
			{DayIdentifier: "light", PercentageModifier: 70},
		},
	}

	tests := []struct {
		name       string
		weekNumber int
		setNumber  int
		daySlug    string
		expected   float64
	}{
		{
			name:       "week 1 set 1 heavy day",
			weekNumber: 1,
			setNumber:  1,
			daySlug:    "heavy",
			expected:   65, // 65 * (100/100) = 65
		},
		{
			name:       "week 1 set 1 light day",
			weekNumber: 1,
			setNumber:  1,
			daySlug:    "light",
			expected:   45.5, // 65 * (70/100) = 45.5
		},
		{
			name:       "week 1 set 3 light day",
			weekNumber: 1,
			setNumber:  3,
			daySlug:    "light",
			expected:   59.5, // 85 * (70/100) = 59.5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &LookupContext{
				WeeklyLookup: weeklyLookup,
				WeekNumber:   tt.weekNumber,
				SetNumber:    tt.setNumber,
				DailyLookup:  dailyLookup,
				DaySlug:      tt.daySlug,
			}

			result := ctx.ApplyModifiers(100) // Base doesn't matter with set-specific
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestLookupContext_ApplyModifiers_WeeklyModifierAndDaily(t *testing.T) {
	// Combined: weekly percentage modifier + daily intensity modifier
	modifier90 := 90.0
	weeklyLookup := &weeklylookup.WeeklyLookup{
		Entries: []weeklylookup.WeeklyLookupEntry{
			{WeekNumber: 1, PercentageModifier: &modifier90},
		},
	}

	dailyLookup := &dailylookup.DailyLookup{
		Entries: []dailylookup.DailyLookupEntry{
			{DayIdentifier: "light", PercentageModifier: 70},
		},
	}

	ctx := &LookupContext{
		WeeklyLookup: weeklyLookup,
		WeekNumber:   1,
		DailyLookup:  dailyLookup,
		DaySlug:      "light",
	}

	// Base 100 -> weekly 90% -> 90 -> daily 70% -> 63
	result := ctx.ApplyModifiers(100)
	expected := 63.0
	if math.Abs(result-expected) > 0.0001 {
		t.Errorf("expected %f, got %f", expected, result)
	}
}

func TestLookupContext_GetRepsForSet(t *testing.T) {
	weeklyLookup := &weeklylookup.WeeklyLookup{
		Entries: []weeklylookup.WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65, 75, 85}, Reps: []int{5, 5, 5}},
			{WeekNumber: 2, Percentages: []float64{70, 80, 90}, Reps: []int{3, 3, 3}},
			{WeekNumber: 3, Percentages: []float64{75, 85, 95}, Reps: []int{5, 3, 1}},
		},
	}

	tests := []struct {
		name       string
		context    *LookupContext
		expected   int
	}{
		{
			name:     "nil context",
			context:  nil,
			expected: -1,
		},
		{
			name: "zero set number",
			context: &LookupContext{
				WeeklyLookup: weeklyLookup,
				WeekNumber:   1,
				SetNumber:    0,
			},
			expected: -1,
		},
		{
			name: "week 1 set 1",
			context: &LookupContext{
				WeeklyLookup: weeklyLookup,
				WeekNumber:   1,
				SetNumber:    1,
			},
			expected: 5,
		},
		{
			name: "week 3 set 3",
			context: &LookupContext{
				WeeklyLookup: weeklyLookup,
				WeekNumber:   3,
				SetNumber:    3,
			},
			expected: 1,
		},
		{
			name: "set out of range",
			context: &LookupContext{
				WeeklyLookup: weeklyLookup,
				WeekNumber:   1,
				SetNumber:    10,
			},
			expected: -1,
		},
		{
			name: "week not found",
			context: &LookupContext{
				WeeklyLookup: weeklyLookup,
				WeekNumber:   5,
				SetNumber:    1,
			},
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.context.GetRepsForSet()
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestLookupContext_ApplyModifiers_SetIndexOutOfRange(t *testing.T) {
	weeklyLookup := &weeklylookup.WeeklyLookup{
		Entries: []weeklylookup.WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65, 75, 85}, Reps: []int{5, 5, 5}},
		},
	}

	ctx := &LookupContext{
		WeeklyLookup: weeklyLookup,
		WeekNumber:   1,
		SetNumber:    10, // Out of range
	}

	// When set number is out of range, base percentage should be used
	result := ctx.ApplyModifiers(90)
	if result != 90 {
		t.Errorf("expected 90 (base percentage), got %f", result)
	}
}

func TestLookupContext_ApplyModifiers_DailyZeroModifier(t *testing.T) {
	dailyLookup := &dailylookup.DailyLookup{
		Entries: []dailylookup.DailyLookupEntry{
			{DayIdentifier: "rest", PercentageModifier: 0}, // 0 should not modify
		},
	}

	ctx := &LookupContext{
		DailyLookup: dailyLookup,
		DaySlug:     "rest",
	}

	// When modifier is 0, it should not be applied (no modification)
	result := ctx.ApplyModifiers(85)
	if result != 85 {
		t.Errorf("expected 85 (no modification for 0 modifier), got %f", result)
	}
}
