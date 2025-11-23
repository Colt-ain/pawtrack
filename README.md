# pawtrack — small Go service for pet event tracking 🐾

API for tracking walks, feedings, meds, and any other pet "events".
Stack: **Go 1.23**, **Gin**, **GORM + SQLite/Postgres**.

## Quick Start (Local, SQLite)
```bash
make deps
make run
curl http://localhost:8080/health
```

## Docker
```bash
docker build -t pawtrack:local .
docker run -p 8080:8080 pawtrack:local
```

## Docker Compose (PostgreSQL)
```bash
docker compose up --build
# API: http://localhost:8080
```

## Migrations (golang-migrate)
- Migrations are located in `./migrations`.
- Run on start is enabled by `RUN_MIGRATIONS=true` (default `false`).
- Migration path can be overridden via `MIGRATIONS_DIR`.

Examples:
```bash
# SQLite (local)
export DB_TYPE=sqlite
export RUN_MIGRATIONS=true
go run ./main.go

# Postgres via docker-compose
docker compose up --build
# (RUN_MIGRATIONS=true is already set in compose)
```

## Seed (demo data)
Set `SEED_ON_START=true` — if table is empty, 3 records will be added.
Do not enable in production :)

## Testing Scripts
Verification scripts are located in `./scripts/`:
- `scripts/verify_api.sh` - Basic API health check
- `scripts/verify_dogs.sh` - Test all Dog CRUD endpoints
- `scripts/verify_users.sh` - Test all User CRUD endpoints
- `scripts/verify_event_filters.sh` - Test event filtering and authorization

Usage:
```bash
bash scripts/verify_dogs.sh
bash scripts/verify_users.sh
```

## ENV
- `ADDR` — HTTP server address (default `:8080`)
- `DB_TYPE` — `sqlite` (default) or `postgres`
- `SQLITE_DSN` — for SQLite (`file:pawtrack.db?_busy_timeout=5000&_fk=1`)
- `DATABASE_URL` — for Postgres (e.g., `postgres://pawtrack:pawtrack@db:5432/pawtrack?sslmode=disable`)
- `RUN_MIGRATIONS` — `true/false` (run migrations on start)
- `MIGRATIONS_DIR` — path to migrations (default `./migrations`)
- `SEED_ON_START` — `true/false` (add demo records if table is empty)
