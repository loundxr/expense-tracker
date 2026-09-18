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

### 1. Authentication & RBAC
- **Secure Registration & Login**: Password hashing via `bcrypt` (default cost: 10).
- **Stateless JWT Authorization**: Cryptographically signed HMAC-SHA256 tokens (`golang-jwt/jwt/v5`) containing `uid` and `role` claims with configurable TTL.
- **Role-Based Access Control (RBAC)**: Distinct permissions for `user` and `admin` roles, enforced via lightweight middleware and service-level authorization guards.

### 2. User Management
- Full user profile lifecycle with administrative controls.
- **Access Boundary Enforcement**: Regular users can only access and modify their own profile; administrators can manage any user.
- **Independent Role Elevation**: Dedicated administrative endpoint (`PATCH /api/v1/users/{id}/role`) separated from profile editing to prevent privilege escalation.

### 3. Accounts & Shared Budgets
- **Multi-Wallet Support**: Create distinct accounts (e.g., "Personal", "Family Budget", "Travel").
- **Automatic Default Account**: Automatically provisions a default "Personal" account upon initial user registration.
- **Shared Access (Multi-User Collaboration)**: Account owners can grant (`POST /share`) or revoke (`DELETE /share`) access to other users via their email addresses, backed by a composite primary key in PostgreSQL (`account_users`).

### 4. Categories (System & Custom)
- **Hybrid Scoping**: Seamlessly combines system-wide default categories (`user_id IS NULL`, seeded via SQL migrations) with custom user-created categories (`user_id = uid`).
- **Data Integrity Protection**: Regular users cannot alter or delete system categories; administrators maintain full moderation capabilities.

### 5. Expenses & Financial Operations
- **High-Precision Money Handling**: Financial amounts are stored as integer cents (`BIGINT`) in PostgreSQL to eliminate floating-point rounding errors (`0.1 + 0.2 != 0.3`). Values are converted to decimal format (`float64`) exclusively at the transport boundary.
- **Strict Relationship Validation**: Expenses can only be charged to accounts the caller has verified access to, and categorized using accessible categories.
- **Deterministic Filtering & Pagination**: Filter transactions by account, category, date range (`from` / `to` with timezone awareness), and paginated via `limit` and `offset`. Uses composite sorting (`ORDER BY date DESC, id DESC`) to eliminate pagination flicker.

### 6. Real-Time Financial Analytics
- **Financial Summary Dashboard**: Instant overview computing total expenditure, transaction counts, average purchase size, largest single purchase, and top spending category.
- **Category Breakdown (Donut/Pie Chart Data)**: Single-pass SQL aggregation utilizing PostgreSQL window functions (`SUM(SUM(amount)) OVER()`) to compute the exact percentage share per category.
- **Spending Trends (Bar/Line Chart Data)**: Time-series analysis aggregating expenses chronologically across configurable intervals (`day` or `month`) using PostgreSQL's `DATE_TRUNC`.
- **Shared Budget Member Contribution**: Visual breakdown revealing the exact financial contribution and percentage of total expenditure for every member in a shared account.

---

## Technical Stack

| Component              | Technology                                                  | Purpose                                                |
| :--------------------- | :---------------------------------------------------------- | :----------------------------------------------------- |
| **Language**           | [Go 1.26.2](https://go.dev/)                                | Core runtime                                           |
| **HTTP Router**        | [go-chi/chi v5](https://github.com/go-chi/chi)              | Fast, idiomatic, lightweight HTTP routing              |
| **Database**           | [PostgreSQL 18.6](https://www.postgresql.org/)              | Relational ACID persistence                            |
| **DB Driver & Pool**   | [pgx/v5 (pgxpool)](https://github.com/jackc/pgx)            | High-performance PostgreSQL driver and connection pool |
| **Cache Store**        | [Redis 8+](https://redis.io/)                               | In-memory cache & future token invalidation store      |
| **Migrations**         | [golang-migrate](https://github.com/golang-migrate/migrate) | Database schema version control                        |
| **Config Loader**      | [envconfig](https://github.com/kelseyhightower/envconfig)   | Type-safe environment variable parsing                 |
| **Structured Logging** | Standard Library `log/slog`                                 | High-performance JSON and pretty-printed logs          |
| **Task Runner**        | [go-task/task](https://taskfile.dev/)                       | Cross-platform build and execution automation          |
| **Containerization**   | Docker & Docker Compose                                     | Multi-stage containerized environment                  |

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
- **Go**: Version `1.26` or higher installed locally.
- **Docker & Docker Compose**: Installed and running.
- **Task**: Installed (`brew install go-task` on macOS or via [taskfile.dev](https://taskfile.dev/installation/)).

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
# Apply all pending migrations
task migrate-up

# Roll back the last migration
task migrate-down

# Check current migration version
task migrate-version
```

### Running the Application

#### Mode 1: Hybrid Development (Recommended for daily coding)
Spins up PostgreSQL and Redis in Docker, while running the Go application locally with live file logging:

```bash
task dev
```

#### Mode 2: Full Dockerized Deployment (Production-like)
Builds and runs all components (Go binary, PostgreSQL, Redis) inside an isolated Docker bridge network:

```bash
# Build and run all services
task compose-up

# Gracefully shut down and purge local application images
task compose-down
```

#### Utility Commands
```bash
# Clean up local log files with confirmation
task logs-clean
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

- [ ] **Cross-Feature Database Transactions:** Wrap `User` creation and default `Account` provisioning in a single atomic database transaction (`WithinTransaction`).
- [ ] **Dynamic Token Invalidation via Redis:** Invalidate active JWTs upon role downgrade using Redis-backed timestamp checks (`iat`).
- [ ] **Multi-Currency & Exchange Rates:** Integrate an external exchange rate API and cache conversion rates in Redis (1h TTL).