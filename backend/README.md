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

When you start the application with an **empty database**, sample users and travel plans are automatically created for testing and development.

#### Seeded Users

| Username | Email | Password | Role | Bio |
|----------|-------|----------|------|-----|
| `admin` | admin@travellink.local | AdminPass123! | Admin | System administrator |
| `traveller` | traveller@travellink.local | TravellerPass123! | Traveller | Loves exploring new places |
| `user` | user@travellink.local | UserPass123! | Simple | Just browsing travel plans |
| `alice` | alice@travellink.local | AlicePass123! | Traveller | Adventure seeker and photographer |
| `bob` | bob@travellink.local | BobPass123! | Traveller | Cultural explorer |
| `charlie` | charlie@travellink.local | CharliePass123! | Simple | Beach and nature lover |

**Role Capabilities:**
- **Admin**: Full system access, user management, content moderation, promotion request review
- **Traveller**: Create and publish travel plans, add/edit nodes, rate/comment on plans, submit for promotion
- **Simple**: Browse and view published plans, rate/comment (limited creation rights)

**Password Hashing**:
All passwords are securely hashed using bcrypt with default cost factor (10 iterations) during seeding. Passwords are logged in startup logs for reference but never stored in plaintext.

#### Seeded Travel Plans

**Published Plans (6 total):**

| Title | Destination | Author | Description |
|-------|-------------|--------|-------------|
| Summer European Adventure | Europe | alice | 3-week itinerary covering Paris, Rome, and Barcelona |
| Tokyo Food Tour | Tokyo, Japan | bob | Culinary journey through traditional and modern cuisine |
| New Zealand Road Trip | New Zealand | alice | Epic road trip around both islands with outdoor activities |
| Machu Picchu and Peruvian Highlands | Peru | traveller | Trek through Andes and explore ancient Inca ruins |
| Morocco Cultural Immersion | Morocco | bob | Markets, deserts, and mountain villages exploration |
| Iceland Ring Road Adventure | Iceland | charlie | Complete circle around Iceland with waterfalls and geysers |

**Draft Plans (4 total):**

| Title | Destination | Author | Description |
|-------|-------------|--------|-------------|
| Budget Southeast Asia Backpacking | Southeast Asia | charlie | Affordable route through Thailand, Vietnam, Cambodia |
| Caribbean Island Hopping | Caribbean | bob | Multiple islands for beaches, diving, and paradise |
| Norway Fjords Road Trip | Norway | alice | Scenic drive with hiking and photography opportunities |
| India Taj Mahal and Beyond | India | traveller | From Taj Mahal to Kerala backwaters cultural heritage |

**Plan Status Explanation:**
- **Published**: Visible to all users on the browse page, can be rated and commented on
- **Draft**: Only visible to plan creator, not listed publicly

**Note**: Seeding only happens when the database is empty. If users already exist, seeding is skipped automatically.

### Testing with Seeded Users

After starting the application, you can immediately test endpoints with the seeded credentials:

```bash
# Login as alice (traveller with published plans)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@travellink.local",
    "password": "AlicePass123!"
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

#### Test Scenarios with Seeded Data

Browse published plans:
```bash
curl -X GET http://localhost:8080/api/v1/plans
```

View alice's published plans:
```bash
curl -X GET http://localhost:8080/api/v1/plans/search?q=alice
```

View plan details:
```bash
curl -X GET http://localhost:8080/api/v1/plans/{plan-id}
```

Create a new plan (requires authentication):
```bash
curl -X POST http://localhost:8080/api/v1/plans \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Custom Plan",
    "destination": "Thailand",
    "nodes": []
  }'
```

### Resetting the Database

To reset and re-seed the database:

```bash
# SQLite - removes database and WAL files
rm -f backend/travellink.db backend/travellink.db-shm backend/travellink.db-wal

# Then restart the application
DB_TYPE=sqlite3 ./api
```

### Troubleshooting: UNIQUE constraint failed: schema_migrations.version

**Error Message:**
```
UNIQUE constraint failed: schema_migrations.version
failed to run database migrations: failed to apply migration NNN: failed to record migration
```

**Cause:**
This occurs when a migration was partially applied (schema change executed but record insertion failed), causing the migration version to be marked as applied when it actually failed. Subsequent runs attempt to re-apply the same migration but the version already exists.

**Solution:**
Delete the database files to force a clean migration run:

```bash
cd backend
rm -f travellink.db travellink.db-shm travellink.db-wal
./api
```

**Prevention:**
Migrations now use atomic transactions (since this fix) - both the schema change AND the migration record are committed together, or both are rolled back on failure. This prevents incomplete migration states.

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
