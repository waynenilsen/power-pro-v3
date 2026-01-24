package juggernaut

import (
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/domain/weeklylookup"
)

// Create531WeeklyLookup creates the standard 5/3/1 weekly lookup.
// This defines the 5/3/1 main work sets that are consistent across all waves.
// The 4-week pattern repeats for all 4 waves (weeks 1-4 = same as 5-8 = 9-12 = 13-16).
// When resolving, use WeekInWave (1-4) rather than CurrentWeek (1-16) to look up the correct entry.
func Create531WeeklyLookup(id string, programID *string) *weeklylookup.WeeklyLookup {
	now := time.Now()

	return &weeklylookup.WeeklyLookup{
		ID:        id,
		Name:      "Inverted Juggernaut 5/3/1",
		ProgramID: programID,
		Entries: []weeklylookup.WeeklyLookupEntry{
			// Week 1 (Accumulation): 65/75/85/75/65 @ 5/5/5+/5/5+
			{
				WeekNumber:  1,
				Percentages: []float64{65.0, 75.0, 85.0, 75.0, 65.0},
				Reps:        []int{5, 5, -5, 5, -5}, // Negative reps indicate AMRAP with minimum
			},
			// Week 2 (Intensification): 70/80/90/80/70 @ 3/3/3+/3/3+
			{
				WeekNumber:  2,
				Percentages: []float64{70.0, 80.0, 90.0, 80.0, 70.0},
				Reps:        []int{3, 3, -3, 3, -3},
			},
			// Week 3 (Realization): 75/85/95/85/75 @ 5/3/1+/3/5+
			{
				WeekNumber:  3,
				Percentages: []float64{75.0, 85.0, 95.0, 85.0, 75.0},
				Reps:        []int{5, 3, -1, 3, -5},
			},
			// Week 4 (Deload): 40/50/60 @ 5/5/5 (only 3 sets)
			{
				WeekNumber:  4,
				Percentages: []float64{40.0, 50.0, 60.0},
				Reps:        []int{5, 5, 5},
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}
