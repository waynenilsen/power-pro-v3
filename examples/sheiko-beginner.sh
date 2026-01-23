#!/bin/bash
# Sheiko Beginner Program Configuration Example
# Demonstrates: High-volume training, RAMP sets, NO automatic progression (manual only)
#
# Usage: ./examples/sheiko-beginner.sh
# Prerequisites: PowerPro server running on localhost:8080

set -e

API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
ADMIN_USER_ID="${ADMIN_USER_ID:-admin-user}"
TEST_USER_ID="${TEST_USER_ID:-test-user-001}"

echo "=== Sheiko Beginner Program Configuration ==="
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

DEADLIFT_ID=$(api POST "/lifts" '{"name":"Deadlift","slug":"deadlift","isCompetitionLift":true}' | extract_id)
echo "  Deadlift: $DEADLIFT_ID"

# Variations
DL_KNEES_ID=$(api POST "/lifts" '{"name":"Deadlift Up to Knees","slug":"deadlift-up-to-knees","isCompetitionLift":false,"parentLiftId":"'"$DEADLIFT_ID"'"}' | extract_id)
echo "  Deadlift Up to Knees: $DL_KNEES_ID"

DL_BOXES_ID=$(api POST "/lifts" '{"name":"Deadlift from Boxes","slug":"deadlift-from-boxes","isCompetitionLift":false,"parentLiftId":"'"$DEADLIFT_ID"'"}' | extract_id)
echo "  Deadlift from Boxes: $DL_BOXES_ID"

INCLINE_ID=$(api POST "/lifts" '{"name":"Incline Bench Press","slug":"incline-bench-press","isCompetitionLift":false,"parentLiftId":"'"$BENCH_ID"'"}' | extract_id)
echo "  Incline Bench: $INCLINE_ID"

GM_ID=$(api POST "/lifts" '{"name":"Good Morning","slug":"good-morning","isCompetitionLift":false}' | extract_id)
echo "  Good Morning: $GM_ID"

echo ""
echo "Step 2: Creating prescriptions (high-volume Sheiko style)..."

# Day 1: Squat (ramping to 80%, high volume)
RX_SQUAT_D1=$(api POST "/prescriptions" '{
    "liftId":"'"$SQUAT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"ONE_RM","percentage":100.0,"roundingIncrement":2.5,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"RAMP","steps":[
        {"percentage":50,"reps":5},
        {"percentage":60,"reps":4},
        {"percentage":70,"reps":3},
        {"percentage":70,"reps":3},
        {"percentage":75,"reps":3},
        {"percentage":75,"reps":2},
        {"percentage":80,"reps":2},
        {"percentage":80,"reps":2},
        {"percentage":75,"reps":3},
        {"percentage":75,"reps":4}
    ],"workSetThreshold":70},
    "order":1,
    "notes":"Focus on technique. No grinding reps.",
    "restSeconds":180
}' | extract_id)
echo "  Squat Day 1 (27 reps): $RX_SQUAT_D1"

# Day 1: Bench (ramping to 75%)
RX_BENCH_D1=$(api POST "/prescriptions" '{
    "liftId":"'"$BENCH_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"ONE_RM","percentage":100.0,"roundingIncrement":2.5,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"RAMP","steps":[
        {"percentage":50,"reps":5},
        {"percentage":60,"reps":4},
        {"percentage":70,"reps":3},
        {"percentage":70,"reps":3},
        {"percentage":70,"reps":3},
        {"percentage":75,"reps":2},
        {"percentage":75,"reps":2},
        {"percentage":70,"reps":3},
        {"percentage":70,"reps":4}
    ],"workSetThreshold":70},
    "order":2,
    "notes":"Paused reps. Build bar speed.",
    "restSeconds":150
}' | extract_id)
echo "  Bench Day 1 (26 reps): $RX_BENCH_D1"

# Day 3: Deadlift Up to Knees (technique work)
RX_DL_KNEES=$(api POST "/prescriptions" '{
    "liftId":"'"$DL_KNEES_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"ONE_RM","percentage":100.0,"roundingIncrement":2.5,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"RAMP","steps":[
        {"percentage":50,"reps":3},
        {"percentage":60,"reps":3},
        {"percentage":65,"reps":3},
        {"percentage":65,"reps":3},
        {"percentage":65,"reps":3}
    ],"workSetThreshold":60},
    "order":1,
    "notes":"Stop at knees. Focus on leg drive.",
    "restSeconds":120
}' | extract_id)
echo "  DL Up to Knees (15 reps): $RX_DL_KNEES"

# Day 3: Bench (medium volume)
RX_BENCH_D3=$(api POST "/prescriptions" '{
    "liftId":"'"$BENCH_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"ONE_RM","percentage":100.0,"roundingIncrement":2.5,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"RAMP","steps":[
        {"percentage":50,"reps":5},
        {"percentage":60,"reps":4},
        {"percentage":70,"reps":3},
        {"percentage":75,"reps":3},
        {"percentage":75,"reps":3},
        {"percentage":80,"reps":2},
        {"percentage":80,"reps":2},
        {"percentage":75,"reps":3}
    ],"workSetThreshold":70},
    "order":2,
    "notes":"Build to heavier weights.",
    "restSeconds":180
}' | extract_id)
echo "  Bench Day 3 (22 reps): $RX_BENCH_D3"

# Day 3: Deadlift from Boxes
RX_DL_BOXES=$(api POST "/prescriptions" '{
    "liftId":"'"$DL_BOXES_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"ONE_RM","percentage":100.0,"roundingIncrement":2.5,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"RAMP","steps":[
        {"percentage":50,"reps":3},
        {"percentage":60,"reps":3},
        {"percentage":70,"reps":3},
        {"percentage":70,"reps":3},
        {"percentage":75,"reps":2},
        {"percentage":75,"reps":2}
    ],"workSetThreshold":65},
    "order":3,
    "notes":"Deficit pulling. Maintain position.",
    "restSeconds":150
}' | extract_id)
echo "  DL from Boxes (16 reps): $RX_DL_BOXES"

# Day 5: Squat (lighter, technique)
RX_SQUAT_D5=$(api POST "/prescriptions" '{
    "liftId":"'"$SQUAT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"ONE_RM","percentage":100.0,"roundingIncrement":2.5,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"RAMP","steps":[
        {"percentage":50,"reps":5},
        {"percentage":60,"reps":4},
        {"percentage":70,"reps":3},
        {"percentage":70,"reps":3},
        {"percentage":70,"reps":3},
        {"percentage":75,"reps":2},
        {"percentage":75,"reps":2}
    ],"workSetThreshold":70},
    "order":1,
    "notes":"Recovery day. Perfect technique.",
    "restSeconds":150
}' | extract_id)
echo "  Squat Day 5 (19 reps): $RX_SQUAT_D5"

# Day 5: Incline Bench
RX_INCLINE=$(api POST "/prescriptions" '{
    "liftId":"'"$INCLINE_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"ONE_RM","percentage":100.0,"roundingIncrement":2.5,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"RAMP","steps":[
        {"percentage":50,"reps":5},
        {"percentage":55,"reps":4},
        {"percentage":60,"reps":3},
        {"percentage":60,"reps":3},
        {"percentage":60,"reps":3}
    ],"workSetThreshold":55},
    "order":2,
    "notes":"Upper chest development.",
    "restSeconds":120
}' | extract_id)
echo "  Incline Bench (15 reps): $RX_INCLINE"

# Day 5: Good Morning
RX_GM=$(api POST "/prescriptions" '{
    "liftId":"'"$GM_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"ONE_RM","percentage":100.0,"roundingIncrement":2.5,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"FIXED","sets":4,"reps":6,"isAmrap":false},
    "order":3,
    "notes":"50-60% of squat. Posterior chain.",
    "restSeconds":90
}' | extract_id)
echo "  Good Morning (24 reps): $RX_GM"

echo ""
echo "Step 3: Creating days..."

DAY_1=$(api POST "/days" '{"name":"Day 1 (Monday)","slug":"sheiko-day-1","metadata":{"focus":"squat-bench","volume":"high","program":"sheiko-beginner"}}' | extract_id)
echo "  Day 1 (Squat + Bench): $DAY_1"

DAY_3=$(api POST "/days" '{"name":"Day 3 (Wednesday)","slug":"sheiko-day-3","metadata":{"focus":"deadlift-bench","volume":"moderate","program":"sheiko-beginner"}}' | extract_id)
echo "  Day 3 (Deadlift + Bench): $DAY_3"

DAY_5=$(api POST "/days" '{"name":"Day 5 (Friday)","slug":"sheiko-day-5","metadata":{"focus":"squat-variations","volume":"moderate","program":"sheiko-beginner"}}' | extract_id)
echo "  Day 5 (Squat + Variations): $DAY_5"

# Add prescriptions to days
api POST "/days/$DAY_1/prescriptions" '{"prescriptionId":"'"$RX_SQUAT_D1"'","order":1}' > /dev/null
api POST "/days/$DAY_1/prescriptions" '{"prescriptionId":"'"$RX_BENCH_D1"'","order":2}' > /dev/null
echo "  Added prescriptions to Day 1"

api POST "/days/$DAY_3/prescriptions" '{"prescriptionId":"'"$RX_DL_KNEES"'","order":1}' > /dev/null
api POST "/days/$DAY_3/prescriptions" '{"prescriptionId":"'"$RX_BENCH_D3"'","order":2}' > /dev/null
api POST "/days/$DAY_3/prescriptions" '{"prescriptionId":"'"$RX_DL_BOXES"'","order":3}' > /dev/null
echo "  Added prescriptions to Day 3"

api POST "/days/$DAY_5/prescriptions" '{"prescriptionId":"'"$RX_SQUAT_D5"'","order":1}' > /dev/null
api POST "/days/$DAY_5/prescriptions" '{"prescriptionId":"'"$RX_INCLINE"'","order":2}' > /dev/null
api POST "/days/$DAY_5/prescriptions" '{"prescriptionId":"'"$RX_GM"'","order":3}' > /dev/null
echo "  Added prescriptions to Day 5"

echo ""
echo "Step 4: Creating cycle..."

CYCLE_ID=$(api POST "/cycles" '{"name":"Sheiko Beginner Prep Phase","lengthWeeks":4}' | extract_id)
echo "  Cycle: $CYCLE_ID"

echo ""
echo "Step 5: Creating weeks..."

# For simplicity, all 4 weeks use same day structure
# In full Sheiko, weeks vary slightly
for week_num in 1 2 3 4; do
    WEEK_ID=$(api POST "/weeks" '{"cycleId":"'"$CYCLE_ID"'","weekNumber":'"$week_num"',"name":"Week '"$week_num"'"}' | extract_id)
    echo "  Week $week_num: $WEEK_ID"
    api POST "/weeks/$WEEK_ID/days" '{"dayId":"'"$DAY_1"'","position":0}' > /dev/null
    api POST "/weeks/$WEEK_ID/days" '{"dayId":"'"$DAY_3"'","position":1}' > /dev/null
    api POST "/weeks/$WEEK_ID/days" '{"dayId":"'"$DAY_5"'","position":2}' > /dev/null
done
echo "  Added days to all weeks"

echo ""
echo "Step 6: Creating program (NO automatic progressions)..."

PROGRAM_ID=$(api POST "/programs" '{
    "name":"Sheiko Beginner",
    "slug":"sheiko-beginner",
    "description":"Boris Sheiko'\''s beginner powerlifting program. High-volume, moderate-intensity training focused on technical perfection. Manual progression only - athlete updates maxes when ready.",
    "cycleId":"'"$CYCLE_ID"'",
    "defaultRounding":2.5
}' | extract_id)
echo "  Program: $PROGRAM_ID"
echo ""
echo "  NOTE: No progressions linked - Sheiko uses MANUAL max updates"

echo ""
echo "Step 7: Setting up test user..."

# Create 1RM estimates for test user (Sheiko uses estimated 1RM, not training max)
curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$SQUAT_ID"'","type":"ONE_RM","value":180.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Squat 1RM: 180 kg"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$BENCH_ID"'","type":"ONE_RM","value":120.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Bench 1RM: 120 kg"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$DEADLIFT_ID"'","type":"ONE_RM","value":200.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Deadlift 1RM: 200 kg"

# Variation maxes inherit from parent or are set separately
curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$DL_KNEES_ID"'","type":"ONE_RM","value":200.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  DL Up to Knees 1RM: 200 kg (same as deadlift)"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$DL_BOXES_ID"'","type":"ONE_RM","value":180.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  DL from Boxes 1RM: 180 kg (harder variation)"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$INCLINE_ID"'","type":"ONE_RM","value":100.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Incline 1RM: 100 kg"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$GM_ID"'","type":"ONE_RM","value":100.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Good Morning 1RM: 100 kg"

echo ""
echo "Step 8: Enrolling user in program..."

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/program" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"programId":"'"$PROGRAM_ID"'"}' > /dev/null
echo "  User enrolled in Sheiko Beginner"

echo ""
echo "Step 9: Generating workout (Week 1, Day 1)..."
echo ""

WORKOUT=$(curl -s -X GET "$API_BASE_URL/users/$TEST_USER_ID/workout" \
    -H "X-User-ID: $TEST_USER_ID")

echo "=== Generated Workout (Week 1, Day 1 - High Volume) ==="
echo "$WORKOUT"

echo ""
echo "=== Sheiko Beginner Configuration Complete ==="
echo ""
echo "Program ID: $PROGRAM_ID"
echo "Test User: $TEST_USER_ID"
echo ""
echo "Key features demonstrated:"
echo "  - High-volume RAMP prescriptions (20-30+ reps per lift)"
echo "  - Uses ONE_RM (not TRAINING_MAX)"
echo "  - NO automatic progressions (Sheiko philosophy)"
echo "  - Lift variations with parent relationships"
echo ""
echo "Manual progression example:"
echo "  When the athlete is ready to increase maxes, call:"
echo "  curl -X POST '\$API_BASE_URL/users/$TEST_USER_ID/lift-maxes' \\"
echo "    -H 'X-User-ID: $TEST_USER_ID' \\"
echo "    -H 'Content-Type: application/json' \\"
echo "    -d '{\"liftId\":\"$SQUAT_ID\",\"type\":\"ONE_RM\",\"value\":185.0,...}'"
