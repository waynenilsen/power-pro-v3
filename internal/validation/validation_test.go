package validation

import (
	"errors"
	"testing"
)

func TestNewResult(t *testing.T) {
	r := NewResult()
	if !r.Valid {
		t.Error("new result should be valid")
	}
	if len(r.Errors) != 0 {
		t.Error("new result should have no errors")
	}
	if len(r.Warnings) != 0 {
		t.Error("new result should have no warnings")
	}
}

func TestResult_AddError(t *testing.T) {
	r := NewResult()
	err := errors.New("test error")
	r.AddError(err)

	if r.Valid {
		t.Error("result should be invalid after adding error")
	}
	if len(r.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(r.Errors))
	}
	if r.Errors[0] != err {
		t.Error("error should be stored")
	}
}

func TestResult_AddError_Multiple(t *testing.T) {
	r := NewResult()
	r.AddError(errors.New("error 1"))
	r.AddError(errors.New("error 2"))
	r.AddError(errors.New("error 3"))

	if r.Valid {
		t.Error("result should be invalid")
	}
	if len(r.Errors) != 3 {
		t.Errorf("expected 3 errors, got %d", len(r.Errors))
	}
}

func TestResult_Error_Valid(t *testing.T) {
	r := NewResult()
	if r.Error() != nil {
		t.Error("valid result should return nil error")
	}
}

func TestResult_Error_Invalid(t *testing.T) {
	r := NewResult()
	r.AddError(errors.New("first error"))
	r.AddError(errors.New("second error"))

	err := r.Error()
	if err == nil {
		t.Fatal("invalid result should return error")
	}

	errMsg := err.Error()
	if errMsg != "validation failed: first error; second error" {
		t.Errorf("unexpected error message: %s", errMsg)
	}
}

func TestResult_Merge(t *testing.T) {
	r1 := NewResult()
	r1.AddError(errors.New("error 1"))

	r2 := NewResult()
	r2.AddError(errors.New("error 2"))
	r2.AddError(errors.New("error 3"))

	r1.Merge(r2)

	if r1.Valid {
		t.Error("merged result should be invalid")
	}
	if len(r1.Errors) != 3 {
		t.Errorf("expected 3 errors after merge, got %d", len(r1.Errors))
	}
}

func TestResult_Merge_Nil(t *testing.T) {
	r := NewResult()
	r.Merge(nil)

	if !r.Valid {
		t.Error("merging nil should not affect validity")
	}
}

func TestResult_Merge_Empty(t *testing.T) {
	r1 := NewResult()
	r2 := NewResult()
	r1.Merge(r2)

	if !r1.Valid {
		t.Error("merging valid result should not affect validity")
	}
}

func TestResult_AddWarning(t *testing.T) {
	r := NewResult()
	r.AddWarning("test warning")

	if !r.Valid {
		t.Error("adding warning should not affect validity")
	}
	if len(r.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(r.Warnings))
	}
	if r.Warnings[0] != "test warning" {
		t.Errorf("expected 'test warning', got %s", r.Warnings[0])
	}
}

func TestResult_HasWarnings(t *testing.T) {
	r := NewResult()
	if r.HasWarnings() {
		t.Error("new result should not have warnings")
	}

	r.AddWarning("warning")
	if !r.HasWarnings() {
		t.Error("result with warning should have warnings")
	}
}

func TestResult_Merge_WithWarnings(t *testing.T) {
	r1 := NewResult()
	r1.AddWarning("warning 1")

	r2 := NewResult()
	r2.AddWarning("warning 2")

	r1.Merge(r2)

	if len(r1.Warnings) != 2 {
		t.Errorf("expected 2 warnings after merge, got %d", len(r1.Warnings))
	}
}
