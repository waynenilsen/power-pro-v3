#!/bin/bash
# Bill Starr 5x5 Program Configuration Example
# Demonstrates: Ramping sets, Daily lookups (H/L/M), weekly progression
#
# Usage: ./examples/bill-starr-5x5.sh
# Prerequisites: PowerPro server running on localhost:8080

set -e

API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
ADMIN_USER_ID="${ADMIN_USER_ID:-admin-user}"
TEST_USER_ID="${TEST_USER_ID:-test-user-001}"

echo "=== Bill Starr 5x5 Program Configuration ==="
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

echo "Step 1: Creating lifts..."

SQUAT_ID=$(api POST "/lifts" '{"name":"Squat","slug":"squat","isCompetitionLift":true}' | extract_id)
echo "  Squat: $SQUAT_ID"

BENCH_ID=$(api POST "/lifts" '{"name":"Bench Press","slug":"bench-press","isCompetitionLift":true}' | extract_id)
echo "  Bench Press: $BENCH_ID"

ROW_ID=$(api POST "/lifts" '{"name":"Bent Over Row","slug":"bent-over-row","isCompetitionLift":false}' | extract_id)
echo "  Bent Over Row: $ROW_ID"

INCLINE_ID=$(api POST "/lifts" '{"name":"Incline Press","slug":"incline-press","isCompetitionLift":false}' | extract_id)
echo "  Incline Press: $INCLINE_ID"

DEADLIFT_ID=$(api POST "/lifts" '{"name":"Deadlift","slug":"deadlift","isCompetitionLift":true}' | extract_id)
echo "  Deadlift: $DEADLIFT_ID"

echo ""
echo "Step 2: Creating daily lookup (H/L/M intensity)..."

DAILY_LOOKUP_ID=$(api POST "/daily-lookups" '{
    "name":"Bill Starr H/L/M Daily Intensity",
    "entries":[
        {"dayIdentifier":"heavy","percentageModifier":100.0,"intensityLevel":"HEAVY"},
        {"dayIdentifier":"light","percentageModifier":80.0,"intensityLevel":"LIGHT"},
        {"dayIdentifier":"medium","percentageModifier":90.0,"intensityLevel":"MEDIUM"}
    ]
}' | extract_id)
echo "  Daily Lookup: $DAILY_LOOKUP_ID"

echo ""
echo "Step 3: Creating prescriptions (ramping sets)..."

# Squat 5x5 Ramping for Heavy Day
RX_SQUAT_HEAVY=$(api POST "/prescriptions" '{
    "liftId":"'"$SQUAT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"day"},
    "setScheme":{"type":"RAMP","steps":[
        {"percentage":50,"reps":5},
        {"percentage":63,"reps":5},
        {"percentage":75,"reps":5},
        {"percentage":88,"reps":5},
        {"percentage":100,"reps":5}
    ],"workSetThreshold":80},
    "order":1,
    "notes":"Ramping to top set. Add weight each week.",
    "restSeconds":180
}' | extract_id)
echo "  Squat 5x5 Ramp: $RX_SQUAT_HEAVY"

# Bench Press 5x5 Ramping
RX_BENCH_HEAVY=$(api POST "/prescriptions" '{
    "liftId":"'"$BENCH_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"day"},
    "setScheme":{"type":"RAMP","steps":[
        {"percentage":50,"reps":5},
        {"percentage":63,"reps":5},
        {"percentage":75,"reps":5},
        {"percentage":88,"reps":5},
        {"percentage":100,"reps":5}
    ],"workSetThreshold":80},
    "order":2,
    "notes":"Ramping to top set",
    "restSeconds":180
}' | extract_id)
echo "  Bench 5x5 Ramp: $RX_BENCH_HEAVY"

# Bent Over Row 5x5 Ramping
RX_ROW_HEAVY=$(api POST "/prescriptions" '{
    "liftId":"'"$ROW_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"day"},
    "setScheme":{"type":"RAMP","steps":[
        {"percentage":50,"reps":5},
        {"percentage":63,"reps":5},
        {"percentage":75,"reps":5},
        {"percentage":88,"reps":5},
        {"percentage":100,"reps":5}
    ],"workSetThreshold":80},
    "order":3,
    "notes":"Ramping to top set",
    "restSeconds":180
}' | extract_id)
echo "  Row 5x5 Ramp: $RX_ROW_HEAVY"

# Light Day: Squat 4x5 capped ramp
RX_SQUAT_LIGHT=$(api POST "/prescriptions" '{
    "liftId":"'"$SQUAT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"day"},
    "setScheme":{"type":"RAMP","steps":[
        {"percentage":50,"reps":5},
        {"percentage":63,"reps":5},
        {"percentage":75,"reps":5},
        {"percentage":75,"reps":5}
    ],"workSetThreshold":70},
    "order":1,
    "notes":"Light day - capped at 75%",
    "restSeconds":120
}' | extract_id)
echo "  Squat 4x5 Light: $RX_SQUAT_LIGHT"

# Light Day: Incline Press 4x5
RX_INCLINE=$(api POST "/prescriptions" '{
    "liftId":"'"$INCLINE_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"day"},
    "setScheme":{"type":"RAMP","steps":[
        {"percentage":50,"reps":5},
        {"percentage":63,"reps":5},
        {"percentage":75,"reps":5},
        {"percentage":75,"reps":5}
    ],"workSetThreshold":70},
    "order":2,
    "notes":"Lighter pressing variant",
    "restSeconds":120
}' | extract_id)
echo "  Incline 4x5: $RX_INCLINE"

# Light Day: Deadlift 4x5
RX_DEADLIFT=$(api POST "/prescriptions" '{
    "liftId":"'"$DEADLIFT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"day"},
    "setScheme":{"type":"RAMP","steps":[
        {"percentage":50,"reps":5},
        {"percentage":63,"reps":5},
        {"percentage":75,"reps":5},
        {"percentage":75,"reps":5}
    ],"workSetThreshold":70},
    "order":3,
    "notes":"Light pulling day",
    "restSeconds":120
}' | extract_id)
echo "  Deadlift 4x5: $RX_DEADLIFT"

echo ""
echo "Step 4: Creating days..."

# Heavy Day (Monday)
DAY_HEAVY=$(api POST "/days" '{"name":"Heavy Day","slug":"bs5x5-heavy","metadata":{"focus":"squat-bench-row","intensity":"heavy","program":"bill-starr-5x5"}}' | extract_id)
echo "  Heavy Day: $DAY_HEAVY"

# Light Day (Wednesday)
DAY_LIGHT=$(api POST "/days" '{"name":"Light Day","slug":"bs5x5-light","metadata":{"focus":"squat-incline-deadlift","intensity":"light","program":"bill-starr-5x5"}}' | extract_id)
echo "  Light Day: $DAY_LIGHT"

# Medium Day (Friday)
DAY_MEDIUM=$(api POST "/days" '{"name":"Medium Day","slug":"bs5x5-medium","metadata":{"focus":"squat-bench-row","intensity":"medium","program":"bill-starr-5x5"}}' | extract_id)
echo "  Medium Day: $DAY_MEDIUM"

# Add prescriptions to Heavy Day
api POST "/days/$DAY_HEAVY/prescriptions" '{"prescriptionId":"'"$RX_SQUAT_HEAVY"'","order":1}' > /dev/null
api POST "/days/$DAY_HEAVY/prescriptions" '{"prescriptionId":"'"$RX_BENCH_HEAVY"'","order":2}' > /dev/null
api POST "/days/$DAY_HEAVY/prescriptions" '{"prescriptionId":"'"$RX_ROW_HEAVY"'","order":3}' > /dev/null
echo "  Added prescriptions to Heavy Day"

# Add prescriptions to Light Day
api POST "/days/$DAY_LIGHT/prescriptions" '{"prescriptionId":"'"$RX_SQUAT_LIGHT"'","order":1}' > /dev/null
api POST "/days/$DAY_LIGHT/prescriptions" '{"prescriptionId":"'"$RX_INCLINE"'","order":2}' > /dev/null
api POST "/days/$DAY_LIGHT/prescriptions" '{"prescriptionId":"'"$RX_DEADLIFT"'","order":3}' > /dev/null
echo "  Added prescriptions to Light Day"

# Add prescriptions to Medium Day (same as heavy but will be scaled by daily lookup)
api POST "/days/$DAY_MEDIUM/prescriptions" '{"prescriptionId":"'"$RX_SQUAT_HEAVY"'","order":1}' > /dev/null
api POST "/days/$DAY_MEDIUM/prescriptions" '{"prescriptionId":"'"$RX_BENCH_HEAVY"'","order":2}' > /dev/null
api POST "/days/$DAY_MEDIUM/prescriptions" '{"prescriptionId":"'"$RX_ROW_HEAVY"'","order":3}' > /dev/null
echo "  Added prescriptions to Medium Day"

echo ""
echo "Step 5: Creating cycle..."

CYCLE_ID=$(api POST "/cycles" '{"name":"Bill Starr 5x5 1-Week Cycle","lengthWeeks":1}' | extract_id)
echo "  Cycle: $CYCLE_ID"

echo ""
echo "Step 6: Creating week..."

WEEK_ID=$(api POST "/weeks" '{"cycleId":"'"$CYCLE_ID"'","weekNumber":1,"name":"Week 1"}' | extract_id)
echo "  Week 1: $WEEK_ID"

# Add days to week (Heavy/Light/Medium)
api POST "/weeks/$WEEK_ID/days" '{"dayId":"'"$DAY_HEAVY"'","position":0}' > /dev/null
api POST "/weeks/$WEEK_ID/days" '{"dayId":"'"$DAY_LIGHT"'","position":1}' > /dev/null
api POST "/weeks/$WEEK_ID/days" '{"dayId":"'"$DAY_MEDIUM"'","position":2}' > /dev/null
echo "  Added days to week (Heavy/Light/Medium)"

echo ""
echo "Step 7: Creating progressions..."

# Weekly linear +5lb
PROG_WEEKLY=$(api POST "/progressions" '{
    "name":"Weekly Linear +5lb",
    "type":"LINEAR_PROGRESSION",
    "parameters":{"increment":5.0,"maxType":"TRAINING_MAX","triggerType":"AFTER_WEEK"}
}' | extract_id)
echo "  Weekly Linear +5lb: $PROG_WEEKLY"

echo ""
echo "Step 8: Creating program..."

PROGRAM_ID=$(api POST "/programs" '{
    "name":"Bill Starr 5x5",
    "slug":"bill-starr-5x5",
    "description":"Classic intermediate program with ramping sets and Heavy/Light/Medium weekly structure. Weekly linear progression on all lifts.",
    "cycleId":"'"$CYCLE_ID"'",
    "dailyLookupId":"'"$DAILY_LOOKUP_ID"'",
    "defaultRounding":5.0
}' | extract_id)
echo "  Program: $PROGRAM_ID"

echo ""
echo "Step 9: Linking progressions to program..."

api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_WEEKLY"'","liftId":"'"$SQUAT_ID"'","priority":1}' > /dev/null
echo "  Squat -> Weekly +5lb"
api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_WEEKLY"'","liftId":"'"$BENCH_ID"'","priority":1}' > /dev/null
echo "  Bench -> Weekly +5lb"
api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_WEEKLY"'","liftId":"'"$ROW_ID"'","priority":1}' > /dev/null
echo "  Row -> Weekly +5lb"
api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_WEEKLY"'","liftId":"'"$INCLINE_ID"'","priority":1}' > /dev/null
echo "  Incline -> Weekly +5lb"
api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_WEEKLY"'","liftId":"'"$DEADLIFT_ID"'","priority":1}' > /dev/null
echo "  Deadlift -> Weekly +5lb"

echo ""
echo "Step 10: Setting up test user..."

# Create training maxes for test user
curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$SQUAT_ID"'","type":"TRAINING_MAX","value":300.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Squat TM: 300 lbs"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$BENCH_ID"'","type":"TRAINING_MAX","value":200.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Bench TM: 200 lbs"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$ROW_ID"'","type":"TRAINING_MAX","value":185.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Row TM: 185 lbs"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$INCLINE_ID"'","type":"TRAINING_MAX","value":160.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Incline TM: 160 lbs"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$DEADLIFT_ID"'","type":"TRAINING_MAX","value":350.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Deadlift TM: 350 lbs"

echo ""
echo "Step 11: Enrolling user in program..."

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/program" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"programId":"'"$PROGRAM_ID"'"}' > /dev/null
echo "  User enrolled in Bill Starr 5x5"

echo ""
echo "Step 12: Generating workout (Heavy Day)..."
echo ""

WORKOUT=$(curl -s -X GET "$API_BASE_URL/users/$TEST_USER_ID/workout" \
    -H "X-User-ID: $TEST_USER_ID")

echo "=== Generated Workout (Heavy Day) ==="
echo "$WORKOUT"

echo ""
echo "=== Bill Starr 5x5 Configuration Complete ==="
echo ""
echo "Program ID: $PROGRAM_ID"
echo "Test User: $TEST_USER_ID"
echo ""
echo "Key features demonstrated:"
echo "  - RAMP set scheme (weights increase across sets)"
echo "  - Daily lookup (H/L/M intensity modifiers)"
echo "  - Weekly linear progression"
