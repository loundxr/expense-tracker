# Expense Tracker API

A production-grade, highly scalable RESTful API for managing personal and shared finances, built with **Go (Golang)**, **Chi**, **PostgreSQL**, **Redis**, and structured around **Clean Architecture / Feature-Driven Modular Design**.

---

## Architectural Overview

The application follows a **Modular Monolith** pattern incorporating **Clean Architecture** principles within isolated feature domains:

```
                  ┌─────────────────────────────────────────────────┐
                  │                 HTTP Transport                  │
                  │   (go-chi/chi v5, DTOs, Path/Query Parsers)     │
                  └────────────────────────┬────────────────────────┘
                                           │
                                           ▼
                  ┌─────────────────────────────────────────────────┐
                  │                  Service Layer                  │
                  │      (Business Rules, RBAC, Domain Audits)      │
                  └──────────────┬───────────────────┬──────────────┘
                                 │                   │
                                 ▼                   ▼
                  ┌──────────────────────┐   ┌──────────────────────┐
                  │   Domain Invariants  │   │   Repository Layer   │
                  │ (Entities, Nullables)│   │  (pgxpool, Adapters) │
                  └──────────────────────┘   └──────────┬───────────┘
                                                        │
                                                        ▼
                                             ┌──────────────────────┐
                                             │ PostgreSQL 16+ / DB  │
                                             └──────────────────────┘
```

- **Feature Modules (`internal/features/*`)**: Each business capability (`auth`, `users`, `accounts`, `categories`, `expenses`, `analytics`) is an encapsulated module exposing a clean constructor factory (`New(deps)`) and explicit route registration.
- **Composition Root (`cmd/expense-tracker/main.go`)**: Only orchestrates top-level dependency injection, database pooling, configuration loading, and graceful shutdown.
- **Consumer-Driven Interfaces**: Service and repository interfaces are defined locally where they are consumed, eliminating circular import cycles and adhering to the *Interface Segregation Principle (ISP)*.

---

## Key Features

- **Authentication & RBAC:** User registration, login, and JWT-based authentication with `user` and `admin` roles.
- **Cross-Feature Atomic Transactions:** User creation and default "Personal" wallet provisioning execute inside an atomic database transaction (`WithinTransaction`) — zero orphan users on failures.
- **Dynamic Token Invalidation (Redis):** Immediate token revocation using Redis-backed timestamp checks (`iat`). If an admin is demoted, their existing tokens are rejected with `401 Unauthorized` instantly.
- **Accounts & Shared Budgets:** Create wallets, share access with other users via email, and revoke permissions.
- **Categories:** System default categories (seeded via migrations) + custom user categories.
- **Expenses & Transactions:** Record spending with high-precision integer-cent math (`BIGINT` in DB) and atomic partial updates (`PATCH`) with defensive copying.
- **Optimistic Locking:** Entity updates protected against concurrent edit collisions via a `version` column (`409 Conflict`).
- **Standardized Response Envelopes:** All collection endpoints return structured JSON envelopes (e.g., `{"expenses": [...]}`) with constructor protection against `null` slices.
- **Real-Time Financial Analytics:**
  - **Summary Dashboard:** Total spending, average transaction, transaction count, largest purchase, top category.
  - **Category Breakdown:** Aggregated spending per category with SQL window-function percentage shares (`SUM(SUM) OVER()`).
  - **Spending Trends:** Chronological time series grouped by `day` or `month` via PostgreSQL `DATE_TRUNC`.
  - **Members Contribution:** Spending distribution among participants of a shared account.

---

## Technical Stack

| Component                   | Technology                                                  | Purpose                                                                              |
| :-------------------------- | :---------------------------------------------------------- | :----------------------------------------------------------------------------------- |
| **Language**                | [Go 1.26.2](https://go.dev/)                                | Core runtime                                                                         |
| **HTTP Router**             | [go-chi/chi v5](https://github.com/go-chi/chi)              | Fast, idiomatic, lightweight HTTP routing                                            |
| **Database**                | [PostgreSQL 18.6](https://www.postgresql.org/)              | Relational ACID persistence                                                          |
| **DB Driver & Pool**        | [pgx/v5 (pgxpool)](https://github.com/jackc/pgx)            | High-performance PostgreSQL driver and connection pool                               |
| **Cache Store**             | [Redis 8+](https://redis.io/)                               | In-memory cache & future token invalidation store                                    |
| **Migrations**              | [golang-migrate](https://github.com/golang-migrate/migrate) | Database schema version control                                                      |
| **Config Loader**           | [envconfig](https://github.com/kelseyhightower/envconfig)   | Type-safe environment variable parsing                                               |
| **Structured Logging**      | Standard Library `log/slog`                                 | High-performance JSON and pretty-printed logs                                        |
| **Security**                | `golang-jwt/jwt/v5`, `bcrypt`                               | Stateless JWT, password hashing                                                      |
| **Automation & Containers** | Docker, Docker Compose, [Taskfile](https://taskfile.dev/)   | Multi-stage containerized environment, cross-platform build and execution automation |

---

## Design Patterns & Engineering Principles

- **Optimistic Locking**: Handled via a dedicated `version` column on mutable tables (`users`, `accounts`, `categories`, `expenses`). Updates verify `WHERE id = $id AND version = $version`, returning a `409 Conflict` on concurrent modifications.
- **Tri-State Nullable / Presence Pattern**: Custom `core_types.Nullable[T]` implementation that distinguishes between:
  1. *Omitted field* (`Set: false`) — do not update.
  2. *Explicit `null`* (`Set: true, Val: nil`) — rejected on non-nullable domain fields.
  3. *New value provided* (`Set: true, Val: &v`) — perform partial update.
- **Defensive Copying in Domain**: Entity mutations in `ApplyPatch` operate on atomic temporary copies (`tmp := *entity`). If domain validation fails, the in-memory entity remains strictly unmodified.
- **Smart Database Adapter (`mapErrors`)**: Centralizes the mapping of driver-specific PostgreSQL error codes (e.g., `23505 Unique Violation`, `23503 Foreign Key Violation`, `pgx.ErrNoRows`) into decoupled core errors (`ErrAlreadyExists`, `ErrNotFound`), preventing infrastructure leakage into service layers.
- **Dual-Stream Structured Logging (`TeeHandler`)**: Log output is split dynamically:
  - Console (`stdout`): Formatted with human-readable ANSI colors via `slogpretty`.
  - File (`out/logs/`): Clean, structured JSON/text records with timestamped file rotation for auditability.
- **Fail-Fast Initialization (`Must` Pattern)**: Universal generic `utils.Must[T]` wrapper ensures misconfigured environment variables abort execution immediately at startup rather than failing silently at runtime.

---

## Directory Structure

```text
.
├── api/
│   └── expense-tracker.postman_collection.json     # Pre-configured Postman collection
├── cmd/
│   └── expense-tracker/
│       ├── main.go               # Application entry point & Composition Root
│       └── Dockerfile            # Multi-stage production container build
├── internal/
│   ├── core/
│   │   ├── auth/jwt/             # JWT token creation, claims & validation 
│   │   ├── cache/redis/          # App caching
│   │   ├── config/               # Timezone, server, and core configs
│   │   ├── domain/               # Pure business models (User, Account, Expense, etc.)
│   │   ├── errors/               # Domain-wide sentinel error definitions
│   │   ├── logger/               # slog setup, handlers (slogpretty, teehandler, slogdiscard)
│   │   ├── repository/postgres/  # Connection pool interfaces & pgx adapters
│   │   ├── transport/http/       # HTTP middleware, context helpers, request/response decoders
│   │   └── types/                # Core utility types (Nullable[T])
│   ├── features/                 # Modular business capabilities
│   │   ├── accounts/             # Account management & shared access
│   │   ├── analytics/            # Financial reporting, trends & summary dashboards
│   │   ├── auth/                 # Authentication, signup & signin
│   │   ├── categories/           # System and personal categories
│   │   ├── expenses/             # Expense tracking & filtering
│   │   └── users/                # User profile management & admin controls
│   └── utils/                    # Generic utilities (Must[T]). Decimal <-> Integer helpers
├── migrations/                   # Sequential SQL migrations (up/down)
├── out/logs/                     # Local file-based log outputs (git-ignored)
├── docker-compose.yaml           # Multi-container orchestration (App, PG, Redis, Migrations)
├── Taskfile.yaml                 # Automation commands (Task runner)
├── .env.example                  # Documented template for environment variables
└── go.mod                        # Go module dependencies
```

---

## Getting Started

### Prerequisites
- [Go 1.26+](https://go.dev/)
- [Docker & Docker Compose](https://www.docker.com/)
- [Task](https://taskfile.dev/)

### Environment Configuration
Clone the repository and initialize your `.env` file:

```bash
cp .env.example .env
```

Review `.env` and configure your credentials:
```env
ENV=local
LOG_LEVEL=debug
TIME_ZONE=Europe/Moscow
LOGGER_FOLDER=./out/logs

SERVER_PORT=8787
HTTP_ADDR=:8787
HTTP_SHUTDOWN_TIMEOUT=10s

POSTGRES_HOST=expense-tracker-postgres
POSTGRES_PORT=5432
POSTGRES_USER=expense_admin
POSTGRES_PASSWORD=your_secure_password
POSTGRES_DB=expense_tracker
POSTGRES_TIMEOUT=5s

REDIS_HOST=expense-tracker-cache
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password
REDIS_DB=0

JWT_SECRET=your_super_secret_jwt_key_at_least_32_chars_long
JWT_TTL=24h
```

### Database Migrations
Initialize the schema using the containerized migration tool:

```bash
task migrate-up         # Apply all pending migrations
task migrate-down       # Roll back the last migration
task migrate-version    # Check current migration version
```

### Running the Application

#### Mode 1: Hybrid Development (Recommended for daily coding)
Spins up PostgreSQL and Redis in Docker, while running the Go application locally with live file logging:

```bash
task dev
```

#### Mode 2: Full Dockerized Deployment (Production-like)
Builds and runs all components inside isolated Docker containers:
```bash
task compose-up       # Build & start all containers
task compose-down     # Gracefully stop containers and purge local images
```

#### Utility Commands
```bash
task logs-clean       # Clean up local log files with confirmation
```

### Testing with Postman
A pre-configured Postman collection is included in the repository:
1. Import `api/expense_tracker.postman_collection.json` into Postman.
2. Run the `POST /auth/signin` request — the JWT token will be automatically captured and applied to all protected endpoints via collection variables.
3. Test any protected endpoints (`Users`, `Accounts`, `Expenses`, `Analytics`).

---

## API Reference Overview

All protected endpoints require the HTTP header:  
`Authorization: Bearer <JWT_TOKEN>`

### Authentication
- `POST /api/v1/auth/signup` — Register a new user & auto-provision a default "Personal" account.
- `POST /api/v1/auth/signin` — Authenticate credentials and receive a JWT.

### Users (Protected)
- `GET /api/v1/users` — *(Admin Only)* List users with optional pagination (`?limit=10&offset=0`).
- `GET /api/v1/users/{id}` — Fetch user details (Self or Admin).
- `PATCH /api/v1/users/{id}` — Partially update email or password with optimistic locking.
- `PATCH /api/v1/users/{id}/role` — *(Admin Only)* Promote or demote user roles (`user` or `admin`).
- `DELETE /api/v1/users/{id}` — Delete user account (Self or Admin).

### Accounts & Wallets (Protected)
- `POST /api/v1/accounts` — Create a new financial wallet.
- `GET /api/v1/accounts` — List all accessible accounts (Owned + Shared).
- `GET /api/v1/accounts/{id}` — Fetch account details (Owner, Member, or Admin).
- `PATCH /api/v1/accounts/{id}` — Rename an account (Owner or Admin only).
- `DELETE /api/v1/accounts/{id}` — Cascade delete an account (Owner or Admin only).
- `POST /api/v1/accounts/{id}/share` — Grant access to another user via email (Owner or Admin only).
- `DELETE /api/v1/accounts/{id}/share` — Revoke access from a member (Owner or Admin only).

### Categories (Protected)
- `POST /api/v1/categories` — Create a custom personal category.
- `GET /api/v1/categories` — List available categories (System defaults + Personal).
- `GET /api/v1/categories/{id}` — Get single category details (System or Owned).
- `PATCH /api/v1/categories/{id}` — Rename a personal category (Admin for system).
- `DELETE /api/v1/categories/{id}` — Delete a personal category (Admin for system).

### Expenses (Protected)
- `POST /api/v1/expenses` — Log a new transaction (validates account access, category access, and positive amount).
- `GET /api/v1/accounts/{account_id}/expenses` — Query transactions of a specific account with filters:
  - `?category_id=X`
  - `?from=YYYY-MM-DD&to=YYYY-MM-DD`
  - `?limit=20&offset=0`
- `GET /api/v1/expenses/{id}` — Fetch details of a single transaction.
- `PATCH /api/v1/expenses/{id}` — Partially update transaction (`amount`, `category_id`, `description`, `date`).
- `DELETE /api/v1/expenses/{id}` — Remove a transaction (Author, Account Owner, or Admin).

### Analytics & Reporting (Protected)
- `GET /api/v1/accounts/{account_id}/analytics/summary` — Financial summary dashboard (Total, Count, Average, Max single expense, Top category).
- `GET /api/v1/accounts/{account_id}/analytics/categories` — Category spending breakdown with SQL window-function percentage shares (for Donut/Pie charts).
- `GET /api/v1/accounts/{account_id}/analytics/trends` — Chronological time series grouped by `?interval=day` or `?interval=month` (for Line/Bar charts).
- `GET /api/v1/accounts/{account_id}/analytics/members` — Member spending distribution for shared accounts.

---

## Roadmap

Planned improvements for upcoming iterations:

- [ ] **Interactive Swagger / OpenAPI Specification:**
  - Generate automated OpenAPI v2/v3 documentation using [swaggo/swag](https://github.com/swaggo/swag).
  - Expose an interactive Swagger UI directly from the server at `/swagger/index.html`.
- [ ] **External Currency Exchange Rates API & Redis Caching:**
  - Integrate an external exchange rate API (e.g., ExchangeRate-API).
  - Cache real-time conversion rates in Redis with a 1-hour TTL to enable multi-currency balance reporting.