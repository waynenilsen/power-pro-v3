# 003: Example Responses

## ERD Reference
Implements: REQ-DOC-003
Related to: NFR-005, NFR-006

## Description
Provide example responses for all API endpoints including success responses and common error responses. Response formats must match actual API behavior.

## Context / Background
Developers need to know what to expect from the API in order to parse responses correctly. Documented response examples enable developers to write client code before testing against the actual API.

## Acceptance Criteria
- [ ] Each endpoint has example success response
- [ ] Common error responses documented (400, 401, 404, etc.)
- [ ] Response formats match actual API behavior (NFR-005)
- [ ] Examples produce expected outputs when executed (NFR-006)

## Technical Notes
- Include full JSON response bodies
- Show HTTP status codes
- Include response headers where relevant
- Document pagination response format for list endpoints
- Error responses should show:
  - 400 Bad Request (validation errors)
  - 401 Unauthorized (missing/invalid auth)
  - 404 Not Found (resource doesn't exist)
  - 500 Internal Server Error (general format)
- Success responses should show:
  - 200 OK (read/update operations)
  - 201 Created (create operations)
  - 204 No Content (delete operations, if applicable)

## Dependencies
- Blocks: 005-workflow-documentation
- Blocked by: 001-endpoint-documentation
- Related: 002-example-requests, 004-error-documentation

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/006-api-developer-experience/erd.md
