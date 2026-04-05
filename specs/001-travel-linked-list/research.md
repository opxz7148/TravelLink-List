# Research: Travel Linked List Feature

**Branch**: `001-travel-linked-list` | **Date**: 2026-04-05 | **Phase**: 0 (Research & Unknowns Resolution)

## Overview

This document consolidates findings from research into key technical and design unknowns extracted from the plan.md Technical Context. All items marked "NEEDS CLARIFICATION" in the plan have been researched and validated.

---

## Research Results

### 1. Database Schema & Migration Strategy

**Unknown**: How to structure SQLite3 schema for linked list representation?

**Decision**: Use explicit junction table (`PlanNode`) with sequence position rather than embedded document approach. SQLite3 supports this via relational schema with AUTOINCREMENT and UNIQUE constraints.

**Implementation**:
- Base tables: `users`, `travel_plans`, `nodes`, `comments`, `ratings`
- Node type discrimination: Use single `nodes` table with `type` TEXT column (discriminator: 'attraction' or 'transition')
- Attraction-specific columns in separate `attraction_node_details` table (1:1 relationship with `nodes` via FOREIGN KEY)
- Transition-specific columns in separate `transition_node_details` table (1:1 relationship with `nodes` via FOREIGN KEY)
- Junction table: `plan_nodes` with columns `(plan_id TEXT, node_id TEXT, sequence_position INTEGER, created_at TIMESTAMP)`
- SQLite3 supports: TEXT for UUIDs, INTEGER PRIMARY KEY, INTEGER for counters, TIMESTAMP for dates, FOREIGN KEY constraints (with PRAGMA foreign_keys=ON)

**Rationale**: Avoids schema bloat from optional columns; maintains normalization while supporting inheritance pattern; SQLite3's simplicity remains advantage over PostgreSQL.

**Alternatives Considered**:
- ❌ JSON columns (SQLite3 supports JSON but less normalized; loses ability to query nodes independently)
- ❌ Strict two-table inheritance: Leads to duplicate schemas; complicates shared node queries
- ✅ Relational with discriminator: Clean, testable, portable to PostgreSQL if needed

**Migrations**: 7 migrations plus 1 seed data migration (008+); all use standard SQL compatible with SQLite3.

**File Size**: Typical SQLite3 db for 1k plans + 10k nodes: ~50-100MB (single file deployment advantage)

---

### 2. Authentication & Token Management

**Unknown**: JWT implementation details (expiration, refresh mechanism, role encoding)?

**Decision**: Standard JWT with 1-hour expiration + optional refresh token strategy.

**Implementation**:
- **Access Token**: 1-hour expiration; encodes `sub` (user_id), `role` (simple/traveller/admin), `exp`, `iat`, `iss`
- **Refresh Token** (optional): 7-day expiration; stored in HTTP-only cookie; used to silently renew access token without re-login
- **Token Issuance**: On login (POST /auth/login); credentials validated against bcrypt-hashed password
- **Token Validation**: Middleware (AuthMiddleware) checks Authorization header, parses JWT, verifies signature + expiration
- **Role Encoding**: `"role"` claim set to string: `"simple"` / `"traveller"` / `"admin"`

**Rationale**: 1-hour balances security (limited exposure if token leaked) and UX (doesn't force frequent re-login).

**Alternatives Considered**:
- ❌ No refresh token: Forces user re-login every hour during browsing (poor UX)
- ❌ 24-hour access token: Increases risk window if token compromised

**Middleware Order**: AuthMiddleware (verify token exists) → RBACMiddleware (check role requirements) → then route handler.

---

### 3. Search & Filtering Performance

**Unknown**: How to efficiently search 10k travel plans by destination, keyword, or author?

**Decision**: Database indexes + query optimization; optional Redis caching for popular searches.

**Implementation**:
- **Indexes**:
  - `travel_plans(status, created_at)` - browse published plans chronologically
  - `travel_plans(title, destination)` - full-text search (PostgreSQL tsvector or trigram index)
  - `travel_plans(author_id)` - user's own plans
  - `nodes(name, type)` - search nodes by name/category when creating plans
  - `comments(plan_id, created_at)` - retrieve comments for plan view

- **Query Strategy**:
  - Full-text search using SQLite3 `LIKE` or `GLOB` pattern matching with indexes for MVP performance
  - Pagination with OFFSET/LIMIT (default 20 items per page)
  - Aggregate rating/comment counts in same query (avoid N+1)
  - Soft limit: SQLite3 single-writer design acceptable for MVP scale (<100 concurrent); migrate to PostgreSQL if needed

- **Caching** (optional, Phase 2):
  - Redis TTL 1 hour for popular search queries
  - Invalidate cache on plan create/update/delete
  - Cache key: `search:{query_hash}:{sort_by}:{page}`

**Rationale**: SQLite3 indexes sufficient for MVP browse/search; single-writer limitation acceptable for collaborative features with moderate concurrency (~5-10 simultaneous edits); full-text search via pattern matching + indexes meets 500ms target.

**Alternatives Considered**:
- ❌ Elasticsearch: Overkill for MVP; adds operational complexity
- ❌ No indexing: 10k plan table scans would timeout

**Target Metric**: Browse/search response < 500ms (p95) with default indexes.

---

### 4. Node Approval Workflow

**Unknown**: How to handle user-created nodes that need admin approval?

**Decision**: Separate `is_approved` flag; unapproved nodes visible only to creator until admin approval.

**Implementation**:
- **Node Table**: Add `created_by` (user_id) and `is_approved` (boolean, default true for system nodes, false for user-created)
- **Visibility Rules**:
  - System nodes (created by seed data): Always visible
  - User-created nodes: Only visible to creator until approved by admin
- **Admin Moderation Flow**:
  - New endpoint: GET `/api/admin/nodes/pending` (returns user-created, unapproved nodes)
  - Approve action: PATCH `/api/admin/nodes/{id}/approve` (sets is_approved=true)
  - Reject action: DELETE `/api/admin/nodes/{id}` with notification to creator

**Rationale**: Balances content quality (admins curate) with creator autonomy (can see own drafted nodes).

**Alternatives Considered**:
- ❌ Auto-approve all user nodes: Allows spam/inappropriate content
- ❌ Queue all nodes for approval: Slows content creation in MVP; not scalable

**Implications**: Service layer (NodeService) enforces visibility rules; repositories filter queries based on calling user's role.

---

### 5. API Response Versioning

**Unknown**: How to version API responses for future compatibility?

**Decision**: Accept header versioning with v1 baseline; all responses include `api_version` field.

**Implementation**:
- **Baseline**: All endpoints default to `v1` (application/vnd.travellink.v1+json` implicit)
- **Response Envelope**:
  ```json
  {
    "success": true,
    "api_version": "1.0",
    "data": { /* response payload */ },
    "error": null,
    "timestamp": "2026-04-05T10:30:00Z"
  }
  ```
- **Error Envelope**:
  ```json
  {
    "success": false,
    "api_version": "1.0",
    "data": null,
    "error": {
      "code": "VALIDATION_ERROR",
      "message": "Email is required",
      "details": { /* optional */ }
    },
    "timestamp": "2026-04-05T10:30:00Z"
  }
  ```

**Rationale**: Explicit api_version in response aids client debugging and logging; timestamp useful for request tracing.

**Alternatives Considered**:
- ❌ HTTP Accept header only: Clients might not set it; harder to debug in browser DevTools
- ❌ URL path versioning (/v1/plans): Cumbersome; versioning often happens alongside Accept header anyway

---

### 6. Frontend State Management

**Unknown**: Pinia vs Vuex for state management?

**Decision**: Pinia (modern Composition API-first store, recommended by Vue 3 team).

**Stores**:
- `authStore`: Current user, roles, JWT token, login/logout actions
- `planStore`: Selected plan, plans list, search filters, loading states
- `uiStore`: Global modals, notifications, loading indicators

**Rationale**: Pinia is Vue 3's native choice; simpler syntax than Vuex; better TypeScript support.

**Alternatives Considered**:
- ❌ Vuex: Still supported but legacy pattern; more boilerplate
- ❌ Prop-drilling only: Works for small apps; scales poorly with deeply nested components

---

### 7. Frontend Build & Deployment

**Unknown**: Vite dev server optimization and production build strategy?

**Decision**: Vite for dev (already configured); production build outputs to `dist/` with gzip compression.

**Implementation**:
- **vite.config.ts**: Already exists; ensure `base: "/"` for root path deployment
- **Build command**: `npm run build` outputs optimized bundle to `frontend/dist/`
- **Asset optimization**:
  - Code splitting: Separate main, vendor, route chunks
  - CSS autoprefixing (browserslist includes modern browsers)
  - Image optimization: Vite defaults to inlining small assets
  - Source maps: Enabled in dev, disabled in production builds

**Deployment**: Serve `frontend/dist/` via web server (nginx/Apache) or Node.js static server; backend API available at `/api`.

**Rationale**: Vite is production-ready; no additional build optimization tools needed for MVP.

---

### 8. Testing Strategy

**Unknown**: Unit vs Integration vs E2E test split?

**Decision**: Start with unit + integration; E2E deferred to Phase 2.

**Backend Tests**:
- **Unit**: Service layer tests (mock repositories); Repository tests (test database queries)
- **Integration**: Controller + Service flow tests (start in-memory SQLite or test SQLite database)
- **Tool**: Go built-in `testing` + `testify/assert/require` + prepared statements for database test fixtures

**Frontend Tests**:
- **Unit**: Component rendering tests (@vue/test-utils)
- **Integration**: Store + Component interaction tests
- **Tool**: Vitest or Jest; Vue Test Utils; optional Mock Service Worker (msw) for API mocking

**E2E** (Phase 2): Playwright or Cypress for full user journeys (browse → search → login → create).

**Rationale**: Unit+Integration cover MVP quickly; E2E expensive to maintain early; defer until features stabilize.

---

## Summary of Decisions

| Task | Found | Decision | Status |
|------|-------|----------|--------|
| Database schema for linked lists | ✓ | SQLite3 with single node table + junction (PlanNode) with sequence position | RESOLVED |
| JWT implementation | ✓ | 1-hour access token + optional refresh; role encoded in claims | RESOLVED |
| Search performance | ✓ | SQLite3 indexes + pagination; Redis cache Phase 2 | RESOLVED |
| Node approval workflow | ✓ | `is_approved` flag; user-created nodes queue for admin | RESOLVED |
| API versioning | ✓ | v1 baseline + `api_version` in response envelope | RESOLVED |
| Frontend state management | ✓ | Pinia store (authStore, planStore, uiStore) | RESOLVED |
| Build & deployment | ✓ | Vite dev, optimized static dist/ build | RESOLVED |
| Testing strategy | ✓ | Unit + Integration MVP (SQLite3 testing); E2E Phase 2 | RESOLVED |

All unknowns resolved. Ready for Phase 1: Design & Contracts.
