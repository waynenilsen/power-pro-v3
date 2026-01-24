// Package schedule provides domain logic for program scheduling.
// This package contains pure business logic with no database dependencies,
// making it testable in isolation.
package schedule

import (
	"errors"
	"time"
)

// Phase represents a training phase in a peaking program.
type Phase string

const (
	// PhasePrep1 is the first preparatory phase (weeks 1-4 in Sheiko).
	PhasePrep1 Phase = "prep1"
	// PhasePrep2 is the second preparatory phase (weeks 5-8 in Sheiko).
	PhasePrep2 Phase = "prep2"
	// PhaseCompetition is the competition/peaking phase (weeks 9-13 in Sheiko).
	PhaseCompetition Phase = "competition"
)

// PhaseDurations defines the length of each phase in weeks.
type PhaseDurations struct {
	Prep1       int // weeks in Prep1 phase
	Prep2       int // weeks in Prep2 phase
	Competition int // weeks in Competition phase
}

// DefaultPhaseDurations returns the standard Sheiko phase durations (13 weeks total).
// Prep1: 4 weeks, Prep2: 4 weeks, Competition: 5 weeks.
func DefaultPhaseDurations() PhaseDurations {
	return PhaseDurations{
		Prep1:       4,
		Prep2:       4,
		Competition: 5,
	}
}

// TotalWeeks returns the total number of weeks across all phases.
func (p PhaseDurations) TotalWeeks() int {
	return p.Prep1 + p.Prep2 + p.Competition
}

// Validation errors
var (
	ErrMeetDateRequired = errors.New("meet_date is required")
	ErrMeetDateInPast   = errors.New("meet_date must be in the future or today")
	ErrInvalidPhase     = errors.New("invalid phase")
)

// DaysOutResult contains the result of a days out calculation.
type DaysOutResult struct {
	DaysOut          int   // Number of days until meet
	Phase            Phase // Current training phase
	WeekWithinPhase  int   // 1-based week number within the current phase
	WeekOverall      int   // 1-based week number in the overall program
	TotalProgramDays int   // Total days in the program
}

// GetDaysOut calculates the number of days from now until the meet date.
// The meet date is inclusive, so if now is the meet date, daysOut is 0.
// Returns a negative number if the meet date has passed.
func GetDaysOut(meetDate, now time.Time) int {
	// Normalize both dates to start of day (midnight) in their respective timezones
	// to ensure consistent day calculations
	meetDay := time.Date(meetDate.Year(), meetDate.Month(), meetDate.Day(), 0, 0, 0, 0, meetDate.Location())
	nowDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Calculate difference in days
	diff := meetDay.Sub(nowDay)
	return int(diff.Hours() / 24)
}

// GetCurrentPhase determines the training phase based on days remaining until meet.
// Uses standard Sheiko phase durations:
// - Competition: 0-34 days out (weeks 9-13, 5 weeks)
// - Prep2: 35-62 days out (weeks 5-8, 4 weeks)
// - Prep1: 63-90 days out (weeks 1-4, 4 weeks)
// - Before program: 91+ days out (returns Prep1 as default)
func GetCurrentPhase(meetDate, now time.Time) Phase {
	return GetCurrentPhaseWithDurations(meetDate, now, DefaultPhaseDurations())
}

// GetCurrentPhaseWithDurations determines the training phase based on days remaining
// and custom phase durations.
func GetCurrentPhaseWithDurations(meetDate, now time.Time, durations PhaseDurations) Phase {
	daysOut := GetDaysOut(meetDate, now)

	// Competition phase: 0 to (competition_weeks * 7 - 1) days out
	competitionDays := durations.Competition * 7
	if daysOut < competitionDays {
		return PhaseCompetition
	}

	// Prep2 phase: competition_days to (competition_days + prep2_weeks * 7 - 1) days out
	prep2Days := durations.Prep2 * 7
	if daysOut < competitionDays+prep2Days {
		return PhasePrep2
	}

	// Prep1 phase: everything else (including before the program starts)
	return PhasePrep1
}

// GetWeekWithinPhase returns the 1-based week number within the current phase.
// Week 1 is the first week of the phase (furthest from meet for that phase).
func GetWeekWithinPhase(meetDate, now time.Time) int {
	return GetWeekWithinPhaseWithDurations(meetDate, now, DefaultPhaseDurations())
}

// GetWeekWithinPhaseWithDurations returns the 1-based week number within the current
// phase using custom durations.
func GetWeekWithinPhaseWithDurations(meetDate, now time.Time, durations PhaseDurations) int {
	daysOut := GetDaysOut(meetDate, now)
	phase := GetCurrentPhaseWithDurations(meetDate, now, durations)

	competitionDays := durations.Competition * 7
	prep2Days := durations.Prep2 * 7
	prep1Days := durations.Prep1 * 7

	switch phase {
	case PhaseCompetition:
		// Days 0-34 map to weeks 5-1 (reversed: closer to meet = higher week in phase)
		// Actually, week 1 of competition is days 28-34, week 5 is days 0-6
		// We want week 1 to be the START of the phase (furthest from meet within phase)
		// So days 28-34 -> week 1, days 21-27 -> week 2, etc.
		weekFromEnd := daysOut / 7
		return durations.Competition - weekFromEnd

	case PhasePrep2:
		// Prep2 starts at competitionDays and ends at competitionDays + prep2Days - 1
		daysIntoPhase := competitionDays + prep2Days - 1 - daysOut
		return daysIntoPhase/7 + 1

	case PhasePrep1:
		// Prep1 starts at competitionDays + prep2Days and goes back
		// If we're before the program starts (daysOut >= totalDays), we're in "week 1"
		totalDays := (durations.Prep1 + durations.Prep2 + durations.Competition) * 7
		if daysOut >= totalDays {
			return 1
		}
		daysIntoPhase := competitionDays + prep2Days + prep1Days - 1 - daysOut
		week := daysIntoPhase/7 + 1
		if week > durations.Prep1 {
			return durations.Prep1
		}
		if week < 1 {
			return 1
		}
		return week
	}

	return 1
}

// GetOverallWeek returns the 1-based week number in the overall program.
// Week 1 is the first week of Prep1 (furthest from meet).
func GetOverallWeek(meetDate, now time.Time) int {
	return GetOverallWeekWithDurations(meetDate, now, DefaultPhaseDurations())
}

// GetOverallWeekWithDurations returns the 1-based week number in the overall program
// using custom durations.
func GetOverallWeekWithDurations(meetDate, now time.Time, durations PhaseDurations) int {
	daysOut := GetDaysOut(meetDate, now)
	totalDays := durations.TotalWeeks() * 7

	// If before program starts, return week 1
	if daysOut >= totalDays {
		return 1
	}

	// If past meet date (negative days), return the last week
	if daysOut < 0 {
		return durations.TotalWeeks()
	}

	// Calculate week from start of program
	// daysOut = totalDays - 1 is the first day of the program
	daysIntoProgram := totalDays - 1 - daysOut
	week := daysIntoProgram/7 + 1

	// Clamp to valid range
	if week < 1 {
		return 1
	}
	if week > durations.TotalWeeks() {
		return durations.TotalWeeks()
	}

	return week
}

// Calculate computes the complete schedule information for a given date relative to meet.
func Calculate(meetDate, now time.Time) (*DaysOutResult, error) {
	return CalculateWithDurations(meetDate, now, DefaultPhaseDurations())
}

// CalculateWithDurations computes the complete schedule information using custom durations.
func CalculateWithDurations(meetDate, now time.Time, durations PhaseDurations) (*DaysOutResult, error) {
	daysOut := GetDaysOut(meetDate, now)
	phase := GetCurrentPhaseWithDurations(meetDate, now, durations)
	weekWithinPhase := GetWeekWithinPhaseWithDurations(meetDate, now, durations)
	weekOverall := GetOverallWeekWithDurations(meetDate, now, durations)
	totalDays := durations.TotalWeeks() * 7

	return &DaysOutResult{
		DaysOut:          daysOut,
		Phase:            phase,
		WeekWithinPhase:  weekWithinPhase,
		WeekOverall:      weekOverall,
		TotalProgramDays: totalDays,
	}, nil
}

// PhaseInfo returns information about a specific phase.
type PhaseInfo struct {
	Phase         Phase
	StartDaysOut  int // Days out when phase starts (inclusive, furthest from meet)
	EndDaysOut    int // Days out when phase ends (inclusive, closest to meet)
	WeeksInPhase  int
	StartWeek     int // Overall week number when phase starts
	EndWeek       int // Overall week number when phase ends
}

// GetPhaseInfo returns detailed information about each phase.
func GetPhaseInfo(durations PhaseDurations) []PhaseInfo {
	totalWeeks := durations.TotalWeeks()
	competitionDays := durations.Competition * 7
	prep2Days := durations.Prep2 * 7

	return []PhaseInfo{
		{
			Phase:        PhasePrep1,
			StartDaysOut: competitionDays + prep2Days + durations.Prep1*7 - 1,
			EndDaysOut:   competitionDays + prep2Days,
			WeeksInPhase: durations.Prep1,
			StartWeek:    1,
			EndWeek:      durations.Prep1,
		},
		{
			Phase:        PhasePrep2,
			StartDaysOut: competitionDays + prep2Days - 1,
			EndDaysOut:   competitionDays,
			WeeksInPhase: durations.Prep2,
			StartWeek:    durations.Prep1 + 1,
			EndWeek:      durations.Prep1 + durations.Prep2,
		},
		{
			Phase:        PhaseCompetition,
			StartDaysOut: competitionDays - 1,
			EndDaysOut:   0,
			WeeksInPhase: durations.Competition,
			StartWeek:    durations.Prep1 + durations.Prep2 + 1,
			EndWeek:      totalWeeks,
		},
	}
}

// ScheduleType represents how a program's schedule is determined.
type ScheduleType string

const (
	// ScheduleTypeRotation means the program follows a rotating schedule (default).
	ScheduleTypeRotation ScheduleType = "rotation"
	// ScheduleTypeDaysOut means the program schedule is determined by days until meet date.
	ScheduleTypeDaysOut ScheduleType = "days_out"
)

// EffectiveScheduleInput contains all inputs needed to determine the effective week.
type EffectiveScheduleInput struct {
	ScheduleType     ScheduleType   // How schedule is determined
	CurrentWeek      int            // Current week for rotation-based programs
	MeetDate         *time.Time     // Meet date for days_out programs
	Now              time.Time      // Current time (for testability)
	PhaseDurations   PhaseDurations // Phase durations for days_out calculation
	CycleLengthWeeks int            // Total weeks in the program cycle
}

// EffectiveScheduleResult contains the result of effective schedule calculation.
type EffectiveScheduleResult struct {
	WeekNumber      int   // The week number to use for workout generation
	Phase           Phase // Current phase (meaningful for days_out, PhasePrep1 for rotation)
	WeekWithinPhase int   // Week within the current phase
	DaysOut         int   // Days until meet (-1 for rotation schedule)
	IsPeaking       bool  // True if in competition phase (for taper application)
}

// GetEffectiveSchedule determines the effective week and phase based on schedule type.
// For rotation schedules: uses CurrentWeek directly.
// For days_out schedules: derives week from meet date.
func GetEffectiveSchedule(input EffectiveScheduleInput) (*EffectiveScheduleResult, error) {
	switch input.ScheduleType {
	case ScheduleTypeRotation, "":
		// Rotation mode: use CurrentWeek directly, no phase-based logic
		return &EffectiveScheduleResult{
			WeekNumber:      input.CurrentWeek,
			Phase:           PhasePrep1, // Default phase for rotation (not meaningful)
			WeekWithinPhase: input.CurrentWeek,
			DaysOut:         -1, // No meet date
			IsPeaking:       false,
		}, nil

	case ScheduleTypeDaysOut:
		// Days out mode: derive week from meet date
		if input.MeetDate == nil {
			return nil, ErrMeetDateRequired
		}

		// Use provided durations or default
		durations := input.PhaseDurations
		if durations.TotalWeeks() == 0 {
			durations = DefaultPhaseDurations()
		}

		result, err := CalculateWithDurations(*input.MeetDate, input.Now, durations)
		if err != nil {
			return nil, err
		}

		// Clamp week to cycle length if specified
		weekNumber := result.WeekOverall
		if input.CycleLengthWeeks > 0 && weekNumber > input.CycleLengthWeeks {
			weekNumber = input.CycleLengthWeeks
		}

		return &EffectiveScheduleResult{
			WeekNumber:      weekNumber,
			Phase:           result.Phase,
			WeekWithinPhase: result.WeekWithinPhase,
			DaysOut:         result.DaysOut,
			IsPeaking:       result.Phase == PhaseCompetition,
		}, nil

	default:
		return nil, errors.New("invalid schedule type")
	}
}

// ShouldTransitionPhase checks if a phase transition would occur between two points in time.
// This is useful for detecting when the user moves between Prep1 -> Prep2 -> Competition.
// Returns the new phase and whether a transition occurred.
func ShouldTransitionPhase(meetDate time.Time, previousTime, currentTime time.Time) (newPhase Phase, transitioned bool) {
	prevPhase := GetCurrentPhase(meetDate, previousTime)
	currPhase := GetCurrentPhase(meetDate, currentTime)

	return currPhase, prevPhase != currPhase
}

// ShouldTransitionPhaseWithDurations checks for phase transition with custom durations.
func ShouldTransitionPhaseWithDurations(meetDate time.Time, previousTime, currentTime time.Time, durations PhaseDurations) (newPhase Phase, transitioned bool) {
	prevPhase := GetCurrentPhaseWithDurations(meetDate, previousTime, durations)
	currPhase := GetCurrentPhaseWithDurations(meetDate, currentTime, durations)

	return currPhase, prevPhase != currPhase
}

// ValidateMeetDateForSchedule validates the meet date based on schedule requirements.
// For days_out schedule, meet date must be set and in the future (or today).
// For rotation schedule, meet date is optional.
func ValidateMeetDateForSchedule(scheduleType ScheduleType, meetDate *time.Time, now time.Time) error {
	if scheduleType == ScheduleTypeDaysOut {
		if meetDate == nil {
			return ErrMeetDateRequired
		}
		daysOut := GetDaysOut(*meetDate, now)
		if daysOut < 0 {
			return ErrMeetDateInPast
		}
	}
	return nil
}

// HandleMeetDateChange handles the scenario where a meet date is changed mid-program.
// Returns the new effective schedule information based on the new meet date.
func HandleMeetDateChange(oldMeetDate, newMeetDate *time.Time, now time.Time, durations PhaseDurations) (*EffectiveScheduleResult, error) {
	if newMeetDate == nil {
		// Meet date cleared - switch to rotation mode behavior
		return &EffectiveScheduleResult{
			WeekNumber:      1, // Start from week 1
			Phase:           PhasePrep1,
			WeekWithinPhase: 1,
			DaysOut:         -1,
			IsPeaking:       false,
		}, nil
	}

	// Calculate new position based on new meet date
	result, err := CalculateWithDurations(*newMeetDate, now, durations)
	if err != nil {
		return nil, err
	}

	return &EffectiveScheduleResult{
		WeekNumber:      result.WeekOverall,
		Phase:           result.Phase,
		WeekWithinPhase: result.WeekWithinPhase,
		DaysOut:         result.DaysOut,
		IsPeaking:       result.Phase == PhaseCompetition,
	}, nil
}
