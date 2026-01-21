# 006: RESTful Conventions Audit

## ERD Reference
Implements: REQ-API-001

## Description
Audit and ensure all API endpoints follow RESTful conventions. Consistent conventions make the API predictable and easier to use.

## Context / Background
RESTful APIs follow established conventions that developers expect. Inconsistent use of HTTP methods, status codes, or URL patterns creates confusion and increases integration effort.

## Acceptance Criteria
- [ ] GET used for reads (retrieving resources)
- [ ] POST used for creates (creating new resources)
- [ ] PUT used for updates (replacing resources)
- [ ] DELETE used for deletes (removing resources)
- [ ] Resource-oriented URLs (nouns, not verbs)
- [ ] Proper HTTP status codes used:
  - 200 OK for successful reads/updates
  - 201 Created for successful creates
  - 204 No Content for successful deletes (if applicable)
  - 400 Bad Request for validation errors
  - 401 Unauthorized for auth failures
  - 404 Not Found for missing resources
  - 500 Internal Server Error for server errors
- [ ] Any non-compliant endpoints identified and fixed or documented

## Technical Notes
- Audit process:
  1. List all endpoints in the API
  2. For each endpoint, verify HTTP method matches operation semantics
  3. Verify URL uses nouns (e.g., `/lifts`, `/programs`) not verbs
  4. Verify status codes are appropriate
  5. Document any exceptions with rationale
- Special cases that may not be pure REST:
  - Workout generation (action endpoint)
  - Progression trigger (action endpoint)
  - These may use POST even though they're not "creates" - document rationale
- Create audit report documenting findings and any fixes made

## Dependencies
- Blocks: 007-response-envelope-standardization
- Blocked by: None
- Related: All API documentation tickets

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/006-api-developer-experience/erd.md
