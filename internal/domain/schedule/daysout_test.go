package schedule

import (
	"testing"
	"time"
)

// Helper function to create a date
func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// ==================== GetDaysOut Tests ====================

func TestGetDaysOut_MeetDateIsToday(t *testing.T) {
	now := date(2024, time.March, 15)
	meetDate := date(2024, time.March, 15)

	daysOut := GetDaysOut(meetDate, now)

	if daysOut != 0 {
		t.Errorf("GetDaysOut() = %d, want 0 for meet date today", daysOut)
	}
}

func TestGetDaysOut_MeetDateTomorrow(t *testing.T) {
	now := date(2024, time.March, 15)
	meetDate := date(2024, time.March, 16)

	daysOut := GetDaysOut(meetDate, now)

	if daysOut != 1 {
		t.Errorf("GetDaysOut() = %d, want 1 for meet date tomorrow", daysOut)
	}
}

func TestGetDaysOut_MeetDateInPast(t *testing.T) {
	now := date(2024, time.March, 15)
	meetDate := date(2024, time.March, 14)

	daysOut := GetDaysOut(meetDate, now)

	if daysOut != -1 {
		t.Errorf("GetDaysOut() = %d, want -1 for meet date yesterday", daysOut)
	}
}

func TestGetDaysOut_90DaysOut(t *testing.T) {
	meetDate := date(2024, time.June, 15)
	now := date(2024, time.March, 17) // 90 days before June 15

	daysOut := GetDaysOut(meetDate, now)

	if daysOut != 90 {
		t.Errorf("GetDaysOut() = %d, want 90", daysOut)
	}
}

func TestGetDaysOut_91DaysOut(t *testing.T) {
	meetDate := date(2024, time.June, 15)
	now := date(2024, time.March, 16) // 91 days before June 15

	daysOut := GetDaysOut(meetDate, now)

	if daysOut != 91 {
		t.Errorf("GetDaysOut() = %d, want 91", daysOut)
	}
}

func TestGetDaysOut_ExactBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		meetDate time.Time
		now      time.Time
		want     int
	}{
		{"0 days", date(2024, 3, 15), date(2024, 3, 15), 0},
		{"1 day", date(2024, 3, 16), date(2024, 3, 15), 1},
		{"7 days", date(2024, 3, 22), date(2024, 3, 15), 7},
		{"35 days", date(2024, 4, 19), date(2024, 3, 15), 35},
		{"63 days", date(2024, 5, 17), date(2024, 3, 15), 63},
		{"91 days", date(2024, 6, 14), date(2024, 3, 15), 91},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDaysOut(tt.meetDate, tt.now)
			if got != tt.want {
				t.Errorf("GetDaysOut() = %d, want %d", got, tt.want)
			}
		})
	}
}

// ==================== GetCurrentPhase Tests ====================

func TestGetCurrentPhase_CompetitionPhase(t *testing.T) {
	// Competition phase: 0-34 days out (5 weeks)
	tests := []struct {
		name    string
		daysOut int
	}{
		{"meet day (0 days)", 0},
		{"1 day out", 1},
		{"7 days out (1 week)", 7},
		{"14 days out (2 weeks)", 14},
		{"21 days out (3 weeks)", 21},
		{"28 days out (4 weeks)", 28},
		{"34 days out (last day of comp)", 34},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meetDate := date(2024, time.June, 15)
			now := meetDate.AddDate(0, 0, -tt.daysOut)

			phase := GetCurrentPhase(meetDate, now)

			if phase != PhaseCompetition {
				t.Errorf("GetCurrentPhase() = %q at %d days out, want %q", phase, tt.daysOut, PhaseCompetition)
			}
		})
	}
}

func TestGetCurrentPhase_Prep2Phase(t *testing.T) {
	// Prep2 phase: 35-62 days out (4 weeks)
	tests := []struct {
		name    string
		daysOut int
	}{
		{"35 days out (first day of prep2)", 35},
		{"42 days out (1 week into prep2)", 42},
		{"49 days out (2 weeks into prep2)", 49},
		{"56 days out (3 weeks into prep2)", 56},
		{"62 days out (last day of prep2)", 62},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meetDate := date(2024, time.June, 15)
			now := meetDate.AddDate(0, 0, -tt.daysOut)

			phase := GetCurrentPhase(meetDate, now)

			if phase != PhasePrep2 {
				t.Errorf("GetCurrentPhase() = %q at %d days out, want %q", phase, tt.daysOut, PhasePrep2)
			}
		})
	}
}

func TestGetCurrentPhase_Prep1Phase(t *testing.T) {
	// Prep1 phase: 63-90 days out (4 weeks)
	tests := []struct {
		name    string
		daysOut int
	}{
		{"63 days out (first day of prep1)", 63},
		{"70 days out (1 week into prep1)", 70},
		{"77 days out (2 weeks into prep1)", 77},
		{"84 days out (3 weeks into prep1)", 84},
		{"90 days out (last day of prep1)", 90},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meetDate := date(2024, time.June, 15)
			now := meetDate.AddDate(0, 0, -tt.daysOut)

			phase := GetCurrentPhase(meetDate, now)

			if phase != PhasePrep1 {
				t.Errorf("GetCurrentPhase() = %q at %d days out, want %q", phase, tt.daysOut, PhasePrep1)
			}
		})
	}
}

func TestGetCurrentPhase_BeforeProgram(t *testing.T) {
	// Before program starts (91+ days out) should default to Prep1
	meetDate := date(2024, time.June, 15)
	now := meetDate.AddDate(0, 0, -100)

	phase := GetCurrentPhase(meetDate, now)

	if phase != PhasePrep1 {
		t.Errorf("GetCurrentPhase() = %q before program, want %q", phase, PhasePrep1)
	}
}

func TestGetCurrentPhase_AfterMeet(t *testing.T) {
	// After meet day (negative days out)
	meetDate := date(2024, time.June, 15)
	now := meetDate.AddDate(0, 0, 5) // 5 days after meet

	phase := GetCurrentPhase(meetDate, now)

	if phase != PhaseCompetition {
		t.Errorf("GetCurrentPhase() = %q after meet, want %q", phase, PhaseCompetition)
	}
}

// ==================== GetWeekWithinPhase Tests ====================

func TestGetWeekWithinPhase_CompetitionPhase(t *testing.T) {
	// Competition phase is 5 weeks (35 days)
	// Week 1: days 28-34, Week 2: days 21-27, Week 3: days 14-20, Week 4: days 7-13, Week 5: days 0-6
	tests := []struct {
		name    string
		daysOut int
		want    int
	}{
		{"day 0 (meet day)", 0, 5},
		{"day 6", 6, 5},
		{"day 7", 7, 4},
		{"day 13", 13, 4},
		{"day 14", 14, 3},
		{"day 20", 20, 3},
		{"day 21", 21, 2},
		{"day 27", 27, 2},
		{"day 28", 28, 1},
		{"day 34", 34, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meetDate := date(2024, time.June, 15)
			now := meetDate.AddDate(0, 0, -tt.daysOut)

			week := GetWeekWithinPhase(meetDate, now)

			if week != tt.want {
				t.Errorf("GetWeekWithinPhase() = %d at %d days out, want %d", week, tt.daysOut, tt.want)
			}
		})
	}
}

func TestGetWeekWithinPhase_Prep2Phase(t *testing.T) {
	// Prep2 phase is 4 weeks (28 days): days 35-62
	// Week 1: days 56-62, Week 2: days 49-55, Week 3: days 42-48, Week 4: days 35-41
	tests := []struct {
		name    string
		daysOut int
		want    int
	}{
		{"day 35 (last day of prep2 week 4)", 35, 4},
		{"day 41", 41, 4},
		{"day 42 (prep2 week 3)", 42, 3},
		{"day 48", 48, 3},
		{"day 49 (prep2 week 2)", 49, 2},
		{"day 55", 55, 2},
		{"day 56 (prep2 week 1)", 56, 1},
		{"day 62 (first day of prep2)", 62, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meetDate := date(2024, time.June, 15)
			now := meetDate.AddDate(0, 0, -tt.daysOut)

			week := GetWeekWithinPhase(meetDate, now)

			if week != tt.want {
				t.Errorf("GetWeekWithinPhase() = %d at %d days out, want %d", week, tt.daysOut, tt.want)
			}
		})
	}
}

func TestGetWeekWithinPhase_Prep1Phase(t *testing.T) {
	// Prep1 phase is 4 weeks (28 days): days 63-90
	// Week 1: days 84-90, Week 2: days 77-83, Week 3: days 70-76, Week 4: days 63-69
	tests := []struct {
		name    string
		daysOut int
		want    int
	}{
		{"day 63 (last day of prep1 week 4)", 63, 4},
		{"day 69", 69, 4},
		{"day 70 (prep1 week 3)", 70, 3},
		{"day 76", 76, 3},
		{"day 77 (prep1 week 2)", 77, 2},
		{"day 83", 83, 2},
		{"day 84 (prep1 week 1)", 84, 1},
		{"day 90 (first day of prep1)", 90, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meetDate := date(2024, time.June, 15)
			now := meetDate.AddDate(0, 0, -tt.daysOut)

			week := GetWeekWithinPhase(meetDate, now)

			if week != tt.want {
				t.Errorf("GetWeekWithinPhase() = %d at %d days out, want %d", week, tt.daysOut, tt.want)
			}
		})
	}
}

// ==================== GetOverallWeek Tests ====================

func TestGetOverallWeek_AllWeeks(t *testing.T) {
	// Total program is 13 weeks (91 days)
	// Week 1: days 84-90, Week 2: days 77-83, ..., Week 13: days 0-6
	tests := []struct {
		name    string
		daysOut int
		want    int
	}{
		{"day 90 (first day, week 1)", 90, 1},
		{"day 84 (end of week 1)", 84, 1},
		{"day 83 (start of week 2)", 83, 2},
		{"day 77 (end of week 2)", 77, 2},
		{"day 76 (start of week 3)", 76, 3},
		{"day 70 (end of week 3)", 70, 3},
		{"day 69 (start of week 4)", 69, 4},
		{"day 63 (end of week 4)", 63, 4},
		{"day 62 (start of week 5)", 62, 5},
		{"day 56 (end of week 5)", 56, 5},
		{"day 55 (start of week 6)", 55, 6},
		{"day 49 (end of week 6)", 49, 6},
		{"day 48 (start of week 7)", 48, 7},
		{"day 42 (end of week 7)", 42, 7},
		{"day 41 (start of week 8)", 41, 8},
		{"day 35 (end of week 8)", 35, 8},
		{"day 34 (start of week 9)", 34, 9},
		{"day 28 (end of week 9)", 28, 9},
		{"day 27 (start of week 10)", 27, 10},
		{"day 21 (end of week 10)", 21, 10},
		{"day 20 (start of week 11)", 20, 11},
		{"day 14 (end of week 11)", 14, 11},
		{"day 13 (start of week 12)", 13, 12},
		{"day 7 (end of week 12)", 7, 12},
		{"day 6 (start of week 13)", 6, 13},
		{"day 0 (meet day, week 13)", 0, 13},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meetDate := date(2024, time.June, 15)
			now := meetDate.AddDate(0, 0, -tt.daysOut)

			week := GetOverallWeek(meetDate, now)

			if week != tt.want {
				t.Errorf("GetOverallWeek() = %d at %d days out, want %d", week, tt.daysOut, tt.want)
			}
		})
	}
}

func TestGetOverallWeek_BeforeProgram(t *testing.T) {
	// Before program starts should return week 1
	meetDate := date(2024, time.June, 15)
	now := meetDate.AddDate(0, 0, -100) // 100 days out

	week := GetOverallWeek(meetDate, now)

	if week != 1 {
		t.Errorf("GetOverallWeek() = %d before program, want 1", week)
	}
}

func TestGetOverallWeek_AfterMeet(t *testing.T) {
	// After meet should return week 13 (last week)
	meetDate := date(2024, time.June, 15)
	now := meetDate.AddDate(0, 0, 5) // 5 days after meet

	week := GetOverallWeek(meetDate, now)

	if week != 13 {
		t.Errorf("GetOverallWeek() = %d after meet, want 13", week)
	}
}

// ==================== Calculate Tests ====================

func TestCalculate_MeetDay(t *testing.T) {
	meetDate := date(2024, time.June, 15)
	now := meetDate

	result, err := Calculate(meetDate, now)

	if err != nil {
		t.Fatalf("Calculate() returned error: %v", err)
	}

	if result.DaysOut != 0 {
		t.Errorf("DaysOut = %d, want 0", result.DaysOut)
	}
	if result.Phase != PhaseCompetition {
		t.Errorf("Phase = %q, want %q", result.Phase, PhaseCompetition)
	}
	if result.WeekWithinPhase != 5 {
		t.Errorf("WeekWithinPhase = %d, want 5", result.WeekWithinPhase)
	}
	if result.WeekOverall != 13 {
		t.Errorf("WeekOverall = %d, want 13", result.WeekOverall)
	}
	if result.TotalProgramDays != 91 {
		t.Errorf("TotalProgramDays = %d, want 91", result.TotalProgramDays)
	}
}

func TestCalculate_FirstDayOfProgram(t *testing.T) {
	meetDate := date(2024, time.June, 15)
	now := meetDate.AddDate(0, 0, -90) // 90 days out = first day of program

	result, err := Calculate(meetDate, now)

	if err != nil {
		t.Fatalf("Calculate() returned error: %v", err)
	}

	if result.DaysOut != 90 {
		t.Errorf("DaysOut = %d, want 90", result.DaysOut)
	}
	if result.Phase != PhasePrep1 {
		t.Errorf("Phase = %q, want %q", result.Phase, PhasePrep1)
	}
	if result.WeekWithinPhase != 1 {
		t.Errorf("WeekWithinPhase = %d, want 1", result.WeekWithinPhase)
	}
	if result.WeekOverall != 1 {
		t.Errorf("WeekOverall = %d, want 1", result.WeekOverall)
	}
}

func TestCalculate_MiddleOfPrep2(t *testing.T) {
	meetDate := date(2024, time.June, 15)
	now := meetDate.AddDate(0, 0, -49) // 49 days out = middle of prep2

	result, err := Calculate(meetDate, now)

	if err != nil {
		t.Fatalf("Calculate() returned error: %v", err)
	}

	if result.DaysOut != 49 {
		t.Errorf("DaysOut = %d, want 49", result.DaysOut)
	}
	if result.Phase != PhasePrep2 {
		t.Errorf("Phase = %q, want %q", result.Phase, PhasePrep2)
	}
	// Week 6 overall, week 2 within prep2
	if result.WeekOverall != 6 {
		t.Errorf("WeekOverall = %d, want 6", result.WeekOverall)
	}
}

// ==================== DefaultPhaseDurations Tests ====================

func TestDefaultPhaseDurations(t *testing.T) {
	durations := DefaultPhaseDurations()

	if durations.Prep1 != 4 {
		t.Errorf("Prep1 = %d, want 4", durations.Prep1)
	}
	if durations.Prep2 != 4 {
		t.Errorf("Prep2 = %d, want 4", durations.Prep2)
	}
	if durations.Competition != 5 {
		t.Errorf("Competition = %d, want 5", durations.Competition)
	}
}

func TestPhaseDurations_TotalWeeks(t *testing.T) {
	durations := DefaultPhaseDurations()

	if durations.TotalWeeks() != 13 {
		t.Errorf("TotalWeeks() = %d, want 13", durations.TotalWeeks())
	}
}

// ==================== Custom Durations Tests ====================

func TestGetCurrentPhaseWithDurations_CustomDurations(t *testing.T) {
	// Custom 10-week program: prep1=2, prep2=3, comp=5
	durations := PhaseDurations{Prep1: 2, Prep2: 3, Competition: 5}

	meetDate := date(2024, time.June, 15)

	tests := []struct {
		name    string
		daysOut int
		want    Phase
	}{
		{"day 0 (competition)", 0, PhaseCompetition},
		{"day 34 (competition boundary)", 34, PhaseCompetition},
		{"day 35 (prep2 start)", 35, PhasePrep2},
		{"day 55 (prep2 boundary)", 55, PhasePrep2},
		{"day 56 (prep1 start)", 56, PhasePrep1},
		{"day 69 (prep1 boundary)", 69, PhasePrep1},
		{"day 70+ (before program)", 70, PhasePrep1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := meetDate.AddDate(0, 0, -tt.daysOut)
			phase := GetCurrentPhaseWithDurations(meetDate, now, durations)

			if phase != tt.want {
				t.Errorf("GetCurrentPhaseWithDurations() = %q at %d days out, want %q", phase, tt.daysOut, tt.want)
			}
		})
	}
}

// ==================== GetPhaseInfo Tests ====================

func TestGetPhaseInfo_DefaultDurations(t *testing.T) {
	durations := DefaultPhaseDurations()
	infos := GetPhaseInfo(durations)

	if len(infos) != 3 {
		t.Fatalf("GetPhaseInfo() returned %d phases, want 3", len(infos))
	}

	// Prep1
	if infos[0].Phase != PhasePrep1 {
		t.Errorf("infos[0].Phase = %q, want %q", infos[0].Phase, PhasePrep1)
	}
	if infos[0].StartWeek != 1 || infos[0].EndWeek != 4 {
		t.Errorf("Prep1 weeks = %d-%d, want 1-4", infos[0].StartWeek, infos[0].EndWeek)
	}

	// Prep2
	if infos[1].Phase != PhasePrep2 {
		t.Errorf("infos[1].Phase = %q, want %q", infos[1].Phase, PhasePrep2)
	}
	if infos[1].StartWeek != 5 || infos[1].EndWeek != 8 {
		t.Errorf("Prep2 weeks = %d-%d, want 5-8", infos[1].StartWeek, infos[1].EndWeek)
	}

	// Competition
	if infos[2].Phase != PhaseCompetition {
		t.Errorf("infos[2].Phase = %q, want %q", infos[2].Phase, PhaseCompetition)
	}
	if infos[2].StartWeek != 9 || infos[2].EndWeek != 13 {
		t.Errorf("Competition weeks = %d-%d, want 9-13", infos[2].StartWeek, infos[2].EndWeek)
	}
}

// ==================== Edge Cases ====================

func TestEdgeCases_LeapYear(t *testing.T) {
	// Test across a leap year boundary
	meetDate := date(2024, time.March, 1) // 2024 is a leap year
	now := date(2024, time.February, 29)

	daysOut := GetDaysOut(meetDate, now)

	if daysOut != 1 {
		t.Errorf("GetDaysOut() across leap year = %d, want 1", daysOut)
	}
}

func TestEdgeCases_YearBoundary(t *testing.T) {
	// Test across a year boundary
	meetDate := date(2025, time.January, 1)
	now := date(2024, time.December, 31)

	daysOut := GetDaysOut(meetDate, now)

	if daysOut != 1 {
		t.Errorf("GetDaysOut() across year boundary = %d, want 1", daysOut)
	}
}

func TestEdgeCases_ExactPhaseBoundaries(t *testing.T) {
	meetDate := date(2024, time.June, 15)

	// Test exact phase boundaries
	tests := []struct {
		name    string
		daysOut int
		want    Phase
	}{
		// Competition phase: 0-34 days (5 weeks * 7 days = 35 days, so 0-34)
		{"day 34 (last day of competition)", 34, PhaseCompetition},
		// Prep2 phase: 35-62 days (4 weeks * 7 = 28 days, so 35 to 35+27=62)
		{"day 35 (first day of prep2)", 35, PhasePrep2},
		{"day 62 (last day of prep2)", 62, PhasePrep2},
		// Prep1 phase: 63-90 days (4 weeks * 7 = 28 days, so 63 to 63+27=90)
		{"day 63 (first day of prep1)", 63, PhasePrep1},
		{"day 90 (last day of prep1)", 90, PhasePrep1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := meetDate.AddDate(0, 0, -tt.daysOut)
			phase := GetCurrentPhase(meetDate, now)

			if phase != tt.want {
				t.Errorf("GetCurrentPhase() = %q at %d days out, want %q", phase, tt.daysOut, tt.want)
			}
		})
	}
}
