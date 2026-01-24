package juggernaut

import (
	"testing"
)

func TestCalculateNewTM_ExactRepStandard(t *testing.T) {
	tests := []struct {
		name        string
		currentTM   float64
		waveIndex   int
		amrapReps   int
		isUpperBody bool
		expected    float64
	}{
		{
			name:        "10s wave exact (lower body)",
			currentTM:   200,
			waveIndex:   0,
			amrapReps:   10,
			isUpperBody: false,
			expected:    200, // 10 - 10 = 0, 0 * 5 = 0
		},
		{
			name:        "8s wave exact (lower body)",
			currentTM:   200,
			waveIndex:   1,
			amrapReps:   8,
			isUpperBody: false,
			expected:    200, // 8 - 8 = 0, 0 * 5 = 0
		},
		{
			name:        "5s wave exact (upper body)",
			currentTM:   100,
			waveIndex:   2,
			amrapReps:   5,
			isUpperBody: true,
			expected:    100, // 5 - 5 = 0, 0 * 2.5 = 0
		},
		{
			name:        "3s wave exact (upper body)",
			currentTM:   100,
			waveIndex:   3,
			amrapReps:   3,
			isUpperBody: true,
			expected:    100, // 3 - 3 = 0, 0 * 2.5 = 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateNewTM(tt.currentTM, tt.waveIndex, tt.amrapReps, tt.isUpperBody)
			if result != tt.expected {
				t.Errorf("CalculateNewTM(%v, %d, %d, %v) = %v, want %v",
					tt.currentTM, tt.waveIndex, tt.amrapReps, tt.isUpperBody, result, tt.expected)
			}
		})
	}
}

func TestCalculateNewTM_OverPerformance(t *testing.T) {
	tests := []struct {
		name        string
		currentTM   float64
		waveIndex   int
		amrapReps   int
		isUpperBody bool
		expected    float64
	}{
		{
			name:        "10s wave +3 reps (lower body squat)",
			currentTM:   200,
			waveIndex:   0,
			amrapReps:   13,
			isUpperBody: false,
			expected:    215, // 13 - 10 = 3, 3 * 5 = 15
		},
		{
			name:        "8s wave +2 reps (lower body deadlift)",
			currentTM:   300,
			waveIndex:   1,
			amrapReps:   10,
			isUpperBody: false,
			expected:    310, // 10 - 8 = 2, 2 * 5 = 10
		},
		{
			name:        "5s wave +5 reps (upper body bench)",
			currentTM:   150,
			waveIndex:   2,
			amrapReps:   10,
			isUpperBody: true,
			expected:    162.5, // 10 - 5 = 5, 5 * 2.5 = 12.5
		},
		{
			name:        "3s wave +4 reps (upper body OHP)",
			currentTM:   80,
			waveIndex:   3,
			amrapReps:   7,
			isUpperBody: true,
			expected:    90, // 7 - 3 = 4, 4 * 2.5 = 10
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateNewTM(tt.currentTM, tt.waveIndex, tt.amrapReps, tt.isUpperBody)
			if result != tt.expected {
				t.Errorf("CalculateNewTM(%v, %d, %d, %v) = %v, want %v",
					tt.currentTM, tt.waveIndex, tt.amrapReps, tt.isUpperBody, result, tt.expected)
			}
		})
	}
}

func TestCalculateNewTM_UnderPerformance(t *testing.T) {
	tests := []struct {
		name        string
		currentTM   float64
		waveIndex   int
		amrapReps   int
		isUpperBody bool
		expected    float64
	}{
		{
			name:        "10s wave -2 reps (lower body)",
			currentTM:   200,
			waveIndex:   0,
			amrapReps:   8,
			isUpperBody: false,
			expected:    190, // 8 - 10 = -2, -2 * 5 = -10
		},
		{
			name:        "8s wave -3 reps (lower body)",
			currentTM:   300,
			waveIndex:   1,
			amrapReps:   5,
			isUpperBody: false,
			expected:    285, // 5 - 8 = -3, -3 * 5 = -15
		},
		{
			name:        "5s wave -2 reps (upper body)",
			currentTM:   150,
			waveIndex:   2,
			amrapReps:   3,
			isUpperBody: true,
			expected:    145, // 3 - 5 = -2, -2 * 2.5 = -5
		},
		{
			name:        "3s wave -1 rep (upper body)",
			currentTM:   80,
			waveIndex:   3,
			amrapReps:   2,
			isUpperBody: true,
			expected:    77.5, // 2 - 3 = -1, -1 * 2.5 = -2.5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateNewTM(tt.currentTM, tt.waveIndex, tt.amrapReps, tt.isUpperBody)
			if result != tt.expected {
				t.Errorf("CalculateNewTM(%v, %d, %d, %v) = %v, want %v",
					tt.currentTM, tt.waveIndex, tt.amrapReps, tt.isUpperBody, result, tt.expected)
			}
		})
	}
}

func TestCalculateNewTM_UpperBodyIncrement(t *testing.T) {
	// Verify upper body uses 2.5 increment
	currentTM := 100.0
	waveIndex := 0 // 10s wave, rep standard = 10
	amrapReps := 12
	isUpperBody := true

	result := CalculateNewTM(currentTM, waveIndex, amrapReps, isUpperBody)
	expected := 105.0 // 12 - 10 = 2, 2 * 2.5 = 5

	if result != expected {
		t.Errorf("Upper body increment: got %v, want %v", result, expected)
	}
}

func TestCalculateNewTM_LowerBodyIncrement(t *testing.T) {
	// Verify lower body uses 5.0 increment
	currentTM := 200.0
	waveIndex := 0 // 10s wave, rep standard = 10
	amrapReps := 12
	isUpperBody := false

	result := CalculateNewTM(currentTM, waveIndex, amrapReps, isUpperBody)
	expected := 210.0 // 12 - 10 = 2, 2 * 5 = 10

	if result != expected {
		t.Errorf("Lower body increment: got %v, want %v", result, expected)
	}
}

func TestCalculateNewTM_WaveRepStandards(t *testing.T) {
	// Verify each wave uses correct rep standard (10/8/5/3)
	// All tests use same TM, +1 over standard, lower body
	currentTM := 200.0
	isUpperBody := false

	tests := []struct {
		waveIndex int
		amrapReps int
		expected  float64
	}{
		{0, 11, 205}, // 10s wave: 11-10=1, 1*5=5
		{1, 9, 205},  // 8s wave: 9-8=1, 1*5=5
		{2, 6, 205},  // 5s wave: 6-5=1, 1*5=5
		{3, 4, 205},  // 3s wave: 4-3=1, 1*5=5
	}

	for _, tt := range tests {
		result := CalculateNewTM(currentTM, tt.waveIndex, tt.amrapReps, isUpperBody)
		if result != tt.expected {
			t.Errorf("Wave %d (rep std %d): got %v, want %v",
				tt.waveIndex, RepStandards[tt.waveIndex], result, tt.expected)
		}
	}
}

func TestCalculateNewTM_InvalidWaveIndex(t *testing.T) {
	currentTM := 200.0

	tests := []struct {
		name      string
		waveIndex int
	}{
		{"negative wave index", -1},
		{"wave index too high", 4},
		{"wave index way too high", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateNewTM(currentTM, tt.waveIndex, 10, false)
			if result != currentTM {
				t.Errorf("Invalid wave index %d: got %v, want %v (unchanged)",
					tt.waveIndex, result, currentTM)
			}
		})
	}
}

func TestCalculateCycleIncrement(t *testing.T) {
	tests := []struct {
		name        string
		isUpperBody bool
		expected    float64
	}{
		{"upper body (bench/OHP)", true, 5.0},
		{"lower body (squat/deadlift)", false, 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateCycleIncrement(tt.isUpperBody)
			if result != tt.expected {
				t.Errorf("CalculateCycleIncrement(%v) = %v, want %v",
					tt.isUpperBody, result, tt.expected)
			}
		})
	}
}

func TestRepStandardsConstant(t *testing.T) {
	expected := []int{10, 8, 5, 3}
	if len(RepStandards) != len(expected) {
		t.Fatalf("RepStandards length: got %d, want %d", len(RepStandards), len(expected))
	}
	for i, v := range expected {
		if RepStandards[i] != v {
			t.Errorf("RepStandards[%d] = %d, want %d", i, RepStandards[i], v)
		}
	}
}

func TestIncrementConstants(t *testing.T) {
	if UpperBodyIncrement != 2.5 {
		t.Errorf("UpperBodyIncrement = %v, want 2.5", UpperBodyIncrement)
	}
	if LowerBodyIncrement != 5.0 {
		t.Errorf("LowerBodyIncrement = %v, want 5.0", LowerBodyIncrement)
	}
}
