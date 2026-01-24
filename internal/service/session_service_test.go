package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/waynenilsen/power-pro-v3/internal/domain/loggedset"
	"github.com/waynenilsen/power-pro-v3/internal/domain/prescription"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
	"github.com/waynenilsen/power-pro-v3/internal/service"
)

// mockPrescriptionRepo implements service.PrescriptionRepository for testing.
type mockPrescriptionRepo struct {
	prescriptions map[string]*prescription.Prescription
}

func (m *mockPrescriptionRepo) GetByID(id string) (*prescription.Prescription, error) {
	p, ok := m.prescriptions[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

// mockLoggedSetLister implements service.LoggedSetLister for testing.
type mockLoggedSetLister struct {
	sets map[string][]loggedset.LoggedSet // key: sessionID:prescriptionID
}

func (m *mockLoggedSetLister) ListBySessionAndPrescription(sessionID, prescriptionID string) ([]loggedset.LoggedSet, error) {
	key := sessionID + ":" + prescriptionID
	sets, ok := m.sets[key]
	if !ok {
		return nil, nil
	}
	return sets, nil
}

func TestSessionService_GetNextSet_NotFound(t *testing.T) {
	prescRepo := &mockPrescriptionRepo{prescriptions: make(map[string]*prescription.Prescription)}
	loggedSetLister := &mockLoggedSetLister{sets: make(map[string][]loggedset.LoggedSet)}

	svc := service.NewSessionService(prescRepo, loggedSetLister)

	req := service.NextSetRequest{
		SessionID:      "session-1",
		PrescriptionID: "nonexistent",
		UserID:         "user-1",
	}

	_, err := svc.GetNextSet(context.Background(), req)
	assert.ErrorIs(t, err, service.ErrPrescriptionNotFound)
}

func TestSessionService_GetNextSet_NotVariableScheme(t *testing.T) {
	// Create a fixed scheme prescription (not variable)
	fixed, _ := setscheme.NewFixedSetScheme(5, 5)
	presc := &prescription.Prescription{
		ID:        "presc-1",
		SetScheme: fixed,
	}

	prescRepo := &mockPrescriptionRepo{
		prescriptions: map[string]*prescription.Prescription{"presc-1": presc},
	}
	loggedSetLister := &mockLoggedSetLister{sets: make(map[string][]loggedset.LoggedSet)}

	svc := service.NewSessionService(prescRepo, loggedSetLister)

	req := service.NextSetRequest{
		SessionID:      "session-1",
		PrescriptionID: "presc-1",
		UserID:         "user-1",
	}

	_, err := svc.GetNextSet(context.Background(), req)
	assert.ErrorIs(t, err, service.ErrNotVariableScheme)
}

func TestSessionService_GetNextSet_NoSetsLogged(t *testing.T) {
	// Create an MRS scheme prescription (variable)
	mrs, _ := setscheme.NewMRS(25, 3, 10, 1)
	presc := &prescription.Prescription{
		ID:        "presc-1",
		SetScheme: mrs,
	}

	prescRepo := &mockPrescriptionRepo{
		prescriptions: map[string]*prescription.Prescription{"presc-1": presc},
	}
	loggedSetLister := &mockLoggedSetLister{sets: make(map[string][]loggedset.LoggedSet)}

	svc := service.NewSessionService(prescRepo, loggedSetLister)

	req := service.NextSetRequest{
		SessionID:      "session-1",
		PrescriptionID: "presc-1",
		UserID:         "user-1",
	}

	_, err := svc.GetNextSet(context.Background(), req)
	assert.ErrorIs(t, err, service.ErrNoSetsLogged)
}

func TestSessionService_GetNextSet_MRS_ContinueAfterFirstSet(t *testing.T) {
	// Create an MRS scheme: target 25 total reps, min 3 reps per set
	mrs, _ := setscheme.NewMRS(25, 3, 10, 1)
	presc := &prescription.Prescription{
		ID:        "presc-1",
		SetScheme: mrs,
	}

	prescRepo := &mockPrescriptionRepo{
		prescriptions: map[string]*prescription.Prescription{"presc-1": presc},
	}

	// User has logged 1 set with 10 reps (total: 10, need 25)
	loggedSetLister := &mockLoggedSetLister{
		sets: map[string][]loggedset.LoggedSet{
			"session-1:presc-1": {
				{SetNumber: 1, Weight: 225, TargetReps: 3, RepsPerformed: 10},
			},
		},
	}

	svc := service.NewSessionService(prescRepo, loggedSetLister)

	req := service.NextSetRequest{
		SessionID:      "session-1",
		PrescriptionID: "presc-1",
		UserID:         "user-1",
	}

	result, err := svc.GetNextSet(context.Background(), req)
	require.NoError(t, err)
	assert.False(t, result.IsComplete)
	assert.NotNil(t, result.NextSet)
	assert.Equal(t, 2, result.NextSet.SetNumber)
	assert.Equal(t, 225.0, result.NextSet.Weight) // MRS keeps same weight
	assert.Equal(t, 1, result.TotalSetsCompleted)
	assert.Equal(t, 10, result.TotalRepsCompleted)
}

func TestSessionService_GetNextSet_MRS_TargetReached(t *testing.T) {
	// Create an MRS scheme: target 25 total reps, min 3 reps per set
	mrs, _ := setscheme.NewMRS(25, 3, 10, 1)
	presc := &prescription.Prescription{
		ID:        "presc-1",
		SetScheme: mrs,
	}

	prescRepo := &mockPrescriptionRepo{
		prescriptions: map[string]*prescription.Prescription{"presc-1": presc},
	}

	// User has logged 3 sets with 27 total reps (exceeds 25)
	loggedSetLister := &mockLoggedSetLister{
		sets: map[string][]loggedset.LoggedSet{
			"session-1:presc-1": {
				{SetNumber: 1, Weight: 225, TargetReps: 3, RepsPerformed: 10},
				{SetNumber: 2, Weight: 225, TargetReps: 3, RepsPerformed: 9},
				{SetNumber: 3, Weight: 225, TargetReps: 3, RepsPerformed: 8},
			},
		},
	}

	svc := service.NewSessionService(prescRepo, loggedSetLister)

	req := service.NextSetRequest{
		SessionID:      "session-1",
		PrescriptionID: "presc-1",
		UserID:         "user-1",
	}

	result, err := svc.GetNextSet(context.Background(), req)
	require.NoError(t, err)
	assert.True(t, result.IsComplete)
	assert.Nil(t, result.NextSet)
	assert.Equal(t, 3, result.TotalSetsCompleted)
	assert.Equal(t, 27, result.TotalRepsCompleted)
	assert.Contains(t, result.TerminationReason, "Target total reps reached")
}

func TestSessionService_GetNextSet_MRS_RepFailure(t *testing.T) {
	// Create an MRS scheme: target 25 total reps, min 3 reps per set
	mrs, _ := setscheme.NewMRS(25, 3, 10, 1)
	presc := &prescription.Prescription{
		ID:        "presc-1",
		SetScheme: mrs,
	}

	prescRepo := &mockPrescriptionRepo{
		prescriptions: map[string]*prescription.Prescription{"presc-1": presc},
	}

	// User has logged 3 sets, last one only got 2 reps (below min 3)
	loggedSetLister := &mockLoggedSetLister{
		sets: map[string][]loggedset.LoggedSet{
			"session-1:presc-1": {
				{SetNumber: 1, Weight: 225, TargetReps: 3, RepsPerformed: 8},
				{SetNumber: 2, Weight: 225, TargetReps: 3, RepsPerformed: 6},
				{SetNumber: 3, Weight: 225, TargetReps: 3, RepsPerformed: 2}, // Failed to hit min
			},
		},
	}

	svc := service.NewSessionService(prescRepo, loggedSetLister)

	req := service.NextSetRequest{
		SessionID:      "session-1",
		PrescriptionID: "presc-1",
		UserID:         "user-1",
	}

	result, err := svc.GetNextSet(context.Background(), req)
	require.NoError(t, err)
	assert.True(t, result.IsComplete)
	assert.Nil(t, result.NextSet)
	assert.Contains(t, result.TerminationReason, "Failed to hit minimum reps")
}

func TestSessionService_GetNextSet_FatigueDrop_ContinueAfterFirstSet(t *testing.T) {
	// Create a FatigueDrop scheme: 3 reps, start RPE 8, stop RPE 10, 5% drop
	fd, _ := setscheme.NewFatigueDrop(3, 8.0, 10.0, 0.05, 10)
	presc := &prescription.Prescription{
		ID:        "presc-1",
		SetScheme: fd,
	}

	prescRepo := &mockPrescriptionRepo{
		prescriptions: map[string]*prescription.Prescription{"presc-1": presc},
	}

	rpe8 := 8.0
	// User has logged 1 set at RPE 8 (below stop threshold of 10)
	loggedSetLister := &mockLoggedSetLister{
		sets: map[string][]loggedset.LoggedSet{
			"session-1:presc-1": {
				{SetNumber: 1, Weight: 315, TargetReps: 3, RepsPerformed: 3, RPE: &rpe8},
			},
		},
	}

	svc := service.NewSessionService(prescRepo, loggedSetLister)

	req := service.NextSetRequest{
		SessionID:      "session-1",
		PrescriptionID: "presc-1",
		UserID:         "user-1",
	}

	result, err := svc.GetNextSet(context.Background(), req)
	require.NoError(t, err)
	assert.False(t, result.IsComplete)
	assert.NotNil(t, result.NextSet)
	assert.Equal(t, 2, result.NextSet.SetNumber)
	// Weight should be 315 * 0.95 = 299.25, rounded down to 299 or 297.5
	assert.Less(t, result.NextSet.Weight, 315.0)
	assert.Equal(t, 3, result.NextSet.TargetReps)
}

func TestSessionService_GetNextSet_FatigueDrop_RPEThresholdReached(t *testing.T) {
	// Create a FatigueDrop scheme: 3 reps, start RPE 8, stop RPE 10, 5% drop
	fd, _ := setscheme.NewFatigueDrop(3, 8.0, 10.0, 0.05, 10)
	presc := &prescription.Prescription{
		ID:        "presc-1",
		SetScheme: fd,
	}

	prescRepo := &mockPrescriptionRepo{
		prescriptions: map[string]*prescription.Prescription{"presc-1": presc},
	}

	rpe8 := 8.0
	rpe9 := 9.0
	rpe10 := 10.0
	// User has logged sets until RPE 10
	loggedSetLister := &mockLoggedSetLister{
		sets: map[string][]loggedset.LoggedSet{
			"session-1:presc-1": {
				{SetNumber: 1, Weight: 315, TargetReps: 3, RepsPerformed: 3, RPE: &rpe8},
				{SetNumber: 2, Weight: 299, TargetReps: 3, RepsPerformed: 3, RPE: &rpe9},
				{SetNumber: 3, Weight: 284, TargetReps: 3, RepsPerformed: 3, RPE: &rpe10},
			},
		},
	}

	svc := service.NewSessionService(prescRepo, loggedSetLister)

	req := service.NextSetRequest{
		SessionID:      "session-1",
		PrescriptionID: "presc-1",
		UserID:         "user-1",
	}

	result, err := svc.GetNextSet(context.Background(), req)
	require.NoError(t, err)
	assert.True(t, result.IsComplete)
	assert.Nil(t, result.NextSet)
	assert.Contains(t, result.TerminationReason, "Target RPE reached")
}
