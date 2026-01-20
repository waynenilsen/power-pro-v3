# 010: Prescription Resolution API Endpoints

## ERD Reference
Implements: REQ-PRSC-008, REQ-PRSC-009, REQ-PRSC-010
Related to: NFR-001, NFR-002, NFR-004, NFR-005

## Description
Implement the API endpoints for resolving prescriptions to concrete sets. This includes single prescription resolution and batch resolution with caching optimization.

## Context / Background
API consumers need the actual workout instructions, not abstract prescriptions. Resolution takes a prescription and a user, fetches the user's maxes, applies the LoadStrategy to calculate weights, and uses the SetScheme to generate the concrete sets.

## Acceptance Criteria
- [ ] Implement `POST /prescriptions/{id}/resolve` endpoint:
  - Input: userId in request body (for max lookup)
  - Calls prescription.Resolve() with user context
  - Returns resolved prescription with concrete sets
  - Returns 404 if prescription not found
  - Returns 422 if max not found for user/lift
  - Performance: < 100ms p95 (NFR-001)
- [ ] Implement `POST /prescriptions/resolve-batch` endpoint:
  - Input: array of prescription IDs + userId
  - Resolves all prescriptions for the user
  - Returns array of resolved prescriptions
  - Returns partial results on partial failure (NFR-005)
  - Failed items include error message
  - Performance: < 500ms p95 for 20 prescriptions (NFR-002)
- [ ] Implement max lookup caching (REQ-PRSC-010):
  - Single max lookup per (user, lift, type) combination per batch
  - Cache lives only for duration of batch request
  - Reduces database queries for workouts with multiple exercises using same lift
- [ ] Graceful failure handling (NFR-004, NFR-005):
  - Clear error message if max not found
  - Batch continues processing after individual failures
  - Response indicates which items succeeded/failed
- [ ] Test coverage > 90% including:
  - Single resolution happy path
  - Batch resolution happy path
  - Max not found error handling
  - Partial batch failure
  - Cache hit verification (no duplicate max queries)

## Technical Notes
- Use domain layer Resolve() method from ticket 008
- Max caching can use simple map keyed by `(userID, liftID, refType)`
- Batch endpoint should process concurrently where possible
- Consider returning lift info in resolved output for display purposes

## Resolved Prescription Response Format
```json
{
  "prescriptionId": "uuid",
  "lift": {
    "id": "uuid",
    "name": "Squat",
    "slug": "squat"
  },
  "sets": [
    {
      "setNumber": 1,
      "weight": 265.0,
      "targetReps": 5,
      "isWorkSet": true
    }
  ],
  "notes": "Focus on depth",
  "restSeconds": 180
}
```

## Batch Resolution Response Format
```json
{
  "results": [
    {
      "prescriptionId": "uuid",
      "status": "success",
      "resolved": { ... }
    },
    {
      "prescriptionId": "uuid",
      "status": "error",
      "error": "No training max found for Squat"
    }
  ]
}
```

## Dependencies
- Blocks: None
- Blocked by: 008 (Domain logic with Resolve method), 009 (CRUD API for prescription retrieval)
- Related: Sprint 001 (LiftMax repository for max lookup)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/002-prescription-system/erd.md
- NFR-001, NFR-002: Performance requirements
- NFR-004, NFR-005: Reliability requirements
