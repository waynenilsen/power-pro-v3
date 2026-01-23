#!/bin/bash
# Starting Strength Program Configuration Example
# Demonstrates: Linear per-session progression, simple 3x5 sets
#
# Usage: ./examples/starting-strength.sh
# Prerequisites: PowerPro server running on localhost:8080

set -e

API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
ADMIN_USER_ID="${ADMIN_USER_ID:-admin-user}"
TEST_USER_ID="${TEST_USER_ID:-test-user-001}"

echo "=== Starting Strength Program Configuration ==="
echo "API: $API_BASE_URL"
echo ""

# Helper function for API calls
api() {
    local method="$1"
    local path="$2"
    local data="$3"

    if [ -n "$data" ]; then
        curl -s -X "$method" "$API_BASE_URL$path" \
            -H "X-User-ID: $ADMIN_USER_ID" \
            -H "X-Admin: true" \
            -H "Content-Type: application/json" \
            -d "$data"
    else
        curl -s -X "$method" "$API_BASE_URL$path" \
            -H "X-User-ID: $ADMIN_USER_ID" \
            -H "X-Admin: true"
    fi
}

# Extract ID from response
extract_id() {
    grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4
}

# Get or create a lift by slug
get_or_create_lift() {
    local name="$1"
    local slug="$2"
    local is_comp="$3"

    # Try to get by slug first
    local existing_id
    existing_id=$(api GET "/lifts/by-slug/$slug" | extract_id)

    if [ -n "$existing_id" ]; then
        echo "$existing_id"
        return
    fi

    # Create new
    api POST "/lifts" "{\"name\":\"$name\",\"slug\":\"$slug\",\"isCompetitionLift\":$is_comp}" | extract_id
}

echo "Step 1: Creating lifts..."

SQUAT_ID=$(get_or_create_lift "Squat" "squat" "true")
echo "  Squat: $SQUAT_ID"

BENCH_ID=$(get_or_create_lift "Bench Press" "bench-press" "true")
echo "  Bench Press: $BENCH_ID"

PRESS_ID=$(get_or_create_lift "Overhead Press" "overhead-press" "false")
echo "  Overhead Press: $PRESS_ID"

DEADLIFT_ID=$(get_or_create_lift "Deadlift" "deadlift" "true")
echo "  Deadlift: $DEADLIFT_ID"

CLEAN_ID=$(get_or_create_lift "Power Clean" "power-clean" "false")
echo "  Power Clean: $CLEAN_ID"

echo ""
echo "Step 2: Creating prescriptions..."

# Squat 3x5
RX_SQUAT=$(api POST "/prescriptions" '{
    "liftId":"'"$SQUAT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"FIXED","sets":3,"reps":5,"isAmrap":false},
    "order":1,
    "notes":"Focus on depth - crease of hip below top of knee",
    "restSeconds":180
}' | extract_id)
echo "  Squat 3x5: $RX_SQUAT"

# Bench 3x5
RX_BENCH=$(api POST "/prescriptions" '{
    "liftId":"'"$BENCH_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"FIXED","sets":3,"reps":5,"isAmrap":false},
    "order":2,
    "notes":"Touch chest, pause briefly, press",
    "restSeconds":180
}' | extract_id)
echo "  Bench 3x5: $RX_BENCH"

# Press 3x5
RX_PRESS=$(api POST "/prescriptions" '{
    "liftId":"'"$PRESS_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"FIXED","sets":3,"reps":5,"isAmrap":false},
    "order":2,
    "notes":"Strict press - no leg drive",
    "restSeconds":180
}' | extract_id)
echo "  Press 3x5: $RX_PRESS"

# Deadlift 1x5
RX_DEADLIFT=$(api POST "/prescriptions" '{
    "liftId":"'"$DEADLIFT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"FIXED","sets":1,"reps":5,"isAmrap":false},
    "order":3,
    "notes":"Reset each rep - no touch and go",
    "restSeconds":180
}' | extract_id)
echo "  Deadlift 1x5: $RX_DEADLIFT"

# Power Clean 5x3
RX_CLEAN=$(api POST "/prescriptions" '{
    "liftId":"'"$CLEAN_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"FIXED","sets":5,"reps":3,"isAmrap":false},
    "order":3,
    "notes":"Explosive triple extension",
    "restSeconds":120
}' | extract_id)
echo "  Power Clean 5x3: $RX_CLEAN"

echo ""
echo "Step 3: Creating days..."

# Day A
DAY_A=$(api POST "/days" '{"name":"Day A","slug":"ss-day-a","metadata":{"focus":"squat-bench-deadlift","program":"starting-strength"}}' | extract_id)
echo "  Day A: $DAY_A"

# Day B
DAY_B=$(api POST "/days" '{"name":"Day B","slug":"ss-day-b","metadata":{"focus":"squat-press-clean","program":"starting-strength"}}' | extract_id)
echo "  Day B: $DAY_B"

# Add prescriptions to Day A
api POST "/days/$DAY_A/prescriptions" '{"prescriptionId":"'"$RX_SQUAT"'","order":1}' > /dev/null
api POST "/days/$DAY_A/prescriptions" '{"prescriptionId":"'"$RX_BENCH"'","order":2}' > /dev/null
api POST "/days/$DAY_A/prescriptions" '{"prescriptionId":"'"$RX_DEADLIFT"'","order":3}' > /dev/null
echo "  Added prescriptions to Day A"

# Add prescriptions to Day B
api POST "/days/$DAY_B/prescriptions" '{"prescriptionId":"'"$RX_SQUAT"'","order":1}' > /dev/null
api POST "/days/$DAY_B/prescriptions" '{"prescriptionId":"'"$RX_PRESS"'","order":2}' > /dev/null
api POST "/days/$DAY_B/prescriptions" '{"prescriptionId":"'"$RX_CLEAN"'","order":3}' > /dev/null
echo "  Added prescriptions to Day B"

echo ""
echo "Step 4: Creating cycle..."

CYCLE_ID=$(api POST "/cycles" '{"name":"Starting Strength 1-Week Cycle","lengthWeeks":1}' | extract_id)
echo "  Cycle: $CYCLE_ID"

echo ""
echo "Step 5: Creating week..."

WEEK_ID=$(api POST "/weeks" '{"cycleId":"'"$CYCLE_ID"'","weekNumber":1,"name":"Week 1"}' | extract_id)
echo "  Week 1: $WEEK_ID"

# Add days to week (A/B/A pattern)
api POST "/weeks/$WEEK_ID/days" '{"dayId":"'"$DAY_A"'","position":0}' > /dev/null
api POST "/weeks/$WEEK_ID/days" '{"dayId":"'"$DAY_B"'","position":1}' > /dev/null
api POST "/weeks/$WEEK_ID/days" '{"dayId":"'"$DAY_A"'","position":2}' > /dev/null
echo "  Added days to week (A/B/A pattern)"

echo ""
echo "Step 6: Creating progressions..."

# Linear +5lb for upper body
PROG_UPPER=$(api POST "/progressions" '{
    "name":"Linear +5lb Upper Body",
    "type":"LINEAR_PROGRESSION",
    "parameters":{"increment":5.0,"maxType":"TRAINING_MAX","triggerType":"AFTER_SESSION"}
}' | extract_id)
echo "  Linear +5lb Upper: $PROG_UPPER"

# Linear +5lb for squat
PROG_SQUAT=$(api POST "/progressions" '{
    "name":"Linear +5lb Squat",
    "type":"LINEAR_PROGRESSION",
    "parameters":{"increment":5.0,"maxType":"TRAINING_MAX","triggerType":"AFTER_SESSION"}
}' | extract_id)
echo "  Linear +5lb Squat: $PROG_SQUAT"

# Linear +10lb for deadlift
PROG_DEADLIFT=$(api POST "/progressions" '{
    "name":"Linear +10lb Deadlift",
    "type":"LINEAR_PROGRESSION",
    "parameters":{"increment":10.0,"maxType":"TRAINING_MAX","triggerType":"AFTER_SESSION"}
}' | extract_id)
echo "  Linear +10lb Deadlift: $PROG_DEADLIFT"

echo ""
echo "Step 7: Creating program..."

PROGRAM_ID=$(api POST "/programs" '{
    "name":"Starting Strength",
    "slug":"starting-strength",
    "description":"Mark Rippetoe'\''s novice linear progression program. Three days per week, alternating A/B workouts with linear progression every session.",
    "cycleId":"'"$CYCLE_ID"'",
    "defaultRounding":5.0
}' | extract_id)
echo "  Program: $PROGRAM_ID"

echo ""
echo "Step 8: Linking progressions to program..."

api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_SQUAT"'","liftId":"'"$SQUAT_ID"'","priority":1}' > /dev/null
echo "  Squat -> Linear +5lb"
api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_UPPER"'","liftId":"'"$BENCH_ID"'","priority":1}' > /dev/null
echo "  Bench -> Linear +5lb"
api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_UPPER"'","liftId":"'"$PRESS_ID"'","priority":1}' > /dev/null
echo "  Press -> Linear +5lb"
api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_DEADLIFT"'","liftId":"'"$DEADLIFT_ID"'","priority":1}' > /dev/null
echo "  Deadlift -> Linear +10lb"
api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_UPPER"'","liftId":"'"$CLEAN_ID"'","priority":1}' > /dev/null
echo "  Power Clean -> Linear +5lb"

echo ""
echo "Step 9: Setting up test user..."

# Create training maxes for test user
curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$SQUAT_ID"'","type":"TRAINING_MAX","value":225.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Squat TM: 225 lbs"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$BENCH_ID"'","type":"TRAINING_MAX","value":155.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Bench TM: 155 lbs"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$PRESS_ID"'","type":"TRAINING_MAX","value":95.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Press TM: 95 lbs"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$DEADLIFT_ID"'","type":"TRAINING_MAX","value":275.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Deadlift TM: 275 lbs"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$CLEAN_ID"'","type":"TRAINING_MAX","value":135.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Power Clean TM: 135 lbs"

echo ""
echo "Step 10: Enrolling user in program..."

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/program" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"programId":"'"$PROGRAM_ID"'"}' > /dev/null
echo "  User enrolled in Starting Strength"

echo ""
echo "Step 11: Generating workout..."
echo ""

WORKOUT=$(curl -s -X GET "$API_BASE_URL/users/$TEST_USER_ID/workout" \
    -H "X-User-ID: $TEST_USER_ID")

echo "=== Generated Workout (Day A) ==="
echo "$WORKOUT"

echo ""
echo "=== Starting Strength Configuration Complete ==="
echo ""
echo "Program ID: $PROGRAM_ID"
echo "Test User: $TEST_USER_ID"
echo ""
echo "To generate the next workout, advance the day and call /workout again."
