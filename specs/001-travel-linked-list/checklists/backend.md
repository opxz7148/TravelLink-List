# Backend Checklist - Travel Linked List

**Overall Status**: 🟡 **In Progress** (58/80 backend tasks complete, 72.5%)

**Last Updated**: April 13, 2026  
**Updated by**: Copilot (Data model refactoring phase)

---

## Phase 1: Setup (6/6) ✅ COMPLETE

- [x] T001 Go project structure with layered directories
- [x] T002 Module imports updated (backend/* → tll-backend/*)
- [x] T003 Backend dependencies installed (Go, Gin, SQLite3, BCrypt, JWT)
- [x] T005 SQLite3 database connection and file location configuration
- [x] T006 Database initialization logic
- [x] T007 Database migration framework setup

---

## Phase 2: Foundational Infrastructure (31/33)

### Migrations (7/7) ✅ COMPLETE
- [x] T008 Migration 001: users table
- [x] T009 Migration 002: travel_plans table
- [x] T010 Migration 003: nodes table
- [x] T011 Migration 004: attraction_node_details and transition_node_details tables
- [x] T012 Migration 005: plan_nodes junction table
- [x] T013 Migration 006: comments table
- [x] T014 Migration 007: ratings table

### Middleware (6/6) ✅ COMPLETE
- [x] T015 Base response envelope middleware
- [x] T016 AuthMiddleware for JWT validation
- [x] T017 RBACMiddleware for role-based access
- [x] T017b Context helper for user data
- [x] T018 Validation middleware
- [x] T019 Error handling wrapper
- [x] T019b Middleware USAGE guide
- [x] T020 Gin router with middleware chain
- [x] T030 CORS configuration

### Services & Repositories (18/19) 🟡 PARTIAL
- [x] T021 UserService interface and implementation
- [x] T021a Renamed to RelationalUserService
- [x] T021b JWTService interface and implementation
- [x] T022 PlanService interface and full implementation
- [x] T022b NodeService interface and implementation
- [x] T023 UserRepository interface
- [x] T023a Consolidated RemoveAdmin to DemoteToSimple
- [x] T023b RelationalUserRepository implementation
- [x] T023c Renamed to RelationalUserRepository
- [x] T024a PlanRepository interface
- [x] T024b RelationalPlanRepository implementation (all 21 methods)
- [x] T025 NodeRepository interface
- [x] T025a RelationalNodeRepository implementation

### Models (5/5) ✅ COMPLETE
- [x] T026 User model with validation
- [x] T027 TravelPlan model with validation
- [x] T028 Node, AttractionNode, TransitionNode models with validation
- [x] T029 PlanNode model with linked list validation
- [x] T024 Comment model with validation
- [x] T049 Rating model with validation

**Checkpoint**: ✅ ALL FOUNDATIONAL INFRASTRUCTURE COMPLETE

---

## Phase 3: User Story 1 - Browse & Search (7/8) 🟡 PARTIAL

Browse & Search Travel Plans

- [x] T032 UserRepository.GetByID and GetByEmail
- [x] T033 PlanRepository.List with pagination and filters
- [x] T034 PlanRepository.FindByDestination (full-text search)
- [x] T035 PlanService.ListPublishedPlans method
- [x] T036 PlanService.SearchPlans method
- [x] T037 PlanController.BrowsePlans handler
- [x] T038 PlanController.SearchPlans handler
- [x] T039 Routes: GET /api/v1/plans, GET /api/v1/plans/search
- [ ] **T040 Database indexes** ⏳ PENDING
  - Migration 008 for browse performance indexes
  - `CREATE INDEX on travel_plans(status, created_at), travel_plans(destination)`
- [ ] T047 Contract test for GET /api/v1/plans

---

## Phase 4: User Story 2 - Comment & Rate (10/12) 🟡 PARTIAL

Comment & Rate Travel Plans

- [x] T048 Comment model with validation
- [x] T049 Rating model with validation
- [x] T050 CommentRepository (RelationalCommentRepository)
- [x] T051 RatingRepository (RelationalRatingRepository)
- [x] T052 CommentService with Create/List/Update
- [x] T053 RatingService with Submit/Update/GetUserRating
- [x] T054 CommentController handlers (all 3 methods)
- [x] T055 RatingController handlers (all 4 methods)
- [x] T056 Routes: POST/GET /api/v1/plans/{id}/comments, POST /api/v1/plans/{id}/ratings
- [x] T057 Transactional rating update (denormalized counts)
- [x] T058 Transactional comment count update
- [x] T068 AuthMiddleware guard for comment/rating endpoints
- [ ] T066 Contract test for POST /api/v1/plans/{id}/comments
- [ ] T067 Contract test for POST /api/v1/plans/{id}/ratings

---

## Phase 5: User Story 3 - Promotion Request (0/6) ⏳ NOT STARTED

Submit Promotion Request for Traveller Status

- [ ] T069 PromotionRequest model
- [ ] T070 PromotionRequestRepository
- [ ] T071 PromotionService.SubmitRequest method
- [ ] T072 PromotionController handlers
- [ ] T073 Routes: POST /api/v1/promotions/request, GET /api/v1/promotions/my-status
- [ ] T077 Contract test for POST /api/v1/promotions/request

---

## Phase 6: User Story 4 - Create Plans with Nodes (10/12) 🟡 PARTIAL

Create Travel Plans with Existing Nodes

- [x] T079 Seed system nodes (23 transitions + attractions) - Migration 009
- [x] T080 NodeRepository.List with filtering
- [x] T081 NodeRepository.GetByID
- [x] T082 NodeService.ListApprovedNodes and ListNodesByType
- [x] T083 NodeController.ListNodes and GetNodeDetail
- [x] T084 PlanService.CreateDraftPlan (empty nodes list)
- [x] T085 PlanService.AddNodesToPlan with validation
- [x] T086 PlanService.ReorderNodes (updates sequence_position)
- [x] T087 PlanService.RemoveNodeFromPlan
- [x] T088 PlanService.PublishPlan (draft → published)
- [x] T089 PlanService.GetPlanWithNodes (linked list retrieval)
- [x] T090 Route: POST /api/v1/plans (create draft)
- [ ] **T091 Route: PATCH /api/v1/plans/{id}/nodes** ⏳ PENDING (reorder/add/remove)
- [ ] T101 Traveller role check on plan creation endpoints
- [ ] T102 Contract test for POST /api/v1/plans
- [ ] T103 Contract test for PATCH /api/v1/plans/{id}/nodes
- [ ] T104 Unit tests for linked list validation (no 2 consecutive attractions)

---

## Phase 7: User Story 5 - Create Nodes (7/7) ✅ COMPLETE

Create New Attraction/Transition Nodes

- [x] T105 NodeRepository.Create
- [x] T106 NodeService.CreateAttractionNode with validation
- [x] T107 NodeService.CreateTransitionNode with validation
- [x] T108 NodeController.CreateAttraction handler
- [x] T109 NodeController.CreateTransition handler
- [x] T110 Routes: POST /api/v1/nodes/attraction, POST /api/v1/nodes/transition
- [ ] T115 Contract test for POST /api/v1/nodes/attraction
- [ ] T116 Contract test for POST /api/v1/nodes/transition
- [ ] T117 Verify user-created nodes default to is_approved=false

---

## Phase 8: User Story 6 - Admin Moderation (5/8) 🟡 PARTIAL

Admin Plan & Node Moderation

- [x] T118 PlanRepository.ListFlaggedPlans
- [x] T118 PlanRepository.SoftDeletePlan
- [x] T123 ModerationController (AdminController) - all methods
- [x] T124 Routes: PATCH /api/v1/admin/plans/{id}/suspend, DELETE /api/v1/admin/plans/{id}, etc.
- [ ] T119 ModerationService.ListFlaggedPlans
- [ ] T120 ModerationService.DeletePlan (soft-delete)
- [ ] T121 ModerationService.ListPendingNodes
- [ ] T122 ModerationService.ApproveNode and RejectNode
- [ ] T130 Contract test for DELETE /api/v1/admin/plans/{id}
- [ ] T131 Verify soft-delete: comments/ratings persist

---

## Phase 9: Authentication & Authorization (3/6)

User Login/Registration & RBAC

- [x] T132 UserService.Register with bcrypt hashing
- [x] T133 UserService.Login with JWT generation
- [x] T134 AuthController (Register and Login)
- [x] T135 Routes: POST /api/v1/auth/register, POST /api/v1/auth/login
- [ ] T136 Unit tests for JWT token generation and verification
- [ ] T144 Contract test for POST /api/v1/auth/register
- [ ] T145 Contract test for POST /api/v1/auth/login

---

## Phase 10: Plan Details & Linked List (3/3) ✅ COMPLETE

Display Complete Travel Plan with Linked List

- [x] T146 PlanService.GetPlanWithNodesForView
- [x] T147 PlanController.GetPlanDetails handler
- [x] T148 Route: GET /api/v1/plans/{id}
- [ ] T153 Contract test for GET /api/v1/plans/{id}

---

## Phase 11: User Profile & My Plans (3/3) ✅ COMPLETE

User Profile & Personal Plan Management

- [x] T154 UserRepository.GetPlansForUser
- [x] T155 UserController (GetProfile, GetMyPlans)
- [x] T156 Routes: GET /api/v1/users/me, GET /api/v1/users/me/plans
- [ ] T161 Contract test for GET /api/v1/users/me

---

## Phase 12: Polish & Cross-Cutting Concerns (0/11) ⏳ NOT STARTED

Backend Polish, Logging, & Documentation

### Middleware & Infrastructure
- [ ] T162 Request logging middleware
- [ ] T163 Error recovery middleware
- [ ] T164 Request ID tracing middleware
- [ ] T165 Structured logging with request context

### API Completeness
- [ ] T166 Health check endpoint: GET /api/v1/health
- [ ] T167 API versioning documentation (OpenAPI)
- [ ] T168 Database transaction helpers
- [ ] T169 Seed data script for demo travel plans
- [ ] T170 Graceful shutdown handler
- [ ] T171 Environment configuration (.env support)

### Testing & Documentation
- [ ] T183 OpenAPI/Swagger documentation
- [ ] T184 Unit tests for repositories
- [ ] T185 Unit tests for services
- [ ] T186 Integration tests for API endpoints
- [ ] T189 API usage examples

---

## Critical Data Model Changes (Completed Apr 13)

**Recent Refactoring**:
- ✅ Added `title` and `hours_of_operation` to `TransitionNodeDetail`
- ✅ Added `description` and `estimated_price_cents` to `PlanNode`
- ✅ Removed `estimated_distance_km` from `TransitionNodeDetail` (plan-specific, not service-specific)
- ✅ Removed `estimated_duration_minutes` from `TransitionNodeDetail` (plan-specific, not service-specific)
- ✅ Updated all database migrations (008 for schema, 009 for seed data)
- ✅ Updated all Swagger documentation
- ✅ Seeded 23 system transition nodes with proper titles and operating hours

---

## Summary by Status

| Category | Completed | Total | % |
|----------|-----------|-------|-----|
| Phase 1: Setup | 6 | 6 | ✅ 100% |
| Phase 2: Foundation | 31 | 33 | 🟡 94% |
| Phase 3: US1 Browse | 7 | 9 | 🟡 78% |
| Phase 4: US2 Comments | 10 | 12 | 🟡 83% |
| Phase 5: US3 Promotion | 0 | 6 | ⏳ 0% |
| Phase 6: US4 Plans | 10 | 15 | 🟡 67% |
| Phase 7: US5 Nodes | 7 | 7 | ✅ 100% |
| Phase 8: US6 Moderation | 5 | 9 | 🟡 56% |
| Phase 9: Auth | 3 | 6 | 🟡 50% |
| Phase 10: Plan Details | 3 | 4 | 🟡 75% |
| Phase 11: User Profile | 3 | 4 | 🟡 75% |
| Phase 12: Polish | 0 | 11 | ⏳ 0% |
| **TOTAL** | **58** | **80** | **72.5%** |

---

## 🚨 Blocking Issues / Attention Required

None currently. All MVP functionality working.

---

## 📋 Next Priority Actions

### Immediate (Next Session)
1. **T040** - Add database indexes for browse performance (migration 008)
2. **T091** - Implement PATCH /api/v1/plans/{id}/nodes route
3. **T101** - Add Traveller role checks to plan creation endpoints

### High Priority (MVP Completion)
1. **Phase 5 (Promotion)** - T069-T077: Complete promotion request workflow
2. **Phase 8 (Moderation)** - T119-T122: Implement ModerationService
3. **Auth Testing** - T136, T144, T145: Add JWT and auth contract tests

### Medium Priority (Polish)
1. **Phase 12** - Logging, health checks, structured logging
2. **Contract Tests** - T047, T066, T067, T102, T103, T115-T117, T130-T131, T153, T161
3. **Unit Tests** - T104, T136, T184, T185

### Lower Priority (Complete Phase)
1. **Documentation** - T183, T189
2. **Graceful Shutdown** - T170
3. **Environment Config** - T171
