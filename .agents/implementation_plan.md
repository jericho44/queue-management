# Implementation Plan — Multi-Tenant Queue Management System (SaaS)

Build a production-grade **Multi-Tenant Queue Management System (QMS) SaaS** with Go (Fiber framework), PostgreSQL database, Redis pub/sub & caching, WebSocket realtime communication, background worker, and a React + TypeScript frontend.

## User Review Required

> [!IMPORTANT]
> - **Architecture**: Modular Monolith written in Go 1.25+ (Fiber framework) with Clean Architecture, separating HTTP handlers, service domain logic, repository data access, and models.
> - **Concurrency Assurance**: Queue sequence numbers use PostgreSQL `FOR UPDATE` row locks to prevent duplicate tickets. `Call Next` counter operations use `FOR UPDATE SKIP LOCKED` to guarantee two staff members never receive the same ticket.
> - **Multi-Tenancy**: Strict `organization_id` isolation at database queries, service context, JWT middleware, and WebSocket subscription channels.
> - **Public Interfaces**: Dedicated high-visibility fullscreen Public Display Screen with sound notifications and live WebSocket updates, plus a Customer Public Tracking Page (`/ticket/:publicToken`).
> - **Compliance**: Full adherence to project standard rules (BIGINT internal IDs, UUID public identifiers, standardized API response envelopes with `meta.request_id` and `timestamp`, pagination schemas, soft deletes, UTC timestamps, audit fields).

---

## Architecture & Technology Stack

```text
                           ┌─────────────────────────────────┐
                           │      React + TypeScript SPA     │
                           └────────────────┬────────────────┘
                                            │
                                  REST API + WebSockets
                                            │
                                            ▼
                           ┌─────────────────────────────────┐
                           │           Nginx Proxy           │
                           └────────────────┬────────────────┘
                                            │
                                            ▼
                   ┌──────────────────────────────────────────────────┐
                   │               Go (Fiber) Backend                 │
                   │                                                  │
                   │ ├── Auth & RBAC         ├── Staff & Counter      │
                   │ ├── Organization        ├── Queue Engine         │
                   │ ├── Branch              ├── Notification Worker  │
                   │ ├── Service             ├── Reporting            │
                   │ └── Subscription        └── WebSocket Hub        │
                   └───────────┬──────────────────────────┬───────────┘
                               │                          │
                               ▼                          ▼
                   ┌───────────────────────┐  ┌───────────────────────┐
                   │      PostgreSQL       │  │         Redis         │
                   │                       │  │                       │
                   │ Source of Truth       │  │ Pub/Sub Event Bus     │
                   │ Concurrency Locks     │  │ Cache & Rate Limit    │
                   └───────────────────────┘  └───────────────────────┘
```

---

## Database Design

PostgreSQL normalized schema meeting all standard rules (BIGINT primary keys, UUID public IDs, audit columns `created_at`, `updated_at`, `deleted_at`, `created_by`, `updated_by`, `deleted_by`).

```sql
-- Organizations (Tenants)
CREATE TABLE organizations (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    settings JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    deleted_by BIGINT
);

-- Subscriptions & SaaS Plans
CREATE TABLE subscriptions (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    plan VARCHAR(50) NOT NULL DEFAULT 'FREE', -- FREE, STARTER, BUSINESS, ENTERPRISE
    max_branches INT NOT NULL DEFAULT 1,
    max_counters INT NOT NULL DEFAULT 2,
    max_staff INT NOT NULL DEFAULT 5,
    max_monthly_tickets INT NOT NULL DEFAULT 1000,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    starts_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    deleted_by BIGINT
);

-- Users & RBAC
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    organization_id BIGINT REFERENCES organizations(id),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL, -- SUPER_ADMIN, OWNER, MANAGER, STAFF, RECEPTIONIST
    phone VARCHAR(50),
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    deleted_by BIGINT
);

-- Branches
CREATE TABLE branches (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    address TEXT,
    phone VARCHAR(50),
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    deleted_by BIGINT
);

-- Services
CREATE TABLE services (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    branch_id BIGINT NOT NULL REFERENCES branches(id),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    prefix VARCHAR(10) NOT NULL DEFAULT 'A', -- A, B, C, P
    avg_service_time_sec INT NOT NULL DEFAULT 480,
    priority_weight INT NOT NULL DEFAULT 1,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    deleted_by BIGINT
);

-- Counters
CREATE TABLE counters (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    branch_id BIGINT NOT NULL REFERENCES branches(id),
    counter_number VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'CLOSED', -- CLOSED, OPEN, BUSY, PAUSED
    current_staff_id BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    deleted_by BIGINT
);

-- Counter Services Pivot
CREATE TABLE counter_services (
    counter_id BIGINT NOT NULL REFERENCES counters(id) ON DELETE CASCADE,
    service_id BIGINT NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    PRIMARY KEY (counter_id, service_id)
);

-- Customers
CREATE TABLE customers (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    email VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    deleted_by BIGINT
);

-- Queue Sequences (Pessimistic Locking Table for Concurrency Safety)
CREATE TABLE queue_sequences (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    branch_id BIGINT NOT NULL REFERENCES branches(id),
    service_id BIGINT NOT NULL REFERENCES services(id),
    sequence_date DATE NOT NULL,
    last_number INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(branch_id, service_id, sequence_date)
);

-- Queue Tickets
CREATE TABLE queue_tickets (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    branch_id BIGINT NOT NULL REFERENCES branches(id),
    service_id BIGINT NOT NULL REFERENCES services(id),
    customer_id BIGINT REFERENCES customers(id),
    counter_id BIGINT REFERENCES counters(id),
    staff_id BIGINT REFERENCES users(id),
    ticket_number VARCHAR(50) NOT NULL, -- e.g. A101
    sequence_number INT NOT NULL,
    queue_date DATE NOT NULL,
    priority VARCHAR(50) NOT NULL DEFAULT 'NORMAL', -- NORMAL, PRIORITY, EMERGENCY
    status VARCHAR(50) NOT NULL DEFAULT 'WAITING', -- WAITING, CALLED, SERVING, COMPLETED, SKIPPED, CANCELLED, NO_SHOW, TRANSFERRED
    public_token UUID NOT NULL DEFAULT gen_random_uuid(),
    called_at TIMESTAMPTZ,
    serving_started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    estimated_wait_seconds INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    deleted_by BIGINT
);

-- Queue Ticket Events (State Machine Audit)
CREATE TABLE queue_ticket_events (
    id BIGSERIAL PRIMARY KEY,
    ticket_id BIGINT NOT NULL REFERENCES queue_tickets(id) ON DELETE CASCADE,
    from_status VARCHAR(50),
    to_status VARCHAR(50) NOT NULL,
    actor_id BIGINT REFERENCES users(id),
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Notifications
CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    channel VARCHAR(50) NOT NULL, -- IN_APP, EMAIL, SMS, WHATSAPP
    recipient VARCHAR(255) NOT NULL,
    subject VARCHAR(255),
    content TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING', -- PENDING, SENT, FAILED
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Audit Logs
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    organization_id BIGINT REFERENCES organizations(id),
    branch_id BIGINT REFERENCES branches(id),
    user_id BIGINT REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id BIGINT,
    old_values JSONB,
    new_values JSONB,
    ip_address VARCHAR(50),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

---

## Proposed Changes

### Backend Infrastructure (`backend/`)

#### [NEW] [`go.mod`](file:///c:/laragon/www/queue-management-tenant/backend/go.mod)
#### [NEW] [`cmd/api/main.go`](file:///c:/laragon/www/queue-management-tenant/backend/cmd/api/main.go)
#### [NEW] [`cmd/worker/main.go`](file:///c:/laragon/www/queue-management-tenant/backend/cmd/worker/main.go)
#### [NEW] [`pkg/response/response.go`](file:///c:/laragon/www/queue-management-tenant/backend/pkg/response/response.go)
#### [NEW] [`pkg/database/postgres.go`](file:///c:/laragon/www/queue-management-tenant/backend/pkg/database/postgres.go)
#### [NEW] [`pkg/redis/redis.go`](file:///c:/laragon/www/queue-management-tenant/backend/pkg/redis/redis.go)
#### [NEW] [`pkg/jwt/jwt.go`](file:///c:/laragon/www/queue-management-tenant/backend/pkg/jwt/jwt.go)
#### [NEW] [`pkg/logger/logger.go`](file:///c:/laragon/www/queue-management-tenant/backend/pkg/logger/logger.go)
#### [NEW] [`internal/auth/`](file:///c:/laragon/www/queue-management-tenant/backend/internal/auth/) (Handler, Service, Repository, Model, DTO)
#### [NEW] [`internal/organization/`](file:///c:/laragon/www/queue-management-tenant/backend/internal/organization/)
#### [NEW] [`internal/branch/`](file:///c:/laragon/www/queue-management-tenant/backend/internal/branch/)
#### [NEW] [`internal/service/`](file:///c:/laragon/www/queue-management-tenant/backend/internal/service/)
#### [NEW] [`internal/counter/`](file:///c:/laragon/www/queue-management-tenant/backend/internal/counter/)
#### [NEW] [`internal/queue/`](file:///c:/laragon/www/queue-management-tenant/backend/internal/queue/) (Sequencer, Lock engine, State machine, Counter call-next)
#### [NEW] [`internal/websocket/`](file:///c:/laragon/www/queue-management-tenant/backend/internal/websocket/) (Gateway, Auth, PubSub Subscriber)
#### [NEW] [`internal/notification/`](file:///c:/laragon/www/queue-management-tenant/backend/internal/notification/) (Notifier, Providers)
#### [NEW] [`internal/report/`](file:///c:/laragon/www/queue-management-tenant/backend/internal/report/) (Analytics & Metrics)
#### [NEW] [`internal/subscription/`](file:///c:/laragon/www/queue-management-tenant/backend/internal/subscription/) (Quota enforcement)
#### [NEW] [`internal/audit/`](file:///c:/laragon/www/queue-management-tenant/backend/internal/audit/)

### Frontend Infrastructure (`frontend/`)

#### [NEW] [`package.json`](file:///c:/laragon/www/queue-management-tenant/frontend/package.json)
#### [NEW] [`src/App.tsx`](file:///c:/laragon/www/queue-management-tenant/frontend/src/App.tsx)
#### [NEW] [`src/features/auth/`](file:///c:/laragon/www/queue-management-tenant/frontend/src/features/auth/)
#### [NEW] [`src/features/dashboard/`](file:///c:/laragon/www/queue-management-tenant/frontend/src/features/dashboard/)
#### [NEW] [`src/features/counter/`](file:///c:/laragon/www/queue-management-tenant/frontend/src/features/counter/) (Staff operational dashboard)
#### [NEW] [`src/features/reception/`](file:///c:/laragon/www/queue-management-tenant/frontend/src/features/reception/) (Ticket issuance)
#### [NEW] [`src/features/display/`](file:///c:/laragon/www/queue-management-tenant/frontend/src/features/display/) (Public Display Screen with animations & audio alerts)
#### [NEW] [`src/features/public-ticket/`](file:///c:/laragon/www/queue-management-tenant/frontend/src/features/public-ticket/) (Customer live tracking)
#### [NEW] [`src/features/reports/`](file:///c:/laragon/www/queue-management-tenant/frontend/src/features/reports/)

### Infrastructure & Deployments (`docker/`, `migrations/`)

#### [NEW] [`docker-compose.yml`](file:///c:/laragon/www/queue-management-tenant/docker-compose.yml)
#### [NEW] [`Dockerfile.backend`](file:///c:/laragon/www/queue-management-tenant/Dockerfile.backend)
#### [NEW] [`Dockerfile.frontend`](file:///c:/laragon/www/queue-management-tenant/Dockerfile.frontend)
#### [NEW] [`.env.example`](file:///c:/laragon/www/queue-management-tenant/.env.example)
#### [NEW] [`migrations/000001_init_schema.up.sql`](file:///c:/laragon/www/queue-management-tenant/migrations/000001_init_schema.up.sql)

---

## Verification Plan

### Automated Tests
1. **Unit & State Machine Tests**: `go test ./internal/queue/...` to verify valid vs invalid state transitions (e.g. `WAITING` -> `CALLED` ok, `COMPLETED` -> `SERVING` fails).
2. **Concurrency Tests**: `go test -v -run TestConcurrentTicketGeneration ./tests/concurrency_test.go` (Spawns 100 concurrent goroutines requesting tickets to guarantee zero sequence duplicates).
3. **Counter Lock Concurrency Test**: `go test -v -run TestConcurrentCallNext ./tests/concurrency_test.go` (Spawns 10 counters calling next simultaneously to ensure no ticket is double-assigned).
4. **Tenant Isolation Test**: `go test -v ./tests/tenant_isolation_test.go` (Verifies Org A cannot view or act on Org B tickets).

### Manual Verification
1. Open Staff Dashboard & Public Display Screen in side-by-side browser windows.
2. Click **Call Next** on Counter 1. Verify Display Screen updates ticket number instantly via WebSocket without page refresh and triggers chime audio alert.
3. Test Customer Public Token link on mobile view to confirm live queue position countdown.
