# pawtrack — маленький Go‑сервис для учёта собачьих событий 🐾

API для записи прогулок, кормлений, лекарств и любых "событий" питомца.
Стек: **Go 1.23**, **Gin**, **GORM + SQLite/Postgres**.

## Быстрый старт (локально, SQLite)
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

## Миграции (golang-migrate)
- Миграции лежат в `./migrations`.
- Запуск при старте включается `RUN_MIGRATIONS=true` (по умолчанию `false`).
- Путь к миграциям можно переопределить через `MIGRATIONS_DIR`.

Примеры:
```bash
# SQLite (локально)
export DB_TYPE=sqlite
export RUN_MIGRATIONS=true
go run ./main.go

# Postgres через docker-compose
docker compose up --build
# (в compose уже выставлено RUN_MIGRATIONS=true)
```

## Seed (демо-данные)
Поставь `SEED_ON_START=true` — при пустой таблице в неё добавятся 3 записи.
Не включай в продакшене :)

## Testing Scripts
Verification scripts are located in `./scripts/`:
- `scripts/verify_api.sh` - Basic API health check
- `scripts/verify_dogs.sh` - Test all Dog CRUD endpoints
- `scripts/verify_users.sh` - Test all User CRUD endpoints

Usage:
```bash
bash scripts/verify_dogs.sh
bash scripts/verify_users.sh
```

## ENV
- `ADDR` — адрес HTTP сервера (по умолчанию `:8080`)
- `DB_TYPE` — `sqlite` (по умолчанию) или `postgres`
- `SQLITE_DSN` — для SQLite (`file:pawtrack.db?_busy_timeout=5000&_fk=1`)
- `DATABASE_URL` — для Postgres (например, `postgres://pawtrack:pawtrack@db:5432/pawtrack?sslmode=disable`)
- `RUN_MIGRATIONS` — `true/false` (запуск миграций при старте)
- `MIGRATIONS_DIR` — путь к миграциям (по умолчанию `./migrations`)
- `SEED_ON_START` — `true/false` (добавить демо-записи при пустой таблице)
