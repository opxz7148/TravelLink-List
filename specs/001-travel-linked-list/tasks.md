# Tasks: Travel Linked List (Organized by Backend/Frontend)

**Input**: Design documents from `/specs/001-travel-linked-list/`  
**Prerequisites**: plan.md (✓), spec.md (✓), data-model.md (✓), contracts/ (✓), research.md (✓)

**Note**: Tasks reorganized by Backend/Frontend to enable clear team assignment while maintaining User Story context for integration.

---

## Format Reference

Each task follows: `- [ ] [TaskID] [P?] [Story] Description with exact file path`

- **[P]**: Can run in parallel (different files, no blocking dependencies)
- **[Story]**: User story label (US1, US2, etc.)
- **ID**: Sequential (T001, T002, T003...)

---

# BACKEND TASKS

## Backend Phase 1: Setup (Project Initialization)

**Purpose**: Backend project structure, dependencies, database setup - prerequisite for all user stories

- [x] T001 Setup Go project structure with layered directories: cmd/api, internal/{controllers,services,models,repositories,middleware,database} in backend/
- [x] T002 Update module imports from backend/* to tll-backend/* across all Go files (cmd/api/main.go, internal/server/server.go)
- [x] T003 [P] Install backend dependencies: Go 1.x, Gin, SQLite3 driver (gorm.io/driver/sqlite), BCrypt, JWT (github.com/golang-jwt/jwt/v5) in backend/go.mod
- [x] T005 Configure SQLite3 database connection and file location (backend/travellink.db) in backend/internal/database/database.go
- [x] T006 [P] Create database initialization logic in backend/internal/database/database.go (auto-creates DB file on first run)
- [x] T007 Create database migration framework in backend/internal/database/migrations/

**Checkpoint**: Backend project structure ready; dependencies installed; SQLite3 database auto-created on first run; migration framework operational

---

## Backend Phase 2: Foundational Infrastructure (Blocking Prerequisites)

**Purpose**: Core systems all user stories depend on - MUST complete before user story implementation

### Database Migrations

- [x] T008 Create migration 001: users table (id TEXT PK, email TEXT UNIQUE, username TEXT UNIQUE, role TEXT, password_hash TEXT, is_active BOOL, created_at TIMESTAMP) in backend/internal/database/migrations/001_create_users_table.sql
- [x] T009 [P] Create migration 002: travel_plans table (id TEXT PK, title TEXT, destination TEXT, author_id TEXT, status TEXT, rating_count INT, rating_sum INT, comment_count INT, created_at TIMESTAMP, updated_at TIMESTAMP) in backend/internal/database/migrations/002_create_travel_plans_table.sql
- [x] T010 [P] Create migration 003: nodes table with discriminator (id TEXT PK, type TEXT, created_by TEXT, is_approved BOOL, created_at TIMESTAMP) in backend/internal/database/migrations/003_create_nodes_table.sql
- [x] T011 [P] Create migration 004: attraction_node_details and transition_node_details tables in backend/internal/database/migrations/004_create_node_details_tables.sql
- [x] T012 [P] Create migration 005: plan_nodes junction table (plan_id TEXT, node_id TEXT, sequence_position INT) in backend/internal/database/migrations/005_create_plan_nodes_table.sql
- [x] T013 [P] Create migration 006: comments table (id TEXT PK, plan_id TEXT, author_id TEXT, text TEXT, created_at TIMESTAMP) in backend/internal/database/migrations/006_create_comments_table.sql
- [x] T014 [P] Create migration 007: ratings table (id TEXT PK, plan_id TEXT, user_id TEXT, stars INT, created_at TIMESTAMP) in backend/internal/database/migrations/007_create_ratings_table.sql

### Middleware

- [x] T015 Create base response envelope middleware in backend/internal/middleware/response_middleware.go
- [x] T016 [P] Implement AuthMiddleware for JWT token validation in backend/internal/middleware/auth_middleware.go
- [x] T017 [P] Implement RBACMiddleware for role-based access control (RequireRole) in backend/internal/middleware/auth_middleware.go
- [x] T017b Create context helper for accessing user data in handlers in backend/internal/middleware/context_helper.go
- [x] T018 Create validation middleware for JSON request validation in backend/internal/middleware/validation_middleware.go

### Error Handling & Routing

- [x] T019 Implement error handling wrapper in backend/internal/middleware/error_handler.go with standard error codes
- [x] T019b Create USAGE guide for middleware package in backend/internal/middleware/USAGE.md
- [x] T020 Setup Gin router with middleware chain in backend/internal/server/routes.go
- [x] T030 Implement CORS configuration in backend/internal/server/routes.go

### Services (Skeleton)

- [x] T021 Create UserService interface and implementation in backend/internal/services/user_service.go (Register, Login, GetUser, UpdateProfile, ChangePassword, PromoteToTraveller, DemoteToSimple, MakeAdmin, DeleteUser)
- [x] T021a Rename SQLiteUserService to RelationalUserService for backend-agnostic naming in backend/internal/services/user_service.go
- [x] T021b Create JWTService interface and implementation in backend/internal/services/jwt_service.go (GenerateToken, ValidateToken, ParseToken)
- [x] T022 [P] Create PlanService interface and FULL IMPLEMENTATION in backend/internal/services/plan_service.go (all 20 methods with RelationalPlanService)
- [x] T022b Create NodeService interface and implementation in backend/internal/services/node_service.go (CreateAttraction, CreateTransition, GetNodeByID, ListApprovedNodes, ListNodesByType, ListApprovedAttractions, ListApprovedTransitions, SearchAttractions, ListNodesByCreator, ApproveNode, DisapproveNode, DeleteNode, Count methods)

### Repositories (Interfaces)

- [x] T023 [P] Create UserRepository interface in backend/internal/repositories/user_repository.go
- [x] T023a Remove redundant RemoveAdmin method and consolidate to DemoteToSimple in UserRepository interface
- [x] T023b [P] Create RelationalUserRepository implementation in backend/internal/repositories/user_repository_sqlite.go (with all CRUD operations and optimizations)
- [x] T023c Rename SQLiteUserRepository to RelationalUserRepository for backend-agnostic naming in backend/internal/repositories/user_repository_sqlite.go
- [x] T024a [P] Create PlanRepository interface in backend/internal/repositories/plan_repository.go (all CRUD, linked-list, denormalization methods)
- [x] T024b [P] Create RelationalPlanRepository implementation in backend/internal/repositories/plan_repository_relational.go (all 21 methods with atomic transactions)
- [x] T025 [P] Create NodeRepository interface in backend/internal/repositories/node_repository.go (CreateAttraction, CreateTransition, GetNodeByID, GetAttractionByID, GetTransitionByID, ListApprovedAttractions, ListApprovedTransitions, SearchAttractions, ListNodesByCreator, ApproveNode, DisapproveNode, DeleteNode, Count methods)
- [x] T025a [P] Create RelationalNodeRepository implementation in backend/internal/repositories/node_repository_relational.go (CreateAttraction, CreateTransition, GetNodeByID - partial)

### Models

- [x] T026 Create User model with validation in backend/internal/models/user.go
- [x] T027 [P] Create TravelPlan model with validation in backend/internal/models/travel_plan.go
- [x] T028 [P] Create Node, AttractionNode, TransitionNode models with validation in backend/internal/models/node.go (Node, AttractionNodeDetail, TransitionNodeDetail)
- [x] T029 [P] Create PlanNode model with linked list validation in backend/internal/models/plan_node.go
- [x] T024 [P] Create Comment model with validation in backend/internal/models/comment.go

**Checkpoint**: All foundational backend infrastructure in place; databases migrated; middleware chain working; ready for user story implementations

---

## ✅ READINESS STATUS FOR USER STORY IMPLEMENTATION

**Phase 2 Completion Status**: 100% Complete (33/33 tasks done) 🎉

**All Blocking Tasks Completed** ✅:
1. ✅ T014: Create migration 007 (ratings table) - DONE
2. ✅ T018: Validation middleware - DONE
3. ✅ T020: Router setup with middleware chain - DONE
4. ✅ T030: CORS configuration - DONE

**Current State**:
- ✅ All models fully implemented (User, TravelPlan, PlanNode, Node, Comment)
- ✅ All repositories fully implemented (PlanRepository, NodeRepository, UserRepository)
- ✅ All services fully implemented (PlanService, UserService, NodeService)
- ✅ All migrations 001-007 complete (users, plans, nodes, node_details, plan_nodes, comments, ratings)
- ✅ All middleware working (response envelope, auth, RBAC, validation)
- ✅ Router/Routes fully configured with 40+ endpoints
- ✅ CORS properly configured for frontend integration

**✅ PHASE 2 CHECKPOINT ACHIEVED - READY FOR USER STORY IMPLEMENTATION**

Backend infrastructure fully operational. Ready to implement:
- US1: Browse & Search Travel Plans
- US2: Comment & Rate Travel Plans
- US3: Submit Promotion Requests
- US4: Create Travel Plans with Existing Nodes

---

## Backend Phase 3: User Story 1 - Browse & Search Travel Plans (Priority: P1) 🎯 MVP

**Goal**: Simple users can discover and search published travel plans with filtering and pagination

### Implementation - User Story 1 Backend

- [x] T032 [US1] Implement UserRepository.GetByID and GetByEmail in backend/internal/repositories/user_repository.go ✓ Implemented
- [x] T033 [US1] Implement PlanRepository.List with pagination and filters in backend/internal/repositories/plan_repository.go ✓ Implemented (ListPublishedPlans)
- [x] T034 [US1] Implement PlanRepository.FindByDestination (full-text search) in backend/internal/repositories/plan_repository.go ✓ Implemented (SearchPlans)
- [x] T035 [US1] Implement PlanService.ListPublishedPlans method in backend/internal/services/plan_service.go ✓ Implemented
- [x] T036 [US1] Implement PlanService.SearchPlans method (destination, keyword, sort) in backend/internal/services/plan_service.go ✓ Implemented
- [x] T037 [US1] Create PlanController.BrowsePlans handler in backend/internal/controllers/plan_controller.go ✓ Implemented
- [x] T038 [US1] Create PlanController.SearchPlans handler in backend/internal/controllers/plan_controller.go ✓ Implemented
- [x] T039 [US1] Register routes: GET /api/v1/plans and GET /api/v1/plans/search in backend/internal/server/routes.go ✓ Wired
- [ ] T040 [P] [US1] Create database indexes for browse performance: CREATE INDEX on travel_plans(status, created_at), travel_plans(destination) in backend/internal/database/migrations/008_add_browse_indexes.sql
- [ ] T047 [US1] Contract test for GET /api/v1/plans endpoint in backend/tests/integration/plan_integration_test.go

**Checkpoint**: Backend browse and search endpoints working; indexed for performance

---

## Backend Phase 4: User Story 2 - Comment & Rate Travel Plans (Priority: P1) 🎯 MVP

**Goal**: Simple users can engage with travel plans by leaving comments and ratings

### Implementation - User Story 2 Backend

- [x] T048 [US2] Create Comment model with validation in backend/internal/models/comment.go ✓ Implemented
- [x] T049 [US2] Create Rating model with validation in backend/internal/models/rating.go ✓ Implemented
- [x] T050 [US2] Implement CommentRepository in backend/internal/repositories/comment_repository.go ✓ Implemented (RelationalCommentRepository)
- [x] T051 [US2] Implement RatingRepository in backend/internal/repositories/rating_repository.go ✓ Implemented (RelationalRatingRepository)
- [x] T052 [US2] Create CommentService with Create/List/Update methods in backend/internal/services/comment_service.go ✓ Implemented (RelationalCommentService)
- [x] T053 [US2] Create RatingService with Submit/Update/GetUserRating methods in backend/internal/services/rating_service.go ✓ Implemented (RelationalRatingService)
- [x] T054 [US2] Create CommentController handlers in backend/internal/controllers/comment_controller.go ✓ All 3 methods implemented
- [x] T055 [US2] Create RatingController handlers in backend/internal/controllers/rating_controller.go ✓ All 4 methods implemented
- [x] T056 [US2] Register routes: POST/GET /api/v1/plans/{id}/comments and POST /api/v1/plans/{id}/ratings in backend/internal/server/routes.go ✓ Wired
- [x] T057 [US2] Implement transactional rating update (denormalized rating_count, rating_sum) in RatingRepository ✓ Implemented
- [x] T058 [US2] Implement transactional comment count update in CommentRepository ✓ Implemented
- [x] T068 [US2] Implement @AuthMiddleware guard decorator for comment/rating endpoints (require login) ✓ Wired in routes
- [ ] T066 [US2] Contract test for POST /api/v1/plans/{id}/comments in backend/tests/integration/comment_integration_test.go
- [ ] T067 [US2] Contract test for POST /api/v1/plans/{id}/ratings in backend/tests/integration/rating_integration_test.go

**Checkpoint**: Backend comment and rating endpoints working; authentication enforced

---

## Backend Phase 5: User Story 3 - Submit Travel Plan for Promotion (Priority: P2)

**Goal**: Simple users can submit careers to request elevation to traveller status

### Implementation - User Story 3 Backend

- [ ] T069 [US3] Create PromotionRequest model in backend/internal/models/promotion_request.go
- [ ] T070 [US3] Implement PromotionRequestRepository in backend/internal/repositories/promotion_request_repository.go
- [ ] T071 [US3] Create PromotionService.SubmitRequest method in backend/internal/services/promotion_service.go
- [ ] T072 [US3] Create PromotionController handlers in backend/internal/controllers/promotion_controller.go
- [ ] T073 [US3] Register routes: POST /api/v1/promotions/request and GET /api/v1/promotions/my-status in backend/internal/server/routes.go
- [ ] T077 [US3] Contract test for POST /api/v1/promotions/request in backend/tests/integration/promotion_integration_test.go

**Checkpoint**: Backend promotion request endpoints working

---

## Backend Phase 6: User Story 4 - Create Travel Plans with Existing Nodes (Priority: P1) 🎯 MVP

**Goal**: Traveller users can create travel itineraries by selecting and ordering existing nodes

### Implementation - User Story 4 Backend

- [ ] T079 [P] [US4] Seed system nodes (attractions and transitions) in backend/internal/database/migrations/009_seed_system_nodes.sql (~50 seed INSERT statements)
- [x] T080 [US4] Implement NodeRepository.List with filtering (approved only, or user's own unapproved) in backend/internal/repositories/node_repository.go ✓ Implemented (ListApprovedAttractions, ListApprovedTransitions)
- [x] T081 [US4] Implement NodeRepository.GetByID in backend/internal/repositories/node_repository.go ✓ Implemented
- [x] T082 [US4] Create NodeService.ListApprovedNodes and ListNodesByType in backend/internal/services/node_service.go
- [x] T083 [US4] Create NodeController.ListNodes and NodeController.GetNodeDetail in backend/internal/controllers/node_controller.go ✓ Wired in routes
- [x] T084 [US4] Implement PlanService.CreateDraftPlan with empty nodes list in backend/internal/services/plan_service.go ✓ Implemented (CreatePlan)
- [x] T085 [US4] Implement PlanService.AddNodesToPlan with validation (no 2 consecutive attractions) in backend/internal/services/plan_service.go ✓ Implemented (AddNodeToPlan with ValidatePlanNodeSequence)
- [x] T086 [US4] Implement PlanService.ReorderNodes (updates sequence_position atomically) in backend/internal/services/plan_service.go ✓ Implemented (ReorderNodeInPlan)
- [x] T087 [US4] Implement PlanService.RemoveNodeFromPlan in backend/internal/services/plan_service.go ✓ Implemented
- [x] T088 [US4] Implement PlanService.PublishPlan (draft → published transition) in backend/internal/services/plan_service.go ✓ Implemented
- [x] T089 [US4] Implement PlanService.GetPlanWithNodes (retrieves full linked list) in backend/internal/services/plan_service.go ✓ Implemented (GetPlanNodeDetails)
- [x] T090 [US4] Register @Traveller route: POST /api/v1/plans (create draft) in backend/internal/server/routes.go ✓ Wired
- [ ] T091 [US4] Register @Traveller route: PATCH /api/v1/plans/{id}/nodes (reorder/add/remove) in backend/internal/server/routes.go
- [x] T092 [US4] Register @Traveller route: PATCH /api/v1/plans/{id}/publish in backend/internal/server/routes.go ✓ Wired
- [ ] T101 [US4] Add Traveller role check on plan creation endpoints via @RBACMiddleware
- [ ] T102 [US4] Contract test for POST /api/v1/plans in backend/tests/integration/plan_integration_test.go
- [ ] T103 [US4] Contract test for PATCH /api/v1/plans/{id}/nodes in backend/tests/integration/plan_integration_test.go
- [ ] T104 [US4] Test linked list validation (no 2 consecutive attractions) in backend/tests/unit/plan_service_test.go

**Checkpoint**: Backend plan creation endpoints working; linked list validation enforced

---

## Backend Phase 7: User Story 5 - Create New Nodes (Priority: P2)

**Goal**: Traveller users can create new attraction or transition nodes if existing nodes don't satisfy requirements

### Implementation - User Story 5 Backend

- [x] T105 [US5] Implement NodeRepository.Create in backend/internal/repositories/node_repository.go ✓ Implemented (CreateNodeAndSave, CreateAttractionAndSave, CreateTransitionAndSave)
- [x] T106 [US5] Create NodeService.CreateAttractionNode with validation in backend/internal/services/node_service.go
- [x] T107 [US5] Create NodeService.CreateTransitionNode with validation in backend/internal/services/node_service.go
- [x] T108 [US5] Create NodeController.CreateAttraction handler in backend/internal/controllers/node_controller.go ✓ Wired
- [x] T109 [US5] Create NodeController.CreateTransition handler in backend/internal/controllers/node_controller.go ✓ Wired
- [x] T110 [US5] Register @Traveller routes: POST /api/v1/nodes/attraction and POST /api/v1/nodes/transition in backend/internal/server/routes.go ✓ Wired
- [ ] T115 [US5] Contract test for POST /api/v1/nodes/attraction in backend/tests/integration/node_integration_test.go
- [ ] T116 [US5] Contract test for POST /api/v1/nodes/transition in backend/tests/integration/node_integration_test.go
- [ ] T117 [US5] Verify user-created nodes default to is_approved=false and hidden from public search

**Checkpoint**: Backend node creation endpoints working; approval workflow enforced

---

## Backend Phase 8: User Story 6 - Admin Moderation (Priority: P2)

**Goal**: Admins can review and remove inappropriate travel plans and nodes from the platform

### Implementation - User Story 6 Backend

- [x] T118 [US6] Implement PlanRepository.ListFlaggedPlans in backend/internal/repositories/plan_repository.go ✓ Implemented
- [x] T118 [US6] Implement PlanRepository.SoftDeletePlan in backend/internal/repositories/plan_repository.go ✓ Implemented (DeletePlan)
- [ ] T119 [US6] Create ModerationService.ListFlaggedPlans in backend/internal/services/moderation_service.go
- [ ] T120 [US6] Create ModerationService.DeletePlan (soft-delete) in backend/internal/services/moderation_service.go
- [ ] T121 [US6] Create ModerationService.ListPendingNodes in backend/internal/services/moderation_service.go
- [ ] T122 [US6] Create ModerationService.ApproveNode and RejectNode in backend/internal/services/moderation_service.go
- [x] T123 [US6] Create ModerationController handlers (AdminController) in backend/internal/controllers/admin_controller.go ✓ All methods implemented
- [x] T124 [US6] Register @Admin routes: PATCH /api/v1/admin/plans/{id}/suspend, DELETE /api/v1/admin/plans/{id}, PATCH /api/v1/admin/nodes/{id}/approve, DELETE /api/v1/admin/nodes/{id} in backend/internal/server/routes.go ✓ Wired
- [ ] T130 [US6] Contract test for DELETE /api/v1/admin/plans/{id} in backend/tests/integration/moderation_integration_test.go
- [ ] T131 [US6] Verify soft-delete: comments/ratings persist when plan deleted

**Checkpoint**: Backend moderation endpoints working; admin-only access enforced

---

## Backend Phase 9: Authentication & Authorization (Cross-Cutting)

**Goal**: User login/registration and role-based access control for all user stories

### Implementation - Authentication Backend

- [x] T132 Implement AuthService.Register with bcrypt hashing in backend/internal/services/auth_service.go ✓ Implemented via UserService.Register
- [x] T133 Implement AuthService.Login with JWT generation in backend/internal/services/auth_service.go ✓ Implemented via UserService.Login
- [x] T134 Create AuthController.Register and AuthController.Login in backend/internal/controllers/auth_controller.go ✓ Implemented
- [x] T135 Register routes: POST /api/v1/auth/register and POST /api/v1/auth/login in backend/internal/server/routes.go ✓ Wired
- [ ] T136 Test JWT token generation and verification in backend/tests/unit/auth_service_test.go
- [ ] T144 Contract test for POST /api/v1/auth/register in backend/tests/integration/auth_integration_test.go
- [ ] T145 Contract test for POST /api/v1/auth/login in backend/tests/integration/auth_integration_test.go

**Checkpoint**: Backend authentication endpoints working; JWT generation and validation verified

---

## Backend Phase 10: View Plan Details & Linked List Display

**Goal**: Display complete travel plan with linked list of nodes in correct sequence

### Implementation - Plan Details View Backend

- [x] T146 Implement PlanService.GetPlanWithNodesForView in backend/internal/services/plan_service.go ✓ Implemented (GetPlanNodeDetails)
- [x] T147 Create PlanController.GetPlanDetails handler in backend/internal/controllers/plan_controller.go ✓ Implemented
- [x] T148 Register route: GET /api/v1/plans/{id} (public for published, protected for draft/admin) in backend/internal/server/routes.go ✓ Wired
- [ ] T153 Contract test for GET /api/v1/plans/{id} in backend/tests/integration/plan_integration_test.go

**Checkpoint**: Backend plan details endpoint working; linked list retrieval verified

---

## Backend Phase 11: User Profile & My Plans

**Goal**: Users can view their profile, see their own travel plans, and manage account

### Implementation - User Profile Backend

- [x] T154 Implement UserRepository.GetPlansForUser in backend/internal/repositories/user_repository.go ✓ Implemented (GetPlansByAuthor)
- [x] T155 Create UserController.GetProfile and GetMyPlans handlers in backend/internal/controllers/user_controller.go ✓ Implemented
- [x] T156 Register @Auth routes: GET /api/v1/users/me and GET /api/v1/users/me/plans in backend/internal/server/routes.go ✓ Wired (GET /api/v1/users/:id)
- [ ] T161 Contract test for GET /api/v1/users/me in backend/tests/integration/user_integration_test.go

**Checkpoint**: Backend profile endpoints working; user plan list endpoint verified

---

## Backend Phase 12: Polish & Cross-Cutting Concerns

**Goal**: Complete backend with polish, error handling, and cross-cutting concerns

### Backend Polish

- [ ] T162 Add request logging middleware in backend/internal/middleware/logging_middleware.go
- [ ] T163 Add error recovery middleware in backend/internal/middleware/recovery_middleware.go
- [ ] T164 Add request ID tracing middleware in backend/internal/middleware/trace_middleware.go
- [ ] T165 Setup structured logging (with request context) in backend/internal/server/server.go
- [ ] T166 Add health check endpoint GET /api/v1/health in backend/internal/controllers/health_controller.go
- [ ] T167 Add API versioning documentation endpoint in OpenAPI/Swagger format
- [ ] T168 Implement database transaction helpers for consistency in backend/internal/database/transactions.go
- [ ] T169 Create seed data script for demo travel plans in backend/internal/database/migrations/010_seed_travel_plans.sql
- [ ] T170 Add graceful shutdown handler in backend/cmd/api/main.go
- [ ] T171 Add environment configuration (.env support) in backend/internal/config/config.go

### Backend Testing & Documentation

- [ ] T183 Document API contracts in OpenAPI format in backend/docs/openapi.yaml
- [ ] T184 Add unit tests for all repository implementations in backend/tests/unit/
- [ ] T185 Add unit tests for all service layer business logic in backend/tests/unit/
- [ ] T186 Add integration tests for all API endpoints in backend/tests/integration/
- [ ] T189 Create API usage examples in docs/API_EXAMPLES.md

**Checkpoint**: Backend complete, documented, tested, polished

---

# FRONTEND TASKS

## Frontend Phase 1: Setup (Project Initialization)

**Purpose**: Frontend project structure, dependencies setup - prerequisite for all user stories

- [ ] T002 Setup Vue 3 Composition API project structure: src/{components,pages,services,stores,types,router} in frontend/
- [ ] T004 [P] Install frontend dependencies: Vue 3, Vite, Pinia, Axios, TypeScript, @vue/test-utils in frontend/package.json

**Checkpoint**: Frontend project structure ready; dependencies installed

---

## Frontend Phase 2: Foundational Infrastructure (Blocking Prerequisites)

**Purpose**: Core frontend systems all user stories depend on

### Frontend Services & State Management

- [ ] T031 Connect frontend HTTP client to backend /api endpoint in frontend/src/services/api.ts

**Checkpoint**: Frontend can connect to backend API

---

## Frontend Phase 3: User Story 1 - Browse & Search Travel Plans (Priority: P1) 🎯 MVP

**Goal**: Simple users can discover and search published travel plans with filtering and pagination

### Implementation - User Story 1 Frontend

- [ ] T041 [US1] Create BrowsePage.vue component in frontend/src/pages/BrowsePage.vue
- [ ] T042 [US1] Create PlanCard.vue component for plan preview in frontend/src/components/PlanCard.vue
- [ ] T043 [US1] Implement plan_service.ts with ListPlans and SearchPlans methods in frontend/src/services/plan_service.ts
- [ ] T044 [US1] Implement planStore (Pinia) with browse/search state in frontend/src/stores/plan_store.ts
- [ ] T045 [US1] Setup Vue Router with /browse route in frontend/src/router/index.ts
- [ ] T046 [US1] Add search filters UI (destination, sort, pagination) in BrowsePage.vue

**Checkpoint**: Frontend browse and search pages working; filters and pagination functional

---

## Frontend Phase 4: User Story 2 - Comment & Rate Travel Plans (Priority: P1) 🎯 MVP

**Goal**: Simple users can engage with travel plans by leaving comments and ratings

### Implementation - User Story 2 Frontend

- [ ] T059 [P] [US2] Create ViewPlanPage.vue component with plan details and linked list nodes in frontend/src/pages/ViewPlanPage.vue
- [ ] T060 [P] [US2] Create CommentList.vue component in frontend/src/components/CommentList.vue
- [ ] T061 [P] [US2] Create CommentForm.vue component in frontend/src/components/CommentForm.vue
- [ ] T062 [P] [US2] Create RatingStars.vue component (input/display) in frontend/src/components/RatingStars.vue
- [ ] T063 [US2] Implement comment_service.ts in frontend/src/services/comment_service.ts
- [ ] T064 [US2] Implement rating_service.ts in frontend/src/services/rating_service.ts
- [ ] T065 [US2] Add comment submission and rating UI to ViewPlanPage.vue

**Checkpoint**: Frontend comment and rating UI components working; can submit to backend

---

## Frontend Phase 5: User Story 3 - Submit Travel Plan for Promotion (Priority: P2)

**Goal**: Simple users can submit their plans to request elevation to traveller status

### Implementation - User Story 3 Frontend

- [ ] T074 [US3] Create PromotionRequestForm.vue in frontend/src/components/PromotionRequestForm.vue
- [ ] T075 [US3] Add promotion request submission UI to user menu/profile in frontend/src/components/Navigation.vue
- [ ] T076 [US3] Implement promotion_service.ts in frontend/src/services/promotion_service.ts
- [ ] T078 [US3] Add promotion request status display to user profile

**Checkpoint**: Frontend promotion request submission UI working

---

## Frontend Phase 6: User Story 4 - Create Travel Plans with Existing Nodes (Priority: P1) 🎯 MVP

**Goal**: Traveller users can create travel itineraries by selecting and ordering existing nodes

### Implementation - User Story 4 Frontend

- [ ] T093 [US4] Create NodeCard.vue component in frontend/src/components/NodeCard.vue
- [ ] T094 [US4] Create CreatePlanPage.vue component in frontend/src/pages/CreatePlanPage.vue
- [ ] T095 [US4] Create NodeSelector.vue (browse and select nodes) in frontend/src/components/NodeSelector.vue
- [ ] T096 [US4] Create LinkedListEditor.vue (drag-drop reordering) in frontend/src/components/LinkedListEditor.vue
- [ ] T097 [US4] Implement node_service.ts in frontend/src/services/node_service.ts
- [ ] T098 [US4] Add createPlanDraft, addNodesToPlan, reorderNodes to plan_service.ts
- [ ] T099 [US4] Add plan creation flow to CreatePlanPage.vue (browse nodes → add → reorder → publish)
- [ ] T100 [US4] Setup router with /create-plan route in frontend/src/router/index.ts

**Checkpoint**: Frontend plan creation flow working; linked list editor functional

---

## Frontend Phase 7: User Story 5 - Create New Nodes (Priority: P2)

**Goal**: Traveller users can create new attraction or transition nodes if existing nodes don't satisfy requirements

### Implementation - User Story 5 Frontend

- [ ] T111 [US5] Create CreateAttractionNodeForm.vue in frontend/src/components/CreateAttractionNodeForm.vue
- [ ] T112 [US5] Create CreateTransitionNodeForm.vue in frontend/src/components/CreateTransitionNodeForm.vue
- [ ] T113 [US5] Add node creation forms to CreatePlanPage.vue (quick create option)
- [ ] T114 [US5] Implement node creation methods in node_service.ts

**Checkpoint**: Frontend node creation forms working; can submit to backend

---

## Frontend Phase 8: User Story 6 - Admin Moderation (Priority: P2)

**Goal**: Admins can review and remove inappropriate travel plans and nodes from the platform

### Implementation - User Story 6 Frontend

- [ ] T125 [US6] Create AdminDashboard.vue in frontend/src/pages/AdminDashboard.vue
- [ ] T126 [US6] Create FlaggedPlansList.vue component in frontend/src/components/FlaggedPlansList.vue
- [ ] T127 [US6] Create PendingNodesList.vue component in frontend/src/components/PendingNodesList.vue
- [ ] T128 [US6] Implement moderation_service.ts in frontend/src/services/moderation_service.ts
- [ ] T129 [US6] Add admin dashboard route with @Admin guard in frontend/src/router/index.ts

**Checkpoint**: Frontend admin moderation dashboard working

---

## Frontend Phase 9: Authentication & Authorization (Cross-Cutting)

**Goal**: User login/registration and role-based access control for all user stories

### Implementation - Authentication Frontend

- [ ] T137 Create LoginPage.vue in frontend/src/pages/LoginPage.vue
- [ ] T138 Create RegisterPage.vue in frontend/src/pages/RegisterPage.vue
- [ ] T139 Implement auth_service.ts with login/register/logout in frontend/src/services/auth_service.ts
- [ ] T140 Implement authStore (Pinia) with user, token, roles state in frontend/src/stores/auth_store.ts
- [ ] T141 Add authentication to HTTP client (Axios/Fetch) - attach JWT to requests in frontend/src/services/api.ts
- [ ] T142 Create router guards for @AuthRequired, @TravellerOnly, @AdminOnly in frontend/src/router/index.ts
- [ ] T143 Add login/logout links to Navigation.vue

**Checkpoint**: Frontend authentication working; JWT integration, role guards, and login/register UI complete

---

## Frontend Phase 10: View Plan Details & Linked List Display

**Goal**: Display complete travel plan with linked list of nodes in correct sequence

### Implementation - Plan Details View Frontend

- [ ] T149 Create LinkedListViewer.vue component to display nodes in sequence in frontend/src/components/LinkedListViewer.vue
- [ ] T150 Add node navigation (prev/next buttons) to LinkedListViewer.vue
- [ ] T151 Update ViewPlanPage.vue to display linked list with LinkedListViewer component
- [ ] T152 Add plan metadata (author, rating, destination, created date) to ViewPlanPage.vue

**Checkpoint**: Frontend plan details display complete with linked list viewer

---

## Frontend Phase 11: User Profile & My Plans

**Goal**: Users can view their profile, see their own travel plans, and manage account

### Implementation - User Profile Frontend

- [ ] T157 Create ProfilePage.vue in frontend/src/pages/ProfilePage.vue
- [ ] T158 Display user info and list of user's travel plans in ProfilePage.vue
- [ ] T159 Add profile link to Navigation.vue
- [ ] T160 Setup router with /profile route in frontend/src/router/index.ts

**Checkpoint**: Frontend user profile and my plans display working

---

## Frontend Phase 12: Polish & Cross-Cutting Concerns

**Goal**: Complete frontend with polish, error handling, and cross-cutting concerns

### Frontend Polish

- [ ] T172 Add global error handling in frontend/src/main.ts (error boundaries)
- [ ] T173 Add loading states and spinners for all API calls
- [ ] T174 Add success/error toast notifications in frontend/src/stores/ui_store.ts
- [ ] T175 Implement responsive design for mobile (375px+), tablet (768px+), desktop in all components
- [ ] T176 Add dark mode support (optional, using CSS variables)
- [ ] T177 Add keyboard navigation (Tab, Enter, Escape) to all interactive components
- [ ] T178 Add ARIA labels for accessibility
- [ ] T179 Create NotFoundPage.vue (404) and ErrorPage.vue (error boundary)
- [ ] T180 Add navigation breadcrumbs for plan creation flow
- [ ] T181 Implement empty state messages (no plans found, etc.) across pages

### Frontend Testing & Documentation

- [ ] T182 Write README.md for quickstart setup (already exists as quickstart.md)
- [ ] T187 Add component tests for all Vue components in frontend/tests/unit/
- [ ] T188 Add user flow integration tests in frontend/tests/integration/

**Checkpoint**: Frontend complete, documented, tested, polished

---

## Dependency Graph & Parallel Execution Strategy

### Critical Path for MVP

**Must Complete in Order**:
1. Backend Phase 1 (Setup) ✅ → All others blocked
2. Backend Phase 2 (Foundational) → Backend Phase 6 (Create Plans - critical for MVP)
3. Backend Phase 9 (Auth) → Protects Phase 4, 6, 7, 8 endpoints
4. Frontend Phase 1 (Setup) → All frontend blocked
5. Frontend Phase 3 (Browse) can start immediately after Frontend Phase 1
6. Frontend Phase 9 (Auth) required before Phase 4, 6, 7, 8 (protected endpoints)

### Parallelizable Work (After Blocking Phases)

**Can Run in Parallel**:
- Backend Phase 3 (US1) and Phase 4 (US2) - independent repositories
- Frontend Phase 3 (US1) and Phase 4 (US2) - independent components
- Backend Phase 5, 7, 8 (P2 features) after Phase 2 foundational complete
- Frontend Phase 5, 7, 8 (P2 features) after Frontend Phase 1 complete

### MVP Scope (Recommended First Release)

**Must Include (P1 User Stories)**:
- Backend Phase 1-2: Setup ✅ & Foundational ✅
- Backend Phase 3: Browse & Search (US1) 
- Backend Phase 4: Comment & Rate (US2) 
- Backend Phase 6: Create Plans (US4) 
- Backend Phase 9: Authentication 
- Backend Phase 10: Plan Details 

- Frontend Phase 1: Setup 
- Frontend Phase 3: Browse & Search (US1) 
- Frontend Phase 4: Comment & Rate (US2) 
- Frontend Phase 6: Create Plans (US4) 
- Frontend Phase 9: Authentication
- Frontend Phase 10: Plan Details

**Can Defer to Phase 2 (P2 Features)**:
- Backend/Frontend Phase 5: Submit Promotion (US3)
- Backend/Frontend Phase 7: Create Nodes (US5)
- Backend/Frontend Phase 8: Admin Moderation (US6)
- Phase 11: Profile (nice-to-have)
- Phase 12: Polish (iterate after MVP launch)

**MVP Estimated Effort**: 
- ~3-4 weeks with 2-3 developers
- Backend Critical Path: Phase 1-2 (2 days) → Phase 6 (3-4 days) = 5-6 days sequential
- Frontend Critical Path: Phase 1 (1 day) → Phase 3 (1 day) → Phase 6 (2 days) = 4 days sequential
- Can parallelize backend Phase 3-5, 9 while frontend Phase 6 in progress

---

## Quality Gates & Verification Checklist

**Before Each Phase**:
- [ ] All new code follows layered architecture (Controller → Service → Model → Repository)
- [ ] All API endpoints have documented contracts in comments or OpenAPI
- [ ] All authentication/authorization requirements met (@AuthMiddleware, @RBACMiddleware)
- [ ] Database migrations tested in fresh database
- [ ] Tests written before implementation (TDD via contract tests)

**Before Merge to Main**:
- [ ] All Phase tasks completed
- [ ] All unit tests pass (`go test ./...` and `npm run test`)
- [ ] All integration tests pass
- [ ] Code review verifies architectural compliance
- [ ] Constitution checklist satisfied
- [ ] API contract compliance verified
- [ ] No SQL injection or security vulnerabilities (code scan)

**Before Production Release**:
- [ ] Performance tests pass (<500ms browse, 100 concurrent users)
- [ ] Load test with production-like data
- [ ] Security audit (password hashing, JWT validation, CORS)
- [ ] Database backup strategy tested
- [ ] Monitoring/logging configured (structured logs)
- [ ] Graceful degradation tested (DB down, etc.)

---

## Summary

**Total Tasks**: 190 actionable, specific, file-path-documented tasks

**Organization**:
- **Backend Tasks**: 111 tasks across 12 phases (setup, foundational, 6 user stories, auth, details, profile, polish)
- **Frontend Tasks**: 79 tasks across 12 phases (setup, foundational, 6 user stories, auth, details, profile, polish)
- Tasks organized by Backend/Frontend for clear team assignment
- User Story context maintained for integration testing
- Blocked dependencies clearly marked (→ indicates sequence requirement)
- [P] indicates parallelizable tasks (different files, no blocking deps)

**MVP Scope**: 
- Backend Critical Path: Phases 1, 2, 3, 4, 6, 9, 10 = ~5-6 days
- Frontend Critical Path: Phases 1, 3, 4, 6, 9, 10 = ~4 days  
- Estimated 3-4 weeks total with 2-3 developers working in parallel
- Complete travel plan creation, browsing, commenting, rating

**Constitutional Compliance**:
- All tasks maintain layered architecture across controller/service/model/repository
- All API endpoints documented in contracts/ before implementation
- Database technology (SQLite3) isolated in repository layer
- Role-based access control (RBAC) enforced at middleware layer
