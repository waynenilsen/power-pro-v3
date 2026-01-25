# Test: 014 Inverted Juggernaut 5/3/1

## Task
Create an E2E test for the Inverted Juggernaut 5/3/1 program in `internal/api/e2e/inverted_juggernaut_test.go`

## Program Characteristics
- **16-week cycle**: Four 4-week waves (10s, 8s, 5s, 3s)
- **Wave periodization**: Each wave has Accumulation, Intensification, Realization, Deload
- **5/3/1 percentage structure**: FS (65%), SS (75%), Top (85%) + AMRAP
- **Training Max = 90% of 1RM**
- **Back-off sets**: Volume work after top set

## Key Features to Test
1. **Rep wave progression**: 10s -> 8s -> 5s -> 3s across 16 weeks
2. **4-week phase structure within each wave**: Accum, Intens, Real, Deload
3. **5/3/1 percentages**: 65%, 75%, 85% within each session
4. **AMRAP on top set**: Verify AMRAP prescription
5. **Back-off volume**: Sets after AMRAP
6. **Deload weeks**: Reduced volume in weeks 4, 8, 12, 16

## Test Template
Complex multi-wave program similar to `wendler_531_test.go`:
- Create WeeklyLookup for 16-week progression
- Create prescriptions with 5/3/1 structure
- Create 16-week cycle with 4 waves
- Test week 1 vs week 5 vs week 9 vs week 13 (different waves)
- Test deload week volume reduction
- Test CycleProgression after 16 weeks
