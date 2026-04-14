# Feature Specification: Travel Linked List

**Feature Branch**: `001-travel-linked-list`  
**Created**: 2026-04-05  
**Status**: Draft  
**Input**: Travel plan sharing application with linked list node structure

## Overview

Travel Linked List is a collaborative travel planning application where users can create, share, and organize travel plans as linked lists. Each plan consists of nodes representing either attractions (tourist sites, restaurants, hotels) or transitions (walking, driving, public transit). Users can build travel itineraries by linking these nodes in sequence.

## User Scenarios & Testing

### User Story 1 - Browse & Search Travel Plans (Priority: P1)

Simple users can discover and search existing travel plans shared by travelers.

**Why this priority**: Core discovery mechanism; enables user engagement with published content; essential for MVP

**Independent Test**: Can be fully tested by browsing available travel plans, searching by destination/keyword, viewing plan details, reading comments

**Acceptance Scenarios**:

1. **Given** simple user is on homepage, **When** they browse available travel plans, **Then** they see paginated list of published plans with title, author, destination, and rating
2. **Given** simple user enters search term, **When** they search, **Then** they see filtered plans matching destination, attraction names, or author
3. **Given** simple user views a travel plan, **When** they open it, **Then** they see the linked list nodes (attractions/transitions) in sequence with details (name, description, type, estimated time)
4. **Given** simple user navigates between plan nodes, **When** they click "next" or "previous", **Then** they see the connected node in the linked list chain

---

### User Story 2 - Comment & Rate Travel Plans (Priority: P1)

Simple users can engage with travel plans by leaving comments and ratings to provide feedback.

**Why this priority**: Enables community engagement; critical for MVP to show user interaction; builds trust via ratings

**Independent Test**: Can be fully tested by commenting on a plan, viewing comments, submitting ratings independent of plan creation

**Acceptance Scenarios**:

1. **Given** simple user views a travel plan, **When** they scroll to comments section, **Then** they see existing comments with author, timestamp, and comment text
2. **Given** simple user is logged in, **When** they write and submit a comment, **Then** their comment appears in the comments list with timestamp
3. **Given** simple user views a travel plan, **When** they submit a rating (1-5 stars), **Then** the rating is recorded and average rating updates
4. **Given** simple user has already rated, **When** they update their rating, **Then** the new rating replaces the old one

---

### User Story 3 - Submit Travel Plan for Promotion (Priority: P2)

Simple users can submit their travel plans to request elevation to "Traveller" status or gain visibility.

**Why this priority**: Enables user progression; provides path for content creators; adds growth mechanism

**Independent Test**: Can be fully tested by submitting a private plan for promotion review independent of plan publication

**Acceptance Scenarios**:

1. **Given** simple user has created a private travel plan, **When** they click "Request Promotion", **Then** they see a submission form
2. **Given** simple user submits a plan for promotion, **When** they provide description/rationale, **Then** the submission is queued for admin review
3. **Given** simple user submits a plan, **When** admin approves, **Then** user status upgrades or plan visibility increases
4. **Given** simple user submitted plan, **When** they check status, **Then** they see "Pending Review", "Approved", or "Rejected" with admin feedback

---

### User Story 4 - Create Travel Plan with Existing Nodes (Priority: P1)

Traveller users can create travel itineraries by selecting existing attraction and transition nodes and linking them into a plan.

**Why this priority**: Core feature of travellers; enables content creation; P1 because it's the primary traveller capability

**Independent Test**: Can be fully tested by creating a complete travel plan from existing nodes, viewing it, and publishing independent of node creation

**Acceptance Scenarios**:

1. **Given** traveller is in plan creation view, **When** they search/browse existing attractions, **Then** they see available attraction nodes (name, category, location, rating)
2. **Given** traveller selects an attraction node, **When** they add it to the plan, **Then** it appears as a node in their linked list
3. **Given** traveller has added an attraction, **When** they add a transition node, **Then** the transition appears between two attractions showing estimated travel time
4. **Given** traveller has created a linked list chain, **When** they arrange nodes, **Then** they can reorder or remove nodes; transitions adjust automatically
5. **Given** traveller completes plan editing, **When** they click "Publish", **Then** the plan becomes visible to other users with their author profile

---

### User Story 5 - Create New Nodes for Travel Plans (Priority: P2)

Traveller users can create new attraction or transition nodes if existing nodes don't satisfy their travel plan requirements.

**Why this priority**: Enables content expansion; allows travelers to fill gaps; lower priority than using existing nodes but important for completeness

**Independent Test**: Can be fully tested by creating a new node and using it in a plan independent of existing node features

**Acceptance Scenarios**:

1. **Given** traveller cannot find desired attraction, **When** they click "Create New Attraction", **Then** they see form for name, category (tourist/restaurant/hotel/etc), location, description, contact info
2. **Given** traveller fills attraction details, **When** they save, **Then** node is created and available for their use; pending admin approval for public visibility
3. **Given** traveller needs transition between nodes, **When** they create transition, **Then** they specify mode (walking/car/bus/train), estimated duration, and optional route notes
4. **Given** traveller saves new node, **When** they return to plan creation, **Then** the new node appears in their node selection list

---

### User Story 6 - Admin Moderation (Priority: P2)

Admins can review and remove inappropriate travel plans from the platform.

**Why this priority**: Essential for platform governance; prevents abuse; lower priority than core user features but required for production

**Independent Test**: Can be fully tested by viewing moderation queue, reviewing reports, and removing plans independent of user creation features

**Acceptance Scenarios**:

1. **Given** admin is in moderation dashboard, **When** they view flagged plans, **Then** they see list with report reason, author, creation date, and plan preview
2. **Given** admin reviews a plan, **When** they determine it violates guidelines, **Then** they click "Delete Plan" with optional reason message
3. **Given** admin deletes a plan, **When** deletion is confirmed, **Then** plan is removed and author receives notification
4. **Given** admin addresses multiple reports, **When** they bulk action, **Then** they can delete multiple plans and send notifications simultaneously

---

### Edge Cases

- What happens when a traveller creates a plan with only one node? (Should be valid; linked list can be single node)
- How does the system handle orphaned nodes if a plan is deleted? (Nodes remain; are marked as "unused"; can be reused in other plans)
- What if a user upgrades from simple user to traveller? (All previously private plans remain private; new plans default to published if user is traveller)
- How are deleted plans from other users treated when referenced in comments? (Comments remain; display "Plan Deleted" gracefully)
- Can users delete their own comments? (Yes; admins can also delete inappropriate comments independently)
- What if all nodes in a plan are marked as inappropriate? (Plan enters "suspended" state; author notified; cannot be viewed until resolved)

---

## Requirements

### Functional Requirements

- **FR-001**: System MUST allow simple users to browse published travel plans with title, author, destination summary, and rating
- **FR-002**: System MUST provide search and filter functionality for travel plans by destination, attraction type, date created, or author
- **FR-003**: System MUST allow users to view a complete travel plan as a linked list displaying all nodes (attractions/transitions) in sequence
- **FR-004**: System MUST allow simple users to submit comments on travel plans with timestamp and author attribution
- **FR-005**: System MUST allow users to rate travel plans (1-5 stars) and display average rating
- **FR-006**: System MUST enable simple users to submit their travel plans for promotion review by admins
- **FR-007**: System MUST allow traveller users to create new travel plans by linking existing or newly-created nodes
- **FR-008**: System MUST enable travellers to create new attraction nodes (with fields: name, category, location, description, contact info)
- **FR-009**: System MUST enable travellers to create new transition nodes (with fields: mode, estimated duration, route notes)
- **FR-010**: System MUST support reordering nodes within a plan and automatic transition adjustment
- **FR-011**: System MUST publish/unpublish travel plans; published plans visible to others, private plans only visible to creator
- **FR-012**: System MUST allow admin users to view moderation dashboard with flagged/reported plans
- **FR-013**: System MUST allow admins to delete inappropriate plans with optional reason notification to author
- **FR-014**: System MUST manage user roles (simple user / traveller / admin) with role-based access control
- **FR-015**: System MUST persist user authentication state and enforce login for plan creation/modification

### Key Entities

- **User**: id, email, username, password_hash, role (simple/traveller/admin), created_at, bio, profile_picture
- **TravelPlan**: id, title, description, author_id, destination, status (draft/published/suspended), created_at, updated_at, rating_avg, comment_count
- **Node**: id, type (attraction/transition), created_by (user_id), created_at, is_approved (for user-created nodes)
- **AttractionNode** (extends Node): name, category (tourist/restaurant/hotel/etc), location, description, contact_info, hours_of_operation
- **TransitionNode** (extends Node): mode (walking/car/bus/train), route_notes, start_node_id, end_node_id
- **PlanNode** (association): plan_id, node_id, sequence_position (for linked list order)
- **Comment**: id, plan_id, user_id, text, created_at, updated_at
- **Rating**: id, plan_id, user_id, stars (1-5), created_at
- **PromotionRequest**: id, plan_id, user_id, status (pending/approved/rejected), admin_notes, created_at, reviewed_at

### Non-Functional Requirements

- System MUST support at least 100 concurrent users during browsing/search operations
- Search queries MUST return results within 500ms for typical travel plan datasets (10k plans)
- User authentication MUST be compliant with OWASP guidelines (JWT tokens, bcrypt hashing, secure storage)
- API responses MUST follow documented contract schema with consistent error handling
- UI MUST be responsive and usable on mobile (375px+), tablet (768px+), and desktop (1440px+) viewports
