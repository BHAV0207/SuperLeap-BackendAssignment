# Mini Lead CRM — Superleap Backend Intern Assessment

A clean, production-inspired RESTful Lead Management CRM built using Go, Gin, PostgreSQL, GORM, Docker, and Redis.

This project implements:

* Level 1 — Core CRUD + Workflow State Machine
* Level 2 — Bulk Operations with Partial Success Handling
* Level 3 — Redis Caching with Graceful Fallback

---

# Tech Stack

## Language — Go

Go was chosen because of:

* simplicity, as i have worked in Go before.
* strong concurrency support
* excellent performance at scale 
* clean standard library

It also encourages clean architecture and explicit error handling.

---

## Framework — Gin

Gin was chosen because it is:

* lightweight
* production-ready

It provides:

* routing
* middleware support
* JSON serialization

without unnecessary complexity.

---

## Database — PostgreSQL

PostgreSQL was chosen because:
* The main reason I chose SQL is because our schema was not changing frequently and we had a fixed schema for leads 
* works well with structured business workflows 
* it is reliable and production-proven

The Lead entity and workflow transitions fit naturally into a relational model.

---

## ORM — GORM

GORM was used for:

* clean database interaction
* automatic timestamps
* soft deletes
* migrations
* model mapping

while still allowing readable SQL-like behavior.

---

## Cache — Redis

Redis was added for:

* caching read-heavy endpoints
* reducing repeated database lookups
* demonstrating cache invalidation strategies

The application gracefully falls back to database access if Redis is unavailable.

---

# Features

## Level 1 — Core CRM

Implemented:

* Create Lead
* Get All Leads
* Get Lead By ID
* Update Lead
* Delete Lead
* Status Transition API

Additional features:

* input validation
* email validation
* phone validation
* UUID validation
* centralized response helpers
* centralized error handling
* DTO-based request/response structure
* workflow transition validation

---

## Level 2 — Bulk Operations

Implemented:

* Bulk Create Leads
* Bulk Update Leads

Features:

* partial success handling
* per-record error reporting
* independent record processing
* validation reuse from single operations

One invalid record does not fail the entire batch.

---

## Level 3 — Redis Caching

Implemented:

* cache GET `/leads/:id`
* invalidate cache on update
* invalidate cache on delete

Features:

* graceful fallback if Redis unavailable
* automatic cache repopulation
* optional caching layer

---

# Seed Data

The application provides a manual seeding script to populate the database with sample leads.

The seeding process is:

* manual (run `bash seed.sh`)
* idempotent (handled by the API)
* duplicate-safe

Sample leads include all workflow states:

* NEW
* CONTACTED
* QUALIFIED
* CONVERTED
* LOST

---

# Setup Instructions

## Prerequisites

Install:

* Go
* Docker
* Docker Compose

---

# Environment Variables

Create a `.env` file:

```env
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=superleap
DB_SSLMODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379
```

---

# Start PostgreSQL + Redis

```bash
docker compose up -d
```

---

```bash
go run cmd/server/main.go
```

# Seed Database

To populate the database with sample leads while the server is running, use the provided seed script:

```bash
bash seed.sh
```

Server runs on:

```txt
http://localhost:8080
```

---

# Postman Collection

Access the Postman collection to test the APIs:
[Postman Collection Link](https://www.postman.com/bhav0207/workspace/my-workspace/collection/45988199-764f0d51-9398-4bc0-bab2-d2fb27ebf7ce?action=share&creator=45988199)

---

# Project Structure

```txt
.
├── cmd
│   └── server
│       └── main.go
│
├── internal
│   ├── config
│   ├── database
│   ├── dto
│   ├── handlers
│   ├── models
│   ├── repositories
│   ├── routes
│   ├── services
│   ├── utils
│   └── validators
│
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

---

# Architecture

The application follows a layered architecture:

```txt
Handler Layer
↓
Service Layer
↓
Repository Layer
↓
Database
```

---

## Handler Layer

Responsible for:

* request parsing
* validation handling
* HTTP status codes
* response formatting

Handlers do NOT contain business logic.

---

## Service Layer

Responsible for:

* business logic
* workflow rules
* status transition validation
* cache orchestration

This layer acts as the core business layer of the application.

---

## Repository Layer

Responsible for:

* database operations
* persistence logic
* querying

Repositories are intentionally kept database-focused.

---

## DTO Layer

DTOs were used to:

* separate API contracts from database models
* support request validation
* prevent leaking internal model structure

Separate DTOs were created for:

* create requests
* update requests
* status updates
* responses
* bulk operations

---

# Lead Workflow Rules

The following workflow rules are enforced:

```txt
NEW → CONTACTED → QUALIFIED → CONVERTED
 ↘ LOST     |-> LOST   |->LOST    
```

Rules:

* A lead starts as `NEW`
* A lead can move forward one step at a time
* A lead can move to `LOST` from any status except `CONVERTED`
* `CONVERTED` and `LOST` are terminal states
* Invalid transitions return HTTP 400

---

# Validation Strategy

Validation is handled at multiple layers.

---

## DTO Validation

Handled using Gin + go-playground/validator.

Includes:

* required fields
* valid email format
* enum validation
* phone validation

Example:

* invalid email rejected before reaching service layer
* invalid status values rejected before business logic execution

---

## Service Validation

Business workflow validation is handled inside the service layer.

Example:

* `NEW → CONVERTED` is rejected
* `LOST → CONTACTED` is rejected

This separates:

* input validation
  from:
* business rule validation

---

# Caching Strategy

Redis caching was implemented only for:

```txt
GET /leads/:id
```

This was intentionally scoped narrowly to:

* keep cache invalidation simple
* avoid stale list/query caches
* reduce unnecessary complexity

---

## Cache Flow

```txt
Request
↓
Check Redis
↓
Cache Hit?
  YES → Return Cached Data
  NO  → Query Database
          ↓
       Store in Redis
          ↓
       Return Response
```

---

## Cache Invalidation

Cache is invalidated when:

* lead updated
* lead status changed
* lead deleted

Invalidation strategy:

* delete cache entry
* repopulate on next read

This approach was chosen because it is:

* simple
* safe
* consistent

---

## Graceful Fallback

If Redis is unavailable:

* the application still works normally
* requests fall back to PostgreSQL
* caching becomes optional

This avoids application failure due to cache outages.

---

# Error Handling

The project uses:

* centralized response helpers
* centralized custom errors
* proper HTTP status codes
* graceful validation handling

Examples:

* `400` → invalid request
* `404` → resource not found
* `500` → unexpected server errors

Custom domain errors are implemented for:

* invalid status transitions


---

# API Endpoints

## Core CRUD

| Method | Endpoint            | Description        |
| ------ | ------------------- | ------------------ |
| POST   | `/leads`            | Create lead        |
| GET    | `/leads`            | Get all leads      |
| GET    | `/leads/:id`        | Get lead by ID     |
| PUT    | `/leads/:id`        | Update lead        |
| DELETE | `/leads/:id`        | Delete lead        |
| PATCH  | `/leads/:id/status` | Update lead status |

---

## Bulk Operations

| Method | Endpoint      | Description       |
| ------ | ------------- | ----------------- |
| POST   | `/leads/bulk` | Bulk create leads |
| PUT    | `/leads/bulk` | Bulk update leads |

---

# Example Requests

## Create Lead

```bash
curl -X POST http://localhost:8080/leads \
-H "Content-Type: application/json" \
-d '{
  "name": "Aman Gupta",
  "email": "aman@example.com",
  "phone": "+91-9876543210",
  "source": "website"
}'
```

---

## Update Status

```bash
curl -X PATCH http://localhost:8080/leads/<id>/status \
-H "Content-Type: application/json" \
-d '{
  "status": "CONTACTED"
}'
```

---

## Bulk Create

```bash
curl -X POST http://localhost:8080/leads/bulk \
-H "Content-Type: application/json" \
-d '{
  "leads": [
    {
      "name": "Lead One",
      "email": "lead1@example.com"
    },
    {
      "name": "Lead Two",
      "email": "lead2@example.com"
    }
  ]
}'
```

---

# Design Decisions

## Why Email Duplication Allowed?  
For historical Relevance 

* it is possible a user might visit few times so rather than deleting the data we can keep the historical data for our analysis later for user behaviour
* thus added a filter to search by email with status 

## Why Full Name Validation?

The assessment specification requires a "Full name of the lead". To satisfy this while maintaining a balance between strictness and flexibility, a "middle ground" approach was taken:

* **Implementation**: The validator ensures the name contains at least two space-separated words (e.g., "Aman Gupta").
* **Reasoning**: This prevents single-word entries (like just "Aman") which are often placeholders in a CRM, while remaining culturally inclusive of various naming patterns.
* **Fallback**: It avoids complex regex that might fail on valid names with special characters or varying lengths.

## Why Layered Architecture?

Layered architecture was chosen to:

* separate responsibilities
* improve maintainability
* reduce coupling
* simplify testing

Each layer has a clear responsibility.

---

## Why DTOs?

DTOs were used to:

* avoid exposing database models directly
* support request validation
* support partial updates cleanly
* keep API contracts explicit

---

## Why Soft Deletes?

Soft deletes were chosen because:

* deleted leads may still be useful historically
* avoids permanent accidental data loss
* supported naturally by GORM

---

## Why Redis Optional?

Redis was intentionally designed as optional because:

* cache outages should not break APIs
* database remains source of truth
* graceful degradation is important in distributed systems

---

## Why Cache Only `GET /leads/:id`?

Caching list endpoints introduces:

* invalidation complexity
* filtering consistency problems
* pagination cache issues

Caching single-record reads provides:

* meaningful performance improvement
* simple invalidation
* safer consistency guarantees

---

## How Concurrent Status Transitions Would Be Handled At Scale

At larger scale, concurrent transitions could cause race conditions.

Possible production solutions:

* optimistic locking
* row-level database locking
* version-based updates
* transactional workflow enforcement

For assignment scope, sequential updates are sufficient.

---