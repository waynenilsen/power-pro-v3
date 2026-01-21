# 002: Example Requests

## ERD Reference
Implements: REQ-DOC-002
Related to: NFR-002, NFR-006

## Description
Provide example requests for all API endpoints to help developers understand expected formats. Examples should use realistic data and be copy-paste ready.

## Context / Background
Examples help developers quickly understand how to interact with the API without having to parse schema definitions. Good examples reduce onboarding time and support burden.

## Acceptance Criteria
- [ ] Each endpoint has at least one example request
- [ ] Examples use realistic data (not "foo" or "test123")
- [ ] Examples are copy-paste ready (valid JSON, proper formatting)
- [ ] Examples shall be copy-paste ready (NFR-002)
- [ ] Examples shall produce expected outputs when executed (NFR-006)

## Technical Notes
- Use curl commands or HTTP request format
- Include headers (Content-Type, Authorization)
- Use realistic data that demonstrates the endpoint's purpose
- Consider providing multiple examples for endpoints with optional parameters
- Example categories:
  - Create operations (POST)
  - Read operations (GET single, GET list)
  - Update operations (PUT/PATCH)
  - Delete operations (DELETE)
  - Special endpoints (workout generation, progression trigger)

## Dependencies
- Blocks: 005-workflow-documentation
- Blocked by: 001-endpoint-documentation
- Related: 003-example-responses

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/006-api-developer-experience/erd.md
