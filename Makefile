APP_NAME=pawtrack
PORT?=8080

.PHONY: deps run build test tidy docker-build docker-up docker-down seed

deps:
	go mod tidy

run:
	go run ./main.go

build:
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME) ./main.go

test:
	go test ./...

tidy:
	go fmt ./...
	go vet ./...

# Docker helpers
docker-build:
	docker build -t pawtrack:local .

docker-up:
	docker compose up --build

docker-down:
	docker compose down -v

# Dev seed (local): populates some events if table is empty
seed:
	DB_TYPE=sqlite SEED_ON_START=true go run ./main.go & 
	sleep 1 && curl -s http://localhost:8080/health >/dev/null || true 
	kill $$! || true

# Run E2E tests (requires docker-up)
test-e2e:
	docker run --rm -v $$(pwd):/app -w /app --network pawtrack_default \
		-e E2E_BASE_URL="http://app:8080/api/v1" \
		-e E2E_DB_DSN="postgres://pawtrack:pawtrack@db:5432/pawtrack?sslmode=disable" \
		golang:1.23 go test -v ./tests/e2e/...
