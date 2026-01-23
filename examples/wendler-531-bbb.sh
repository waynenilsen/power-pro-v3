#!/bin/bash
# Wendler 5/3/1 BBB Program Configuration Example
# Demonstrates: Weekly lookups, AMRAP sets, cycle-end progression
#
# Usage: ./examples/wendler-531-bbb.sh
# Prerequisites: PowerPro server running on localhost:8080

set -e

API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
ADMIN_USER_ID="${ADMIN_USER_ID:-admin-user}"
TEST_USER_ID="${TEST_USER_ID:-test-user-001}"

echo "=== Wendler 5/3/1 BBB Program Configuration ==="
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

PRESS_ID=$(api POST "/lifts" '{"name":"Overhead Press","slug":"overhead-press","isCompetitionLift":false}' | extract_id)
echo "  Overhead Press: $PRESS_ID"

echo ""
echo "Step 2: Creating weekly lookup (5/3/1 percentages)..."

# Weekly lookup defines the percentages for each week of the cycle
WEEKLY_LOOKUP_ID=$(api POST "/weekly-lookups" '{
    "name":"Wendler 5/3/1 Weekly Percentages",
    "entries":[
        {"weekNumber":1,"set1Percentage":65.0,"set2Percentage":75.0,"set3Percentage":85.0},
        {"weekNumber":2,"set1Percentage":70.0,"set2Percentage":80.0,"set3Percentage":90.0},
        {"weekNumber":3,"set1Percentage":75.0,"set2Percentage":85.0,"set3Percentage":95.0},
        {"weekNumber":4,"set1Percentage":40.0,"set2Percentage":50.0,"set3Percentage":60.0}
    ]
}' | extract_id)
echo "  Weekly Lookup: $WEEKLY_LOOKUP_ID"

echo ""
echo "Step 3: Creating prescriptions..."

# Main work sets for 5/3/1 (uses RAMP scheme that references weekly lookup)
# Week 1: 65x5, 75x5, 85x5+ (AMRAP)
# Week 2: 70x3, 80x3, 90x3+ (AMRAP)
# Week 3: 75x5, 85x3, 95x1+ (AMRAP)
# Week 4: 40x5, 50x5, 60x5 (Deload)

RX_SQUAT_MAIN=$(api POST "/prescriptions" '{
    "liftId":"'"$SQUAT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"week"},
    "setScheme":{"type":"RAMP","steps":[
        {"lookupField":"set1Percentage","reps":5},
        {"lookupField":"set2Percentage","reps":5},
        {"lookupField":"set3Percentage","reps":5,"isAmrap":true}
    ],"workSetThreshold":60},
    "order":1,
    "notes":"5/3/1 main work - AMRAP on last set",
    "restSeconds":180
}' | extract_id)
echo "  Squat Main: $RX_SQUAT_MAIN"

RX_SQUAT_BBB=$(api POST "/prescriptions" '{
    "liftId":"'"$SQUAT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":50.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"FIXED","sets":5,"reps":10,"isAmrap":false},
    "order":2,
    "notes":"BBB supplemental - 5x10 at 50%",
    "restSeconds":90
}' | extract_id)
echo "  Squat BBB: $RX_SQUAT_BBB"

RX_BENCH_MAIN=$(api POST "/prescriptions" '{
    "liftId":"'"$BENCH_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"week"},
    "setScheme":{"type":"RAMP","steps":[
        {"lookupField":"set1Percentage","reps":5},
        {"lookupField":"set2Percentage","reps":5},
        {"lookupField":"set3Percentage","reps":5,"isAmrap":true}
    ],"workSetThreshold":60},
    "order":1,
    "notes":"5/3/1 main work - AMRAP on last set",
    "restSeconds":180
}' | extract_id)
echo "  Bench Main: $RX_BENCH_MAIN"

RX_BENCH_BBB=$(api POST "/prescriptions" '{
    "liftId":"'"$BENCH_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":50.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"FIXED","sets":5,"reps":10,"isAmrap":false},
    "order":2,
    "notes":"BBB supplemental - 5x10 at 50%",
    "restSeconds":90
}' | extract_id)
echo "  Bench BBB: $RX_BENCH_BBB"

RX_DEADLIFT_MAIN=$(api POST "/prescriptions" '{
    "liftId":"'"$DEADLIFT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"week"},
    "setScheme":{"type":"RAMP","steps":[
        {"lookupField":"set1Percentage","reps":5},
        {"lookupField":"set2Percentage","reps":5},
        {"lookupField":"set3Percentage","reps":5,"isAmrap":true}
    ],"workSetThreshold":60},
    "order":1,
    "notes":"5/3/1 main work - AMRAP on last set",
    "restSeconds":180
}' | extract_id)
echo "  Deadlift Main: $RX_DEADLIFT_MAIN"

RX_DEADLIFT_BBB=$(api POST "/prescriptions" '{
    "liftId":"'"$DEADLIFT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":50.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"FIXED","sets":5,"reps":10,"isAmrap":false},
    "order":2,
    "notes":"BBB supplemental - 5x10 at 50%",
    "restSeconds":90
}' | extract_id)
echo "  Deadlift BBB: $RX_DEADLIFT_BBB"

RX_PRESS_MAIN=$(api POST "/prescriptions" '{
    "liftId":"'"$PRESS_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"week"},
    "setScheme":{"type":"RAMP","steps":[
        {"lookupField":"set1Percentage","reps":5},
        {"lookupField":"set2Percentage","reps":5},
        {"lookupField":"set3Percentage","reps":5,"isAmrap":true}
    ],"workSetThreshold":60},
    "order":1,
    "notes":"5/3/1 main work - AMRAP on last set",
    "restSeconds":180
}' | extract_id)
echo "  Press Main: $RX_PRESS_MAIN"

RX_PRESS_BBB=$(api POST "/prescriptions" '{
    "liftId":"'"$PRESS_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":50.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"FIXED","sets":5,"reps":10,"isAmrap":false},
    "order":2,
    "notes":"BBB supplemental - 5x10 at 50%",
    "restSeconds":90
}' | extract_id)
echo "  Press BBB: $RX_PRESS_BBB"

echo ""
echo "Step 4: Creating days..."

DAY_SQUAT=$(api POST "/days" '{"name":"Squat Day","slug":"531-squat","metadata":{"focus":"squat","program":"wendler-531-bbb"}}' | extract_id)
echo "  Squat Day: $DAY_SQUAT"

DAY_BENCH=$(api POST "/days" '{"name":"Bench Day","slug":"531-bench","metadata":{"focus":"bench","program":"wendler-531-bbb"}}' | extract_id)
echo "  Bench Day: $DAY_BENCH"

DAY_DEADLIFT=$(api POST "/days" '{"name":"Deadlift Day","slug":"531-deadlift","metadata":{"focus":"deadlift","program":"wendler-531-bbb"}}' | extract_id)
echo "  Deadlift Day: $DAY_DEADLIFT"

DAY_PRESS=$(api POST "/days" '{"name":"Press Day","slug":"531-press","metadata":{"focus":"press","program":"wendler-531-bbb"}}' | extract_id)
echo "  Press Day: $DAY_PRESS"

# Add prescriptions to days
api POST "/days/$DAY_SQUAT/prescriptions" '{"prescriptionId":"'"$RX_SQUAT_MAIN"'","order":1}' > /dev/null
api POST "/days/$DAY_SQUAT/prescriptions" '{"prescriptionId":"'"$RX_SQUAT_BBB"'","order":2}' > /dev/null
echo "  Added prescriptions to Squat Day"

api POST "/days/$DAY_BENCH/prescriptions" '{"prescriptionId":"'"$RX_BENCH_MAIN"'","order":1}' > /dev/null
api POST "/days/$DAY_BENCH/prescriptions" '{"prescriptionId":"'"$RX_BENCH_BBB"'","order":2}' > /dev/null
echo "  Added prescriptions to Bench Day"

api POST "/days/$DAY_DEADLIFT/prescriptions" '{"prescriptionId":"'"$RX_DEADLIFT_MAIN"'","order":1}' > /dev/null
api POST "/days/$DAY_DEADLIFT/prescriptions" '{"prescriptionId":"'"$RX_DEADLIFT_BBB"'","order":2}' > /dev/null
echo "  Added prescriptions to Deadlift Day"

api POST "/days/$DAY_PRESS/prescriptions" '{"prescriptionId":"'"$RX_PRESS_MAIN"'","order":1}' > /dev/null
api POST "/days/$DAY_PRESS/prescriptions" '{"prescriptionId":"'"$RX_PRESS_BBB"'","order":2}' > /dev/null
echo "  Added prescriptions to Press Day"

echo ""
echo "Step 5: Creating cycle..."

CYCLE_ID=$(api POST "/cycles" '{"name":"Wendler 5/3/1 4-Week Cycle","lengthWeeks":4}' | extract_id)
echo "  Cycle: $CYCLE_ID"

echo ""
echo "Step 6: Creating weeks..."

# Week 1 (5s week)
WEEK_1=$(api POST "/weeks" '{"cycleId":"'"$CYCLE_ID"'","weekNumber":1,"name":"Week 1 (5s)"}' | extract_id)
echo "  Week 1 (5s): $WEEK_1"
api POST "/weeks/$WEEK_1/days" '{"dayId":"'"$DAY_SQUAT"'","position":0}' > /dev/null
api POST "/weeks/$WEEK_1/days" '{"dayId":"'"$DAY_BENCH"'","position":1}' > /dev/null
api POST "/weeks/$WEEK_1/days" '{"dayId":"'"$DAY_DEADLIFT"'","position":2}' > /dev/null
api POST "/weeks/$WEEK_1/days" '{"dayId":"'"$DAY_PRESS"'","position":3}' > /dev/null

# Week 2 (3s week)
WEEK_2=$(api POST "/weeks" '{"cycleId":"'"$CYCLE_ID"'","weekNumber":2,"name":"Week 2 (3s)"}' | extract_id)
echo "  Week 2 (3s): $WEEK_2"
api POST "/weeks/$WEEK_2/days" '{"dayId":"'"$DAY_SQUAT"'","position":0}' > /dev/null
api POST "/weeks/$WEEK_2/days" '{"dayId":"'"$DAY_BENCH"'","position":1}' > /dev/null
api POST "/weeks/$WEEK_2/days" '{"dayId":"'"$DAY_DEADLIFT"'","position":2}' > /dev/null
api POST "/weeks/$WEEK_2/days" '{"dayId":"'"$DAY_PRESS"'","position":3}' > /dev/null

# Week 3 (5/3/1 week)
WEEK_3=$(api POST "/weeks" '{"cycleId":"'"$CYCLE_ID"'","weekNumber":3,"name":"Week 3 (5/3/1)"}' | extract_id)
echo "  Week 3 (5/3/1): $WEEK_3"
api POST "/weeks/$WEEK_3/days" '{"dayId":"'"$DAY_SQUAT"'","position":0}' > /dev/null
api POST "/weeks/$WEEK_3/days" '{"dayId":"'"$DAY_BENCH"'","position":1}' > /dev/null
api POST "/weeks/$WEEK_3/days" '{"dayId":"'"$DAY_DEADLIFT"'","position":2}' > /dev/null
api POST "/weeks/$WEEK_3/days" '{"dayId":"'"$DAY_PRESS"'","position":3}' > /dev/null

# Week 4 (Deload)
WEEK_4=$(api POST "/weeks" '{"cycleId":"'"$CYCLE_ID"'","weekNumber":4,"name":"Week 4 (Deload)"}' | extract_id)
echo "  Week 4 (Deload): $WEEK_4"
api POST "/weeks/$WEEK_4/days" '{"dayId":"'"$DAY_SQUAT"'","position":0}' > /dev/null
api POST "/weeks/$WEEK_4/days" '{"dayId":"'"$DAY_BENCH"'","position":1}' > /dev/null
api POST "/weeks/$WEEK_4/days" '{"dayId":"'"$DAY_DEADLIFT"'","position":2}' > /dev/null
api POST "/weeks/$WEEK_4/days" '{"dayId":"'"$DAY_PRESS"'","position":3}' > /dev/null

echo "  Added days to all weeks"

echo ""
echo "Step 7: Creating progressions..."

# Lower body: +10lb per cycle
PROG_LOWER=$(api POST "/progressions" '{
    "name":"5/3/1 Lower Body +10lb",
    "type":"CYCLE_PROGRESSION",
    "parameters":{"increment":10.0,"maxType":"TRAINING_MAX"}
}' | extract_id)
echo "  Lower Body +10lb: $PROG_LOWER"

# Upper body: +5lb per cycle
PROG_UPPER=$(api POST "/progressions" '{
    "name":"5/3/1 Upper Body +5lb",
    "type":"CYCLE_PROGRESSION",
    "parameters":{"increment":5.0,"maxType":"TRAINING_MAX"}
}' | extract_id)
echo "  Upper Body +5lb: $PROG_UPPER"

echo ""
echo "Step 8: Creating program..."

PROGRAM_ID=$(api POST "/programs" '{
    "name":"Wendler 5/3/1 BBB",
    "slug":"wendler-531-bbb",
    "description":"Jim Wendler'\''s 5/3/1 with Boring But Big supplemental work. 4-week mesocycles with AMRAP sets and cycle-end progression.",
    "cycleId":"'"$CYCLE_ID"'",
    "weeklyLookupId":"'"$WEEKLY_LOOKUP_ID"'",
    "defaultRounding":5.0
}' | extract_id)
echo "  Program: $PROGRAM_ID"

echo ""
echo "Step 9: Linking progressions to program..."

api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_LOWER"'","liftId":"'"$SQUAT_ID"'","priority":1}' > /dev/null
echo "  Squat -> +10lb per cycle"
api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_LOWER"'","liftId":"'"$DEADLIFT_ID"'","priority":1}' > /dev/null
echo "  Deadlift -> +10lb per cycle"
api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_UPPER"'","liftId":"'"$BENCH_ID"'","priority":1}' > /dev/null
echo "  Bench -> +5lb per cycle"
api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_UPPER"'","liftId":"'"$PRESS_ID"'","priority":1}' > /dev/null
echo "  Press -> +5lb per cycle"

echo ""
echo "Step 10: Setting up test user..."

# Create training maxes (typically 90% of actual 1RM for 5/3/1)
curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$SQUAT_ID"'","type":"TRAINING_MAX","value":315.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Squat TM: 315 lbs (90% of 350 1RM)"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$BENCH_ID"'","type":"TRAINING_MAX","value":225.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Bench TM: 225 lbs (90% of 250 1RM)"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$DEADLIFT_ID"'","type":"TRAINING_MAX","value":360.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Deadlift TM: 360 lbs (90% of 400 1RM)"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$PRESS_ID"'","type":"TRAINING_MAX","value":135.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Press TM: 135 lbs (90% of 150 1RM)"

echo ""
echo "Step 11: Enrolling user in program..."

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/program" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"programId":"'"$PROGRAM_ID"'"}' > /dev/null
echo "  User enrolled in Wendler 5/3/1 BBB"

echo ""
echo "Step 12: Generating workout (Week 1, Squat Day)..."
echo ""

WORKOUT=$(curl -s -X GET "$API_BASE_URL/users/$TEST_USER_ID/workout" \
    -H "X-User-ID: $TEST_USER_ID")

echo "=== Generated Workout (Week 1 - 5s Week, Squat Day) ==="
echo "$WORKOUT"

echo ""
echo "=== Wendler 5/3/1 BBB Configuration Complete ==="
echo ""
echo "Program ID: $PROGRAM_ID"
echo "Test User: $TEST_USER_ID"
echo ""
echo "Key features demonstrated:"
echo "  - Weekly lookup tables (5/3/1 percentages per week)"
echo "  - AMRAP final sets"
echo "  - 4-week mesocycle"
echo "  - Cycle-end progression (not per session)"
echo "  - Different increments for upper/lower body"
