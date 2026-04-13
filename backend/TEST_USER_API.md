# User API Test Script Documentation

## Overview

`test_user_api.sh` is a comprehensive bash script that tests all user-related API endpoints using curl. It covers authentication, profile management, password changes, and error handling scenarios.

## Prerequisites

- **Bash** (macOS/Linux)
- **curl** (should be pre-installed)
- **jq** (for JSON parsing)
  ```bash
  # macOS
  brew install jq
  
  # Ubuntu/Debian
  sudo apt-get install jq
  ```
- **Backend server running** on `http://localhost:8080`

## Quick Start

### 1. Start the Backend Server

```bash
cd backend
go run cmd/api/main.go
```

Expected output:
```
[INFO] Server starting on http://localhost:8080
[INFO] Database: travellink.db
```

### 2. Run the Test Script

```bash
cd backend
./test_user_api.sh
```

## Test Coverage

The script tests the following scenarios:

### Happy Path ✓
1. **User Registration** - Register a new user with email, username, password
2. **User Login** - Authenticate with email and password, receive JWT token
3. **Get User Profile** - Fetch user profile (public endpoint, no auth required)
4. **Update User Profile** - Update display_name and bio (auth required)
5. **Change Password** - Change password with old password verification
6. **Login with New Password** - Verify new password works
7. **Logout** - Logout endpoint (returns 204 No Content)

### Admin Operations ✓
8. **Create Test Admin User** - Setup for admin permission tests
9. **Admin Operation Without Permission** - Verify non-admin request is rejected

### Error Handling ✓
10. **Invalid Token** - Request with malformed token is rejected
11. **Weak Password** - Registration with weak password is rejected
12. **Duplicate Email** - Registration with existing email is rejected

## Output Format

The script provides color-coded output for easy reading:

```
→ [Step execution messages] - Yellow
✓ [Success messages] - Green
✗ [Error messages] - Red
```

Each test includes:
- Step description
- Full JSON response (formatted with jq)
- Success/failure indicator

## Sample Output

```
===============================================
  TEST 1: User Registration
===============================================
→ Registering user: john_doe...
Response: {
  "success": true,
  "api_version": "1.0",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "john_doe@example.com",
      "username": "john_doe",
      "display_name": "John Doe",
      "role": "simple",
      "created_at": "2026-04-13T10:30:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600
  },
  "error": null,
  "timestamp": "2026-04-13T10:30:00Z"
}
✓ User registered successfully
→ Access Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
→ User ID: 550e8400-e29b-41d4-a716-446655440000
```

## Manual Testing with curl

### Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "john_doe",
    "password": "SecurePass123!",
    "display_name": "John Doe"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!"
  }'
```

### Get User Profile

```bash
# Public - no auth required
curl -X GET http://localhost:8080/api/v1/users/{user_id}

# Or with auth (optional)
curl -X GET http://localhost:8080/api/v1/users/{user_id} \
  -H "Authorization: Bearer {access_token}"
```

### Update Profile

```bash
curl -X PUT http://localhost:8080/api/v1/users/{user_id} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {access_token}" \
  -d '{
    "display_name": "Updated Name",
    "bio": "My bio here"
  }'
```

### Change Password

```bash
curl -X POST http://localhost:8080/api/v1/users/{user_id}/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {access_token}" \
  -d '{
    "old_password": "SecurePass123!",
    "new_password": "NewPass456!"
  }'
```

### Logout

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer {access_token}"
```

## Troubleshooting

### "Server is not running"
- Make sure backend server is started: `go run cmd/api/main.go`
- Verify it's accessible: `curl http://localhost:8080/health`

### "jq: command not found"
- Install jq: `brew install jq` (macOS) or `apt-get install jq` (Linux)

### "Permission denied" when running script
- Make it executable: `chmod +x test_user_api.sh`

### Tests fail with "User already exists"
- Run the script again (it uses unique timestamps or adjust email addresses)
- Or clear the database and restart the server

### Invalid token error
- JWT tokens expire after 1 hour
- Re-run the script or login again to get a new token

## Environment Variables

You can customize the script by modifying these variables at the top:

```bash
API_BASE_URL="http://localhost:8080/api/v1"    # API base URL
CONTENT_TYPE="Content-Type: application/json"  # Content type header
```

## Running Individual Tests

To run specific tests, you can call individual functions:

```bash
source test_user_api.sh
check_server
test_register
test_login
test_get_profile
```

## Continuous Integration

This script can be integrated into CI/CD pipelines:

```bash
#!/bin/bash
set -e  # Exit on any error

# Start server in background
go run cmd/api/main.go &
SERVER_PID=$!

# Wait for server to start
sleep 2

# Run tests
./test_user_api.sh

# Cleanup
kill $SERVER_PID
```

## Notes

- The script is **idempotent** - it uses unique identifiers (emails, usernames) for each run
- All sensitive data (tokens, passwords) are handled securely
- The script provides detailed error messages for debugging
- No cleanup is performed - test data remains in the database for inspection

## Next Steps

- Create similar test scripts for Plan, Node, Comment, and Rating endpoints
- Automate script execution in CI/CD pipeline
- Add performance testing and load testing scenarios
- Integrate with automated testing frameworks (e.g., Postman, k6)
