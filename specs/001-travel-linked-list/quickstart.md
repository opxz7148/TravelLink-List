# Quickstart: Travel Linked List Developer Guide

**Feature**: Travel Linked List  
**Date**: 2026-04-05  
**Target Audience**: Backend & Frontend developers setting up the feature locally

---

## Prerequisites

- Go 1.x installed (backend)
- Node.js 18+ with npm (frontend)
- SQLite3 (comes with most systems; verify: `sqlite3 --version`)
- Git (version control)
- VS Code or preferred IDE

---

## Project Setup

### 1. Clone Repository & Prepare Branch

```bash
# Already on feature branch
git status
# On branch 001-travel-linked-list

# Install dependencies
cd backend && go mod download
cd ../frontend && npm install
```

### 2. Create SQLite3 Database File

```bash
cd backend
# Database is auto-created on first server run, OR manually create:
sqlite3 travellink.db ".tables"
# (Just opens an empty DB in interactive mode; press Ctrl+D to exit)
```

### 3. Run Database Migrations

```bash
cd backend
# Migrations run automatically on server startup
# OR manually trigger:
go run cmd/api/main.go --migrate
```

Expected output:
```
[INFO] Running migration 001_create_users_table.sql
[INFO] Running migration 002_create_travel_plans_table.sql
... (7 migrations total)
[INFO] All migrations applied successfully
[INFO] Database: /path/to/backend/travellink.db
```

Verify:
```bash
sqlite3 travellink.db "SELECT count(*) FROM users;"
# Expected: 0
```

### 4. Seed Initial Data (Optional)

```bash
# Load system nodes (attractions/transitions)
go run cmd/api/main.go --seed-data

# Check database
sqlite3 travellink.db "SELECT COUNT(*) FROM nodes WHERE is_approved = true;"
# Expected: ~50 system nodes
```

### 5. Start Backend Server

```bash
cd backend
go run cmd/api/main.go
# [INFO] Server starting on http://localhost:8080
# [INFO] Database: travellink.db (file-based SQLite3)
```

### 6. Start Frontend Dev Server (in new terminal)

```bash
cd frontend
npm run dev
# ➜  Local:   http://localhost:5173/
```

### 7. Verify Setup

- Backend health: `curl http://localhost:8080/api/v1/health`
- Frontend: Open http://localhost:5173 in browser
- Expected: Browse travel plans page loads
- Database file: `backend/travellink.db` exists and is ~100KB+

---

## Project Structure Quick Reference

### Backend File Organization

```
backend/
├── cmd/api/main.go              # Entry point
├── internal/
│   ├── controllers/             # HTTP handlers
│   │   ├── auth_controller.go
│   │   ├── plan_controller.go
│   ├── services/                # Business logic
│   │   ├── auth_service.go      # Auth/registration
│   │   ├── plan_service.go      # Plan CRUD
│   ├── models/                  # Domain entities
│   │   ├── user.go
│   │   ├── travel_plan.go
│   ├── repositories/            # Data access
│   │   ├── user_repository.go
│   │   ├── plan_repository.go
│   ├── database/
│   │   ├── database.go          # DB connection
│   │   ├── migrations/          # SQL files
```

**Layer Dependencies** (architecture principle):
- Controllers depend on Middleware + Services
- Services depend on Models + Repositories
- Repositories depend on database/sql + Models
- Models have no external dependencies

### Frontend File Organization

```
frontend/src/
├── components/              # Reusable Vue components
│   ├── PlanCard.vue
│   ├── Navigation.vue
├── pages/                   # Route pages
│   ├── BrowsePage.vue       # Browse/search plans
│   ├── ViewPlanPage.vue     # View single plan
│   ├── CreatePlanPage.vue   # Create/edit plan
├── services/                # API client
│   ├── api.ts               # HTTP client setup
│   ├── plan_service.ts      # Plan API calls
├── stores/                  # Pinia state
│   ├── plan_store.ts
│   ├── auth_store.ts
```

---

## Common Development Tasks

### Task 1: Add a New API Endpoint

**Example**: Add endpoint to get user's own travel plans

1. **Define contract** in `/specs/001-travel-linked-list/contracts/plan-contract.md`
   ```
   GET /api/v1/plans/my-plans
   Authorization: Bearer <token>
   Response: { plans: [...] }
   ```

2. **Create/update service** in `backend/internal/services/plan_service.go`
   ```go
   func (s *PlanService) GetUserPlans(ctx context.Context, userID string) ([]TravelPlan, error) {
       return s.planRepo.FindByAuthor(ctx, userID)
   }
   ```

3. **Add controller** in `backend/internal/controllers/plan_controller.go`
   ```go
   func (c *PlanController) GetUserPlans(w http.ResponseWriter, r *http.Request) {
       // Extract userID from JWT token (via middleware)
       userID := r.Context().Value("user_id").(string)
       plans, err := c.planService.GetUserPlans(r.Context(), userID)
       // Write response JSON
   }
   ```

4. **Register route** in `backend/internal/server/routes.go`
   ```go
   router.GET("/api/v1/plans/my-plans", authMiddleware, planController.GetUserPlans)
   ```

5. **Test endpoint** in `backend/tests/integration/plan_integration_test.go`
   ```go
   func TestGetUserPlans(t *testing.T) {
       // Create test user, plans, then call endpoint
   }
   ```

### Task 2: Add a Frontend Component

**Example**: Create a plan card component

1. **Create component** `frontend/src/components/PlanCard.vue`
   ```vue
   <template>
     <div class="plan-card">
       <h3>{{ plan.title }}</h3>
       <p>{{ plan.destination }}</p>
       <span class="rating">⭐ {{ plan.rating_average }}</span>
     </div>
   </template>

   <script setup lang="ts">
   import { PropType } from 'vue';
   interface Plan { ... }
   defineProps({ plan: Object as PropType<Plan> });
   </script>

   <style scoped>
   .plan-card { /* styles */ }
   </style>
   ```

2. **Use in page** `frontend/src/pages/BrowsePage.vue`
   ```vue
   <PlanCard v-for="plan in plans" :key="plan.id" :plan="plan" />
   ```

3. **Test component** `frontend/tests/unit/PlanCard.test.ts`
   ```ts
   import { mount } from '@vue/test-utils';
   import PlanCard from '@/components/PlanCard.vue';
   // Test rendering, props, events
   ```

### Task 3: Write a Database Migration

**Example**: Add table for travel plan favorites (Phase 2 feature)

1. **Create migration file** `backend/internal/database/migrations/008_create_favorites_table.sql`
   ```sql
   CREATE TABLE favorites (
       id UUID PRIMARY KEY,
       user_id UUID NOT NULL REFERENCES users(id),
       plan_id UUID NOT NULL REFERENCES travel_plans(id),
       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
       UNIQUE(user_id, plan_id)
   );

   CREATE INDEX idx_favorites_user_id ON favorites(user_id);
   ```

2. **Migration runner** automatically discovers and executes in `internal/database/database.go`

3. **Test migration** by running:
   ```bash
   docker-compose down
   docker-compose up -d
   go run cmd/api/main.go --migrate
   # Verify new table exists
   psql -U postgres -d travellink -c "\dt favorites"
   ```

---

## Testing Locally

### Backend Tests

```bash
cd backend

# Run all tests
go test ./...

# Run specific test file
go test ./internal/services -v -run TestPlanService

# Run with coverage
go test -cover ./...
```

### Frontend Tests

```bash
cd frontend

# Run all tests
npm run test

# Run with watch mode
npm run test:watch

# Generate coverage report
npm run test:coverage
```

### Manual API Testing

**Option 1: cURL**

```bash
# Register a new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "Password123!"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123!"
  }'

# Browse plans (use token from login response)
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/plans?destination=Paris
```

**Option 2: Postman/Thunder Client**

1. Import collection from `specs/001-travel-linked-list/contracts/` markdown files
2. Set Postman variables: `BASE_URL=http://localhost:8080`, `TOKEN=<from login>`
3. Run requests

**Option 3: Browser DevTools**

1. Open http://localhost:5173
2. Register/login via UI
3. Browse plans in browser
4. Network tab shows all API requests

---

## Constitution Compliance Checklist

When adding new features, verify:

- [ ] **Layered Architecture**: Does the code follow Controller → Service → Model → Repository?
- [ ] **API Contracts**: Is the endpoint documented in `contracts/` with request/response schema?
- [ ] **Authentication**: Are protected endpoints guarded by `@AuthMiddleware` and role-checked?
- [ ] **Services**: Are services stateless and receive fully-validated objects?
- [ ] **Repositories**: Do all database queries go through repository interfaces?
- [ ] **Tests**: Are unit tests written for services, integration tests for endpoints?

---

## Common Issues & Fixes

| Issue | Cause | Fix |
|-------|-------|-----|
| No database file created | SQLite3 file not yet created | Run `go run cmd/api/main.go --migrate` or start server once |
| "database is locked" error | SQLite3 write contention (multiple processes writing) | Kill other processes; SQLite3 single-writer limitation noted for MVP |
| Migrations fail with "table already exists" | Migrations already ran | Delete `backend/travellink.db` and restart server to re-run |
| Cannot connect to database | File permissions issue | Check `backend/travellink.db` permissions: `ls -la backend/travellink.db` |
| Frontend can't reach backend API | CORS not configured | Check `backend/internal/server/routes.go` for CORS middleware |
| JWT token invalid | Token expired (1 hour TTL) | Login again to get new token |
| Tests failing due to database state | Previous test data persists | Delete `backend/travellink.db` before running tests; tests create fresh DB |

---

## Next Steps

1. **Review architecture**: Read [data-model.md](data-model.md) to understand entities
2. **Understand contracts**: Review [contracts/](contracts/) for API specifications
3. **Start implementing**: Pick a P1 user story from [../../spec.md](../../spec.md)
4. **Follow TDD**: Write tests first, then implementation
5. **Verify compliance**: Ensure layered architecture is maintained

---

## Quick Links

- **Specification**: [../spec.md](../spec.md)
- **Implementation Plan**: [../plan.md](../plan.md)
- **Data Model**: [../data-model.md](../data-model.md)
- **API Contracts**: [../contracts/](../contracts/)
- **Research**: [../research.md](../research.md)
- **Constitution**: [../../.specify/memory/constitution.md](../../.specify/memory/constitution.md)

---

## Support

- For architecture questions, refer to [../../.specify/memory/constitution.md](../../.specify/memory/constitution.md)
- For API design issues, review [../contracts/](../contracts/)
- For code structure, check layered architecture layout in `backend/internal/`
- For database queries, see migration files in `backend/internal/database/migrations/`
