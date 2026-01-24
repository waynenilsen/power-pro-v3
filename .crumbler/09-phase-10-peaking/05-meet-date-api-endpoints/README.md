# Meet Date API Endpoints

Implement API endpoints for meet date management.

## Endpoints

### PUT /api/v1/users/{userId}/programs/{programId}/state/meet-date
Set or update meet date for enrolled program.

**Request Body:**
```json
{
  "meet_date": "2024-06-15"  // ISO 8601 date, or null to clear
}
```

**Response:**
```json
{
  "meet_date": "2024-06-15",
  "days_out": 91,
  "current_phase": "prep_1",
  "weeks_to_meet": 13
}
```

### GET /api/v1/users/{userId}/programs/{programId}/state/countdown
Get meet countdown information.

**Response:**
```json
{
  "meet_date": "2024-06-15",
  "days_out": 45,
  "current_phase": "prep_2",
  "phase_week": 3,
  "taper_multiplier": 1.0
}
```

## Tasks

1. Add handler for PUT meet-date endpoint
2. Add handler for GET countdown endpoint
3. Validate meet date is in future
4. Update state service to support meet date operations
5. Add API integration tests

## Acceptance Criteria

- [ ] Can set meet date via API
- [ ] Can clear meet date via API
- [ ] Can retrieve countdown information
- [ ] Invalid dates rejected with appropriate errors
- [ ] API tests verify all endpoints
