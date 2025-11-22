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

# Dev seed (локально): заполняет несколько событий, если таблица пуста
seed:
	DB_TYPE=sqlite SEED_ON_START=true go run ./main.go & 
	sleep 1 && curl -s http://localhost:8080/health >/dev/null || true 
	kill $$! || true
