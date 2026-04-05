# Implementation Plan: Travel Linked List

**Branch**: `001-travel-linked-list` | **Date**: 2026-04-05 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-travel-linked-list/spec.md`

**Note**: This plan implements the Travel Linked List feature—a collaborative platform enabling users to create, share, and discover travel itineraries built as linked lists of attraction and transition nodes.

## Summary

Travel Linked List is a full-stack collaborative travel planning application. Core value: users create travel itineraries as linked lists of individual nodes (attractions like restaurants/hotels/museums and transitions like walking/driving), then share with a community for discovery, rating, and commentary. The platform supports three user roles (simple user / traveller / admin) with escalating permissions. Technical approach: layered backend (Go/Gin) with separate controller, service, model, and repository layers; Vue 3 frontend for responsive UX; SQLite3 file-based database for persistence. MVP focuses on P1 features: browse/search plans, create plans from existing nodes, comment/rate, with moderation capability locked to admin.

## Technical Context

**Language/Version**: Go 1.x (backend), Vue 3 Composition API (frontend)  
**Primary Dependencies**: Gin web framework (backend routing), Vue 3 + Vite (frontend build), SQLite3 driver (github.com/mattn/go-sqlite3)  
**Storage**: SQLite3 (file-based at backend/travellink.db; abstracted via repository layer); single-file embedded database, no external infrastructure required  
**Testing**: Go built-in testing package + testify/assert (backend unit tests); Vitest or Jest (frontend unit tests); REST client for API integration tests (Postman/Thunder Client)  
**Target Platform**: Linux/macOS/Windows (backend runs as standalone binary with embedded SQLite3); modern browsers (frontend, responsive from 375px+)
**Project Type**: Full-stack web service (backend API + Vue SPA frontend)  
**Performance Goals**: Sub-500ms search response time for 10k travel plans; 100 concurrent users during browse/search; sub-100ms API endpoint latency (p95)  
**Constraints**: JWT tokens with 1-hour expiration; plan publication atomic (no partial state); node orphaning allowed (no cascade delete); comments/ratings persist even if plan deleted; single SQLite3 writer (acceptable for MVP)  
**Scale/Scope**: MVP supports ~100 travellers producing ~1k travel plans; ~1k simple users browsing; initial launch: 6 major screens (Browse, View Plan, Create Plan, Profile, Admin Dashboard, Settings)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

✅ **GATE PASS** — All constitution principles satisfied:

| Principle | Compliance | Evidence |
|-----------|-----------|----------|
| I. Layered Architecture | ✅ PASS | Backend separates: Controller layer (HTTP handlers + middleware for auth), Service layer (PlanService, NodeService, CommentService, AuthService), Model layer (User, TravelPlan, Node entities), Repository layer (abstracts SQLite3 access). |
| II. API Contracts | ✅ PASS | All endpoints define request/response schemas in OpenAPI contract. Standard envelope: `{ "success": bool, "data": T, "error": { "code", "message" } }`. HTTP semantics enforced (200/201 success, 400 validation, 401 auth, 403 forbidden, 404 not found). |
| III. Auth & Authorization | ✅ PASS | JWT middleware enforces authentication before route handlers. RBAC middleware checks roles (simple/traveller/admin) before service invocation. Admin-only endpoints (delete plan, bulk moderation) guarded. Passwords bcrypt-hashed. |
| IV. Service-Based Use Cases | ✅ PASS | Six coarse-grained services: AuthService (login/register), PlanService (CRUD plans), NodeService (manage attraction/transition nodes), CommentService (create/list comments), RatingService (rate plans), ModerationService (admin deletion/review). Each stateless, receives fully-validated objects. |
| V. Repository Abstraction | ✅ PASS | Repositories (UserRepository, PlanRepository, NodeRepository, CommentRepository, RatingRepository) abstract all DB queries. Services depend only on interfaces, not concrete implementations. Database technology change (SQLite3 → PostgreSQL → MongoDB) impacts only repository layer. |

**Summary**: Feature design is CONSTITUTIONAL COMPLIANT. No violations. No justifications required.

## Project Structure

### Documentation (this feature)

```text
specs/001-travel-linked-list/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command) - research unknowns
├── data-model.md        # Phase 1 output (/speckit.plan command) - entities and relationships
├── quickstart.md        # Phase 1 output (/speckit.plan command) - developer onboarding
├── contracts/           # Phase 1 output (/speckit.plan command) - API contract specs
│   ├── auth-contract.md
│   ├── plan-contract.md
│   ├── node-contract.md
│   ├── comment-contract.md
│   └── rating-contract.md
└── tasks.md             # Phase 2 output (/speckit.tasks command) - implementation tasks
```

### Source Code (repository root)

Web application structure (frontend + backend):

**Backend** (Go/Gin):
```text
backend/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── controllers/                # Controller layer (HTTP handlers + middleware)
│   │   ├── auth_controller.go
│   │   ├── plan_controller.go
│   │   ├── node_controller.go
│   │   ├── comment_controller.go
│   │   ├── rating_controller.go
│   │   └── moderation_controller.go
│   ├── middleware/                 # Middleware for auth, RBAC, validation
│   │   ├── auth_middleware.go
│   │   ├── rbac_middleware.go
│   │   └── validation_middleware.go
│   ├── services/                   # Service layer (business logic)
│   │   ├── auth_service.go
│   │   ├── plan_service.go
│   │   ├── node_service.go
│   │   ├── comment_service.go
│   │   ├── rating_service.go
│   │   └── moderation_service.go
│   ├── models/                     # Model layer (domain entities)
│   │   ├── base.go                 # Already exists
│   │   ├── user.go
│   │   ├── travel_plan.go
│   │   ├── node.go
│   │   ├── attraction_node.go
│   │   ├── transition_node.go
│   │   ├── comment.go
│   │   ├── rating.go
│   │   └── promotion_request.go
│   ├── repositories/               # Repository layer (data access)
│   │   ├── user_repository.go
│   │   ├── plan_repository.go
│   │   ├── node_repository.go
│   │   ├── comment_repository.go
│   │   ├── rating_repository.go
│   │   └── promotion_request_repository.go
│   ├── database/                   # Already exists + migrations
│   │   ├── database.go
│   │   ├── database_test.go
│   │   └── migrations/
│   │       ├── 001_create_users_table.sql
│   │       ├── 002_create_travel_plans_table.sql
│   │       ├── 003_create_nodes_table.sql
│   │       ├── 004_create_plan_nodes_association_table.sql
│   │       ├── 005_create_comments_table.sql
│   │       ├── 006_create_ratings_table.sql
│   │       └── 007_create_promotion_requests_table.sql
│   ├── server/                     # Already exists + updates
│   │   ├── server.go
│   │   ├── server_test.go
│   │   └── routes.go               # Will be extended with new routes
│   └── config/                     # Configuration management
│       └── config.go
├── tests/
│   ├── integration/                # API integration tests
│   │   ├── auth_integration_test.go
│   │   ├── plan_integration_test.go
│   │   └── ...
│   └── unit/                       # Service/Repository unit tests
│       ├── auth_service_test.go
│       ├── plan_service_test.go
│       └── ...
├── go.mod                          # Already exists
├── Makefile                        # Already exists
├── docker-compose.yml              # Already exists
└── README.md                       # Already exists
```

**Frontend** (Vue 3):
```text
frontend/
├── src/
│   ├── App.vue                     # Already exists
│   ├── main.ts                     # Already exists
│   ├── components/
│   │   ├── HelloWorld.vue          # Already exists (remove in MVP)
│   │   ├── TheWelcome.vue          # Already exists (remove in MVP)
│   │   ├── Navigation.vue          # Header/navbar
│   │   ├── PlanCard.vue            # Reusable plan preview card
│   │   ├── NodeCard.vue            # Reusable node display
│   │   ├── CommentItem.vue         # Comment display
│   │   ├── RatingStars.vue         # Rating input/display
│   │   └── icons/                  # Already exists
│   ├── pages/
│   │   ├── BrowsePage.vue          # Browse/search travel plans (P1)
│   │   ├── ViewPlanPage.vue        # View single plan + comments (P1)
│   │   ├── CreatePlanPage.vue      # Create/edit plan (P1)
│   │   ├── LoginPage.vue           # User login
│   │   ├── RegisterPage.vue        # User registration
│   │   ├── ProfilePage.vue         # User profile + my plans
│   │   ├── AdminDashboard.vue      # Moderation interface (P2)
│   │   └── NotFoundPage.vue        # 404 error
│   ├── services/
│   │   ├── api.ts                  # HTTP client setup (Axios or Fetch)
│   │   ├── auth_service.ts         # Login/logout/token management
│   │   ├── plan_service.ts         # Plan CRUD operations
│   │   ├── node_service.ts         # Node operations
│   │   ├── comment_service.ts      # Comment operations
│   │   └── rating_service.ts       # Rating operations
│   ├── stores/                     # Pinia/Vuex state management
│   │   ├── auth_store.ts           # Current user, roles, token
│   │   ├── plan_store.ts           # Plans list, current plan
│   │   └── ui_store.ts             # Loading, modals, notifications
│   ├── types/                      # TypeScript interfaces
│   │   ├── index.ts
│   │   └── api.ts                  # API request/response types
│   ├── assets/                     # Already exists
│   │   ├── base.css
│   │   └── main.css
│   └── router/                     # Vue Router configuration
│       └── index.ts
├── tests/
│   ├── unit/
│   │   ├── plan_service.test.ts
│   │   └── auth_store.test.ts
│   └── integration/
│       └── browse_flow.test.ts     # End-to-end browse flow
├── vite.config.ts                  # Already exists
├── tsconfig.json                   # Already exists
├── package.json                    # Already exists
├── index.html                      # Already exists
└── README.md                       # Already exists
```

**Structure Decision**: Full-stack web application (Option 2) with separate backend (Go/Gin) and frontend (Vue 3) directories. Backend uses layered architecture per constitution. Frontend uses component-based architecture with Composition API. Shared API contract documentation in `/specs/001-travel-linked-list/contracts/`.
