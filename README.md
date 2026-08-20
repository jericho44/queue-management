# Multi-Tenant Queue Management System (QMS) SaaS

Production-grade **Multi-Tenant Queue Management System (SaaS)** built with **Go (Fiber framework)**, **PostgreSQL**, **Redis**, **WebSockets**, and a modern **React + TypeScript SPA**.

Designed with Clean Architecture, Domain-Driven Modular Monolith, strict database tenant isolation, pessimistic concurrency locking, and realtime WebSocket event distribution.

---

## Key Features & Architecture

- **Multi-Tenant Isolation**: Every entity is scoped by `organization_id` at the database constraints, repository query filters, JWT middleware context, and WebSocket subscription channels.
- **Concurrency-Safe Sequence Generation**: Queue sequence numbers use PostgreSQL pessimistic row locking (`SELECT ... FOR UPDATE`) on `queue_sequences` to guarantee **zero duplicate ticket numbers** (e.g., `A101`, `A102`), even under 100+ simultaneous requests.
- **Race-Condition-Free Counter Operations**: Counter `Call Next` operations execute within PostgreSQL transactions using `FOR UPDATE SKIP LOCKED`, guaranteeing that 10 counters calling next simultaneously will **never claim the same customer ticket**.
- **Realtime WebSockets & Pub/Sub**: State transitions publish to Redis Pub/Sub channels (`org:{org_id}:branch:{branch_id}`) and update connected Staff Dashboards, Public Display Screens (with synthesized audio chimes), and Customer Live Tracking pages.
- **Strict Controlled State Machine**:
  ```text
  WAITING  ──►  CALLED  ──►  SERVING  ──►  COMPLETED
                   │           │
                   ├── SKIPPED ├── SKIPPED
                   ├── NO_SHOW ├── NO_SHOW
                   └── CANCEL  └── TRANSFERRED
  ```
- **High-Visibility Fullscreen Public Display**: Dedicated display interface with large typography, live ticket announcements, chime audio notifications, and active counter boards.
- **Customer Live Tracking**: Mobile public url (`/ticket/:publicToken`) showing live estimated wait countdown, people ahead, and counter call notifications.

---

## Tech Stack

- **Backend**: Go 1.25+ with Fiber framework, Clean Architecture (Handlers, Services, Repositories, Models, DTOs).
- **Database**: PostgreSQL (Source of truth, normalized, UUID public IDs, BIGINT internal keys, audit columns, soft deletes).
- **Cache & Pub/Sub**: Redis 7+ (Pub/Sub event bus, rate limiting, temporary state caching).
- **Realtime**: WebSocket Gateway with channel authorization (`org:{org_id}:branch:{branch_id}`).
- **Frontend**: React 18, TypeScript, Vite, Tailwind CSS, TanStack Query, Zustand, Lucide icons, Recharts.
- **Containerization**: Docker Compose (`api`, `worker`, `frontend`, `postgres`, `redis`).

---

## Database Schema Highlights

Standard compliant PostgreSQL schema:
- `organizations`, `subscriptions`, `users`, `branches`, `services`, `counters`, `counter_services`
- `queue_sequences` (Concurrency sequencer lock table)
- `queue_tickets` (Main ticket entity with `public_token`)
- `queue_ticket_events` (State transition audit log)
- `notifications`, `announcements`, `audit_logs`

All business tables include: `id`, `uuid`, `created_at`, `updated_at`, `deleted_at`, `created_by`, `updated_by`, `deleted_by`.

---

## API Standard & Response Envelope

All API endpoints return standard response envelopes with unique `request_id` and `timestamp`:

```json
{
  "success": true,
  "message": "Ticket issued successfully",
  "data": {
    "ticket_number": "A101",
    "status": "WAITING",
    "estimated_wait_seconds": 240
  },
  "meta": {
    "request_id": "a3b1c2d4-e5f6-7890-abcd-ef1234567890",
    "timestamp": "2026-08-19T10:55:00Z"
  }
}
```

---

## Running with Docker Compose

1. Clone repository and start containers:
   ```bash
   docker compose up -d
   ```
2. Services will be available at:
   - **Frontend SPA**: `http://localhost`
   - **Backend API**: `http://localhost:8080/api/v1`
   - **Public Display**: `http://localhost/display`
   - **Health Checks**: `http://localhost:8080/health`

---

## Local Development (Without Docker)

### Backend:
```bash
cd backend
go run cmd/api/main.go
```

### Frontend:
```bash
cd frontend
npm install
npm run dev
```

---

## Demo Credentials

- **Owner**: `owner@healthcare.com` / `password123`
- **Manager**: `manager@healthcare.com` / `password123`
- **Staff Counter**: `staff@healthcare.com` / `password123`
- **Receptionist**: `receptionist@healthcare.com` / `password123`
