# Test: 016 GZCL Method Compendium

## Task
Create an E2E test for the GZCL Method Compendium in `internal/api/e2e/gzcl_compendium_test.go`

## Program Characteristics
- **3-tier system**: T1 (competition lifts), T2 (variations), T3 (accessories)
- **Training Max = 85-90% of 1RM**: Conservative daily 2RM estimate
- **Multiple variants**: VDIP, UHF, Rippler, etc.
- **Intensity ranges**: T1 (80-100%), T2 (60-80%), T3 (RPE-based)
- **Rep ranges**: T1 (1-5 reps), T2 (5-12 reps), T3 (10-20 reps)

## Key Features to Test
1. **3-tier structure**: Verify T1, T2, T3 exercise ordering
2. **Training Max calculations**: Weights based on TM percentages
3. **Volume-dependent progression**: Progress based on total reps achieved
4. **MRS (Max Rep Sets)**: T3 exercises use MRS protocol
5. **Multi-variant support**: Test at least one program variant (VDIP)

## Test Template
Similar to existing `gzclp_t1_test.go` and `gzclp_t2_test.go`:
- Create tiered prescriptions for T1, T2, T3
- Create multi-lift days with all tiers
- Test percentage-based weight calculations
- Test MRS protocol for T3
- Test progression after achieving volume thresholds
