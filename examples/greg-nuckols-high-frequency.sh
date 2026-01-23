#!/bin/bash
# Greg Nuckols High Frequency Program Configuration Example
# Demonstrates: Daily undulating periodization, multi-week cycles, combined lookups
#
# Usage: ./examples/greg-nuckols-high-frequency.sh
# Prerequisites: PowerPro server running on localhost:8080

set -e

API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
ADMIN_USER_ID="${ADMIN_USER_ID:-admin-user}"
TEST_USER_ID="${TEST_USER_ID:-test-user-001}"

echo "=== Greg Nuckols High Frequency Program Configuration ==="
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

BENCH_ID=$(api POST "/lifts" '{"name":"Bench Press","slug":"bench-press","isCompetitionLift":true}' | extract_id)
echo "  Bench Press: $BENCH_ID"

SQUAT_ID=$(api POST "/lifts" '{"name":"Squat","slug":"squat","isCompetitionLift":true}' | extract_id)
echo "  Squat: $SQUAT_ID"

DEADLIFT_ID=$(api POST "/lifts" '{"name":"Deadlift","slug":"deadlift","isCompetitionLift":true}' | extract_id)
echo "  Deadlift: $DEADLIFT_ID"

PRESS_ID=$(api POST "/lifts" '{"name":"Overhead Press","slug":"overhead-press","isCompetitionLift":false}' | extract_id)
echo "  Overhead Press: $PRESS_ID"

ROW_ID=$(api POST "/lifts" '{"name":"T-Bar Row","slug":"t-bar-row","isCompetitionLift":false}' | extract_id)
echo "  T-Bar Row: $ROW_ID"

echo ""
echo "Step 2: Creating daily lookup (intensity by day)..."

DAILY_LOOKUP_ID=$(api POST "/daily-lookups" '{
    "name":"Nuckols Daily Intensity",
    "entries":[
        {"dayIdentifier":"monday","percentageModifier":75.0,"intensityLevel":"MODERATE"},
        {"dayIdentifier":"tuesday","percentageModifier":80.0,"intensityLevel":"HEAVY"},
        {"dayIdentifier":"wednesday","percentageModifier":70.0,"intensityLevel":"LIGHT"},
        {"dayIdentifier":"thursday","percentageModifier":85.0,"intensityLevel":"HEAVY"},
        {"dayIdentifier":"friday","percentageModifier":85.0,"intensityLevel":"HEAVY"},
        {"dayIdentifier":"saturday","percentageModifier":65.0,"intensityLevel":"LIGHT"}
    ]
}' | extract_id)
echo "  Daily Lookup: $DAILY_LOOKUP_ID"

echo ""
echo "Step 3: Creating weekly lookup (volume progression)..."

# Weekly lookup tracks sets/reps progression across 3 weeks
WEEKLY_LOOKUP_ID=$(api POST "/weekly-lookups" '{
    "name":"Nuckols 3-Week Volume Progression",
    "entries":[
        {"weekNumber":1,"set1Percentage":100.0,"set2Percentage":100.0,"set3Percentage":100.0},
        {"weekNumber":2,"set1Percentage":100.0,"set2Percentage":100.0,"set3Percentage":100.0},
        {"weekNumber":3,"set1Percentage":100.0,"set2Percentage":100.0,"set3Percentage":100.0}
    ]
}' | extract_id)
echo "  Weekly Lookup: $WEEKLY_LOOKUP_ID"

echo ""
echo "Step 4: Creating prescriptions..."

# Monday: Bench + Squat at 75% (4x3 wk1, 5x3 wk2, 5x4 wk3)
RX_BENCH_MON=$(api POST "/prescriptions" '{
    "liftId":"'"$BENCH_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"day"},
    "setScheme":{"type":"FIXED","sets":4,"reps":3,"isAmrap":false},
    "order":1,
    "notes":"Moderate day - focus on bar speed",
    "restSeconds":150
}' | extract_id)
echo "  Bench Mon (75% 4x3): $RX_BENCH_MON"

RX_SQUAT_MON=$(api POST "/prescriptions" '{
    "liftId":"'"$SQUAT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"day"},
    "setScheme":{"type":"FIXED","sets":4,"reps":3,"isAmrap":false},
    "order":2,
    "notes":"Technical practice",
    "restSeconds":180
}' | extract_id)
echo "  Squat Mon (75% 4x3): $RX_SQUAT_MON"

# Tuesday: Bench + Squat at 80% (3x2 wk1, 4x2 wk2, 4x3 wk3)
RX_BENCH_TUE=$(api POST "/prescriptions" '{
    "liftId":"'"$BENCH_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"day"},
    "setScheme":{"type":"FIXED","sets":3,"reps":2,"isAmrap":false},
    "order":1,
    "notes":"Strength building day",
    "restSeconds":180
}' | extract_id)
echo "  Bench Tue (80% 3x2): $RX_BENCH_TUE"

RX_SQUAT_TUE=$(api POST "/prescriptions" '{
    "liftId":"'"$SQUAT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"day"},
    "setScheme":{"type":"FIXED","sets":3,"reps":2,"isAmrap":false},
    "order":2,
    "notes":"Heavy doubles",
    "restSeconds":180
}' | extract_id)
echo "  Squat Tue (80% 3x2): $RX_SQUAT_TUE"

RX_ROW_TUE=$(api POST "/prescriptions" '{
    "liftId":"'"$ROW_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"FIXED","sets":3,"reps":8,"isAmrap":false},
    "order":3,
    "notes":"Back work at 6RM weight",
    "restSeconds":90
}' | extract_id)
echo "  Row Tue (6RM 3x8): $RX_ROW_TUE"

# Wednesday: Bench + Squat at 70% (4x4 wk1, 5x4 wk2, 5x5 wk3)
RX_BENCH_WED=$(api POST "/prescriptions" '{
    "liftId":"'"$BENCH_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"day"},
    "setScheme":{"type":"FIXED","sets":4,"reps":4,"isAmrap":false},
    "order":1,
    "notes":"Light recovery day",
    "restSeconds":120
}' | extract_id)
echo "  Bench Wed (70% 4x4): $RX_BENCH_WED"

RX_SQUAT_WED=$(api POST "/prescriptions" '{
    "liftId":"'"$SQUAT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"day"},
    "setScheme":{"type":"FIXED","sets":4,"reps":4,"isAmrap":false},
    "order":2,
    "notes":"Volume accumulation",
    "restSeconds":150
}' | extract_id)
echo "  Squat Wed (70% 4x4): $RX_SQUAT_WED"

# Thursday: Bench + Squat at 85% (3x1 wk1, 4x1 wk2, AMAP wk3)
RX_BENCH_THU=$(api POST "/prescriptions" '{
    "liftId":"'"$BENCH_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"day"},
    "setScheme":{"type":"FIXED","sets":3,"reps":1,"isAmrap":false},
    "order":1,
    "notes":"Heavy singles - intensity day",
    "restSeconds":240
}' | extract_id)
echo "  Bench Thu (85% 3x1): $RX_BENCH_THU"

RX_SQUAT_THU=$(api POST "/prescriptions" '{
    "liftId":"'"$SQUAT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"day"},
    "setScheme":{"type":"FIXED","sets":3,"reps":1,"isAmrap":false},
    "order":2,
    "notes":"Heavy singles",
    "restSeconds":240
}' | extract_id)
echo "  Squat Thu (85% 3x1): $RX_SQUAT_THU"

# Friday: Deadlift + OHP at 85% (3x1 wk1, 4x1 wk2, AMAP wk3)
RX_DEADLIFT=$(api POST "/prescriptions" '{
    "liftId":"'"$DEADLIFT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":85.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"FIXED","sets":3,"reps":1,"isAmrap":false},
    "order":1,
    "notes":"Heavy deadlift singles",
    "restSeconds":300
}' | extract_id)
echo "  Deadlift Fri (85% 3x1): $RX_DEADLIFT"

RX_OHP=$(api POST "/prescriptions" '{
    "liftId":"'"$PRESS_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":85.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"FIXED","sets":3,"reps":1,"isAmrap":false},
    "order":2,
    "notes":"Heavy overhead work",
    "restSeconds":180
}' | extract_id)
echo "  OHP Fri (85% 3x1): $RX_OHP"

RX_ROW_FRI=$(api POST "/prescriptions" '{
    "liftId":"'"$ROW_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST"},
    "setScheme":{"type":"FIXED","sets":3,"reps":8,"isAmrap":false},
    "order":3,
    "notes":"Back work",
    "restSeconds":90
}' | extract_id)
echo "  Row Fri (6RM 3x8): $RX_ROW_FRI"

# Saturday: Bench + Squat at 65% (5x5 wk1, 6x5 wk2, 6x6 wk3)
RX_BENCH_SAT=$(api POST "/prescriptions" '{
    "liftId":"'"$BENCH_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"day"},
    "setScheme":{"type":"FIXED","sets":5,"reps":5,"isAmrap":false},
    "order":1,
    "notes":"High volume hypertrophy",
    "restSeconds":90
}' | extract_id)
echo "  Bench Sat (65% 5x5): $RX_BENCH_SAT"

RX_SQUAT_SAT=$(api POST "/prescriptions" '{
    "liftId":"'"$SQUAT_ID"'",
    "loadStrategy":{"type":"PERCENT_OF","referenceType":"TRAINING_MAX","percentage":100.0,"roundingIncrement":5.0,"roundingDirection":"NEAREST","lookupKey":"day"},
    "setScheme":{"type":"FIXED","sets":5,"reps":5,"isAmrap":false},
    "order":2,
    "notes":"Volume accumulation",
    "restSeconds":120
}' | extract_id)
echo "  Squat Sat (65% 5x5): $RX_SQUAT_SAT"

echo ""
echo "Step 5: Creating days..."

DAY_MON=$(api POST "/days" '{"name":"Monday","slug":"nuckols-mon","metadata":{"focus":"bench-squat","intensity":"moderate","program":"greg-nuckols"}}' | extract_id)
echo "  Monday: $DAY_MON"

DAY_TUE=$(api POST "/days" '{"name":"Tuesday","slug":"nuckols-tue","metadata":{"focus":"bench-squat-row","intensity":"heavy","program":"greg-nuckols"}}' | extract_id)
echo "  Tuesday: $DAY_TUE"

DAY_WED=$(api POST "/days" '{"name":"Wednesday","slug":"nuckols-wed","metadata":{"focus":"bench-squat","intensity":"light","program":"greg-nuckols"}}' | extract_id)
echo "  Wednesday: $DAY_WED"

DAY_THU=$(api POST "/days" '{"name":"Thursday","slug":"nuckols-thu","metadata":{"focus":"bench-squat","intensity":"heavy","program":"greg-nuckols"}}' | extract_id)
echo "  Thursday: $DAY_THU"

DAY_FRI=$(api POST "/days" '{"name":"Friday","slug":"nuckols-fri","metadata":{"focus":"deadlift-ohp-row","intensity":"heavy","program":"greg-nuckols"}}' | extract_id)
echo "  Friday: $DAY_FRI"

DAY_SAT=$(api POST "/days" '{"name":"Saturday","slug":"nuckols-sat","metadata":{"focus":"bench-squat","intensity":"light","program":"greg-nuckols"}}' | extract_id)
echo "  Saturday: $DAY_SAT"

# Add prescriptions to days
api POST "/days/$DAY_MON/prescriptions" '{"prescriptionId":"'"$RX_BENCH_MON"'","order":1}' > /dev/null
api POST "/days/$DAY_MON/prescriptions" '{"prescriptionId":"'"$RX_SQUAT_MON"'","order":2}' > /dev/null
echo "  Added prescriptions to Monday"

api POST "/days/$DAY_TUE/prescriptions" '{"prescriptionId":"'"$RX_BENCH_TUE"'","order":1}' > /dev/null
api POST "/days/$DAY_TUE/prescriptions" '{"prescriptionId":"'"$RX_SQUAT_TUE"'","order":2}' > /dev/null
api POST "/days/$DAY_TUE/prescriptions" '{"prescriptionId":"'"$RX_ROW_TUE"'","order":3}' > /dev/null
echo "  Added prescriptions to Tuesday"

api POST "/days/$DAY_WED/prescriptions" '{"prescriptionId":"'"$RX_BENCH_WED"'","order":1}' > /dev/null
api POST "/days/$DAY_WED/prescriptions" '{"prescriptionId":"'"$RX_SQUAT_WED"'","order":2}' > /dev/null
echo "  Added prescriptions to Wednesday"

api POST "/days/$DAY_THU/prescriptions" '{"prescriptionId":"'"$RX_BENCH_THU"'","order":1}' > /dev/null
api POST "/days/$DAY_THU/prescriptions" '{"prescriptionId":"'"$RX_SQUAT_THU"'","order":2}' > /dev/null
echo "  Added prescriptions to Thursday"

api POST "/days/$DAY_FRI/prescriptions" '{"prescriptionId":"'"$RX_DEADLIFT"'","order":1}' > /dev/null
api POST "/days/$DAY_FRI/prescriptions" '{"prescriptionId":"'"$RX_OHP"'","order":2}' > /dev/null
api POST "/days/$DAY_FRI/prescriptions" '{"prescriptionId":"'"$RX_ROW_FRI"'","order":3}' > /dev/null
echo "  Added prescriptions to Friday"

api POST "/days/$DAY_SAT/prescriptions" '{"prescriptionId":"'"$RX_BENCH_SAT"'","order":1}' > /dev/null
api POST "/days/$DAY_SAT/prescriptions" '{"prescriptionId":"'"$RX_SQUAT_SAT"'","order":2}' > /dev/null
echo "  Added prescriptions to Saturday"

echo ""
echo "Step 6: Creating cycle..."

CYCLE_ID=$(api POST "/cycles" '{"name":"Nuckols 3-Week High Frequency Cycle","lengthWeeks":3}' | extract_id)
echo "  Cycle: $CYCLE_ID"

echo ""
echo "Step 7: Creating weeks..."

for week_num in 1 2 3; do
    WEEK_ID=$(api POST "/weeks" '{"cycleId":"'"$CYCLE_ID"'","weekNumber":'"$week_num"',"name":"Week '"$week_num"'"}' | extract_id)
    echo "  Week $week_num: $WEEK_ID"
    api POST "/weeks/$WEEK_ID/days" '{"dayId":"'"$DAY_MON"'","position":0}' > /dev/null
    api POST "/weeks/$WEEK_ID/days" '{"dayId":"'"$DAY_TUE"'","position":1}' > /dev/null
    api POST "/weeks/$WEEK_ID/days" '{"dayId":"'"$DAY_WED"'","position":2}' > /dev/null
    api POST "/weeks/$WEEK_ID/days" '{"dayId":"'"$DAY_THU"'","position":3}' > /dev/null
    api POST "/weeks/$WEEK_ID/days" '{"dayId":"'"$DAY_FRI"'","position":4}' > /dev/null
    api POST "/weeks/$WEEK_ID/days" '{"dayId":"'"$DAY_SAT"'","position":5}' > /dev/null
done
echo "  Added days to all weeks (6 days per week)"

echo ""
echo "Step 8: Creating progressions..."

# Main lifts: cycle-end progression
PROG_MAIN=$(api POST "/progressions" '{
    "name":"Nuckols Cycle Progression",
    "type":"CYCLE_PROGRESSION",
    "parameters":{"increment":5.0,"maxType":"TRAINING_MAX"}
}' | extract_id)
echo "  Main Lifts +5lb/cycle: $PROG_MAIN"

# Back work: separate progression
PROG_BACK=$(api POST "/progressions" '{
    "name":"Back Work Progression",
    "type":"CYCLE_PROGRESSION",
    "parameters":{"increment":5.0,"maxType":"TRAINING_MAX"}
}' | extract_id)
echo "  Back Work +5lb/cycle: $PROG_BACK"

echo ""
echo "Step 9: Creating program..."

PROGRAM_ID=$(api POST "/programs" '{
    "name":"Greg Nuckols High Frequency",
    "slug":"greg-nuckols-high-frequency",
    "description":"High frequency program with daily undulating periodization. Bench and Squat 5x per week with varying intensity. 3-week cycles with progressive volume.",
    "cycleId":"'"$CYCLE_ID"'",
    "dailyLookupId":"'"$DAILY_LOOKUP_ID"'",
    "weeklyLookupId":"'"$WEEKLY_LOOKUP_ID"'",
    "defaultRounding":5.0
}' | extract_id)
echo "  Program: $PROGRAM_ID"

echo ""
echo "Step 10: Linking progressions to program..."

api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_MAIN"'","liftId":"'"$BENCH_ID"'","priority":1}' > /dev/null
echo "  Bench -> +5lb per cycle"
api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_MAIN"'","liftId":"'"$SQUAT_ID"'","priority":1}' > /dev/null
echo "  Squat -> +5lb per cycle"
api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_MAIN"'","liftId":"'"$DEADLIFT_ID"'","priority":1}' > /dev/null
echo "  Deadlift -> +5lb per cycle"
api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_MAIN"'","liftId":"'"$PRESS_ID"'","priority":1}' > /dev/null
echo "  OHP -> +5lb per cycle"
api POST "/programs/$PROGRAM_ID/progressions" '{"progressionId":"'"$PROG_BACK"'","liftId":"'"$ROW_ID"'","priority":1}' > /dev/null
echo "  Row -> +5lb per cycle (6RM)"

echo ""
echo "Step 11: Setting up test user..."

# Create training maxes
curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$BENCH_ID"'","type":"TRAINING_MAX","value":225.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Bench TM: 225 lbs"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$SQUAT_ID"'","type":"TRAINING_MAX","value":315.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Squat TM: 315 lbs"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$DEADLIFT_ID"'","type":"TRAINING_MAX","value":405.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Deadlift TM: 405 lbs"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$PRESS_ID"'","type":"TRAINING_MAX","value":135.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  OHP TM: 135 lbs"

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/lift-maxes" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"liftId":"'"$ROW_ID"'","type":"TRAINING_MAX","value":185.0,"effectiveDate":"2024-01-15T00:00:00Z"}' > /dev/null
echo "  Row 6RM: 185 lbs"

echo ""
echo "Step 12: Enrolling user in program..."

curl -s -X POST "$API_BASE_URL/users/$TEST_USER_ID/program" \
    -H "X-User-ID: $TEST_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"programId":"'"$PROGRAM_ID"'"}' > /dev/null
echo "  User enrolled in Greg Nuckols High Frequency"

echo ""
echo "Step 13: Generating workout (Week 1, Monday)..."
echo ""

WORKOUT=$(curl -s -X GET "$API_BASE_URL/users/$TEST_USER_ID/workout" \
    -H "X-User-ID: $TEST_USER_ID")

echo "=== Generated Workout (Week 1, Monday - Moderate Day 75%) ==="
echo "$WORKOUT"

echo ""
echo "=== Greg Nuckols High Frequency Configuration Complete ==="
echo ""
echo "Program ID: $PROGRAM_ID"
echo "Test User: $TEST_USER_ID"
echo ""
echo "Key features demonstrated:"
echo "  - Daily undulating periodization (65-85% intensity)"
echo "  - High frequency (bench/squat 5x/week)"
echo "  - 3-week progression cycle"
echo "  - Combined daily + weekly lookups"
echo "  - Multiple max types (TRAINING_MAX + TRAINING_MAX)"
echo "  - 6 training days per week"
