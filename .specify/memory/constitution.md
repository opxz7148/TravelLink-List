# TravelLink Constitution

<!-- 
Sync Impact Report (v1.1.0 Amendment - DRY Principle):
- Version: 1.0.1 → 1.1.0 (MINOR: new principle VI added for code quality)
- Added: Principle VI. Code Quality & DRY - guides agents and developers to avoid code duplication
- Rationale: Prevents maintenance burden, reduces bugs, improves consistency
- Impact: All future code generation must check for and refactor duplicated patterns
- Date: 2026-04-12
-->

## Core Principles

### I. Layered Architecture
The application MUST maintain strict separation of concerns across four architectural layers:
- **Controller Layer** (HTTP entry point, API routing, middleware): HTTP handlers, request validation, routing, authentication/authorization checks via middleware
- **Service Layer** (business logic): Coarse-grained, use-case-based services; each service encapsulates a distinct business capability; pure business logic free of framework dependencies
- **Model Layer** (domain entities): Domain models instantiated from data without side effects; value objects and aggregates for domain concepts
- **Repository Layer** (data access abstraction): All database operations abstracted behind repository interfaces; enables testability and persistence agnostic business logic

Cross-layer dependencies MUST follow unidirectional flow: Controllers → Services → Models → Repositories. No upward or circular dependencies permitted.

### II. API Request/Response Contracts
All HTTP API endpoints MUST define explicit contracts via request and response schemas:
- Request validation at controller layer (JSON schema or middleware validation)
- Consistent response envelope: `{ "success": bool, "data": T | null, "error": { "code": string, "message": string } | null }`
- HTTP status codes align to semantics: 200/201 (success), 400 (validation), 401 (auth), 403 (forbidden), 404 (not found), 500 (server error)
- All contracts documented in OpenAPI/Swagger specification or code comments

### III. Authentication & Authorization (NON-NEGOTIABLE)
Security operations MUST be enforced via middleware before business logic executes:
- Authentication: JWT or session-based token validation at controller layer via middleware; tokens must include user context (ID, roles, permissions)
- Authorization: Role-based or permission-based checks MUST occur in controller middleware before service invocation
- No business logic layer shall perform authentication/authorization decisions; those belong in controller middleware only
- Sensitive operations (delete, admin actions) MUST verify explicit authorization

### IV. Service-Based Use Cases
Services MUST align to coarse-grained business use cases, not technical operations:
- One service per major user-facing or business operation (e.g., `BookingService`, `PaymentService`, `NotificationService`)
- Services MUST be stateless and reusable across multiple controllers or contexts
- Service methods MUST receive fully-validated domain objects and repositories as parameters
- Complex workflows spanning multiple services MUST be orchestrated at service layer or via a dedicated orchestrator, NOT in controller

### V. Repository Abstraction
Data access MUST be abstracted behind repository interfaces:
- Repositories expose query methods (Get, List, Save, Delete) without exposing underlying persistence mechanism
- All database operations isolated in repository implementations (concrete queries, transaction handling)
- Services and upper layers interact only with repository interfaces; database technology (SQLite3, PostgreSQL, MongoDB) change MUST NOT ripple into business logic
- Repository MUST not contain business logic; it is a pure data accessor

### VI. Code Quality & DRY Principle
Code MUST follow the Don't Repeat Yourself (DRY) principle to minimize duplication and maintenance burden:
- **Pattern Extraction**: When identical or near-identical code blocks appear 2+ times, refactor into a reusable helper, utility function, or shared method
- **Semantic Duplication**: Repeated logic patterns (e.g., validation, authorization checks, status comparisons) MUST be consolidated into single implementation; extract to private helpers or well-named methods
- **Documentation Duplication**: Avoid copy-pasting identical documentation across multiple methods; use shared explanations or link to common patterns
- **Scope of Extraction**: Helpers MUST be named clearly to reflect their purpose; extraction threshold is readability vs. performance trade-off (prefer clarity and maintainability)
- **Interface Consistency**: Related methods in the same entity or service MUST follow identical validation, error handling, and response patterns
- **Code Review**: All pull requests MUST be reviewed for DRY violations; duplications MUST be consolidated before merge
- **Scope Exception**: Documentation and comments may repeat for clarity in API contracts and method signatures (e.g., pagination parameter docs in repository methods); this is accepted for API clarity and IDE hover hints

## Technology Stack
Backend: Go 1.x with Gin web framework  
Frontend: Vue 3 (Composition API preferred)  
Database: SQLite3 (file-based, embeddable; abstracted via repository layer)  
API Protocol: RESTful JSON over HTTP/HTTPS  
Testing: Go built-in testing package (unit); API integration tests via REST client or testing framework  

## Storage & Persistence
- **Primary Store**: SQLite3 embedded database (travellink.db file)
- **Location**: `backend/travellink.db` (included in application binary directory)
- **Schema**: Managed via SQL migrations in `backend/internal/database/migrations/`
- **Advantages**: Single-file deployment, no external DB infrastructure, automatic backup via file copy, suitable for MVP and small-to-medium scale
- **Constraints**: Single-writer limitation noted; acceptable for MVP (~100 concurrent users); horizontal scaling deferred to Phase 2+
- **Migration**: SQLite3 → PostgreSQL/MySQL possible by updating repository layer + connection string; business logic unaffected

## Security & Authentication
- JWT tokens issued on successful login; tokens include standard claims (iss, sub, exp, iat) plus custom claims for roles/permissions
- Token validation enforced in authentication middleware on all protected endpoints
- Role-based access control (RBAC) enforced in authorization middleware before handler execution
- Passwords hashed with bcrypt or equivalent; no plaintext storage
- SQLite3 file secured via OS file permissions; sensitive data (passwords, tokens, API keys) MUST NOT be logged or exposed in error messages
- CORS, CSRF protections configured per deployment environment

## Development Workflow
- All features MUST be developed on feature branches and merged via pull requests
- Code review MUST verify architecture layer compliance and contract definitions
- Tests added for all new features and bug fixes before merge (unit tests, integration tests for cross-layer scenarios)
- Deployment MUST verify all tests pass and no contract violations exist
- Database file (travellink.db) committed to git during MVP; backup strategy defined for production

## Governance
This constitution defines non-negotiable architectural rules. All decisions regarding feature design, component structure, and implementation patterns MUST align with these principles.

**Amendment Process**: Changes to principles MUST be documented with rationale; all contributors MUST approve amendments. Version MUST be incremented: major (principle removal/redefinition), minor (new principle/section), patch (clarifications, technology specifications).

**Compliance Verification**: All pull requests reviewed for principle adherence; violations block merge until corrected.

**Version**: 1.1.0 | **Ratified**: 2026-04-05 | **Last Amended**: 2026-04-12
