// Package repository provides database repository implementations.
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ProgressionHistoryEntry represents a progression log entry with joined data.
type ProgressionHistoryEntry struct {
	ID              string          `json:"id"`
	ProgressionID   string          `json:"progressionId"`
	ProgressionName string          `json:"progressionName"`
	ProgressionType string          `json:"progressionType"`
	LiftID          string          `json:"liftId"`
	LiftName        string          `json:"liftName"`
	PreviousValue   float64         `json:"previousValue"`
	NewValue        float64         `json:"newValue"`
	Delta           float64         `json:"delta"`
	TriggerType     string          `json:"triggerType"`
	TriggerContext  json.RawMessage `json:"triggerContext"`
	AppliedAt       time.Time       `json:"appliedAt"`
}

// ProgressionHistoryFilter contains filter parameters for querying progression history.
type ProgressionHistoryFilter struct {
	UserID          string
	LiftID          *string
	ProgressionType *string
	TriggerType     *string
	StartDate       *time.Time
	EndDate         *time.Time
	Limit           int64
	Offset          int64
}

// ProgressionHistoryRepository handles progression history queries.
type ProgressionHistoryRepository struct {
	db *sql.DB
}

// NewProgressionHistoryRepository creates a new ProgressionHistoryRepository.
func NewProgressionHistoryRepository(sqlDB *sql.DB) *ProgressionHistoryRepository {
	return &ProgressionHistoryRepository{
		db: sqlDB,
	}
}

// List retrieves progression history entries with filters.
func (r *ProgressionHistoryRepository) List(ctx context.Context, filter ProgressionHistoryFilter) ([]ProgressionHistoryEntry, int64, error) {
	// Set defaults
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	// Build WHERE clause dynamically
	var whereConditions []string
	var args []interface{}

	whereConditions = append(whereConditions, "pl.user_id = ?")
	args = append(args, filter.UserID)

	if filter.LiftID != nil {
		whereConditions = append(whereConditions, "pl.lift_id = ?")
		args = append(args, *filter.LiftID)
	}

	if filter.ProgressionType != nil {
		whereConditions = append(whereConditions, "p.type = ?")
		args = append(args, *filter.ProgressionType)
	}

	if filter.TriggerType != nil {
		whereConditions = append(whereConditions, "pl.trigger_type = ?")
		args = append(args, *filter.TriggerType)
	}

	if filter.StartDate != nil {
		whereConditions = append(whereConditions, "pl.applied_at >= ?")
		args = append(args, filter.StartDate.Format(time.RFC3339))
	}

	if filter.EndDate != nil {
		whereConditions = append(whereConditions, "pl.applied_at <= ?")
		args = append(args, filter.EndDate.Format(time.RFC3339))
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Get total count
	countQuery := `
		SELECT COUNT(*) FROM progression_logs pl
		JOIN progressions p ON pl.progression_id = p.id
		WHERE ` + whereClause

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count progression history: %w", err)
	}

	// Get entries with pagination
	listQuery := `
		SELECT
			pl.id,
			pl.user_id,
			pl.progression_id,
			pl.lift_id,
			pl.previous_value,
			pl.new_value,
			pl.delta,
			pl.trigger_type,
			pl.trigger_context,
			pl.applied_at,
			p.name AS progression_name,
			p.type AS progression_type,
			l.name AS lift_name
		FROM progression_logs pl
		JOIN progressions p ON pl.progression_id = p.id
		JOIN lifts l ON pl.lift_id = l.id
		WHERE ` + whereClause + `
		ORDER BY pl.applied_at DESC
		LIMIT ? OFFSET ?`

	// Add pagination args
	listArgs := append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list progression history: %w", err)
	}
	defer rows.Close()

	var entries []ProgressionHistoryEntry
	for rows.Next() {
		var entry ProgressionHistoryEntry
		var userID, appliedAtStr, triggerContext string

		err := rows.Scan(
			&entry.ID,
			&userID,
			&entry.ProgressionID,
			&entry.LiftID,
			&entry.PreviousValue,
			&entry.NewValue,
			&entry.Delta,
			&entry.TriggerType,
			&triggerContext,
			&appliedAtStr,
			&entry.ProgressionName,
			&entry.ProgressionType,
			&entry.LiftName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan progression history entry: %w", err)
		}

		appliedAt, _ := time.Parse(time.RFC3339, appliedAtStr)
		entry.AppliedAt = appliedAt
		entry.TriggerContext = json.RawMessage(triggerContext)

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating progression history: %w", err)
	}

	return entries, total, nil
}
