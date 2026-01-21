# 001: Endpoint Documentation

## ERD Reference
Implements: REQ-DOC-001
Related to: NFR-001, NFR-003, NFR-005

## Description
Document all API endpoints with HTTP method, path, description, request body schema, and response schema for success and error cases. This is the foundation for all API documentation.

## Context / Background
PowerPro is a headless API. External developers need clear documentation of what endpoints exist and how to use them. This ticket establishes the core endpoint documentation that all other documentation builds upon.

## Acceptance Criteria
- [ ] Every endpoint documented with HTTP method and path
- [ ] Every endpoint has a clear description of its purpose
- [ ] Request body schema documented where applicable
- [ ] Response schema documented for success cases
- [ ] Response schema documented for error cases
- [ ] Documentation is clear and accessible to developers unfamiliar with PowerPro (NFR-001)
- [ ] Documentation can be maintained/generated from code where possible (NFR-003)
- [ ] Documentation accurately reflects current API behavior (NFR-005)

## Technical Notes
- Review all existing endpoints in the codebase
- Use OpenAPI/Swagger specification for machine-readable docs
- Use Markdown for human-readable documentation
- Consider generating documentation from code/schemas for maintainability
- Endpoints to document include:
  - Lift CRUD endpoints
  - LiftMax CRUD endpoints
  - Prescription CRUD endpoints
  - Day/Week/Cycle CRUD endpoints
  - Program CRUD endpoints
  - User Program State endpoints
  - Workout generation endpoint
  - Progression endpoints
  - Authentication endpoints

## Dependencies
- Blocks: 002-example-requests, 003-example-responses, 004-error-documentation, 005-workflow-documentation
- Blocked by: None
- Related: All other documentation tickets

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/006-api-developer-experience/erd.md
- PRD: phases/in-progress/001-core-foundation/sprints/in-progress/006-api-developer-experience/prd.md
