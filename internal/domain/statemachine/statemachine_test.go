package statemachine

import (
	"testing"
)

func TestInvalidTransitionError(t *testing.T) {
	err := NewInvalidTransitionError("STATE_A", "STATE_B")

	if err.From != "STATE_A" {
		t.Errorf("expected From to be STATE_A, got %s", err.From)
	}

	if err.To != "STATE_B" {
		t.Errorf("expected To to be STATE_B, got %s", err.To)
	}

	expected := "invalid transition from STATE_A to STATE_B"
	if err.Error() != expected {
		t.Errorf("expected error message %q, got %q", expected, err.Error())
	}
}

func TestTransition(t *testing.T) {
	transition := Transition{From: "A", To: "B"}

	if transition.From != "A" {
		t.Errorf("expected From to be A, got %s", transition.From)
	}

	if transition.To != "B" {
		t.Errorf("expected To to be B, got %s", transition.To)
	}
}
