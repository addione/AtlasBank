# Authentication Strategy for AtlasBank

## Overview
For a secure banking application, we need a robust multi-layered authentication approach. Here's a comprehensive strategy:

## Recommended Authentication Flow

### 🔐 **Two-Factor Authentication (2FA) with JWT**

#### Phase 1: Initial Login (Username + Password)
1. User submits email/username + password
2. Server validates credentials
3. If valid, generate and send OTP to user's email/phone
4. Return a temporary session token (short-lived, 5 minutes)
5. User must verify OTP within time limit

#### Phase 2: OTP Verification
1. User submits OTP with temporary session token
2. Server validates OTP
3. If valid, issue JWT access token + refresh token
4. Return tokens to client

#### Phase 3: Authenticated Requests
1. Client includes JWT in Authorization header
2. Server validates JWT for each request
3. Refresh token used to get new access token when expired

---

## Proposed Implementation

### Authentication Methods

#### 1. **JWT (JSON Web Tokens)** ✅ Recommended
**Pros:**
- Stateless authentication
- Scalable (no server-side session storage)
- Can include user claims (roles, permissions)
- Industry standard

**Cons:**
- Cannot be revoked easily (use short expiry + refresh tokens)
- Token size can be large

**Implementation:**
```
Access Token: 15 minutes expiry
Refresh Token: 7 days expiry
Store refresh tokens in database for revocation capability
```

#### 2. **Session-Based Authentication**
**Pros:**
- Easy to revoke
- Server has full control

**Cons:**
- Requires server-side storage (Redis)
- Less scalable
- Sticky sessions needed for load balancing

#### 3. **OAuth 2.0 / OpenID Connect**
**Pros:**
- Industry standard
- Supports third-party login (Google, etc.)
- Delegated authorization

**Cons:**
- Complex implementation
- Overkill for internal banking app

---

## Recommended Security Features

### 1. **Multi-Factor Authentication (MFA)**
- ✅ Password (something you know)
- ✅ OTP via Email/SMS (something you have)
- 🔄 Biometric (future: something you are)

### 2. **Rate Limiting**
- Limit login attempts: 5 attempts per 15 minutes
- Lock account after 10 failed attempts
- CAPTCHA after 3 failed attempts

### 3. **Password Security**
- ✅ Bcrypt hashing (already implemented)
- Minimum 8 characters
- Require: uppercase, lowercase, number, special char
- Password history (prevent reuse of last 5 passwords)
- Password expiry (90 days for banking)

### 4. **Token Security**
- Short-lived access tokens (15 min)
- Refresh token rotation
- Store refresh tokens in database
- Revoke tokens on logout
- Blacklist compromised tokens

### 5. **Additional Security Layers**
- IP whitelisting for sensitive operations
- Device fingerprinting
- Geolocation checks
- Transaction signing with OTP
- Session timeout (15 min inactivity)

---

## Proposed Login Flow

### Step 1: Login Request
```
POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}

Response (200 OK):
{
  "message": "OTP sent to your email",
  "temp_token": "eyJhbGc...",
  "expires_in": 300,
  "user_id": 123
}
```

### Step 2: OTP Verification
```
POST /api/v1/auth/verify-otp
{
  "temp_token": "eyJhbGc...",
  "otp": "0000"
}

Response (200 OK):
{
  "message": "Login successful",
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "token_type": "Bearer",
  "expires_in": 900,
  "user": {
    "id": 123,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "is_verified": true
  }
}
```

### Step 3: Authenticated Request
```
GET /api/v1/users/me
Authorization: Bearer eyJhbGc...

Response (200 OK):
{
  "user": { ... }
}
```

### Step 4: Token Refresh
```
POST /api/v1/auth/refresh
{
  "refresh_token": "eyJhbGc..."
}

Response (200 OK):
{
  "access_token": "eyJhbGc...",
  "expires_in": 900
}
```

### Step 5: Logout
```
POST /api/v1/auth/logout
Authorization: Bearer eyJhbGc...

Response (200 OK):
{
  "message": "Logged out successfully"
}
```

---

## Database Schema Updates Needed

### 1. Refresh Tokens Table
```sql
CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    token VARCHAR(500) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### 2. Login Attempts Table (Rate Limiting)
```sql
CREATE TABLE login_attempts (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    attempted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    success BOOLEAN DEFAULT FALSE
);
```

### 3. User Sessions Table (Optional)
```sql
CREATE TABLE user_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    device_info TEXT,
    ip_address VARCHAR(45),
    last_activity TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

---

## Required Go Packages

```bash
# JWT handling
go get github.com/golang-jwt/jwt/v5

# Rate limiting
go get golang.org/x/time/rate

# Password validation
go get github.com/go-playground/validator/v10
```

---

## Implementation Checklist

### Phase 1: Basic Auth (Current Sprint)
- [x] User registration with password hashing
- [x] OTP table and service
- [ ] Login endpoint (email + password)
- [ ] OTP generation on login
- [ ] OTP verification endpoint
- [ ] JWT token generation
- [ ] JWT middleware for protected routes

### Phase 2: Enhanced Security
- [ ] Refresh token mechanism
- [ ] Rate limiting on login
- [ ] Account lockout after failed attempts
- [ ] Password strength validation
- [ ] Session management

### Phase 3: Advanced Features
- [ ] Device fingerprinting
- [ ] IP whitelisting
- [ ] Geolocation checks
- [ ] Audit logging
- [ ] Suspicious activity detection

---

## Security Best Practices

1. **Never log sensitive data** (passwords, tokens, OTPs)
2. **Use HTTPS only** in production
3. **Implement CORS** properly
4. **Sanitize all inputs** to prevent injection
5. **Use prepared statements** for database queries
6. **Implement request signing** for critical operations
7. **Regular security audits** and penetration testing
8. **Keep dependencies updated**
9. **Use environment variables** for secrets
10. **Implement proper error handling** (don't leak info)

---

## Recommended Approach for AtlasBank

**I recommend implementing:**

✅ **JWT-based authentication** with:
- Email + Password login
- OTP verification (2FA)
- Short-lived access tokens (15 min)
- Long-lived refresh tokens (7 days)
- Token revocation capability
- Rate limiting on login attempts

This provides:
- ✅ Strong security (2FA)
- ✅ Scalability (stateless JWT)
- ✅ Good UX (refresh tokens reduce re-login)
- ✅ Compliance (banking security standards)

**Next Steps:**
1. Implement JWT service
2. Create auth controller with login/verify/refresh/logout
3. Add JWT middleware for protected routes
4. Implement rate limiting
5. Add comprehensive logging

Would you like me to proceed with this implementation?
