# Test: 006 Texas Method

## Task
Create an E2E test for the Texas Method program in `internal/api/e2e/texas_method_test.go`

## Program Characteristics
- **3-day weekly cycle**: Monday (Volume), Wednesday (Recovery), Friday (Intensity)
- **DailyLookup pattern**: Mon=90% of Fri, Wed=80% of Mon
- **Linear Progression**: Weekly progression based on Friday PR attempts
- **Different rep schemes by day**: Mon 5x5, Wed 2x5, Fri 1x5

## Key Features to Test
1. **Daily intensity variation**: Mon 90%, Wed ~72%, Fri 100%
2. **Volume Day structure**: 5x5 at 90% of intensity day
3. **Recovery Day structure**: 2x5 at 80% of volume day
4. **Intensity Day structure**: 1x5 PR attempt
5. **Weekly progression**: +5lb lower body, +2.5lb upper body after successful Friday
6. **Stall handling** (if implemented)

## Test Template
Similar to `bill_starr_test.go` which also uses daily intensity patterns:
- Create DailyLookup for volume/recovery/intensity percentages
- Create prescriptions for each day
- Create 1-week cycle with Mon/Wed/Fri
- Test weight calculations at different daily intensities
- Test progression after completing the week
