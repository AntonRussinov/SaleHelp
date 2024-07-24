### Self-documented makefile
help: ### Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

run: ### run app 
	go run cmd/salehelp/main.go -c=configs/config.yml
.PHONY: run	

build: ### builds docker image
	docker build -f build/package/Dockerfile .
.PHONY: build

lint: ### run golangci-lint
	golangci-lint run ./...
.PHONY: lint

fix: ### fix all files with golangci-lint
	golangci-lint run --fix ./...
.PHONY: fix	

test: ### run all tests
	go test -v ./...
.PHONY: test

race: ### run test with rave
	CGO_ENABLED=1 go test -v -cover -race ./...
.PHONY: race

coverage: ### run all tests and generate coverage report
###	go test -v -coverpkg=./internal/... ./internal/... -covermode=atomic -coverprofile=coverage.out
	go test ./internal/... -coverpkg=./internal/... -covermode=atomic -coverprofile=coverage.out
	go tool cover -func=coverage.out 
	go tool cover -html=coverage.out
.PHONY: coverage

swagger: ### generate swagger docs
	swag init -g cmd/swagger/main.go
.PHONY: swagger

compose: ### run docker-compose
	docker compose -f deployments/docker-compose.local.yml up -d
.PHONY: compose

compose-down: ### down docker-compose
	docker compose -f deployments/docker-compose.local.postgres.yml down --remove-orphans
.PHONY: compose-down

compose-build: ### rebuild and run docker-compose
	docker compose -f deployments/docker-compose.local.postgres.yml up --build -d
.PHONY: compose-build

compose-db: ### run compose with db only
	docker compose -f deployments/docker-compose.local.postgres.db.yml up -d
.PHONY: compose-db

compose-db-down: ### downb db-only compose
	docker compose -f deployments/docker-compose.local.postgres.db.yml down --remove-orphans
.PHONY: compose-db-down

# migrate: ### run migrations
# 	go run cmd/migration/main.go -c=configs/config.yml
# .PHONY: migrate

# new-migration: ### create new migration, call make new-migration name=your_migration_name
# 	echo "package db_migration" > "internal/db_migration/$$(date +%s)_$(name).go"
# .PHONY: new-migration

migrate-create: ### create new migration, call make migrate create name=your_migration_name
	goose  -dir db/migrations create $(name) sql
.PHONY: migrate-create

migrate-down: ### run migrations
	goose -dir db/migrations postgres "host=192.168.64.5 port=31416 user=postgres password=postgres dbname=salehelp sslmode=disable" down

migrate-reset: ### rollback all migrations
	goose -dir db/migrations postgres "host=192.168.64.5 port=31416 user=postgres password=postgres dbname=salehelp sslmode=disable" reset 

migrate-up: ### run migrations
	goose -dir db/migrations postgres "host=192.168.64.5 port=31416 user=postgres password=postgres dbname=salehelp sslmode=disable" up