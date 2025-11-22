# Go Backend Template

A production-ready, opinionated Go backend starter template designed for scalability, observability, and developer experience.

## Features

- **Framework**: [Echo](https://echo.labstack.com/) v4 - High performance, extensible, minimalist Go web framework.
- **Database**: PostgreSQL 16+ with [sqlx](https://github.com/jmoiron/sqlx) and [squirrel](https://github.com/Masterminds/squirrel) for type-safe query building.
- **Authentication**: JWT-based auth with Refresh Token Rotation and Family Tracking.
- **Password Recovery**: Email-based password recovery flow.
- **Caching & Rate Limiting**: Redis 7.
- **Observability**: Full OpenTelemetry (OTel) integration with the LGTM stack (Loki, Grafana, Tempo, Prometheus).
- **Logging**: Structured logging with `slog`.
- **Validation**: Request validation using `go-playground/validator`.
- **Reverse Proxy**: [Traefik](https://traefik.io/) with automatic HTTPS and load balancing.
- **Development**: Hot-reloading with Air, Docker Compose overrides.
- **Tools**: pgAdmin (DB GUI) and MailHog (Email testing).
- **Migrations**: Database migrations with `golang-migrate`.
- **UUIDv7**: Modern, time-sortable UUIDs for primary keys.

## Prerequisites

- [Go 1.24+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- [Make](https://www.gnu.org/software/make/)

## Quick Start

### Local Development

Start the stack with hot-reloading enabled:

```bash
make up
```

This uses `docker-compose.yml` (base) and `docker-compose.override.yml` (dev overrides). It spins up a local Traefik instance to simulate the production routing.

**Access Points:**
- **API**: `http://localhost:8080` (Direct) or `http://api.localhost.tiangolo.com` (via Traefik)
- **Traefik Dashboard**: `http://localhost:8090`
- **pgAdmin**: `http://pgadmin.localhost.tiangolo.com` (Email: `admin@example.com`, Pass: `admin`)
- **MailHog**: `http://mailhog.localhost.tiangolo.com` (SMTP: 1025, UI: 8025)

*Note: `localhost.tiangolo.com` is a public domain that points to `127.0.0.1`, useful for testing subdomains locally.*

### Production Deployment

In production, the setup is split into two parts: the Ingress Proxy (Traefik) and the Application Stack.

1.  **Start the Traefik Proxy**:
    This runs the public-facing Traefik instance that handles HTTPS (Let's Encrypt) and routing.
    ```bash
    make up-traefik
    ```

2.  **Start the Application**:
    This runs the core services (API, Postgres, Redis, pgAdmin) without development overrides.
    ```bash
    make up-prod
    ```

## Configuration

Configuration is managed via environment variables. The `internal/config` package loads these from the `.env` file or the system environment.

### Docker Compose Strategy

- **`docker-compose.yml`**: Base configuration for all environments. Defines core services (`api`, `postgres`, `redis`, `pgadmin`) and their production settings (restart policy, networks, labels).
- **`docker-compose.override.yml`**: Development overrides. Adds a local `proxy` (Traefik), `mailhog`, exposes ports, and configures `air` for hot-reloading.
- **`docker-compose.traefik.yml`**: Production Traefik configuration. Defines the main ingress proxy and `traefik-public` network.
- **`docker-compose.observability.yml`**: Optional observability stack (LGTM).

To run with observability in dev:
```bash
make up-full
```

## Project Structure

```
.
├── api/                # API definitions (Bruno collection)
├── cmd/
│   └── api/            # Main entry point
├── internal/
│   ├── auth/           # Auth handlers
│   ├── config/         # Configuration loading
│   ├── database/       # Database connection & helpers
│   ├── email/          # Email sender
│   ├── jwt/            # JWT logic
│   ├── middleware/     # Custom middleware (Auth, Logger, RateLimit)
│   ├── redis/          # Redis client
│   ├── response/       # Standardized API responses
│   ├── server/         # Server setup & routes
│   ├── telemetry/      # OpenTelemetry setup
│   ├── user/           # User domain (Handler, Service, Repo, Model)
│   └── validator/      # Input validation
├── migrations/         # SQL migrations
├── docker-compose.yml  # Base services
├── docker-compose.override.yml # Dev overrides
├── docker-compose.traefik.yml # Prod Proxy
├── docker-compose.observability.yml # Observability stack
└── Makefile            # Project commands
```

## API Documentation

A [Bruno](https://www.usebruno.com/) collection is available in `api/collection/collection.json`. Import this into Bruno to test the endpoints.

### Key Endpoints

- `POST /api/v1/auth/register`: Register a new user.
- `POST /api/v1/auth/login`: Login and receive Access/Refresh tokens.
- `POST /api/v1/auth/refresh`: Refresh access token.
- `POST /api/v1/auth/recover-password`: Request password reset email.
- `POST /api/v1/auth/reset-password`: Reset password with token.
- `GET /api/v1/users/me`: Get current user profile (Protected).
- `GET /health`: Health check.

## Commands

Run `make help` to see all available commands.

- `make up`: Start dev stack.
- `make up-prod`: Start prod stack.
- `make up-traefik`: Start prod proxy.
- `make up-full`: Start dev stack + observability.
- `make down`: Stop all services.
- `make logs`: View logs.
- `make test`: Run tests.
- `make lint`: Run linter.
- `make build`: Build binary.
- `make migrate-create name=my_migration`: Create a new migration.
- `make migrate-up`: Apply migrations.
