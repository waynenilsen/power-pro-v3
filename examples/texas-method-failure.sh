#!/bin/bash
# Texas Method Failure Handling Demo
# Demonstrates the DeloadOnFailure progression system
#
# Texas Method uses consecutive failure tracking:
#   - 2 consecutive failures on intensity day -> 10% deload
#   - Success resets the failure counter
#
# Weekly structure:
#   Monday: Volume Day (5x5 at 90%)
#   Wednesday: Recovery Day (light)
#   Friday: Intensity Day (1x5 at 100%) <- progression happens here
#
# Usage: ./examples/texas-method-failure.sh
# Prerequisites: PowerPro server running on localhost:8080

set -e

API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
ADMIN_USER_ID="${ADMIN_USER_ID:-admin-user}"
TEST_USER_ID="texas-demo-$(date +%s)"

echo "=== Texas Method Failure Handling Demo ==="
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
echo "Step 2: Creating Texas Method DeloadOnFailure Progression..."

DELOAD_PROG_ID=$(api POST "/progressions" '{
    "name":"Texas Method Squat Deload",
    "type":"DELOAD_ON_FAILURE",
    "parameters":{
        "failureThreshold":2,
        "deloadType":"percent",
        "deloadPercent":0.10,
        "resetOnDeload":true,
        "maxType":"TRAINING_MAX"
    }
}' | extract_id)
echo "  DeloadOnFailure Progression: $DELOAD_PROG_ID"

# Also create a linear progression for success
LINEAR_PROG_ID=$(api POST "/progressions" '{
    "name":"Texas Method Linear +5lb",
    "type":"LINEAR_PROGRESSION",
    "parameters":{
        "increment":5.0,
        "maxType":"TRAINING_MAX",
        "triggerType":"AFTER_SESSION"
    }
}' | extract_id)
echo "  Linear Progression: $LINEAR_PROG_ID"

echo ""
echo "Step 3: Creating intensity day prescription..."

RX_SQUAT=$(api POST "/prescriptions" '{
    "liftId":"'"$SQUAT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0},
    "setScheme":{"type":"FIXED","sets":1,"reps":5},
    "order":1,
    "notes":"Texas Method Intensity Day - 1x5 PR attempt"
}' | extract_id)
echo "  Prescription: $RX_SQUAT"

echo ""
echo "Step 4: Creating Friday (Intensity Day)..."

DAY_ID=$(api POST "/days" '{"name":"Friday - Intensity","slug":"texas-friday-demo"}' | extract_id)
api POST "/days/$DAY_ID/prescriptions" '{"prescriptionId":"'"$RX_SQUAT"'"}' > /dev/null
echo "  Day: $DAY_ID"

echo ""
echo "Step 5: Creating cycle and week..."

CYCLE_ID=$(api POST "/cycles" '{"name":"Texas Method 1-Week Demo","lengthWeeks":1}' | extract_id)
echo "  Cycle: $CYCLE_ID"

WEEK_ID=$(api POST "/weeks" '{"cycleId":"'"$CYCLE_ID"'","weekNumber":1}' | extract_id)
api POST "/weeks/$WEEK_ID/days" '{"dayId":"'"$DAY_ID"'","dayOfWeek":"FRIDAY"}' > /dev/null
echo "  Week: $WEEK_ID"

echo ""
echo "Step 6: Creating program..."

PROGRAM_ID=$(api POST "/programs" '{
    "name":"Texas Method Demo",
    "slug":"texas-method-demo",
    "cycleId":"'"$CYCLE_ID"'"
}' | extract_id)
echo "  Program: $PROGRAM_ID"

api POST "/programs/$PROGRAM_ID/progressions" '{
    "progressionId":"'"$DELOAD_PROG_ID"'",
    "liftId":"'"$SQUAT_ID"'",
    "priority":1,
    "enabled":true
}' > /dev/null
api POST "/programs/$PROGRAM_ID/progressions" '{
    "progressionId":"'"$LINEAR_PROG_ID"'",
    "liftId":"'"$SQUAT_ID"'",
    "priority":2,
    "enabled":true
}' > /dev/null
echo "  Linked progressions to program"

echo ""
echo "Step 7: Setting up user..."

user_api POST "/users/$TEST_USER_ID/lift-maxes" '{
    "liftId":"'"$SQUAT_ID"'",
    "type":"TRAINING_MAX",
    "value":315.0,
    "effectiveDate":"2024-01-15T00:00:00Z"
}' > /dev/null
echo "  Squat Training Max: 315 lbs"

user_api POST "/users/$TEST_USER_ID/program" '{"programId":"'"$PROGRAM_ID"'"}' > /dev/null
echo "  Enrolled in Texas Method program"

echo ""
echo "=== Simulation: Week 1 (Intensity Day - 1x5 at 315 lbs) ==="
echo ""

echo "Attempt: 1x5 at 315 lbs"
echo "  Result: FAILURE - Only got 4 reps"
echo "  Failure counter: 0 -> 1"
echo ""

SESSION1_ID=$(uuidgen 2>/dev/null || cat /proc/sys/kernel/random/uuid)
user_api POST "/sessions/$SESSION1_ID/sets" '{
    "sets":[
        {"prescriptionId":"'"$RX_SQUAT"'","liftId":"'"$SQUAT_ID"'","setNumber":1,"weight":315.0,"targetReps":5,"repsPerformed":4,"isAmrap":false}
    ]
}' > /dev/null
echo "  Logged failure to session $SESSION1_ID"

echo "Triggering deload progression..."
TRIGGER1=$(user_api POST "/users/$TEST_USER_ID/progressions/trigger" '{
    "progressionId":"'"$DELOAD_PROG_ID"'",
    "liftId":"'"$SQUAT_ID"'",
    "force":true
}')
echo "Result: Progression NOT applied (threshold not met)"
echo "$TRIGGER1" | python3 -m json.tool 2>/dev/null || echo "$TRIGGER1"

# Advance state
user_api POST "/users/$TEST_USER_ID/program-state/advance" > /dev/null

echo ""
echo "=== Simulation: Week 2 (Intensity Day - 1x5 at 315 lbs) ==="
echo ""

echo "Attempt: 1x5 at 315 lbs"
echo "  Result: FAILURE - Only got 3 reps"
echo "  Failure counter: 1 -> 2 (THRESHOLD MET!)"
echo ""

SESSION2_ID=$(uuidgen 2>/dev/null || cat /proc/sys/kernel/random/uuid)
user_api POST "/sessions/$SESSION2_ID/sets" '{
    "sets":[
        {"prescriptionId":"'"$RX_SQUAT"'","liftId":"'"$SQUAT_ID"'","setNumber":1,"weight":315.0,"targetReps":5,"repsPerformed":3,"isAmrap":false}
    ]
}' > /dev/null
echo "  Logged failure to session $SESSION2_ID"

echo "Triggering deload progression..."
TRIGGER2=$(user_api POST "/users/$TEST_USER_ID/progressions/trigger" '{
    "progressionId":"'"$DELOAD_PROG_ID"'",
    "liftId":"'"$SQUAT_ID"'",
    "force":true
}')
echo "Result:"
echo "$TRIGGER2" | python3 -m json.tool 2>/dev/null || echo "$TRIGGER2"

echo ""
echo "=== Result ==="
echo "2 consecutive failures -> 10% DELOAD APPLIED"
echo "New Training Max: 315 - 31.5 = 283.5 lbs (rounded to 285 lbs)"
echo "Failure counter: RESET to 0"
echo ""

echo "=== Demo Complete ==="
echo ""
echo "Texas Method Failure Protocol:"
echo "  1. First stall: Keep weight, try again next week"
echo "  2. Second consecutive stall: Deload 10%"
echo "  3. Any success: Resets failure counter"
echo ""
echo "This prevents lifters from grinding too long at a weight they can't progress."
