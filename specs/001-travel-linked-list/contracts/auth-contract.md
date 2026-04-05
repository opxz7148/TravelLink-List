# API Contract: Authentication Endpoints

**Feature**: Travel Linked List  
**Domain**: User Management & Authentication  
**Version**: 1.0  
**Last Updated**: 2026-04-05

---

## Overview

Authentication endpoints handle user registration, login, token refresh, and logout. All endpoints use REST conventions with JSON payloads. Authentication state is managed via JWT tokens.

---

## Endpoint 1: User Registration

### Request

```
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "john_doe",
  "password": "SecurePass123!",
  "display_name": "John Doe"
}
```

**Field Validation**:
- `email`: RFC 5322 format, max 255 chars, must not already exist
- `username`: 3-50 chars, alphanumeric + underscore, must not already exist
- `password`: min 8 chars, at least one uppercase, one digit, one special character
- `display_name`: optional, max 100 chars

### Response: Success (201 Created)

```json
{
  "success": true,
  "api_version": "1.0",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "username": "john_doe",
      "display_name": "John Doe",
      "role": "simple",
      "created_at": "2026-04-05T10:30:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600
  },
  "error": null,
  "timestamp": "2026-04-05T10:30:00Z"
}
```

### Response: Validation Error (400 Bad Request)

```json
{
  "success": false,
  "api_version": "1.0",
  "data": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Email already registered",
    "details": {
      "field": "email",
      "reason": "Email must be unique"
    }
  },
  "timestamp": "2026-04-05T10:30:00Z"
}
```

---

## Endpoint 2: User Login

### Request

```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Field Validation**:
- `email`: max 255 chars, required
- `password`: required, compared against bcrypt hash

### Response: Success (200 OK)

```json
{
  "success": true,
  "api_version": "1.0",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "username": "john_doe",
      "display_name": "John Doe",
      "role": "traveller",
      "created_at": "2026-04-05T10:30:00Z",
      "last_login_at": "2026-04-05T10:35:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9... (optional, if refresh token strategy enabled)",
    "token_type": "Bearer",
    "expires_in": 3600
  },
  "error": null,
  "timestamp": "2026-04-05T10:35:00Z"
}
```

**Note**: Refresh token, if implemented, is sent as HTTP-only cookie; not included in JSON response.

### Response: Authentication Failure (401 Unauthorized)

```json
{
  "success": false,
  "api_version": "1.0",
  "data": null,
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Email or password is incorrect"
  },
  "timestamp": "2026-04-05T10:30:00Z"
}
```

### Response: User Inactive (401 Unauthorized)

```json
{
  "success": false,
  "api_version": "1.0",
  "data": null,
  "error": {
    "code": "ACCOUNT_INACTIVE",
    "message": "Your account has been deactivated by an administrator"
  },
  "timestamp": "2026-04-05T10:30:00Z"
}
```

---

## Endpoint 3: Token Refresh (Optional)

### Request

```
POST /api/v1/auth/refresh
Authorization: Bearer <refresh_token>
```

**Note**: Refresh token sent as HTTP-only cookie or in Authorization header (if cookie-based auth not available).

### Response: Success (200 OK)

```json
{
  "success": true,
  "api_version": "1.0",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9... (new token)",
    "token_type": "Bearer",
    "expires_in": 3600
  },
  "error": null,
  "timestamp": "2026-04-05T10:30:00Z"
}
```

### Response: Invalid Token (401 Unauthorized)

```json
{
  "success": false,
  "api_version": "1.0",
  "data": null,
  "error": {
    "code": "INVALID_REFRESH_TOKEN",
    "message": "Refresh token expired or invalid"
  },
  "timestamp": "2026-04-05T10:30:00Z"
}
```

---

## Endpoint 4: Logout

### Request

```
POST /api/v1/auth/logout
Authorization: Bearer <access_token>
```

### Response: Success (204 No Content)

```
204 No Content
```

**Note**: Client should discard access_token and refresh_token locally. Server may invalidate refresh_token if persistent token store is maintained (Phase 2 optimization).

---

## JWT Token Structure

**Access Token Claims**:
```json
{
  "sub": "550e8400-e29b-41d4-a716-446655440000",        // User ID
  "role": "traveller",                                   // User role
  "email": "user@example.com",                           // User email
  "iss": "travellink",                                   // Issuer
  "iat": 1712309400,                                     // Issued at
  "exp": 1712313000                                      // Expires in (1 hour)
}
```

**Refresh Token** (if implemented):
- Same structure as access token
- Expiration: 7 days
- Issued only on login, not automatically

---

## Request/Response Headers

| Header | Example | Scope |
|--------|---------|-------|
| `Authorization` | `Bearer eyJ...` | Request (protected endpoints) |
| `Content-Type` | `application/json` | Request/Response |
| `Accept` | `application/vnd.travellink.v1+json` | Request (optional version signal) |

---

## Error Codes

| Code | HTTP Status | Meaning |
|------|------------|---------|
| `VALIDATION_ERROR` | 400 | Invalid input data (email exists, weak password, etc.) |
| `INVALID_CREDENTIALS` | 401 | Email or password mismatch |
| `ACCOUNT_INACTIVE` | 401 | User account deactivated by admin |
| `INVALID_TOKEN` | 401 | JWT token invalid, expired, or malformed |
| `INVALID_REFRESH_TOKEN` | 401 | Refresh token expired or not found |
| `INTERNAL_SERVER_ERROR` | 500 | Unexpected server error |

---

## Summary

All authentication endpoints follow:
- REST conventions (POST for state change, GET for retrieval)
- Consistent response envelope with success/error structure
- JWT-based stateless authentication
- API versioning via /v1/ path prefix
- Proper HTTP status codes (201 for created, 400 for client error, 401 for auth)
- JSON request/response with UTF-8 encoding
