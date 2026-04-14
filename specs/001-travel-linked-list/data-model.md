# Data Model: Travel Linked List

**Branch**: `001-travel-linked-list` | **Date**: 2026-04-05 | **Phase**: 1 (Design & Contracts)

## Overview

This document defines all domain entities, their attributes, relationships, and validation rules for the Travel Linked List feature. Entities are organized by functional domain (User Management, Travel Planning, Content Association, Community Engagement, Moderation).

---

## Domain Entities

### User Management Domain

#### Entity: User

Represents a platform user with role-based permissions.

**Fields**:
- `id` (UUID, primary key)
- `email` (string, unique, max 255 chars, validated as RFC 5322)
- `username` (string, unique, max 50 chars, alphanumeric + underscore)
- `password_hash` (string, bcrypt hash, never stored plaintext)
- `role` (enum: `simple` | `traveller` | `admin`, default: `simple`)
- `display_name` (string, max 100 chars, optional, defaults to username)
- `bio` (string, max 500 chars, optional)
- `profile_picture_url` (string, URL validation, optional)
- `is_active` (boolean, default: true, soft-delete via false)
- `created_at` (timestamp, UTC, immutable)
- `updated_at` (timestamp, UTC, on modification)
- `last_login_at` (timestamp, UTC, nullable, on successful auth)

**Validation Rules**:
- Email must be unique; duplicate registration rejected
- Username must be unique; alphanumeric + underscore only; 3-50 chars
- Password: minimum 8 characters, at least one uppercase, one digit, one special char
- Role can only be elevated by admin (simple → traveller), never downgraded except for violations
- Soft-delete: Setting is_active=false hides user from public listings but preserves authored content

**Relationships**:
- One-to-Many: User → TravelPlan (as author)
- One-to-Many: User → Node (as creator)
- One-to-Many: User → Comment (as author)
- One-to-Many: User → Rating (as rater)
- One-to-Many: User → PromotionRequest (as requester)

---

### Travel Planning Domain

#### Entity: TravelPlan

Represents a complete travel itinerary as a linked list of nodes.

**Fields**:
- `id` (UUID, primary key)
- `title` (string, max 150 chars, required)
- `description` (string, max 1000 chars, optional)
- `destination` (string, max 200 chars, required; searchable, not normalized location)
- `author_id` (UUID, foreign key to User, denormalized for query efficiency)
- `status` (enum: `draft` | `published` | `suspended`, default: `draft`)
- `rating_count` (integer, default: 0, denormalized for fast retrieval)
- `rating_sum` (integer, default: 0, used to calculate average)
- `comment_count` (integer, default: 0, denormalized for display)
- `is_deleted_by_admin` (boolean, default: false, soft-delete flag)
- `created_at` (timestamp, UTC, immutable)
- `updated_at` (timestamp, UTC, on modification or node reordering)

**Validation Rules**:
- Title: Required, 1-150 chars, non-empty after trim
- Destination: Required, 1-200 chars, enables search but not strict geocoding
- Status transitions: draft → published (author-triggered) | any → suspended (admin-triggered)
- Only published plans visible to non-authors
- Draft plans visible only to author and admin
- Plan creation requires author_id; immutable after creation

**Denormalization Justification**:
- `rating_count`, `rating_sum`, `comment_count` cached for browse/list queries; updated transactiona\
lly when Rating/Comment created/deleted
- `author_id` denormalized for efficient author filtering without join

**Relationships**:
- Many-to-One: TravelPlan → User (as author)
- One-to-Many: TravelPlan → PlanNode (linked list association)
- One-to-Many: TravelPlan → Comment
- One-to-Many: TravelPlan → Rating

**State Transitions**:
```
draft → published (author via UI)
       ↓
       suspended (admin, user notified)
       ↓
       draft (admin restore, user notified)
```

---

#### Entity: Node (Base/Abstract)

Represents a single point of interest or transition in a travel plan. Uses single-table inheritance with type discriminator.

**Fields**:
- `id` (UUID, primary key)
- `type` (enum: `attraction` | `transition`, discriminator for polymorphism)
- `created_by` (UUID, foreign key to User, creator for moderation)
- `is_approved` (boolean, default: true for system nodes, false for user-created)
- `created_at` (timestamp, UTC, immutable)
- `updated_at` (timestamp, UTC, nullable, updated only by admin approval)

**Validation Rules**:
- `type` must be set on creation; immutable
- `is_approved` defaults to true only if created_by is admin; user-created nodes require approval
- Approval cannot be revoked once set (is_approved can only go false → true)

**Subtype Discrimination**:
- If `type == "attraction"`: attraction-specific fields in `attraction_node_details` table
- If `type == "transition"`: transition-specific fields in `transition_node_details` table

---

#### Entity: AttractionNode

Specialization of Node representing a point of interest.

**Additional Fields**:
- `name` (string, max 200 chars, required)
- `category` (enum: `tourist_attraction` | `restaurant` | `hotel` | `museum` | `park` | `shopping` | `entertainment` | `other`, required)
- `location` (string, max 300 chars, required, e.g., "123 Main St, City, Country" or coordinates)
- `description` (string, max 1000 chars, optional)
- `contact_info` (string, max 200 chars, optional, e.g., phone/email)
- `hours_of_operation` (string, max 200 chars, optional, e.g., "09:00-18:00, Mon-Sun")
- `estimated_visit_duration_minutes` (integer, optional, e.g., 60, 120)

**Validation Rules**:
- Name: Required, 1-200 chars
- Category: Required, must be one of enum values; changeable by creator or admin
- Location: Required, used for display; not geocoded
- Contact info: Optional but if provided, must be valid email or phone format
- Hours: ISO 8601 time format or human-readable text

**Relationships**:
- Extends Node
- Referenced by PlanNode (sequence in travel plans)

**Search/Index**: Name and category indexed for node selection during plan creation

---

#### Entity: TransitionNode

Specialization of Node representing movement between attractions. Stores immutable information about a transportation mode or service.

**Additional Fields**:
- `title` (string, max 200 chars, required, e.g., "Bus Line 5", "M1 Train", "Walking")
- `mode` (enum: `walking` | `car` | `bus` | `train` | `bike` | `taxi` | `flight` | `other`, required)
- `description` (string, max 1000 chars, optional, general information independent of plan, e.g., "Public transit connecting downtown to airport")
- `hours_of_operation` (string, max 200 chars, optional, e.g., "Mon-Fri 6:00-23:00, Sat-Sun 7:00-22:00")
- `route_notes` (string, max 500 chars, optional, e.g., "Take line 5 north, exit at Main St")

**Validation Rules**:
- Title: Required, 1-200 chars, identifies the specific service/line (immutable)
- Mode: Required, must be one of enum values
- Description: Optional, general information independent of any specific plan
- Hours of operation: Optional, applies to the service itself (e.g., public transit schedules)
- Route notes: Optional, max 500 chars
- All transition node fields are immutable after creation (system nodes only)

**Relationships**:
- Extends Node
- Referenced by PlanNode (transition between attractions in linked list)

**Implicit Relationships**:
- TransitionNodes represent movement between attractions; multiple consecutive transitions allowed (e.g., walk to station, take train, exit at destination)
- AttractionNodes cannot be consecutive; each must be separated by at least one transition
- Validated in service layer: no two consecutive attractions only

---

### Content Association Domain

#### Entity: PlanNode

Association entity representing the linked list structure of a TravelPlan with plan-specific customizations.

**Fields**:
- `id` (UUID, primary key)
- `plan_id` (UUID, foreign key to TravelPlan)
- `node_id` (UUID, foreign key to Node)
- `sequence_position` (integer, 1-indexed, required)
- `description` (string, max 500 chars, optional, plan-specific context independent of node, e.g., "Try the house special pasta", "Scenic route with views")
- `estimated_price_cents` (integer, optional, > 0, cost in cents for this leg in the journey, e.g., 1500 = $15.00)
- `duration_minutes` (integer, optional, plan-specific duration customization)
- `created_at` (timestamp, UTC, immutable)
 
**Validation Rules**:
- Composite unique constraint: (plan_id, node_id) no duplicates
- Composite unique constraint: (plan_id, sequence_position) no gaps; positions must be contiguous 1..N
- When inserting at position K, all positions >= K are incremented atomically (transaction)
- When deleting position K, all positions > K are decremented atomically
- Cannot have two consecutive attraction nodes (but multiple consecutive transitions allowed)

**Ordering**:
- Plans are queries as: `SELECT node_id FROM plan_nodes WHERE plan_id=? ORDER BY sequence_position ASC`
- Reordering via UPDATE plan_nodes SET sequence_position=? WHERE plan_id=? AND node_id=?
- Drag-and-drop UI reordering triggers batch update within single transaction

**Relationships**:
- Many-to-One: PlanNode → TravelPlan
- Many-to-One: PlanNode → Node

---

### Community Engagement Domain

#### Entity: Comment

Represents user commentary on a TravelPlan.

**Fields**:
- `id` (UUID, primary key)
- `plan_id` (UUID, foreign key to TravelPlan)
- `author_id` (UUID, foreign key to User)
- `text` (string, max 1000 chars, required, min 1 char)
- `is_deleted_by_admin` (boolean, default: false, soft-delete)
- `created_at` (timestamp, UTC, immutable)
- `updated_at` (timestamp, UTC, nullable, only for non-deleted comments)

**Validation Rules**:
- Text: Required, 1-1000 chars, trimmed
- Plan must exist and not be deleted
- Author must be active user
- One comment per (plan_id, author_id) can be edited; multiple comments allowed per user per plan
- Updates allowed only by author or admin; cannot edit older than 30 days (preserve history intent)

**Relationships**:
- Many-to-One: Comment → TravelPlan
- Many-to-One: Comment → User (as author)

**Denormalization**:
- TravelPlan.comment_count updated on create/delete

---

#### Entity: Rating

Represents user rating of a TravelPlan.

**Fields**:
- `id` (UUID, primary key)
- `plan_id` (UUID, foreign key to TravelPlan)
- `user_id` (UUID, foreign key to User)
- `stars` (integer, 1-5, required)
- `created_at` (timestamp, UTC, immutable)
- `updated_at` (timestamp, UTC, shows when rating last updated)

**Validation Rules**:
- Stars: Required, must be 1, 2, 3, 4, or 5 (integer)
- Unique constraint: (plan_id, user_id) - one rating per user per plan
- Users can change their rating; new rating replaces previous
- Plan must be published or draft (if user is author)

**Relationships**:
- Many-to-One: Rating → TravelPlan
- Many-to-One: Rating → User

**Denormalization**:
- TravelPlan.rating_count and TravelPlan.rating_sum updated on create/update/delete; average = sum/count

---

### Platform Moderation Domain

#### Entity: PromotionRequest

Represents a simple user's request to become a traveller or for plan promotion.

**Fields**:
- `id` (UUID, primary key)
- `user_id` (UUID, foreign key to User requesting)
- `plan_id` (UUID, foreign key to TravelPlan being promoted, nullable if upgrading user role)
- `status` (enum: `pending` | `approved` | `rejected`, default: `pending`)
- `admin_notes` (string, max 500 chars, optional, filled on rejection or approval)
- `created_at` (timestamp, UTC, immutable)
- `reviewed_at` (timestamp, UTC, nullable, set when admin approves/rejects)

**Validation Rules**:
- Status transitions: pending → approved | pending → rejected (not reversible without creating new request)
- Admin can only review once; reviewed_at immutable after set
- Plan optional; if provided, plan must be authored by requesting user

**Relationships**:
- Many-to-One: PromotionRequest → User (as requester)
- Many-to-One: PromotionRequest → TravelPlan (optional, as subject)

---

## Aggregate Roots & Consistency Boundaries

**TravelPlan Aggregate**:
- Aggregate root: TravelPlan
- Members: PlanNode, Comment, Rating, PromotionRequest (for this plan)
- Invariants:
  - Plan has contiguous sequence of PlanNodes (1..N)
  - Cannot have two consecutive attraction nodes (transitions can be consecutive)
  - Ratings and comments persist even if plan deleted (soft-delete only)
  - Only published or admin-suspended plans visible to non-authors

**Node Aggregate**:
- Aggregate root: Node (with AttractionNode or TransitionNode specialization)
- Members: None (leaf aggregate)
- Invariants:
  - Immutable after creation except for is_approved flag
  - User-created nodes require admin approval before public visibility
  - Nodes reusable across multiple plans (no ownership)

**User Aggregate**:
- Aggregate root: User
- Members: Comment, Rating, PromotionRequest authored by user
- Invariants:
  - Email and username unique
  - Role changes only via admin (no self-promotion)
  - Active user (is_active=true) required to author new content

---

## Database Schema Outline

**Tables**:
1. `users` (id PK, email UNIQUE, username UNIQUE, role, ...)
2. `travel_plans` (id PK, author_id FK, status, ...)
3. `nodes` (id PK, type DISCRIMINATOR, created_by FK, is_approved, ...)
4. `attraction_node_details` (node_id PK FK, name, category, ...)
5. `transition_node_details` (node_id PK FK, mode, duration, ...)
6. `plan_nodes` (id PK, plan_id FK, node_id FK, sequence_position UNIQUE w/ plan_id, ...)
7. `comments` (id PK, plan_id FK, author_id FK, is_deleted_by_admin, ...)
8. `ratings` (id PK, plan_id FK, user_id FK UNIQUE w/ plan_id, stars, ...)
9. `promotion_requests` (id PK, user_id FK, plan_id FK nullable, status, ...)

**Indexes**:
- `users`: (email), (username)
- `travel_plans`: (author_id), (status, created_at), (destination)
- `nodes`: (type), (created_by), (is_approved)
- `plan_nodes`: (plan_id, sequence_position), (node_id)
- `comments`: (plan_id, created_at), (author_id)
- `ratings`: (plan_id), (user_id)
- `promotion_requests`: (user_id, status), (plan_id)

---

## Relationship Diagram (ERD-style text)

```
┌─────────────────────┐
│        User         │
├─────────────────────┤
│ id (PK)             │
│ email (UNIQUE)      │
│ username (UNIQUE)   │
│ role                │
│ ...                 │
└──────────┬──────────┘
           │
           ├─────────────────────────────────────────┐
           │                                         │
    ┌──────▼──────────────┐    ┌────────────────────▼─┐
    │   TravelPlan        │    │      Node (base)     │
    ├─────────────────────┤    ├──────────────────────┤
    │ id (PK)             │    │ id (PK)              │
    │ author_id (FK→User) │    │ type (DISCRIMINATOR) │
    │ title               │    │ created_by (FK→User) │
    │ status              │    │ is_approved          │
    │ ...                 │    │ ...                  │
    └──────────┬──────────┘    └─────┬────────────┬───┘
               │                     │            │
               │                     │      ┌─────▼──────────────────┐
               │         ┌───────────┴──┐   │ TransitionNode         │
               │         │ PlanNode     │   │ ├────────────────────┤ │
               │         ├──────────────┤   │ │ mode               │ │
               │         │plan_id (FK)  │   │ │ duration_minutes   │ │
               │         │node_id (FK)  │   │ │ route_notes        │ │
               │         │sequence_pos  │   │ └────────────────────┘ │
               │         └──────────────┘   └────────────────────────┘
               │
         ┌─────┴──────┬──────────┬───────────────────────┐
         │            │          │                       │
    ┌────▼────────┐ ┌─▼──────┐ ┌┴───────────┐ ┌─────────▼────┐
    │ Comment     │ │ Rating │ │Attraction  │ │Promotion     │
    ├─────────────┤ ├────────┤ │Node        │ │Request       │
    │ id (PK)     │ │id (PK) │ ├────────────┤ ├──────────────┤
    │ plan_id (FK)│ │plan_id │ │ name       │ │ user_id (FK) │
    │ author_id   │ │user_id │ │ category   │ │ plan_id (FK) │
    │ text        │ │stars   │ │ location   │ │ status       │
    └─────────────┘ └────────┘ │ hours      │ └──────────────┘
                               └────────────┘
```

---

## Summary

All entities follow the constitution's architecture principles:
- Models are pure domain objects (no framework coupling)
- Repositories abstract all database access to top-level service layer
- Entities validated at boundary (service layer) before persistence
- Soft-delete patterns preserve data while enabling logical removal
- Denormalization used judiciously for common queries (rating/comment counts)
- Single-table inheritance pattern supports polymorphic Node types
- Composite constraints ensure linked list integrity without external validation
