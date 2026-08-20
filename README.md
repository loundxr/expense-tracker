# Expense Tracker API

A powerful and scalable RESTful API for managing personal and shared finances, built with Go.

## Features
- **User Authentication:** Secure registration and login using JWT (JSON Web Tokens).
- **Shared Accounts:** Create financial accounts and share access with other users (e.g., family or project budgets).
- **Expense Management:** Track spending with detailed records, including amounts, categories, and timestamps.
- **Dynamic Categories:** Use system-default categories or create your own.
- **Analytics & Statistics:** Get insights into spending habits by category or time periods.
- **Multi-currency Support (Roadmap):** Default currency in USD with real-time conversion capabilities via external APIs.

## Tech Stack
- **Language:** Go 1.21+
- **Router:** [go-chi/chi](https://github.com/go-chi/chi) - lightweight and idiomatic.
- **Database:** PostgreSQL.
- **SQL Wrapper:** [sqlx](https://github.com/jmoiron/sqlx) - for flexible SQL mapping.
- **Logging:** `log/slog` - structured logging from the standard library.
- **Configuration:** [Viper](https://github.com/spf13/viper).
- **Containerization:** Docker & Docker Compose.
- **Migrations:** [golang-migrate](https://github.com/golang-migrate/migrate).

## Architecture
The project follows a **Layered Architecture** pattern:
1. **Handler Layer:** Manages HTTP requests, input validation, and responses.
2. **Service Layer:** Contains core business logic (e.g., permission checks for shared accounts).
3. **Repository Layer:** Handles direct database interactions.
4. **Domain Layer:** Defines core entities and interface contracts.

## Getting Started

### Prerequisites
- Docker and Docker Compose installed.

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/expense-tracker.git
   ```
2. Set up your environment variables in `.env` (refer to `.env.example`).
3. Spin up the infrastructure:
   ```bash
   docker-compose up -d
   ```