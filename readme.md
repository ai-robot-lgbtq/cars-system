# Cars System

A full-stack **used car trading platform** built as a learning project, with production-grade structure and engineering practices.

> 📘 **Full design specification**: see [`开发说明书.md`](./开发说明书.md) (中文) — architecture, data model, API design, milestones, etc.

## 🚀 Quick Start

```bash
git clone git@github.com:ai-robot-lgbtq/cars-system.git
cd cars-system
cp .env.example .env
docker compose up -d --build
```

Once services are up (≈ 30 seconds), verify:

| URL | Purpose |
|---|---|
| http://localhost | User app (Next.js web) |
| http://localhost/admin | Admin app (Next.js admin) |
| http://localhost/api/v1/health | Backend health check |
| http://localhost:5432 | PostgreSQL (zhparser enabled) |
| http://localhost:6379 | Redis |

Expected response from `/api/v1/health`:
```json
{ "code": 0, "message": "ok", "data": { "status": "ok", "service": "cars-system", "db": "ok" } }
```

## 🛠️ Tech Stack

**Backend**
- Go 1.22 + Gin + GORM
- PostgreSQL 16 (with `zhparser` for Chinese full-text search)
- Redis 7 (asynq task queue + WebSocket pub/sub)

**Frontend**
- Next.js 15 (App Router) + TypeScript 5
- Ant Design 5 + Zustand + TanStack Query
- Two apps: `web` (port 3000) and `admin` (port 3001), shared `packages/shared-types`

**Infrastructure**
- Docker Compose + Nginx (single host)
- GitHub Actions CI (lint + build + test)
- pnpm workspaces (monorepo)

## 📐 Architecture

```
Browser (User / Admin)
        │
        ▼  HTTP/REST + WebSocket
   ┌────────────────────────────────────┐
   │  Nginx (:80)                       │
   │  ├─ /         → web:3000          │
   │  ├─ /admin/   → admin:3001        │
   │  └─ /api/     → backend:8080      │
   └────────────────────────────────────┘
                │
                ▼
   ┌────────────────────────────────────┐
   │  Go Backend (Gin + GORM)           │
   │  ├─ Auth, User, Catalog, Order     │
   │  ├─ Payment (Local Mock + Gateway) │
   │  ├─ Chat (WebSocket Hub)           │
   │  └─ Review, Aftersales, Admin      │
   │                                    │
   │  asynq Worker (separate container) │
   └────────────────────────────────────┘
        │           │
        ▼           ▼
   PostgreSQL    Redis
```

WebSocket Hub uses a broker interface with two implementations:
- `MemoryBroker` (default, single-process, zero dependencies) — for development
- `RedisBroker` (Pub/Sub) — for production multi-process. Switch via `WS_BROKER=memory|redis`.

Payment gateway similarly abstracts `LocalMockGateway` (default) from real integrations like `WeChatPayV3Gateway` (stub ready).

## 📂 Project Structure

```
cars-system/
├── backend/                  # Go backend (Gin + GORM)
│   ├── cmd/{api,worker}/
│   ├── internal/{auth,catalog,order,payment,...}/
│   └── migrations/
├── apps/
│   ├── web/                  # User-facing Next.js app
│   └── admin/                # Admin Next.js app
├── packages/
│   └── shared-types/         # Shared TS types (API, enums, errors)
├── nginx/                    # Reverse proxy config
├── docs/superpowers/         # Specs and implementation plans
├── 开发说明书.md              # Full design spec (Chinese)
└── docker-compose.yml
```

## 🧪 Testing

```bash
# Backend (unit tests for config, db, response, handler)
cd backend && go test ./... -race

# Frontend (Vitest for shared types, web, admin)
pnpm -r test

# Type checking
pnpm -r typecheck
```

## 🛣️ Roadmap

| Phase | Scope | Status |
|---|---|---|
| **P0 Foundation** | docker-compose, backend skeleton, web/admin skeletons, CI, migrations | ✅ Done |
| **P1 User & Catalog** | Auth (email + OAuth), user profiles, car CRUD, image upload, admin audit | ⏳ Pending |
| **P2 Order & Payment** | Order state machine, local mock payment, asynq timeout cancel | ⏳ Pending |
| **P3 Communication** | WebSocket chat, message center, two-way reviews, aftersales | ⏳ Pending |
| **P4 Polish** | SMS verification, OAuth real, full-text search tuning, dashboard, E2E | ⏳ Pending |

## 📖 Learning Goals

This project intentionally covers a wide technical surface for learning purposes:

- **REST API design** with versioning and unified response envelope
- **JWT authentication** + role-based access control (buyer/seller/admin)
- **State machines** for order lifecycle with timeout cancellation
- **Payment gateway abstraction** — interface-first design so future integrations are drop-in
- **WebSocket Hub** with broker abstraction (memory ↔ Redis Pub/Sub)
- **Async tasks** with asynq (delayed + scheduled jobs)
- **Full-text search** in PostgreSQL with `zhparser` for Chinese
- **Monorepo** with pnpm workspaces and shared TypeScript types
- **Docker Compose orchestration** for one-command dev environment

## 📄 License

MIT