# Cars System P0 — Foundation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 搭建可一键启动的空骨架 —— `docker compose up` 起 Postgres + Redis + Go 后端 + Next.js 用户端 + Next.js 管理端 + Nginx，`GET /api/v1/health` 返回 `{"code":0,"message":"ok"}`。

**Architecture:** 单仓多包（monorepo），pnpm workspaces 管理前端 + 共享 types；Go 后端单进程跑 HTTP + 同镜像分容器跑 asynq worker；Nginx 统一对外暴露 80 端口；所有依赖通过 docker-compose 编排，本地开发 `docker compose up -d` 即可。

**Tech Stack:**
- Backend: Go 1.22 + Gin + GORM + Viper + zap + golang-migrate
- Frontend: Next.js 15 (App Router) + TypeScript 5 + Ant Design 5 + Zustand + TanStack Query
- Monorepo: pnpm workspaces
- Infra: Docker + Docker Compose + Nginx + PostgreSQL 16 (with zhparser) + Redis 7
- CI: GitHub Actions (golangci-lint, ESLint, Vitest, build)

---

## Global Constraints

- Go 版本 ≥ 1.22
- Node 版本 ≥ 20
- pnpm 版本 ≥ 9
- Docker + Docker Compose v2
- 所有金额字段必须用 `NUMERIC(12,2)`，禁止 FLOAT
- 所有时间字段必须用 `TIMESTAMPTZ`，禁止 TIMESTAMP
- API 统一响应：`{ code: number, message: string, data?: any }`，成功 code=0
- 所有 commit message 用英文，遵循 conventional commits（`feat:` / `chore:` / `docs:` / `test:` / `fix:`）
- 文件路径统一使用 kebab-case（Go 例外用 snake_case 不强制）
- 所有环境变量从 `.env` 读取，禁止硬编码敏感信息
- 中文用户界面文案允许（学习项目），但代码注释、日志、commit message 用英文

---

## File Structure

整个仓库的目录布局（P0 完成后）：

```
cars-system/
├── backend/
│   ├── cmd/
│   │   ├── api/main.go              # HTTP 入口
│   │   └── worker/main.go           # asynq worker 入口（占位）
│   ├── internal/
│   │   ├── config/                  # viper 配置加载
│   │   │   ├── config.go
│   │   │   └── config_test.go
│   │   ├── database/                # GORM 连接
│   │   │   ├── database.go
│   │   │   └── database_test.go
│   │   └── shared/
│   │       ├── errors/              # 错误码
│   │       │   └── errors.go
│   │       ├── middleware/          # 中间件
│   │       │   ├── logger.go
│   │       │   ├── recovery.go
│   │       │   └── cors.go
│   │       ├── response/            # 统一响应
│   │       │   ├── response.go
│   │       │   └── response_test.go
│   │       └── handler/             # 通用 handler
│   │           ├── health.go
│   │           └── health_test.go
│   ├── migrations/
│   │   └── 000001_init_schema.up.sql    # P0 仅占位，无业务表
│   │   └── 000001_init_schema.down.sql
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── Makefile
├── apps/
│   ├── web/                         # 用户端 :3000
│   │   ├── app/
│   │   │   ├── layout.tsx
│   │   │   ├── page.tsx
│   │   │   └── globals.css
│   │   ├── lib/
│   │   │   ├── api.ts
│   │   │   └── api.test.ts
│   │   ├── next.config.ts
│   │   ├── package.json
│   │   ├── tsconfig.json
│   │   ├── vitest.config.ts
│   │   └── Dockerfile
│   └── admin/                       # 管理端 :3001
│       ├── app/
│       │   ├── layout.tsx
│       │   └── page.tsx
│       ├── lib/
│       │   ├── api.ts
│       │   └── api.test.ts
│       ├── next.config.ts
│       ├── package.json
│       ├── tsconfig.json
│       ├── vitest.config.ts
│       └── Dockerfile
├── packages/
│   └── shared-types/                # 共享 TS 类型
│       ├── src/
│       │   ├── index.ts
│       │   ├── api.ts               # APIResponse, PageResult
│       │   ├── errors.ts            # ErrorCode
│       │   └── enums.ts             # OrderStatus, PaymentStatus, UserRole
│       ├── tests/
│       │   └── enums.test.ts
│       ├── package.json
│       ├── tsconfig.json
│       └── vitest.config.ts
├── nginx/
│   └── nginx.conf
├── .github/
│   └── workflows/
│       └── ci.yml
├── data/
│   └── uploads/                     # 空目录，git 忽略，挂载图片
├── docker-compose.yml
├── .env.example
├── .gitignore
├── .dockerignore
├── README.md
├── 开发说明书.md                     # 已有
└── pnpm-workspace.yaml
```

---

## Task 1: Git 仓库初始化 + 目录结构

**Files:**
- Create: `.gitignore`
- Create: `.dockerignore`
- Create: `.env.example`
- Create: `README.md`
- Create: `data/uploads/.gitkeep`
- Create: `backend/migrations/.gitkeep`
- Create: `backend/cmd/api/main.go`（占位）
- Create: `backend/cmd/worker/main.go`（占位）
- Create: `backend/go.mod`
- Create: `backend/Makefile`
- Create: `pnpm-workspace.yaml`

**Interfaces:**
- Produces: `backend/go.mod`（module path = `github.com/scutech/cars-system/backend`）
- Produces: 仓库根目录结构

---

- [ ] **Step 1.1: 初始化 git 仓库**

```bash
cd /home/scutech/桌面/study/cars-system
git init
git config user.email "you@example.com"   # 如未配置
git config user.name "Your Name"          # 如未配置
git branch -M main
```

预期：`Initialized empty Git repository in /home/scutech/桌面/study/cars-system/.git/`

- [ ] **Step 1.2: 创建 .gitignore**

写入 `/home/scutech/桌面/study/cars-system/.gitignore`：

```gitignore
# --- Go ---
backend/bin/
backend/tmp/
*.exe
*.test
*.out
coverage.txt
cover.out

# --- Node / Next.js ---
node_modules/
.pnpm-store/
.next/
.turbo/
dist/
build/
out/
*.tsbuildinfo

# --- Env ---
.env
.env.local
.env.*.local
!.env.example

# --- IDE ---
.idea/
.vscode/
*.swp
*.swo
.DS_Store
Thumbs.db

# --- Logs ---
*.log
logs/
npm-debug.log*
yarn-debug.log*
pnpm-debug.log*

# --- Data ---
data/uploads/*
!data/uploads/.gitkeep
backend/migrations/*.sql.bak

# --- Tests ---
.playwright/
playwright-report/
test-results/
coverage/

# --- OS ---
.DS_Store
```

- [ ] **Step 1.3: 创建 .dockerignore**

写入 `/home/scutech/桌面/study/cars-system/.dockerignore`：

```
.git
.github
docs
data
*.md
.idea
.vscode
node_modules
.next
**/node_modules
**/.next
**/coverage
**/dist
```

- [ ] **Step 1.4: 创建 .env.example**

写入 `/home/scutech/桌面/study/cars-system/.env.example`：

```bash
# Application
APP_ENV=development
APP_PORT=8080

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=cars
DB_PASSWORD=cars_pass
DB_NAME=cars_db
DB_SSLMODE=disable

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=dev-secret-change-me-in-production
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h

# WebSocket
WS_BROKER=memory

# Payment
PAYMENT_GATEWAY=local

# Frontend
NEXT_PUBLIC_API_BASE_URL=http://localhost/api/v1
NEXT_PUBLIC_WS_URL=ws://localhost/api/v1/ws
```

- [ ] **Step 1.5: 创建 pnpm-workspace.yaml**

写入 `/home/scutech/桌面/study/cars-system/pnpm-workspace.yaml`：

```yaml
packages:
  - 'apps/*'
  - 'packages/*'
```

- [ ] **Step 1.6: 创建 backend go.mod**

```bash
cd /home/scutech/桌面/study/cars-system/backend
go mod init github.com/scutech/cars-system/backend
```

预期：创建 `backend/go.mod`，内容含 `module github.com/scutech/cars-system/backend` 和 `go 1.22`

- [ ] **Step 1.7: 创建 backend 占位文件**

写入 `/home/scutech/桌面/study/cars-system/backend/cmd/api/main.go`：

```go
package main

import "fmt"

func main() {
    fmt.Println("api server starting (stub)")
}
```

写入 `/home/scutech/桌面/study/cars-system/backend/cmd/worker/main.go`：

```go
package main

import "fmt"

func main() {
    fmt.Println("worker starting (stub)")
}
```

- [ ] **Step 1.8: 创建 backend/Makefile**

写入 `/home/scutech/桌面/study/cars-system/backend/Makefile`：

```makefile
.PHONY: build run-api run-worker test lint migrate-up migrate-down

build:
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker

run-api:
	go run ./cmd/api

run-worker:
	go run ./cmd/worker

test:
	go test ./... -race

lint:
	golangci-lint run ./...

migrate-up:
	migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

migrate-down:
	migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down 1
```

- [ ] **Step 1.9: 创建占位目录与 .gitkeep**

```bash
cd /home/scutech/桌面/study/cars-system
mkdir -p data/uploads backend/migrations
touch data/uploads/.gitkeep backend/migrations/.gitkeep
```

- [ ] **Step 1.10: 创建 README.md**

写入 `/home/scutech/桌面/study/cars-system/README.md`：

```markdown
# Cars System — 二手车交易平台

A full-stack used car trading platform built with React + TypeScript + Go.

## Quick Start

```bash
cp .env.example .env
docker compose up -d
```

- Web (user): http://localhost
- Web (admin): http://localhost/admin
- API: http://localhost/api/v1
- API health: http://localhost/api/v1/health

## Development

See [开发说明书.md](./开发说明书.md) for full design spec.

## Tech Stack

- Backend: Go 1.22 + Gin + GORM + PostgreSQL 16 + Redis 7
- Frontend: Next.js 15 + TypeScript + Ant Design 5
- Infra: Docker Compose + Nginx

## License

MIT
```

- [ ] **Step 1.11: 第一次 commit**

```bash
cd /home/scutech/桌面/study/cars-system
git add .gitignore .dockerignore .env.example README.md pnpm-workspace.yaml backend/go.mod backend/Makefile backend/cmd/ data/ 开发说明书.md
git commit -m "chore: initialize monorepo skeleton with go backend stub"
```

预期：commit 成功，工作区干净。

---

## Task 2: Docker Compose — Postgres + Redis 基础服务

**Files:**
- Create: `docker-compose.yml`
- Create: `backend/Dockerfile`
- Create: `nginx/nginx.conf`（占位）
- Create: `apps/web/Dockerfile`（占位）
- Create: `apps/admin/Dockerfile`（占位）

**Interfaces:**
- Produces: `docker compose up -d postgres redis` 可启动且健康

---

- [ ] **Step 2.1: 创建 nginx/nginx.conf 占位**

写入 `/home/scutech/桌面/study/cars-system/nginx/nginx.conf`：

```nginx
events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    server {
        listen 80;
        server_name localhost;

        location / {
            return 200 'cars-system nginx placeholder\n';
            add_header Content-Type text/plain;
        }
    }
}
```

- [ ] **Step 2.2: 创建 backend/Dockerfile（多阶段）**

写入 `/home/scutech/桌面/study/cars-system/backend/Dockerfile`：

```dockerfile
# syntax=docker/dockerfile:1.6

# ---- Build stage ----
FROM golang:1.22-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/api ./cmd/api && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/worker ./cmd/worker

# ---- Runtime stage ----
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /out/api /app/api
COPY --from=build /out/worker /app/worker
COPY migrations /app/migrations
EXPOSE 8080
CMD ["/app/api"]
```

- [ ] **Step 2.3: 创建 apps/web/Dockerfile 占位**

写入 `/home/scutech/桌面/study/cars-system/apps/web/Dockerfile`：

```dockerfile
# syntax=docker/dockerfile:1.6
FROM node:20-alpine
WORKDIR /app
# P0 占位 — 后续 Task 添加 pnpm install / build
RUN echo "console.log('web placeholder')" > index.js
EXPOSE 3000
CMD ["node", "index.js"]
```

- [ ] **Step 2.4: 创建 apps/admin/Dockerfile 占位**

写入 `/home/scutech/桌面/study/cars-system/apps/admin/Dockerfile`：

```dockerfile
# syntax=docker/dockerfile:1.6
FROM node:20-alpine
WORKDIR /app
# P0 占位 — 后续 Task 添加 pnpm install / build
RUN echo "console.log('admin placeholder')" > index.js
EXPOSE 3001
CMD ["node", "index.js"]
```

- [ ] **Step 2.5: 创建 docker-compose.yml**

写入 `/home/scutech/桌面/study/cars-system/docker-compose.yml`：

```yaml
services:
  postgres:
    image: ghcr.io/amutu/zhparser:16
    container_name: cars-postgres
    environment:
      POSTGRES_USER: ${DB_USER:-cars}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-cars_pass}
      POSTGRES_DB: ${DB_NAME:-cars_db}
    ports:
      - "5432:5432"
    volumes:
      - pg_data:/var/lib/postgresql/data
      - ./backend/migrations:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-cars}"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: cars-redis
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5

  nginx:
    image: nginx:alpine
    container_name: cars-nginx
    ports:
      - "80:80"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - postgres
      - redis
    profiles: ["full"]   # 默认不启，Task 8 启用

volumes:
  pg_data:
```

- [ ] **Step 2.6: 复制 .env.example 为 .env**

```bash
cd /home/scutech/桌面/study/cars-system
cp .env.example .env
```

- [ ] **Step 2.7: 启动 postgres + redis**

```bash
cd /home/scutech/桌面/study/cars-system
docker compose up -d postgres redis
```

预期：两个容器状态 `healthy` 或 `running`。

- [ ] **Step 2.8: 验证 Postgres 可连接**

```bash
docker exec -it cars-postgres psql -U cars -d cars_db -c "SELECT version();"
```

预期：输出 PostgreSQL 16 版本信息。

- [ ] **Step 2.9: 验证 Redis 可连接**

```bash
docker exec -it cars-redis redis-cli ping
```

预期：`PONG`

- [ ] **Step 2.10: commit**

```bash
cd /home/scutech/桌面/study/cars-system
git add docker-compose.yml backend/Dockerfile apps/web/Dockerfile apps/admin/Dockerfile nginx/nginx.conf
git commit -m "chore: add docker-compose with postgres and redis"
```

---

## Task 3: 共享 TypeScript Types 包

**Files:**
- Create: `packages/shared-types/package.json`
- Create: `packages/shared-types/tsconfig.json`
- Create: `packages/shared-types/vitest.config.ts`
- Create: `packages/shared-types/src/index.ts`
- Create: `packages/shared-types/src/api.ts`
- Create: `packages/shared-types/src/errors.ts`
- Create: `packages/shared-types/src/enums.ts`
- Create: `packages/shared-types/tests/enums.test.ts`

**Interfaces:**
- Produces: 包名 `@cars-system/shared-types`，导出 `APIResponse`, `PageResult`, `ErrorCode`, `OrderStatus`, `PaymentStatus`, `UserRole`

---

- [ ] **Step 3.1: 创建 package.json**

写入 `/home/scutech/桌面/study/cars-system/packages/shared-types/package.json`：

```json
{
  "name": "@cars-system/shared-types",
  "version": "0.1.0",
  "private": true,
  "main": "./src/index.ts",
  "types": "./src/index.ts",
  "scripts": {
    "test": "vitest run",
    "test:watch": "vitest",
    "typecheck": "tsc --noEmit"
  },
  "devDependencies": {
    "typescript": "^5.4.0",
    "vitest": "^1.6.0"
  }
}
```

- [ ] **Step 3.2: 创建 tsconfig.json**

写入 `/home/scutech/桌面/study/cars-system/packages/shared-types/tsconfig.json`：

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "ESNext",
    "moduleResolution": "Bundler",
    "lib": ["ES2022"],
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "noUncheckedIndexedAccess": true,
    "isolatedModules": true,
    "resolveJsonModule": true
  },
  "include": ["src/**/*", "tests/**/*"]
}
```

- [ ] **Step 3.3: 创建 vitest.config.ts**

写入 `/home/scutech/桌面/study/cars-system/packages/shared-types/vitest.config.ts`：

```typescript
import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    environment: 'node',
    globals: false,
  },
})
```

- [ ] **Step 3.4: 创建 src/api.ts**

写入 `/home/scutech/桌面/study/cars-system/packages/shared-types/src/api.ts`：

```typescript
/**
 * Standard API response envelope.
 */
export interface APIResponse<T = unknown> {
  code: number
  message: string
  data?: T
}

/**
 * Pagination result wrapper.
 */
export interface PageResult<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

/**
 * Pagination query params accepted by list endpoints.
 */
export interface PageQuery {
  page?: number
  page_size?: number
}
```

- [ ] **Step 3.5: 创建 src/errors.ts**

写入 `/home/scutech/桌面/study/cars-system/packages/shared-types/src/errors.ts`：

```typescript
/**
 * Centralized error codes. Keep in sync with backend/internal/shared/errors/errors.go.
 *
 * 10xxx  Generic
 * 20xxx  Auth
 * 30xxx  User / 31xxx Catalog / 32xxx Order / 33xxx Payment
 * 34xxx  Chat / 35xxx Review / Aftersales
 * 40xxx  Admin
 */
export const ErrorCode = {
  OK: 0,
  PARAM_INVALID: 10001,
  SYSTEM_ERROR: 10002,

  UNAUTHORIZED: 20001,
  TOKEN_EXPIRED: 20002,
  FORBIDDEN: 20003,

  USER_NOT_FOUND: 30001,
  CAR_NOT_FOUND: 31001,
  CAR_ALREADY_SOLD: 31002,
  ORDER_STATE_INVALID: 32001,
  ORDER_TIMEOUT: 32002,
  PAYMENT_FAILED: 33001,
  REFUND_FAILED: 33002,
} as const

export type ErrorCodeType = (typeof ErrorCode)[keyof typeof ErrorCode]
```

- [ ] **Step 3.6: 创建 src/enums.ts**

写入 `/home/scutech/桌面/study/cars-system/packages/shared-types/src/enums.ts`：

```typescript
/**
 * User roles. Stored as bit flags in DB (1=buyer, 2=seller, 4=admin).
 */
export enum UserRole {
  GUEST = 0,
  BUYER = 1,
  SELLER = 2,
  ADMIN = 4,
}

export function hasRole(userRole: number, required: UserRole): boolean {
  return (userRole & required) === required
}

/**
 * Car status.
 */
export enum CarStatus {
  DRAFT = 0,
  PENDING = 1,
  ONLINE = 2,
  SOLD = 3,
  OFFLINE = 4,
}

/**
 * Order status.
 */
export enum OrderStatus {
  CREATED = 0,
  PAID = 1,
  SHIPPING = 2,
  TRANSFERRING = 3,
  COMPLETED = 4,
  CANCELLED = 5,
  REFUNDED = 6,
}

/**
 * Payment status.
 */
export enum PaymentStatus {
  PENDING = 0,
  SUCCESS = 1,
  FAILED = 2,
  REFUNDED = 3,
}
```

- [ ] **Step 3.7: 创建 src/index.ts**

写入 `/home/scutech/桌面/study/cars-system/packages/shared-types/src/index.ts`：

```typescript
export * from './api.js'
export * from './errors.js'
export * from './enums.js'
```

注意：TypeScript ESM + Bundler 模式可省略 `.js`，但显式写出来更安全（也兼容 NodeNext）。如使用 `moduleResolution: "Bundler"` 可省略，本计划保留 `.js` 以兼容未来切换。

实际 P0 阶段可改成省略 `.js`：

```typescript
export * from './api'
export * from './errors'
export * from './enums'
```

- [ ] **Step 3.8: 编写 enums 测试**

写入 `/home/scutech/桌面/study/cars-system/packages/shared-types/tests/enums.test.ts`：

```typescript
import { describe, expect, it } from 'vitest'
import { UserRole, hasRole, OrderStatus, CarStatus } from '../src/enums.js'
import { ErrorCode } from '../src/errors.js'

describe('UserRole.hasRole', () => {
  it('returns true when user has required role', () => {
    // user is seller+admin (2|4 = 6)
    expect(hasRole(6, UserRole.SELLER)).toBe(true)
    expect(hasRole(6, UserRole.ADMIN)).toBe(true)
  })

  it('returns false when user lacks required role', () => {
    expect(hasRole(UserRole.BUYER, UserRole.SELLER)).toBe(false)
    expect(hasRole(UserRole.SELLER, UserRole.ADMIN)).toBe(false)
  })

  it('returns false for guest', () => {
    expect(hasRole(UserRole.GUEST, UserRole.BUYER)).toBe(false)
  })
})

describe('OrderStatus transitions', () => {
  it('CREATED is initial state', () => {
    expect(OrderStatus.CREATED).toBe(0)
  })

  it('happy path ends at COMPLETED', () => {
    const happyPath = [
      OrderStatus.CREATED,
      OrderStatus.PAID,
      OrderStatus.SHIPPING,
      OrderStatus.TRANSFERRING,
      OrderStatus.COMPLETED,
    ]
    expect(happyPath.length).toBe(5)
  })
})

describe('CarStatus', () => {
  it('DRAFT is 0', () => {
    expect(CarStatus.DRAFT).toBe(0)
  })
})

describe('ErrorCode', () => {
  it('OK is 0', () => {
    expect(ErrorCode.OK).toBe(0)
  })
})
```

- [ ] **Step 3.9: 安装依赖并运行测试**

```bash
cd /home/scutech/桌面/study/cars-system
# 安装 pnpm（如未安装）
npm install -g pnpm@9

pnpm install
pnpm --filter @cars-system/shared-types test
```

预期：所有测试 PASS，输出形如：

```
✓ tests/enums.test.ts (5 tests) passed
```

- [ ] **Step 3.10: commit**

```bash
cd /home/scutech/桌面/study/cars-system
git add packages/shared-types pnpm-lock.yaml
git commit -m "feat(shared-types): add API types, error codes, and enums with tests"
```

---

## Task 4: 后端 — 配置加载（TDD）

**Files:**
- Create: `backend/internal/config/config.go`
- Create: `backend/internal/config/config_test.go`

**Interfaces:**
- Produces: `config.Load() (*Config, error)`，读取环境变量返回强类型 Config

---

- [ ] **Step 4.1: 添加依赖**

```bash
cd /home/scutech/桌面/study/cars-system/backend
go get github.com/spf13/viper@latest
go get github.com/stretchr/testify@latest
```

- [ ] **Step 4.2: 编写失败测试**

写入 `/home/scutech/桌面/study/cars-system/backend/internal/config/config_test.go`：

```go
package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear env vars that might interfere
	clearEnv(t)

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "development", cfg.AppEnv)
	assert.Equal(t, 8080, cfg.AppPort)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, 5432, cfg.DBPort)
	assert.Equal(t, "memory", cfg.WSBroker)
	assert.Equal(t, "local", cfg.PaymentGateway)
	assert.Equal(t, 15*time.Minute, cfg.JWTAccessTTL)
	assert.Equal(t, 7*24*time.Hour, cfg.JWTRefreshTTL)
}

func TestLoad_FromEnv(t *testing.T) {
	os.Setenv("APP_ENV", "production")
	os.Setenv("APP_PORT", "9090")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("JWT_SECRET", "supersecret")
	defer clearEnv(t)

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "production", cfg.AppEnv)
	assert.Equal(t, 9090, cfg.AppPort)
	assert.Equal(t, "db.example.com", cfg.DBHost)
	assert.Equal(t, "supersecret", cfg.JWTSecret)
}

func clearEnv(t *testing.T) {
	t.Helper()
	vars := []string{
		"APP_ENV", "APP_PORT", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD",
		"DB_NAME", "DB_SSLMODE", "REDIS_HOST", "REDIS_PORT", "REDIS_PASSWORD",
		"REDIS_DB", "JWT_SECRET", "JWT_ACCESS_TTL", "JWT_REFRESH_TTL",
		"WS_BROKER", "PAYMENT_GATEWAY",
	}
	for _, v := range vars {
		os.Unsetenv(v)
	}
}
```

- [ ] **Step 4.3: 运行测试确认失败**

```bash
cd /home/scutech/桌面/study/cars-system/backend
go test ./internal/config/... -v
```

预期：编译错误（`Load` 未定义）或测试 FAIL。

- [ ] **Step 4.4: 实现 Config**

写入 `/home/scutech/桌面/study/cars-system/backend/internal/config/config.go`：

```go
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv   string
	AppPort  int

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int

	JWTSecret      string
	JWTAccessTTL   time.Duration
	JWTRefreshTTL  time.Duration

	WSBroker      string
	PaymentGateway string
}

// Load reads configuration from environment variables (with viper).
func Load() (*Config, error) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Defaults
	setDefaults(v)

	cfg := &Config{
		AppEnv:         v.GetString("APP_ENV"),
		AppPort:        v.GetInt("APP_PORT"),
		DBHost:         v.GetString("DB_HOST"),
		DBPort:         v.GetInt("DB_PORT"),
		DBUser:         v.GetString("DB_USER"),
		DBPassword:     v.GetString("DB_PASSWORD"),
		DBName:         v.GetString("DB_NAME"),
		DBSSLMode:      v.GetString("DB_SSLMODE"),
		RedisHost:      v.GetString("REDIS_HOST"),
		RedisPort:      v.GetInt("REDIS_PORT"),
		RedisPassword:  v.GetString("REDIS_PASSWORD"),
		RedisDB:        v.GetInt("REDIS_DB"),
		JWTSecret:      v.GetString("JWT_SECRET"),
		JWTAccessTTL:   v.GetDuration("JWT_ACCESS_TTL"),
		JWTRefreshTTL:  v.GetDuration("JWT_REFRESH_TTL"),
		WSBroker:       v.GetString("WS_BROKER"),
		PaymentGateway: v.GetString("PAYMENT_GATEWAY"),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config invalid: %w", err)
	}
	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_PORT", 8080)

	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_USER", "cars")
	v.SetDefault("DB_PASSWORD", "cars_pass")
	v.SetDefault("DB_NAME", "cars_db")
	v.SetDefault("DB_SSLMODE", "disable")

	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", 6379)
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	v.SetDefault("JWT_SECRET", "dev-secret-change-me-in-production")
	v.SetDefault("JWT_ACCESS_TTL", "15m")
	v.SetDefault("JWT_REFRESH_TTL", "168h")

	v.SetDefault("WS_BROKER", "memory")
	v.SetDefault("PAYMENT_GATEWAY", "local")
}

func (c *Config) validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}
```

- [ ] **Step 4.5: 运行测试确认通过**

```bash
cd /home/scutech/桌面/study/cars-system/backend
go test ./internal/config/... -v
```

预期：3 个测试全部 PASS。

- [ ] **Step 4.6: commit**

```bash
cd /home/scutech/桌面/study/cars-system/backend
git add go.mod go.sum internal/config/
cd ..
git commit -m "feat(backend): add config loading with viper and tests"
```

---

## Task 5: 后端 — 错误码 + 统一响应（TDD）

**Files:**
- Create: `backend/internal/shared/errors/errors.go`
- Create: `backend/internal/shared/response/response.go`
- Create: `backend/internal/shared/response/response_test.go`

**Interfaces:**
- Produces: `response.OK(c, data)`、`response.Fail(c, code, message)`、`errors.New(code, msg)` 返回 Gin handler 友好的 error

---

- [ ] **Step 5.1: 创建 errors.go**

写入 `/home/scutech/桌面/study/cars-system/backend/internal/shared/errors/errors.go`：

```go
package errors

import "fmt"

// Error codes. Keep in sync with packages/shared-types/src/errors.ts.
const (
	CodeOK             = 0
	CodeParamInvalid   = 10001
	CodeSystemError    = 10002
	CodeUnauthorized   = 20001
	CodeTokenExpired   = 20002
	CodeForbidden      = 20003
	CodeUserNotFound   = 30001
	CodeCarNotFound    = 31001
	CodeCarAlreadySold = 31002
	CodeOrderState     = 32001
	CodeOrderTimeout   = 32002
	CodePaymentFailed  = 33001
	CodeRefundFailed   = 33002
)

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code=%d message=%s", e.Code, e.Message)
}

func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}
```

- [ ] **Step 5.2: 编写 response 测试**

写入 `/home/scutech/桌面/study/cars-system/backend/internal/shared/response/response_test.go`：

```go
package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "github.com/scutech/cars-system/backend/internal/shared/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/ok", func(c *gin.Context) { OK(c, gin.H{"foo": "bar"}) })
	r.GET("/fail", func(c *gin.Context) { Fail(c, apperrors.CodeCarNotFound, "car not found") })
	return r
}

func TestOK(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ok", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.EqualValues(t, 0, body["code"])
	assert.Equal(t, "ok", body["message"])
	assert.NotNil(t, body["data"])
}

func TestFail(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/fail", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code) // 业务错误仍用 200，code 字段标识

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.EqualValues(t, 31001, body["code"])
	assert.Equal(t, "car not found", body["message"])
	assert.Nil(t, body["data"])
}
```

- [ ] **Step 5.3: 实现 response.go**

写入 `/home/scutech/桌面/study/cars-system/backend/internal/shared/response/response.go`：

```go
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/scutech/cars-system/backend/internal/shared/errors"
)

type Envelope struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Envelope{Code: apperrors.CodeOK, Message: "ok", Data: data})
}

func Fail(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(http.StatusOK, Envelope{Code: code, Message: message})
}
```

- [ ] **Step 5.4: 添加 gin 依赖**

```bash
cd /home/scutech/桌面/study/cars-system/backend
go get github.com/gin-gonic/gin@latest
```

- [ ] **Step 5.5: 运行测试**

```bash
cd /home/scutech/桌面/study/cars-system/backend
go test ./internal/shared/... -v
```

预期：2 个测试 PASS。

- [ ] **Step 5.6: commit**

```bash
cd /home/scutech/桌面/study/cars-system/backend
git add go.mod go.sum internal/shared/
cd ..
git commit -m "feat(backend): add error codes and unified response envelope with tests"
```

---

## Task 6: 后端 — 中间件（logger / recovery / cors）

**Files:**
- Create: `backend/internal/shared/middleware/logger.go`
- Create: `backend/internal/shared/middleware/recovery.go`
- Create: `backend/internal/shared/middleware/cors.go`

**Interfaces:**
- Produces: `middleware.Logger()`, `middleware.Recovery()`, `middleware.CORS()`

---

- [ ] **Step 6.1: 添加 zap 依赖**

```bash
cd /home/scutech/桌面/study/cars-system/backend
go get go.uber.org/zap@latest
go get github.com/gin-contrib/cors@latest
```

- [ ] **Step 6.2: 创建 logger.go**

写入 `/home/scutech/桌面/study/cars-system/backend/internal/shared/middleware/logger.go`：

```go
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger returns a Gin middleware that logs each request with zap.
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		logger.Info("http request",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Int("body_size", c.Writer.Size()),
		)
	}
}
```

- [ ] **Step 6.3: 创建 recovery.go**

写入 `/home/scutech/桌面/study/cars-system/backend/internal/shared/middleware/recovery.go`：

```go
package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	apperrors "github.com/scutech/cars-system/backend/internal/shared/errors"
	"github.com/scutech/cars-system/backend/internal/shared/response"
)

// Recovery returns a Gin middleware that recovers from panics and returns 500.
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.ByteString("stack", debug.Stack()),
				)
				response.Fail(c, apperrors.CodeSystemError, "internal server error")
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
```

- [ ] **Step 6.4: 创建 cors.go**

写入 `/home/scutech/桌面/study/cars-system/backend/internal/shared/middleware/cors.go`：

```go
package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS returns a permissive CORS middleware suitable for development.
// In production, restrict AllowOrigins to known domains.
func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
}
```

- [ ] **Step 6.5: 验证编译**

```bash
cd /home/scutech/桌面/study/cars-system/backend
go build ./...
```

预期：无错误。

- [ ] **Step 6.6: commit**

```bash
cd /home/scutech/桌面/study/cars-system/backend
git add go.mod go.sum internal/shared/middleware/
cd ..
git commit -m "feat(backend): add logger, recovery, and cors middleware"
```

---

## Task 7: 后端 — 数据库连接（TDD）

**Files:**
- Create: `backend/internal/database/database.go`
- Create: `backend/internal/database/database_test.go`

**Interfaces:**
- Produces: `database.Connect(cfg) (*gorm.DB, error)`，连接 PostgreSQL 并返回 GORM 实例

---

- [ ] **Step 7.1: 添加 GORM 依赖**

```bash
cd /home/scutech/桌面/study/cars-system/backend
go get gorm.io/gorm@latest
go get gorm.io/driver/postgres@latest
```

- [ ] **Step 7.2: 编写测试**

写入 `/home/scutech/桌面/study/cars-system/backend/internal/database/database_test.go`：

```go
package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/scutech/cars-system/backend/internal/config"
)

// TestConnect_RequiresConfig is a smoke test that Connect returns an error
// when given a bad config (unreachable host). It does not require a real DB.
func TestConnect_InvalidHost(t *testing.T) {
	os.Setenv("DB_HOST", "nonexistent.invalid.host")
	os.Setenv("DB_PORT", "1")
	os.Setenv("APP_ENV", "test")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PORT")
	defer os.Unsetenv("APP_ENV")

	cfg, err := config.Load()
	require.NoError(t, err)

	db, err := Connect(cfg)
	assert.Error(t, err, "expected connection to fail with invalid host")
	assert.Nil(t, db)
}
```

- [ ] **Step 7.3: 运行测试确认失败**

```bash
cd /home/scutech/桌面/study/cars-system/backend
go test ./internal/database/... -v
```

预期：编译错误（`Connect` 未定义）。

- [ ] **Step 7.4: 实现 database.go**

写入 `/home/scutech/桌面/study/cars-system/backend/internal/database/database.go`：

```go
package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/scutech/cars-system/backend/internal/config"
)

// Connect opens a GORM connection to PostgreSQL and verifies it with Ping.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return db, nil
}
```

- [ ] **Step 7.5: 运行测试**

```bash
cd /home/scutech/桌面/study/cars-system/backend
go test ./internal/database/... -v
```

预期：测试 PASS（因为故意连不上，验证错误返回）。

- [ ] **Step 7.6: commit**

```bash
cd /home/scutech/桌面/study/cars-system/backend
git add go.mod go.sum internal/database/
cd ..
git commit -m "feat(backend): add gorm postgres connection with tests"
```

---

## Task 8: 后端 — Health Check Handler（TDD）

**Files:**
- Create: `backend/internal/shared/handler/health.go`
- Create: `backend/internal/shared/handler/health_test.go`

**Interfaces:**
- Produces: `handler.Health(db)` 返回 `gin.HandlerFunc`，检查数据库 ping 后返回 status

---

- [ ] **Step 8.1: 编写测试**

写入 `/home/scutech/桌面/study/cars-system/backend/internal/shared/handler/health_test.go`：

```go
package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestHealth_NoDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/health", Health(nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.EqualValues(t, 0, body["code"])
	assert.Equal(t, "ok", body["message"])

	data := body["data"].(map[string]interface{})
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "cars-system", data["service"])
	assert.Equal(t, "unavailable", data["db"])
}

func TestHealth_WithDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Pass a non-nil *gorm.DB to exercise the DB branch.
	// We can't ping without a real DB, so we expect db status to remain unavailable
	// even with non-nil db, but the handler must not panic.
	r := gin.New()
	r.GET("/api/v1/health", Health(&gorm.DB{}))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
```

- [ ] **Step 8.2: 运行测试确认失败**

```bash
cd /home/scutech/桌面/study/cars-system/backend
go test ./internal/shared/handler/... -v
```

预期：编译错误（`Health` 未定义）。

- [ ] **Step 8.3: 实现 health.go**

写入 `/home/scutech/桌面/study/cars-system/backend/internal/shared/handler/health.go`：

```go
package handler

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/scutech/cars-system/backend/internal/shared/response"
)

// Health returns a Gin handler that reports service health.
// If db is provided, it pings the database with a short timeout.
func Health(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		dbStatus := "unavailable"
		if db != nil {
			sqlDB, err := db.DB()
			if err == nil {
				ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
				defer cancel()
				if pingErr := sqlDB.PingContext(ctx); pingErr == nil {
					dbStatus = "ok"
				}
			}
		}

		data := gin.H{
			"status":  "ok",
			"service": "cars-system",
			"db":      dbStatus,
		}
		response.OK(c, data)
	}
}
```

- [ ] **Step 8.4: 运行测试**

```bash
cd /home/scutech/桌面/study/cars-system/backend
go test ./internal/shared/handler/... -v
```

预期：2 个测试 PASS。

- [ ] **Step 8.5: commit**

```bash
cd /home/scutech/桌面/study/cars-system/backend
git add internal/shared/handler/
cd ..
git commit -m "feat(backend): add health check handler with db ping"
```

---

## Task 9: 后端 — 整合 cmd/api/main.go

**Files:**
- Modify: `backend/cmd/api/main.go`（替换占位）

**Interfaces:**
- Produces: `main()` 启动 Gin，监听 `:8080`，挂载 `/api/v1/health`

---

- [ ] **Step 9.1: 编写 main.go**

写入 `/home/scutech/桌面/study/cars-system/backend/cmd/api/main.go`：

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/scutech/cars-system/backend/internal/config"
	"github.com/scutech/cars-system/backend/internal/database"
	"github.com/scutech/cars-system/backend/internal/shared/handler"
	"github.com/scutech/cars-system/backend/internal/shared/middleware"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("logger init: %v", err)
	}
	defer logger.Sync() //nolint:errcheck

	// Connect to database (non-fatal: log warning if unavailable, health endpoint will report it)
	db, dbErr := database.Connect(cfg)
	if dbErr != nil {
		logger.Warn("database unavailable on startup", zap.Error(dbErr))
	} else {
		logger.Info("database connected")
	}

	// Setup Gin
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middleware.Logger(logger), middleware.Recovery(logger), middleware.CORS())

	// Health check
	r.GET("/api/v1/health", handler.Health(db))

	// HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.AppPort),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		logger.Info("api server starting", zap.Int("port", cfg.AppPort))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("forced shutdown", zap.Error(err))
	}
}
```

- [ ] **Step 9.2: 验证编译**

```bash
cd /home/scutech/桌面/study/cars-system/backend
go build ./cmd/api
```

预期：生成 `bin/api` 可执行文件，无错误。

- [ ] **Step 9.3: 本地启动并 curl 测试**

```bash
cd /home/scutech/桌面/study/cars-system/backend
./bin/api &
sleep 2
curl -s http://localhost:8080/api/v1/health | jq
kill %1
```

预期：返回 `{"code":0,"message":"ok","data":{"status":"ok","service":"cars-system","db":"ok"}}`

（如未配置 jq，可去掉 `| jq` 直接看 raw JSON）

- [ ] **Step 9.4: commit**

```bash
cd /home/scutech/桌面/study/cars-system/backend
git add cmd/api/main.go
cd ..
git commit -m "feat(backend): wire gin server with config, db, and health endpoint"
```

---

## Task 10: 前端 — Next.js 用户端骨架

**Files:**
- Create: `apps/web/package.json`
- Create: `apps/web/tsconfig.json`
- Create: `apps/web/next.config.ts`
- Create: `apps/web/vitest.config.ts`
- Create: `apps/web/app/layout.tsx`
- Create: `apps/web/app/page.tsx`
- Create: `apps/web/app/globals.css`
- Create: `apps/web/lib/api.ts`
- Create: `apps/web/lib/api.test.ts`

**Interfaces:**
- Produces: `pnpm --filter web dev` 起 :3000 服务
- Produces: `apps/web` 引用 `@cars-system/shared-types`

---

- [ ] **Step 10.1: 创建 package.json**

写入 `/home/scutech/桌面/study/cars-system/apps/web/package.json`：

```json
{
  "name": "web",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev -p 3000",
    "build": "next build",
    "start": "next start -p 3000",
    "lint": "next lint",
    "typecheck": "tsc --noEmit",
    "test": "vitest run"
  },
  "dependencies": {
    "@cars-system/shared-types": "workspace:*",
    "next": "^15.0.0",
    "react": "^18.3.0",
    "react-dom": "^18.3.0",
    "antd": "^5.20.0",
    "@ant-design/icons": "^5.4.0",
    "@tanstack/react-query": "^5.51.0",
    "zustand": "^5.0.0"
  },
  "devDependencies": {
    "@types/node": "^20.14.0",
    "@types/react": "^18.3.0",
    "@types/react-dom": "^18.3.0",
    "@vitejs/plugin-react": "^4.3.0",
    "eslint": "^8.57.0",
    "eslint-config-next": "^15.0.0",
    "typescript": "^5.4.0",
    "vitest": "^1.6.0",
    "happy-dom": "^15.0.0"
  }
}
```

- [ ] **Step 10.2: 创建 tsconfig.json**

写入 `/home/scutech/桌面/study/cars-system/apps/web/tsconfig.json`：

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [{ "name": "next" }],
    "paths": {
      "@/*": ["./*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}
```

- [ ] **Step 10.3: 创建 next.config.ts**

写入 `/home/scutech/桌面/study/cars-system/apps/web/next.config.ts`：

```typescript
import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
  reactStrictMode: true,
  transpilePackages: ['@cars-system/shared-types'],
  experimental: {
    typedRoutes: false,
  },
}

export default nextConfig
```

- [ ] **Step 10.4: 创建 vitest.config.ts**

写入 `/home/scutech/桌面/study/cars-system/apps/web/vitest.config.ts`：

```typescript
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'happy-dom',
    globals: false,
  },
})
```

- [ ] **Step 10.5: 创建 app/layout.tsx**

写入 `/home/scutech/桌面/study/cars-system/apps/web/app/layout.tsx`：

```tsx
import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'Cars System — User',
  description: '二手车交易平台',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="zh-CN">
      <body>{children}</body>
    </html>
  )
}
```

- [ ] **Step 10.6: 创建 app/page.tsx**

写入 `/home/scutech/桌面/study/cars-system/apps/web/app/page.tsx`：

```tsx
'use client'

import { ConfigProvider, Typography } from 'antd'

const { Title, Paragraph } = Typography

export default function HomePage() {
  return (
    <ConfigProvider>
      <main style={{ padding: 48 }}>
        <Title>Cars System</Title>
        <Paragraph>二手车交易平台 — 用户端骨架</Paragraph>
        <Paragraph type="secondary">P0 阶段：基础设施就绪，业务功能将在后续里程碑实现。</Paragraph>
      </main>
    </ConfigProvider>
  )
}
```

- [ ] **Step 10.7: 创建 app/globals.css**

写入 `/home/scutech/桌面/study/cars-system/apps/web/app/globals.css`：

```css
html, body {
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}
```

- [ ] **Step 10.8: 创建 lib/api.ts**

写入 `/home/scutech/桌面/study/cars-system/apps/web/lib/api.ts`：

```typescript
import type { APIResponse } from '@cars-system/shared-types'

const BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? 'http://localhost:8080/api/v1'

export class APIError extends Error {
  constructor(public code: number, message: string) {
    super(message)
    this.name = 'APIError'
  }
}

export async function apiGet<T>(path: string): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    method: 'GET',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
  })
  return handleResponse<T>(res)
}

async function handleResponse<T>(res: Response): Promise<T> {
  const json = (await res.json()) as APIResponse<T>
  if (json.code !== 0) {
    throw new APIError(json.code, json.message)
  }
  return json.data as T
}
```

- [ ] **Step 10.9: 创建 lib/api.test.ts**

写入 `/home/scutech/桌面/study/cars-system/apps/web/lib/api.test.ts`：

```typescript
import { describe, expect, it, vi, afterEach } from 'vitest'
import { apiGet, APIError } from './api'

describe('apiGet', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('returns data when response code is 0', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      json: async () => ({ code: 0, message: 'ok', data: { foo: 'bar' } }),
    }))

    const data = await apiGet<{ foo: string }>('/health')
    expect(data).toEqual({ foo: 'bar' })
  })

  it('throws APIError when response code is non-zero', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      json: async () => ({ code: 31001, message: 'car not found' }),
    }))

    await expect(apiGet('/cars/1')).rejects.toThrow(APIError)
    await expect(apiGet('/cars/1')).rejects.toMatchObject({
      code: 31001,
      message: 'car not found',
    })
  })
})
```

- [ ] **Step 10.10: 安装依赖**

```bash
cd /home/scutech/桌面/study/cars-system
pnpm install
```

- [ ] **Step 10.11: 运行测试**

```bash
pnpm --filter web test
```

预期：2 个测试 PASS。

- [ ] **Step 10.12: 本地启动验证**

```bash
pnpm --filter web dev
# 另一终端：
curl -s http://localhost:3000 | grep -o 'Cars System' | head -1
```

预期：输出 `Cars System`。

- [ ] **Step 10.13: commit**

```bash
cd /home/scutech/桌面/study/cars-system
git add apps/web pnpm-lock.yaml
git commit -m "feat(web): scaffold next.js 15 user app with ant design and api client"
```

---

## Task 11: 前端 — Next.js 管理端骨架

**Files:** 与 Task 10 同结构，但路径在 `apps/admin/`，端口 3001

---

- [ ] **Step 11.1: 创建 package.json**

写入 `/home/scutech/桌面/study/cars-system/apps/admin/package.json`：

```json
{
  "name": "admin",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev -p 3001",
    "build": "next build",
    "start": "next start -p 3001",
    "lint": "next lint",
    "typecheck": "tsc --noEmit",
    "test": "vitest run"
  },
  "dependencies": {
    "@cars-system/shared-types": "workspace:*",
    "next": "^15.0.0",
    "react": "^18.3.0",
    "react-dom": "^18.3.0",
    "antd": "^5.20.0",
    "@ant-design/icons": "^5.4.0",
    "@tanstack/react-query": "^5.51.0",
    "zustand": "^5.0.0"
  },
  "devDependencies": {
    "@types/node": "^20.14.0",
    "@types/react": "^18.3.0",
    "@types/react-dom": "^18.3.0",
    "@vitejs/plugin-react": "^4.3.0",
    "eslint": "^8.57.0",
    "eslint-config-next": "^15.0.0",
    "typescript": "^5.4.0",
    "vitest": "^1.6.0",
    "happy-dom": "^15.0.0"
  }
}
```

- [ ] **Step 11.2: 创建 tsconfig.json**

写入 `/home/scutech/桌面/study/cars-system/apps/admin/tsconfig.json`：

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [{ "name": "next" }],
    "paths": {
      "@/*": ["./*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}
```

- [ ] **Step 11.3: 创建 next.config.ts**

写入 `/home/scutech/桌面/study/cars-system/apps/admin/next.config.ts`：

```typescript
import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
  reactStrictMode: true,
  transpilePackages: ['@cars-system/shared-types'],
  experimental: {
    typedRoutes: false,
  },
}

export default nextConfig
```

- [ ] **Step 11.4: 创建 vitest.config.ts**

写入 `/home/scutech/桌面/study/cars-system/apps/admin/vitest.config.ts`：

```typescript
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'happy-dom',
    globals: false,
  },
})
```

- [ ] **Step 11.5: 创建 app/layout.tsx**

写入 `/home/scutech/桌面/study/cars-system/apps/admin/app/layout.tsx`：

```tsx
import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'Cars System — Admin',
  description: '二手车交易平台管理后台',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="zh-CN">
      <body>{children}</body>
    </html>
  )
}
```

- [ ] **Step 11.6: 创建 app/globals.css**

写入 `/home/scutech/桌面/study/cars-system/apps/admin/app/globals.css`：

```css
html, body {
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}
```

- [ ] **Step 11.7: 创建 next-env.d.ts 占位**

```bash
echo '/// <reference types="next" />
/// <reference types="next/image-types/global" />' > /home/scutech/桌面/study/cars-system/apps/admin/next-env.d.ts
```

预期：生成 `apps/admin/next-env.d.ts` 文件。

具体每个文件内容相同（路径前缀改了）。

- [ ] **Step 11.8: 创建 app/page.tsx（admin 定制）**

写入 `/home/scutech/桌面/study/cars-system/apps/admin/app/page.tsx`：

```tsx
'use client'

import { ConfigProvider, Typography } from 'antd'

const { Title, Paragraph } = Typography

export default function AdminHomePage() {
  return (
    <ConfigProvider>
      <main style={{ padding: 48 }}>
        <Title level={2}>Cars System Admin</Title>
        <Paragraph>管理后台骨架</Paragraph>
        <Paragraph type="secondary">P0 阶段：基础设施就绪，审核/看板等功能将在后续里程碑实现。</Paragraph>
      </main>
    </ConfigProvider>
  )
}
```

- [ ] **Step 11.9: 创建 lib/api.ts**

写入 `/home/scutech/桌面/study/cars-system/apps/admin/lib/api.ts`：

```typescript
import type { APIResponse } from '@cars-system/shared-types'

const BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? 'http://localhost:8080/api/v1'

export class APIError extends Error {
  constructor(public code: number, message: string) {
    super(message)
    this.name = 'APIError'
  }
}

export async function apiGet<T>(path: string): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    method: 'GET',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
  })
  return handleResponse<T>(res)
}

async function handleResponse<T>(res: Response): Promise<T> {
  const json = (await res.json()) as APIResponse<T>
  if (json.code !== 0) {
    throw new APIError(json.code, json.message)
  }
  return json.data as T
}
```

- [ ] **Step 11.10: 创建 lib/api.test.ts**

写入 `/home/scutech/桌面/study/cars-system/apps/admin/lib/api.test.ts`：

```typescript
import { describe, expect, it, vi, afterEach } from 'vitest'
import { apiGet, APIError } from './api'

describe('apiGet', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('returns data when response code is 0', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      json: async () => ({ code: 0, message: 'ok', data: { foo: 'bar' } }),
    }))

    const data = await apiGet<{ foo: string }>('/health')
    expect(data).toEqual({ foo: 'bar' })
  })

  it('throws APIError when response code is non-zero', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      json: async () => ({ code: 31001, message: 'car not found' }),
    }))

    await expect(apiGet('/cars/1')).rejects.toThrow(APIError)
    await expect(apiGet('/cars/1')).rejects.toMatchObject({
      code: 31001,
      message: 'car not found',
    })
  })
})
```

- [ ] **Step 11.11: 安装 + 测试 + 启动**

```bash
cd /home/scutech/桌面/study/cars-system
pnpm install
pnpm --filter admin test
pnpm --filter admin dev
# 另一终端：
curl -s http://localhost:3001 | grep -o 'Admin' | head -1
```

预期：测试 PASS，curl 输出包含 `Admin`。

- [ ] **Step 11.12: commit**

```bash
cd /home/scutech/桌面/study/cars-system
git add apps/admin
git commit -m "feat(admin): scaffold next.js 15 admin app with ant design"
```

---

## Task 12: Docker Compose — 集成所有服务

**Files:**
- Modify: `docker-compose.yml`（添加 backend / worker / web / admin service）
- Modify: `nginx/nginx.conf`（完整路由）
- Modify: `apps/web/Dockerfile`（替换占位为真实 Next.js 镜像）
- Modify: `apps/admin/Dockerfile`（同上）

---

- [ ] **Step 12.1: 完善 web Dockerfile**

写入 `/home/scutech/桌面/study/cars-system/apps/web/Dockerfile`：

```dockerfile
# syntax=docker/dockerfile:1.6

# ---- Deps stage ----
FROM node:20-alpine AS deps
WORKDIR /app
RUN corepack enable
COPY package.json pnpm-lock.yaml* ./
COPY ../packages/shared-types/package.json ../packages/shared-types/
RUN pnpm install --frozen-lockfile

# ---- Build stage ----
FROM node:20-alpine AS build
WORKDIR /app
RUN corepack enable
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN pnpm --filter web build

# ---- Runtime stage ----
FROM node:20-alpine AS runtime
WORKDIR /app
RUN corepack enable
ENV NODE_ENV=production
COPY --from=build /app/apps/web/.next ./apps/web/.next
COPY --from=build /app/apps/web/public ./apps/web/public
COPY --from=build /app/apps/web/package.json ./apps/web/
COPY --from=build /app/apps/web/next.config.ts ./apps/web/
COPY --from=build /app/node_modules ./node_modules
WORKDIR /app/apps/web
EXPOSE 3000
CMD ["pnpm", "start"]
```

> 注意：上面的 `COPY ../packages/...` 在 build context 是 monorepo 根时才能用。需要在 `docker-compose.yml` 中把 build context 设为 `.`（仓库根）。

- [ ] **Step 12.2: 完善 admin Dockerfile**

同 Task 12.1，把 `web` 替换为 `admin`，端口 3001。

写入 `/home/scutech/桌面/study/cars-system/apps/admin/Dockerfile`：

```dockerfile
# syntax=docker/dockerfile:1.6

FROM node:20-alpine AS deps
WORKDIR /app
RUN corepack enable
COPY package.json pnpm-lock.yaml* ./
COPY ../packages/shared-types/package.json ../packages/shared-types/
RUN pnpm install --frozen-lockfile

FROM node:20-alpine AS build
WORKDIR /app
RUN corepack enable
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN pnpm --filter admin build

FROM node:20-alpine AS runtime
WORKDIR /app
RUN corepack enable
ENV NODE_ENV=production
COPY --from=build /app/apps/admin/.next ./apps/admin/.next
COPY --from=build /app/apps/admin/public ./apps/admin/public
COPY --from=build /app/apps/admin/package.json ./apps/admin/
COPY --from=build /app/apps/admin/next.config.ts ./apps/admin/
COPY --from=build /app/node_modules ./node_modules
WORKDIR /app/apps/admin
EXPOSE 3001
CMD ["pnpm", "start"]
```

- [ ] **Step 12.3: 完善 nginx/nginx.conf**

写入 `/home/scutech/桌面/study/cars-system/nginx/nginx.conf`：

```nginx
events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    upstream backend_api {
        server backend:8080;
    }

    upstream web_app {
        server web:3000;
    }

    upstream admin_app {
        server admin:3001;
    }

    server {
        listen 80;
        server_name localhost;

        # User app
        location / {
            proxy_pass http://web_app;
            proxy_http_version 1.1;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # Admin app
        location /admin/ {
            proxy_pass http://admin_app/;
            proxy_http_version 1.1;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }

        # API + WebSocket
        location /api/ {
            proxy_pass http://backend_api;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_read_timeout 60s;
        }

        # Static uploads (future)
        location /static/ {
            alias /var/www/uploads/;
            expires 30d;
        }
    }
}
```

- [ ] **Step 12.4: 完善 docker-compose.yml**

写入 `/home/scutech/桌面/study/cars-system/docker-compose.yml`：

```yaml
services:
  postgres:
    image: ghcr.io/amutu/zhparser:16
    container_name: cars-postgres
    environment:
      POSTGRES_USER: ${DB_USER:-cars}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-cars_pass}
      POSTGRES_DB: ${DB_NAME:-cars_db}
    ports:
      - "5432:5432"
    volumes:
      - pg_data:/var/lib/postgresql/data
      - ./backend/migrations:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-cars}"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: cars-redis
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5

  backend:
    build:
      context: .
      dockerfile: backend/Dockerfile
    container_name: cars-backend
    environment:
      APP_ENV: ${APP_ENV:-development}
      APP_PORT: 8080
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: ${DB_USER:-cars}
      DB_PASSWORD: ${DB_PASSWORD:-cars_pass}
      DB_NAME: ${DB_NAME:-cars_db}
      DB_SSLMODE: disable
      REDIS_HOST: redis
      REDIS_PORT: 6379
      JWT_SECRET: ${JWT_SECRET:-dev-secret-change-me-in-production}
      WS_BROKER: ${WS_BROKER:-memory}
      PAYMENT_GATEWAY: ${PAYMENT_GATEWAY:-local}
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    restart: unless-stopped

  worker:
    build:
      context: .
      dockerfile: backend/Dockerfile
    container_name: cars-worker
    command: ["/app/worker"]
    environment:
      APP_ENV: ${APP_ENV:-development}
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: ${DB_USER:-cars}
      DB_PASSWORD: ${DB_PASSWORD:-cars_pass}
      DB_NAME: ${DB_NAME:-cars_db}
      DB_SSLMODE: disable
      REDIS_HOST: redis
      REDIS_PORT: 6379
      JWT_SECRET: ${JWT_SECRET:-dev-secret-change-me-in-production}
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    restart: unless-stopped

  web:
    build:
      context: .
      dockerfile: apps/web/Dockerfile
    container_name: cars-web
    environment:
      NEXT_PUBLIC_API_BASE_URL: http://localhost/api/v1
      NEXT_PUBLIC_WS_URL: ws://localhost/api/v1/ws
    depends_on:
      - backend
    restart: unless-stopped

  admin:
    build:
      context: .
      dockerfile: apps/admin/Dockerfile
    container_name: cars-admin
    environment:
      NEXT_PUBLIC_API_BASE_URL: http://localhost/api/v1
      NEXT_PUBLIC_WS_URL: ws://localhost/api/v1/ws
    depends_on:
      - backend
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    container_name: cars-nginx
    ports:
      - "80:80"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./data/uploads:/var/www/uploads:ro
    depends_on:
      - backend
      - web
      - admin
    restart: unless-stopped

volumes:
  pg_data:
```

- [ ] **Step 12.5: 删除 nginx profiles 标记**

确认 docker-compose.yml 没有 `profiles: ["full"]` 标记（已删除 Step 2.5 中的 profiles 写法）。

- [ ] **Step 12.6: 一键启动所有服务**

```bash
cd /home/scutech/桌面/study/cars-system
docker compose down -v   # 清掉 P0 早期启动的孤立容器
docker compose up -d --build
```

预期：6 个容器（postgres / redis / backend / worker / web / admin / nginx）全部 running。

- [ ] **Step 12.7: 验证端到端**

```bash
# API 健康检查
curl -s http://localhost/api/v1/health | python3 -m json.tool

# 用户端首页
curl -s http://localhost | grep -o 'Cars System' | head -1

# 管理端首页
curl -s http://localhost/admin/ | grep -o 'Admin' | head -1
```

预期：
- API 返回 `{"code": 0, "message": "ok", "data": {"status": "ok", "service": "cars-system", "db": "ok"}}`
- 用户端 grep 命中 `Cars System`
- 管理端 grep 命中 `Admin`

- [ ] **Step 12.8: 查看日志确认无错误**

```bash
docker compose logs --tail=50 backend
```

预期：日志含 `api server starting` 与 `database connected`，无 panic / fatal。

- [ ] **Step 12.9: commit**

```bash
cd /home/scutech/桌面/study/cars-system
git add docker-compose.yml nginx/nginx.conf apps/web/Dockerfile apps/admin/Dockerfile
git commit -m "feat(infra): full docker-compose with backend, worker, web, admin, nginx"
```

---

## Task 13: CI 基础 — GitHub Actions

**Files:**
- Create: `.github/workflows/ci.yml`

---

- [ ] **Step 13.1: 创建 CI 配置**

写入 `/home/scutech/桌面/study/cars-system/.github/workflows/ci.yml`：

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  backend:
    name: Backend (Go)
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: backend
    services:
      postgres:
        image: ghcr.io/amutu/zhparser:16
        env:
          POSTGRES_USER: cars
          POSTGRES_PASSWORD: cars_pass
          POSTGRES_DB: cars_test
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 5s
          --health-timeout 5s
          --health-retries 5
      redis:
        image: redis:7-alpine
        ports:
          - 6379:6379
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 5s
          --health-timeout 3s
          --health-retries 5
    env:
      DB_HOST: localhost
      DB_PORT: 5432
      DB_USER: cars
      DB_PASSWORD: cars_pass
      DB_NAME: cars_test
      DB_SSLMODE: disable
      REDIS_HOST: localhost
      REDIS_PORT: 6379
      JWT_SECRET: ci-test-secret
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache-dependency-path: backend/go.sum
      - name: Download deps
        run: go mod download
      - name: Build
        run: go build ./...
      - name: Test
        run: go test ./... -race -coverprofile=cover.out
      - name: Upload coverage
        uses: actions/upload-artifact@v4
        with:
          name: backend-coverage
          path: backend/cover.out
        if: always()

  frontend:
    name: Frontend (TS)
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: pnpm/action-setup@v4
        with:
          version: 9
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'pnpm'
      - name: Install
        run: pnpm install --frozen-lockfile
      - name: Lint
        run: pnpm -r lint
      - name: Typecheck
        run: pnpm -r typecheck
      - name: Test
        run: pnpm -r test
      - name: Build
        run: pnpm --filter web build && pnpm --filter admin build
```

- [ ] **Step 13.2: commit**

```bash
cd /home/scutech/桌面/study/cars-system
git add .github
git commit -m "ci: add github actions for backend and frontend"
```

> 注：CI 是否能在 GitHub 上实际跑通，取决于你后续是否把仓库 push 到 GitHub。本任务仅定义工作流文件，**不在本地执行**。

---

## Task 14: 数据库迁移框架 + 占位迁移

**Files:**
- Create: `backend/migrations/000001_init_schema.up.sql`
- Create: `backend/migrations/000001_init_schema.down.sql`

---

- [ ] **Step 14.1: 创建占位 up 迁移**

写入 `/home/scutech/桌面/study/cars-system/backend/migrations/000001_init_schema.up.sql`：

```sql
-- P0 placeholder migration.
-- P1 will add: users, user_profiles, user_addresses, cars, car_images, etc.
-- P2 will add: orders, payments, reviews, aftersales.
-- P3 will add: conversations, messages, notifications.

-- Verify zhparser extension is available (bundled with image).
CREATE EXTENSION IF NOT EXISTS zhparser;

-- Create text search config (idempotent).
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_ts_config WHERE cfgname = 'chinese_zh') THEN
        EXECUTE 'CREATE TEXT SEARCH CONFIGURATION chinese_zh (PARSER = zhparser)';
        EXECUTE 'ALTER TEXT SEARCH CONFIGURATION chinese_zh ADD MAPPING FOR n,v,a,i,e,l,t WITH simple';
    END IF;
END $$;
```

- [ ] **Step 14.2: 创建 down 迁移**

写入 `/home/scutech/桌面/study/cars-system/backend/migrations/000001_init_schema.down.sql`：

```sql
DROP TEXT SEARCH CONFIGURATION IF EXISTS chinese_zh;
-- Extension is not dropped to avoid affecting other databases.
```

- [ ] **Step 14.3: 删除占位 .gitkeep**

```bash
cd /home/scutech/桌面/study/cars-system
rm backend/migrations/.gitkeep
```

- [ ] **Step 14.4: 重启 postgres 让迁移生效**

```bash
cd /home/scutech/桌面/study/cars-system
docker compose down postgres
docker volume rm cars-system_pg_data 2>/dev/null || true
docker compose up -d postgres
sleep 5
docker exec -it cars-postgres psql -U cars -d cars_db -c "SELECT cfgname FROM pg_ts_config WHERE cfgname='chinese_zh';"
```

预期：返回一行 `chinese_zh`，表示分词配置已创建。

- [ ] **Step 14.5: commit**

```bash
cd /home/scutech/桌面/study/cars-system
git add backend/migrations
git commit -m "feat(db): init schema migration with zhparser configuration"
```

---

## Task 15: 收尾 — README 完善 + 最终验证

**Files:**
- Modify: `README.md`

---

- [ ] **Step 15.1: 完善 README.md**

写入 `/home/scutech/桌面/study/cars-system/README.md`：

```markdown
# Cars System — 二手车交易平台

A full-stack used car trading platform built with **React + TypeScript + Go**.

## 技术栈

- **后端**：Go 1.22 + Gin + GORM + PostgreSQL 16 + Redis 7
- **前端**：Next.js 15 + TypeScript 5 + Ant Design 5
- **基础设施**：Docker Compose + Nginx + GitHub Actions
- **详细设计**：见 [`开发说明书.md`](./开发说明书.md)

## 一键启动

```bash
# 克隆（首次）
git clone <repo-url>
cd cars-system
cp .env.example .env

# 启动所有服务
docker compose up -d --build

# 等待 ~30 秒后验证
curl http://localhost/api/v1/health
```

启动后可访问：

| 地址 | 用途 |
|---|---|
| http://localhost | 用户端 (Next.js web) |
| http://localhost/admin | 管理端 (Next.js admin) |
| http://localhost/api/v1/health | 后端健康检查 |
| http://localhost:5432 | PostgreSQL |
| http://localhost:6379 | Redis |

## 本地开发（不使用 Docker）

```bash
# 后端
cd backend
go run ./cmd/api

# 前端用户端
cd apps/web
pnpm install
pnpm dev

# 前端管理端
cd apps/admin
pnpm install
pnpm dev
```

## 项目结构

```
cars-system/
├── backend/                  Go 后端
├── apps/web/                 用户端 (Next.js :3000)
├── apps/admin/               管理端 (Next.js :3001)
├── packages/shared-types/    共享 TypeScript 类型
├── nginx/                    反向代理配置
├── 开发说明书.md              设计文档
└── docker-compose.yml
```

## 测试

```bash
# 后端
cd backend && go test ./... -race

# 前端
pnpm -r test

# 端到端（P4 阶段添加 Playwright）
```

## 当前进度

- [x] **P0 基础设施** ← 当前阶段
- [ ] P1 用户与车辆
- [ ] P2 订单与本地支付
- [ ] P3 通信与评价
- [ ] P4 增强与打磨

## License

MIT
```

- [ ] **Step 15.2: 最终一键启动验证**

```bash
cd /home/scutech/桌面/study/cars-system
docker compose down -v
docker compose up -d --build
sleep 30
curl -s http://localhost/api/v1/health
echo ""
curl -sI http://localhost | head -1
curl -sI http://localhost/admin/ | head -1
```

预期：
- API 返回 code=0, status=ok
- 用户端 200 OK
- 管理端 200 OK

- [ ] **Step 15.3: 最终 commit**

```bash
cd /home/scutech/桌面/study/cars-system
git add README.md
git commit -m "docs: comprehensive README for P0 milestone"
```

- [ ] **Step 15.4: 推送到 GitHub**

```bash
# 第一次推送（如已创建 GitHub 仓库）
git remote add origin <github-repo-url>
git push -u origin main

# 或在 GitHub 上创建空仓库后：
# 按 GitHub 提示执行
```

---

## Self-Review Checklist（执行前自查）

完成后请逐条核对：

- [ ] **目录结构**符合 §File Structure 描述
- [ ] 所有 commit message 遵循 conventional commits
- [ ] 所有 API 响应符合 `{ code, message, data }` 格式
- [ ] 所有金额字段用 `NUMERIC(12,2)`（P0 无金额字段，约定带入 P1）
- [ ] 所有时间字段用 `TIMESTAMPTZ`（同上）
- [ ] `go test ./...` 全部通过且无 race condition
- [ ] `pnpm -r test` 全部通过
- [ ] `pnpm -r typecheck` 全部通过
- [ ] `docker compose up -d --build` 一次性启动成功
- [ ] `curl http://localhost/api/v1/health` 返回 `code: 0`

## P0 完成的定义（Definition of Done）

- ✅ 上述 15 个任务全部完成
- ✅ 所有 commit 在 main 分支
- ✅ 一次性 `docker compose up -d --build` 起服务
- ✅ 健康检查、用户端、管理端都可访问
- ✅ 本地 `go test` 与 `pnpm test` 全绿
- ✅ CI 配置文件已存在（即使未实际跑通）

## 完成后进入 P1

P0 通过后，下一步：
1. 创建 P1 feature 分支：`git checkout -b feature/p1-user-and-cars`
2. 调用 writing-plans skill 生成 P1 实施计划（用户与车辆模块）
3. 按计划执行

---

## 附录：常见问题排查

| 症状 | 排查 |
|---|---|
| `docker compose up` 后 backend 一直 restart | `docker compose logs backend`，多半是连不上 Postgres —— 等待 postgres healthy |
| `pnpm install` 报 lockfile 不一致 | `pnpm install --no-frozen-lockfile` 更新后重 commit lockfile |
| `curl http://localhost/api/v1/health` 返回 502 | nginx 还没准备好，backend 也可能没起来；`docker compose ps` 看状态 |
| zhparser 镜像拉取失败 | 改用 `postgres:16` + 单独安装 zhparser 扩展（需修改镜像） |
| 端口冲突 | 修改 `.env` 中 `APP_PORT` 和 docker-compose 中 `ports` 映射 |
