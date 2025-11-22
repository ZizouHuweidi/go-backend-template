# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## up: start core services (App, DB, Redis, Traefik) in Dev mode
.PHONY: up
up:
	docker-compose up -d

## up-prod: start core services in Prod mode (no overrides)
.PHONY: up-prod
up-prod:
	docker-compose -f docker-compose.yml up -d

## up-traefik: start production proxy
.PHONY: up-traefik
up-traefik:
	docker-compose -f docker-compose.traefik.yml up -d

## up-full: start all services including observability stack
.PHONY: up-full
up-full:
	docker-compose -f docker-compose.yml -f docker-compose.override.yml -f docker-compose.observability.yml up -d

## down: stop all services
.PHONY: down
down:
	docker-compose -f docker-compose.yml -f docker-compose.observability.yml down --remove-orphans

## logs: view logs for all services
.PHONY: logs
logs:
	docker-compose -f docker-compose.yml -f docker-compose.observability.yml logs -f

## quick-start: start core services and view logs
.PHONY: quick-start
quick-start: up logs

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	go test -v -race ./...

## lint: run linters
.PHONY: lint
lint:
	golangci-lint run

## lint-fix: run linters and fix issues
.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy
	go fmt ./...

# ==================================================================================== #
# BUILD & RUN
# ==================================================================================== #

## build: build the binary
.PHONY: build
build:
	go build -o bin/server ./cmd/api

## run: run the binary (requires DB on localhost)
.PHONY: run
run: build
	./bin/server

## docker-up-prod: start production containers
.PHONY: docker-up-prod
docker-up-prod:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

# ==================================================================================== #
# MIGRATIONS
# ==================================================================================== #

## migrate-create name=$1: create a new migration file
.PHONY: migrate-create
migrate-create:
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## migrate-up: apply all up migrations
.PHONY: migrate-up
migrate-up:
	migrate -path=./migrations -database="${DB_DSN}" up

## migrate-down steps=$1: rollback the last migration (default 1 step)
.PHONY: migrate-down
migrate-down:
	migrate -path=./migrations -database="${DB_DSN}" down ${steps}

## migrate-status: check migration status
.PHONY: migrate-status
migrate-status:
	migrate -path=./migrations -database="${DB_DSN}" version

## migrate-goto version=$1: migrate to a specific version
.PHONY: migrate-goto
migrate-goto:
	migrate -path=./migrations -database="${DB_DSN}" goto ${version}

## migrate-force version=$1: force a specific version
.PHONY: migrate-force
migrate-force:
	migrate -path=./migrations -database="${DB_DSN}" force ${version}
