### Self-documented makefile
CONFIG       ?= configs/config.yml
BIN          ?= bin/salehelp
IMAGE        ?= salehelp:local
# Default matches docker-compose.yml. Override: make migrate-up DATABASE_URL=postgres://...
DATABASE_URL ?= postgres://salehelp:salehelp@localhost:5432/salehelp?sslmode=disable
MIGRATIONS   ?= db/migrations

help: ### Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

run: ### run app (CONFIG=configs/config.yml; copy configs/config.example.yml first)
	go run ./cmd/salehelp -c=$(CONFIG)
.PHONY: run

build: ### build the binary into bin/salehelp
	CGO_ENABLED=0 go build -o $(BIN) ./cmd/salehelp
.PHONY: build

docker-build: ### build docker image (IMAGE=salehelp:local)
	docker build -t $(IMAGE) .
.PHONY: docker-build

lint: ### run golangci-lint
	golangci-lint run ./...
.PHONY: lint

fix: ### fix all files with golangci-lint
	golangci-lint run --fix ./...
.PHONY: fix

test: ### run all tests
	go test -v ./...
.PHONY: test

race: ### run tests with the race detector
	CGO_ENABLED=1 go test -v -cover -race ./...
.PHONY: race

coverage: ### run tests and generate coverage report (coverage.out)
	go test ./internal/... -coverpkg=./internal/... -covermode=atomic -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
.PHONY: coverage

swagger: ### generate swagger docs (go install github.com/swaggo/swag/cmd/swag@v1.8.12)
	swag init -g cmd/swagger/main.go -o docs --parseDependency --parseInternal
.PHONY: swagger

compose: ### start local PostgreSQL via docker compose
	docker compose up -d
.PHONY: compose

compose-down: ### stop local PostgreSQL (data volume is kept)
	docker compose down --remove-orphans
.PHONY: compose-down

migrate-create: ### create a new migration: make migrate-create name=add_something
	goose -dir $(MIGRATIONS) create $(name) sql
.PHONY: migrate-create

migrate-up: ### apply all migrations (DATABASE_URL)
	goose -dir $(MIGRATIONS) postgres "$(DATABASE_URL)" up
.PHONY: migrate-up

migrate-down: ### roll back the last migration (DATABASE_URL)
	goose -dir $(MIGRATIONS) postgres "$(DATABASE_URL)" down
.PHONY: migrate-down

migrate-reset: ### roll back all migrations (DATABASE_URL)
	goose -dir $(MIGRATIONS) postgres "$(DATABASE_URL)" reset
.PHONY: migrate-reset

db-forward: ### port-forward Postgres from the minikube cluster to localhost:5432
	kubectl port-forward service/postgres-salehelp 5432:5432
.PHONY: db-forward
