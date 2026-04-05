# Tasks: Travel Linked List

**Input**: Design documents from `/specs/001-travel-linked-list/`  
**Prerequisites**: plan.md (✓), spec.md (✓), data-model.md (✓), contracts/ (✓), research.md (✓)

**Note**: Tasks organized by user story to enable independent implementation and testing. P1 stories are MVP priority.

---

## Format Reference

Each task follows: `- [ ] [TaskID] [P?] [Story] Description with exact file path`

- **[P]**: Can run in parallel (different files, no blocking dependencies)
- **[Story]**: User story label (US1, US2, etc.)
- **ID**: Sequential (T001, T002, T003...)

---

## Phase 1: Setup (Project Initialization)

**Purpose**: Project structure, dependencies, database setup - prerequisite for all user stories

- [x] T001 Setup Go project structure with layered directories: cmd/api, internal/{controllers,services,models,repositories,middleware,database} in backend/
- [ ] T002 Setup Vue 3 Composition API project structure: src/{components,pages,services,stores,types,router} in frontend/
- [x] T003 [P] Install backend dependencies: Go 1.x, Gin, SQLite3 driver (gorm.io/driver/sqlite), BCrypt, JWT in backend/go.mod
- [ ] T004 [P] Install frontend dependencies: Vue 3, Vite, Pinia, Axios, TypeScript, @vue/test-utils in frontend/package.json
- [x] T005 Configure SQLite3 database connection and file location (backend/travellink.db) in backend/internal/database/database.go
- [x] T006 [P] Create database initialization logic in backend/internal/database/database.go (auto-creates DB file on first run)
- [ ] T007 Create database migration framework in backend/internal/database/migrations/

**Checkpoint**: Project structure ready; dependencies installed; SQLite3 database file auto-created on first run

---

## Phase 2: Foundational Infrastructure (Blocking Prerequisites)

**Purpose**: Core systems all user stories depend on - MUST complete before user story implementation

- [ ] T008 Create migration 001: users table (id TEXT PK, email TEXT UNIQUE, username TEXT UNIQUE, role TEXT, password_hash TEXT, is_active BOOL, created_at TIMESTAMP) in backend/internal/database/migrations/001_create_users_table.sql
- [ ] T009 [P] Create migration 002: travel_plans table (id TEXT PK, title TEXT, destination TEXT, author_id TEXT, status TEXT, rating_count INT, rating_sum INT, comment_count INT, created_at TIMESTAMP, updated_at TIMESTAMP) in backend/internal/database/migrations/002_create_travel_plans_table.sql
- [ ] T010 [P] Create migration 003: nodes table with discriminator (id TEXT PK, type TEXT, created_by TEXT, is_approved BOOL, created_at TIMESTAMP) in backend/internal/database/migrations/003_create_nodes_table.sql
- [ ] T011 [P] Create migration 004: attraction_node_details and transition_node_details tables in backend/internal/database/migrations/004_create_node_details_tables.sql
- [ ] T012 [P] Create migration 005: plan_nodes junction table (plan_id TEXT, node_id TEXT, sequence_position INT) in backend/internal/database/migrations/005_create_plan_nodes_table.sql
- [ ] T013 [P] Create migration 006: comments table (id TEXT PK, plan_id TEXT, author_id TEXT, text TEXT, created_at TIMESTAMP) in backend/internal/database/migrations/006_create_comments_table.sql
- [ ] T014 [P] Create migration 007: ratings table (id TEXT PK, plan_id TEXT, user_id TEXT, stars INT, created_at TIMESTAMP) in backend/internal/database/migrations/007_create_ratings_table.sql
- [ ] T015 Create base response envelope middleware in backend/internal/middleware/response_middleware.go
- [ ] T016 [P] Implement AuthMiddleware for JWT token validation in backend/internal/middleware/auth_middleware.go
- [ ] T017 [P] Implement RBACMiddleware for role-based access control in backend/internal/middleware/rbac_middleware.go
- [ ] T018 Create validation middleware for JSON request validation in backend/internal/middleware/validation_middleware.go
- [ ] T019 Implement error handling wrapper in backend/internal/server/errors.go with standard error codes
- [ ] T020 Setup Gin router with middleware chain in backend/internal/server/routes.go
- [ ] T021 Create AuthService interface and errors in backend/internal/services/auth_service.go (skeleton only)
- [ ] T022 [P] Create PlanService interface and errors in backend/internal/services/plan_service.go (skeleton only)
- [ ] T023 [P] Create UserRepository interface in backend/internal/repositories/user_repository.go
- [ ] T024 [P] Create PlanRepository interface in backend/internal/repositories/plan_repository.go
- [ ] T025 [P] Create NodeRepository interface in backend/internal/repositories/node_repository.go
- [ ] T026 Create User model with validation in backend/internal/models/user.go
- [ ] T027 [P] Create TravelPlan model with validation in backend/internal/models/travel_plan.go
- [ ] T028 [P] Create Node, AttractionNode, TransitionNode models with validation in backend/internal/models/node.go
- [ ] T029 [P] Create PlanNode model with linked list validation in backend/internal/models/plan_node.go
- [ ] T030 Implement CORS configuration in backend/internal/server/routes.go
- [ ] T031 Connect frontend HTTP client to backend /api endpoint in frontend/src/services/api.ts

**Checkpoint**: All foundational infrastructure in place; databases migrated; middleware chain working; ready for user story implementations

---

## Phase 3: User Story 1 - Browse & Search Travel Plans (Priority: P1) 🎯 MVP

**Goal**: Simple users can discover and search published travel plans with filtering and pagination

**Independent Test**: 
- Can browse paginated list of published plans without login
- Can search plans by destination/keyword
- Can sort by recent/popular/rating
- Displays plan summary (title, author, rating, comment count)

### Implementation - User Story 1

- [ ] T032 [US1] Implement UserRepository.GetByID and GetByEmail in backend/internal/repositories/user_repository.go
- [ ] T033 [US1] Implement PlanRepository.List with pagination and filters in backend/internal/repositories/plan_repository.go
- [ ] T034 [US1] Implement PlanRepository.FindByDestination (full-text search) in backend/internal/repositories/plan_repository.go
- [ ] T035 [US1] Implement PlanService.ListPublishedPlans method in backend/internal/services/plan_service.go
- [ ] T036 [US1] Implement PlanService.SearchPlans method (destination, keyword, sort) in backend/internal/services/plan_service.go
- [ ] T037 [US1] Create PlanController.BrowsePlans handler in backend/internal/controllers/plan_controller.go
- [ ] T038 [US1] Create PlanController.SearchPlans handler in backend/internal/controllers/plan_controller.go
- [ ] T039 [US1] Register routes: GET /api/v1/plans and GET /api/v1/plans?search=... in backend/internal/server/routes.go
- [ ] T040 [P] [US1] Create database indexes for browse performance: CREATE INDEX on travel_plans(status, created_at), travel_plans(destination) in backend/internal/database/migrations/008_add_browse_indexes.sql
- [ ] T041 [US1] Create BrowsePage.vue component in frontend/src/pages/BrowsePage.vue
- [ ] T042 [US1] Create PlanCard.vue component for plan preview in frontend/src/components/PlanCard.vue
- [ ] T043 [US1] Implement plan_service.ts with ListPlans and SearchPlans methods in frontend/src/services/plan_service.ts
- [ ] T044 [US1] Implement planStore (Pinia) with browse/search state in frontend/src/stores/plan_store.ts
- [ ] T045 [US1] Setup Vue Router with /browse route in frontend/src/router/index.ts
- [ ] T046 [US1] Add search filters UI (destination, sort, pagination) in BrowsePage.vue
- [ ] T047 [US1] Contract test for GET /api/v1/plans endpoint in backend/tests/integration/plan_integration_test.go

**Checkpoint**: Users can browse and search published travel plans independently

---

## Phase 4: User Story 2 - Comment & Rate Travel Plans (Priority: P1) 🎯 MVP

**Goal**: Simple users can engage with travel plans by leaving comments and ratings

**Independent Test**:
- Can submit 1-5 star rating on a published plan (requires login)
- Can leave comments on plans with timestamp
- Can view comments and ratings on plan details
- Rating average recalculates after new rating
- Can update own rating

### Implementation - User Story 2

- [ ] T048 [US2] Create Comment model with validation in backend/internal/models/comment.go
- [ ] T049 [US2] Create Rating model with validation in backend/internal/models/rating.go
- [ ] T050 [US2] Implement CommentRepository in backend/internal/repositories/comment_repository.go
- [ ] T051 [US2] Implement RatingRepository in backend/internal/repositories/rating_repository.go
- [ ] T052 [US2] Create CommentService with Create/List/Update methods in backend/internal/services/comment_service.go
- [ ] T053 [US2] Create RatingService with Submit/Update/GetUserRating methods in backend/internal/services/rating_service.go
- [ ] T054 [US2] Create CommentController handlers in backend/internal/controllers/comment_controller.go
- [ ] T055 [US2] Create RatingController handlers in backend/internal/controllers/rating_controller.go
- [ ] T056 [US2] Register routes: POST/GET /api/v1/plans/{id}/comments and POST /api/v1/plans/{id}/ratings in backend/internal/server/routes.go
- [ ] T057 [US2] Implement transactional rating update (denormalized rating_count, rating_sum) in RatingRepository
- [ ] T058 [US2] Implement transactional comment count update in CommentRepository
- [ ] T059 [P] [US2] Create ViewPlanPage.vue component with plan details and linked list nodes in frontend/src/pages/ViewPlanPage.vue
- [ ] T060 [P] [US2] Create CommentList.vue component in frontend/src/components/CommentList.vue
- [ ] T061 [P] [US2] Create CommentForm.vue component in frontend/src/components/CommentForm.vue
- [ ] T062 [P] [US2] Create RatingStars.vue component (input/display) in frontend/src/components/RatingStars.vue
- [ ] T063 [US2] Implement comment_service.ts in frontend/src/services/comment_service.ts
- [ ] T064 [US2] Implement rating_service.ts in frontend/src/services/rating_service.ts
- [ ] T065 [US2] Add comment submission and rating UI to ViewPlanPage.vue
- [ ] T066 [US2] Contract test for POST /api/v1/plans/{id}/comments in backend/tests/integration/comment_integration_test.go
- [ ] T067 [US2] Contract test for POST /api/v1/plans/{id}/ratings in backend/tests/integration/rating_integration_test.go
- [ ] T068 [US2] Implement @AuthMiddleware guard decorator for comment/rating endpoints (require login)

**Checkpoint**: Authenticated users can comment and rate plans; ratings and comments persist and count updates

---

## Phase 5: User Story 3 - Submit Travel Plan for Promotion (Priority: P2)

**Goal**: Simple users can submit their plans to request elevation to traveller status

**Independent Test**:
- Simple user can submit promotion request with optional description
- Admin sees pending promotion requests in dashboard
- Promotion request shows pending/approved/rejected status
- User receives notification when request processed

### Implementation - User Story 3

- [ ] T069 [US3] Create PromotionRequest model in backend/internal/models/promotion_request.go
- [ ] T070 [US3] Implement PromotionRequestRepository in backend/internal/repositories/promotion_request_repository.go
- [ ] T071 [US3] Create PromotionService.SubmitRequest method in backend/internal/services/promotion_service.go
- [ ] T072 [US3] Create PromotionController handlers in backend/internal/controllers/promotion_controller.go
- [ ] T073 [US3] Register routes: POST /api/v1/promotions/request and GET /api/v1/promotions/my-status in backend/internal/server/routes.go
- [ ] T074 [US3] Create PromotionRequestForm.vue in frontend/src/components/PromotionRequestForm.vue
- [ ] T075 [US3] Add promotion request submission UI to user menu/profile in frontend/src/components/Navigation.vue
- [ ] T076 [US3] Implement promotion_service.ts in frontend/src/services/promotion_service.ts
- [ ] T077 [US3] Contract test for POST /api/v1/promotions/request in backend/tests/integration/promotion_integration_test.go
- [ ] T078 [US3] Add promotion request status display to user profile

**Checkpoint**: Simple users can submit promotion requests; requests queue for admin review

---

## Phase 6: User Story 4 - Create Travel Plans with Existing Nodes (Priority: P1) 🎯 MVP

**Goal**: Traveller users can create travel itineraries by selecting and ordering existing nodes

**Independent Test**:
- Traveller can browse existing attraction and transition nodes
- Can add nodes to plan in sequence (verified linked list constraint: no 2 consecutive attractions)
- Can reorder nodes (drag-and-drop)
- Can remove nodes
- Can save plan as draft
- Can publish plan to make it visible to others

### Implementation - User Story 4

- [ ] T079 [P] [US4] Seed system nodes (attractions and transitions) in backend/internal/database/migrations/009_seed_system_nodes.sql (~50 seed INSERT statements)
- [ ] T080 [US4] Implement NodeRepository.List with filtering (approved only, or user's own unapproved) in backend/internal/repositories/node_repository.go
- [ ] T081 [US4] Implement NodeRepository.GetByID in backend/internal/repositories/node_repository.go
- [ ] T082 [US4] Create NodeService.ListApprovedNodes and ListNodesByType in backend/internal/services/node_service.go
- [ ] T083 [US4] Create NodeController.ListNodes and NodeController.GetNodeDetail in backend/internal/controllers/node_controller.go
- [ ] T084 [US4] Implement PlanService.CreateDraftPlan with empty nodes list in backend/internal/services/plan_service.go
- [ ] T085 [US4] Implement PlanService.AddNodesToPlan with validation (no 2 consecutive attractions) in backend/internal/services/plan_service.go
- [ ] T086 [US4] Implement PlanService.ReorderNodes (updates sequence_position atomically) in backend/internal/services/plan_service.go
- [ ] T087 [US4] Implement PlanService.RemoveNodeFromPlan in backend/internal/services/plan_service.go
- [ ] T088 [US4] Implement PlanService.PublishPlan (draft → published transition) in backend/internal/services/plan_service.go
- [ ] T089 [US4] Implement PlanService.GetPlanWithNodes (retrieves full linked list) in backend/internal/services/plan_service.go
- [ ] T090 [US4] Register @Traveller route: POST /api/v1/plans (create draft) in backend/internal/server/routes.go
- [ ] T091 [US4] Register @Traveller route: PATCH /api/v1/plans/{id}/nodes (reorder/add/remove) in backend/internal/server/routes.go
- [ ] T092 [US4] Register @Traveller route: PATCH /api/v1/plans/{id}/publish in backend/internal/server/routes.go
- [ ] T093 [US4] Create NodeCard.vue component in frontend/src/components/NodeCard.vue
- [ ] T094 [US4] Create CreatePlanPage.vue component in frontend/src/pages/CreatePlanPage.vue
- [ ] T095 [US4] Create NodeSelector.vue (browse and select nodes) in frontend/src/components/NodeSelector.vue
- [ ] T096 [US4] Create LinkedListEditor.vue (drag-drop reordering) in frontend/src/components/LinkedListEditor.vue
- [ ] T097 [US4] Implement node_service.ts in frontend/src/services/node_service.ts
- [ ] T098 [US4] Add createPlanDraft, addNodesToPlan, reorderNodes to plan_service.ts
- [ ] T099 [US4] Add plan creation flow to CreatePlanPage.vue (browse nodes → add → reorder → publish)
- [ ] T100 [US4] Setup router with /create-plan route in frontend/src/router/index.ts
- [ ] T101 [US4] Add Traveller role check on plan creation endpoints via @RBACMiddleware
- [ ] T102 [US4] Contract test for POST /api/v1/plans in backend/tests/integration/plan_integration_test.go
- [ ] T103 [US4] Contract test for PATCH /api/v1/plans/{id}/nodes in backend/tests/integration/plan_integration_test.go
- [ ] T104 [US4] Test linked list validation (no 2 consecutive attractions) in backend/tests/unit/plan_service_test.go

**Checkpoint**: Travellers can create, edit, and publish travel plans from existing nodes; linked list integrity maintained

---

## Phase 7: User Story 5 - Create New Nodes (Priority: P2)

**Goal**: Traveller users can create new attraction or transition nodes if existing nodes don't satisfy requirements

**Independent Test**:
- Traveller can create new attraction node (name, category, location, etc.)
- Traveller can create new transition node (mode, duration, route notes)
- New user-created nodes default to is_approved=false
- Nodes visible to creator immediately; only visible to others after admin approval
- Admin receives notification of pending nodes

### Implementation - User Story 5

- [ ] T105 [US5] Implement NodeRepository.Create in backend/internal/repositories/node_repository.go
- [ ] T106 [US5] Create NodeService.CreateAttractionNode with validation in backend/internal/services/node_service.go
- [ ] T107 [US5] Create NodeService.CreateTransitionNode with validation in backend/internal/services/node_service.go
- [ ] T108 [US5] Create NodeController.CreateAttraction handler in backend/internal/controllers/node_controller.go
- [ ] T109 [US5] Create NodeController.CreateTransition handler in backend/internal/controllers/node_controller.go
- [ ] T110 [US5] Register @Traveller routes: POST /api/v1/nodes/attraction and POST /api/v1/nodes/transition in backend/internal/server/routes.go
- [ ] T111 [US5] Create CreateAttractionNodeForm.vue in frontend/src/components/CreateAttractionNodeForm.vue
- [ ] T112 [US5] Create CreateTransitionNodeForm.vue in frontend/src/components/CreateTransitionNodeForm.vue
- [ ] T113 [US5] Add node creation forms to CreatePlanPage.vue (quick create option)
- [ ] T114 [US5] Implement node creation methods in node_service.ts
- [ ] T115 [US5] Contract test for POST /api/v1/nodes/attraction in backend/tests/integration/node_integration_test.go
- [ ] T116 [US5] Contract test for POST /api/v1/nodes/transition in backend/tests/integration/node_integration_test.go
- [ ] T117 [US5] Verify user-created nodes default to is_approved=false and hidden from public search

**Checkpoint**: Travellers can create new nodes; user-created nodes queue for admin approval

---

## Phase 8: User Story 6 - Admin Moderation (Priority: P2)

**Goal**: Admins can review and remove inappropriate travel plans and nodes from the platform

**Independent Test**:
- Admin can view list of flagged/reported plans
- Admin can delete plan with reason notification to author
- Admin can view pending user-created nodes
- Admin can approve/reject user nodes
- Deleted plans soft-deleted but comments/ratings preserved
- Plan author receives deletion notification

### Implementation - User Story 6

- [ ] T118 [US6] Implement PlanRepository.ListFlaggedPlans in backend/internal/repositories/plan_repository.go
- [ ] T118 [US6] Implement PlanRepository.SoftDeletePlan in backend/internal/repositories/plan_repository.go
- [ ] T119 [US6] Create ModerationService.ListFlaggedPlans in backend/internal/services/moderation_service.go
- [ ] T120 [US6] Create ModerationService.DeletePlan (soft-delete) in backend/internal/services/moderation_service.go
- [ ] T121 [US6] Create ModerationService.ListPendingNodes in backend/internal/services/moderation_service.go
- [ ] T122 [US6] Create ModerationService.ApproveNode and RejectNode in backend/internal/services/moderation_service.go
- [ ] T123 [US6] Create ModerationController handlers in backend/internal/controllers/moderation_controller.go
- [ ] T124 [US6] Register @Admin routes: GET /api/v1/admin/plans/flagged, DELETE /api/v1/admin/plans/{id}, GET /api/v1/admin/nodes/pending in backend/internal/server/routes.go
- [ ] T125 [US6] Create AdminDashboard.vue in frontend/src/pages/AdminDashboard.vue
- [ ] T126 [US6] Create FlaggedPlansList.vue component in frontend/src/components/FlaggedPlansList.vue
- [ ] T127 [US6] Create PendingNodesList.vue component in frontend/src/components/PendingNodesList.vue
- [ ] T128 [US6] Implement moderation_service.ts in frontend/src/services/moderation_service.ts
- [ ] T129 [US6] Add admin dashboard route with @Admin guard in frontend/src/router/index.ts
- [ ] T130 [US6] Contract test for DELETE /api/v1/admin/plans/{id} in backend/tests/integration/moderation_integration_test.go
- [ ] T131 [US6] Verify soft-delete: comments/ratings persist when plan deleted

**Checkpoint**: Admins can moderate inappropriate content; user-created nodes approved by admins

---

## Phase 9: Authentication & Authorization (Cross-Cutting)

**Goal**: User login/registration and role-based access control for all user stories

**Independent Test**:
- New user can register with email, username, password
- Registered user can login and receive JWT token
- Protected endpoints reject requests without valid token
- Token includes user role (simple/traveller/admin)
- Traveller-only endpoints reject simple users
- Admin-only endpoints reject non-admins

### Implementation - Authentication

- [ ] T132 Implement AuthService.Register with bcrypt hashing in backend/internal/services/auth_service.go
- [ ] T133 Implement AuthService.Login with JWT generation in backend/internal/services/auth_service.go
- [ ] T134 Create AuthController.Register and AuthController.Login in backend/internal/controllers/auth_controller.go
- [ ] T135 Register routes: POST /api/v1/auth/register and POST /api/v1/auth/login in backend/internal/server/routes.go
- [ ] T136 Test JWT token generation and verification in backend/tests/unit/auth_service_test.go
- [ ] T137 Create LoginPage.vue in frontend/src/pages/LoginPage.vue
- [ ] T138 Create RegisterPage.vue in frontend/src/pages/RegisterPage.vue
- [ ] T139 Implement auth_service.ts with login/register/logout in frontend/src/services/auth_service.ts
- [ ] T140 Implement authStore (Pinia) with user, token, roles state in frontend/src/stores/auth_store.ts
- [ ] T141 Add authentication to HTTP client (Axios/Fetch) - attach JWT to requests in frontend/src/services/api.ts
- [ ] T142 Create router guards for @AuthRequired, @TravellerOnly, @AdminOnly in frontend/src/router/index.ts
- [ ] T143 Add login/logout links to Navigation.vue
- [ ] T144 Contract test for POST /api/v1/auth/register in backend/tests/integration/auth_integration_test.go
- [ ] T145 Contract test for POST /api/v1/auth/login in backend/tests/integration/auth_integration_test.go

**Checkpoint**: Users can register/login; JWT auth working; role-based access control enforced

---

## Phase 10: View Plan Details & Linked List Display (Foundation for all user stories)

**Goal**: Display complete travel plan with linked list of nodes in correct sequence

**Prerequisite**: Depends on Phase 3 (User Story 1 - browse) and Phase 4 (comments/ratings)

### Implementation - Plan Details View

- [ ] T146 Implement PlanService.GetPlanWithNodesForView in backend/internal/services/plan_service.go
- [ ] T147 Create PlanController.GetPlanDetails handler in backend/internal/controllers/plan_controller.go
- [ ] T148 Register route: GET /api/v1/plans/{id} (public for published, protected for draft/admin) in backend/internal/server/routes.go
- [ ] T149 Create LinkedListViewer.vue component to display nodes in sequence in frontend/src/components/LinkedListViewer.vue
- [ ] T150 Add node navigation (prev/next buttons) to LinkedListViewer.vue
- [ ] T151 Update ViewPlanPage.vue to display linked list with LinkedListViewer component
- [ ] T152 Add plan metadata (author, rating, destination, created date) to ViewPlanPage.vue
- [ ] T153 Contract test for GET /api/v1/plans/{id} in backend/tests/integration/plan_integration_test.go

**Checkpoint**: Users can view complete travel plans with all nodes in linked list order

---

## Phase 11: User Profile & My Plans (Navigation)

**Goal**: Users can view their profile, see their own travel plans, and manage account

### Implementation - User Profile

- [ ] T154 Implement UserRepository.GetPlansForUser in backend/internal/repositories/user_repository.go
- [ ] T155 Create UserController.GetProfile and GetMyPlans handlers in backend/internal/controllers/user_controller.go
- [ ] T156 Register @Auth routes: GET /api/v1/users/me and GET /api/v1/users/me/plans in backend/internal/server/routes.go
- [ ] T157 Create ProfilePage.vue in frontend/src/pages/ProfilePage.vue
- [ ] T158 Display user info and list of user's travel plans in ProfilePage.vue
- [ ] T159 Add profile link to Navigation.vue
- [ ] T160 Setup router with /profile route in frontend/src/router/index.ts
- [ ] T161 Contract test for GET /api/v1/users/me in backend/tests/integration/user_integration_test.go

**Checkpoint**: Users can view their profile and own travel plans

---

## Phase 12: Polish & Cross-Cutting Concerns

**Goal**: Complete feature with polish, error handling, and cross-cutting concerns

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

### Testing & Documentation

- [ ] T182 Write README.md for quickstart setup (already exists as quickstart.md)
- [ ] T183 Document API contracts in OpenAPI format in backend/docs/openapi.yaml
- [ ] T184 Add unit tests for all repository implementations in backend/tests/unit/
- [ ] T185 Add unit tests for all service layer business logic in backend/tests/unit/
- [ ] T186 Add integration tests for all API endpoints in backend/tests/integration/
- [ ] T187 Add component tests for all Vue components in frontend/tests/unit/
- [ ] T188 Add user flow integration tests in frontend/tests/integration/
- [ ] T189 Create architecture decision record (ADR) for linked list design in docs/ADR.md
- [ ] T190 Create API usage examples in docs/API_EXAMPLES.md

**Checkpoint**: Feature complete, documented, tested, polished

---

## Dependency Graph & Parallel Execution Strategy

### Phase Breakdown by Parallelization

**Phase 1 (Setup)**: All tasks can run in parallel (T001-T007)
- Estimated: 1-2 sockets (dev session, setup only)

**Phase 2 (Foundational)**: Mostly parallel, some sequencing
- Sequence: T008 → T009-T014 [P] → T015-T031
- Estimated: 2-3 days parallel work

**Phase 3 (US1)**: Can start after Phase 2
- Backend: T032-T040 parallel after T024
- Frontend: T041-T046 parallel; no backend dependency
- Estimated: 1-2 days parallel

**Phase 4 (US2)**: Can start after Phase 2, independent of US1
- Backend: T048-T057 parallel
- Frontend: T059-T068 parallel after Phase 3 (reuses components)
- Estimated: 1-2 days parallel

**Phase 6 (US4)**: Critical path for MVP
- Backend: T079 (seed) → T080-T091 sequential (repository → service → controller)
- Frontend: T093-T104 parallel
- Estimated: 2-3 days sequential

**Phases 5, 7, 8**: Can start after Phase 6 completes
- Each phase 1-2 days

**Phase 9 (Auth)**: Can start immediately after Phase 2
- Critical for protecting US4, US5, US6 endpoints
- Estimated: 1-2 days

**Phase 11-12**: Final polish after all features complete
- Estimated: 1-2 days

### MVP Scope (Recommended First Release)

**Must Include (P1 User Stories)**:
- Phase 1: Setup ✅
- Phase 2: Foundational ✅
- Phase 3: Browse & Search (US1) ✅
- Phase 4: Comment & Rate (US2) ✅
- Phase 6: Create Plans (US4) ✅
- Phase 9: Authentication ✅
- Phase 10: View Plan Details ✅

**Can Defer to Phase 2 (P2 Features)**:
- Phase 5: Submit Promotion (US3)
- Phase 7: Create Nodes (US5)
- Phase 8: Admin Moderation (US6)
- Phase 11: Profile (nice-to-have)
- Phase 12: Polish (iterate after MVP launch)

**MVP Estimated Effort**: ~3-4 weeks (2-3 developers, parallel work)

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
- [ ] All unit tests pass (`go test ./...`)
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
- 12 phases (1 setup, 1 foundational, 6 user stories, 1 auth, 1 details, 1 profile, 1 polish)
- Organized by user story for independent team-level parallelization
- Blocked dependencies clearly marked (→ indicates sequence requirement)
- [P] indicates parallelizable tasks (different files, no blocking deps)

**MVP Scope**: 
- Phases 1, 2, 3, 4, 6, 9, 10 = Complete travel plan creation, browsing, commenting, rating
- Estimated 3-4 weeks
- 2-3 developers can work independently on different user stories

**Constitutional Compliance**:
- All tasks maintain layered architecture across controller/service/model/repository
- All API endpoints documented in contracts/ before implementation
- Authentication + Authorization enforced on protected endpoints
- Services implement coarse-grained use cases (Auth, Plan, Node, Comment, Rating, Moderation)
- Repositories abstract all database access

**Ready to Execute**: Pick a task, follow the checklist format, implement, test, PR.
