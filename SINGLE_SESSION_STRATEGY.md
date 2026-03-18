# Single Session Strategy - Prevent Multiple Device Login

## Overview
This document outlines the strategy to prevent users from logging in from multiple devices simultaneously. Only one active session per user will be allowed.

## Implementation Approach

### Option 1: Session Table with Device Tracking ✅ Recommended

**How it works:**
1. When user logs in, check if they have an active session
2. If active session exists, invalidate it and create new session
3. Store session info with device details
4. On each request, validate session is still active
5. On logout, mark session as inactive

**Pros:**
- Full control over sessions
- Can track device information
- Easy to implement "logout from all devices"
- Can show user their active sessions
- Audit trail of login history

**Cons:**
- Requires database query on each request (can be cached in Redis)
- Slightly more complex than token-only approach

### Option 2: Single Refresh Token Approach

**How it works:**
1. Store only one refresh token per user in database
2. When user logs in from new device, revoke old refresh token
3. Old device's access token expires naturally (15 min)
4. Old device cannot refresh token

**Pros:**
- Simpler implementation
- Less database overhead

**Cons:**
- Old device stays logged in until access token expires (up to 15 min)
- No device tracking
- No session history

---

## Recommended Implementation

### Database Schema

#### User Sessions Table
```sql
CREATE TABLE user_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    session_token VARCHAR(500) UNIQUE NOT NULL,
    device_info TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    CONSTRAINT fk_session_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_is_active ON user_sessions(is_active);
CREATE INDEX idx_user_sessions_session_token ON user_sessions(session_token);
```

### Login Flow with Single Session

```
1. User submits email + password
   ↓
2. Validate credentials
   ↓
3. Check for active sessions
   ↓
4. If active session exists:
   - Invalidate old session (set is_active = false)
   - Optionally notify user via email
   ↓
5. Generate OTP
   ↓
6. Create temporary session
   ↓
7. Return temp token + OTP sent message

8. User submits OTP + temp token
   ↓
9. Validate OTP
   ↓
10. Create new active session
    - Generate session token (UUID)
    - Store device info, IP, user agent
    - Set expiry (7 days)
    ↓
11. Generate JWT access token
    - Include session_token in claims
    ↓
12. Return access token + session token
```

### Session Validation Middleware

```go
// On each authenticated request:
1. Extract JWT from Authorization header
2. Validate JWT signature and expiry
3. Extract session_token from JWT claims
4. Check if session is still active in database
   - Query: SELECT is_active FROM user_sessions WHERE session_token = ?
5. If session inactive, return 401 Unauthorized
6. Update last_activity timestamp
7. Proceed with request
```

### Logout Flow

```
1. Extract session_token from JWT
   ↓
2. Mark session as inactive
   UPDATE user_sessions SET is_active = false WHERE session_token = ?
   ↓
3. Return success
```

### Logout from All Devices

```
1. Get user_id from JWT
   ↓
2. Invalidate all sessions for user
   UPDATE user_sessions SET is_active = false WHERE user_id = ?
   ↓
3. Return success
```

---

## Enhanced Features

### 1. Session Management Endpoints

```
GET /api/v1/auth/sessions
- List all active sessions for current user
- Show device info, IP, last activity

DELETE /api/v1/auth/sessions/:session_id
- Logout from specific session

DELETE /api/v1/auth/sessions/all
- Logout from all devices except current
```

### 2. Session Notifications

When a new login invalidates an old session:
- Send email notification to user
- Include device info and location
- Provide "This wasn't me?" link for security

### 3. Trusted Devices (Future Enhancement)

- Allow users to mark devices as "trusted"
- Trusted devices can have multiple concurrent sessions
- Require additional verification for untrusted devices

### 4. Session Limits

- Maximum session duration: 7 days
- Automatic cleanup of expired sessions
- Configurable session timeout for inactivity

---

## Security Considerations

### 1. Session Token Generation
```go
// Use cryptographically secure random token
sessionToken := uuid.New().String()
```

### 2. Session Hijacking Prevention
- Store session token hash in database (not plain text)
- Validate IP address consistency (optional, can break with mobile networks)
- Validate user agent consistency
- Implement session fingerprinting

### 3. Concurrent Login Attempts
- If user tries to login while OTP verification pending, invalidate old OTP
- Prevent session fixation attacks

### 4. Session Cleanup
```go
// Periodic cleanup job (run daily)
DELETE FROM user_sessions 
WHERE expires_at < NOW() 
   OR (is_active = false AND created_at < NOW() - INTERVAL '30 days')
```

---

## Implementation Checklist

### Phase 1: Basic Single Session
- [ ] Create user_sessions table migration
- [ ] Add UserSession model
- [ ] Create SessionService
- [ ] Update login flow to check/invalidate existing sessions
- [ ] Store session on successful OTP verification
- [ ] Add session validation to JWT middleware
- [ ] Implement logout endpoint

### Phase 2: Session Management
- [ ] List active sessions endpoint
- [ ] Logout from specific session
- [ ] Logout from all devices
- [ ] Session cleanup job

### Phase 3: Enhanced Security
- [ ] Email notifications on new login
- [ ] Device fingerprinting
- [ ] Session activity tracking
- [ ] Suspicious activity detection

---

## User Experience

### Scenario 1: User logs in from new device
```
1. User logs in from Device B
2. System detects active session on Device A
3. System invalidates Device A session
4. User on Device A gets 401 on next request
5. Device A shows: "You've been logged out because you logged in from another device"
6. Device B proceeds normally
```

### Scenario 2: User wants multiple devices
```
Option A: Implement "trusted devices" feature
Option B: Allow user to configure in settings
Option C: Banking standard - enforce single session for security
```

---

## Recommended Configuration

For a banking application, I recommend:

✅ **Enforce single active session** (highest security)
✅ **Email notification on new login**
✅ **7-day session expiry**
✅ **15-minute inactivity timeout** (configurable)
✅ **Session activity logging**
✅ **Ability to view and manage sessions**

This provides the best balance of security and user experience for a banking application.

---

## Next Steps

Would you like me to implement:
1. ✅ Basic single session enforcement
2. Session management endpoints
3. Email notifications
4. All of the above

Please confirm and I'll proceed with the implementation!
