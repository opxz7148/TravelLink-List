# API Contract: Travel Plan Endpoints

**Feature**: Travel Linked List  
**Domain**: Travel Planning & Content Management  
**Version**: 1.0  
**Last Updated**: 2026-04-05

---

## Overview

Travel Plan endpoints handle browsing shared plans, viewing plan details, creating new plans, and publishing. Plans are ordered linked lists of nodes (attractions/transitions).

---

## Endpoint 1: Browse & Search Travel Plans

### Request

```
GET /api/v1/plans?destination=Paris&sort=recent&page=1&limit=20
Authorization: Bearer <access_token> (optional for published plans)
```

**Query Parameters**:
- `destination` (string, optional): Search by destination name (case-insensitive)
- `search` (string, optional): Full-text search on title + description
- `sort` (enum, default: `recent`): `recent` | `popular` | `rating`
  - `recent`: newest plans first (created_at DESC)
  - `popular`: most comments + ratings first
  - `rating`: highest average rating first
- `page` (integer, default: 1): Pagination offset
- `limit` (integer, default: 20, max: 100): Items per page
- `author` (string, optional): Filter by author username

### Response: Success (200 OK)

```json
{
  "success": true,
  "api_version": "1.0",
  "data": {
    "plans": [
      {
        "id": "01234567-89ab-cdef-0123-456789abcdef",
        "title": "Week in Paris: Food & Culture",
        "description": "A 7-day itinerary exploring Parisian museums, cafes, and neighborhoods",
        "destination": "Paris, France",
        "author": {
          "id": "550e8400-e29b-41d4-a716-446655440000",
          "username": "sarah_travel",
          "display_name": "Sarah T."
        },
        "rating_average": 4.5,
        "rating_count": 8,
        "comment_count": 3,
        "node_count": 12,
        "created_at": "2026-03-15T10:30:00Z",
        "status": "published"
      },
      {
        "id": "76543210-fedc-ba98-7654-3210fedcba98",
        "title": "Tokyo Adventure: 5 Days",
        "description": "Modern meets traditional: Shibuya, temples, street food",
        "destination": "Tokyo, Japan",
        "author": {
          "id": "660e8400-e29b-41d4-a716-446655440001",
          "username": "alex_nomad",
          "display_name": "Alex N."
        },
        "rating_average": 4.8,
        "rating_count": 15,
        "comment_count": 5,
        "node_count": 10,
        "created_at": "2026-02-28T14:20:00Z",
        "status": "published"
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 45,
      "total_items": 890,
      "limit": 20
    }
  },
  "error": null,
  "timestamp": "2026-04-05T10:30:00Z"
}
```

### Response: Not Found (200 OK, empty results)

```json
{
  "success": true,
  "api_version": "1.0",
  "data": {
    "plans": [],
    "pagination": {
      "current_page": 1,
      "total_pages": 0,
      "total_items": 0,
      "limit": 20
    }
  },
  "error": null,
  "timestamp": "2026-04-05T10:30:00Z"
}
```

---

## Endpoint 2: Get Travel Plan Details

### Request

```
GET /api/v1/plans/{plan_id}
Authorization: Bearer <access_token> (optional for published plans, required for draft/admin access)
```

**Path Parameters**:
- `plan_id` (UUID, required): The plan's unique ID

### Response: Success (200 OK)

```json
{
  "success": true,
  "api_version": "1.0",
  "data": {
    "plan": {
      "id": "01234567-89ab-cdef-0123-456789abcdef",
      "title": "Week in Paris: Food & Culture",
      "description": "A 7-day itinerary exploring Parisian museums, cafes, and neighborhoods",
      "destination": "Paris, France",
      "author": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "username": "sarah_travel",
        "display_name": "Sarah T.",
        "profile_picture_url": "https://..."
      },
      "rating_average": 4.5,
      "rating_count": 8,
      "comment_count": 3,
      "nodes": [
        {
          "id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
          "type": "attraction",
          "sequence": 1,
          "name": "Louvre Museum",
          "category": "museum",
          "location": "Rue de Rivoli, 75001 Paris",
          "description": "World's largest art museum",
          "hours_of_operation": "09:00-18:00, Closed Tuesdays",
          "estimated_visit_duration_minutes": 120
        },
        {
          "id": "bbbbbbbb-cccc-dddd-eeee-ffffffffffff",
          "type": "transition",
          "sequence": 2,
          "mode": "walking",
          "estimated_duration_minutes": 15,
          "route_notes": "Walk along the Seine toward Pont des Arts"
        },
        {
          "id": "cccccccc-dddd-eeee-ffff-gggggggggggg",
          "type": "attraction",
          "sequence": 3,
          "name": "Café de Flore",
          "category": "restaurant",
          "location": "172 Boulevard Saint-Germain, 75006 Paris",
          "description": "Historic café, famous for pastries",
          "contact_info": "+33 1 45 48 55 26",
          "hours_of_operation": "07:00-01:30 daily"
        }
      ],
      "created_at": "2026-03-15T10:30:00Z",
      "updated_at": "2026-04-01T14:20:00Z",
      "status": "published"
    },
    "user_rating": 5,  // Current user's rating, null if not rated
    "has_commented": true  // Whether current user has commented
  },
  "error": null,
  "timestamp": "2026-04-05T10:30:00Z"
}
```

### Response: Not Found (404)

```json
{
  "success": false,
  "api_version": "1.0",
  "data": null,
  "error": {
    "code": "PLAN_NOT_FOUND",
    "message": "Travel plan with ID '01234567-89ab-cdef-0123-456789abcdef' not found"
  },
  "timestamp": "2026-04-05T10:30:00Z"
}
```

### Response: Forbidden (403 for draft plans not owned by user)

```json
{
  "success": false,
  "api_version": "1.0",
  "data": null,
  "error": {
    "code": "FORBIDDEN",
    "message": "You do not have permission to view this draft plan"
  },
  "timestamp": "2026-04-05T10:30:00Z"
}
```

---

## Endpoint 3: Create Travel Plan

### Request

```
POST /api/v1/plans
Authorization: Bearer <access_token> (REQUIRED, role: traveller or admin)
Content-Type: application/json

{
  "title": "Week in Paris: Food & Culture",
  "description": "A 7-day itinerary exploring Parisian museums, cafes, and neighborhoods",
  "destination": "Paris, France",
  "nodes": [
    {
      "node_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
      "sequence": 1
    },
    {
      "node_id": "bbbbbbbb-cccc-dddd-eeee-ffffffffffff",
      "sequence": 2
    }
  ]
}
```

**Field Validation**:
- `title`: Required, 1-150 chars
- `description`: Optional, max 1000 chars
- `destination`: Required, 1-200 chars
- `nodes`: Array of node references (can be empty for draft)
  - `node_id`: UUID of existing node
  - `sequence`: 1-indexed position in plan (must be contiguous)

**Authorization**: Only `traveller` or `admin` role can create plans.

### Response: Success (201 Created)

```json
{
  "success": true,
  "api_version": "1.0",
  "data": {
    "plan": {
      "id": "01234567-89ab-cdef-0123-456789abcdef",
      "title": "Week in Paris: Food & Culture",
      "destination": "Paris, France",
      "status": "draft",
      "node_count": 2,
      "created_at": "2026-04-05T10:30:00Z",
      "updated_at": "2026-04-05T10:30:00Z"
    }
  },
  "error": null,
  "timestamp": "2026-04-05T10:30:00Z"
}
```

### Response: Validation Error (400)

```json
{
  "success": false,
  "api_version": "1.0",
  "data": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid node sequence: gaps in sequence positions",
    "details": {
      "field": "nodes",
      "reason": "Sequence positions must be contiguous 1..N"
    }
  },
  "timestamp": "2026-04-05T10:30:00Z"
}
```

### Response: Unauthorized (401)

```json
{
  "success": false,
  "api_version": "1.0",
  "data": null,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Only travellers and admins can create plans"
  },
  "timestamp": "2026-04-05T10:30:00Z"
}
```

---

## Endpoint 4: Publish Travel Plan

### Request

```
PATCH /api/v1/plans/{plan_id}/publish
Authorization: Bearer <access_token> (REQUIRED, author or admin)
Content-Type: application/json

{}
```

### Response: Success (200 OK)

```json
{
  "success": true,
  "api_version": "1.0",
  "data": {
    "plan": {
      "id": "01234567-89ab-cdef-0123-456789abcdef",
      "status": "published",
      "updated_at": "2026-04-05T10:35:00Z"
    }
  },
  "error": null,
  "timestamp": "2026-04-05T10:35:00Z"
}
```

### Response: Unauthorized (403)

```json
{
  "success": false,
  "api_version": "1.0",
  "data": null,
  "error": {
    "code": "FORBIDDEN",
    "message": "Only the plan author or admin can publish this plan"
  },
  "timestamp": "2026-04-05T10:30:00Z"
}
```

---

## Endpoint 5: Update Travel Plan Nodes (Reorder)

### Request

```
PATCH /api/v1/plans/{plan_id}/nodes
Authorization: Bearer <access_token> (REQUIRED, author or admin)
Content-Type: application/json

{
  "nodes": [
    {
      "node_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
      "sequence": 1
    },
    {
      "node_id": "cccccccc-dddd-eeee-ffff-gggggggggggg",
      "sequence": 2
    },
    {
      "node_id": "bbbbbbbb-cccc-dddd-eeee-ffffffffffff",
      "sequence": 3
    }
  ]
}
```

**Validation**:
- Sequence must be contiguous 1..N
- No two consecutive attraction nodes (transitions must separate attractions)
- Multiple consecutive transitions allowed (e.g., walk → train, or walk → bus → train)
- All nodes must exist and be approved (or user-created)

### Response: Success (200 OK)

```json
{
  "success": true,
  "api_version": "1.0",
  "data": {
    "plan": {
      "id": "01234567-89ab-cdef-0123-456789abcdef",
      "node_count": 3,
      "updated_at": "2026-04-05T10:35:00Z"
    }
  },
  "error": null,
  "timestamp": "2026-04-05T10:35:00Z"
}
```

---

## Error Codes

| Code | HTTP Status | Meaning |
|------|------------|---------|
| `PLAN_NOT_FOUND` | 404 | Plan does not exist or is not accessible |
| `VALIDATION_ERROR` | 400 | Invalid input (gaps in sequence, consecutive nodes, etc.) |
| `UNAUTHORIZED` | 401 | Invalid or missing authentication token |
| `FORBIDDEN` | 403 | User lacks permission (not author, only traveller can create) |
| `INTERNAL_SERVER_ERROR` | 500 | Unexpected server error |

---

## Summary

Travel plan endpoints follow:
- RESTful conventions (GET for read, POST for create, PATCH for update)
- Consistent response envelope
- Proper status codes (201 for created, 200 for update, 404 for not found)
- Authorization checks (author-only for draft plans, traveller-only for creation)
- Pagination for list endpoints
- Linked list node ordering with validation
