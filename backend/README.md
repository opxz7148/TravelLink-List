# Project backend

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```

## Database Setup

### Automatic Seeding

When you start the application with an **empty database**, sample users are automatically created for testing and development:

| Username | Email | Password | Role | Purpose |
|----------|-------|----------|------|---------|
| `admin` | admin@travellink.local | AdminPass123! | Admin | Full system access, user management, moderation |
| `traveller` | traveller@travellink.local | TravellerPass123! | Traveller | Create and publish travel plans, add nodes, rate/comment |
| `user` | user@travellink.local | UserPass123! | Simple | Browse and view travel plans, rate/comment (read-only) |

**Note**: Seeding only happens when the database is empty. If users already exist, seeding is skipped automatically.

### Testing with Seeded Users

After starting the application, you can immediately test endpoints with the seeded credentials:

```bash
# Login as admin
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@travellink.local",
    "password": "AdminPass123!"
  }'

# Response includes JWT token
{
  "success": true,
  "data": {
    "access_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

Use the returned `access_token` with Bearer authentication for protected endpoints:

```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer <access_token>"
```

### Resetting the Database

To reset and re-seed the database:

```bash
# SQLite
rm -f backend/travellink.db backend/travellink.db-shm backend/travellink.db-wal

# Then restart the application
DB_TYPE=sqlite3 ./api
```

## API Documentation

### Swagger UI

Interactive API documentation is available at:
```
http://localhost:8080/swagger/index.html
```

All endpoints are documented with:
- Request/response schemas
- Parameter descriptions
- Authentication requirements
- Example values

#### Bearer Token Authentication in Swagger UI

When testing protected endpoints in Swagger UI:

1. Click the **"Authorize"** button (🔒)
2. In the "Bearer (apiKey)" field, enter your full token **including the "Bearer " prefix**:
   ```
   Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
   ```
3. Click **"Authorize"** to apply authentication
4. Try protected endpoints - the Bearer token will be automatically added to requests

**Note**: Unlike some modern API documentation tools, Swagger UI requires you to manually enter the `Bearer ` prefix as part of the token. Always use the format: `Bearer <token>`
