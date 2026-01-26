# Authentication System

This document describes the authentication system used in this application. All agents working on authentication-related features must follow these specifications.

## Authentication Method

### Email/Password Authentication
- Users authenticate using **email and password**
- **Name is optional** - users may provide a name during registration, but it's not required
- Email serves as the primary identifier

## Session Management

### Session-Based Authentication
- The system uses **session-based authentication** (not JWT-based)
- Each authenticated user receives a **session ID**
- Sessions are stored server-side and can be revoked

### Why Not JWTs?
JWTs are **not used** for the following reasons:
1. **RPC Architecture**: The frontend/backend split and RPC nature of the system makes JWT revocation difficult
2. **Revocation is Critical**: The system requires the ability to revoke sessions immediately
3. **Server-Side Control**: Session-based auth provides better control over session lifecycle

## Token Transmission

### Bearer Token Authentication
- The **session ID is sent as a bearer token** to the backend
- Format: `Authorization: Bearer <session-id>`
- The backend validates the session ID and retrieves session information

## Session Lifecycle

### Session Creation
- Created when user successfully authenticates (email/password)
- Session ID is generated and stored server-side
- Session ID is returned to the client

### Session Validation
- Backend validates bearer token (session ID) on each RPC request
- Validates that session exists and is active
- Retrieves user information from session

### Session Revocation
- Sessions can be revoked server-side
- Revocation is immediate (no token expiration wait)
- Revoked sessions are rejected on subsequent requests

### Session Cleanup Strategy

The system provides multiple mechanisms for session cleanup:

#### 1. Logout Invalidation
- When a user logs out, only that specific session is deleted
- Other active sessions for the same user remain valid
- Logout is token-specific, not user-wide

#### 2. Session Expiration
- Sessions have a fixed lifetime (7 days from creation)
- Expired sessions are rejected during validation
- The `ValidateSession` method checks `now().After(session.ExpiresAt)`

#### 3. Expired Session Purge
- The `CleanupExpiredSessions()` method removes all sessions where `expires_at < current_time`
- Returns the count of sessions deleted
- This should be called periodically (e.g., via a cron job or scheduled task)
- Uses the `idx_sessions_expires_at` index for efficient querying

#### 4. Cascade Deletion on User Delete
- The sessions table has a foreign key constraint: `FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE`
- When a user is deleted, all their sessions are automatically deleted
- No manual cleanup required for user deletion scenarios

#### 5. User Session Management
- `DeleteUserSessions(userID)` removes all sessions for a specific user
- Useful for "logout from all devices" functionality
- Returns the count of sessions deleted

#### Best Practices
- Run `CleanupExpiredSessions()` periodically (daily recommended)
- Monitor session table growth
- Consider implementing session limits per user if needed

## User Registration

### Required Fields
- **Email**: Required, serves as unique identifier
- **Password**: Required, must meet security requirements

### Optional Fields
- **Name**: Optional, users may provide a name but it's not required

## Security Considerations

### Password Requirements
- Password requirements should be defined (minimum length, complexity, etc.)
- Passwords must be hashed (never stored in plain text)
- Use secure password hashing (e.g., bcrypt, argon2)

### Session Security
- Session IDs must be cryptographically secure (random, unpredictable)
- Session IDs should be long enough to prevent brute force
- Implement session timeout/expiration policies
- Protect against session fixation attacks

### Bearer Token Security
- Transmit bearer tokens over HTTPS only
- Validate bearer tokens on every request
- Handle invalid/expired tokens gracefully

## Implementation Notes

### RPC Integration
- Authentication middleware validates bearer token before RPC handler execution
- Session information is made available to RPC handlers
- Failed authentication returns appropriate error response

### Database Schema
- Users table: email (unique), password_hash, name (nullable), created_at, etc.
- Sessions table: session_id (primary key), user_id (foreign key), created_at, expires_at, revoked_at, etc.

### API Endpoints (RPC Methods)
- `Register(email, password, name?)` - Create new user account
- `Login(email, password)` - Authenticate and create session
- `Logout(session_id)` - Revoke session
- `ValidateSession(session_id)` - Check if session is valid
- `GetUserInfo(session_id)` - Get user information from session

## Testing

### E2E Test Considerations
- Tests must create users and sessions
- Tests must validate bearer token authentication
- Tests must verify session revocation works
- Each test gets a fresh database (no session leakage between tests)

## Future Considerations

### Potential Improvements
- Evaluate if bearer token approach is optimal (noted as "maybe that's a bad idea")
- Consider session refresh mechanisms
- Consider multi-device session management
- Consider "remember me" functionality

## Integration with Tech Stack

This authentication system integrates with:
- **Go backend**: Session management implemented in Go
- **SQLite database**: Sessions stored in database
- **sqlc**: Type-safe database queries for session operations
- **goose migrations**: Schema changes for users and sessions tables
- **RPC system**: Authentication middleware for RPC handlers

## Related Documents

- See `tech-stack.md` for technology stack details
- See `erd.md` for requirements document structure (if creating auth ERD)
