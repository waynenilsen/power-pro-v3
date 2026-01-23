// Package repository provides database repository implementations.
package repository

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
)

// setupTestDB creates a temporary SQLite database with the required schema for testing.
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	// Create temp file for SQLite
	tmpFile, err := os.CreateTemp("", "failure_counter_test_*.db")
	if err != nil {
		t.Fatalf("failed to create temp db: %v", err)
	}
	dbPath := tmpFile.Name()
	tmpFile.Close()

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		os.Remove(dbPath)
		t.Fatalf("failed to open db: %v", err)
	}

	// Create required tables
	schema := `
		CREATE TABLE users (
			id TEXT PRIMARY KEY
		);

		CREATE TABLE lifts (
			id TEXT PRIMARY KEY
		);

		CREATE TABLE progressions (
			id TEXT PRIMARY KEY
		);

		CREATE TABLE failure_counters (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			lift_id TEXT NOT NULL,
			progression_id TEXT NOT NULL,
			consecutive_failures INT NOT NULL DEFAULT 0,
			last_failure_at TEXT,
			last_success_at TEXT,
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT NOT NULL DEFAULT (datetime('now')),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (lift_id) REFERENCES lifts(id) ON DELETE CASCADE,
			FOREIGN KEY (progression_id) REFERENCES progressions(id) ON DELETE CASCADE,
			UNIQUE(user_id, lift_id, progression_id)
		);
	`

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		os.Remove(dbPath)
		t.Fatalf("failed to create schema: %v", err)
	}

	// Insert test data for foreign keys
	_, err = db.Exec(`
		INSERT INTO users (id) VALUES ('user-1'), ('user-2');
		INSERT INTO lifts (id) VALUES ('lift-1'), ('lift-2');
		INSERT INTO progressions (id) VALUES ('prog-1'), ('prog-2');
	`)
	if err != nil {
		db.Close()
		os.Remove(dbPath)
		t.Fatalf("failed to insert test data: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.Remove(dbPath)
	}

	return db, cleanup
}

func TestFailureCounterRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewFailureCounterRepository(db)

	now := time.Now()
	fc := &progression.FailureCounter{
		ID:                  "fc-001",
		UserID:              "user-1",
		LiftID:              "lift-1",
		ProgressionID:       "prog-1",
		ConsecutiveFailures: 0,
		LastFailureAt:       nil,
		LastSuccessAt:       nil,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	err := repo.Create(fc)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Verify it was created
	retrieved, err := repo.GetByID("fc-001")
	if err != nil {
		t.Fatalf("GetByID() failed: %v", err)
	}
	if retrieved == nil {
		t.Fatal("expected counter to be found")
	}
	if retrieved.UserID != "user-1" {
		t.Errorf("UserID = %s, want user-1", retrieved.UserID)
	}
}

func TestFailureCounterRepository_GetByKey(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewFailureCounterRepository(db)

	now := time.Now()
	fc := &progression.FailureCounter{
		ID:                  "fc-001",
		UserID:              "user-1",
		LiftID:              "lift-1",
		ProgressionID:       "prog-1",
		ConsecutiveFailures: 3,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := repo.Create(fc); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Test GetByKey
	retrieved, err := repo.GetByKey("user-1", "lift-1", "prog-1")
	if err != nil {
		t.Fatalf("GetByKey() failed: %v", err)
	}
	if retrieved == nil {
		t.Fatal("expected counter to be found")
	}
	if retrieved.ConsecutiveFailures != 3 {
		t.Errorf("ConsecutiveFailures = %d, want 3", retrieved.ConsecutiveFailures)
	}

	// Test GetByKey with non-existent key
	notFound, err := repo.GetByKey("user-1", "lift-1", "prog-2")
	if err != nil {
		t.Fatalf("GetByKey() failed: %v", err)
	}
	if notFound != nil {
		t.Error("expected nil for non-existent key")
	}
}

func TestFailureCounterRepository_IncrementOnFailure(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewFailureCounterRepository(db)

	// Test upsert on non-existent counter (should create with count=1)
	count, err := repo.IncrementOnFailure("user-1", "lift-1", "prog-1")
	if err != nil {
		t.Fatalf("IncrementOnFailure() failed: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}

	// Verify the counter was created
	fc, err := repo.GetByKey("user-1", "lift-1", "prog-1")
	if err != nil {
		t.Fatalf("GetByKey() failed: %v", err)
	}
	if fc == nil {
		t.Fatal("expected counter to be created")
	}
	if fc.ConsecutiveFailures != 1 {
		t.Errorf("ConsecutiveFailures = %d, want 1", fc.ConsecutiveFailures)
	}
	if fc.LastFailureAt == nil {
		t.Error("expected LastFailureAt to be set")
	}

	// Increment again
	count, err = repo.IncrementOnFailure("user-1", "lift-1", "prog-1")
	if err != nil {
		t.Fatalf("IncrementOnFailure() failed: %v", err)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}

	// Third increment
	count, err = repo.IncrementOnFailure("user-1", "lift-1", "prog-1")
	if err != nil {
		t.Fatalf("IncrementOnFailure() failed: %v", err)
	}
	if count != 3 {
		t.Errorf("count = %d, want 3", count)
	}
}

func TestFailureCounterRepository_ResetOnSuccess(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewFailureCounterRepository(db)

	// First create a counter with some failures
	count, err := repo.IncrementOnFailure("user-1", "lift-1", "prog-1")
	if err != nil {
		t.Fatalf("IncrementOnFailure() failed: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}

	// Increment twice more
	repo.IncrementOnFailure("user-1", "lift-1", "prog-1")
	repo.IncrementOnFailure("user-1", "lift-1", "prog-1")

	// Verify we have 3 failures
	fc, _ := repo.GetByKey("user-1", "lift-1", "prog-1")
	if fc.ConsecutiveFailures != 3 {
		t.Errorf("before reset: ConsecutiveFailures = %d, want 3", fc.ConsecutiveFailures)
	}

	// Reset on success
	err = repo.ResetOnSuccess("user-1", "lift-1", "prog-1")
	if err != nil {
		t.Fatalf("ResetOnSuccess() failed: %v", err)
	}

	// Verify the counter was reset
	fc, err = repo.GetByKey("user-1", "lift-1", "prog-1")
	if err != nil {
		t.Fatalf("GetByKey() failed: %v", err)
	}
	if fc.ConsecutiveFailures != 0 {
		t.Errorf("ConsecutiveFailures = %d, want 0", fc.ConsecutiveFailures)
	}
	if fc.LastSuccessAt == nil {
		t.Error("expected LastSuccessAt to be set")
	}
}

func TestFailureCounterRepository_GetFailureCount(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewFailureCounterRepository(db)

	// Test with non-existent counter (should return 0)
	count, err := repo.GetFailureCount("user-1", "lift-1", "prog-1")
	if err != nil {
		t.Fatalf("GetFailureCount() failed: %v", err)
	}
	if count != 0 {
		t.Errorf("count = %d, want 0", count)
	}

	// Create some failures
	repo.IncrementOnFailure("user-1", "lift-1", "prog-1")
	repo.IncrementOnFailure("user-1", "lift-1", "prog-1")

	// Get the count
	count, err = repo.GetFailureCount("user-1", "lift-1", "prog-1")
	if err != nil {
		t.Fatalf("GetFailureCount() failed: %v", err)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
}

func TestFailureCounterRepository_ListByUser(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewFailureCounterRepository(db)

	// Create multiple counters for user-1
	repo.IncrementOnFailure("user-1", "lift-1", "prog-1")
	repo.IncrementOnFailure("user-1", "lift-2", "prog-1")
	repo.IncrementOnFailure("user-1", "lift-1", "prog-2")

	// Create a counter for user-2
	repo.IncrementOnFailure("user-2", "lift-1", "prog-1")

	// List for user-1
	counters, err := repo.ListByUser("user-1")
	if err != nil {
		t.Fatalf("ListByUser() failed: %v", err)
	}
	if len(counters) != 3 {
		t.Errorf("got %d counters, want 3", len(counters))
	}

	// List for user-2
	counters, err = repo.ListByUser("user-2")
	if err != nil {
		t.Fatalf("ListByUser() failed: %v", err)
	}
	if len(counters) != 1 {
		t.Errorf("got %d counters, want 1", len(counters))
	}
}

func TestFailureCounterRepository_ListByUserAndLift(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewFailureCounterRepository(db)

	// Create counters
	repo.IncrementOnFailure("user-1", "lift-1", "prog-1")
	repo.IncrementOnFailure("user-1", "lift-1", "prog-2")
	repo.IncrementOnFailure("user-1", "lift-2", "prog-1")

	// List for user-1, lift-1
	counters, err := repo.ListByUserAndLift("user-1", "lift-1")
	if err != nil {
		t.Fatalf("ListByUserAndLift() failed: %v", err)
	}
	if len(counters) != 2 {
		t.Errorf("got %d counters, want 2", len(counters))
	}

	// List for user-1, lift-2
	counters, err = repo.ListByUserAndLift("user-1", "lift-2")
	if err != nil {
		t.Fatalf("ListByUserAndLift() failed: %v", err)
	}
	if len(counters) != 1 {
		t.Errorf("got %d counters, want 1", len(counters))
	}
}

func TestFailureCounterRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewFailureCounterRepository(db)

	// Create a counter
	now := time.Now()
	fc := &progression.FailureCounter{
		ID:            "fc-001",
		UserID:        "user-1",
		LiftID:        "lift-1",
		ProgressionID: "prog-1",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	repo.Create(fc)

	// Delete it
	err := repo.Delete("fc-001")
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	// Verify it's gone
	retrieved, _ := repo.GetByID("fc-001")
	if retrieved != nil {
		t.Error("expected counter to be deleted")
	}
}

func TestFailureCounterRepository_DeleteByKey(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewFailureCounterRepository(db)

	// Create a counter via increment
	repo.IncrementOnFailure("user-1", "lift-1", "prog-1")

	// Verify it exists
	fc, _ := repo.GetByKey("user-1", "lift-1", "prog-1")
	if fc == nil {
		t.Fatal("expected counter to exist before delete")
	}

	// Delete by key
	err := repo.DeleteByKey("user-1", "lift-1", "prog-1")
	if err != nil {
		t.Fatalf("DeleteByKey() failed: %v", err)
	}

	// Verify it's gone
	fc, _ = repo.GetByKey("user-1", "lift-1", "prog-1")
	if fc != nil {
		t.Error("expected counter to be deleted")
	}
}

func TestFailureCounterRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewFailureCounterRepository(db)

	// Create a counter
	now := time.Now()
	fc := &progression.FailureCounter{
		ID:                  "fc-001",
		UserID:              "user-1",
		LiftID:              "lift-1",
		ProgressionID:       "prog-1",
		ConsecutiveFailures: 0,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	repo.Create(fc)

	// Update it
	failureTime := time.Now()
	fc.ConsecutiveFailures = 5
	fc.LastFailureAt = &failureTime
	fc.UpdatedAt = time.Now()

	err := repo.Update(fc)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	// Verify the update
	retrieved, _ := repo.GetByID("fc-001")
	if retrieved.ConsecutiveFailures != 5 {
		t.Errorf("ConsecutiveFailures = %d, want 5", retrieved.ConsecutiveFailures)
	}
	if retrieved.LastFailureAt == nil {
		t.Error("expected LastFailureAt to be set")
	}
}

func TestFailureCounterRepository_ConcurrentUpserts(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewFailureCounterRepository(db)

	// Simulate multiple failures in sequence
	for i := 1; i <= 10; i++ {
		count, err := repo.IncrementOnFailure("user-1", "lift-1", "prog-1")
		if err != nil {
			t.Fatalf("IncrementOnFailure() iteration %d failed: %v", i, err)
		}
		if count != i {
			t.Errorf("iteration %d: count = %d, want %d", i, count, i)
		}
	}

	// Reset
	repo.ResetOnSuccess("user-1", "lift-1", "prog-1")

	// Verify reset
	fc, _ := repo.GetByKey("user-1", "lift-1", "prog-1")
	if fc.ConsecutiveFailures != 0 {
		t.Errorf("after reset: ConsecutiveFailures = %d, want 0", fc.ConsecutiveFailures)
	}

	// Increment again after reset
	count, _ := repo.IncrementOnFailure("user-1", "lift-1", "prog-1")
	if count != 1 {
		t.Errorf("after reset and increment: count = %d, want 1", count)
	}
}
