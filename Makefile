# ==============================================================================
# Boilerplate Clean Architecture + DDD — Makefile
# ==============================================================================
# Variabel default — dapat di-override: make migrate-up DB_URL=postgres://...
# ==============================================================================

COMPOSE_FILE     = deploy/docker/docker-compose.yml
COMPOSE          = docker compose -f $(COMPOSE_FILE)

# golang-migrate database URL (baca dari .config.yml, override saat dibutuhkan)
DB_HOST         ?= localhost
DB_PORT         ?= 5432
DB_USER         ?= postgres
DB_PASS         ?= postgres
DB_NAME         ?= boilerplate
DB_SSLMODE      ?= disable
DB_URL          ?= postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

MIGRATIONS_DIR   = migrations
BINARY_NAME      = boilerplate
BUILD_DIR        = .dist

.PHONY: help \
        dev build up down stop restart logs ps \
        migrate-install migrate-create migrate-up migrate-down migrate-force migrate-version \
        tidy lint test

# ==============================================================================
# Help
# ==============================================================================

## help: tampilkan semua target yang tersedia
help:
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@echo "  Docker:"
	@awk '/^## docker/,0' $(MAKEFILE_LIST) | grep '^##' | sed 's/## /  /'
	@echo ""
	@echo "  Migrate:"
	@awk '/^## migrate/,0' $(MAKEFILE_LIST) | grep '^##' | sed 's/## /  /'
	@echo ""
	@echo "  Dev:"
	@awk '/^## dev/,0' $(MAKEFILE_LIST) | grep '^##' | sed 's/## /  /'
	@echo ""

# ==============================================================================
# Docker — Development
# ==============================================================================

## docker dev: jalankan semua service dalam mode development (dengan Air hot-reload)
dev:
	$(COMPOSE) up

## docker dev detached: jalankan di background
dev-d:
	$(COMPOSE) up -d

## docker build: build ulang image app (berguna setelah ubah Dockerfile/go.mod)
build:
	$(COMPOSE) build app

## docker build no-cache: build ulang tanpa cache layer
build-nocache:
	$(COMPOSE) build --no-cache app

## docker up: jalankan semua service (tidak rebuild image)
up:
	$(COMPOSE) up -d

## docker down: hentikan dan hapus semua container + network (volume tetap aman)
down:
	$(COMPOSE) down

## docker down volumes: hentikan dan hapus semua container + volume (HATI-HATI: data DB hilang)
down-v:
	$(COMPOSE) down -v

## docker stop: hentikan semua container tanpa menghapus
stop:
	$(COMPOSE) stop

## docker stop app: hentikan hanya container app
stop-app:
	$(COMPOSE) stop app

## docker restart: restart semua service
restart:
	$(COMPOSE) restart

## docker restart app: restart hanya container app
restart-app:
	$(COMPOSE) restart app

## docker logs: tampilkan log semua service secara realtime
logs:
	$(COMPOSE) logs -f

## docker logs app: tampilkan log hanya app
logs-app:
	$(COMPOSE) logs -f app

## docker ps: lihat status semua container
ps:
	$(COMPOSE) ps

## docker shell: masuk ke shell container app
shell:
	$(COMPOSE) exec app sh

# ==============================================================================
# Build lokal (tanpa Docker)
# ==============================================================================

## dev build-local: build binary lokal
build-local:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server/main.go

## dev run-local: build dan jalankan binary lokal
run-local: build-local
	./$(BUILD_DIR)/$(BINARY_NAME)

# ==============================================================================
# golang-migrate
# Dokumentasi: https://github.com/golang-migrate/migrate
# ==============================================================================

## migrate-install: install golang-migrate CLI
migrate-install:
	@echo "Installing golang-migrate..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Done. Pastikan $(shell go env GOPATH)/bin ada di PATH."

## migrate-create name=<nama>: buat pasangan file migration baru
## Contoh: make migrate-create name=create_users_table
migrate-create:
	@if [ -z "$(name)" ]; then echo "ERROR: Berikan nama migration. Contoh: make migrate-create name=create_users_table"; exit 1; fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

## migrate-up: jalankan semua migration yang belum dieksekusi
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

## migrate-up step=N: jalankan N migration ke atas
migrate-up-n:
	@if [ -z "$(step)" ]; then echo "ERROR: Berikan step. Contoh: make migrate-up-n step=1"; exit 1; fi
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up $(step)

## migrate-down: rollback satu migration terakhir
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

## migrate-down-all: rollback semua migration (HATI-HATI: data hilang)
migrate-down-all:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

## migrate-force version=N: paksa set versi migration (gunakan saat dirty state)
## Contoh: make migrate-force version=1
migrate-force:
	@if [ -z "$(version)" ]; then echo "ERROR: Berikan version. Contoh: make migrate-force version=1"; exit 1; fi
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(version)

## migrate-version: tampilkan versi migration saat ini
migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

## migrate-drop: hapus semua tabel (HATI-HATI: destructive)
migrate-drop:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" drop -f

# ==============================================================================
# Dev tools
# ==============================================================================

## dev tidy: sinkronisasi go.mod dan go.sum
tidy:
	go mod tidy

## dev lint: jalankan golangci-lint
lint:
	golangci-lint run ./...

## dev test: jalankan semua unit test
test:
	go test -v -race -cover ./...
