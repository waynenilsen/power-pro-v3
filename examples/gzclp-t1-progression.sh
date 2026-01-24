#!/bin/bash
# GZCLP T1 Stage Progression Demo
# Demonstrates the stage-based progression system with failure handling
#
# GZCLP T1 uses stage progression through failure:
#   Stage 0: 5x3+ (minimum 15 reps to pass)
#   Stage 1: 6x2+ (minimum 12 reps to pass)
#   Stage 2: 10x1+ (minimum 10 reps to pass)
#
# On failure: advance to next stage, keep weight
# On final stage failure: reset to stage 0 with 15% deload
#
# Usage: ./examples/gzclp-t1-progression.sh
# Prerequisites: PowerPro server running on localhost:8080

set -e

API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
ADMIN_USER_ID="${ADMIN_USER_ID:-admin-user}"
TEST_USER_ID="gzclp-t1-demo-$(date +%s)"

echo "=== GZCLP T1 Stage Progression Demo ==="
echo "API: $API_BASE_URL"
echo "Test User: $TEST_USER_ID"
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

user_api() {
    local method="$1"
    local path="$2"
    local data="$3"

    if [ -n "$data" ]; then
        curl -s -X "$method" "$API_BASE_URL$path" \
            -H "X-User-ID: $TEST_USER_ID" \
            -H "Content-Type: application/json" \
            -d "$data"
    else
        curl -s -X "$method" "$API_BASE_URL$path" \
            -H "X-User-ID: $TEST_USER_ID"
    fi
}

# Extract ID from response
extract_id() {
    grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4
}

echo "Step 1: Getting seeded squat lift..."
SQUAT_ID="00000000-0000-0000-0000-000000000001"
echo "  Squat ID: $SQUAT_ID"

echo ""
echo "Step 2: Creating GZCLP T1 Stage Progression..."

STAGE_PROG_ID=$(api POST "/progressions" '{
    "name":"GZCLP T1 Squat Demo",
    "type":"STAGE_PROGRESSION",
    "parameters":{
        "stages":[
            {"name":"5x3+","sets":5,"reps":3,"isAmrap":true,"minVolume":15},
            {"name":"6x2+","sets":6,"reps":2,"isAmrap":true,"minVolume":12},
            {"name":"10x1+","sets":10,"reps":1,"isAmrap":true,"minVolume":10}
        ],
        "currentStage":0,
        "resetOnExhaustion":true,
        "deloadOnReset":true,
        "deloadPercent":0.15,
        "maxType":"TRAINING_MAX"
    }
}' | extract_id)
echo "  Stage Progression: $STAGE_PROG_ID"

echo ""
echo "Step 3: Creating prescription..."

RX_SQUAT=$(api POST "/prescriptions" '{
    "liftId":"'"$SQUAT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0},
    "setScheme":{"type":"AMRAP","sets":5,"reps":3},
    "order":1,
    "notes":"GZCLP T1 Squat - Stage-based progression"
}' | extract_id)
echo "  Prescription: $RX_SQUAT"

echo ""
echo "Step 4: Creating day..."

DAY_ID=$(api POST "/days" '{"name":"GZCLP Day A","slug":"gzclp-day-a-demo"}' | extract_id)
api POST "/days/$DAY_ID/prescriptions" '{"prescriptionId":"'"$RX_SQUAT"'"}' > /dev/null
echo "  Day: $DAY_ID"

echo ""
echo "Step 5: Creating cycle and week..."

CYCLE_ID=$(api POST "/cycles" '{"name":"GZCLP 1-Week Cycle Demo","lengthWeeks":1}' | extract_id)
echo "  Cycle: $CYCLE_ID"

WEEK_ID=$(api POST "/weeks" '{"cycleId":"'"$CYCLE_ID"'","weekNumber":1}' | extract_id)
api POST "/weeks/$WEEK_ID/days" '{"dayId":"'"$DAY_ID"'","dayOfWeek":"MONDAY"}' > /dev/null
echo "  Week: $WEEK_ID"

echo ""
echo "Step 6: Creating program..."

PROGRAM_ID=$(api POST "/programs" '{
    "name":"GZCLP Demo",
    "slug":"gzclp-demo",
    "cycleId":"'"$CYCLE_ID"'"
}' | extract_id)
echo "  Program: $PROGRAM_ID"

api POST "/programs/$PROGRAM_ID/progressions" '{
    "progressionId":"'"$STAGE_PROG_ID"'",
    "liftId":"'"$SQUAT_ID"'",
    "priority":1,
    "enabled":true
}' > /dev/null
echo "  Linked progression to program"

echo ""
echo "Step 7: Setting up user..."

user_api POST "/users/$TEST_USER_ID/lift-maxes" '{
    "liftId":"'"$SQUAT_ID"'",
    "type":"TRAINING_MAX",
    "value":200.0,
    "effectiveDate":"2024-01-15T00:00:00Z"
}' > /dev/null
echo "  Squat Training Max: 200 lbs"

user_api POST "/users/$TEST_USER_ID/program" '{"programId":"'"$PROGRAM_ID"'"}' > /dev/null
echo "  Enrolled in GZCLP program"

echo ""
echo "=== Simulation: Week 1 (5x3+ at 200 lbs) ==="
echo ""

# Get initial workout
WORKOUT=$(user_api GET "/users/$TEST_USER_ID/workout")
echo "Initial workout:"
echo "$WORKOUT" | python3 -m json.tool 2>/dev/null || echo "$WORKOUT"
echo ""

echo "Simulating FAILURE: Total 13 reps (below 15 minimum)"
echo "  Set 1: 3 reps ✓"
echo "  Set 2: 3 reps ✓"
echo "  Set 3: 3 reps ✓"
echo "  Set 4: 3 reps ✓"
echo "  Set 5 (AMRAP): 1 rep ✗ (failed to get minimum)"
echo "  Total: 13 reps < 15 minimum"
echo ""

# Log the failure sets
SESSION_ID=$(uuidgen 2>/dev/null || cat /proc/sys/kernel/random/uuid)
user_api POST "/sessions/$SESSION_ID/sets" '{
    "sets":[
        {"prescriptionId":"'"$RX_SQUAT"'","liftId":"'"$SQUAT_ID"'","setNumber":1,"weight":200.0,"targetReps":3,"repsPerformed":3,"isAmrap":false},
        {"prescriptionId":"'"$RX_SQUAT"'","liftId":"'"$SQUAT_ID"'","setNumber":2,"weight":200.0,"targetReps":3,"repsPerformed":3,"isAmrap":false},
        {"prescriptionId":"'"$RX_SQUAT"'","liftId":"'"$SQUAT_ID"'","setNumber":3,"weight":200.0,"targetReps":3,"repsPerformed":3,"isAmrap":false},
        {"prescriptionId":"'"$RX_SQUAT"'","liftId":"'"$SQUAT_ID"'","setNumber":4,"weight":200.0,"targetReps":3,"repsPerformed":3,"isAmrap":false},
        {"prescriptionId":"'"$RX_SQUAT"'","liftId":"'"$SQUAT_ID"'","setNumber":5,"weight":200.0,"targetReps":3,"repsPerformed":1,"isAmrap":true}
    ]
}' > /dev/null
echo "  Logged sets to session $SESSION_ID"

echo ""
echo "Triggering stage progression..."
TRIGGER_RESULT=$(user_api POST "/users/$TEST_USER_ID/progressions/trigger" '{
    "progressionId":"'"$STAGE_PROG_ID"'",
    "liftId":"'"$SQUAT_ID"'",
    "force":true
}')
echo "Trigger result:"
echo "$TRIGGER_RESULT" | python3 -m json.tool 2>/dev/null || echo "$TRIGGER_RESULT"

echo ""
echo "=== Result ==="
echo "Stage 0 (5x3+) FAILED -> Advanced to Stage 1 (6x2+)"
echo "Weight remains: 200 lbs (stages change rep scheme, not weight)"
echo ""

echo "=== Demo Complete ==="
echo ""
echo "Next steps in a real scenario:"
echo "  1. Week 2: Attempt 6x2+ at 200 lbs"
echo "  2. If fail: Advance to Stage 2 (10x1+)"
echo "  3. Week 3: Attempt 10x1+ at 200 lbs"
echo "  4. If fail: Reset to Stage 0 (5x3+) with 15% deload = 170 lbs"
